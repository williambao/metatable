package model

import (
	"errors"

	"github.com/rs/xid"
)

type TemplateColumn struct {
	BaseModel       `xorm:"extends"`
	TemplateId      string         `json:"template_id" form:"template_id" xorm:"varchar(20)"`
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
func (c *TemplateColumn) isValidData(data interface{}) error {

	return nil
}

func GetTemplateColumn(id string) *TemplateColumn {
	if len(id) == 0 {
		return nil
	}
	var item TemplateColumn
	GetById(id, &item)
	return &item
}

func GetTemplateColumns(templateId string) ([]TemplateColumn, error) {
	list := make([]TemplateColumn, 0)
	err := db.Where("template_id = ?", templateId).Asc("order").Find(&list)
	if err != nil {
		return list, err
	}

	return list, nil
}

func CreateTemplateColumn(column *TemplateColumn) error {
	if len(column.Name) == 0 {
		return errors.New("请输入列名")
	}

	if len(column.TemplateId) == 0 {
		return errors.New("无效的模板Id")
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

	session.Commit()

	UpdateTemplateColumnCount(column.TemplateId)

	return nil
}

func UpdateTemplateColumn(column *TemplateColumn, cols ...string) error {
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

	session.Commit()
	return nil
}

func DeleteTemplateColumn(columnId string, user *User, template *Template) error {

	column := GetTemplateColumn(columnId)
	if column == nil {
		return errors.New("无效的列")
	}

	session := db.NewSession()
	err := session.Begin()
	if err != nil {
		session.Rollback()
		return err
	}

	err = DeleteById(columnId, &TemplateColumn{})
	if err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		return err
	}

	UpdateTemplateColumnCount(template.Id)

	return nil
}

// 修改列顺序, data为列id列表, 按顺序更新order字段为 0 ~ len(data)
func UpdateTemplateColumnOrders(template *Template, data []string) error {
	session := db.NewSession()
	err := session.Begin()

	for index, id := range data {
		obj := new(TemplateColumn)
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
func isValidTemplateColumnType(columnType ColumnType) bool {
	for _, t := range allColumnTypes {
		if t == columnType {
			return true
		}
	}
	return false
}
