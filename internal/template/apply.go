package template

import (
	"fmt"
	"path"

	"templater/internal/executor"
	"templater/internal/fs"
)

type ApplyResult struct {
	Applied        []string
	AlreadyApplied []string
}

func ApplyFeature(fileSystem fs.FileSystem, exec executor.Executor, templatePath, targetPath, feature string) error {
	patchPath := path.Join(templatePath, feature, "base.patch")
	cmd := fmt.Sprintf("git apply --directory=%s %s", targetPath, patchPath)

	_, stderr, exitCode, err := exec.Execute(cmd, "30s", nil)
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("failed to apply %s: %s", feature, stderr)
	}

	return nil
}

func ApplyFeatures(fileSystem fs.FileSystem, exec executor.Executor, templatePath, targetPath string, features []string) (*ApplyResult, error) {
	available, err := ListFeatures(fileSystem, templatePath)
	if err != nil {
		return nil, err
	}

	alreadyApplied, err := ReadApplied(fileSystem.(fs.WritableFS), targetPath)
	if err != nil {
		return nil, err
	}
	alreadySet := toSet(alreadyApplied)

	hasRoot := hasRootPatch(fileSystem, templatePath)

	result := &ApplyResult{}
	var toApply []string
	seen := make(map[string]bool)

	for _, feature := range features {
		deps := ResolveDependencies(feature, available, hasRoot)
		for _, dep := range deps {
			if seen[dep] {
				continue
			}
			seen[dep] = true
			if alreadySet[dep] {
				result.AlreadyApplied = append(result.AlreadyApplied, dep)
			} else {
				toApply = append(toApply, dep)
			}
		}
	}

	var applied []string
	for _, feature := range toApply {
		if err := ApplyFeature(fileSystem, exec, templatePath, targetPath, feature); err != nil {
			rollback(exec, templatePath, targetPath, applied)
			return nil, err
		}
		applied = append(applied, feature)
	}

	result.Applied = applied
	return result, nil
}

func rollback(exec executor.Executor, templatePath, targetPath string, applied []string) {
	for i := len(applied) - 1; i >= 0; i-- {
		reverseFeature(exec, templatePath, targetPath, applied[i])
	}
}

func reverseFeature(exec executor.Executor, templatePath, targetPath, feature string) {
	patchPath := path.Join(templatePath, feature, "base.patch")
	cmd := fmt.Sprintf("git apply --reverse --directory=%s %s", targetPath, patchPath)
	exec.Execute(cmd, "30s", nil)
}

func hasRootPatch(fileSystem fs.FileSystem, templatePath string) bool {
	_, err := fileSystem.Stat(path.Join(templatePath, "base.patch"))
	return err == nil
}
