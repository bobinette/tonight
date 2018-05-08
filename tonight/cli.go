package tonight

import (
	"context"
	"fmt"
)

func RegisterCLI(repo TaskRepository, index TaskIndex) map[string]func() {
	return map[string]func(){
		"reindex.all": func() {
			ctx := context.Background()

			tasks, err := repo.All(ctx)
			if err != nil {
				panic(err)
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

			for _, task := range tasks {
				if err := index.Index(ctx, task); err != nil {
					panic(err)
				}
			}

			fmt.Printf("%d tasks reindexed\n", len(tasks))
		},
	}
}
