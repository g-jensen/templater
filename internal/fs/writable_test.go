package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOSWritableFS_WriteFile(t *testing.T) {
	tempDir := t.TempDir()
	fs := OSWritableFS{BaseDir: tempDir}

	err := fs.WriteFile("test.txt", []byte("hello"))
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(tempDir, "test.txt"))
	require.NoError(t, err)
	assert.Equal(t, "hello", string(content))
}

func TestOSWritableFS_AppendFile(t *testing.T) {
	tempDir := t.TempDir()
	fs := OSWritableFS{BaseDir: tempDir}

	err := fs.WriteFile("test.txt", []byte("hello"))
	require.NoError(t, err)

	err = fs.AppendFile("test.txt", []byte(" world"))
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(tempDir, "test.txt"))
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(content))
}
