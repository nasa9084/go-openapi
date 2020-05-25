package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestOAuthFlowExampleUnmarshalYAML(t *testing.T) {
	yml := `type: oauth2
flows:
  implicit:
    authorizationUrl: https://example.com/api/oauth/dialog
    scopes:
      write:pets: modify pets in your account
      read:pets: read your pets
  authorizationCode:
    authorizationUrl: https://example.com/api/oauth/dialog
    tokenUrl: https://example.com/api/oauth/token
    scopes:
      write:pets: modify pets in your account
      read:pets: read your pets `

	var got SecurityScheme
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := SecurityScheme{
		type_: "oauth2",
		flows: &OAuthFlows{
			implicit: &OAuthFlow{
				authorizationURL: "https://example.com/api/oauth/dialog",
				scopes: map[string]string{
					"write:pets": "modify pets in your account",
					"read:pets":  "read your pets",
				},
			},
			authorizationCode: &OAuthFlow{
				authorizationURL: "https://example.com/api/oauth/dialog",
				tokenURL:         "https://example.com/api/oauth/token",
				scopes: map[string]string{
					"write:pets": "modify pets in your account",
					"read:pets":  "read your pets",
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func TestOAuthFlowUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want OAuthFlow
	}{
		{
			yml: `refreshUrl: https://example.com`,
			want: OAuthFlow{
				refreshURL: "https://example.com",
			},
		},
		{
			yml: `x-foo: bar`,
			want: OAuthFlow{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got OAuthFlow
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestOAuthFlowUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml: `authorizationUrl: foobar`,
			// authorizationUrl expects URI
			want: errors.New(`parse "foobar": invalid URI for request`),
		},
		{
			yml: `tokenUrl: foobar`,
			// tokenUrl expects URI
			want: errors.New(`parse "foobar": invalid URI for request`),
		},
		{
			yml: `refreshUrl: foobar`,
			// refreshUrl expects URI
			want: errors.New(`parse "foobar": invalid URI for request`),
		},
		{
			yml: `scopes: foobar`,
			// scopes expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &OAuthFlow{})
			assertSameError(t, got, tt.want)
		})
	}
}
