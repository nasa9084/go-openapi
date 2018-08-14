package openapi_test

import (
	"encoding/json"
	"reflect"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
	yaml "gopkg.in/yaml.v2"
)

func TestSecurityRequirementValidate(t *testing.T) {
	candidates := []struct {
		label  string
		in     openapi.SecurityRequirement
		hasErr bool
	}{
		{"empty", openapi.SecurityRequirement{}, true},
	}
	for _, c := range candidates {
		if err := c.in.Validate(); (err != nil) == c.hasErr {
			t.Error(err)
			return
		}
	}
}

func TestSecurityRequirementUnmarshalYAML(t *testing.T) {
	yml := `apiKey: []`
	secReq := openapi.SecurityRequirement{}
	if err := yaml.Unmarshal([]byte(yml), &secReq); err != nil {
		t.Error(err)
		return
	}
	ar := secReq.Get("apiKey")
	if ar == nil {
		t.Error("securityRequirement.Get(`apiKey`) should be [], not nil")
		return
	}
	if len(ar) != 0 {
		t.Error("securityRequirement.Get(`apiKey`) should be zero length array")
		return
	}
	yml = `apiKey:
  - foo
  - bar
`
	secReq = openapi.SecurityRequirement{}
	if err := yaml.Unmarshal([]byte(yml), &secReq); err != nil {
		t.Error(err)
		return
	}
	ar = secReq.Get("apiKey")
	if ar == nil {
		t.Error("securityRequirement.Get(`apiKey`) should not be nil, but nil")
		return
	}
	if !reflect.DeepEqual(ar, []string{"foo", "bar"}) {
		t.Errorf("securityRequirement.Get(`apiKey`) should be [foo bar], but %v", ar)
		return
	}
}

func TestSecurityRequirementUnmarshalJSON(t *testing.T) {
	jsn := `{"apiKey": []}`
	secReq := openapi.SecurityRequirement{}
	if err := json.Unmarshal([]byte(jsn), &secReq); err != nil {
		t.Error(err)
		return
	}
	ar := secReq.Get("apiKey")
	if ar == nil {
		t.Error("securityRequirement.Get(`apiKey`) should be [], not nil")
		return
	}
	if len(ar) != 0 {
		t.Error("securityRequirement.Get(`apiKey`) should be zero length array")
		return
	}
	jsn = `{"apiKey": ["foo", "bar"]}`
	secReq = openapi.SecurityRequirement{}
	if err := yaml.Unmarshal([]byte(jsn), &secReq); err != nil {
		t.Error(err)
		return
	}
	ar = secReq.Get("apiKey")
	if ar == nil {
		t.Error("securityRequirement.Get(`apiKey`) should not be nil, but nil")
		return
	}
	if !reflect.DeepEqual(ar, []string{"foo", "bar"}) {
		t.Errorf("securityRequirement.Get(`apiKey`) should be [foo bar], but %v", ar)
		return
	}
}
