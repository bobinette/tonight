package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/bobinette/tonight/tonight"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Get(ctx context.Context, id uint) (tonight.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, username FROM users WHERE id = ?", id)
	return r.get(ctx, row)
}

func (r *UserRepository) GetByName(ctx context.Context, username string) (tonight.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, username FROM users WHERE username = ?", username)
	return r.get(ctx, row)
}

func (r *UserRepository) get(ctx context.Context, row *sql.Row) (tonight.User, error) {
	var id uint
	var username string
	if err := row.Scan(&id, &username); err != nil {
		if err == sql.ErrNoRows {
			return tonight.User{}, nil
		}
		return tonight.User{}, err
	}

	rows, err := r.db.QueryContext(ctx, "SELECT task_id FROM user_has_tasks WHERE user_id = ?", id)
	if err != nil {
		return tonight.User{}, err
	}
	defer rows.Close()

	taskIDs := make([]uint, 0)
	for rows.Next() {
		var taskID uint
		if err := rows.Scan(&taskID); err != nil {
			return tonight.User{}, err
		}

		taskIDs = append(taskIDs, taskID)
	}

	if err := rows.Close(); err != nil {
		return tonight.User{}, err
	}

	rows, err = r.db.QueryContext(ctx, "SELECT tag, colour FROM user_customs_tags WHERE user_id = ?", id)
	if err != nil {
		return tonight.User{}, err
	}
	defer rows.Close()

	tagColours := make(map[string]string)
	for rows.Next() {
		var tag string
		var colour string
		if err := rows.Scan(&tag, &colour); err != nil {
			return tonight.User{}, err
		}

		tagColours[tag] = colour
	}

	if err := rows.Close(); err != nil {
		return tonight.User{}, err
	}

	return tonight.User{
		ID:         id,
		Name:       username,
		TaskIDs:    taskIDs,
		TagColours: tagColours,
	}, nil
}

func (r *UserRepository) Insert(ctx context.Context, user *tonight.User) error {
	now := time.Now()
	res, err := r.db.ExecContext(ctx, "INSERT INTO users (username, created_at) VALUES (?, ?)", user.Name, now)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint(id)
	return nil
}

func (r *UserRepository) AddTaskToUser(ctx context.Context, userID uint, taskID uint) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_has_tasks (user_id, task_id, created_at)
		VALUES (?, ?, ?)
	`, userID, taskID, time.Now())
	return err
}

func (r *UserRepository) UpdateTagColor(ctx context.Context, userID uint, tag string, colour string) error {
	now := time.Now()
	row := r.db.QueryRowContext(ctx, "SELECT NULL FROM user_customs_tags WHERE user_id = ? AND tag = ?", userID, tag)
	var useless interface{}
	if err := row.Scan(&useless); err != nil {
		if err == sql.ErrNoRows {
			_, err := r.db.ExecContext(ctx, `
				INSERT INTO user_customs_tags (user_id, tag, colour, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?)
			`, userID, tag, colour, now, now)
			return err
		}
		return err
	}

	_, err := r.db.ExecContext(ctx, `
		UPDATE user_customs_tags SET colour = ?, updated_at = ? WHERE user_id = ? AND tag = ?
	`, colour, now, userID, tag)
	return err
}
