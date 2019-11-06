package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/nasa9084/go-openapi/internal/astutil"
	"github.com/nasa9084/go-openapi/internal/generator"
)

func init() {
	flag.Parse()
}

func main() {
	g := generator.New("mkresolver.go")

	f, err := parser.ParseFile(token.NewFileSet(), "interfaces.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if !isOpenAPIObject(genDecl) {
			log.Printf("%v is not an openapi object. skip", genDecl.Specs[0].(*ast.TypeSpec).Name.Name)
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

			var hasReference bool
			for _, field := range st.Fields.List {
				if astutil.ParseTags(field).Get("yaml") == "$ref" {
					hasReference = true
					break
				}
			}

			if !hasReference {
				continue
			}

			log.Printf("generate %s.resolve()", typ.Name.Name)
			g.Printf("\n\nfunc (v *%s) resolve() (*%[1]s, error) {", typ.Name.Name)
			g.Printf("\nif v.reference == \"\" {")
			g.Printf("\nreturn v, nil")
			g.Printf("\n}")
			g.Printf("\nresolvedInterface, err := resolve(v.root, v.reference)")
			g.Printf("\nif err != nil {")
			g.Printf("\nreturn nil, err")
			g.Printf("\n}")
			g.Printf("\nif resolved, ok := resolvedInterface.(*%s); ok {", typ.Name.Name)
			g.Printf("\nreturn resolved, nil")
			g.Printf("\n}")
			g.Import("", "fmt")
			g.Printf("\npanic(fmt.Sprintf(\"type assertion error in resolving %%s\", v.reference))")
			g.Printf("\n}")
		}
	}
	if err := g.Save("resolve_gen.go"); err != nil {
		log.Fatal(err)
	}

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
