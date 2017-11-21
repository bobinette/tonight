package tonight

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskSorter(t *testing.T) {
	tests := map[string]struct {
		tasks    []taskDuration
		expected []taskDuration
	}{
		"by prio only": {
			tasks:    []taskDuration{{Priority: 4}, {Priority: 5}, {Priority: 0}, {Priority: 2}},
			expected: []taskDuration{{Priority: 5}, {Priority: 4}, {Priority: 2}, {Priority: 0}},
		},
		"prio and duration": {
			tasks: []taskDuration{
				{Priority: 4, Duration: 1 * time.Hour},
				{Priority: 4, Duration: 2 * time.Hour},
				{Priority: 2, Duration: 2 * time.Hour},
				{Priority: 2, Duration: 1 * time.Hour},
			},
			expected: []taskDuration{
				{Priority: 4, Duration: 2 * time.Hour},
				{Priority: 4, Duration: 1 * time.Hour},
				{Priority: 2, Duration: 2 * time.Hour},
				{Priority: 2, Duration: 1 * time.Hour},
			},
		},
		"Keep order": {
			tasks: []taskDuration{
				{ID: 1, Priority: 2, Duration: 1 * time.Hour},
				{ID: 2, Priority: 2, Duration: 1 * time.Hour},
				{ID: 3, Priority: 2, Duration: 1 * time.Hour},
				{ID: 4, Priority: 2, Duration: 1 * time.Hour},
			},
			expected: []taskDuration{
				{ID: 1, Priority: 2, Duration: 1 * time.Hour},
				{ID: 2, Priority: 2, Duration: 1 * time.Hour},
				{ID: 3, Priority: 2, Duration: 1 * time.Hour},
				{ID: 4, Priority: 2, Duration: 1 * time.Hour},
			},
		},
	}

	for name, test := range tests {
		sort.Stable(taskSorter(test.tasks))
		assert.Equal(t, test.expected, test.tasks, name)
	}
}
