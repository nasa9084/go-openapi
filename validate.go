package openapi

import (
	"net/url"
	"strings"
)

func validateURLTemplate(s string) error {
	var scheme, address, port, path string
	rest := s
	splitScheme := strings.SplitN(rest, "://", 2)
	if len(splitScheme) == 2 {
		scheme = splitScheme[0]
		rest = splitScheme[1]
	} else {
		rest = splitScheme[0]
	}

	splitPath := strings.SplitN(rest, "/", 2)
	if len(splitPath) == 2 {
		rest = splitPath[0]
		path = splitPath[1]
	} else {
		rest = splitPath[0]
	}

	splitPort := strings.SplitN(rest, ":", 2)
	if len(splitPort) == 2 {
		address = splitPort[0]
		port = splitPort[1]
	} else {
		address = splitPort[0]
	}
	scheme = urlTemplateVarRegexp.ReplaceAllLiteralString(scheme, "http")
	address = urlTemplateVarRegexp.ReplaceAllLiteralString(address, "placeholder")
	port = urlTemplateVarRegexp.ReplaceAllLiteralString(port, "80")
	path = urlTemplateVarRegexp.ReplaceAllLiteralString(path, "placeholder")

	s = ""
	if scheme != "" {
		s += scheme + "://"
	}
	if address != "" {
		s += address
	}
	if port != "" {
		s += ":" + port
	}
	if path != "" {
		s += "/" + path
	}

	_, err := url.Parse(s)
	return err
}
