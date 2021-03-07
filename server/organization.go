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

func GetUserOrganizations(c *gin.Context) {
	var (
		user = session.User(c)
	)

	organizations, err := model.GetUserOrganizations(user)
	if err != nil {
		abort(c, err.Error())
		return
	}

	c.JSON(200, organizations)

}

func CreateOrganization(c *gin.Context) {
	in := &model.Organization{}
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
	in.MemberCount = 1
	in.MemberLimit = 10
	in.TableLimit = 10

	err = model.CreateOrgnization(in)
	if err != nil {
		logrus.Errorf("create organization failed: %s", err.Error())
		// abort(c, "创建表格出错!")
		abort(c, err.Error())
		return
	}

	c.JSON(200, in)
}

func UpdateOrganization(c *gin.Context) {
	in := map[string]interface{}{}
	err := c.BindJSON(&in)
	if err != nil {
		abort(c, "无效数据")
		return
	}

	obj := &model.Organization{}
	b, err := json.Marshal(in)
	err = json.Unmarshal(b, obj)
	if err != nil {
		abort(c, "无效的Json数据格式")
		return
	}

	organizationId := c.Param("organizationId")

	user := session.User(c)
	oUser, err := model.GetOrganizationUser(organizationId, user.Id)
	if err != nil {
		abort(c, err.Error())
		return
	}

	if !oUser.IsAdmin && !oUser.IsOwner {
		abort(c, "您无权限修改")
		return
	}

	obj.Id = organizationId
	obj.UserId = user.Id

	cols := make([]string, 0)
	for k, _ := range in {
		cols = append(cols, k)
	}

	err = model.UpdateOrganization(obj, cols...)
	if err != nil {
		abort(c, err.Error())
		return
	}

	success(c)
}

func DeleteOrganization(c *gin.Context) {
	id := c.Param("organizationId")
	user := session.User(c)

	err := model.DeleteOrganization(id, user)
	if err != nil {
		abort(c, err.Error())
		return
	}
	success(c)
}
