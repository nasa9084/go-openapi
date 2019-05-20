package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestCallback_Validate(t *testing.T) {
	candidates := []candidate{
		{"empty", &openapi.Callback{}, nil},
		{"zero string", &openapi.Callback{"": &openapi.PathItem{}}, openapi.ErrRuntimeExprFormat},
		{"example.com", &openapi.Callback{"http://example.com": &openapi.PathItem{}}, nil},
		{"with url expr", &openapi.Callback{"http://example.com/{$url}": &openapi.PathItem{}}, nil},
		{"with method expr", &openapi.Callback{"https://example.com/{$method}": &openapi.PathItem{}}, nil},
		{"with body and flag", &openapi.Callback{"http://example.com/{$request.body#/user/uuid}": &openapi.PathItem{}}, nil},
		{"with two expr", &openapi.Callback{"http://example.com/{$request.body#/user/uuid}/{$method}": &openapi.PathItem{}}, nil},
		{"with invalid expr", &openapi.Callback{"https://example.com/{foo}": &openapi.PathItem{}}, openapi.ErrRuntimeExprFormat},
		{"with empty flag", &openapi.Callback{"http://example.com/{$request.body#}": &openapi.PathItem{}}, openapi.ErrRuntimeExprFormat},
		{"with second invalid", &openapi.Callback{"http://example.com/{$request.body#/user/uuid}/{$value}": &openapi.PathItem{}}, openapi.ErrRuntimeExprFormat},
	}
	testValidater(t, candidates)
}
