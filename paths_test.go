package openapi

import "testing"

func TestPathsValidate(t *testing.T) {
	t.Run("duplicate pathItem", testPathItemDuplicate)
	t.Run("test path", testPaths)
}

func getPaths(path, id1, id2 string) Paths {
	return Paths{
		path: &PathItem{
			Get:  &Operation{OperationID: id1, Responses: Responses{"200": &Response{Description: "foo"}}},
			Post: &Operation{OperationID: id2, Responses: Responses{"200": &Response{Description: "foo"}}},
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
