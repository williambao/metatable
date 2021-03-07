package model

import "errors"
import "github.com/rs/xid"

type TableDataPermission string

const (
	TableDataPermissionAll  TableDataPermission = "all"  // 全部数据权限
	TableDataPermissionTeam                     = "team" // 本组权限
	TableDataPermissionUser                     = "user" // 本人增加的数据权限
	TableDataPermissionNone                     = "none" // 无权限
)

type TableUser struct {
	BaseModel  `xorm:"extends"`
	TableId    string              `json:"-" xorm:"UNIQUE(uk_table_user) NOT NULL"`
	UserId     string              `json:"user_id" xorm:"UNIQUE(uk_table_user)  NOT NULL"`
	Nickname   string              `json:"nickname" xorm:"-"`
	Username   string              `json:"username" xorm:"-"`
	UserAvatar string              `json:"user_avatar" xorm:"-"`
	Remark     string              `json:"remark"`
	ViewData   TableDataPermission `json:"view_data"`
	EditData   TableDataPermission `json:"edit_data"`
	// IsAllowInsertData bool                `json:"is_allow_insert_data"`
	// DeleteData        TableDataPermission `json:"delete_data"`

	IsTableAdmin bool `json:"is_table_admin"`
}

func CreateTableUser(item *TableUser) error {
	// 检查下是否已存在
	cnt, _ := db.Where("table_id = ? and user_id = ?", item.TableId, item.UserId).Count(&TableUser{})
	if cnt > 0 {
		return errors.New("用户已存在, 无法重复增加")
	}

	item.Id = xid.New().String()
	if item.ViewData == "" {
		item.ViewData = TableDataPermissionNone
	}
	if item.EditData == "" {
		item.EditData = TableDataPermissionNone
	}
	// if item.DeleteData == "" {
	// 	item.DeleteData = TableDataPermissionNone
	// }

	err := Insert(item)
	return err
}

func GetTableUsers(tableId string) ([]TableUser, error) {
	list := make([]TableUser, 0)
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

func GetTableUser(tableId, userId string) (*TableUser, error) {
	var table TableUser
	exist, err := db.Where("table_id = ? and user_id = ?", tableId, userId).Get(&table)
	if err != nil || !exist {
		return nil, err
	}

	return &table, nil
}

func UpdateTableUser(user *TableUser, cols ...string) error {
	err := UpdateById(user.Id, user, cols...)
	return err
}
func DeleteTableUser(id string) error {
	var table TableUser
	err := GetById(id, &table)
	if err != nil {
		return err
	}

	if table.IsTableAdmin {
		return errors.New("无法删除表格管理员")
	}

	err = DeleteById(id, &table)
	return err
}
