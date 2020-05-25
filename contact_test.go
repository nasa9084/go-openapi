package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestContactExampleUnmarshalYAML(t *testing.T) {
	yml := `name: API Support
url: http://www.example.com/support
email: support@example.com`

	var got Contact
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := Contact{
		name:  "API Support",
		url:   "http://www.example.com/support",
		email: "support@example.com",
	}
	assertEqual(t, got, want)
}

func TestContactUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Contact
	}{
		{
			yml: `email: foo@example.com`,
			want: Contact{
				email: "foo@example.com",
			},
		},
		{
			yml: `x-foo: bar`,
			want: Contact{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Contact
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestContactUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `email: invalidEmail`,
			want: errors.New(`"email" field must be an email address`),
		},
		{
			yml:  `url: foobar`,
			want: errors.New(`parse "foobar": invalid URI for request`),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Contact{})
			assertSameError(t, got, tt.want)
		})
	}
}
