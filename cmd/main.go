package main

import (
	"dirsync/internal/logger"
	"dirsync/service"
	"flag"
	"fmt"
	"os"
)

func main() {
	deleteMissing := flag.Bool("delete-missing", false, "Delete files in target that are missing from source")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--delete-missing] <source> <target>\n", os.Args[0])
		os.Exit(1)
	}

	source := args[0]
	target := args[1]

	svc, err := service.NewSyncSvc(source, target, *deleteMissing)
	if err != nil {
		logger.Error("Initialization failed: %v", err)
		os.Exit(1)
	}

	if err := svc.Execute(); err != nil {
		logger.Error("Synchronization failed: %v", err)
		os.Exit(1)
	}
}
