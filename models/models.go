package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"time"
)


/*
User 结构体  		表名 ： user
 */
type User struct {
	Id 				int64 		`json:"id"`
	Email			string 		`orm:"unique";json:"email"`
	PasswordHash	string 		`json:"password_hash"`
	AimEmail		string 		`json:"aim_email"`
	PushEmail		string		`json:"push_email"`
	PushAuto		bool		`json:"push_auto"`
	Balance			uint 		`json:"balance"`
	FeedIdList		[]int64		`json:"feed_id_list"`
	CreateTime		time.Time	`json:"create_time"`
}
func (*User) TableName() string {
	return "user"
}

/*
RssFeed结构体			表名：rss_feed
 */
type RssFeed struct {
	Id 				int64 		`json:"id"`
	Name 			string 		`json:"name"`
	Describe		string 		`json:"describe"`
	Value 			string 		`json:"value"`
	SubscriberCount	int64		`json:"subscriber_count"`
}
func (*RssFeed) TableName() string  {
	return "rss_feed"
}

type PushEmail struct {
	Id 				int64 		`json:"id"`
	Address 		string 		`json:"address"`
	UnderUser		[]*User		`orm:"reverse(many)";json:"under_user"`
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

	e := orm.RegisterDataBase("default", "mysql",DBStr)
	if e != nil {
		beego.Error("注册Database错误 ： ",e.Error())
		return
	}

	orm.RegisterModel(new(User))
	orm.RegisterModel(new(RssFeed))
	orm.RegisterModel(new(PushEmail))

	e = orm.RunSyncdb("default",true,true)
	if e != nil {
		beego.Error("同步错误 ： ",e.Error())
		return
	}


}
