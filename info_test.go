package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestInfoValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Info{}, openapi.ErrRequired{Target: "info.title"}},
		{"withTitle", openapi.Info{Title: "foo"}, openapi.ErrRequired{Target: "info.version"}},
		{"withTitleAndVersion", openapi.Info{Title: "foo", Version: "1.0"}, nil},
		{"withInvalidToS", openapi.Info{Title: "foo", TermsOfService: "foobar", Version: "1.0"}, openapi.ErrFormatInvalid{Target: "info.termsOfService", Format: "URL"}},
		{"withInvalidLicense", openapi.Info{Title: "foo", Version: "1.0", License: &openapi.License{URL: "foobar"}}, openapi.ErrRequired{Target: "license.name"}},
	}
	testValidater(t, candidates)
}
