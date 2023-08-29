package service_test

import (
	"context"
	"testing"

	"github.com/psxzz/backend-trainee-assignment/internal/app/model"
	"github.com/psxzz/backend-trainee-assignment/internal/app/service"
	"github.com/psxzz/backend-trainee-assignment/internal/app/storage"
	"github.com/psxzz/backend-trainee-assignment/internal/app/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestSegments(t *testing.T) {
	t.Run("creates new segment", func(t *testing.T) {
		var (
			db  = memory.New()
			svc = service.New(db)
		)

		resp, err := svc.CreateSegment(context.Background(), "Hello")
		assert.NoError(t, err)

		assert.Equal(t, resp, &model.Segment{ID: 0, Name: "Hello"})
	})

	t.Run("returns error if duplicate", func(t *testing.T) {
		var (
			db  = memory.New()
			svc = service.New(db)
		)

		svc.CreateSegment(context.Background(), "Hello")
		_, err := svc.CreateSegment(context.Background(), "Hello")
		assert.ErrorIs(t, err, storage.ErrSegmentExists)
	})

	t.Run("deletes existing segment", func(t *testing.T) {
		var (
			db  = memory.New()
			svc = service.New(db)
		)

		respCreate, err := svc.CreateSegment(context.Background(), "Hello")
		assert.NoError(t, err)
		respDelete, err := svc.DeleteSegment(context.Background(), "Hello")
		assert.NoError(t, err)

		assert.Equal(t, respCreate, respDelete)
	})

	t.Run("returns error if segment not exists", func(t *testing.T) {
		var (
			db  = memory.New()
			svc = service.New(db)
		)

		_, err := svc.DeleteSegment(context.Background(), "Hello")
		assert.ErrorIs(t, err, storage.ErrSegmentNotFound)
	})
}

func TestUserExperiments(t *testing.T) {
	t.Run("creates new user experiment", func(t *testing.T) {
		var (
			db  = memory.New()
			svc = service.New(db)
		)

		_, err := svc.CreateSegment(context.Background(), "Hello")
		assert.NoError(t, err)

		exp, err := svc.AddUserExperiments(context.Background(), 1010, []string{"Hello"})
		assert.NoError(t, err)

		assert.Contains(t, db.Experiments()[1010], struct {
			ID        int64
			UserID    int64
			SegmentID int64
		}{ID: exp[0].ID, UserID: 1010, SegmentID: exp[0].Segment.ID})
	})

	t.Run("skips insert if experiment exists", func(t *testing.T) {
		var (
			db  = memory.New()
			svc = service.New(db)
		)

		_, err := svc.CreateSegment(context.Background(), "Hello")
		assert.NoError(t, err)

		_, err = svc.AddUserExperiments(context.Background(), 1010, []string{"Hello"})
		assert.NoError(t, err)

		resp, err := svc.AddUserExperiments(context.Background(), 1010, []string{"Hello"})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(resp))
	})

	t.Run("deletes existing user experiment", func(t *testing.T) {
		var (
			db  = memory.New()
			svc = service.New(db)
		)

		svc.CreateSegment(context.Background(), "Hello")
		respCreate, err := svc.AddUserExperiments(context.Background(), 1010, []string{"Hello"})
		assert.NoError(t, err)
		respDelete, err := svc.RemoveUserExperiments(context.Background(), 1010, []string{"Hello"})
		assert.NoError(t, err)

		assert.ElementsMatch(t, respCreate, respDelete)
	})

	t.Run("skip delete if user experiment not exists", func(t *testing.T) {
		var (
			db  = memory.New()
			svc = service.New(db)
		)

		resp, err := svc.RemoveUserExperiments(context.Background(), 1010, []string{"Hello"})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(resp))
	})
}

func TestListExperiments(t *testing.T) {
	t.Run("lists all user experiments", func(t *testing.T) {
		var (
			db  = memory.New()
			svc = service.New(db)
		)

		resp, err := svc.ListUserSegments(context.Background(), 1010)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(resp.Segments))

		seg1, err := svc.CreateSegment(context.Background(), "Hello")
		assert.NoError(t, err)

		seg2, err := svc.CreateSegment(context.Background(), "World")
		assert.NoError(t, err)
		_, err = svc.AddUserExperiments(context.Background(), 1010, []string{"Hello", "World"})
		assert.NoError(t, err)

		resp, err = svc.ListUserSegments(context.Background(), 1010)
		assert.NoError(t, err)
		assert.ElementsMatch(t, resp.Segments, []model.Segment{*seg1, *seg2})
	})
}
