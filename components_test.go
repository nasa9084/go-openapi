package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestComponentsExampleUnmarshalYAML(t *testing.T) {
	yml := `components:
  schemas:
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

	var target struct {
		Components *Components
	}

	if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
		t.Fatal(err)
	}

	got := target.Components
	want := &Components{
		schemas: map[string]*Schema{
			"GeneralError": {
				type_: "object",
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
			"Category": {
				type_: "object",
				properties: map[string]*Schema{
					"id": {
						type_:  "integer",
						format: "int64",
					},
					"name": {
						type_: "string",
					},
				},
			},
			"Tag": {
				type_: "object",
				properties: map[string]*Schema{
					"id": {
						type_:  "integer",
						format: "int64",
					},
					"name": {
						type_: "string",
					},
				},
			},
		},
		parameters: map[string]*Parameter{
			"skipParam": {
				name:        "skip",
				in:          "query",
				description: "number of items to skip",
				required:    true,
				schema: &Schema{
					type_:  "integer",
					format: "int32",
				},
			},
			"limitParam": {
				name:        "limit",
				in:          "query",
				description: "max records to return",
				required:    true,
				schema: &Schema{
					type_:  "integer",
					format: "int32",
				},
			},
		},
		responses: map[string]*Response{
			"NotFound": {
				description: "Entity not found.",
			},
			"IllegalInput": {
				description: "Illegal input for operation.",
			},
			"GeneralError": {
				description: "General Error",
				content: map[string]*MediaType{
					"application/json": {
						schema: &Schema{
							reference: "#/components/schemas/GeneralError",
						},
					},
				},
			},
		},
		securitySchemes: map[string]*SecurityScheme{
			"api_key": {
				type_: "apiKey",
				name:  "api_key",
				in:    "header",
			},
			"petstore_auth": {
				type_: "oauth2",
				flows: &OAuthFlows{
					implicit: &OAuthFlow{
						authorizationURL: "http://example.org/api/oauth/dialog",
						scopes: map[string]string{
							"write:pets": "modify pets in your account",
							"read:pets":  "read your pets",
						},
					},
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func TestComponentsUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Components
	}{
		{
			yml: `x-foo: bar`,
			want: Components{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
		{
			yml: `examples:
  fooExample:
    value: foo`,
			want: Components{
				examples: map[string]*Example{
					"fooExample": {
						value: "foo",
					},
				},
			},
		},
		{
			yml: `requestBodies:
  FooRequest:
    content:
      application/json: {}`,
			want: Components{
				requestBodies: map[string]*RequestBody{
					"FooRequest": {
						content: map[string]*MediaType{
							"application/json": {},
						},
					},
				},
			},
		},
		{
			yml: `headers:
  FooHeader:
    schema:
      type: string`,
			want: Components{
				headers: map[string]*Header{
					"FooHeader": {
						schema: &Schema{
							type_: "string",
						},
					},
				},
			},
		},
		{
			yml: `callbacks:
  CallBack:
    $request.body#/url:
      summary: foobar`,
			want: Components{
				callbacks: map[string]*Callback{
					"CallBack": {
						callback: map[string]*PathItem{
							"$request.body#/url": {
								summary: "foobar",
							},
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Components
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestComponentsUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
		{
			yml:  `schemas: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `responses: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `parameters: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `requestBodies: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `headers: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `links: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `callbacks: foo`,
			want: errors.New("String node doesn't MapNode"),
		},

		{
			yml: `examples:
  fooExample:
    foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Components{})
			assertSameError(t, got, tt.want)
		})
	}
}
