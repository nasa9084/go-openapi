package openapi_test

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

var doc, _ = openapi.LoadFile("test/testspec.yaml")

func TestResolveSchema(t *testing.T) {
	schema, err := openapi.ResolveSchema(doc, "#/components/schemas/definition")
	if err != nil {
		t.Error(err)
		return
	}
	if schema == nil {
		t.Error("schema is nil")
		return
	}
	if schema.Title != "schemaTitle" {
		t.Errorf("%s != schemaTitle", schema.Title)
		return
	}
	if schema.Type != "string" {
		t.Errorf("%s != string", schema.Type)
		return
	}
}

func TestResolveResponse(t *testing.T) {
	response, err := openapi.ResolveResponse(doc, "#/components/responses/notFound")
	if err != nil {
		t.Error(err)
		return
	}
	if response == nil {
		t.Error("response is nil")
		return
	}
	if response.Description != "not found" {
		t.Errorf("%s != not found", response.Description)
		return
	}
	appJSON, ok := response.Content["application/json"]
	if !ok {
		t.Error("content.application/json is invalid")
		return
	}
	if appJSON.Schema.Type != "string" {
		t.Errorf("%s != string", appJSON.Schema.Type)
		return
	}
}

func TestResolveParameter(t *testing.T) {
	parameter, err := openapi.ResolveParameter(doc, "#/components/parameters/pathParam")
	if err != nil {
		t.Error(err)
		return
	}
	if parameter == nil {
		t.Error("parameter is nil")
		return
	}
	if parameter.Name != "id" {
		t.Errorf("%s != id", parameter.Name)
		return
	}
	if parameter.In != "path" {
		t.Errorf("%s != path", parameter.In)
		return
	}
	if parameter.Description != "user id" {
		t.Errorf("%s != user id", parameter.Description)
		return
	}
	if parameter.Required != true {
		t.Errorf("%t != true", parameter.Required)
		return
	}
	if parameter.Schema == nil {
		t.Errorf("parameter.Schema is nil")
		return
	}
	if parameter.Schema.Type != "string" {
		t.Errorf("%s != string", parameter.Schema.Type)
		return
	}
}
