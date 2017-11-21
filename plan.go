package tonight

import (
	"sort"
	"time"
)

type taskDuration struct {
	ID       uint
	Priority int
	Duration time.Duration
}

func plan(tasks []Task, d time.Duration) ([]Task, time.Duration) {
	durations := make([]taskDuration, len(tasks))
	taskMapping := make(map[uint]Task)
	for i, task := range tasks {
		td, err := time.ParseDuration(task.Duration)
		if err != nil {
			td = 1 * time.Hour
		}

		durations[i] = taskDuration{
			ID:       task.ID,
			Priority: task.Priority,
			Duration: td,
		}

		taskMapping[task.ID] = task
	}

	sort.Stable(taskSorter(durations))
	var cumDur time.Duration = 0
	planned := make([]Task, 0)
	for _, task := range durations {
		if cumDur+task.Duration > d {
			continue
		}

		planned = append(planned, taskMapping[task.ID])
		cumDur += task.Duration
	}

	return planned, cumDur
}

type taskSorter []taskDuration

func (t taskSorter) Len() int      { return len(t) }
func (t taskSorter) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t taskSorter) Less(i, j int) bool {
	if t[i].Priority != t[j].Priority {
		return t[i].Priority > t[j].Priority
	}

	return t[i].Duration > t[j].Duration
}
