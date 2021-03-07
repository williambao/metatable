package model

import (
	"errors"

	"github.com/rs/xid"
)

type ColumnOption struct {
	Id              string `json:"id" xorm:"varchar(20) pk"`
	Name            string `json:"name" form:"name" xorm:"varchar(100) NOT NULL"`
	TextColor       string `json:"text_color,omitempty" form:"text_color" xorm:"varchar(20)"`
	BackgroundColor string `json:"background_color,omitempty" form:"background_color" xorm:"varchar(20)"`
	URL             string `json:"url,omitempty" form:"url" xorm:"url"`
	IsDefault       bool   `json:"is_default" form:"is_default"`
}

func GetColumnOptions(tableId string, columnId string) ([]ColumnOption, error) {
	list := make([]ColumnOption, 0)
	err := db.Where("table_id = ? and column_id = ?", tableId, columnId).Find(&list)
	return list, err
}

func CreateColumnOption(option *ColumnOption) error {
	if len(option.Name) == 0 {
		return errors.New("请输入选项名称")
	}

	id := xid.New()

	option.Id = id.String()

	err := Insert(option)

	return err
}

func UpdateColumnOption(option *ColumnOption, cols ...string) error {

	err := UpdateById(option.Id, option, cols...)

	return err
}
