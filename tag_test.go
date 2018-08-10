package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestTagValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Tag{}, true},
		{"withEmptyExternalDocs", openapi.Tag{ExternalDocs: &openapi.ExternalDocumentation{}}, true},
		{"withValidExternalDocs", openapi.Tag{ExternalDocs: &openapi.ExternalDocumentation{URL: exampleCom}}, true},

		{"withName", openapi.Tag{Name: "foo"}, false},
		{"withNameAndEmptyExternalDocs", openapi.Tag{Name: "foo", ExternalDocs: &openapi.ExternalDocumentation{}}, true},
		{"withNameAndValidExternalDocs", openapi.Tag{Name: "foo", ExternalDocs: &openapi.ExternalDocumentation{URL: exampleCom}}, false},
	}
	testValidater(t, candidates)
}
