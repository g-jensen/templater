package template

import (
	"strings"

	"templater/internal/fs"
)

func ParseFeaturesFile(fileSystem fs.WritableFS, path string) ([]string, error) {
	data, err := fileSystem.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var features []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			features = append(features, trimmed)
		}
	}

	return features, nil
}
