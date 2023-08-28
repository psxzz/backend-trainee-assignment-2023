package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/psxzz/backend-trainee-assignment/internal/app/model"
	"github.com/psxzz/backend-trainee-assignment/internal/app/storage"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) AddSegment(ctx context.Context, name string) (*model.Segment, error) {
	op := "storage.postgresql.AddSegment"

	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	row := conn.QueryRowContext(ctx,
		"INSERT INTO Segments(segment_name) VALUES ($1) RETURNING *;", name)

	if err := row.Err(); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrSegmentExists)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var segment model.Segment
	if err := row.Scan(&segment.ID, &segment.Name); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &segment, nil
}

func (s *Storage) DeleteSegment(ctx context.Context, name string) (*model.Segment, error) {
	op := "storage.postgresql.DeleteSegment"

	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	id, err := s.getSegmentID(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%s.getSegmentID: %w", op, err)
	}

	row := conn.QueryRowContext(ctx,
		"DELETE FROM Segments WHERE id = $1 RETURNING *;", id)

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var deleted model.Segment
	if err := row.Scan(&deleted.ID, &deleted.Name); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &deleted, nil
}

func (s *Storage) AddUserToSegment(ctx context.Context, userID int64, segmentName string) (*model.UserExperiment, error) {
	op := "storage.postgresql.AddUserToSegment"

	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	segmentID, err := s.getSegmentID(ctx, segmentName)
	if err != nil {
		return nil, fmt.Errorf("%s.getSegmentID: %w", op, err)
	}

	row := conn.QueryRowContext(ctx,
		"INSERT INTO user_experiments(user_id, segment_id) VALUES ($1, $2) RETURNING id;",
		userID, segmentID)

	if err := row.Err(); err != nil {
		if err := err.(*pq.Error); err != nil && err.Code == "23505" {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrAlreadyInExperiment)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &model.UserExperiment{
		ID:        id,
		UserID:    userID,
		SegmentID: segmentID,
	}, nil
}

func (s *Storage) DeleteUserFromSegment(ctx context.Context, userID int64, segmentName string) (*model.UserExperiment, error) {
	op := "storage.postgresql.DeleteUserFromSegment"

	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	segmentID, err := s.getSegmentID(ctx, segmentName)
	if err != nil {
		return nil, fmt.Errorf("%s.getSegmentID: %w", op, err)
	}

	row := conn.QueryRowContext(ctx,
		"DELETE FROM user_experiments WHERE user_id = $1 AND segment_id = $2 RETURNING *;",
		userID, segmentID)

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	deleted := &model.UserExperiment{}
	if err := row.Scan(deleted.ID, deleted.UserID, deleted.SegmentID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserExperimentNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return deleted, nil
}

func (s *Storage) UserSegments(ctx context.Context, userID int64) (*model.UserExperimentList, error) {
	op := "storage.postgresql.UserSegments"
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx,
		"SELECT s.id, s.segment_name FROM user_experiments u JOIN segments s "+
			"ON u.segment_id = s.id WHERE u.user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	expList := model.UserExperimentList{
		UserID: userID,
	}

	for rows.Next() {
		var seg model.Segment
		if err := rows.Scan(&seg.ID, &seg.Name); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		expList.Segments = append(expList.Segments, seg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &expList, nil
}

func (s *Storage) getSegmentID(ctx context.Context, name string) (int64, error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var id int64
	row := conn.QueryRowContext(ctx,
		"SELECT id FROM Segments WHERE segment_name = $1;", name)

	if err := row.Err(); err != nil {
		return 0, err
	}

	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, storage.ErrSegmentNotFound
		}
		return 0, err
	}

	return id, nil
}
