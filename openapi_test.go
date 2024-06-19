package openapi_test

import (
	"os"
	"reflect"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestLoadFile(t *testing.T) {
	t.Run("testspec.yaml", testTestSpec)
	t.Run("securityRequirementRoot", testSecurityRequirementRootDoc)
	t.Run("petstore.yaml", testPetStore)
	t.Run("petstore-expanded.yaml", testPetStoreExpanded)
	t.Run("callback-example.yaml", testCallBackExample)
	t.Run("link-example.yaml", testLinkExample)
	t.Run("api-with-example.yaml", testAPIWithExample)
	t.Run("uspto.yaml", testUspto)
}

func testTestSpec(t *testing.T) {
	doc, err := openapi.LoadFile("testdata/testspec.yaml")
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

func eqDocument(t *testing.T, a, b openapi.Document) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("document is not valid: %+v != %+v", a, b)
		if !reflect.DeepEqual(a.Version, b.Version) {
			t.Log("document.Version is not valid")
		}
		if !reflect.DeepEqual(a.Info, b.Info) {
			t.Log("document.Info is not valid")
		}
		if !reflect.DeepEqual(a.Servers, b.Servers) {
			t.Log("document.Servers is not valid")
		}
		if !reflect.DeepEqual(a.Paths, b.Paths) {
			t.Log("document.Paths is not valid")
		}
		if !reflect.DeepEqual(a.Components, b.Components) {
			t.Log("document.Components is not valid")
		}
		if !reflect.DeepEqual(a.Security, b.Security) {
			t.Log("document.Security is not valid")
		}
		if !reflect.DeepEqual(a.Tags, b.Tags) {
			t.Log("document.Tags is not valid")
		}
		if !reflect.DeepEqual(a.ExternalDocs, b.ExternalDocs) {
			t.Log("document.ExternalDocs is not valid")
		}
	}
}

func testSecurityRequirementRootDoc(t *testing.T) {
	document := []byte(`---
openapi: 3.0.2
info:
  title: none
  version: 1.0.0
paths:
  /:
    get:
      description: foobarbaz
      security:
      - apikey: []
      responses:
        '200':
          description: ok
components:
  securitySchemes:
    apikey:
      name: X-API-Key
      type: apiKey
      in: header
`)
	doc, err := openapi.Load(document)
	if err != nil {
		t.Fatal(err)
	}
	if err := doc.Validate(); err != nil {
		t.Error(err)
	}
}

func testPetStore(t *testing.T) {
	doc, err := openapi.LoadFile("testdata/petstore.yaml")
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
	}
	expectPaths := openapi.Paths{
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
	}
	expectComponents := &openapi.Components{
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
	}
	expect.Paths = expectPaths
	expect.Components = expectComponents
	eqDocument(t, *doc, expect)
}

func testPetStoreExpanded(t *testing.T) {
	doc, err := openapi.LoadFile("testdata/petstore-expanded.yaml")
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
			TermsOfService: "http://swagger.io/terms/",
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
	}
	expectPaths := openapi.Paths{
		"/pets": &openapi.PathItem{
			Get: &openapi.Operation{
				Description: `Returns all pets from the system that the user has access to
Nam sed condimentum est. Maecenas tempor sagittis sapien, nec rhoncus sem sagittis sit amet. Aenean at gravida augue, ac iaculis sem. Curabitur odio lorem, ornare eget elementum nec, cursus id lectus. Duis mi turpis, pulvinar ac eros ac, tincidunt varius justo. In hac habitasse platea dictumst. Integer at adipiscing ante, a sagittis ligula. Aenean pharetra tempor ante molestie imperdiet. Vivamus id aliquam diam. Cras quis velit non tortor eleifend sagittis. Praesent at enim pharetra urna volutpat venenatis eget eget mauris. In eleifend fermentum facilisis. Praesent enim enim, gravida ac sodales sed, placerat id erat. Suspendisse lacus dolor, consectetur non augue vel, vehicula interdum libero. Morbi euismod sagittis libero sed lacinia.

Sed tempus felis lobortis leo pulvinar rutrum. Nam mattis velit nisl, eu condimentum ligula luctus nec. Phasellus semper velit eget aliquet faucibus. In a mattis elit. Phasellus vel urna viverra, condimentum lorem id, rhoncus nibh. Ut pellentesque posuere elementum. Sed a varius odio. Morbi rhoncus ligula libero, vel eleifend nunc tristique vitae. Fusce et sem dui. Aenean nec scelerisque tortor. Fusce malesuada accumsan magna vel tempus. Quisque mollis felis eu dolor tristique, sit amet auctor felis gravida. Sed libero lorem, molestie sed nisl in, accumsan tempor nisi. Fusce sollicitudin massa ut lacinia mattis. Sed vel eleifend lorem. Pellentesque vitae felis pretium, pulvinar elit eu, euismod sapien.
`,
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
				Description: "Creates a new pet in the store.  Duplicates are allowed",
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
									Ref: "#/components/schemas/Error",
								},
							},
						},
					},
				},
			},
		},
	}
	expectComponents := &openapi.Components{
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
	}
	expect.Paths = expectPaths
	expect.Components = expectComponents
	eqDocument(t, *doc, expect)
}

func testCallBackExample(t *testing.T) {
	doc, err := openapi.LoadFile("testdata/callback-example.yaml")
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
			Title:   "Callback Example",
			Version: "1.0.0",
		},
		Servers: []*openapi.Server{
			&openapi.Server{
				URL: "/",
			},
		},
	}
	expectPaths := openapi.Paths{
		"/streams": &openapi.PathItem{
			Post: &openapi.Operation{
				Description: "subscribes a client to receive out-of-band data",
				Parameters: []*openapi.Parameter{
					&openapi.Parameter{
						Name:     "callbackUrl",
						In:       "query",
						Required: true,
						Description: `the location where data will be sent.  Must be network accessible
by the source server
`,
						Schema: &openapi.Schema{
							Type:    "string",
							Format:  "uri",
							Example: "https://tonys-server.com",
						},
					},
				},
				Responses: openapi.Responses{
					"201": &openapi.Response{
						Description: "subscription successfully created",
						Content: map[string]*openapi.MediaType{
							"application/json": &openapi.MediaType{
								Schema: &openapi.Schema{
									Description: "subscription information",
									Required:    []string{"subscriptionId"},
									Properties: map[string]*openapi.Schema{
										"subscriptionId": &openapi.Schema{
											Description: "this unique identifier allows management of the subscription",
											Type:        "string",
											Example:     "2531329f-fb09-4ef7-887e-84e648214436",
										},
									},
								},
							},
						},
					},
				},
				Callbacks: map[string]*openapi.Callback{
					"onData": &openapi.Callback{
						"{$request.query.callbackUrl}/data": &openapi.PathItem{
							Post: &openapi.Operation{
								RequestBody: &openapi.RequestBody{
									Description: "subscription payload",
									Content: map[string]*openapi.MediaType{
										"application/json": &openapi.MediaType{
											Schema: &openapi.Schema{
												Properties: map[string]*openapi.Schema{
													"timestamp": &openapi.Schema{
														Type:   "string",
														Format: "date-time",
													},
													"userData": &openapi.Schema{
														Type: "string",
													},
												},
											},
										},
									},
								},
								Responses: openapi.Responses{
									"202": &openapi.Response{
										Description: `Your server implementation should return this HTTP status code
if the data was received successfully
`,
									},
									"204": &openapi.Response{
										Description: `Your server should return this HTTP status code if no longer interested
in further updates
`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	expect.Paths = expectPaths
	eqDocument(t, *doc, expect)
}

func testLinkExample(t *testing.T) {
	doc, err := openapi.LoadFile("testdata/link-example.yaml")
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
			Title:   "Link Example",
			Version: "1.0.0",
		},
		Servers: []*openapi.Server{
			&openapi.Server{
				URL: "/",
			},
		},
	}
	expectPaths := openapi.Paths{
		"/2.0/users/{username}": &openapi.PathItem{
			Get: &openapi.Operation{
				OperationID: "getUserByName",
				Parameters: []*openapi.Parameter{
					&openapi.Parameter{
						Name:     "username",
						In:       "path",
						Required: true,
						Schema: &openapi.Schema{
							Type: "string",
						},
					},
				},
				Responses: openapi.Responses{
					"200": &openapi.Response{
						Description: "The User",
						Content: map[string]*openapi.MediaType{
							"application/json": &openapi.MediaType{
								Schema: &openapi.Schema{
									Ref: "#/components/schemas/user",
								},
							},
						},
						Links: map[string]*openapi.Link{
							"userRepositories": &openapi.Link{
								Ref: "#/components/links/UserRepositories",
							},
						},
					},
				},
			},
		},
		"/2.0/repositories/{username}": &openapi.PathItem{
			Get: &openapi.Operation{
				OperationID: "getRepositoriesByOwner",
				Parameters: []*openapi.Parameter{
					&openapi.Parameter{
						Name:     "username",
						In:       "path",
						Required: true,
						Schema: &openapi.Schema{
							Type: "string",
						},
					},
				},
				Responses: openapi.Responses{
					"200": &openapi.Response{
						Description: "repositories owned by the supplied user",
						Content: map[string]*openapi.MediaType{
							"application/json": &openapi.MediaType{
								Schema: &openapi.Schema{
									Type: "array",
									Items: &openapi.Schema{
										Ref: "#/components/schemas/repository",
									},
								},
							},
						},
						Links: map[string]*openapi.Link{
							"userRepository": &openapi.Link{
								Ref: "#/components/links/UserRepository",
							},
						},
					},
				},
			},
		},
		"/2.0/repositories/{username}/{slug}": &openapi.PathItem{
			Get: &openapi.Operation{
				OperationID: "getRepository",
				Parameters: []*openapi.Parameter{
					&openapi.Parameter{
						Name:     "username",
						In:       "path",
						Required: true,
						Schema: &openapi.Schema{
							Type: "string",
						},
					},
					&openapi.Parameter{
						Name:     "slug",
						In:       "path",
						Required: true,
						Schema: &openapi.Schema{
							Type: "string",
						},
					},
				},
				Responses: openapi.Responses{
					"200": &openapi.Response{
						Description: "The repository",
						Content: map[string]*openapi.MediaType{
							"application/json": &openapi.MediaType{
								Schema: &openapi.Schema{
									Ref: "#/components/schemas/repository",
								},
							},
						},
						Links: map[string]*openapi.Link{
							"repositoryPullRequests": &openapi.Link{
								Ref: "#/components/links/RepositoryPullRequests",
							},
						},
					},
				},
			},
		},
		"/2.0/repositories/{username}/{slug}/pullrequests": &openapi.PathItem{
			Get: &openapi.Operation{
				OperationID: "getPullRequestsByRepository",
				Parameters: []*openapi.Parameter{
					&openapi.Parameter{
						Name:     "username",
						In:       "path",
						Required: true,
						Schema: &openapi.Schema{
							Type: "string",
						},
					},
					&openapi.Parameter{
						Name:     "slug",
						In:       "path",
						Required: true,
						Schema: &openapi.Schema{
							Type: "string",
						},
					},
					&openapi.Parameter{
						Name: "state",
						In:   "query",
						Schema: &openapi.Schema{
							Type: "string",
							Enum: []string{
								"open",
								"merged",
								"declined",
							},
						},
					},
				},
				Responses: openapi.Responses{
					"200": &openapi.Response{
						Description: "an array of pull request objects",
						Content: map[string]*openapi.MediaType{
							"application/json": &openapi.MediaType{
								Schema: &openapi.Schema{
									Type: "array",
									Items: &openapi.Schema{
										Ref: "#/components/schemas/pullrequest",
									},
								},
							},
						},
					},
				},
			},
		},
		"/2.0/repositories/{username}/{slug}/pullrequests/{pid}": &openapi.PathItem{
			Get: &openapi.Operation{
				OperationID: "getPullRequestsById",
				Parameters: []*openapi.Parameter{
					&openapi.Parameter{
						Name:     "username",
						In:       "path",
						Required: true,
						Schema: &openapi.Schema{
							Type: "string",
						},
					},
					&openapi.Parameter{
						Name:     "slug",
						In:       "path",
						Required: true,
						Schema: &openapi.Schema{
							Type: "string",
						},
					},
					&openapi.Parameter{
						Name:     "pid",
						In:       "path",
						Required: true,
						Schema: &openapi.Schema{
							Type: "string",
						},
					},
				},
				Responses: openapi.Responses{
					"200": &openapi.Response{
						Description: "a pull request object",
						Content: map[string]*openapi.MediaType{
							"application/json": &openapi.MediaType{
								Schema: &openapi.Schema{
									Ref: "#/components/schemas/pullrequest",
								},
							},
						},
						Links: map[string]*openapi.Link{
							"pullRequestMerge": &openapi.Link{
								Ref: "#/components/links/PullRequestMerge",
							},
						},
					},
				},
			},
		},
		"/2.0/repositories/{username}/{slug}/pullrequests/{pid}/merge": &openapi.PathItem{
			Post: &openapi.Operation{
				OperationID: "mergePullRequest",
				Parameters: []*openapi.Parameter{
					&openapi.Parameter{
						Name:     "username",
						In:       "path",
						Required: true,
						Schema: &openapi.Schema{
							Type: "string",
						},
					},
					&openapi.Parameter{
						Name:     "slug",
						In:       "path",
						Required: true,
						Schema: &openapi.Schema{
							Type: "string",
						},
					},
					&openapi.Parameter{
						Name:     "pid",
						In:       "path",
						Required: true,
						Schema: &openapi.Schema{
							Type: "string",
						},
					},
				},
				Responses: openapi.Responses{
					"204": &openapi.Response{
						Description: "the PR was successfully merged",
					},
				},
			},
		},
	}
	expectComponents := &openapi.Components{
		Links: map[string]*openapi.Link{
			"UserRepositories": &openapi.Link{
				OperationID: "getRepositoriesByOwner",
				Parameters: map[string]interface{}{
					"username": "$response.body#/username",
				},
			},
			"UserRepository": &openapi.Link{
				OperationID: "getRepository",
				Parameters: map[string]interface{}{
					"username": "$response.body#/owner/username",
					"slug":     "$response.body#/slug",
				},
			},
			"RepositoryPullRequests": &openapi.Link{
				OperationID: "getPullRequestsByRepository",
				Parameters: map[string]interface{}{
					"username": "$response.body#/owner/username",
					"slug":     "$response.body#/slug",
				},
			},
			"PullRequestMerge": &openapi.Link{
				OperationID: "mergePullRequest",
				Parameters: map[string]interface{}{
					"username": "$response.body#/author/username",
					"slug":     "$response.body#/repository/slug",
					"pid":      "$response.body#/id",
				},
			},
		},
		Schemas: map[string]*openapi.Schema{
			"user": &openapi.Schema{
				Type: "object",
				Properties: map[string]*openapi.Schema{
					"username": &openapi.Schema{
						Type: "string",
					},
					"uuid": &openapi.Schema{
						Type: "string",
					},
				},
			},
			"repository": &openapi.Schema{
				Type: "object",
				Properties: map[string]*openapi.Schema{
					"slug": &openapi.Schema{
						Type: "string",
					},
					"owner": &openapi.Schema{
						Ref: "#/components/schemas/user",
					},
				},
			},
			"pullrequest": &openapi.Schema{
				Type: "object",
				Properties: map[string]*openapi.Schema{
					"id": &openapi.Schema{
						Type: "integer",
					},
					"title": &openapi.Schema{
						Type: "string",
					},
					"repository": &openapi.Schema{
						Ref: "#/components/schemas/repository",
					},
					"author": &openapi.Schema{
						Ref: "#/components/schemas/user",
					},
				},
			},
		},
	}
	expect.Paths = expectPaths
	expect.Components = expectComponents
	eqDocument(t, *doc, expect)
}

func testAPIWithExample(t *testing.T) {
	doc, err := openapi.LoadFile("testdata/api-with-example.yaml")
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
			Title:   "Simple API overview",
			Version: "v2",
		},
		Servers: []*openapi.Server{
			&openapi.Server{
				URL: "/",
			},
		},
	}
	expectPaths := openapi.Paths{
		"/": &openapi.PathItem{
			Get: &openapi.Operation{
				OperationID: "listVersionsv2",
				Summary:     "List API versions",
				Responses: openapi.Responses{
					"200": &openapi.Response{
						Description: `200 response`,
						Content: map[string]*openapi.MediaType{
							"application/json": &openapi.MediaType{},
						},
					},
					"300": &openapi.Response{
						Description: `300 response`,
						Content: map[string]*openapi.MediaType{
							"application/json": &openapi.MediaType{},
						},
					},
				},
			},
		},
		"/v2": &openapi.PathItem{
			Get: &openapi.Operation{
				OperationID: "getVersionDetailsv2",
				Summary:     "Show API version details",
				Responses: openapi.Responses{
					"200": &openapi.Response{
						Description: `200 response`,
						Content: map[string]*openapi.MediaType{
							"application/json": &openapi.MediaType{},
						},
					},
					"203": &openapi.Response{
						Description: `203 response`,
						Content: map[string]*openapi.MediaType{
							"application/json": &openapi.MediaType{},
						},
					},
				},
			},
		},
	}
	expect.Paths = expectPaths
	root200Examples := map[string]*openapi.Example{
		"foo": &openapi.Example{
			Value: map[interface{}]interface{}{
				"versions": []interface{}{
					map[interface{}]interface{}{
						"status":  "CURRENT",
						"updated": "2011-01-21T11:33:21Z",
						"id":      "v2.0",
						"links": []interface{}{
							map[interface{}]interface{}{
								"href": "http://127.0.0.1:8774/v2/",
								"rel":  "self",
							},
						},
					},
					map[interface{}]interface{}{
						"status":  "EXPERIMENTAL",
						"updated": "2013-07-23T11:33:21Z",
						"id":      "v3.0",
						"links": []interface{}{
							map[interface{}]interface{}{
								"href": "http://127.0.0.1:8774/v3/",
								"rel":  "self",
							},
						},
					},
				},
			},
		},
	}
	root300Examples := map[string]*openapi.Example{
		"foo": &openapi.Example{
			Value: `{
 "versions": [
       {
         "status": "CURRENT",
         "updated": "2011-01-21T11:33:21Z",
         "id": "v2.0",
         "links": [
             {
                 "href": "http://127.0.0.1:8774/v2/",
                 "rel": "self"
             }
         ]
     },
     {
         "status": "EXPERIMENTAL",
         "updated": "2013-07-23T11:33:21Z",
         "id": "v3.0",
         "links": [
             {
                 "href": "http://127.0.0.1:8774/v3/",
                 "rel": "self"
             }
         ]
     }
 ]
}
`,
		},
	}
	v2200Examples := map[string]*openapi.Example{
		"foo": &openapi.Example{
			Value: map[interface{}]interface{}{
				"version": map[interface{}]interface{}{
					"status":  "CURRENT",
					"updated": "2011-01-21T11:33:21Z",
					"media-types": []interface{}{
						map[interface{}]interface{}{
							"base": "application/xml",
							"type": "application/vnd.openstack.compute+xml;version=2",
						},
						map[interface{}]interface{}{
							"base": "application/json",
							"type": "application/vnd.openstack.compute+json;version=2",
						},
					},
					"id": "v2.0",
					"links": []interface{}{
						map[interface{}]interface{}{
							"href": "http://127.0.0.1:8774/v2/",
							"rel":  "self",
						},
						map[interface{}]interface{}{
							"href": "http://docs.openstack.org/api/openstack-compute/2/os-compute-devguide-2.pdf",
							"type": "application/pdf",
							"rel":  "describedby",
						},
						map[interface{}]interface{}{
							"href": "http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl",
							"type": "application/vnd.sun.wadl+xml",
							"rel":  "describedby",
						},
						map[interface{}]interface{}{
							"href": "http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl",
							"type": "application/vnd.sun.wadl+xml",
							"rel":  "describedby",
						},
					},
				},
			},
		},
	}
	v2203Examples := map[string]*openapi.Example{
		"foo": &openapi.Example{
			Value: map[interface{}]interface{}{
				"version": map[interface{}]interface{}{
					"status":  "CURRENT",
					"updated": "2011-01-21T11:33:21Z",
					"media-types": []interface{}{
						map[interface{}]interface{}{
							"base": "application/xml",
							"type": "application/vnd.openstack.compute+xml;version=2",
						},
						map[interface{}]interface{}{
							"base": "application/json",
							"type": "application/vnd.openstack.compute+json;version=2",
						},
					},
					"id": "v2.0",
					"links": []interface{}{
						map[interface{}]interface{}{
							"href": "http://23.253.228.211:8774/v2/",
							"rel":  "self",
						},
						map[interface{}]interface{}{
							"href": "http://docs.openstack.org/api/openstack-compute/2/os-compute-devguide-2.pdf",
							"type": "application/pdf",
							"rel":  "describedby",
						},
						map[interface{}]interface{}{
							"href": "http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl",
							"type": "application/vnd.sun.wadl+xml",
							"rel":  "describedby",
						},
					},
				},
			},
		},
	}
	expect.Paths["/"].Get.Responses["200"].Content["application/json"].Examples = root200Examples
	expect.Paths["/"].Get.Responses["300"].Content["application/json"].Examples = root300Examples
	expect.Paths["/v2"].Get.Responses["200"].Content["application/json"].Examples = v2200Examples
	expect.Paths["/v2"].Get.Responses["203"].Content["application/json"].Examples = v2203Examples
	eqDocument(t, *doc, expect)
}

func testUspto(t *testing.T) {
	doc, err := openapi.LoadFile("testdata/uspto.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if err := doc.Validate(); err != nil {
		t.Error(err)
		return
	}
	expect := openapi.Document{
		Version: "3.0.1",
		Servers: []*openapi.Server{
			&openapi.Server{
				URL: "{scheme}://developer.uspto.gov/ds-api",
				Variables: map[string]*openapi.ServerVariable{
					"scheme": &openapi.ServerVariable{
						Description: "The Data Set API is accessible via https and http",
						Enum: []string{
							"https",
							"http",
						},
						Default: "https",
					},
				},
			},
		},
		Info: &openapi.Info{
			Description: `The Data Set API (DSAPI) allows the public users to discover and search USPTO exported data sets. This is a generic API that allows USPTO users to make any CSV based data files searchable through API. With the help of GET call, it returns the list of data fields that are searchable. With the help of POST call, data can be fetched based on the filters on the field names. Please note that POST call is used to search the actual data. The reason for the POST call is that it allows users to specify any complex search criteria without worry about the GET size limitations as well as encoding of the input parameters.`,
			Version:     "1.0.0",
			Title:       "USPTO Data Set API",
			Contact: &openapi.Contact{
				Name:  "Open Data Portal",
				URL:   "https://developer.uspto.gov",
				Email: "developer@uspto.gov",
			},
		},
		Tags: []*openapi.Tag{
			&openapi.Tag{
				Name:        "metadata",
				Description: "Find out about the data sets",
			},
			&openapi.Tag{
				Name:        "search",
				Description: "Search a data set",
			},
		},
	}
	expectPaths := openapi.Paths{
		"/": &openapi.PathItem{
			Get: &openapi.Operation{
				Tags:        []string{"metadata"},
				OperationID: "list-data-sets",
				Summary:     "List available data sets",
				Responses: openapi.Responses{
					"200": &openapi.Response{
						Description: "Returns a list of data sets",
						Content: map[string]*openapi.MediaType{
							"application/json": &openapi.MediaType{
								Schema: &openapi.Schema{
									Ref: "#/components/schemas/dataSetList",
								},
								Example: map[interface{}]interface{}{
									"total": 2,
									"apis": []interface{}{
										map[interface{}]interface{}{
											"apiKey":              "oa_citations",
											"apiVersionNumber":    "v1",
											"apiUrl":              "https://developer.uspto.gov/ds-api/oa_citations/v1/fields",
											"apiDocumentationUrl": "https://developer.uspto.gov/ds-api-docs/index.html?url=https://developer.uspto.gov/ds-api/swagger/docs/oa_citations.json",
										},
										map[interface{}]interface{}{
											"apiKey":              "cancer_moonshot",
											"apiVersionNumber":    "v1",
											"apiUrl":              "https://developer.uspto.gov/ds-api/cancer_moonshot/v1/fields",
											"apiDocumentationUrl": "https://developer.uspto.gov/ds-api-docs/index.html?url=https://developer.uspto.gov/ds-api/swagger/docs/cancer_moonshot.json",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"/{dataset}/{version}/fields": &openapi.PathItem{
			Get: &openapi.Operation{
				Tags:        []string{"metadata"},
				Summary:     `Provides the general information about the API and the list of fields that can be used to query the dataset.`,
				Description: `This GET API returns the list of all the searchable field names that are in the oa_citations. Please see the 'fields' attribute which returns an array of field names. Each field or a combination of fields can be searched using the syntax options shown below.`,
				OperationID: "list-searchable-fields",
				Parameters: []*openapi.Parameter{
					&openapi.Parameter{
						Name:        "dataset",
						In:          "path",
						Description: "Name of the dataset.",
						Required:    true,
						Example:     "oa_citations",
						Schema: &openapi.Schema{
							Type: "string",
						},
					},
					&openapi.Parameter{
						Name:        "version",
						In:          "path",
						Description: "Version of the dataset.",
						Required:    true,
						Example:     "v1",
						Schema: &openapi.Schema{
							Type: "string",
						},
					},
				},
				Responses: openapi.Responses{
					"200": &openapi.Response{
						Description: `The dataset API for the given version is found and it is accessible to consume.`,
						Content: map[string]*openapi.MediaType{
							"application/json": &openapi.MediaType{
								Schema: &openapi.Schema{
									Type: "string",
								},
							},
						},
					},
					"404": &openapi.Response{
						Description: `The combination of dataset name and version is not found in the system or it is not published yet to be consumed by public.`,
						Content: map[string]*openapi.MediaType{
							"application/json": &openapi.MediaType{
								Schema: &openapi.Schema{
									Type: "string",
								},
							},
						},
					},
				},
			},
		},
		"/{dataset}/{version}/records": &openapi.PathItem{
			Post: &openapi.Operation{
				Tags:        []string{"search"},
				Summary:     `Provides search capability for the data set with the given search criteria.`,
				Description: `This API is based on Solr/Lucense Search. The data is indexed using SOLR. This GET API returns the list of all the searchable field names that are in the Solr Index. Please see the 'fields' attribute which returns an array of field names. Each field or a combination of fields can be searched using the Solr/Lucene Syntax. Please refer https://lucene.apache.org/core/3_6_2/queryparsersyntax.html#Overview for the query syntax. List of field names that are searchable can be determined using above GET api.`,
				OperationID: "perform-search",
				Parameters: []*openapi.Parameter{
					&openapi.Parameter{
						Name:        "version",
						In:          "path",
						Description: "Version of the dataset.",
						Required:    true,
						Schema: &openapi.Schema{
							Type:    "string",
							Default: "v1",
						},
					},
					&openapi.Parameter{
						Name:        "dataset",
						In:          "path",
						Description: "Name of the dataset. In this case, the default value is oa_citations",
						Required:    true,
						Schema: &openapi.Schema{
							Type:    "string",
							Default: "oa_citations",
						},
					},
				},
				Responses: openapi.Responses{
					"200": &openapi.Response{
						Description: "successful operation",
						Content: map[string]*openapi.MediaType{
							"application/json": &openapi.MediaType{
								Schema: &openapi.Schema{
									Type: "array",
									Items: &openapi.Schema{
										Type: "object",
										AdditionalProperties: &openapi.Schema{
											Type: "object",
										},
									},
								},
							},
						},
					},
					"404": &openapi.Response{
						Description: "No matching record found for the given criteria.",
					},
				},
				RequestBody: &openapi.RequestBody{
					Content: map[string]*openapi.MediaType{
						"application/x-www-form-urlencoded": &openapi.MediaType{
							Schema: &openapi.Schema{
								Type: "object",
								Properties: map[string]*openapi.Schema{
									"criteria": &openapi.Schema{
										Description: `Uses Lucene Query Syntax in the format of propertyName:value, propertyName:[num1 TO num2] and date range format: propertyName:[yyyyMMdd TO yyyyMMdd]. In the response please see the 'docs' element which has the list of record objects. Each record structure would consist of all the fields and their corresponding values.`,
										Type:        "string",
										Default:     "*:*",
									},
									"start": &openapi.Schema{
										Description: "Starting record number. Default value is 0.",
										Type:        "integer",
										Default:     0,
									},
									"rows": &openapi.Schema{
										Description: `Specify number of rows to be returned. If you run the search with default values, in the response you will see 'numFound' attribute which will tell the number of records available in the dataset.`,
										Type:        "integer",
										Default:     100,
									},
								},
								Required: []string{"criteria"},
							},
						},
					},
				},
			},
		},
	}
	expect.Paths = expectPaths
	expectComponents := &openapi.Components{
		Schemas: map[string]*openapi.Schema{
			"dataSetList": &openapi.Schema{
				Type: "object",
				Properties: map[string]*openapi.Schema{
					"total": &openapi.Schema{
						Type: "integer",
					},
					"apis": &openapi.Schema{
						Type: "array",
						Items: &openapi.Schema{
							Type: "object",
							Properties: map[string]*openapi.Schema{
								"apiKey": &openapi.Schema{
									Type:        "string",
									Description: "To be used as a dataset parameter value",
								},
								"apiVersionNumber": &openapi.Schema{
									Type:        "string",
									Description: "To be used as a version parameter value",
								},
								"apiUrl": &openapi.Schema{
									Type:        "string",
									Format:      "uriref",
									Description: "The URL describing the dataset's fields",
								},
								"apiDocumentationUrl": &openapi.Schema{
									Type:        "string",
									Format:      "uriref",
									Description: "A URL to the API console for each API",
								},
							},
						},
					},
				},
			},
		},
	}
	expect.Components = expectComponents
	eqDocument(t, *doc, expect)
}
