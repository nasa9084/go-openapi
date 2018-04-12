package openapi_test

import "testing"

func TestSuccessResponse(t *testing.T) {
	resp, status, ok := doc.Paths["/"].Get.SuccessResponse()
	if !ok {
		t.Error("cannot find success response")
	}
	if status != 200 {
		t.Error("%d != 200", status)
	}
	if resp == nil {
		t.Error("resp is error")
	}
}
