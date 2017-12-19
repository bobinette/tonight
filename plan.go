package tonight

import (
	"sort"
	"time"
)

func plan(tasks []Task, d time.Duration, strict bool) []Task {
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
