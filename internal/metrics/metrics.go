package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var IncomeCounter = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "income_count",
		Help: "Count of income earned",
	},
)
