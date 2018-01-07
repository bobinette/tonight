package tonight

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTransitionAllowed(t *testing.T) {
	tests := map[string]struct {
		task       Task
		logType    LogType
		shouldFail bool
	}{
		// START
		"start + empty: ok": {
			task:       Task{Log: nil},
			logType:    LogTypeStart,
			shouldFail: false,
		},
		"start + started: nok": {
			task:       Task{Log: []Log{{Type: LogTypeStart}, {Type: LogTypeCompletion}}},
			logType:    LogTypeStart,
			shouldFail: true,
		},
		"start + paused: ok": {
			task:       Task{Log: []Log{{Type: LogTypeStart}, {Type: LogTypePause}}},
			logType:    LogTypeStart,
			shouldFail: false,
		},
		"start + wont do: nok": {
			task:       Task{Log: []Log{{Type: LogTypeWontDo}}},
			logType:    LogTypeStart,
			shouldFail: true,
		},
		"start + done: nok": {
			task:       Task{Log: []Log{{Type: LogTypeCompletion, Completion: 100}}},
			logType:    LogTypeStart,
			shouldFail: true,
		},

		// PAUSE
		"pause + empty: nok": {
			task:       Task{Log: nil},
			logType:    LogTypePause,
			shouldFail: true,
		},
		"pause + started: ok": {
			task:       Task{Log: []Log{{Type: LogTypeStart}, {Type: LogTypeCompletion}}},
			logType:    LogTypePause,
			shouldFail: false,
		},
		"pause + paused: nok": {
			task:       Task{Log: []Log{{Type: LogTypeStart}, {Type: LogTypePause}}},
			logType:    LogTypePause,
			shouldFail: true,
		},
		"pause + wont do: nok": {
			task:       Task{Log: []Log{{Type: LogTypeWontDo}}},
			logType:    LogTypePause,
			shouldFail: true,
		},
		"pause + done: nok": {
			task:       Task{Log: []Log{{Type: LogTypeCompletion, Completion: 100}}},
			logType:    LogTypePause,
			shouldFail: true,
		},

		// WONT_DO
		"won't do + empty: ok": {
			task:       Task{Log: nil},
			logType:    LogTypeWontDo,
			shouldFail: false,
		},
		"won't do + started: ok": {
			task:       Task{Log: []Log{{Type: LogTypeStart}, {Type: LogTypeCompletion}}},
			logType:    LogTypeWontDo,
			shouldFail: false,
		},
		"won't do + paused: ok": {
			task:       Task{Log: []Log{{Type: LogTypeStart}, {Type: LogTypePause}}},
			logType:    LogTypeWontDo,
			shouldFail: false,
		},
		"won't do + wont do: nok": {
			task:       Task{Log: []Log{{Type: LogTypeWontDo}}},
			logType:    LogTypeWontDo,
			shouldFail: true,
		},
		"won't do + done: nok": {
			task:       Task{Log: []Log{{Type: LogTypeCompletion, Completion: 100}}},
			logType:    LogTypeWontDo,
			shouldFail: true,
		},

		// PROGRESS
		"progress + empty: ok": {
			task:       Task{Log: nil},
			logType:    LogTypeCompletion,
			shouldFail: false,
		},
		"progress + started: ok": {
			task:       Task{Log: []Log{{Type: LogTypeStart}, {Type: LogTypeCompletion}}},
			logType:    LogTypeCompletion,
			shouldFail: false,
		},
		"progress + paused: ok": {
			task:       Task{Log: []Log{{Type: LogTypeStart}, {Type: LogTypePause}}},
			logType:    LogTypeCompletion,
			shouldFail: false,
		},
		"progress + wont do: nok": {
			task:       Task{Log: []Log{{Type: LogTypeWontDo}}},
			logType:    LogTypeCompletion,
			shouldFail: true,
		},
		"progress + done: nok": {
			task:       Task{Log: []Log{{Type: LogTypeCompletion, Completion: 100}}},
			logType:    LogTypeCompletion,
			shouldFail: true,
		},
	}

	for name, test := range tests {
		can := isTransitionAllowed(test.task, test.logType)
		assert.Equal(t, test.shouldFail, !can, "%s - %s", name, test.logType)
	}
}
