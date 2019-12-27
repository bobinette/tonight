package mysql

import (
	"context"
	"database/sql"

	"github.com/bobinette/tonight"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) UserStore {
	return UserStore{db: db}
}

func (s UserStore) Ensure(ctx context.Context, user *tonight.User) error {
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

func (s UserStore) insertUser(ctx context.Context, user *tonight.User) error {
	query := `
INSERT INTO users (id, name)
VALUES (?, ?)
`
	_, err := s.db.ExecContext(ctx, query, user.ID, user.Name)
	return err
}
