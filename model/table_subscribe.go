package model

import (
	"errors"

	"github.com/rs/xid"
)

type TableSubscribe struct {
	BaseModel  `xorm:"extends"`
	TableId    string `json:"-" xorm:"UNIQUE(uk_table_subscribe) NOT NULL"`
	UserId     string `json:"user_id" xorm:"UNIQUE(uk_table_subscribe)  NOT NULL"`
	Nickname   string `json:"nickname" xorm:"-"`
	Username   string `json:"username" xorm:"-"`
	UserAvatar string `json:"user_avatar" xorm:"-"`
	Remark     string `json:"remark"`
	CreatedBy  string `json:"created_by,omitempty" xorm:"varchar(20)"`
}

func CreateTableSubscribe(item *TableSubscribe) error {
	// 检查下是否已存在
	cnt, _ := db.Where("table_id = ? and user_id = ?", item.TableId, item.UserId).Count(&TableSubscribe{})
	if cnt > 0 {
		return errors.New("您已经关注此表格, 无法重复关注")
	}

	item.Id = xid.New().String()

	err := Insert(item)
	return err
}

func GetTableSubscribes(tableId string) ([]TableSubscribe, error) {
	list := make([]TableSubscribe, 0)
	err := db.Where("table_id = ?", tableId).Find(&list)

	// 更新下nickname
	if len(list) > 0 {
		ids := make([]string, 0)
		for _, item := range list {
			ids = append(ids, item.UserId)
		}

		cols := []string{"id", "nickname", "username", "avatar"}
		users, _ := getUsersByIds(ids, cols)

		for _, user := range users {
			for idx, u1 := range list {
				if user.Id == u1.UserId {
					list[idx].Nickname = user.Nickname
					list[idx].Username = user.Username
					list[idx].UserAvatar = user.Avatar
				}
			}
		}
	}

	return list, err

}

func DeleteTableSubscribe(id string) error {
	var table TableSubscribe
	err := DeleteById(id, &table)
	return err
}
