package template

import (
	"slices"
	"strings"
)

func ResolveDependencies(feature string, available []string, hasRoot bool) []string {
	var result []string

	if hasRoot {
		result = append(result, "")
	}

	parts := strings.Split(feature, "/")
	for i := range parts {
		ancestor := strings.Join(parts[:i+1], "/")
		if slices.Contains(available, ancestor) {
			result = append(result, ancestor)
		}
	}

	return result
}
