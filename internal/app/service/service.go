package service

import (
	"context"
	"errors"

	"github.com/psxzz/backend-trainee-assignment/internal/app/model"
	"github.com/psxzz/backend-trainee-assignment/internal/app/storage"
)

type Storage interface {
	AddSegment(context.Context, string) (*storage.SegmentDTO, error)
	DeleteSegment(context.Context, string) (*storage.SegmentDTO, error)
	AddUserToSegment(context.Context, int64, string) (*storage.UserExperimentDTO, error)
	DeleteUserFromSegment(context.Context, int64, string) (*storage.UserExperimentDTO, error)
	UserSegments(context.Context, int64) (*storage.UserExperimentListDTO, error)
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (svc *Service) CreateSegment(ctx context.Context, name string) (*model.Segment, error) {
	segment, err := svc.storage.AddSegment(ctx, name)
	if err != nil {
		return nil, err
	}

	return (*model.Segment)(segment), nil
}

func (svc *Service) DeleteSegment(ctx context.Context, name string) (*model.Segment, error) {
	segment, err := svc.storage.DeleteSegment(ctx, name)
	if err != nil {
		return nil, err
	}

	return (*model.Segment)(segment), nil
}

func (svc *Service) AddUserExperiments(ctx context.Context, userID int64, segmentNames []string) ([]*model.UserExperiment, error) {
	var experiments []*model.UserExperiment

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
	var experiments []*model.UserExperiment
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
