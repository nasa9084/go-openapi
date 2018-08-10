package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

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
	invalidURL := openapi.OAuthFlow{AuthorizationURL: "foobar", TokenURL: "foobar", RefreshURL: "foobar", Scopes: mockScopes}
	zeroMap := openapi.OAuthFlow{AuthorizationURL: exampleCom, TokenURL: exampleCom, RefreshURL: exampleCom, Scopes: map[string]string{}}

	candidates := []struct {
		label   string
		in      openapi.OAuthFlow
		haveErr [4]bool
	}{
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

		{"invalidURL", invalidURL, [4]bool{true, true, true, true}},
		{"zero length map", zeroMap, [4]bool{true, true, true, true}},
	}
	for _, c := range candidates {
		testOAuthFlowValidate(t, c.label, c.in, c.haveErr)
	}
}

var flowTypes = []string{"implicit", "password", "clientCredentials", "authorizationCode"}

func testOAuthFlowValidate(t *testing.T, label string, oauthFlow openapi.OAuthFlow, haveErr [4]bool) {
	if err := oauthFlow.Validate(""); err == nil {
		t.Logf("%s-empty", label)
		t.Error("error should be occurred, but not")
	}
	if err := oauthFlow.Validate("foobar"); err == nil {
		t.Logf("%s-wrongtype", label)
		t.Error("error should be occurred, but not")
	}
	for i, flowType := range flowTypes {
		if err := oauthFlow.Validate(flowType); (err != nil) != haveErr[i] {
			t.Logf("%s-%s", label, flowType)
			if haveErr[i] {
				t.Error("error should be occurred, but not")
				continue
			}
			t.Error("error should not be occurred, but occurred")
			t.Log(err)
		}
	}
}
