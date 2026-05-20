package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	GRPCRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "streaming_grpc_requests_total",
			Help: "Total gRPC requests received",
		}, []string{"method"},
	)

	PlaysTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "streaming_plays_total",
			Help: "Total number of recorded plays",
		},
	)
)

func Register() {
	prometheus.MustRegister(GRPCRequests)
	prometheus.MustRegister(PlaysTotal)
}
