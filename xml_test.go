package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestXML_Validate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.XML{}, openapi.ErrRequired{Target: "xml.namespace"}},
		{"invalidURLNamespace", openapi.XML{Namespace: "foobar"}, openapi.ErrFormatInvalid{Target: "xml.namespace", Format: "URL"}},
		{"withNamespace", openapi.XML{Namespace: exampleCom}, nil},
	}
	testValidater(t, candidates)
}
