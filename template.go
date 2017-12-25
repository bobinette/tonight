package tonight

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/labstack/echo"
)

type templateRenderer struct {
	names     map[string]string
	templates map[string]*template.Template
}

type templateTree struct {
	name     string
	filename string
	children []templateTree
}

func templateTreeFromName(name string) templateTree {
	filename := name
	if !strings.HasSuffix(filename, ".tmpl") && !strings.HasSuffix(filename, ".js") {
		filename = fmt.Sprintf("%s.tmpl", filename)
	}
	return templateTree{
		name:     name,
		filename: filename,
	}
}

func (t *templateTree) load(dir string, funcMap template.FuncMap, tmpl *template.Template) error {
	data, err := ioutil.ReadFile(path.Join(dir, t.filename))
	if err != nil {
		return err
	}

	if _, err := tmpl.Parse(string(data)); err != nil {
		return err
	}

	for _, child := range t.children {
		childTemplate := tmpl.New(child.name).Funcs(funcMap)
		if err := child.load(dir, funcMap, childTemplate); err != nil {
			return err
		}
	}

	return nil
}

func RegisterTemplateRenderer(e *echo.Echo, dir string) error {
	trees := map[string]templateTree{
		"home": {
			name:     "home",
			filename: "home.tmpl",
			children: []templateTree{
				{
					name:     "planning",
					filename: "planning.tmpl",
					children: []templateTree{
						templateTreeFromName("task-row-buttons-planning"),
						templateTreeFromName("timeline"),
						templateTreeFromName("task-row"),
						templateTreeFromName("task-list"),
					},
				},
				{
					name:     "pending-tasks",
					filename: "task-list.tmpl",
					children: []templateTree{
						templateTreeFromName("task-row-buttons-pending"),
						templateTreeFromName("timeline"),
						templateTreeFromName("task-row"),
					},
				},
				templateTreeFromName("timeline.tmpl"),
			},
		},
		"planning": { // @TODO: factorize the trees?
			name:     "planning",
			filename: "planning.tmpl",
			children: []templateTree{
				templateTreeFromName("task-row-buttons-planning"),
				templateTreeFromName("timeline"),
				templateTreeFromName("task-row"),
				templateTreeFromName("task-list"),
			},
		},
		"pending-tasks": {
			name:     "pending-tasks",
			filename: "task-list.tmpl",
			children: []templateTree{
				templateTreeFromName("task-row-buttons-pending"),
				templateTreeFromName("timeline"),
				templateTreeFromName("task-row"),
			},
		},
		"complete-tasks": {
			name:     "complete-tasks",
			filename: "task-list.tmpl",
			children: []templateTree{
				templateTreeFromName("task-row-buttons-complete"),
				templateTreeFromName("timeline"),
				templateTreeFromName("task-row"),
			},
		},
		"login": {
			name:     "login",
			filename: "login.tmpl",
		},
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

	names := make(map[string]string)
	templates := make(map[string]*template.Template)
	for name, tree := range trees {
		tmpl := template.New(name).Funcs(funcMap)
		if err := tree.load(dir, funcMap, tmpl); err != nil {
			return err
		}

		names[name] = tree.name
		templates[name] = tmpl
	}

	e.Renderer = &templateRenderer{
		names:     names,
		templates: templates,
	}

	return nil
}

func (t *templateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("unknown template: %s", name)
	}

	name = t.names[name]
	return tmpl.ExecuteTemplate(w, name, data)
}
