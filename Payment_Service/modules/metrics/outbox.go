package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	OutboxPublished = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "payment", Subsystem: "outbox", Name: "published_total",
		Help: "Total outbox events published to Kafka",
	})
	OutboxFailed = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "payment", Subsystem: "outbox", Name: "failed_total",
		Help: "Total outbox publish failures",
	})
	OutboxLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "payment", Subsystem: "outbox", Name: "publish_seconds",
		Help:    "Latency seconds to publish messages",
		Buckets: prometheus.DefBuckets,
	})
)

func MustRegister() {
	prometheus.MustRegister(OutboxPublished, OutboxFailed, OutboxLatency)
}
