package openapi

import (
	"reflect"
	"testing"
)

func TestCar(t *testing.T) {
	ss := []string{"foo", "bar", "baz"}
	s := car(ss)
	if s != "foo" {
		t.Errorf("%s != foo", s)
		return
	}
}

func TestCdr(t *testing.T) {
	ss := []string{"foo", "bar", "baz"}
	expected := []string{"bar", "baz"}
	s := cdr(ss)
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("%+v != %+v", s, expected)
		return
	}
}
