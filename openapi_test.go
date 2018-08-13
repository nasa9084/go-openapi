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
	if err := doc.Validate(); err != nil {
		t.Error(err)
		return
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

func testPetStoreExpanded(t *testing.T) {
	doc, err := openapi.LoadFile("test/petstore-expanded.yaml")
	if err != nil {
		t.Fatal(err)
		return
	}
	if err := doc.Validate(); err != nil {
		t.Error(err)
		return
	}
	expect := openapi.Document{
		Version: "3.0.0",
		Info: &openapi.Info{
			Version:        "1.0.0",
			Title:          "Swagger Petstore",
			Description:    "A sample API that uses a petstore as an example to demonstrate features in the OpenAPI 3.0 specification",
			TermsOfService: "http://swagger.io/terms",
			Contact: &openapi.Contact{
				Name:  "Swagger API Team",
				Email: "apiteam@swagger.io",
				URL:   "http://swagger.io",
			},
			License: &openapi.License{
				Name: "Apache 2.0",
				URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
			},
		},
		Servers: []*openapi.Server{
			&openapi.Server{URL: "http://petstore.swagger.io/api"},
		},
		Paths: openapi.Paths{
			"/pets": &openapi.PathItem{
				Get: &openapi.Operation{
					Description: ` Returns all pets from the system that the user has access to
        Nam sed condimentum est. Maecenas tempor sagittis sapien, nec rhoncus sem sagittis sit amet. Aenean at gravida augue, ac iaculis sem. Curabitur odio lorem, ornare eget elementum nec, cursus id lectus. Duis mi turpis, pulvinar ac eros ac, tincidunt varius justo. In hac habitasse platea dictumst. Integer at adipiscing ante, a sagittis ligula. Aenean pharetra tempor ante molestie imperdiet. Vivamus id aliquam diam. Cras quis velit non tortor eleifend sagittis. Praesent at enim pharetra urna volutpat venenatis eget eget mauris. In eleifend fermentum facilisis. Praesent enim enim, gravida ac sodales sed, placerat id erat. Suspendisse lacus dolor, consectetur non augue vel, vehicula interdum libero. Morbi euismod sagittis libero sed lacinia.

        Sed tempus felis lobortis leo pulvinar rutrum. Nam mattis velit nisl, eu condimentum ligula luctus nec. Phasellus semper velit eget aliquet faucibus. In a mattis elit. Phasellus vel urna viverra, condimentum lorem id, rhoncus nibh. Ut pellentesque posuere elementum. Sed a varius odio. Morbi rhoncus ligula libero, vel eleifend nunc tristique vitae. Fusce et sem dui. Aenean nec scelerisque tortor. Fusce malesuada accumsan magna vel tempus. Quisque mollis felis eu dolor tristique, sit amet auctor felis gravida. Sed libero lorem, molestie sed nisl in, accumsan tempor nisi. Fusce sollicitudin massa ut lacinia mattis. Sed vel eleifend lorem. Pellentesque vitae felis pretium, pulvinar elit eu, euismod sapien.`,
					OperationID: "findPets",
					Parameters: []*openapi.Parameter{
						&openapi.Parameter{
							Name:        "tags",
							In:          "query",
							Description: "tags to filter by",
							Required:    false,
							Style:       "form",
							Schema: &openapi.Schema{
								Type: "array",
								Items: &openapi.Schema{
									Type: "string",
								},
							},
						},
						&openapi.Parameter{
							Name:        "limit",
							In:          "query",
							Description: "maximum number of results to return",
							Required:    false,
							Schema: &openapi.Schema{
								Type:   "integer",
								Format: "int32",
							},
						},
					},
					Responses: openapi.Responses{
						"200": &openapi.Response{
							Description: "pet response",
							Content: map[string]*openapi.MediaType{
								"application/json": &openapi.MediaType{
									Schema: &openapi.Schema{
										Type: "array",
										Items: &openapi.Schema{
											Ref: "#/components/schemas/Pet",
										},
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
					Description: "Creates a new pet in the store. Duplicates are allowed",
					OperationID: "addPet",
					RequestBody: &openapi.RequestBody{
						Description: "Pet to add to the store",
						Required:    true,
						Content: map[string]*openapi.MediaType{
							"application/json": &openapi.MediaType{
								Schema: &openapi.Schema{
									Ref: "#/components/schemas/NewPet",
								},
							},
						},
					},
					Responses: openapi.Responses{
						"200": &openapi.Response{
							Description: "pet response",
							Content: map[string]*openapi.MediaType{
								"application/json": &openapi.MediaType{
									Schema: &openapi.Schema{
										Ref: "#/components/schemas/Pet",
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
			"/pets/{id}": &openapi.PathItem{
				Get: &openapi.Operation{
					Description: "Returns a user based on a single ID, if the user does not have access to the pet",
					OperationID: "find pet by id",
					Parameters: []*openapi.Parameter{
						&openapi.Parameter{
							Name:        "id",
							In:          "path",
							Description: "ID of pet to fetch",
							Required:    true,
							Schema: &openapi.Schema{
								Type:   "integer",
								Format: "int64",
							},
						},
					},
					Responses: openapi.Responses{
						"200": &openapi.Response{
							Description: "pet response",
							Content: map[string]*openapi.MediaType{
								"application/json": &openapi.MediaType{
									Schema: &openapi.Schema{
										Ref: "#/component/schemas/Pet",
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
				Delete: &openapi.Operation{
					Description: "deletes a single pet based on the ID supplied",
					OperationID: "deletePet",
					Parameters: []*openapi.Parameter{
						&openapi.Parameter{
							Name:        "id",
							In:          "path",
							Description: "ID of pet to delete",
							Required:    true,
							Schema: &openapi.Schema{
								Type:   "integer",
								Format: "int64",
							},
						},
					},
					Responses: openapi.Responses{
						"204": &openapi.Response{
							Description: "pet deleted",
						},
						"default": &openapi.Response{
							Description: "unexpected error",
							Content: map[string]*openapi.MediaType{
								"application/json": &openapi.MediaType{
									Schema: &openapi.Schema{
										Ref: "#/componentes/schemas/Error",
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
					AllOf: []*openapi.Schema{
						&openapi.Schema{
							Ref: "#/components/schemas/NewPet",
						},
						&openapi.Schema{
							Required: []string{"id"},
							Properties: map[string]*openapi.Schema{
								"id": &openapi.Schema{
									Type:   "integer",
									Format: "int64",
								},
							},
						},
					},
				},
				"NewPet": &openapi.Schema{
					Required: []string{"name"},
					Properties: map[string]*openapi.Schema{
						"name": &openapi.Schema{
							Type: "string",
						},
						"tag": &openapi.Schema{
							Type: "string",
						},
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
