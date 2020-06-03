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
	tests := []struct {
		label     string
		reference string
		msg       string
	}{
		{
			label:     "invalid component type",
			reference: "#/components/unknown/Unknown",
			msg:       "local Parameter reference must begin with `#/components/parameters/`",
		},
		{
			label:     "not local reference",
			reference: "https://example.com/example.json#/components/parameters/FooParameter",
			msg:       "not supported reference type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			parameter := Parameter{
				reference: tt.reference,
				root:      root,
			}

			_, got := parameter.resolve()

			want := ErrCannotResolved(tt.reference, tt.msg)

			assertSameError(t, got, want)
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

func testRequestBodyResolveError(t *testing.T, root *OpenAPI) {
	tests := []struct {
		label     string
		reference string
		msg       string
	}{
		{
			label:     "invalid component type",
			reference: "#/components/unknown/Unknown",
			msg:       "local RequestBody reference must begin with `#/components/requestBodies/`",
		},
		{
			label:     "not local reference",
			reference: "https://example.com/example.json#/components/requestBodies/FooRequestBody",
			msg:       "not supported reference type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			requestBody := RequestBody{
				reference: tt.reference,
				root:      root,
			}

			_, got := requestBody.resolve()

			want := ErrCannotResolved(tt.reference, tt.msg)

			assertSameError(t, got, want)
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

func testResponseResolveError(t *testing.T, root *OpenAPI) {
	tests := []struct {
		label     string
		reference string
		msg       string
	}{
		{
			label:     "invalid component type",
			reference: "#/components/unknown/Unknown",
			msg:       "local Response reference must begin with `#/components/responses/`",
		},
		{
			label:     "not local reference",
			reference: "https://example.com/example.json#/components/responses/FooResponse",
			msg:       "not supported reference type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			response := Response{
				reference: tt.reference,
				root:      root,
			}

			_, got := response.resolve()

			want := ErrCannotResolved(tt.reference, tt.msg)

			assertSameError(t, got, want)
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

func testCallbackResolveError(t *testing.T, root *OpenAPI) {
	tests := []struct {
		label     string
		reference string
		msg       string
	}{
		{
			label:     "invalid component type",
			reference: "#/components/unknown/Unknown",
			msg:       "local Callback reference must begin with `#/components/callbacks/`",
		},
		{
			label:     "not local reference",
			reference: "https://example.com/example.json#/components/callbacks/FooCallback",
			msg:       "not supported reference type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			callback := Callback{
				reference: tt.reference,
				root:      root,
			}

			_, got := callback.resolve()

			want := ErrCannotResolved(tt.reference, tt.msg)

			assertSameError(t, got, want)
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

func testExampleResolveError(t *testing.T, root *OpenAPI) {
	tests := []struct {
		label     string
		reference string
		msg       string
	}{
		{
			label:     "invalid component type",
			reference: "#/components/unknown/Unknown",
			msg:       "local Example reference must begin with `#/components/examples/`",
		},
		{
			label:     "not local reference",
			reference: "https://example.com/example.json#/components/examples/FooExample",
			msg:       "not supported reference type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			example := Example{
				reference: tt.reference,
				root:      root,
			}

			_, got := example.resolve()

			want := ErrCannotResolved(tt.reference, tt.msg)

			assertSameError(t, got, want)
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

func testLinkResolveError(t *testing.T, root *OpenAPI) {
	tests := []struct {
		label     string
		reference string
		msg       string
	}{
		{
			label:     "invalid component type",
			reference: "#/components/unknown/Unknown",
			msg:       "local Link reference must begin with `#/components/links/`",
		},
		{
			label:     "not local reference",
			reference: "https://example.com/example.json#/components/links/FooLink",
			msg:       "not supported reference type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			link := Link{
				reference: tt.reference,
				root:      root,
			}

			_, got := link.resolve()

			want := ErrCannotResolved(tt.reference, tt.msg)

			assertSameError(t, got, want)
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

func testHeaderResolveError(t *testing.T, root *OpenAPI) {
	tests := []struct {
		label     string
		reference string
		msg       string
	}{
		{
			label:     "invalid component type",
			reference: "#/components/unknown/Unknown",
			msg:       "local Header reference must begin with `#/components/headers/`",
		},
		{
			label:     "not local reference",
			reference: "https://example.com/example.json#/components/headers/FooHeader",
			msg:       "not supported reference type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			header := Header{
				reference: tt.reference,
				root:      root,
			}

			_, got := header.resolve()

			want := ErrCannotResolved(tt.reference, tt.msg)

			assertSameError(t, got, want)
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

func testSchemaResolveError(t *testing.T, root *OpenAPI) {
	tests := []struct {
		label     string
		reference string
		msg       string
	}{
		{
			label:     "invalid component type",
			reference: "#/components/unknown/Unknown",
			msg:       "local Schema reference must begin with `#/components/schemas/`",
		},
		{
			label:     "not local reference",
			reference: "https://example.com/example.json#/components/schemas/FooSchema",
			msg:       "not supported reference type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			schema := Schema{
				reference: tt.reference,
				root:      root,
			}

			_, got := schema.resolve()

			want := ErrCannotResolved(tt.reference, tt.msg)

			assertSameError(t, got, want)
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

func testSecuritySchemeResolveError(t *testing.T, root *OpenAPI) {
	tests := []struct {
		label     string
		reference string
		msg       string
	}{
		{
			label:     "invalid component type",
			reference: "#/components/unknown/Unknown",
			msg:       "local SecurityScheme reference must begin with `#/components/securitySchemes/`",
		},
		{
			label:     "not local reference",
			reference: "https://example.com/example.json#/components/securitySchemes/FooSecurityScheme",
			msg:       "not supported reference type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			securityScheme := SecurityScheme{
				reference: tt.reference,
				root:      root,
			}

			_, got := securityScheme.resolve()

			want := ErrCannotResolved(tt.reference, tt.msg)

			assertSameError(t, got, want)
		})
	}
}
