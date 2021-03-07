package server

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/williambao/metatable/model"
	"github.com/williambao/metatable/router/middleware/session"
)

func GetRecords(c *gin.Context) {
	table := session.Table(c)
	// user := session.User(c)

	page := queryInt(c, "page", 1)
	pagesize := queryInt(c, "pagesize", 100)
	sort := c.Query("sort")
	query := c.Query("query")

	view := &model.TableView{}
	view.TableId = table.Id

	if len(sort) > 0 {
		sortArr := strings.Split(sort, ",")
		view.Sorts = sortArr
	}

	user := session.User(c)
	perm := session.TablePermission(c)

	list, err := model.GetViewRecords(view, perm, user, query, page, pagesize)
	if err != nil {
		abort(c, err.Error())
		return
	}

	// maps := make([]map[string]interface{}, 0)
	// for index, item := range list {
	// 	cells, _ := model.GetRecordCells(item.Id)

	// 	for _, item2 := range cells {
	// 		maps = append(maps, item2.ToMap())
	// 	}
	// 	list[index].Cells = maps
	// }

	c.JSON(200, list)
}

func CreateRecord(c *gin.Context) {
	in := &model.Record{}
	err := c.Bind(in)
	if err != nil {
		logrus.Errorf("create record invalid data: %s", err.Error())
		abort(c, "无效的数据")
		return
	}

	table := session.Table(c)
	user := session.User(c)
	in.UserId = user.Id
	in.UpdatedBy = user.Id
	in.TableId = table.Id
	err = model.CreateRecord(in)
	if err != nil {
		abort(c, err.Error())
		return
	}
	in.UserName = user.Nickname
	in.UserAvatar = user.Avatar

	c.JSON(200, in)

}

func DeleteRecord(c *gin.Context) {
	// id := c.Param("recordId")
	abort(c, "此接口不可用")
	return
	// err := model.DeleteRecord(id)
	// if err != nil {
	// 	abort(c, err.Error())
	// 	return
	// }
	// success(c)
}

func DeleteRecordByBatch(c *gin.Context) {
	recordIds := c.Query("ids")
	if len(recordIds) == 0 {
		abort(c, "请传入ids")
		return
	}
	ids := strings.Split(recordIds, ",")
	perm := session.TablePermission(c)
	user := session.User(c)

	err := model.DeleteRecordByBatch(ids, perm, user)
	if err != nil {
		abort(c, err.Error())
		return
	}
	success(c)
}

// 修改列
func EditReocrdColumn(c *gin.Context) {
	in := &model.Record{}
	err := c.Bind(in)
	if err != nil {
		abort(c, "无效的数据")
		return
	}

	id := c.Param("recordId")

	err = model.UpdateById(id, in, "cells")
	if err != nil {
		abort(c, err.Error())
		return
	}

	// id := c.Param("recordId")
	c.JSON(200, in)
}

// 修改某一个数据字段
func PatchRecordColumn(c *gin.Context) {
	in := &model.Record{}
	err := c.Bind(in)
	if err != nil {
		abort(c, "无效的数据")
		return
	}

	if in.Cells == nil {
		abort(c, "请传入正确的cells字段内容")
		return
	}
	user := session.User(c)

	id := c.Param("recordId")

	record, err := model.GetRecord(id)
	if err != nil {
		abort(c, err.Error())
		return
	}

	record.UpdatedBy = user.Id

	for k, v := range in.Cells {
		record.Cells[k] = v
	}

	err = model.UpdateRecordCells(record)
	if err != nil {
		abort(c, err.Error())
		return
	}

	success(c)
}
