package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestValidateStatusCode(t *testing.T) {
	candidates := []struct {
		label  string
		in     string
		hasErr bool
	}{
		{"empty", "", true},
		{"default", "default", false},
		{"200", "200", false},
		{"tooLow", "99", true},
		{"tooHigh", "600", true},
		{"wildcard", "2XX", false},
		{"invalid", "foobar", true},
	}
	for _, c := range candidates {
		if err := openapi.ValidateStatusCode(c.in); (err != nil) != c.hasErr {
			if c.hasErr {
				t.Error("error should be occurred, but not")
				continue
			}
			t.Errorf("error should not be occurred, but occurred: %s", err)
		}
	}
}
