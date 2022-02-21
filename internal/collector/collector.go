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
		[]string{"category_group_name", "category_name"},
		nil,
	)
	categoryActivityDesc = prometheus.NewDesc(
		"ynab_category_activity",
		"Amount of activity in category",
		[]string{"category_group_name", "category_name"},
		nil,
	)
	categoryBalanceDesc = prometheus.NewDesc(
		"ynab_category_balance",
		"Category balance",
		[]string{"category_group_name", "category_name"},
		nil,
	)
)

type YnabCollector struct {
	Client client.YnabApiClient
	Logger *zap.SugaredLogger
}

func New(ynabToken string, logger *zap.SugaredLogger) YnabCollector {
	httpClient := http.DefaultClient
	ynabClient := client.YnabApiClient{Token: ynabToken, Client: httpClient, Logger: logger}
	return YnabCollector{Client: ynabClient, Logger: logger}
}

func (ynabCollector YnabCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(ynabCollector, ch)
}

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
		categoryGroups, err := ynabCollector.Client.GetCategories(budget.Id)
		if err != nil {
			ynabCollector.Logger.Error("Failed to get category data")
			ch <- prometheus.NewInvalidMetric(
				prometheus.NewDesc("ynab_exporter_error",
					"Failed to get category data", nil, nil),
				err)
			return
		}
		ynabCollector.Logger.Debugf("Retrieved %d category groups for budget %s", len(categoryGroups.CategoryGroups), budget.Id)

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
					[]string{categoryGroup.Name, category.Name}...,
				)

				ch <- prometheus.MustNewConstMetric(
					categoryActivityDesc,
					prometheus.GaugeValue,
					float64(category.Activity)/float64(1000),
					[]string{categoryGroup.Name, category.Name}...,
				)

				ch <- prometheus.MustNewConstMetric(
					categoryBalanceDesc,
					prometheus.GaugeValue,
					float64(category.Balance)/float64(1000),
					[]string{categoryGroup.Name, category.Name}...,
				)
			}
		}

		ynabCollector.Logger.Debugf("Retrieved %d accounts for budget %s", len(budget.Accounts), budget.Id)
		for _, account := range budget.Accounts {
			if account.Closed || account.Deleted {
				continue
			}

			ch <- prometheus.MustNewConstMetric(
				accountClearedBalanceDesc,
				prometheus.GaugeValue,
				float64(account.ClearedBalance)/float64(1000),
				[]string{budget.Id, budget.Name, account.Name, account.Type}...,
			)

			ch <- prometheus.MustNewConstMetric(
				accountUnclearedBalanceDesc,
				prometheus.GaugeValue,
				float64(account.UnclearedBalance)/float64(1000),
				[]string{budget.Id, budget.Name, account.Name, account.Type}...,
			)
		}
	}
}
