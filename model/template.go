package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"
)

type Template struct {
	BaseModel   `xorm:"extends"`
	Name        string     `json:"name" form:"name" xorm:"varchar(200) NOT NULL"`
	Description string     `json:"description" form:"Text"`
	Icon        string     `json:"icon"`
	Version     string     `json:"version"`
	UserId      string     `json:"user_id" xorm:"varchar(20) not null"`
	UseCount    int        `json:"use_count"`
	Price       float64    `json:"float64"`
	ColumnCount int        `json:"column_count"`
	IsActive    bool       `json:"is_active"`
	IsInit      bool       `json:"is_init"`
	DeletedAt   *time.Time `json:"-" xorm:"deleted" `
	UserName    string     `json:"username" xorm:"-"`
}

func GetTemplates(name string, isActive bool, isInit bool, page int, pagesize int) ([]Template, error) {
	list := make([]Template, 0)
	where := db.Where("")
	// where := db.Where("")
	if isActive {
		where = where.And("is_active = ?", true)
	}
	if isInit {
		where = where.And("is_init = ?", true)
	}
	if name != "" {
		where = where.And("name like '%?%'", name)
	}

	err := where.Limit(pagesize, (page-1)*pagesize).Desc("created_at").Find(&list)
	if len(list) > 0 {
		userIds := make([]string, 0)
		for _, record := range list {
			if !isExists(userIds, record.UserId) {
				userIds = append(userIds, record.UserId)
			}
			// if !isExists(userIds, record.UpdatedBy) {
			// 	userIds = append(userIds, record.UpdatedBy)
			// }
		}
		fmt.Println(userIds)
		users := make([]User, 0)
		err := db.In("id", userIds).Cols("nickname", "avatar", "id").Find(&users)
		if err != nil {
			return list, err
		}
		for _, user := range users {
			for index, record := range list {
				if record.UserId == user.Id {
					list[index].UserName = user.Nickname
					// list[index].UserAvatar = user.Avatar
				}
				// if record.UpdatedBy == user.Id {
				// 	list[index].UpdateUserName = user.Nickname
				// 	list[index].UpdateUserAvatar = user.Avatar
				// }
			}
		}
	}

	return list, err
}

func UpdateTemplateColumnCount(id string) error {
	cnt, _ := db.Where("template_id = ?", id).Count(&TemplateColumn{})
	obj := new(Template)
	obj.ColumnCount = int(cnt)
	err := UpdateById(id, obj, "column_count")
	return err
}

func GetTemplateById(id string) (*Template, error) {
	if len(id) == 0 {
		return nil, errors.New("无效的模板id")
	}
	var item Template
	err := GetById(id, &item)

	return &item, err
}

func CreateTemplate(template *Template) error {
	if len(template.Name) == 0 {
		return errors.New("请输入模板名")
	}
	if template.UserId == "" {
		return errors.New("无效用户id")
	}

	template.Id = xid.New().String()
	template.IsActive = true

	err := Insert(template)

	return err
}

func UpdateTemplate(template *Template, cols ...string) error {
	err := UpdateById(template.Id, template, cols...)
	return err
}

func DeleteTemplate(id string, user *User) error {
	session := db.NewSession()
	err := session.Begin()

	_, err = session.Where("template_id = ?", id).Delete(&TemplateColumn{})
	if err != nil {
		session.Rollback()
		return err
	}

	// 删除表格
	err = DeleteById(id, &Template{})
	if err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	return err
}
