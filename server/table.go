package server

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/williambao/metatable/model"
	"github.com/williambao/metatable/router/middleware/session"
	"github.com/williambao/metatable/utils"
)

func GetTable(c *gin.Context) {
	table := session.Table(c)
	if table != nil && table.Permission == nil || table.Permission.Id == "" {
		table.Permission = session.TablePermission(c)
	}
	c.JSON(200, table)
}

func CreateTable(c *gin.Context) {
	in := &model.Table{}
	err := c.Bind(in)
	if err != nil {
		abort(c, "无效数据")
		return
	}

	user := session.User(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, utils.NewLoginRequiredError())
		c.Abort()
		return
	}

	in.UserId = user.Id

	err = model.CreateTable(in)
	if err != nil {
		logrus.Errorf("创建表格出错: %s", err.Error())
		// abort(c, "创建表格出错!")
		abort(c, err.Error())
		return
	}

	c.JSON(200, in)
}

type updateTable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
}

func UpdateTable(c *gin.Context) {
	in := map[string]interface{}{}
	err := c.BindJSON(&in)
	if err != nil {
		abort(c, "无效数据")
		logrus.Debugf("update table data err: " + err.Error())
		return
	}

	table := session.Table(c)

	obj := &model.Table{}
	b, err := json.Marshal(in)
	err = json.Unmarshal(b, obj)
	if err != nil {
		abort(c, "无效的Json数据格式")
		return
	}
	obj.Id = table.Id

	cols := make([]string, 0)
	for k, _ := range in {
		cols = append(cols, k)
	}

	logrus.Debugf("update table cols: %v", cols)
	logrus.Debugf("update table cover: %v", obj.CoverURL)

	err = model.UpdateTable(obj, cols...)
	if err != nil {
		abort(c, err.Error())
		return
	}
	success(c)
}

func GetUserTables(c *gin.Context) {
	var (
		user     = session.User(c)
		page     = queryInt(c, "page", 1)
		pagesize = queryInt(c, "pagesize", 30)
		sort     = c.Query("sort")
	)

	tables, err := model.GetUserTables(user, page, pagesize, sort)
	if err != nil {
		abort(c, err.Error())
		return
	}

	c.JSON(200, tables)

}

// 获取表格所有列信息
func GetTableColumns(c *gin.Context) {
	// user := session.User(c)

	list, err := model.GetTableColumns(c.Param("tableId"))
	if err != nil {
		abort(c, err.Error())
		return
	}

	c.JSON(200, list)
}

func DeleteTable(c *gin.Context) {
	id := c.Param("tableId")
	user := session.User(c)

	err := model.DeleteTable(id, user)
	if err != nil {
		abort(c, err.Error())
		return
	}
	success(c)
}
