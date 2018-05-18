package openapi

import "testing"

func TestTagValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", Tag{}, true},
		{"withEmptyExternalDocs", Tag{ExternalDocs: &ExternalDocumentation{}}, true},
		{"withValidExternalDocs", Tag{ExternalDocs: &ExternalDocumentation{URL: exampleCom}}, true},

		{"withName", Tag{Name: "foo"}, false},
		{"withNameAndEmptyExternalDocs", Tag{Name: "foo", ExternalDocs: &ExternalDocumentation{}}, true},
		{"withNameAndValidExternalDocs", Tag{Name: "foo", ExternalDocs: &ExternalDocumentation{URL: exampleCom}}, false},
	}
	testValidater(t, candidates)
}
