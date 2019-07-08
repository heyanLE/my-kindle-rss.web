package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"my-kindle-rss/models"
	"my-kindle-rss/web/controller/api"
	"my-kindle-rss/web/router"
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

	beego.Run()
}
