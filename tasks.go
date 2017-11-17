package tonight

import (
	"context"
	"time"
)

type Task struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Done   bool       `json:"done"`
	DoneAt *time.Time `json:"doneAt"`

	CreatedAt time.Time `json:"createdAt"`
}

type TaskRepository interface {
	List(ctx context.Context, done bool) ([]Task, error)
	Create(ctx context.Context, t *Task) error

	MarkDone(ctx context.Context, taskID uint) error
	UpdateRanks(ctx context.Context, ranks map[uint]uint) error

	Delete(ctx context.Context, taskID uint) error
}

type TaskService struct {
	repo TaskRepository
}
