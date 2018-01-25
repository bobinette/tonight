package tonight

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
)

type Indexer struct {
	Repository interface {
		All(ctx context.Context) ([]Task, error)
	}
	Index TaskIndex
}

func (i *Indexer) IndexAll(c echo.Context) error {
	defer c.Request().Body.Close()

	tasks, err := i.Repository.All(c.Request().Context())
	if err != nil {
		return err
	}

	scores := scoreMany(tasks, score)
	for taskID, s := range scores {
		for i, task := range tasks {
			if task.ID != taskID {
				continue
			}

			tasks[i].Score = s
		}
	}

	count := 0
	for _, task := range tasks {
		if err := i.Index.Index(c.Request().Context(), task); err != nil {
			return err
		}
		count++
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"count": count})
}
