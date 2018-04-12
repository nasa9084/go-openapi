package openapi_test

import (
	"os"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

var doc *openapi.Document

func TestMain(m *testing.M) {
	var err error
	doc, err = openapi.LoadFile("test/testspec.yaml")
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestLoadFile(t *testing.T) {
	doc, err := openapi.LoadFile("test/testspec.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if doc.Version != "3.0.0" {
		t.Errorf("api version is not valid")
		return
	}
	info := doc.Info
	if info.Title != "openapi specification test" {
		t.Errorf("info.title is not valid")
		return
	}
	if info.Version != "1.0" {
		t.Errorf("info.version is not valid")
		return
	}
	paths := doc.Paths
	if paths["/"].Get.Responses["200"].Description != "ok" {
		t.Errorf("paths./.get.responses.200.description is not valid")
		return
	}
}
