package openapi_test

import (
	"reflect"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestComponents(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Components{}, nil},
	}
	testValidater(t, candidates)
}

func TestComponentsValidateKeys(t *testing.T) {
	candidates := []struct {
		label string
		in    openapi.Components
		err   error
	}{
		{"empty", openapi.Components{}, nil},
		{"invalidKey", openapi.Components{Parameters: map[string]*openapi.Parameter{"@": &openapi.Parameter{}}}, openapi.MapKeyFormatError},
		{"validKey", openapi.Components{Parameters: map[string]*openapi.Parameter{"foo": &openapi.Parameter{}}}, nil},
	}
	for _, c := range candidates {
		if err := openapi.ValidateComponentKeys(c.in); err != c.err {
			t.Log(c.label)
			t.Errorf("error should be %s, but %s", c.err, err)
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
