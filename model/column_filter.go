package model

type ColumnFilterType string

// 筛选条件
const (
	ColumnFilterTypeEqual              ColumnFilterType = "equal"
	ColumnFilterTypeNotEqual           ColumnFilterType = "not_equal"
	ColumnFilterTypeContains           ColumnFilterType = "contains"
	ColumnFilterTypeNotContains        ColumnFilterType = "not_contains"
	ColumnFilterTypeBetween            ColumnFilterType = "between"
	ColumnFilterTypeEmpty              ColumnFilterType = "empty"
	ColumnFilterTypeNotEmpty           ColumnFilterType = "not_empty"
	ColumnFilterTypeIn                 ColumnFilterType = "in"
	ColumnFilterTypeNotIn              ColumnFilterType = "not_in"
	ColumnFilterTypeSelected           ColumnFilterType = "selected"
	ColumnFilterTypeNotSelected        ColumnFilterType = "not_selected"
	ColumnFilterTypeLessThan           ColumnFilterType = "less_than"
	ColumnFilterTypeLessOrEqualThan    ColumnFilterType = "less_or_equal_than"
	ColumnFilterTypeGreaterThan        ColumnFilterType = "greater_than"
	ColumnFilterTypeGreaterOrEqualThan ColumnFilterType = "greater_or_equal_than"
	ColumnFilterTypeToday              ColumnFilterType = "today"
	ColumnFilterTypeTomorrow           ColumnFilterType = "tomorrow"
	ColumnFilterTypeThisWeek           ColumnFilterType = "this_week"
	ColumnFilterTypeLastWeek           ColumnFilterType = "last_week"
	ColumnFilterTypeThisMonth          ColumnFilterType = "this_month"
	ColumnFilterTypeLastMonth          ColumnFilterType = "last_month"
	ColumnFilterTypeThisYear           ColumnFilterType = "this_year"
	ColumnFilterTypeLastYear           ColumnFilterType = "last_year"
)

type ColumnFilter struct {
	Name  string           `json:"name"`
	Type  ColumnFilterType `json:"type"`
	Value string           `json:"value,omitempty"`
}
