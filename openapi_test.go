package openapi_test

import (
	"os"
	"reflect"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

var doc *openapi.Document

func TestMain(m *testing.M) {
	var err error
	doc, err = openapi.LoadFile("test/testspec.yaml")
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestLoadFile(t *testing.T) {
	t.Run("testspec.yml", testTestSpec)
	t.Run("petstore.yml", testPetStore)
}

func testTestSpec(t *testing.T) {
	doc, err := openapi.LoadFile("test/testspec.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if doc.Version != "3.0.0" {
		t.Errorf("api version is not valid")
		return
	}
	info := doc.Info
	if info.Title != "openapi specification test" {
		t.Errorf("info.title is not valid")
		return
	}
	if info.Version != "1.0" {
		t.Errorf("info.version is not valid")
		return
	}
	paths := doc.Paths
	if paths["/"].Get.Responses["200"].Description != "ok" {
		t.Errorf("paths./.get.responses.200.description is not valid")
		return
	}
}

func testPetStore(t *testing.T) {
	doc, err := openapi.LoadFile("test/petstore.yaml")
	if err != nil {
		t.Fatal(err)
	}
	expect := openapi.Document{
		Version: "3.0.0",
		Info: &openapi.Info{
			Version: "1.0.0",
			Title:   "Swagger Petstore",
			License: &openapi.License{
				Name: "MIT",
			},
		},
		Servers: []*openapi.Server{
			&openapi.Server{
				URL: "http://petstore.swagger.io/v1",
			},
		},
		Paths: openapi.Paths{
			"/pets": &openapi.PathItem{
				Get: &openapi.Operation{
					Summary:     "List all pets",
					OperationID: "listPets",
					Tags:        []string{"pets"},
					Parameters: []*openapi.Parameter{
						&openapi.Parameter{
							Name:        "limit",
							In:          "query",
							Description: "How many items to return at one time (max 100)",
							Required:    false,
							Schema: &openapi.Schema{
								Type:   "integer",
								Format: "int32",
							},
						},
					},
					Responses: openapi.Responses{
						"200": &openapi.Response{
							Description: "A paged array of pets",
							Headers: map[string]*openapi.Header{
								"x-next": &openapi.Header{
									Description: "A link to the next page of responses",
									Schema: &openapi.Schema{
										Type: "string",
									},
								},
							},
							Content: map[string]*openapi.MediaType{
								"application/json": &openapi.MediaType{
									Schema: &openapi.Schema{
										Ref: "#/components/schemas/Pets",
									},
								},
							},
						},
						"default": &openapi.Response{
							Description: "unexpected error",
							Content: map[string]*openapi.MediaType{
								"application/json": &openapi.MediaType{
									Schema: &openapi.Schema{
										Ref: "#/components/schemas/Error",
									},
								},
							},
						},
					},
				},
				Post: &openapi.Operation{
					Summary:     "Create a pet",
					OperationID: "createPets",
					Tags:        []string{"pets"},
					Responses: openapi.Responses{
						"201": &openapi.Response{
							Description: "Null response",
						},
						"default": &openapi.Response{
							Description: "unexpected error",
							Content: map[string]*openapi.MediaType{
								"application/json": &openapi.MediaType{
									Schema: &openapi.Schema{
										Ref: "#/components/schemas/Error",
									},
								},
							},
						},
					},
				},
			},
			"/pets/{petId}": &openapi.PathItem{
				Get: &openapi.Operation{
					Summary:     "Info for a specific pet",
					OperationID: "showPetById",
					Tags:        []string{"pets"},
					Parameters: []*openapi.Parameter{
						&openapi.Parameter{
							Name:        "petId",
							In:          "path",
							Required:    true,
							Description: "The id of the pet to retrieve",
							Schema: &openapi.Schema{
								Type: "string",
							},
						},
					},
					Responses: openapi.Responses{
						"200": &openapi.Response{
							Description: "Expected response to a valid request",
							Content: map[string]*openapi.MediaType{
								"application/json": &openapi.MediaType{
									Schema: &openapi.Schema{
										Ref: "#/components/schemas/Pets",
									},
								},
							},
						},
						"default": &openapi.Response{
							Description: "unexpected error",
							Content: map[string]*openapi.MediaType{
								"application/json": &openapi.MediaType{
									Schema: &openapi.Schema{
										Ref: "#/components/schemas/Error",
									},
								},
							},
						},
					},
				},
			},
		},
		Components: &openapi.Components{
			Schemas: map[string]*openapi.Schema{
				"Pet": &openapi.Schema{
					Required: []string{"id", "name"},
					Properties: map[string]*openapi.Schema{
						"id": &openapi.Schema{
							Type:   "integer",
							Format: "int64",
						},
						"name": &openapi.Schema{
							Type: "string",
						},
						"tag": &openapi.Schema{
							Type: "string",
						},
					},
				},
				"Pets": &openapi.Schema{
					Type: "array",
					Items: &openapi.Schema{
						Ref: "#/components/schemas/Pet",
					},
				},
				"Error": &openapi.Schema{
					Required: []string{"code", "message"},
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
			},
		},
	}
	if !reflect.DeepEqual(*doc, expect) {
		t.Errorf("document is not valid: %+v != %+v", doc, expect)
		if !reflect.DeepEqual(doc.Version, expect.Version) {
			t.Error("document.Version is not valid")
		}
		if !reflect.DeepEqual(doc.Info, expect.Info) {
			t.Error("document.Info is not valid")
		}
		if !reflect.DeepEqual(doc.Servers, expect.Servers) {
			t.Error("document.Servers is not valid")
		}
		if !reflect.DeepEqual(doc.Paths, expect.Paths) {
			t.Error("document.Paths is not valid")
		}
		if !reflect.DeepEqual(doc.Components, expect.Components) {
			t.Error("document.Components is not valid")
		}
		if !reflect.DeepEqual(doc.Security, expect.Security) {
			t.Error("document.Security is not valid")
		}
		if !reflect.DeepEqual(doc.Tags, expect.Tags) {
			t.Error("document.Tags is not valid")
		}
		if !reflect.DeepEqual(doc.ExternalDocs, expect.ExternalDocs) {
			t.Error("document.ExternalDocs is not valid")
		}
		return
	}
}
