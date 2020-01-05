package inmem

import (
	"testing"

	"github.com/bobinette/tonight/tonighttest"
)

func TestStore(t *testing.T) {
	store := NewStore()
	tonighttest.TestStores(
		t,
		store.ProjectStore(),
		store.TaskStore(),
		store.UserStore(),
	)
}
