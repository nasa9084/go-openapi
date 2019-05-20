package openapi_test

import (
	"reflect"
	"strconv"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
	yaml "gopkg.in/yaml.v2"
)

func TestPaths_Validate(t *testing.T) {
	candidates := []candidate{
		{
			"has duplicated op id",
			openapi.Paths{
				"/foo": &openapi.PathItem{
					Get:  &openapi.Operation{OperationID: "foo", Responses: openapi.Responses{"200": &openapi.Response{Description: "foo"}}},
					Post: &openapi.Operation{OperationID: "foo", Responses: openapi.Responses{"200": &openapi.Response{Description: "foo"}}},
				},
			},
			openapi.ErrOperationIDDuplicated,
		},
		{
			"valid",
			openapi.Paths{
				"/foo/bar": &openapi.PathItem{
					Get:  &openapi.Operation{OperationID: "foo", Responses: openapi.Responses{"200": &openapi.Response{Description: "foo"}}},
					Post: &openapi.Operation{OperationID: "bar", Responses: openapi.Responses{"200": &openapi.Response{Description: "foo"}}},
				},
			},
			nil,
		},
		{
			"rel path",
			openapi.Paths{
				"foo/bar": &openapi.PathItem{
					Get:  &openapi.Operation{OperationID: "foo", Responses: openapi.Responses{"200": &openapi.Response{Description: "foo"}}},
					Post: &openapi.Operation{OperationID: "bar", Responses: openapi.Responses{"200": &openapi.Response{Description: "foo"}}},
				},
			},
			openapi.ErrPathFormat,
		},
		{
			"duplicated path",
			openapi.Paths{
				"/foo/{bar}": &openapi.PathItem{
					Get: &openapi.Operation{OperationID: "foo1", Responses: openapi.Responses{"200": &openapi.Response{Description: "foo1"}}},
				},
				"/foo/{baz}": &openapi.PathItem{
					Get: &openapi.Operation{OperationID: "foo2", Responses: openapi.Responses{"200": &openapi.Response{Description: "foo2"}}},
				},
			},
			openapi.ErrPathsDuplicated,
		},
	}
	testValidater(t, candidates)
}

func TestPaths_GetOperationByID(t *testing.T) {
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
