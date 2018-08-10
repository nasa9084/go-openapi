package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestPathsValidate(t *testing.T) {
	t.Run("duplicate pathItem", testPathItemDuplicate)
	t.Run("test path", testPaths)
}

func getPaths(path, id1, id2 string) openapi.Paths {
	return openapi.Paths{
		path: &openapi.PathItem{
			Get:  &openapi.Operation{OperationID: id1, Responses: openapi.Responses{"200": &openapi.Response{Description: "foo"}}},
			Post: &openapi.Operation{OperationID: id2, Responses: openapi.Responses{"200": &openapi.Response{Description: "foo"}}},
		},
	}
}

func testPathItemDuplicate(t *testing.T) {
	candidates := []candidate{
		{"invalid", getPaths("/foo/bar", "foobar", "foobar"), true},
		{"valid", getPaths("/foo/bar", "foo", "bar"), false},
	}
	testValidater(t, candidates)
}

func testPaths(t *testing.T) {
	candidates := []candidate{
		{"abs path", getPaths("/foo/bar", "foo", "bar"), false},
		{"rel path", getPaths("foo/bar", "foo", "bar"), true},
	}
	testValidater(t, candidates)
}
