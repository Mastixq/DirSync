package fileutils

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type OsFS struct{}

func (OsFS) Stat(name string) (os.FileInfo, error)        { return os.Stat(name) }
func (OsFS) Open(name string) (*os.File, error)           { return os.Open(name) }
func (OsFS) Create(name string) (*os.File, error)         { return os.Create(name) }
func (OsFS) MkdirAll(path string, perm fs.FileMode) error { return os.MkdirAll(path, perm) }
func (OsFS) Chmod(name string, mode fs.FileMode) error    { return os.Chmod(name, mode) }
func (OsFS) ReadDir(name string) ([]os.DirEntry, error)   { return os.ReadDir(name) }
func (OsFS) WalkDir(root string, fn fs.WalkDirFunc) error {
	return fs.WalkDir(os.DirFS(root), ".", func(path string, d fs.DirEntry, err error) error {
		absPath := filepath.Join(root, path) // <- poprawka
		return fn(absPath, d, err)
	})
}
func (OsFS) RemoveAll(path string) error                       { return os.RemoveAll(path) }
func (OsFS) Chtimes(name string, atime, mtime time.Time) error { return os.Chtimes(name, atime, mtime) }
