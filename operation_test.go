package openapi_test

import (
	"reflect"
	"strconv"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestOperation_Validate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Operation{}, openapi.ErrRequired{Target: "operation.responses"}},
		{"duplicatedParameter", openapi.Operation{Responses: openapi.Responses{}, Parameters: []*openapi.Parameter{&openapi.Parameter{Name: "foo", In: "query"}, &openapi.Parameter{Name: "foo", In: "query"}}}, openapi.ErrParameterDuplicated},
		{"valid", openapi.Operation{Responses: openapi.Responses{}}, nil},
	}
	testValidater(t, candidates)
}

func TestOperation_SuccessResponse(t *testing.T) {
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
	for i, c := range candidates {
		t.Run(strconv.Itoa(i)+"/"+c.label, func(t *testing.T) {
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
		})
	}
}
