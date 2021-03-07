package model

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/williambao/metatable/config"
	"xorm.io/xorm"
	"xorm.io/xorm/caches"
)

var (
	db *xorm.Engine
)

func DB() *xorm.Engine {
	return db
}

func OpenDB(cfg *config.Config) error {
	var err error
	if cfg.Database.IsMySQL {
		server := fmt.Sprintf("%s:%s", cfg.Database.Host, cfg.Database.Port)
		connString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&&loc=Local",
			cfg.Database.Username, cfg.Database.Password, server, cfg.Database.DatabaseName)
		db, err = xorm.NewEngine("mysql", connString)
	} else {
		db, err = xorm.NewEngine("sqlite3", cfg.Database.Host)

	}
	if err != nil {
		return err
	}

	db.ShowSQL(true)

	return nil
}

// 同步数结构到数据库
func Migrate() error {
	err := db.Sync2(
		new(User),
		new(Invite),
		new(Table),
		new(TableUser),
		new(TableView),
		new(Column),
		new(Record),
		new(File),
		new(Sms),
		new(Template),
		new(TemplateColumn),
		new(Organization),
		new(OrganizationUser),
	)

	return err
}

// 打开缓存
func EnableCache() error {
	cacher := caches.NewLRUCacher(caches.NewMemoryStore(), 1000)
	db.MapCacher(&User{}, cacher)
	db.MapCacher(&Table{}, cacher)
	db.MapCacher(&TableUser{}, cacher)
	db.MapCacher(&Column{}, cacher)
	db.MapCacher(&Organization{}, cacher)

	return nil
}

type BaseModel struct {
	Id string `json:"id" xorm:"varchar(20) pk"`

	CreatedAt   *time.Time `json:"created_at,omitempty" xorm:"created"`
	CreatedBy   string     `json:"created_by" xorm:"varchar(20)"`
	CreatedUser *User      `json:"created_user,omitempty" xorm:"-"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" xorm:"updated"`
	UpdatedBy   string     `json:"updated_by" xorm:"varchar(20)"`
	UpdatedUser *User      `json:"Updated_user,omitempty" xorm:"-"`
	// DeletedAt *time.Time `json:"-" xorm:"deleted" `
}

// 判断某数组是否包含某字符值
func StringIn(array []string, item string) bool {
	for _, str := range array {
		if item == str {
			return true
		}
	}
	return false
}
