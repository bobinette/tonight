package tonighttest

import (
	"context"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	"github.com/bobinette/tonight"
)

func TestStores(
	t *testing.T,
	projectStore tonight.ProjectStore,
	taskStore tonight.TaskStore,
	userStore tonight.UserStore,
) {
	ctx := context.Background()

	user := tonight.User{ID: "testuser"}
	require.NoError(t, userStore.Ensure(ctx, &user))

	project := tonight.Project{
		UUID:      uuid.NewV1(),
		Name:      "Test project",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, projectStore.Upsert(ctx, project, user))
	// Create another project to make sure it is retrieved as well
	require.NoError(t, projectStore.Upsert(ctx, tonight.Project{
		UUID:      uuid.NewV1(),
		Name:      "Test other project",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, user))

	task := tonight.Task{
		UUID:      uuid.NewV1(),
		Title:     "Test task",
		Status:    tonight.TaskStatusTODO,
		Project:   project,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, taskStore.Upsert(ctx, task))

	otherTask := tonight.Task{
		UUID:      uuid.NewV1(),
		Title:     "Test task other",
		Status:    tonight.TaskStatusTODO,
		Project:   project,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, taskStore.Upsert(ctx, otherTask))

	ps, err := projectStore.List(ctx, user)
	require.NoError(t, err)
	require.Len(t, ps, 2)

	titles := make([]string, len(ps[0].Tasks))
	for i, task := range ps[0].Tasks {
		titles[i] = task.Title
		require.Equal(t, tonight.TaskStatusTODO, task.Status)
	}
	require.Equal(t, []string{task.Title, otherTask.Title}, titles)

	ps, err = projectStore.List(ctx, tonight.User{ID: "unkown"})
	require.NoError(t, err)
	require.Len(t, ps, 0)

	retrievedTask, err := taskStore.Get(ctx, task.UUID, user)
	require.NoError(t, err)
	require.Equal(t, task.Title, retrievedTask.Title)
	require.Equal(t, task.Status, retrievedTask.Status)
}
