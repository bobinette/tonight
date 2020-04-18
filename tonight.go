package tonight

import (
	"context"
	"time"

	"github.com/bobinette/tonight/auth"
	"github.com/google/uuid"
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
	Upsert(ctx context.Context, p Project, u auth.User) error
	List(ctx context.Context, u auth.User) ([]Project, error)
	Get(ctx context.Context, uuid uuid.UUID, u auth.User) (Project, error)

	Find(ctx context.Context, slug string, u auth.User) (Project, error)
}

type UserStore interface {
	Ensure(ctx context.Context, user *auth.User) error
	Permission(ctx context.Context, user auth.User, projectUUID uuid.UUID) (string, error)
}
