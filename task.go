package tonight

import (
	"context"
	"time"
)

type LogType string

const (
	LogTypeCompletion LogType = "COMPLETION"
	LogTypePause              = "PAUSE"
	LogTypeStart              = "START"
	LogTypeLog                = "LOG"
)

type Task struct {
	ID          uint
	Title       string
	Description string

	Priority int
	Rank     uint
	Tags     []string

	Duration string
	Deadline *time.Time

	Log []Log

	Dependencies []Dependency

	CreatedAt time.Time
}

func (t Task) Completion() int {
	c := 0
	for _, log := range t.Log {
		if log.Completion > c {
			c = log.Completion
		}
	}
	return c
}

func (t Task) Done() bool {
	for _, log := range t.Log {
		if log.Completion == 100 {
			return true
		}
	}
	return false
}

func (t Task) DoneAt() *time.Time {
	for _, log := range t.Log {
		if log.Completion == 100 {
			doneAt := log.CreatedAt
			return &doneAt
		}
	}
	return nil
}

type Log struct {
	Type        LogType
	Completion  int
	Description string

	CreatedAt time.Time
}

type Dependency struct {
	ID   uint
	Done bool
}

type TaskRepository interface {
	List(ctx context.Context, ids []uint) ([]Task, error)
	Create(ctx context.Context, t *Task) error
	Update(ctx context.Context, t *Task) error

	Log(ctx context.Context, taskID uint, log Log) error
	UpdateRanks(ctx context.Context, ranks map[uint]uint) error

	Delete(ctx context.Context, taskID uint) error
}

type TaskIndex interface {
	Search(ctx context.Context, q string, done bool, allowedIDs []uint) ([]uint, error)
	Index(ctx context.Context, task Task) error
	Delete(ctx context.Context, taskID uint) error
}

type taskService struct {
	repo  TaskRepository
	index TaskIndex

	userRepo UserRepository
}

func (ts *taskService) list(ctx context.Context, user User, q string, done bool) ([]Task, error) {
	ids, err := ts.index.Search(ctx, q, done, user.TaskIDs)
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
	task := parse(input)

	if err := ts.repo.Create(ctx, &task); err != nil {
		return Task{}, err
	}

	if err := ts.index.Index(ctx, task); err != nil {
		return Task{}, err
	}

	if err := ts.userRepo.AddTaskToUser(ctx, user.ID, task.ID); err != nil {
		return Task{}, err
	}

	return task, nil
}

func (ts *taskService) update(ctx context.Context, taskID uint, input string) (Task, error) {
	task := parse(input)
	task.ID = taskID

	if err := ts.repo.Update(ctx, &task); err != nil {
		return Task{}, err
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
	}

	// Ensure completion does not go down
	task := tasks[0]
	for _, l := range task.Log {
		if l.Completion > log.Completion {
			log.Completion = l.Completion
		}
	}

	if err := ts.repo.Log(ctx, taskID, log); err != nil {
		return Task{}, err
	}

	task.Log = append(task.Log, log)
	if err := ts.index.Index(ctx, task); err != nil {
		return Task{}, err
	}

	return task, nil
}
