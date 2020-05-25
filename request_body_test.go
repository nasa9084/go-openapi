package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestRequestBodyExampleUnmarshalYAML(t *testing.T) {
	t.Run("with a referenced model", func(t *testing.T) {
		yml := `description: user to add to the system
content:
  'application/json':
    schema:
      $ref: '#/components/schemas/User'
    examples:
      user:
        summary: User Example
        externalValue: 'http://foo.bar/examples/user-example.json'
  'application/xml':
    schema:
      $ref: '#/components/schemas/User'
    examples:
      user:
        summary: User Example in XML
        externalValue: 'http://foo.bar/examples/user-example.xml'
  'text/plain':
    examples:
      user:
        summary: User example in text plain format
        externalValue: 'http://foo.bar/examples/user-example.txt'
  '*/*':
    examples:
      user:
        summary: User example in other format
        externalValue: 'http://foo.bar/examples/user-example.whatever'`
		var got RequestBody
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := RequestBody{
			description: "user to add to the system",
			content: map[string]*MediaType{
				"application/json": {
					schema: &Schema{
						reference: "#/components/schemas/User",
					},
					examples: map[string]*Example{
						"user": {
							summary:       "User Example",
							externalValue: "http://foo.bar/examples/user-example.json",
						},
					},
				},
				"application/xml": {
					schema: &Schema{
						reference: "#/components/schemas/User",
					},
					examples: map[string]*Example{
						"user": {
							summary:       "User Example in XML",
							externalValue: "http://foo.bar/examples/user-example.xml",
						},
					},
				},
				"text/plain": {
					examples: map[string]*Example{
						"user": {
							summary:       "User example in text plain format",
							externalValue: "http://foo.bar/examples/user-example.txt",
						},
					},
				},
				"*/*": {
					examples: map[string]*Example{
						"user": {
							summary:       "User example in other format",
							externalValue: "http://foo.bar/examples/user-example.whatever",
						},
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("array of string", func(t *testing.T) {
		yml := `description: user to add to the system
required: true
content:
  text/plain:
    schema:
      type: array
      items:
        type: string`
		var got RequestBody
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := RequestBody{
			description: "user to add to the system",
			required:    true,
			content: map[string]*MediaType{
				"text/plain": {
					schema: &Schema{
						type_: "array",
						items: &Schema{
							type_: "string",
						},
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
}

func TestRequestBodyUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want RequestBody
	}{
		{
			yml: `$ref: "#/components/requestBodies/foo"`,
			want: RequestBody{
				reference: "#/components/requestBodies/foo",
			},
		},
		{
			yml: `content: {}
x-foo: bar`,
			want: RequestBody{
				content: map[string]*MediaType{},
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got RequestBody
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestRequestBodyUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `description: foo`,
			want: ErrRequired("content"),
		},
		{
			yml:  `content: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `content:
  application/json: {}
foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &RequestBody{})
			assertSameError(t, got, tt.want)
		})
	}
}
