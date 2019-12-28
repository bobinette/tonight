package tonight

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

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
	userStore UserStore,
) error {
	s := service{
		eventStore:   eventStore,
		taskStore:    taskStore,
		projectStore: projectStore,
		userStore:    userStore,
	}

	srv.POST("/tasks", s.createTask)
	srv.GET("/tasks", s.listTasks)

	srv.POST("/projects", s.createProject)
	srv.GET("/projects", s.listProjects)

	return nil
}

type service struct {
	eventStore   EventStore
	taskStore    TaskStore
	projectStore ProjectStore
	userStore    UserStore
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

	perm, err := s.userStore.Permission(ctx, user, t.Project.UUID.String())
	if err != nil {
		return err
	}
	if perm != "owner" {
		return errors.New("insufficient permissions")
	}

	id := uuid.NewV1()
	now := time.Now()
	evt := Event{
		UUID:      id,
		Type:      TaskCreate,
		UserID:    user.ID,
		Payload:   interceptor.raw,
		CreatedAt: now,
	}
	if err := s.eventStore.Store(ctx, evt); err != nil {
		return err
	}

	t.UUID = id
	t.CreatedAt = now
	t.UpdatedAt = now
	if err := s.taskStore.Upsert(ctx, t); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": t,
	})
}

func (s service) listTasks(c echo.Context) error {
	fmt.Println(c.Request().Header)
	ctx := c.Request().Context()
	tasks, err := s.taskStore.List(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": tasks,
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
		UUID:      id,
		Type:      ProjectCreate,
		UserID:    user.ID,
		Payload:   interceptor.raw,
		CreatedAt: now,
	}
	if err := s.eventStore.Store(ctx, evt); err != nil {
		return fmt.Errorf("error storing event: %w", err)
	}

	project.UUID = id
	project.CreatedAt = now
	project.UpdatedAt = now
	if err := s.projectStore.Upsert(ctx, project, user); err != nil {
		return fmt.Errorf("error storing project: %w", err)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
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
