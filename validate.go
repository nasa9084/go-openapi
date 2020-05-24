package openapi

import (
	"net/url"
	"strings"
)

func validateURLTemplate(s string) error {
	scheme, addr, port, path := splitURLTemplate(s)

	scheme = urlTemplateVarRegexp.ReplaceAllLiteralString(scheme, "http")
	addr = urlTemplateVarRegexp.ReplaceAllLiteralString(addr, "placeholder")
	port = urlTemplateVarRegexp.ReplaceAllLiteralString(port, "80")
	path = urlTemplateVarRegexp.ReplaceAllLiteralString(path, "placeholder")

	_, err := url.Parse(buildURL(scheme, addr, port, path))

	return err
}

func splitURLTemplate(s string) (scheme, addr, port, path string) {
	if strings.Contains(s, "://") {
		splitScheme := strings.SplitN(s, "://", 2)
		scheme = splitScheme[0]
		s = splitScheme[1]
	}

	if strings.Contains(s, "/") {
		splitPath := strings.SplitN(s, "/", 2)
		s = splitPath[0]
		path = splitPath[1]
	}

	if strings.Contains(s, ":") {
		splitPort := strings.SplitN(s, ":", 2)
		addr = splitPort[0]
		port = splitPort[1]
	} else {
		addr = s
	}

	return
}

func buildURL(scheme, addr, port, path string) string {
	var ret string

	if scheme != "" {
		ret += scheme + "://"
	}

	if addr != "" {
		ret += addr
	}

	if port != "" {
		ret += ":" + port
	}

	if path != "" {
		ret += "/" + path
	}

	return ret
}
