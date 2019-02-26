package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestResponsesValidate(t *testing.T) {
	validResp := &openapi.Response{Description: "foobar"}
	candidates := []candidate{
		{"empty", openapi.Responses{}, nil},
		{"hasInvalidStatus", openapi.Responses{"foobar": validResp}, openapi.ErrInvalidStatusCode},
		{"hasDefaultStatus", openapi.Responses{"default": validResp}, nil},
		{"hasWildcardStatus", openapi.Responses{"2XX": validResp}, nil},
		{"hasOKStatus", openapi.Responses{"200": validResp}, nil},
		{"hasEmptyResponse", openapi.Responses{"200": &openapi.Response{}}, openapi.ErrRequired{Target: "response.description"}},
	}
	testValidater(t, candidates)
}

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
