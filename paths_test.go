package openapi_test

import (
	"reflect"
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
		{"invalid", getPaths("/foo/bar", "foobar", "foobar"), openapi.ErrOperationIDDuplicated},
		{"valid", getPaths("/foo/bar", "foo", "bar"), nil},
	}
	testValidater(t, candidates)
}

func testPaths(t *testing.T) {
	candidates := []candidate{
		{"abs path", getPaths("/foo/bar", "foo", "bar"), nil},
		{"rel path", getPaths("foo/bar", "foo", "bar"), openapi.PathFormatError},
	}
	testValidater(t, candidates)
}

func TestPathsGetOperationByID(t *testing.T) {
	target := openapi.Paths{
		"/": &openapi.PathItem{
			Get: &openapi.Operation{OperationID: "bar"},
		},
	}
	candidates := []struct {
		paths  openapi.Paths
		opID   string
		expect *openapi.Operation
	}{
		{openapi.Paths{}, "", nil},
		{openapi.Paths{}, "foo", nil},
		{openapi.Paths{}, "bar", nil},
		{target, "", nil},
		{target, "foo", nil},
		{target, "bar", &openapi.Operation{OperationID: "bar"}},
	}
	for _, c := range candidates {
		op := c.paths.GetOperationByID(c.opID)
		if !reflect.DeepEqual(op, c.expect) {
			t.Errorf("%+v != %+v", op, c.expect)
		}
	}
}
