package tonight

import (
	"fmt"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"
)

func formatDescription(desc string) string {
	out := blackfriday.Run([]byte(desc))
	return bluemonday.UGCPolicy().Sanitize(string(out))
}

func formatDuration(dur time.Duration) string {
	str := ""
	h := int(dur / time.Hour)
	if h > 0 {
		str = fmt.Sprintf("%dh", h)
	}
	dur = dur - time.Duration(h)*time.Hour

	m := int(dur / time.Minute)
	if m > 0 {
		str = fmt.Sprintf("%s%dm", str, m)
	}
	dur = dur - time.Duration(m)*time.Minute

	s := int(dur / time.Second)
	if s > 0 {
		str = fmt.Sprintf("%s%ds", str, s)
	}

	return str
}
