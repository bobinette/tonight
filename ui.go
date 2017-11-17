package tonight

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/labstack/echo"
)

type templateRenderer struct {
	templates map[string]*template.Template
}

func RegisterTemplateRenderer(e *echo.Echo, dir string) error {
	tmpls := map[string][]string{
		"home":  {"home.tmpl", "tasks.tmpl"},
		"tasks": {"tasks.tmpl"},
	}

	funcMap := map[string]interface{}{
		"formatDate": humanize.Time,
	}

	renderer := &templateRenderer{
		templates: make(map[string]*template.Template),
	}
	for name, filenames := range tmpls {
		files := make([]string, len(filenames))
		for i, filename := range filenames {
			files[i] = path.Join(dir, filename)
		}

		t, err := template.New(name).Funcs(funcMap).ParseFiles(files...)
		if err != nil {
			return err
		}

		renderer.templates[name] = t
	}

	e.Renderer = renderer

	return nil
}

func (t *templateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("unknown template: %s", name)
	}

	return tmpl.ExecuteTemplate(w, fmt.Sprintf("%s.tmpl", name), data)
}

type UIService struct {
	repo TaskRepository
}

func NewUIService(repo TaskRepository) *UIService {
	return &UIService{
		repo: repo,
	}
}

func (us *UIService) Home(c echo.Context) error {
	tasks, err := us.repo.List(c.Request().Context(), false)
	if err != nil {
		return err
	}

	data := struct {
		Tasks []Task
	}{
		Tasks: tasks,
	}
	return c.Render(http.StatusOK, "home", data)
}

func (us *UIService) CreateTask(c echo.Context) error {
	defer c.Request().Body.Close()

	var t struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&t); err != nil {
		return err
	}

	task := parse(t.Content)

	if err := us.repo.Create(c.Request().Context(), &task); err != nil {
		return err
	}

	return us.tasks(c)
}

func (us *UIService) MarkDone(c echo.Context) error {
	defer c.Request().Body.Close()

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		return err
	}

	if err := us.repo.MarkDone(c.Request().Context(), uint(taskID)); err != nil {
		return err
	}

	return us.tasks(c)
}

func (us *UIService) tasks(c echo.Context) error {
	tasks, err := us.repo.List(c.Request().Context(), false)
	if err != nil {
		return err
	}

	data := struct {
		Tasks []Task
	}{
		Tasks: tasks,
	}
	return c.Render(http.StatusOK, "tasks", data)
}

func (us *UIService) DoneTasks(c echo.Context) error {
	tasks, err := us.repo.List(c.Request().Context(), true)
	if err != nil {
		return err
	}

	data := struct {
		Tasks []Task
	}{
		Tasks: tasks,
	}
	return c.Render(http.StatusOK, "tasks", data)
}

func (us *UIService) Delete(c echo.Context) error {
	defer c.Request().Body.Close()

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		return err
	}

	if err := us.repo.Delete(c.Request().Context(), uint(taskID)); err != nil {
		return err
	}

	return us.tasks(c)
}

func (us *UIService) UpdateRanks(c echo.Context) error {
	defer c.Request().Body.Close()

	var body struct {
		Ranks map[uint]uint `json:"ranks"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return err
	}

	if err := us.repo.UpdateRanks(c.Request().Context(), body.Ranks); err != nil {
		return err
	}

	return us.tasks(c)
}
