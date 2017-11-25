package mysql

import (
	"strings"
)

func join(s, sep string, n int) string {
	a := make([]string, n)
	for i := 0; i < n; i++ {
		a[i] = s
	}
	return strings.Join(a, sep)
}
