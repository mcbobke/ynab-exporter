package metrics

// Contains the metrics for this exporter itself
// Meta-metrics, if you will

import (
	"github.com/mcbobke/ynab-exporter/cmd/ynab-exporter/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	namespace        string = "ynab"
	subsystem        string = "exporter"
	BuildVersionInfo prometheus.Gauge
	ApiCallCounter   prometheus.Counter
)

func init() {
	BuildVersionInfo = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "ynab_exporter_build_info",
			Help:      "Build info for this instance of ynab-exporter",
			ConstLabels: prometheus.Labels{
				"build_version": version.BuildVersion,
				"build_time":    version.BuildTime,
			},
		},
	)

	ApiCallCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "api_calls_count",
			Help:      "Count of calls to the YNAB API",
		},
	)
}
