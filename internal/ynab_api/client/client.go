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

func (apiClient YnabApiClient) httpRequest(method string, path string, body io.Reader, params map[string]string) *http.Response {
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

	metrics.ApiCallCounter.Inc()
	return response
}

func (apiClient YnabApiClient) GetBudgets() budget.BudgetData {
	var budgetData budget.BudgetResponseData

	params := map[string]string{
		"include_accounts": "true",
	}

	response := apiClient.httpRequest("GET", "/budgets", nil, params)

	rawResponseData, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		log.Fatalf("Error getting budgets [%s]", err)
	}

	err = json.Unmarshal(rawResponseData, &budgetData)
	if err != nil {
		log.Fatalf("Error unmarshaling budget data [%s]", err)
	}

	return budgetData.Data
}
