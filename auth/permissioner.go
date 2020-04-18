package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Permission string

const (
	ProjectEdit Permission = "project.edit"
)

type Permissioner struct {
	s Store
}

func NewPermissioner(s Store) Permissioner {
	return Permissioner{s: s}
}

func (p Permissioner) HasPermission(ctx context.Context, user User, projectUUID uuid.UUID, perm Permission) error {
	userPerm, err := p.s.Permission(ctx, user, projectUUID)
	if err != nil {
		return err
	}

	switch perm {
	case ProjectEdit:
		if userPerm != "owner" {
			fmt.Println(user, projectUUID, userPerm)
			return ErrInsufficientPermissions
		}
	default:
		return fmt.Errorf("unknow permission %s", p)
	}

	return nil
}
