package tonight

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

type TaskStatus string

const (
	TaskStatusTODO TaskStatus = "TODO"
	TaskStatusDONE TaskStatus = "DONE"
)

// A Task is the basic object of Tonight.
type Task struct {
	UUID uuid.UUID `json:"uuid"`

	Title  string     `json:"title"`
	Status TaskStatus `json:"status"`

	Release Release `json:"release"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// A TaskStore is responsible for storing tasks, typically in a
// database.
type TaskStore interface {
	Get(ctx context.Context, uuid uuid.UUID, u User) (Task, error)
	Upsert(ctx context.Context, t Task) error
	Delete(ctx context.Context, uuid uuid.UUID) error

	Reorder(ctx context.Context, rankedUUIDs []uuid.UUID) error
}

type taskService struct {
	taskStore    TaskStore
	releaseStore ReleaseStore
	eventStore   EventStore
	userStore    UserStore
}

func (s taskService) delete(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := uuid.FromString(c.Param("uuid"))
	if err != nil {
		return err
	}

	ctx := c.Request().Context()

	user, err := userFromHeader(c)
	if err != nil {
		return err
	}
	if err := s.userStore.Ensure(ctx, &user); err != nil {
		return fmt.Errorf("error ensuring user: %w", err)
	}

	task, err := s.taskStore.Get(ctx, id, user)
	if err != nil {
		return fmt.Errorf("error retrieving task: %w", err)
	}
	if task.UUID.String() == emptyUUID {
		return fmt.Errorf("task %s not found", id)
	}

	release, err := s.releaseStore.Get(ctx, task.Release.UUID)
	if err != nil {
		return err
	}

	perm, err := s.userStore.Permission(ctx, user, release.Project.UUID.String())
	if err != nil {
		return err
	}
	if perm != "owner" {
		return errors.New("insufficient permissions")
	}

	now := time.Now()
	evt := Event{
		UUID:       uuid.NewV1(),
		Type:       TaskDelete,
		EntityUUID: id,
		UserID:     user.ID,
		Payload:    []byte("{}"),
		CreatedAt:  now,
	}
	if err := s.eventStore.Store(ctx, evt); err != nil {
		return fmt.Errorf("error storing event: %w", err)
	}

	if err := s.taskStore.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting task: %w", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "ok",
	})
}
