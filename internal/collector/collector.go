package collector

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mcbobke/ynab-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type YnabCollector struct {
	Token  string
	client http.Client
}

func New(ynabToken string) YnabCollector {
	client := http.Client{}
	return YnabCollector{Token: ynabToken, client: client}
}

func convertToUnits(milliunits int) int {
	return milliunits / 1000
}

func (ynabCollector YnabCollector) getRequest(path string) map[string]interface{} {
	var responseData []byte
	var returnData map[string]interface{}
	baseUrl := "https://api.youneedabudget.com/v1"

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", baseUrl, path), nil)
	if err != nil {
		log.Fatalln("Could not create a Request [%w]", err)
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ynabCollector.Token))

	response, err := ynabCollector.client.Do(request)
	if err != nil {
		log.Fatalf("Error sending GET to %s [%s]", path, err)
	}
	defer response.Body.Close()
	metrics.ApiCallCounter.Inc()

	_, err = response.Body.Read(responseData)
	if err != nil {
		log.Fatalln("Error reading response data [%w]", err)
	}

	err = json.Unmarshal(responseData, &returnData)
	if err != nil {
		log.Fatalln("Error unmarshling response data [%w]", err)
	}

	return returnData
}

func (ynabCollector YnabCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(ynabCollector, ch)
}

func (ynabCollector YnabCollector) Collect(chan<- prometheus.Metric) {
	// Get budgets
	var budgets []string
	budgetsResponse := ynabCollector.getRequest("/budgets")
	extractedBudgets := budgetsResponse["data"].(map[string]interface{})["budgets"]
	for budget, _ := range extractedBudgets {
		budgets = append(budgets, budget["id"])
	}
	log.Println("Budgets found: %s", budgets["data"])
}
