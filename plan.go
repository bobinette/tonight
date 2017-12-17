package tonight

import (
	"fmt"
	"sort"
	"time"
)

type taskDuration struct {
	ID       uint
	Priority int
	Duration time.Duration
}

func plan(tasks []Task, d time.Duration) ([]Task, time.Duration) {
	if len(tasks) == 0 {
		return nil, 0
	}

	durations := make([]taskDuration, 0, len(tasks))
	taskMapping := make(map[uint]Task)
	for _, task := range tasks {
		ready := true
		for _, dep := range task.Dependencies {
			if !dep.Done {
				ready = false
				break
			}
		}
		if !ready {
			continue
		}

		durations = append(durations, taskDuration{
			ID:       task.ID,
			Priority: task.Priority,
			Duration: task.LeftDuration(),
		})

		taskMapping[task.ID] = task
	}

	sort.Stable(taskSorter(durations))

	var cumDur time.Duration = 0
	quickestTask := durations[0] // will be used if the planning is empty
	planned := make([]Task, 0)
	for _, task := range durations {
		if task.Duration < quickestTask.Duration {
			quickestTask = task
		}

		if cumDur+task.Duration > d {
			continue
		}

		planned = append(planned, taskMapping[task.ID])
		cumDur += task.Duration
	}

	if len(planned) > 0 {
		return planned, cumDur
	}

	// No task could fit, we select the quickest one
	return []Task{taskMapping[quickestTask.ID]}, quickestTask.Duration
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
