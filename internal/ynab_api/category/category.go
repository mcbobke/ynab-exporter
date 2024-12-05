// Package category declares structs to model the response data of the YNAB API /budgets/{budget_id}/categories endpoint.
package category

// Category describes a single budget category within a group.
type Category struct {
	ID                      string `json:"id,omitempty"`
	CategoryGroupID         string `json:"category_group_id,omitempty"`
	Name                    string `json:"name,omitempty"`
	Hidden                  bool   `json:"hidden,omitempty"`
	OriginalCategoryGroupID string `json:"original_category_group_id,omitempty"`
	Note                    string `json:"note,omitempty"`
	Budgeted                int    `json:"budgeted,omitempty"`
	Activity                int    `json:"activity,omitempty"`
	Balance                 int    `json:"balance,omitempty"`
	GoalType                string `json:"goal_type,omitempty"`
	GoalCreationMonth       string `json:"goal_creation_month,omitempty"`
	GoalTarget              int    `json:"goal_target,omitempty"`
	GoalTargetMonth         string `json:"goal_target_month,omitempty"`
	GoalPercentageComplete  int    `json:"goal_percentage_complete,omitempty"`
	GoalMonthsToBudget      int    `json:"goal_months_to_budget,omitempty"`
	GoalUnderFunded         int    `json:"goal_under_funded,omitempty"`
	GoalOverallFunded       int    `json:"goal_overall_funded,omitempty"`
	GoalOverallLeft         int    `json:"goal_overall_left,omitempty"`
	Deleted                 bool   `json:"deleted,omitempty"`
}

// Group contains the list of all categories within a specific group.
type Group struct {
	ID         string     `json:"id,omitempty"`
	Name       string     `json:"name,omitempty"`
	Hidden     bool       `json:"hidden,omitempty"`
	Deleted    bool       `json:"deleted,omitempty"`
	Categories []Category `json:"categories,omitempty"`
}

// Data contains the list of all category groups for a specific budget.
type Data struct {
	CategoryGroups  []Group `json:"category_groups,omitempty"`
	ServerKnowledge int     `json:"server_knowledge,omitempty"`
}

// Response is the top-level struct that the API response is unmarshaled into.
type Response struct {
	Data Data `json:"data,omitempty"`
}
