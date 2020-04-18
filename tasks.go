package tonight

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bobinette/tonight/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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
	Get(ctx context.Context, uuid uuid.UUID, u auth.User) (Task, error)
	Upsert(ctx context.Context, t Task) error
	Delete(ctx context.Context, uuid uuid.UUID) error

	GetProjectUUIDs(ctx context.Context, uuids []uuid.UUID) ([]uuid.UUID, error)

	Reorder(ctx context.Context, rankedUUIDs []uuid.UUID) error
}

type taskCRUD struct {
	taskStore    TaskStore
	permissioner Permissioner
}

func (s taskCRUD) create(c echo.Context) error {
	defer c.Request().Body.Close()

	var task Task
	if err := json.NewDecoder(c.Request().Body).Decode(&task); err != nil {
		return err
	}

	ctx := c.Request().Context()

	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}

	projectUUID := task.Release.Project.UUID
	if err := s.permissioner.HasPermission(ctx, user, projectUUID, auth.ProjectEdit); err != nil {
		return err
	}

	eventUUID, ok := c.Get("event_uuid").(uuid.UUID)
	if !ok {
		return errors.New("cannot set uuid")
	}

	now := time.Unix(eventUUID.Time().UnixTime())
	task.UUID = eventUUID
	task.Status = TaskStatusTODO
	task.CreatedAt = now
	task.UpdatedAt = now
	if err := s.taskStore.Upsert(ctx, task); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": task,
	})
}

func (s taskCRUD) delete(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return err
	}

	ctx := c.Request().Context()

	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}

	task, err := s.taskStore.Get(ctx, id, user)
	if err != nil {
		return fmt.Errorf("error retrieving task: %w", err)
	}
	if task.UUID.String() == emptyUUID {
		return fmt.Errorf("task %s not found", id)
	}

	if err := s.permissioner.HasPermission(ctx, user, task.Release.Project.UUID, auth.ProjectEdit); err != nil {
		return err
	}

	if err := s.taskStore.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting task: %w", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "ok",
	})
}

func (s taskCRUD) update(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return err
	}

	var task Task
	if err := json.NewDecoder(c.Request().Body).Decode(&task); err != nil {
		return err
	}

	if task.UUID.String() != id.String() {
		return fmt.Errorf("invalid data: %w", errors.New("uuids should be the same"))
	}

	if task.Title == "" {
		return errors.New("title cannot be empty")
	}

	ctx := c.Request().Context()

	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}

	existing, err := s.taskStore.Get(ctx, id, user)
	if err != nil {
		return fmt.Errorf("error retrieving task: %w", err)
	}
	if existing.UUID.String() == emptyUUID {
		return fmt.Errorf("task %s not found", id)
	}

	if err = s.permissioner.HasPermission(ctx, user, existing.Release.Project.UUID, auth.ProjectEdit); err != nil {
		return err
	}

	task.UpdatedAt = time.Now()
	if err := s.taskStore.Upsert(ctx, task); err != nil {
		return fmt.Errorf("error updating task: %w", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": task,
	})
}

func (s taskCRUD) markAsDone(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return err
	}

	ctx := c.Request().Context()

	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}

	task, err := s.taskStore.Get(ctx, id, user)
	if err != nil {
		return fmt.Errorf("error retrieving task: %w", err)
	}
	if task.UUID.String() == emptyUUID {
		return fmt.Errorf("task %s not found", id)
	}

	if err = s.permissioner.HasPermission(ctx, user, task.Release.Project.UUID, auth.ProjectEdit); err != nil {
		return err
	}

	task.Status = TaskStatusDONE
	task.UpdatedAt = time.Now()
	if err := s.taskStore.Upsert(ctx, task); err != nil {
		return fmt.Errorf("error updating task: %w", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": task,
	})
}

func (s taskCRUD) reorder(c echo.Context) error {
	defer c.Request().Body.Close()

	var body struct {
		Ranks []uuid.UUID `json:"ranks"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return fmt.Errorf("error deconding request: %w", err)
	}

	if len(body.Ranks) == 0 {
		return c.NoContent(http.StatusOK)
	}

	ctx := c.Request().Context()

	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}

	projectUUIDs, err := s.taskStore.GetProjectUUIDs(ctx, body.Ranks)
	if len(projectUUIDs) == 0 {
		return auth.ErrInsufficientPermissions
	}

	projectUUID := projectUUIDs[0]
	for _, o := range projectUUIDs {
		if o != projectUUID {
			return errors.New("tasks do not belong to the same project") // 400
		}
	}

	if err = s.permissioner.HasPermission(ctx, user, projectUUID, auth.ProjectEdit); err != nil {
		return err
	}

	if err := s.taskStore.Reorder(ctx, body.Ranks); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
