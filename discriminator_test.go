package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestDiscriminatorValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Discriminator{}, true},
		{"withPropertyName", openapi.Discriminator{PropertyName: "foobar"}, false},
	}
	testValidater(t, candidates)
}
