package model

import (
	"errors"

	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

type ColumnConfig struct {
	Type          string  `json:"type,omitempty" form:"type"`
	Max           float64 `json:"max,omitempty" form:"max"`
	Min           float64 `json:"min,omitempty" form:"min"`
	Precision     int     `json:"precision,omitempty" form:"precision"`
	Format        string  `json:"format,omitempty" form:"format"` // 有些字段需要格式化显示, 如日期
	IsShowPrecent bool    `json:"is_show_precent,omitempty" form:"is_show_precent"`
	IsContainTime bool    `json:"is_contain_time,omitempty" form:"is_contain_time"`
	IsMultiple    bool    `json:"is_multiple,omitempty" form:"is_multiple"` // 是否多选
}

type Column struct {
	BaseModel       `xorm:"extends"`
	TableId         string         `json:"table_id" form:"table_id" xorm:"varchar(20)"`
	UserId          string         `json:"user_id" form:"user_id" xorm:"varchar(20)"`
	Name            string         `json:"name" form:"name" xorm:"varchar(200)"`
	Type            ColumnType     `json:"type" form:"type" xorm:"varchar(50) NOT NULL"`
	Order           int            `json:"order" form:"order"`
	Width           int            `json:"width" form:"width"`
	TextAlign       string         `json:"text_align" xorm:"varchar(100)"`
	IsPrimary       bool           `json:"is_primary" form:"is_primary"` // 主Column, 用于一些地方显示
	IsFreezed       bool           `json:"is_freezed" form:"is_freezed"` // 是否已冻结, 用于前端浏览时固定
	IsEditable      bool           `json:"is_editable" form:"is_editable"`
	DefaultValue    string         `json:"default_value" form:"default_value" xorm:"varchar(200)"`
	IsRequired      bool           `json:"is_required" form:"is_required"`
	ReferenceTable  string         `json:"reference_table" form:"reference_table" xorm:"varchar(100)"`
	ReferenceColumn string         `json:"reference_column" form:"reference_column" xorm:"varchar(100)"`
	Config          ColumnConfig   `json:"config,omitempty" form:"config" xorm:"text json"`
	Options         []ColumnOption `json:"options,omitempty" form:"options" xorm:"text json"`
}

// 判断下数据是否符合column相关标准.
// 如数值型的不能以字符串格式存入数据库中
func (c *Column) isValidData(data interface{}) error {

	return nil
}

func GetColumnById(id string) *Column {
	if len(id) == 0 {
		return nil
	}
	var item Column
	GetById(id, &item)
	return &item
}

func GetTableColumns(tableId string) ([]Column, error) {
	list := make([]Column, 0)
	err := db.Where("table_id = ?", tableId).Asc("order").Find(&list)
	if err != nil {
		return list, err
	}

	// 若是下拉选项, 把相关options拉取过来
	// for index, item := range list {
	// 	if item.Type == ColumnTypeSelection {
	// 		options, err := GetColumnOptions(item.TableId, item.Id)
	// 		if err != nil {
	// 			logrus.Errorf("查询列Options出错误: %s", err.Error())
	// 		}
	// 		list[index].Options = options
	// 	}
	// }

	return list, nil
}

func CreateColumn(column *Column) error {
	if len(column.Name) == 0 {
		return errors.New("请输入列名")
	}

	if len(column.TableId) == 0 {
		return errors.New("无效的表格Id")
	}

	if len(column.UserId) == 0 {
		return errors.New("无效的用户Id")
	}

	if !isValidColumnType(column.Type) {
		return errors.New("无效的列类型")
	}

	// 给一些初始值
	if column.Width == 0 {
		column.Width = 180
	}

	session := db.NewSession()
	defer session.Close()
	id := xid.New()

	column.Id = id.String()

	err := session.Begin()

	// 当是下拉选项字段时, 把options存起来
	for index, option := range column.Options {
		if len(option.Id) == 0 {
			column.Options[index].Id = xid.New().String()
		}
	}

	// 更新下对齐方式
	switch column.Type {
	case ColumnTypeBoolean, ColumnTypeDatetime, ColumnTypeSystem:
		column.TextAlign = "center"
	case ColumnTypeNumber:
		column.TextAlign = "right"
	default:
		column.TextAlign = "left"
	}
	_, err = session.Insert(column)
	if err != nil {
		session.Rollback()
		return err
	}

	// 更新已有数据里的cells字段, 把当前column的默认值补上
	err = UpdateRecordCellValue(column)
	if err != nil {
		logrus.Errorf("新增column时更新cells默认值出错: %s", err.Error())
	}

	session.Commit()
	return nil
}

func UpdateColumn(column *Column, cols ...string) error {
	session := db.NewSession()
	defer session.Close()

	// 当是下拉选项字段时, 把options存起来
	for index, option := range column.Options {
		if len(option.Id) == 0 {
			column.Options[index].Id = xid.New().String()
		}
	}

	err := session.Begin()
	_, err = session.Cols(cols...).ID(column.Id).Update(column)
	if err != nil {
		session.Rollback()
		return err
	}

	// if column.Type == ColumnTypeSelection {
	// 	for index, option := range column.Options {
	// 		option.TableId = column.TableId
	// 		option.ColumnId = column.Id
	// 		option.UserId = column.UserId
	// 		err = CreateColumnOption(&option)
	// 		if err != nil {
	// 			session.Rollback()
	// 			return err
	// 		}
	// 		column.Options[index] = option
	// 	}
	// }

	session.Commit()
	return nil
}

func DeleteColumn(columnId string, user *User, table *Table) error {

	column := GetColumnById(columnId)
	if column == nil {
		return errors.New("无效的列")
	}

	session := db.NewSession()
	err := session.Begin()
	if err != nil {
		session.Rollback()
		return err
	}

	// 把record中跟此列相关数据删掉
	err = RemoveRecordCellValue(column)

	err = DeleteById(columnId, &Column{})
	if err != nil {
		session.Rollback()
		return err
	}

	//todo:
	// 把数据里的此字段相关值删除掉
	// records, _ := GetTableColumns(table.Id)
	// for index, record := range records {

	// }

	err = session.Commit()
	if err != nil {
		session.Rollback()
		return err
	}

	return nil
}

// 修改列顺序, data为列id列表, 按顺序更新order字段为 0 ~ len(data)
func UpdateColumnOrders(table *Table, data []string) error {
	session := db.NewSession()
	err := session.Begin()

	for index, id := range data {
		obj := new(Column)
		obj.Id = id
		obj.Order = index
		err = UpdateById(id, obj, "order")
		if err != nil {
			session.Rollback()
			return err
		}
	}

	err = session.Commit()
	return err
}

// 检查字段类型是否有效
func isValidColumnType(columnType ColumnType) bool {
	for _, t := range allColumnTypes {
		if t == columnType {
			return true
		}
	}
	return false
}
