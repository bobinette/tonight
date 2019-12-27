package mysql

import "fmt"

import "strings"

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
		default:
			qArgs = append(qArgs, "?")
			args = append(args, p)
		}
	}
	return qArgs, args
}
