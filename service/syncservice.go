package service

import (
	"dirsync/internal/logger"
	fileutils "dirsync/internal/utils"
	"fmt"
	"path/filepath"
)

type SyncSvc struct {
	source        string
	target        string
	deleteMissing bool
	fs            fileutils.FS
}

func NewSyncSvc(source, target string, deleteMissing bool) (*SyncSvc, error) {
	return NewSyncSvcWithFS(source, target, deleteMissing, fileutils.OsFS{})
}

func NewSyncSvcWithFS(source, target string, deleteMissing bool, fs fileutils.FS) (*SyncSvc, error) {
	sourceAbs, err := filepath.Abs(source)
	if err != nil {
		return nil, fmt.Errorf("invalid source path: %w", err)
	}
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return nil, fmt.Errorf("invalid target path: %w", err)
	}

	logger.Info("source path: %s", sourceAbs)
	logger.Info("target path: %s", targetAbs)

	return &SyncSvc{
		source:        sourceAbs,
		target:        targetAbs,
		deleteMissing: deleteMissing,
		fs:            fs,
	}, nil
}

func (s *SyncSvc) Execute() error {
	if err := s.SyncFiles(); err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}
	if s.deleteMissing {
		if err := s.CleanMissingFiles(); err != nil {
			logger.Warn("Cleanup warning: %v", err)
		}
	}
	return nil
}

func (s *SyncSvc) SyncFiles() error {
	files, err := fileutils.ListFiles(s.fs, s.source)
	if err != nil {
		return fmt.Errorf("cannot list source files: %w", err)
	}

	for _, file := range files {
		if err := fileutils.CopyFilePreserveTree(s.fs, file, s.source, s.target); err != nil {
			logger.Error("Copy error for %s: %v", file, err)
			continue
		}
	}
	return nil
}

func (s *SyncSvc) CleanMissingFiles() error {
	if err := fileutils.DeleteMissing(s.fs, s.source, s.target); err != nil {
		return fmt.Errorf("failed to delete missing files: %w", err)
	}
	return nil
}
