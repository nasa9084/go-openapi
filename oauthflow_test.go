package openapi_test

import (
	"fmt"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

type candidateBase struct {
	label string
	in    openapi.OAuthFlow
	err   [4]error
}

func TestOAuthFlowValidate(t *testing.T) {
	mockScopes := map[string]string{"foo": "bar"}

	empty := openapi.OAuthFlow{}
	aURL := openapi.OAuthFlow{AuthorizationURL: exampleCom}
	tURL := openapi.OAuthFlow{TokenURL: exampleCom}
	rURL := openapi.OAuthFlow{RefreshURL: exampleCom}
	scopes := openapi.OAuthFlow{Scopes: mockScopes}
	atURL := openapi.OAuthFlow{AuthorizationURL: exampleCom, TokenURL: exampleCom}
	arURL := openapi.OAuthFlow{AuthorizationURL: exampleCom, RefreshURL: exampleCom}
	aURLscopes := openapi.OAuthFlow{AuthorizationURL: exampleCom, Scopes: mockScopes}
	trURL := openapi.OAuthFlow{TokenURL: exampleCom, RefreshURL: exampleCom}
	tURLscopes := openapi.OAuthFlow{TokenURL: exampleCom, Scopes: mockScopes}
	rURLscopes := openapi.OAuthFlow{RefreshURL: exampleCom, Scopes: mockScopes}
	atrURL := openapi.OAuthFlow{AuthorizationURL: exampleCom, TokenURL: exampleCom, RefreshURL: exampleCom}
	atURLscopes := openapi.OAuthFlow{AuthorizationURL: exampleCom, TokenURL: exampleCom, Scopes: mockScopes}
	arURLscopes := openapi.OAuthFlow{AuthorizationURL: exampleCom, RefreshURL: exampleCom, Scopes: mockScopes}
	trURLscopes := openapi.OAuthFlow{TokenURL: exampleCom, RefreshURL: exampleCom, Scopes: mockScopes}
	atrURLscopes := openapi.OAuthFlow{AuthorizationURL: exampleCom, TokenURL: exampleCom, RefreshURL: exampleCom, Scopes: mockScopes}
	invalidaURL := openapi.OAuthFlow{AuthorizationURL: "foobar", TokenURL: exampleCom, RefreshURL: exampleCom, Scopes: mockScopes}
	invalidtURL := openapi.OAuthFlow{AuthorizationURL: exampleCom, TokenURL: "foobar", RefreshURL: exampleCom, Scopes: mockScopes}
	invalidrURL := openapi.OAuthFlow{AuthorizationURL: exampleCom, TokenURL: exampleCom, RefreshURL: "foobar", Scopes: mockScopes}
	zeroMap := openapi.OAuthFlow{AuthorizationURL: exampleCom, TokenURL: exampleCom, RefreshURL: exampleCom, Scopes: map[string]string{}}

	aURLRequiredError := openapi.ErrRequired{Target: "oauthFlow.authorizationUrl"}
	tURLRequiredError := openapi.ErrRequired{Target: "oauthFlow.tokenUrl"}
	scopesRequiredError := openapi.ErrRequired{Target: "oauthFlow.scopes"}
	aURLFormatError := openapi.ErrFormatInvalid{Target: "oauthFlow.authorizationUrl", Format: "URL"}
	tURLFormatError := openapi.ErrFormatInvalid{Target: "oauthFlow.tokenUrl", Format: "URL"}
	rURLFormatError := openapi.ErrFormatInvalid{Target: "oauthFlow.refreshUrl", Format: "URL"}

	candidatebases := []candidateBase{
		{"empty", empty, [4]error{aURLRequiredError, tURLRequiredError, tURLRequiredError, aURLRequiredError}},
		{"aURL", aURL, [4]error{scopesRequiredError, tURLRequiredError, tURLRequiredError, tURLRequiredError}},
		{"tURL", tURL, [4]error{aURLRequiredError, scopesRequiredError, scopesRequiredError, aURLRequiredError}},
		{"rURL", rURL, [4]error{aURLRequiredError, tURLRequiredError, tURLRequiredError, aURLRequiredError}},
		{"scopes", scopes, [4]error{aURLRequiredError, tURLRequiredError, tURLRequiredError, aURLRequiredError}},
		{"aURL/tURL", atURL, [4]error{scopesRequiredError, scopesRequiredError, scopesRequiredError, scopesRequiredError}},
		{"aURL/rURL", arURL, [4]error{scopesRequiredError, tURLRequiredError, tURLRequiredError, tURLRequiredError}},
		{"aURL/scopes", aURLscopes, [4]error{nil, tURLRequiredError, tURLRequiredError, tURLRequiredError}},
		{"tURL/rURL", trURL, [4]error{aURLRequiredError, scopesRequiredError, scopesRequiredError, aURLRequiredError}},
		{"tURL/scopes", tURLscopes, [4]error{aURLRequiredError, nil, nil, aURLRequiredError}},
		{"rURL/scopes", rURLscopes, [4]error{aURLRequiredError, tURLRequiredError, tURLRequiredError, aURLRequiredError}},
		{"aURL/tURL/rURL", atrURL, [4]error{scopesRequiredError, scopesRequiredError, scopesRequiredError, scopesRequiredError}},
		{"aURL/tURL/scopes", atURLscopes, [4]error{nil, nil, nil, nil}},
		{"aURL/rURL/scopes", arURLscopes, [4]error{nil, tURLRequiredError, tURLRequiredError, tURLRequiredError}},
		{"tURL/rURL/scopes", trURLscopes, [4]error{aURLRequiredError, nil, nil, aURLRequiredError}},
		{"aURL/tURL/rURL/scopes", atrURLscopes, [4]error{nil, nil, nil, nil}},

		{"invalidaURL", invalidaURL, [4]error{aURLFormatError, nil, nil, aURLFormatError}},
		{"invalidtURL", invalidtURL, [4]error{nil, tURLFormatError, tURLFormatError, tURLFormatError}},
		{"invalidrURL", invalidrURL, [4]error{rURLFormatError, rURLFormatError, rURLFormatError, rURLFormatError}},
		{"zero length map", zeroMap, [4]error{scopesRequiredError, scopesRequiredError, scopesRequiredError, scopesRequiredError}},
	}
	candidates := generateCandidates(candidatebases)
	testValidater(t, candidates)
}

var flowTypes = []string{"implicit", "password", "clientCredentials", "authorizationCode"}

func generateCandidates(base []candidateBase) []candidate {
	candidates := []candidate{}
	for _, c := range base {
		c.in.SetFlowType("")
		candidates = append(candidates, candidate{fmt.Sprintf("%s-empty", c.label), c.in, openapi.InvalidFlowTypeError})
		c.in.SetFlowType("foobar")
		candidates = append(candidates, candidate{fmt.Sprintf("%s-wrongtype", c.label), c.in, openapi.InvalidFlowTypeError})
		for i, flowType := range flowTypes {
			c.in.SetFlowType(flowType)
			candidates = append(candidates, candidate{fmt.Sprintf("%s-%s", c.label, flowType), c.in, c.err[i]})
		}
	}
	return candidates
}
