package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestRequestBodyValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.RequestBody{}, true},
		{"emptyContent", openapi.RequestBody{Content: map[string]*openapi.MediaType{}}, true},
		{"valid", openapi.RequestBody{Content: map[string]*openapi.MediaType{"application/json": &openapi.MediaType{}}}, false},
	}
	testValidater(t, candidates)
}
