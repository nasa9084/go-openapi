package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestPathItemExampleUnmarshalYAML(t *testing.T) {
	yml := `get:
  description: Returns pets based on ID
  summary: Find pets by ID
  operationId: getPetsById
  responses:
    '200':
      description: pet response
      content:
        '*/*' :
          schema:
            type: array
            items:
              $ref: '#/components/schemas/Pet'
    default:
      description: error payload
      content:
        'text/html':
          schema:
            $ref: '#/components/schemas/ErrorModel'
parameters:
- name: id
  in: path
  description: ID of pet to use
  required: true
  schema:
    type: array
    # This is in example but maybe mistake
    # style: simple
    items:
      type: string  `

	var got PathItem
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := PathItem{
		get: &Operation{
			description: "Returns pets based on ID",
			summary:     "Find pets by ID",
			operationID: "getPetsById",
			responses: &Responses{
				responses: map[string]*Response{
					"200": {
						description: "pet response",
						content: map[string]*MediaType{
							"*/*": {
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
						description: "error payload",
						content: map[string]*MediaType{
							"text/html": {
								schema: &Schema{
									reference: "#/components/schemas/ErrorModel",
								},
							},
						},
					},
				},
			},
		},
		parameters: []*Parameter{
			{
				name:        "id",
				in:          "path",
				description: "ID of pet to use",
				required:    true,
				schema: &Schema{
					type_: "array",
					items: &Schema{
						type_: "string",
					},
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func TestPathItemUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want PathItem
	}{
		{
			yml: `summary: this is summary text
description: this is description`,
			want: PathItem{
				summary:     "this is summary text",
				description: "this is description",
			},
		},
		{
			yml: `x-foo: bar`,
			want: PathItem{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
		{
			yml: `put:
  responses:
    "200":
      description: foobar`,
			want: PathItem{
				put: &Operation{
					responses: &Responses{
						responses: map[string]*Response{
							"200": {
								description: "foobar",
							},
						},
					},
				},
			},
		},
		{
			yml: `options:
  responses:
    "200":
      description: foobar`,
			want: PathItem{
				options: &Operation{
					responses: &Responses{
						responses: map[string]*Response{
							"200": {
								description: "foobar",
							},
						},
					},
				},
			},
		},
		{
			yml: `head:
  responses:
    "200":
      description: foobar`,
			want: PathItem{
				head: &Operation{
					responses: &Responses{
						responses: map[string]*Response{
							"200": {
								description: "foobar",
							},
						},
					},
				},
			},
		},
		{
			yml: `patch:
  responses:
    "200":
      description: foobar`,
			want: PathItem{
				patch: &Operation{
					responses: &Responses{
						responses: map[string]*Response{
							"200": {
								description: "foobar",
							},
						},
					},
				},
			},
		},
		{
			yml: `trace:
  responses:
    "200":
      description: foobar`,
			want: PathItem{
				trace: &Operation{
					responses: &Responses{
						responses: map[string]*Response{
							"200": {
								description: "foobar",
							},
						},
					},
				},
			},
		},
		{
			yml: `servers:
- url: example.com`,
			want: PathItem{
				servers: []*Server{
					{
						url: "example.com",
					},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got PathItem
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestPathItemUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `get: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `put: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `post: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `delete: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `options: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `head: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `trace: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `servers: foo`,
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &PathItem{})
			assertSameError(t, got, tt.want)
		})
	}
}
