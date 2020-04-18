package events

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var letters = []rune(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ&é"'(§è!çà)@#1234567890°_-$^ù£`)

func randBody(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func withTimeout(t *testing.T, f func(), d time.Duration) {
	done := make(chan struct{}, 0)
	go func() {
		f()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(d):
		t.Errorf("test timed out after %v", d)
	}
}

// TestStore will test the Store passed as parameter against a standard
// scenario.
//
// It will use user ids "events.user-1", "events.user-2" and "events.user-3". Make sure
// to add those in your database to avoid foreign key issues.
func TestStore(t *testing.T, store Store) {
	require := require.New(t)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	events := make([]Event, 20)
	for i := range events {
		body := randBody(10 + rng.Intn(300))

		eventUUID := uuid.Must(uuid.NewUUID())
		evt := Event{
			UUID:       eventUUID,
			Type:       ProjectCreate,
			EntityUUID: uuid.Must(uuid.NewUUID()),
			UserID:     fmt.Sprintf("events.user-%d", i%3+1),
			Payload:    []byte(body),
			CreatedAt:  time.Unix(eventUUID.Time().UnixTime()),
		}

		events[i] = evt
	}

	ctx := context.Background()
	for _, evt := range events {
		err := store.Store(ctx, evt)
		require.NoError(err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan Event, 0)
	go func() {
		defer close(ch)
		err := store.List(ctx, ch)
		require.NoError(err)
	}()

	withTimeout(t, func() {
		i := 0
		for evt := range ch {
			events[i].CreatedAt = events[i].CreatedAt.Truncate(time.Second)
			require.Equal(events[i], evt)
			i++
		}
	}, 2000*time.Millisecond)
}
