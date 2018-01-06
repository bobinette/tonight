package tonight

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type LogType string

const (
	LogTypeCompletion LogType = "COMPLETION"
	LogTypePause              = "PAUSE"
	LogTypeStart              = "START"
	LogTypeLog                = "LOG"
	LogTypeWontDo             = "WONT_DO"
)

type DoneStatus int

const (
	DoneStatusPending DoneStatus = iota
	DoneStatusDone
	DoneStatusWontDo
)

func DoneStatusFromString(s string) DoneStatus {
	switch s {
	case "pending":
		return DoneStatusPending
	case "done":
		return DoneStatusDone
	case "won't do":
		return DoneStatusWontDo
	}
	return DoneStatusPending
}

func (ds DoneStatus) String() string {
	switch ds {
	case DoneStatusPending:
		return "pending"
	case DoneStatusDone:
		return "done"
	case DoneStatusWontDo:
		return "won't do"
	}
	return ""
}

type Task struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Priority int      `json:"priority"`
	Rank     uint     `json:"rank"`
	Tags     []string `json:"tags"`

	Duration string     `json:"duration"`
	Deadline *time.Time `json:"deadline"`

	Score float64 `json:"score"`

	Log []Log `json:"log"`

	Dependencies []Dependency `json:"dependencies"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (t Task) LeftDuration() time.Duration {
	d, err := time.ParseDuration(t.Duration)
	if err != nil {
		d = time.Hour
	}

	return time.Duration(100-t.Completion()) * d / time.Duration(100)
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

func (t Task) Done() DoneStatus {
	for _, log := range t.Log {
		if log.Completion == 100 {
			return DoneStatusDone
		}

		if log.Type == LogTypeWontDo {
			return DoneStatusWontDo
		}
	}

	return DoneStatusPending
}

func (t Task) DoneAt() *time.Time {
	for _, log := range t.Log {
		if log.Completion == 100 || log.Type == LogTypeWontDo {
			doneAt := log.CreatedAt
			return &doneAt
		}
	}
	return nil
}

type Log struct {
	Type        LogType `json:"type"`
	Completion  int     `json:"completion"`
	Description string  `json:"description"`

	CreatedAt time.Time `json:"createdAt"`
}

type Dependency struct {
	ID    uint   `json:"id"`
	Done  bool   `json:"done"`
	Title string `json:"title"`
}

type TaskRepository interface {
	List(ctx context.Context, ids []uint) ([]Task, error)
	Create(ctx context.Context, t *Task) error
	Update(ctx context.Context, t *Task) error

	Log(ctx context.Context, taskID uint, log Log) error
	UpdateRanks(ctx context.Context, ranks map[uint]uint) error

	DependencyTrees(ctx context.Context, taskID uint) ([]Task, error)

	Delete(ctx context.Context, taskID uint) error
}

type TaskSearchParameters struct {
	IDs      []uint
	Q        string
	Statuses []DoneStatus
	SortBy   string
}

type TaskIndex interface {
	Search(ctx context.Context, p TaskSearchParameters) ([]uint, error)
	Index(ctx context.Context, task Task) error
	Delete(ctx context.Context, taskID uint) error
}

type taskService struct {
	repo  TaskRepository
	index TaskIndex

	userRepo UserRepository
}

func (ts *taskService) list(ctx context.Context, user User, q string, doneStatuses []DoneStatus, sortBy string) ([]Task, error) {
	ids, err := ts.index.Search(ctx, TaskSearchParameters{
		Q:        q,
		Statuses: doneStatuses,
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

	tasks, err = ts.repo.List(ctx, []uint{taskID})
	if err != nil {
		return Task{}, err
	}

	// Ensure completion does not go down
	task = tasks[0]
	if err := ts.index.Index(ctx, task); err != nil {
		return Task{}, err
	}

	return task, nil
}
