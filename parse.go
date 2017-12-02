package tonight

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// Task
	priorityRegex         = regexp.MustCompile(`^(!*)`)
	dependenciesRegex     = regexp.MustCompile(`needs:((?:\d+,?)+)`)
	tagsRegex             = regexp.MustCompile(`\B\#(\w+\b)`)
	durationRegex         = regexp.MustCompile(`\B~([0-9hms]+)`)
	deadlineRegex         = regexp.MustCompile(`\B>(\d{4}-\d{1,2}-\d{1,2})`)
	titleDescriptionRegex = regexp.MustCompile(`([^:]*)(?::(.*))?`)

	// Log
	logKeywordRegex    = regexp.MustCompile(`^(pause|stop|start|resume|done)`)
	logCompletionRegex = regexp.MustCompile(`^(\d+)%`)
	logFractionRegex   = regexp.MustCompile(`^(\d+)/([1-9]\d*)`)

	logKeywordMapping = map[string]LogType{
		"pause":  LogTypePause,
		"stop":   LogTypePause,
		"start":  LogTypeStart,
		"resume": LogTypeStart,
		"done":   LogTypeCompletion,
	}
)

func parse(content string) Task {
	task := Task{}

	parseFunctions := []func(s string) string{
		// extract the priority
		func(s string) string {
			matches := priorityRegex.FindStringSubmatch(s)
			task.Priority = len(matches[1])
			if task.Priority > 5 {
				task.Priority = 5
			}

			return priorityRegex.ReplaceAllString(s, "")
		},
		// extract dependencies
		func(s string) string {
			matches := dependenciesRegex.FindStringSubmatch(s)
			if len(matches) == 0 {
				return s
			}

			ids := strings.Split(matches[1], ",")
			task.Dependencies = make([]Dependency, 0, len(ids))
			for _, idStr := range ids {
				if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
					task.Dependencies = append(task.Dependencies, Dependency{ID: uint(id)})
				}
			}

			return dependenciesRegex.ReplaceAllString(s, "")
		},
		// extract tags
		func(s string) string {
			matches := tagsRegex.FindAllStringSubmatch(s, -1)
			task.Tags = make([]string, len(matches))
			for i, match := range matches {
				task.Tags[i] = match[1]
			}
			return tagsRegex.ReplaceAllString(s, "")
		},
		// save duration: starting with '~'
		func(s string) string {
			matches := durationRegex.FindStringSubmatch(s)
			if len(matches) == 0 {
				return s
			}

			task.Duration = matches[1]
			return durationRegex.ReplaceAllString(s, "")
		},
		// save deadline: starting with '>'
		func(s string) string {
			matches := deadlineRegex.FindStringSubmatch(s)
			if len(matches) == 0 {
				return s
			}

			if deadline, err := time.Parse("2006-01-02", matches[1]); err == nil {
				deadline = deadline.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
				task.Deadline = &deadline
			}

			return deadlineRegex.ReplaceAllString(s, "")
		},
		// extract the title:
		func(s string) string {
			s = strings.TrimSpace(s)

			matches := titleDescriptionRegex.FindStringSubmatch(s)
			if len(matches) == 0 {
				return s
			}

			task.Title = matches[1]
			task.Description = matches[2]
			return deadlineRegex.ReplaceAllString(s, "")
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
	log := Log{
		Type:       LogTypeCompletion,
		Completion: 0,
	}

	if keywordMatch := logKeywordRegex.FindStringSubmatch(desc); len(keywordMatch) > 0 {
		log.Type = logKeywordMapping[keywordMatch[1]]
		desc = logKeywordRegex.ReplaceAllString(desc, "")

		if keywordMatch[1] == "done" {
			log.Completion = 100
		}
	} else if completionMatch := logCompletionRegex.FindStringSubmatch(desc); len(completionMatch) > 0 {
		log.Completion, _ = strconv.Atoi(completionMatch[1])
		desc = logCompletionRegex.ReplaceAllString(desc, "")
	} else if fractionMatch := logFractionRegex.FindStringSubmatch(desc); len(fractionMatch) > 0 {
		num, _ := strconv.Atoi(fractionMatch[1])
		den, _ := strconv.Atoi(fractionMatch[2])
		log.Completion = (num * 100) / den
		desc = logFractionRegex.ReplaceAllString(desc, "")
	}

	if log.Completion > 100 {
		log.Completion = 100
	}

	log.Description = strings.TrimSpace(desc)

	return log
}
