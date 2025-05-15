// --- fileutils/fs.go ---
package fileutils

import (
	"io"
	"io/fs"
	"time"
)

// FS defines an abstract file system interface for dependency injection and testing.
type FS interface {
	Stat(name string) (fs.FileInfo, error)
	Open(name string) (io.ReadCloser, error)
	Create(name string) (io.WriteCloser, error)
	MkdirAll(path string, perm fs.FileMode) error
	Chmod(name string, mode fs.FileMode) error
	ReadDir(name string) ([]fs.DirEntry, error)
	WalkDir(root string, fn fs.WalkDirFunc) error
	RemoveAll(path string) error
	Chtimes(name string, atime, mtime time.Time) error
}
