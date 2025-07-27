package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	NatsMessagesReceived = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "nats_messages_received_total",
			Help: "Total number of messages received from NATS",
		},
	)

	NatsMessageLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "nats_message_latency_us",
			Help: "Latency from price generation to NATS reader (in microseconds)",
			Buckets: []float64{
				10, 20, 50, 100, 200, 400, 800, 1600, 3200, 6400, 10000, 20000, 40000, 80000, 160000, 320000, 640000, 1000000,
			},
		},
	)

	WsServerDroppedMessages = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ws_server_dropped_messages_total",
			Help: "Total number of messages dropped due to full reader channel",
		},
	)

	EndToEndLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "ws_end_to_end_latency_us",
		Help:    "End-to-end latency from price generation to WS client receive (in µs)",
		Buckets: []float64{10, 50, 100, 200, 400, 800, 1600, 3200, 6400, 12800, 25600, 51200, 100000, 200000},
	})
)

func InitNatsMetrics() {
	prometheus.MustRegister(NatsMessagesReceived)
	prometheus.MustRegister(NatsMessageLatency)
	prometheus.MustRegister(WsServerDroppedMessages)
	prometheus.MustRegister(EndToEndLatency)
}
