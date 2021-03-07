package server

import (
	"github.com/gin-gonic/gin"
	"github.com/williambao/metatable/model"
	"github.com/williambao/metatable/router/middleware/session"
)

func GetTableSubscribes(c *gin.Context) {
	table := session.Table(c)
	list, err := model.GetTableSubscribes(table.Id)
	if err != nil {
		abort(c, err.Error())
		return
	}

	c.JSON(200, list)
}

func AddTableSubscribe(c *gin.Context) {

	table := session.Table(c)
	suser := session.User(c)

	tableSubscribe := new(model.TableSubscribe)
	tableSubscribe.UserId = suser.Id
	tableSubscribe.CreatedBy = suser.Id
	tableSubscribe.TableId = table.Id

	err := model.CreateTableSubscribe(tableSubscribe)
	if err != nil {
		abort(c, err.Error())
		return
	}

	c.JSON(200, tableSubscribe)
}

func DeleteTableSubscribe(c *gin.Context) {
	err := model.DeleteTableSubscribe(c.Param("tableUserId"))
	if err != nil {
		abort(c, err.Error())
		return
	}

	success(c)
}
