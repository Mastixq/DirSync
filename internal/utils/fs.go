package fileutils

import (
	"io/fs"
	"os"
)

type FS interface {
	Stat(name string) (os.FileInfo, error)
	Open(name string) (*os.File, error)
	Create(name string) (*os.File, error)
	MkdirAll(path string, perm fs.FileMode) error
	Chmod(name string, mode fs.FileMode) error
	ReadDir(name string) ([]os.DirEntry, error)
	WalkDir(root string, fn fs.WalkDirFunc) error
}
