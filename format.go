package tonight

import (
	"html/template"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"
)

func formatDescription(desc string) template.HTML {
	out := blackfriday.Run([]byte(desc))
	return template.HTML(bluemonday.UGCPolicy().Sanitize(string(out)))
}

func formatDuration(dur time.Duration) string {
	return dur.String()
}
