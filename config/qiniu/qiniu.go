package qiniu

import "context"

const key = "qiniu"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}
type QiniuConfig struct {
	AccessKey   string
	SecretKey   string
	Bucket      string
	CallbackURL string
	ExpiresIn   int64
	UploadURL   string
	URL         string
}

func FromContext(c context.Context) *QiniuConfig {
	return c.Value(key).(*QiniuConfig)
}

func ToContext(c Setter, config *QiniuConfig) {
	c.Set(key, config)
}
