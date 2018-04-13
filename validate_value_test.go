package openapi

import "testing"

func TestHasDuplicatedParameter(t *testing.T) {
	t.Run("no duplicated param", testHasDuplicatedParameterFalse)
	t.Run("there's duplicated param", testHasDuplicatedParameterTrue)
}

func testHasDuplicatedParameterFalse(t *testing.T) {
	params := []*Parameter{
		&Parameter{Name: "foo", In: "header"},
		&Parameter{Name: "foo", In: "path", Required: true},
		&Parameter{Name: "bar", In: "path", Required: true},
	}
	if hasDuplicatedParameter(params) {
		t.Error("should return false")
	}
}

func testHasDuplicatedParameterTrue(t *testing.T) {
	params := []*Parameter{
		&Parameter{Name: "foo", In: "header"},
		&Parameter{Name: "foo", In: "header"},
	}
	if !hasDuplicatedParameter(params) {
		t.Error("should return true")
	}
}
