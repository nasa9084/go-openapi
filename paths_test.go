package openapi_test

import (
	"reflect"
	"strconv"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
	yaml "gopkg.in/yaml.v2"
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
		{"rel path", getPaths("foo/bar", "foo", "bar"), openapi.ErrPathFormat},
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
	for i, c := range candidates {
		t.Run(strconv.Itoa(i)+"/"+c.opID, func(t *testing.T) {
			op := c.paths.GetOperationByID(c.opID)
			if !reflect.DeepEqual(op, c.expect) {
				t.Errorf("%+v != %+v", op, c.expect)
			}
		})
	}
}

func TestPathsByExample(t *testing.T) {
	example := `/pets:
  get:
    description: Returns all pets from the system that the user has access to
    responses:
      '200':
        description: A list of pets.
        content:
          application/json:
            schema:
              type: array
              items:
                $ref: '#/components/schemas/pet'`
	var paths openapi.Paths
	if err := yaml.Unmarshal([]byte(example), &paths); err != nil {
		t.Error(err)
		return
	}
	expect := openapi.Paths{
		"/pets": &openapi.PathItem{
			Get: &openapi.Operation{
				Description: "Returns all pets from the system that the user has access to",
				Responses: openapi.Responses{
					"200": &openapi.Response{
						Description: "A list of pets.",
						Content: map[string]*openapi.MediaType{
							"application/json": &openapi.MediaType{
								Schema: &openapi.Schema{
									Type: "array",
									Items: &openapi.Schema{
										Ref: "#/components/schemas/pet",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(paths, expect) {
		t.Errorf("%+v != %+v", paths, expect)
		return
	}
}
