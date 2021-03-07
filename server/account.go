package server

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"github.com/williambao/metatable/config/wechat"
	"github.com/williambao/metatable/model"
	"github.com/williambao/metatable/router/middleware/session"
	"github.com/williambao/metatable/shared/token"
	"github.com/williambao/metatable/small"
)

type register struct {
	OrganizationName string `json:"organization_name"`
	Nickname         string `json:"nickname"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Phone            string `json:"phone"`
	InviteCode       string `json:"invite_code"`
}

func Register(c *gin.Context) {
	in := &register{}
	err := c.Bind(in)
	if err != nil {
		abort(c, "无效数据")
		return
	}

	isValid, err := model.IsValidInviteCode(in.InviteCode)
	if err != nil {
		abort(c, err.Error())
		return
	}
	if !isValid {
		abort(c, "无效邀请码!")
		return
	}

	org, err := model.GetOrganizationByName(in.OrganizationName)
	if err == nil && org.Id != "" {
		// abort(c, err.Error())

		abort(c, "此公司已经申请账号")
		return
	}

	orgId := xid.New().String()

	user := new(model.User)
	user.Id = xid.New().String()
	user.Nickname = in.Nickname
	user.Phone = in.Phone
	user.Password = in.Password
	user.Email = in.Email
	user.Username = in.Email
	user.OrganizationId = orgId
	user.InviteCode = in.InviteCode

	err = model.CreateUser(user)
	if err != nil {
		log.Errorf("register user failed. %s", err.Error())
		abort(c, err.Error())
		return
	}

	org = new(model.Organization)
	org.Id = orgId
	org.Name = in.OrganizationName
	org.UserId = user.Id
	org.CreatedBy = user.Id
	org.UpdatedBy = user.Id

	err = model.CreateOrgnization(org)
	if err != nil {
		log.Errorf("register user create organization failed. %s", err.Error())
		abort(c, err.Error())
		return
	}

	exp := time.Now().Add(time.Hour * 24 * 30).Unix()
	token := token.New(token.UserToken, user.Id)
	tokenstr, err := token.SignExpires(user.Hash, exp)
	if err != nil {
		abort(c, "生成token失败")
		return
	}
	user.Token = tokenstr
	user.ExpiredIn = exp

	// 创建下默认表格
	templates, _ := model.GetTemplates("", true, true, 1, 9999)
	for _, tt := range templates {
		tb := new(model.Table)
		tb.UserId = user.Id
		tb.Name = tt.Name
		tb.IsPrivate = true
		tb.OrganizationId = org.Id
		tb.TemplateId = tt.Id
		err := model.CreateTable(tb)
		if err != nil {
			log.Errorf("register user create table faild. %s", err.Error())
		}
	}

	c.JSON(200, user)
}

func Login(c *gin.Context) {
	in := &model.Login{}
	err := c.Bind(in)
	if err != nil {
		abort(c, "无效数据")
		return
	}

	if len(in.UserName) == 0 {
		abort(c, "请输入用户名")
		return
	}

	if len(in.Password) == 0 {
		abort(c, "请输入密码")
		return
	}

	user := model.GetUserByUserName(in.UserName)
	if user == nil {
		abort(c, "未找到用户")
		return
	}

	// 检查下密码是否正确
	if !user.IsValidPassword(in.Password) {
		abort(c, "密码不正确")
		return
	}

	exp := time.Now().Add(time.Hour * 24 * 30).Unix()
	token := token.New(token.UserToken, user.Id)
	tokenstr, err := token.SignExpires(user.Hash, exp)
	if err != nil {
		abort(c, "生成token失败")
		return
	}
	user.Token = tokenstr
	user.ExpiredIn = exp

	c.JSON(200, user)
}

type loginWx struct {
	Code          string `json:"code"`
	EncryptedData string `json:"encryptedData"`
	IV            string `json:"iv"`
	Signature     string `json:"signature"`
	Nickname      string `json:"nickname"`
	Avatar        string `json:"avatar"`
	Sex           int    `json:"sex"`
	Country       string `json:"country"`
	Province      string `json:"province"`
	City          string `json:"city"`
}

func LoginWX(c *gin.Context) {
	abort(c, "微信自动登陆暂时无法使用")
	return

	in := &loginWx{}
	err := c.Bind(in)
	if err != nil {
		abort(c, "无效的code")
		return
	}

	logrus.Debugf("request data: %v", in)

	wechatConfig := wechat.FromContext(c)

	wx := small.NewWx(wechatConfig.AppId, wechatConfig.AppSecret)

	sessions, err := wx.GetWxSessionKey(in.Code)
	if err != nil {
		abort(c, err.Error())
		return
	}
	if sessions.ErrMsg != "" {
		abort(c, "当前code过期, 请重新刷新登录")
		return
	}

	userInfo, err := small.GetWxUserInfo(sessions.SessionKey, in.EncryptedData, in.IV)
	if err != nil {
		logrus.Debugf("get weixin user info failed: %s", err.Error())
	}

	logrus.Debugf("user info : %v", userInfo)

	user, err := model.GetUserByOpenId(sessions.Openid)
	if err != nil {

		user = new(model.User)
		user.SessionKey = sessions.SessionKey
		user.OpenId = sessions.Openid
		if &userInfo != nil {
			user.UnionId = userInfo.UnionId
		}
		user.Password = xid.New().String()
		user.Nickname = in.Nickname
		user.Sex = in.Sex
		user.Province = in.Province
		user.City = in.City
		user.Country = in.Country
		user.Avatar = in.Avatar
		err = model.CreateUser(user)
		if err != nil {
			logrus.Errorf("create wx user failed: %s", err.Error())
			abort(c, "创建用户失败")
			return
		}
	}

	exp := time.Now().Add(time.Hour * 24 * 30).Unix()
	token := token.New(token.UserToken, user.Id)
	tokenstr, err := token.SignExpires(user.Hash, exp)
	if err != nil {
		abort(c, "生成token失败")
		return
	}
	user.Token = tokenstr
	user.ExpiredIn = exp

	c.JSON(200, user)

}

type updateAccount struct {
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func UpdateAccount(c *gin.Context) {
	in := &updateAccount{}
	err := c.Bind(in)
	if err != nil {
		abort(c, "无效的传入格式")
		return
	}
	if in.Phone == "" {
		abort(c, "请输入手机号")
		return
	}

	suser := session.User(c)
	userId := c.Param("id")
	if userId != suser.Id && suser.IsAdmin {
		abort(c, "您无权修改此用户密码")
	}

	ok := model.IsValidSms(in.Phone, in.Code)
	if !ok {
		abort(c, "无效的验证码")
		return
	}

	if suser.Username != suser.Id {
		abort(c, "您已经注册过账号. 无法重复注册!")
		return
	}

	phoneCheck := model.GetUserByPhone(in.Phone)
	if phoneCheck != nil {
		abort(c, "此手机号已注册过账号, 无法重复注册!")
		return
	}
	usernameCheck := model.GetUserByUserName(in.Username)
	if usernameCheck != nil {
		abort(c, "此用户名已注册过账号, 无法重复注册!")
		return
	}

	err = suser.UpdateAccount(in.Phone, in.Username, in.Password)
	if err != nil {
		abort(c, err.Error())
		return
	}

	err = model.SetSmsInActive(in.Phone, in.Code)
	if err != nil {
		logrus.Errorf("set inactive sms code failed. phone: %s code: %s\n", in.Phone, in.Code)
	}

	success(c)

}

type updatePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func UpdatePassword(c *gin.Context) {
	in := &updatePassword{}
	err := c.Bind(in)
	if err != nil {
		abort(c, "无效的传入格式")
		return
	}

	suser := session.User(c)
	userId := c.Param("id")
	if userId != suser.Id && suser.IsAdmin {
		abort(c, "您无权修改此用户密码")
	}

	user, err := model.GetUserById(userId)
	if err != nil {
		abort(c, err.Error())
		return
	}

	err = user.UpdatePassword(in.OldPassword, in.NewPassword)
	if err != nil {
		abort(c, err.Error())
		return
	}

	success(c)
}

type resetPassword struct {
	Phone       string `json:"phone"`
	Username    string `json:"username"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

func ResetPassword(c *gin.Context) {
	in := &resetPassword{}
	err := c.Bind(in)
	if err != nil {
		abort(c, "无效的传入格式")
		return
	}

	suser := session.User(c)
	userId := c.Param("id")
	if userId != suser.Id && suser.IsAdmin {
		abort(c, "您无权修改此用户密码")
	}

	ok := model.IsValidSms(in.Phone, in.Code)
	if !ok {
		abort(c, "无效的验证码")
		return
	}
}
