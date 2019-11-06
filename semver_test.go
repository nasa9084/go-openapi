package openapi

import (
	"strconv"
	"testing"
)

func TestSplitSemVer(t *testing.T) {
	type want struct {
		ver  string
		pre  string
		meta string
	}
	tests := []struct {
		input string
		want  want
	}{
		{
			input: "",
			want:  want{},
		},
		{
			input: "1",
			want: want{
				ver: "1",
			},
		},
		{
			input: "-1.2.3",
			want: want{
				ver: "-1.2.3",
			},
		},
		{
			input: "1.9.0",
			want: want{
				ver: "1.9.0",
			},
		},
		{
			input: "1.0.0-alpha",
			want: want{
				ver: "1.0.0",
				pre: "alpha",
			},
		},
		{
			input: "1.0.0-alpha.1",
			want: want{
				ver: "1.0.0",
				pre: "alpha.1",
			},
		},
		{
			input: "1.0.0-0.3.7",
			want: want{
				ver: "1.0.0",
				pre: "0.3.7",
			},
		},
		{
			input: "1.0.0-x.7.z.92",
			want: want{
				ver: "1.0.0",
				pre: "x.7.z.92",
			},
		},
		{
			input: "1.0.0-alpha+001",
			want: want{
				ver:  "1.0.0",
				pre:  "alpha",
				meta: "001",
			},
		},
		{
			input: "1.0.0+20130313144700",
			want: want{
				ver:  "1.0.0",
				meta: "20130313144700",
			},
		},
		{
			input: "1.0.0-beta+exp.sha.5114f85",
			want: want{
				ver:  "1.0.0",
				pre:  "beta",
				meta: "exp.sha.5114f85",
			},
		},
		{
			input: "1.0.0+exp.sha.5114f85-beta",
			want: want{
				ver:  "1.0.0",
				pre:  "beta",
				meta: "exp.sha.5114f85",
			},
		},
		{
			input: "1.0.0-alpha.beta",
			want: want{
				ver: "1.0.0",
				pre: "alpha.beta",
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"/"+tt.input, func(t *testing.T) {
			ver, pre, meta := splitSemVer(tt.input)
			if ver != tt.want.ver {
				t.Errorf("version mismatch:\n  got:  %s\n  want: %s", ver, tt.want.ver)
				return
			}
			if pre != tt.want.pre {
				t.Errorf("prerelease mismatch:\n  got:  %s\n  want: %s", pre, tt.want.pre)
				return
			}
			if meta != tt.want.meta {
				t.Errorf("buildmetadata mismatch:\n  got:  %s\n  want: %s", meta, tt.want.meta)
				return
			}
		})
	}
}

func TestIsValidSemVer(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{
			input: "",
			want:  false,
		},
		{
			input: "1",
			want:  false,
		},
		{
			input: "3.0.0",
			want:  true,
		},
		{
			input: "a.0.1",
			want:  false,
		},
		{
			input: "0.b.2",
			want:  false,
		},
		{
			input: "0.3.c",
			want:  false,
		},
		{
			input: "-1.2.3",
			want:  false,
		},
		{
			input: "1.-2.3",
			want:  false,
		},
		{
			input: "1.2.-3",
			want:  false,
		},
		{
			input: "1.9.0",
			want:  true,
		},
		{
			input: "01.9.0",
			want:  false,
		},
		{
			input: "1.09.0",
			want:  false,
		},
		{
			input: "1.9.00",
			want:  false,
		},
		{
			input: "1.0.0-alpha",
			want:  true,
		},
		{
			input: "1.0.0-alpha.1",
			want:  true,
		},
		{
			input: "1.0.0-0.3.7",
			want:  true,
		},
		{
			input: "1.0.0-x.7.z.92",
			want:  true,
		},
		{
			input: "1.0.0-alpha+001",
			want:  true,
		},
		{
			input: "1.0.0+20130313144700",
			want:  true,
		},
		{
			input: "1.0.0-beta+exp.sha.5114f85",
			want:  true,
		},
		{
			input: "1.0.0+exp.sha.5114f85-beta",
			want:  true,
		},
		{
			input: "1.0.0-",
			want:  false,
		},
		{
			input: "1.0.0-alpha..gamma",
			want:  false,
		},
		{
			input: "1.0.0+",
			want:  false,
		},
		{
			input: "1.0.0+alpha..gamma",
			want:  false,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"/"+tt.input, func(t *testing.T) {
			got := isValidSemVer(tt.input)
			if got != tt.want {
				t.Errorf("unexpected result:\n  got:  %t\n  want: %t", got, tt.want)
				return
			}
		})
	}
}

func TestIsValidPrerelease(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{
			input: "",
			want:  true,
		},
		{
			input: "alpha.beta.gamma",
			want:  true,
		},
		{
			input: "01",
			want:  false,
		},
		{
			input: "あるふぁ",
			want:  false,
		},
		{
			input: "alpha..gamma",
			want:  false,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"/"+tt.input, func(t *testing.T) {
			got := isValidPrerelease(tt.input)
			if got != tt.want {
				t.Errorf("unexpected result:\n  got:  %t\n  want: %t", got, tt.want)
				return
			}
		})
	}
}

func TestIsValidBuildmetadata(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{
			input: "",
			want:  true,
		},
		{
			input: "alpha.beta.gamma",
			want:  true,
		},
		{
			input: "びるど",
			want:  false,
		},
		{
			input: "alpha..gamma",
			want:  false,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"/"+tt.input, func(t *testing.T) {
			got := isValidBuildmetadata(tt.input)
			if got != tt.want {
				t.Errorf("unexpected result:\n  got:  %t\n  want: %t", got, tt.want)
				return
			}
		})
	}
}
