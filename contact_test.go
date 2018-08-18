package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestContactValidate(t *testing.T) {
	urlRequiredError := openapi.ErrRequired{Target: "contact.url"}
	candidates := []candidate{
		{"empty", openapi.Contact{}, urlRequiredError},
		{"withURL", openapi.Contact{URL: exampleCom}, nil},
		{"invalidURL", openapi.Contact{URL: "foobar"}, openapi.ErrFormatInvalid{Target: "contact.url"}},
		{"withEmail", openapi.Contact{Email: exampleMail}, urlRequiredError},
		{"valid", openapi.Contact{URL: exampleCom, Email: exampleMail}, nil},
		{"invalidEmail", openapi.Contact{URL: exampleCom, Email: "foobar"}, openapi.EmailFormatError},
	}

	testValidater(t, candidates)
}
