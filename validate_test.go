package openapi

import (
	"strconv"
	"testing"
)

func TestValidateURLTemplate(t *testing.T) {
	tests := []struct {
		url  string
		want error
	}{
		{
			url:  "https://development.gigantic-server.com/v1",
			want: nil,
		},
		{
			url:  "https://{username}.gigantic-server.com:{port}/{basePath}",
			want: nil,
		},
		{
			url:  "{scheme}://developer.uspto.gov/ds-api",
			want: nil,
		},
		{
			url:  "example.com/foo/bar",
			want: nil,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"/"+tt.url, func(t *testing.T) {
			if err := validateURLTemplate(tt.url); err != tt.want {
				t.Errorf("unexpected:\n  got:  %v\n  want: %v", err, tt.want)
				return
			}
		})
	}
}
