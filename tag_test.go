package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestTagExampleUnmarshal(t *testing.T) {
	yml := `name: pet
description: Pets operations`

	var got Tag
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := Tag{
		name:        "pet",
		description: "Pets operations",
	}
	assertEqual(t, got, want)
}

func TestTagUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Tag
	}{
		{
			yml: `name: theName
x-foo: bar`,
			want: Tag{
				name: "theName",
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
		{
			yml: `name: foo
externalDocs:
  url: https://example.com`,
			want: Tag{
				name: "foo",
				externalDocs: &ExternalDocumentation{
					url: "https://example.com",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Tag
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestTagUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml: `name: foo
externalDocs: bar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `description: foobar`,
			want: ErrRequired("name"),
		},
		{
			yml: `name: tagName
foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Tag{})
			assertSameError(t, got, tt.want)
		})
	}
}
