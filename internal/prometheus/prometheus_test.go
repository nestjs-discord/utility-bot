package prometheus_test

import (
	"github.com/nestjs-discord/utility-bot/internal/prometheus"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestListenAndServe(t *testing.T) {
	p := prometheus.New()

	// Set a value for the heartbeat latency gauge
	p.SetHeartbeatLatency(2.5)

	// Listen address
	addr := "127.0.0.1:2112"

	// Start listening on the test server
	go func() {
		err := p.ListenAndServe(addr)
		assert.NoError(t, err)
	}()

	// Wait for the server to start listening
	time.Sleep(time.Second)

	// Send a GET request to the /metrics endpoint
	resp, err := http.Get("http://" + addr + "/metrics")
	assert.NoError(t, err)

	// Close response body
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Test response header status
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read body
	body, err := io.ReadAll(resp.Body)

	assert.NoError(t, err)
	assert.Contains(t, string(body), "nestjs_discord_utility_bot_heartbeat_latency 2.5")
}
