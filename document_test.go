package openapi_test

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestDocument_Validate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Document{}, openapi.ErrRequired{Target: "openapi"}},
		{"withInvalidVersion",
			openapi.Document{
				Version: "1.0",
				Info:    &openapi.Info{},
				Paths:   openapi.Paths{},
			},
			openapi.ErrFormatInvalid{Target: "openapi version", Format: "X.Y.Z"},
		},
		{"withVersion",
			openapi.Document{
				Version: "3.0.0",
			},
			openapi.ErrRequired{Target: "info"},
		},
		{"valid",
			openapi.Document{
				Version: "3.0.0",
				Info:    &openapi.Info{Title: "foo", TermsOfService: exampleCom, Version: "1.0"},
				Paths:   openapi.Paths{},
			},
			nil,
		},
		{"noPaths",
			openapi.Document{
				Version: "3.0.0",
				Info:    &openapi.Info{Title: "foo", TermsOfService: exampleCom, Version: "1.0"},
			},
			openapi.ErrRequired{Target: "paths"},
		},
	}
	testValidater(t, candidates)
}

func TestOASVersion(t *testing.T) {
	candidates := []struct {
		label string
		in    string
		err   error
	}{
		{"empty", "", openapi.ErrRequired{Target: "openapi"}},
		{"invalidVersion", "foobar", openapi.ErrFormatInvalid{Target: "openapi version", Format: "X.Y.Z"}},
		{"swagger", "2.0", openapi.ErrFormatInvalid{Target: "openapi version", Format: "X.Y.Z"}},
		{"valid", "3.0.0", nil},
		{"unsupportedVersion", "4.0.0", openapi.ErrUnsupportedVersion},
		{"invalidMajorVersion", "foo.0.0", openapi.ErrFormatInvalid{Target: "major part of openapi version"}},
		{"invalidMinorVersion", "0.bar.0", openapi.ErrFormatInvalid{Target: "minor part of openapi version"}},
		{"invalidPatchVersion", "0.0.baz", openapi.ErrFormatInvalid{Target: "patch part of openapi version"}},
	}
	for i, c := range candidates {
		t.Run(strconv.Itoa(i)+"/"+c.label, func(t *testing.T) {
			doc := openapi.Document{
				Version: c.in,
				Info:    &openapi.Info{Title: "foo", Version: "1.0"},
				Paths:   openapi.Paths{},
			}
			if err := doc.Validate(); err != c.err {
				if c.err != nil {
					t.Error("error should be occurred, but not")
					return
				}
				t.Errorf("error should not be occured: %s", err)
			}
		})
	}
}

func TestDocument_Walk(t *testing.T) {
	// doc is not valid openapi document but no problem in this test
	rootGetOp := &openapi.Operation{OperationID: "rootGet"}
	fooGetOp := &openapi.Operation{OperationID: "fooGet"}
	fooPostOp := &openapi.Operation{OperationID: "fooPost"}
	barGetOp := &openapi.Operation{OperationID: "barGet"}
	barPostOp := &openapi.Operation{OperationID: "barPost"}
	barPutOp := &openapi.Operation{OperationID: "barPut"}
	barDeleteOp := &openapi.Operation{OperationID: "barDelete"}

	rootPI := &openapi.PathItem{
		Get: rootGetOp,
	}
	fooPI := &openapi.PathItem{
		Get:  fooGetOp,
		Post: fooPostOp,
	}
	barPI := &openapi.PathItem{
		Get:    barGetOp,
		Post:   barPostOp,
		Put:    barPutOp,
		Delete: barDeleteOp,
	}

	doc := &openapi.Document{
		Version: "3.0.1",
		Info: &openapi.Info{
			Title:   "title",
			Version: "v0.0.0",
		},
		Paths: openapi.Paths{
			"/":    rootPI,
			"/foo": fooPI,
			"/bar": barPI,
		},
	}
	getWalkFn := func() openapi.WalkFunc {
		wants := []struct {
			doc       *openapi.Document
			method    string
			path      string
			pathItem  *openapi.PathItem
			operation *openapi.Operation
		}{
			{doc, http.MethodGet, "/", rootPI, rootGetOp},
			{doc, http.MethodDelete, "/bar", barPI, barDeleteOp},
			{doc, http.MethodGet, "/bar", barPI, barGetOp},
			{doc, http.MethodPost, "/bar", barPI, barPostOp},
			{doc, http.MethodPut, "/bar", barPI, barPutOp},
			{doc, http.MethodGet, "/foo", fooPI, fooGetOp},
			{doc, http.MethodPost, "/foo", fooPI, fooPostOp},
		}
		var idx int
		return func(doc *openapi.Document, method, path string, pathItem *openapi.PathItem, op *openapi.Operation) (err error) {
			want := wants[idx]
			defer func() {
				if err != nil {
					t.Logf("\ndoc: %+v\nmethod: %s\npath: %s\npathItem: %+v\nop: %+v\n", doc, method, path, pathItem, op)

				}
			}()
			if !reflect.DeepEqual(doc, want.doc) {
				return fmt.Errorf("unexpected document:\n  got:\t%+v\n  want:\t%+v", doc, want.doc)
			}
			if method != want.method {
				return fmt.Errorf("unexpected method:\n  got:\t%s\n  want:\t%s", method, want.method)
			}
			if path != want.path {
				return fmt.Errorf("unexpected path:\n  got:\t%s\n  want:\t%s", path, want.path)
			}
			if !reflect.DeepEqual(pathItem, want.pathItem) {
				return fmt.Errorf("unexpected path item:\n  got:\t%+v\n  want:\t%+v", pathItem, want.pathItem)
			}
			if !reflect.DeepEqual(op, want.operation) {
				return fmt.Errorf("unexpected operation:\n  got:\t%+v\n  want:\t%+v", op, want.operation)
			}
			idx++
			return nil
		}
	}
	if err := doc.Walk(getWalkFn()); err != nil {
		t.Error(err)
	}
}
