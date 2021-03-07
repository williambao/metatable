package session

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/williambao/metatable/model"
	"github.com/williambao/metatable/shared/token"
	"github.com/williambao/metatable/utils"
)

func User(c *gin.Context) *model.User {
	v, ok := c.Get("user")
	if !ok {
		return nil
	}
	u, ok := v.(*model.User)
	if !ok {
		return nil
	}
	return u
}

func Token(c *gin.Context) *token.Token {
	v, ok := c.Get("token")
	if !ok {
		return nil
	}
	u, ok := v.(*token.Token)
	if !ok {
		return nil
	}
	return u
}

func SetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *model.User

		t, err := token.ParseRequest(c.Request, func(t *token.Token) (string, error) {

			var err error
			user, err = model.GetUserById(t.Text)
			if err != nil {
				return "", err
			}
			return user.Hash, nil
		})
		// fmt.Printf("token: %v, error: %v", t, err)
		// log.Debugf("get token %v, error: %v", t.Text, err)
		if err == nil {
			// confv := c.MustGet("config")
			// if conf, ok := confv.(*model.Config); ok {
			// 	user.Admin = conf.IsAdmin(user)
			// }
			c.Set("user", user)

			// if this is a session token (ie not the API token)
			// this means the user is accessing with a web browser,
			// so we should implement CSRF protection measures.
			if t.Kind == token.SessToken {
				err = token.CheckCsrf(c.Request, func(t *token.Token) (string, error) {
					return user.Hash, nil
				})
				// if csrf token validation fails, exit immediately
				// with a not authorized error.
				if err != nil {
					c.JSON(http.StatusUnauthorized, utils.NewNoAccessPermissionError(""))
					return
				}
			}
		}
		c.Next()
	}
}
func MustAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := User(c)
		switch {
		case user == nil:
			c.JSON(http.StatusUnauthorized, utils.NewLoginRequiredError())
			c.Abort()
		case user.IsAdmin == false:
			c.JSON(http.StatusUnauthorized, utils.NewNoAccessPermissionError(""))
			c.Abort()
		default:
			c.Next()
		}
	}
}

func MustUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := User(c)
		switch {
		case user == nil:
			c.JSON(http.StatusUnauthorized, utils.NewLoginRequiredError())
			c.Abort()
		default:
			c.Next()
		}
	}
}

// 当是POST,PUT,DELETE等非查询操作时, 检查用户状态,
// 冻结的用户不允许做任何操作
func MustActiveUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := User(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, utils.NewLoginRequiredError())
			c.Abort()
			return
		}
		if strings.ToUpper(c.Request.Method) != "GET" && !user.IsActive {
			c.JSON(400, utils.NewError("帐号被冻结, 无法操作"))
			c.Abort()
			return

		}
		c.Next()
	}
}

func MustTableViewDataPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		table := Table(c)
		if table.IsPrivate {
			user := User(c)
			perm := TablePermission(c)
			if user == nil {
				c.JSON(http.StatusUnauthorized, utils.NewLoginRequiredError())
				c.Abort()
				return
			}
			if !user.IsAdmin && !perm.IsTableAdmin && perm.ViewData == model.TableDataPermissionNone {
				c.JSON(400, utils.NewError("您无权限查看数据"))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
func MustTableInsertDataPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := User(c)
		perm := TablePermission(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, utils.NewLoginRequiredError())
			c.Abort()
			return
		}
		if !user.IsAdmin && !perm.IsTableAdmin && perm.EditData == model.TableDataPermissionNone {
			c.JSON(400, utils.NewError("您无权限操作"))
			c.Abort()
			return
		}

		c.Next()
	}
}
func MustTableEditDataPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := User(c)
		perm := TablePermission(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, utils.NewLoginRequiredError())
			c.Abort()
			return
		}
		if !user.IsAdmin && !perm.IsTableAdmin && perm.EditData == model.TableDataPermissionNone {
			c.JSON(400, utils.NewError("您无权限操作"))
			c.Abort()
			return
		}

		c.Next()
	}
}
func MustTableDeleteDataPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := User(c)
		perm := TablePermission(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, utils.NewLoginRequiredError())
			c.Abort()
			return
		}
		if !user.IsAdmin && !perm.IsTableAdmin && perm.EditData == model.TableDataPermissionNone {
			c.JSON(400, utils.NewError("您无权限操作"))
			c.Abort()
			return
		}

		c.Next()
	}
}
func MustTableAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := User(c)
		perm := TablePermission(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, utils.NewLoginRequiredError())
			c.Abort()
			return
		}
		if perm == nil {
			c.JSON(400, utils.NewError("无表格权限"))
			c.Abort()
			return
		}
		if !user.IsAdmin && !perm.IsTableAdmin {
			c.JSON(400, utils.NewError("您无权限操作"))
			c.Abort()
			return
		}

		c.Next()
	}
}
