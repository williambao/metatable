package model

import "time"
import "strings"

type RecordCell struct {
	BaseModel   `xorm:"extends"`
	TableId     string     `json:"table_id" xorm:"varchar(20)  NOT NULL"`
	ColumnId    string     `json:"column_id" xorm:"varchar(20) NOT NULL"`
	ColumnType  ColumnType `json:"column_type" xorm:"varchar(20) NOT NULL"`
	RecordId    string     `json:"record_id" xorm:"varchar(20) NOT NULL"`
	UserId      string     `json:"user_id" xorm:"varchar(20) NOT NULL"`
	Key         string     `json:"key" xorm:"varchar(20)"`
	Value       string     `json:"value"`
	DateValue   *time.Time `json:"-"`
	NumberValue float64    `json:"-"`
	BoolValue   bool       `json:"-"`
}

func GetRecordCells(recordId string) ([]RecordCell, error) {
	list := make([]RecordCell, 0)
	err := db.Where("record_id = ?", recordId).Find(&list)
	return list, err
}

// 返回数据给前台, 有些字段可能要处理下

func (r *RecordCell) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	var value interface{}
	switch r.ColumnType {
	case ColumnTypeBoolean:
		value = r.BoolValue
	case ColumnTypeDatetime:
		value = r.DateValue
	case ColumnTypeNumber:
		value = r.NumberValue
	case ColumnTypeSelection:
		value = strings.Split(r.Value, ",")
	default:
		value = r.Value
	}
	result[r.Key] = value
	return result
}
