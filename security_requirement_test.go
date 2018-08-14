package openapi_test

import (
	"encoding/json"
	"reflect"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
	yaml "gopkg.in/yaml.v2"
)

func TestSecurityRequirementValidate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.SecurityRequirement{}, false},
	}
	testValidater(t, candidates)
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
	// empty list case
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
	// not empty list case
	jsn = `{"apiKey": ["foo", "bar"]}`
	secReq = openapi.SecurityRequirement{}
	if err := json.Unmarshal([]byte(jsn), &secReq); err != nil {
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
	// invalid case
	jsn = `{"apiKey": "foo"}` // value should be array
	secReq = openapi.SecurityRequirement{}
	if err := json.Unmarshal([]byte(jsn), &secReq); err == nil {
		t.Error("error should be occurred, but not")
		return
	}
}
