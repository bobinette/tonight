package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bobinette/tonight"
)

type ProjectStore struct {
	db *sql.DB
}

func NewProjectStore(db *sql.DB) ProjectStore {
	return ProjectStore{db: db}
}

func (s ProjectStore) Upsert(ctx context.Context, p tonight.Project) error {
	query := `
INSERT INTO projects (uuid, name, created_at, updated_at)
VALUE (?, ?, ?, ?)
`
	_, err := s.db.ExecContext(ctx, query, p.UUID, p.Name, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s ProjectStore) List(ctx context.Context) ([]tonight.Project, error) {
	query := `
SELECT uuid, name, created_at, updated_at
FROM projects
ORDER BY created_at
`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]tonight.Project, 0)
	uuids := make([]string, 0)
	for rows.Next() {
		var p tonight.Project
		err := rows.Scan(
			&p.UUID,
			&p.Name,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		projects = append(projects, p)
		uuids = append(uuids, p.UUID.String())
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return projects, nil
	}

	tasks, err := s.loadTasks(ctx, uuids)
	if err != nil {
		return nil, err
	}

	for i, p := range projects {
		p.Tasks = tasks[p.UUID.String()]
		if p.Tasks == nil {
			p.Tasks = make([]tonight.Task, 0)
		}
		projects[i] = p
	}

	return projects, nil
}

func (s ProjectStore) loadTasks(ctx context.Context, uuids []string) (map[string][]tonight.Task, error) {
	qArgs, args := prepareArgs(uuids)
	query := fmt.Sprintf(`
SELECT uuid, title, project_uuid, created_at, updated_at
FROM tasks
WHERE project_uuid IN %s
ORDER BY created_at
`, qArgs...)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasksByProjectUUID := make(map[string][]tonight.Task)
	for rows.Next() {
		var t tonight.Task
		var projectUUID string
		err := rows.Scan(
			&t.UUID,
			&t.Title,
			&projectUUID,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		tasksByProjectUUID[projectUUID] = append(tasksByProjectUUID[projectUUID], t)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return tasksByProjectUUID, nil
}
