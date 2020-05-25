package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestExternalDocumentationExampleUnmarshalYAML(t *testing.T) {
	yml := `description: Find more info here
url: https://example.com`

	var got ExternalDocumentation
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := ExternalDocumentation{
		description: "Find more info here",
		url:         "https://example.com",
	}
	assertEqual(t, got, want)
}

func TestExternalDocumentationUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want ExternalDocumentation
	}{
		{
			yml: `url: https://example.com`,
			want: ExternalDocumentation{
				url: "https://example.com",
			},
		},
		{
			yml: `url: https://example.com
x-foo: bar`,
			want: ExternalDocumentation{
				url: "https://example.com",
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got ExternalDocumentation
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestExternalDocumentationUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `description: foo`,
			want: ErrRequired("url"),
		},
		{
			yml:  `url: foobar`,
			want: errors.New(`parse "foobar": invalid URI for request`),
		},
		{
			yml: `url: https://example.com
foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &ExternalDocumentation{})
			assertSameError(t, got, tt.want)
		})
	}
}
