package budget

type DateFormatStruct struct {
	Format string `json:"format"`
}

type CurrencyFormatStruct struct {
	IsoCode          string `json:"iso_code,omitempty"`
	ExampleFormat    string `json:"example_format,omitempty"`
	DecimalDigits    int    `json:"decimal_digits,omitempty"`
	DecimalSeparator string `json:"decimal_separator,omitempty"`
	SymbolFirst      bool   `json:"symbol_first,omitempty"`
	GroupSeparator   string `json:"group_separator,omitempty"`
	CurrencySymbol   string `json:"currency_symbol,omitempty"`
	DisplaySymbol    bool   `json:"display_symbol,omitempty"`
}

type Account struct {
	Id                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	Type                string `json:"type,omitempty"`
	OnBudget            bool   `json:"on_budget,omitempty"`
	Closed              bool   `json:"closed,omitempty"`
	Note                string `json:"note,omitempty"`
	Balance             int    `json:"balance,omitempty"`
	ClearedBalance      int    `json:"cleared_balance,omitempty"`
	UnclearedBalance    int    `json:"uncleared_balance,omitempty"`
	TransferPayeeId     string `json:"transfer_payee_id,omitempty"`
	DirectImportLinked  bool   `json:"direct_import_linked,omitempty"`
	DirectImportInError bool   `json:"direct_import_in_error,omitempty"`
	Deleted             bool   `json:"deleted,omitempty"`
}

type Budget struct {
	Id             string               `json:"id,omitempty"`
	Name           string               `json:"name,omitempty"`
	LastModifiedOn string               `json:"last_modified_on,omitempty"`
	FirstMonth     string               `json:"first_month,omitempty"`
	LastMonth      string               `json:"last_month,omitempty"`
	DateFormat     DateFormatStruct     `json:"date_format,omitempty"`
	CurrencyFormat CurrencyFormatStruct `json:"currency_format,omitempty"`
	Accounts       []Account            `json:"accounts,omitempty"`
}

type BudgetData struct {
	Budgets       []Budget `json:"budgets,omitempty"`
	DefaultBudget Budget   `json:"default_budget,omitempty"`
}

type BudgetResponseData struct {
	Data BudgetData `json:"data,omitempty"`
}
