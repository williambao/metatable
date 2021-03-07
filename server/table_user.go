package server

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/williambao/metatable/model"
	"github.com/williambao/metatable/router/middleware/session"
)

func GetTableUsers(c *gin.Context) {
	table := session.Table(c)
	list, err := model.GetTableUsers(table.Id)
	if err != nil {
		abort(c, err.Error())
		return
	}

	c.JSON(200, list)
}

type addTableUser struct {
	Username string `json:"username"`
}

func AddTableUser(c *gin.Context) {
	in := &addTableUser{}
	err := c.Bind(in)
	if err != nil {
		abort(c, "数据格式不对")
		return
	}

	user := model.GetUserByUserName(in.Username)
	if user == nil {
		abort(c, "未找到用户")
		return
	}

	table := session.Table(c)
	suser := session.User(c)

	tableUser := new(model.TableUser)
	tableUser.UserId = user.Id
	tableUser.CreatedBy = suser.Id
	tableUser.TableId = table.Id

	err = model.CreateTableUser(tableUser)
	if err != nil {
		abort(c, err.Error())
		return
	}

	c.JSON(200, in)
}

func UpdateTableUser(c *gin.Context) {
	in := map[string]interface{}{}
	err := c.BindJSON(&in)
	if err != nil {
		abort(c, "无效数据")
		return
	}

	obj := &model.TableUser{}
	b, err := json.Marshal(in)
	err = json.Unmarshal(b, obj)
	if err != nil {
		abort(c, "无效的Json数据格式")
		return
	}

	table := session.Table(c)

	obj.TableId = table.Id
	obj.Id = c.Param("tableUserId")
	cols := make([]string, 0)
	for k, _ := range in {
		if k == "user_id" {
			continue
		}
		cols = append(cols, k)
	}

	err = model.UpdateTableUser(obj, cols...)
	if err != nil {
		abort(c, err.Error())
		return
	}

	success(c)
}
func DeleteTableUser(c *gin.Context) {
	err := model.DeleteTableUser(c.Param("tableUserId"))
	if err != nil {
		abort(c, err.Error())
		return
	}

	success(c)
}
