package openapi

import (
	"errors"
	"reflect"
	"testing"
)

const (
	exampleCom  = "https://example.com"
	exampleMail = "foo@example.com"
)

type candidate struct {
	label  string
	in     validater
	hasErr bool
}

func testValidater(t *testing.T, candidates []candidate) {
	t.Helper()
	for _, c := range candidates {
		if err := c.in.Validate(); (err != nil) != c.hasErr {
			if c.hasErr {
				t.Error("error should be occurred, but not")
				continue
			}
			t.Errorf("error is occurred: %s", err)
		}
	}
}

type mockValidater struct {
	err error
}

func (v mockValidater) Validate() error {
	return v.err
}

func TestValidateAll(t *testing.T) {
	valid := mockValidater{}
	invalid := mockValidater{errors.New("err")}

	candidates := []struct {
		label  string
		in     []validater
		hasErr bool
	}{
		{"nil", nil, false},
		{"empty", []validater{}, false},
		{"all valid", []validater{valid, valid, valid}, false},
		{"have invalid", []validater{valid, invalid, valid}, true},
		{"have nil", []validater{valid, nil, valid, valid}, false},
	}
	for _, c := range candidates {
		if err := validateAll(c.in); (err != nil) != c.hasErr {
			t.Log(c.label)
			t.Errorf("error: %s", err)
		}
	}
}

func TestHasDuplicatedParameter(t *testing.T) {
	t.Run("no duplicated param", testHasDuplicatedParameterFalse)
	t.Run("there's duplicated param", testHasDuplicatedParameterTrue)
}

func testHasDuplicatedParameterFalse(t *testing.T) {
	params := []*Parameter{
		&Parameter{Name: "foo", In: "header"},
		&Parameter{Name: "foo", In: "path", Required: true},
		&Parameter{Name: "bar", In: "path", Required: true},
	}
	if hasDuplicatedParameter(params) {
		t.Error("should return false")
	}
}

func testHasDuplicatedParameterTrue(t *testing.T) {
	params := []*Parameter{
		&Parameter{Name: "foo", In: "header"},
		&Parameter{Name: "foo", In: "header"},
	}
	if !hasDuplicatedParameter(params) {
		t.Error("should return true")
	}
}

func TestMustURL(t *testing.T) {
	candidates := []struct {
		label  string
		in     string
		hasErr bool
	}{
		{"empty", "", true},
		{"valid HTTP url", "http://example.com", false},
		{"allowed relative path", "foo/bar/baz", true},
		{"absolute path", "/foo/bar/baz", false},
		{"plain string", "foobarbaz", true},
	}
	for _, c := range candidates {
		if err := mustURL(c.label, c.in); (err != nil) != c.hasErr {
			t.Logf("error occured at %s", c.label)
			if c.hasErr {
				t.Error("error should occured, but not")
				return
			}
			t.Error("error should not occurred, but occurred")
			return
		}
	}
}

func TestDocumentValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", Document{}, true},
		{"withInvalidVersion", Document{Version: "1.0"}, true},
		{"withVersion", Document{Version: "3.0.0"}, true},
		{"valid", Document{Version: "3.0.0", Info: &Info{Title: "foo", TermsOfService: exampleCom, Version: "1.0"}, Paths: Paths{}}, false},
	}
	testValidater(t, candidates)
}

func TestValidateOASVersion(t *testing.T) {
	candidates := []struct {
		label  string
		in     string
		hasErr bool
	}{
		{"empty", "", true},
		{"invalidVersion", "foobar", true},
		{"swagger", "2.0", true},
		{"valid", "3.0.0", false},
	}
	for _, c := range candidates {
		if err := validateOASVersion(c.in); (err != nil) != c.hasErr {
			t.Log(c.label)
			if c.hasErr {
				t.Error("error should be occurred, but not")
				continue
			}
			t.Errorf("error should not be occured: %s", err)
		}
	}
}

func TestInfoValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", Info{}, true},
	}
	testValidater(t, candidates)
}

func TestContactValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", Contact{}, true},
		{"withURL", Contact{URL: exampleCom}, false},
		{"invalidURL", Contact{URL: "foobar"}, true},
		{"withEmail", Contact{Email: exampleMail}, true},
		{"valid", Contact{URL: exampleCom, Email: exampleMail}, false},
		{"invalidEmail", Contact{URL: exampleCom, Email: "foobar"}, true},
	}

	testValidater(t, candidates)
}

func TestLicenseValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", License{}, true},
		{"withName", License{Name: "foobar"}, true},
		{"withURL", License{URL: exampleCom}, true},
		{"invalidURL", License{Name: "foobar", URL: "foobar"}, true},
		{"valid", License{Name: "foobar", URL: exampleCom}, false},
	}
	testValidater(t, candidates)
}

func TestServerValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", Server{}, true},
		{"invalidURL", Server{URL: "foobar%"}, true},
		{"withURL", Server{URL: exampleCom}, false},
	}
	testValidater(t, candidates)
}

func TestServerVariableValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", ServerVariable{}, true},
		{"withDefault", ServerVariable{Default: "default"}, false},
	}
	testValidater(t, candidates)
}

func TestComponents(t *testing.T) {
	candidates := []candidate{
		{"empty", Components{}, false},
	}
	testValidater(t, candidates)
}

func TestComponentsValidateKeys(t *testing.T) {
	candidates := []struct {
		label  string
		in     Components
		hasErr bool
	}{
		{"empty", Components{}, false},
	}
	for _, c := range candidates {
		if err := c.in.validateKeys(); (err != nil) != c.hasErr {
			t.Log(c.label)
			if c.hasErr {
				t.Error("error should be occurred, but not")
				continue
			}
			t.Errorf("error should not be occurred: %s", err)
		}
	}
}

func TestReduceComponentKeys(t *testing.T) {
	candidates := []struct {
		label    string
		in       Components
		expected []string
	}{
		{"empty", Components{}, []string{}},
	}
	for _, c := range candidates {
		keys := reduceComponentKeys(c.in)
		if !reflect.DeepEqual(keys, c.expected) {
			t.Log(c.label)
			t.Errorf("%+v != %+v", keys, c.expected)
		}
	}
}

func TestReduceComponentObjects(t *testing.T) {
	candidates := []struct {
		label    string
		in       Components
		expected []validater
	}{
		{"empty", Components{}, []validater{}},
	}
	for _, c := range candidates {
		objects := reduceComponentObjects(c.in)
		if !reflect.DeepEqual(objects, c.expected) {
			t.Log(c.label)
			t.Errorf("%+v != %+v", objects, c.expected)
		}
	}
}

func TestPathsValidate(t *testing.T) {
	t.Run("duplicate pathItem", testPathItemDuplicate)
}

func getPaths(id1, id2 string) Paths {
	return Paths{
		"/foo/bar": &PathItem{
			Get:  &Operation{OperationID: id1, Responses: Responses{"200": &Response{Description: "foo"}}},
			Post: &Operation{OperationID: id2, Responses: Responses{"200": &Response{Description: "foo"}}},
		},
	}
}

func testPathItemDuplicate(t *testing.T) {
	candidates := []candidate{
		{"invalid", getPaths("foobar", "foobar"), true},
		{"valid", getPaths("foo", "bar"), false},
	}
	testValidater(t, candidates)
}

func TestExternalDocumentationValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", ExternalDocumentation{}, true},
		{"invalidURL", ExternalDocumentation{URL: "foobar"}, true},
		{"valid", ExternalDocumentation{URL: exampleCom}, false},
	}
	testValidater(t, candidates)
}

func TestTagValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", Tag{}, true},
		{"withEmptyExternalDocs", Tag{ExternalDocs: &ExternalDocumentation{}}, true},
		{"withValidExternalDocs", Tag{ExternalDocs: &ExternalDocumentation{URL: exampleCom}}, true},

		{"withName", Tag{Name: "foo"}, false},
	}
	testValidater(t, candidates)
}

func TestSchemaValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", Schema{}, false},
	}
	testValidater(t, candidates)
}

func TestDiscriminatorValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", Discriminator{}, true},
		{"withPropertyName", Discriminator{PropertyName: "foobar"}, false},
	}
	testValidater(t, candidates)
}

func TestXMLValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", XML{}, true},
		{"invalidURLNamespace", XML{Namespace: "foobar"}, true},
		{"withNamespace", XML{Namespace: exampleCom}, false},
	}
	testValidater(t, candidates)
}

func TestOAuthFlowValidate(t *testing.T) {
	mockScopes := map[string]string{"foo": "bar"}

	empty := OAuthFlow{}
	aURL := OAuthFlow{AuthorizationURL: exampleCom}
	tURL := OAuthFlow{TokenURL: exampleCom}
	rURL := OAuthFlow{RefreshURL: exampleCom}
	scopes := OAuthFlow{Scopes: mockScopes}
	atURL := OAuthFlow{AuthorizationURL: exampleCom, TokenURL: exampleCom}
	arURL := OAuthFlow{AuthorizationURL: exampleCom, RefreshURL: exampleCom}
	aURLscopes := OAuthFlow{AuthorizationURL: exampleCom, Scopes: mockScopes}
	trURL := OAuthFlow{TokenURL: exampleCom, RefreshURL: exampleCom}
	tURLscopes := OAuthFlow{TokenURL: exampleCom, Scopes: mockScopes}
	rURLscopes := OAuthFlow{RefreshURL: exampleCom, Scopes: mockScopes}
	atrURL := OAuthFlow{AuthorizationURL: exampleCom, TokenURL: exampleCom, RefreshURL: exampleCom}
	atURLscopes := OAuthFlow{AuthorizationURL: exampleCom, TokenURL: exampleCom, Scopes: mockScopes}
	arURLscopes := OAuthFlow{AuthorizationURL: exampleCom, RefreshURL: exampleCom, Scopes: mockScopes}
	trURLscopes := OAuthFlow{TokenURL: exampleCom, RefreshURL: exampleCom, Scopes: mockScopes}
	atrURLscopes := OAuthFlow{AuthorizationURL: exampleCom, TokenURL: exampleCom, RefreshURL: exampleCom, Scopes: mockScopes}
	invalidURL := OAuthFlow{AuthorizationURL: "foobar", TokenURL: "foobar", RefreshURL: "foobar", Scopes: mockScopes}
	zeroMap := OAuthFlow{AuthorizationURL: exampleCom, TokenURL: exampleCom, RefreshURL: exampleCom, Scopes: map[string]string{}}

	candidates := []struct {
		label   string
		in      OAuthFlow
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

func testOAuthFlowValidate(t *testing.T, label string, oauthFlow OAuthFlow, haveErr [4]bool) {
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
