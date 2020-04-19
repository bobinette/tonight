package auth

import (
	"context"
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserStore struct {
	AllPermissionsFunc func(ctx context.Context, user User) (map[uuid.UUID]Permission, error)
}

func (m mockUserStore) Ensure(ctx context.Context, user *User) error {
	return nil
}

func (m mockUserStore) SetPermission(ctx context.Context, user User, projectUUID uuid.UUID, perm Permission) error {
	return nil
}

func (m mockUserStore) AllPermissions(ctx context.Context, user User) (map[uuid.UUID]Permission, error) {
	return m.AllPermissionsFunc(ctx, user)
}

func (m mockUserStore) Permission(ctx context.Context, user User, projectUUID uuid.UUID) (Permission, error) {
	return "", nil
}

func TestPermission_Includes(t *testing.T) {
	assert := assert.New(t)

	assert.True(ProjectOwn.Includes(ProjectOwn))
	assert.True(ProjectEdit.Includes(ProjectEdit))

	assert.True(ProjectOwn.Includes(ProjectEdit))
	assert.True(ProjectOwn.Includes(ProjectView))
	assert.False(ProjectView.Includes(ProjectOwn))
}
func TestPermission_AllowedProjects(t *testing.T) {
	require := require.New(t)

	edit := uuid.MustParse("7e242326-824b-11ea-bc55-0242ac130003")
	own := uuid.MustParse("d071db20-824e-11ea-bc55-0242ac130003")
	view := uuid.MustParse("0af31ff5-ee78-4d5a-b196-3c386b14666a")
	nothing := uuid.MustParse("28fbc636-fd91-4afd-a64a-6b9467632de2")

	uuids := []uuid.UUID{
		edit,
		own,
		view,
		nothing,
	}

	perms := map[uuid.UUID]Permission{
		uuids[0]: ProjectEdit,
		uuids[1]: ProjectOwn,
		uuids[2]: ProjectView,
	}

	store := mockUserStore{
		AllPermissionsFunc: func(ctx context.Context, user User) (map[uuid.UUID]Permission, error) {
			return perms, nil
		},
	}

	user := User{}
	ctx := context.Background()

	permissioner := Permissioner{store: store}

	projectUUIDs, err := permissioner.AllowedProjects(ctx, user, ProjectEdit)
	require.NoError(err)
	expected := asStrings(edit, own)
	require.Equal(expected, asStrings(projectUUIDs...))

	projectUUIDs, err = permissioner.AllowedProjects(ctx, user, ProjectView)
	require.NoError(err)
	expected = asStrings(edit, own, view)
	require.Equal(expected, asStrings(projectUUIDs...))

	projectUUIDs, err = permissioner.AllowedProjects(ctx, user, Permission("unkown"))
	require.NoError(err)
	expected = asStrings()
	require.Equal(expected, asStrings(projectUUIDs...))
}

func asStrings(uuids ...uuid.UUID) []string {
	strs := make([]string, len(uuids))
	for i, u := range uuids {
		strs[i] = u.String()
	}
	sort.Strings(strs)
	return strs
}
