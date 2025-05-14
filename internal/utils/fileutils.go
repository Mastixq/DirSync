package fileutils

import (
	"dirsync/internal/logger"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ListFiles(srcRoot string) ([]string, error) {
	var files []string

	err := filepath.Walk(srcRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error("dupadupad134]n", err.Error())
		}
		fmt.Println("its a loop")
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func CopyFilePreserveTree(srcFilePath, srcRoot, dstRoot string) error {
	srcInfo, err := getFileInfo(srcFilePath)
	if err != nil {
		return err
	}

	if !isRegularFile(srcInfo) {
		return nil
	}

	relPath, err := getRelativePath(srcFilePath, srcRoot)
	if err != nil {
		return err
	}
	dstPath := filepath.Join(dstRoot, relPath)

	if err := createDirForFile(dstPath); err != nil {
		return err
	}

	if err := copyFile(srcFilePath, dstPath); err != nil {
		return err
	}

	if err := setFilePermissions(srcInfo, dstPath); err != nil {
		return err
	}

	logger.Info("Copied: %s → %s", srcFilePath, dstPath)
	return nil
}

func getFileInfo(srcFilePath string) (os.FileInfo, error) {
	srcInfo, err := os.Stat(srcFilePath)
	if err != nil {
		return nil, fmt.Errorf("can't access source file %s: %v", srcFilePath, err)
	}
	return srcInfo, nil
}

func isRegularFile(srcInfo os.FileInfo) bool {
	return srcInfo.Mode().IsRegular()
}

func getRelativePath(srcFilePath, srcRoot string) (string, error) {
	relPath, err := filepath.Rel(srcRoot, srcFilePath)
	if err != nil {
		return "", fmt.Errorf("error during getting relative path: %v", err)
	}
	return relPath, nil
}

func createDirForFile(dstPath string) error {
	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("can't create dir %s: %v", dstDir, err)
	}
	return nil
}

func copyFile(srcFilePath, dstPath string) error {
	in, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := in.Close(); cerr != nil {
			logger.Error("Error during closing file %s: %v", srcFilePath, cerr)
		}
	}()

	out, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("can't create target file %s: %v", dstPath, err)
	}
	defer func() {
		if cerr := out.Close(); cerr != nil {
			logger.Error("can't close %s: %v", dstPath, cerr)
		}
	}()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("error during copying %s → %s: %v", srcFilePath, dstPath, err)
	}
	return nil
}

func setFilePermissions(srcInfo os.FileInfo, dstPath string) error {
	if err := os.Chmod(dstPath, srcInfo.Mode()); err != nil {
		return fmt.Errorf("can't change permissions %s: %v", dstPath, err)
	}
	return nil
}

func IsDir(name string) (bool, error) {
	fileInfo, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}
