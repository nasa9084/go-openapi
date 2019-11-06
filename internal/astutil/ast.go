package astutil

import (
	"go/ast"
	"log"
	"reflect"
	"strings"
)

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
