package utils

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *Error) Error() string {
	buf, _ := json.Marshal(e)
	return string(buf)
}

func NewError(msg string, v ...interface{}) error {
	return NewErrorWithCode(10, msg, v...)
}

func NewErrorWithCode(code int, msg string, v ...interface{}) error {
	text := fmt.Sprintf(msg, v...)
	return &Error{Code: code, Message: text}
}
func NewInternelError() error {
	return NewErrorWithCode(500, "内部错误")
}

func NewParamRequiredError(name string) error {
	return NewError("参数不能为空 - %s", name)
}

func NewNotFoundError() error {
	return NewErrorWithCode(404, "未找到")
}

func NewSensitiveWordError() error {
	return NewError("您录入的信息无效，当前无法进行操作，请修改后再次提交！")
}

func NewLoginRequiredError() error {
	return NewErrorWithCode(401, "请先登录")
}
func NewNoAccessPermissionError(msg string) error {
	if msg == "" {
		msg = "no  permissions"
	}
	return NewErrorWithCode(403, msg)
}
