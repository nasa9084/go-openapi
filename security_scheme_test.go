package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestSecurityScheme(t *testing.T) {
	t.Run("validate", testSecuritySchemeValidate)
}

func testSecuritySchemeValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.SecurityScheme{}, true},
		{"apiKey/noName", openapi.SecurityScheme{Type: "apiKey", In: "query"}, true},
		{"apiKey/noIn", openapi.SecurityScheme{Type: "apiKey", Name: "foo"}, true},
		{"apiKey/invalidIn", openapi.SecurityScheme{Type: "apiKey", Name: "foo", In: "bar"}, true},
		{"apiKey/valid", openapi.SecurityScheme{Type: "apiKey", Name: "foo", In: "query"}, false},
		{"http/noScheme", openapi.SecurityScheme{Type: "http"}, true},
		{"http/valid", openapi.SecurityScheme{Type: "http", Scheme: "Basic"}, false},
		{"oauth2/noFlows", openapi.SecurityScheme{Type: "oauth2"}, true},
		{"oauth2/emptyFlows", openapi.SecurityScheme{Type: "oauth2", Flows: &openapi.OAuthFlows{}}, false},
		{"oauth2/valid", openapi.SecurityScheme{Type: "oauth2", Flows: &openapi.OAuthFlows{Implicit: &openapi.OAuthFlow{AuthorizationURL: "http://example.com", Scopes: map[string]string{"foo": "bar"}}}}, false},
		{"openIdConnect/noOIDCURL", openapi.SecurityScheme{Type: "openIdConnect"}, true},
		{"openIdConnect/valid", openapi.SecurityScheme{Type: "openIdConnect", OpenIDConnectURL: "http://example.com"}, false},
	}
	testValidater(t, candidates)
}
