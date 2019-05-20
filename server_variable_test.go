package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestServerVariable_Validate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.ServerVariable{}, openapi.ErrRequired{Target: "serverVariable.default"}},
		{"withDefault", openapi.ServerVariable{Default: "default"}, nil},
	}
	testValidater(t, candidates)
}
