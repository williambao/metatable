package server

import (
	"github.com/gin-gonic/gin"
	"github.com/williambao/metatable/router/middleware/session"
)

func GetSelf(c *gin.Context) {
	user := session.User(c)
	if user.RoleIds == nil {
		user.RoleIds = []string{}
	}
	c.JSON(200, user)
}
