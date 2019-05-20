package openapi_test

import (
	"strconv"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestGetOperationByMethod(t *testing.T) {
	pathItem := openapi.PathItem{
		Get:     &openapi.Operation{OperationID: "get"},
		Put:     &openapi.Operation{OperationID: "put"},
		Post:    &openapi.Operation{OperationID: "post"},
		Delete:  &openapi.Operation{OperationID: "delete"},
		Options: &openapi.Operation{OperationID: "options"},
		Head:    &openapi.Operation{OperationID: "head"},
		Patch:   &openapi.Operation{OperationID: "patch"},
		Trace:   &openapi.Operation{OperationID: "trace"},
	}
	candidates := []struct {
		method string
		opID   string
		isNil  bool
	}{
		{"", "", true},
		{"0", "", true},
		{"foobar", "", true},
		{"get", "get", false},
		{"put", "put", false},
		{"POST", "post", false},
		{"Delete", "delete", false},
	}

	for i, c := range candidates {
		t.Run(strconv.Itoa(i)+"/"+c.method+"."+c.opID, func(t *testing.T) {
			op := pathItem.GetOperationByMethod(c.method)
			if !c.isNil && op == nil {
				t.Errorf("operation should be returned: %s", c.method)
				return
			}
			if c.isNil {
				return
			}
			if op.OperationID != c.opID {
				t.Errorf("%s != %s", op.OperationID, c.opID)
			}
		})
	}
}

func TestOperations(t *testing.T) {
	pathItem := openapi.PathItem{
		Get:     &openapi.Operation{OperationID: "get"},
		Put:     &openapi.Operation{OperationID: "put"},
		Post:    &openapi.Operation{OperationID: "post"},
		Delete:  &openapi.Operation{OperationID: "delete"},
		Options: &openapi.Operation{OperationID: "options"},
		Head:    &openapi.Operation{OperationID: "head"},
		Patch:   &openapi.Operation{OperationID: "patch"},
		Trace:   &openapi.Operation{OperationID: "trace"},
	}
	ops := pathItem.Operations()
	if len(ops) != 8 {
		t.Errorf("size of operations is invalid: %d != 8", len(ops))
	}
}
