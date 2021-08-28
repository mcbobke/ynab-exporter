package collector

import (
	"net/http"

	"github.com/mcbobke/ynab-exporter/logging"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ModCollector YnabCollector
)

type YnabCollector struct {
}

func (YnabCollector) Collect(chan<- prometheus.Metric) {
	_, err := http.Get("http://www.google.com/robots.txt")
	if err != nil {
		logging.ModLogger.Println("Failed")
	} else {
		logging.ModLogger.Println("Success")
	}
}
