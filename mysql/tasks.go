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
		orderBy = "done_at DESC"
	}

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT id, title, description, priority, duration, done, done_at, created_at
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
		var priority int
		var duration string
		var done bool
		var doneAt *time.Time
		var createdAt time.Time
		if err := rows.Scan(&id, &title, &description, &priority, &duration, &done, &doneAt, &createdAt); err != nil {
			return nil, err
		}

		task := tonight.Task{
			ID:          id,
			Title:       title,
			Description: description,

			Priority: priority,
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
		INSERT INTO tasks (title, description, priority, duration, rank, done, created_at, updated_at)
		     VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, t.Title, t.Description, t.Priority, t.Duration, 999, t.Done, now, now)
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

func (r *TaskRepository) MarkDone(ctx context.Context, taskID uint, description string) error {
	now := time.Now()
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE tasks SET done = ?, done_description = ?, done_at = ?, updated_at = ? WHERE id = ?",
		true, description, now, now, taskID,
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

func (r *TaskRepository) StartPlanning(ctx context.Context, duration string, taskIDs []uint) (tonight.Planning, error) {
	now := time.Now()
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO planning (duration, startedAt, dismissed) VALUES (?, ?, ?)
		`, duration, now, false)
	if err != nil {
		return tonight.Planning{}, err
	}

	planningID, err := res.LastInsertId()
	if err != nil {
		return tonight.Planning{}, err
	}

	for rank, taskID := range taskIDs {
		_, err := r.db.ExecContext(
			ctx,
			"INSERT INTO planning_has_task (planning_id, rank, task_id) VALUE (?, ?, ?)",
			planningID, rank, taskID,
		)
		if err != nil {
			return tonight.Planning{}, err
		}
	}

	tasks, err := r.tasks(ctx, taskIDs)
	if err != nil {
		return tonight.Planning{}, err
	}

	planning := tonight.Planning{
		ID: uint(planningID),

		Duration: duration,

		StartedAt: now,
		Dismissed: false,

		Tasks: tasks,
	}

	return planning, nil
}

func (r *TaskRepository) CurrentPlanning(ctx context.Context) (tonight.Planning, error) {
	row := r.db.QueryRowContext(
		ctx,
		"SELECT id, duration, startedAt, dismissed FROM planning ORDER BY startedAt DESC LIMIT 1",
	)

	var id uint
	var duration string
	var startedAt time.Time
	var dismissed bool
	if err := row.Scan(&id, &duration, &startedAt, &dismissed); err != nil {
		if err == sql.ErrNoRows {
			return tonight.Planning{}, nil
		}
		return tonight.Planning{}, err
	}

	if dismissed {
		return tonight.Planning{}, nil
	}

	planning := tonight.Planning{
		ID:        id,
		Duration:  duration,
		StartedAt: startedAt,
		Dismissed: dismissed,
	}

	rows, err := r.db.QueryContext(ctx, "SELECT task_id FROM planning_has_task WHERE planning_id = ?", id)
	if err != nil {
		return tonight.Planning{}, err
	}
	defer rows.Close()

	taskIDs := make([]uint, 0)
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return tonight.Planning{}, err
		}

		taskIDs = append(taskIDs, id)
	}

	if err := rows.Close(); err != nil {
		return tonight.Planning{}, err
	}

	tasks, err := r.tasks(ctx, taskIDs)
	if err != nil {
		return tonight.Planning{}, err
	}

	planning.Tasks = tasks
	return planning, nil
}

func (r *TaskRepository) DismissPlanning(ctx context.Context) error {
	row := r.db.QueryRowContext(
		ctx,
		"SELECT id FROM planning ORDER BY startedAt DESC LIMIT 1",
	)

	var id uint
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	_, err := r.db.ExecContext(ctx, "UPDATE planning SET dismissed = ? WHERE id = ?", true, id)
	return err
}

func (r *TaskRepository) tasks(ctx context.Context, ids []uint) ([]tonight.Task, error) {
	params := make([]interface{}, len(ids))
	for i, id := range ids {
		params[i] = id
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT id, title, description, priority, duration, done, done_at, created_at
		  FROM tasks
		 WHERE id IN (%s)
`, join("?", ",", len(ids))), params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	taskMap := make(map[uint]tonight.Task, 0)
	for rows.Next() {
		var id uint
		var title string
		var description string
		var priority int
		var duration string
		var done bool
		var doneAt *time.Time
		var createdAt time.Time
		if err := rows.Scan(
			&id,
			&title,
			&description,
			&priority,
			&duration,
			&done,
			&doneAt,
			&createdAt,
		); err != nil {
			return nil, err
		}

		task := tonight.Task{
			ID:          id,
			Title:       title,
			Description: description,

			Priority: priority,
			Duration: duration,

			Done:   done,
			DoneAt: doneAt,

			CreatedAt: createdAt,
		}
		taskMap[task.ID] = task
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return nil, nil
	}

	rows, err = r.db.QueryContext(ctx, fmt.Sprintf(
		fmt.Sprintf("SELECT task_id, tag FROM tags WHERE task_id IN (%s)", join("?", ",", len(ids))),
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

func join(s, sep string, n int) string {
	a := make([]string, n)
	for i := 0; i < n; i++ {
		a[i] = s
	}
	return strings.Join(a, sep)
}
