package openapi

import (
	"strconv"
	"testing"
)

//nolint:funlen // most of code is testcases
func TestIsRuntimerExpr(t *testing.T) {
	tests := []struct {
		expr string
		want bool
	}{
		{
			expr: "$method",
			want: true,
		},
		{
			expr: "$request.header.accept",
			want: true,
		},
		{
			expr: "$request.path.id",
			want: true,
		},
		{
			expr: "$request.body#/user/uuid",
			want: true,
		},
		{
			expr: "$url",
			want: true,
		},
		{
			expr: "$response.body#/status",
			want: true,
		},
		{
			expr: "$response.header.Server",
			want: true,
		},
		{
			expr: "invalid.expr",
			want: false,
		},
		{
			expr: "$neither.request.response",
			want: false,
		},
		{
			expr: "$request",
			want: false,
		},
		{
			expr: "$request.header.",
			want: false,
		},
		{
			expr: "$request.body.foo#/fuga",
			want: false,
		},
		{
			expr: "$request.query.value",
			want: true,
		},
		{
			expr: "$request.",
			want: false,
		},
		{
			expr: "https://example.com/{$request.body#foo}",
			want: false,
		},
		{
			expr: "$response.header.ほげ",
			want: false,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"/"+tt.expr, func(t *testing.T) {
			got := IsRuntimeExpr(tt.expr)
			if got != tt.want {
				t.Errorf("unexpected: %t != %t", got, tt.want)
				return
			}
		})
	}
}
