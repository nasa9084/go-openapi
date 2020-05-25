package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestOperationExampleUnmarshalYAML(t *testing.T) {
	yml := `tags:
- pet
summary: Updates a pet in the store with form data
operationId: updatePetWithForm
parameters:
- name: petId
  in: path
  description: ID of pet that needs to be updated
  required: true
  schema:
    type: string
requestBody:
  content:
    'application/x-www-form-urlencoded':
      schema:
       properties:
          name:
            description: Updated name of the pet
            type: string
          status:
            description: Updated status of the pet
            type: string
       required:
         - status
responses:
  '200':
    description: Pet updated.
    content:
      'application/json': {}
      'application/xml': {}
  '405':
    description: Method Not Allowed
    content:
      'application/json': {}
      'application/xml': {}
security:
- petstore_auth:
  - write:pets
  - read:pets`

	var got Operation
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := Operation{
		tags:        []string{"pet"},
		summary:     "Updates a pet in the store with form data",
		operationID: "updatePetWithForm",
		parameters: []*Parameter{
			{
				name:        "petId",
				in:          "path",
				description: "ID of pet that needs to be updated",
				required:    true,
				schema: &Schema{
					type_: "string",
				},
			},
		},
		requestBody: &RequestBody{
			content: map[string]*MediaType{
				"application/x-www-form-urlencoded": {
					schema: &Schema{
						properties: map[string]*Schema{
							"name": {
								description: "Updated name of the pet",
								type_:       "string",
							},
							"status": {
								description: "Updated status of the pet",
								type_:       "string",
							},
						},
						required: []string{"status"},
					},
				},
			},
		},
		responses: &Responses{
			responses: map[string]*Response{
				"200": {
					description: "Pet updated.",
					content: map[string]*MediaType{
						"application/json": {},
						"application/xml":  {},
					},
				},
				"405": {
					description: "Method Not Allowed",
					content: map[string]*MediaType{
						"application/json": {},
						"application/xml":  {},
					},
				},
			},
		},
		security: []*SecurityRequirement{
			{
				securityRequirement: map[string][]string{
					"petstore_auth": {"write:pets", "read:pets"},
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func TestOperationUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Operation
	}{
		{
			yml: `responses: {}
externalDocs:
  url: https://example.com`,
			want: Operation{
				responses: &Responses{},
				externalDocs: &ExternalDocumentation{
					url: "https://example.com",
				},
			},
		},
		{
			yml: `responses: {}
deprecated: true`,
			want: Operation{
				responses:  &Responses{},
				deprecated: true,
			},
		},
		{
			yml: `responses: {}
servers:
- url: example.com`,
			want: Operation{
				responses: &Responses{},
				servers: []*Server{
					{url: "example.com"},
				},
			},
		},
		{
			yml: `responses: {}
x-foo: bar`,
			want: Operation{
				responses: &Responses{},
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Operation
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestOperationUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `tags: ["foo"]`,
			want: ErrRequired("responses"),
		},
		{
			yml: `externalDocs: foo
responses: {}`,
			// externalDocs expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `parameters: foo
responses: {}`,
			// parameters expects an array
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml: `requestBody: foo
responses: {}`,
			// requestBody expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `responses: foo`,
			// responses expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `security: foo
responses: {}`,
			// security expects an array
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml: `servers: foo
responses: {}`,
			// servers expects an array
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml: `responses:
  "200":
    description: foobar
foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Operation{})
			assertSameError(t, got, tt.want)
		})
	}
}
