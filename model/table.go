package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/xid"
)

type Table struct {
	BaseModel         `xorm:"extends"`
	Name              string `json:"name" form:"name" xorm:"varchar(200) NOT NULL"`
	Description       string `json:"description" form:"Text"`
	CoverURL          string `json:"cover_url" form:"cover_url" xorm:"cover_url"`
	OrganizationId    string `json:"organization_id" xorm:"varchar(20)"`
	OrganizationName  string `json:"organization_name"`
	IsPrivate         bool   `json:"is_private" form:"is_private"` // 私有表格
	IsSystem          bool   `json:"is_system" form:"is_system"`   // 私有表格
	IsCheckPermission bool   `json:""`
	UserId            string `json:"user_id" xorm:"varchar(20)"`
	IsEditable        bool   `json:"is_editable"`
	RecordCount       int64  `json:"record_count"`

	TemplateId      string `json:"template_id,omitempty"`
	TemplateVersion string `json:"string"`

	Columns    []Column    `json:"columns,omitempty" xorm:"-"`
	Views      []TableView `json:"views" xorm:"-"`
	Permission *TableUser  `json:"permission,omitempty" xorm:"-"`

	Organization *Organization `json:"organization,omitempty" xorm:"-"`
	// Conditions map[ColumnType][]ColumnFilter `json:"conditions" xorm:"-"`
}

func GetTableById(id string, user *User) (*Table, error) {
	if len(id) == 0 {
		return nil, errors.New("无效的表格id")
	}
	var item Table
	err := GetById(id, &item)
	if err != nil {
		return &item, err
	}

	// item.Conditions = GetColumnTypeFilters()

	// 把一些字列表填上
	columns, err := GetTableColumns(id)
	if err == nil {
		item.Columns = columns
	}

	views, err := GetTableViews(id, user)
	if err == nil {
		item.Views = views
	}

	if user != nil {
		perm, _ := GetTableUser(id, user.Id)
		// 有些字段不返回给前台了
		if perm != nil {
			perm.UpdatedAt = nil
			perm.CreatedAt = nil
			item.Permission = perm
		}

	}

	if item.OrganizationId != "" {
		item.Organization, _ = GetOrganizationById(item.OrganizationId)
	}

	return &item, nil
}

func CreateTable(table *Table) error {
	if len(table.Name) == 0 {
		return errors.New("请输入表格名")
	}

	session := db.NewSession()
	defer session.Close()
	id := xid.New()

	table.Id = id.String()

	if table.OrganizationId != "" {
		org, _ := GetOrganizationById(table.OrganizationId)
		table.OrganizationName = org.Name
	}

	err := session.Begin()

	_, err = session.Insert(table)
	if err != nil {
		session.Rollback()
		return err
	}

	// 创建表格时, 把当前用户存入到表格用户表, 并设置成表格管理员
	user := &TableUser{}
	user.Id = xid.New().String()
	user.TableId = table.Id
	user.UserId = table.UserId
	user.ViewData = TableDataPermissionAll
	// user.DeleteData = TableDataPermissionAll
	user.EditData = TableDataPermissionAll
	// user.IsAllowInsertData = true
	user.IsTableAdmin = true
	_, err = session.Insert(user)
	if err != nil {
		session.Rollback()
		return err
	}

	// 创建默认view
	view := &TableView{}
	view.Id = xid.New().String()
	view.TableId = table.Id
	view.UserId = table.UserId
	view.Name = "全部数据"
	view.FilterType = ViewFilterTypeAnd
	view.IsDefault = true
	_, err = session.Insert(view)
	if err != nil {
		session.Rollback()
		return err
	}

	if table.TemplateId != "" {
		err = addColumnsFromTemplate(table, table.TemplateId)
		if err != nil {
			session.Rollback()
			return err
		}
	} else {

		// 创建默认column
		column := &Column{}
		column.Id = xid.New().String()
		column.TableId = table.Id
		column.UserId = table.UserId
		column.Name = "标题"
		column.Width = 180
		column.Type = ColumnTypeString
		column.IsPrimary = true
		err = CreateColumn(column)
		if err != nil {
			session.Rollback()
			return err
		}
	}

	err = session.Commit()
	if err != nil {
		return err
	}

	if table.OrganizationId != "" {
		UpdateOrganizationTableCount(table.OrganizationId)
	}

	return nil
}

func addColumnsFromTemplate(table *Table, templateId string) error {
	columns, err := GetTemplateColumns(templateId)
	if err != nil {
		fmt.Println("XXXX:" + err.Error())
		return err
	}

	columnIds := make([]string, 0)

	for _, col := range columns {
		newCol := new(Column)
		newCol.Name = col.Name
		newCol.Type = col.Type
		newCol.TableId = table.Id
		newCol.Options = col.Options
		newCol.Config = col.Config
		newCol.UserId = table.UserId
		newCol.IsEditable = col.IsEditable
		newCol.IsPrimary = col.IsPrimary
		newCol.IsFreezed = col.IsFreezed
		newCol.IsRequired = col.IsRequired

		err = CreateColumn(newCol)
		if err != nil {
			removeColumns(columnIds)
			return err
		} else {
			columnIds = append(columnIds, newCol.Id)
		}

	}

	return nil
}

func removeColumns(ids []string) error {

	return nil
}

func UpdateTable(table *Table, cols ...string) error {
	err := UpdateById(table.Id, table, cols...)
	return err
}

type UserTable struct {
	Table        `xorm:"extends"`
	Username     string    `json:"username"`
	SubscribedAt time.Time `json:"subscribed_at"`
}

func GetUserTables(user *User, page int, pagesize int, sort string) ([]UserTable, error) {
	tables := make([]UserTable, 0)
	sortColumn := "t.id"
	if len(sort) > 0 {
		sortColumn = "t." + strings.TrimPrefix(sort, "-")
		if strings.HasPrefix(sort, "-") {
			sortColumn = sortColumn + " desc "
		}
	}

	err := db.SQL("select t.*, u.created_at as subscribed_at, ui.nickname as username "+
		"from `table` t "+
		"join user ui on ui.id = t.user_id "+
		"join table_user u on u.table_id = t.id "+
		"where u.user_id = ? order by ? limit ?, ?",
		user.Id,
		sortColumn,
		(page-1)*pagesize, pagesize,
	).Find(&tables)

	systemTables := make([]Table, 0)
	db.Where("is_system = 1").Find(&systemTables)
	if len(systemTables) > 0 {
		for _, t := range systemTables {
			exist := false
			for _, x := range tables {
				if x.Id == t.Id {
					exist = true
					break
				}
			}
			if !exist {
				t2 := UserTable{}
				t2.Table = t
				tables = append(tables, t2)
			}
		}
	}

	// for index, _ := range tables {
	// 	tables[index].Username = user.Nickname
	// }

	return tables, err
}

func DeleteTable(id string, user *User) error {
	session := db.NewSession()
	err := session.Begin()

	table, _ := GetTableById(id, user)

	// 删除用户
	_, err = session.Where("table_id = ?", id).Delete(&TableUser{})
	if err != nil {
		session.Rollback()
		return err
	}

	// 删除视图
	_, err = session.Where("table_id = ?", id).Delete(&TableView{})
	if err != nil {
		session.Rollback()
		return err
	}

	_, err = session.Where("table_id = ?", id).Delete(&Column{})
	if err != nil {
		session.Rollback()
		return err
	}

	_, err = session.Where("table_id = ?", id).Delete(&Record{})
	if err != nil {
		session.Rollback()
		return err
	}

	// 删除表格
	err = DeleteById(id, &Table{})
	if err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	if err != nil {
		return err
	}
	if table.OrganizationId != "" {
		UpdateOrganizationTableCount(table.OrganizationId)
	}

	return err
}

func updateTableRecordCount(tableId string) error {
	cnt, err := db.Where("table_id = ?", tableId).Count(&Record{})
	if err != nil {
		return err
	}

	table := new(Table)
	table.Id = tableId
	table.RecordCount = cnt
	err = UpdateById(tableId, table, "record_count")
	if err != nil {
		return err
	}
	return nil
}
