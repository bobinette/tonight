package tonight

import (
	"sort"
	"time"
)

func plan(tasks []Task, d time.Duration, strict bool) []Task {
	tasks = filterUndoneDependencies(tasks)
	if len(tasks) == 0 {
		return nil
	}

	sort.Stable(byScore(tasks))

	var cumDur time.Duration
	planned := make([]Task, 0, len(tasks))
	for i := 0; cumDur < d && i < len(tasks); i++ {
		task := tasks[i]
		if strict && cumDur+task.LeftDuration() > d {
			// This task does not fit
			continue
		}

		planned = append(planned, task)
		cumDur += task.LeftDuration()
	}

	return planned
}

func planNext(tasks []Task, planning Planning, afterID uint) []Task {
	tasks = filterUndoneDependencies(tasks)
	if len(tasks) == 0 {
		return nil
	}

	sort.Stable(byScore(tasks))

	var cumDur time.Duration
	planned := make([]Task, 0, len(tasks))
	isAfter := false
	for i := 0; cumDur < planning.Duration && i < len(tasks); i++ {
		task := tasks[i]
		if task.ID == afterID {
			isAfter = true
			continue
		}

		if isPlanned(task, planning) {
			cumDur += task.LeftDuration()
			continue
		}

		if !isAfter {
			continue
		}

		if planning.Strict && cumDur+task.LeftDuration() > planning.Duration {
			// This task does not fit
			continue
		}

		planned = append(planned, task)
		cumDur += task.LeftDuration()
	}

	return planned
}

func isPlanned(task Task, planning Planning) bool {
	for _, t := range planning.Tasks {
		if t.ID == task.ID {
			return true
		}
	}
	return false
}

func filterUndoneDependencies(tasks []Task) []Task {
	filtered := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		hasUndone := false

		for _, dep := range task.Dependencies {
			if !dep.Done {
				hasUndone = true
				break
			}
		}

		if !hasUndone {
			filtered = append(filtered, task)
		}
	}

	return filtered
}

type byScoreSorter struct {
	tasks  []Task
	scores []float64
}

func byScore(tasks []Task) *byScoreSorter {
	scores := make([]float64, len(tasks))
	for i, task := range tasks {
		scores[i] = score(task)
	}

	return &byScoreSorter{
		tasks:  tasks,
		scores: scores,
	}
}

func (s *byScoreSorter) Len() int { return len(s.tasks) }
func (s *byScoreSorter) Swap(i, j int) {
	s.tasks[i], s.tasks[j] = s.tasks[j], s.tasks[i]
	s.scores[i], s.scores[j] = s.scores[j], s.scores[i]
}
func (s *byScoreSorter) Less(i, j int) bool { return s.scores[i] > s.scores[j] }
