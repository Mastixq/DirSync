package fileutils

import (
	"bytes"
	"crypto/sha256"
	"dirsync/internal/logger"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ListFiles(fsys FS, root string) ([]string, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("invalid root path: %w", err)
	}

	var files []string

	err = fsys.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Warn("Cannot access path %s: %v", path, err)
			return nil
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func CopyFilePreserveTree(fsys FS, srcFile, srcRoot, dstRoot string) error {
	srcInfo, err := fsys.Stat(srcFile)
	if err != nil || !srcInfo.Mode().IsRegular() {
		return err
	}

	relPath, err := filepath.Rel(srcRoot, srcFile)
	if err != nil {
		return err
	}
	dstPath := filepath.Join(dstRoot, relPath)
	dstDir := filepath.Dir(dstPath)

	if err := ensureDirFromSourcePerms(fsys, filepath.Dir(srcFile), dstDir); err != nil {
		return err
	}

	if _, err := fsys.Stat(dstPath); err == nil {
		same, err := filesAreEqual(fsys, srcFile, dstPath)
		if err != nil {
			return err
		}
		if same {
			logger.Info("Skipping identical file: %s", dstPath)
			return nil
		}
		logger.Info("Overwriting changed file: %s", dstPath)
	}

	if err := copyFileContents(fsys, srcFile, dstPath, srcInfo.Mode()); err != nil {
		return err
	}

	if err := preserveModTime(fsys, dstPath, srcInfo.ModTime()); err != nil {
		return err
	}

	logger.Info("Copied: %s → %s", srcFile, dstPath)
	return nil
}

func ensureDirFromSourcePerms(fsys FS, srcDir, dstDir string) error {
	srcInfo, err := fsys.Stat(srcDir)
	if err != nil {
		return err
	}
	if err := fsys.MkdirAll(dstDir, srcInfo.Mode()); err != nil {
		return err
	}
	return fsys.Chmod(dstDir, srcInfo.Mode())
}

func filesAreEqual(fsys FS, pathA, pathB string) (bool, error) {
	infoA, err := fsys.Stat(pathA)
	if err != nil {
		return false, err
	}
	infoB, err := fsys.Stat(pathB)
	if err != nil {
		return false, err
	}

	if !infoA.ModTime().Equal(infoB.ModTime()) {
		return false, nil
	}

	hashA, err := hashFile(fsys, pathA)
	if err != nil {
		return false, err
	}
	hashB, err := hashFile(fsys, pathB)
	if err != nil {
		return false, err
	}
	return bytes.Equal(hashA, hashB), nil
}

func hashFile(fsys FS, path string) ([]byte, error) {
	f, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil

}

func copyFileContents(fsys FS, src, dst string, perm fs.FileMode) error {
	in, err := fsys.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := fsys.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return fsys.Chmod(dst, perm)
}

func preserveModTime(fsys FS, path string, modTime time.Time) error {
	return fsys.Chtimes(path, modTime, modTime)
}

func DeleteMissing(fsys FS, srcRoot, dstRoot string) error {
	absSrcRoot, err := filepath.Abs(srcRoot)
	if err != nil {
		return fmt.Errorf("cannot resolve absolute source path: %w", err)
	}
	absDstRoot, err := filepath.Abs(dstRoot)
	if err != nil {
		return fmt.Errorf("cannot resolve absolute target path: %w", err)
	}

	return fsys.WalkDir(absDstRoot, func(dstPath string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Warn("Failed to access: %s: %v", dstPath, err)
			return nil
		}

		absDstPath, err := filepath.Abs(dstPath)
		if err != nil {
			logger.Warn("Invalid path: %s", dstPath)
			return nil
		}

		relPath, err := filepath.Rel(absDstRoot, absDstPath)
		if err != nil || isOutsideBase(relPath) {
			logger.Warn("Skipping path outside target: %s", dstPath)
			return nil
		}

		srcPath := filepath.Join(absSrcRoot, relPath)

		if _, err := fsys.Stat(srcPath); isNotExist(err) {
			logger.Info("Removing missing: %s", dstPath)
			if removeErr := fsys.RemoveAll(dstPath); removeErr != nil {
				logger.Error("Failed to remove %s: %v", dstPath, removeErr)
			}
		}

		return nil
	})
}
func isNotExist(err error) bool {
	return err != nil && (fs.ErrNotExist == err || os.IsNotExist(err))
}

func isOutsideBase(rel string) bool {
	clean := filepath.Clean(rel)
	return clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator))
}
