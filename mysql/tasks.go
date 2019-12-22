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
INSERT INTO tasks (uuid, title, created_at, updated_at)
VALUE (?, ?, ?, ?)
`
	_, err := s.db.ExecContext(ctx, query, t.UUID, t.Title, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s TaskStore) Get(ctx context.Context, uuid uuid.UUID) (tonight.Task, error) {
	return tonight.Task{}, nil
}

func (s TaskStore) List(ctx context.Context) ([]tonight.Task, error) {
	query := `
SELECT uuid, title, created_at, updated_at
FROM tasks
ORDER BY created_at
`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]tonight.Task, 0)
	for rows.Next() {
		var t tonight.Task
		err := rows.Scan(
			&t.UUID,
			&t.Title,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	return tasks, nil
}
