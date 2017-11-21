package tonight

import (
	"strings"
)

func parse(content string) Task {
	task := Task{}

	parseFunctions := []func(s string) string{
		// extract the priority
		func(s string) string {
			for strings.HasPrefix(s, "!") && task.Priority <= 5 {
				task.Priority++
				s = s[1:]
			}

			return s
		},
		// extract the title: everything before any of the following symbols:
		// [:, #]
		func(s string) string {
			if idx := strings.IndexAny(s, ":#~"); idx >= 0 {
				task.Title = strings.TrimSpace(s[:idx])
				return s[idx:]
			}

			task.Title = strings.TrimSpace(s)
			return ""
		},
		// save description: starting with ':'
		func(s string) string {
			if !strings.HasPrefix(s, ":") {
				return s
			}

			s = s[1:] // Remove ':'

			if idx := strings.IndexAny(s, "#~"); idx >= 0 {
				task.Description = strings.TrimSpace(s[:idx])
				return s[idx:]
			}

			task.Description = strings.TrimSpace(s)
			return ""
		},
		// save tags: starting with '#'
		func(s string) string {
			for strings.HasPrefix(s, "#") {
				s = s[1:] // Remove '#'

				remaining := ""
				if idx := strings.IndexAny(s, " #~"); idx >= 0 {
					remaining = strings.TrimSpace(s[idx:])
					s = s[:idx]
				}

				task.Tags = append(task.Tags, strings.TrimSpace(s))
				s = remaining
			}
			return s
		},
		// save description: starting with ':'
		func(s string) string {
			if !strings.HasPrefix(s, "~") {
				return s
			}

			s = s[1:] // Remove '~'

			remaining := ""
			if idx := strings.Index(s, " "); idx >= 0 {
				remaining = strings.TrimSpace(s[idx:])
				s = s[:idx]
			}

			task.Duration = strings.TrimSpace(s)
			return remaining
		},
	}

	for i := 0; len(content) > 0 && i < len(parseFunctions); i++ {
		content = parseFunctions[i](content)
	}

	// Clean fields
	task.Title = strings.Trim(task.Title, " ")
	task.Description = strings.Trim(task.Description, " ")

	return task
}
