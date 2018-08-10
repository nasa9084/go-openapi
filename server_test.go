package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestServerValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Server{}, true},
		{"invalidURL", openapi.Server{URL: "foobar%"}, true},
		{"withURL", openapi.Server{URL: exampleCom}, false},
	}
	testValidater(t, candidates)
}
