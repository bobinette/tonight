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
	percentageCompletionRegex = regexp.MustCompile(`(?:([0-9]*%|pause|stop|start|resume))?(.*)`)
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
	matches := percentageCompletionRegex.FindStringSubmatch(desc)

	log := Log{
		Type:        LogTypeCompletion,
		Completion:  0,
		Description: strings.TrimSpace(matches[2]),
	}

	if len(matches[1]) > 0 {
		v := matches[1]
		if v == "pause" || v == "stop" {
			log.Type = LogTypePause
		} else if v == "start" || v == "resume" {
			log.Type = LogTypeStart
		} else {
			v = v[:len(v)-1]
			log.Completion, _ = strconv.Atoi(v)
			if log.Completion > 100 {
				log.Completion = 100
			} else if log.Completion == 0 {
				log.Type = LogTypeStart
			}
		}
	} else {
		log.Completion = 100
	}

	return log
}
