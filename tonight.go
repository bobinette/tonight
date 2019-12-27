package tonight

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

// A Task is the basic object of Tonight.
type Task struct {
	UUID uuid.UUID `json:"uuid"`

	Title string `json:"title"`

	Project Project `json:"project"`

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

// A Project groups tasks.
type Project struct {
	UUID uuid.UUID `json:"uuid"`

	Name string `json:"name"`

	Tasks []Task `json:"tasks"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// A ProjectStore is responsible for storing projects, typically in a
// database.
type ProjectStore interface {
	Upsert(ctx context.Context, t Project) error
	List(ctx context.Context) ([]Project, error)
	// Get(ctx context.Context, uuid uuid.UUID) (Project, error)
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserStore interface {
	Ensure(ctx context.Context, user *User) error
}
