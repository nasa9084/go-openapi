package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestServerExampleUnmarshalYAML(t *testing.T) {
	t.Run("single server", func(t *testing.T) {
		yml := `url: https://development.gigantic-server.com/v1
description: Development server`

		var got Server
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}

		want := Server{
			url:         "https://development.gigantic-server.com/v1",
			description: "Development server",
		}
		assertEqual(t, got, want)
	})
	t.Run("servers", func(t *testing.T) {
		yml := `servers:
- url: https://development.gigantic-server.com/v1
  description: Development server
- url: https://staging.gigantic-server.com/v1
  description: Staging server
- url: https://api.gigantic-server.com/v1
  description: Production server`

		var target struct {
			Servers []*Server
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		got := target.Servers
		want := []*Server{
			{
				url:         "https://development.gigantic-server.com/v1",
				description: "Development server",
			},
			{
				url:         "https://staging.gigantic-server.com/v1",
				description: "Staging server",
			},
			{
				url:         "https://api.gigantic-server.com/v1",
				description: "Production server",
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("with variables", func(t *testing.T) {
		yml := `servers:
- url: https://{username}.gigantic-server.com:{port}/{basePath}
  description: The production API server
  variables:
    username:
      # note! no enum here means it is an open value
      default: demo
      description: this value is assigned by the service provider, in this example "gigantic-server.com"
    port:
      enum:
        - '8443'
        - '443'
      default: '8443'
    basePath:
      # open meaning there is the opportunity to use special base paths as assigned by the provider, default is "v2"
      default: v2`

		var target struct {
			Servers []*Server
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}

		got := target.Servers[0]
		want := &Server{
			url:         "https://{username}.gigantic-server.com:{port}/{basePath}",
			description: "The production API server",
			variables: map[string]*ServerVariable{
				"username": {
					default_:    "demo",
					description: `this value is assigned by the service provider, in this example "gigantic-server.com"`,
				},
				"port": {
					enum:     []string{"8443", "443"},
					default_: "8443",
				},
				"basePath": {
					default_: "v2",
				},
			},
		}
		assertEqual(t, got, want)
	})
}

func TestServerUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Server
	}{
		{
			yml: `url: https://example.com`,
			want: Server{
				url: "https://example.com",
			},
		},
		{
			yml: `url: https://example.com
x-foo: bar`,
			want: Server{
				url: "https://example.com",
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Server
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestServerUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `description: foobar`,
			want: ErrRequired("url"),
		},
		{
			yml: `description: foobar
url: https%20://example.com`,
			want: errors.New(`parse "https%20://example.com": first path segment in URL cannot contain colon`),
		},
		{
			yml: `description: foobar
url: example.com
variables: hoge`,
			// variables expexts an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `url: example.com
foo: bar`,
			want: errors.New(`unknown key: foo`),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Server{})
			assertSameError(t, got, tt.want)
		})
	}
}
