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
) error {
	s := service{
		eventStore: eventStore,
		taskStore:  taskStore,
	}

	srv.POST("/tasks", s.create)
	srv.GET("/tasks", s.get)

	return nil
}

type service struct {
	eventStore EventStore
	taskStore  TaskStore
}

func (s service) create(c echo.Context) error {
	defer c.Request().Body.Close()

	var t Task
	interceptor := payloadInterceptor{
		v: &t,
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&interceptor); err != nil {
		return err
	}

	if t.UUID.String() != "" && t.UUID.String() != emptyUUID {
		fmt.Println(t.UUID)
		return errors.New("uuid should be empty")
	}

	ctx := c.Request().Context()

	id := uuid.NewV1()
	now := time.Now()
	evt := Event{
		UUID:      id,
		Type:      TaskCreate,
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

func (s service) get(c echo.Context) error {
	ctx := c.Request().Context()
	tasks, err := s.taskStore.List(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": tasks,
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
