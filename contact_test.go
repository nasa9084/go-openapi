package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestContactValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Contact{}, true},
		{"withURL", openapi.Contact{URL: exampleCom}, false},
		{"invalidURL", openapi.Contact{URL: "foobar"}, true},
		{"withEmail", openapi.Contact{Email: exampleMail}, true},
		{"valid", openapi.Contact{URL: exampleCom, Email: exampleMail}, false},
		{"invalidEmail", openapi.Contact{URL: exampleCom, Email: "foobar"}, true},
	}

	testValidater(t, candidates)
}
