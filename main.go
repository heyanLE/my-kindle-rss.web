package main

import (
	"encoding/gob"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/astaxie/beego/session/mysql"
	"my-kindle-rss/controller/api"
	"my-kindle-rss/models"
	"my-kindle-rss/router"
)

func main() {

	/*
		初始化Models
	*/
	models.InitModels()
	orm.RegisterModel(new(api.FileItem))
	models.RunSyncDb()

	/*
		初始化Router
	*/
	router.InitRouter()

	/*
		初始化Session
	*/
	InitSession()

	beego.Run()
}

func InitSession() {
	//beego的session序列号是用gob的方式，因此需要将注册models.User
	gob.Register(new(models.User))
	//https://beego.me/docs/mvc/controller/session.md
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "mykindlerss-key"
	beego.BConfig.WebConfig.Session.SessionProvider = "mysql"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = "root:root@heyanle...@tcp(localhost:3306)/my_kindle_rss?charset=utf8"
}
