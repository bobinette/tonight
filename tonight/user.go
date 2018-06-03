package tonight

import (
	"context"
)

type User struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"isAdmin"`

	TaskIDs []uint `json:"-"`

	TagColours map[string]string `json:"tagColours"`
}

type UserRepository interface {
	Get(ctx context.Context, id uint) (User, error)
	GetByName(ctx context.Context, name string) (User, error)
	Insert(ctx context.Context, user *User) error

	AddTaskToUser(ctx context.Context, userID uint, taskID uint) error
	UpdateTagColor(ctx context.Context, userID uint, tag string, color string) error
}

type userService struct {
	repo UserRepository
}

func (us *userService) getOrCreate(ctx context.Context, username string) (User, error) {
	// get...
	user, err := us.repo.GetByName(ctx, username)
	if err != nil {
		return User{}, err
	}

	// ...or create
	if user.ID == 0 {
		user = User{
			Name: username,
		}

		if err := us.repo.Insert(ctx, &user); err != nil {
			return User{}, err
		}
	}

	return user, nil
}

func (us *userService) customizeColour(ctx context.Context, user User, tag, colour string) (User, error) {
	err := us.repo.UpdateTagColor(ctx, user.ID, tag, colour)
	if err != nil {
		return User{}, err
	}

	return us.repo.Get(ctx, user.ID)
}
