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
		{"empty", openapi.SecurityScheme{}, openapi.ErrRequired{Target: "securityScheme.type"}},
		{"apiKey/noName", openapi.SecurityScheme{Type: "apiKey", In: "query"}, openapi.ErrRequired{Target: "securityScheme.name"}},
		{"apiKey/noIn", openapi.SecurityScheme{Type: "apiKey", Name: "foo"}, openapi.ErrRequired{Target: "securityScheme.in"}},
		{"apiKey/invalidIn", openapi.SecurityScheme{Type: "apiKey", Name: "foo", In: "bar"}, openapi.ErrMustOneOf{Object: "securityScheme.in", ValidValues: []string{"query", "header", "cookie"}}},
		{"apiKey/valid", openapi.SecurityScheme{Type: "apiKey", Name: "foo", In: "query"}, nil},
		{"http/noScheme", openapi.SecurityScheme{Type: "http"}, openapi.ErrRequired{Target: "securityScheme.scheme"}},
		{"http/valid", openapi.SecurityScheme{Type: "http", Scheme: "Basic"}, nil},
		{"oauth2/noFlows", openapi.SecurityScheme{Type: "oauth2"}, openapi.ErrRequired{Target: "securityScheme.flows"}},
		{"oauth2/emptyFlows", openapi.SecurityScheme{Type: "oauth2", Flows: &openapi.OAuthFlows{}}, nil},
		{"oauth2/valid", openapi.SecurityScheme{Type: "oauth2", Flows: &openapi.OAuthFlows{Implicit: &openapi.OAuthFlow{AuthorizationURL: "http://example.com", Scopes: map[string]string{"foo": "bar"}}}}, nil},
		{"openIdConnect/noOIDCURL", openapi.SecurityScheme{Type: "openIdConnect"}, openapi.ErrRequired{Target: "securityScheme.openIdConnectUrl"}},
		{"openIdConnect/valid", openapi.SecurityScheme{Type: "openIdConnect", OpenIDConnectURL: "http://example.com"}, nil},
		{"invalidType", openapi.SecurityScheme{Type: "foo"}, openapi.ErrMustOneOf{Object: "securityScheme.type", ValidValues: openapi.SecuritySchemeTypeList}},
	}
	testValidater(t, candidates)
}
