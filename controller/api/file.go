package api

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type FileItem struct {
	Id       int    `json:"id"`
	Ip       string `json:"ip"`
	Url      string `json:"url"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
	TimeUnix int64  `json:"time_unix"`
	Method   string `json:"method"`
}

func (*FileItem) TableName() string {
	return "file"
}

func setUp(item *FileItem) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		beego.Error("Set Up Item Err #1 :", err.Error())
		_ = o.Rollback()
		return err
	}
	_, err = o.Insert(item)
	if err != nil {
		beego.Error("Set Up Item Err #2 :", err.Error())
		_ = o.Rollback()
		return err
	}
	err = o.Commit()
	if err != nil {
		beego.Error("Set Up Item Err #3 :", err.Error())
	}
	return err
}

func isTest(ip string) bool {
	o := orm.NewOrm()
	qs := o.QueryTable("file")

	i, e := qs.Filter("Url", "/api/v1/register").Filter("Ip", ip).Filter("Code__gt", 200).Filter("TimeUnix__gt", int64(time.Now().Unix()-int64(60*60))).Count()
	beego.Info("IsTest I =>", i)
	if i <= 5 {
		i, e = qs.Filter("Url", "/api/v1/user").Filter("Method", "POST").Filter("Ip", ip).Filter("Code__gt", 200).Filter("TimeUnix__gt", int64(time.Now().Unix()-int64(60*60))).Count()
	}

	if e != nil {
		beego.Error("IsTest Err :", e.Error())
	}

	beego.Info("IsTest I =>", i)
	beego.Info("isTest Time =>", time.Now().Unix()-60*60)
	return i > 5
}
