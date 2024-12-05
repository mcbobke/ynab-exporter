// Package metrics contains meta-metrics for this exporter itself.
package metrics

import (
	"github.com/mcbobke/ynab-exporter/cmd/ynab-exporter/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	namespace        = "ynab"
	subsystem        = "exporter"
	buildVersionInfo prometheus.Gauge

	// APICallCounter is incremented every time the YNAB API is called.
	APICallCounter prometheus.Counter
)

func init() {
	buildVersionInfo = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "build_info",
			Help:      "Build info for this instance of ynab-exporter",
			ConstLabels: prometheus.Labels{
				"build_version": version.BuildVersion,
				"build_time":    version.BuildTime,
			},
		},
	)

	buildVersionInfo.Set(float64(1))

	APICallCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "api_calls_count",
			Help:      "Count of calls to the YNAB API",
		},
	)
}
