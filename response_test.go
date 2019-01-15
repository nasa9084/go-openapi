package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestResponseValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Response{}, openapi.ErrRequired{Target: "response.description"}},
		{"withDescription", openapi.Response{Description: "foobar"}, nil},
		{"withRef", openapi.Response{Ref: "#/component/responses/foo"}, nil},
	}
	testValidater(t, candidates)
}
