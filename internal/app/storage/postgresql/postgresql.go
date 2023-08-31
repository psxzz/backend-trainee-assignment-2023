package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/psxzz/backend-trainee-assignment/internal/app/storage"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) AddSegment(ctx context.Context, name string) (*storage.SegmentDTO, error) {
	op := "storage.postgresql.AddSegment"

	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	row := conn.QueryRowContext(ctx,
		"INSERT INTO Segments(segment_name) VALUES ($1) RETURNING *;", name)

	if err := row.Err(); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" { //nolint:errorlint
			return nil, fmt.Errorf("%s: %w", op, storage.ErrSegmentExists)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var segment storage.SegmentDTO
	if err := row.Scan(&segment.ID, &segment.Name); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &segment, nil
}

func (s *Storage) DeleteSegment(ctx context.Context, name string) (*storage.SegmentDTO, error) {
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

	var deleted storage.SegmentDTO
	if err := row.Scan(&deleted.ID, &deleted.Name); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &deleted, nil
}

func (s *Storage) AddUserToSegment(ctx context.Context, userID int64, segmentName string) (*storage.UserExperimentDTO, error) {
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
		if err := err.(*pq.Error); err != nil && err.Code == "23505" { //nolint:errorlint
			return nil, fmt.Errorf("%s: %w", op, storage.ErrAlreadyInExperiment)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.logExperiment(ctx, userID, segmentName, "add"); err != nil {
		log.Println(err)
	}

	return &storage.UserExperimentDTO{
		ID:     id,
		UserID: userID,
		Segment: storage.SegmentDTO{
			ID:   segmentID,
			Name: segmentName,
		},
	}, nil
}

func (s *Storage) DeleteUserFromSegment(ctx context.Context, userID int64, segmentName string) (*storage.UserExperimentDTO, error) {
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

	deleted := storage.UserExperimentDTO{
		Segment: storage.SegmentDTO{
			Name: segmentName,
		},
	}

	if err := row.Scan(&deleted.ID, &deleted.UserID, &deleted.Segment.ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserExperimentNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.logExperiment(ctx, userID, segmentName, "remove"); err != nil {
		log.Println(err)
	}

	return &deleted, nil
}

func (s *Storage) UserSegments(ctx context.Context, userID int64) (*storage.UserExperimentListDTO, error) {
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

	expList := &storage.UserExperimentListDTO{
		UserID: userID,
	}

	for rows.Next() {
		var seg storage.SegmentDTO
		if err := rows.Scan(&seg.ID, &seg.Name); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		expList.Segments = append(expList.Segments, seg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return expList, nil
}

func (s *Storage) UserExperimentLogs(ctx context.Context, userID int64, start time.Time) ([]*storage.UserExperimentLogRecordDTO, error) {
	op := "storage.postgresql.UserExperimentLogs"
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx,
		"SELECT user_id, segment_name, op_type, added_at FROM "+
			"log_user_experiments WHERE user_id = $1 AND "+
			"added_at BETWEEN $2 AND $2 + INTERVAL '1 month'", userID, start)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var records []*storage.UserExperimentLogRecordDTO

	for rows.Next() {
		var rec storage.UserExperimentLogRecordDTO

		if err := rows.Scan(&rec.UserID, &rec.SegmentName,
			&rec.Operation, &rec.AddedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		records = append(records, &rec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return records, nil
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

func (s *Storage) logExperiment(ctx context.Context, userID int64, segmentName, opType string) error {
	op := "storage.postgresql.logExperiment"
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		"INSERT INTO log_user_experiments(user_id, segment_name, op_type)"+
			"VALUES ($1, $2, $3);", userID, segmentName, opType)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
