package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestLicenseValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.License{}, true},
		{"withName", openapi.License{Name: "foobar"}, false},
		{"withURL", openapi.License{URL: exampleCom}, true},
		{"invalidURL", openapi.License{Name: "foobar", URL: "foobar"}, true},
		{"full", openapi.License{Name: "foobar", URL: exampleCom}, false},
	}
	testValidater(t, candidates)
}
