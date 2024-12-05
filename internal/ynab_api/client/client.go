// Package client provides a client for the YNAB API.
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

// YnabAPIClient is the API client.
type YnabAPIClient struct {
	Token  string
	Client *http.Client
	Logger *zap.SugaredLogger
}

func (apiClient YnabAPIClient) httpRequest(method string, path string, body io.Reader, params map[string]string) ([]byte, error) {
	baseURL := "https://api.youneedabudget.com/v1"
	apiErrorStatusCodes := map[int]string{
		400: "Bad API request",
		401: "Not authorized",
		403: "Subscription has lapsed",
		404: "Not found",
		429: "Rate limit reached",
		500: "Internal server error",
		503: "API unavailable",
	}

	request, err := http.NewRequest(method, fmt.Sprintf("%s%s", baseURL, path), body)
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

	metrics.APICallCounter.Inc()

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

// GetBudgets returns the budget data for a YNAB account.
func (apiClient YnabAPIClient) GetBudgets() (budget.Data, error) {
	var unmarshaledResponse budget.Response

	params := map[string]string{
		"include_accounts": "true",
	}

	response, err := apiClient.httpRequest("GET", "/budgets", nil, params)
	if err != nil {
		apiClient.Logger.Error("Failed to get budget data")
		return budget.Data{}, errors.New("failed to get budget data")
	}

	err = json.Unmarshal(response, &unmarshaledResponse)
	if err != nil {
		apiClient.Logger.Errorf("Error unmarshaling budget data [%s]", err)
		return budget.Data{}, fmt.Errorf("error unmarshaling budget data [%w]", err)
	}

	return unmarshaledResponse.Data, nil
}

// GetCategories returns the category data for a budget.
func (apiClient YnabAPIClient) GetCategories(budgetID string) (category.Data, error) {
	var unmarshaledResponse category.Response

	response, err := apiClient.httpRequest(
		"GET",
		fmt.Sprintf("/budgets/%s/categories", budgetID),
		nil,
		nil,
	)
	if err != nil {
		apiClient.Logger.Error("Failed to get category data")
		return category.Data{}, errors.New("failed to get category data")
	}

	err = json.Unmarshal(response, &unmarshaledResponse)
	if err != nil {
		apiClient.Logger.Errorf("Error unmarshaling category data [%s]", err)
		return category.Data{}, fmt.Errorf("error unmarshaling category data [%w]", err)
	}

	return unmarshaledResponse.Data, nil
}
