package tonight

import (
	"context"
	"time"
)

type Task struct {
	ID          uint
	Title       string
	Description string

	Tags     []string
	Duration string

	Done   bool
	DoneAt *time.Time

	CreatedAt time.Time
}

type TaskRepository interface {
	List(ctx context.Context, done bool) ([]Task, error)
	Create(ctx context.Context, t *Task) error

	MarkDone(ctx context.Context, taskID uint, description string) error
	UpdateRanks(ctx context.Context, ranks map[uint]uint) error

	Delete(ctx context.Context, taskID uint) error
}

type TaskService struct {
	repo TaskRepository
}
