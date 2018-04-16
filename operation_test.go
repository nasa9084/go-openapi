package openapi_test

import (
	"reflect"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestSuccessResponse(t *testing.T) {
	candidates := []struct {
		label  string
		in     *openapi.Operation
		resp   *openapi.Response
		status int
		ok     bool
	}{
		{"nil", nil, nil, -1, false},
		{"empty", &openapi.Operation{}, nil, -1, false},
		{"haveInvalid", &openapi.Operation{Responses: openapi.Responses{"foo": &openapi.Response{}}}, nil, 0, false},
		{"haveNilResp", &openapi.Operation{Responses: openapi.Responses{"200": nil}}, nil, 0, false},
		{"have100", &openapi.Operation{Responses: openapi.Responses{"100": &openapi.Response{}}}, nil, 0, false},
		{"have200", &openapi.Operation{Responses: openapi.Responses{"200": &openapi.Response{}}}, &openapi.Response{}, 200, true},
		{"haveDefault", &openapi.Operation{Responses: openapi.Responses{"default": &openapi.Response{}}}, &openapi.Response{}, 0, true},
		{"have200andDefault", &openapi.Operation{Responses: openapi.Responses{"200": &openapi.Response{}, "default": &openapi.Response{}}}, &openapi.Response{}, 200, true},
		{"have2XX", &openapi.Operation{Responses: openapi.Responses{"2XX": &openapi.Response{}}}, &openapi.Response{}, 0, true},
		{"have1XX", &openapi.Operation{Responses: openapi.Responses{"1XX": &openapi.Response{}}}, nil, 0, false},
	}
	for _, c := range candidates {
		resp, status, ok := c.in.SuccessResponse()
		if c.resp != nil && resp == nil {
			t.Error("resp should not be nil")
			return
		}
		if !reflect.DeepEqual(c.resp, resp) {
			t.Errorf("%+v != %+v", c.resp, resp)
			return
		}
		if status != c.status {
			t.Errorf("%d != %d", status, c.status)
			return
		}
		if ok != c.ok {
			t.Errorf("%t != %t", ok, c.ok)
			return
		}
	}
}
