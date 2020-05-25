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

func TestIsOneOf(t *testing.T) {
	tests := []struct {
		s    string
		list []string
		want bool
	}{
		{
			s:    "",
			list: []string{},
			want: false,
		},
		{
			s:    "a",
			list: []string{"a", "b"},
			want: true,
		},
		{
			s:    "c",
			list: []string{"a", "b"},
			want: false,
		},
		{
			s:    "a",
			list: nil,
			want: false,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := isOneOf(tt.s, tt.list)
			if got != tt.want {
				t.Errorf("unexpected: %t != %t", got, tt.want)
				return
			}
		})
	}
}
