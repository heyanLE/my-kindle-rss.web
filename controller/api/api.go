package api

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/memcache"
	_ "github.com/astaxie/beego/cache/redis"
	"math/rand"
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
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Value   interface{} `json:"value"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
	DateUnix int64  `json:"date_unix"`
	Captcha  string `json:"captcha"`
}

func (c *Controller) Get() {
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
			if CaptchaVerify(captchaId, captcha) {
				c.RegisterGet()
			} else {
				c.ResponseJson(412, "需要验证且验证失败")
			}
		} else {
			c.RegisterGet()
		}
		break
	case "/api/v1/register-verify":
		c.RegisterVerifyGet()
		break
	case "/api/v1/feed-list":
		c.FeedListGet()
		break
	case "/api/v1/feed":
		c.FeedGet()
		break
	}

}

func (c *Controller) Post() {
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
		c.ResponseJson(415, "接受到的表示不受支持")
		return
	}

	beego.Info("Post 请求 :", path)

	switch path {
	case "/api/v1/user":
		if isTest(ip) {

			captchaId := c.Ctx.Input.Header("X-Captcha-Id")
			captcha := c.Ctx.Input.Header("X-Captcha")

			if CaptchaVerify(captchaId, captcha) {
				c.UserPost()
			} else {
				c.ResponseJson(412, "需要验证且验证失败")
			}

		} else {
			c.UserPost()
		}
		break
	case "/api/v1/register":
		c.RegisterPost()
		break
	case "/api/v1/feed":
		c.FeedPost()
		break
	}
}

func (c *Controller) Delete() {

	accept := c.Ctx.Input.Header("Accept")
	acceptLanguage := c.Ctx.Input.Header("Accept-Language")
	contentType := c.Ctx.Input.Header("Content-Type")

	if contentType != "application/json" || accept != "application/json" || acceptLanguage != "zh-cn,zh;q=0.5" {
		c.ResponseJson(415, "接受到的表示不受支持")
		return
	}

	path := c.Ctx.Request.URL.Path
	switch path {
	case "/api/v1/user":
		c.UserDelete()
		break
	}

}

func (c *Controller) UserGet() {
	u := c.GetSession(UserSessionKey)
	if u != nil {
		c.ResponseJsonWithValue(200, "获取当前登录用户信息成功", &u)
	} else {
		c.ResponseJson(404, "当前没有登录用户")
	}
}

func (c *Controller) UserPost() {
	lq := LoginRequest{}
	body := c.Ctx.Input.RequestBody
	e := json.Unmarshal(body, &lq)
	if e == nil {
		u, e := models.Login(lq.Email, lq.Password)
		if e == nil {
			c.ResponseJsonWithValue(200, "登陆成功", &u)
			c.SetSession(UserSessionKey, &lq)
		} else if e == models.IncorrectUserOPassErr {
			c.ResponseJson(404, "用户名或密码错误")
		} else {
			beego.Error("UserPostErr :", e.Error())
			c.ResponseJson(500, "未知错误 ："+e.Error())
		}
	} else {
		c.ResponseJson(400, "参数错误")
	}
}

func (c *Controller) UserDelete() {
	c.DelSession(UserSessionKey)
	c.ResponseJson(200, "已删除用户Session")
}

func (c *Controller) RegisterGet() {
	email := c.GetString("_email")
	if email == "" || !utils.VerifyEmailFormat(email) {
		c.ResponseJson(400, "参数错误")
	} else {
		rand.Seed(time.Now().Unix())
		captcha := strconv.FormatInt(rand.Int63n(99999), 10)
		dateUnix := time.Now().Unix() + 60*20
		dateString := time.Unix(dateUnix, 0).Format("2006-01-02 15:04:05")
		token := utils.Hash(email + strconv.FormatInt(dateUnix, 10) + captcha)
		v := map[string]interface{}{
			"email":     email,
			"token":     token,
			"date_unix": dateUnix}
		c.ResponseJsonWithValue(200, "开始注册成功", &v)
		utils.EmailCaptcha(email, captcha, dateString)
	}
}

func (c *Controller) RegisterPost() {
	rq := RegisterRequest{}
	body := c.Ctx.Input.RequestBody
	e := json.Unmarshal(body, &rq)
	if e == nil {
		if !utils.VerifyEmailFormat(rq.Email) {
			c.ResponseJson(400, "参数错误")
		}
		n := time.Now().Unix()
		if n >= rq.DateUnix {
			tokenN := utils.Hash(rq.Email + strconv.FormatInt(rq.DateUnix, 10) + rq.Captcha)
			if tokenN == rq.Token {
				u, e := models.Register(rq.Email, rq.Password)
				if e == models.EmailExistErr {
					c.ResponseJson(404, "邮箱已存在")
				} else if e == nil {
					c.ResponseJsonWithValue(200, "注册成功", &u)
				} else {
					beego.Error("RegisterPostErr :", e.Error())
					c.ResponseJson(500, "未知错误 ："+e.Error())
				}
			} else {
				c.ResponseJson(404, "验证码错误")
			}
		} else {
			c.ResponseJson(303, "验证码过时")
		}
	} else {
		c.ResponseJson(400, "参数错误")
	}
}

func (c *Controller) RegisterVerifyGet() {
	email := c.GetString("_email")
	token := c.GetString("_token")
	dateUnix, e := c.GetInt64("_date_unix")
	captcha := c.GetString("_captcha")
	if e != nil || email == "" || token == "" || captcha == "" {
		c.ResponseJson(400, "参数错误")
	}
	n := time.Now().Unix()
	if n >= dateUnix {
		c.ResponseJson(303, "邮箱验证码过时")
	}
	tokenN := utils.Hash(email + strconv.FormatInt(dateUnix, 10) + captcha)
	if token == tokenN {
		c.ResponseJson(200, "邮箱验证码正确")
	} else {
		c.ResponseJson(404, "验证码错误")
	}
}

/*
	404		请登录
*/
func (c *Controller) FeedGet() {
	u := c.GetSession(UserSessionKey)
	if u == nil {
		c.ResponseJson(404, "请登录")
	} else {
		lr := u.(LoginRequest)
		u, e := models.Login(lr.Email, lr.Password)
		if e == nil {
			lrr, e := models.UserFeed(&u)
			if e == models.UserNotFound {
				c.DelSession(UserSessionKey)
				c.ResponseJson(404, "登录信息错误，请重新登陆")
			} else {
				c.ResponseJsonWithValue(200, "获取用户订阅列表成功", &lrr)
			}
		} else if e == models.IncorrectUserOPassErr {
			c.DelSession(UserSessionKey)
			c.ResponseJson(404, "登录信息错误，请重新登陆")
		} else {
			c.ResponseJson(500, "未知错误："+e.Error())
		}
	}
}

func (c *Controller) FeedPost() {
	u := c.GetSession(UserSessionKey)
	if u == nil {
		c.ResponseJson(404, "请登录")
	} else {
		lr := u.(LoginRequest)
		u, e := models.Login(lr.Email, lr.Password)
		if e == nil {
			idM := make(map[string]int64)
			body := c.Ctx.Input.RequestBody
			e = json.Unmarshal(body, &idM)
			if e == nil {
				if id, ok := idM["id"]; ok {
					e = models.UserFeedPost(&u, id)
					if e == models.FeedNotFound {
						c.ResponseJson(400, "Feed不存在")
					} else if e == models.UserNotFound {
						c.DelSession(UserSessionKey)
						c.ResponseJson(404, "登录信息错误，请重新登陆")
					} else if e == nil {
						c.ResponseJson(200, "订阅成功")
					} else {
						c.ResponseJson(500, "未知错误："+e.Error())
					}
				} else {
					c.ResponseJson(400, "参数错误")
				}
			} else {
				c.ResponseJson(400, "参数错误")
			}
		} else if e == models.IncorrectUserOPassErr {
			c.DelSession(UserSessionKey)
			c.ResponseJson(404, "登录信息错误，请重新登陆")
		} else {
			c.ResponseJson(500, "未知错误："+e.Error())
		}
	}
}

func (c *Controller) FeedListGet() {
	ca, e := cache.NewCache("redis", `{"key":"MyKindleRss","conn":"127.0.0.1:6379"}`)
	if e == nil {
		if flS := ca.Get("FeedListAll"); flS != nil {
			var ss []models.RssFeed
			e = json.Unmarshal(flS.([]byte), &ss)
			if e == nil {
				c.ResponseJsonWithValue(200, "获取Feed列表成功-缓存", &ss)
				return
			}
		}
	}
	var rl []models.RssFeed
	e = models.FeedList(&rl)
	if e == nil {
		c.ResponseJsonWithValue(200, "获取Feed列表成功-数据库", &rl)
		jsBy, e := json.Marshal(rl)
		jsStr := string(jsBy)
		beego.Info(jsStr)
		if e == nil {
			e = ca.Put("FeedListAll", jsStr, time.Hour*1)
			if e != nil {
				beego.Error("FeedListGetErr #1 : ", e.Error())
			}
		} else {
			beego.Error("FeedListGetErr #2 : ", e.Error())
		}
	} else {
		beego.Error("FeedListGetErr #3 : ", e.Error())
		c.ResponseJson(500, "未知错误："+e.Error())
	}
}

func (c *Controller) ResponseJson(code int, message string) {
	rb := ResponseBody{}
	rb.Code = code
	rb.Message = message
	c.Data["json"] = &rb
	c.ServeJSON()
	SetUp(c.Ctx.Input.IP(), c.Ctx.Request.Method, c.Ctx.Request.URL.Path, &rb)
}

func (c *Controller) ResponseJsonWithValue(code int, message string, value interface{}) {
	rb := ResponseBody{}
	rb.Code = code
	rb.Message = message
	rb.Value = value
	c.Data["json"] = &rb
	c.ServeJSON()
	SetUp(c.Ctx.Input.IP(), c.Ctx.Request.Method, c.Ctx.Request.URL.Path, &rb)
}

func SetUp(ip string, method string, path string, rb *ResponseBody) {
	i := FileItem{}
	i.Ip = ip
	i.Message = (*rb).Message
	i.Code = (*rb).Code
	i.TimeUnix = time.Now().Unix()
	i.Method = method
	i.Url = path
	_ = setUp(&i)
}
