package aliyun

import "context"

const key = "aliyun"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}
type AliyunConfig struct {
	AccessKey    string
	SecretKey    string
	SignName     string
	TemplateCode string
}

func FromContext(c context.Context) *AliyunConfig {
	return c.Value(key).(*AliyunConfig)
}

func ToContext(c Setter, config *AliyunConfig) {
	c.Set(key, config)
}
