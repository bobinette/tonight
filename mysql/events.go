package mysql

import (
	"context"
	"database/sql"

	"github.com/bobinette/tonight"
)

type EventStore struct {
	db *sql.DB
}

func NewEventStore(db *sql.DB) EventStore {
	return EventStore{db: db}
}

func (s EventStore) Store(ctx context.Context, e tonight.Event) error {
	query := `
INSERT INTO events (uuid, type, user_id, payload, created_at)
VALUES (?, ?, ?, ?, ?)
`
	_, err := s.db.ExecContext(ctx, query, e.UUID, e.Type, e.UserID, e.Payload, e.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s EventStore) List(ctx context.Context, ch chan<- tonight.Event) error {
	return nil
}
