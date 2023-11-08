package uptime

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUptime(t *testing.T) {
	startTime = time.Now().Add(-time.Hour * 24) // Set start time to 24 hours ago
	expected := "1 day ago"

	// Call the Uptime function
	actual := Uptime()

	// Check that the result contains the expected string
	assert.Equal(t, expected, actual)
}
