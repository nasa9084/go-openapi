package main

import (
	"flag"
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

		generateGetter(g, object)
	}

	if err := g.Save("getter_gen.go"); err != nil {
		log.Fatal(err)
	}
}

func generateGetter(g *generator.Generator, object astutil.OpenAPIObject) {
	var hasReference bool

	for _, field := range object.Fields {
		if field.Name == "reference" {
			hasReference = true
			break
		}
	}

	if hasReference {
		log.Printf("  %s has reference field", object.Name)
	}

	for _, field := range object.Fields {
		ft := astutil.TypeString(field.Type)

		if object.Name == expose(field.Name) {
			continue
		}

		log.Printf("  generate %s.%s()", object.Name, expose(field.Name))

		g.Printf("\n\nfunc (v *%s) %s() %s {", object.Name, expose(field.Name), ft)

		if field.Name != "root" && hasReference {
			g.Printf("\nif v.reference != \"\" {")
			g.Printf("\nresolved, err := v.resolve()")
			g.Printf("\nif err != nil {")
			g.Printf("\npanic(err)")
			g.Printf("\n}")
			g.Printf("\nreturn resolved.%s", field.Name)
			g.Printf("\n}")
		}

		if _, ok := field.Type.(*ast.StarExpr); ok {
			g.Printf("\nif v.%s == nil {", field.Name)
			// trim "*"
			g.Printf("\nreturn &%s{}", ft[1:])
			g.Printf("\n}")
		}

		g.Printf("\nreturn v.%s", field.Name)
		g.Printf("\n}")
	}
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
