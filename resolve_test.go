package openapi

import (
	"reflect"
	"strconv"
	"testing"
)

func TestResolveParameter(t *testing.T) {
	root := &OpenAPI{
		components: &Components{
			parameters: map[string]*Parameter{
				"FooParameter": {
					name: "FooParameter",
				},
			},
		},
	}

	root.setRoot(root)

	tests := []struct {
		ref  string
		want *Parameter
	}{
		{
			ref:  "#/components/parameters/FooParameter",
			want: root.components.parameters["FooParameter"],
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			got, err := resolve(root, tt.ref)
			if err != nil {
				t.Fatal(err)
			}
			tt.want.setRoot(root)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected:\ngot:  %+v\nwant: %+v", got, tt.want)
				return
			}
		})
	}
}

func TestResolveRequestBody(t *testing.T) {
	root := &OpenAPI{
		components: &Components{
			requestBodies: map[string]*RequestBody{
				"FooRequest": {
					description: "FooRequest",
				},
			},
		},
	}

	root.setRoot(root)

	tests := []struct {
		ref  string
		want *RequestBody
	}{
		{
			ref:  "#/components/requestBodies/FooRequest",
			want: root.components.requestBodies["FooRequest"],
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			got, err := resolve(root, tt.ref)
			if err != nil {
				t.Fatal(err)
			}
			tt.want.setRoot(root)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected:\ngot:  %+v\nwant: %+v", got, tt.want)
				return
			}
		})
	}
}

func TestResolveResponse(t *testing.T) {
	root := &OpenAPI{
		components: &Components{
			responses: map[string]*Response{
				"FooResponse": {
					description: "FooResponse",
				},
			},
		},
	}

	root.setRoot(root)

	tests := []struct {
		ref  string
		want *Response
	}{
		{
			ref:  "#/components/responses/FooResponse",
			want: root.components.responses["FooResponse"],
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			got, err := resolve(root, tt.ref)
			if err != nil {
				t.Fatal(err)
			}
			tt.want.setRoot(root)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected:\ngot:  %+v\nwant: %+v", got, tt.want)
				return
			}
		})
	}
}

func TestResolveCallback(t *testing.T) {
	root := &OpenAPI{
		components: &Components{
			callbacks: map[string]*Callback{
				"FooCallback": {
					callback: map[string]*PathItem{
						"/v1": {
							summary: "FooCallbackPathItem",
						},
					},
				},
			},
		},
	}

	root.setRoot(root)

	tests := []struct {
		ref  string
		want *Callback
	}{
		{
			ref:  "#/components/callbacks/FooCallback",
			want: root.components.callbacks["FooCallback"],
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			got, err := resolve(root, tt.ref)
			if err != nil {
				t.Fatal(err)
			}
			tt.want.setRoot(root)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected:\ngot:  %+v\nwant: %+v", got, tt.want)
				return
			}
		})
	}
}

func TestResolveExample(t *testing.T) {
	root := &OpenAPI{
		components: &Components{
			examples: map[string]*Example{
				"FooExample": {
					summary: "FooExample",
				},
			},
		},
	}

	root.setRoot(root)

	tests := []struct {
		ref  string
		want *Example
	}{
		{
			ref:  "#/components/examples/FooExample",
			want: root.components.examples["FooExample"],
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			got, err := resolve(root, tt.ref)
			if err != nil {
				t.Fatal(err)
			}
			tt.want.setRoot(root)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected:\ngot:  %+v\nwant: %+v", got, tt.want)
				return
			}
		})
	}
}

func TestResolveLink(t *testing.T) {
	root := &OpenAPI{
		components: &Components{
			links: map[string]*Link{
				"FooLink": {
					operationID: "FooLink",
				},
			},
		},
	}

	root.setRoot(root)

	tests := []struct {
		ref  string
		want *Link
	}{
		{
			ref:  "#/components/links/FooLink",
			want: root.components.links["FooLink"],
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			got, err := resolve(root, tt.ref)
			if err != nil {
				t.Fatal(err)
			}
			tt.want.setRoot(root)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected:\ngot:  %+v\nwant: %+v", got, tt.want)
				return
			}
		})
	}
}

func TestResolveHeader(t *testing.T) {
	root := &OpenAPI{
		components: &Components{
			headers: map[string]*Header{
				"FooHeader": {
					description: "FooHeader",
				},
			},
		},
	}

	root.setRoot(root)

	tests := []struct {
		ref  string
		want *Header
	}{
		{
			ref:  "#/components/headers/FooHeader",
			want: root.components.headers["FooHeader"],
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			got, err := resolve(root, tt.ref)
			if err != nil {
				t.Fatal(err)
			}
			tt.want.setRoot(root)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected:\ngot:  %+v\nwant: %+v", got, tt.want)
				return
			}
		})
	}
}

func TestResolveSchema(t *testing.T) {
	root := &OpenAPI{
		components: &Components{
			schemas: map[string]*Schema{
				"FooSchema": {
					title: "FooSchema",
				},
			},
		},
	}

	root.setRoot(root)

	tests := []struct {
		ref  string
		want *Schema
	}{
		{
			ref:  "#/components/schemas/FooSchema",
			want: root.components.schemas["FooSchema"],
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			got, err := resolve(root, tt.ref)
			if err != nil {
				t.Fatal(err)
			}
			tt.want.setRoot(root)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected:\ngot:  %+v\nwant: %+v", got, tt.want)
				return
			}
		})
	}
}

func TestResolveSecurityScheme(t *testing.T) {
	root := &OpenAPI{
		components: &Components{
			securitySchemes: map[string]*SecurityScheme{
				"FooSecurityScheme": {
					name: "FooSecurityScheme",
				},
			},
		},
	}

	root.setRoot(root)

	tests := []struct {
		ref  string
		want *SecurityScheme
	}{
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
			tt.want.setRoot(root)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected:\ngot:  %+v\nwant: %+v", got, tt.want)
				return
			}
		})
	}
}

func TestResolveError(t *testing.T) {
	root := &OpenAPI{
		components: &Components{
			parameters:      map[string]*Parameter{},
			requestBodies:   map[string]*RequestBody{},
			responses:       map[string]*Response{},
			callbacks:       map[string]*Callback{},
			examples:        map[string]*Example{},
			links:           map[string]*Link{},
			headers:         map[string]*Header{},
			schemas:         map[string]*Schema{},
			securitySchemes: map[string]*SecurityScheme{},
		},
	}
	tests := []struct {
		ref  string
		want error
	}{
		{
			ref:  "",
			want: ErrInvalidReference(""),
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
			ref:  "#/components/schemas/FooSchema",
			want: ErrCannotResolved("#/components/schemas/FooSchema", "not found"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, got := resolve(root, tt.ref)
			if got == nil {
				t.Error("error should not be nil")
				return
			}
			if got != tt.want {
				t.Errorf("unexpected error:\ngot:  %+v\nwant: %+v", got, tt.want)
				return
			}
		})
	}
}
