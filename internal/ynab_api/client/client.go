package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mcbobke/ynab-exporter/internal/metrics"
	"github.com/mcbobke/ynab-exporter/internal/ynab_api/budget"
	"github.com/mcbobke/ynab-exporter/internal/ynab_api/category"
	"go.uber.org/zap"
)

type YnabApiClient struct {
	Token  string
	Client *http.Client
	Logger *zap.SugaredLogger
}

func (apiClient YnabApiClient) httpRequest(method string, path string, body io.Reader, params map[string]string) ([]byte, error) {
	baseUrl := "https://api.youneedabudget.com/v1"
	apiErrorStatusCodes := map[int]string{
		400: "Bad API request",
		401: "Not authorized",
		403: "Subscription has lapsed",
		404: "Not found",
		429: "Rate limit reached",
		500: "Internal server error",
		503: "API unavailable",
	}

	request, err := http.NewRequest(method, fmt.Sprintf("%s%s", baseUrl, path), body)
	if err != nil {
		apiClient.Logger.Errorf("Error creating a new request [%s]", err)
		return []byte{}, fmt.Errorf("error creating a new request [%w]", err)
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiClient.Token))

	if params != nil {
		query := request.URL.Query()
		for key, val := range params {
			query.Add(key, val)
		}
		request.URL.RawQuery = query.Encode()
	}

	apiClient.Logger.Debugf("Sending %s to %s", request.Method, request.URL)
	response, err := apiClient.Client.Do(request)
	if err != nil {
		apiClient.Logger.Errorf("Error sending %s to %s [%s]", request.Method, request.URL, err)
		return []byte{}, fmt.Errorf("error sending %s to %s [%w]", request.Method, request.URL, err)
	}
	defer response.Body.Close()

	metrics.ApiCallCounter.Inc()

	logLine, ok := apiErrorStatusCodes[response.StatusCode]
	if ok {
		apiClient.Logger.Errorf("API response code %d indicates an error [%s]", response.StatusCode, logLine)
		return []byte{}, fmt.Errorf("API response code %d indicates an error [%s]", response.StatusCode, logLine)
	}

	rawResponseData, err := io.ReadAll(response.Body)
	if err != nil {
		apiClient.Logger.Errorf("Error reading API response data [%s]", err)
		return []byte{}, fmt.Errorf("error reading API response data [%w]", err)
	}

	return rawResponseData, nil
}

func (apiClient YnabApiClient) GetBudgets() (budget.BudgetData, error) {
	var budgetData budget.BudgetResponseData

	params := map[string]string{
		"include_accounts": "true",
	}

	response, err := apiClient.httpRequest("GET", "/budgets", nil, params)
	if err != nil {
		apiClient.Logger.Error("Failed to get budget data")
		return budget.BudgetData{}, errors.New("failed to get budget data")
	}

	err = json.Unmarshal(response, &budgetData)
	if err != nil {
		apiClient.Logger.Errorf("Error unmarshaling budget data [%s]", err)
		return budget.BudgetData{}, fmt.Errorf("error unmarshaling budget data [%w]", err)
	}

	return budgetData.Data, nil
}

func (apiClient YnabApiClient) GetCategories(budgetId string) (category.CategoryData, error) {
	var categoryData category.CategoryResponseData

	response, err := apiClient.httpRequest(
		"GET",
		fmt.Sprintf("/budgets/%s/categories", budgetId),
		nil,
		nil,
	)
	if err != nil {
		apiClient.Logger.Error("Failed to get category data")
		return category.CategoryData{}, errors.New("failed to get category data")
	}

	err = json.Unmarshal(response, &categoryData)
	if err != nil {
		apiClient.Logger.Errorf("Error unmarshaling category data [%s]", err)
		return category.CategoryData{}, fmt.Errorf("error unmarshaling category data [%w]", err)
	}

	return categoryData.Data, nil
}
