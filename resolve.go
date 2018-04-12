package openapi

import (
	"errors"
	"strings"
)

var ErrTypeAssertion = errors.New("type assertion error")

func resolve(root *Document, ref string) (interface{}, error) {
	switch {
	case strings.HasPrefix(ref, "#/"):
		path := strings.Split(ref, "/")
		if len(path) < 2 {
			return nil, errors.New("ref string length invalid")
		}
		return root.resolve(cdr(path))
	default:
		return nil, errors.New("cannot resolve relative document")
	}
}

func (doc *Document) resolve(path []string) (interface{}, error) {
	switch s := car(path); s {
	case "components":
		return doc.Components.resolve(cdr(path))
	default:
		return nil, errors.New("unknown reference path: " + s)
	}
}

func (components *Components) resolve(path []string) (interface{}, error) {
	if len(path) != 2 {
		return nil, errors.New("cannot resolve")
	}
	var ret interface{}
	var ok bool
	switch s := car(path); s {
	case "schemas":
		ret, ok = components.Schemas[car(cdr(path))]
	case "responses":
		ret, ok = components.Responses[car(cdr(path))]
	case "parameters":
		ret, ok = components.Parameters[car(cdr(path))]
	case "examples":
		ret, ok = components.Examples[car(cdr(path))]
	case "requestBodies":
		ret, ok = components.RequestBodies[car(cdr(path))]
	case "headers":
		ret, ok = components.Headers[car(cdr(path))]
	case "securitySchemes":
		ret, ok = components.SecuritySchemes[car(cdr(path))]
	case "links":
		ret, ok = components.Links[car(cdr(path))]
	case "callbacks":
		ret, ok = components.Callbacks[car(cdr(path))]
	default:
		return nil, errors.New("unknown reference path: " + s)
	}
	if !ok {
		return nil, errors.New("not found: " + car(cdr(path)))
	}
	return ret, nil
}

// ResolveSchema resolves a schema reference string.
func ResolveSchema(root *Document, ref string) (*Schema, error) {
	si, err := resolve(root, ref)
	if err != nil {
		return nil, err
	}
	if s, ok := si.(*Schema); ok {
		return s, nil
	}
	return nil, ErrTypeAssertion
}

// ResolveResponse resolves a response reference string.
func ResolveResponse(root *Document, ref string) (*Response, error) {
	ri, err := resolve(root, ref)
	if err != nil {
		return nil, err
	}
	if r, ok := ri.(*Response); ok {
		return r, nil
	}
	return nil, ErrTypeAssertion
}

// ResolveParameter resolves a response reference string.
func ResolveParameter(root *Document, ref string) (*Parameter, error) {
	pi, err := resolve(root, ref)
	if err != nil {
		return nil, err
	}
	if p, ok := pi.(*Parameter); ok {
		return p, nil
	}
	return nil, ErrTypeAssertion
}

// ResolveExample resolves an example reference string.
func ResolveExample(root *Document, ref string) (*Example, error) {
	ei, err := resolve(root, ref)
	if err != nil {
		return nil, err
	}
	if e, ok := ei.(*Example); ok {
		return e, nil
	}
	return nil, ErrTypeAssertion
}

// ResolveRequestBody resolves a requestBody reference string.
func ResolveRequestBody(root *Document, ref string) (*RequestBody, error) {
	ri, err := resolve(root, ref)
	if err != nil {
		return nil, err
	}
	if r, ok := ri.(*RequestBody); ok {
		return r, nil
	}
	return nil, ErrTypeAssertion
}

// ResolveHeader resolves a header reference string.
func ResolveHeader(root *Document, ref string) (*Header, error) {
	hi, err := resolve(root, ref)
	if err != nil {
		return nil, err
	}
	if h, ok := hi.(*Header); ok {
		return h, nil
	}
	return nil, ErrTypeAssertion
}

// ResolveSecurityScheme resolves a securityScheme reference string.
func ResolveSecurityScheme(root *Document, ref string) (*SecurityScheme, error) {
	si, err := resolve(root, ref)
	if err != nil {
		return nil, err
	}
	if s, ok := si.(*SecurityScheme); ok {
		return s, nil
	}
	return nil, ErrTypeAssertion
}

// ResolveLink resolves a link reference string.
func ResolveLink(root *Document, ref string) (*Link, error) {
	li, err := resolve(root, ref)
	if err != nil {
		return nil, err
	}
	if l, ok := li.(*Link); ok {
		return l, nil
	}
	return nil, ErrTypeAssertion
}

// ResolveCallback resolves a callback reference string.
func ResolveCallback(root *Document, ref string) (*Callback, error) {
	ci, err := resolve(root, ref)
	if err != nil {
		return nil, err
	}
	if c, ok := ci.(*Callback); ok {
		return c, nil
	}
	return nil, ErrTypeAssertion
}
