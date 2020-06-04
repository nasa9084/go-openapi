//nolint:funlen
package openapi

import (
	"io/ioutil"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestAPIWithExample(t *testing.T) {
	b, err := ioutil.ReadFile("test/testdata/api-with-examples.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var got OpenAPI
	if err := yaml.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}

	want := OpenAPI{
		openapi: "3.0.0",
		info: &Info{
			title:   "Simple API overview",
			version: "2.0.0",
		},
		paths: &Paths{
			paths: map[string]*PathItem{
				"/": {
					get: &Operation{
						operationID: "listVersionsv2",
						summary:     "List API versions",
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "200 response",
									content: map[string]*MediaType{
										"application/json": {
											examples: map[string]*Example{
												"foo": {
													value: map[string]interface{}{
														"versions": []interface{}{
															map[string]interface{}{
																"status":  "CURRENT",
																"updated": "2011-01-21T11:33:21Z",
																"id":      "v2.0",
																"links": []interface{}{
																	map[string]interface{}{
																		"href": "http://127.0.0.1:8774/v2/",
																		"rel":  "self",
																	},
																},
															},
															map[string]interface{}{
																"status":  "EXPERIMENTAL",
																"updated": "2013-07-23T11:33:21Z",
																"id":      "v3.0",
																"links": []interface{}{
																	map[string]interface{}{
																		"href": "http://127.0.0.1:8774/v3/",
																		"rel":  "self",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
								"300": {
									description: "300 response",
									content: map[string]*MediaType{
										"application/json": {
											examples: map[string]*Example{
												"foo": {
													value: `{
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
											},
										},
									},
								},
							},
						},
					},
				},
				"/v2": {
					get: &Operation{
						operationID: "getVersionDetailsv2",
						summary:     "Show API version details",
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "200 response",
									content: map[string]*MediaType{
										"application/json": {
											examples: map[string]*Example{
												"foo": {
													value: map[string]interface{}{
														"version": map[string]interface{}{
															"status":  "CURRENT",
															"updated": "2011-01-21T11:33:21Z",
															"media-types": []interface{}{
																map[string]interface{}{
																	"base": "application/xml",
																	"type": "application/vnd.openstack.compute+xml;version=2",
																},
																map[string]interface{}{
																	"base": "application/json",
																	"type": "application/vnd.openstack.compute+json;version=2",
																},
															},
															"id": "v2.0",
															"links": []interface{}{
																map[string]interface{}{
																	"href": "http://127.0.0.1:8774/v2/",
																	"rel":  "self",
																},
																map[string]interface{}{
																	"href": "http://docs.openstack.org/api/openstack-compute/2/os-compute-devguide-2.pdf",
																	"type": "application/pdf",
																	"rel":  "describedby",
																},
																map[string]interface{}{
																	"href": "http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl",
																	"type": "application/vnd.sun.wadl+xml",
																	"rel":  "describedby",
																},
																map[string]interface{}{
																	"href": "http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl",
																	"type": "application/vnd.sun.wadl+xml",
																	"rel":  "describedby",
																},
															},
														},
													},
												},
											},
										},
									},
								},
								"203": {
									description: "203 response",
									content: map[string]*MediaType{
										"application/json": {
											examples: map[string]*Example{
												"foo": {
													value: map[string]interface{}{
														"version": map[string]interface{}{
															"status":  "CURRENT",
															"updated": "2011-01-21T11:33:21Z",
															"media-types": []interface{}{
																map[string]interface{}{
																	"base": "application/xml",
																	"type": "application/vnd.openstack.compute+xml;version=2",
																},
																map[string]interface{}{
																	"base": "application/json",
																	"type": "application/vnd.openstack.compute+json;version=2",
																},
															},
															"id": "v2.0",
															"links": []interface{}{
																map[string]interface{}{
																	"href": "http://23.253.228.211:8774/v2/",
																	"rel":  "self",
																},
																map[string]interface{}{
																	"href": "http://docs.openstack.org/api/openstack-compute/2/os-compute-devguide-2.pdf",
																	"type": "application/pdf",
																	"rel":  "describedby",
																},
																map[string]interface{}{
																	"href": "http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl",
																	"type": "application/vnd.sun.wadl+xml",
																	"rel":  "describedby",
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	want.setRoot(&want)

	assertEqual(t, got, want)
}

func TestCallbackExample(t *testing.T) {
	b, err := ioutil.ReadFile("test/testdata/callback-example.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var got OpenAPI
	if err := yaml.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}

	//nolint:lll // cannot do better
	want := OpenAPI{
		openapi: "3.0.0",
		info: &Info{
			title:   "Callback Example",
			version: "1.0.0",
		},
		paths: &Paths{
			paths: map[string]*PathItem{
				"/streams": {
					post: &Operation{
						description: "subscribes a client to receive out-of-band data",
						parameters: []*Parameter{
							{
								name:        "callbackUrl",
								in:          "query",
								required:    true,
								description: "the location where data will be sent.  Must be network accessible\nby the source server",
								schema: &Schema{
									type_:   "string",
									format:  "uri",
									example: "https://tonys-server.com",
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"201": {
									description: "subscription successfully created",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												description: "subscription information",
												required: []string{
													"subscriptionId",
												},
												properties: map[string]*Schema{
													"subscriptionId": {
														description: "this unique identifier allows management of the subscription",
														type_:       "string",
														example:     "2531329f-fb09-4ef7-887e-84e648214436",
													},
												},
											},
										},
									},
								},
							},
						},
						callbacks: map[string]*Callback{
							"onData": {
								callback: map[string]*PathItem{
									"{$request.query.callbackUrl}/data": {
										post: &Operation{
											requestBody: &RequestBody{
												description: "subscription payload",
												content: map[string]*MediaType{
													"application/json": {
														schema: &Schema{
															type_: "object",
															properties: map[string]*Schema{
																"timestamp": {
																	type_:  "string",
																	format: "date-time",
																},
																"userData": {
																	type_: "string",
																},
															},
														},
													},
												},
											},
											responses: &Responses{
												responses: map[string]*Response{
													"202": {
														description: "Your server implementation should return this HTTP status code\nif the data was received successfully",
													},
													"204": {
														description: "Your server should return this HTTP status code if no longer interested\nin further updates",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	want.setRoot(&want)

	assertEqual(t, got, want)
}

func TestLinkExample(t *testing.T) {
	b, err := ioutil.ReadFile("test/testdata/link-example.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var got OpenAPI
	if err := yaml.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}

	want := OpenAPI{
		openapi: "3.0.0",
		info: &Info{
			title:   "Link Example",
			version: "1.0.0",
		},
		paths: &Paths{
			paths: map[string]*PathItem{
				"/2.0/users/{username}": {
					get: &Operation{
						operationID: "getUserByName",
						parameters: []*Parameter{
							{
								name:     "username",
								in:       "path",
								required: true,
								schema: &Schema{
									type_: "string",
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "The User",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												reference: "#/components/schemas/user",
											},
										},
									},
									links: map[string]*Link{
										"userRepositories": {
											reference: "#/components/links/UserRepositories",
										},
									},
								},
							},
						},
					},
				},
				"/2.0/repositories/{username}": {
					get: &Operation{
						operationID: "getRepositoriesByOwner",
						parameters: []*Parameter{
							{
								name:     "username",
								in:       "path",
								required: true,
								schema: &Schema{
									type_: "string",
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "repositories owned by the supplied user",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												type_: "array",
												items: &Schema{
													reference: "#/components/schemas/repository",
												},
											},
										},
									},
									links: map[string]*Link{
										"userRepository": {
											reference: "#/components/links/UserRepository",
										},
									},
								},
							},
						},
					},
				},
				"/2.0/repositories/{username}/{slug}": {
					get: &Operation{
						operationID: "getRepository",
						parameters: []*Parameter{
							{
								name:     "username",
								in:       "path",
								required: true,
								schema: &Schema{
									type_: "string",
								},
							},
							{
								name:     "slug",
								in:       "path",
								required: true,
								schema: &Schema{
									type_: "string",
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "The repository",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												reference: "#/components/schemas/repository",
											},
										},
									},
									links: map[string]*Link{
										"repositoryPullRequests": {
											reference: "#/components/links/RepositoryPullRequests",
										},
									},
								},
							},
						},
					},
				},
				"/2.0/repositories/{username}/{slug}/pullrequests": {
					get: &Operation{
						operationID: "getPullRequestsByRepository",
						parameters: []*Parameter{
							{
								name:     "username",
								in:       "path",
								required: true,
								schema: &Schema{
									type_: "string",
								},
							},
							{
								name:     "slug",
								in:       "path",
								required: true,
								schema: &Schema{
									type_: "string",
								},
							},
							{
								name: "state",
								in:   "query",
								schema: &Schema{
									type_: "string",
									enum: []string{
										"open",
										"merged",
										"declined",
									},
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "an array of pull request objects",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												type_: "array",
												items: &Schema{
													reference: "#/components/schemas/pullrequest",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				"/2.0/repositories/{username}/{slug}/pullrequests/{pid}": {
					get: &Operation{
						operationID: "getPullRequestsById",
						parameters: []*Parameter{
							{
								name:     "username",
								in:       "path",
								required: true,
								schema: &Schema{
									type_: "string",
								},
							},
							{
								name:     "slug",
								in:       "path",
								required: true,
								schema: &Schema{
									type_: "string",
								},
							},
							{
								name:     "pid",
								in:       "path",
								required: true,
								schema: &Schema{
									type_: "string",
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "a pull request object",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												reference: "#/components/schemas/pullrequest",
											},
										},
									},
									links: map[string]*Link{
										"pullRequestMerge": {
											reference: "#/components/links/PullRequestMerge",
										},
									},
								},
							},
						},
					},
				},
				"/2.0/repositories/{username}/{slug}/pullrequests/{pid}/merge": {
					post: &Operation{
						operationID: "mergePullRequest",
						parameters: []*Parameter{
							{
								name:     "username",
								in:       "path",
								required: true,
								schema: &Schema{
									type_: "string",
								},
							},
							{
								name:     "slug",
								in:       "path",
								required: true,
								schema: &Schema{
									type_: "string",
								},
							},
							{
								name:     "pid",
								in:       "path",
								required: true,
								schema: &Schema{
									type_: "string",
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"204": {
									description: "the PR was successfully merged",
								},
							},
						},
					},
				},
			},
		},
		components: &Components{
			links: map[string]*Link{
				"UserRepositories": {
					operationID: "getRepositoriesByOwner",
					parameters: map[string]interface{}{
						"username": "$response.body#/username",
					},
				},
				"UserRepository": {
					operationID: "getRepository",
					parameters: map[string]interface{}{
						"username": "$response.body#/owner/username",
						"slug":     "$response.body#/slug",
					},
				},
				"RepositoryPullRequests": {
					operationID: "getPullRequestsByRepository",
					parameters: map[string]interface{}{
						"username": "$response.body#/owner/username",
						"slug":     "$response.body#/slug",
					},
				},
				"PullRequestMerge": {
					operationID: "mergePullRequest",
					parameters: map[string]interface{}{
						"username": "$response.body#/author/username",
						"slug":     "$response.body#/repository/slug",
						"pid":      "$response.body#/id",
					},
				},
			},
			schemas: map[string]*Schema{
				"user": {
					type_: "object",
					properties: map[string]*Schema{
						"username": {
							type_: "string",
						},
						"uuid": {
							type_: "string",
						},
					},
				},
				"repository": {
					type_: "object",
					properties: map[string]*Schema{
						"slug": {
							type_: "string",
						},
						"owner": {
							reference: "#/components/schemas/user",
						},
					},
				},
				"pullrequest": {
					type_: "object",
					properties: map[string]*Schema{
						"id": {
							type_: "integer",
						},
						"title": {
							type_: "string",
						},
						"repository": {
							reference: "#/components/schemas/repository",
						},
						"author": {
							reference: "#/components/schemas/user",
						},
					},
				},
			},
		},
	}
	want.setRoot(&want)

	assertEqual(t, got, want)
}

func TestPetstoreExpanded(t *testing.T) {
	b, err := ioutil.ReadFile("test/testdata/petstore-expanded.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var got OpenAPI
	if err := yaml.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}

	//nolint:lll // cannot do better
	want := OpenAPI{
		openapi: "3.0.0",
		info: &Info{
			version:        "1.0.0",
			title:          "Swagger Petstore",
			description:    "A sample API that uses a petstore as an example to demonstrate features in the OpenAPI 3.0 specification",
			termsOfService: "http://swagger.io/terms/",
			contact: &Contact{
				name:  "Swagger API Team",
				email: "apiteam@swagger.io",
				url:   "http://swagger.io",
			},
			license: &License{
				name: "Apache 2.0",
				url:  "https://www.apache.org/licenses/LICENSE-2.0.html",
			},
		},
		servers: []*Server{
			{
				url: "http://petstore.swagger.io/api",
			},
		},
		paths: &Paths{
			paths: map[string]*PathItem{
				"/pets": {
					get: &Operation{
						description: `Returns all pets from the system that the user has access to
Nam sed condimentum est. Maecenas tempor sagittis sapien, nec rhoncus sem sagittis sit amet. Aenean at gravida augue, ac iaculis sem. Curabitur odio lorem, ornare eget elementum nec, cursus id lectus. Duis mi turpis, pulvinar ac eros ac, tincidunt varius justo. In hac habitasse platea dictumst. Integer at adipiscing ante, a sagittis ligula. Aenean pharetra tempor ante molestie imperdiet. Vivamus id aliquam diam. Cras quis velit non tortor eleifend sagittis. Praesent at enim pharetra urna volutpat venenatis eget eget mauris. In eleifend fermentum facilisis. Praesent enim enim, gravida ac sodales sed, placerat id erat. Suspendisse lacus dolor, consectetur non augue vel, vehicula interdum libero. Morbi euismod sagittis libero sed lacinia.

Sed tempus felis lobortis leo pulvinar rutrum. Nam mattis velit nisl, eu condimentum ligula luctus nec. Phasellus semper velit eget aliquet faucibus. In a mattis elit. Phasellus vel urna viverra, condimentum lorem id, rhoncus nibh. Ut pellentesque posuere elementum. Sed a varius odio. Morbi rhoncus ligula libero, vel eleifend nunc tristique vitae. Fusce et sem dui. Aenean nec scelerisque tortor. Fusce malesuada accumsan magna vel tempus. Quisque mollis felis eu dolor tristique, sit amet auctor felis gravida. Sed libero lorem, molestie sed nisl in, accumsan tempor nisi. Fusce sollicitudin massa ut lacinia mattis. Sed vel eleifend lorem. Pellentesque vitae felis pretium, pulvinar elit eu, euismod sapien.`,
						operationID: "findPets",
						parameters: []*Parameter{
							{
								name:        "tags",
								in:          "query",
								description: "tags to filter by",
								required:    false,
								style:       "form",
								schema: &Schema{
									type_: "array",
									items: &Schema{
										type_: "string",
									},
								},
							},
							{
								name:        "limit",
								in:          "query",
								description: "maximum number of results to return",
								required:    false,
								schema: &Schema{
									type_:  "integer",
									format: "int32",
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "pet response",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												type_: "array",
												items: &Schema{
													reference: "#/components/schemas/Pet",
												},
											},
										},
									},
								},
								"default": {
									description: "unexpected error",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												reference: "#/components/schemas/Error",
											},
										},
									},
								},
							},
						},
					},
					post: &Operation{
						description: "Creates a new pet in the store.  Duplicates are allowed",
						operationID: "addPet",
						requestBody: &RequestBody{
							description: "Pet to add to the store",
							required:    true,
							content: map[string]*MediaType{
								"application/json": {
									schema: &Schema{
										reference: "#/components/schemas/NewPet",
									},
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "pet response",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												reference: "#/components/schemas/Pet",
											},
										},
									},
								},
								"default": {
									description: "unexpected error",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												reference: "#/components/schemas/Error",
											},
										},
									},
								},
							},
						},
					},
				},
				"/pets/{id}": {
					get: &Operation{
						description: "Returns a user based on a single ID, if the user does not have access to the pet",
						operationID: "find pet by id",
						parameters: []*Parameter{
							{
								name:        "id",
								in:          "path",
								description: "ID of pet to fetch",
								required:    true,
								schema: &Schema{
									type_:  "integer",
									format: "int64",
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "pet response",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												reference: "#/components/schemas/Pet",
											},
										},
									},
								},
								"default": {
									description: "unexpected error",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												reference: "#/components/schemas/Error",
											},
										},
									},
								},
							},
						},
					},
					delete: &Operation{
						description: "deletes a single pet based on the ID supplied",
						operationID: "deletePet",
						parameters: []*Parameter{
							{
								name:        "id",
								in:          "path",
								description: "ID of pet to delete",
								required:    true,
								schema: &Schema{
									type_:  "integer",
									format: "int64",
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"204": {
									description: "pet deleted",
								},
								"default": {
									description: "unexpected error",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												reference: "#/components/schemas/Error",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		components: &Components{
			schemas: map[string]*Schema{
				"Pet": {
					allOf: []*Schema{
						{
							reference: "#/components/schemas/NewPet",
						},
						{
							type_:    "object",
							required: []string{"id"},
							properties: map[string]*Schema{
								"id": {
									type_:  "integer",
									format: "int64",
								},
							},
						},
					},
				},

				"NewPet": {
					type_:    "object",
					required: []string{"name"},
					properties: map[string]*Schema{
						"name": {
							type_: "string",
						},
						"tag": {
							type_: "string",
						},
					},
				},
				"Error": {
					type_:    "object",
					required: []string{"code", "message"},
					properties: map[string]*Schema{
						"code": {
							type_:  "integer",
							format: "int32",
						},
						"message": {
							type_: "string",
						},
					},
				},
			},
		},
	}
	want.setRoot(&want)

	assertEqual(t, got, want)
}

func TestPetstore(t *testing.T) {
	b, err := ioutil.ReadFile("test/testdata/petstore.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var got OpenAPI
	if err := yaml.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}

	want := OpenAPI{
		openapi: "3.0.0",
		info: &Info{
			version: "1.0.0",
			title:   "Swagger Petstore",
			license: &License{
				name: "MIT",
			},
		},
		servers: []*Server{
			{
				url: "http://petstore.swagger.io/v1",
			},
		},
		paths: &Paths{
			paths: map[string]*PathItem{
				"/pets": {
					get: &Operation{
						summary:     "List all pets",
						operationID: "listPets",
						tags:        []string{"pets"},
						parameters: []*Parameter{
							{
								name:        "limit",
								in:          "query",
								description: "How many items to return at one time (max 100)",
								required:    false,
								schema: &Schema{
									type_:  "integer",
									format: "int32",
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "A paged array of pets",
									headers: map[string]*Header{
										"x-next": {
											description: "A link to the next page of responses",
											schema: &Schema{
												type_: "string",
											},
										},
									},
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												reference: "#/components/schemas/Pets",
											},
										},
									},
								},
								"default": {
									description: "unexpected error",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												reference: "#/components/schemas/Error",
											},
										},
									},
								},
							},
						},
					},
					post: &Operation{
						summary:     "Create a pet",
						operationID: "createPets",
						tags:        []string{"pets"},
						responses: &Responses{
							responses: map[string]*Response{
								"201": {
									description: "Null response",
								},
								"default": {
									description: "unexpected error",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												reference: "#/components/schemas/Error",
											},
										},
									},
								},
							},
						},
					},
				},
				"/pets/{petId}": {
					get: &Operation{
						summary:     "Info for a specific pet",
						operationID: "showPetById",
						tags:        []string{"pets"},
						parameters: []*Parameter{
							{
								name:        "petId",
								in:          "path",
								required:    true,
								description: "The id of the pet to retrieve",
								schema: &Schema{
									type_: "string",
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "Expected response to a valid request",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												reference: "#/components/schemas/Pet",
											},
										},
									},
								},
								"default": {
									description: "unexpected error",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												reference: "#/components/schemas/Error",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		components: &Components{
			schemas: map[string]*Schema{
				"Pet": {
					type_:    "object",
					required: []string{"id", "name"},
					properties: map[string]*Schema{
						"id": {
							type_:  "integer",
							format: "int64",
						},
						"name": {
							type_: "string",
						},
						"tag": {
							type_: "string",
						},
					},
				},
				"Pets": {
					type_: "array",
					items: &Schema{
						reference: "#/components/schemas/Pet",
					},
				},
				"Error": {
					type_:    "object",
					required: []string{"code", "message"},
					properties: map[string]*Schema{
						"code": {
							type_:  "integer",
							format: "int32",
						},
						"message": {
							type_: "string",
						},
					},
				},
			},
		},
	}
	want.setRoot(&want)

	assertEqual(t, got, want)
}

func TestUspto(t *testing.T) {
	b, err := ioutil.ReadFile("test/testdata/uspto.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var got OpenAPI
	if err := yaml.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}

	//nolint:lll // cannot do better
	want := OpenAPI{
		openapi: "3.0.1",
		servers: []*Server{
			{
				url: "{scheme}://developer.uspto.gov/ds-api",
				variables: map[string]*ServerVariable{
					"scheme": {
						description: "The Data Set API is accessible via https and http",
						enum: []string{
							"https",
							"http",
						},
						default_: "https",
					},
				},
			},
		},
		info: &Info{
			description: `The Data Set API (DSAPI) allows the public users to discover and search USPTO exported data sets. This is a generic API that allows USPTO users to make any CSV based data files searchable through API. With the help of GET call, it returns the list of data fields that are searchable. With the help of POST call, data can be fetched based on the filters on the field names. Please note that POST call is used to search the actual data. The reason for the POST call is that it allows users to specify any complex search criteria without worry about the GET size limitations as well as encoding of the input parameters.`,
			version:     "1.0.0",
			title:       "USPTO Data Set API",
			contact: &Contact{
				name:  "Open Data Portal",
				url:   "https://developer.uspto.gov",
				email: "developer@uspto.gov",
			},
		},
		tags: []*Tag{
			{
				name:        "metadata",
				description: "Find out about the data sets",
			},
			{
				name:        "search",
				description: "Search a data set",
			},
		},
		paths: &Paths{
			paths: map[string]*PathItem{
				"/": {
					get: &Operation{
						tags:        []string{"metadata"},
						operationID: "list-data-sets",
						summary:     "List available data sets",
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "Returns a list of data sets",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												reference: "#/components/schemas/dataSetList",
											},
											example: map[string]interface{}{
												"total": uint64(2),
												"apis": []interface{}{
													map[string]interface{}{
														"apiKey":              "oa_citations",
														"apiVersionNumber":    "v1",
														"apiUrl":              "https://developer.uspto.gov/ds-api/oa_citations/v1/fields",
														"apiDocumentationUrl": "https://developer.uspto.gov/ds-api-docs/index.html?url=https://developer.uspto.gov/ds-api/swagger/docs/oa_citations.json",
													},
													map[string]interface{}{
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
				},
				"/{dataset}/{version}/fields": {
					get: &Operation{
						tags:        []string{"metadata"},
						summary:     `Provides the general information about the API and the list of fields that can be used to query the dataset.`,
						description: `This GET API returns the list of all the searchable field names that are in the oa_citations. Please see the 'fields' attribute which returns an array of field names. Each field or a combination of fields can be searched using the syntax options shown below.`,
						operationID: "list-searchable-fields",
						parameters: []*Parameter{
							{
								name:        "dataset",
								in:          "path",
								description: "Name of the dataset.",
								required:    true,
								example:     "oa_citations",
								schema: &Schema{
									type_: "string",
								},
							},
							{
								name:        "version",
								in:          "path",
								description: "Version of the dataset.",
								required:    true,
								example:     "v1",
								schema: &Schema{
									type_: "string",
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: `The dataset API for the given version is found and it is accessible to consume.`,
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												type_: "string",
											},
										},
									},
								},
								"404": {
									description: `The combination of dataset name and version is not found in the system or it is not published yet to be consumed by public.`,
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												type_: "string",
											},
										},
									},
								},
							},
						},
					},
				},
				"/{dataset}/{version}/records": {
					post: &Operation{
						tags:        []string{"search"},
						summary:     `Provides search capability for the data set with the given search criteria.`,
						description: `This API is based on Solr/Lucense Search. The data is indexed using SOLR. This GET API returns the list of all the searchable field names that are in the Solr Index. Please see the 'fields' attribute which returns an array of field names. Each field or a combination of fields can be searched using the Solr/Lucene Syntax. Please refer https://lucene.apache.org/core/3_6_2/queryparsersyntax.html#Overview for the query syntax. List of field names that are searchable can be determined using above GET api.`,
						operationID: "perform-search",
						parameters: []*Parameter{
							{
								name:        "version",
								in:          "path",
								description: "Version of the dataset.",
								required:    true,
								schema: &Schema{
									type_:    "string",
									default_: "v1",
								},
							},
							{
								name:        "dataset",
								in:          "path",
								description: "Name of the dataset. In this case, the default value is oa_citations",
								required:    true,
								schema: &Schema{
									type_:    "string",
									default_: "oa_citations",
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "successful operation",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												type_: "array",
												items: &Schema{
													type_: "object",
													additionalProperties: &Schema{
														type_: "object",
													},
												},
											},
										},
									},
								},
								"404": {
									description: "No matching record found for the given criteria.",
								},
							},
						},
						requestBody: &RequestBody{
							content: map[string]*MediaType{
								"application/x-www-form-urlencoded": {
									schema: &Schema{
										type_: "object",
										properties: map[string]*Schema{

											"criteria": {
												description: `Uses Lucene Query Syntax in the format of propertyName:value, propertyName:[num1 TO num2] and date range format: propertyName:[yyyyMMdd TO yyyyMMdd]. In the response please see the 'docs' element which has the list of record objects. Each record structure would consist of all the fields and their corresponding values.`, type_: "string",
												default_: "*:*",
											},
											"start": {
												description: "Starting record number. Default value is 0.",
												type_:       "integer",
												default_:    "0",
											},
											"rows": {description: `Specify number of rows to be returned. If you run the search with default values, in the response you will see 'numFound' attribute which will tell the number of records available in the dataset.`, type_: "integer", default_: "100"},
										},
										required: []string{"criteria"},
									},
								},
							},
						},
					},
				},
			},
		},
		components: &Components{
			schemas: map[string]*Schema{
				"dataSetList": {
					type_: "object",
					properties: map[string]*Schema{
						"total": {
							type_: "integer",
						},
						"apis": {
							type_: "array",
							items: &Schema{
								type_: "object",
								properties: map[string]*Schema{
									"apiKey": {
										type_:       "string",
										description: "To be used as a dataset parameter value",
									},
									"apiVersionNumber": {
										type_:       "string",
										description: "To be used as a version parameter value",
									},
									"apiUrl": {
										type_:       "string",
										format:      "uriref",
										description: "The URL describing the dataset's fields",
									},
									"apiDocumentationUrl": {
										type_:       "string",
										format:      "uriref",
										description: "A URL to the API console for each API",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	want.setRoot(&want)

	assertEqual(t, got, want)
}
