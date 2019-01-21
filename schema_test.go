package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestSchemaValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Schema{}, nil},
	}
	testValidater(t, candidates)
}
