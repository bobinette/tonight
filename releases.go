package tonight

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
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

type releaseService struct {
	store      ReleaseStore
	userStore  UserStore
	eventStore EventStore
}

func (s *releaseService) create(c echo.Context) error {
	defer c.Request().Body.Close()

	var release Release
	interceptor := payloadInterceptor{
		v: &release,
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&interceptor); err != nil {
		fmt.Println("trololo")
		return err
	}

	if release.UUID.String() != "" && release.UUID.String() != emptyUUID {
		return errors.New("uuid should be empty")
	}

	ctx := c.Request().Context()

	user, err := userFromHeader(c)
	if err != nil {
		return err
	}
	if err := s.userStore.Ensure(ctx, &user); err != nil {
		return fmt.Errorf("error ensuring user: %w", err)
	}

	projectUUID, err := uuid.FromString(c.Param("project_uuid"))
	if err != nil {
		return err
	}

	perm, err := s.userStore.Permission(ctx, user, projectUUID.String())
	if err != nil {
		return err
	}
	if perm != "owner" {
		return errors.New("insufficient permissions")
	}

	id := uuid.NewV1()
	now := time.Now()
	evt := Event{
		UUID:       id,
		Type:       ReleaseCreate,
		EntityUUID: id,
		UserID:     user.ID,
		Payload:    interceptor.raw,
		CreatedAt:  now,
	}
	if err := s.eventStore.Store(ctx, evt); err != nil {
		return err
	}

	release.UUID = id
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
