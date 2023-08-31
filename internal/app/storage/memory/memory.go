package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/psxzz/backend-trainee-assignment/internal/app/storage"
)

var (
	segmentsIdx        int64 = 0
	userExperimentsIdx int64 = 0
)

type Storage struct {
	segments map[string]struct {
		ID   int64
		Name string
	}
	userExperiments map[int64][]struct {
		ID        int64
		UserID    int64
		SegmentID int64
	}
}

func New() *Storage {
	return &Storage{
		segments: make(map[string]struct {
			ID   int64
			Name string
		}),
		userExperiments: make(map[int64][]struct {
			ID        int64
			UserID    int64
			SegmentID int64
		}),
	}
}

func (s *Storage) Segments() map[string]struct {
	ID   int64
	Name string
} {
	return s.segments
}

func (s *Storage) Experiments() map[int64][]struct {
	ID        int64
	UserID    int64
	SegmentID int64
} {
	return s.userExperiments
}

func (s *Storage) AddSegment(ctx context.Context, name string) (*storage.SegmentDTO, error) {
	if _, ok := s.segments[name]; ok {
		return nil, fmt.Errorf("mock storage add: %w", storage.ErrSegmentExists)
	}

	s.segments[name] = struct {
		ID   int64
		Name string
	}{
		segmentsIdx,
		name,
	}

	res := &storage.SegmentDTO{
		ID:   segmentsIdx,
		Name: name,
	}
	segmentsIdx++

	return res, nil
}

func (s *Storage) DeleteSegment(ctx context.Context, name string) (*storage.SegmentDTO, error) {
	if _, ok := s.segments[name]; !ok {
		return nil, fmt.Errorf("mock storage delete: %w", storage.ErrSegmentNotFound)
	}

	res := &storage.SegmentDTO{
		ID:   s.segments[name].ID,
		Name: s.segments[name].Name,
	}
	delete(s.segments, name)

	return res, nil
}

func (s *Storage) AddUserToSegment(ctx context.Context, userID int64, segmentName string) (*storage.UserExperimentDTO, error) {
	segment, ok := s.segments[segmentName]
	if !ok {
		return nil, storage.ErrSegmentNotFound
	}

	for _, record := range s.userExperiments[userID] {
		if record.SegmentID == segment.ID {
			return nil, fmt.Errorf("mock stogare add user to segment: %w", storage.ErrAlreadyInExperiment)
		}
	}

	s.userExperiments[userID] = append(s.userExperiments[userID], struct {
		ID        int64
		UserID    int64
		SegmentID int64
	}{
		ID:        userExperimentsIdx,
		UserID:    userID,
		SegmentID: segment.ID,
	})

	res := &storage.UserExperimentDTO{
		ID:     userExperimentsIdx,
		UserID: userID,
		Segment: storage.SegmentDTO{
			ID:   segment.ID,
			Name: segmentName,
		},
	}
	userExperimentsIdx++

	return res, nil
}

func (s *Storage) DeleteUserFromSegment(ctx context.Context, userID int64, segmentName string) (*storage.UserExperimentDTO, error) {
	segment, ok := s.segments[segmentName]
	if !ok {
		return nil, storage.ErrSegmentNotFound
	}
	var (
		idx int                        = -1
		res *storage.UserExperimentDTO = &storage.UserExperimentDTO{
			Segment: storage.SegmentDTO{},
		}
	)

	if records, ok := s.userExperiments[userID]; ok {
		for i, record := range records {
			if record.SegmentID == segment.ID {
				idx = i
				res.ID = record.ID
				res.UserID = userID
				res.Segment.ID = segment.ID
				res.Segment.Name = segmentName
				break
			}
		}

		if idx == -1 {
			return nil, fmt.Errorf("mock storage delete user from segment: %w", storage.ErrUserExperimentNotFound)
		}
		s.userExperiments[userID] = append(s.userExperiments[userID][:idx],
			s.userExperiments[userID][idx+1:]...)
	}

	return res, nil
}

func (s *Storage) UserSegments(ctx context.Context, userID int64) (*storage.UserExperimentListDTO, error) {
	res := &storage.UserExperimentListDTO{
		UserID: userID,
	}

	segmentKeys := make([]string, 0, len(s.segments))
	for k := range s.segments {
		segmentKeys = append(segmentKeys, k)
	}

	for _, record := range s.userExperiments[userID] {

		for _, key := range segmentKeys {
			if s.segments[key].ID == record.SegmentID {
				res.Segments = append(res.Segments, storage.SegmentDTO{
					ID:   s.segments[key].ID,
					Name: s.segments[key].Name,
				})
				break
			}
		}
	}

	return res, nil
}

func (s *Storage) UserExperimentLogs(ctx context.Context, userID int64, start time.Time) ([]*storage.UserExperimentLogRecordDTO, error) {
	return nil, nil
}

func (s *Storage) DeleteOldExperiments(ctx context.Context) error {
	return nil
}

func (s *Storage) AddUserToSegmentWithExpiracy(ctx context.Context, userID int64, segmentName string, expiredAt time.Time) (*storage.UserExperimentDTO, error) {
	return nil, nil
}
