package openapi

import (
	"errors"
	"fmt"
)

// codebeat:disable[TOO_MANY_IVARS]

// SecurityRequirement Object
type SecurityRequirement struct {
	document *Document
	mp       map[string][]string
}

// Get returns required security schemes. If there is not given name,
// this function returns nil.
func (secReq SecurityRequirement) Get(name string) []string {
	val, ok := secReq.mp[name]
	if !ok {
		return nil
	}
	return val
}

// Validate the values of SecurityRequirement object.
func (secReq SecurityRequirement) Validate() error {
	if len(secReq.mp) == 0 {
		return nil
	}
	if secReq.document == nil {
		return errors.New("missing root document for security requirement")
	}
	components := secReq.document.Components
	if components == nil {
		return errors.New("components object in parent document is nil")
	}
	for name, arr := range secReq.mp {
		secScheme, ok := components.SecuritySchemes[name]
		if !ok {
			return fmt.Errorf("%s is not declared in securitySchemes under the components object", name)
		}
		switch secScheme.Type {
		case "oauth2":
			for _, scope := range arr {
				_, implicit := secScheme.Flows.Implicit.Scopes[scope]
				_, password := secScheme.Flows.Password.Scopes[scope]
				_, cc := secScheme.Flows.ClientCredentials.Scopes[scope]
				_, ac := secScheme.Flows.AuthorizationCode.Scopes[scope]
				if !implicit && !password && !cc && !ac {
					return fmt.Errorf("%s is not defined in securitySchemes under the components object", scope)
				}
			}
		default:
			if arr != nil && len(arr) != 0 {
				return fmt.Errorf("securityRequirements for %s type must be empty", secScheme.Type)
			}
		}
	}
	return nil
}

func (secReq *SecurityRequirement) setDocument(doc *Document) {
	secReq.document = doc
}
