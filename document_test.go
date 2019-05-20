package openapi_test

import (
	"strconv"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestDocument_Validate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Document{}, openapi.ErrRequired{Target: "openapi"}},
		{"withInvalidVersion",
			openapi.Document{
				Version: "1.0",
				Info:    &openapi.Info{},
				Paths:   openapi.Paths{},
			},
			openapi.ErrFormatInvalid{Target: "openapi version", Format: "X.Y.Z"},
		},
		{"withVersion",
			openapi.Document{
				Version: "3.0.0",
			},
			openapi.ErrRequired{Target: "info"},
		},
		{"valid",
			openapi.Document{
				Version: "3.0.0",
				Info:    &openapi.Info{Title: "foo", TermsOfService: exampleCom, Version: "1.0"},
				Paths:   openapi.Paths{},
			},
			nil,
		},
		{"noPaths",
			openapi.Document{
				Version: "3.0.0",
				Info:    &openapi.Info{Title: "foo", TermsOfService: exampleCom, Version: "1.0"},
			},
			openapi.ErrRequired{Target: "paths"},
		},
	}
	testValidater(t, candidates)
}

func TestOASVersion(t *testing.T) {
	candidates := []struct {
		label string
		in    string
		err   error
	}{
		{"empty", "", openapi.ErrRequired{Target: "openapi"}},
		{"invalidVersion", "foobar", openapi.ErrFormatInvalid{Target: "openapi version", Format: "X.Y.Z"}},
		{"swagger", "2.0", openapi.ErrFormatInvalid{Target: "openapi version", Format: "X.Y.Z"}},
		{"valid", "3.0.0", nil},
		{"unsupportedVersion", "4.0.0", openapi.ErrUnsupportedVersion},
		{"invalidMajorVersion", "foo.0.0", openapi.ErrFormatInvalid{Target: "major part of openapi version"}},
		{"invalidMinorVersion", "0.bar.0", openapi.ErrFormatInvalid{Target: "minor part of openapi version"}},
		{"invalidPatchVersion", "0.0.baz", openapi.ErrFormatInvalid{Target: "patch part of openapi version"}},
	}
	for i, c := range candidates {
		t.Run(strconv.Itoa(i)+"/"+c.label, func(t *testing.T) {
			doc := openapi.Document{
				Version: c.in,
				Info:    &openapi.Info{Title: "foo", Version: "1.0"},
				Paths:   openapi.Paths{},
			}
			if err := doc.Validate(); err != c.err {
				if c.err != nil {
					t.Error("error should be occurred, but not")
					return
				}
				t.Errorf("error should not be occured: %s", err)
			}
		})
	}
}
