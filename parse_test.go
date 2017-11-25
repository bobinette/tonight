package tonight

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tomorrow := time.Date(2017, 11, 24, 23, 59, 59, 0, time.UTC)

	tests := map[string]struct {
		content  string
		expected Task
	}{
		"only a title": {
			content:  "This is a title",
			expected: Task{Title: "This is a title"},
		},
		"with a description": {
			content:  "This is a title: now is the description",
			expected: Task{Title: "This is a title", Description: "now is the description"},
		},
		"with a description and a tag": {
			content: "This is a title: now is the description #tag",
			expected: Task{
				Title:       "This is a title",
				Description: "now is the description",
				Tags:        []string{"tag"},
			},
		},
		"with no description and a tag": {
			content: "This is a title #tag",
			expected: Task{
				Title: "This is a title",
				Tags:  []string{"tag"},
			},
		},
		"with a description and 3 tags": {
			content: "This is a title: now is the description #tag1 #tag2 #tag3",
			expected: Task{
				Title:       "This is a title",
				Description: "now is the description",
				Tags:        []string{"tag1", "tag2", "tag3"},
			},
		},
		"with a duration": {
			content:  "This is a title ~2h30m",
			expected: Task{Title: "This is a title", Duration: "2h30m"},
		},
		"with a description, 2 tags and the duration": {
			content: "This is a title: now is the description #tag1 #tag2 ~45m",
			expected: Task{
				Title:       "This is a title",
				Description: "now is the description",
				Tags:        []string{"tag1", "tag2"},
				Duration:    "45m",
			},
		},
		"with a priority": {
			content: "!!This is a title: now is the description #tag ~2h30m",
			expected: Task{
				Title:       "This is a title",
				Description: "now is the description",

				Priority: 2,
				Tags:     []string{"tag"},
				Duration: "2h30m",
			},
		},
		"with a deadline": {
			content: "This is a title >2017-11-24",
			expected: Task{
				Title:    "This is a title",
				Deadline: &tomorrow,
			},
		},
		"with dependencies": {
			content: "This is a title needs:1,2,3",
			expected: Task{
				Title: "This is a title",
				Dependencies: []Dependency{
					{ID: 1}, {ID: 2}, {ID: 3},
				},
			},
		},
	}

	for name, test := range tests {
		task := parse(test.content)
		assert.Equal(t, test.expected, task, name)
	}
}

func TestParseLog(t *testing.T) {
	tests := map[string]struct {
		content  string
		expected Log
	}{
		"without completion: 100%% by default": {
			content:  "this is the description",
			expected: Log{Description: "this is the description", Completion: 100},
		},
		"with completion": {
			content:  "25% this is the description",
			expected: Log{Description: "this is the description", Completion: 25},
		},
	}

	for name, test := range tests {
		log := parseLog(test.content)
		assert.Equal(t, test.expected.Completion, log.Completion, name)
		assert.True(t, test.expected.Description == log.Description, name)
	}
}
