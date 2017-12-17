package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bobinette/tonight"
)

func TestTaskRepository(t *testing.T, repo tonight.TaskRepository) {
	ctx := context.Background()
	task := tonight.Task{
		Title:        "test task",
		Priority:     3,
		Tags:         []string{"tag", "test"},
		Duration:     "2m",
		Dependencies: []tonight.Dependency{},
	}

	// Create the task
	err := repo.Create(ctx, &task)
	require.NoError(t, err)

	assert.NotEqual(t, uint(0), task.ID)
	assert.Equal(t, uint(1), task.Rank)
	assert.False(t, task.CreatedAt.IsZero())

	// Retrieve it
	tasks, err := repo.List(ctx, []uint{task.ID})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, task, tasks[0])

	// Update the description
	task.Description = "description"
	task.Tags = append(task.Tags, "zzz")
	taskCopy := task
	taskCopy.Rank = 0 // reset rank to verify that it is still reloaded
	err = repo.Update(ctx, &taskCopy)
	assert.NoError(t, err)
	assert.Equal(t, task, taskCopy)

	// Add a log
	log := tonight.Log{
		Type:        tonight.LogTypeCompletion,
		Completion:  17,
		Description: "log",
	}
	err = repo.Log(ctx, task.ID, log)
	assert.NoError(t, err)

	// Update the rank
	ranks := map[uint]uint{task.ID: 2}
	err = repo.UpdateRanks(ctx, ranks)
	assert.NoError(t, err)

	// Verify that everything has indeed been updated
	task.Log = []tonight.Log{log}
	task.Rank = 2
	tasks, err = repo.List(ctx, []uint{task.ID})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tasks))
	// Fake the log time...
	task.Log[0].CreatedAt = tasks[0].Log[0].CreatedAt
	assert.Equal(t, task, tasks[0])

	// Delete the task
	err = repo.Delete(ctx, task.ID)
	assert.NoError(t, err)

	// Make sure I cannot get it back
	tasks, err = repo.List(ctx, []uint{task.ID})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(tasks))
}
