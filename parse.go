package tonight

import (
	"fmt"
	"strings"
	"time"
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
			if idx := strings.IndexAny(s, ":#~>"); idx >= 0 {
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
		// save duration: starting with '~'
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
		// save deadline: starting with '>'
		func(s string) string {
			if !strings.HasPrefix(s, ">") {
				return s
			}

			s = s[1:] // Remove '>'

			remaining := ""
			if idx := strings.Index(s, " "); idx >= 0 {
				remaining = strings.TrimSpace(s[idx:])
				s = s[:idx]
			}

			deadline, err := time.Parse("2006-01-02", s)
			if err != nil {
				// Fail early
				return remaining
			}

			deadline = deadline.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			task.Deadline = &deadline
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

func formatRaw(t Task) string {
	out := fmt.Sprintf("%s%s", strings.Repeat("!", t.Priority), t.Title)

	if t.Description != "" {
		out = fmt.Sprintf("%s: %s", out, t.Description)
	}

	for _, tag := range t.Tags {
		out = fmt.Sprintf("%s #%s", out, tag)
	}

	if t.Duration != "" {
		out = fmt.Sprintf("%s ~%s", out, t.Duration)
	}

	if t.Deadline != nil {
		out = fmt.Sprintf("%s ~%s", out, t.Deadline.Format("2006-01-02"))
	}

	return out
}
