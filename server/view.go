package server

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/williambao/metatable/model"
	"github.com/williambao/metatable/router/middleware/session"
)

func GetTableView(c *gin.Context) {

}

func CreateTableView(c *gin.Context) {
	in := &model.TableView{}
	err := c.Bind(in)
	if err != nil {
		abort(c, "无效数据")
		return
	}

	table := session.Table(c)
	user := session.User(c)
	in.TableId = table.Id
	in.UserId = user.Id
	err = model.CreateTableView(in)
	if err != nil {
		abort(c, err.Error())
		return
	}

	c.JSON(200, in)
}

func UpdateTableView(c *gin.Context) {
	in := map[string]interface{}{}
	err := c.BindJSON(&in)
	if err != nil {
		abort(c, "无效数据")
		return
	}

	obj := &model.TableView{}
	b, err := json.Marshal(in)
	err = json.Unmarshal(b, obj)
	if err != nil {
		abort(c, "无效的Json数据格式")
		return
	}

	// user := session.User(c)
	id := c.Param("viewId")

	obj.Id = id
	cols := make([]string, 0)
	for k, _ := range in {
		cols = append(cols, k)
	}

	err = model.UpdateTableView(obj, cols...)
	if err != nil {
		logrus.Errorf("update table view failed : %s\n", err.Error())
		abort(c, "更新失败")
		return
	}

	view, err := model.GetTableView(id)
	if err != nil {
		abort(c, err.Error())
		return
	}

	c.JSON(200, view)
}

// 按view所设的条件来查询
func GetViewRecords(c *gin.Context) {
	page := queryInt(c, "page", 1)
	pagesize := queryInt(c, "pagesize", 100)
	query := c.Query("query")

	id := c.Param("viewId")

	view, err := model.GetTableView(id)
	if err != nil {
		abort(c, "查询失败")
		return
	}
	perm := session.TablePermission(c)

	user := session.User(c)

	list, err := model.GetViewRecords(view, perm, user, query, page, pagesize)
	if err != nil {
		logrus.Errorf("get view records failed : %s\n", err.Error())
		abort(c, "查询失败")
		return
	}

	c.JSON(200, list)
}
