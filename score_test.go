package tonight

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScore(t *testing.T) {
	in3Days := time.Now().Add(3 * 24 * time.Hour)
	tasks := []Task{
		{Deadline: &in3Days, Duration: "2h"},
		{Priority: 5, Duration: "2h"},
		{Deadline: &in3Days},
		{Priority: 5},
		{Priority: 3, Duration: "2h"},
		{Priority: 3},
		{Duration: "2h"},
		{Duration: "15m"},
	}

	scores := make([]float64, len(tasks))
	for i, task := range tasks {
		scores[i] = score(task)
	}

	for i := 0; i < len(scores)-2; i++ {
		assert.True(t, scores[i] > scores[i+1], "%d: %f <= %f\n%+v\n%+v", i, scores[i], scores[i+1], tasks[i], tasks[i+1])
	}
}
