package openapi

import (
	"strings"
)

func resolve(root *OpenAPI, ref string) (interface{}, error) {
	if !strings.HasPrefix(ref, "#/") {
		// currently only supports local reference
		return nil, ErrInvalidReference(ref)
	}

	parts := strings.Split(ref, "/")

	if len(parts) != 4 {
		return nil, ErrInvalidReference(ref)
	}

	if parts[1] != "components" {
		return nil, ErrCannotResolved(ref, "only supports to resolve under #/components")
	}

	var (
		ret interface{}
		ok  bool
	)

	next := parts[3]

	switch parts[2] {
	case "schemas":
		ret, ok = root.components.schemas[next]
	case "responses":
		ret, ok = root.components.responses[next]
	case "parameters":
		ret, ok = root.components.parameters[next]
	case "examples":
		ret, ok = root.components.examples[next]
	case "requestBodies":
		ret, ok = root.components.requestBodies[next]
	case "headers":
		ret, ok = root.components.headers[next]
	case "securitySchemes":
		ret, ok = root.components.securitySchemes[next]
	case "links":
		ret, ok = root.components.links[next]
	case "callbacks":
		ret, ok = root.components.callbacks[next]
	default:
		return nil, ErrCannotResolved(ref, "unknown component type")
	}

	if !ok {
		return nil, ErrCannotResolved(ref, "not found")
	}

	return ret, nil
}
