package openapi_test

import (
	"fmt"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

type candidateBase struct {
	label   string
	in      openapi.OAuthFlow
	haveErr [4]bool
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

	candidatebases := []candidateBase{
		{"empty", empty, [4]bool{true, true, true, true}},
		{"aURL", aURL, [4]bool{true, true, true, true}},
		{"tURL", tURL, [4]bool{true, true, true, true}},
		{"rURL", rURL, [4]bool{true, true, true, true}},
		{"scopes", scopes, [4]bool{true, true, true, true}},
		{"aURL/tURL", atURL, [4]bool{true, true, true, true}},
		{"aURL/rURL", arURL, [4]bool{true, true, true, true}},
		{"aURL/scopes", aURLscopes, [4]bool{false, true, true, true}},
		{"tURL/rURL", trURL, [4]bool{true, true, true, true}},
		{"tURL/scopes", tURLscopes, [4]bool{true, false, false, true}},
		{"rURL/scopes", rURLscopes, [4]bool{true, true, true, true}},
		{"aURL/tURL/rURL", atrURL, [4]bool{true, true, true, true}},
		{"aURL/tURL/scopes", atURLscopes, [4]bool{false, false, false, false}},
		{"aURL/rURL/scopes", arURLscopes, [4]bool{false, true, true, true}},
		{"tURL/rURL/scopes", trURLscopes, [4]bool{true, false, false, true}},
		{"aURL/tURL/rURL/scopes", atrURLscopes, [4]bool{false, false, false, false}},

		{"invalidaURL", invalidaURL, [4]bool{true, false, false, true}},
		{"invalidtURL", invalidtURL, [4]bool{false, true, true, true}},
		{"invalidrURL", invalidrURL, [4]bool{true, true, true, true}},
		{"zero length map", zeroMap, [4]bool{true, true, true, true}},
	}
	candidates := generateCandidates(candidatebases)
	testValidater(t, candidates)
}

var flowTypes = []string{"implicit", "password", "clientCredentials", "authorizationCode"}

func generateCandidates(base []candidateBase) []candidate {
	candidates := []candidate{}
	for _, c := range base {
		c.in.SetFlowType("")
		candidates = append(candidates, candidate{fmt.Sprintf("%s-empty", c.label), c.in, true})
		c.in.SetFlowType("foobar")
		candidates = append(candidates, candidate{fmt.Sprintf("%s-wrongtype", c.label), c.in, true})
		for i, flowType := range flowTypes {
			c.in.SetFlowType(flowType)
			candidates = append(candidates, candidate{fmt.Sprintf("%s-%s", c.label, flowType), c.in, c.haveErr[i]})
		}
	}
	return candidates
}
