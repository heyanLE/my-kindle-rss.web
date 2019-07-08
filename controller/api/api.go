package api

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"my-kindle-rss/models"
	"my-kindle-rss/utils"
	"strconv"
	"time"
)

const (
	UserSessionKey = "UserSession"
)

type Controller struct {
	beego.Controller
}

type ResponseBody struct {
	Code 		int				`json:"code"`
	Message 	string			`json:"message"`
	Value 		interface{}		`json:"value"`
}

type LoginRequest struct {
	Email		string		`json:"email"`
	Password	string		`json:"password"`
}

type RegisterRequest struct {
	Email    		string 		`json:"email"`
	Password 		string 		`json:"password"`
	Token 			string		`json:"token"`
	DateUnix		int64 		`json:"date_unix"`
	Captcha 		string 		`json:"captcha"`
}

func (c *Controller) Get () {
	path := c.Ctx.Request.URL.Path
	switch path {
	case "/api/v1/user":
		c.UserGet()
		break
	case "/api/v1/register":
		ip := c.Ctx.Input.IP()
		if isTest(ip) {
			captchaId := c.GetString("_captcha_id")
			captcha := c.GetString("_captcha")
			if CaptchaVerify(captchaId,captcha) {
				c.RegisterGet()
			}else {
				c.ResponseJson(412,"需要验证且验证失败")
			}
		}else{
			c.RegisterGet()
		}
		break
	case "/api/v1/register-verify":
		c.RegisterVerifyGet()
		break
	}


}

func (c *Controller) Post () {

	method := c.Ctx.Input.Header("X-HTTP-Method-Override")
	if method == "Delete" {
		c.Delete()
		return
	}

	accept := c.Ctx.Input.Header("Accept")
	acceptLanguage := c.Ctx.Input.Header("Accept-Language")
	contentType := c.Ctx.Input.Header("Content-Type")

	ip := c.Ctx.Input.IP()
	path := c.Ctx.Request.URL.Path

	if contentType != "application/json" || accept != "application/json" || acceptLanguage != "zh-cn,zh;q=0.5" {
		c.ResponseJson(415,"接受到的表示不受支持")
		return
	}

	beego.Info("Post 请求 :",path)

	switch path {
	case "/api/v1/user":
		if isTest(ip) {

			captchaId := c.Ctx.Input.Header("X-Captcha-Id")
			captcha := c.Ctx.Input.Header("X-Captcha")

			if CaptchaVerify(captchaId,captcha) {
				c.UserPost()
			}else {
				c.ResponseJson(412,"需要验证且验证失败")
			}

		}else{
			c.UserPost()
		}
		break
	case "/api/v1/register":
		c.RegisterPost()
		break

	}
}

func (c *Controller) Delete() {

	accept := c.Ctx.Input.Header("Accept")
	acceptLanguage := c.Ctx.Input.Header("Accept-Language")
	contentType := c.Ctx.Input.Header("Content-Type")

	if contentType != "application/json" || accept != "application/json" || acceptLanguage != "zh-cn,zh;q=0.5" {
		c.ResponseJson(415,"接受到的表示不受支持")
		return
	}

	path := c.Ctx.Request.URL.Path
	switch path {
	case "/api/v1/user":
		c.UserDelete()
		break
	}

}

func (c *Controller) UserGet(){
	u := c.GetSession(UserSessionKey)
	if u != nil {
		c.ResponseJsonWithValue(200,"获取当前登录用户信息成功",&u)
	}else{
		c.ResponseJson(404,"当前没有登录用户")
	}
}

func (c *Controller) UserPost(){
	lq := LoginRequest{}
	body := c.Ctx.Input.RequestBody
	e := json.Unmarshal(body, &lq)
	if e == nil {
		u,e := models.Login(lq.Email,lq.Password)
		if e == nil {
			c.ResponseJsonWithValue(200,"登陆成功",&u)
		}else if e == models.IncorrectUserOPassErr {
			c.ResponseJson(404,"用户名或密码错误")
		}else {
			c.ResponseJson(500,"未知错误 ：" + e.Error())
		}
	} else {
		c.ResponseJson(400,"参数错误")
	}
}

func (c *Controller) UserDelete(){
	c.DelSession(UserSessionKey)
	c.ResponseJson(200,"已删除用户Session")
}

func (c *Controller) RegisterGet(){
}

func (c *Controller) RegisterPost(){
	rq := RegisterRequest{}
	body := c.Ctx.Input.RequestBody
	e := json.Unmarshal(body, &rq)
	if e == nil {
		n := time.Now().Unix()
		if n >= rq.DateUnix {
			tokenN := utils.Hash(rq.Email + strconv.FormatInt(rq.DateUnix,10) + rq.Captcha)
			if tokenN == rq.Token {
				u,e := models.Register(rq.Email,rq.Password)
				if e == models.EmailExistErr{
					c.ResponseJson(404,"邮箱已存在")
				}else if e == nil {
					c.ResponseJsonWithValue(200,"注册成功",&u)
				}else {
					c.ResponseJson(500,"未知错误 ："+e.Error())
				}
			}else {
				c.ResponseJson(404,"验证码错误")
			}
		}else{
			c.ResponseJson(303,"验证码过时")
		}
	}else {
		c.ResponseJson(400,"参数错误")
	}
}

func (c *Controller) RegisterVerifyGet(){
	email := c.GetString("_email")
	token := c.GetString("_token")
	dateUnix,e := c.GetInt64("_date_unix")
	captcha := c.GetString("_captcha")
	if e != nil || email == "" || token == "" || captcha == ""{
		c.ResponseJson(400,"参数错误")
	}
	n := time.Now().Unix()
	if n >= dateUnix{
		c.ResponseJson(303,"邮箱验证码过时")
	}
	tokenN := utils.Hash(email + strconv.FormatInt(dateUnix,10) + captcha)
	if token == tokenN{
		c.ResponseJson(200,"邮箱验证码正确")
	}else{
		c.ResponseJson(404,"验证码错误")
	}
}

func (c *Controller) ResponseJson (code int , message string){
	rb := ResponseBody{}
	rb.Code = code
	rb.Message = message
	c.Data["json"] = rb
	c.ServeJSON()
	SetUp(c.Ctx.Input.IP(),c.Ctx.Request.Method,c.Ctx.Request.URL.Path,&rb)
}

func (c *Controller) ResponseJsonWithValue (code int , message string, value interface{}){
	rb := ResponseBody{}
	rb.Code = code
	rb.Message = message
	rb.Value = value
	c.Data["json"] = rb
	c.ServeJSON()
	SetUp(c.Ctx.Input.IP(),c.Ctx.Request.Method,c.Ctx.Request.URL.Path,&rb)
}

func SetUp(ip string ,method string , path string , rb *ResponseBody){
	i := FileItem{}
	i.Ip = ip
	i.Message = (*rb).Message
	i.Code = (*rb).Code
	i.TimeUnix = time.Now().Unix()
	i.Method = method
	i.Url = path
	_ = setUp(&i)
}
