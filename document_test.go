package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestDocumentValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Document{}, openapi.ErrRequired{Target: "openapi"}},
		{"withInvalidVersion", openapi.Document{Version: "1.0"}, openapi.ErrFormatInvalid{Target: "openapi version", Format: "X.Y.Z"}},
		{"withVersion", openapi.Document{Version: "3.0.0"}, openapi.ErrRequired{Target: "info"}},
		{"valid", openapi.Document{Version: "3.0.0", Info: &openapi.Info{Title: "foo", TermsOfService: exampleCom, Version: "1.0"}, Paths: openapi.Paths{}}, nil},
		{"noPaths", openapi.Document{Version: "3.0.0", Info: &openapi.Info{Title: "foo", TermsOfService: exampleCom, Version: "1.0"}}, openapi.ErrRequired{Target: "paths"}},
	}
	testValidater(t, candidates)
}

func TestValidateOASVersion(t *testing.T) {
	candidates := []struct {
		label  string
		in     string
		hasErr bool
	}{
		{"empty", "", true},
		{"invalidVersion", "foobar", true},
		{"swagger", "2.0", true},
		{"valid", "3.0.0", false},
		{"unsupportedVersion", "4.0.0", true},
		{"invalidMajorVersion", "foo.0.0", true},
		{"invalidMinorVersion", "0.bar.0", true},
		{"invalidPatchVersion", "0.0.baz", true},
	}
	for _, c := range candidates {
		if err := openapi.ValidateOASVersion(c.in); (err != nil) != c.hasErr {
			t.Log(c.label)
			if c.hasErr {
				t.Error("error should be occurred, but not")
				continue
			}
			t.Errorf("error should not be occured: %s", err)
		}
	}
}
