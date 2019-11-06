package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strconv"
	"strings"

	"github.com/nasa9084/go-openapi/internal/astutil"
	"github.com/nasa9084/go-openapi/internal/generator"
)

func init() {
	flag.Parse()
}

func main() {
	g := generator.New("mkunmarshalyaml.go")

	f, err := parser.ParseFile(token.NewFileSet(), "interfaces.go", nil, parser.ParseComments)
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

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if !isOpenAPIObject(genDecl) {
			log.Printf("%v is not an openapi object. skip.", genDecl.Specs[0].(*ast.TypeSpec).Name.Name)
			continue
		}

		for _, spec := range genDecl.Specs {
			typ, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := typ.Type.(*ast.StructType)
			if !ok {
				continue
			}

			log.Printf("generate %s.Unmarshal()", typ.Name.Name)
			g.Printf("\n\nfunc (v *%s) UnmarshalYAML(b []byte) error {", typ.Name.Name)
			g.Printf("\nvar proxy map[string]raw")
			g.Import("yaml", "github.com/goccy/go-yaml")
			g.Printf("\nif err := yaml.Unmarshal(b, &proxy); err != nil {")
			g.Printf("\nreturn err")
			g.Printf("\n}")

			var noUnknown bool
			for _, field := range st.Fields.List {
				if yamlName(field, astutil.ParseTags(field)) == "$ref" {
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

			for _, field := range st.Fields.List {
				fn := field.Names[0].Name
				tag := astutil.ParseTags(field)
				yn := yamlName(field, tag)
				required := isRequired(tag)

				if yn == "-" || yn == "$ref" {
					continue
				}

				if isInline(tag) {
					ft, ok := field.Type.(*ast.MapType)
					if !ok {
						log.Fatalf("expected map for inline %s but %s", yn, field.Type)
					}
					formatTag := tag["format"]
					g.Printf("\n%s := map[string]%s{}", fn, astutil.TypeString(ft.Value))
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
							g.Printf("\n%sRegexp := regexp.MustCompile(`%s`)", fn, formatTag[1])
							g.Printf("\nif !%sRegexp.MatchString(key) {", fn)
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
					g.Printf("\nvar %sv %s", fn, strings.TrimPrefix(astutil.TypeString(ft.Value), "*"))
					g.Printf("\nif err := yaml.Unmarshal(val, &%sv); err != nil {", fn)
					g.Printf("\nreturn err")
					g.Printf("\n}")
					g.Printf("\n%s[key] = ", fn)
					if _, ok := ft.Value.(*ast.StarExpr); ok {
						g.Printf("&")
					}
					g.Printf("%sv", fn)
					g.Printf("\ndelete(proxy, key)")
					g.Printf("\n}")
					g.Printf("\nif len(%s) != 0 {", fn)
					g.Printf("\nv.%s = %s", fn, fn)
					g.Printf("\n}")
					continue
				}

				g.Printf("\n\n")
				if required {
					g.Printf("%sBytes, ok := proxy[\"%s\"]", fn, yn)
					g.Printf("\nif !ok {")
					g.Printf("\nreturn ErrRequired(%s)", strconv.Quote(yn))
					g.Printf("\n}")
				} else {
					g.Printf("if %sBytes, ok := proxy[\"%s\"]; ok {", fn, yn)
				}

				typName := strings.TrimPrefix(astutil.TypeString(field.Type), "*")
				g.Printf("\nvar %sVal %s", fn, typName)
				if typName == "string" {
					g.Printf("\nif err := yaml.Unmarshal(q(%sBytes), &%[1]sVal); err != nil {", fn)
				} else {
					g.Printf("\nif err := yaml.Unmarshal(%sBytes, &%[1]sVal); err != nil {", fn)
				}
				g.Printf("\nreturn err")
				g.Printf("\n}")
				g.Printf("\nv.%s = ", fn)
				if _, ok := field.Type.(*ast.StarExpr); ok {
					g.Printf("&")
				}
				g.Printf("%sVal", fn)
				g.Printf("\ndelete(proxy, `%s`)", yn)
				if !required {
					g.Printf("\n}")
				}

				switch tag.Get("format") {
				case "semver":
					g.Printf("\n\nif !isValidSemVer(v.%s) {", fn)
					g.Import("", "errors")
					g.Printf("\nreturn errors.New(`\"%s\" field must be a valid semantic version but not`)", yn)
					g.Printf("\n}")
				case "url":
					g.Printf("\n")
					if !required {
						g.Printf("\nif v.%s != \"\" {", fn)
					}

					if len(tag["format"]) > 1 && tag["format"][1] == "template" {
						g.Printf("\nif err := validateURLTemplate(v.%s)", fn)
					} else {
						g.Import("", "net/url")
						g.Printf("\nif _, err := url.ParseRequestURI(v.%s)", fn)
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
						g.Printf("\nif v.%s != \"\" {", fn)
					}
					g.Printf("\n\nif v.%s != \"\" && !emailRegexp.MatchString(v.%[1]s) {", fn)
					g.Import("", "errors")
					g.Printf("\nreturn errors.New(`\"%s\" field must be an email address`)", yn)
					g.Printf("\n}")
					if !required {
						g.Printf("\n}")
					}
				case "runtime":
					if _, ok := field.Type.(*ast.MapType); ok {
						g.Printf("\n\nfor key := range v.%s {", fn)
						g.Printf("\nif !matchRuntimeExpr(key) {")
						g.Import("", "errors")
						g.Printf("\nreturn errors.New(`the keys of \"%s\" must be a runtime expression`)", yn)
						g.Printf("\n}")
						g.Printf("\n}")
					}
				case "regexp":
					if _, ok := field.Type.(*ast.MapType); ok {
						g.Import("", "regexp")
						g.Printf("\n\n%sRegexp := regexp.MustCompile(`%s`)", fn, tag["format"][1])
						g.Printf("\nfor key := range v.%s {", fn)
						g.Printf("\nif !%sRegexp.MatchString(v.%s) {", fn, fn)
						g.Import("", "errors")
						g.Printf("\nreturn errors.New(`the keys of \"%s\" must be match \"%s\"`)", yn, tag["format"][1])
						g.Printf("\n}")
					}
				}
				if list, ok := tag["oneof"]; ok {
					g.Printf("\n")
					if !required {
						g.Printf("\nif v.%s != \"\" {", fn)
					}
					g.Printf("\nif !isOneOf(v.%s, %#v) {", fn, list)
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

			if typ.Name.Name == "OpenAPI" {
				g.Printf("\nv.setRoot(v)")
			}

			g.Printf("\nreturn nil")
			g.Printf("\n}")
		}
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

func yamlName(field *ast.Field, t astutil.Tags) string {
	yn := t.Get("yaml")
	if yn != "" {
		return yn
	}
	return field.Names[0].Name
}

func quoteEachString(list []string) []string {
	ret := make([]string, len(list))
	for i := range list {
		ret[i] = strconv.Quote(list[i])
	}
	return ret
}

func isOpenAPIObject(genDecl *ast.GenDecl) bool {
	if genDecl.Doc == nil || len(genDecl.Doc.List) == 0 {
		return false
	}
	for _, doc := range genDecl.Doc.List {
		if doc.Text == "//+object" {
			return true
		}
	}
	return false
}
