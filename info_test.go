package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestInfoExampleUnmarshalYAML(t *testing.T) {
	yml := `title: Sample Pet Store App
description: This is a sample server for a pet store.
termsOfService: http://example.com/terms/
contact:
  name: API Support
  url: http://www.example.com/support
  email: support@example.com
license:
  name: Apache 2.0
  url: https://www.apache.org/licenses/LICENSE-2.0.html
version: 1.0.1`

	var got Info
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := Info{
		title:          "Sample Pet Store App",
		description:    "This is a sample server for a pet store.",
		termsOfService: "http://example.com/terms/",
		contact: &Contact{
			name:  "API Support",
			url:   "http://www.example.com/support",
			email: "support@example.com",
		},
		license: &License{
			name: "Apache 2.0",
			url:  "https://www.apache.org/licenses/LICENSE-2.0.html",
		},
		version: "1.0.1",
	}
	assertEqual(t, got, want)
}

func TestInfoUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Info
	}{
		{
			yml: `title: this is title
version: 1.0.0`,
			want: Info{
				title:   "this is title",
				version: "1.0.0",
			},
		},
		{
			yml: `title: this is title
version: 1.0.0
x-foo: bar`,
			want: Info{
				title:   "this is title",
				version: "1.0.0",
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Info
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestInfoUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `version: 1.0.0`,
			want: ErrRequired("title"),
		},
		{
			yml:  `title: this is title`,
			want: ErrRequired("version"),
		},
		{
			yml: `ersion: 1.0.0
title: foobar
termsOfService: hoge`,
			// termsOfService expects URI
			want: errors.New(`parse "hoge": invalid URI for request`),
		},
		{
			yml: `ersion: 1.0.0
title: foobar
contact: hoge`,
			// contact expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `ersion: 1.0.0
title: foobar
license: hoge`,
			// license expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `title: this is title
version: 1.0.0
foo: bar`,
			want: errors.New(`unknown key: foo`),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Info{})
			assertSameError(t, got, tt.want)
		})
	}
}
