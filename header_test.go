package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestHeader_Validate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Header{}, nil},
		{"2 contents", openapi.Header{
			Content: map[string]*openapi.MediaType{
				"application/json": &openapi.MediaType{},
				"image/png":        &openapi.MediaType{},
			},
		}, openapi.ErrTooManyHeaderContent},
	}
	testValidater(t, candidates)
}
