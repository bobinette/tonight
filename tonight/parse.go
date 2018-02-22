package tonight

import (
	"errors"
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
	tagsRegex             = regexp.MustCompile(`\B\#((?:\w|\-|\:)+)\b`)
	durationRegex         = regexp.MustCompile(`\B~([0-9a-zA-Z]+)`) // a-zA-Z to have a meaningful error message
	deadlineRegex         = regexp.MustCompile(`\B>(\d{4}-\d{1,2}-\d{1,2})`)
	titleDescriptionRegex = regexp.MustCompile(`([^:]*)(?::(.*))?`)

	// Log
	logCompletionRegex = regexp.MustCompile(`^(\d+)%`)
	logFractionRegex   = regexp.MustCompile(`^(\d+)/([1-9]\d*)`)

	logKeywordMapping = map[string]LogType{
		"pause":    LogTypePause,
		"stop":     LogTypePause,
		"start":    LogTypeStart,
		"resume":   LogTypeStart,
		"done":     LogTypeProgress,
		"won't do": LogTypeWontDo,
	}
	logKeywordRegex = func(m map[string]LogType) *regexp.Regexp {
		keyWords := make([]string, len(m))
		i := 0
		for k := range m {
			keyWords[i] = k
			i++
		}

		return regexp.MustCompile(fmt.Sprintf(`^(%s)`, strings.Join(keyWords, "|")))
	}(logKeywordMapping)

	planningRegex = regexp.MustCompile(`^(.* for )?(!)?([0-9a-zA-Z]+)$`)
)

func parse(content string) (Task, error) {
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
			tags := make(map[string]struct{})
			for _, match := range matches {
				tags[match[1]] = struct{}{}
			}

			task.Tags = make([]string, len(tags))
			i := 0
			for tag := range tags {
				task.Tags[i] = tag
				i++
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

	if task.Duration != "" {
		dur, err := time.ParseDuration(task.Duration)
		if err != nil {
			return Task{}, err
		}

		// Reformat it for consistency
		task.Duration = formatDuration(dur)
	}

	return task, nil
}

func parseLog(desc string) Log {
	log := Log{
		Type:       LogTypeComment,
		Completion: 0,
	}

	if keywordMatch := logKeywordRegex.FindStringSubmatch(desc); len(keywordMatch) > 0 {
		log.Type = logKeywordMapping[keywordMatch[1]]
		desc = logKeywordRegex.ReplaceAllString(desc, "")

		if strings.HasPrefix(desc, ":") {
			desc = desc[1:]
		}

		if keywordMatch[1] == "done" {
			log.Completion = 100
			log.Type = LogTypeProgress
		}
	} else if completionMatch := logCompletionRegex.FindStringSubmatch(desc); len(completionMatch) > 0 {
		log.Completion, _ = strconv.Atoi(completionMatch[1])
		log.Type = LogTypeProgress
		desc = logCompletionRegex.ReplaceAllString(desc, "")
	} else if fractionMatch := logFractionRegex.FindStringSubmatch(desc); len(fractionMatch) > 0 {
		num, _ := strconv.Atoi(fractionMatch[1])
		den, _ := strconv.Atoi(fractionMatch[2])
		log.Completion = (num * 100) / den
		log.Type = LogTypeProgress
		desc = logFractionRegex.ReplaceAllString(desc, "")
	}

	if log.Completion > 100 {
		log.Completion = 100
	}

	log.Description = strings.TrimSpace(desc)

	return log
}

func parsePlanning(input string) (string, time.Duration, bool, error) {
	matches := planningRegex.FindStringSubmatch(input)
	if len(matches) != 4 {
		return "", 0, false, errors.New("incorrect format")
	}

	q := matches[1]
	if q != "" {
		q = q[0 : len(q)-5] // len(" for ")
	}
	strict := matches[2] == "!"
	duration, err := time.ParseDuration(matches[3])
	if err != nil {
		return "", 0, false, err
	}

	return q, duration, strict, nil
}
