package template

import (
	"os"
	"path"
	"sort"

	"gopkg.in/yaml.v3"
	"templater/internal/fs"
)

type appliedYml struct {
	Applied []string `yaml:"applied"`
}

func ReadApplied(fileSystem fs.WritableFS, targetPath string) ([]string, error) {
	data, err := fileSystem.ReadFile(path.Join(targetPath, ".templater/applied.yml"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var applied appliedYml
	if err := yaml.Unmarshal(data, &applied); err != nil {
		return nil, err
	}

	return applied.Applied, nil
}

func WriteApplied(fileSystem fs.WritableFS, targetPath string, features []string) error {
	sorted := make([]string, len(features))
	copy(sorted, features)
	sort.Strings(sorted)

	data, err := yaml.Marshal(appliedYml{Applied: sorted})
	if err != nil {
		return err
	}

	return fileSystem.WriteFile(path.Join(targetPath, ".templater/applied.yml"), data)
}
