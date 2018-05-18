package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bobinette/tonight/tonight"
)

func TestTagReader(
	t *testing.T,
	tagReader tonight.TagReader,
	taskRepo tonight.TaskRepository,
	userRepo tonight.UserRepository,
) {
	ctx := context.Background()

	tasks := map[string][]tonight.Task{
		"user 1": []tonight.Task{
			{Tags: []string{"tag", "hello", "world"}},
			{Tags: []string{"tag", "pizza", "yolo"}},
			{Tags: []string{"yolo", "yoloy"}},
		},
		"user 2": []tonight.Task{
			{Tags: []string{"hello", "world", "war", "z"}},
		},
	}

	for userName, userTasks := range tasks {
		user := tonight.User{Name: userName}
		err := userRepo.Insert(ctx, &user)
		require.NoError(t, err)

		for _, task := range userTasks {
			err := taskRepo.Create(ctx, &task)
			require.NoError(t, err)

			err = userRepo.AddTaskToUser(ctx, user.ID, task.ID)
			require.NoError(t, err)
		}
	}

	tags, err := tagReader.Tags(ctx, "user 1", "")
	require.NoError(t, err)
	expected := []string{"tag", "yolo", "hello", "pizza", "world", "yoloy"}
	require.Equal(t, expected, tags)

	tags, err = tagReader.Tags(ctx, "user 1", "y")
	require.NoError(t, err)
	expected = []string{"yolo", "yoloy"}
	require.Equal(t, expected, tags)
}
