package fileutils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func ListFiles(dirPath string) []string {
	var files []string
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println(fmt.Errorf("nie można odczytać katalogu %s: %w", dirPath, err))
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		files = append(files, fullPath)
	}
	return files
}

func CopyFileToTargetDir(src, targetDir string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("nie mogę otworzyć pliku źródłowego: %w", err)
	}
	defer func() {
		if err := in.Close(); err != nil {
			log.Printf("Błąd przy zamykaniu pliku: %v", err)
		}
	}()

	fileName := filepath.Base(src)

	out, err := os.Create(fmt.Sprintf("%s/%s", targetDir, fileName))
	if err != nil {
		return fmt.Errorf("nie mogę utworzyć pliku docelowego: %w", err)
	}
	defer func() {
		cerr := out.Close()
		if cerr != nil {
			fmt.Printf("Błąd przy zamykaniu pliku docelowego: %v\n", cerr)
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("błąd kopiowania danych: %w", err)
	}

	err = out.Sync()
	if err != nil {
		return fmt.Errorf("błąd zapisu na dysk: %w", err)
	}

	return nil
}
