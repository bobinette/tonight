package tonighttest

import (
	"testing"

	"github.com/bobinette/tonight"
)

func TestStores(
	t *testing.T,
	projectStore tonight.ProjectStore,
	taskStore tonight.TaskStore,
	userStore tonight.UserStore,
) {
	// ctx := context.Background()

	// user := tonight.User{ID: "testuser"}
	// require.NoError(t, userStore.Ensure(ctx, &user))

	// project := tonight.Project{
	// 	UUID:      uuid.NewUUID(),
	// 	Name:      "Test project",
	// 	Slug:      "slug ",
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// }
	// require.NoError(t, projectStore.Upsert(ctx, project, user))
	// // Create another project to make sure it is retrieved as well
	// require.NoError(t, projectStore.Upsert(ctx, tonight.Project{
	// 	UUID:      uuid.NewUUID(),
	// 	Name:      "Test other project",
	// 	Slug:      "slug other",
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// }, user))

	// firstTask := tonight.Task{
	// 	UUID:      uuid.NewUUID(),
	// 	Title:     "Test task",
	// 	Status:    tonight.TaskStatusTODO,
	// 	Project:   project,
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// }
	// require.NoError(t, taskStore.Upsert(ctx, firstTask))

	// otherTask := tonight.Task{
	// 	UUID:      uuid.NewUUID(),
	// 	Title:     "Test task other",
	// 	Status:    tonight.TaskStatusTODO,
	// 	Project:   project,
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// }
	// require.NoError(t, taskStore.Upsert(ctx, otherTask))

	// ps, err := projectStore.List(ctx, user)
	// require.NoError(t, err)
	// require.Len(t, ps, 2)

	// titles := make([]string, len(ps[0].Tasks))
	// for i, task := range ps[0].Tasks {
	// 	titles[i] = task.Title
	// 	require.Equal(t, tonight.TaskStatusTODO, task.Status)
	// }
	// require.Equal(t, []string{firstTask.Title, otherTask.Title}, titles)

	// anotherTask := tonight.Task{
	// 	UUID:      uuid.NewUUID(),
	// 	Title:     "Another task",
	// 	Status:    tonight.TaskStatusTODO,
	// 	Project:   project,
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// }
	// require.NoError(t, taskStore.Upsert(ctx, anotherTask))
	// require.Equal(t, []string{firstTask.Title, otherTask.Title}, titles)

	// lastTask := tonight.Task{
	// 	UUID:      uuid.NewUUID(),
	// 	Title:     "Last task",
	// 	Status:    tonight.TaskStatusTODO,
	// 	Project:   project,
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// }
	// require.NoError(t, taskStore.Upsert(ctx, lastTask))

	// require.NoError(t, taskStore.Reorder(ctx, nil))
	// rankedUUIDs := []uuid.UUID{
	// 	otherTask.UUID,
	// 	firstTask.UUID,
	// }
	// require.NoError(t, taskStore.Reorder(ctx, rankedUUIDs))

	// ps, err = projectStore.List(ctx, user)
	// require.NoError(t, err)
	// require.Len(t, ps, 2)

	// titles = make([]string, len(ps[0].Tasks))
	// for i, task := range ps[0].Tasks {
	// 	titles[i] = task.Title
	// }
	// require.Equal(t, []string{otherTask.Title, firstTask.Title, anotherTask.Title, lastTask.Title}, titles)

	// ps, err = projectStore.List(ctx, tonight.User{ID: "unkown"})
	// require.NoError(t, err)
	// require.Len(t, ps, 0)

	// retrievedTask, err := taskStore.Get(ctx, firstTask.UUID, user)
	// require.NoError(t, err)
	// require.Equal(t, firstTask.Title, retrievedTask.Title)
	// require.Equal(t, firstTask.Status, retrievedTask.Status)

	// lastTask.Title = "Updated title"
	// lastTask.Status = tonight.TaskStatusDONE
	// require.NoError(t, taskStore.Upsert(ctx, lastTask))

	// retrievedTask, err = taskStore.Get(ctx, lastTask.UUID, user)
	// require.NoError(t, err)
	// require.Equal(t, lastTask.Title, retrievedTask.Title)
	// require.Equal(t, lastTask.Status, retrievedTask.Status)
}
