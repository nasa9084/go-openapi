package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestResponsesExampleUnmarshalYAML(t *testing.T) {
	yml := `'200':
  description: a pet to be returned
  content:
    application/json:
      schema:
        $ref: '#/components/schemas/Pet'
default:
  description: Unexpected error
  content:
    application/json:
      schema:
        $ref: '#/components/schemas/ErrorModel'`

	var got Responses
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := Responses{
		responses: map[string]*Response{
			"200": {
				description: "a pet to be returned",
				content: map[string]*MediaType{
					"application/json": {
						schema: &Schema{
							reference: "#/components/schemas/Pet",
						},
					},
				},
			},
			"default": {
				description: "Unexpected error",
				content: map[string]*MediaType{
					"application/json": {
						schema: &Schema{
							reference: "#/components/schemas/ErrorModel",
						},
					},
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func TestResponsesUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Responses
	}{
		{
			yml: `"200":
  description: foobar`,
			want: Responses{
				responses: map[string]*Response{
					"200": {
						description: "foobar",
					},
				},
			},
		},
		{
			yml: `x-foo: bar`,
			want: Responses{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Responses
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestResponsesUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `"200": foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `"600":
  description: foobar`,
			want: ErrUnknownKey("600"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Responses{})
			assertSameError(t, got, tt.want)
		})
	}
}
