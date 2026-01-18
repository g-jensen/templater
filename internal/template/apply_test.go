package template

import (
	"testing"

	"templater/internal/testutil/executor"
	"templater/internal/testutil/fs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyFeature_ExecutesGitApply(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddDir("templates")
	memfs.AddDir("templates/auth")
	memfs.AddFile("templates/auth/base.patch", []byte("patch content"))
	memfs.AddDir("project")

	exec := &executor.FakeExecutor{}

	err := ApplyFeature(memfs, exec, "templates", "project", "auth")
	require.NoError(t, err)

	require.Len(t, exec.Commands, 1)
	expected_command := "git apply --directory=project templates/auth/base.patch"
	assert.Equal(t, exec.Commands[0].Command, expected_command)
}

func TestApplyFeature_ReturnsErrorOnFailure(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddDir("templates")
	memfs.AddDir("templates/auth")
	memfs.AddFile("templates/auth/base.patch", []byte("patch content"))
	memfs.AddDir("project")

	exec := &executor.FakeExecutor{
		DefaultExitCode: 1,
		Stderr:          "patch does not apply",
	}

	err := ApplyFeature(memfs, exec, "templates", "project", "auth")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "patch does not apply")
}

func TestApplyFeatures_AppliesDependenciesInOrder(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddDir("templates")
	memfs.AddDir("templates/auth")
	memfs.AddDir("templates/auth/oauth")
	memfs.AddFile("templates/auth/base.patch", []byte("auth patch"))
	memfs.AddFile("templates/auth/oauth/base.patch", []byte("oauth patch"))
	memfs.AddDir("project")

	exec := &executor.FakeExecutor{}

	result, err := ApplyFeatures(memfs, exec, "templates", "project", []string{"auth/oauth"})
	require.NoError(t, err)

	assert.Equal(t, []string{"auth", "auth/oauth"}, result.Applied)
	require.Len(t, exec.Commands, 2)
	assert.Contains(t, exec.Commands[0].Command, "auth/base.patch")
	assert.Contains(t, exec.Commands[1].Command, "auth/oauth/base.patch")
}

func TestApplyFeatures_SkipsAlreadyApplied(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddDir("templates")
	memfs.AddDir("templates/auth")
	memfs.AddDir("templates/auth/oauth")
	memfs.AddFile("templates/auth/base.patch", []byte("auth patch"))
	memfs.AddFile("templates/auth/oauth/base.patch", []byte("oauth patch"))
	memfs.AddDir("project")
	memfs.AddFile("project/.templater/applied.yml", []byte("applied:\n  - auth\n"))

	exec := &executor.FakeExecutor{}

	result, err := ApplyFeatures(memfs, exec, "templates", "project", []string{"auth/oauth"})
	require.NoError(t, err)

	assert.Equal(t, []string{"auth/oauth"}, result.Applied)
	assert.Equal(t, []string{"auth"}, result.AlreadyApplied)
	require.Len(t, exec.Commands, 1)
}

func TestApplyFeatures_RollsBackOnFailure(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddDir("templates")
	memfs.AddDir("templates/auth")
	memfs.AddDir("templates/auth/oauth")
	memfs.AddFile("templates/auth/base.patch", []byte("auth patch"))
	memfs.AddFile("templates/auth/oauth/base.patch", []byte("oauth patch"))
	memfs.AddDir("project")

	exec := &executor.FakeExecutor{
		ExitCodes: map[string]int{
			"git apply --directory=project templates/auth/oauth/base.patch": 1,
		},
		Stderr: "patch does not apply",
	}

	_, err := ApplyFeatures(memfs, exec, "templates", "project", []string{"auth/oauth"})
	require.Error(t, err)

	assert.Equal(t, 3, len(exec.Commands))
	lastCmd := exec.Commands[len(exec.Commands)-1]
	expectedCommand := "git apply --reverse --directory=project templates/auth/base.patch"
	assert.Equal(t, expectedCommand, lastCmd.Command)
}

func TestApplyFeatures_RollsBackMultipleOnFailure(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddDir("templates")
	memfs.AddDir("templates/auth")
	memfs.AddDir("templates/auth/oauth")
	memfs.AddFile("templates/base.patch", []byte("base patch"))
	memfs.AddFile("templates/auth/base.patch", []byte("auth patch"))
	memfs.AddFile("templates/auth/oauth/base.patch", []byte("oauth patch"))
	memfs.AddDir("project")

	exec := &executor.FakeExecutor{
		ExitCodes: map[string]int{
			"git apply --directory=project templates/auth/oauth/base.patch": 1,
		},
		Stderr: "patch does not apply",
	}

	_, err := ApplyFeatures(memfs, exec, "templates", "project", []string{"auth/oauth"})
	require.Error(t, err)

	assert.Equal(t, 5, len(exec.Commands))

	baseCommand := exec.Commands[len(exec.Commands)-1]
	expectedBaseCommand := "git apply --reverse --directory=project templates/base.patch"
	authCommand := exec.Commands[len(exec.Commands)-2]
	expectedAuthCommand := "git apply --reverse --directory=project templates/auth/base.patch"
	assert.Equal(t, expectedBaseCommand, baseCommand.Command)
	assert.Equal(t, expectedAuthCommand, authCommand.Command)
}

func TestDryRun_ReturnsWhatWouldBeApplied(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddDir("templates")
	memfs.AddDir("templates/auth")
	memfs.AddDir("templates/auth/oauth")
	memfs.AddFile("templates/auth/base.patch", []byte("auth patch"))
	memfs.AddFile("templates/auth/oauth/base.patch", []byte("oauth patch"))
	memfs.AddDir("project")

	result, err := DryRun(memfs, "templates", "project", []string{"auth/oauth"})
	require.NoError(t, err)

	assert.Equal(t, []string{"auth", "auth/oauth"}, result.WouldApply)
	assert.Empty(t, result.AlreadyApplied)
}

func TestDryRun_ExcludesAlreadyApplied(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddDir("templates")
	memfs.AddDir("templates/auth")
	memfs.AddDir("templates/auth/oauth")
	memfs.AddFile("templates/auth/base.patch", []byte("auth patch"))
	memfs.AddFile("templates/auth/oauth/base.patch", []byte("oauth patch"))
	memfs.AddDir("project")
	memfs.AddFile("project/.templater/applied.yml", []byte("applied:\n  - auth\n"))

	result, err := DryRun(memfs, "templates", "project", []string{"auth/oauth"})
	require.NoError(t, err)

	assert.Equal(t, []string{"auth/oauth"}, result.WouldApply)
	assert.Equal(t, []string{"auth"}, result.AlreadyApplied)
}
