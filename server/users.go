package server

import (
	"github.com/gin-gonic/gin"
	"github.com/williambao/metatable/model"
	"github.com/williambao/metatable/router/middleware/session"
)

func GetUsers(c *gin.Context) {

}

func GetUser(c *gin.Context) {
	user, _ := model.GetUserById(c.Param("id"))
	if user == nil {
		abort(c, "未找到用户")
		return
	}

	if !user.IsActive {
		abort(c, "用户帐号被冻结")
		return
	}

	c.JSON(200, user)
}

func UpdateUser(c *gin.Context) {
	in := &model.User{}
	err := c.Bind(in)
	if err != nil {
		abort(c, "无效数据")
		return
	}

	user, _ := model.GetUserById(c.Param("id"))
	if user == nil {
		abort(c, "未找到用户")
		return
	}

	sessionUser := session.User(c)
	if !sessionUser.IsAdmin && sessionUser.Id != c.Param("id") {
		abort(c, "无权限操作")
		return
	}

	cols := make([]string, 0)
	for k, _ := range c.Request.PostForm {
		cols = append(cols, k)
	}

	err = model.UpdateById(user.Id, in, cols...)
	if err != nil {
		abort(c, err.Error())
		return
	}

	success(c)
}
