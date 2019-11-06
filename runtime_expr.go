package openapi

import (
	"strings"
)

func IsRuntimeExpr(expr string) bool {
	if !strings.Contains(expr, "{") {
		return isRuntimeExpr(expr)
	}
	for i, b := range expr {
		if b != '{' {
			continue
		}
		for j, e := range expr[i+1:] {
			if e != '}' {
				continue
			}
			if !isRuntimeExpr(expr[i+1 : j+i]) {
				return false
			}
		}
	}
	return true
}

func isRuntimeExpr(expr string) bool {
	// expression = ( "$url" | "$method" | "$statusCode" | "$request." source | "$response." source )
	if isOneOf(expr, []string{"$url", "$method", "$statusCode"}) {
		return true
	}
	var source string
	switch {
	case strings.HasPrefix(expr, "$request."):
		source = expr[9:]
	case strings.HasPrefix(expr, "$response."):
		source = expr[10:]
	default:
		return false
	}
	// source = ( header-reference | query-reference | path-reference | body-reference )
	switch {
	case strings.HasPrefix(source, "header."):
		// header-reference = "header." token
		return isRFC7230Token(source[7:])
	case strings.HasPrefix(source, "query."):
		// query-reference = "query." name
		return source[6:] != ""
	case strings.HasPrefix(source, "path."):
		// path-reference = "path." name
		return source[5:] != ""
	case strings.HasPrefix(source, "body#"):
		// body-reference = "body" ["#" fragment]
		// fragment = a JSON Pointer [RFC 6901](https://tools.ietf.org/html/rfc6901)
		return isJSONPointer(source[5:])
	}

	return false
}

func isJSONPointer(s string) bool {
	return strings.HasPrefix(s, "/")
}

const rfc7230tchar = "-!#$%&'*+.^_`|~0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func isRFC7230Token(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !strings.ContainsRune(rfc7230tchar, r) {
			return false
		}
	}
	return true
}
