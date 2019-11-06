package openapi

import (
	"reflect"
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
						"/v1": &PathItem{
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
	root.setRoot(root)
	tests := []struct {
		ref  string
		want interface{}
	}{
		{
			ref: "#/components/parameters/FooParameter",
			want: &Parameter{
				name: "FooParameter",
			},
		},
		{
			ref: "#/components/requestBodies/FooRequest",
			want: &RequestBody{
				description: "FooRequest",
			},
		},
		{
			ref: "#/components/responses/FooResponse",
			want: &Response{
				description: "FooResponse",
			},
		},
		{
			ref: "#/components/callbacks/FooCallback",
			want: &Callback{
				callback: map[string]*PathItem{
					"/v1": &PathItem{
						summary: "FooCallbackPathItem",
					},
				},
			},
		},
		{
			ref: "#/components/examples/FooExample",
			want: &Example{
				summary: "FooExample",
			},
		},
		{
			ref: "#/components/links/FooLink",
			want: &Link{
				operationID: "FooLink",
			},
		},
		{
			ref: "#/components/headers/FooHeader",
			want: &Header{
				description: "FooHeader",
			},
		},
		{
			ref: "#/components/schemas/FooSchema",
			want: &Schema{
				title: "FooSchema",
			},
		},
		{
			ref: "#/components/securitySchemes/FooSecurityScheme",
			want: &SecurityScheme{
				name: "FooSecurityScheme",
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.ref, func(t *testing.T) {
			got, err := resolve(root, tt.ref)
			if err != nil {
				t.Fatal(err)
			}
			tt.want.(interface{ setRoot(*OpenAPI) }).setRoot(root)
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
