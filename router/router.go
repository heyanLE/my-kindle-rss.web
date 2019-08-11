package router

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"log"
	"my-kindle-rss/controller/api"
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
	beego.Router("/api/v1/register", &api.Controller{})
	beego.Router("/api/v1/user", &api.Controller{})
	beego.Router("/api/v1/register-verify", &api.Controller{})
	beego.Router("/api/v1/feed-list", &api.Controller{})
	beego.Router("/api/v1/feed", &api.Controller{})
	beego.Router("/api/v1/feed-refresh", &api.Controller{})
	beego.Router("/api/v1/feed-from-list", &api.Controller{})
	beego.Router("/api/v1/feed-class-list", &api.Controller{})

	beego.Router("/api/v1/article", &api.Controller{})
	beego.Router("/api/v1/change", &api.Controller{})

	beego.Router("/api/v1/captcha", &api.CaptchaController{})
	beego.Handler("/api/v1/captcha/img/*", api.InitHandler())

}

func TransparentStatic(ctx *context.Context) {
	orPath := ctx.Request.URL.Path
	log.Println(orPath)
	//如果url中有api，则取消静态路由重定向
	if strings.Index(orPath, "api") >= 0 {
		return
	}
	beego.Info("OrPath:", orPath)
	switch orPath {
	case "/home":
		orPath = orPath + ".html"
	}
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/"+orPath)
}
