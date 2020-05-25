package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestSecuritySchemeExampleUnmarshalYAML(t *testing.T) {
	t.Run("basic auth", func(t *testing.T) {
		yml := `type: http
scheme: basic`
		var got SecurityScheme
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := SecurityScheme{
			type_:  "http",
			scheme: "basic",
		}
		assertEqual(t, got, want)
	})
	t.Run("api key", func(t *testing.T) {
		yml := `type: apiKey
name: api_key
in: header`
		var got SecurityScheme
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := SecurityScheme{
			type_: "apiKey",
			name:  "api_key",
			in:    "header",
		}
		assertEqual(t, got, want)
	})
	t.Run("JWT bearer", func(t *testing.T) {
		yml := `type: http
scheme: bearer
bearerFormat: JWT`
		var got SecurityScheme
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := SecurityScheme{
			type_:        "http",
			scheme:       "bearer",
			bearerFormat: "JWT",
		}
		assertEqual(t, got, want)
	})
	t.Run("implicit oauth2", func(t *testing.T) {
		yml := `type: oauth2
flows:
  implicit:
    authorizationUrl: https://example.com/api/oauth/dialog
    scopes:
      write:pets: modify pets in your account
      read:pets: read your pets`
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
			},
		}
		assertEqual(t, got, want)
	})
}

func TestSecuritySchemeUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want SecurityScheme
	}{
		{
			yml: `type: apiKey`,
			want: SecurityScheme{
				type_: "apiKey",
			},
		},
		{
			yml: `type: http`,
			want: SecurityScheme{
				type_: "http",
			},
		},
		{
			yml: `type: oauth2`,
			want: SecurityScheme{
				type_: "oauth2",
			},
		},
		{
			yml: `type: openIdConnect`,
			want: SecurityScheme{
				type_: "openIdConnect",
			},
		},
		{
			yml: `description: foobar`,
			want: SecurityScheme{
				description: "foobar",
			},
		},
		{
			yml: `in: query`,
			want: SecurityScheme{
				in: "query",
			},
		},
		{
			yml: `in: header`,
			want: SecurityScheme{
				in: "header",
			},
		},
		{
			yml: `in: cookie`,
			want: SecurityScheme{
				in: "cookie",
			},
		},
		{
			yml: `openIdConnectUrl: https://example.com`,
			want: SecurityScheme{
				openIDConnectURL: "https://example.com",
			},
		},
		{
			yml: `$ref: '#/components/securitySchemes/Foo'`,
			want: SecurityScheme{
				reference: "#/components/securitySchemes/Foo",
			},
		},
		{
			yml: `x-foo: bar`,
			want: SecurityScheme{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got SecurityScheme
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestSecuritySchemeUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `type: foo`,
			want: errors.New(`"type" field must be one of ["apiKey", "http", "oauth2", "openIdConnect"]`),
		},
		{
			yml:  `in: "foo"`,
			want: errors.New(`"in" field must be one of ["query", "header", "cookie"]`),
		},
		{
			yml:  `flows: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `openIdConnectUrl: foo`,
			want: errors.New(`parse "foo": invalid URI for request`),
		},

		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &SecurityScheme{})
			assertSameError(t, got, tt.want)
		})
	}
}
