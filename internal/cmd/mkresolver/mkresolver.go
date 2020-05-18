package main

import (
	"flag"
	"log"

	"github.com/nasa9084/go-openapi/internal/astutil"
	"github.com/nasa9084/go-openapi/internal/generator"
)

func main() {
	flag.Parse()
	g := generator.New("mkresolver.go")

	objects, err := astutil.ParseOpenAPIObjects("interfaces.go")
	if err != nil {
		log.Fatal(err)
	}

	for _, object := range objects {
		var hasReference bool
		for _, field := range object.Fields {
			if field.Tags.Get("yaml") == "$ref" {
				hasReference = true
				break
			}
		}

		if !hasReference {
			continue
		}

		log.Printf("generate %s.resolve()", object.Name)
		g.Printf("\n\nfunc (v *%s) resolve() (*%[1]s, error) {", object.Name)
		g.Printf("\nif v.reference == \"\" {")
		g.Printf("\nreturn v, nil")
		g.Printf("\n}")
		g.Printf("\nresolvedInterface, err := resolve(v.root, v.reference)")
		g.Printf("\nif err != nil {")
		g.Printf("\nreturn nil, err")
		g.Printf("\n}")
		g.Printf("\nif resolved, ok := resolvedInterface.(*%s); ok {", object.Name)
		g.Printf("\nreturn resolved, nil")
		g.Printf("\n}")
		g.Import("", "fmt")
		g.Printf("\npanic(fmt.Sprintf(\"type assertion error in resolving %%s\", v.reference))")
		g.Printf("\n}")
	}
	if err := g.Save("resolve_gen.go"); err != nil {
		log.Fatal(err)
	}
}
