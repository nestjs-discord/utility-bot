package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Prometheus struct {
	heartbeatLatencyGauge prometheus.Gauge
}

func New() *Prometheus {
	hlg := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "nestjs_discord",
		Subsystem: "utility_bot",
		Name:      "heartbeat_latency",
		Help:      "The latency between heartbeat acknowledgement and heartbeat send in seconds",
	})

	return &Prometheus{
		heartbeatLatencyGauge: hlg,
	}
}

func (p *Prometheus) SetHeartbeatLatency(value float64) {
	p.heartbeatLatencyGauge.Set(value)
}

func (p *Prometheus) ListenAndServe(addr string) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, nil)
}
