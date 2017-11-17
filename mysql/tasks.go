package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/bobinette/tonight"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(addr string) (*TaskRepository, error) {
	db, err := sql.Open("mysql", addr)
	if err != nil {
		return nil, err
	}

	return &TaskRepository{db: db}, nil
}

func (r *TaskRepository) Close() error {
	return r.db.Close()
}

func (r *TaskRepository) List(ctx context.Context, done bool) ([]tonight.Task, error) {
	orderBy := "rank"
	if done {
		orderBy = "doneAt DESC"
	}

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT id, title, description, done, doneAt, created_at
		  FROM tasks
		 WHERE done = ?
		   AND deleted = ?
	  ORDER BY %s
`, orderBy), done, false)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]tonight.Task, 0)
	for rows.Next() {
		var id uint
		var title string
		var description string
		var done bool
		var doneAt *time.Time
		var createdAt time.Time
		if err := rows.Scan(&id, &title, &description, &done, &doneAt, &createdAt); err != nil {
			return nil, err
		}

		task := tonight.Task{
			ID:          id,
			Title:       title,
			Description: description,

			Done:   done,
			DoneAt: doneAt,

			CreatedAt: createdAt,
		}
		tasks = append(tasks, task)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) Create(ctx context.Context, t *tonight.Task) error {
	if t.ID != 0 {
		return errors.New("cannot update a task")
	}

	now := time.Now()
	res, err := r.db.ExecContext(
		ctx,
		"INSERT INTO tasks (title, description, rank, done, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		t.Title, t.Description, 999, t.Done, now, now,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	t.ID = uint(id)
	return nil
}

func (r *TaskRepository) MarkDone(ctx context.Context, taskID uint) error {
	now := time.Now()
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE tasks SET done = ?, doneAt = ?, updated_at = ? WHERE id = ?",
		true, now, now, taskID,
	)
	return err
}

func (r *TaskRepository) Delete(ctx context.Context, taskID uint) error {
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE tasks SET deleted = ? WHERE id = ?",
		true, taskID,
	)
	return err
}

func (r *TaskRepository) UpdateRanks(ctx context.Context, ranks map[uint]uint) error {
	for id, rank := range ranks {
		_, err := r.db.ExecContext(ctx, "UPDATE tasks SET rank = ? WHERE id = ?", rank, id)
		if err != nil {
			return err
		}
	}
	return nil
}
