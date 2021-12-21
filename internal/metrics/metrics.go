package metrics

// Contains the metrics for this exporter itself
// Meta-metrics, if you will

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	namespace string = "ynab"
	subsystem string = "exporter"
	ApiCallCounter prometheus.Counter
)

func init() {
	ApiCallCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name: "api_calls_count",
			Help: "Count of calls to the YNAB API",
		},
	)
}
