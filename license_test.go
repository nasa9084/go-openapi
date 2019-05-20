package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestLicense_Validate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.License{}, openapi.ErrRequired{Target: "license.name"}},
		{"withName", openapi.License{Name: "foobar"}, nil},
		{"withURL", openapi.License{URL: exampleCom}, openapi.ErrRequired{Target: "license.name"}},
		{"invalidURL", openapi.License{Name: "foobar", URL: "foobar"}, openapi.ErrFormatInvalid{Target: "license.url", Format: "URL"}},
		{"full", openapi.License{Name: "foobar", URL: exampleCom}, nil},
	}
	testValidater(t, candidates)
}
