package template

import (
	"path"
	"sort"

	"templater/internal/fs"
)

func ListFeatures(fileSystem fs.FileSystem, repoPath string) ([]string, error) {
	var features []string
	queue := []string{""}

	for len(queue) > 0 {
		relPath := queue[0]
		queue = queue[1:]

		fullPath := repoPath
		if relPath != "" {
			fullPath = path.Join(repoPath, relPath)
		}

		entries, err := fileSystem.ReadDir(fullPath)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			entryRel := entry.Name()
			if relPath != "" {
				entryRel = path.Join(relPath, entry.Name())
			}

			if entry.IsDir() {
				queue = append(queue, entryRel)
			} else if entry.Name() == "base.patch" && relPath != "" {
				features = append(features, relPath)
			}
		}
	}

	sort.Strings(features)
	return features, nil
}
