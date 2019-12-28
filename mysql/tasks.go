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
INSERT INTO tasks (uuid, title, status, project_uuid, created_at, updated_at)
VALUE (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	status = ?
`
	_, err := s.db.ExecContext(
		ctx,
		query,
		t.UUID,
		t.Title,
		t.Status,
		t.Project.UUID,
		t.CreatedAt,
		t.UpdatedAt,
		t.Status,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s TaskStore) Get(ctx context.Context, uuid uuid.UUID, user tonight.User) (tonight.Task, error) {
	query := `
SELECT tasks.uuid, tasks.title, tasks.status, tasks.project_uuid, tasks.created_at, tasks.updated_at
FROM tasks
JOIN user_permission_on_project ON user_permission_on_project.project_uuid = tasks.project_uuid
WHERE user_permission_on_project.user_id = ? AND tasks.uuid = ?
`
	row := s.db.QueryRowContext(ctx, query, user.ID, uuid)
	var t tonight.Task
	err := row.Scan(
		&t.UUID,
		&t.Title,
		&t.Status,
		&t.Project.UUID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return tonight.Task{}, err
	}
	return t, nil
}
