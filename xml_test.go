package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestXMLValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.XML{}, true},
		{"invalidURLNamespace", openapi.XML{Namespace: "foobar"}, true},
		{"withNamespace", openapi.XML{Namespace: exampleCom}, false},
	}
	testValidater(t, candidates)
}
