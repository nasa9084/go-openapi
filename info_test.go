package openapi_test

import (
	"reflect"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
	yaml "gopkg.in/yaml.v2"
)

func TestInfoValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Info{}, openapi.ErrRequired{Target: "info.title"}},
		{"withTitle", openapi.Info{Title: "foo"}, openapi.ErrRequired{Target: "info.version"}},
		{"withTitleAndVersion", openapi.Info{Title: "foo", Version: "1.0"}, nil},
		{"withInvalidToS", openapi.Info{Title: "foo", TermsOfService: "foobar", Version: "1.0"}, openapi.ErrFormatInvalid{Target: "info.termsOfService", Format: "URL"}},
		{"withInvalidLicense", openapi.Info{Title: "foo", Version: "1.0", License: &openapi.License{URL: "foobar"}}, openapi.ErrRequired{Target: "license.name"}},
	}
	testValidater(t, candidates)
}

func TestInfoByExample(t *testing.T) {
	example := `title: Sample Pet Store App
description: This is a sample server for a pet store.
termsOfService: http://example.com/terms/
contact:
  name: API Support
  url: http://www.example.com/support
  email: support@example.com
license:
  name: Apache 2.0
  url: https://www.apache.org/licenses/LICENSE-2.0.html
version: 1.0.1`
	info := openapi.Info{}
	if err := yaml.Unmarshal([]byte(example), &info); err != nil {
		t.Error(err)
		return
	}
	if err := info.Validate(); err != nil {
		t.Error(err)
		return
	}
	expected := openapi.Info{
		Title:          "Sample Pet Store App",
		Description:    "This is a sample server for a pet store.",
		TermsOfService: "http://example.com/terms/",
		Contact: &openapi.Contact{
			Name:  "API Support",
			URL:   "http://www.example.com/support",
			Email: "support@example.com",
		},
		License: &openapi.License{
			Name: "Apache 2.0",
			URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
		},
		Version: "1.0.1",
	}
	if !reflect.DeepEqual(info, expected) {
		t.Errorf("%+v != %+v", info, expected)
		return
	}
}
