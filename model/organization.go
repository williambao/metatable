package model

import (
	"errors"

	"github.com/rs/xid"
)

type OrganizationTeam struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Remark string `json:"remark"`
}

type Organization struct {
	BaseModel `xorm:"extends"`
	UserId    string `json:"user_id" xorm:"varchar(20)"`
	// Username    string `json:"username" xorm:"UNIQUE NOT NULL"`
	Name        string             `json:"name"  xorm:"varchar(200) NOT NULL"`
	Description string             `json:"description" `
	Avatar      string             `json:"avatar"`
	MemberCount int                `json:"member_count"`
	TableCount  int                `json:"table_count"`
	TableLimit  int                `json:"table_limit"`  // 可创建表格数
	MemberLimit int                `json:"member_limit"` // 可有成员数
	Teams       []OrganizationTeam `json:"teams" xorm:"text json"`
}

func UpdateOrganizationTableCount(id string) error {
	cnt, err := db.Where("organization_id = ?", id).Count(&Table{})
	if err != nil {
		return err
	}
	org := new(Organization)
	org.Id = id
	org.TableCount = int(cnt)
	err = UpdateById(id, org, "table_count")
	return err
}

func GetOrganizationByName(name string) (*Organization, error) {
	if len(name) == 0 {
		return nil, errors.New("无效的团队id")
	}
	var item Organization
	has, err := db.Where("name = ?", name).Get(&item)
	if err != nil {
		return &item, err
	}
	if !has {
		return &item, errors.New("未找到团队")
	}
	if item.Teams == nil {
		item.Teams = []OrganizationTeam{}
	}

	return &item, nil
}

func GetOrganizationById(id string) (*Organization, error) {
	if len(id) == 0 {
		return nil, errors.New("无效的团队id")
	}
	var item Organization
	err := GetById(id, &item)
	if err != nil {
		return &item, err
	}
	if item.Teams == nil {
		item.Teams = []OrganizationTeam{}
	}

	return &item, nil
}

func CreateOrgnization(organization *Organization) error {
	if len(organization.Name) == 0 {
		return errors.New("请输入团队名")
	}

	organization.MemberCount = 1
	organization.MemberLimit = 10
	organization.TableLimit = 10

	session := db.NewSession()
	defer session.Close()
	if organization.Id == "" {
		organization.Id = xid.New().String()

	}
	err := session.Begin()

	_, err = session.Insert(organization)
	if err != nil {
		session.Rollback()
		return err
	}

	// 创建表格时, 把当前用户存入到表格用户表, 并设置成表格管理员
	user := &OrganizationUser{}
	user.Id = xid.New().String()
	user.OrganizationId = organization.Id
	user.UserId = organization.UserId
	user.IsAdmin = true
	user.IsOwner = true
	_, err = session.Insert(user)
	if err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	if err != nil {
		return err
	}
	return nil
}

func UpdateOrganization(organization *Organization, cols ...string) error {
	session := db.NewSession()
	defer session.Close()

	// 当是下拉选项字段时, 把options存起来
	for index, team := range organization.Teams {
		if len(team.Id) == 0 {
			organization.Teams[index].Id = xid.New().String()
		}
	}

	err := session.Begin()
	_, err = session.Cols(cols...).ID(organization.Id).Update(organization)
	if err != nil {
		session.Rollback()
		return err
	}

	if StringIn(cols, "teams") {
		users, _ := GetOrganizationUsers(organization.Id)
		ids := make([]string, 0)
		for _, u := range users {
			if u.TeamId == "" {
				continue
			}

			exist := false
			for _, t := range organization.Teams {
				if t.Id == u.TeamId {
					exist = true
				}
			}
			if !exist {
				ids = append(ids, u.Id)
			}
		}
		if len(ids) > 0 {
			obj := new(OrganizationUser)
			obj.TeamId = ""
			db.In("user_id = ?", ids).Cols("team_id").Update(obj)
		}
	}

	session.Commit()
	return nil
}

func GetUserOrganizations(user *User) ([]Organization, error) {
	list := make([]Organization, 0)
	users := make([]OrganizationUser, 0)
	err := db.Where("user_id = ?", user.Id).Find(&users)
	if err != nil {
		return list, err
	}
	ids := make([]string, 0)
	for _, u := range users {
		ids = append(ids, u.OrganizationId)
	}

	err = db.In("id", ids).Find(&list)

	for idx, item := range list {
		if item.Teams == nil {
			list[idx].Teams = []OrganizationTeam{}
		}
	}

	return list, err
}

func DeleteOrganization(id string, user *User) error {
	u, err := GetOrganizationUser(id, user.Id)
	if err != nil {
		return err
	}
	if !u.IsOwner {
		return errors.New("您无法删除此团队. 只有团队所有人可删除.")
	}

	// 若有表格则不能删除
	cnt, _ := db.Where("organization_id = ?", id).Count(&Table{})
	if cnt > 0 {
		return errors.New("此团队下有表格, 请先删除表格后再试~")
	}

	session := db.NewSession()
	err = session.Begin()

	// 删除用户
	_, err = session.Where("organization_id = ?", id).Delete(&OrganizationUser{})
	if err != nil {
		session.Rollback()
		return err
	}

	// 删除团队
	// _, err = session.Where("organization_id = ?", id).Delete(&OrganizationTeam{})
	// if err != nil {
	// 	session.Rollback()
	// 	return err
	// }

	// 删除表格
	err = DeleteById(id, &Organization{})
	if err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	return err
}
