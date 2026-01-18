package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"templater/internal/testutil/fs"
)

func TestParseFeaturesFile(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddFile("features.txt", []byte("auth/oauth/google\nauth/oauth/github\ndatabase/migrations\n"))

	features, err := ParseFeaturesFile(memfs, "features.txt")
	require.NoError(t, err)

	assert.Equal(t, []string{"auth/oauth/google", "auth/oauth/github", "database/migrations"}, features)
}

func TestParseFeaturesFile_IgnoresEmptyLines(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddFile("features.txt", []byte("auth\n\ndatabase\n\n"))

	features, err := ParseFeaturesFile(memfs, "features.txt")
	require.NoError(t, err)

	assert.Equal(t, []string{"auth", "database"}, features)
}

func TestParseFeaturesFile_TrimsWhitespace(t *testing.T) {
	memfs := fs.NewMemoryFS()
	memfs.AddFile("features.txt", []byte("  auth  \n  database  \n"))

	features, err := ParseFeaturesFile(memfs, "features.txt")
	require.NoError(t, err)

	assert.Equal(t, []string{"auth", "database"}, features)
}
