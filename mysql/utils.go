package mysql

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func prepareArgs(params ...interface{}) ([]interface{}, []interface{}) {
	qArgs := make([]interface{}, 0)
	args := make([]interface{}, 0)
	for _, p := range params {
		switch p := p.(type) {
		case []string:
			s := make([]string, len(p))
			for i, e := range p {
				s[i] = "?"
				args = append(args, e)
			}
			qArgs = append(qArgs, fmt.Sprintf("(%s)", strings.Join(s, ",")))
		case []uuid.UUID:
			s := make([]string, len(p))
			for i, e := range p {
				s[i] = "?"
				args = append(args, e.String())
			}
			qArgs = append(qArgs, fmt.Sprintf("(%s)", strings.Join(s, ",")))
		default:
			qArgs = append(qArgs, "?")
			args = append(args, p)
		}
	}
	return qArgs, args
}
