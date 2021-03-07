package model

type ColumnType string

const (
	ColumnTypeString    ColumnType = "string"
	ColumnTypeText      ColumnType = "text"
	ColumnTypeNumber    ColumnType = "number"
	ColumnTypeBoolean   ColumnType = "boolean"
	ColumnTypeImage     ColumnType = "image"
	ColumnTypeSelection ColumnType = "selection"
	ColumnTypeDatetime  ColumnType = "datetime"
	ColumnTypeSystem    ColumnType = "system"
	ColumnTypeRelation  ColumnType = "relation"
	ColumnTypeMember    ColumnType = "member"
)

var allColumnTypes = []ColumnType{
	ColumnTypeString, ColumnTypeText, ColumnTypeNumber,
	ColumnTypeBoolean, ColumnTypeImage, ColumnTypeSelection, ColumnTypeDatetime,
	ColumnTypeSystem, ColumnTypeRelation, ColumnTypeMember,
}

// 获取每个数据类型的查询条件
func GetColumnTypeFilters() map[ColumnType][]ColumnFilter {
	result := map[ColumnType][]ColumnFilter{}
	for _, col := range allColumnTypes {
		result[col] = col.filters()
	}
	return result
}

func (t ColumnType) filters() []ColumnFilter {
	list := make([]ColumnFilter, 0)
	switch t {
	case ColumnTypeString, ColumnTypeText:
		list = append(list, ColumnFilter{Name: "等于", Type: ColumnFilterTypeEqual, Value: "string"})
		list = append(list, ColumnFilter{Name: "不等于", Type: ColumnFilterTypeNotEqual, Value: "string"})
		list = append(list, ColumnFilter{Name: "包含", Type: ColumnFilterTypeContains, Value: "string"})
		list = append(list, ColumnFilter{Name: "不包含", Type: ColumnFilterTypeNotContains, Value: "string"})
		list = append(list, ColumnFilter{Name: "为空", Type: ColumnFilterTypeEmpty})
		list = append(list, ColumnFilter{Name: "不为空", Type: ColumnFilterTypeNotEmpty})
	case ColumnTypeNumber:
		list = append(list, ColumnFilter{Name: "等于", Type: ColumnFilterTypeEqual, Value: "number"})
		list = append(list, ColumnFilter{Name: "不等于", Type: ColumnFilterTypeNotEqual, Value: "number"})
		list = append(list, ColumnFilter{Name: "大于", Type: ColumnFilterTypeGreaterThan, Value: "number"})
		list = append(list, ColumnFilter{Name: "大于等于", Type: ColumnFilterTypeGreaterOrEqualThan, Value: "number"})
		list = append(list, ColumnFilter{Name: "小于", Type: ColumnFilterTypeLessThan, Value: "number"})
		list = append(list, ColumnFilter{Name: "小于等于", Type: ColumnFilterTypeLessOrEqualThan, Value: "number"})
		list = append(list, ColumnFilter{Name: "介于", Type: ColumnFilterTypeBetween, Value: "number"})
		list = append(list, ColumnFilter{Name: "为空", Type: ColumnFilterTypeEmpty})
		list = append(list, ColumnFilter{Name: "不为空", Type: ColumnFilterTypeNotEmpty})
	case ColumnTypeDatetime, ColumnTypeSystem:
		list = append(list, ColumnFilter{Name: "等于", Type: ColumnFilterTypeEqual, Value: "string"})
		list = append(list, ColumnFilter{Name: "不等于", Type: ColumnFilterTypeNotEqual, Value: "string"})
		list = append(list, ColumnFilter{Name: "大于", Type: ColumnFilterTypeGreaterThan, Value: "string"})
		list = append(list, ColumnFilter{Name: "大于等于", Type: ColumnFilterTypeGreaterOrEqualThan, Value: "string"})
		list = append(list, ColumnFilter{Name: "小于", Type: ColumnFilterTypeLessThan, Value: "string"})
		list = append(list, ColumnFilter{Name: "小于等于", Type: ColumnFilterTypeLessOrEqualThan, Value: "string"})
	// list = append(list, ColumnFilter{Name: "是今天", Type: ColumnFilterTypeToday})
	// list = append(list, ColumnFilter{Name: "是明天", Type: ColumnFilterTypeTomorrow})
	// list = append(list, ColumnFilter{Name: "上周", Type: ColumnFilterTypeLastWeek})
	case ColumnTypeSelection:
		list = append(list, ColumnFilter{Name: "等于", Type: ColumnFilterTypeEqual, Value: "option"})
		list = append(list, ColumnFilter{Name: "不等于", Type: ColumnFilterTypeNotEqual, Value: "option"})
		list = append(list, ColumnFilter{Name: "包含", Type: ColumnFilterTypeContains, Value: "option"})
		list = append(list, ColumnFilter{Name: "不包含", Type: ColumnFilterTypeNotContains, Value: "option"})
		list = append(list, ColumnFilter{Name: "为空", Type: ColumnFilterTypeEmpty})
		list = append(list, ColumnFilter{Name: "不为空", Type: ColumnFilterTypeNotEmpty})
	case ColumnTypeBoolean:
		list = append(list, ColumnFilter{Name: "已勾选", Type: ColumnFilterTypeSelected})
		list = append(list, ColumnFilter{Name: "未勾选", Type: ColumnFilterTypeNotSelected})
	case ColumnTypeImage:
		list = append(list, ColumnFilter{Name: "为空", Type: ColumnFilterTypeEmpty})
		list = append(list, ColumnFilter{Name: "不为空", Type: ColumnFilterTypeNotEmpty})
	}
	return list
}
