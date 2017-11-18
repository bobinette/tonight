package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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
		SELECT id, title, description, duration, done, doneAt, created_at
		  FROM tasks
		 WHERE done = ?
		   AND deleted = ?
	  ORDER BY %s
`, orderBy), done, false)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	taskMap := make(map[uint]tonight.Task, 0)
	ids := make([]uint, 0)
	for rows.Next() {
		var id uint
		var title string
		var description string
		var duration string
		var done bool
		var doneAt *time.Time
		var createdAt time.Time
		if err := rows.Scan(&id, &title, &description, &duration, &done, &doneAt, &createdAt); err != nil {
			return nil, err
		}

		task := tonight.Task{
			ID:          id,
			Title:       title,
			Description: description,

			Duration: duration,

			Done:   done,
			DoneAt: doneAt,

			CreatedAt: createdAt,
		}
		taskMap[task.ID] = task
		ids = append(ids, task.ID)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return nil, nil
	}

	marks := make([]string, len(ids))
	params := make([]interface{}, len(ids))
	for i, id := range ids {
		marks[i] = "?"
		params[i] = id
	}
	rows, err = r.db.QueryContext(ctx, fmt.Sprintf(
		fmt.Sprintf("SELECT task_id, tag FROM tags WHERE task_id IN (%s)", strings.Join(marks, ",")),
	), params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make(map[uint][]string)
	for rows.Next() {
		var taskID uint
		var tag string
		if err := rows.Scan(&taskID, &tag); err != nil {
			return nil, err
		}

		tags[taskID] = append(tags[taskID], tag)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	tasks := make([]tonight.Task, len(ids))
	for i, id := range ids {
		task := taskMap[id]
		task.Tags = tags[task.ID]

		tasks[i] = task
	}

	return tasks, nil
}

func (r *TaskRepository) Create(ctx context.Context, t *tonight.Task) error {
	if t.ID != 0 {
		return errors.New("cannot update a task")
	}

	now := time.Now()
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO tasks (title, description, duration, rank, done, created_at, updated_at)
		     VALUES (?, ?, ?, ?, ?, ?, ?)
	`, t.Title, t.Description, t.Duration, 999, t.Done, now, now)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	taskID := uint(id)

	if len(t.Tags) > 0 {
		values := make([]string, len(t.Tags))
		params := make([]interface{}, 2*len(t.Tags))
		for i, tag := range t.Tags {
			values[i] = "(?, ?)"
			params[i*2] = taskID
			params[i*2+1] = tag
		}
		res, err = r.db.ExecContext(
			ctx,
			fmt.Sprintf("INSERT INTO tags (task_id, tag) VALUES %s", strings.Join(values, ",")),
			params...,
		)
		if err != nil {
			return err
		}
	}

	t.ID = taskID
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
