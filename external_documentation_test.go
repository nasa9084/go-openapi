package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestExternalDocumentationValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.ExternalDocumentation{}, openapi.ErrRequired{Target: "externalDocumentation.url"}},
		{"invalidURL", openapi.ExternalDocumentation{URL: "foobar"}, openapi.ErrFormatInvalid{Target: "externalDocumentation.url"}},
		{"valid", openapi.ExternalDocumentation{URL: exampleCom}, nil},
	}
	testValidater(t, candidates)
}
