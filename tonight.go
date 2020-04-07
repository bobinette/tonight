package tonight

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

// A Project groups tasks.
type Project struct {
	UUID uuid.UUID `json:"uuid"`

	Name string `json:"name"`
	Slug string `json:"slug"`

	Description string `json:"description"`

	Releases []Release `json:"releases"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// A ProjectStore is responsible for storing projects, typically in a
// database.
type ProjectStore interface {
	Upsert(ctx context.Context, p Project, u User) error
	List(ctx context.Context, u User) ([]Project, error)
	Get(ctx context.Context, uuid uuid.UUID, u User) (Project, error)

	Find(ctx context.Context, slug string, u User) (Project, error)
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserStore interface {
	Ensure(ctx context.Context, user *User) error
	Permission(ctx context.Context, user User, projectUUID string) (string, error)
}
