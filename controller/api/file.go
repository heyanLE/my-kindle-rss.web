package api

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type FileItem struct {
	Id 			int			`json:"id"`
	Ip 			string 		`json:"ip"`
	Url 		string		`json:"url"`
	Code    	int			`json:"code"`
	Message 	string 		`json:"message"`
	TimeUnix 	int64		`json:"time_unix"`
	Method 		string 		`json:"method"`
}
func (*FileItem) TableName () string {
	return "file"
}

func setUp(item *FileItem) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		beego.Error("Set Up Item Err #1 :",err.Error())
		_ = o.Rollback()
		return err
	}
	_,err = o.Insert(item)
	if err != nil {
		beego.Error("Set Up Item Err #2 :",err.Error())
		_ = o.Rollback()
		return err
	}
	err = o.Commit()
	if err != nil {
		beego.Error("Set Up Item Err #3 :",err.Error())
	}
	return err
}

func isTest (ip string) bool {
	o := orm.NewOrm()
	qs := o.QueryTable("file")

	cond := orm.NewCondition().Or("Path","/api/v1/register").Or("Path","/api/v1/login")

	i,e := qs.Filter("Ip",ip).Filter("TimeUnix__gte",time.Now().Unix() - 60*60).SetCond(cond).Count()
	if e != nil {
		beego.Error("IsTest Err :",e.Error())
	}
	return i >= 5
}
