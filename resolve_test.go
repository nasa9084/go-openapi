package openapi_test

import (
	"reflect"
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
	if !parameter.Required {
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

func TestResolveExample(t *testing.T) {
	example, err := openapi.ResolveExample(doc, "#/components/examples/eg")
	if err != nil {
		t.Error(err)
		return
	}
	if example == nil {
		t.Error("example is nil")
		return
	}
	if example.Summary != "a simple example" {
		t.Errorf("%s != a simple example", example.Summary)
		return
	}
	expected := map[interface{}]interface{}{"foo": "bar"}
	if !reflect.DeepEqual(example.Value, expected) {
		t.Errorf("%+v != %+v", example.Value, expected)
		return
	}
}

func TestResolveRequestBody(t *testing.T) {
	requestBody, err := openapi.ResolveRequestBody(doc, "#/components/requestBodies/user")
	if err != nil {
		t.Error(err)
		return
	}
	if requestBody == nil {
		t.Error("requestBody is nil")
		return
	}
	appJSON, ok := requestBody.Content["application/json"]
	if !ok {
		t.Error("content.application/json not found")
		return
	}
	if appJSON.Schema.Type != "object" {
		t.Errorf("%s != object", appJSON.Schema.Type)
		return
	}
	id, ok := appJSON.Schema.Properties["id"]
	if !ok {
		t.Error("properties.id not found")
		return
	}
	if id.Type != "string" {
		t.Errorf("%s != string", id.Type)
		return
	}
	password, ok := appJSON.Schema.Properties["password"]
	if !ok {
		t.Error("properties.password not found")
		return
	}
	if password.Type != "string" {
		t.Errorf("%s != string", password.Type)
		return
	}
	if password.Format != "password" {
		t.Errorf("%s != password", password.Format)
		return
	}
}

func TestResolveHeader(t *testing.T) {
	header, err := openapi.ResolveHeader(doc, "#/components/headers/x-session")
	if err != nil {
		t.Error(err)
		return
	}
	if header == nil {
		t.Error("header is nil")
		return
	}
	if header.Description != "session token" {
		t.Errorf("%s != session token", header.Description)
		return
	}
	if header.Schema == nil {
		t.Error("header.schema is nil")
		return
	}
	if header.Schema.Type != "string" {
		t.Errorf("%s != string", header.Schema.Type)
		return
	}
}

func TestResolveSecurityScheme(t *testing.T) {
	securityScheme, err := openapi.ResolveSecurityScheme(doc, "#/components/securitySchemes/basicAuth")
	if err != nil {
		t.Error(err)
		return
	}
	if securityScheme == nil {
		t.Error("securityScheme is nil")
		return
	}
	if securityScheme.Type != "http" {
		t.Errorf("%s != http", securityScheme.Type)
		return
	}
	if securityScheme.Scheme != "basic" {
		t.Errorf("%s != basic", securityScheme.Scheme)
		return
	}
}

func TestResolveLink(t *testing.T) {
	link, err := openapi.ResolveLink(doc, "#/components/links/someLink")
	if err != nil {
		t.Error(err)
		return
	}
	if link == nil {
		t.Error("link is nil")
		return
	}
	if link.Description != "a link" {
		t.Errorf("%s != a link", link.Description)
		return
	}
}

func TestResolveCallback(t *testing.T) {
	callback, err := openapi.ResolveCallback(doc, "#/components/callbacks/cb")
	if err != nil {
		t.Error(err)
		return
	}
	if callback == nil {
		t.Error("callback is nil")
		return
	}
	pi, ok := (*callback)["http://example.com"]
	if !ok {
		t.Error("callback.http://example.com not found")
		return
	}
	if pi.Get == nil {
		t.Error("get is nil")
		return
	}
	if pi.Get.Description != "callback" {
		t.Errorf("%s != callback", pi.Get.Description)
		return
	}
}
