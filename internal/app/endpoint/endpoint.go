package endpoint

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/psxzz/backend-trainee-assignment/internal/app/model"
	"github.com/psxzz/backend-trainee-assignment/internal/app/storage"
)

type Service interface {
	CreateSegment(context.Context, string) (*model.Segment, error)
	DeleteSegment(context.Context, string) (*model.Segment, error)
	AddUserExperiments(context.Context, int64, []string) ([]*model.UserExperiment, error)
	RemoveUserExperiments(context.Context, int64, []string) ([]*model.UserExperiment, error)
	ListUserSegments(context.Context, int64) (*model.UserExperimentList, error)
	CreateLog(context.Context, int64, string) (*model.LogInfo, error)
}

type Endpoint struct {
	svc Service
}

func New(svc Service) *Endpoint {
	return &Endpoint{
		svc: svc,
	}
}

func (e *Endpoint) HandleCreate(ctx echo.Context) error {
	var req segmentRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusMethodNotAllowed, errorResponse{
			Message: "Validation error: field 'name' not found",
		})
	}

	segment, err := e.svc.CreateSegment(ctx.Request().Context(), req.Name)
	if err != nil {
		if errors.Is(err, storage.ErrSegmentExists) {
			return ctx.JSON(http.StatusBadRequest, errorResponse{
				Message: errors.Unwrap(err).Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			Message: "Internal error",
		})
	}

	return ctx.JSON(http.StatusOK, segment)
}

func (e *Endpoint) HandleDelete(ctx echo.Context) error {
	var req segmentRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			Message: "Internal error",
		})
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusMethodNotAllowed, errorResponse{
			Message: "Validation error: field 'name' not found",
		})
	}

	segment, err := e.svc.DeleteSegment(ctx.Request().Context(), req.Name)
	if err != nil {
		if errors.Is(err, storage.ErrSegmentNotFound) {
			return ctx.JSON(http.StatusNotFound, errorResponse{
				Message: errors.Unwrap(err).Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			Message: "Internal error",
		})
	}

	return ctx.JSON(http.StatusOK, segment)
}

func (e *Endpoint) HandleExperiments(ctx echo.Context) error {
	var req userExperimentRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			Message: "Internal error",
		})
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusMethodNotAllowed, errorResponse{
			Message: "Validation error: invalid request body",
		})
	}

	added, err := e.svc.AddUserExperiments(ctx.Request().Context(), req.UserID, req.ToAdd)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			Message: "Internal error",
		})
	}

	removed, err := e.svc.RemoveUserExperiments(ctx.Request().Context(), req.UserID, req.ToRemove)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			Message: "Internal error",
		})
	}

	return ctx.JSON(http.StatusOK, userExperimentResponse{
		UserID:  req.UserID,
		Added:   added,
		Removed: removed,
	})
}

func (e *Endpoint) HandleUserExperimentList(ctx echo.Context) error {
	var req experimentListRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			Message: "Internal error",
		})
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusMethodNotAllowed, errorResponse{
			Message: "Validation error: invalid request body",
		})
	}

	list, err := e.svc.ListUserSegments(ctx.Request().Context(), req.UserID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			Message: "Internal error",
		})
	}

	return ctx.JSON(http.StatusOK, list)
}

func (e *Endpoint) HandleCreateLog(ctx echo.Context) error {
	var req userLogRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			Message: "Internal error",
		})
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusMethodNotAllowed, errorResponse{
			Message: "Validation error: invalid request body",
		})
	}

	info, err := e.svc.CreateLog(ctx.Request().Context(), req.UserID, req.From)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			Message: "Internal error",
		})
	}

	return ctx.JSON(http.StatusOK, info)
}

type errorResponse struct {
	Message string `json:"message"`
}

type segmentRequest struct {
	Name string `json:"name" validate:"required"`
}

type userExperimentRequest struct {
	UserID   int64    `json:"user_id" validate:"required"`
	ToAdd    []string `json:"to_add" validate:"required"`
	ToRemove []string `json:"to_remove" validate:"required"`
}

type userExperimentResponse struct {
	UserID  int64                   `json:"user_id"`
	Added   []*model.UserExperiment `json:"added"`
	Removed []*model.UserExperiment `json:"removed"`
}

type experimentListRequest struct {
	UserID int64 `json:"user_id" validate:"required"`
}

type userLogRequest struct {
	UserID int64  `json:"user_id" validate:"required"`
	From   string `json:"from" validate:"required"`
}
