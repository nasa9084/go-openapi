package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
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
