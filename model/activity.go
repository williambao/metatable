package model

type ActivityAction string

const (
	ActivityActionCreate ActivityAction = "create"
	ActivityActionUpdate ActivityAction = "update"
	ActivityActionDelete ActivityAction = "delete"
	ActivityActionMove   ActivityAction = "move"
)

func (a ActivityAction) Message() string {
	switch a {
	case ActivityActionCreate:
		return "创建了"
	case ActivityActionUpdate:
		return "更新了"
	case ActivityActionDelete:
		return "删除了"
	case ActivityActionMove:
		return "移动了"
	}
	return ""
}

type Activity struct {
	BaseModel  `xorm:"extends"`
	TableId    string         `json:"table_id" form:"table_id" xorm:"varchar(20)"`
	UserId     string         `json:"user_id" form:"user_id" xorm:"varchar(20)"`
	Action     ActivityAction `json:"action" xorm:"varchar(20)"`
	Content    string         `json:"content"`
	TargetType string         `json:"target_type"`
}
