package tonight

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"

	"github.com/bobinette/tonight/auth"
	"github.com/bobinette/tonight/events"
)

const (
	emptyUUID = "00000000-0000-0000-0000-000000000000"
)

func RegisterHTTP(
	srv *echo.Group,
	eventStore events.Store,
	taskStore TaskStore,
	projectStore ProjectStore,
	releaseStore ReleaseStore,
	userStore UserStore,
	permissioner auth.Permissioner,
) error {
	s := newService(eventStore, taskStore, projectStore, releaseStore, userStore)

	withEvent := func(t events.EventType) echo.MiddlewareFunc {
		return events.Middleware(t, eventStore)
	}

	srv.POST("/projects", s.createProject)
	srv.GET("/projects", s.listProjects)
	srv.GET("/projects/:uuid", s.getProject)
	srv.GET("/projects/slug/:slug", s.findProject)
	srv.POST("/projects/:uuid", s.updateProject)
	srv.POST("/projects/:uuid/tasks/ranks", s.rankTasks)

	// releases.go
	r := releaseCRUD{
		store:        releaseStore,
		permissioner: permissioner,
	}
	srv.POST("/projects/:project_uuid/releases", r.create, withEvent(events.ReleaseCreate))

	// tasks.go
	t := taskCRUD{
		taskStore:    taskStore,
		permissioner: permissioner,
	}
	srv.POST("/tasks", t.create, withEvent(events.TaskCreate))
	srv.POST("/tasks/:uuid", t.update, withEvent(events.TaskUpdate))
	srv.DELETE("/tasks/:uuid", t.delete, withEvent(events.TaskDelete))
	srv.POST("/tasks/:uuid/done", t.markAsDone, withEvent(events.TaskDone))
	srv.POST("/tasks/ranks", t.reorder, withEvent(events.TasksReorder))

	return nil
}

type service struct {
	eventStore   events.Store
	taskStore    TaskStore
	projectStore ProjectStore
	releaseStore ReleaseStore
	userStore    UserStore
}

func newService(
	eventStore events.Store,
	taskStore TaskStore,
	projectStore ProjectStore,
	releaseStore ReleaseStore,
	userStore UserStore,
) service {
	return service{
		eventStore:   eventStore,
		taskStore:    taskStore,
		projectStore: projectStore,
		releaseStore: releaseStore,
		userStore:    userStore,
	}
}

func (s service) createTask(c echo.Context) error {
	defer c.Request().Body.Close()

	var t Task
	interceptor := payloadInterceptor{
		v: &t,
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&interceptor); err != nil {
		return err
	}

	if t.UUID.String() != "" && t.UUID.String() != emptyUUID {
		return errors.New("uuid should be empty")
	}

	ctx := c.Request().Context()

	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}
	if err := s.userStore.Ensure(ctx, &user); err != nil {
		return fmt.Errorf("error ensuring user: %w", err)
	}

	projectUUID, err := uuid.Parse(c.Param("project_uuid"))
	if err != nil {
		return err
	}

	perm, err := s.userStore.Permission(ctx, user, projectUUID)
	if err != nil {
		return err
	}
	if perm != "owner" {
		return errors.New("insufficient permissions")
	}

	releaseUUID, err := uuid.Parse(c.Param("release_uuid"))
	if err != nil {
		return err
	}
	release, err := s.releaseStore.Get(ctx, releaseUUID)
	if err != nil {
		return err
	}
	if release.Project.UUID.String() != projectUUID.String() {
		return errors.New("release not found")
	}

	eventUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	now := time.Unix(eventUUID.Time().UnixTime())
	evt := events.Event{
		UUID:       eventUUID,
		Type:       events.TaskCreate,
		EntityUUID: eventUUID,
		UserID:     user.ID,
		Payload:    interceptor.raw,
		CreatedAt:  now,
	}
	if err := s.eventStore.Store(ctx, evt); err != nil {
		return err
	}

	t.UUID = eventUUID
	t.Status = TaskStatusTODO
	t.Release.UUID = releaseUUID
	t.CreatedAt = now
	t.UpdatedAt = now
	if err := s.taskStore.Upsert(ctx, t); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": t,
	})
}

func (s service) createProject(c echo.Context) error {
	defer c.Request().Body.Close()

	var project Project
	interceptor := payloadInterceptor{
		v: &project,
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&interceptor); err != nil {
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
	if err := s.userStore.Ensure(ctx, &user); err != nil {
		return fmt.Errorf("error ensuring user: %w", err)
	}

	eventUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	now := time.Unix(eventUUID.Time().UnixTime())
	evt := events.Event{
		UUID:       eventUUID,
		Type:       events.ProjectCreate,
		EntityUUID: eventUUID,
		UserID:     user.ID,
		Payload:    interceptor.raw,
		CreatedAt:  now,
	}
	if err := s.eventStore.Store(ctx, evt); err != nil {
		return fmt.Errorf("error storing event: %w", err)
	}

	project.UUID = eventUUID
	project.Slug = fmt.Sprintf("%s-%s", slug.Make(project.Name), eventUUID.String()[:8])
	project.CreatedAt = now
	project.UpdatedAt = now
	if err := s.projectStore.Upsert(ctx, project, user); err != nil {
		return fmt.Errorf("error storing project: %w", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": project,
	})
}

func (s service) updateProject(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return err
	}

	var project Project
	interceptor := payloadInterceptor{
		v: &project,
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&interceptor); err != nil {
		return fmt.Errorf("error deconding request: %w", err)
	}

	if project.UUID.String() != id.String() {
		return fmt.Errorf("invalid data: %w", errors.New("uuids should be the same"))
	}

	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()

	_, err = s.projectStore.Get(ctx, id, user)
	if err != nil {
		return err
	}

	eventUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	now := time.Unix(eventUUID.Time().UnixTime())
	evt := events.Event{
		UUID:       eventUUID,
		Type:       events.ProjectUpdate,
		EntityUUID: id,
		UserID:     user.ID,
		Payload:    interceptor.raw,
		CreatedAt:  now,
	}
	if err := s.eventStore.Store(ctx, evt); err != nil {
		return fmt.Errorf("error storing event: %w", err)
	}

	project.UpdatedAt = now
	if err := s.projectStore.Upsert(ctx, project, user); err != nil {
		return fmt.Errorf("error storing project: %w", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": project,
	})
}

func (s service) findProject(c echo.Context) error {
	slug := c.Param("slug")
	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	project, err := s.projectStore.Find(ctx, slug, user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": project,
	})
}

func (s service) getProject(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return err
	}

	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	project, err := s.projectStore.Get(ctx, id, user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": project,
	})
}

func (s service) listProjects(c echo.Context) error {
	ctx := c.Request().Context()

	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}
	if err := s.userStore.Ensure(ctx, &user); err != nil {
		return fmt.Errorf("error ensuring user: %w", err)
	}

	projects, err := s.projectStore.List(ctx, user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": projects,
	})
}

func (s service) rankTasks(c echo.Context) error {
	defer c.Request().Body.Close()

	projectUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return err
	}

	var body struct {
		Ranks []uuid.UUID `json:"ranks"`
	}
	interceptor := payloadInterceptor{
		v: &body,
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&interceptor); err != nil {
		return fmt.Errorf("error deconding request: %w", err)
	}

	ctx := c.Request().Context()

	user, err := auth.ExtractUser(c)
	if err != nil {
		return err
	}
	if err := s.userStore.Ensure(ctx, &user); err != nil {
		return fmt.Errorf("error ensuring user: %w", err)
	}

	perm, err := s.userStore.Permission(ctx, user, projectUUID)
	if err != nil {
		return err
	}
	if perm != "owner" {
		return errors.New("insufficient permissions")
	}

	eventUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	evt := events.Event{
		UUID:       eventUUID,
		Type:       events.ProjectReorderTasks,
		EntityUUID: projectUUID,
		UserID:     user.ID,
		Payload:    interceptor.raw,
		CreatedAt:  time.Unix(eventUUID.Time().UnixTime()),
	}
	if err := s.eventStore.Store(ctx, evt); err != nil {
		return fmt.Errorf("error storing event: %w", err)
	}

	// Check that all uuids belong to the same project, at that all the uuids
	// of todo tasks have been sent

	if err := s.taskStore.Reorder(ctx, body.Ranks); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": "done"})
}

type payloadInterceptor struct {
	raw []byte

	v interface{}
}

func (i *payloadInterceptor) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, i.v); err != nil {
		return err
	}

	i.raw = b
	return nil
}
