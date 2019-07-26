package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"gopkg.in/gomail.v2"
	"regexp"
	"strings"
	"time"
)

const (
	EmailModel = `<html>
<body>
<div style="padding-top: 80px ; padding-bottom: 80px ; background: #f6f6f6 ; justify-xontent:center;display:flex; align_items:center;">
	<div style="padding: 80px ; background: white ; border-width: 10px ; border-color: #e9e9e9 ; margin:0px auto;text-align:center;">
		<div style="font-size: 30px ; font-weight:bold" >邮箱验证</div>
		<div style="text-align: left ; margin-top: 50px ; font-size: 14px">$Hello, $Email</div>
		<div style="text-align: left ; padding-right: 20px; font-size: 14px ;line-height: 23px;">
			您的验证码为：<span style="color: #000;font-weight:bold">$Captcha</span>
			<br>请在&#160;$Date&#160;前完成验证</div>
		<div style="height: 80px"></div>
		<a style="text-align: center ; font-size: 14px " href="http://www.mykindlerss.cn">MyKindleRss</a>
	</div>
</div>
</body></html>`
)

func Hash(str string) string {
	beego.Info("Hash:", str)
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func EmailCaptcha(email string, captcha string, date string) {

	const (
		MORNING   = "早上好"
		NOON      = "中午好"
		AFTERNOON = "下午好"
		NIGHT     = "晚上好"
	)
	hello := MORNING
	n := time.Now().Hour()
	switch {
	case n > 17:
		hello = NIGHT
		break
	case n > 14:
		hello = AFTERNOON
		break
	case n > 10:
		hello = NOON
		break
	case n > 6:
		hello = MORNING
		break
	default:
		hello = NIGHT
	}

	em := EmailModel
	ee := strings.NewReplacer("$Hello", hello, "$Email", email, "$Captcha", captcha, "$Date", date).Replace(em)
	beego.Info(em)
	m := gomail.NewMessage()
	m.SetHeader("From", "MyKindleRss"+"<"+"mykindlerss@163.com"+">")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "【邮箱验证】 MyKindleRss  - 每天推送好文章")
	m.SetBody("text/html", ee)

	d := gomail.NewDialer("smtp.163.com", 25, "mykindlerss@163.com", "heyanle137")
	_ = d.DialAndSend(m)

}

func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
