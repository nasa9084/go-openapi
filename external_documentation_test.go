package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestExternalDocumentationValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.ExternalDocumentation{}, true},
		{"invalidURL", openapi.ExternalDocumentation{URL: "foobar"}, true},
		{"valid", openapi.ExternalDocumentation{URL: exampleCom}, false},
	}
	testValidater(t, candidates)
}
