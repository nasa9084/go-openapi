package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestOAuthFlowsUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want OAuthFlows
	}{
		{
			yml: `password: {}`,
			want: OAuthFlows{
				password: &OAuthFlow{},
			},
		},
		{
			yml: `clientCredentials: {}`,
			want: OAuthFlows{
				clientCredentials: &OAuthFlow{},
			},
		},
		{
			yml: `x-foo: bar`,
			want: OAuthFlows{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got OAuthFlows
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestOAuthFlowsUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml: `implicit: foobar`,
			// implicit expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `password: foobar`,
			// password expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `clientCredentials: foobar`,
			// clientCredentials expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `authorizationCode: foobar`,
			// authorizationCode expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &OAuthFlows{})
			assertSameError(t, got, tt.want)
		})
	}
}
