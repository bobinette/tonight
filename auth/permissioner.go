package auth

import (
	"context"

	"github.com/google/uuid"
)

type Permission string

// List of available project permissions
const (
	ProjectOwn  Permission = "project.own"
	ProjectEdit Permission = "project.edit"
	ProjectView Permission = "project.view"
)

var projectPermissionsInclusion = map[Permission]uint{
	ProjectView: 0,
	ProjectEdit: 1,
	ProjectOwn:  2,
}

func (p Permission) Includes(o Permission) bool {
	pv, ok := projectPermissionsInclusion[p]
	if !ok {
		return false
	}

	ov, ok := projectPermissionsInclusion[o]
	if !ok {
		return false
	}

	return pv >= ov
}

type Store interface {
	Ensure(ctx context.Context, user *User) error

	SetPermission(ctx context.Context, user User, projectUUID uuid.UUID, perm Permission) error
	AllPermissions(ctx context.Context, user User) (map[uuid.UUID]Permission, error)
	Permission(ctx context.Context, user User, projectUUID uuid.UUID) (Permission, error)
}

type Permissioner struct {
	store Store
}

func NewPermissioner(s Store) Permissioner {
	return Permissioner{store: s}
}

func (p Permissioner) HasPermission(ctx context.Context, user User, projectUUID uuid.UUID, perm Permission) error {
	userPerm, err := p.store.Permission(ctx, user, projectUUID)
	if err != nil {
		return err
	}

	if !userPerm.Includes(perm) {
		return ErrInsufficientPermissions
	}
	return nil
}

func (p Permissioner) AllowProject(ctx context.Context, user User, projectUUID uuid.UUID, perm Permission) error {
	return p.store.SetPermission(ctx, user, projectUUID, perm)
}

func (p Permissioner) AllowedProjects(ctx context.Context, user User, perm Permission) ([]uuid.UUID, error) {
	permissions, err := p.store.AllPermissions(ctx, user)
	if err != nil {
		return nil, err
	}

	projectUUIDs := make([]uuid.UUID, 0, len(permissions))
	for projectUUID, projectPerm := range permissions {
		if projectPerm.Includes(perm) {
			projectUUIDs = append(projectUUIDs, projectUUID)
		}
	}

	return projectUUIDs, nil
}
