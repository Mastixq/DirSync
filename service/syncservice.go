package service

import (
	"dirsync/internal/logger"
	fileutils "dirsync/internal/utils"
)

type SyncSvc struct {
	source        string
	target        string
	deleteMissing bool
	fs            fileutils.FS
}

func NewSyncSvc(source, target string, deleteMissing bool, fs fileutils.FS) *SyncSvc {
	return &SyncSvc{
		source:        source,
		target:        target,
		deleteMissing: deleteMissing,
		fs:            fs,
	}
}

func (s *SyncSvc) Sync() error {
	files, err := fileutils.ListFiles(s.fs, s.source)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := fileutils.CopyFilePreserveTree(s.fs, file, s.source, s.target)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
	}

	// opcjonalna logika dla deleteMissing...

	return nil
}
