package main

import (
	"flag"
	"go/ast"
	"log"
	"strconv"
	"strings"

	"github.com/nasa9084/go-openapi/internal/astutil"
	"github.com/nasa9084/go-openapi/internal/generator"
)

func main() {
	flag.Parse()
	g := generator.New("mkunmarshalyaml.go")

	objects, err := astutil.ParseOpenAPIObjects("interfaces.go")
	if err != nil {
		log.Fatal(err)
	}

	g.Printf("\nfunc q(b []byte) []byte {")
	g.Import("", "bytes")
	g.Printf("\nif !bytes.HasPrefix(b, []byte(\"|\")) {")
	g.Printf("\nif bytes.ContainsRune(b, '\\'') {")
	g.Printf("\nreturn append([]byte{'\"'}, append(b, '\"')...)")
	g.Printf("\n}")
	g.Printf("\nreturn append([]byte{'\\''}, append(b, '\\'')...)")
	g.Printf("\n}")
	g.Printf("\nreturn b")
	g.Printf("\n}")

	for _, object := range objects {
		log.Printf("generate %s.Unmarshal()", object.Name)
		g.Printf("\n\nfunc (v *%s) UnmarshalYAML(b []byte) error {", object.Name)
		g.Printf("\nvar proxy map[string]raw")
		g.Import("yaml", "github.com/goccy/go-yaml")
		g.Printf("\nif err := yaml.Unmarshal(b, &proxy); err != nil {")
		g.Printf("\nreturn err")
		g.Printf("\n}")

		var noUnknown bool
		for _, field := range object.Fields {
			if yamlName(field) == "$ref" {
				g.Printf("\nif referenceBytes, ok := proxy[\"$ref\"]; ok {")
				g.Printf("\nvar referenceVal string")
				g.Import("yaml", "github.com/goccy/go-yaml")
				g.Printf("\nif err := yaml.Unmarshal(q(referenceBytes), &referenceVal); err != nil {")
				g.Printf("\nreturn err")
				g.Printf("\n}")
				g.Printf("\nv.reference = referenceVal")
				g.Printf("\ndelete(proxy, \"$ref\")")
				g.Printf("\nreturn nil")
				g.Printf("\n}")
				continue
			}
		}

		for _, field := range object.Fields {
			yn := yamlName(field)
			required := isRequired(field.Tags)

			if yn == "-" || yn == "$ref" {
				continue
			}

			if isInline(field.Tags) {
				ft, ok := field.Type.(*ast.MapType)
				if !ok {
					log.Fatalf("expected map for inline %s but %s", yn, field.Type)
				}
				formatTag := field.Tags["format"]
				g.Printf("\n%s := map[string]%s{}", field.Name, astutil.TypeString(ft.Value))
				g.Printf("\nfor key, val := range proxy {")
				if len(formatTag) > 0 {
					switch formatTag[0] {
					case "prefix":
						g.Import("", "strings")
						g.Printf("\nif !strings.HasPrefix(key, \"%s\") {", formatTag[1])
						g.Printf("\ncontinue")
						g.Printf("\n}")
					case "regexp":
						g.Import("", "regexp")
						g.Printf("\n%sRegexp := regexp.MustCompile(`%s`)", field.Name, formatTag[1])
						g.Printf("\nif !%sRegexp.MatchString(key) {", field.Name)
						g.Printf("\ncontinue")
						g.Printf("\n}")
					case "runtime":
						g.Printf("\nif !IsRuntimeExpr(key) {")
						g.Printf("\ncontinue")
						g.Printf("\n}")
					}
				} else {
					noUnknown = true
				}
				g.Printf("\nvar %sv %s", field.Name, strings.TrimPrefix(astutil.TypeString(ft.Value), "*"))
				g.Printf("\nif err := yaml.Unmarshal(val, &%sv); err != nil {", field.Name)
				g.Printf("\nreturn err")
				g.Printf("\n}")
				g.Printf("\n%s[key] = ", field.Name)
				if _, ok := ft.Value.(*ast.StarExpr); ok {
					g.Printf("&")
				}
				g.Printf("%sv", field.Name)
				g.Printf("\ndelete(proxy, key)")
				g.Printf("\n}")
				g.Printf("\nif len(%s) != 0 {", field.Name)
				g.Printf("\nv.%s = %s", field.Name, field.Name)
				g.Printf("\n}")
				continue
			}

			g.Printf("\n\n")
			if required {
				g.Printf("%sBytes, ok := proxy[\"%s\"]", field.Name, yn)
				g.Printf("\nif !ok {")
				g.Printf("\nreturn ErrRequired(%s)", strconv.Quote(yn))
				g.Printf("\n}")
			} else {
				g.Printf("if %sBytes, ok := proxy[\"%s\"]; ok {", field.Name, yn)
			}

			typName := strings.TrimPrefix(astutil.TypeString(field.Type), "*")
			g.Printf("\nvar %sVal %s", field.Name, typName)
			if typName == "string" {
				g.Printf("\nif err := yaml.Unmarshal(q(%sBytes), &%[1]sVal); err != nil {", field.Name)
			} else {
				g.Printf("\nif err := yaml.Unmarshal(%sBytes, &%[1]sVal); err != nil {", field.Name)
			}
			g.Printf("\nreturn err")
			g.Printf("\n}")
			g.Printf("\nv.%s = ", field.Name)
			if _, ok := field.Type.(*ast.StarExpr); ok {
				g.Printf("&")
			}
			g.Printf("%sVal", field.Name)
			g.Printf("\ndelete(proxy, `%s`)", yn)
			if !required {
				g.Printf("\n}")
			}

			switch field.Tags.Get("format") {
			case "semver":
				g.Printf("\n\nif !isValidSemVer(v.%s) {", field.Name)
				g.Import("", "errors")
				g.Printf("\nreturn errors.New(`\"%s\" field must be a valid semantic version but not`)", yn)
				g.Printf("\n}")
			case "url":
				g.Printf("\n")
				if !required {
					g.Printf("\nif v.%s != \"\" {", field.Name)
				}

				if len(field.Tags["format"]) > 1 && field.Tags["format"][1] == "template" {
					g.Printf("\nif err := validateURLTemplate(v.%s)", field.Name)
				} else {
					g.Import("", "net/url")
					g.Printf("\nif _, err := url.ParseRequestURI(v.%s)", field.Name)
				}
				g.Printf("; err != nil {")
				g.Printf("\nreturn err")
				g.Printf("\n}")
				if !required {
					g.Printf("\n}")
				}
			case "email":
				g.Printf("\n")
				if !required {
					g.Printf("\nif v.%s != \"\" {", field.Name)
				}
				g.Printf("\n\nif v.%s != \"\" && !emailRegexp.MatchString(v.%[1]s) {", field.Name)
				g.Import("", "errors")
				g.Printf("\nreturn errors.New(`\"%s\" field must be an email address`)", yn)
				g.Printf("\n}")
				if !required {
					g.Printf("\n}")
				}
			case "runtime":
				if _, ok := field.Type.(*ast.MapType); ok {
					g.Printf("\n\nfor key := range v.%s {", field.Name)
					g.Printf("\nif !matchRuntimeExpr(key) {")
					g.Import("", "errors")
					g.Printf("\nreturn errors.New(`the keys of \"%s\" must be a runtime expression`)", yn)
					g.Printf("\n}")
					g.Printf("\n}")
				}
			case "regexp":
				if _, ok := field.Type.(*ast.MapType); ok {
					g.Import("", "regexp")
					g.Printf("\n\n%sRegexp := regexp.MustCompile(`%s`)", field.Name, field.Tags["format"][1])
					g.Printf("\nfor key := range v.%s {", field.Name)
					g.Printf("\nif !%sRegexp.MatchString(v.%s) {", field.Name, field.Name)
					g.Import("", "errors")
					g.Printf("\nreturn errors.New(`the keys of \"%s\" must be match \"%s\"`)", yn, field.Tags["format"][1])
					g.Printf("\n}")
				}
			}
			if list, ok := field.Tags["oneof"]; ok {
				g.Printf("\n")
				if !required {
					g.Printf("\nif v.%s != \"\" {", field.Name)
				}
				g.Printf("\nif !isOneOf(v.%s, %#v) {", field.Name, list)
				g.Import("", "errors")
				g.Printf("\nreturn errors.New(`\"%s\" field must be one of [%s]`)", yn, strings.Join(quoteEachString(list), ", "))
				g.Printf("\n}")
				if !required {
					g.Printf("\n}")
				}
			}
		}
		if !noUnknown {
			g.Printf("\nif len(proxy) != 0 {")
			g.Printf("\nfor k := range proxy {")
			g.Printf("\nreturn ErrUnknownKey(k)")
			g.Printf("\n}")
			g.Printf("\n}")
		}

		if object.Name == "OpenAPI" {
			g.Printf("\nv.setRoot(v)")
		}

		g.Printf("\nreturn nil")
		g.Printf("\n}")
	}
	if err := g.Save("unmarshalyaml_gen.go"); err != nil {
		log.Fatal(err)
	}
}

func isInline(t astutil.Tags) bool {
	vs := t["yaml"]
	if len(vs) < 2 {
		return false
	}
	for _, v := range vs[1:] {
		if v == "inline" {
			return true
		}
	}
	return false
}

func isRequired(t astutil.Tags) bool {
	return t.Get("required") != ""
}

func yamlName(field astutil.OpenAPIObjectField) string {
	yn := field.Tags.Get("yaml")
	if yn != "" {
		return yn
	}
	return field.Name
}

func quoteEachString(list []string) []string {
	ret := make([]string, len(list))
	for i := range list {
		ret[i] = strconv.Quote(list[i])
	}
	return ret
}
