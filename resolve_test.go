package openapi

import (
	"strconv"
	"testing"
)

func TestResolve(t *testing.T) {
	root := &OpenAPI{
		components: &Components{
			parameters: map[string]*Parameter{
				"FooParameter": {
					name: "FooParameter",
				},
			},
			requestBodies: map[string]*RequestBody{
				"FooRequest": {
					description: "FooRequest",
				},
			},
			responses: map[string]*Response{
				"FooResponse": {
					description: "FooResponse",
				},
			},
			callbacks: map[string]*Callback{
				"FooCallback": {
					callback: map[string]*PathItem{
						"/v1": {
							summary: "FooCallbackPathItem",
						},
					},
				},
			},
			examples: map[string]*Example{
				"FooExample": {
					summary: "FooExample",
				},
			},
			links: map[string]*Link{
				"FooLink": {
					operationID: "FooLink",
				},
			},
			headers: map[string]*Header{
				"FooHeader": {
					description: "FooHeader",
				},
			},
			schemas: map[string]*Schema{
				"FooSchema": {
					title: "FooSchema",
				},
			},
			securitySchemes: map[string]*SecurityScheme{
				"FooSecurityScheme": {
					name: "FooSecurityScheme",
				},
			},
		},
	}

	testParameterResolve(t, root)
	testRequestBodyResolve(t, root)
	testResponseResolve(t, root)
	testCallbackResolve(t, root)
	testExampleResolve(t, root)
	testLinkResolve(t, root)
	testHeaderResolve(t, root)
	testSchemaResolve(t, root)
	testSecuritySchemeResolve(t, root)

	testParameterResolveError(t, root)
	testRequestBodyResolveError(t, root)
	testResponseResolveError(t, root)
	testCallbackResolveError(t, root)
	testExampleResolveError(t, root)
	testLinkResolveError(t, root)
	testHeaderResolveError(t, root)
	testSchemaResolveError(t, root)
	testSecuritySchemeResolveError(t, root)
}

func testParameterResolve(t *testing.T, root *OpenAPI) {
	tests := []struct {
		ref  string
		want *Parameter
	}{
		{
			ref:  "#/components/parameters/FooParameter",
			want: root.components.parameters["FooParameter"],
		},
		{
			ref: "",
			want: &Parameter{
				root: root,
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			v := Parameter{
				reference: tt.ref,
			}
			v.setRoot(root)
			got, err := v.resolve()
			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, got, tt.want)
		})
	}
}

func testParameterResolveError(t *testing.T, root *OpenAPI) {
	t.Run("parameter.resolve error", func(t *testing.T) {
		parameter := Parameter{
			reference: "#/components/unknown/Unknown",
			root:      root,
		}

		_, got := parameter.resolve()

		want := ErrCannotResolved(parameter.reference, "local Parameter reference must begin with `#/components/parameters/`")

		assertSameError(t, got, want)
	})
}

func testRequestBodyResolve(t *testing.T, root *OpenAPI) {
	tests := []struct {
		ref  string
		want *RequestBody
	}{
		{
			ref:  "#/components/requestBodies/FooRequest",
			want: root.components.requestBodies["FooRequest"],
		},
		{
			ref: "",
			want: &RequestBody{
				root: root,
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			v := RequestBody{
				reference: tt.ref,
			}
			v.setRoot(root)
			got, err := v.resolve()
			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, got, tt.want)
		})
	}
}

func testRequestBodyResolveError(t *testing.T, root *OpenAPI) {
	t.Run("requestBody.resolve error", func(t *testing.T) {
		requestBody := RequestBody{
			reference: "#/components/unknown/Unknown",
			root:      root,
		}

		_, got := requestBody.resolve()

		want := ErrCannotResolved(
			requestBody.reference,
			"local RequestBody reference must begin with `#/components/requestBodies/`",
		)

		assertSameError(t, got, want)
	})
}

func testResponseResolve(t *testing.T, root *OpenAPI) {
	tests := []struct {
		ref  string
		want *Response
	}{
		{
			ref:  "#/components/responses/FooResponse",
			want: root.components.responses["FooResponse"],
		},
		{
			ref: "",
			want: &Response{
				root: root,
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			v := Response{
				reference: tt.ref,
			}
			v.setRoot(root)
			got, err := v.resolve()
			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, got, tt.want)
		})
	}
}

func testResponseResolveError(t *testing.T, root *OpenAPI) {
	t.Run("response.resolve error", func(t *testing.T) {
		response := Response{
			reference: "#/components/unknown/Unknown",
			root:      root,
		}

		_, got := response.resolve()

		want := ErrCannotResolved(response.reference, "local Response reference must begin with `#/components/responses/`")

		assertSameError(t, got, want)
	})
}

func testCallbackResolve(t *testing.T, root *OpenAPI) {
	tests := []struct {
		ref  string
		want *Callback
	}{
		{
			ref:  "#/components/callbacks/FooCallback",
			want: root.components.callbacks["FooCallback"],
		},
		{
			ref: "",
			want: &Callback{
				root: root,
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			v := Callback{
				reference: tt.ref,
			}
			v.setRoot(root)
			got, err := v.resolve()
			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, got, tt.want)
		})
	}
}

func testCallbackResolveError(t *testing.T, root *OpenAPI) {
	t.Run("callback.resolve error", func(t *testing.T) {
		callback := Callback{
			reference: "#/components/unknown/Unknown",
			root:      root,
		}

		_, got := callback.resolve()

		want := ErrCannotResolved(callback.reference, "local Callback reference must begin with `#/components/callbacks/`")

		assertSameError(t, got, want)
	})
}

func testExampleResolve(t *testing.T, root *OpenAPI) {
	tests := []struct {
		ref  string
		want *Example
	}{
		{
			ref:  "#/components/examples/FooExample",
			want: root.components.examples["FooExample"],
		},
		{
			ref: "",
			want: &Example{
				root: root,
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			v := Example{
				reference: tt.ref,
			}
			v.setRoot(root)
			got, err := v.resolve()
			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, got, tt.want)
		})
	}
}

func testExampleResolveError(t *testing.T, root *OpenAPI) {
	t.Run("example.resolve error", func(t *testing.T) {
		example := Example{
			reference: "#/components/unknown/Unknown",
			root:      root,
		}

		_, got := example.resolve()

		want := ErrCannotResolved(example.reference, "local Example reference must begin with `#/components/examples/`")

		assertSameError(t, got, want)
	})
}

func testLinkResolve(t *testing.T, root *OpenAPI) {
	tests := []struct {
		ref  string
		want *Link
	}{
		{
			ref:  "#/components/links/FooLink",
			want: root.components.links["FooLink"],
		},
		{
			ref: "",
			want: &Link{
				root: root,
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			v := Link{
				reference: tt.ref,
			}
			v.setRoot(root)
			got, err := v.resolve()
			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, got, tt.want)
		})
	}
}

func testLinkResolveError(t *testing.T, root *OpenAPI) {
	t.Run("link.resolve error", func(t *testing.T) {
		link := Link{
			reference: "#/components/unknown/Unknown",
			root:      root,
		}

		_, got := link.resolve()

		want := ErrCannotResolved(link.reference, "local Link reference must begin with `#/components/links/`")

		assertSameError(t, got, want)
	})
}

func testHeaderResolve(t *testing.T, root *OpenAPI) {
	tests := []struct {
		ref  string
		want *Header
	}{
		{
			ref:  "#/components/headers/FooHeader",
			want: root.components.headers["FooHeader"],
		},
		{
			ref: "",
			want: &Header{
				root: root,
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			v := Header{
				reference: tt.ref,
			}
			v.setRoot(root)
			got, err := v.resolve()
			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, got, tt.want)
		})
	}
}

func testHeaderResolveError(t *testing.T, root *OpenAPI) {
	t.Run("header.resolve error", func(t *testing.T) {
		header := Header{
			reference: "#/components/unknown/Unknown",
			root:      root,
		}

		_, got := header.resolve()

		want := ErrCannotResolved(header.reference, "local Header reference must begin with `#/components/headers/`")

		assertSameError(t, got, want)
	})
}

func testSchemaResolve(t *testing.T, root *OpenAPI) {
	tests := []struct {
		ref  string
		want *Schema
	}{
		{
			ref:  "#/components/schemas/FooSchema",
			want: root.components.schemas["FooSchema"],
		},
		{
			ref: "",
			want: &Schema{
				root: root,
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			v := Schema{
				reference: tt.ref,
			}
			v.setRoot(root)
			got, err := v.resolve()
			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, got, tt.want)
		})
	}
}

func testSchemaResolveError(t *testing.T, root *OpenAPI) {
	t.Run("schema.resolve error", func(t *testing.T) {
		schema := Schema{
			reference: "#/components/unknown/Unknown",
			root:      root,
		}

		_, got := schema.resolve()

		want := ErrCannotResolved(schema.reference, "local Schema reference must begin with `#/components/schemas/`")

		assertSameError(t, got, want)
	})
}

func testSecuritySchemeResolve(t *testing.T, root *OpenAPI) {
	tests := []struct {
		ref  string
		want *SecurityScheme
	}{
		{
			ref:  "#/components/securitySchemes/FooSecurityScheme",
			want: root.components.securitySchemes["FooSecurityScheme"],
		},
		{
			ref: "",
			want: &SecurityScheme{
				root: root,
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			v := SecurityScheme{
				reference: tt.ref,
			}
			v.setRoot(root)
			got, err := v.resolve()
			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, got, tt.want)
		})
	}
}

func testSecuritySchemeResolveError(t *testing.T, root *OpenAPI) {
	t.Run("securityScheme.resolve error", func(t *testing.T) {
		securityScheme := SecurityScheme{
			reference: "#/components/unknown/Unknown",
			root:      root,
		}

		_, got := securityScheme.resolve()

		want := ErrCannotResolved(
			securityScheme.reference,
			"local SecurityScheme reference must begin with `#/components/securitySchemes/`",
		)

		assertSameError(t, got, want)
	})
}
