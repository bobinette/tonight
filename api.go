package tonight

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
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
		userRepository: userRepo,
	}

	apiGroup := e.Group("/api")
	apiGroup.Use(JWTMiddleware(jwtKey))
	apiGroup.Use(UserMiddleware(jwtKey, userRepo))

	apiGroup.GET("/me", h.me)
	apiGroup.GET("/tasks", h.searchTasks)
	apiGroup.POST("/tasks", h.createTask)
	apiGroup.POST("/tasks/:id/log", h.log)
}

type apiHandler struct {
	repo  TaskRepository
	index TaskIndex

	userRepository UserRepository

	taskService     taskService
	planningService planningService
}

func (h *apiHandler) me(c echo.Context) error {
	user, err := loadUser(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h *apiHandler) searchTasks(c echo.Context) error {
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
