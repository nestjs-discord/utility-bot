package ratelimit_test

import (
	"github.com/nestjs-discord/utility-bot/pkg/ratelimit"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTTLMap(t *testing.T) {
	// Create a new TTLMap with maxTTL 1 second
	m := ratelimit.New(1)

	// Test IncrementUsage and Get methods
	m.IncrementUsage("key1")
	v := m.GetUsageCount("key1")
	assert.Equal(t, 1, v, "Expected value for key1 to be 1, but got %d", v)

	// Test that value is incremented properly
	m.IncrementUsage("key1")
	v = m.GetUsageCount("key1")
	assert.Equal(t, 2, v, "Expected value for key1 to be 2, but got %d", v)

	// Test that value is deleted after maxTTL
	time.Sleep(3 * time.Second)
	v = m.GetUsageCount("key1")
	assert.Equal(t, 0, v, "Expected value for key1 to be 2, but got %d", v)
}
