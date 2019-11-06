package openapi

import (
	"regexp"
)

var (
	urlTemplateVarRegexp = regexp.MustCompile("{[^}]+}") // nolint[gocheckonglobals]

	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$") //nolint[lll]
)

type raw []byte

func (v *raw) UnmarshalYAML(b []byte) error {
	*v = b
	return nil
}

func isOneOf(s string, list []string) bool {
	for _, t := range list {
		if t == s {
			return true
		}
	}
	return false
}
