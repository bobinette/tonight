package tonight

import (
	"context"
	"time"
)

type Task struct {
	ID          uint
	Title       string
	Description string

	Priority int
	Tags     []string
	Duration string

	Done   bool
	DoneAt *time.Time

	CreatedAt time.Time
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

	MarkDone(ctx context.Context, taskID uint, description string) error
	UpdateRanks(ctx context.Context, ranks map[uint]uint) error

	Delete(ctx context.Context, taskID uint) error

	StartPlanning(ctx context.Context, duration string, taskIDs []uint) (Planning, error)
	DismissPlanning(ctx context.Context) error
	CurrentPlanning(ctx context.Context) (Planning, error)
}
