// Package budget declares structs to model the response data of the YNAB API /budgets endpoint.
package budget

// DateFormat describes the interpretation of timestamp data returned with a budget.
type DateFormat struct {
	Format string `json:"format"`
}

// CurrencyFormat describes the interpretation of the currency configured for a budget.
type CurrencyFormat struct {
	IsoCode          string `json:"iso_code,omitempty"`
	ExampleFormat    string `json:"example_format,omitempty"`
	DecimalDigits    int    `json:"decimal_digits,omitempty"`
	DecimalSeparator string `json:"decimal_separator,omitempty"`
	SymbolFirst      bool   `json:"symbol_first,omitempty"`
	GroupSeparator   string `json:"group_separator,omitempty"`
	CurrencySymbol   string `json:"currency_symbol,omitempty"`
	DisplaySymbol    bool   `json:"display_symbol,omitempty"`
}

// Account represents a single account within a budget in YNAB.
type Account struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	Type                string `json:"type,omitempty"`
	OnBudget            bool   `json:"on_budget,omitempty"`
	Closed              bool   `json:"closed,omitempty"`
	Note                string `json:"note,omitempty"`
	Balance             int    `json:"balance,omitempty"`
	ClearedBalance      int    `json:"cleared_balance,omitempty"`
	UnclearedBalance    int    `json:"uncleared_balance,omitempty"`
	TransferPayeeID     string `json:"transfer_payee_id,omitempty"`
	DirectImportLinked  bool   `json:"direct_import_linked,omitempty"`
	DirectImportInError bool   `json:"direct_import_in_error,omitempty"`
	Deleted             bool   `json:"deleted,omitempty"`
}

// Budget represents a single budget configured in YNAB.
type Budget struct {
	ID             string         `json:"id,omitempty"`
	Name           string         `json:"name,omitempty"`
	LastModifiedOn string         `json:"last_modified_on,omitempty"`
	FirstMonth     string         `json:"first_month,omitempty"`
	LastMonth      string         `json:"last_month,omitempty"`
	DateFormat     DateFormat     `json:"date_format,omitempty"`
	CurrencyFormat CurrencyFormat `json:"currency_format,omitempty"`
	Accounts       []Account      `json:"accounts,omitempty"`
}

// Data contains the list of all budgets account-wide and the account's default budget.
type Data struct {
	Budgets       []Budget `json:"budgets,omitempty"`
	DefaultBudget Budget   `json:"default_budget,omitempty"`
}

// Response is the top-level struct that the API response is unmarshaled into.
type Response struct {
	Data Data `json:"data,omitempty"`
}
