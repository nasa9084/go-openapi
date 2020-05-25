package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestResponseExampleUnmarshalYAML(t *testing.T) {
	t.Run("array of complex type", func(t *testing.T) {
		yml := `description: A complex object array response
content:
  application/json:
    schema:
      type: array
      items:
        $ref: '#/components/schemas/VeryComplexType'`

		var got Response
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}

		want := Response{
			description: "A complex object array response",
			content: map[string]*MediaType{
				"application/json": {
					schema: &Schema{
						type_: "array",
						items: &Schema{
							reference: "#/components/schemas/VeryComplexType",
						},
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("string", func(t *testing.T) {
		yml := `description: A simple string response
content:
  text/plain:
    schema:
      type: string`

		var got Response
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}

		want := Response{
			description: "A simple string response",
			content: map[string]*MediaType{
				"text/plain": {
					schema: &Schema{
						type_: "string",
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("plain text with headers", func(t *testing.T) {
		yml := `description: A simple string response
content:
  text/plain:
    schema:
      type: string
    example: 'whoa!'
headers:
  X-Rate-Limit-Limit:
    description: The number of allowed requests in the current period
    schema:
      type: integer
  X-Rate-Limit-Remaining:
    description: The number of remaining requests in the current period
    schema:
      type: integer
  X-Rate-Limit-Reset:
    description: The number of seconds left in the current period
    schema:
      type: integer`

		var got Response
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}

		want := Response{
			description: "A simple string response",
			content: map[string]*MediaType{
				"text/plain": {
					schema: &Schema{
						type_: "string",
					},
					example: "whoa!",
				},
			},
			headers: map[string]*Header{
				"X-Rate-Limit-Limit": {
					description: "The number of allowed requests in the current period",
					schema: &Schema{
						type_: "integer",
					},
				},
				"X-Rate-Limit-Remaining": {
					description: "The number of remaining requests in the current period",
					schema: &Schema{
						type_: "integer",
					},
				},
				"X-Rate-Limit-Reset": {
					description: "The number of seconds left in the current period",
					schema: &Schema{
						type_: "integer",
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("no return value", func(t *testing.T) {
		yml := `description: object created`
		var got Response
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}

		want := Response{
			description: "object created",
		}
		assertEqual(t, got, want)
	})
}

func TestResponseUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Response
	}{
		{
			yml: `$ref: "#/components/responses/Foo"`,
			want: Response{
				reference: "#/components/responses/Foo",
			},
		},
		{
			yml: `description: foobar
x-foo: bar`,
			want: Response{
				description: "foobar",
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Response
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestResponseUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `headers: {}`,
			want: ErrRequired("description"),
		},
		{
			yml: `description: foo
headers: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `description: foo
content: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `description: foo
links: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `description: foo
foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Response{})
			assertSameError(t, got, tt.want)
		})
	}
}
