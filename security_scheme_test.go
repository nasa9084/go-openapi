package openapi

import (
	"testing"
)

func TestSecurityScheme(t *testing.T) {
	t.Run("validate", testSecuritySchemeValidate)
}

func testSecuritySchemeValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", SecurityScheme{}, true},
		{"apiKey/noName", SecurityScheme{Type: "apiKey", In: "query"}, true},
		{"apiKey/noIn", SecurityScheme{Type: "apiKey", Name: "foo"}, true},
		{"apiKey/invalidIn", SecurityScheme{Type: "apiKey", Name: "foo", In: "bar"}, true},
		{"apiKey/valid", SecurityScheme{Type: "apiKey", Name: "foo", In: "query"}, false},
		{"http/noScheme", SecurityScheme{Type: "http"}, true},
		{"http/valid", SecurityScheme{Type: "http", Scheme: "Basic"}, false},
		{"oauth2/noFlows", SecurityScheme{Type: "oauth2"}, true},
		{"oauth2/emptyFlows", SecurityScheme{Type: "oauth2", Flows: &OAuthFlows{}}, false},
		{"oauth2/valid", SecurityScheme{Type: "oauth2", Flows: &OAuthFlows{Implicit: &OAuthFlow{AuthorizationURL: "http://example.com", Scopes: map[string]string{"foo": "bar"}}}}, false},
		{"openIdConnect/noOIDCURL", SecurityScheme{Type: "openIdConnect"}, true},
		{"openIdConnect/valid", SecurityScheme{Type: "openIdConnect", OpenIDConnectURL: "http://example.com"}, false},
	}
	testValidater(t, candidates)
}
