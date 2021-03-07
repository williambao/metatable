package model

import (
	"errors"

	"github.com/rs/xid"
)

type OrganizationUser struct {
	BaseModel      `xorm:"extends"`
	OrganizationId string `json:"organization_id" xorm:"UNIQUE(uk_organization_user) NOT NULL"`
	TeamId         string `json:"team_id" xorm:"varchar(20)"`
	UserId         string `json:"user_id" xorm:"UNIQUE(uk_organization_user)  NOT NULL"`
	Remark         string `json:"remark"`
	Nickname       string `json:"nickname"`
	TeamName       string `json:"team_name" xorm:"-"`
	Username       string `json:"username" xorm:"-"`
	UserNickname   string `json:"user_nickname" xorm:"-"`
	UserAvatar     string `json:"user_avatar" xorm:"-"`
	UserSex        int    `json:"user_sex" xorm:"-"`
	IsAdmin        bool   `json:"is_admin"`
	IsOwner        bool   `json:"is_owner"`
}

func CreateOrganizationUser(item *OrganizationUser) error {
	// 检查下是否已存在
	cnt, _ := db.Where("organization_id = ? and user_id = ?", item.OrganizationId, item.UserId).Count(&OrganizationUser{})
	if cnt > 0 {
		return errors.New("用户已存在, 无法重复增加")
	}

	item.Id = xid.New().String()

	err := Insert(item)
	return err
}

func GetOrganizationUsers(organizationId string) ([]OrganizationUser, error) {

	list := make([]OrganizationUser, 0)
	organization, err := GetOrganizationById(organizationId)
	if err != nil {
		return list, err
	}
	err = db.Where("organization_id = ?", organizationId).Find(&list)

	for idx, item := range list {
		for _, tt := range organization.Teams {
			if tt.Id == item.TeamId {
				list[idx].TeamName = tt.Name
				break
			}
			list[idx].TeamName = "未分组"
		}
	}

	// 更新下nickname
	if len(list) > 0 {
		ids := make([]string, 0)
		for _, item := range list {
			ids = append(ids, item.UserId)
		}

		cols := []string{"id", "nickname", "username", "avatar", "sex"}
		users, _ := getUsersByIds(ids, cols)

		for _, user := range users {
			for idx, u1 := range list {
				if user.Id == u1.UserId {
					list[idx].Nickname = user.Nickname
					list[idx].Username = user.Username
					list[idx].UserAvatar = user.Avatar
					list[idx].UserSex = user.Sex
				}
			}
		}
	}

	return list, err

}

func GetOrganizationUser(organizationId, userId string) (*OrganizationUser, error) {
	var table OrganizationUser
	exist, err := db.Where("organization_id = ? and user_id = ?", organizationId, userId).Get(&table)
	if err != nil || !exist {
		return nil, err
	}

	return &table, nil
}

func UpdateOrganizationUser(user *OrganizationUser, cols ...string) error {
	err := UpdateById(user.Id, user, cols...)
	return err
}
func DeleteOrganizationUser(id string) error {
	var organizatonUser OrganizationUser
	err := GetById(id, &organizatonUser)
	if err != nil {
		return err
	}

	if organizatonUser.IsAdmin || organizatonUser.IsOwner {
		return errors.New("无法删除表格管理员")
	}

	err = DeleteById(id, &organizatonUser)
	return err
}
