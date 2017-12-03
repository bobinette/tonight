package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/bobinette/tonight"
)

type PlanningRepository struct {
	db *sql.DB

	taskRepo tonight.TaskRepository
}

func NewPlanningRepository(db *sql.DB, taskRepo tonight.TaskRepository) *PlanningRepository {
	return &PlanningRepository{
		db:       db,
		taskRepo: taskRepo,
	}
}

func (pr *PlanningRepository) Create(ctx context.Context, userID uint, duration string, taskIDs []uint) (tonight.Planning, error) {
	now := time.Now()
	res, err := pr.db.ExecContext(ctx, `
        INSERT INTO planning (user_id, duration, startedAt, dismissed) VALUES (?, ?, ?, ?)
        `, userID, duration, now, false)
	if err != nil {
		return tonight.Planning{}, err
	}

	planningID, err := res.LastInsertId()
	if err != nil {
		return tonight.Planning{}, err
	}

	for rank, taskID := range taskIDs {
		_, err := pr.db.ExecContext(
			ctx,
			"INSERT INTO planning_has_task (planning_id, rank, task_id) VALUE (?, ?, ?)",
			planningID, rank, taskID,
		)
		if err != nil {
			return tonight.Planning{}, err
		}
	}

	tasks, err := pr.taskRepo.List(ctx, taskIDs)
	if err != nil {
		return tonight.Planning{}, err
	}

	dur, _ := time.ParseDuration(duration)

	planning := tonight.Planning{
		ID: uint(planningID),

		Duration: dur,

		StartedAt: now,
		Dismissed: false,

		Tasks: tasks,
	}

	return planning, nil
}

func (pr *PlanningRepository) Get(ctx context.Context, userID uint) (tonight.Planning, error) {
	row := pr.db.QueryRowContext(
		ctx,
		`SELECT id, duration, startedAt, dismissed FROM planning
        WHERE user_id = ?
        ORDER BY startedAt DESC LIMIT 1
        `, userID,
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

	dur, _ := time.ParseDuration(duration)

	planning := tonight.Planning{
		ID:        id,
		Duration:  dur,
		StartedAt: startedAt,
		Dismissed: dismissed,
	}

	rows, err := pr.db.QueryContext(ctx, "SELECT task_id FROM planning_has_task WHERE planning_id = ?", id)
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

	tasks, err := pr.taskRepo.List(ctx, taskIDs)
	if err != nil {
		return tonight.Planning{}, err
	}

	planning.Tasks = tasks
	return planning, nil
}

func (pr *PlanningRepository) Dismiss(ctx context.Context, userID uint) error {
	row := pr.db.QueryRowContext(
		ctx,
		"SELECT id FROM planning WHERE user_id = ? ORDER BY startedAt DESC LIMIT 1",
		userID,
	)

	var id uint
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	_, err := pr.db.ExecContext(ctx, "UPDATE planning SET dismissed = ? WHERE id = ?", true, id)
	return err
}
