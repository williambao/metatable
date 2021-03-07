package model

type Role struct {
	BaseModel `xorm:"extends"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	IsActive  bool   `json:"is_active"`
}
