package model

import (
	"encoding/base32"
	"errors"
	"fmt"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/rs/xid"

	"github.com/williambao/metatable/utils"
)

type User struct {
	BaseModel      `xorm:"extends"`
	OrganizationId string `json:"organization_id" xorm:"unique(organization_user) NOT NULL"`
	Username       string `json:"username" form:"username" xorm:"unique(organization_user) NOT NULL"`
	Nickname       string `json:"nickname" form:"nickname"`
	Password       string `json:"-" form:"password" xorm:"NOT NULL"`
	OpenId         string `json:"-" `
	UnionId        string `json:"-" `
	SessionKey     string `json:"-" `
	Phone          string `json:"phone" form:"phone"`
	Email          string `json:"email" form:"email"`
	Avatar         string `json:"avatar"`
	Sex            int    `json:"sex"`
	Country        string `json:"country"`
	Province       string `json:"province"`
	City           string `json:"city"`

	RoleIds []string `json:"role_ids" xorm:"text json"`

	IsActive bool `json:"is_active"`
	IsAdmin  bool `json:"is_admin"`

	InviteCode string `json:"invite_code"`

	Token     string `xorm:"-" json:"token,omitempty"`
	ExpiredIn int64  `xorm:"-" json:"expired_in,omitempty"`

	TeamId string `json:"team_id" `

	// Hash is a unique token used to sign tokens.
	Hash string `json:"-" xorm:"varchar(128) UNIQUE NOT NULL"`
}

func (u *User) InRole(name string) bool {
	if u.RoleIds == nil {
		return false
	}

	return StringIn(u.RoleIds, name)
}

func (u *User) EncodePassword() error {
	u.Password = utils.EncodeMd5(u.Password)
	return nil
}

//  compare user password
func (u *User) IsValidPassword(pass string) bool {
	newUser := &User{Password: pass}
	newUser.EncodePassword()

	return u.Password == newUser.Password
}

// 更新用户密码
func (u *User) UpdatePassword(old string, pass string) error {
	if !u.IsValidPassword(old) {
		return errors.New("密码验证失败")
	}
	user := &User{Password: pass}
	user.EncodePassword()

	return UpdateById(u.Id, user, "password")
}

func (u *User) UpdateAccount(phone, username, password string) error {
	user := &User{Password: password}
	user.EncodePassword()
	user.Phone = phone
	user.Username = username

	return UpdateById(u.Id, user, "username", "phone", "password")
}

func GetUserByOpenId(openid string) (*User, error) {
	if len(openid) == 0 {
		return nil, errors.New("invalid openid")
	}
	var user User
	has, err := db.Where("open_id = ?", openid).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("未找到用户")
	}
	return &user, nil
}

func GetUserById(id string) (*User, error) {
	if len(id) == 0 {
		return nil, errors.New("未找到用户")
	}
	var user User
	err := GetById(id, &user)
	return &user, err
}

func GetUserByHash(hash string) (*User, error) {
	if len(hash) == 0 {
		return nil, errors.New("invalid hash")
	}
	var user User
	has, err := db.Where("hash = ?", hash).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("未找到用户")
	}
	return &user, nil
}

func GetUserByPhone(phone string) *User {
	var user User
	has, err := db.Where("phone = ?", phone).Get(&user)
	if !has || err != nil {
		return nil
	}
	return &user
}

func GetUserByUserName(username string) *User {
	var user User
	has, err := db.Where("username = ?", strings.ToLower(username)).Get(&user)
	if !has || err != nil {
		return nil
	}
	return &user
}

// 检查用户名是否已被注册
func IsUserNameExist(username string) bool {
	// if username == "" {
	// 	return false
	// }
	username = strings.ToLower(username)
	has, _ := db.Where("username = ?", username).Get(&User{})
	return has
}

func CreateUser(user *User) error {

	// 用户的手机号就是用户名
	// if user.Username == "" {
	// 	return fmt.Errorf("请输入username字段")
	// }

	if user.Id == "" {
		user.Id = xid.New().String()
	}

	if user.Username == "" {
		user.Username = user.Id
		user.Password = xid.New().String()
	}

	if user.Nickname == "" {
		user.Nickname = user.Username
	}

	if user.OrganizationId == "" {
		return errors.New("未知团队")
	}

	if IsUserNameExist(user.Username) {
		return fmt.Errorf("此账号(%s)已经注册,无法重复注册", user.Username)
	}
	// user.Username = strings.ToLower(user.Username)
	//user.Email = strings.ToLower(user.Email)
	user.IsActive = true
	user.Hash = base32.StdEncoding.EncodeToString(
		securecookie.GenerateRandomKey(32),
	)

	user.EncodePassword()

	err := Insert(user)
	if err != nil {
		return err
	}

	// 把邀请码过期掉. 一个邀请码只能用一个人

	return nil
}

func getUsersByIds(ids []string, cols []string) ([]User, error) {
	list := make([]User, 0)
	err := db.In("id", ids).Cols(cols...).Find(&list)
	return list, err
}
