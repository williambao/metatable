package model

import (
	"errors"

	"github.com/rs/xid"
)

type Invite struct {
	BaseModel `xorm:"extends"`
	UserId    string `json:"user_id"`
	Code      string `json:"code"`
	IsActive  bool   `json:"is_active"`
}

// 创建新的邀请码
func CreateInvite(user *User) (*Invite, error) {
	data := &Invite{}
	data.UserId = user.Id
	data.Code = xid.New().String()

	err := Insert(data)

	return data, err
}

// 检查邀请码是否有效
func IsValidInviteCode(code string) (bool, error) {
	data := new(Invite)
	has, err := db.Where("code = ?", code).Get(data)
	// fmt.Printf("has :%v, err: %v", has, err)
	if err != nil {
		return false, err
	}
	if !has {
		return false, errors.New("无效的邀请码")
	}

	if !data.IsActive {
		return false, errors.New("邀请码已过期")
	}

	return true, nil
}

func SetInviteCodeExpires(code string) error {
	data := new(Invite)
	_, err := db.Where("code = ?", code).Get(data)
	if data == nil {
		return errors.New("无效邀请码")
	}
	data.IsActive = false
	err = UpdateById(data.Id, data, "code")

	return err
}
