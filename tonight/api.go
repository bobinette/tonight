package tonight

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/labstack/echo"
)

var (
	colourRegex = regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)
)

func RegisterAPIHandler(
	e *echo.Echo,
	jwtKey []byte,
	repo TaskRepository,
	index TaskIndex,
	planningRepo PlanningRepository,
	userRepo UserRepository,
) {
	h := apiHandler{
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
		userService: userService{
			jwtKey: jwtKey,
			repo:   userRepo,
		},
	}

	apiGroup := e.Group("/api")
	apiGroup.Use(JWTMiddleware(jwtKey))
	apiGroup.Use(UserMiddleware(jwtKey, userRepo))

	// User
	apiGroup.GET("/me", h.me)
	apiGroup.POST("/tags/:tag", h.customizeColour)

	// Tasks
	apiGroup.GET("/tasks", h.searchTasks)
	apiGroup.POST("/tasks", h.createTask)
	apiGroup.POST("/tasks/:id", h.update)
	apiGroup.DELETE("/tasks/:id", h.delete)
	apiGroup.POST("/tasks/:id/log", h.log)

	// Planning
	apiGroup.GET("/planning", h.currentPlanning)
	apiGroup.POST("/planning", h.createPlanning)
	apiGroup.DELETE("/planning", h.dismissPlanning)

	// Admin
	adminGroup := e.Group("/admin", AdminMiddleware())
	adminGroup.POST("/reindex", h.indexAll)
}

type apiHandler struct {
	repo  TaskRepository
	index TaskIndex

	taskService     taskService
	planningService planningService
	userService     userService
}

func (h *apiHandler) me(c echo.Context) error {
	user, err := loadUser(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h *apiHandler) customizeColour(c echo.Context) error {
	defer c.Request().Body.Close()
	user, err := loadUser(c)

	tag := c.Param("tag")

	var body struct {
		Colour string `json:"colour"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return err
	}

	if !colourRegex.MatchString(body.Colour) {
		return errors.New("invalid colour, expecting an hexadecimal representation such as #A2B5F7")
	}

	user, err = h.userService.customizeColour(c.Request().Context(), user, tag, body.Colour)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h *apiHandler) searchTasks(c echo.Context) error {
	q := c.QueryParam("q")
	sortBy := c.QueryParam("sortBy")
	stringStatuses := func(a []string) []string { // Unique
		res := make([]string, 0, len(a))
		set := make(map[string]struct{})

		for _, s := range a {
			if _, ok := set[s]; ok {
				continue
			}

			res = append(res, s)
			set[s] = struct{}{}
		}

		return res
	}(c.QueryParams()["statuses"])

	ctx := c.Request().Context()

	user, err := loadUser(c)
	if err != nil {
		return err
	}

	statuses := make([]Status, len(stringStatuses))
	if len(stringStatuses) == 0 {
		statuses = []Status{StatusPending}
	} else {
		for i, s := range stringStatuses {
			statuses[i] = StatusFromString(s)
		}
	}

	tasks, err := h.taskService.list(ctx, user, q, statuses, sortBy)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"tasks": tasks})
}

func (h *apiHandler) createTask(c echo.Context) error {
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

	task, err := h.taskService.create(ctx, user, body.Content)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, task)
}

func (h *apiHandler) update(c echo.Context) error {
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
	task, err := h.taskService.update(ctx, taskID, body.Content)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, task)
}

func (h *apiHandler) log(c echo.Context) error {
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
		Log string `json:"log"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return err
	}

	ctx := c.Request().Context()

	task, err := h.taskService.log(ctx, taskID, body.Log)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, task)
}

func (h *apiHandler) delete(c echo.Context) error {
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

	if err := h.taskService.delete(ctx, taskID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *apiHandler) currentPlanning(c echo.Context) error {
	ctx := c.Request().Context()

	user, err := loadUser(c)
	if err != nil {
		return err
	}

	planning, err := h.planningService.current(ctx, user)
	if err != nil {
		return err
	}

	if planning.ID == 0 {
		return c.JSON(http.StatusOK, nil)
	}

	return c.JSON(http.StatusOK, planning)
}

func (h *apiHandler) createPlanning(c echo.Context) error {
	defer c.Request().Body.Close()

	var body struct {
		Input string `json:"input"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return err
	}

	ctx := c.Request().Context()

	user, err := loadUser(c)
	if err != nil {
		return err
	}

	planning, err := h.planningService.plan(ctx, user, body.Input)
	if err != nil {
		return err
	}

	if planning.ID == 0 {
		return c.JSON(http.StatusOK, nil)
	}

	return c.JSON(http.StatusOK, planning)
}

func (h *apiHandler) dismissPlanning(c echo.Context) error {
	ctx := c.Request().Context()

	user, err := loadUser(c)
	if err != nil {
		return err
	}

	if err := h.planningService.dismiss(ctx, user); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *apiHandler) indexAll(c echo.Context) error {
	defer c.Request().Body.Close()

	tasks, err := h.repo.All(c.Request().Context())
	if err != nil {
		return err
	}

	scores := scoreMany(tasks, score)
	for taskID, s := range scores {
		for i, task := range tasks {
			if task.ID != taskID {
				continue
			}

			tasks[i].Score = s
		}
	}

	for _, task := range tasks {
		if err := h.index.Index(c.Request().Context(), task); err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"count": len(tasks)})
}
