package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

/*
User 结构体  		表名 ： user
*/
type User struct {
	Id           int64      `json:"id,omitempty"`
	Email        string     `orm:"unique";json:"email,omitempty"`
	PasswordHash string     `json:"password_hash,omitempty"`
	AimEmail     string     `json:"aim_email,omitempty"`
	PushEmail    *PushEmail `json:"push_email,omitempty,omitempty" orm:"rel(fk)"`
	PushAuto     bool       `json:"push_auto,omitempty"`
	Balance      uint       `json:"balance,omitempty"`
	FeedIdList   []*RssFeed `orm:"reverse(many)" json:"feed_id_list,omitempty"`
	CreateTime   time.Time  `json:"create_time,omitempty"`
	PushTime     uint       `json:"push_time,omitempty"`
}

func (*User) TableName() string {
	return "user"
}

/*
RssFeed结构体			表名：rss_feed
*/
type RssFeed struct {
	Id         int64   `json:"id,omitempty"`
	Name       string  `json:"name,omitempty"`
	Describe   string  `json:"describe,omitempty"`
	Value      string  `json:"value,omitempty"`
	Type       int     `json:"type"`
	Subscriber []*User `orm:"rel(m2m)" json:"subscriber,omitempty"`
}

func (*RssFeed) TableName() string {
	return "rss_feed"
}

type PushEmail struct {
	Id        int64   `json:"id,omitempty"`
	Address   string  `json:"address,omitempty"`
	UnderUser []*User `orm:"reverse(many)" json:"under_user,omitempty"`
	OverLoad  bool    `json:"over_load,omitempty"`
}

func (*PushEmail) TableName() string {
	return "push_email"
}

func InitModels() {
	DBHost := "localhost"
	DBPort := "3306"
	DBUser := "root"
	DBPassword := "root@heyanle..."
	DBName := "my_kindle_rss"
	DBStr := DBUser + ":" + DBPassword + "@tcp(" + DBHost + ":" + DBPort + ")/" + DBName + "?charset=utf8"
	e := orm.RegisterDataBase("default", "mysql", DBStr)
	if e != nil {
		beego.Error("注册Database错误 ： ", e.Error())
		return
	}
	orm.RegisterModel(new(PushEmail), new(User), new(RssFeed))
}

func RunSyncDb() {
	e := orm.RunSyncdb("default", true, true)
	if e != nil {
		beego.Error("同步错误 ： ", e.Error())
		return
	}
}
