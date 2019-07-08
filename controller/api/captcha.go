package api

import (
	"github.com/astaxie/beego"
	"github.com/dchest/captcha"
	"net/http"
)

type CaptchaController struct {
	beego.Controller
}
func (c *CaptchaController) Get() {
	s := NewCaptcha()
	rb := ResponseBody{}
	rb.Code = 200
	rb.Message = "获取验证码成功"
	v := make(map[string]string)
	v["captcha_id"] = s
	rb.Value = v
	c.Data["json"] = rb
	c.ServeJSON()
}

func CaptchaVerify(captchaId string ,digits string) bool {
	return captcha.VerifyString(captchaId,digits)
}

func NewCaptcha() string {
	id := captcha.NewLen(4)
	b := captcha.RandomDigits(4)
	captcha.NewImage(id,b,100,36)
	return id
}

func InitHandler() http.Handler{
	return captcha.Server(100,36)
}
