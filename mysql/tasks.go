package mysql

import (
	"context"
	"database/sql"
	"fmt"

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
INSERT INTO tasks (uuid, title, project_uuid, created_at, updated_at)
VALUE (?, ?, ?, ?, ?)
`
	_, err := s.db.ExecContext(ctx, query, t.UUID, t.Title, t.Project.UUID, t.CreatedAt, t.UpdatedAt)
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
SELECT uuid, title, project_uuid, created_at, updated_at
FROM tasks
ORDER BY created_at
`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]tonight.Task, 0)
	projectUUIDsSet := make(map[string]struct{})
	projectUUIDs := make([]string, 0)
	for rows.Next() {
		var t tonight.Task
		err := rows.Scan(
			&t.UUID,
			&t.Title,
			&t.Project.UUID,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		uuidStr := t.Project.UUID.String()
		if _, ok := projectUUIDsSet[uuidStr]; !ok {
			projectUUIDsSet[uuidStr] = struct{}{}
			projectUUIDs = append(projectUUIDs, uuidStr)
		}
		tasks = append(tasks, t)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	projects, err := s.loadProjects(ctx, projectUUIDs)
	if err != nil {
		return nil, err
	}

	findProject := func(uuid string, ps []tonight.Project) *tonight.Project {
		for _, p := range ps {
			if p.UUID.String() == uuid {
				p := p
				return &p
			}
		}
		return nil
	}

	for i, task := range tasks {
		p := findProject(task.Project.UUID.String(), projects)
		if p != nil {
			task.Project = *p
		}
		tasks[i] = task
	}

	return tasks, nil
}

func (s TaskStore) loadProjects(ctx context.Context, uuids []string) ([]tonight.Project, error) {
	qArgs, args := prepareArgs(uuids)
	query := fmt.Sprintf(`
SELECT uuid, name, created_at, updated_at
FROM projects
WHERE uuid IN %s
`, qArgs...)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]tonight.Project, 0)
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
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	return projects, nil
}
