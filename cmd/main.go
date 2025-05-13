package main

import (
	fileutils "dirsync/internal/utils"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Printf("Hello %s %s", os.Args[1], os.Args[2])
		for _, file := range fileutils.ListFiles(os.Args[1]) {
			fmt.Println(file)
			err := fileutils.CopyFileToTargetDir(file, os.Args[2])
			if err != nil {
				fmt.Println(err)
			}
		}

	} else {
		fmt.Println("Hello world")
	}
}
