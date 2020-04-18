package mysql

import (
	"context"
	"database/sql"

	"github.com/bobinette/tonight/events"
)

type EventStore struct {
	db *sql.DB
}

func NewEventStore(db *sql.DB) EventStore {
	return EventStore{db: db}
}

func (s EventStore) Store(ctx context.Context, e events.Event) error {
	query := `
INSERT INTO events (uuid, type, entity_uuid, user_id, payload, created_at)
VALUES (?, ?, ?, ?, ?, ?)
`
	_, err := s.db.ExecContext(
		ctx,
		query,
		e.UUID,
		e.Type,
		e.EntityUUID,
		e.UserID,
		e.Payload,
		e.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s EventStore) List(ctx context.Context, ch chan<- events.Event) error {
	stmt, err := s.db.PrepareContext(ctx, `
SELECT uuid, type, entity_uuid, user_id, payload, created_at
FROM events
ORDER BY created_at ASC
LIMIT ?
OFFSET ?
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	limit := 50
	offset := 0

mainloop:
	for {
		select {
		case <-ctx.Done():
			break mainloop
		default:
		}
		sent, err := s.scrollAt(ctx, ch, stmt, limit, offset)
		if err != nil {
			return err
		}

		if sent == 0 {
			break
		}

		offset += limit
	}

	if err := stmt.Close(); err != nil {
		return err
	}

	return nil
}

func (s EventStore) scrollAt(ctx context.Context, ch chan<- events.Event, stmt *sql.Stmt, args ...interface{}) (int, error) {
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	sent := 0
	for rows.Next() {
		var evt events.Event
		err := rows.Scan(
			&evt.UUID,
			&evt.Type,
			&evt.EntityUUID,
			&evt.UserID,
			&evt.Payload,
			&evt.CreatedAt,
		)
		if err != nil {
			return 0, err
		}

		select {
		case <-ctx.Done():
			return 0, nil
		default:
		}
		ch <- evt
		sent++
	}

	if err := rows.Close(); err != nil {
		return 0, err
	}

	return sent, nil
}
