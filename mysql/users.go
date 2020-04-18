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

func (s UserStore) Permission(ctx context.Context, user auth.User, projectUUID uuid.UUID) (string, error) {
	query := `
SELECT permission
FROM user_permission_on_project
WHERE user_id = ? AND project_uuid = ?
`
	row := s.db.QueryRowContext(ctx, query, user.ID, projectUUID)
	var perm string
	if err := row.Scan(&perm); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return perm, nil
}
