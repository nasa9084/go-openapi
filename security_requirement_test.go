package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestSecurityRequirementValidate(t *testing.T) {
	candidates := []struct {
		label  string
		in     openapi.SecurityRequirement
		hasErr bool
	}{
		{"empty", openapi.SecurityRequirement{}, true},
	}
	for _, c := range candidates {
		if err := c.in.Validate(); (err != nil) == c.hasErr {
			t.Error(err)
			return
		}
	}
}
