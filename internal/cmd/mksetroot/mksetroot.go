package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"

	"github.com/nasa9084/go-openapi/internal/generator"
)

func init() {
	flag.Parse()
}

func main() {
	g := generator.New("mksetroot.go")

	f, err := parser.ParseFile(token.NewFileSet(), "interfaces.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Doc == nil || len(genDecl.Doc.List) == 0 || genDecl.Doc.List[0].Text != "//+object" {
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

			log.Printf("generate %s.setRoot()", typ.Name.Name)
			g.Printf("\n\nfunc (v *%s) setRoot(root *OpenAPI) {", typ.Name.Name)

			for _, field := range st.Fields.List {
				switch t := field.Type.(type) {
				case *ast.Ident, *ast.InterfaceType:
					// nothing to do
				case *ast.StarExpr: // pointer to struct
					if field.Names[0].Name == "root" {
						g.Printf("\nv.root = root")
						continue
					}
					g.Printf("\nif v.%s != nil {", field.Names[0].Name)
					g.Printf("\nv.%s.setRoot(root)", field.Names[0].Name)
					g.Printf("\n}")
				case *ast.ArrayType:
					switch t.Elt.(type) {
					case *ast.Ident, *ast.InterfaceType:
						// nothing to do
					case *ast.StarExpr:
						g.Printf("\nfor i := range v.%s {", field.Names[0].Name)
						g.Printf("\nv.%s[i].setRoot(root)", field.Names[0].Name)
						g.Printf("\n}")
					default:
						log.Print(reflect.TypeOf(t.Elt))
					}
				case *ast.MapType:
					switch tv := t.Value.(type) {
					case *ast.Ident, *ast.InterfaceType:
						// nothing to do
					case *ast.StarExpr:
						g.Printf("\nfor k := range v.%s {", field.Names[0].Name)
						g.Printf("\nv.%s[k].setRoot(root)", field.Names[0].Name)
						g.Printf("\n}")
					case *ast.ArrayType:
						switch tv.Elt.(type) {
						case *ast.StarExpr:
							g.Printf("\nfor k := range v.%s {", field.Names[0].Name)
							g.Printf("\nfor i := range v.%s[k] {", field.Names[0].Name)
							g.Printf("\nv.%s[k][i].setRoot(root)", field.Names[0].Name)
							g.Printf("\n}")
							g.Printf("\n}")
						}
					default:
						log.Print(reflect.TypeOf(t.Value))
					}
				default:
					log.Printf("%s %s", field.Type, reflect.TypeOf(field.Type))
				}
			}

			g.Printf("\n}")
		}
	}
	if err := g.Save("setroot_gen.go"); err != nil {
		log.Fatal(err)
	}
}
