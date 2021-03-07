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

func GetTemplates(c *gin.Context) {
	page := queryInt(c, "page", 1)
	pagesize := queryInt(c, "pagesize", 100)
	name := c.Query("name")
	isActive := c.Query("is_active")

	list, err := model.GetTemplates(name, isActive == "1" || isActive == "true", false, page, pagesize)
	if err != nil {
		abort(c, err.Error())
		return
	}

	c.JSON(200, list)
}

func GetTemplate(c *gin.Context) {
	id := c.Param("templateId")
	obj, err := model.GetTemplateById(id)
	if err != nil {
		abort(c, err.Error())
		return
	}
	c.JSON(200, obj)
}

func CreateTemplate(c *gin.Context) {
	in := &model.Template{}
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

	if !user.InRole("offical") && !user.IsAdmin {
		abort(c, "无权限")
		return
	}

	in.UserId = user.Id

	err = model.CreateTemplate(in)
	if err != nil {
		logrus.Errorf("创建表格出错: %s", err.Error())
		abort(c, "创建表格出错!")
		return
	}

	c.JSON(200, in)
}

type updateTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	IsActive    bool   `json:"is_active"`
}

func UpdateTemplate(c *gin.Context) {
	templateId := c.Param("templateId")
	in := map[string]interface{}{}
	err := c.BindJSON(&in)
	if err != nil {
		abort(c, "无效数据")
		logrus.Debugf("update template data err: " + err.Error())
		return
	}

	user := session.User(c)

	template, err := model.GetTemplateById(templateId)
	if err != nil {
		abort(c, err.Error())
		return
	}

	if !user.IsAdmin && template.UserId != user.Id {
		abort(c, "无权限")
		return
	}

	obj := &model.Template{}
	b, err := json.Marshal(in)
	err = json.Unmarshal(b, obj)
	if err != nil {
		abort(c, "无效的Json数据格式")
		return
	}
	obj.Id = template.Id

	cols := make([]string, 0)
	for k, _ := range in {
		cols = append(cols, k)
	}

	err = model.UpdateTemplate(obj, cols...)
	if err != nil {
		abort(c, err.Error())
		return
	}
	success(c)
}

func DeleteTemplate(c *gin.Context) {
	id := c.Param("templateId")
	user := session.User(c)

	err := model.DeleteTemplate(id, user)
	if err != nil {
		abort(c, err.Error())
		return
	}
	success(c)
}

// 获取表格所有列信息
func GetTemplateColumns(c *gin.Context) {
	// user := session.User(c)
	templateId := c.Param("templateId")
	list, err := model.GetTemplateColumns(templateId)
	if err != nil {
		abort(c, err.Error())
		return
	}

	c.JSON(200, list)
}

func CreateTemplateColumn(c *gin.Context) {

	in := &model.TemplateColumn{}
	err := c.Bind(in)
	if err != nil {
		abort(c, "无效数据")
		return
	}

	user := session.User(c)
	in.UserId = user.Id
	in.TemplateId = c.Param("templateId")
	in.Order = 99

	err = model.CreateTemplateColumn(in)
	if err != nil {
		abort(c, err.Error())
		return
	}

	c.JSON(200, in)
}

func UpdateTemplateColumn(c *gin.Context) {
	in := map[string]interface{}{}
	err := c.BindJSON(&in)
	if err != nil {
		abort(c, "无效数据")
		return
	}

	obj := &model.TemplateColumn{}
	b, err := json.Marshal(in)
	err = json.Unmarshal(b, obj)
	if err != nil {
		abort(c, "无效的Json数据格式")
		return
	}

	user := session.User(c)

	obj.TemplateId = c.Param("templateId")
	obj.Id = c.Param("columnId")
	obj.UserId = user.Id

	cols := make([]string, 0)
	for k, _ := range in {
		cols = append(cols, k)
	}

	err = model.UpdateTemplateColumn(obj, cols...)
	if err != nil {
		abort(c, err.Error())
		return
	}

	success(c)
}

func DeleteTemplateColumn(c *gin.Context) {
	templateId := c.Param("templateId")
	columnId := c.Param("columnId")
	user := session.User(c)
	template, _ := model.GetTemplateById(templateId)
	err := model.DeleteTemplateColumn(columnId, user, template)
	if err != nil {
		abort(c, err.Error())
		return
	}

	success(c)
}

func UpdateTemplateColumnOrders(c *gin.Context) {
	templateId := c.Param("templateId")
	in := new(columnOrders)
	err := c.BindJSON(&in)
	if err != nil {
		abort(c, "无效数据")
		return
	}

	template, _ := model.GetTemplateById(templateId)
	err = model.UpdateTemplateColumnOrders(template, in.Ids)
	if err != nil {
		abort(c, err.Error())
		return
	}

	success(c)
}
