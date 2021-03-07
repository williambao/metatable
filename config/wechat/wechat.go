package wechat

import "context"

const key = "wechat"

type Setter interface {
	Set(string, interface{})
}
type WechatConfig struct {
	AppId     string
	AppSecret string
}

func FromContext(c context.Context) *WechatConfig {
	return c.Value(key).(*WechatConfig)
}

func ToContext(c Setter, config *WechatConfig) {
	c.Set(key, config)
}
