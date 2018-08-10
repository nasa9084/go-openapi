package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestLicenseValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.License{}, true},
		{"withName", openapi.License{Name: "foobar"}, true},
		{"withURL", openapi.License{URL: exampleCom}, true},
		{"invalidURL", openapi.License{Name: "foobar", URL: "foobar"}, true},
		{"valid", openapi.License{Name: "foobar", URL: exampleCom}, false},
	}
	testValidater(t, candidates)
}
