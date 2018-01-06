package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bobinette/tonight"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, t *tonight.Task) error {
	if t.ID != 0 {
		return errors.New("cannot update a task")
	}

	row := r.db.QueryRowContext(ctx, "SELECT max(rank) FROM tasks")
	var rankp *uint
	if err := row.Scan(&rankp); err != nil {
		return err
	}
	rank := uint(0)
	if rankp != nil {
		rank = *rankp
	}
	rank++

	now := time.Now()
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO tasks (title, description, priority, duration, deadline, rank, created_at, updated_at)
		     VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, t.Title, t.Description, t.Priority, t.Duration, t.Deadline, rank, now, now)
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
		_, err := r.db.ExecContext(
			ctx,
			fmt.Sprintf("INSERT INTO tags (task_id, tag) VALUES %s", strings.Join(values, ",")),
			params...,
		)
		if err != nil {
			return err
		}
	}

	if len(t.Dependencies) > 0 {
		values := join("(?, ?, ?)", ",", len(t.Dependencies))
		params := make([]interface{}, 3*len(t.Dependencies))
		for i, dep := range t.Dependencies {
			params[i*3+0] = taskID
			params[i*3+1] = dep.ID
			params[i*3+2] = now
		}

		_, err := r.db.ExecContext(
			ctx,
			fmt.Sprintf(`
				INSERT INTO task_dependencies (task_id, dependency_task_id, created_at)
				VALUES %s`,
				values,
			), params...,
		)
		if err != nil {
			return err
		}
	}

	t.ID = taskID
	t.Rank = rank
	t.CreatedAt = now.Round(time.Second)
	t.UpdatedAt = now.Round(time.Second)
	return nil
}

func (r *TaskRepository) Update(ctx context.Context, t *tonight.Task) error {
	if t.ID == 0 {
		return errors.New("cannot insert a task")
	}

	now := time.Now()
	_, err := r.db.ExecContext(ctx, `
		UPDATE tasks
		SET
			title = ?,
			description = ?,
			priority = ?,
			duration = ?,
			deadline = ?,
			updated_at = ?
		WHERE id = ?
	`, t.Title, t.Description, t.Priority, t.Duration, t.Deadline, now, t.ID)
	if err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx, "DELETE FROM tags WHERE task_id = ?", t.ID); err != nil {
		return err
	}

	if len(t.Tags) > 0 {
		values := make([]string, len(t.Tags))
		params := make([]interface{}, 2*len(t.Tags))
		for i, tag := range t.Tags {
			values[i] = "(?, ?)"
			params[i*2] = t.ID
			params[i*2+1] = tag
		}
		_, err := r.db.ExecContext(
			ctx,
			fmt.Sprintf("INSERT INTO tags (task_id, tag) VALUES %s", strings.Join(values, ",")),
			params...,
		)
		if err != nil {
			return err
		}
	}

	if _, err := r.db.ExecContext(ctx, "DELETE FROM task_dependencies WHERE task_id = ?", t.ID); err != nil {
		return err
	}

	if len(t.Dependencies) > 0 {
		values := join("(?, ?, ?)", ",", len(t.Dependencies))
		params := make([]interface{}, 3*len(t.Dependencies))
		for i, dep := range t.Dependencies {
			params[i*3+0] = t.ID
			params[i*3+1] = dep.ID
			params[i*3+2] = now
		}

		_, err := r.db.ExecContext(
			ctx,
			fmt.Sprintf(`
				INSERT INTO task_dependencies (task_id, dependency_task_id, created_at)
				VALUES %s`,
				values,
			), params...,
		)
		if err != nil {
			return err
		}
	}

	tasks, err := r.List(ctx, []uint{t.ID})
	if err != nil {
		return err
	}
	*t = tasks[0]

	return nil
}

func (r *TaskRepository) Log(ctx context.Context, taskID uint, log tonight.Log) error {
	now := time.Now()
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO task_logs (task_id, type, completion, description, created_at)
			VALUES (?, ?, ?, ?, ?)`,
		taskID, string(log.Type), log.Completion, log.Description, now,
	)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `UPDATE tasks SET updated_at = ? WHERE id = ?`, now, taskID)
	if err != nil {
		return err
	}

	return nil
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

func (r *TaskRepository) List(ctx context.Context, ids []uint) ([]tonight.Task, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	params := make([]interface{}, len(ids))
	for i, id := range ids {
		params[i] = id
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT id, title, description, priority, rank, duration, deadline, created_at, updated_at
		  FROM tasks
		 WHERE id IN (%s)
		   AND deleted = 0
`, join("?", ",", len(ids))), params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	idOrder := make(map[uint]int)
	for i, id := range ids {
		idOrder[id] = i
	}

	tasks, err := r.loadTasks(ctx, rows)
	if err != nil {
		return nil, err
	}

	sort.Sort(&keepOrder{
		idOrder: idOrder,
		tasks:   tasks,
	})

	return tasks, nil
}

func (r *TaskRepository) loadTasks(ctx context.Context, rows *sql.Rows) ([]tonight.Task, error) {
	taskMap := make(map[uint]tonight.Task, 0)
	ids := make([]uint, 0)
	for rows.Next() {
		var id uint
		var title string
		var description string
		var priority int
		var rank uint
		var duration string
		var deadline *time.Time
		var createdAt time.Time
		var updatedAt time.Time
		if err := rows.Scan(&id, &title, &description, &priority, &rank, &duration, &deadline, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		task := tonight.Task{
			ID:          id,
			Title:       title,
			Description: description,

			Priority: priority,
			Rank:     rank,

			Duration: duration,
			Deadline: deadline,

			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
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

	marks := join("?", ",", len(ids))
	params := make([]interface{}, len(ids))
	for i, id := range ids {
		params[i] = id
	}

	// Fetch tags
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(
		"SELECT task_id, tag FROM tags WHERE task_id IN (%s)",
		marks,
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

	// Fetch logs
	rows, err = r.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT task_id, type, completion, description, created_at
		FROM task_logs
		WHERE task_id IN (%s)
		ORDER BY task_id, created_at
	`, marks,
	), params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make(map[uint][]tonight.Log)
	for rows.Next() {
		var taskID uint
		var logType tonight.LogType
		var completion int
		var description string
		var createdAt time.Time
		if err := rows.Scan(&taskID, &logType, &completion, &description, &createdAt); err != nil {
			return nil, err
		}

		logs[taskID] = append(logs[taskID], tonight.Log{
			Type:        logType,
			Completion:  completion,
			Description: description,
			CreatedAt:   createdAt,
		})
	}

	// Fetch dependencies
	rows, err = r.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT task_id, dependency_task_id, tasks.title
		FROM task_dependencies
		JOIN tasks ON tasks.id = dependency_task_id
		WHERE task_id IN (%s) AND tasks.deleted = 0
	`, marks,
	), append(params)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dependencies := make(map[uint][]tonight.Dependency)
	dependencyIDs := make([]uint, 0)
	for rows.Next() {
		var taskID uint
		var dependencyID uint
		var title string
		if err := rows.Scan(&taskID, &dependencyID, &title); err != nil {
			return nil, err
		}

		dependencies[taskID] = append(dependencies[taskID], tonight.Dependency{
			ID:    dependencyID,
			Title: title,
		})
		dependencyIDs = append(dependencyIDs, dependencyID)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	dependencyLogs, err := r.loadLogs(ctx, dependencyIDs...)
	if err != nil {
		return nil, err
	}

	tasks := make([]tonight.Task, len(ids))
	for i, id := range ids {
		task := taskMap[id]
		task.Tags = tags[task.ID]
		task.Log = logs[task.ID]

		task.Dependencies = make([]tonight.Dependency, len(dependencies[task.ID]))
		for i, dep := range dependencies[task.ID] {
			dep.Done = false

			for _, log := range dependencyLogs[dep.ID] {
				if log.Completion == 100 {
					dep.Done = true
				}

				if log.Type == tonight.LogTypeWontDo {
					break
				}
			}

			task.Dependencies[i] = dep
		}

		tasks[i] = task
	}

	return tasks, nil
}

func (r *TaskRepository) loadLogs(ctx context.Context, taskIDs ...uint) (map[uint][]tonight.Log, error) {
	if len(taskIDs) == 0 {
		return nil, nil
	}

	marks := join("?", ",", len(taskIDs))
	params := make([]interface{}, len(taskIDs))
	for i, id := range taskIDs {
		params[i] = id
	}

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT task_id, type, completion, description, created_at
		FROM task_logs
		WHERE task_id IN (%s)
		ORDER BY task_id, created_at DESC
	`, marks,
	), params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make(map[uint][]tonight.Log)
	for rows.Next() {
		var taskID uint
		var logType tonight.LogType
		var completion int
		var description string
		var createdAt time.Time
		if err := rows.Scan(&taskID, &logType, &completion, &description, &createdAt); err != nil {
			return nil, err
		}

		logs[taskID] = append(logs[taskID], tonight.Log{
			Type:        logType,
			Completion:  completion,
			Description: description,
			CreatedAt:   createdAt,
		})
	}

	return logs, nil
}

func (r *TaskRepository) All(ctx context.Context) ([]tonight.Task, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT id, title, description, priority, rank, duration, deadline, created_at, updated_at
		  FROM tasks
		  WHERE deleted = ?
		`,
		false,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.loadTasks(ctx, rows)
}

func (r *TaskRepository) DependencyTrees(ctx context.Context, taskID uint) ([]tonight.Task, error) {
	buffer := make(map[uint]struct{})

	ids, err := r.dependentTasksIDs(ctx, []uint{taskID}, buffer)
	if err != nil {
		return nil, err
	}

	return r.List(ctx, ids)
}

func (r *TaskRepository) dependentTasksIDs(ctx context.Context, ids []uint, buffer map[uint]struct{}) ([]uint, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	params := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		if _, ok := buffer[id]; !ok {
			params = append(params, id)
			buffer[id] = struct{}{}
		}
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT task_id
		FROM task_dependencies
		JOIN tasks ON tasks.id = task_id
		WHERE dependency_task_id IN (%s) AND tasks.deleted = 0
	`, join("?", ",", len(ids)),
	), append(params)...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	dependencyIDs := make([]uint, 0)
	dependencyIDsSet := make(map[uint]struct{})
	for rows.Next() {
		var taskID uint
		if err := rows.Scan(&taskID); err != nil {
			return nil, err
		}

		dependencyIDs = append(dependencyIDs, taskID)
		dependencyIDsSet[taskID] = struct{}{}
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	deeperDependencies, err := r.dependentTasksIDs(ctx, dependencyIDs, buffer)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		if _, ok := dependencyIDsSet[id]; !ok {
			dependencyIDs = append(dependencyIDs, id)
		}
	}

	return append(dependencyIDs, deeperDependencies...), nil
}
