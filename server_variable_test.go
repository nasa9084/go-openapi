package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestServerVariableUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want ServerVariable
	}{
		{
			yml: `default: defaultValue`,
			want: ServerVariable{
				default_: "defaultValue",
			},
		},
		{
			yml: `default: defaultValue
x-foo: bar`,
			want: ServerVariable{
				default_: "defaultValue",
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got ServerVariable
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestServerVariableUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `description: foobar`,
			want: ErrRequired("default"),
		},
		{
			yml: `default: foo
enum: bar`,
			// enum expects an array
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml: `default: defaultValue
foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &ServerVariable{})
			assertSameError(t, got, tt.want)
		})
	}
}
