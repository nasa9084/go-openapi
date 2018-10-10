package openapi_test

import (
	"reflect"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
	yaml "gopkg.in/yaml.v2"
)

func TestContactValidate(t *testing.T) {
	urlRequiredError := openapi.ErrRequired{Target: "contact.url"}
	candidates := []candidate{
		{"empty", openapi.Contact{}, urlRequiredError},
		{"withURL", openapi.Contact{URL: exampleCom}, nil},
		{"invalidURL", openapi.Contact{URL: "foobar"}, openapi.ErrFormatInvalid{Target: "contact.url", Format: "URL"}},
		{"withEmail", openapi.Contact{Email: exampleMail}, urlRequiredError},
		{"valid", openapi.Contact{URL: exampleCom, Email: exampleMail}, nil},
		{"invalidEmail", openapi.Contact{URL: exampleCom, Email: "foobar"}, openapi.ErrFormatInvalid{Target: "contact.email", Format: "email"}},
	}

	testValidater(t, candidates)
}

func TestContactByExample(t *testing.T) {
	example := `name: API Support
url: http://www.example.com/support
email: support@example.com`
	contact := openapi.Contact{}
	if err := yaml.Unmarshal([]byte(example), &contact); err != nil {
		t.Error(err)
		return
	}
	if err := contact.Validate(); err != nil {
		t.Error(err)
		return
	}
	expected := openapi.Contact{
		Name:  "API Support",
		URL:   "http://www.example.com/support",
		Email: "support@example.com",
	}
	if !reflect.DeepEqual(contact, expected) {
		t.Errorf("%+v != %+v", contact, expected)
		return
	}
}
