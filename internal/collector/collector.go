package collector

import (
	"log"

	"github.com/mcbobke/ynab-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type YnabCollector struct {
	Token string
}

func New(ynabToken string) YnabCollector {
	return YnabCollector{Token: ynabToken}
}

func (ynabCollector YnabCollector) Describe(chan<- prometheus.Desc) {
	log.Println("Not implemented")
}

func (ynabcollection YnabCollector) Collect(chan<- prometheus.Metric) {
	metrics.ApiCallCounter.Inc()
	log.Println("Not implemented")
}
