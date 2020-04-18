package events

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

// Event types
const (
	TaskCreate   EventType = "TaskCreate"
	TaskUpdate   EventType = "TaskUpdate"
	TaskDone     EventType = "TaskDone"
	TaskDelete   EventType = "TaskDelete"
	TasksReorder EventType = "TasksReorder"

	ReleaseCreate EventType = "ReleaseCreate"

	ProjectCreate       EventType = "ProjectCreate"
	ProjectUpdate       EventType = "ProjectUpdate"
	ProjectReorderTasks EventType = "ProjectReorderTasks"
)

// An Event is used to record every mutation requested
// by users.
type Event struct {
	UUID       uuid.UUID
	Type       EventType
	EntityUUID uuid.UUID
	UserID     string
	Payload    []byte
	CreatedAt  time.Time
}
