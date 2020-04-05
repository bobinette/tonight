package mysql

import (
	"context"
	"database/sql"

	uuid "github.com/satori/go.uuid"

	"github.com/bobinette/tonight"
)

type TaskStore struct {
	db *sql.DB
}

func NewTaskStore(db *sql.DB) TaskStore {
	return TaskStore{db: db}
}

func (s TaskStore) Upsert(ctx context.Context, t tonight.Task) error {
	query := `
INSERT INTO tasks (uuid, title, status, release_uuid, created_at, updated_at)
VALUE (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	status = ?,
	title = ?
`
	_, err := s.db.ExecContext(
		ctx,
		query,
		t.UUID,
		t.Title,
		t.Status,
		t.Release.UUID,
		t.CreatedAt,
		t.UpdatedAt,
		t.Status,
		t.Title,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s TaskStore) Get(ctx context.Context, uuid uuid.UUID, user tonight.User) (tonight.Task, error) {
	query := `
SELECT tasks.uuid, tasks.title, tasks.status, tasks.release_uuid, tasks.created_at, tasks.updated_at
FROM tasks
JOIN releases ON releases.uuid = tasks.release_uuid
JOIN user_permission_on_project ON user_permission_on_project.project_uuid = releases.project_uuid
WHERE user_permission_on_project.user_id = ? AND tasks.uuid = ?
`
	row := s.db.QueryRowContext(ctx, query, user.ID, uuid)
	var t tonight.Task
	err := row.Scan(
		&t.UUID,
		&t.Title,
		&t.Status,
		&t.Release.UUID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return tonight.Task{}, err
	}
	return t, nil
}

func (s TaskStore) Reorder(ctx context.Context, rankedUUIDs []uuid.UUID) (err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		e := tx.Rollback()
		if err == nil && e != sql.ErrTxDone {
			err = e
		}
	}()

	query := "UPDATE tasks SET rank = ? WHERE uuid = ?"
	for rank, taskUUID := range rankedUUIDs {
		if _, err := tx.ExecContext(ctx, query, rank, taskUUID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
