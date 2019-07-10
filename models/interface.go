package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"my-kindle-rss/utils"
	"strconv"
	"time"
)

const (
	PushEmailSuffix = "@ikndle.cn"
)

func Login(email string, password string) (User, error) {
	o := orm.NewOrm()
	u := User{Email: email}
	e := o.Read(&u, "Email")
	if e == orm.ErrNoRows {
		return u, IncorrectUserOPassErr
	} else if e == nil {
		if u.PasswordHash == utils.Hash(password) {
			return u, nil
		} else {
			return User{Email: email}, IncorrectUserOPassErr
		}
	} else {
		beego.Error("Login Error #1 : ", e.Error())
		return u, e
	}
}

func Register(email string, password string) (User, error) {
	o := orm.NewOrm()
	u := User{Email: email}
	e := o.Read(&u, "Email")
	if e == orm.ErrNoRows {
		pe, e := GetPushEmail()
		if e != nil {
			beego.Error("Register Error #1 : ", e.Error())
			return u, e
		}
		u.Email = email
		u.PasswordHash = utils.Hash(password)
		u.Balance = 14
		u.PushAuto = false
		u.CreateTime = time.Now()
		u.PushTime = 18
		u.PushEmail = &pe
		pe.UnderUser[len(pe.UnderUser)] = &u
		if len(pe.UnderUser) >= 16 {
			pe.OverLoad = true
		}
		err := o.Begin()
		if err != nil {
			beego.Error("Register Error #2 : ", err.Error())
			_ = o.Rollback()
			return u, err
		}
		_, err = o.Insert(&u)
		if err != nil {
			beego.Error("Register Error #3 : ", err.Error())
			_ = o.Rollback()
			return u, err
		}
		_, err = o.Update(&pe)
		if err != nil {
			_ = o.Rollback()
			return u, err
		}
		err = o.Commit()
		if err != nil {
			beego.Error("Register Error #4 : ", err.Error())
		}
		return u, err

	} else if e == nil {
		return User{Email: email}, EmailExistErr
	} else {
		beego.Error("Register Error #4 : ", e.Error())
		return u, e
	}
}

func GetPushEmail() (PushEmail, error) {
	pe := PushEmail{OverLoad: false}
	o := orm.NewOrm()
	e := o.Read(&pe, "OverLoad")
	if e == orm.ErrNoRows {
		err := o.Begin()
		if err != nil {
			_ = o.Rollback()
			return pe, err
		}
		id, err := o.Insert(&pe)
		if err != nil {
			_ = o.Rollback()
			return pe, err
		}
		pe.Address = strconv.FormatInt(id, 10) + PushEmailSuffix
		_, err = o.Update(&pe)
		if err != nil {
			err = o.Rollback()
			return pe, err
		}
		err = o.Commit()
		return pe, err
	} else {
		return pe, e
	}
}

func FeedList(feed *[]RssFeed) error {
	o := orm.NewOrm()
	_, e := o.QueryTable(new(RssFeed)).All(feed)
	return e
}

func UserFeed(user *User) ([]*RssFeed, error) {
	o := orm.NewOrm()
	e := o.Read(user)
	if e != nil {
		return *new([]*RssFeed), UserNotFound
	}
	lr := user.FeedIdList

	for i := 0; i < len(lr); i++ {
		_ = o.Read(&lr[i])
	}
	return lr, nil
}

func UserFeedPost(user *User, id int64) error {
	o := orm.NewOrm()
	feed := RssFeed{Id: id}
	e := o.Read(feed)
	if e == orm.ErrMissPK || e == orm.ErrNoRows {
		return FeedNotFound
	} else if e == nil {
		e = o.Read(&user)
		if e != nil {
			return UserNotFound
		}
		user.FeedIdList[len(user.FeedIdList)] = &feed
		feed.Subscriber[len(feed.Subscriber)] = user
		e = o.Begin()
		if e != nil {
			_ = o.Rollback()
			return e
		}
		_, e = o.Update(&feed)
		if e != nil {
			_ = o.Rollback()
			return e
		}
		_, e = o.Update(user)
		if e != nil {
			_ = o.Rollback()
			return e
		}
		e = o.Commit()
		if e != nil {
			_ = o.Rollback()
			return e
		}
		return nil
	} else {
		return e
	}
}
