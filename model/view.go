package model

import (
	"github.com/rs/xid"
)

const (
	ViewFilterTypeAnd = "and"
	ViewFilterTypeOr  = "or"
)

type TableViewFilter struct {
	ColumnId string           `json:"column_id" form:"column_id" `
	Operator ColumnFilterType `json:"operator" form:"operator"`
	Value    string           `json:"value" form:"value"`
}

type TableView struct {
	BaseModel   `xorm:"extends"`
	TableId     string            `json:"-" form:"table_id" xorm:"varchar(20)"`
	UserId      string            `json:"-" form:"user_id" xorm:"varchar(20)"`
	Name        string            `json:"name" form:"name"`
	IsDefault   bool              `json:"is_default"`
	IsPrivate   bool              `json:"is_private" form:"is_private"`
	Order       int               `json:"order" form:"order"`
	FilterType  string            `json:"filter_type" form:"filter_type"` // and / or
	HideColumns []string          `json:"hide_columns" form:"hide_columns" xorm:"text json"`
	Sorts       []string          `json:"sorts" form:"sorts" xorm:"text json"` // -column(desc) / column(asc)
	Filters     []TableViewFilter `json:"filters" form:"filters" xorm:"text json"`
}

func (t *TableView) BindMap(data map[string]interface{}) error {

	for k, v := range data {
		switch k {
		case "name":
			t.Name = v.(string)
		case "is_private":
			t.IsPrivate = v.(bool)
		case "order":
			t.Order = v.(int)
		case "filter_type":
			t.FilterType = v.(string)
		case "hide_columns":
			t.HideColumns = v.([]string)
		case "sorts":
			t.Sorts = v.([]string)
		case "filters":
			t.Filters = v.([]TableViewFilter)
		}
	}

	return nil
}

func GetTableView(viewId string) (*TableView, error) {
	var view TableView
	err := GetById(viewId, &view)
	return &view, err
}

func GetTableViews(tableId string, user *User) ([]TableView, error) {
	list := make([]TableView, 0)
	userId := ""
	if user != nil {
		userId = user.Id
	}
	err := db.SQL("select * from table_view where table_id = ? and (is_private = ? or user_id = ? and is_private = ?) order by `order`, created_at", tableId, false, userId, true).Find(&list)
	if err != nil {
		return list, err
	}

	return list, nil
}

func CreateTableView(view *TableView) error {
	if view.FilterType == "" {
		view.FilterType = ViewFilterTypeAnd
	}
	view.Id = xid.New().String()
	err := Insert(view)
	return err
}

func UpdateTableView(view *TableView, cols ...string) error {
	session := db.NewSession()
	defer session.Close()

	err := session.Begin()
	_, err = session.Cols(cols...).ID(view.Id).Update(view)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

func UpdateTableViewFrom(id string, data map[string]interface{}) error {
	_, err := db.Table(new(TableView)).ID(id).Update(data)
	return err
}
