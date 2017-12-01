package tonight

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatDuration(t *testing.T) {
	tests := map[time.Duration]string{
		1 * time.Hour:                                 "1h",
		2*time.Hour + 31*time.Minute + 43*time.Second: "2h31m43s",
	}

	for dur, expected := range tests {
		actual := formatDuration(dur)
		assert.Equal(t, expected, actual, "%v", dur)
	}
}
