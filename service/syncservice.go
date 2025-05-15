package service

import (
	"dirsync/internal/logger"
	fileutils "dirsync/internal/utils"
	"path/filepath"
)

type SyncSvc struct {
	source        string
	target        string
	deleteMissing bool
	fs            fileutils.FS
}

func NewSyncSvc(source, target string, deleteMissing bool) (*SyncSvc, error) {
	sourceAbs, err := filepath.Abs(source)
	if err != nil {
		return nil, err
	}
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return nil, err
	}
	logger.Info("source path: %s", sourceAbs)
	logger.Info("target path: %s", targetAbs)
	return &SyncSvc{
		source:        sourceAbs,
		target:        targetAbs,
		deleteMissing: deleteMissing,
		fs:            fileutils.OsFS{},
	}, nil
}

func (s *SyncSvc) Execute() error {
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
	if s.deleteMissing {
		err := fileutils.DeleteMissing(s.fs, s.source, s.target)
		if err != nil {
			logger.Error("Error removing missing files: %v", err)
		}
	}

	return nil
}
