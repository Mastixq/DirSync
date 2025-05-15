package fileutils

import (
	"bytes"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

type fakeFile struct {
	name    string
	data    *bytes.Buffer
	mode    fs.FileMode
	modTime time.Time
}

func (f *fakeFile) Read(p []byte) (int, error)  { return f.data.Read(p) }
func (f *fakeFile) Write(p []byte) (int, error) { return f.data.Write(p) }
func (f *fakeFile) Close() error                { return nil }

type FakeFS struct {
	files map[string]*fakeFile
	dirs  map[string]fs.FileMode
}

func NewFakeFS() *FakeFS {
	return &FakeFS{
		files: make(map[string]*fakeFile),
		dirs:  make(map[string]fs.FileMode),
	}
}

func (f *FakeFS) Stat(name string) (fs.FileInfo, error) {
	if file, ok := f.files[name]; ok {
		return fakeFileInfo{file.name, file.mode, file.modTime}, nil
	}
	if mode, ok := f.dirs[name]; ok {
		return fakeFileInfo{name, mode | fs.ModeDir, time.Now()}, nil
	}
	return nil, fs.ErrNotExist
}

func (f *FakeFS) Open(name string) (io.ReadCloser, error) {
	file, ok := f.files[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(file.data.Bytes())), nil
}

func (f *FakeFS) Create(name string) (io.WriteCloser, error) {
	buf := &bytes.Buffer{}
	f.files[name] = &fakeFile{
		name:    filepath.Base(name),
		data:    buf,
		mode:    0644,
		modTime: time.Now(),
	}
	return nopWriteCloser{Buffer: buf}, nil
}

type nopWriteCloser struct {
	*bytes.Buffer
}

func (n nopWriteCloser) Close() error { return nil }

func (f *FakeFS) MkdirAll(path string, perm fs.FileMode) error {
	f.dirs[path] = perm
	return nil
}

func (f *FakeFS) Chmod(name string, mode fs.FileMode) error {
	if file, ok := f.files[name]; ok {
		file.mode = mode
		return nil
	}
	if _, ok := f.dirs[name]; ok {
		f.dirs[name] = mode
		return nil
	}
	return fs.ErrNotExist
}

func (f *FakeFS) ReadDir(name string) ([]fs.DirEntry, error) {
	var entries []fs.DirEntry
	for path := range f.files {
		if strings.HasPrefix(path, name+"/") {
			rel := strings.TrimPrefix(path, name+"/")
			if !strings.Contains(rel, "/") {
				entries = append(entries, fakeDirEntry{name: rel})
			}
		}
	}
	return entries, nil
}

func (f *FakeFS) WalkDir(root string, fn fs.WalkDirFunc) error {
	for path := range f.dirs {
		if strings.HasPrefix(path, root) {
			de := fakeDirEntry{name: filepath.Base(path)}
			if err := fn(path, de, nil); err != nil {
				return err
			}
		}
	}
	for path := range f.files {
		if strings.HasPrefix(path, root) {
			de := fakeDirEntry{name: filepath.Base(path)}
			if err := fn(path, de, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *FakeFS) RemoveAll(path string) error {
	delete(f.files, path)
	delete(f.dirs, path)
	return nil
}

func (f *FakeFS) Chtimes(name string, atime, mtime time.Time) error {
	if file, ok := f.files[name]; ok {
		file.modTime = mtime
		return nil
	}
	return fs.ErrNotExist
}

type fakeFileInfo struct {
	name    string
	mode    fs.FileMode
	modTime time.Time
}

func (f fakeFileInfo) Name() string       { return f.name }
func (f fakeFileInfo) Size() int64        { return 0 }
func (f fakeFileInfo) Mode() fs.FileMode  { return f.mode }
func (f fakeFileInfo) ModTime() time.Time { return f.modTime }
func (f fakeFileInfo) IsDir() bool        { return f.mode.IsDir() }
func (f fakeFileInfo) Sys() interface{}   { return nil }

type fakeDirEntry struct {
	name string
}

func (f fakeDirEntry) Name() string               { return f.name }
func (f fakeDirEntry) IsDir() bool                { return false }
func (f fakeDirEntry) Type() fs.FileMode          { return 0 }
func (f fakeDirEntry) Info() (fs.FileInfo, error) { return nil, nil }
