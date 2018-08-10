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
