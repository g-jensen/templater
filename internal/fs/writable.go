package fs

import (
	"os"
	"path/filepath"
)

type OSWritableFS struct {
	BaseDir string
}

func NewOSWritableFS(baseDir string) *OSWritableFS {
	return &OSWritableFS{BaseDir: baseDir}
}

func (fs *OSWritableFS) WriteFile(path string, data []byte) error {
	fullPath := filepath.Join(fs.BaseDir, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(fullPath, data, 0644)
}

func (fs *OSWritableFS) AppendFile(path string, data []byte) error {
	fullPath := filepath.Join(fs.BaseDir, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

func (fs *OSWritableFS) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(filepath.Join(fs.BaseDir, path))
}
