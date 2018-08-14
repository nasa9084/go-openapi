package openapi_test

import (
	"reflect"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestComponents(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Components{}, false},
	}
	testValidater(t, candidates)
}

func TestComponentsValidateKeys(t *testing.T) {
	candidates := []struct {
		label  string
		in     openapi.Components
		hasErr bool
	}{
		{"empty", openapi.Components{}, false},
		{"invalidKey", openapi.Components{Parameters: map[string]*openapi.Parameter{"@": &openapi.Parameter{}}}, true},
		{"validKey", openapi.Components{Parameters: map[string]*openapi.Parameter{"foo": &openapi.Parameter{}}}, false},
	}
	for _, c := range candidates {
		if err := openapi.ValidateComponentKeys(c.in); (err != nil) != c.hasErr {
			t.Log(c.label)
			if c.hasErr {
				t.Error("error should be occurred, but not")
				continue
			}
			t.Errorf("error should not be occurred: %s", err)
		}
	}
}
func TestReduceComponentKeys(t *testing.T) {
	candidates := []struct {
		label    string
		in       openapi.Components
		expected []string
	}{
		{"empty", openapi.Components{}, []string{}},
	}
	for _, c := range candidates {
		keys := openapi.ReduceComponentKeys(c.in)
		if !reflect.DeepEqual(keys, c.expected) {
			t.Log(c.label)
			t.Errorf("%+v != %+v", keys, c.expected)
		}
	}
}

func TestReduceComponentObjects(t *testing.T) {
	candidates := []struct {
		label    string
		in       openapi.Components
		expected []openapi.Validater
	}{
		{"empty", openapi.Components{}, []openapi.Validater{}},
	}
	for _, c := range candidates {
		objects := openapi.ReduceComponentObjects(c.in)
		if !reflect.DeepEqual(objects, c.expected) {
			t.Log(c.label)
			t.Errorf("%+v != %+v", objects, c.expected)
		}
	}
}
