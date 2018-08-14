package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestInfoValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Info{}, true},
		{"withTitle", openapi.Info{Title: "foo"}, true},
		{"withTitleAndVersion", openapi.Info{Title: "foo", Version: "1.0"}, false},
		{"withInvalidToS", openapi.Info{Title: "foo", TermsOfService: "foobar", Version: "1.0"}, true},
		{"withInvalidLicense", openapi.Info{Title: "foo", Version: "1.0", License: &openapi.License{URL: "foobar"}}, true},
	}
	testValidater(t, candidates)
}
