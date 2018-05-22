package tonight

import (
	"context"
	"errors"
	"regexp"
	"time"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type LogType string

const (
	LogTypeProgress LogType = "PROGRESS"
	LogTypePause            = "PAUSE"
	LogTypeStart            = "START"
	LogTypeWontDo           = "WONT_DO"
	LogTypeComment          = "COMMENT"
	LogTypePostpone         = "POSTPONE"
)

type Status int

const (
	StatusPending Status = iota
	StatusDone
	StatusWontDo
)

var (
	postponedUntilRegex = regexp.MustCompile(`postponed until (\d{4}-\d{2}-\d{2})`)
)

func StatusFromString(s string) Status {
	switch s {
	case "pending":
		return StatusPending
	case "done":
		return StatusDone
	case "won't do":
		return StatusWontDo
	}
	return StatusPending
}

func (ds Status) String() string {
	switch ds {
	case StatusPending:
		return "pending"
	case StatusDone:
		return "done"
	case StatusWontDo:
		return "won't do"
	}
	return ""
}

type Task struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Priority int      `json:"priority"`
	Rank     uint     `json:"rank"`
	Tags     []string `json:"tags"`

	Duration string     `json:"duration"`
	Deadline *time.Time `json:"deadline"`

	Score float64 `json:"score"`

	Log []Log `json:"log"`

	Dependencies []Dependency `json:"dependencies"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (t Task) LeftDuration() time.Duration {
	d, err := time.ParseDuration(t.Duration)
	if err != nil {
		d = time.Hour
	}

	return time.Duration(100-t.Completion()) * d / time.Duration(100)
}

func (t Task) Completion() int {
	c := 0
	for _, log := range t.Log {
		if log.Completion > c {
			c = log.Completion
		}
	}
	return c
}

func (t Task) Done() Status {
	for _, log := range t.Log {
		if log.Completion == 100 {
			return StatusDone
		}

		if log.Type == LogTypeWontDo {
			return StatusWontDo
		}
	}

	return StatusPending
}

func (t Task) DoneAt() *time.Time {
	for _, log := range t.Log {
		if log.Completion == 100 || log.Type == LogTypeWontDo {
			doneAt := log.CreatedAt
			return &doneAt
		}
	}
	return nil
}

func (t Task) PostponedUntil() *time.Time {
	for i := len(t.Log); i > 0; i-- {
		log := t.Log[i-1]
		if log.Type != LogTypePostpone {
			continue
		}

		match := postponedUntilRegex.FindStringSubmatch(log.Description)
		if len(match) == 0 {
			continue
		}

		if t, err := time.Parse("2006-01-02", match[1]); err == nil {
			return &t
		}
	}

	return nil
}

type Log struct {
	Type        LogType `json:"type"`
	Completion  int     `json:"completion"`
	Description string  `json:"description"`

	CreatedAt time.Time `json:"createdAt"`
}

type Dependency struct {
	ID    uint   `json:"id"`
	Done  bool   `json:"done"`
	Title string `json:"title"`
}

type TaskRepository interface {
	List(ctx context.Context, ids []uint) ([]Task, error)
	Create(ctx context.Context, t *Task) error
	Update(ctx context.Context, t *Task) error

	Log(ctx context.Context, taskID uint, log Log) error
	UpdateRanks(ctx context.Context, ranks map[uint]uint) error

	DependencyTrees(ctx context.Context, taskID uint) ([]Task, error)

	Delete(ctx context.Context, taskID uint) error

	All(ctx context.Context) ([]Task, error)
}

type TaskSearchParameters struct {
	IDs      []uint
	Q        string
	Statuses []Status
	SortBy   string
}

type TaskIndex interface {
	Search(ctx context.Context, p TaskSearchParameters) ([]uint, error)
	Index(ctx context.Context, task Task) error
	Delete(ctx context.Context, taskID uint) error
}
