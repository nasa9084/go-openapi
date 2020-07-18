package unmarshalyaml

import (
	"errors"
	"go/ast"
	"log"
	"strconv"
	"strings"

	"github.com/nasa9084/go-openapi/internal/astutil"
	"github.com/nasa9084/go-openapi/internal/generator"
)

const (
	generatorName = "UnmarshalYAMLGenerator"
	saveTo        = "unmarshalyaml_gen.go"
)

type Generator struct {
	*generator.Generator

	objects []astutil.OpenAPIObject
}

func NewGenerator(objects []astutil.OpenAPIObject) *Generator {
	return &Generator{
		Generator: generator.New(generatorName),

		objects: objects,
	}
}

func (g *Generator) Generate() error {
	for _, object := range g.objects {
		log.Printf("generate %s.Unmarshal()", object.Name)

		if err := g.generateUnmarshalYAML(object); err != nil {
			return err
		}
	}

	return g.Save(saveTo)
}

func (g *Generator) generateUnmarshalYAML(object astutil.OpenAPIObject) error {
	g.Printf("\n\nfunc (v *%s) UnmarshalYAML(unmarshal func(interface{}) error) error {", object.Name)
	g.Printf("\nvar proxy map[string]rawMessage")
	g.Printf("\nif err := unmarshal(&proxy); err != nil {")
	g.Printf("\nreturn err")
	g.Printf("\n}")

	for _, field := range object.Fields {
		if field.YAMLName() == "$ref" {
			g.generateReferenceUnmarshal()
			break
		}
	}

	var noUnknown bool

	for _, field := range object.Fields {
		if yamlName := field.YAMLName(); yamlName == "-" || yamlName == "$ref" {
			continue
		}

		if isInline(field.Tags) {
			if g.generateInlineUnmarshal(field) {
				// only overwrites when true
				noUnknown = true
			}

			continue
		}

		g.Printf("\n\n")
		g.generateUnmarshalField(field)

		if err := g.generateFormatValidation(field); err != nil {
			return err
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

	return nil
}

func (g *Generator) generateReferenceUnmarshal() {
	g.Printf("\nif p, ok := proxy[\"$ref\"]; ok {")
	g.Printf("\nvar referenceVal string")
	g.Printf("\nif err := p.unmarshal(&referenceVal); err != nil {")
	g.Printf("\nreturn err")
	g.Printf("\n}")
	g.Printf("\nv.reference = referenceVal")
	g.Printf("\ndelete(proxy, \"$ref\")")
	g.Printf("\nreturn nil")
	g.Printf("\n}")
}

func (g *Generator) generateInlineUnmarshal(field astutil.OpenAPIObjectField) (noUnknown bool) {
	ft, ok := field.Type.(*ast.MapType)
	if !ok {
		log.Fatalf("expected map for inline %s but %s", field.YAMLName(), field.Type)
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
	g.Printf("\nif err := val.unmarshal(&%sv); err != nil {", field.Name)
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

	return noUnknown
}

func (g *Generator) generateUnmarshalField(field astutil.OpenAPIObjectField) {
	if field.IsRequired() {
		g.Printf("%sUnmarshal, ok := proxy[\"%s\"]", field.Name, field.YAMLName())
		g.Printf("\nif !ok {")
		g.Printf("\nreturn ErrRequired(%s)", strconv.Quote(field.YAMLName()))
		g.Printf("\n}")
	} else {
		g.Printf("if %sUnmarshal, ok := proxy[\"%s\"]; ok {", field.Name, field.YAMLName())
		defer g.Printf("\n}")
	}

	g.Printf("\nvar %sVal %s", field.Name, strings.TrimPrefix(field.TypeString(), "*"))

	g.Printf("\nif err := %sUnmarshal.unmarshal(&%[1]sVal); err != nil {", field.Name)
	g.Printf("\nreturn err")
	g.Printf("\n}")

	g.Printf("\nv.%s = ", field.Name)

	if field.IsPointerType() {
		g.Printf("&")
	}

	if field.IsStringType() {
		g.Printf("strings.TrimSuffix(")
	}

	g.Printf("%sVal", field.Name)

	if field.IsStringType() {
		g.Printf(`, "\n")`)
	}

	g.Printf("\ndelete(proxy, `%s`)", field.YAMLName())
}

func (g *Generator) generateFormatValidation(field astutil.OpenAPIObjectField) error {
	switch field.Tags.Get("format") {
	case "semver":
		g.Printf("\n\nif !isValidSemVer(v.%s) {", field.Name)
		g.Import("", "errors")
		g.Printf("\nreturn errors.New(`\"%s\" field must be a valid semantic version but not`)", field.YAMLName())
		g.Printf("\n}")
	case "url":
		g.generateURLValidation(field)
	case "email":
		g.generateEmailValidation(field)
	case "runtime":
		if _, ok := field.Type.(*ast.MapType); !ok {
			return errors.New("`runtime` validation constraints can only be used for map value")
		}

		g.Printf("\n\nfor key := range v.%s {", field.Name)
		g.Printf("\nif !matchRuntimeExpr(key) {")
		g.Import("", "errors")
		g.Printf("\nreturn errors.New(`the keys of \"%s\" must be a runtime expression`)", field.YAMLName())
		g.Printf("\n}")
		g.Printf("\n}")
	case "regexp":
		if _, ok := field.Type.(*ast.MapType); !ok {
			return errors.New("`regexp` validation contraints can only be used for map value")
		}

		g.Import("", "regexp")
		g.Printf("\n\n%sRegexp := regexp.MustCompile(`%s`)", field.Name, field.Tags["format"][1])
		g.Printf("\nfor key := range v.%s {", field.Name)
		g.Printf("\nif !%sRegexp.MatchString(v.%s) {", field.Name, field.Name)
		g.Import("", "errors")
		g.Printf("\nreturn errors.New(`the keys of \"%s\" must be match \"%s\"`)", field.YAMLName(), field.Tags["format"][1])
		g.Printf("\n}")
	}

	if list, ok := field.Tags["oneof"]; ok {
		g.generateOneOfValidation(field, list)
	}

	return nil
}

func (g *Generator) generateURLValidation(field astutil.OpenAPIObjectField) {
	g.Printf("\n")

	if !field.IsRequired() {
		g.Printf("\nif v.%s != \"\" {", field.Name)
		defer g.Printf("\n}")
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
}

func (g *Generator) generateEmailValidation(field astutil.OpenAPIObjectField) {
	g.Printf("\n")

	if !field.IsRequired() {
		g.Printf("\nif v.%s != \"\" {", field.Name)
		defer g.Printf("\n}")
	}

	g.Printf("\nif v.%s != \"\" && !emailRegexp.MatchString(v.%[1]s) {", field.Name)
	g.Import("", "errors")
	g.Printf("\nreturn errors.New(`\"%s\" field must be an email address`)", field.YAMLName())
	g.Printf("\n}")
}

func (g *Generator) generateOneOfValidation(field astutil.OpenAPIObjectField, list []string) {
	g.Printf("\n")

	if !field.IsRequired() {
		g.Printf("\nif v.%s != \"\" {", field.Name)
		defer g.Printf("\n}")
	}

	g.Printf("\nif !isOneOf(v.%s, %#v) {", field.Name, list)
	g.Import("", "errors")
	g.Printf("\nreturn errors.New(`\"%s\" field must be one of [%s]`)",
		field.YAMLName(), strings.Join(quoteEachString(list), ", "))
	g.Printf("\n}")
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

func quoteEachString(list []string) []string {
	ret := make([]string, len(list))

	for i := range list {
		ret[i] = strconv.Quote(list[i])
	}

	return ret
}