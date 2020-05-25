package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestParameterExampleUnmarshalYAML(t *testing.T) {
	t.Run("header parameter", func(t *testing.T) {
		yml := `name: token
in: header
description: token to be passed as a header
required: true
schema:
  type: array
  items:
    type: integer
    format: int64
style: simple`
		var got Parameter
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := Parameter{
			name:        "token",
			in:          "header",
			description: "token to be passed as a header",
			required:    true,
			schema: &Schema{
				type_: "array",
				items: &Schema{
					type_:  "integer",
					format: "int64",
				},
			},
			style: "simple",
		}
		assertEqual(t, got, want)
	})
	t.Run("path parameter", func(t *testing.T) {
		yml := `name: username
in: path
description: username to fetch
required: true
schema:
  type: string`
		var got Parameter
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := Parameter{
			name:        "username",
			in:          "path",
			description: "username to fetch",
			required:    true,
			schema: &Schema{
				type_: "string",
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("optional query parameter", func(t *testing.T) {
		yml := `name: id
in: query
description: ID of the object to fetch
required: false
schema:
  type: array
  items:
    type: string
style: form
explode: true`
		var got Parameter
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := Parameter{
			name:        "id",
			in:          "query",
			description: "ID of the object to fetch",
			required:    false,
			schema: &Schema{
				type_: "array",
				items: &Schema{
					type_: "string",
				},
			},
			style:   "form",
			explode: true,
		}
		assertEqual(t, got, want)
	})
	t.Run("free form", func(t *testing.T) {
		yml := `in: query
name: freeForm
schema:
  type: object
  additionalProperties:
    type: integer
style: form`
		var got Parameter
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := Parameter{
			in:   "query",
			name: "freeForm",
			schema: &Schema{
				type_: "object",
				additionalProperties: &Schema{
					type_: "integer",
				},
			},
			style: "form",
		}
		assertEqual(t, got, want)
	})
	t.Run("complex parameter", func(t *testing.T) {
		yml := `in: query
name: coordinates
content:
  application/json:
    schema:
      type: object
      required:
        - lat
        - long
      properties:
        lat:
          type: number
        long:
          type: number`
		var got Parameter
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := Parameter{
			in:   "query",
			name: "coordinates",
			content: map[string]*MediaType{
				"application/json": {
					schema: &Schema{
						type_:    "object",
						required: []string{"lat", "long"},
						properties: map[string]*Schema{
							"lat": {
								type_: "number",
							},
							"long": {
								type_: "number",
							},
						},
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
}

func TestParameterUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Parameter
	}{
		{
			yml: `name: foo
in: path`,
			want: Parameter{
				name: "foo",
				in:   "path",
			},
		},
		{
			yml: `name: foo
in: path
deprecated: true`,
			want: Parameter{
				name:       "foo",
				in:         "path",
				deprecated: true,
			},
		},
		{
			yml: `name: foo
in: path
allowReserved: true`,
			want: Parameter{
				name:          "foo",
				in:            "path",
				allowReserved: true,
			},
		},
		{
			yml: `name: foo
in: path
examples:
  foo:
    value: bar`,
			want: Parameter{
				name: "foo",
				in:   "path",
				examples: map[string]*Example{
					"foo": {
						value: "bar",
					},
				},
			},
		},
		{
			yml: `name: foo
in: path
allowEmptyValue: true`,
			want: Parameter{
				name:            "foo",
				in:              "path",
				allowEmptyValue: true,
			},
		},
		{
			yml: `name: foo
in: path
x-foo: bar`,
			want: Parameter{
				name: "foo",
				in:   "path",
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
		{
			yml: `$ref: "#/components/parameters/Foo"`,
			want: Parameter{
				reference: "#/components/parameters/Foo",
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Parameter
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestParameterUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `in: path`,
			want: ErrRequired("name"),
		},
		{
			yml:  `name: foo`,
			want: ErrRequired("in"),
		},
		{
			yml: `name: foo
in: bar`,
			want: errors.New(`"in" field must be one of ["query", "header", "path", "cookie"]`),
		},
		{
			yml: `in: path
name: foo
schema: hoge`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `in: path
name: foo
examples: hoge`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `in: path
name: foo
content: hoge`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `name: namename
in: query
foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Parameter{})
			assertSameError(t, got, tt.want)
		})
	}
}
