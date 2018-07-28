package tonight

import (
	"context"
	"fmt"
	"time"
)

type Planning struct {
	ID uint `json:"id"`

	Q        string        `json:"q"`
	Duration time.Duration `json:"duration"`
	Strict   bool          `json:"strict"`

	Dismissed bool      `json:"dismissed"`
	StartedAt time.Time `json:"startedAt"`

	Tasks []Task `json:"tasks"`
}

func (p Planning) Done() bool {
	for _, task := range p.Tasks {
		if task.Done() == StatusPending {
			return false
		}
	}
	return true
}

func (p Planning) TotalDuration() time.Duration {
	var totalDuration time.Duration = 0
	for _, task := range p.Tasks {
		totalDuration += task.LeftDuration()
	}

	return totalDuration
}

type PlanningRepository interface {
	Get(ctx context.Context, userID uint) (Planning, error)
	Create(ctx context.Context, userID uint, planning *Planning) error
	Update(ctx context.Context, userID uint, planning *Planning) error
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

	if len(planning.Tasks) > 0 && planning.Done() {
		// recreate a planning with the same search
		input := planning.Duration.String()
		if planning.Q != "" {
			input = fmt.Sprintf("%s for %s", planning.Q, input)
		}
		return ps.plan(ctx, user, input)
	}

	return planning, nil
}

func (ps *planningService) plan(ctx context.Context, user User, input string) (Planning, error) {
	q, d, strict, err := parsePlanning(input)
	if err != nil {
		return Planning{}, err
	}

	ids, err := ps.taskIndex.Search(ctx, TaskSearchParameters{
		Q:        q,
		Statuses: []Status{StatusPending},
		IDs:      user.TaskIDs,
	})
	if err != nil {
		return Planning{}, err
	}

	tasks, err := ps.taskRepo.List(ctx, ids)
	if err != nil {
		return Planning{}, err
	}

	planned := plan(tasks, d, strict)

	planning := Planning{
		Q:        q,
		Duration: d,
		Strict:   strict,

		Tasks:     planned,
		StartedAt: time.Now(),
	}

	if err := ps.repo.Create(ctx, user.ID, &planning); err != nil {
		return Planning{}, err
	}

	return planning, err
}

func (ps *planningService) doLater(ctx context.Context, user User, taskID uint) (Planning, error) {
	planning, err := ps.current(ctx, user)
	if err != nil {
		return Planning{}, err
	}

	ids, err := ps.taskIndex.Search(ctx, TaskSearchParameters{
		Q:        planning.Q,
		Statuses: []Status{StatusPending},
		IDs:      user.TaskIDs,
	})
	if err != nil {
		return Planning{}, err
	}

	tasks, err := ps.taskRepo.List(ctx, ids)
	if err != nil {
		return Planning{}, err
	}

	planned := make([]Task, 0)
	for _, task := range planning.Tasks {
		if task.ID != taskID {
			planned = append(planned, task)
		}
	}

	tasks = planNext(tasks, planning, taskID)
	planned = append(planned, tasks...)
	planning.Tasks = planned

	if err := ps.repo.Update(ctx, user.ID, &planning); err != nil {
		return Planning{}, err
	}

	return planning, nil
}

func (ps *planningService) dismiss(ctx context.Context, user User) error {
	if err := ps.repo.Dismiss(ctx, user.ID); err != nil {
		return err
	}
	return nil
}
