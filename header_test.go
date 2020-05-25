package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestHeaderExampleUnmarshalYAML(t *testing.T) {
	yml := `description: The number of allowed requests in the current period
schema:
  type: integer`

	var got Header
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := Header{
		description: "The number of allowed requests in the current period",
		schema: &Schema{
			type_: "integer",
		},
	}
	assertEqual(t, got, want)
}

func TestHeaderUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Header
	}{
		{
			yml: `required: true`,
			want: Header{
				required: true,
			},
		},
		{
			yml: `deprecated: true`,
			want: Header{
				deprecated: true,
			},
		},
		{
			yml: `allowEmptyValue: true`,
			want: Header{
				allowEmptyValue: true,
			},
		},
		{
			yml: `style: foo`,
			want: Header{
				style: "foo",
			},
		},
		{
			yml: `explode: true`,
			want: Header{
				explode: true,
			},
		},
		{
			yml: `allowReserved: true`,
			want: Header{
				allowReserved: true,
			},
		},
		{
			yml: `x-foo: bar`,
			want: Header{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
		{
			yml: `example: foo`,
			want: Header{
				example: "foo",
			},
		},
		{
			yml: `examples:
  foo:
    value: bar`,
			want: Header{
				examples: map[string]*Example{
					"foo": {
						value: "bar",
					},
				},
			},
		},
		{
			yml: `content:
  application/json: {}`,
			want: Header{
				content: map[string]*MediaType{
					"application/json": {},
				},
			},
		},
		{
			yml: `$ref: '#/components/headers/foo'`,
			want: Header{
				reference: "#/components/headers/foo",
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Header
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestHeaderUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `schema: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `examples: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `content: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Header{})
			assertSameError(t, got, tt.want)
		})
	}
}
