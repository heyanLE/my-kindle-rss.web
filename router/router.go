package router

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"log"
	"my-kindle-rss/web/controller/api"
	"net/http"
	"strings"
)

func InitRouter() {
	/*
		过滤路由 不包含api的请求重定向至static/html/*
	*/
	beego.InsertFilter("/", beego.BeforeRouter, TransparentStatic)
	beego.InsertFilter("/*", beego.BeforeRouter, TransparentStatic)

	/*
		Api路由
	*/
	beego.Router("/api/v1/register",&api.Controller{})
	beego.Router("/api/v1/user",&api.Controller{})
	beego.Router("/api/v1/register-verify",&api.Controller{})

	beego.Router("/api/v1/captcha",&api.CaptchaController{})
	beego.Handler("/api/v1/captcha/img/*",api.InitHandler())

}

func TransparentStatic(ctx *context.Context){
	orPath := ctx.Request.URL.Path
	log.Println(orPath)
	//如果url中有api，则取消静态路由重定向
	if strings.Index(orPath, "api") >= 0 {
		return
	}
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/html/" + orPath)
}