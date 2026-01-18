package template

import (
	"testing"

	"templater/internal/testutil/fs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadApplied_NoTemplaterDir(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddDir("project")

	applied, err := ReadApplied(memfs, "project")
	require.NoError(t, err)
	assert.Empty(t, applied)
}

func TestReadApplied_EmptyAppliedYml(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddDir("project")
	memfs.AddFile("project/.templater/applied.yml", []byte("applied:\n"))

	applied, err := ReadApplied(memfs, "project")
	require.NoError(t, err)
	assert.Empty(t, applied)
}

func TestReadApplied_WithFeatures(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddDir("project")
	memfs.AddFile("project/.templater/applied.yml",
		[]byte("applied:\n"+
			"  - auth\n"+
			"  - auth/oauth\n"+
			"  - auth/oauth/google"))

	applied, err := ReadApplied(memfs, "project")
	require.NoError(t, err)
	assert.Equal(t, []string{"auth", "auth/oauth", "auth/oauth/google"}, applied)
}

func TestWriteApplied(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddDir("project")

	err := WriteApplied(memfs, "project", []string{"auth", "auth/oauth"})
	require.NoError(t, err)

	data, err := memfs.ReadFile("project/.templater/applied.yml")
	require.NoError(t, err)
	assert.Equal(t,
		"applied:\n"+
			"    - auth\n"+
			"    - auth/oauth\n",
		string(data))
}

func TestWriteApplied_SortsAlphabetically(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddDir("project")

	err := WriteApplied(memfs, "project", []string{"database", "auth", "auth/oauth"})
	require.NoError(t, err)

	data, err := memfs.ReadFile("project/.templater/applied.yml")
	require.NoError(t, err)
	assert.Equal(t,
		"applied:\n"+
			"    - auth\n"+
			"    - auth/oauth\n"+
			"    - database\n",
		string(data))
}
