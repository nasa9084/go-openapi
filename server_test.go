package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestServerValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Server{}, openapi.ErrRequired{Target: "server.url"}},
		{"invalidURL", openapi.Server{URL: "foobar%"}, openapi.ErrFormatInvalid{Target: "server.url"}},
		{"withURL", openapi.Server{URL: exampleCom}, nil},
	}
	testValidater(t, candidates)
}
