package openapi_test

import (
	"reflect"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
	yaml "gopkg.in/yaml.v2"
)

func TestComponents(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Components{}, nil},
	}
	testValidater(t, candidates)
}

func TestComponentsValidateKeys(t *testing.T) {
	candidates := []struct {
		label string
		in    openapi.Components
		err   error
	}{
		{"empty", openapi.Components{}, nil},
		{"invalidKey", openapi.Components{Parameters: map[string]*openapi.Parameter{"@": &openapi.Parameter{}}}, openapi.ErrMapKeyFormat},
		{"validKey", openapi.Components{Parameters: map[string]*openapi.Parameter{"foo": &openapi.Parameter{}}}, nil},
	}
	for _, c := range candidates {
		if err := openapi.ValidateComponentKeys(c.in); err != c.err {
			t.Log(c.label)
			t.Errorf("error should be %s, but %s", c.err, err)
		}
	}
}
func TestReduceComponentKeys(t *testing.T) {
	candidates := []struct {
		label    string
		in       openapi.Components
		expected []string
	}{
		{"empty", openapi.Components{}, []string{}},
	}
	for _, c := range candidates {
		keys := openapi.ReduceComponentKeys(c.in)
		if !reflect.DeepEqual(keys, c.expected) {
			t.Log(c.label)
			t.Errorf("%+v != %+v", keys, c.expected)
		}
	}
}

func TestReduceComponentObjects(t *testing.T) {
	candidates := []struct {
		label    string
		in       openapi.Components
		expected []openapi.Validater
	}{
		{"empty", openapi.Components{}, []openapi.Validater{}},
	}
	for _, c := range candidates {
		objects := openapi.ReduceComponentObjects(c.in)
		if !reflect.DeepEqual(objects, c.expected) {
			t.Log(c.label)
			t.Errorf("%+v != %+v", objects, c.expected)
		}
	}
}

func TestComponentsByExample(t *testing.T) {
	example := `schemas:
  GeneralError:
    type: object
    properties:
      code:
        type: integer
        format: int32
      message:
        type: string
  Category:
    type: object
    properties:
      id:
        type: integer
        format: int64
      name:
        type: string
  Tag:
    type: object
    properties:
      id:
        type: integer
        format: int64
      name:
        type: string
parameters:
  skipParam:
    name: skip
    in: query
    description: number of items to skip
    required: true
    schema:
      type: integer
      format: int32
  limitParam:
    name: limit
    in: query
    description: max records to return
    required: true
    schema:
      type: integer
      format: int32
responses:
  NotFound:
    description: Entity not found.
  IllegalInput:
    description: Illegal input for operation.
  GeneralError:
    description: General Error
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/GeneralError'
securitySchemes:
  api_key:
    type: apiKey
    name: api_key
    in: header
  petstore_auth:
    type: oauth2
    flows:
      implicit:
        authorizationUrl: http://example.org/api/oauth/dialog
        scopes:
          write:pets: modify pets in your account
          read:pets: read your pets`
	var components openapi.Components
	if err := yaml.Unmarshal([]byte(example), &components); err != nil {
		t.Error(err)
		return
	}
	expect := openapi.Components{
		Schemas: map[string]*openapi.Schema{
			"GeneralError": &openapi.Schema{
				Type: "object",
				Properties: map[string]*openapi.Schema{
					"code": &openapi.Schema{
						Type:   "integer",
						Format: "int32",
					},
					"message": &openapi.Schema{
						Type: "string",
					},
				},
			},
			"Category": &openapi.Schema{
				Type: "object",
				Properties: map[string]*openapi.Schema{
					"id": &openapi.Schema{
						Type:   "integer",
						Format: "int64",
					},
					"name": &openapi.Schema{
						Type: "string",
					},
				},
			},
			"Tag": &openapi.Schema{
				Type: "object",
				Properties: map[string]*openapi.Schema{
					"id": &openapi.Schema{
						Type:   "integer",
						Format: "int64",
					},
					"name": &openapi.Schema{
						Type: "string",
					},
				},
			},
		},
		Parameters: map[string]*openapi.Parameter{
			"skipParam": &openapi.Parameter{
				Name:        "skip",
				In:          "query",
				Description: "number of items to skip",
				Required:    true,
				Schema: &openapi.Schema{
					Type:   "integer",
					Format: "int32",
				},
			},
			"limitParam": &openapi.Parameter{
				Name:        "limit",
				In:          "query",
				Description: "max records to return",
				Required:    true,
				Schema: &openapi.Schema{
					Type:   "integer",
					Format: "int32",
				},
			},
		},
		Responses: openapi.Responses{
			"NotFound": &openapi.Response{
				Description: "Entity not found.",
			},
			"IllegalInput": &openapi.Response{
				Description: "Illegal input for operation.",
			},
			"GeneralError": &openapi.Response{
				Description: "General Error",
				Content: map[string]*openapi.MediaType{
					"application/json": &openapi.MediaType{
						Schema: &openapi.Schema{
							Ref: "#/components/schemas/GeneralError",
						},
					},
				},
			},
		},
		SecuritySchemes: map[string]*openapi.SecurityScheme{
			"api_key": &openapi.SecurityScheme{
				Type: "apiKey",
				Name: "api_key",
				In:   "header",
			},
			"petstore_auth": &openapi.SecurityScheme{
				Type: "oauth2",
				Flows: &openapi.OAuthFlows{
					Implicit: &openapi.OAuthFlow{
						AuthorizationURL: "http://example.org/api/oauth/dialog",
						Scopes: map[string]string{
							"write:pets": "modify pets in your account",
							"read:pets":  "read your pets",
						},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(components, expect) {
		t.Errorf("%+v != %+v", components, expect)
		return
	}
}
