package openapi

import "testing"

func TestInfoValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", Info{}, true},
		{"withTitle", Info{Title: "foo"}, true},
		{"withTitleAndVersion", Info{Title: "foo", Version: "1.0"}, false},
		{"withInvalidToS", Info{Title: "foo", TermsOfService: "foobar", Version: "1.0"}, true},
	}
	testValidater(t, candidates)
}
