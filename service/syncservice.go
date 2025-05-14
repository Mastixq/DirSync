package service

import (
	"dirsync/internal/logger"
	fileutils "dirsync/internal/utils"
	"fmt"
)

type SyncSvc struct {
	source        string
	target        string
	deleteMissing bool
}

func NewSyncSvc(source, target string, deleteMissing bool) SyncSvc {
	return SyncSvc{
		source:        source,
		target:        target,
		deleteMissing: deleteMissing,
	}
}

func (s *SyncSvc) Execute() error {
	// source is directory
	if isDir, err := fileutils.IsDir(s.source); !isDir || err != nil {
		fmt.Println("" +
			"\nhero")
		return err
	}

	sourceDirfiles, err := fileutils.ListFiles(s.source)
	if err != nil {
		fmt.Println("herherhere\n\n")
		logger.Error(err.Error())
		return err
	}

	for _, file := range sourceDirfiles {
		fmt.Println("\n", file)
		err := fileutils.CopyFilePreserveTree(file, s.source, s.target)
		if err != nil {
			fmt.Println("dupa")
			logger.Error(err.Error())
		}
	}

	//cleanup() if flaga
	return nil
}
