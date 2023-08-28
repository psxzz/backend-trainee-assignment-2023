package storage

import "errors"

var (
	// TODO: declare errors
	ErrSegmentExists          = errors.New("segment with current name already exists")
	ErrSegmentNotFound        = errors.New("segment with current name not found")
	ErrAlreadyInExperiment    = errors.New("current user is already in segment")
	ErrUserExperimentNotFound = errors.New("user experiment not found")
)
