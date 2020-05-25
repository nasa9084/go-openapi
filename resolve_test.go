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

	testResolve(t, root)
	testResolveError(t, root)

	testParameterResolve(t, root)
	testRequestBodyResolve(t, root)
	testResponseResolve(t, root)
	testCallbackResolve(t, root)
	testExampleResolve(t, root)
	testLinkResolve(t, root)
	testHeaderResolve(t, root)
	testSchemaResolve(t, root)
	testSecuritySchemeResolve(t, root)

	testResolveTypeAssertionPanicParameter(t, root)
	testResolveTypeAssertionPanicRequestBody(t, root)
	testResolveTypeAssertionPanicResponse(t, root)
	testResolveTypeAssertionPanicCallback(t, root)
	testResolveTypeAssertionPanicExample(t, root)
	testResolveTypeAssertionPanicLink(t, root)
	testResolveTypeAssertionPanicHeader(t, root)
	testResolveTypeAssertionPanicSchema(t, root)
	testResolveTypeAssertionPanicSecurityScheme(t, root)
}

func testResolve(t *testing.T, root *OpenAPI) {
	tests := []struct {
		ref  string
		want interface{}
	}{
		{
			ref:  "#/components/parameters/FooParameter",
			want: root.components.parameters["FooParameter"],
		},
		{
			ref:  "#/components/requestBodies/FooRequest",
			want: root.components.requestBodies["FooRequest"],
		},
		{
			ref:  "#/components/responses/FooResponse",
			want: root.components.responses["FooResponse"],
		},
		{
			ref:  "#/components/callbacks/FooCallback",
			want: root.components.callbacks["FooCallback"],
		},
		{
			ref:  "#/components/examples/FooExample",
			want: root.components.examples["FooExample"],
		},
		{
			ref:  "#/components/links/FooLink",
			want: root.components.links["FooLink"],
		},
		{
			ref:  "#/components/headers/FooHeader",
			want: root.components.headers["FooHeader"],
		},
		{
			ref:  "#/components/schemas/FooSchema",
			want: root.components.schemas["FooSchema"],
		},
		{
			ref:  "#/components/securitySchemes/FooSecurityScheme",
			want: root.components.securitySchemes["FooSecurityScheme"],
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			got, err := resolve(root, tt.ref)
			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, got, tt.want)
		})
	}
}

func testResolveError(t *testing.T, root *OpenAPI) {
	tests := []struct {
		ref  string
		want error
	}{
		{
			ref:  "",
			want: ErrInvalidReference(""),
		},
		{
			ref:  "#/short/reference",
			want: ErrInvalidReference("#/short/reference"),
		},
		{
			ref:  "not/begin/with/sharp",
			want: ErrInvalidReference("not/begin/with/sharp"),
		},
		{
			ref:  "#/foo/bar/baz",
			want: ErrCannotResolved("#/foo/bar/baz", "only supports to resolve under #/components"),
		},
		{
			ref:  "#/components/schemas/UnknownSchema",
			want: ErrCannotResolved("#/components/schemas/UnknownSchema", "not found"),
		},
		{
			ref:  "#/components/unknown/reference",
			want: ErrCannotResolved("#/components/unknown/reference", "unknown component type"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, got := resolve(root, tt.ref)
			if got == nil {
				t.Error("error should not be nil")
				return
			}

			assertSameError(t, got, tt.want)
		})
	}
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

func testResolveTypeAssertionPanicParameter(t *testing.T, root *OpenAPI) {
	t.Run("type assertion panic parameter", func(t *testing.T) {
		parameter := Parameter{
			reference: "#/components/requestBodies/FooRequest",
			root:      root,
		}

		defer func() {
			if err := recover(); err == nil {
				t.Error("panic is expected but not")
			}
		}()

		parameter.resolve() //nolint:errcheck // it panics
	})
}

func testResolveTypeAssertionPanicRequestBody(t *testing.T, root *OpenAPI) {
	t.Run("type assertion panic requestBody", func(t *testing.T) {
		requestBody := RequestBody{
			reference: "#/components/responses/FooResponse",
			root:      root,
		}

		defer func() {
			if err := recover(); err == nil {
				t.Error("panic is expected but not")
			}
		}()

		requestBody.resolve() //nolint:errcheck // it panics
	})
}

func testResolveTypeAssertionPanicResponse(t *testing.T, root *OpenAPI) {
	t.Run("type assertion panic response", func(t *testing.T) {
		response := Response{
			reference: "#/components/callbacks/FooCallback",
			root:      root,
		}

		defer func() {
			if err := recover(); err == nil {
				t.Error("panic is expected but not")
			}
		}()

		response.resolve() //nolint:errcheck // it panics
	})
}

func testResolveTypeAssertionPanicCallback(t *testing.T, root *OpenAPI) {
	t.Run("type assertion panic callback", func(t *testing.T) {
		callback := Callback{
			reference: "#/components/examples/FooExample",
			root:      root,
		}

		defer func() {
			if err := recover(); err == nil {
				t.Error("panic is expected but not")
			}
		}()

		callback.resolve() //nolint:errcheck // it panics
	})
}

func testResolveTypeAssertionPanicExample(t *testing.T, root *OpenAPI) {
	t.Run("type assertion panic example", func(t *testing.T) {
		example := Example{
			reference: "#/components/links/FooLink",
			root:      root,
		}

		defer func() {
			if err := recover(); err == nil {
				t.Error("panic is expected but not")
			}
		}()

		example.resolve() //nolint:errcheck // it panics
	})
}

func testResolveTypeAssertionPanicLink(t *testing.T, root *OpenAPI) {
	t.Run("type assertion panic link", func(t *testing.T) {
		link := Link{
			reference: "#/components/headers/FooHeader",
			root:      root,
		}

		defer func() {
			if err := recover(); err == nil {
				t.Error("panic is expected but not")
			}
		}()

		link.resolve() //nolint:errcheck // it panics
	})
}

func testResolveTypeAssertionPanicHeader(t *testing.T, root *OpenAPI) {
	t.Run("type assertion panic header", func(t *testing.T) {
		header := Header{
			reference: "#/components/schemas/FooSchema",
			root:      root,
		}

		defer func() {
			if err := recover(); err == nil {
				t.Error("panic is expected but not")
			}
		}()

		header.resolve() //nolint:errcheck // it panics
	})
}

func testResolveTypeAssertionPanicSchema(t *testing.T, root *OpenAPI) {
	t.Run("type assertion panic schema", func(t *testing.T) {
		schema := Schema{
			reference: "#/components/securitySchemes/FooSecurityScheme",
			root:      root,
		}

		defer func() {
			if err := recover(); err == nil {
				t.Error("panic is expected but not")
			}
		}()

		schema.resolve() //nolint:errcheck // it panics
	})
}

func testResolveTypeAssertionPanicSecurityScheme(t *testing.T, root *OpenAPI) {
	t.Run("type assertion panic securityScheme", func(t *testing.T) {
		securityScheme := SecurityScheme{
			reference: "#/components/parameters/FooParameter",
			root:      root,
		}

		defer func() {
			if err := recover(); err == nil {
				t.Error("panic is expected but not")
			}
		}()

		securityScheme.resolve() //nolint:errcheck // it panics
	})
}
