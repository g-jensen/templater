package template

import (
	"testing"

	"templater/internal/testutil/fs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListFeatures(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddDir("repo")
	memfs.AddDir("repo/auth")
	memfs.AddDir("repo/auth/oauth")
	memfs.AddDir("repo/auth/oauth/google")
	memfs.AddDir("repo/db")
	memfs.AddFile("repo/auth/base.patch", []byte{})
	memfs.AddFile("repo/auth/oauth/base.patch", []byte{})
	memfs.AddFile("repo/auth/oauth/google/base.patch", []byte{})
	memfs.AddFile("repo/db/base.patch", []byte{})

	features, err := ListFeatures(memfs, "repo")
	require.NoError(t, err)

	expected := []string{"auth", "auth/oauth", "auth/oauth/google", "db"}
	assert.Equal(t, expected, features)
}

func TestListFeatures_ExcludesRootBasePatch(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddDir("repo")
	memfs.AddFile("repo/base.patch", []byte{})

	features, err := ListFeatures(memfs, "repo")
	require.NoError(t, err)

	assert.Empty(t, features)
}

func TestListFeatures_MissingBasePatch(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddDir("repo")
	memfs.AddDir("repo/auth")
	memfs.AddDir("repo/auth/oauth")
	memfs.AddDir("repo/auth/oauth/google")
	memfs.AddDir("repo/auth/oauth/github")
	memfs.AddFile("repo/auth/base.patch", []byte{})
	memfs.AddFile("repo/auth/oauth/google/base.patch", []byte{})
	memfs.AddFile("repo/auth/oauth/github/base.patch", []byte{})

	features, err := ListFeatures(memfs, "repo")
	require.NoError(t, err)

	expected := []string{"auth", "auth/oauth/github", "auth/oauth/google"}
	assert.Equal(t, expected, features)
}
