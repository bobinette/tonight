package tonight

import (
	"context"
	"time"
)

type Planning struct {
	ID uint

	Duration  time.Duration
	Dismissed bool
	StartedAt time.Time

	Tasks []Task
}

func (p Planning) Done() bool {
	for _, task := range p.Tasks {
		if task.Done() == DoneStatusNotDone {
			return false
		}
	}
	return true
}

func (p Planning) TotalDuration() time.Duration {
	var totalDuration time.Duration = 0
	for _, task := range p.Tasks {
		if dur, err := time.ParseDuration(task.Duration); err == nil {
			totalDuration += dur
		} else {
			// By default use 1h for tasks without duration
			// (Yes, it is totally arbitrary)
			totalDuration += 1 * time.Hour
		}
	}

	return totalDuration
}

type PlanningRepository interface {
	Get(ctx context.Context, userID uint) (Planning, error)
	Create(ctx context.Context, userID uint, duration string, taskIDs []uint) (Planning, error)
	Dismiss(ctx context.Context, userID uint) error
}

type planningService struct {
	repo PlanningRepository

	taskRepo  TaskRepository
	taskIndex TaskIndex
}

func (ps *planningService) current(ctx context.Context, user User) (Planning, error) {
	planning, err := ps.repo.Get(ctx, user.ID)
	if err != nil {
		return Planning{}, err
	}

	if planning.Done() {
		return Planning{}, nil
	}

	return planning, nil
}

func (ps *planningService) plan(ctx context.Context, user User, d time.Duration) (Planning, error) {
	ids, err := ps.taskIndex.Search(ctx, "", DoneStatusNotDone, user.TaskIDs)
	if err != nil {
		return Planning{}, err
	}

	tasks, err := ps.taskRepo.List(ctx, ids)
	if err != nil {
		return Planning{}, err
	}

	planned, _ := plan(tasks, d)

	taskIDs := make([]uint, len(planned))
	for i, task := range planned {
		taskIDs[i] = task.ID
	}

	planning, err := ps.repo.Create(ctx, user.ID, formatDuration(d), taskIDs)
	if err != nil {
		return Planning{}, err
	}

	return planning, err
}

func (ps *planningService) dismiss(ctx context.Context, user User) error {
	if err := ps.repo.Dismiss(ctx, user.ID); err != nil {
		return err
	}
	return nil
}
