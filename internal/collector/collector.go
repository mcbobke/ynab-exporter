// Package collector provides an implementation of prometheus.Collector for metrics collection.
package collector

import (
	"net/http"

	"github.com/mcbobke/ynab-exporter/internal/ynab_api/client"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var (
	accountClearedBalanceDesc = prometheus.NewDesc(
		"ynab_account_cleared_balance",
		"Cleared balance of account",
		[]string{"budget_id", "budget_name", "account_name", "type"},
		nil,
	)
	accountUnclearedBalanceDesc = prometheus.NewDesc(
		"ynab_account_uncleared_balance",
		"Uncleared balance of account",
		[]string{"budget_id", "budget_name", "account_name", "type"},
		nil,
	)
	categoryBudgetedDesc = prometheus.NewDesc(
		"ynab_category_budgeted",
		"Amount budgeted to category",
		[]string{"budget_id", "budget_name", "category_group_name", "category_name"},
		nil,
	)
	categoryActivityDesc = prometheus.NewDesc(
		"ynab_category_activity",
		"Amount of activity in category",
		[]string{"budget_id", "budget_name", "category_group_name", "category_name"},
		nil,
	)
	categoryBalanceDesc = prometheus.NewDesc(
		"ynab_category_balance",
		"Category balance",
		[]string{"budget_id", "budget_name", "category_group_name", "category_name"},
		nil,
	)
)

// YnabCollector implements the prometheus.Collector interface.
type YnabCollector struct {
	Client client.YnabAPIClient
	Logger *zap.SugaredLogger
}

// New returns an instance of the YnabCollector struct.
func New(ynabToken string, logger *zap.SugaredLogger) YnabCollector {
	httpClient := http.DefaultClient
	ynabClient := client.YnabAPIClient{Token: ynabToken, Client: httpClient, Logger: logger}
	return YnabCollector{Client: ynabClient, Logger: logger}
}

// Describe delegates to prometheus.DescribeByCollect to send metrics descriptors to ch.
func (ynabCollector YnabCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(ynabCollector, ch)
}

// Collect is called by the default prometheus.Registry to collect data from the YNAB API and return timeseries.
func (ynabCollector YnabCollector) Collect(ch chan<- prometheus.Metric) {
	budgets, err := ynabCollector.Client.GetBudgets()
	if err != nil {
		ynabCollector.Logger.Error("Failed to get budget data")
		ch <- prometheus.NewInvalidMetric(
			prometheus.NewDesc("ynab_exporter_error",
				"Failed to get budget data", nil, nil),
			err)
		return
	}
	ynabCollector.Logger.Infof("Retrieved %d budgets", len(budgets.Budgets))

	for _, budget := range budgets.Budgets {
		categoryGroups, err := ynabCollector.Client.GetCategories(budget.ID)
		if err != nil {
			ynabCollector.Logger.Error("Failed to get category data")
			ch <- prometheus.NewInvalidMetric(
				prometheus.NewDesc("ynab_exporter_error",
					"Failed to get category data", nil, nil),
				err)
			return
		}
		ynabCollector.Logger.Debugf("Retrieved %d category groups for budget %s", len(categoryGroups.CategoryGroups), budget.ID)

		for _, categoryGroup := range categoryGroups.CategoryGroups {
			if categoryGroup.Hidden || categoryGroup.Deleted {
				continue
			}
			ynabCollector.Logger.Debugf("Retrieved %d categories for category group %s", len(categoryGroup.Categories), categoryGroup.Name)

			for _, category := range categoryGroup.Categories {
				if category.Hidden || category.Deleted {
					continue
				}

				ch <- prometheus.MustNewConstMetric(
					categoryBudgetedDesc,
					prometheus.GaugeValue,
					float64(category.Budgeted)/float64(1000),
					[]string{budget.ID, budget.Name, categoryGroup.Name, category.Name}...,
				)

				ch <- prometheus.MustNewConstMetric(
					categoryActivityDesc,
					prometheus.GaugeValue,
					float64(category.Activity)/float64(1000),
					[]string{budget.ID, budget.Name, categoryGroup.Name, category.Name}...,
				)

				ch <- prometheus.MustNewConstMetric(
					categoryBalanceDesc,
					prometheus.GaugeValue,
					float64(category.Balance)/float64(1000),
					[]string{budget.ID, budget.Name, categoryGroup.Name, category.Name}...,
				)
			}
		}

		ynabCollector.Logger.Debugf("Retrieved %d accounts for budget %s", len(budget.Accounts), budget.ID)
		for _, account := range budget.Accounts {
			if account.Closed || account.Deleted {
				continue
			}

			ch <- prometheus.MustNewConstMetric(
				accountClearedBalanceDesc,
				prometheus.GaugeValue,
				float64(account.ClearedBalance)/float64(1000),
				[]string{budget.ID, budget.Name, account.Name, account.Type}...,
			)

			ch <- prometheus.MustNewConstMetric(
				accountUnclearedBalanceDesc,
				prometheus.GaugeValue,
				float64(account.UnclearedBalance)/float64(1000),
				[]string{budget.ID, budget.Name, account.Name, account.Type}...,
			)
		}
	}
}
