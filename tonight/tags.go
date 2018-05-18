package tonight

import (
	"context"
)

type TagReader interface {
	Tags(ctx context.Context, user, q string) ([]string, error)
}
