package tonight

import (
	"math"
	"time"
)

func scoreMany(tasks []Task, scoreFunc func(task Task) float64) map[uint]float64 {
	// Compute unaries
	scores := make(map[uint]float64)
	for _, task := range tasks {
		scores[task.ID] = scoreFunc(task)
	}

	trees := buildDependencyTrees(tasks)
	for _, tree := range trees {
		tree.traverseBottomUp(func(t *dependencyTree) {
			for _, child := range t.children {
				if t.node.PostponedUntil != nil && t.node.PostponedUntil.After(time.Now()) {
					continue
				}
				scores[t.node.ID] += scores[child.node.ID] + 1
			}
		})
	}

	return scores
}

// score gives a score to task to be used when ranking for planning.
func score(task Task) float64 {
	// Postponed tasks have a score of 0. This way, if a postponed task is
	// a dependency of a very urgent task, it will still be scored highly, but
	// postponed tasks that are not depended on will be at the bottom of the ocean
	if task.PostponedUntil != nil && task.PostponedUntil.After(time.Now()) {
		return 0
	}

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
