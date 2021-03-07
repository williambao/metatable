package session

import (
	"github.com/gin-gonic/gin"
	"github.com/williambao/metatable/model"
	"github.com/williambao/metatable/utils"
)

func Table(c *gin.Context) *model.Table {
	v, ok := c.Get("table")
	if !ok {
		return nil
	}
	r, ok := v.(*model.Table)
	if !ok {
		return nil
	}
	return r
}

func SetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			tableId = c.Param("tableId")
			user    = User(c)
		)

		table, err := model.GetTableById(tableId, user)
		if err != nil {
			c.JSON(400, utils.NewError(err.Error()))
			c.Abort()
			return
		}

		c.Set("table", table)
		c.Next()
	}
}

func TablePermission(c *gin.Context) *model.TableUser {
	v, ok := c.Get("tableuser")
	if !ok {
		return nil
	}
	u, ok := v.(*model.TableUser)
	if !ok {
		return nil
	}
	return u
}

func SetTablePermission() gin.HandlerFunc {

	return func(c *gin.Context) {
		user := User(c)
		table := Table(c)

		if table == nil {
			c.JSON(400, utils.NewError("无效表格id"))
			c.Abort()
			return
		}

		if table.IsPrivate {
			if user == nil {
				c.JSON(401, utils.NewLoginRequiredError())
				c.Abort()
				return
			}
			perm, err := model.GetTableUser(table.Id, user.Id)
			if err != nil {
				c.JSON(400, utils.NewError(err.Error()))
				c.Abort()
				return
			}
			c.Set("tableuser", perm)
		} else {
			// 公开的table, 任何人都应有查看权限
			perm := new(model.TableUser)
			perm.ViewData = model.TableDataPermissionUser
			perm.UserId = "public"
			if user != nil {
				perm2, _ := model.GetTableUser(table.Id, user.Id)
				if perm2 != nil {
					perm = perm2
				}
			}

			c.Set("tableuser", perm)
		}

		c.Next()
	}
}
