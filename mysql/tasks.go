package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/bobinette/tonight"
	"github.com/bobinette/tonight/auth"
)

type TaskStore struct {
	db *sql.DB
}

func NewTaskStore(db *sql.DB) TaskStore {
	return TaskStore{db: db}
}

func (s TaskStore) Get(ctx context.Context, uuid uuid.UUID, user auth.User) (tonight.Task, error) {
	query := `
SELECT tasks.uuid, tasks.title, tasks.status, tasks.release_uuid, tasks.created_at, tasks.updated_at, releases.project_uuid
FROM tasks
JOIN releases ON releases.uuid = tasks.release_uuid
JOIN user_permission_on_project ON user_permission_on_project.project_uuid = releases.project_uuid
WHERE user_permission_on_project.user_id = ? AND tasks.uuid = ? AND tasks.deleted = 0
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
		&t.Release.Project.UUID,
	)
	if err != nil {
		return tonight.Task{}, err
	}
	return t, nil
}

func (s TaskStore) GetProjectUUIDs(ctx context.Context, uuids []uuid.UUID) ([]uuid.UUID, error) {
	if len(uuids) == 0 {
		return nil, nil
	}

	qArgs, args := prepareArgs(uuids)
	query := fmt.Sprintf(`
SELECT releases.project_uuid
FROM releases
JOIN tasks ON tasks.release_uuid = releases.uuid
WHERE tasks.uuid IN %s
`, qArgs...)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projectUUIDs := make([]uuid.UUID, 0, len(uuids))
	for rows.Next() {
		var projectUUID uuid.UUID
		if err := rows.Scan(&projectUUID); err != nil {
			return nil, err
		}
		projectUUIDs = append(projectUUIDs, projectUUID)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return projectUUIDs, nil
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

func (s TaskStore) Delete(ctx context.Context, taskUUID uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, "UPDATE tasks SET deleted = 1 where uuid = ?", taskUUID)
	return err
}
