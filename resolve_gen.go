// Code generated by mkresolver.go. DO NOT EDIT.

package openapi

import (
	"fmt"
	"strings"
)

func (v *Parameter) resolve() (*Parameter, error) {
	if v.reference == "" {
		return v, nil
	}

	if v.resolved != nil {
		return v.resolved, nil
	}

	if strings.HasPrefix(v.reference, `#/`) {
		return v.resolveLocal()
	}

	return nil, ErrCannotResolved(v.reference, `not supported reference type`)
}

func (v *Parameter) resolveLocal() (*Parameter, error) {
	prefix := `#/components/parameters/`

	if !strings.HasPrefix(v.reference, prefix) {
		return nil, ErrCannotResolved(v.reference, fmt.Sprintf("local Parameter reference must begin with `%s`", prefix))
	}

	key := strings.TrimPrefix(v.reference, prefix)

	if resolved, ok := v.root.components.parameters[key]; ok {
		return resolved, nil
	}

	return nil, ErrCannotResolved(v.reference, `not found`)
}

func (v *RequestBody) resolve() (*RequestBody, error) {
	if v.reference == "" {
		return v, nil
	}

	if v.resolved != nil {
		return v.resolved, nil
	}

	if strings.HasPrefix(v.reference, `#/`) {
		return v.resolveLocal()
	}

	return nil, ErrCannotResolved(v.reference, `not supported reference type`)
}

func (v *RequestBody) resolveLocal() (*RequestBody, error) {
	prefix := `#/components/requestBodies/`

	if !strings.HasPrefix(v.reference, prefix) {
		return nil, ErrCannotResolved(v.reference, fmt.Sprintf("local RequestBody reference must begin with `%s`", prefix))
	}

	key := strings.TrimPrefix(v.reference, prefix)

	if resolved, ok := v.root.components.requestBodies[key]; ok {
		return resolved, nil
	}

	return nil, ErrCannotResolved(v.reference, `not found`)
}

func (v *Response) resolve() (*Response, error) {
	if v.reference == "" {
		return v, nil
	}

	if v.resolved != nil {
		return v.resolved, nil
	}

	if strings.HasPrefix(v.reference, `#/`) {
		return v.resolveLocal()
	}

	return nil, ErrCannotResolved(v.reference, `not supported reference type`)
}

func (v *Response) resolveLocal() (*Response, error) {
	prefix := `#/components/responses/`

	if !strings.HasPrefix(v.reference, prefix) {
		return nil, ErrCannotResolved(v.reference, fmt.Sprintf("local Response reference must begin with `%s`", prefix))
	}

	key := strings.TrimPrefix(v.reference, prefix)

	if resolved, ok := v.root.components.responses[key]; ok {
		return resolved, nil
	}

	return nil, ErrCannotResolved(v.reference, `not found`)
}

func (v *Callback) resolve() (*Callback, error) {
	if v.reference == "" {
		return v, nil
	}

	if v.resolved != nil {
		return v.resolved, nil
	}

	if strings.HasPrefix(v.reference, `#/`) {
		return v.resolveLocal()
	}

	return nil, ErrCannotResolved(v.reference, `not supported reference type`)
}

func (v *Callback) resolveLocal() (*Callback, error) {
	prefix := `#/components/callbacks/`

	if !strings.HasPrefix(v.reference, prefix) {
		return nil, ErrCannotResolved(v.reference, fmt.Sprintf("local Callback reference must begin with `%s`", prefix))
	}

	key := strings.TrimPrefix(v.reference, prefix)

	if resolved, ok := v.root.components.callbacks[key]; ok {
		return resolved, nil
	}

	return nil, ErrCannotResolved(v.reference, `not found`)
}

func (v *Example) resolve() (*Example, error) {
	if v.reference == "" {
		return v, nil
	}

	if v.resolved != nil {
		return v.resolved, nil
	}

	if strings.HasPrefix(v.reference, `#/`) {
		return v.resolveLocal()
	}

	return nil, ErrCannotResolved(v.reference, `not supported reference type`)
}

func (v *Example) resolveLocal() (*Example, error) {
	prefix := `#/components/examples/`

	if !strings.HasPrefix(v.reference, prefix) {
		return nil, ErrCannotResolved(v.reference, fmt.Sprintf("local Example reference must begin with `%s`", prefix))
	}

	key := strings.TrimPrefix(v.reference, prefix)

	if resolved, ok := v.root.components.examples[key]; ok {
		return resolved, nil
	}

	return nil, ErrCannotResolved(v.reference, `not found`)
}

func (v *Link) resolve() (*Link, error) {
	if v.reference == "" {
		return v, nil
	}

	if v.resolved != nil {
		return v.resolved, nil
	}

	if strings.HasPrefix(v.reference, `#/`) {
		return v.resolveLocal()
	}

	return nil, ErrCannotResolved(v.reference, `not supported reference type`)
}

func (v *Link) resolveLocal() (*Link, error) {
	prefix := `#/components/links/`

	if !strings.HasPrefix(v.reference, prefix) {
		return nil, ErrCannotResolved(v.reference, fmt.Sprintf("local Link reference must begin with `%s`", prefix))
	}

	key := strings.TrimPrefix(v.reference, prefix)

	if resolved, ok := v.root.components.links[key]; ok {
		return resolved, nil
	}

	return nil, ErrCannotResolved(v.reference, `not found`)
}

func (v *Header) resolve() (*Header, error) {
	if v.reference == "" {
		return v, nil
	}

	if v.resolved != nil {
		return v.resolved, nil
	}

	if strings.HasPrefix(v.reference, `#/`) {
		return v.resolveLocal()
	}

	return nil, ErrCannotResolved(v.reference, `not supported reference type`)
}

func (v *Header) resolveLocal() (*Header, error) {
	prefix := `#/components/headers/`

	if !strings.HasPrefix(v.reference, prefix) {
		return nil, ErrCannotResolved(v.reference, fmt.Sprintf("local Header reference must begin with `%s`", prefix))
	}

	key := strings.TrimPrefix(v.reference, prefix)

	if resolved, ok := v.root.components.headers[key]; ok {
		return resolved, nil
	}

	return nil, ErrCannotResolved(v.reference, `not found`)
}

func (v *Schema) resolve() (*Schema, error) {
	if v.reference == "" {
		return v, nil
	}

	if v.resolved != nil {
		return v.resolved, nil
	}

	if strings.HasPrefix(v.reference, `#/`) {
		return v.resolveLocal()
	}

	return nil, ErrCannotResolved(v.reference, `not supported reference type`)
}

func (v *Schema) resolveLocal() (*Schema, error) {
	prefix := `#/components/schemas/`

	if !strings.HasPrefix(v.reference, prefix) {
		return nil, ErrCannotResolved(v.reference, fmt.Sprintf("local Schema reference must begin with `%s`", prefix))
	}

	key := strings.TrimPrefix(v.reference, prefix)

	if resolved, ok := v.root.components.schemas[key]; ok {
		return resolved, nil
	}

	return nil, ErrCannotResolved(v.reference, `not found`)
}

func (v *SecurityScheme) resolve() (*SecurityScheme, error) {
	if v.reference == "" {
		return v, nil
	}

	if v.resolved != nil {
		return v.resolved, nil
	}

	if strings.HasPrefix(v.reference, `#/`) {
		return v.resolveLocal()
	}

	return nil, ErrCannotResolved(v.reference, `not supported reference type`)
}

func (v *SecurityScheme) resolveLocal() (*SecurityScheme, error) {
	prefix := `#/components/securitySchemes/`

	if !strings.HasPrefix(v.reference, prefix) {
		return nil, ErrCannotResolved(v.reference, fmt.Sprintf("local SecurityScheme reference must begin with `%s`", prefix))
	}

	key := strings.TrimPrefix(v.reference, prefix)

	if resolved, ok := v.root.components.securitySchemes[key]; ok {
		return resolved, nil
	}

	return nil, ErrCannotResolved(v.reference, `not found`)
}
