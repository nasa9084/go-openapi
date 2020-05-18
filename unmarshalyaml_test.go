//nolint:dupl,goconst,funlen
package openapi

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected openapi:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var target OpenAPI
			got := yaml.Unmarshal([]byte(tt.yml), &target)
			if tt.want != nil {
				if got == nil {
					t.Errorf("error `%s` is expected but not", tt.want.Error())
					return
				}
				if got.Error() != tt.want.Error() {
					t.Errorf("unexpected error:\n  got:  %s\n  want: %s", got.Error(), tt.want.Error())
					return
				}
			} else {
				if got != nil {
					t.Errorf("error: %s", got)
					return
				}
			}
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

	var info Info
	if err := yaml.Unmarshal([]byte(yml), &info); err != nil {
		t.Fatal(err)
	}

	if info.title != "Sample Pet Store App" {
		t.Errorf("unexpected info.title: %s", info.title)
		return
	}
	if info.description != "This is a sample server for a pet store." {
		t.Errorf("unexpected info.description: %s", info.description)
		return
	}
	if info.termsOfService != "http://example.com/terms/" {
		t.Errorf("unexpected info.termsOfService: %s", info.termsOfService)
		return
	}
	if info.contact.name != "API Support" {
		t.Errorf("unexpected info.contact.name: %s", info.contact.name)
		return
	}
	if info.contact.url != "http://www.example.com/support" {
		t.Errorf("unexpected info.contact.url: %s", info.contact.url)
		return
	}
	if info.contact.email != "support@example.com" {
		t.Errorf("unexpected info.contact.email: %s", info.contact.email)
		return
	}
	if info.license.name != "Apache 2.0" {
		t.Errorf("unexpected info.license.name: %s", info.license.name)
		return
	}
	if info.license.url != "https://www.apache.org/licenses/LICENSE-2.0.html" {
		t.Errorf("unexpected info.license.url: %s", info.license.url)
		return
	}
	if info.version != "1.0.1" {
		t.Errorf("unexpected info.version: %s", info.version)
		return
	}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected info:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var info Info
			got := yaml.Unmarshal([]byte(tt.yml), &info)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
		})
	}
}

func TestContactExampleUnmarshalYAML(t *testing.T) {
	yml := `name: API Support
url: http://www.example.com/support
email: support@example.com`

	var contact Contact
	if err := yaml.Unmarshal([]byte(yml), &contact); err != nil {
		t.Fatal(err)
	}

	if contact.name != "API Support" {
		t.Errorf("unexpected contact.name: %s", contact.name)
		return
	}
	if contact.url != "http://www.example.com/support" {
		t.Errorf("unexpected contact.url: %s", contact.url)
		return
	}
	if contact.email != "support@example.com" {
		t.Errorf("unexpected contact.email: %s", contact.email)
		return
	}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected contact:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var contact Contact
			got := yaml.Unmarshal([]byte(tt.yml), &contact)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
		})
	}
}

func TestLicenseExampleUnmarshalYAML(t *testing.T) {
	yml := `name: Apache 2.0
url: https://www.apache.org/licenses/LICENSE-2.0.html`

	var license License
	if err := yaml.Unmarshal([]byte(yml), &license); err != nil {
		t.Fatal(err)
	}

	if license.name != "Apache 2.0" {
		t.Errorf("unexpected license.name: %s", license.name)
		return
	}
	if license.url != "https://www.apache.org/licenses/LICENSE-2.0.html" {
		t.Errorf("unexpected license.url: %s", license.url)
		return
	}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected license:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var license License
			got := yaml.Unmarshal([]byte(tt.yml), &license)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
		})
	}
}

func TestServerExampleUnmarshalYAML(t *testing.T) {
	t.Run("single server", func(t *testing.T) {
		yml := `url: https://development.gigantic-server.com/v1
description: Development server`

		var server Server
		if err := yaml.Unmarshal([]byte(yml), &server); err != nil {
			t.Fatal(err)
		}

		if server.url != "https://development.gigantic-server.com/v1" {
			t.Errorf("unexpected server.url: %s", server.url)
			return
		}
		if server.description != "Development server" {
			t.Errorf("unexpected server.description: %s", server.description)
			return
		}
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
		servers := target.Servers
		t.Run("0", func(t *testing.T) {
			server := servers[0]
			if server.url != "https://development.gigantic-server.com/v1" {
				t.Errorf("unexpected server.url: %s", server.url)
				return
			}
			if server.description != "Development server" {
				t.Errorf("unexpected server.description: %s", server.description)
				return
			}
		})
		t.Run("1", func(t *testing.T) {
			server := servers[1]
			if server.url != "https://staging.gigantic-server.com/v1" {
				t.Errorf("unexpected server.url: %s", server.url)
				return
			}
			if server.description != "Staging server" {
				t.Errorf("unexpected server.description: %s", server.description)
				return
			}
		})
		t.Run("2", func(t *testing.T) {
			server := servers[2]
			if server.url != "https://api.gigantic-server.com/v1" {
				t.Errorf("unexpected server.url: %s", server.url)
				return
			}
			if server.description != "Production server" {
				t.Errorf("unexpected server.description: %s", server.description)
				return
			}
		})
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

		server := target.Servers[0]
		if server.url != "https://{username}.gigantic-server.com:{port}/{basePath}" {
			t.Errorf("unexpected server url: %s", server.url)
			return
		}
		if server.description != "The production API server" {
			t.Errorf("unexpected server.descripion: %s", server.description)
			return
		}
		if server.variables["username"].default_ != "demo" {
			t.Errorf("unexpected server.variables.username.default: %s", server.variables["username"].default_)
			return
		}
		if server.variables["username"].description != `this value is assigned by the service provider, in this example "gigantic-server.com"` {
			t.Errorf("unexpected server.variables.username.description: %s", server.variables["username"].description)
			return
		}
		if len(server.variables["port"].enum) != 2 {
			t.Errorf("unexpected length of server.variables.port.enum: %d", len(server.variables["port"].enum))
			return
		}
		if !reflect.DeepEqual(server.variables["port"].enum, []string{"8443", "443"}) {
			t.Errorf("unexpected server.variables.port.enum: %#v", server.variables["port"].enum)
			return
		}
		if server.variables["port"].default_ != "8443" {
			t.Errorf("unexpected server.variables.port.default: %s", server.variables["port"].default_)
			return
		}
		if server.variables["basePath"].default_ != "v2" {
			t.Errorf("unexpected server.variables.basepath.default: %s", server.variables["basePath"].default_)
			return
		}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected server:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var server Server
			got := yaml.Unmarshal([]byte(tt.yml), &server)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected serverVariable:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var serverVariable ServerVariable
			got := yaml.Unmarshal([]byte(tt.yml), &serverVariable)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
	components := target.Components
	schemas := components.schemas
	t.Run("schemas.GeneralError", func(t *testing.T) {
		generalError, ok := schemas["GeneralError"]
		if !ok {
			t.Error("components.schemas.GeneralError is not found")
			return
		}
		if generalError.type_ != "object" {
			t.Errorf("unexpected components.schema.type: %s", generalError.type_)
			return
		}
		code, ok := generalError.properties["code"]
		if !ok {
			t.Error("components.schemas.GeneralError.properties.code is not found")
			return
		}
		if code.type_ != "integer" {
			t.Errorf("unexpected components.schemas.GeneralError.properties.code.type: %s", code.type_)
			return
		}
		if code.format != "int32" {
			t.Errorf("unexpected components.schemas.GeneralError.properties.code.format: %s", code.format)
			return
		}
		message, ok := generalError.properties["message"]
		if !ok {
			t.Error("components.schemas.GeneralError.properties.message is not found")
			return
		}
		if message.type_ != "string" {
			t.Errorf("unexpected components.schemas.GeneralError.properties.message.type: %s", message.type_)
			return
		}
	})
	t.Run("schemas.Category", func(t *testing.T) {
		category, ok := schemas["Category"]
		if !ok {
			t.Error("components.schemas.Category is not found")
			return
		}
		if category.type_ != "object" {
			t.Errorf("unexpected components.schema.Category.type: %s", category.type_)
			return
		}
		id, ok := category.properties["id"]
		if !ok {
			t.Error("components.schemas.Category.properties.id is not found")
			return
		}
		if id.type_ != "integer" {
			t.Errorf("unexpected components.schemas.Category.properties.id.type: %s", id.type_)
			return
		}
		if id.format != "int64" {
			t.Errorf("unexpected components.schemas.Category.properties.id.format: %s", id.format)
			return
		}
		name, ok := category.properties["name"]
		if !ok {
			t.Error("components.schemas.Category.properties.name is not found")
			return
		}
		if name.type_ != "string" {
			t.Errorf("unexpected components.schemas.Category.properties.name.type: %s", name.type_)
			return
		}
	})
	parameters := components.parameters
	t.Run("parameters.skipParam", func(t *testing.T) {
		skipParam, ok := parameters["skipParam"]
		if !ok {
			t.Error("components.parameters.skipParam is not found")
			return
		}
		if skipParam.name != "skip" {
			t.Errorf("unexpected components.parameters.skipParam.name: %s", skipParam.name)
			return
		}
		if skipParam.in != "query" {
			t.Errorf("unexpected components.parameters.skipParam.in: %s", skipParam.in)
			return
		}
		if skipParam.description != "number of items to skip" {
			t.Errorf("unexpected components.parameters.skipParam.description: %s", skipParam.description)
			return
		}
		if skipParam.required != true {
			t.Errorf("unexpected components.parameters.skipParam.required: %t", skipParam.required)
			return
		}
		schema := skipParam.schema
		if schema.type_ != "integer" {
			t.Errorf("unexpected components.parameters.skipParam.schema.type: %s", schema.type_)
			return
		}
		if schema.format != "int32" {
			t.Errorf("unexpected components.parameters.skipParam.schema.format: %s", schema.format)
			return
		}
	})
	t.Run("parameters.limitParam", func(t *testing.T) {
		limitParam, ok := parameters["limitParam"]
		if !ok {
			t.Error("components.parameters.limitParam is not found")
			return
		}
		if limitParam.name != "limit" {
			t.Errorf("unexpected components.parameters.limitParam.name: %s", limitParam.name)
			return
		}
		if limitParam.in != "query" {
			t.Errorf("unexpected components.parameters.limitParam.in: %s", limitParam.in)
			return
		}
		if limitParam.description != "max records to return" {
			t.Errorf("unexpected components.parameters.limitParam.description: %s", limitParam.description)
			return
		}
		if limitParam.required != true {
			t.Errorf("unexpected components.parameters.limitParam.required: %t", limitParam.required)
			return
		}
		schema := limitParam.schema
		if schema.type_ != "integer" {
			t.Errorf("unexpected components.parameters.limitParam.schema.type: %s", schema.type_)
			return
		}
		if schema.format != "int32" {
			t.Errorf("unexpected components.parameters.limitParam.schema.format: %s", schema.format)
			return
		}
	})
	responses := components.responses
	t.Run("responses.NotFound", func(t *testing.T) {
		notFound, ok := responses["NotFound"]
		if !ok {
			t.Error("components.responses.NotFound is not found")
			return
		}
		if notFound.description != "Entity not found." {
			t.Errorf("unexpected components.responses.NotFound.description: %s", notFound.description)
			return
		}
	})
	t.Run("responses.IllegalInput", func(t *testing.T) {
		illegalInput, ok := responses["IllegalInput"]
		if !ok {
			t.Error("components.responses.IllegalInput is not found")
			return
		}
		if illegalInput.description != "Illegal input for operation." {
			t.Errorf("unexpected components.responses.IllegalInput.description: %s", illegalInput.description)
			return
		}
	})
	t.Run("responses.GeneralError", func(t *testing.T) {
		generalError, ok := responses["GeneralError"]
		if !ok {
			t.Error("components.responses.GeneralError is not found")
			return
		}
		if generalError.description != "General Error" {
			t.Errorf("unexpected components.responses.GeneralError.description: %s", generalError.description)
			return
		}
		mediaType, ok := generalError.content["application/json"]
		if !ok {
			t.Error("components.responses.GeneralError.content.application/json is not found")
			return
		}
		if mediaType.schema.reference != "#/components/schemas/GeneralError" {
			t.Errorf("unexpected components.responses.GeneralError.content.application/json.schema.$ref")
			return
		}
	})
	securitySchemes := components.securitySchemes
	t.Run("securitySchemes.api_key", func(t *testing.T) {
		apiKey, ok := securitySchemes["api_key"]
		if !ok {
			t.Error("components.securitySchemes.api_key is not found")
			return
		}
		if apiKey.type_ != "apiKey" {
			t.Errorf("unexpected components.securitySchemes.api_key.type: %s", apiKey.type_)
			return
		}
		if apiKey.name != "api_key" {
			t.Errorf("unexpected components.securitySchemes.api_key.name: %s", apiKey.name)
			return
		}
		if apiKey.in != "header" {
			t.Errorf("unexpected components.securitySchemes.api_key.in: %s", apiKey.in)
			return
		}
	})
	t.Run("securitySchemes.petstore_auth", func(t *testing.T) {
		petstoreAuth, ok := securitySchemes["petstore_auth"]
		if !ok {
			t.Error("components.securitySchemes.petstore_auth is not found")
			return
		}
		if petstoreAuth.type_ != "oauth2" {
			t.Errorf("unexpected components.securitySchemes.petstore_auth.type: %s", petstoreAuth.type_)
			return
		}
		if petstoreAuth.flows.implicit.authorizationURL != "http://example.org/api/oauth/dialog" {
			t.Errorf("unexpected components.securitySchemes.petstore_auth.flows.implicit.authorizationURL: %s", petstoreAuth.flows.implicit.authorizationURL)
			return
		}
		scopes := petstoreAuth.flows.implicit.scopes
		write, ok := scopes["write:pets"]
		if !ok {
			t.Error("components.securitySchemes.petstore_auth.flows.implicit.scopes.write:pets is not found")
			return
		}
		if write != "modify pets in your account" {
			t.Errorf("unexpected components.securitySchemes.petstore_auth.flows.implicit.scopes.write:pets: %s", write)
		}
		read, ok := scopes["read:pets"]
		if !ok {
			t.Error("components.securitySchemes.petstore_auth.flows.implicit.scopes.read:pets is not found")
			return
		}
		if read != "read your pets" {
			t.Errorf("unexpected components.securitySchemes.petstore_auth.flows.implicit.scopes.read:pets: %s", write)
		}
	})
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected components:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var components Components
			got := yaml.Unmarshal([]byte(tt.yml), &components)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
	var paths Paths
	if err := yaml.Unmarshal([]byte(yml), &paths); err != nil {
		t.Fatal(err)
	}
	if _, ok := paths.paths["/pets"]; !ok {
		t.Error("paths./pets is not found")
		return
	}
	op := paths.paths["/pets"].get
	if op.description != "Returns all pets from the system that the user has access to" {
		t.Errorf("unexpected paths./pets.get.description: %s", op.description)
		return
	}
	response, ok := op.responses.responses["200"]
	if !ok {
		t.Error("paths./pets.get.responses.200 is not found")
		return
	}
	if _, ok := response.content["application/json"]; !ok {
		t.Error("paths./pets.get.responses.200.content.application/json is not found")
		return
	}
	schema := response.content["application/json"].schema
	if schema.type_ != "array" {
		t.Errorf("unexpected paths./pets.get.responses.200.content.application/json.schema.type: %s", schema.type_)
		return
	}
	if schema.items.reference != "#/components/schemas/pet" {
		t.Errorf("unexpected paths./pets.get.responses.200.content.application/json.schema.items.$ref: %s", schema.reference)
		return
	}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected paths:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var paths Paths
			got := yaml.Unmarshal([]byte(tt.yml), &paths)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
	var pathItem PathItem
	if err := yaml.Unmarshal([]byte(yml), &pathItem); err != nil {
		t.Fatal(err)
	}
	t.Run("get", func(t *testing.T) {
		operation := pathItem.get
		if operation.description != "Returns pets based on ID" {
			t.Errorf("unexpected pathItem.get.description: %s", operation.description)
			return
		}
		if operation.summary != "Find pets by ID" {
			t.Errorf("unexpected pathItem.get.summary: %s", operation.summary)
			return
		}
		if operation.operationID != "getPetsById" {
			t.Errorf("unexpected pathItem.get.operationId: %s", operation.operationID)
			return
		}
		t.Run("200", func(t *testing.T) {
			response, ok := pathItem.get.responses.responses["200"]
			if !ok {
				t.Error("pathItem.get.responses.200 is not found")
				return
			}
			if response.description != "pet response" {
				t.Errorf("unexpected pathItem.get.responses.200.description: %s", response.description)
				return
			}
			if _, ok := response.content["*/*"]; !ok {
				t.Error("pathItem.get.responses.200.content.*/* is not found")
				return
			}
			schema := response.content["*/*"].schema
			if schema.type_ != "array" {
				t.Errorf("unexpected pathItem.get.responses.200.content.*/*.schema.type: %s", schema.type_)
				return
			}
			if schema.items.reference != "#/components/schemas/Pet" {
				t.Errorf("unexpected pathItem.get.responses.200.content.*/*.schema.items.$ref: %s", schema.items.reference)
				return
			}
		})
		t.Run("default", func(t *testing.T) {
			response, ok := pathItem.get.responses.responses["default"]
			if !ok {
				t.Error("pathItem.get.responses.default is not found")
				return
			}
			if response.description != "error payload" {
				t.Errorf("unexpected pathItem.get.responses.default.description: %s", response.description)
				return
			}
			if _, ok := response.content["text/html"]; !ok {
				t.Error("pathItem.get.responses.default.content.text/html is not found")
				return
			}
			if response.content["text/html"].schema.reference != "#/components/schemas/ErrorModel" {
				t.Errorf("unexpected pathItem.get.responses.default.content.text/html.schema.$ref: %s", response.content["text/html"].schema.reference)
				return
			}
		})
	})
	t.Run("parameters", func(t *testing.T) {
		parameters := pathItem.parameters
		id := parameters[0]
		if id.name != "id" {
			t.Errorf("unexpected pathItem.parameters.0.name: %s", id.name)
			return
		}
		if id.in != "path" {
			t.Errorf("unexpected pathItem.parameters.0.in: %s", id.in)
			return
		}
		if id.description != "ID of pet to use" {
			t.Errorf("unexpected pathItem.parameters.0.description: %s", id.description)
			return
		}
		if id.required != true {
			t.Errorf("unexpected pathItem.parameters.0.required: %t", id.required)
			return
		}
		if id.schema.type_ != "array" {
			t.Errorf("unexpected pathItem.parameters.0.schema.type: %s", id.schema.type_)
			return
		}
		if id.schema.items.type_ != "string" {
			t.Errorf("unexpected pathItem.parameters.0.schema.items.type: %s", id.schema.items.type_)
			return
		}
	})
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected pathItem:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var pathItem PathItem
			got := yaml.Unmarshal([]byte(tt.yml), &pathItem)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
	var operation Operation
	if err := yaml.Unmarshal([]byte(yml), &operation); err != nil {
		t.Fatal(err)
	}
	if operation.tags[0] != "pet" {
		t.Errorf("unexpected operation.tags.0: %s", operation.tags[0])
		return
	}
	if operation.summary != "Updates a pet in the store with form data" {
		t.Errorf("unexpected opration.summary: %s", operation.summary)
		return
	}
	if operation.operationID != "updatePetWithForm" {
		t.Errorf("unexpected operation.operationId: %s", operation.operationID)
		return
	}
	parameter := operation.parameters[0]
	if parameter.name != "petId" {
		t.Errorf("unexpected operation.parameters.0.name: %s", parameter.name)
		return
	}
	if parameter.in != "path" {
		t.Errorf("unexpected operation.parameters.0.in: %s", parameter.in)
		return
	}
	if parameter.description != "ID of pet that needs to be updated" {
		t.Errorf("unexpected operation.parameters.0.description: %s", parameter.description)
		return
	}
	if parameter.required != true {
		t.Errorf("unexpected operation.parameters.0.required: %t", parameter.required)
		return
	}
	if parameter.schema.type_ != "string" {
		t.Errorf("unexpected operation.parameters.0.schema.type: %s", parameter.schema.type_)
		return
	}
	if _, ok := operation.requestBody.content["application/x-www-form-urlencoded"]; !ok {
		t.Error("operation.requestBody.content.application/x-www-form-urlencoded is not found")
		return
	}
	schema := operation.requestBody.content["application/x-www-form-urlencoded"].schema
	name, ok := schema.properties["name"]
	if !ok {
		t.Error("operation.requestBody.content.application/x-www-form-urlencoded.schema.properties.name is not found")
		return
	}
	if name.description != "Updated name of the pet" {
		t.Errorf("unexpected operation.requestBody.content.application/x-www-form-urlencoded.schema.properties.name.description: %s", name.description)
		return
	}
	if name.type_ != "string" {
		t.Errorf("unexpected operation.requestBody.content.application/x-www-form-urlencoded.schema.properties.name.type: %s", name.type_)
		return
	}
	status, ok := schema.properties["status"]
	if !ok {
		t.Error("operation.requestBody.content.application/x-www-form-urlencoded.schema.properties.status is not found")
		return
	}
	if status.description != "Updated status of the pet" {
		t.Errorf("unexpected operation.requestBody.content.application/x-www-form-urlencoded.schema.properties.status.description: %s", status.description)
		return
	}
	if status.type_ != "string" {
		t.Errorf("unexpected operation.requestBody.content.application/x-www-form-urlencoded.schema.properties.status.type: %s", status.type_)
		return
	}
	if schema.required[0] != "status" {
		t.Errorf("unexpected operation.requestBody.content.application/x-www-form-urlencoded.schema.required.0: %s", schema.required[0])
		return
	}
	if _, ok := operation.responses.responses["200"]; !ok {
		t.Error("operation.responses.200 is not found")
		return
	}
	if operation.responses.responses["200"].description != "Pet updated." {
		t.Errorf("unexpected operation.responses.200.description: %s", operation.responses.responses["200"].description)
		return
	}
	if _, ok := operation.responses.responses["200"].content["application/json"]; !ok {
		t.Error("operation.responses.200.content.application/json is not found")
	}
	if _, ok := operation.responses.responses["200"].content["application/xml"]; !ok {
		t.Error("operation.responses.200.content.application/xml is not found")
		return
	}
	if _, ok := operation.responses.responses["405"]; !ok {
		t.Error("operation.responses.405 is not found")
		return
	}
	if operation.responses.responses["405"].description != "Method Not Allowed" {
		t.Errorf("unexpected operation.responses.405.description: %s", operation.responses.responses["405"].description)
		return
	}
	if _, ok := operation.responses.responses["405"].content["application/json"]; !ok {
		t.Error("operation.responses.405.content.application/json is not found")
	}
	if _, ok := operation.responses.responses["405"].content["application/xml"]; !ok {
		t.Error("operation.responses.405.content.application/xml is not found")
		return
	}
	securityRequirement, ok := operation.security[0].securityRequirement["petstore_auth"]
	if !ok {
		t.Error("operation.security.0.petstore_auth is not found")
		return
	}
	if securityRequirement[0] != "write:pets" {
		t.Errorf("unexpected operation.security.0.petstore_auth.0: %s", securityRequirement[0])
		return
	}
	if securityRequirement[1] != "read:pets" {
		t.Errorf("unexpected operation.security.0.petstore_auth.1: %s", securityRequirement[1])
		return
	}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected operation:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var operation Operation
			got := yaml.Unmarshal([]byte(tt.yml), &operation)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
		})
	}
}

func TestExternalDocumentationExampleUnmarshalYAML(t *testing.T) {
	yml := `description: Find more info here
url: https://example.com`
	var externalDocumentation ExternalDocumentation
	if err := yaml.Unmarshal([]byte(yml), &externalDocumentation); err != nil {
		t.Fatal(err)
	}
	if externalDocumentation.description != "Find more info here" {
		t.Errorf("unexpected externalDocumentation.description: %s", externalDocumentation.description)
		return
	}
	if externalDocumentation.url != "https://example.com" {
		t.Errorf("unexpected externalDocumentation.url: %s", externalDocumentation.url)
		return
	}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected externalDocumentation:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var externalDocumentation ExternalDocumentation
			got := yaml.Unmarshal([]byte(tt.yml), &externalDocumentation)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
		var parameter Parameter
		if err := yaml.Unmarshal([]byte(yml), &parameter); err != nil {
			t.Fatal(err)
		}
		if parameter.name != "token" {
			t.Errorf("unexpected paramater.name: %s", parameter.name)
			return
		}
		if parameter.in != "header" {
			t.Errorf("unexpected parameter.in: %s", parameter.in)
			return
		}
		if parameter.description != "token to be passed as a header" {
			t.Errorf("unexpected parameter.description: %s", parameter.description)
			return
		}
		if parameter.required != true {
			t.Errorf("unexpected parameter.required: %t", parameter.required)
			return
		}
		if parameter.schema.type_ != "array" {
			t.Errorf("unexpected parameter.schema.type: %s", parameter.schema.type_)
			return
		}
		if parameter.schema.items.type_ != "integer" {
			t.Errorf("unexpected parameter.schema.items.type: %s", parameter.schema.items.type_)
			return
		}
		if parameter.schema.items.format != "int64" {
			t.Errorf("unexpected parameter.schema.items.format: %s", parameter.schema.items.format)
			return
		}
		if parameter.style != "simple" {
			t.Errorf("unexpected paarameter.style: %s", parameter.style)
			return
		}
	})
	t.Run("path parameter", func(t *testing.T) {
		yml := `name: username
in: path
description: username to fetch
required: true
schema:
  type: string`
		var parameter Parameter
		if err := yaml.Unmarshal([]byte(yml), &parameter); err != nil {
			t.Fatal(err)
		}
		if parameter.name != "username" {
			t.Errorf("unexpected parameter.name: %s", parameter.name)
			return
		}
		if parameter.in != "path" {
			t.Errorf("unexpected parameter.in: %s", parameter.in)
			return
		}
		if parameter.description != "username to fetch" {
			t.Errorf("unexpected parameter.description: %s", parameter.description)
			return
		}
		if parameter.required != true {
			t.Errorf("unexpected parameter.required: %t", parameter.required)
			return
		}
		if parameter.schema.type_ != "string" {
			t.Errorf("unexpected parameter.schema.type: %s", parameter.schema.type_)
			return
		}
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
		var parameter Parameter
		if err := yaml.Unmarshal([]byte(yml), &parameter); err != nil {
			t.Fatal(err)
		}
		if parameter.name != "id" {
			t.Errorf("unexpected parameter.name: %s", parameter.name)
			return
		}
		if parameter.in != "query" {
			t.Errorf("unexpected parameter.in: %s", parameter.in)
			return
		}
		if parameter.description != "ID of the object to fetch" {
			t.Errorf("unexpected parameter.description: %s", parameter.description)
			return
		}
		if parameter.required != false {
			t.Errorf("unexpected parameter.required: %t", parameter.required)
			return
		}
		if parameter.schema.type_ != "array" {
			t.Errorf("unexpected parameter.schema.type: %s", parameter.schema.type_)
			return
		}
		if parameter.schema.items.type_ != "string" {
			t.Errorf("unexpected parameter.schema.items.type: %s", parameter.schema.items.type_)
			return
		}
		if parameter.style != "form" {
			t.Errorf("unexpected parameter.style: %s", parameter.style)
			return
		}
		if parameter.explode != true {
			t.Errorf("unexpected parameter.explode: %t", parameter.explode)
			return
		}
	})
	t.Run("free form", func(t *testing.T) {
		yml := `in: query
name: freeForm
schema:
  type: object
  additionalProperties:
    type: integer
style: form`
		var parameter Parameter
		if err := yaml.Unmarshal([]byte(yml), &parameter); err != nil {
			t.Fatal(err)
		}
		if parameter.in != "query" {
			t.Errorf("unexpected parameter.in: %s", parameter.in)
			return
		}
		if parameter.name != "freeForm" {
			t.Errorf("unexpected parameter.name: %s", parameter.name)
			return
		}
		if parameter.schema.type_ != "object" {
			t.Errorf("unexpected parameter.schema.type_: %s", parameter.schema.type_)
			return
		}
		if parameter.schema.additionalProperties.type_ != "integer" {
			t.Errorf("unexpected parameter.schema.additionalProperties.type_: %s", parameter.schema.additionalProperties.type_)
			return
		}
		if parameter.style != "form" {
			t.Errorf("unexpected parameter.style: %s", parameter.style)
			return
		}
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
		var parameter Parameter
		if err := yaml.Unmarshal([]byte(yml), &parameter); err != nil {
			t.Fatal(err)
		}
		if parameter.in != "query" {
			t.Errorf("unexpected parameter.in: %s", parameter.in)
			return
		}
		if parameter.name != "coordinates" {
			t.Errorf("unexpected parameter.name: %s", parameter.name)
			return
		}
		if parameter.content["application/json"].schema.type_ != "object" {
			t.Errorf("unexpected parameter.content.application/json.schema.type_: %s", parameter.content["application/json"].schema.type_)
			return
		}
		if parameter.content["application/json"].schema.required[0] != "lat" {
			t.Errorf("unexpected parameter.content.application/json.schema.required.0: %s", parameter.content["application/json"].schema.required[0])
			return
		}
		if parameter.content["application/json"].schema.required[1] != "long" {
			t.Errorf("unexpected parameter.content.application/json..schema.required.1: %s", parameter.content["application/json"].schema.required[1])
			return
		}
		if parameter.content["application/json"].schema.properties["lat"].type_ != "number" {
			t.Errorf("unexpected parameter.content.application/json.schema.properties.lat.type: %s", parameter.content["application/json"].schema.properties["lat"].type_)
			return
		}
		if parameter.content["application/json"].schema.properties["long"].type_ != "number" {
			t.Errorf("unexpected parameter.content.application/json.schema.properties.long.type: %s", parameter.content["application/json"].schema.properties["long"].type_)
			return
		}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected parameter:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var parameter Parameter
			got := yaml.Unmarshal([]byte(tt.yml), &parameter)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
		var requestBody RequestBody
		if err := yaml.Unmarshal([]byte(yml), &requestBody); err != nil {
			t.Fatal(err)
		}
		if requestBody.description != "user to add to the system" {
			t.Errorf("unexpected requestBody.description: %s", requestBody.description)
			return
		}
		t.Run("application/json", func(t *testing.T) {
			mediaType, ok := requestBody.content["application/json"]
			if !ok {
				t.Error("requestBody.content.application/json is not found")
				return
			}
			if mediaType.schema.reference != "#/components/schemas/User" {
				t.Errorf("unexpected requestBody.content.application/json.schema.$ref: %s", mediaType.schema.reference)
				return
			}
			example, ok := mediaType.examples["user"]
			if !ok {
				t.Error("requestBody.content.application/json.examples.user is not found")
				return
			}
			if example.summary != "User Example" {
				t.Errorf("unexpected requestBody.content.application/json.examples.user.summary: %s", example.summary)
				return
			}
			if example.externalValue != "http://foo.bar/examples/user-example.json" {
				t.Errorf("unexpected requestBody.content.application/json.examples.user.externalValue: %s", example.externalValue)
				return
			}
		})
		t.Run("application/xml", func(t *testing.T) {
			mediaType, ok := requestBody.content["application/xml"]
			if !ok {
				t.Error("requestBody.content.application/xml is not found")
				return
			}
			if mediaType.schema.reference != "#/components/schemas/User" {
				t.Errorf("unexpected requestBody.content.application/xml.schema.$ref: %s", mediaType.schema.reference)
				return
			}
			example, ok := mediaType.examples["user"]
			if !ok {
				t.Error("requestBody.content.application/xml.examples.user is not found")
				return
			}
			if example.summary != "User Example in XML" {
				t.Errorf("unexpected requestBody.content.application/xml.examples.user.summary: %s", example.summary)
				return
			}
			if example.externalValue != "http://foo.bar/examples/user-example.xml" {
				t.Errorf("unexpected requestBody.content.application/xml.examples.user.externalValue: %s", example.externalValue)
				return
			}
		})
		t.Run("text/plain", func(t *testing.T) {
			mediaType, ok := requestBody.content["text/plain"]
			if !ok {
				t.Error("requestBody.content.text/plain is not found")
				return
			}
			example, ok := mediaType.examples["user"]
			if !ok {
				t.Error("requestBody.content.text/plain.examples.user is not found")
				return
			}
			if example.summary != "User example in text plain format" {
				t.Errorf("unexpected requestBody.content.text/plain.examples.user.summary: %s", example.summary)
				return
			}
			if example.externalValue != "http://foo.bar/examples/user-example.txt" {
				t.Errorf("unexpected requestBody.content.text/plain.examples.user.externalValue: %s", example.externalValue)
				return
			}
		})
		t.Run("*/*", func(t *testing.T) {
			mediaType, ok := requestBody.content["*/*"]
			if !ok {
				t.Error("requestBody.content.*/* is not found")
				return
			}
			example, ok := mediaType.examples["user"]
			if !ok {
				t.Error("requestBody.content.*/*.examples.user is not found")
				return
			}
			if example.summary != "User example in other format" {
				t.Errorf("unexpected requestBody.content.*/*.examples.user.summary: %s", example.summary)
				return
			}
			if example.externalValue != "http://foo.bar/examples/user-example.whatever" {
				t.Errorf("unexpected requestBody.content.*/*.examples.user.externalValue: %s", example.externalValue)
				return
			}
		})
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
		var requestBody RequestBody
		if err := yaml.Unmarshal([]byte(yml), &requestBody); err != nil {
			t.Fatal(err)
		}
		if requestBody.description != "user to add to the system" {
			t.Errorf("unexpected requestBody.description: %s", requestBody.description)
			return
		}
		if requestBody.required != true {
			t.Errorf("unexpected requestBody.required: %t", requestBody.required)
			return
		}
		mediaType, ok := requestBody.content["text/plain"]
		if !ok {
			t.Error("requestBody.content.text/plain is not found")
			return
		}
		if mediaType.schema.type_ != "array" {
			t.Errorf("unexpected mediaType.schema.type: %s", mediaType.schema.type_)
			return
		}
		if mediaType.schema.items.type_ != "string" {
			t.Errorf("unexpected mediaType.schema.items.type: %s", mediaType.schema.items.type_)
			return
		}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected requestBody:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var requestBody RequestBody
			got := yaml.Unmarshal([]byte(tt.yml), &requestBody)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
	var target map[string]*MediaType
	if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
		t.Fatal(err)
	}
	mediaType, ok := target["application/json"]
	if !ok {
		t.Error("application/json is not found")
		return
	}
	if mediaType.schema.reference != "#/components/schemas/Pet" {
		t.Errorf("unexpected mediaType.schema.$ref: %s", mediaType.schema.reference)
		return
	}
	assert := func(t *testing.T, got map[string]interface{}, want map[string]string) {
		name := strings.Split(t.Name(), "/")[1]
		for key := range want {
			if v, ok := got[key]; !ok {
				t.Errorf("mediaType.examples.%s.value.%s is not found", name, key)
				return
			} else if v != want[key] {
				t.Errorf("unexpected mediaType.examples.%s.value.%s: %s != %s", name, key, v, want[key])
			}
		}
	}
	t.Run("cat", func(t *testing.T) {
		example, ok := mediaType.examples["cat"]
		if !ok {
			t.Error("mediaType.examples.cat is not found")
			return
		}
		value, ok := example.value.(map[string]interface{})
		if !ok {
			t.Errorf("mediaType.examples.cat.value is assumed map[string]interface but %v", reflect.TypeOf(example.value))
			return
		}
		want := map[string]string{
			"name":    "Fluffy",
			"petType": "Cat",
			"color":   "White",
			"gender":  "male",
			"breed":   "Persian",
		}
		assert(t, value, want)
	})
	t.Run("dog", func(t *testing.T) {
		example, ok := mediaType.examples["dog"]
		if !ok {
			t.Error("mediaType.examples.dog is not found")
			return
		}
		value, ok := example.value.(map[string]interface{})
		if !ok {
			t.Errorf("mediaType.examples.dog.value is assumed map[string]interface but %v", reflect.TypeOf(example.value))
			return
		}
		want := map[string]string{
			"name":    "Puma",
			"petType": "Dog",
			"color":   "Black",
			"gender":  "Female",
			"breed":   "Mixed",
		}
		assert(t, value, want)
	})
	t.Run("frog", func(t *testing.T) {
		example, ok := mediaType.examples["frog"]
		if !ok {
			t.Error("mediaType.examples.frog is not found")
			return
		}
		if example.reference != "#/components/examples/frog-example" {
			t.Errorf("unexpected mediaType.examples.frog.$ref: %s", example.reference)
			return
		}
	})
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected mediaType:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var mediaType MediaType
			got := yaml.Unmarshal([]byte(tt.yml), &mediaType)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
	if _, ok := target.RequestBody.content["multipart/mixed"]; !ok {
		t.Error("requestBody.content.multipart/mixed is not found")
		return
	}
	schema := target.RequestBody.content["multipart/mixed"].schema
	if schema.type_ != "object" {
		t.Errorf("unexpected requestBody.content.multipart/mixed.schema.type: %s", schema.type_)
		return
	}
	id, ok := schema.properties["id"]
	if !ok {
		t.Error("requestBody.content.multipart/mixed.schema.properties.id is not found")
		return
	}
	if id.type_ != "string" {
		t.Errorf("unexpected id.type: %s", id.type_)
		return
	}
	if id.format != "uuid" {
		t.Errorf("unexpected id.format: %s", id.format)
		return
	}
	address, ok := schema.properties["address"]
	if !ok {
		t.Error("requestBody.content.multipart/mixed.schema.properties.address is not found")
		return
	}
	if address.type_ != "object" {
		t.Errorf("unexpected id.type: %s", id.type_)
		return
	}
	historyMetadata, ok := schema.properties["historyMetadata"]
	if !ok {
		t.Error("requestBody.content.multipart/mixed.schema.properties.historyMetadata is not found")
		return
	}
	if historyMetadata.description != "metadata in XML format" {
		t.Errorf("unexpected historyMetadata.description: %s", historyMetadata.description)
		return
	}
	if historyMetadata.type_ != "object" {
		t.Errorf("unexpected historyMetadata.type: %s", historyMetadata.type_)
		return
	}
	profileImage, ok := schema.properties["profileImage"]
	if !ok {
		t.Error("requestBody.content.multipart/mixed.schema.properties.profileImage is not found")
		return
	}
	if profileImage.type_ != "string" {
		t.Errorf("unexpected id.type: %s", id.type_)
		return
	}
	if profileImage.format != "binary" {
		t.Errorf("unexpected id.format: %s", id.format)
		return
	}
	t.Run("historyMetadata", func(t *testing.T) {
		encoding := target.RequestBody.content["multipart/mixed"].encoding["historyMetadata"]
		if encoding.contentType != "application/xml; charset=utf-8" {
			t.Errorf("unexpected encoding.contentType: %s", encoding.contentType)
			return
		}
	})
	t.Run("profileImage", func(t *testing.T) {
		encoding := target.RequestBody.content["multipart/mixed"].encoding["profileImage"]
		if encoding.contentType != "image/png, image/jpeg" {
			t.Errorf("unexpected encoding.contentType: %s", encoding.contentType)
			return
		}
		xRateLimitLimit, ok := encoding.headers["X-Rate-Limit-Limit"]
		if !ok {
			t.Error("encoding.headers.X-Rate-Limit-Limit is not found")
			return
		}
		if xRateLimitLimit.description != "The number of allowed requests in the current period" {
			t.Errorf("unexpeceted encoding.headers.X-Rate-Limit-Limit.description: %s", xRateLimitLimit.description)
			return
		}
		if xRateLimitLimit.schema.type_ != "integer" {
			t.Errorf("unexpected encoding.headers.X-Rate-Limit-Limit.schema.type: %s", xRateLimitLimit.schema.type_)
			return
		}
	})
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected encoding:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var encoding Encoding
			got := yaml.Unmarshal([]byte(tt.yml), &encoding)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
	var responses Responses
	if err := yaml.Unmarshal([]byte(yml), &responses); err != nil {
		t.Fatal(err)
	}
	t.Run("200", func(t *testing.T) {
		response, ok := responses.responses["200"]
		if !ok {
			t.Error("responses.200 is not found")
			return
		}
		if response.description != "a pet to be returned" {
			t.Errorf("unexpected responses.200.description: %s", response.description)
			return
		}
		mediaType, ok := response.content["application/json"]
		if !ok {
			t.Error("responses.200.content.application/json is not found")
			return
		}
		if mediaType.schema.reference != "#/components/schemas/Pet" {
			t.Errorf("unexpected responses.200.content.application/json.schema.$ref: %s", mediaType.schema.reference)
			return
		}
	})
	t.Run("default", func(t *testing.T) {
		response, ok := responses.responses["default"]
		if !ok {
			t.Error("responses.default is not found")
			return
		}
		if response.description != "Unexpected error" {
			t.Errorf("unexpected responses.default.description: %s", response.description)
			return
		}
		mediaType, ok := response.content["application/json"]
		if !ok {
			t.Error("responses.default.content.application/json is not found")
			return
		}
		if mediaType.schema.reference != "#/components/schemas/ErrorModel" {
			t.Errorf("unexpected responses.default.content.application/json.schema.$ref: %s", mediaType.schema.reference)
			return
		}
	})
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected responses:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var responses Responses
			got := yaml.Unmarshal([]byte(tt.yml), &responses)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
		var response Response
		if err := yaml.Unmarshal([]byte(yml), &response); err != nil {
			t.Fatal(err)
		}
		if response.description != "A complex object array response" {
			t.Errorf("unexpected response.description: %s", response.description)
			return
		}
		mediaType, ok := response.content["application/json"]
		if !ok {
			t.Error("response.content.application/json is not found")
			return
		}
		if mediaType.schema.type_ != "array" {
			t.Errorf("unexpected response.content.application/json.schema.type: %s", mediaType.schema.type_)
			return
		}
		if mediaType.schema.items.reference != "#/components/schemas/VeryComplexType" {
			t.Errorf("unexpected response.content.application/json.schema.items.$ref: %s", mediaType.schema.items.reference)
			return
		}
	})
	t.Run("string", func(t *testing.T) {
		yml := `description: A simple string response
content:
  text/plain:
    schema:
      type: string`
		var response Response
		if err := yaml.Unmarshal([]byte(yml), &response); err != nil {
			t.Fatal(err)
		}
		if response.description != "A simple string response" {
			t.Errorf("unexpected response.description: %s", response.description)
			return
		}
		mediaType, ok := response.content["text/plain"]
		if !ok {
			t.Error("response.content.text/plain is not found")
			return
		}
		if mediaType.schema.type_ != "string" {
			t.Errorf("unexpected  response.content.text/plain.schema.type: %s", mediaType.schema.type_)
			return
		}
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
		var response Response
		if err := yaml.Unmarshal([]byte(yml), &response); err != nil {
			t.Fatal(err)
		}
		if response.description != "A simple string response" {
			t.Errorf("unexpected response.description: %s", response.description)
			return
		}
		mediaType, ok := response.content["text/plain"]
		if !ok {
			t.Error("response.content.text/plain is not found")
			return
		}
		if mediaType.schema.type_ != "string" {
			t.Errorf("unexpected response.content.text/plain.schema.type: %s", mediaType.schema.type_)
			return
		}
		t.Run("X-Rate-Limit-Limit", func(t *testing.T) {
			header, ok := response.headers["X-Rate-Limit-Limit"]
			if !ok {
				t.Error("response.headers.X-Rate-Limit-Limit is not found")
				return
			}
			if header.description != "The number of allowed requests in the current period" {
				t.Errorf("unexpected response.headers.X-Rate-Limit-Limit.description: %s", header.description)
				return
			}
			if header.schema.type_ != "integer" {
				t.Errorf("unexpected response.headers.X-Rate-Limit-Limit.schema.type: %s", header.schema.type_)
				return
			}
		})
		t.Run("X-Rate-Limit-Remaining", func(t *testing.T) {
			header, ok := response.headers["X-Rate-Limit-Remaining"]
			if !ok {
				t.Error("response.headers.X-Rate-Limit-Remaining is not found")
				return
			}
			if header.description != "The number of remaining requests in the current period" {
				t.Errorf("unexpected response.headers.X-Rate-Limit-Remaining.description: %s", header.description)
				return
			}
			if header.schema.type_ != "integer" {
				t.Errorf("unexpected response.headers.X-Rate-Limit-Remaining.schema.type: %s", header.schema.type_)
				return
			}
		})
		t.Run("X-Rate-Limit-Reset", func(t *testing.T) {
			header, ok := response.headers["X-Rate-Limit-Reset"]
			if !ok {
				t.Error("response.headers.X-Rate-Limit-Reset is not found")
				return
			}
			if header.description != "The number of seconds left in the current period" {
				t.Errorf("unexpected response.headers.X-Rate-Limit-Reset.description: %s", header.description)
				return
			}
			if header.schema.type_ != "integer" {
				t.Errorf("unexpected response.headers.X-Rate-Limit-Reset.schema.type: %s", header.schema.type_)
				return
			}
		})
	})
	t.Run("no return value", func(t *testing.T) {
		yml := `description: object created`
		var response Response
		if err := yaml.Unmarshal([]byte(yml), &response); err != nil {
			t.Fatal(err)
		}
		if response.description != "object created" {
			t.Errorf("unexpected response.description: %s", response.description)
			return
		}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected response:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var response Response
			got := yaml.Unmarshal([]byte(tt.yml), &response)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
	var target map[string]*Callback
	if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
		t.Fatal(err)
	}
	callback, ok := target["myWebhook"]
	if !ok {
		t.Error("myWebhook is not found")
		return
	}
	pathItem, ok := callback.callback["http://notificationServer.com?transactionId={$request.body#/id}&email={$request.body#/email}"]
	if !ok {
		t.Error("myWebhook.http://notificationServer.com?transactionId={$request.body#/id}&email={$request.body#/email} is not found")
		return
	}
	if pathItem.post.requestBody.description != "Callback payload" {
	}
	mediaType, ok := pathItem.post.requestBody.content["application/json"]
	if !ok {
		t.Error("myWebhook.http://notificationServer.com?transactionId={$request.body#/id}&email={$request.body#/email}.post.requestBody.content.application/json is not found")
		return
	}
	if mediaType.schema.reference != "#/components/schemas/SomePayload" {
		t.Errorf("unexpected myWebhook.http://notificationServer.com?transactionId={$request.body#/id}&email={$request.body#/email}.post.requestBody.schema.$ref: %s", mediaType.schema.reference)
		return
	}
	response, ok := pathItem.post.responses.responses["200"]
	if !ok {
		t.Error("myWebhook.http://notificationServer.com?transactionId={$request.body#/id}&email={$request.body#/email}.post.responses.200 is not found")
		return
	}
	if response.description != "webhook successfully processed and no retries will be performed" {
		t.Errorf("unexpected myWebhook.http://notificationServer.com?transactionId={$request.body#/id}&email={$request.body#/email}.post.responses.200.description: %s", response.description)
		return
	}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected callback:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var callback Callback
			got := yaml.Unmarshal([]byte(tt.yml), &callback)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
		t.Run("application/json", func(t *testing.T) {
			mediaType, ok := target.RequestBody.content["application/json"]
			if !ok {
				t.Error("requestBody.content.application/json is not found")
				return
			}
			if mediaType.schema.reference != "#/components/schemas/Address" {
				t.Errorf("unexpected requestBody.content.application/json.schema.$ref: %s", mediaType.schema.reference)
				return
			}
			t.Run("foo", func(t *testing.T) {
				example, ok := mediaType.examples["foo"]
				if !ok {
					t.Error("requestBody.content.application/json.examples.foo is not found")
					return
				}
				if example.summary != "A foo example" {
					t.Errorf("unexpected requestBody.content.application/json.examples.foo.summary: %s", example.summary)
					return
				}
				if !reflect.DeepEqual(example.value, map[string]interface{}{"foo": "bar"}) {
					t.Errorf("unexpected requestBody.content.application/json.examples.foo.value: %v", example.value)
					return
				}
			})
			t.Run("bar", func(t *testing.T) {
				example, ok := mediaType.examples["bar"]
				if !ok {
					t.Error("requestBody.content.application/json.examples.bar is not found")
					return
				}
				if example.summary != "A bar example" {
					t.Errorf("unexpected requestBody.content.application/json.examples.bar.summary: %s", example.summary)
					return
				}
				if !reflect.DeepEqual(example.value, map[string]interface{}{"bar": "baz"}) {
					t.Errorf("unexpected requestBody.content.application/json.examples.bar.value: %v", example.value)
					return
				}
			})
		})
		t.Run("application/xml", func(t *testing.T) {
			mediaType, ok := target.RequestBody.content["application/xml"]
			if !ok {
				t.Error("requestBody.content.application/xml is not found")
				return
			}
			example, ok := mediaType.examples["xmlExample"]
			if !ok {
				t.Error("requestBody.content.application/xml.examples.xmlExample is not found")
				return
			}
			if example.summary != "This is an example in XML" {
				t.Errorf("unexpected requestBody.content.application/xml.examples.xmlExample.summary: %s", example.summary)
				return
			}
			if example.externalValue != "http://example.org/examples/address-example.xml" {
				t.Errorf("unexpected requestBody.content.application/xml.examples.xmlExample.externalValue: %s", example.externalValue)
				return
			}
		})
		t.Run("text/plain", func(t *testing.T) {
			mediaType, ok := target.RequestBody.content["text/plain"]
			if !ok {
				t.Error("requestBody.content.text/plain is not found")
				return
			}
			example, ok := mediaType.examples["textExample"]
			if !ok {
				t.Error("requestBody.content.text/plain.examples.textExample is not found")
				return
			}
			if example.summary != "This is a text example" {
				t.Errorf("unexpected requestBody.content.text/plain.examples.textExample.summary: %s", example.summary)
				return
			}
			if example.externalValue != "http://foo.bar/examples/address-example.txt" {
				t.Errorf("unexpected requestBody.content.text/plain.examples.textExample.externalValue: %s", example.externalValue)
				return
			}
		})
	})
	/*
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
					parameter := target.Parameters[0]
					if parameter.name != "zipCode" {
						t.Errorf("unexpected parameters.0.name: %s", parameter.name)
						return
					}
					if parameter.in != "query" {
						t.Errorf("unexpected parameters.0.in: %s", parameter.in)
						return
					}
					if parameter.schema.type_ != "string" {
						t.Errorf("unexpected parameters.0.schema.type: %s", parameter.schema.type_)
						return
					}
					if parameter.schema.format != "zip-code" {
						t.Errorf("unexpected parameters.0.schema.format: %s", parameter.schema.format)
						return
					}
				})
	*/
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
		response, ok := target.Responses.responses["200"]
		if !ok {
			t.Error("responses.200 is not found")
			return
		}
		if response.description != "your car appointment has been booked" {
			t.Errorf("unexpected responses.200.description: %s", response.description)
			return
		}
		mediaType, ok := response.content["application/json"]
		if !ok {
			t.Error("responses.200.content.application/json is not found")
			return
		}
		if mediaType.schema.reference != "#/components/schemas/SuccessResponse" {
			t.Errorf("unexpected responses.200.content.application/json.schema.$ref: %s", mediaType.schema.reference)
			return
		}
		example, ok := mediaType.examples["confirmation-success"]
		if !ok {
			t.Error("responses.200.content.application/json.examples.confirmation-success is not found")
			return
		}
		if example.reference != "#/components/examples/confirmation-success" {
			t.Errorf("unexpected responses.200.content.application/json.examples.confirmation-success.$ref: %s", example.reference)
			return
		}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected example:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var example Example
			got := yaml.Unmarshal([]byte(tt.yml), &example)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
		t.Run("/users/{id}", func(t *testing.T) {
			pathItem, ok := target.Paths.paths["/users/{id}"]
			if !ok {
				t.Error("paths./users/{id} is not found")
				return
			}
			id := pathItem.parameters[0]
			if id.name != "id" {
				t.Errorf("unexpected paths./users/{id}.parameters.0.name: %s", id.name)
				return
			}
			if id.in != "path" {
				t.Errorf("unexpected paths./users/{id}.parameters.0.in: %s", id.in)
				return
			}
			if id.required != true {
				t.Errorf("unexpected paths./users/{id}.parameters.0.required: %t", id.required)
				return
			}
			if id.description != "the user identifier, as userId" {
				t.Errorf("unexpected paths./users/{id}.parameters.0.description: %s", id.description)
				return
			}
			if id.schema.type_ != "string" {
				t.Errorf("unexpected paths./users/{id}.parameters.0.schema.type: %s", id.schema.type_)
				return
			}
			response, ok := pathItem.get.responses.responses["200"]
			if !ok {
				t.Error("paths./users/{id}.get.responses.200 is not found")
				return
			}
			if response.description != "the user being returned" {
				t.Errorf("unexpected paths./users/{id}.get.responses.200.description: %s", response.description)
				return
			}
			mediaType, ok := response.content["application/json"]
			if !ok {
				t.Error("paths./users/{id}.get.responses.200.content.application/json is not found")
				return
			}
			if mediaType.schema.type_ != "object" {
				t.Errorf("unexpected paths./users/{id}.get.responses.200.content.application/json.schema.type: %s", mediaType.schema.type_)
				return
			}
			property, ok := mediaType.schema.properties["uuid"]
			if !ok {
				t.Error("paths./users/{id}.get.responses.200.content.application/json.schema.properties.uuid is not found")
				return
			}
			if property.type_ != "string" {
				t.Errorf("unexpected paths./users/{id}.get.responses.200.content.application/json.schema.properties.uuid.type: %s", property.type_)
				return
			}
			if property.format != "uuid" {
				t.Errorf("unexpected paths./users/{id}.get.responses.200.content.application/json.schema.properties.uuid.format: %s", property.format)
				return
			}
			link, ok := response.links["address"]
			if !ok {
				t.Error("paths./users/{id}.get.responses.200.links.address is not found")
				return
			}
			if link.operationID != "getUserAddress" {
				t.Errorf("unexpected paths./users/{id}.get.responses.200.links.address.operationId: %s", link.operationID)
				return
			}
			parameter, ok := link.parameters["userId"]
			if !ok {
				t.Error("paths./users/{id}.get.responses.200.links.address.parameters.userId is not found")
				return
			}
			if parameter != "$request.path.id" {
				t.Errorf("unexpected paths./users/{id}.get.responses.200.links.address.parameters.userId: %v", parameter)
				return
			}
		})
		t.Run("/users/{userid}/address", func(t *testing.T) {
			pathItem, ok := target.Paths.paths["/users/{userid}/address"]
			if !ok {
				t.Error("paths./users/{userid}/address is not found")
				return
			}
			userid := pathItem.parameters[0]
			if userid.name != "userid" {
				t.Errorf("unexpected paths./users/{userid}/address.name: %s", userid.name)
				return
			}
			if userid.in != "path" {
				t.Errorf("unexpected paths./users/{userid}/address.in: %s", userid.in)
				return
			}
			if userid.required != true {
				t.Errorf("unexpected paths./users/{userid}/address.required: %t", userid.required)
				return
			}
			if userid.description != "the user identifier, as userId" {
				t.Errorf("unexpected paths./users/{userid}/address.dEscription: %s", userid.description)
				return
			}
			if userid.schema.type_ != "string" {
				t.Errorf("unexpected paths./users/{userid}/address.schema.type: %s", userid.schema.type_)
				return
			}
			if pathItem.get.operationID != "getUserAddress" {
				t.Errorf("unexpected paths./users/{userid}/address.get.operationId: %s", pathItem.get.operationID)
				return
			}
			response, ok := pathItem.get.responses.responses["200"]
			if !ok {
				t.Error("paths./users/{userid}/address.get.responses.200 is not found")
				return
			}
			if response.description != "the user's address" {
				t.Errorf("unexpected paths./users/{userid}/address.get.responses.200.description: %s", response.description)
				return
			}
		})
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
		link, ok := target.Links["address"]
		if !ok {
			t.Error("links.address is not found")
			return
		}
		if link.operationID != "getUserAddressByUUID" {
			t.Errorf("unexpected links.address.operationId: %s", link.operationID)
			return
		}
		if link.parameters["userUuid"] != "$response.body#/uuid" {
			t.Errorf("unexpected links.address.parameters.userUuid: %s", link.parameters["userUuid"])
			return
		}
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
		link, ok := target.Links["UserRepositories"]
		if !ok {
			t.Error("links.UserRepositories is not found")
		}
		if link.operationRef != "#/paths/~12.0~1repositories~1{username}/get" {
			t.Errorf("unexpected links.UserRepositories.operationRef: %s", link.operationRef)
			return
		}
		if link.parameters["username"] != "$response.body#/username" {
			t.Errorf("unexpected links.address.parameters.userUuid: %s", link.parameters["username"])
			return
		}
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
		link, ok := target.Links["UserRepositories"]
		if !ok {
			t.Error("links.UserRepositories is not found")
		}
		if link.operationRef != "https://na2.gigantic-server.com/#/paths/~12.0~1repositories~1{username}/get" {
			t.Errorf("unexpected links.UserRepositories.operationRef: %s", link.operationRef)
			return
		}
		if link.parameters["username"] != "$response.body#/username" {
			t.Errorf("unexpected links.address.parameters.userUuid: %s", link.parameters["username"])
			return
		}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected link:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var link Link
			got := yaml.Unmarshal([]byte(tt.yml), &link)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
		})
	}
}

func TestHeaderExampleUnmarshalYAML(t *testing.T) {
	yml := `description: The number of allowed requests in the current period
schema:
  type: integer`
	var header Header
	if err := yaml.Unmarshal([]byte(yml), &header); err != nil {
		t.Fatal(err)
	}
	if header.description != "The number of allowed requests in the current period" {
		t.Errorf("unexpected header.description: %s", header.description)
		return
	}
	if header.schema.type_ != "integer" {
		t.Errorf("unexpected header.schema.type: %s", header.schema.type_)
		return
	}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected header:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var header Header
			got := yaml.Unmarshal([]byte(tt.yml), &header)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
		})
	}
}

func TestTagExampleUnmarshal(t *testing.T) {
	yml := `name: pet
description: Pets operations`
	var tag Tag
	if err := yaml.Unmarshal([]byte(yml), &tag); err != nil {
		t.Fatal(err)
	}
	if tag.name != "pet" {
		t.Errorf("unexpected tag.name: %s", tag.name)
		return
	}
	if tag.description != "Pets operations" {
		t.Errorf("unexpected tag.description: %s", tag.description)
		return
	}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected tag:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var tag Tag
			got := yaml.Unmarshal([]byte(tt.yml), &tag)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
		})
	}
}

func TestSchemaExampleUnmarshalYAML(t *testing.T) {
	t.Run("primitive", func(t *testing.T) {
		yml := `type: string
format: email`
		var schema Schema
		if err := yaml.Unmarshal([]byte(yml), &schema); err != nil {
			t.Fatal(err)
		}
		if schema.type_ != "string" {
			t.Errorf("unexpected schema.type: %s", schema.type_)
			return
		}
		if schema.format != "email" {
			t.Errorf("unexpected schema.format: %s", schema.format)
			return
		}
	})
	t.Run("simple model", func(t *testing.T) {
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
		var schema Schema
		if err := yaml.Unmarshal([]byte(yml), &schema); err != nil {
			t.Fatal(err)
		}
		if schema.type_ != "object" {
			t.Errorf("unexpected schema.type: %s", schema.type_)
			return
		}
		if !reflect.DeepEqual(schema.required, []string{"name"}) {
			t.Errorf("unexpected schema.required: %q", schema.required)
			return
		}
		name, ok := schema.properties["name"]
		if !ok {
			t.Error("schema.properties.name is not found")
			return
		}
		if name.type_ != "string" {
			t.Errorf("unexpected schema.properties.name.type: %s", name.type_)
			return
		}
		address, ok := schema.properties["address"]
		if !ok {
			t.Error("schema.properties.address is not found")
			return
		}
		if address.reference != "#/components/schemas/Address" {
			t.Errorf("unexpected schema.properties.address.$ref: %s", address.reference)
			return
		}
		age, ok := schema.properties["age"]
		if !ok {
			t.Error("schema.properties.age is not found")
			return
		}
		if age.type_ != "integer" {
			t.Errorf("unexpected schema.properties.age.type: %s", age.type_)
			return
		}
		if age.format != "int32" {
			t.Errorf("unexpected schema.properties.age.format: %s", age.format)
			return
		}
		if age.minimum != 0 {
			t.Errorf("unexpected schema.properties.age.minimum: %d", age.minimum)
			return
		}
	})
	t.Run("simple string to string map", func(t *testing.T) {
		yml := `type: object
additionalProperties:
  type: string`
		var schema Schema
		if err := yaml.Unmarshal([]byte(yml), &schema); err != nil {
			t.Fatal(err)
		}
		if schema.type_ != "object" {
			t.Errorf("unexpected schema.type: %s", schema.type_)
			return
		}
		if schema.additionalProperties.type_ != "string" {
			t.Errorf("unexpected schema.additiionalProperties.type: %s", schema.additionalProperties.type_)
			return
		}
	})
	t.Run("string to model map", func(t *testing.T) {
		yml := `type: object
additionalProperties:
  $ref: '#/components/schemas/ComplexModel'`
		var schema Schema
		if err := yaml.Unmarshal([]byte(yml), &schema); err != nil {
			t.Fatal(err)
		}
		if schema.type_ != "object" {
			t.Errorf("unexpected schema.type: %s", schema.type_)
			return
		}
		if schema.additionalProperties.reference != "#/components/schemas/ComplexModel" {
			t.Errorf("unexpected schema.additionalProperties.$ref: %s", schema.additionalProperties.reference)
			return
		}
	})
	t.Run("model with example", func(t *testing.T) {
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
		var schema Schema
		if err := yaml.Unmarshal([]byte(yml), &schema); err != nil {
			t.Fatal(err)
		}
		if schema.type_ != "object" {
			t.Errorf("unexpected schema.type: %s", schema.type_)
			return
		}
		id, ok := schema.properties["id"]
		if !ok {
			t.Error("schema.properties.id is not found")
			return
		}
		if id.type_ != "integer" {
			t.Errorf("unexpected schema.properties.id.type: %s", id.type_)
			return
		}
		if id.format != "int64" {
			t.Errorf("unexpected schema.properties.id.format: %s", id.format)
			return
		}
		name, ok := schema.properties["name"]
		if !ok {
			if !ok {
				t.Error("schema.properties.name is not found")
				return
			}
			if name.type_ != "string" {
				t.Errorf("unexpected schema.properties.name.type: %s", name.type_)
				return
			}
		}
		if !reflect.DeepEqual(schema.required, []string{"name"}) {
			t.Errorf("unexpected schema.required: %q", schema.required)
			return
		}
		if example, ok := schema.example.(map[string]interface{}); ok {
			if name, ok := example["name"]; !ok {
				t.Error("schema.example.name is not found")
				return
			} else if name != "Puma" {
				t.Errorf("unexpected schema.example.name: %s", name)
				return
			}
			if id, ok := example["id"]; !ok {
				t.Error("schema.example.id is not found")
				return
			} else if id != uint64(1) {
				t.Errorf("unexpected schema.example.id: %d", id)
				return
			}
		}
	})
	t.Run("models with composition", func(t *testing.T) {
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
		t.Run("ErrorModel", func(t *testing.T) {
			schema, ok := target.Components.schemas["ErrorModel"]
			if !ok {
				t.Error("components.schemas.ErrorModel is not found")
				return
			}
			if schema.type_ != "object" {
				t.Errorf("unexpected components.schemas.ErrorModel.type: %s", schema.type_)
				return
			}
			if !reflect.DeepEqual(schema.required, []string{"message", "code"}) {
				t.Errorf("unexpected components.schemas.ErrorModel.required: %q", schema.required)
				return
			}
			message, ok := schema.properties["message"]
			if !ok {
				t.Error("components.schemas.ErrorModel.properties.message is not found")
				return
			}
			if message.type_ != "string" {
				t.Errorf("unexpected components.schemas.ErrorModel.properties.message.type: %s", message.type_)
				return
			}
			code, ok := schema.properties["code"]
			if !ok {
				t.Error("components.schemas.ErrorModel.properties.code is not found")
				return
			}
			if code.type_ != "integer" {
				t.Errorf("unexpected components.schemas.ErrorModel.properties.code.type: %s", code.type_)
				return
			}
			if code.minimum != 100 {
				t.Errorf("unexpected components.schemas.ErrorModel.properties.code.minimum: %d", code.minimum)
				return
			}
			if code.maximum != 600 {
				t.Errorf("unexpected components.schemas.ErrorModel.properties.code.maximum: %d", code.maximum)
				return
			}
		})
		t.Run("ExtendedErrorModel", func(t *testing.T) {
			schema, ok := target.Components.schemas["ExtendedErrorModel"]
			if !ok {
				t.Error("components.schemas.ExtendedErrorModel is not found")
				return
			}
			ref := schema.allOf[0]
			if ref.reference != "#/components/schemas/ErrorModel" {
				t.Errorf("unexpected components.schemas.ExtendedErrorModel.allOf.0.$ref: %s", ref.reference)
				return
			}
			ext := schema.allOf[1]
			if ext.type_ != "object" {
				t.Errorf("unexpected components.schemas.ExtendedErrorModel.allOf.1.type: %s", ext.type_)
				return
			}
			if !reflect.DeepEqual(ext.required, []string{"rootCause"}) {
				t.Errorf("unexpected components.schemas.ExtendedErrorModel.allOf.1.required: %q", ext.required)
				return
			}
			rootCause, ok := ext.properties["rootCause"]
			if !ok {
				t.Error("components.schemas.ExtendedErrorModel.allOf.1.properties.rootCause is not found")
				return
			}
			if rootCause.type_ != "string" {
				t.Errorf("unexpected components.schemas.ExtendedErrorModel.allOf.1.properties.rootCause.type: %s", rootCause.type_)
				return
			}
		})
	})
	t.Run("models with polymorphism support", func(t *testing.T) {
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
		t.Run("Pet", func(t *testing.T) {
			pet, ok := target.Components.schemas["Pet"]
			if !ok {
				t.Error("components.schemas.Pet is not found")
				return
			}
			if pet.type_ != "object" {
				t.Errorf("unexpected components.schemas.Pet.type: %s", pet.type_)
				return
			}
			if pet.discriminator.propertyName != "petType" {
				t.Errorf("unexpected components.schemas.discriminator.propertyName: %s", pet.discriminator.propertyName)
				return
			}
			name, ok := pet.properties["name"]
			if !ok {
				t.Error("components.schemas.Pet.properties.name is not found")
				return
			}
			if name.type_ != "string" {
				t.Errorf("unexpected components.schemas.Pet.properties.name.type: %s", name.type_)
				return
			}
			petType, ok := pet.properties["petType"]
			if !ok {
				t.Error("components.schemas.Pet.properties.petType is not found")
				return
			}
			if petType.type_ != "string" {
				t.Errorf("unexpected components.schemas.Pet.properties.petType.type: %s", petType.type_)
				return
			}
			if !reflect.DeepEqual(pet.required, []string{"name", "petType"}) {
				t.Errorf("unexpected componets.schemas.Pet.required: %q", pet.required)
				return
			}
		})
		t.Run("Cat", func(t *testing.T) {
			cat, ok := target.Components.schemas["Cat"]
			if !ok {
				t.Error("components.schemas.Cat is not found")
				return
			}
			if cat.description != "A representation of a cat" {
				t.Errorf("unexpected components.schemas.Cat.description: %s", cat.description)
				return
			}
			if cat.allOf[0].reference != "#/components/schemas/Pet" {
				t.Errorf("unexpected components.schemas.Cat.allOf.0.$ref: %s", cat.allOf[0].reference)
				return
			}
			schema := cat.allOf[1]
			if schema.type_ != "object" {
				t.Errorf("unexpected components.schemas.Cat.allOf.1.type: %s", schema.type_)
				return
			}
			huntingSkill, ok := schema.properties["huntingSkill"]
			if !ok {
				t.Error("components.schemas.Cat.allOf.1.properties.huntingSkill is not found")
				return
			}
			if huntingSkill.type_ != "string" {
				t.Errorf("unexpected components.schemas.Cat.allOf.1.properties.huntingSkill.type: %s", huntingSkill.type_)
				return
			}
			if huntingSkill.description != "The measured skill for hunting" {
				t.Errorf("unexpected components.schemas.Cat.allOf.1.properties.huntingSkill.description: %s", huntingSkill.description)
				return
			}
			if !reflect.DeepEqual(schema.required, []string{"huntingSkill"}) {
				t.Errorf("unexpected components.schemas.Cat.allOf.1.required: %q", schema.required)
				return
			}
		})
		t.Run("Dog", func(t *testing.T) {
			dog, ok := target.Components.schemas["Dog"]
			if !ok {
				t.Error("components.schemas.Dog is not found")
				return
			}
			if dog.description != "A representation of a dog" {
				t.Errorf("unexpected components.schemas.Dog.description: %s", dog.description)
				return
			}
			if dog.allOf[0].reference != "#/components/schemas/Pet" {
				t.Errorf("unexpected components.schemas.Dog.allOf.0.$ref: %s", dog.allOf[0].reference)
				return
			}
			schema := dog.allOf[1]
			if schema.type_ != "object" {
				t.Errorf("unexpected components.schemas.Dog.allOf.1.type: %s", schema.type_)
				return
			}
			packSize, ok := schema.properties["packSize"]
			if !ok {
				t.Error("components.schemas.Dog.allOf.1.properties.packSize is not found")
				return
			}
			if packSize.type_ != "integer" {
				t.Errorf("unexpected components.schemas.Dog.allOf.1.properties.packSize.type: %s", packSize.type_)
				return
			}
			if packSize.format != "int32" {
				t.Errorf("unexpected components.schemas.Dog.allOf.1.properties.packSize.format: %s", packSize.format)
				return
			}
			if packSize.description != "the size of the pack the dog is from" {
				t.Errorf("unexpected components.schemas.Dog.allOf.1.properties.packSize.description: %s", packSize.description)
				return
			}
			if packSize.default_ != "0" {
				t.Errorf("unexpected components.schemas.Dog.allOf.1.properties.packSize.default: %s", packSize.default_)
				return
			}
			if packSize.minimum != 0 {
				t.Errorf("unexpected components.schemas.Dog.allOf.1.properties.packSize.minimum: %d", packSize.minimum)
				return
			}
			if !reflect.DeepEqual(schema.required, []string{"packSize"}) {
				t.Errorf("unexpected components.schemas.Dog.allOf.1.required: %q", schema.required)
				return
			}
		})
	})
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected schema:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var schema Schema
			got := yaml.Unmarshal([]byte(tt.yml), &schema)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
		var target map[string]*Schema
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		schema, ok := target["MyResponseType"]
		if !ok {
			t.Error("MyResponseType is not found")
			return
		}
		if schema.oneOf[0].reference != "#/components/schemas/Cat" {
			t.Errorf("unexpected MyResponseType.oneOf.0.$ref: %s", schema.oneOf[0].reference)
			return
		}
		if schema.oneOf[1].reference != "#/components/schemas/Dog" {
			t.Errorf("unexpected MyResponseType.oneOf.1.$ref: %s", schema.oneOf[1].reference)
			return
		}
		if schema.oneOf[2].reference != "#/components/schemas/Lizard" {
			t.Errorf("unexpected MyResponseType.oneOf.2.$ref: %s", schema.oneOf[2].reference)
			return
		}
		if schema.discriminator.propertyName != "petType" {
			t.Errorf("unexpected MyResponseType.discriminator.propertyName: %s", schema.discriminator.propertyName)
			return
		}
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
		var target map[string]*Schema
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		schema, ok := target["MyResponseType"]
		if !ok {
			t.Error("MyResponseType is not found")
			return
		}
		if schema.oneOf[0].reference != "#/components/schemas/Cat" {
			t.Errorf("unexpected MyResponseType.oneOf.0.$ref: %s", schema.oneOf[0].reference)
			return
		}
		if schema.oneOf[1].reference != "#/components/schemas/Dog" {
			t.Errorf("unexpected MyResponseType.oneOf.1.$ref: %s", schema.oneOf[1].reference)
			return
		}
		if schema.oneOf[2].reference != "#/components/schemas/Lizard" {
			t.Errorf("unexpected MyResponseType.oneOf.2.$ref: %s", schema.oneOf[2].reference)
			return
		}
		if schema.oneOf[3].reference != "https://gigantic-server.com/schemas/Monster/schema.json" {
			t.Errorf("unexpected MyResponseType.oneOf.3.$ref: %s", schema.oneOf[3].reference)
			return
		}
		if schema.discriminator.propertyName != "petType" {
			t.Errorf("unexpected MyResponseType.discriminator.propertyName: %s", schema.discriminator.propertyName)
			return
		}
		dog, ok := schema.discriminator.mapping["dog"]
		if !ok {
			t.Error("MyResponseType.discriminator.mapping.dog is not found")
			return
		}
		if dog != "#/components/schemas/Dog" {
			t.Errorf("unexpected MyResponseType.discriminator.mapping.dog: %s", dog)
			return
		}
		monster, ok := schema.discriminator.mapping["monster"]
		if !ok {
			t.Error("MyResponseType.discriminator.mapping.monster is not found")
			return
		}
		if monster != "https://gigantic-server.com/schemas/Monster/schema.json" {
			t.Errorf("unexpected MyResponseType.discriminator.mapping.monster: %s", monster)
			return
		}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected discriminator:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var discriminator Discriminator
			got := yaml.Unmarshal([]byte(tt.yml), &discriminator)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
		})
	}
}

func TestXMLExampleUnmarshalYAML(t *testing.T) {
	t.Run("basic string", func(t *testing.T) {
		yml := `animals:
  type: string
  xml:
    name: animal`
		var target map[string]*Schema
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		animals, ok := target["animals"]
		if !ok {
			t.Error("animals is not found")
			return
		}
		if animals.type_ != "string" {
			t.Errorf("unexpected animals.type: %s", animals.type_)
			return
		}
		if animals.xml.name != "animal" {
			t.Errorf("unexpected animals.xml.name: %s", animals.xml.name)
			return
		}
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
		var target map[string]*Schema
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		person, ok := target["Person"]
		if !ok {
			t.Error("Person is not found")
			return
		}
		if person.type_ != "object" {
			t.Errorf("unexpected Person.type: %s", person.type_)
			return
		}
		id, ok := person.properties["id"]
		if !ok {
			t.Error("Person.properties.id is not found")
			return
		}
		if id.type_ != "integer" {
			t.Errorf("unexpected Person.properties.id.type: %s", id.type_)
			return
		}
		if id.format != "int32" {
			t.Errorf("unexpected Person.properties.id.format: %s", id.format)
			return
		}
		if id.xml.attribute != true {
			t.Errorf("unexpected Person.properties.id.xml.attribute: %t", id.xml.attribute)
			return
		}
		name, ok := person.properties["name"]
		if !ok {
			t.Error("Person.properties.name is not found")
			return
		}
		if name.type_ != "string" {
			t.Errorf("unexpected Person.properties.name.type: %s", name.type_)
			return
		}
		if name.xml.namespace != "http://example.com/schema/sample" {
			t.Errorf("unexpected Person.properties.name.xml.namespace: %s", name.xml.namespace)
			return
		}
		if name.xml.prefix != "sample" {
			t.Errorf("unexpected Person.properties.name.xml.prefix: %s", name.xml.prefix)
			return
		}
	})
	t.Run("wrapped array", func(t *testing.T) {
		yml := `animals:
  type: array
  items:
    type: string
  xml:
    wrapped: true`
		var target map[string]*Schema
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		animals, ok := target["animals"]
		if !ok {
			t.Error("animals is not found")
			return
		}
		if animals.type_ != "array" {
			t.Errorf("unexpected animals.type: %s", animals.type_)
			return
		}
		if animals.items.type_ != "string" {
			t.Errorf("unexpected animals.items.type: %s", animals.items.type_)
			return
		}
		if animals.xml.wrapped != true {
			t.Errorf("unexpected animals.xml.wrapped: %t", animals.xml.wrapped)
			return
		}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected xml:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var xml XML
			got := yaml.Unmarshal([]byte(tt.yml), &xml)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
		})
	}
}

func TestSecuritySchemeExampleUnmarshalYAML(t *testing.T) {
	t.Run("basic auth", func(t *testing.T) {
		yml := `type: http
scheme: basic`
		var securityScheme SecurityScheme
		if err := yaml.Unmarshal([]byte(yml), &securityScheme); err != nil {
			t.Fatal(err)
		}
		if securityScheme.type_ != "http" {
			t.Errorf("unexpected securityScheme.type: %s", securityScheme.type_)
			return
		}
		if securityScheme.scheme != "basic" {
			t.Errorf("unexpected securityScheme.scheme: %s", securityScheme.scheme)
			return
		}
	})
	t.Run("api key", func(t *testing.T) {
		yml := `type: apiKey
name: api_key
in: header`
		var securityScheme SecurityScheme
		if err := yaml.Unmarshal([]byte(yml), &securityScheme); err != nil {
			t.Fatal(err)
		}
		if securityScheme.type_ != "apiKey" {
			t.Errorf("unexpected securityScheme.type: %s", securityScheme.type_)
			return
		}
		if securityScheme.name != "api_key" {
			t.Errorf("unexpected securityScheme.name: %s", securityScheme.name)
			return
		}
		if securityScheme.in != "header" {
			t.Errorf("unexpected securityScheme.in: %s", securityScheme.in)
			return
		}
	})
	t.Run("JWT bearer", func(t *testing.T) {
		yml := `type: http
scheme: bearer
bearerFormat: JWT`
		var securityScheme SecurityScheme
		if err := yaml.Unmarshal([]byte(yml), &securityScheme); err != nil {
			t.Fatal(err)
		}
		if securityScheme.type_ != "http" {
			t.Errorf("unexpected securityScheme.type: %s", securityScheme.type_)
			return
		}
		if securityScheme.scheme != "bearer" {
			t.Errorf("unexpected securityScheme.scheme: %s", securityScheme.scheme)
			return
		}
		if securityScheme.bearerFormat != "JWT" {
			t.Errorf("unexpected securityScheme.bearerFormat: %s", securityScheme.bearerFormat)
			return
		}
	})
	t.Run("implicit oauth2", func(t *testing.T) {
		yml := `type: oauth2
flows:
  implicit:
    authorizationUrl: https://example.com/api/oauth/dialog
    scopes:
      write:pets: modify pets in your account
      read:pets: read your pets`
		var securityScheme SecurityScheme
		if err := yaml.Unmarshal([]byte(yml), &securityScheme); err != nil {
			t.Fatal(err)
		}
		if securityScheme.type_ != "oauth2" {
			t.Errorf("unexpected securityScheme.type: %s", securityScheme.type_)
			return
		}
		if securityScheme.flows.implicit.authorizationURL != "https://example.com/api/oauth/dialog" {
			t.Errorf("unexpected securityScheme.flows.implicit.authorizationURL: %s", securityScheme.flows.implicit.authorizationURL)
			return
		}
		write, ok := securityScheme.flows.implicit.scopes["write:pets"]
		if !ok {
			t.Error("securityScheme.flows.implicit.scopes.write:pets is not found")
			return
		}
		if write != "modify pets in your account" {
			t.Errorf("unexpected securityScheme.flows.implicit.scopes.write:pets: %s", write)
			return
		}
		read, ok := securityScheme.flows.implicit.scopes["read:pets"]
		if !ok {
			t.Error("securityScheme.flows.implicit.scopes.read:pets is not found")
			return
		}
		if read != "read your pets" {
			t.Errorf("unexpected securityScheme.flows.implicit.scopes.read:pets: %s", read)
			return
		}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected securityScheme:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var securityScheme SecurityScheme
			got := yaml.Unmarshal([]byte(tt.yml), &securityScheme)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected oAuthFlows:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var oAuthFlows OAuthFlows
			got := yaml.Unmarshal([]byte(tt.yml), &oAuthFlows)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
	var securityScheme SecurityScheme
	if err := yaml.Unmarshal([]byte(yml), &securityScheme); err != nil {
		t.Fatal(err)
	}
	if securityScheme.type_ != "oauth2" {
		t.Errorf("unexpected securityScheme.type: %s", securityScheme.type_)
		return
	}
	t.Run("implicit", func(t *testing.T) {
		implicit := securityScheme.flows.implicit
		if implicit.authorizationURL != "https://example.com/api/oauth/dialog" {
			t.Errorf("unexpected securityScheme.flows.implicit.authorizationUrl: %s", implicit.authorizationURL)
			return
		}
		write, ok := implicit.scopes["write:pets"]
		if !ok {
			t.Error("securityScheme.flows.implicit.scopes.write:pets is not found")
			return
		}
		if write != "modify pets in your account" {
			t.Errorf("unexpected securityScheme.flows.implicit.scopes.write:pets: %s", write)
			return
		}
		read, ok := implicit.scopes["read:pets"]
		if !ok {
			t.Error("securityScheme.flows.implicit.scopes.read:pets is not found")
			return
		}
		if read != "read your pets" {
			t.Errorf("unexpected securityScheme.flows.implicit.scopes.read:pets: %s", read)
			return
		}
	})
	t.Run("authorizationCode", func(t *testing.T) {
		authorizationCode := securityScheme.flows.authorizationCode
		if authorizationCode.authorizationURL != "https://example.com/api/oauth/dialog" {
			t.Errorf("unexpected securityScheme.flows.authorizationCode.authorizationUrl: %s", authorizationCode.authorizationURL)
			return
		}
		if authorizationCode.tokenURL != "https://example.com/api/oauth/token" {
			t.Errorf("unexpected securityScheme.flows.authorizationCode.tokenUrl: %s", authorizationCode.tokenURL)
			return
		}
		write, ok := authorizationCode.scopes["write:pets"]
		if !ok {
			t.Error("securityScheme.flows.authorizationCode.scopes.write:pets is not found")
			return
		}
		if write != "modify pets in your account" {
			t.Errorf("unexpected securityScheme.flows.authorizationCode.scopes.write:pets: %s", write)
			return
		}
		read, ok := authorizationCode.scopes["read:pets"]
		if !ok {
			t.Error("securityScheme.flows.authorizationCode.scopes.read:pets is not found")
			return
		}
		if read != "read your pets" {
			t.Errorf("unexpected securityScheme.flows.authorizationCode.scopes.read:pets: %s", read)
			return
		}
	})
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected oAuthFlow:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var oAuthFlow OAuthFlow
			got := yaml.Unmarshal([]byte(tt.yml), &oAuthFlow)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
		})
	}
}

func TestSecurityRequirementExampleUnmarshalYAML(t *testing.T) {
	t.Run("non-oauth2", func(t *testing.T) {
		yml := `api_key: []`
		var securityRequirement SecurityRequirement
		if err := yaml.Unmarshal([]byte(yml), &securityRequirement); err != nil {
			t.Fatal(err)
		}
		if apiKey, ok := securityRequirement.securityRequirement["api_key"]; !ok {
			t.Error("api_key is not found")
			return
		} else if !reflect.DeepEqual(apiKey, []string{}) {
			t.Errorf("unexpected api_key: %q", apiKey)
			return
		}
	})
	t.Run("oauth2", func(t *testing.T) {
		yml := `petstore_auth:
- write:pets
- read:pets`
		var securityRequirement SecurityRequirement
		if err := yaml.Unmarshal([]byte(yml), &securityRequirement); err != nil {
			t.Fatal(err)
		}
		if auth, ok := securityRequirement.securityRequirement["petstore_auth"]; !ok {
			t.Error("petstore_auth is not found")
			return
		} else if auth[0] != "write:pets" {
			t.Errorf("unexpected petstore_auth.0: %s", auth[0])
			return
		} else if auth[1] != "read:pets" {
			t.Errorf("unexpected petstore_auth.1: %s", auth[1])
			return
		}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected securityRequirement:\n  got:  %#v\n  want: %#v", got, tt.want)
				return
			}
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
			var securityRequirement SecurityRequirement
			got := yaml.Unmarshal([]byte(tt.yml), &securityRequirement)
			if got == nil {
				t.Error("error is expected but not")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", got, tt.want)
				return
			}
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
