package atomic

import (
	"strings"

	"github.com/gobwas/glob"
	"golang.org/x/exp/slices"
)

func StringInSliceCaseInsensitive(value string, allowlist []string) bool {
	value = strings.ToLower(value)
	for i, v := range allowlist {
		allowlist[i] = v
	}
	return slices.Contains(allowlist, value)
}

func StringMatchesCaseInsensitivePattern(value, pattern string) bool {
	value = strings.ToLower(value)
	pattern = strings.ToLower(pattern)
	return glob.MustCompile(pattern).Match(value)
}

func StringMatchesAnyCaseInsensitivePattern(value string, patterns []string) bool {
	value = strings.ToLower(value)
	for _, pattern := range patterns {
		pattern = strings.ToLower(pattern)
		g := glob.MustCompile(pattern)
		if g.Match(value) {
			return true
		}
	}
	return false
}

func AnyStringMatchesAnyCaseInsensitivePattern(values, patterns []string) bool {
	for _, value := range values {
		for _, pattern := range patterns {
			if StringMatchesCaseInsensitivePattern(value, pattern) {
				return true
			}
		}
	}
	return false
}
