package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mcbobke/ynab-exporter/internal/metrics"
	"github.com/mcbobke/ynab-exporter/internal/ynab_api/budget"
)

type YnabApiClient struct {
	Token  string
	client http.Client
}

func (apiClient YnabApiClient) httpRequest(method string, path string, body io.Reader, params map[string]string) []byte {
	baseUrl := "https://api.youneedabudget.com/v1"
	request, err := http.NewRequest(method, fmt.Sprintf("%s%s", baseUrl, path), body)
	if err != nil {
		log.Fatalf("Error creating a new request [%s]", err)
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiClient.Token))

	if params != nil {
		query := request.URL.Query()
		for key, val := range params {
			query.Add(key, val)
		}
		request.URL.RawQuery = query.Encode()
	}

	response, err := apiClient.client.Do(request)
	if err != nil {
		log.Fatalf("Error sending %s to %s [%s]", request.Method, request.URL, err)
	}
	defer response.Body.Close()

	rawResponseData, err := io.ReadAll(response.Body)
	metrics.ApiCallCounter.Inc()
	return rawResponseData
}

func (apiClient YnabApiClient) GetBudgets() budget.BudgetData {
	var budgetData budget.BudgetResponseData

	params := map[string]string{
		"include_accounts": "true",
	}

	response := apiClient.httpRequest("GET", "/budgets", nil, params)

	err := json.Unmarshal(response, &budgetData)
	if err != nil {
		log.Fatalf("Error unmarshaling budget data [%s]", err)
	}

	return budgetData.Data
}
