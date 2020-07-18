package setroot

import (
	"go/ast"
	"log"
	"reflect"

	"github.com/nasa9084/go-openapi/internal/astutil"
	"github.com/nasa9084/go-openapi/internal/generator"
)

const (
	generatorName = "SetRootGenerator"
	saveTo        = "setroot_gen.go"
)

var ignoreFields = []string{"resolved"}

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
		log.Printf("generate %s.setRoot()", object.Name)
		g.mkSetRoot(object)
	}

	return g.Save(saveTo)
}

func (g *Generator) mkSetRoot(object astutil.OpenAPIObject) {
	g.Printf("\n\nfunc (v *%s) setRoot(root *OpenAPI) {", object.Name)

	for _, field := range object.Fields {
		if oneOf(field.Name, ignoreFields) {
			continue
		}

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
				if _, ok := tv.Elt.(*ast.StarExpr); ok {
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

func oneOf(s string, list []string) bool {
	for _, t := range list {
		if s == t {
			return true
		}
	}

	return false
}
