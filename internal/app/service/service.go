package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/psxzz/backend-trainee-assignment/internal/app/model"
	"github.com/psxzz/backend-trainee-assignment/internal/app/storage"
)

const (
	logDateFormat       = "2006-01"
	logFilenameTemplate = "log_user_%d_%v.csv"
)

type Storage interface {
	AddSegment(context.Context, string) (*storage.SegmentDTO, error)
	DeleteSegment(context.Context, string) (*storage.SegmentDTO, error)
	AddUserToSegment(context.Context, int64, string) (*storage.UserExperimentDTO, error)
	DeleteUserFromSegment(context.Context, int64, string) (*storage.UserExperimentDTO, error)
	UserSegments(context.Context, int64) (*storage.UserExperimentListDTO, error)
	UserExperimentLogs(context.Context, int64, time.Time) ([]*storage.UserExperimentLogRecordDTO, error)
}

type Service struct {
	storage  Storage
	logsPath string
}

func New(storage Storage, logsPath string) *Service {
	logsPath = strings.TrimRight(logsPath, "/")

	return &Service{
		storage:  storage,
		logsPath: logsPath,
	}
}

func (svc *Service) CreateSegment(ctx context.Context, name string) (*model.Segment, error) {
	segmentDTO, err := svc.storage.AddSegment(ctx, name)
	if err != nil {
		return nil, err
	}

	return &model.Segment{
		ID:   segmentDTO.ID,
		Name: segmentDTO.Name,
	}, nil
}

func (svc *Service) DeleteSegment(ctx context.Context, name string) (*model.Segment, error) {
	segmentDTO, err := svc.storage.DeleteSegment(ctx, name)
	if err != nil {
		return nil, err
	}

	return &model.Segment{
		ID:   segmentDTO.ID,
		Name: segmentDTO.Name,
	}, nil
}

func (svc *Service) AddUserExperiments(ctx context.Context, userID int64, segmentNames []string) ([]*model.UserExperiment, error) {
	experiments := make([]*model.UserExperiment, 0, len(segmentNames))

	for _, segmentName := range segmentNames {
		expDTO, err := svc.storage.AddUserToSegment(ctx, userID, segmentName)

		if err != nil &&
			!(errors.Is(err, storage.ErrSegmentNotFound) ||
				errors.Is(err, storage.ErrAlreadyInExperiment)) {
			return nil, err
		}

		if expDTO == nil {
			continue
		}

		exp := &model.UserExperiment{
			ID:     expDTO.ID,
			UserID: expDTO.UserID,
			Segment: model.Segment{
				ID:   expDTO.Segment.ID,
				Name: expDTO.Segment.Name,
			},
		}

		experiments = append(experiments, exp)
	}

	return experiments, nil
}

func (svc *Service) RemoveUserExperiments(ctx context.Context, userID int64, segmentNames []string) ([]*model.UserExperiment, error) {
	experiments := make([]*model.UserExperiment, 0, len(segmentNames))

	for _, segmentName := range segmentNames {
		expDTO, err := svc.storage.DeleteUserFromSegment(ctx, userID, segmentName)

		if err != nil &&
			!(errors.Is(err, storage.ErrSegmentNotFound) ||
				errors.Is(err, storage.ErrUserExperimentNotFound)) {
			return nil, err
		}

		if expDTO == nil {
			continue
		}

		exp := &model.UserExperiment{
			ID:     expDTO.ID,
			UserID: expDTO.UserID,
			Segment: model.Segment{
				ID:   expDTO.Segment.ID,
				Name: expDTO.Segment.Name,
			},
		}

		experiments = append(experiments, exp)
	}

	return experiments, nil
}

func (svc *Service) ListUserSegments(ctx context.Context, userID int64) (*model.UserExperimentList, error) {
	listDTO, err := svc.storage.UserSegments(ctx, userID)
	if err != nil {
		return nil, err
	}

	list := &model.UserExperimentList{
		UserID:   listDTO.UserID,
		Segments: make([]model.Segment, 0, len(listDTO.Segments)),
	}

	for _, seg := range listDTO.Segments {
		list.Segments = append(list.Segments, (model.Segment)(seg))
	}

	return list, nil
}

func (s *Service) CreateLog(ctx context.Context, userID int64, start string) (*model.LogInfo, error) {
	from, err := time.Parse(logDateFormat, start)
	if err != nil {
		return nil, err
	}

	records, err := s.storage.UserExperimentLogs(ctx, userID, from)
	if err != nil {
		return nil, err
	}

	if err := os.Mkdir(s.logsPath, 0777); err != nil && !errors.Is(err, os.ErrExist) { //nolint:gomnd
		return nil, err
	}

	logName := fmt.Sprintf(logFilenameTemplate, userID, start)
	path := path.Join(s.logsPath, logName)

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = ';'
	defer w.Flush()

	if err := w.Write(
		[]string{"user_id", "segment_name", "operation", "added_at"},
	); err != nil {
		return nil, err
	}

	for _, record := range records {
		row := []string{
			fmt.Sprint(record.UserID),
			record.SegmentName,
			record.Operation,
			record.AddedAt.Local().Format(time.DateTime),
		}
		if err := w.Write(row); err != nil {
			return nil, err
		}
	}

	return &model.LogInfo{
		UserID: userID,
		From:   start,
		Path:   path,
	}, nil
}
