package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestSecurityRequirementExampleUnmarshalYAML(t *testing.T) {
	t.Run("non-oauth2", func(t *testing.T) {
		yml := `api_key: []`
		var got SecurityRequirement
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := SecurityRequirement{
			securityRequirement: map[string][]string{
				"api_key": {},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("oauth2", func(t *testing.T) {
		yml := `petstore_auth:
- write:pets
- read:pets`
		var got SecurityRequirement
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := SecurityRequirement{
			securityRequirement: map[string][]string{
				"petstore_auth": {"write:pets", "read:pets"},
			},
		}
		assertEqual(t, got, want)
	})
}

func TestSecurityRequirementUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want SecurityRequirement
	}{}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got SecurityRequirement
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestSecurityRequirementUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `foo: bar`,
			want: errors.New("String node doesn't ArrayNode"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &SecurityRequirement{})
			assertSameError(t, got, tt.want)
		})
	}
}
