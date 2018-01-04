package tonight

import (
	"math"
	"time"
)

func scoreMany(tasks []Task) map[uint]float64 {
	// Construct the dependency tree
	// Well, not exactly...
	reversedDependencies := make(map[uint][]uint)
	for _, task := range tasks {
		for _, dep := range task.Dependencies {
			deps := reversedDependencies[dep.ID]
			deps = append(deps, task.ID)
			reversedDependencies[dep.ID] = deps
		}
	}

	// Compute all the scores
	scores := make(map[uint]float64)
	for _, task := range tasks {
		scores[task.ID] = score(task)
	}

	// Boost the scores of task listed as dependencies
	for depID, taskIDs := range reversedDependencies {
		for _, taskID := range taskIDs {
			scores[depID] += scores[taskID]
		}
	}

	return scores
}

// score gives a score to task to be used when ranking for planning.
func score(task Task) float64 {
	// Start by using the priority of the task
	s := float64(task.Priority)

	// Then we take the duration into consideration: we want the longer tasks to appear first, but
	// we don't want the duration to completely take over the priority
	if d, err := time.ParseDuration(task.Duration); err == nil && d > 0 {
		// In case there is a task that is supposed to take more than
		// e^5 = 148h...
		s += math.Min(math.Log(1+float64(d)/float64(time.Hour)), 5)
	}

	// Finally, we boost the tasks with a deadline
	if task.Deadline != nil {
		delta := time.Now().Sub(*task.Deadline)
		s += 6 * (1 - 1/(1+math.Exp(3-float64(delta))))
	}

	return s
}
