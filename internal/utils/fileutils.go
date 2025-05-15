package fileutils

import (
	"dirsync/internal/logger"
	"io"
	"io/fs"

	"path/filepath"
)

func ListFiles(fsys FS, root string) ([]string, error) {
	var files []string

	err := fsys.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Warn("Brak dostępu do %s: %v", path, err)
			return nil // pomijamy błędne pliki
		}
		if !d.IsDir() {
			files = append(files, filepath.Join(root, path))
		}
		return nil
	})

	return files, err
}

func CopyFilePreserveTree(fsys FS, srcFile, srcRoot, dstRoot string) error {
	srcInfo, err := fsys.Stat(srcFile)
	if err != nil {
		return err
	}
	if !srcInfo.Mode().IsRegular() {
		return nil
	}

	relPath, err := filepath.Rel(srcRoot, srcFile)
	if err != nil {
		return err
	}
	dstPath := filepath.Join(dstRoot, relPath)
	dstDir := filepath.Dir(dstPath)

	srcDir := filepath.Dir(srcFile)
	srcDirInfo, err := fsys.Stat(srcDir)
	if err != nil {
		return err
	}

	if err := fsys.MkdirAll(dstDir, srcDirInfo.Mode()); err != nil {
		return err
	}
	if err := fsys.Chmod(dstDir, srcDirInfo.Mode()); err != nil {
		return err
	}

	in, err := fsys.Open(srcFile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := fsys.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	if err := fsys.Chmod(dstPath, srcInfo.Mode()); err != nil {
		return err
	}

	logger.Info("Copied: %s → %s", srcFile, dstPath)
	return nil
}
