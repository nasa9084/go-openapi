package main

import (
	"flag"
	"go/ast"
	"log"
	"reflect"

	"github.com/nasa9084/go-openapi/internal/astutil"
	"github.com/nasa9084/go-openapi/internal/generator"
)

func main() {
	flag.Parse()
	g := generator.New("mksetroot.go")

	objects, err := astutil.ParseOpenAPIObjects("interfaces.go")
	if err != nil {
		log.Fatal(err)
	}

	for _, object := range objects {
		mkSetRoot(g, object)
	}
	if err := g.Save("setroot_gen.go"); err != nil {
		log.Fatal(err)
	}
}

func mkSetRoot(g *generator.Generator, object astutil.OpenAPIObject) {
	log.Printf("generate %s.setRoot()", object.Name)
	g.Printf("\n\nfunc (v *%s) setRoot(root *OpenAPI) {", object.Name)

	for _, field := range object.Fields {
		switch t := field.Type.(type) {
		case *ast.Ident, *ast.InterfaceType:
			// nothing to do
		case *ast.StarExpr: // pointer to struct
			if field.Name == "root" {
				g.Printf("\nv.root = root")
				continue
			}
			g.Printf("\nif v.%s != nil {", field.Name)
			g.Printf("\nv.%s.setRoot(root)", field.Name)
			g.Printf("\n}")
		case *ast.ArrayType:
			switch t.Elt.(type) {
			case *ast.Ident, *ast.InterfaceType:
				// nothing to do
			case *ast.StarExpr:
				g.Printf("\nfor i := range v.%s {", field.Name)
				g.Printf("\nv.%s[i].setRoot(root)", field.Name)
				g.Printf("\n}")
			default:
				log.Print(reflect.TypeOf(t.Elt))
			}
		case *ast.MapType:
			switch tv := t.Value.(type) {
			case *ast.Ident, *ast.InterfaceType:
				// nothing to do
			case *ast.StarExpr:
				g.Printf("\nfor k := range v.%s {", field.Name)
				g.Printf("\nv.%s[k].setRoot(root)", field.Name)
				g.Printf("\n}")
			case *ast.ArrayType:
				switch tv.Elt.(type) {
				case *ast.StarExpr:
					g.Printf("\nfor k := range v.%s {", field.Name)
					g.Printf("\nfor i := range v.%s[k] {", field.Name)
					g.Printf("\nv.%s[k][i].setRoot(root)", field.Name)
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
