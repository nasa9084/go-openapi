package openapi

import (
	"regexp"
)

var (
	urlTemplateVarRegexp = regexp.MustCompile("{[^}]+}") // nolint[gocheckonglobals]

	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$") //nolint[lll]
)

func isOneOf(s string, list []string) bool {
	for _, t := range list {
		if t == s {
			return true
		}
	}

	return false
}

type rawMessage struct {
	unmarshal func(interface{}) error
}

func (msg *rawMessage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	msg.unmarshal = unmarshal
	return nil
}
