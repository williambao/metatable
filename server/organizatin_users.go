package server

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/williambao/metatable/model"
	"github.com/williambao/metatable/router/middleware/session"
)

func GetOrganizationUsers(c *gin.Context) {

	organization, err := model.GetOrganizationById(c.Param("organizationId"))
	if err != nil {
		abort(c, err.Error())
		return
	}

	list, err := model.GetOrganizationUsers(organization.Id)
	if err != nil {
		abort(c, err.Error())
		return
	}

	c.JSON(200, list)
}

type addOrginazationUser struct {
	Email    string `json:"email"`
	TeamID   string `json:"team_id"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Sex      int    `json:"sex"`
}

func AddOrganizationUser(c *gin.Context) {
	in := &addOrginazationUser{}
	err := c.Bind(in)
	if err != nil {
		abort(c, "数据格式不对")
		return
	}

	suser := session.User(c)

	organization, err := model.GetOrganizationById(c.Param("organizationId"))
	if err != nil {
		abort(c, err.Error())
		return
	}

	member := new(model.User)
	member.OrganizationId = organization.Id
	member.Username = in.Email
	member.Email = in.Email
	member.Nickname = in.Nickname
	member.Password = in.Password
	member.TeamId = in.TeamID
	member.Sex = in.Sex
	member.IsActive = true
	member.CreatedBy = suser.Id
	member.UpdatedBy = suser.Id

	err = model.CreateUser(member)
	if err != nil {
		abort(c, err.Error())
		return
	}

	organizationUser := new(model.OrganizationUser)
	organizationUser.UserId = member.Id
	organizationUser.CreatedBy = suser.Id
	organizationUser.OrganizationId = organization.Id

	err = model.CreateOrganizationUser(organizationUser)
	if err != nil {
		abort(c, err.Error())
		return
	}

	c.JSON(200, in)
}

func UpdateOrganizationUser(c *gin.Context) {
	in := map[string]interface{}{}
	err := c.BindJSON(&in)
	if err != nil {
		abort(c, "无效数据")
		return
	}

	obj := &model.OrganizationUser{}
	b, err := json.Marshal(in)
	err = json.Unmarshal(b, obj)
	if err != nil {
		abort(c, "无效的Json数据格式")
		return
	}

	organization, err := model.GetOrganizationById(c.Param("organizationId"))
	if err != nil {
		abort(c, err.Error())
		return
	}

	obj.OrganizationId = organization.Id
	obj.Id = c.Param("organizationUserId")
	cols := make([]string, 0)
	for k, _ := range in {
		if k == "user_id" {
			continue
		}
		cols = append(cols, k)
	}

	err = model.UpdateOrganizationUser(obj, cols...)
	if err != nil {
		abort(c, err.Error())
		return
	}

	success(c)
}
func DeleteOrganizationUser(c *gin.Context) {
	err := model.DeleteOrganizationUser(c.Param("organizationUserId"))
	if err != nil {
		abort(c, err.Error())
		return
	}

	success(c)
}
