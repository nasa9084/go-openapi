package main

import (
	"flag"
	"fmt"
	"go/ast"
	"log"
	"strings"
	"unicode"

	"github.com/nasa9084/go-openapi/internal/astutil"
	"github.com/nasa9084/go-openapi/internal/generator"
)

func main() {
	flag.Parse()

	g := generator.New("mkgetter.go")

	objects, err := astutil.ParseOpenAPIObjects("interfaces.go")
	if err != nil {
		log.Fatal(err)
	}

	for _, object := range objects {
		log.Printf("generate getters for %s", object.Name)

		if err := generateGetters(g, object); err != nil {
			log.Fatal(err)
		}
	}

	if err := g.Save("getter_gen.go"); err != nil {
		log.Fatal(err)
	}
}

func generateGetters(g *generator.Generator, object astutil.OpenAPIObject) error {
	for _, field := range object.Fields {
		if object.Name == expose(field.Name) && (object.Name != "OpenAPI" && field.Name != "openapi") {
			if err := generateMapGetter(g, field, "Get"); err != nil {
				return err
			}

			continue
		}

		if field.Name == "extension" {
			if err := generateMapGetter(g, field, "Extension"); err != nil {
				return err
			}

			continue
		}

		if err := generateGetter(g, field); err != nil {
			return err
		}
	}

	return nil
}

func generateGetter(g *generator.Generator, field astutil.OpenAPIObjectField) error {
	log.Printf("  generate %s.%s()", field.ParentObject.Name, expose(field.Name))

	fieldType := astutil.TypeString(field.Type)

	g.Printf("\n\nfunc (v *%s) %s() %s {", field.ParentObject.Name, expose(field.Name), fieldType)
	defer g.Printf("\n}")

	if field.ParentObject.HasReference {
		resolveReference(g)
	}

	if _, ok := field.Type.(*ast.StarExpr); ok {
		g.Printf("\nif v.%s == nil {", field.Name)
		g.Printf("\nreturn &%s{}", strings.TrimPrefix(fieldType, "*"))
		g.Printf("\n}")
	}

	g.Printf("\nreturn v.%s", field.Name)

	return nil
}

func generateMapGetter(g *generator.Generator, field astutil.OpenAPIObjectField, fnName string) error {
	log.Printf("  generate %s.%s()", field.ParentObject.Name, fnName)

	var (
		keyType, valType string
		isValueStarExpr  bool
	)

	m, ok := field.Type.(*ast.MapType)
	if !ok {
		return fmt.Errorf("%s.%s must be a map type", field.ParentObject.Name, field.Name)
	}

	keyType = astutil.TypeString(m.Key)
	valType = astutil.TypeString(m.Value)

	if _, ok := m.Value.(*ast.StarExpr); ok {
		isValueStarExpr = true
	}

	g.Printf("\n\nfunc (v *%s) %s(key %s) %s {", field.ParentObject.Name, fnName, keyType, valType)
	defer g.Printf("\n}")

	if field.ParentObject.HasReference {
		resolveReference(g)
	}

	if valType != "interface{}" {
		g.Printf("\nif val, ok := v.%s[key]; ok {", field.Name)
		g.Printf("\nreturn val")
		g.Printf("\n}")
		g.Printf("\nreturn ")

		if isValueStarExpr {
			g.Printf("&")
		}

		g.Printf("%s{}", strings.TrimPrefix(valType, "*"))
	} else {
		g.Printf("\nreturn v.%s[key]", field.Name)
	}

	return nil
}

func resolveReference(g *generator.Generator) {
	g.Printf("\nif v.reference != \"\" {")
	g.Printf("\nresolved, err := v.resolve()")
	g.Printf("\nif err != nil {")
	g.Printf("\npanic(err)")
	g.Printf("\n}")
	g.Printf("\nv = resolved")
	g.Printf("\n}")
}

func expose(ident string) string {
	if strings.HasSuffix(ident, "_") {
		ident = ident[:len(ident)-1]
	}

	switch ident {
	case "openapi":
		return "OpenAPI"
	case "url":
		return "URL"
	case "xml":
		return "XML"
	}

	rident := []rune(ident)

	return string(append([]rune{unicode.ToUpper(rident[0])}, rident[1:]...))
}
