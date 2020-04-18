package tonight

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/bobinette/tonight/auth"
)

type Release struct {
	UUID uuid.UUID `json:"uuid"`

	Title       string `json:"title"`
	Description string `json:"description"`

	Project Project `json:"project"`
	Tasks   []Task  `json:"tasks"`

	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

type ReleaseStore interface {
	Get(ctx context.Context, id uuid.UUID) (Release, error)
	List(ctx context.Context, projectUUID uuid.UUID) ([]Release, error)

	Upsert(ctx context.Context, release Release) error
}

type Permissioner interface {
	HasPermission(ctx context.Context, user auth.User, projectUUID uuid.UUID, perm auth.Permission) error
}

type releaseCRUD struct {
	store        ReleaseStore
	permissioner Permissioner
}

func (s *releaseCRUD) create(c echo.Context) error {
	defer c.Request().Body.Close()

	var release Release
	if err := json.NewDecoder(c.Request().Body).Decode(&release); err != nil {
		fmt.Println("trololo")
		return err
	}

	ctx := c.Request().Context()

	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}

	projectUUID, err := uuid.Parse(c.Param("project_uuid"))
	if err != nil {
		return err
	}

	if err := s.permissioner.HasPermission(ctx, user, projectUUID, auth.ProjectEdit); err != nil {
		return err
	}

	eventUUID, ok := c.Get("event_uuid").(uuid.UUID)
	if !ok {
		return errors.New("cannot set uuid")
	}

	now := time.Unix(eventUUID.Time().UnixTime())
	release.UUID = eventUUID
	release.Project.UUID = projectUUID
	release.CreatedAt = now
	release.UpdatedAt = now
	if err := s.store.Upsert(ctx, release); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": release,
	})
}
