package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestServerVariableValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.ServerVariable{}, true},
		{"withDefault", openapi.ServerVariable{Default: "default"}, false},
	}
	testValidater(t, candidates)
}
