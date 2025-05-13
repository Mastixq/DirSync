package fileutils

import (
	"fmt"
	"io"
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
	// Otwórz źródłowy plik do odczytu
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("nie mogę otworzyć pliku źródłowego: %w", err)
	}
	defer in.Close()

	fileName := filepath.Base(src)

	// Utwórz plik docelowy do zapisu (nadpisze jeśli istnieje)
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

	// Kopiowanie zawartości
	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("błąd kopiowania danych: %w", err)
	}

	// Zatwierdzenie zawartości na dysku
	err = out.Sync()
	if err != nil {
		return fmt.Errorf("błąd zapisu na dysk: %w", err)
	}

	return nil
}
