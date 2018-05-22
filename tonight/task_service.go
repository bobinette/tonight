package tonight

import (
	"context"
	"fmt"
	"time"
)

type taskService struct {
	repo  TaskRepository
	index TaskIndex

	userRepo UserRepository
}

func (ts *taskService) list(ctx context.Context, user User, q string, Statuses []Status, sortBy string) ([]Task, error) {
	ids, err := ts.index.Search(ctx, TaskSearchParameters{
		Q:        q,
		Statuses: Statuses,
		IDs:      user.TaskIDs,
		SortBy:   sortBy,
	})
	if err != nil {
		return nil, err
	}

	tasks, err := ts.repo.List(ctx, ids)
	if err != nil {
		return nil, err
	}

	return tasks, err
}

func (ts *taskService) create(ctx context.Context, user User, input string) (Task, error) {
	task, err := parse(input)
	if err != nil {
		return Task{}, err
	}

	if err := ts.repo.Create(ctx, &task); err != nil {
		return Task{}, err
	}

	if err := ts.userRepo.AddTaskToUser(ctx, user.ID, task.ID); err != nil {
		return Task{}, err
	}

	tasks, err := ts.repo.DependencyTrees(ctx, task.ID)
	if err != nil {
		return Task{}, err
	}

	scores := scoreMany(tasks, score)
	for taskID, s := range scores {
		for i, task := range tasks {
			if task.ID != taskID {
				continue
			}

			tasks[i].Score = s
		}
	}

	for _, task := range tasks {
		if err := ts.index.Index(ctx, task); err != nil {
			return Task{}, err
		}
	}

	return task, nil
}

func (ts *taskService) update(ctx context.Context, taskID uint, input string) (Task, error) {
	task, err := parse(input)
	if err != nil {
		return Task{}, err
	}
	task.ID = taskID

	if err := ts.repo.Update(ctx, &task); err != nil {
		return Task{}, err
	}

	tasks, err := ts.repo.DependencyTrees(ctx, taskID)
	if err != nil {
		return Task{}, err
	}

	scores := scoreMany(tasks, score)
	for taskID, s := range scores {
		for i, task := range tasks {
			if task.ID != taskID {
				continue
			}

			tasks[i].Score = s
		}
	}

	for _, task := range tasks {
		if err := ts.index.Index(ctx, task); err != nil {
			return Task{}, err
		}
	}

	return task, nil
}

func (ts *taskService) delete(ctx context.Context, taskID uint) error {
	if err := ts.repo.Delete(ctx, taskID); err != nil {
		return err
	}

	if err := ts.index.Delete(ctx, taskID); err != nil {
		return err
	}

	return nil
}

func (ts *taskService) updateRanks(ctx context.Context, ranks map[uint]uint) error {
	if err := ts.repo.UpdateRanks(ctx, ranks); err != nil {
		return err
	}

	for id := range ranks {
		tasks, err := ts.repo.List(ctx, []uint{id})
		if err != nil {
			return err
		}

		if err := ts.index.Index(ctx, tasks[0]); err != nil {
			return err
		}
	}

	return nil
}

func (ts *taskService) log(ctx context.Context, taskID uint, input string) (Task, error) {
	log := parseLog(input)

	tasks, err := ts.repo.List(ctx, []uint{taskID})
	if err != nil {
		return Task{}, err
	} else if len(tasks) == 0 {
		return Task{}, ErrTaskNotFound
	}
	task := tasks[0]

	// Handle special case of postponing...
	if log.Type == LogTypePostpone {

		if _, err := time.Parse("2006-01-02", log.Description); err != nil {
			return Task{}, fmt.Errorf("error decoding date for postponing: %v", err)
		}

		log.Description = fmt.Sprintf("postponed until %s", log.Description)
	} else if log.Type == LogTypeDuration {
		// Check that the duration is valid
		if _, err := time.ParseDuration(log.Description); err != nil {
			return Task{}, err
		}

		oldDuration := task.Duration
		if oldDuration == "" {
			oldDuration = "(none)"
		}
		task.Duration = log.Description
		log.Description = fmt.Sprintf("duration updated: %s -> %s", oldDuration, task.Duration)

		if err := ts.repo.Update(ctx, &task); err != nil {
			return Task{}, err
		}
	}

	if !isTransitionAllowed(task, log.Type) {
		return Task{}, fmt.Errorf("Log type %s not allowed for task %d right now", log.Type, task.ID)
	}

	// Ensure completion does not go down
	for _, l := range task.Log {
		if l.Completion > log.Completion {
			log.Completion = l.Completion
		}
	}

	if err := ts.repo.Log(ctx, taskID, log); err != nil {
		return Task{}, err
	}

	tasks, err = ts.repo.DependencyTrees(ctx, taskID)
	if err != nil {
		return Task{}, err
	}

	scores := scoreMany(tasks, score)
	for taskID, s := range scores {
		for i, task := range tasks {
			if task.ID != taskID {
				continue
			}

			tasks[i].Score = s
		}
	}

	for _, task := range tasks {
		if err := ts.index.Index(ctx, task); err != nil {
			return Task{}, err
		}
	}

	return task, nil
}
