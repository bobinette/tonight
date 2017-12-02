package tonight

import (
	"context"
	"time"
)

type Task struct {
	ID          uint
	Title       string
	Description string

	Priority int
	Rank     uint
	Tags     []string

	Duration string
	Deadline *time.Time

	Done   bool
	DoneAt *time.Time

	Completion int // max([log.Completion for log in Log])
	Log        []Log

	Dependencies []Dependency

	CreatedAt time.Time
}

type LogType string

const (
	LogTypeCompletion LogType = "COMPLETION"
	LogTypePause              = "PAUSE"
	LogTypeStart              = "START"
	LogTypeLog                = "LOG"
)

type Log struct {
	Type        LogType
	Completion  int
	Description string

	CreatedAt time.Time
}

type Dependency struct {
	ID   uint
	Done bool
}

type Planning struct {
	ID uint

	Duration string

	Dismissed bool
	StartedAt time.Time

	Tasks []Task
}

type TaskRepository interface {
	List(ctx context.Context, ids []uint) ([]Task, error)
	Create(ctx context.Context, t *Task) error
	Update(ctx context.Context, t *Task) error

	MarkDone(ctx context.Context, taskID uint, log Log) error
	UpdateRanks(ctx context.Context, ranks map[uint]uint) error

	Delete(ctx context.Context, taskID uint) error

	StartPlanning(ctx context.Context, userID uint, duration string, taskIDs []uint) (Planning, error)
	DismissPlanning(ctx context.Context, userID uint) error
	CurrentPlanning(ctx context.Context, userID uint) (Planning, error)
}

type TaskIndex interface {
	Search(ctx context.Context, q string, done bool, allowedIDs []uint) ([]uint, error)
	Index(ctx context.Context, task Task) error
	Delete(ctx context.Context, taskID uint) error
}
