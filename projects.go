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
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
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
	Upsert(ctx context.Context, p Project) error
	List(ctx context.Context, uuids []uuid.UUID) ([]Project, error)
	Get(ctx context.Context, uuid uuid.UUID) (Project, error)

	Find(ctx context.Context, slug string) (Project, error)
}

type projectCRUD struct {
	store        ProjectStore
	permissioner Permissioner
}

func (p projectCRUD) create(c echo.Context) error {
	defer c.Request().Body.Close()

	var project Project
	if err := json.NewDecoder(c.Request().Body).Decode(&project); err != nil {
		return fmt.Errorf("error deconding request: %w", err)
	}

	if project.UUID.String() != "" && project.UUID.String() != emptyUUID {
		return fmt.Errorf("invalid data: %w", errors.New("uuid should be empty"))
	}

	ctx := c.Request().Context()

	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}

	eventUUID, ok := c.Get("event_uuid").(uuid.UUID)
	if !ok {
		return errors.New("cannot set uuid")
	}

	now := time.Unix(eventUUID.Time().UnixTime())

	project.UUID = eventUUID
	project.Slug = computeSlug(project)
	project.CreatedAt = now
	project.UpdatedAt = now
	if err := p.store.Upsert(ctx, project); err != nil {
		return fmt.Errorf("error storing project: %w", err)
	}

	if err := p.permissioner.AllowProject(ctx, user, project.UUID, auth.ProjectOwn); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": project,
	})
}

func (p projectCRUD) update(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return err
	}

	var project Project
	if err := json.NewDecoder(c.Request().Body).Decode(&project); err != nil {
		return fmt.Errorf("error deconding request: %w", err)
	}

	if project.UUID.String() != id.String() {
		return fmt.Errorf("invalid data: %w", errors.New("uuids should be the same"))
	}
	if project.Name == "" {
		return errors.New("name cannot be empty")
	}

	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()

	if err := p.permissioner.HasPermission(ctx, user, project.UUID, auth.ProjectEdit); err != nil {
		return err
	}

	existing, err := p.store.Get(ctx, project.UUID)
	if err != nil {
		return err
	}

	existing.Name = project.Name
	existing.Description = project.Description
	existing.Slug = computeSlug(existing)
	existing.UpdatedAt = time.Now()
	if err := p.store.Upsert(ctx, existing); err != nil {
		return fmt.Errorf("error storing project: %w", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": project,
	})
}

func (p projectCRUD) find(c echo.Context) error {
	slug := c.Param("slug")

	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	project, err := p.store.Find(ctx, slug)
	if err != nil {
		return err
	}

	if err := p.permissioner.HasPermission(ctx, user, project.UUID, auth.ProjectEdit); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": project,
	})
}

func (p projectCRUD) get(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return err
	}

	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	project, err := p.store.Get(ctx, id)
	if err != nil {
		return err
	}

	if err := p.permissioner.HasPermission(ctx, user, project.UUID, auth.ProjectEdit); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": project,
	})
}

func (p projectCRUD) list(c echo.Context) error {
	ctx := c.Request().Context()

	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}

	allowedUUIDs, err := p.permissioner.AllowedProjects(ctx, user, auth.ProjectEdit)
	if err != nil {
		return err
	}

	projects, err := p.store.List(ctx, allowedUUIDs)
	if err != nil {
		return err
	}

	// Ensure non-nil
	if projects == nil {
		projects = make([]Project, 0)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": projects,
	})
}

func computeSlug(p Project) string {
	return fmt.Sprintf("%s-%s", slug.Make(p.Name), p.UUID.String()[:8])
}
