package openapi

import (
	"net/http"
	"strings"
)

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
