package tonight

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
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
	}

	for name, test := range tests {
		task := parse(test.content)
		assert.Equal(t, test.expected, task, name)
	}
}
