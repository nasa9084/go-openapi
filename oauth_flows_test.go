package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestOAuthFlowsValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.OAuthFlows{}, false},
		{"invalidImplicit", openapi.OAuthFlows{Implicit: &openapi.OAuthFlow{}}, true},
		{"invalidPassword", openapi.OAuthFlows{Password: &openapi.OAuthFlow{}}, true},
		{"invalidClientCredentials", openapi.OAuthFlows{ClientCredentials: &openapi.OAuthFlow{}}, true},
		{"invalidAuthorizationCode", openapi.OAuthFlows{AuthorizationCode: &openapi.OAuthFlow{}}, true},
	}
	testValidater(t, candidates)
}
