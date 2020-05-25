package openapi

import (
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestXMLExampleUnmarshalYAML(t *testing.T) {
	t.Run("basic string", testXMLUnmarshalYAMLBasicString)
	t.Run("attribute prefix namespace", testXMLUnmarshalYAMLAttributePrefixNamespace)
	t.Run("wrapped array", testXMLUnmarshalYAMLWrappedArray)
}

func testXMLUnmarshalYAMLBasicString(t *testing.T) {
	yml := `animals:
  type: string
  xml:
    name: animal`
	var got map[string]*Schema
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := map[string]*Schema{
		"animals": {
			type_: "string",
			xml: &XML{
				name: "animal",
			},
		},
	}
	assertEqual(t, got, want)
}

func testXMLUnmarshalYAMLAttributePrefixNamespace(t *testing.T) {
	yml := `Person:
  type: object
  properties:
    id:
      type: integer
      format: int32
      xml:
        attribute: true
    name:
      type: string
      xml:
        namespace: http://example.com/schema/sample
        prefix: sample`
	var got map[string]*Schema
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := map[string]*Schema{
		"Person": {
			type_: "object",
			properties: map[string]*Schema{
				"id": {
					type_:  "integer",
					format: "int32",
					xml:    &XML{attribute: true},
				},
				"name": {
					type_: "string",
					xml: &XML{
						namespace: "http://example.com/schema/sample",
						prefix:    "sample",
					},
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func testXMLUnmarshalYAMLWrappedArray(t *testing.T) {
	yml := `animals:
  type: array
  items:
    type: string
  xml:
    wrapped: true`
	var got map[string]*Schema
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := map[string]*Schema{
		"animals": {
			type_: "array",
			items: &Schema{type_: "string"},
			xml:   &XML{wrapped: true},
		},
	}
	assertEqual(t, got, want)
}

func TestXMLUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want XML
	}{
		{
			yml: `x-foo: bar`,
			want: XML{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got XML
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestXMLUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &XML{})
			assertSameError(t, got, tt.want)
		})
	}
}
