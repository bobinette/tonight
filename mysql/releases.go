package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bobinette/tonight"
	"github.com/google/uuid"
)

type ReleaseStore struct {
	db *sql.DB
}

func NewReleaseStore(db *sql.DB) ReleaseStore {
	return ReleaseStore{db: db}
}

func (s ReleaseStore) Get(ctx context.Context, id uuid.UUID) (tonight.Release, error) {
	query := `
SELECT uuid, title, description, project_uuid, created_at, updated_at
FROM releases
WHERE uuid = ?
`

	row := s.db.QueryRowContext(ctx, query, id)
	var release tonight.Release
	err := row.Scan(
		&release.UUID,
		&release.Title,
		&release.Description,
		&release.Project.UUID,
		&release.CreatedAt,
		&release.UpdatedAt,
	)
	if err != nil {
		return tonight.Release{}, err
	}

	tasks, err := s.loadTasks(ctx, []string{release.UUID.String()})
	if err != nil {
		return tonight.Release{}, err
	}

	release.Tasks = tasks[release.UUID.String()]
	if release.Tasks == nil {
		release.Tasks = make([]tonight.Task, 0)
	}

	return release, nil
}

func (s ReleaseStore) List(ctx context.Context, projectUUID uuid.UUID) ([]tonight.Release, error) {
	query := `
SELECT uuid, title, description, project_uuid, created_at, updated_at
FROM releases
WHERE project_uuid = ?
ORDER BY
	CASE WHEN project_uuid = uuid
	THEN 1
	OTHER 0
	END
	ASC,
	title ASC
`
	rows, err := s.db.QueryContext(ctx, query, projectUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	releases := make([]tonight.Release, 0)
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
		releases = append(releases, release)
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
		for i, release := range releases {
			if releaseUUID == release.UUID.String() {
				release.Tasks = tasks
				releases[i] = release
				break
			}
		}
	}

	return releases, nil
}

func (s ReleaseStore) Upsert(ctx context.Context, release tonight.Release) error {
	query := `
	INSERT INTO releases (uuid, title, description, project_uuid, created_at, updated_at)
	VALUE (?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
		title = ?
	`
	_, err := s.db.ExecContext(
		ctx,
		query,
		release.UUID,
		release.Title,
		release.Description,
		release.Project.UUID,
		release.CreatedAt,
		release.UpdatedAt,
		release.Title,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s ReleaseStore) loadTasks(ctx context.Context, uuids []string) (map[string][]tonight.Task, error) {
	if len(uuids) == 0 {
		return nil, nil
	}

	qArgs, args := prepareArgs(uuids)
	query := fmt.Sprintf(`
SELECT uuid, title, status, release_uuid, created_at, updated_at
FROM tasks
WHERE release_uuid IN %s AND tasks.deleted = 0
ORDER BY -rank DESC, created_at
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
