package model

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"fmt"

	"github.com/rs/xid"
)

type Record struct {
	BaseModel `xorm:"extends"`
	TableId   string `json:"table_id" xorm:"varchar(20)"`
	UserId    string `json:"user_id"  xorm:"varchar(20)"`

	Version          int                    `json:"version"  xorm:"version"`
	Cells            map[string]interface{} `json:"cells" form:"cells" xorm:"text json"`
	UserName         string                 `json:"user_name" xorm:"-"`
	UpdateUserName   string                 `json:"update_user_name" xorm:"-"`
	UserAvatar       string                 `json:"create_user_avatar" xorm:"-"`
	UpdateUserAvatar string                 `json:"update_user_avatar" xorm:"-"`
}

func GetRecord(id string) (*Record, error) {
	var record Record
	err := GetById(id, &record)
	return &record, err
}

func GetTableRecords(tableId string) ([]Record, error) {
	list := make([]Record, 0)
	err := db.Where("table_id = ?", tableId).Find(&list)
	return list, err
}

func DeleteRecord(id string) error {

	session := db.NewSession()
	err := session.Begin()

	// // 先删除cells
	// _, err := session.Where("record_id = ?", id).Delete(&RecordCell{})
	// if err != nil {
	// 	session.Rollback()
	// 	return err
	// }

	err = DeleteById(id, &Record{})
	if err != nil {
		session.Rollback()
		return err
	}
	err = session.Commit()
	return err
}

// 批量删除
func DeleteRecordByBatch(ids []string, perm *TableUser, user *User) error {
	if perm.EditData == TableDataPermissionNone {
		return errors.New("无权限删除")
	}

	ids2 := make([]string, 0)
	if perm.EditData == TableDataPermissionUser {
		records := make([]Record, 0)
		err := db.In("id", ids).Cols("id", "user_id").Find(&records)
		if err != nil {
			return err
		}
		for _, r := range records {
			if r.UserId == user.Id {
				ids2 = append(ids2, r.Id)
			}
		}
	} else if perm.EditData == TableDataPermissionAll {
		ids2 = ids
	}

	cnt, err := db.In("id", ids2).Delete(&Record{})
	if cnt == 0 {
		return errors.New("无权限")
	}
	if int(cnt) != len(ids) {
		return errors.New("部分数据未能删除(无权限)")
	}
	return err
}

func CreateRecord(data *Record) error {
	//
	if len(data.TableId) == 0 {
		return errors.New("无效的表格id")
	}
	if len(data.UserId) == 0 {
		return errors.New("无效的用户id")
	}

	if data.Cells == nil {
		data.Cells = make(map[string]interface{})
	}

	if err := checkRecordColumnValues(data); err != nil {
		return err
	}

	data.Id = xid.New().String()
	err := Insert(data)
	if err != nil {
		return err
	}

	updateTableRecordCount(data.TableId)

	return nil
}

func UpdateRecordCells(record *Record) error {
	if err := checkRecordColumnValues(record); err != nil {
		return err
	}

	err := UpdateById(record.Id, record, "cells", "updated_by")
	return err
}

// 当增新column时, 更新以有数据里的cells字段
func UpdateRecordCellValue(column *Column) error {
	val := getDefaultRecordValue(column)
	mapVal := map[string]interface{}{column.Id: val}
	jsonB, err := json.Marshal(mapVal)
	if err != nil {
		return err
	}

	sql := "update `record` set cells = JSON_MERGE(cells, ?) where table_id = ?"
	_, err = db.Exec(sql, string(jsonB), column.TableId)

	return err
}

func RemoveRecordCellValue(column *Column) error {
	sql := fmt.Sprintf("update `record` set cells = JSON_REMOVE(cells, '$.%s') where table_id = ?", column.Id)
	_, err := db.Exec(sql, column.TableId)
	return err
}

func getDefaultRecordValue(column *Column) interface{} {
	var value interface{}

	switch column.Type {
	case ColumnTypeBoolean:
		value = false
	case ColumnTypeString, ColumnTypeText:
		value = ""
	case ColumnTypeImage:
		value = []interface{}{}
	case ColumnTypeNumber:
		value = 0
	case ColumnTypeSelection, ColumnTypeMember:
		value = []interface{}{}
	}

	return value
}

func GetViewRecords2(view *TableView, query string, page int, pagesize int) ([]Record, error) {
	list := make([]Record, 0)
	where := db.Where("table_id = ?", view.TableId)

	if len(query) > 0 {
		where = where.And("JSON_SEARCH(cells, 'all', '%?%') is not null", query)
	}

	if len(view.Filters) > 0 {
		// for _, filter := range view.Filters {
		// 	// where = where.Or()
		// }
	}

	// 排序
	if len(view.Sorts) > 0 {
		for _, item := range view.Sorts {
			// 以'-'开头的字段代表倒序, '-'后面跟的是column id
			desc := ""
			if strings.HasPrefix(item, "-") {
				item = strings.TrimPrefix(item, "-")
				desc = " desc "
			}
			where = where.OrderBy(fmt.Sprintf("json_extract(cells, '$.%s') %s", item, desc))
		}
	}

	err := where.Limit(pagesize, (page-1)*pagesize).Find(&list)

	// 如果选反了隐藏一些字段. 则直接不返回给前台了
	if len(view.HideColumns) > 0 {
		for index, item := range list {
			for _, colId := range view.HideColumns {
				delete(item.Cells, colId)
			}
			list[index] = item
		}
	}

	return list, err
}

// 用view筛选条件, 来查询并返回数据
func GetViewRecords(view *TableView, perm *TableUser, user *User, query string, page int, pagesize int) ([]Record, error) {
	list := make([]Record, 0)

	if perm.ViewData == TableDataPermissionNone {
		return list, nil
	}

	sql := "select * from record r where r.table_id = ? "

	// 返回自己创建的数据
	if perm.ViewData == TableDataPermissionUser {
		sql += " and user_id = '" + user.Id + "' "
	}
	if perm.ViewData == TableDataPermissionTeam {
		// teamUsers := make([]OrganizationUser, 0)

		// sql += " and user_id = '" + user.Id + "' "

	}

	if len(query) > 0 {
		query := fmt.Sprintf(" and JSON_SEARCH(cells, 'all', '%%%s%%') is not null", query)
		sql += query
	}

	if len(view.Filters) > 0 {
		wheres := []string{}

		for _, filter := range view.Filters {
			str := ""
			switch filter.Operator {
			case ColumnFilterTypeEqual:
				str = fmt.Sprintf("json_extract(cells, '$.%s') = '%s' ", filter.ColumnId, filter.Value)
			case ColumnFilterTypeNotEmpty: // 不为空
				str = fmt.Sprintf("json_length(cells, '$.%s') > 0 ", filter.ColumnId)
			case ColumnFilterTypeEmpty:
				str = fmt.Sprintf("json_length(cells, '$.%s') = 0 ", filter.ColumnId)
			case ColumnFilterTypeContains:
				str = fmt.Sprintf("json_extract(cells, '$.%s') like '%%%s%%' ", filter.ColumnId, filter.Value)
			case ColumnFilterTypeSelected:
				str = fmt.Sprintf("json_extract(cells, '$.%s') = true", filter.ColumnId)
			case ColumnFilterTypeNotSelected:
				str = fmt.Sprintf("json_extract(cells, '$.%s') = false", filter.ColumnId)
			}
			if len(str) > 0 {
				wheres = append(wheres, str)
			}
		}
		sql += " and ( " + strings.Join(wheres, view.FilterType) + ")"
	}

	if len(view.Sorts) > 0 {
		cols := make([]string, 0)
		for _, col := range view.Sorts {
			realCol := strings.TrimPrefix(col, "-")
			desc := ""
			if strings.HasPrefix(col, "-") {
				desc = " desc "
			}
			cols = append(cols, fmt.Sprintf("json_extract(cells, '$.%s') %s", realCol, desc))

		}
		sortsStr := fmt.Sprintf(" order by %s", strings.Join(cols, ", "))
		sql += sortsStr
	}

	sql += " limit ?, ?"
	err := db.SQL(sql, view.TableId, (page-1)*pagesize, pagesize).Find(&list)
	if err != nil {
		return list, err
	}
	// 如果选反了隐藏一些字段. 则直接不返回给前台了
	if len(view.HideColumns) > 0 {
		for index, item := range list {
			for _, colId := range view.HideColumns {
				delete(item.Cells, colId)
			}
			list[index] = item
		}
	}

	// 拉取一下create/update用户基本信息
	if len(list) > 0 {
		userIds := make([]string, 0)
		for _, record := range list {
			if !isExists(userIds, record.UserId) {
				userIds = append(userIds, record.UserId)
			}
			if !isExists(userIds, record.UpdatedBy) {
				userIds = append(userIds, record.UpdatedBy)
			}
		}
		fmt.Println(userIds)
		users := make([]User, 0)
		err := db.In("id", userIds).Cols("nickname", "avatar", "id").Find(&users)
		if err != nil {
			return list, err
		}
		for _, user := range users {
			for index, record := range list {
				if record.UserId == user.Id {
					list[index].UserName = user.Nickname
					list[index].UserAvatar = user.Avatar
				}
				if record.UpdatedBy == user.Id {
					list[index].UpdateUserName = user.Nickname
					list[index].UpdateUserAvatar = user.Avatar
				}
			}
		}
	}

	return list, nil
}

func isExists(list []string, item string) bool {
	for _, obj := range list {
		if obj == item {
			return true
		}
	}
	return false
}

func checkRecordColumnValues(record *Record) error {
	// 若未填写的字段, 有默认值的给填写上
	columns, err := GetTableColumns(record.TableId)
	if err != nil {
		return err
	}

	for _, col := range columns {
		// 只有前端未传入某字段时, 才做初始化
		val, ok := record.Cells[col.Id]
		if !ok {
			val = getDefaultRecordValue(&col)

		}

		//如果是数值类型, 并且传入的是字符串格式, 帮忙转成数值类型
		switch col.Type {
		case ColumnTypeNumber:
			if v, ok := val.(string); ok {
				if v == "" {
					v = "0"
				}
				v2, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return fmt.Errorf("值:%s 无效的数值类型", val)
				}
				val = v2
			}

		}
		record.Cells[col.Id] = val
	}

	return nil
}
