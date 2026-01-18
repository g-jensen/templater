package fs

import "os"

type FileSystem interface {
	ReadFile(path string) ([]byte, error)
	ReadDir(path string) ([]os.DirEntry, error)
	Stat(path string) (os.FileInfo, error)
}

type OSFileSystem struct{}

func (OSFileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (OSFileSystem) ReadDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

func (OSFileSystem) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}
