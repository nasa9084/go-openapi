package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestLicenseExampleUnmarshalYAML(t *testing.T) {
	yml := `name: Apache 2.0
url: https://www.apache.org/licenses/LICENSE-2.0.html`

	var got License
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := License{
		name: "Apache 2.0",
		url:  "https://www.apache.org/licenses/LICENSE-2.0.html",
	}

	assertEqual(t, got, want)
}

func TestLicenseUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want License
	}{
		{
			yml: `name: licensename`,
			want: License{
				name: "licensename",
			},
		},
		{
			yml: `name: licensename
x-foo: bar`,
			want: License{
				name: "licensename",
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got License
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestLicenseUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `url: example.com`,
			want: ErrRequired("name"),
		},
		{
			yml: `name: foobar
url: hoge`,
			want: errors.New(`parse "hoge": invalid URI for request`),
		},
		{
			yml: `name: licensename
foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &License{})
			assertSameError(t, got, tt.want)
		})
	}
}
