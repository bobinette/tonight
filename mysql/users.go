package mysql

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/bobinette/tonight/auth"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) UserStore {
	return UserStore{db: db}
}

func (s UserStore) Ensure(ctx context.Context, user *auth.User) error {
	query := `
SELECT id, name FROM users
WHERE id = ?
`
	row := s.db.QueryRowContext(ctx, query, user.ID)
	if err := row.Scan(&user.ID, &user.Name); err != nil {
		if err != sql.ErrNoRows {
			return err
		}

		return s.insertUser(ctx, user)
	}

	return nil
}

func (s UserStore) insertUser(ctx context.Context, user *auth.User) error {
	query := `
INSERT INTO users (id, name)
VALUES (?, ?)
`
	_, err := s.db.ExecContext(ctx, query, user.ID, user.Name)
	return err
}

func (s UserStore) Permission(ctx context.Context, user auth.User, projectUUID uuid.UUID) (auth.Permission, error) {
	query := `
SELECT permission
FROM user_permission_on_project
WHERE user_id = ? AND project_uuid = ?
`
	row := s.db.QueryRowContext(ctx, query, user.ID, projectUUID)
	var perm auth.Permission
	if err := row.Scan(&perm); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return perm, nil
}

func (s UserStore) SetPermission(ctx context.Context, user auth.User, projectUUID uuid.UUID, perm auth.Permission) error {
	query := `
INSERT IGNORE INTO user_permission_on_project (user_id, project_uuid, permission)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
	perm = ?
`
	if _, err := s.db.ExecContext(ctx, query, user.ID, projectUUID, perm); err != nil {
		return err
	}

	return nil
}

func (s UserStore) AllPermissions(ctx context.Context, user auth.User) (map[uuid.UUID]auth.Permission, error) {
	query := `
SELECT project_uuid, permission
FROM user_permission_on_project
WHERE user_id = ?
`
	rows, err := s.db.QueryContext(ctx, query, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	perms := make(map[uuid.UUID]auth.Permission)
	for rows.Next() {
		var u uuid.UUID
		var perm auth.Permission
		if err := rows.Scan(&u, &perm); err != nil {
			return nil, err
		}
		perms[u] = perm
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	return perms, nil
}
