package openapi_test

import (
	"reflect"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
	yaml "gopkg.in/yaml.v2"
)

func TestServerValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Server{}, openapi.ErrRequired{Target: "server.url"}},
		{"invalidURL", openapi.Server{URL: "foobar%"}, openapi.ErrFormatInvalid{Target: "server.url", Format: "URL"}},
		{"withURL", openapi.Server{URL: exampleCom}, nil},
	}
	testValidater(t, candidates)
}

func TestServerByExample(t *testing.T) {
	spec := `url: https://development.gigantic-server.com/v1
description: Development server`
	var server openapi.Server
	if err := yaml.Unmarshal([]byte(spec), &server); err != nil {
		t.Error(err)
		return
	}
	expect := openapi.Server{
		URL:         "https://development.gigantic-server.com/v1",
		Description: "Development server",
	}
	if !reflect.DeepEqual(server, expect) {
		t.Errorf("%+v != %+v", server, expect)
		return
	}
}
