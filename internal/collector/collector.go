package collector

import (
	"log"

	"github.com/mcbobke/ynab-exporter/internal/ynab_api/client"
	"github.com/prometheus/client_golang/prometheus"
)

type YnabCollector struct {
	client client.YnabApiClient
}

func New(ynabToken string) YnabCollector {
	client := client.YnabApiClient{Token: ynabToken}
	return YnabCollector{client: client}
}

func convertToUnits(milliunits int) int {
	return milliunits / 1000
}

func (ynabCollector YnabCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(ynabCollector, ch)
}

func (ynabCollector YnabCollector) Collect(ch chan<- prometheus.Metric) {
	budgets := ynabCollector.client.GetBudgets()
	for _, val := range budgets.Budgets {
		log.Printf("Budget ID: %s", val.Id)
		log.Println("Accounts:", val.Accounts)
	}
}
