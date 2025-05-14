package service

import fileutils "dirsync/internal/utils"

type SyncSvc struct {
	source        string
	target        string
	deleteMissing bool
}

func newSyncSvc(source, target string, deleteMissing bool) SyncSvc {
	return SyncSvc{
		source:        source,
		target:        target,
		deleteMissing: deleteMissing,
	}
}

func (s *SyncSvc) Execute() error {
	// source is directory
	// target is directory, if doesnt exists create one
	sourceDirfiles := fileutils.ListFiles(s.source)

	for _, file := range sourceDirfiles {
		fileutils.CopyFileToTargetDir(file, s.target)
	}

	targetDirFiles := fileutils.ListFiles(s.target)

	for _, file := targetDirFiles
}


