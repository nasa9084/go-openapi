package main

import (
	"flag"
	"log"
	"unicode"

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
		if !object.HasReference {
			continue
		}

		generateResolve(g, object)
		generateResolveLocal(g, object)
	}

	if err := g.Save("resolve_gen.go"); err != nil {
		log.Fatal(err)
	}
}

func generateResolve(g *generator.Generator, object astutil.OpenAPIObject) {
	log.Printf("generate %s.resolve()", object.Name)

	g.Printf("\n\nfunc (v *%s) resolve() (*%[1]s, error) {", object.Name)
	defer g.Printf("\n}")

	g.Printf("\nif v.reference == \"\" {")
	g.Printf("\nreturn v, nil")
	g.Printf("\n}")
	g.Printf("\n\nif v.resolved != nil {")
	g.Printf("\nreturn v.resolved, nil")
	g.Printf("\n}")

	g.Printf("\n\nif strings.HasPrefix(v.reference, `#/`) {")
	g.Printf("\nreturn v.resolveLocal()")
	g.Printf("\n}")

	g.Printf("\n\nreturn nil, ErrCannotResolved(v.reference, `not supported reference type`)")
}

func generateResolveLocal(g *generator.Generator, object astutil.OpenAPIObject) {
	log.Printf("  generate %s.resolveLocal()", object.Name)

	g.Printf("\n\nfunc (v *%s) resolveLocal() (*%[1]s, error) {", object.Name)
	defer g.Printf("\n}")

	g.Import("", "strings")
	g.Printf("\nprefix := `#/components/%s/`", plural(unexpose(object.Name)))

	g.Printf("\n\nif !strings.HasPrefix(v.reference, prefix) {")
	g.Import("", "fmt")
	//nolint:lll // cannot fix
	g.Printf("\nreturn nil, ErrCannotResolved(v.reference, fmt.Sprintf(\"local %s reference must begin with `%%s`\", prefix))", object.Name)
	g.Printf("\n}")

	g.Printf("\n\nkey := strings.TrimPrefix(v.reference, prefix)")

	g.Printf("\n\nif resolved, ok := v.root.components.%s[key]; ok {", plural(unexpose(object.Name)))
	g.Printf("\nreturn resolved, nil")
	g.Printf("\n}")

	g.Printf("\n\nreturn nil, ErrCannotResolved(v.reference, `not found`)")
}

func unexpose(ident string) string {
	rident := []rune(ident)

	return string(append([]rune{unicode.ToLower(rident[0])}, rident[1:]...))
}

func plural(ident string) string {
	rident := []rune(ident)

	switch rident[len(rident)-2] {
	case 'a', 'e', 'i', 'o', 'u':
		return ident + "s"
	}

	if rident[len(rident)-1] != 'y' {
		return ident + "s"
	}

	return string(append(rident[:len(rident)-1], []rune("ies")...))
}
