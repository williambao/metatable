package header

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NoCache is a middleware function that appends headers
// to prevent the client from caching the HTTP response.
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}

// Options is a middleware function that appends headers
// for options requests and aborts then exits the middleware
// chain and ends the request.
func Options(c *gin.Context) {
	// if c.Request.Method != "OPTIONS" {
	// 	c.Next()
	// } else {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
	c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Content-Type", "application/json")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	} else {
		c.Next()
	}

	// }
}

// Version is a middleware function that appends the Drone
// version information to the HTTP response. This is intended
// for debugging and troubleshooting.
// func Version(c *gin.Context) {
// 	c.Header("X-TableXX-VERSION", version.Version)
// 	c.Next()
// }
