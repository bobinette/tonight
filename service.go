package tonight

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gosimple/slug"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

const (
	emptyUUID = "00000000-0000-0000-0000-000000000000"
)

func RegisterHTTP(
	srv *echo.Group,
	eventStore EventStore,
	taskStore TaskStore,
	projectStore ProjectStore,
	releaseStore ReleaseStore,
	userStore UserStore,
) error {
	s := newService(eventStore, taskStore, projectStore, releaseStore, userStore)
	releaseSrv := releaseService{
		store:      releaseStore,
		eventStore: eventStore,
		userStore:  userStore,
	}
	taskSrv := taskService{
		taskStore:    taskStore,
		releaseStore: releaseStore,
		eventStore:   eventStore,
		userStore:    userStore,
	}

	srv.POST("/tasks/:uuid", s.updateTask)
	srv.DELETE("/tasks/:uuid", taskSrv.delete)
	srv.POST("/tasks/:uuid/done", s.markAsDone)

	srv.POST("/projects", s.createProject)
	srv.GET("/projects", s.listProjects)
	srv.GET("/projects/:uuid", s.getProject)
	srv.GET("/projects/slug/:slug", s.findProject)
	srv.POST("/projects/:uuid", s.updateProject)
	srv.POST("/projects/:uuid/tasks/ranks", s.rankTasks)

	srv.POST("/projects/:project_uuid/releases", releaseSrv.create)
	srv.POST("/projects/:project_uuid/releases/:release_uuid/tasks", s.createTask)

	return nil
}

type service struct {
	eventStore   EventStore
	taskStore    TaskStore
	projectStore ProjectStore
	releaseStore ReleaseStore
	userStore    UserStore
}

func newService(
	eventStore EventStore,
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

	releaseUUID, err := uuid.FromString(c.Param("release_uuid"))
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

	id := uuid.NewV1()
	now := time.Now()
	evt := Event{
		UUID:       id,
		Type:       TaskCreate,
		EntityUUID: id,
		UserID:     user.ID,
		Payload:    interceptor.raw,
		CreatedAt:  now,
	}
	if err := s.eventStore.Store(ctx, evt); err != nil {
		return err
	}

	t.UUID = id
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

func (s service) updateTask(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := uuid.FromString(c.Param("uuid"))
	if err != nil {
		return err
	}

	var t Task
	interceptor := payloadInterceptor{
		v: &t,
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&interceptor); err != nil {
		return err
	}

	if t.UUID.String() != id.String() {
		return fmt.Errorf("invalid data: %w", errors.New("uuids should be the same"))
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

	if _, err = s.projectStore.Get(ctx, release.Project.UUID, user); err != nil {
		return err
	}

	if t.Title == "" {
		return errors.New("title cannot be empty")
	}

	eventUUID := uuid.NewV1()
	now := time.Now()
	evt := Event{
		UUID:       eventUUID,
		Type:       TaskUpdate,
		EntityUUID: id,
		UserID:     user.ID,
		Payload:    interceptor.raw,
		CreatedAt:  now,
	}
	if err := s.eventStore.Store(ctx, evt); err != nil {
		return fmt.Errorf("error storing event: %w", err)
	}

	t.UpdatedAt = time.Now()
	if err := s.taskStore.Upsert(ctx, t); err != nil {
		return fmt.Errorf("error updating task: %w", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "ok",
	})
}

func (s service) markAsDone(c echo.Context) error {
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

	if _, err = s.projectStore.Get(ctx, release.Project.UUID, user); err != nil {
		return err
	}

	eventUUID := uuid.NewV1()
	now := time.Now()
	evt := Event{
		UUID:       eventUUID,
		Type:       TaskDone,
		EntityUUID: id,
		UserID:     user.ID,
		Payload:    []byte("{}"),
		CreatedAt:  now,
	}
	if err := s.eventStore.Store(ctx, evt); err != nil {
		return fmt.Errorf("error storing event: %w", err)
	}

	task.Status = TaskStatusDONE
	task.UpdatedAt = time.Now()
	if err := s.taskStore.Upsert(ctx, task); err != nil {
		return fmt.Errorf("error updating task: %w", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "ok",
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

	user, err := userFromHeader(c)
	if err != nil {
		return err
	}
	if err := s.userStore.Ensure(ctx, &user); err != nil {
		return fmt.Errorf("error ensuring user: %w", err)
	}

	id := uuid.NewV1()
	now := time.Now()
	evt := Event{
		UUID:       id,
		Type:       ProjectCreate,
		EntityUUID: id,
		UserID:     user.ID,
		Payload:    interceptor.raw,
		CreatedAt:  now,
	}
	if err := s.eventStore.Store(ctx, evt); err != nil {
		return fmt.Errorf("error storing event: %w", err)
	}

	project.UUID = id
	project.Slug = fmt.Sprintf("%s-%s", slug.Make(project.Name), id.String()[:8])
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

	id, err := uuid.FromString(c.Param("uuid"))
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

	user, err := userFromHeader(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()

	_, err = s.projectStore.Get(ctx, id, user)
	if err != nil {
		return err
	}

	now := time.Now()
	evt := Event{
		UUID:       uuid.NewV1(),
		Type:       ProjectUpdate,
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
	user, err := userFromHeader(c)
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
	id, err := uuid.FromString(c.Param("uuid"))
	if err != nil {
		return err
	}

	user, err := userFromHeader(c)
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

	user, err := userFromHeader(c)
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

	projectUUID, err := uuid.FromString(c.Param("uuid"))
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

	user, err := userFromHeader(c)
	if err != nil {
		return err
	}
	if err := s.userStore.Ensure(ctx, &user); err != nil {
		return fmt.Errorf("error ensuring user: %w", err)
	}

	perm, err := s.userStore.Permission(ctx, user, projectUUID.String())
	if err != nil {
		return err
	}
	if perm != "owner" {
		return errors.New("insufficient permissions")
	}

	evt := Event{
		UUID:       uuid.NewV1(),
		Type:       ProjectReorderTasks,
		EntityUUID: projectUUID,
		UserID:     user.ID,
		Payload:    interceptor.raw,
		CreatedAt:  time.Now(),
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

func userFromHeader(c echo.Context) (User, error) {
	id := c.Request().Header.Get("Token-Claim-Sub")
	if id == "" {
		return User{}, errors.New("no user")
	}

	name := c.Request().Header.Get("Token-Claim-Name")

	return User{
		ID:   id,
		Name: name,
	}, nil
}
