package tonight

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

type EventType string

// An Event is used to record every mutation requested
// by users.
type Event struct {
	UUID      uuid.UUID
	Type      EventType
	Payload   []byte
	CreatedAt time.Time
}

// An EventStore should store and retrieve Events.
type EventStore interface {
	// Store e in the database/store.
	Store(ctx context.Context, e Event) error

	// List all the events from the store. List takes
	// a channel as input to make it more convenient
	// to scroll through all the events stored.
	List(ctx context.Context, ch chan<- Event) error
}
