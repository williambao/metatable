package router

import (
	"github.com/gin-gonic/gin"
	"github.com/williambao/metatable/config"
)

const configKey = "config"

func SetConfig(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		config.ToContext(c, cfg)
		c.Next()
	}
}
