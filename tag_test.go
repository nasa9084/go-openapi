package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestTagValidate(t *testing.T) {
	nameRequiredError := openapi.ErrRequired{Target: "tag.name"}
	candidates := []candidate{
		{"empty", openapi.Tag{}, nameRequiredError},
		{"withEmptyExternalDocs", openapi.Tag{ExternalDocs: &openapi.ExternalDocumentation{}}, nameRequiredError},
		{"withValidExternalDocs", openapi.Tag{ExternalDocs: &openapi.ExternalDocumentation{URL: exampleCom}}, nameRequiredError},

		{"withName", openapi.Tag{Name: "foo"}, nil},
		{"withNameAndEmptyExternalDocs", openapi.Tag{Name: "foo", ExternalDocs: &openapi.ExternalDocumentation{}}, openapi.ErrRequired{Target: "externalDocumentation.url"}},
		{"withNameAndValidExternalDocs", openapi.Tag{Name: "foo", ExternalDocs: &openapi.ExternalDocumentation{URL: exampleCom}}, nil},
	}
	testValidater(t, candidates)
}
