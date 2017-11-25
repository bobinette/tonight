package tonight

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	completionRE = regexp.MustCompile(`(?:([0-9]*)?% )?(.*)`)
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
		// extract dependencies (should be last)
		func(s string) string {
			idx := strings.Index(s, "needs:")
			if idx < 0 {
				return s
			}

			var lastIdx int // len("needs:")
			for lastIdx = idx + 6; lastIdx < len(s) && s[lastIdx] != ' '; lastIdx++ {
			}

			ids := strings.Split(s[idx+6:lastIdx], ",")
			task.Dependencies = make([]Dependency, 0, len(ids))
			for _, idStr := range ids {
				if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
					task.Dependencies = append(task.Dependencies, Dependency{ID: uint(id)})
				}
			}

			return s[:idx]
		},
		// extract the title:
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

func parseLog(desc string) Log {
	matches := completionRE.FindStringSubmatch(desc)

	log := Log{
		Completion:  100,
		Description: matches[2],
	}

	if len(matches[1]) > 0 {
		log.Completion, _ = strconv.Atoi(matches[1])
		if log.Completion > 100 {
			log.Completion = 100
		}
	}

	return log
}
