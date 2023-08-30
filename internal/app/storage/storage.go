package storage

import (
	"errors"
	"time"
)

var (
	// TODO: declare errors
	ErrSegmentExists          = errors.New("segment with current name already exists")
	ErrSegmentNotFound        = errors.New("segment with current name not found")
	ErrAlreadyInExperiment    = errors.New("current user is already in segment")
	ErrUserExperimentNotFound = errors.New("user experiment not found")
)

type SegmentDTO struct {
	ID   int64
	Name string
}

type UserExperimentDTO struct {
	ID      int64
	UserID  int64
	Segment SegmentDTO
}

type UserExperimentListDTO struct {
	UserID   int64
	Segments []SegmentDTO
}

type UserExperimentLogRecordDTO struct {
	UserID      int64
	SegmentName string
	Operation   string
	AddedAt     time.Time
}
