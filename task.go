package tonight

import (
	"context"
	"html/template"
	"time"
)

type Task struct {
	ID          uint
	Title       string
	Description string

	DescriptionMD template.HTML

	Priority int
	Tags     []string

	Duration string
	Deadline *time.Time

	Done   bool
	DoneAt *time.Time

	Completion int // max([log.Completion for log in Log])
	Log        []Log

	Dependencies []Dependency

	CreatedAt time.Time
}

type Log struct {
	Completion  int
	Description string

	CreatedAt time.Time
}

type Dependency struct {
	ID   uint
	Done bool
}

type Planning struct {
	ID uint

	Duration string

	Dismissed bool
	StartedAt time.Time

	Tasks []Task
}

type TaskRepository interface {
	List(ctx context.Context, done bool) ([]Task, error)
	Create(ctx context.Context, t *Task) error
	Update(ctx context.Context, t *Task) error

	MarkDone(ctx context.Context, taskID uint, log Log) error
	UpdateRanks(ctx context.Context, ranks map[uint]uint) error

	Delete(ctx context.Context, taskID uint) error

	StartPlanning(ctx context.Context, duration string, taskIDs []uint) (Planning, error)
	DismissPlanning(ctx context.Context) error
	CurrentPlanning(ctx context.Context) (Planning, error)
}
