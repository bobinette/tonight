package tonight

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type LogType string

const (
	LogTypeProgress LogType = "PROGRESS"
	LogTypePause            = "PAUSE"
	LogTypeStart            = "START"
	LogTypeWontDo           = "WONT_DO"
	LogTypeComment          = "COMMENT"
	LogTypePostpone         = "POSTPONE"
)

type Status int

const (
	StatusPending Status = iota
	StatusDone
	StatusWontDo
)

func StatusFromString(s string) Status {
	switch s {
	case "pending":
		return StatusPending
	case "done":
		return StatusDone
	case "won't do":
		return StatusWontDo
	}
	return StatusPending
}

func (ds Status) String() string {
	switch ds {
	case StatusPending:
		return "pending"
	case StatusDone:
		return "done"
	case StatusWontDo:
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

	Duration       string     `json:"duration"`
	Deadline       *time.Time `json:"deadline"`
	PostponedUntil *time.Time `json:"postponedUntil"`

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

func (t Task) Done() Status {
	for _, log := range t.Log {
		if log.Completion == 100 {
			return StatusDone
		}

		if log.Type == LogTypeWontDo {
			return StatusWontDo
		}
	}

	return StatusPending
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
	Statuses []Status
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
		date, err := time.Parse("2006-01-02", log.Description)
		if err != nil {
			return Task{}, fmt.Errorf("error decoding date for postponing: %v", err)
		}

		task.PostponedUntil = &date
		if err := ts.repo.Update(ctx, &task); err != nil {
			return Task{}, err
		}

		log.Description = fmt.Sprintf("postponed until %s", log.Description)
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
