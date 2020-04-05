package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bobinette/tonight"
	uuid "github.com/satori/go.uuid"
)

type ProjectStore struct {
	db *sql.DB
}

func NewProjectStore(db *sql.DB) ProjectStore {
	return ProjectStore{db: db}
}

func (s ProjectStore) Upsert(ctx context.Context, p tonight.Project, u tonight.User) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
INSERT INTO projects (uuid, name, description, slug, created_at, updated_at)
VALUE (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	name = ?,
	description = ?,
	slug = ?,
	updated_at = ?
`
	if _, err := tx.ExecContext(
		ctx,
		query,
		p.UUID,
		p.Name,
		p.Description,
		p.Slug,
		p.CreatedAt,
		p.UpdatedAt,
		// update
		p.Name,
		p.Description,
		p.Slug,
		p.UpdatedAt,
	); err != nil {
		return err
	}

	query = `
INSERT IGNORE INTO user_permission_on_project (user_id, project_uuid, permission)
VALUES (?, ?, ?)
`
	if _, err := tx.ExecContext(ctx, query, u.ID, p.UUID, "owner"); err != nil {
		return err
	}

	query = `
INSERT IGNORE INTO releases (uuid, title, description, project_uuid, created_at, updated_at)
VALUE (?, ?, ?, ?, ?, ?)
`
	if _, err := tx.ExecContext(ctx, query, p.UUID, "Backlog", "", p.UUID, p.CreatedAt, p.UpdatedAt); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s ProjectStore) List(ctx context.Context, u tonight.User) ([]tonight.Project, error) {
	query := `
SELECT projects.uuid, projects.name, projects.description, projects.slug, projects.created_at, projects.updated_at
FROM projects
JOIN user_permission_on_project ON user_permission_on_project.project_uuid = projects.uuid
WHERE user_permission_on_project.user_id = ?
ORDER BY created_at
`
	rows, err := s.db.QueryContext(ctx, query, u.ID)
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
			&p.Description,
			&p.Slug,
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

	releases, err := s.loadReleases(ctx, uuids)
	if err != nil {
		return nil, err
	}

	for i, p := range projects {
		p.Releases = releases[p.UUID.String()]
		if p.Releases == nil {
			p.Releases = make([]tonight.Release, 0)
		}
		projects[i] = p
	}

	return projects, nil
}

func (s ProjectStore) Get(ctx context.Context, uuid uuid.UUID, u tonight.User) (tonight.Project, error) {
	query := `
SELECT projects.uuid, projects.name, projects.description, projects.slug, projects.created_at, projects.updated_at
FROM projects
JOIN user_permission_on_project ON user_permission_on_project.project_uuid = projects.uuid
WHERE projects.uuid = ? AND user_permission_on_project.user_id = ?
ORDER BY created_at
`

	row := s.db.QueryRowContext(ctx, query, uuid, u.ID)
	var p tonight.Project
	err := row.Scan(
		&p.UUID,
		&p.Name,
		&p.Description,
		&p.Slug,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return tonight.Project{}, err
	}

	releases, err := s.loadReleases(ctx, []string{p.UUID.String()})
	if err != nil {
		return tonight.Project{}, err
	}

	p.Releases = releases[p.UUID.String()]
	if p.Releases == nil {
		p.Releases = make([]tonight.Release, 0)
	}

	return p, nil
}

func (s ProjectStore) Find(ctx context.Context, slug string, u tonight.User) (tonight.Project, error) {
	query := `
SELECT projects.uuid, projects.name, projects.description, projects.slug, projects.created_at, projects.updated_at
FROM projects
JOIN user_permission_on_project ON user_permission_on_project.project_uuid = projects.uuid
WHERE projects.slug = ? AND user_permission_on_project.user_id = ?
ORDER BY created_at
`

	row := s.db.QueryRowContext(ctx, query, slug, u.ID)
	var p tonight.Project
	err := row.Scan(
		&p.UUID,
		&p.Name,
		&p.Description,
		&p.Slug,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return tonight.Project{}, err
	}

	releases, err := s.loadReleases(ctx, []string{p.UUID.String()})
	if err != nil {
		return tonight.Project{}, err
	}

	p.Releases = releases[p.UUID.String()]
	if p.Releases == nil {
		p.Releases = make([]tonight.Release, 0)
	}

	return p, nil
}

func (s ProjectStore) loadReleases(ctx context.Context, projectUUIDs []string) (map[string][]tonight.Release, error) {
	if len(projectUUIDs) == 0 {
		return nil, nil
	}

	qArgs, args := prepareArgs(projectUUIDs)
	query := fmt.Sprintf(`
SELECT uuid, title, description, project_uuid, created_at, updated_at
FROM releases
WHERE project_uuid IN %s
ORDER BY
	CASE WHEN project_uuid = uuid
	THEN 1
	ELSE 0
	END ASC,
	title ASC
`, qArgs...)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	releasesByProjectUUID := make(map[string][]tonight.Release, 0)
	releaseUUIDs := make([]string, 0)
	for rows.Next() {
		var release tonight.Release
		err := rows.Scan(
			&release.UUID,
			&release.Title,
			&release.Description,
			&release.Project.UUID,
			&release.CreatedAt,
			&release.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		release.Tasks = make([]tonight.Task, 0)
		releasesByProjectUUID[release.Project.UUID.String()] = append(
			releasesByProjectUUID[release.Project.UUID.String()],
			release,
		)
		releaseUUIDs = append(releaseUUIDs, release.UUID.String())
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	tasksByReleaseUUID, err := s.loadTasks(ctx, releaseUUIDs)
	if err != nil {
		return nil, err
	}

	for releaseUUID, tasks := range tasksByReleaseUUID {
		for _, releases := range releasesByProjectUUID {
			for i, release := range releases {
				if releaseUUID == release.UUID.String() {
					release.Tasks = tasks
					releases[i] = release
					break
				}
			}
		}
	}

	return releasesByProjectUUID, nil
}

func (s ProjectStore) loadTasks(ctx context.Context, uuids []string) (map[string][]tonight.Task, error) {
	qArgs, args := prepareArgs(uuids)
	query := fmt.Sprintf(`
SELECT tasks.uuid, tasks.title, tasks.status, tasks.release_uuid, tasks.created_at, tasks.updated_at
FROM tasks
JOIN releases ON releases.uuid = tasks.release_uuid
WHERE releases.project_uuid IN %s
ORDER BY -tasks.rank DESC, tasks.created_at
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
			&t.Status,
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
