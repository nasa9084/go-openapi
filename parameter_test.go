package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestHasDuplicatedParameter(t *testing.T) {
	t.Run("no duplicated param", testHasDuplicatedParameterFalse)
	t.Run("there's duplicated param", testHasDuplicatedParameterTrue)
}

func testHasDuplicatedParameterFalse(t *testing.T) {
	params := []*openapi.Parameter{
		&openapi.Parameter{Name: "foo", In: "header"},
		&openapi.Parameter{Name: "foo", In: "path", Required: true},
		&openapi.Parameter{Name: "bar", In: "path", Required: true},
	}
	if openapi.HasDuplicatedParameter(params) {
		t.Error("should return false")
	}
}

func testHasDuplicatedParameterTrue(t *testing.T) {
	params := []*openapi.Parameter{
		&openapi.Parameter{Name: "foo", In: "header"},
		&openapi.Parameter{Name: "foo", In: "header"},
	}
	if !openapi.HasDuplicatedParameter(params) {
		t.Error("should return true")
	}
}

func TestParameterValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Parameter{}, true},
		{"withName", openapi.Parameter{Name: "foo"}, true},
		{"withName-in", openapi.Parameter{Name: "foo", In: "path"}, true},
		{"withName-invalidIn", openapi.Parameter{Name: "foo", In: "bar"}, true},
		{"withName-inPath-notRequired", openapi.Parameter{Name: "foo", In: "path"}, true},
		{"withName-inPath-required", openapi.Parameter{Name: "foo", In: "path", Required: true}, false},
		{"allowEmptyValue-notQuery", openapi.Parameter{Name: "foo", In: "header", AllowEmptyValue: true}, true},
		{"allowEmptyValue-query", openapi.Parameter{Name: "foo", In: "query", AllowEmptyValue: true}, false},
	}
	testValidater(t, candidates)
}
