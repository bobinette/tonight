package tonight

import (
	"strings"
)

func parse(content string) Task {
	task := Task{
		Title: content,
	}

	if idx := strings.Index(content, ":"); idx >= 0 {
		task.Title = content[:idx]
		task.Description = content[idx+1:]
	}

	// Clean fields
	task.Title = strings.Trim(task.Title, " ")
	task.Description = strings.Trim(task.Description, " ")

	return task
}
