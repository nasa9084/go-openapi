package openapi

import (
	"strconv"
	"strings"
)

const alnum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-"

func splitSemVer(version string) (ver, pre, meta string) {
	hyphen := strings.Index(version, "-")
	plus := strings.Index(version, "+")
	rversion := []rune(version)

	if hyphen > 0 { // prerelease is defined
		if plus > 0 { // buildmetadata is defined
			if hyphen < plus {
				return string(rversion[:hyphen]), string(rversion[hyphen+1 : plus]), string(rversion[plus+1:])
			}

			return string(rversion[:plus]), string(rversion[hyphen+1:]), string(rversion[plus+1 : hyphen])
		}

		return string(rversion[:hyphen]), string(rversion[hyphen+1:]), ""
	} else if plus > 0 { // only buildmetadata
		return string(rversion[:plus]), "", string(rversion[plus+1:])
	}

	// both prerelease and buildmetadata are not defined
	return version, "", ""
}

func isValidSemVer(version string) bool {
	var ver, pre, meta string
	ver, pre, meta = splitSemVer(version)

	split := strings.Split(ver, ".")
	if len(split) != 3 {
		return false
	}

	major, minor, patch := split[0], split[1], split[2]

	// version number must not contain leading zero
	if !hasLeadingZero(major, minor, patch) {
		return false
	}

	// version number must be non-negative integers
	if major, err := strconv.Atoi(major); err != nil {
		return false
	} else if major < 0 {
		return false
	}

	// minor and patch never be negative: if there's "-", it is parsed as prerelease
	if _, err := strconv.Atoi(minor); err != nil {
		return false
	}

	if _, err := strconv.Atoi(patch); err != nil {
		return false
	}

	// prerelease version is identifiers splitted by dot.
	// identifiers must be composed by ASCCI alpha-numberics and hyphens.
	// identifiers must not be empty, and must not begin with leading zeroes.
	if strings.Contains(version, "-") && pre == "" {
		return false
	}

	if !isValidPrerelease(pre) {
		return false
	}

	// buildmetadata is identifiers splitted by dot.
	// identifiers must be composed by ASSCII alpha-numerics and hyphens.
	// identifiers must not be empty.
	if strings.Contains(version, "+") && meta == "" {
		return false
	}

	if !isValidBuildmetadata(meta) {
		return false
	}

	return true
}

func hasLeadingZero(major, minor, patch string) bool {
	if major != "0" && strings.HasPrefix(major, "0") {
		return true
	}

	if minor != "0" && strings.HasPrefix(minor, "0") {
		return true
	}

	if patch != "0" && strings.HasPrefix(patch, "0") {
		return true
	}

	return false
}

func isValidPrerelease(pre string) bool {
	if pre == "" {
		return true
	}

	for _, ident := range strings.Split(pre, ".") {
		if ident == "" {
			return false
		}

		if _, err := strconv.Atoi(ident); err == nil {
			if ident != "0" && strings.HasPrefix(ident, "0") {
				return false
			}
		}

		for _, c := range ident {
			if !strings.ContainsRune(alnum, c) {
				return false
			}
		}
	}

	return true
}

func isValidBuildmetadata(meta string) bool {
	if meta == "" {
		return true
	}

	for _, ident := range strings.Split(meta, ".") {
		if ident == "" {
			return false
		}

		for _, c := range ident {
			if !strings.ContainsRune(alnum, c) {
				return false
			}
		}
	}

	return true
}
