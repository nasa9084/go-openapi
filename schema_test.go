package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestSchemaValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Schema{}, nil},
	}
	testValidater(t, candidates)
}

func TestGoType(t *testing.T) {
	candidates := []struct {
		schema   openapi.Schema
		expected string
	}{
		{
			openapi.Schema{},
			"",
		},
		{
			openapi.Schema{Type: "string"},
			"string",
		},
		{
			openapi.Schema{Type: "integer"},
			"int",
		},
		{
			openapi.Schema{Type: "integer", Format: "int32"},
			"int32",
		},
		{
			openapi.Schema{Ref: "#/components/schemas/foo"},
			"foo",
		},
		{
			openapi.Schema{
				Type: "object",
				Properties: map[string]*openapi.Schema{
					"foo": &openapi.Schema{
						Type: "string",
					},
				},
			},
			`struct {
foo string
}`,
		},
		{
			openapi.Schema{
				Type: "array",
				Items: &openapi.Schema{
					Type: "string",
				},
			},
			"[]string",
		},
	}
	for _, c := range candidates {
		out := c.schema.GoType()
		if out != c.expected {
			t.Errorf("%s != %s", out, c.expected)
		}
	}
}
