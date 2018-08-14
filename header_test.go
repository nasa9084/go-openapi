package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestHeaderValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Header{}, false},
		{"2 contents", openapi.Header{
			Content: map[string]*openapi.MediaType{
				"application/json": &openapi.MediaType{},
				"image/png":        &openapi.MediaType{},
			},
		}, true},
	}
	testValidater(t, candidates)
}
