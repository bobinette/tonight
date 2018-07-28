package tests

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bobinette/tonight/tonight"
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
	assert.False(t, task.UpdatedAt.IsZero())

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
	task.UpdatedAt = taskCopy.UpdatedAt
	assert.Equal(t, task, taskCopy)

	// Add a log
	log := tonight.Log{
		Type:        tonight.LogTypeProgress,
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

	testDependencies(t, repo)
}

func testDependencies(t *testing.T, repo tonight.TaskRepository) {
	ctx := context.Background()

	dependencies := [][]int{
		{},
		{},
		{1},
		{2},
		{2},
		{1, 0},
	}

	tasks := make([]tonight.Task, 0) // Use append to panic if the dependency build order cannot be respected
	for _, deps := range dependencies {
		task := tonight.Task{}
		for _, dep := range deps {
			task.Dependencies = append(task.Dependencies, tonight.Dependency{ID: tasks[dep].ID})
		}

		require.NoError(t, repo.Create(ctx, &task))
		tasks = append(tasks, task)
	}

	tests := map[int][]int{
		// when updating x, the scores of ys... need to be updated
		0: {0, 5},
		1: {1, 2, 3, 4, 5},
		2: {1, 2, 3, 4},
		3: {1, 2, 3},
		4: {1, 2, 4},
		5: {0, 1, 5},
	}

	for taskID, expectedIDs := range tests {
		retrievedTasks, err := repo.DependencyTrees(ctx, tasks[taskID].ID)
		assert.NoError(t, err)

		taskIDs := make([]int, 0)
		for _, task := range retrievedTasks {
			taskIDs = append(taskIDs, int(task.ID))
		}

		expected := make([]int, len(expectedIDs))
		for i, idx := range expectedIDs {
			expected[i] = int(tasks[idx].ID)
		}
		sort.Ints(expected)
		sort.Ints(taskIDs)
		assert.Equal(t, expected, taskIDs, "%d", taskID)
	}

	// Test the cycle protection
	task := tasks[1]
	task.Dependencies = []tonight.Dependency{tonight.Dependency{ID: tasks[2].ID}}
	tasks[1] = task
	fmt.Println(tasks[1].ID, tasks[2].ID)
	require.NoError(t, repo.Update(ctx, &tasks[1]))

	deps, err := repo.DependencyTrees(ctx, task.ID)
	require.NoError(t, err)
	fmt.Printf("%+v\n", deps)
}
