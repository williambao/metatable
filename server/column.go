package server

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/williambao/metatable/model"
	"github.com/williambao/metatable/router/middleware/session"
)

func GetColumn(c *gin.Context) {
	table := session.Table(c)
	c.JSON(200, table)
}

func CreateTableColumn(c *gin.Context) {

	in := &model.Column{}
	err := c.Bind(in)
	if err != nil {
		abort(c, "无效数据")
		return
	}

	user := session.User(c)
	in.UserId = user.Id
	in.TableId = c.Param("tableId")
	in.Order = 99

	err = model.CreateColumn(in)
	if err != nil {
		abort(c, err.Error())
		return
	}

	c.JSON(200, in)
}

func UpdateColumn(c *gin.Context) {
	in := map[string]interface{}{}
	err := c.BindJSON(&in)
	if err != nil {
		abort(c, "无效数据")
		return
	}

	obj := &model.Column{}
	b, err := json.Marshal(in)
	err = json.Unmarshal(b, obj)
	if err != nil {
		abort(c, "无效的Json数据格式")
		return
	}

	user := session.User(c)

	obj.TableId = c.Param("tableId")
	obj.Id = c.Param("columnId")
	obj.UserId = user.Id

	cols := make([]string, 0)
	for k, _ := range in {
		cols = append(cols, k)
	}

	err = model.UpdateColumn(obj, cols...)
	if err != nil {
		abort(c, err.Error())
		return
	}

	success(c)
}

func DeleteColumn(c *gin.Context) {
	columnId := c.Param("columnId")
	user := session.User(c)
	table := session.Table(c)
	err := model.DeleteColumn(columnId, user, table)
	if err != nil {
		abort(c, err.Error())
		return
	}

	success(c)
}

type columnOrders struct {
	Ids []string `json:"ids"`
}

func UpdateColumnOrders(c *gin.Context) {
	in := new(columnOrders)
	err := c.BindJSON(&in)
	if err != nil {
		abort(c, "无效数据")
		return
	}
	table := session.Table(c)
	err = model.UpdateColumnOrders(table, in.Ids)
	if err != nil {
		abort(c, err.Error())
		return
	}

	success(c)
}

// 针对下拉框字段, 增加下拉项
// func CreateColumnOption(c *gin.Context) {
// 	in := &model.ColumnOption{}
// 	err := c.Bind(in)
// 	if err != nil {
// 		abort(c, "无效数据")
// 		return
// 	}

// 	user := session.User(c)
// 	in.TableId = c.Param("tableId")
// 	in.ColumnId = c.Param("columnId")
// 	in.UserId = user.Id

// 	err = model.CreateColumnOption(in)
// 	if err != nil {
// 		abort(c, err.Error())
// 		return
// 	}

// 	c.JSON(200, in)
// }

// func UpdateColumnOption(c *gin.Context) {
// 	in := &model.ColumnOption{}
// 	err := c.Bind(in)
// 	if err != nil {
// 		abort(c, "无效数据")
// 		return
// 	}

// 	user := session.User(c)
// 	in.TableId = c.Param("tableId")
// 	in.ColumnId = c.Param("columnId")
// 	in.UserId = user.Id
// 	in.Id = c.Param("optionId")

// 	cols := make([]string, 0)
// 	for k, _ := range c.Request.PostForm {
// 		cols = append(cols, k)
// 	}

// 	fmt.Printf("cols:%v", cols)

// 	err = model.UpdateColumnOption(in, cols...)
// 	if err != nil {
// 		abort(c, err.Error())
// 		return
// 	}

// 	c.JSON(200, in)
// }
