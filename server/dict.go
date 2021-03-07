package server

import (
	"github.com/gin-gonic/gin"
	"github.com/williambao/metatable/model"
)

func GetConditions(c *gin.Context) {
	result := model.GetColumnTypeFilters()
	c.JSON(200, result)
}
