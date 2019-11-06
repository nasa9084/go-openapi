package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
	"unicode"

	"github.com/nasa9084/go-openapi/internal/astutil"
	"github.com/nasa9084/go-openapi/internal/generator"
)

func init() {
	flag.Parse()
}

func main() {
	g := generator.New("mkgetter.go")

	f, err := parser.ParseFile(token.NewFileSet(), "interfaces.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
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
			log.Printf("generate getters for %s", typ.Name.Name)

			var hasReference bool

			for _, field := range st.Fields.List {
				if field.Names[0].Name == "reference" {
					hasReference = true
					break
				}
			}

			if hasReference {
				log.Printf("  %s has reference field", typ.Name.Name)
			}

			for _, field := range st.Fields.List {
				fn := field.Names[0].Name
				ft := astutil.TypeString(field.Type)

				if typ.Name.Name == expose(fn) {
					continue
				}

				log.Printf("  generate %s.%s()", typ.Name.Name, expose(fn))

				g.Printf("\n\nfunc (v *%s) %s() %s {", typ.Name.Name, expose(fn), ft)

				if fn != "root" && hasReference {
					g.Printf("\nif v.reference != \"\" {")
					g.Printf("\nresolved, err := v.resolve()")
					g.Printf("\nif err != nil {")
					g.Printf("\npanic(err)")
					g.Printf("\n}")
					g.Printf("\nreturn resolved.%s", fn)
					g.Printf("\n}")
				}

				if strings.HasPrefix(ft, "*") {
					g.Printf("\nif v.%s == nil {", fn)
					// trim "*"
					g.Printf("\nreturn &%s{}", ft[1:])
					g.Printf("\n}")
				}

				g.Printf("\nreturn v.%s", fn)
				g.Printf("\n}")
			}
		}
	}

	if err := g.Save("getter_gen.go"); err != nil {
		log.Fatal(err)
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
