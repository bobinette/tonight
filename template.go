package tonight

import (
	"fmt"
	"html/template"
	"io"
	"path"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/labstack/echo"
)

type templateRenderer struct {
	templates map[string]*template.Template
}

func RegisterTemplateRenderer(e *echo.Echo, dir string) error {
	tmpls := map[string][]string{
		"home": {
			"home.tmpl", "tasks.tmpl", "plan.tmpl", "timeline.tmpl",
			// JS files
			"tasksList.js", "sort.js", "delete.js", "new.js",
			"doneTasksList.js", "plan.js", "utils.js", "home.js",
		},
		"tasks": {"tasks.tmpl", "timeline.tmpl"},
		"plan":  {"plan.tmpl", "timeline.tmpl"},
		"login": {"login.tmpl"},
	}

	funcMap := map[string]interface{}{
		"formatDuration":     formatDuration,
		"formatDateRelative": humanize.Time,
		"formatPriority": func(p int) string {
			return strings.Repeat("!", p)
		},
		"formatDate": func(d time.Time) string {
			return d.Format("2006-01-02")
		},
		"formatDateTime": func(dt time.Time) string {
			return dt.Format("2006-01-02 15:03:04")
		},
		"raw":            formatRaw,
		"formatMarkDown": formatDescription,
		"formatDependencies": func(dependencies []Dependency) string {
			total := len(dependencies)
			done := 0
			for _, d := range dependencies {
				if d.Done {
					done++
				}
			}

			return fmt.Sprintf("%d/%d", done, total)
		},
		"reverseLog": func(a []Log) []Log {
			reversed := make([]Log, len(a))
			l := len(reversed) - 1
			for i, e := range a {
				reversed[l-i] = e
			}
			return reversed
		},
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
