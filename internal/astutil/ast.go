package astutil

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"
	"strings"
)

const openAPIObjectMarker = "//+object"

func isNotOpenAPIObject(genDecl *ast.GenDecl) bool {
	return genDecl.Doc == nil || len(genDecl.Doc.List) == 0 || genDecl.Doc.List[0].Text != openAPIObjectMarker
}

type OpenAPIObject struct {
	Name   string
	Fields []OpenAPIObjectField
}

type OpenAPIObjectField struct {
	Name string
	Type ast.Expr
	Tags Tags
}

func ParseOpenAPIObjects(filename string) ([]OpenAPIObject, error) {
	f, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	var ret []OpenAPIObject
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if isNotOpenAPIObject(genDecl) {
			continue
		}

		for _, spec := range genDecl.Specs {
			typ, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			o := OpenAPIObject{
				Name: typ.Name.Name,
			}

			st, ok := typ.Type.(*ast.StructType)
			if !ok {
				continue
			}

			for _, field := range st.Fields.List {
				o.Fields = append(o.Fields, OpenAPIObjectField{
					Name: field.Names[0].Name,
					Type: field.Type,
					Tags: ParseTags(field),
				})
			}

			ret = append(ret, o)
		}
	}
	return ret, nil
}

func (field OpenAPIObjectField) YAMLName() string {
	if yamlName := field.Tags.Get("yaml"); yamlName != "" {
		return yamlName
	}
	return field.Name
}

func (field OpenAPIObjectField) IsRequired() bool {
	return field.Tags.Get("required") != ""
}

func TypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr: // pointer
		return "*" + t.X.(*ast.Ident).Name
	case *ast.ArrayType:
		return "[]" + TypeString(t.Elt)
	case *ast.MapType:
		return "map[" + TypeString(t.Key) + "]" + TypeString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	default:
		log.Fatalf("unknown type: %s", reflect.TypeOf(t))
	}
	return ""
}

type Tags map[string][]string

func ParseTags(f *ast.Field) Tags {
	if f.Tag == nil {
		return nil
	}
	s := strings.Trim(f.Tag.Value, "`")
	if s == "" {
		return nil
	}
	t := Tags{}
	for _, tt := range strings.Fields(s) {
		kv := strings.Split(tt, ":")
		t[kv[0]] = append(t[kv[0]], strings.Split(strings.Trim(kv[1], `"`), ",")...)
	}
	return t
}

func (t Tags) Get(key string) string {
	if vs, ok := t[key]; ok {
		return vs[0]
	}
	return ""
}
