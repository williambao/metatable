package model

import "github.com/rs/xid"
import "time"

type Sms struct {
	BaseModel `xorm:"extends"`
	UserId    string `json:"user_id" xorm:"varchar(20)"`
	Phone     string `json:"phone" xorm:"varchar(20)"`
	Code      string `json:"code" xorm:"varchar(20)"`
	RequestId string `json:"request_id" xorm:"varchar(100)"`
	IsActive  bool   `json:"is_active"`
}

func CreateSms(sms *Sms) error {
	sms.Id = xid.New().String()
	sms.IsActive = true
	err := Insert(sms)
	return err
}

// 验证码是否OK
func IsValidSms(phone, code string) bool {
	var sms Sms
	has, err := db.Where("phone = ? and code = ?", phone, code).Desc("created_at").Get(&sms)
	if !has || err != nil || &sms == nil {
		return false
	}

	now := time.Now()
	now = now.Add(-10 * time.Minute)

	if !sms.IsActive || sms.CreatedAt.Before(now) {
		return false
	}

	return true
}

// 把短信验证码设置为[已过期]
func SetSmsInActive(phone, code string) error {
	var sms []Sms
	err := db.Where("phone = ? and code = ?", phone, code).Find(&sms)
	if err != nil {
		return err
	}

	for _, c := range sms {
		c.IsActive = false
		err := UpdateById(c.Id, c, "is_active")
		if err != nil {
			// todo: print err to log
		}
	}

	return nil
}
