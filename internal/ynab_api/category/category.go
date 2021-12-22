package category

type Category struct {
	Id                      string `json:"id,omitempty"`
	CategoryGroupId         string `json:"category_group_id,omitempty"`
	Name                    string `json:"name,omitempty"`
	Hidden                  bool   `json:"hidden,omitempty"`
	OriginalCategoryGroupId string `json:"original_category_group_id,omitempty"`
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

type CategoryGroup struct {
	Id         string     `json:"id,omitempty"`
	Name       string     `json:"name,omitempty"`
	Hidden     bool       `json:"hidden,omitempty"`
	Deleted    bool       `json:"deleted,omitempty"`
	Categories []Category `json:"categories,omitempty"`
}

type CategoryData struct {
	CategoryGroups  []CategoryGroup `json:"category_groups,omitempty"`
	ServerKnowledge int             `json:"server_knowledge,omitempty"`
}

type CategoryResponseData struct {
	Data CategoryData `json:"data,omitempty"`
}
