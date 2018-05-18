package openapi

import (
	"errors"
	"net/http"
	"strings"
)

// codebeat:disable[TOO_MANY_IVARS]

// PathItem Object
type PathItem struct {
	Ref string `yaml:"$ref"`

	Summary     string
	Description string
	Get         *Operation
	Put         *Operation
	Post        *Operation
	Delete      *Operation
	Options     *Operation
	Head        *Operation
	Patch       *Operation
	Trace       *Operation
	Servers     []*Server
	Parameters  []*Parameter
}

var methods = []string{
	http.MethodGet,
	http.MethodPut,
	http.MethodPost,
	http.MethodDelete,
	http.MethodOptions,
	http.MethodHead,
	http.MethodPatch,
	http.MethodTrace,
}

// GetOperation returns a operation object associated with given method.
// The method is case-insensitive, converted to upper case in this function.
// If the method is invalid, this function will return nil.
func (pathItem *PathItem) GetOperation(method string) *Operation {
	switch strings.ToUpper(method) {
	case http.MethodGet:
		return pathItem.Get
	case http.MethodPost:
		return pathItem.Post
	case http.MethodPut:
		return pathItem.Put
	case http.MethodDelete:
		return pathItem.Delete
	case http.MethodOptions:
		return pathItem.Options
	case http.MethodHead:
		return pathItem.Head
	case http.MethodPatch:
		return pathItem.Patch
	case http.MethodTrace:
		return pathItem.Trace
	default:
		return nil
	}
}

// Operations returns a map containing operation object as a
// value associated with a HTTP method as a key.
// If an operation is nil, it won't be added returned map, so
// the size of returned map is not same always.
func (pathItem PathItem) Operations() map[string]*Operation {
	ops := map[string]*Operation{}
	for _, method := range methods {
		if op := pathItem.GetOperation(method); op != nil {
			ops[method] = op
		}
	}
	return ops
}

// Validate the values of PathItem object.
func (pathItem PathItem) Validate() error {
	validaters := []validater{}
	for _, op := range pathItem.Operations() {
		validaters = append(validaters, op)
	}
	for _, s := range pathItem.Servers {
		validaters = append(validaters, s)
	}
	if hasDuplicatedParameter(pathItem.Parameters) {
		return errors.New("some parameter is duplicated")
	}
	for _, p := range pathItem.Parameters {
		validaters = append(validaters, p)
	}
	return validateAll(validaters)
}

func hasDuplicatedParameter(parameters []*Parameter) bool {
	for i, p := range parameters {
		for _, q := range parameters[i+1:] {
			if p.Name == q.Name && p.In == q.In {
				return true
			}
		}
	}
	return false
}
