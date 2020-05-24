//nolint:dupl,goconst,funlen
// dupl: test code contains so many duplicate code: e.g. unmarshal and validate the result
// goconst: test code contains so many magic numbers, esp. expected values
// funlen: test function tend to be long, because of definition of test cases
package openapi

import (
	"errors"
	"strconv"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestOpenAPIUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want OpenAPI
	}{
		{
			yml: `openapi: 3.0.2
info:
  title: info.title
  version: 1.0.0
paths:
  /:
    get:
      responses:
        '200':
          description: getOpResp
security:
  - api_key: []
`,
			want: OpenAPI{
				openapi: "3.0.2",
				info: &Info{
					title:   "info.title",
					version: "1.0.0",
				},
				paths: &Paths{
					paths: map[string]*PathItem{
						"/": {
							get: &Operation{
								responses: &Responses{
									responses: map[string]*Response{
										"200": {
											description: "getOpResp",
										},
									},
								},
							},
						},
					},
				},
				security: []*SecurityRequirement{
					{
						securityRequirement: map[string][]string{"api_key": {}},
					},
				},
			},
		},
		{
			yml: `openapi: 3.0.2
info:
  title: info.title
  version: 1.0.0
paths:
  /:
    get:
      responses:
        '200':
          description: getOpResp
tags:
  - name: fooTag
`,
			want: OpenAPI{
				openapi: "3.0.2",
				info: &Info{
					title:   "info.title",
					version: "1.0.0",
				},
				paths: &Paths{
					paths: map[string]*PathItem{
						"/": {
							get: &Operation{
								responses: &Responses{
									responses: map[string]*Response{
										"200": {
											description: "getOpResp",
										},
									},
								},
							},
						},
					},
				},
				tags: []*Tag{
					{
						name: "fooTag",
					},
				},
			},
		},
		{
			yml: `openapi: 3.0.2
info:
  title: info.title
  version: 1.0.0
paths:
  /:
    get:
      responses:
        '200':
          description: getOpResp
x-foo: bar
`,
			want: OpenAPI{
				openapi: "3.0.2",
				info: &Info{
					title:   "info.title",
					version: "1.0.0",
				},
				paths: &Paths{
					paths: map[string]*PathItem{
						"/": {
							get: &Operation{
								responses: &Responses{
									responses: map[string]*Response{
										"200": {
											description: "getOpResp",
										},
									},
								},
							},
						},
					},
				},
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
		{
			yml: `openapi: 3.0.2
info:
  title: title
  version: 1.0.0
paths: {}
externalDocs:
  url: https://example.com`,
			want: OpenAPI{
				openapi: "3.0.2",
				info: &Info{
					title:   "title",
					version: "1.0.0",
				},
				paths: &Paths{},
				externalDocs: &ExternalDocumentation{
					url: "https://example.com",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got OpenAPI
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			tt.want.setRoot(&tt.want)
			assertEqual(t, got, tt.want)
		})
	}
}

func TestOpenAPIUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml: `info:
  title: foobar
  version: 1.0.0`,
			want: ErrRequired("openapi"),
		},
		{
			yml:  `openapi: version1`,
			want: errors.New(`"openapi" field must be a valid semantic version but not`),
		},
		{
			yml:  `openapi: 3.0.2`,
			want: ErrRequired("info"),
		},
		{
			yml: `openapi: 3.0.2
info:
  title: foobar
  version: 1.0.0`,
			want: ErrRequired("paths"),
		},
		{
			yml: `openapi: 3.0.2
info:
  version: 1.0.0`,
			want: ErrRequired("title"),
		},
		{
			yml: `openapi: 3.0.2
info:
  title: foobar
  version: 1.0.0
servers: foo`,
			// servers expects an array
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml: `openapi: 3.0.2
info:
  title: foobar
  version: 1.0.0
paths: foo`,
			// paths expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `openapi: 3.0.2
info:
  title: foobar
  version: 1.0.0
paths: {}
components: foo`,
			// components expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `openapi: 3.0.2
info:
  title: foobar
  version: 1.0.0
paths: {}
security: foobar`,
			// security expects an array
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml: `openapi: 3.0.2
info:
  title: foobar
  version: 1.0.0
paths: {}
tags: foobar`,
			// tags expects an array
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml: `openapi: 3.0.2
info:
  title: foobar
  version: 1.0.0
paths: {}
externalDocs: foobar`,
			// externalDocs expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `openapi: 3.0.2
info:
  title: foobar
  version: 1.0.0
paths:
  /: {}
foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &OpenAPI{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestInfoExampleUnmarshalYAML(t *testing.T) {
	yml := `title: Sample Pet Store App
description: This is a sample server for a pet store.
termsOfService: http://example.com/terms/
contact:
  name: API Support
  url: http://www.example.com/support
  email: support@example.com
license:
  name: Apache 2.0
  url: https://www.apache.org/licenses/LICENSE-2.0.html
version: 1.0.1`

	var got Info
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := Info{
		title:          "Sample Pet Store App",
		description:    "This is a sample server for a pet store.",
		termsOfService: "http://example.com/terms/",
		contact: &Contact{
			name:  "API Support",
			url:   "http://www.example.com/support",
			email: "support@example.com",
		},
		license: &License{
			name: "Apache 2.0",
			url:  "https://www.apache.org/licenses/LICENSE-2.0.html",
		},
		version: "1.0.1",
	}
	assertEqual(t, got, want)
}

func TestInfoUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Info
	}{
		{
			yml: `title: this is title
version: 1.0.0`,
			want: Info{
				title:   "this is title",
				version: "1.0.0",
			},
		},
		{
			yml: `title: this is title
version: 1.0.0
x-foo: bar`,
			want: Info{
				title:   "this is title",
				version: "1.0.0",
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Info
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestInfoUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `version: 1.0.0`,
			want: ErrRequired("title"),
		},
		{
			yml:  `title: this is title`,
			want: ErrRequired("version"),
		},
		{
			yml: `ersion: 1.0.0
title: foobar
termsOfService: hoge`,
			// termsOfService expects URI
			want: errors.New(`parse "hoge": invalid URI for request`),
		},
		{
			yml: `ersion: 1.0.0
title: foobar
contact: hoge`,
			// contact expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `ersion: 1.0.0
title: foobar
license: hoge`,
			// license expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `title: this is title
version: 1.0.0
foo: bar`,
			want: errors.New(`unknown key: foo`),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Info{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestContactExampleUnmarshalYAML(t *testing.T) {
	yml := `name: API Support
url: http://www.example.com/support
email: support@example.com`

	var got Contact
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := Contact{
		name:  "API Support",
		url:   "http://www.example.com/support",
		email: "support@example.com",
	}
	assertEqual(t, got, want)
}

func TestContactUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Contact
	}{
		{
			yml: `email: foo@example.com`,
			want: Contact{
				email: "foo@example.com",
			},
		},
		{
			yml: `x-foo: bar`,
			want: Contact{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Contact
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestContactUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `email: invalidEmail`,
			want: errors.New(`"email" field must be an email address`),
		},
		{
			yml:  `url: foobar`,
			want: errors.New(`parse "foobar": invalid URI for request`),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Contact{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestLicenseExampleUnmarshalYAML(t *testing.T) {
	yml := `name: Apache 2.0
url: https://www.apache.org/licenses/LICENSE-2.0.html`

	var got License
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := License{
		name: "Apache 2.0",
		url:  "https://www.apache.org/licenses/LICENSE-2.0.html",
	}

	assertEqual(t, got, want)
}

func TestLicenseUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want License
	}{
		{
			yml: `name: licensename`,
			want: License{
				name: "licensename",
			},
		},
		{
			yml: `name: licensename
x-foo: bar`,
			want: License{
				name: "licensename",
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got License
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestLicenseUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `url: example.com`,
			want: ErrRequired("name"),
		},
		{
			yml: `name: foobar
url: hoge`,
			want: errors.New(`parse "hoge": invalid URI for request`),
		},
		{
			yml: `name: licensename
foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &License{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestServerExampleUnmarshalYAML(t *testing.T) {
	t.Run("single server", func(t *testing.T) {
		yml := `url: https://development.gigantic-server.com/v1
description: Development server`

		var got Server
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}

		want := Server{
			url:         "https://development.gigantic-server.com/v1",
			description: "Development server",
		}
		assertEqual(t, got, want)
	})
	t.Run("servers", func(t *testing.T) {
		yml := `servers:
- url: https://development.gigantic-server.com/v1
  description: Development server
- url: https://staging.gigantic-server.com/v1
  description: Staging server
- url: https://api.gigantic-server.com/v1
  description: Production server`

		var target struct {
			Servers []*Server
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		got := target.Servers
		want := []*Server{
			{
				url:         "https://development.gigantic-server.com/v1",
				description: "Development server",
			},
			{
				url:         "https://staging.gigantic-server.com/v1",
				description: "Staging server",
			},
			{
				url:         "https://api.gigantic-server.com/v1",
				description: "Production server",
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("with variables", func(t *testing.T) {
		yml := `servers:
- url: https://{username}.gigantic-server.com:{port}/{basePath}
  description: The production API server
  variables:
    username:
      # note! no enum here means it is an open value
      default: demo
      description: this value is assigned by the service provider, in this example "gigantic-server.com"
    port:
      enum:
        - '8443'
        - '443'
      default: '8443'
    basePath:
      # open meaning there is the opportunity to use special base paths as assigned by the provider, default is "v2"
      default: v2`

		var target struct {
			Servers []*Server
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}

		got := target.Servers[0]
		want := &Server{
			url:         "https://{username}.gigantic-server.com:{port}/{basePath}",
			description: "The production API server",
			variables: map[string]*ServerVariable{
				"username": {
					default_:    "demo",
					description: `this value is assigned by the service provider, in this example "gigantic-server.com"`,
				},
				"port": {
					enum:     []string{"8443", "443"},
					default_: "8443",
				},
				"basePath": {
					default_: "v2",
				},
			},
		}
		assertEqual(t, got, want)
	})
}

func TestServerUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Server
	}{
		{
			yml: `url: https://example.com`,
			want: Server{
				url: "https://example.com",
			},
		},
		{
			yml: `url: https://example.com
x-foo: bar`,
			want: Server{
				url: "https://example.com",
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Server
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestServerUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `description: foobar`,
			want: ErrRequired("url"),
		},
		{
			yml: `description: foobar
url: https%20://example.com`,
			want: errors.New(`parse "https%20://example.com": first path segment in URL cannot contain colon`),
		},
		{
			yml: `description: foobar
url: example.com
variables: hoge`,
			// variables expexts an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `url: example.com
foo: bar`,
			want: errors.New(`unknown key: foo`),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Server{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestServerVariableUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want ServerVariable
	}{
		{
			yml: `default: defaultValue`,
			want: ServerVariable{
				default_: "defaultValue",
			},
		},
		{
			yml: `default: defaultValue
x-foo: bar`,
			want: ServerVariable{
				default_: "defaultValue",
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got ServerVariable
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestServerVariableUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `description: foobar`,
			want: ErrRequired("default"),
		},
		{
			yml: `default: foo
enum: bar`,
			// enum expects an array
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml: `default: defaultValue
foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &ServerVariable{})
			assertSameError(t, got, tt.want)
		})
	}
}

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

func TestPathsExampleUnmarshalYAML(t *testing.T) {
	yml := `/pets:
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
	var got Paths
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := Paths{
		paths: map[string]*PathItem{
			"/pets": {
				get: &Operation{
					description: "Returns all pets from the system that the user has access to",
					responses: &Responses{
						responses: map[string]*Response{
							"200": {
								description: "A list of pets.",
								content: map[string]*MediaType{
									"application/json": {
										schema: &Schema{
											type_: "array",
											items: &Schema{
												reference: "#/components/schemas/pet",
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
	assertEqual(t, got, want)
}

func TestPathsUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Paths
	}{
		{
			yml: `x-foo: bar`,
			want: Paths{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Paths
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestPathsUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Paths{})
			assertSameError(t, got, tt.want)
		})
	}
}

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

func TestExternalDocumentationExampleUnmarshalYAML(t *testing.T) {
	yml := `description: Find more info here
url: https://example.com`
	var got ExternalDocumentation
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := ExternalDocumentation{
		description: "Find more info here",
		url:         "https://example.com",
	}
	assertEqual(t, got, want)
}

func TestExternalDocumentationUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want ExternalDocumentation
	}{
		{
			yml: `url: https://example.com`,
			want: ExternalDocumentation{
				url: "https://example.com",
			},
		},
		{
			yml: `url: https://example.com
x-foo: bar`,
			want: ExternalDocumentation{
				url: "https://example.com",
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got ExternalDocumentation
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestExternalDocumentationUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `description: foo`,
			want: ErrRequired("url"),
		},
		{
			yml:  `url: foobar`,
			want: errors.New(`parse "foobar": invalid URI for request`),
		},
		{
			yml: `url: https://example.com
foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &ExternalDocumentation{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestParameterExampleUnmarshalYAML(t *testing.T) {
	t.Run("header parameter", func(t *testing.T) {
		yml := `name: token
in: header
description: token to be passed as a header
required: true
schema:
  type: array
  items:
    type: integer
    format: int64
style: simple`
		var got Parameter
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := Parameter{
			name:        "token",
			in:          "header",
			description: "token to be passed as a header",
			required:    true,
			schema: &Schema{
				type_: "array",
				items: &Schema{
					type_:  "integer",
					format: "int64",
				},
			},
			style: "simple",
		}
		assertEqual(t, got, want)
	})
	t.Run("path parameter", func(t *testing.T) {
		yml := `name: username
in: path
description: username to fetch
required: true
schema:
  type: string`
		var got Parameter
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := Parameter{
			name:        "username",
			in:          "path",
			description: "username to fetch",
			required:    true,
			schema: &Schema{
				type_: "string",
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("optional query parameter", func(t *testing.T) {
		yml := `name: id
in: query
description: ID of the object to fetch
required: false
schema:
  type: array
  items:
    type: string
style: form
explode: true`
		var got Parameter
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := Parameter{
			name:        "id",
			in:          "query",
			description: "ID of the object to fetch",
			required:    false,
			schema: &Schema{
				type_: "array",
				items: &Schema{
					type_: "string",
				},
			},
			style:   "form",
			explode: true,
		}
		assertEqual(t, got, want)
	})
	t.Run("free form", func(t *testing.T) {
		yml := `in: query
name: freeForm
schema:
  type: object
  additionalProperties:
    type: integer
style: form`
		var got Parameter
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := Parameter{
			in:   "query",
			name: "freeForm",
			schema: &Schema{
				type_: "object",
				additionalProperties: &Schema{
					type_: "integer",
				},
			},
			style: "form",
		}
		assertEqual(t, got, want)
	})
	t.Run("complex parameter", func(t *testing.T) {
		yml := `in: query
name: coordinates
content:
  application/json:
    schema:
      type: object
      required:
        - lat
        - long
      properties:
        lat:
          type: number
        long:
          type: number`
		var got Parameter
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := Parameter{
			in:   "query",
			name: "coordinates",
			content: map[string]*MediaType{
				"application/json": {
					schema: &Schema{
						type_:    "object",
						required: []string{"lat", "long"},
						properties: map[string]*Schema{
							"lat": {
								type_: "number",
							},
							"long": {
								type_: "number",
							},
						},
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
}

func TestParameterUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Parameter
	}{
		{
			yml: `name: foo
in: path`,
			want: Parameter{
				name: "foo",
				in:   "path",
			},
		},
		{
			yml: `name: foo
in: path
deprecated: true`,
			want: Parameter{
				name:       "foo",
				in:         "path",
				deprecated: true,
			},
		},
		{
			yml: `name: foo
in: path
allowReserved: true`,
			want: Parameter{
				name:          "foo",
				in:            "path",
				allowReserved: true,
			},
		},
		{
			yml: `name: foo
in: path
examples:
  foo:
    value: bar`,
			want: Parameter{
				name: "foo",
				in:   "path",
				examples: map[string]*Example{
					"foo": {
						value: "bar",
					},
				},
			},
		},
		{
			yml: `name: foo
in: path
allowEmptyValue: true`,
			want: Parameter{
				name:            "foo",
				in:              "path",
				allowEmptyValue: true,
			},
		},
		{
			yml: `name: foo
in: path
x-foo: bar`,
			want: Parameter{
				name: "foo",
				in:   "path",
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
		{
			yml: `$ref: "#/components/parameters/Foo"`,
			want: Parameter{
				reference: "#/components/parameters/Foo",
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Parameter
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestParameterUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `in: path`,
			want: ErrRequired("name"),
		},
		{
			yml:  `name: foo`,
			want: ErrRequired("in"),
		},
		{
			yml: `name: foo
in: bar`,
			want: errors.New(`"in" field must be one of ["query", "header", "path", "cookie"]`),
		},
		{
			yml: `in: path
name: foo
schema: hoge`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `in: path
name: foo
examples: hoge`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `in: path
name: foo
content: hoge`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `name: namename
in: query
foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Parameter{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestRequestBodyExampleUnmarshalYAML(t *testing.T) {
	t.Run("with a referenced model", func(t *testing.T) {
		yml := `description: user to add to the system
content:
  'application/json':
    schema:
      $ref: '#/components/schemas/User'
    examples:
      user:
        summary: User Example
        externalValue: 'http://foo.bar/examples/user-example.json'
  'application/xml':
    schema:
      $ref: '#/components/schemas/User'
    examples:
      user:
        summary: User Example in XML
        externalValue: 'http://foo.bar/examples/user-example.xml'
  'text/plain':
    examples:
      user:
        summary: User example in text plain format
        externalValue: 'http://foo.bar/examples/user-example.txt'
  '*/*':
    examples:
      user:
        summary: User example in other format
        externalValue: 'http://foo.bar/examples/user-example.whatever'`
		var got RequestBody
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := RequestBody{
			description: "user to add to the system",
			content: map[string]*MediaType{
				"application/json": {
					schema: &Schema{
						reference: "#/components/schemas/User",
					},
					examples: map[string]*Example{
						"user": {
							summary:       "User Example",
							externalValue: "http://foo.bar/examples/user-example.json",
						},
					},
				},
				"application/xml": {
					schema: &Schema{
						reference: "#/components/schemas/User",
					},
					examples: map[string]*Example{
						"user": {
							summary:       "User Example in XML",
							externalValue: "http://foo.bar/examples/user-example.xml",
						},
					},
				},
				"text/plain": {
					examples: map[string]*Example{
						"user": {
							summary:       "User example in text plain format",
							externalValue: "http://foo.bar/examples/user-example.txt",
						},
					},
				},
				"*/*": {
					examples: map[string]*Example{
						"user": {
							summary:       "User example in other format",
							externalValue: "http://foo.bar/examples/user-example.whatever",
						},
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("array of string", func(t *testing.T) {
		yml := `description: user to add to the system
required: true
content:
  text/plain:
    schema:
      type: array
      items:
        type: string`
		var got RequestBody
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := RequestBody{
			description: "user to add to the system",
			required:    true,
			content: map[string]*MediaType{
				"text/plain": {
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
	})
}

func TestRequestBodyUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want RequestBody
	}{
		{
			yml: `$ref: "#/components/requestBodies/foo"`,
			want: RequestBody{
				reference: "#/components/requestBodies/foo",
			},
		},
		{
			yml: `content: {}
x-foo: bar`,
			want: RequestBody{
				content: map[string]*MediaType{},
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got RequestBody
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestRequestBodyUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `description: foo`,
			want: ErrRequired("content"),
		},
		{
			yml:  `content: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `content:
  application/json: {}
foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &RequestBody{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestMediaTypeExampleUnmarshalYAML(t *testing.T) {
	yml := `application/json:
  schema:
    $ref: "#/components/schemas/Pet"
  examples:
    cat:
      summary: An example of a cat
      value:
        name: Fluffy
        petType: Cat
        color: White
        gender: male
        breed: Persian
    dog:
      summary: An example of a dog with a cat's name
      value:
        name: Puma
        petType: Dog
        color: Black
        gender: Female
        breed: Mixed
    frog:
      $ref: "#/components/examples/frog-example"`
	var got map[string]*MediaType
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := map[string]*MediaType{
		"application/json": {
			schema: &Schema{
				reference: "#/components/schemas/Pet",
			},
			examples: map[string]*Example{
				"cat": {
					summary: "An example of a cat",
					value: map[string]interface{}{
						"name":    "Fluffy",
						"petType": "Cat",
						"color":   "White",
						"gender":  "male",
						"breed":   "Persian",
					},
				},
				"dog": {
					summary: "An example of a dog with a cat's name",
					value: map[string]interface{}{
						"name":    "Puma",
						"petType": "Dog",
						"color":   "Black",
						"gender":  "Female",
						"breed":   "Mixed",
					},
				},
				"frog": {
					reference: "#/components/examples/frog-example",
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func TestMediaTypeUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want MediaType
	}{
		{
			yml: `x-foo: bar`,
			want: MediaType{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got MediaType
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestMediaTypeUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `examples: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `encoding: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &MediaType{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestEncodingExampleUnmarshalYAML(t *testing.T) {
	yml := `requestBody:
  content:
    multipart/mixed:
      schema:
        type: object
        properties:
          id:
            # default is text/plain
            type: string
            format: uuid
          address:
            # default is application/json
            type: object
            properties: {}
          historyMetadata:
            # need to declare XML format!
            description: metadata in XML format
            type: object
            properties: {}
          profileImage:
            # default is application/octet-stream, need to declare an image type only!
            type: string
            format: binary
      encoding:
        historyMetadata:
          # require XML Content-Type in utf-8 encoding
          contentType: application/xml; charset=utf-8
        profileImage:
          # only accept png/jpeg
          contentType: image/png, image/jpeg
          headers:
            X-Rate-Limit-Limit:
              description: The number of allowed requests in the current period
              schema:
                type: integer`
	var target struct {
		RequestBody RequestBody `yaml:"requestBody"`
	}
	if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
		t.Fatal(err)
	}
	got := target.RequestBody
	want := RequestBody{
		content: map[string]*MediaType{
			"multipart/mixed": {
				schema: &Schema{
					type_: "object",
					properties: map[string]*Schema{
						"id": {
							type_:  "string",
							format: "uuid",
						},
						"address": {
							type_:      "object",
							properties: map[string]*Schema{},
						},
						"historyMetadata": {
							description: "metadata in XML format",
							type_:       "object",
							properties:  map[string]*Schema{},
						},
						"profileImage": {
							type_:  "string",
							format: "binary",
						},
					},
				},
				encoding: map[string]*Encoding{
					"historyMetadata": {
						contentType: "application/xml; charset=utf-8",
					},
					"profileImage": {
						contentType: "image/png, image/jpeg",
						headers: map[string]*Header{
							"X-Rate-Limit-Limit": {
								description: "The number of allowed requests in the current period",
								schema: &Schema{
									type_: "integer",
								},
							},
						},
					},
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func TestEncodingUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Encoding
	}{
		{
			yml: `style: foobar`,
			want: Encoding{
				style: "foobar",
			},
		},
		{
			yml: `explode: true`,
			want: Encoding{
				explode: true,
			},
		},
		{
			yml: `allowReserved: true`,
			want: Encoding{
				allowReserved: true,
			},
		},
		{
			yml: `x-foo: bar`,
			want: Encoding{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Encoding
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestEncodingUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `headers: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Encoding{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestResponsesExampleUnmarshalYAML(t *testing.T) {
	yml := `'200':
  description: a pet to be returned
  content:
    application/json:
      schema:
        $ref: '#/components/schemas/Pet'
default:
  description: Unexpected error
  content:
    application/json:
      schema:
        $ref: '#/components/schemas/ErrorModel'`
	var got Responses
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := Responses{
		responses: map[string]*Response{
			"200": {
				description: "a pet to be returned",
				content: map[string]*MediaType{
					"application/json": {
						schema: &Schema{
							reference: "#/components/schemas/Pet",
						},
					},
				},
			},
			"default": {
				description: "Unexpected error",
				content: map[string]*MediaType{
					"application/json": {
						schema: &Schema{
							reference: "#/components/schemas/ErrorModel",
						},
					},
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func TestResponsesUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Responses
	}{
		{
			yml: `"200":
  description: foobar`,
			want: Responses{
				responses: map[string]*Response{
					"200": {
						description: "foobar",
					},
				},
			},
		},
		{
			yml: `x-foo: bar`,
			want: Responses{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Responses
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestResponsesUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `"200": foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `"600":
  description: foobar`,
			want: ErrUnknownKey("600"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Responses{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestResponseExampleUnmarshalYAML(t *testing.T) {
	t.Run("array of complex type", func(t *testing.T) {
		yml := `description: A complex object array response
content:
  application/json:
    schema:
      type: array
      items:
        $ref: '#/components/schemas/VeryComplexType'`
		var got Response
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := Response{
			description: "A complex object array response",
			content: map[string]*MediaType{
				"application/json": {
					schema: &Schema{
						type_: "array",
						items: &Schema{
							reference: "#/components/schemas/VeryComplexType",
						},
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("string", func(t *testing.T) {
		yml := `description: A simple string response
content:
  text/plain:
    schema:
      type: string`
		var got Response
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := Response{
			description: "A simple string response",
			content: map[string]*MediaType{
				"text/plain": {
					schema: &Schema{
						type_: "string",
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("plain text with headers", func(t *testing.T) {
		yml := `description: A simple string response
content:
  text/plain:
    schema:
      type: string
    example: 'whoa!'
headers:
  X-Rate-Limit-Limit:
    description: The number of allowed requests in the current period
    schema:
      type: integer
  X-Rate-Limit-Remaining:
    description: The number of remaining requests in the current period
    schema:
      type: integer
  X-Rate-Limit-Reset:
    description: The number of seconds left in the current period
    schema:
      type: integer`
		var got Response
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := Response{
			description: "A simple string response",
			content: map[string]*MediaType{
				"text/plain": {
					schema: &Schema{
						type_: "string",
					},
					example: "whoa!",
				},
			},
			headers: map[string]*Header{
				"X-Rate-Limit-Limit": {
					description: "The number of allowed requests in the current period",
					schema: &Schema{
						type_: "integer",
					},
				},
				"X-Rate-Limit-Remaining": {
					description: "The number of remaining requests in the current period",
					schema: &Schema{
						type_: "integer",
					},
				},
				"X-Rate-Limit-Reset": {
					description: "The number of seconds left in the current period",
					schema: &Schema{
						type_: "integer",
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("no return value", func(t *testing.T) {
		yml := `description: object created`
		var got Response
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := Response{
			description: "object created",
		}
		assertEqual(t, got, want)
	})
}

func TestResponseUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Response
	}{
		{
			yml: `$ref: "#/components/responses/Foo"`,
			want: Response{
				reference: "#/components/responses/Foo",
			},
		},
		{
			yml: `description: foobar
x-foo: bar`,
			want: Response{
				description: "foobar",
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Response
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestResponseUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `headers: {}`,
			want: ErrRequired("description"),
		},
		{
			yml: `description: foo
headers: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `description: foo
content: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `description: foo
links: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `description: foo
foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Response{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestCallbackExampleUnmarshalYAML(t *testing.T) {
	yml := `myWebhook:
  'http://notificationServer.com?transactionId={$request.body#/id}&email={$request.body#/email}':
    post:
      requestBody:
        description: Callback payload
        content:
          'application/json':
            schema:
              $ref: '#/components/schemas/SomePayload'
      responses:
        '200':
          description: webhook successfully processed and no retries will be performed`
	var got map[string]*Callback
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := map[string]*Callback{
		"myWebhook": {
			callback: map[string]*PathItem{
				"http://notificationServer.com?transactionId={$request.body#/id}&email={$request.body#/email}": {
					post: &Operation{
						requestBody: &RequestBody{
							description: "Callback payload",
							content: map[string]*MediaType{
								"application/json": {
									schema: &Schema{
										reference: "#/components/schemas/SomePayload",
									},
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "webhook successfully processed and no retries will be performed",
								},
							},
						},
					},
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func TestCallbackUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Callback
	}{
		{
			yml: `x-foo: bar`,
			want: Callback{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
		{
			yml: `$ref: "#/components/callbacks/foo"`,
			want: Callback{
				reference: "#/components/callbacks/foo",
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Callback
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestCallbackUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `$url: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Callback{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestExampleExampleUnmarshalYAML(t *testing.T) {
	t.Run("in a request body", func(t *testing.T) {
		yml := `requestBody:
  content:
    'application/json':
      schema:
        $ref: '#/components/schemas/Address'
      examples:
        foo:
          summary: A foo example
          value: {"foo": "bar"}
        bar:
          summary: A bar example
          value: {"bar": "baz"}
    'application/xml':
      examples:
        xmlExample:
          summary: This is an example in XML
          externalValue: 'http://example.org/examples/address-example.xml'
    'text/plain':
      examples:
        textExample:
          summary: This is a text example
          externalValue: 'http://foo.bar/examples/address-example.txt'`
		var target struct {
			RequestBody RequestBody `yaml:"requestBody"`
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		got := target.RequestBody
		want := RequestBody{
			content: map[string]*MediaType{
				"application/json": {
					schema: &Schema{
						reference: "#/components/schemas/Address",
					},
					examples: map[string]*Example{
						"foo": {
							summary: "A foo example",
							value: map[string]interface{}{
								"foo": "bar",
							},
						},
						"bar": {
							summary: "A bar example",
							value: map[string]interface{}{
								"bar": "baz",
							},
						},
					},
				},
				"application/xml": {
					examples: map[string]*Example{
						"xmlExample": {
							summary:       "This is an example in XML",
							externalValue: "http://example.org/examples/address-example.xml",
						},
					},
				},
				"text/plain": {
					examples: map[string]*Example{
						"textExample": {
							summary:       "This is a text example",
							externalValue: "http://foo.bar/examples/address-example.txt",
						},
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("in a parameter", func(t *testing.T) {
		yml := `parameters:
  - name: 'zipCode'
    in: 'query'
    schema:
      type: 'string'
      format: 'zip-code'
    examples:
      zip-example:
        $ref: '#/components/examples/zip-example'`
		var target struct {
			Parameters []*Parameter
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		got := target.Parameters
		want := []*Parameter{
			{
				name: "zipCode",
				in:   "query",
				schema: &Schema{
					type_:  "string",
					format: "zip-code",
				},
				examples: map[string]*Example{
					"zip-example": {
						reference: "#/components/examples/zip-example",
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("in a response", func(t *testing.T) {
		yml := `responses:
  '200':
    description: your car appointment has been booked
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/SuccessResponse'
        examples:
          confirmation-success:
            $ref: '#/components/examples/confirmation-success'`
		var target struct {
			Responses Responses
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		got := target.Responses
		want := Responses{
			responses: map[string]*Response{
				"200": {
					description: "your car appointment has been booked",
					content: map[string]*MediaType{
						"application/json": {
							schema: &Schema{
								reference: "#/components/schemas/SuccessResponse",
							},
							examples: map[string]*Example{
								"confirmation-success": {
									reference: "#/components/examples/confirmation-success",
								},
							},
						},
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
}

func TestExampleUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Example
	}{
		{
			yml: `description: foobar`,
			want: Example{
				description: "foobar",
			},
		},
		{
			yml: `x-foo: bar`,
			want: Example{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Example
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestExampleUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Example{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestLinkExampleUnmarshalYAML(t *testing.T) {
	t.Run("$request.path.id", func(t *testing.T) {
		yml := `paths:
  /users/{id}:
    parameters:
    - name: id
      in: path
      required: true
      description: the user identifier, as userId
      schema:
        type: string
    get:
      responses:
        '200':
          description: the user being returned
          content:
            application/json:
              schema:
                type: object
                properties:
                  uuid: # the unique user id
                    type: string
                    format: uuid
          links:
            address:
              # the target link operationId
              operationId: getUserAddress
              parameters:
                # get the "id" field from the request path parameter named "id"
                userId: $request.path.id
  # the path item of the linked operation
  /users/{userid}/address:
    parameters:
    - name: userid
      in: path
      required: true
      description: the user identifier, as userId
      schema:
        type: string
    # linked operation
    get:
      operationId: getUserAddress
      responses:
        '200':
          description: the user's address`
		var target struct {
			Paths Paths
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		got := target.Paths
		want := Paths{
			paths: map[string]*PathItem{
				"/users/{id}": {
					parameters: []*Parameter{
						{
							name:        "id",
							in:          "path",
							required:    true,
							description: "the user identifier, as userId",
							schema:      &Schema{type_: "string"},
						},
					},
					get: &Operation{
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "the user being returned",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												type_: "object",
												properties: map[string]*Schema{
													"uuid": {
														type_:  "string",
														format: "uuid",
													},
												},
											},
										},
									},
									links: map[string]*Link{
										"address": {
											operationID: "getUserAddress",
											parameters: map[string]interface{}{
												"userId": "$request.path.id",
											},
										},
									},
								},
							},
						},
					},
				},
				"/users/{userid}/address": {
					parameters: []*Parameter{
						{
							name:        "userid",
							in:          "path",
							required:    true,
							description: "the user identifier, as userId",
							schema:      &Schema{type_: "string"},
						},
					},
					get: &Operation{
						operationID: "getUserAddress",
						responses: &Responses{
							responses: map[string]*Response{
								"200": {description: "the user's address"},
							},
						},
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("can use values from the response body", func(t *testing.T) {
		yml := `links:
  address:
    operationId: getUserAddressByUUID
    parameters:
      # get the "uuid" field from the "uuid" field in the response body
      userUuid: $response.body#/uuid`
		var target struct {
			Links map[string]*Link
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		got := target.Links
		want := map[string]*Link{
			"address": {
				operationID: "getUserAddressByUUID",
				parameters: map[string]interface{}{
					"userUuid": "$response.body#/uuid",
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("relative operationRef", func(t *testing.T) {
		yml := `links:
  UserRepositories:
    # returns array of '#/components/schemas/repository'
    operationRef: '#/paths/~12.0~1repositories~1{username}/get'
    parameters:
      username: $response.body#/username`
		var target struct {
			Links map[string]*Link
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		got := target.Links
		want := map[string]*Link{
			"UserRepositories": {
				operationRef: "#/paths/~12.0~1repositories~1{username}/get",
				parameters: map[string]interface{}{
					"username": "$response.body#/username",
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("absolute operationRef", func(t *testing.T) {
		yml := `links:
  UserRepositories:
    # returns array of '#/components/schemas/repository'
    operationRef: 'https://na2.gigantic-server.com/#/paths/~12.0~1repositories~1{username}/get'
    parameters:
      username: $response.body#/username`
		var target struct {
			Links map[string]*Link
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		got := target.Links
		want := map[string]*Link{
			"UserRepositories": {
				operationRef: "https://na2.gigantic-server.com/#/paths/~12.0~1repositories~1{username}/get",
				parameters: map[string]interface{}{
					"username": "$response.body#/username",
				},
			},
		}
		assertEqual(t, got, want)
	})
}

func TestLinkUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Link
	}{
		{
			yml: `requestBody: {}`,
			want: Link{
				requestBody: map[string]interface{}{},
			},
		},
		{
			yml: `description: foo`,
			want: Link{
				description: "foo",
			},
		},
		{
			yml: `server:
  url: example.com`,
			want: Link{
				server: &Server{
					url: "example.com",
				},
			},
		},
		{
			yml: `x-foo: bar`,
			want: Link{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Link
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestLinkUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `parameters: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `server: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Link{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestHeaderExampleUnmarshalYAML(t *testing.T) {
	yml := `description: The number of allowed requests in the current period
schema:
  type: integer`
	var got Header
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := Header{
		description: "The number of allowed requests in the current period",
		schema: &Schema{
			type_: "integer",
		},
	}
	assertEqual(t, got, want)
}

func TestHeaderUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Header
	}{
		{
			yml: `required: true`,
			want: Header{
				required: true,
			},
		},
		{
			yml: `deprecated: true`,
			want: Header{
				deprecated: true,
			},
		},
		{
			yml: `allowEmptyValue: true`,
			want: Header{
				allowEmptyValue: true,
			},
		},
		{
			yml: `style: foo`,
			want: Header{
				style: "foo",
			},
		},
		{
			yml: `explode: true`,
			want: Header{
				explode: true,
			},
		},
		{
			yml: `allowReserved: true`,
			want: Header{
				allowReserved: true,
			},
		},
		{
			yml: `x-foo: bar`,
			want: Header{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
		{
			yml: `example: foo`,
			want: Header{
				example: "foo",
			},
		},
		{
			yml: `examples:
  foo:
    value: bar`,
			want: Header{
				examples: map[string]*Example{
					"foo": {
						value: "bar",
					},
				},
			},
		},
		{
			yml: `content:
  application/json: {}`,
			want: Header{
				content: map[string]*MediaType{
					"application/json": {},
				},
			},
		},
		{
			yml: `$ref: '#/components/headers/foo'`,
			want: Header{
				reference: "#/components/headers/foo",
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Header
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestHeaderUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `schema: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `examples: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `content: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Header{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestTagExampleUnmarshal(t *testing.T) {
	yml := `name: pet
description: Pets operations`
	var got Tag
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := Tag{
		name:        "pet",
		description: "Pets operations",
	}
	assertEqual(t, got, want)
}

func TestTagUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Tag
	}{
		{
			yml: `name: theName
x-foo: bar`,
			want: Tag{
				name: "theName",
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
		{
			yml: `name: foo
externalDocs:
  url: https://example.com`,
			want: Tag{
				name: "foo",
				externalDocs: &ExternalDocumentation{
					url: "https://example.com",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Tag
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestTagUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml: `name: foo
externalDocs: bar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `description: foobar`,
			want: ErrRequired("name"),
		},
		{
			yml: `name: tagName
foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Tag{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestSchemaExampleUnmarshalYAML(t *testing.T) {
	t.Run("primitive", testSchemaExampleUnmarshalYAMLPrimitive)
	t.Run("simple model", testSchemaExampleUnmarshalYAMLSimpleModel)
	t.Run("simple string to string map", testSchemaExampleUnmarshalYAMLStringToStringMap)
	t.Run("string to model map", testSchemaExampleUnmarshalYAMLStringToModelMap)
	t.Run("model with example", testSchemaExampleUnmarshalYAMLModelExample)
	t.Run("models with composition", testSchemaExampleUnmarshalYAMLComposition)
	t.Run("models with polymorphism support", testSchemaExampleUnmarshalYAMLPolymorphism)
}

func testSchemaExampleUnmarshalYAMLPrimitive(t *testing.T) {
	yml := `type: string
format: email`
	var got Schema
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := Schema{
		type_:  "string",
		format: "email",
	}
	assertEqual(t, got, want)
}

func testSchemaExampleUnmarshalYAMLSimpleModel(t *testing.T) {
	yml := `type: object
required:
- name
properties:
  name:
    type: string
  address:
    $ref: '#/components/schemas/Address'
  age:
    type: integer
    format: int32
    minimum: 0`
	var got Schema
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := Schema{
		type_:    "object",
		required: []string{"name"},
		properties: map[string]*Schema{
			"name":    {type_: "string"},
			"address": {reference: "#/components/schemas/Address"},
			"age": {
				type_:   "integer",
				format:  "int32",
				minimum: 0,
			},
		},
	}
	assertEqual(t, got, want)
}

func testSchemaExampleUnmarshalYAMLStringToStringMap(t *testing.T) {
	yml := `type: object
additionalProperties:
  type: string`
	var got Schema
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := Schema{
		type_: "object",
		additionalProperties: &Schema{
			type_: "string",
		},
	}
	assertEqual(t, got, want)
}

func testSchemaExampleUnmarshalYAMLStringToModelMap(t *testing.T) {
	yml := `type: object
additionalProperties:
  $ref: '#/components/schemas/ComplexModel'`
	var got Schema
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := Schema{
		type_: "object",
		additionalProperties: &Schema{
			reference: "#/components/schemas/ComplexModel",
		},
	}
	assertEqual(t, got, want)
}

func testSchemaExampleUnmarshalYAMLModelExample(t *testing.T) {
	yml := `type: object
properties:
  id:
    type: integer
    format: int64
  name:
    type: string
required:
- name
example:
  name: Puma
  id: 1`
	var got Schema
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := Schema{
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
		required: []string{"name"},
		example: map[string]interface{}{
			"name": "Puma",
			"id":   uint64(1),
		},
	}
	assertEqual(t, got, want)
}

func testSchemaExampleUnmarshalYAMLComposition(t *testing.T) {
	yml := `components:
  schemas:
    ErrorModel:
      type: object
      required:
      - message
      - code
      properties:
        message:
          type: string
        code:
          type: integer
          minimum: 100
          maximum: 600
    ExtendedErrorModel:
      allOf:
      - $ref: '#/components/schemas/ErrorModel'
      - type: object
        required:
        - rootCause
        properties:
          rootCause:
            type: string`
	var target struct {
		Components Components
	}
	if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
		t.Fatal(err)
	}
	got := target.Components
	want := Components{
		schemas: map[string]*Schema{
			"ErrorModel": {
				type_:    "object",
				required: []string{"message", "code"},
				properties: map[string]*Schema{
					"message": {
						type_: "string",
					},
					"code": {
						type_:   "integer",
						minimum: 100,
						maximum: 600,
					},
				},
			},
			"ExtendedErrorModel": {
				allOf: []*Schema{
					{
						reference: "#/components/schemas/ErrorModel",
					},
					{
						type_:    "object",
						required: []string{"rootCause"},
						properties: map[string]*Schema{
							"rootCause": {type_: "string"},
						},
					},
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func testSchemaExampleUnmarshalYAMLPolymorphism(t *testing.T) {
	yml := `components:
  schemas:
    Pet:
      type: object
      discriminator:
        propertyName: petType
      properties:
        name:
          type: string
        petType:
          type: string
      required:
      - name
      - petType
    Cat:  ## "Cat" will be used as the discriminator value
      description: A representation of a cat
      allOf:
      - $ref: '#/components/schemas/Pet'
      - type: object
        properties:
          huntingSkill:
            type: string
            description: The measured skill for hunting
            enum:
            - clueless
            - lazy
            - adventurous
            - aggressive
        required:
        - huntingSkill
    Dog:  ## "Dog" will be used as the discriminator value
      description: A representation of a dog
      allOf:
      - $ref: '#/components/schemas/Pet'
      - type: object
        properties:
          packSize:
            type: integer
            format: int32
            description: the size of the pack the dog is from
            default: 0
            minimum: 0
        required:
        - packSize`
	var target struct {
		Components Components
	}
	if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
		t.Fatal(err)
	}
	got := target.Components
	want := Components{
		schemas: map[string]*Schema{
			"Pet": {
				type_: "object",
				discriminator: &Discriminator{
					propertyName: "petType",
				},
				properties: map[string]*Schema{
					"name":    {type_: "string"},
					"petType": {type_: "string"},
				},
				required: []string{"name", "petType"},
			},
			"Cat": {
				description: "A representation of a cat",
				allOf: []*Schema{
					{reference: "#/components/schemas/Pet"},
					{
						type_: "object",
						properties: map[string]*Schema{
							"huntingSkill": {
								type_:       "string",
								description: "The measured skill for hunting",
								enum:        []string{"clueless", "lazy", "adventurous", "aggressive"},
							},
						},
						required: []string{"huntingSkill"},
					},
				},
			},
			"Dog": {
				description: "A representation of a dog",
				allOf: []*Schema{
					{reference: "#/components/schemas/Pet"},
					{
						type_: "object",
						properties: map[string]*Schema{
							"packSize": {
								type_:       "integer",
								format:      "int32",
								description: "the size of the pack the dog is from",
								default_:    "0",
								minimum:     0,
							},
						},
						required: []string{"packSize"},
					},
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func TestSchemaUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Schema
	}{
		{
			yml: `title: Foo API`,
			want: Schema{
				title: "Foo API",
			},
		},
		{
			yml: `multipleOf: 3`,
			want: Schema{
				multipleOf: 3,
			},
		},
		{
			yml: `exclusiveMaximum: true`,
			want: Schema{
				exclusiveMaximum: true,
			},
		},
		{
			yml: `exclusiveMinimum: true`,
			want: Schema{
				exclusiveMinimum: true,
			},
		},
		{
			yml: `minLength: 3
maxLength: 6`,
			want: Schema{
				minLength: 3,
				maxLength: 6,
			},
		},
		{
			yml: `pattern: ^foo.+$`,
			want: Schema{
				pattern: "^foo.+$",
			},
		},
		{
			yml: `minItems: 3
maxItems: 6`,
			want: Schema{
				minItems: 3,
				maxItems: 6,
			},
		},
		{
			yml: `minProperties: 3
maxProperties: 6`,
			want: Schema{
				minProperties: 3,
				maxProperties: 6,
			},
		},
		{
			yml: `anyOf:
- description: foo`,
			want: Schema{
				anyOf: []*Schema{
					{
						description: "foo",
					},
				},
			},
		},
		{
			yml: `not:
  description: foo`,
			want: Schema{
				not: &Schema{
					description: "foo",
				},
			},
		},
		{
			yml: `nullable: true`,
			want: Schema{
				nullable: true,
			},
		},
		{
			yml: `writeOnly: true
readOnly: true`,
			want: Schema{
				writeOnly: true,
				readOnly:  true,
			},
		},
		{
			yml: `externalDocs:
url: https://example.com`,
			want: Schema{
				externalDocs: &ExternalDocumentation{
					url: "https://example.com",
				},
			},
		},
		{
			yml: `deprecated: true`,
			want: Schema{
				deprecated: true,
			},
		},
		{
			yml: `x-foo: bar`,
			want: Schema{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Schema
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestSchemaUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `enum: foobar`,
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml:  `allOf: foobar`,
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml:  `oneOf: foobar`,
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml:  `anyOf: foobar`,
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml:  `not: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `properties: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `additionalProperties: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `discriminator: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `xml: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `externalDocs: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Schema{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestDiscriminatorExampleUnmarshalYAML(t *testing.T) {
	t.Run("example", func(t *testing.T) {
		yml := `MyResponseType:
  oneOf:
  - $ref: '#/components/schemas/Cat'
  - $ref: '#/components/schemas/Dog'
  - $ref: '#/components/schemas/Lizard'
  discriminator:
    propertyName: petType`
		var got map[string]*Schema
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := map[string]*Schema{
			"MyResponseType": {
				oneOf: []*Schema{
					{reference: "#/components/schemas/Cat"},
					{reference: "#/components/schemas/Dog"},
					{reference: "#/components/schemas/Lizard"},
				},
				discriminator: &Discriminator{
					propertyName: "petType",
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("mapping", func(t *testing.T) {
		yml := `MyResponseType:
  oneOf:
  - $ref: '#/components/schemas/Cat'
  - $ref: '#/components/schemas/Dog'
  - $ref: '#/components/schemas/Lizard'
  - $ref: 'https://gigantic-server.com/schemas/Monster/schema.json'
  discriminator:
    propertyName: petType
    mapping:
      dog: '#/components/schemas/Dog'
      monster: 'https://gigantic-server.com/schemas/Monster/schema.json'`
		var got map[string]*Schema
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := map[string]*Schema{
			"MyResponseType": {
				oneOf: []*Schema{
					{reference: "#/components/schemas/Cat"},
					{reference: "#/components/schemas/Dog"},
					{reference: "#/components/schemas/Lizard"},
					{reference: "https://gigantic-server.com/schemas/Monster/schema.json"},
				},
				discriminator: &Discriminator{
					propertyName: "petType",
					mapping: map[string]string{
						"dog":     "#/components/schemas/Dog",
						"monster": "https://gigantic-server.com/schemas/Monster/schema.json",
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
}

func TestDiscriminatorUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Discriminator
	}{}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Discriminator
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestDiscriminatorUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `mapping: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Discriminator{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestXMLExampleUnmarshalYAML(t *testing.T) {
	t.Run("basic string", func(t *testing.T) {
		yml := `animals:
  type: string
  xml:
    name: animal`
		var got map[string]*Schema
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := map[string]*Schema{
			"animals": {
				type_: "string",
				xml: &XML{
					name: "animal",
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("attribute prefix namespace", func(t *testing.T) {
		yml := `Person:
  type: object
  properties:
    id:
      type: integer
      format: int32
      xml:
        attribute: true
    name:
      type: string
      xml:
        namespace: http://example.com/schema/sample
        prefix: sample`
		var got map[string]*Schema
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := map[string]*Schema{
			"Person": {
				type_: "object",
				properties: map[string]*Schema{
					"id": {
						type_:  "integer",
						format: "int32",
						xml:    &XML{attribute: true},
					},
					"name": {
						type_: "string",
						xml: &XML{
							namespace: "http://example.com/schema/sample",
							prefix:    "sample",
						},
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("wrapped array", func(t *testing.T) {
		yml := `animals:
  type: array
  items:
    type: string
  xml:
    wrapped: true`
		var got map[string]*Schema
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := map[string]*Schema{
			"animals": {
				type_: "array",
				items: &Schema{type_: "string"},
				xml:   &XML{wrapped: true},
			},
		}
		assertEqual(t, got, want)
	})
}

func TestXMLUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want XML
	}{
		{
			yml: `x-foo: bar`,
			want: XML{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got XML
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestXMLUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &XML{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestSecuritySchemeExampleUnmarshalYAML(t *testing.T) {
	t.Run("basic auth", func(t *testing.T) {
		yml := `type: http
scheme: basic`
		var got SecurityScheme
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := SecurityScheme{
			type_:  "http",
			scheme: "basic",
		}
		assertEqual(t, got, want)
	})
	t.Run("api key", func(t *testing.T) {
		yml := `type: apiKey
name: api_key
in: header`
		var got SecurityScheme
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := SecurityScheme{
			type_: "apiKey",
			name:  "api_key",
			in:    "header",
		}
		assertEqual(t, got, want)
	})
	t.Run("JWT bearer", func(t *testing.T) {
		yml := `type: http
scheme: bearer
bearerFormat: JWT`
		var got SecurityScheme
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := SecurityScheme{
			type_:        "http",
			scheme:       "bearer",
			bearerFormat: "JWT",
		}
		assertEqual(t, got, want)
	})
	t.Run("implicit oauth2", func(t *testing.T) {
		yml := `type: oauth2
flows:
  implicit:
    authorizationUrl: https://example.com/api/oauth/dialog
    scopes:
      write:pets: modify pets in your account
      read:pets: read your pets`
		var got SecurityScheme
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := SecurityScheme{
			type_: "oauth2",
			flows: &OAuthFlows{
				implicit: &OAuthFlow{
					authorizationURL: "https://example.com/api/oauth/dialog",
					scopes: map[string]string{
						"write:pets": "modify pets in your account",
						"read:pets":  "read your pets",
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
}

func TestSecuritySchemeUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want SecurityScheme
	}{
		{
			yml: `type: apiKey`,
			want: SecurityScheme{
				type_: "apiKey",
			},
		},
		{
			yml: `type: http`,
			want: SecurityScheme{
				type_: "http",
			},
		},
		{
			yml: `type: oauth2`,
			want: SecurityScheme{
				type_: "oauth2",
			},
		},
		{
			yml: `type: openIdConnect`,
			want: SecurityScheme{
				type_: "openIdConnect",
			},
		},
		{
			yml: `description: foobar`,
			want: SecurityScheme{
				description: "foobar",
			},
		},
		{
			yml: `in: query`,
			want: SecurityScheme{
				in: "query",
			},
		},
		{
			yml: `in: header`,
			want: SecurityScheme{
				in: "header",
			},
		},
		{
			yml: `in: cookie`,
			want: SecurityScheme{
				in: "cookie",
			},
		},
		{
			yml: `openIdConnectUrl: https://example.com`,
			want: SecurityScheme{
				openIDConnectURL: "https://example.com",
			},
		},
		{
			yml: `$ref: '#/components/securitySchemes/Foo'`,
			want: SecurityScheme{
				reference: "#/components/securitySchemes/Foo",
			},
		},
		{
			yml: `x-foo: bar`,
			want: SecurityScheme{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got SecurityScheme
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestSecuritySchemeUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `type: foo`,
			want: errors.New(`"type" field must be one of ["apiKey", "http", "oauth2", "openIdConnect"]`),
		},
		{
			yml:  `in: "foo"`,
			want: errors.New(`"in" field must be one of ["query", "header", "cookie"]`),
		},
		{
			yml:  `flows: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `openIdConnectUrl: foo`,
			want: errors.New(`parse "foo": invalid URI for request`),
		},

		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &SecurityScheme{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestOAuthFlowsUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want OAuthFlows
	}{
		{
			yml: `password: {}`,
			want: OAuthFlows{
				password: &OAuthFlow{},
			},
		},
		{
			yml: `clientCredentials: {}`,
			want: OAuthFlows{
				clientCredentials: &OAuthFlow{},
			},
		},
		{
			yml: `x-foo: bar`,
			want: OAuthFlows{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got OAuthFlows
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestOAuthFlowsUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml: `implicit: foobar`,
			// implicit expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `password: foobar`,
			// password expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `clientCredentials: foobar`,
			// clientCredentials expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml: `authorizationCode: foobar`,
			// authorizationCode expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &OAuthFlows{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestOAuthFlowExampleUnmarshalYAML(t *testing.T) {
	yml := `type: oauth2
flows:
  implicit:
    authorizationUrl: https://example.com/api/oauth/dialog
    scopes:
      write:pets: modify pets in your account
      read:pets: read your pets
  authorizationCode:
    authorizationUrl: https://example.com/api/oauth/dialog
    tokenUrl: https://example.com/api/oauth/token
    scopes:
      write:pets: modify pets in your account
      read:pets: read your pets `
	var got SecurityScheme
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}
	want := SecurityScheme{
		type_: "oauth2",
		flows: &OAuthFlows{
			implicit: &OAuthFlow{
				authorizationURL: "https://example.com/api/oauth/dialog",
				scopes: map[string]string{
					"write:pets": "modify pets in your account",
					"read:pets":  "read your pets",
				},
			},
			authorizationCode: &OAuthFlow{
				authorizationURL: "https://example.com/api/oauth/dialog",
				tokenURL:         "https://example.com/api/oauth/token",
				scopes: map[string]string{
					"write:pets": "modify pets in your account",
					"read:pets":  "read your pets",
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func TestOAuthFlowUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want OAuthFlow
	}{
		{
			yml: `refreshUrl: https://example.com`,
			want: OAuthFlow{
				refreshURL: "https://example.com",
			},
		},
		{
			yml: `x-foo: bar`,
			want: OAuthFlow{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got OAuthFlow
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestOAuthFlowUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml: `authorizationUrl: foobar`,
			// authorizationUrl expects URI
			want: errors.New(`parse "foobar": invalid URI for request`),
		},
		{
			yml: `tokenUrl: foobar`,
			// tokenUrl expects URI
			want: errors.New(`parse "foobar": invalid URI for request`),
		},
		{
			yml: `refreshUrl: foobar`,
			// refreshUrl expects URI
			want: errors.New(`parse "foobar": invalid URI for request`),
		},
		{
			yml: `scopes: foobar`,
			// scopes expects an object
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &OAuthFlow{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestSecurityRequirementExampleUnmarshalYAML(t *testing.T) {
	t.Run("non-oauth2", func(t *testing.T) {
		yml := `api_key: []`
		var got SecurityRequirement
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := SecurityRequirement{
			securityRequirement: map[string][]string{
				"api_key": {},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("oauth2", func(t *testing.T) {
		yml := `petstore_auth:
- write:pets
- read:pets`
		var got SecurityRequirement
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := SecurityRequirement{
			securityRequirement: map[string][]string{
				"petstore_auth": {"write:pets", "read:pets"},
			},
		}
		assertEqual(t, got, want)
	})
}

func TestSecurityRequirementUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want SecurityRequirement
	}{}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got SecurityRequirement
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestSecurityRequirementUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `foo: bar`,
			want: errors.New("String node doesn't ArrayNode"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &SecurityRequirement{})
			assertSameError(t, got, tt.want)
		})
	}
}

func TestIsOneOf(t *testing.T) {
	tests := []struct {
		s    string
		list []string
		want bool
	}{
		{
			s:    "",
			list: []string{},
			want: false,
		},
		{
			s:    "a",
			list: []string{"a", "b"},
			want: true,
		},
		{
			s:    "c",
			list: []string{"a", "b"},
			want: false,
		},
		{
			s:    "a",
			list: nil,
			want: false,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := isOneOf(tt.s, tt.list)
			if got != tt.want {
				t.Errorf("unexpected: %t != %t", got, tt.want)
				return
			}
		})
	}
}
