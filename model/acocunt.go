package model

type Login struct {
	UserName string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
