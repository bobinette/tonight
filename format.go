package tonight

import (
	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"
)

func formatDescription(desc string) string {
	out := blackfriday.Run([]byte(desc))
	return bluemonday.UGCPolicy().Sanitize(string(out))
}
