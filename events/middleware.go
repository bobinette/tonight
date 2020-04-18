package events

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/bobinette/tonight/auth"
)

// A Store should store and retrieve Events.
type Store interface {
	// Store e in the database/store.
	Store(ctx context.Context, e Event) error

	// List all the events from the store. List takes
	// a channel as input to make it more convenient
	// to scroll through all the events stored.
	List(ctx context.Context, ch chan<- Event) error
}

func Middleware(typ EventType, store Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			defer req.Body.Close()

			eventUUID, err := uuid.NewUUID()
			if err != nil {
				return err
			}

			var entityUUID uuid.UUID
			if entityUUIDStr := c.Param("uuid"); entityUUIDStr != "" {
				u, err := uuid.Parse(c.Param("uuid"))
				if err != nil {
					return err
				}
				entityUUID = u
			} else {
				entityUUID = eventUUID
			}

			// So that the downstream can use it as well if needed
			c.Set("event_uuid", eventUUID)

			var buf bytes.Buffer
			teeReader := io.TeeReader(req.Body, &buf)
			body, err := ioutil.ReadAll(teeReader)
			if err != nil {
				return err
			}

			user, err := auth.ExtractUser(c)
			if err != nil {
				return err
			}

			evt := Event{
				UUID:       eventUUID,
				Type:       typ,
				EntityUUID: entityUUID,
				UserID:     user.ID,
				Payload:    body, // Add the params and the query params!
				CreatedAt:  time.Unix(eventUUID.Time().UnixTime()),
			}

			ctx := req.Context()
			if err := store.Store(ctx, evt); err != nil {
				return err
			}

			req.Body = ioutil.NopCloser(&buf)
			return next(c)
		}
	}
}
