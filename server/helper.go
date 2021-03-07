package server

import (
	"strconv"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/williambao/metatable/utils"
)

func abort(c *gin.Context, message string) {
	err := utils.NewError(message)
	c.JSON(400, err)
}

// 返回简单的成功标识
// {"success": true}
func success(c *gin.Context) {
	c.JSON(200, gin.H{"success": true})
}

// func getPostKeys(c *gin.Context) ([]string, error) {
// 	list := make([]string, 0)

// 	if c.Request.Header.Get("content-type") == "application/json" {
// 		var obj map[string]interface{}
// 		err := c.Bind(&obj)
// 		if err != nil {
// 			return list, err
// 		}
// 		for k := range obj {
// 			list = append(list, k)
// 		}
// 	}

// 	return list, nil
// }

func mapToJsonString(data map[string]interface{}) (string, error) {
	b, err := json.Marshal(data)
	return string(b), err
}

func queryInt(c *gin.Context, name string, def ...int) int {
	sv := c.Query(name)
	if sv == "" {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}

	i, _ := strconv.ParseInt(sv, 10, 32)
	return int(i)
}

func queryInt64(c *gin.Context, name string, def ...int64) int64 {
	sv := c.Query(name)
	if sv == "" {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}

	i, _ := strconv.ParseInt(sv, 10, 64)
	return i
}

func formInt64(c *gin.Context, name string, def ...int64) int64 {
	sv := c.PostForm(name)
	if sv == "" {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}

	i, _ := strconv.ParseInt(sv, 10, 32)
	return i
}

// 判断某数组是否包含某字符值
func StringIn(array []string, item string) bool {
	for _, str := range array {
		if item == str {
			return true
		}
	}
	return false
}
