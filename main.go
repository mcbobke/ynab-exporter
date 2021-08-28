package main

import (
	"net/http"

	"github.com/mcbobke/ynab-exporter/logging"
	"github.com/mcbobke/ynab-exporter/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var modVersion string = "0.0.1"

func main() {
	logging.ModLogger.Printf("Starting ynab-exporter version %s", modVersion)
	metrics.IncomeCounter.Inc()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe("localhost:9090", nil)
}
