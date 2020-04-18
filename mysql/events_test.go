package mysql

import (
	"context"
	"testing"

	"github.com/bobinette/tonight/events"
	"github.com/stretchr/testify/require"
)

func TestEventStore(t *testing.T) {
	db, tearDown := setUp(t)
	defer tearDown()

	ctx := context.Background()
	clean := func() {
		db.ExecContext(ctx, `DELETE FROM events`)
		db.ExecContext(ctx, `DELETE FROM users`)
	}

	clean()
	defer clean()

	_, err := db.ExecContext(ctx, `
INSERT INTO users (id, name)
VALUES
	("events.user-1", "events.user-1"),
	("events.user-2", "events.user-2"),
	("events.user-3", "events.user-3")
`)
	require.NoError(t, err)

	store := NewEventStore(db)
	events.TestStore(t, store)
}
