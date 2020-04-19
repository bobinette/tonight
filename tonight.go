package tonight

import (
	"context"

	"github.com/bobinette/tonight/auth"
	"github.com/google/uuid"
)

// The Permissioner is used to control the permissions on projects for users.
type Permissioner interface {
	// AllowProject should register the perm permission for user on the project defined by projectUUID.
	AllowProject(ctx context.Context, user auth.User, projectUUID uuid.UUID, perm auth.Permission) error

	// AllowedProjects should return the list of all the project on which user has the perm permission.
	AllowedProjects(ctx context.Context, user auth.User, perm auth.Permission) ([]uuid.UUID, error)

	// HasPermission should return an error if user does not have the perm permission on the project defined
	// by projectUUID, and nil otherwise.
	HasPermission(ctx context.Context, user auth.User, projectUUID uuid.UUID, perm auth.Permission) error
}
