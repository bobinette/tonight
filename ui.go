package tonight

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

type planningForTemplate struct {
	Tasks            []Task
	StartedAt        time.Time
	RequiredDuration time.Duration
	TotalDuration    time.Duration
}

func RegisterUIHandler(
	e *echo.Echo,
	jwtKey []byte,
	repo TaskRepository,
	index TaskIndex,
	planningRepo PlanningRepository,
	userRepo UserRepository,
) {
	h := uiHandler{
		taskService: taskService{
			repo:     repo,
			index:    index,
			userRepo: userRepo,
		},
		planningService: planningService{
			repo:      planningRepo,
			taskRepo:  repo,
			taskIndex: index,
		},
		userRepository: userRepo,
	}

	e.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusPermanentRedirect, "/ui/home") })
	e.GET("/home", func(c echo.Context) error { return c.Redirect(http.StatusPermanentRedirect, "/ui/home") })

	uiGroup := e.Group("/ui")
	uiGroup.Use(JWTMiddleware(jwtKey))
	uiGroup.Use(UserMiddleware(jwtKey, userRepo))

	uiGroup.GET("/home", h.home)

	// -- Calls serving html to partially update the page
	uiGroup.GET("/tasks", h.search)
	uiGroup.POST("/tasks", h.create)
	uiGroup.POST("/tasks/:id", h.update)
	uiGroup.POST("/tasks/:id/done", h.log)
	uiGroup.DELETE("/tasks/:id", h.delete)
	uiGroup.GET("/done", h.listDone)
	uiGroup.POST("/ranks", h.updateRanks)

	uiGroup.POST("/plan", h.plan)
	uiGroup.GET("/plan", h.currentPlanning)
	uiGroup.DELETE("/plan", h.dismissPlanning)
}

type uiHandler struct {
	repo  TaskRepository
	index TaskIndex

	userRepository UserRepository

	taskService     taskService
	planningService planningService
}

func (h *uiHandler) home(c echo.Context) error {
	q := c.QueryParam("q")
	ctx := c.Request().Context()

	user, err := loadUser(c)
	if err != nil {
		return err
	}

	tasks, err := h.taskService.list(ctx, user, q, []DoneStatus{DoneStatusNotDone})
	if err != nil {
		return err
	}

	planning, err := h.planningService.current(ctx, user)
	if err != nil {
		return err
	}

	data := struct {
		Tasks    []Task
		Sortable bool
		Planning Planning
		User     User
	}{
		Tasks:    tasks,
		Sortable: true,
		Planning: planning,
		User:     user,
	}
	return c.Render(http.StatusOK, "home", data)
}

func (us *uiHandler) search(c echo.Context) error {
	q := c.QueryParam("q")
	ctx := c.Request().Context()

	user, err := loadUser(c)
	if err != nil {
		return err
	}

	tasks, err := us.taskService.list(ctx, user, q, []DoneStatus{DoneStatusNotDone})
	if err != nil {
		return err
	}

	data := struct {
		Tasks    []Task
		Sortable bool
	}{
		Tasks:    tasks,
		Sortable: q == "",
	}
	return c.Render(http.StatusOK, "tasks", data)
}

func (us *uiHandler) create(c echo.Context) error {
	defer c.Request().Body.Close()

	user, err := loadUser(c)
	if err != nil {
		return err
	}

	var body struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return err
	}
	ctx := c.Request().Context()

	if _, err := us.taskService.create(ctx, user, body.Content); err != nil {
		return err
	}

	if _, err := reloadUser(c, us.userRepository); err != nil {
		return err
	}

	return us.search(c)
}

func (us *uiHandler) update(c echo.Context) error {
	defer c.Request().Body.Close()

	taskIDStr := c.Param("id")
	taskID64, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		return err
	}
	taskID := uint(taskID64)

	if err := checkPermission(c, taskID); err != nil {
		return err
	}

	var body struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return err
	}

	ctx := c.Request().Context()
	if _, err := us.taskService.update(ctx, taskID, body.Content); err != nil {
		return err
	}

	return us.search(c)
}

func (us *uiHandler) log(c echo.Context) error {
	defer c.Request().Body.Close()

	taskIDStr := c.Param("id")
	taskID64, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		return err
	}
	taskID := uint(taskID64)

	if err := checkPermission(c, taskID); err != nil {
		return err
	}

	var body struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return err
	}

	ctx := c.Request().Context()

	if _, err := us.taskService.log(ctx, taskID, body.Description); err != nil {
		return err
	}

	return us.search(c)
}

func (us *uiHandler) listDone(c echo.Context) error {
	ctx := c.Request().Context()

	user, err := loadUser(c)
	if err != nil {
		return err
	}

	tasks, err := us.taskService.list(ctx, user, "", []DoneStatus{DoneStatusDone, DoneStatusWontDo})
	if err != nil {
		return err
	}

	data := struct {
		Tasks    []Task
		Sortable bool
	}{
		Tasks:    tasks,
		Sortable: true,
	}
	return c.Render(http.StatusOK, "tasks", data)
}

func (us *uiHandler) delete(c echo.Context) error {
	defer c.Request().Body.Close()

	taskIDStr := c.Param("id")
	taskID64, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		return err
	}
	taskID := uint(taskID64)

	if err := checkPermission(c, taskID); err != nil {
		return err
	}

	ctx := c.Request().Context()

	if err := us.taskService.delete(ctx, taskID); err != nil {
		return err
	}

	return us.search(c)
}

func (us *uiHandler) updateRanks(c echo.Context) error {
	defer c.Request().Body.Close()

	var body struct {
		Ranks map[uint]uint `json:"ranks"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return err
	}

	for id := range body.Ranks {
		if err := checkPermission(c, id); err != nil {
			return err
		}
	}

	ctx := c.Request().Context()
	if err := us.taskService.updateRanks(ctx, body.Ranks); err != nil {
		return err
	}

	return us.search(c)
}

func (us *uiHandler) plan(c echo.Context) error {
	defer c.Request().Body.Close()

	var body struct {
		Duration string `json:"duration"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return err
	}

	d, err := time.ParseDuration(body.Duration)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()

	user, err := loadUser(c)
	if err != nil {
		return err
	}

	planning, err := us.planningService.plan(ctx, user, d)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "plan", planning)
}

func (us *uiHandler) currentPlanning(c echo.Context) error {
	ctx := c.Request().Context()

	user, err := loadUser(c)
	if err != nil {
		return err
	}

	planning, err := us.planningService.current(ctx, user)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "plan", planning)
}

func (us *uiHandler) dismissPlanning(c echo.Context) error {
	ctx := c.Request().Context()

	user, err := loadUser(c)
	if err != nil {
		return err
	}

	if err := us.planningService.dismiss(ctx, user); err != nil {
		return err
	}

	return c.Render(http.StatusOK, "plan", planningForTemplate{})
}
