package config

import (
	"github.com/gin-gonic/gin"
)

const key = "config"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

func FromContext(c *gin.Context) *Config {
	return c.Value(key).(*Config)
}

func ToContext(c Setter, cfg *Config) {
	c.Set(key, cfg)
}
