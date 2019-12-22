package tonight

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	TaskCreate EventType = "TaskCreate"
)

// A Task is the basic object of Tonight.
type Task struct {
	UUID uuid.UUID `json:"uuid"`

	Title string `json:"title"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// A TaskStore is responsible for storing tasks, typically in a
// database.
type TaskStore interface {
	Upsert(ctx context.Context, t Task) error
	List(ctx context.Context) ([]Task, error)
	Get(ctx context.Context, uuid uuid.UUID) (Task, error)
}
