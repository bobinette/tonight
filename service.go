package tonight

import (
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
	permissioner auth.Permissioner,
) error {
	withEvent := func(t events.EventType) echo.MiddlewareFunc {
		return events.Middleware(t, eventStore)
	}

	// Projects
	p := projectCRUD{
		store:        projectStore,
		permissioner: permissioner,
	}
	srv.GET("/projects", p.list)
	srv.POST("/projects", p.create, withEvent(events.ProjectCreate))
	srv.GET("/projects/:uuid", p.get)
	srv.GET("/projects/slug/:slug", p.find)
	srv.POST("/projects/:uuid", p.update, withEvent(events.ProjectUpdate))

	// Releases
	r := releaseCRUD{
		store:        releaseStore,
		permissioner: permissioner,
	}
	srv.POST("/projects/:project_uuid/releases", r.create, withEvent(events.ReleaseCreate))

	// Tasks
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
