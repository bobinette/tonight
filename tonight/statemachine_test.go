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
			task:       Task{Log: []Log{{Type: LogTypeStart}, {Type: LogTypeProgress}}},
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
			task:       Task{Log: []Log{{Type: LogTypeProgress, Completion: 100}}},
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
			task:       Task{Log: []Log{{Type: LogTypeStart}, {Type: LogTypeProgress}}},
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
			task:       Task{Log: []Log{{Type: LogTypeProgress, Completion: 100}}},
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
			task:       Task{Log: []Log{{Type: LogTypeStart}, {Type: LogTypeProgress}}},
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
			task:       Task{Log: []Log{{Type: LogTypeProgress, Completion: 100}}},
			logType:    LogTypeWontDo,
			shouldFail: true,
		},

		// PROGRESS
		"progress + empty: ok": {
			task:       Task{Log: nil},
			logType:    LogTypeProgress,
			shouldFail: false,
		},
		"progress + started: ok": {
			task:       Task{Log: []Log{{Type: LogTypeStart}, {Type: LogTypeProgress}}},
			logType:    LogTypeProgress,
			shouldFail: false,
		},
		"progress + paused: ok": {
			task:       Task{Log: []Log{{Type: LogTypeStart}, {Type: LogTypePause}}},
			logType:    LogTypeProgress,
			shouldFail: false,
		},
		"progress + wont do: nok": {
			task:       Task{Log: []Log{{Type: LogTypeWontDo}}},
			logType:    LogTypeProgress,
			shouldFail: true,
		},
		"progress + done: nok": {
			task:       Task{Log: []Log{{Type: LogTypeProgress, Completion: 100}}},
			logType:    LogTypeProgress,
			shouldFail: true,
		},

		// COMMENT
		"comment + empty: ok": {
			task:       Task{Log: nil},
			logType:    LogTypeComment,
			shouldFail: false,
		},
		"comment + started: ok": {
			task:       Task{Log: []Log{{Type: LogTypeStart}, {Type: LogTypeProgress}}},
			logType:    LogTypeComment,
			shouldFail: false,
		},
		"comment + paused: ok": {
			task:       Task{Log: []Log{{Type: LogTypeStart}, {Type: LogTypePause}}},
			logType:    LogTypeComment,
			shouldFail: false,
		},
		"comment + wont do: nok": {
			task:       Task{Log: []Log{{Type: LogTypeWontDo}}},
			logType:    LogTypeComment,
			shouldFail: false,
		},
		"comment + done: nok": {
			task:       Task{Log: []Log{{Type: LogTypeProgress, Completion: 100}}},
			logType:    LogTypeComment,
			shouldFail: false,
		},
	}

	for name, test := range tests {
		can := isTransitionAllowed(test.task, test.logType)
		assert.Equal(t, test.shouldFail, !can, "%s - %s", name, test.logType)
	}
}
