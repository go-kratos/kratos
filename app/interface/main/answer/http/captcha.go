package http

import (
	"fmt"
	"strconv"
	"strings"

	"go-common/app/interface/main/answer/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func validate(c *bm.Context) {
	var (
		mid, _      = c.Get("mid")
		params      = c.Request.Form
		challenge   = params.Get("geetest_challenge")
		validate    = params.Get("geetest_validate")
		seccode     = params.Get("geetest_seccode")
		success     = params.Get("geetest_success")
		captchaType = params.Get("captcha_type") // typ == "gt"
		cookie      = c.Request.Header.Get("cookie")
		mobile      = strings.Contains(c.Request.UserAgent(), model.MobileUserAgentFlag)
		ct          string
		sid         string
	)
	if mobile {
		ct = model.PlatH5
	} else {
		ct = model.PlatPC
	}
	successi, err := strconv.Atoi(success)
	if err != nil {
		successi = 1
	}
	sidCookie, err := c.Request.Cookie("sid")
	if err != nil {
		log.Warn("cookie do not contains sid error(%v)", err)
	} else {
		sid = sidCookie.Value
	}
	comargs := map[string]string{
		"ua":    c.Request.Header.Get("User-Agent"),
		"buvid": c.Request.Header.Get("Buvid"),
		"refer": c.Request.Header.Get("Referer"),
		"url":   c.Request.Header.Get("URL"),
		"sid":   sid,
	}
	if captchaType == model.BiliCaptcha {
		validate = params.Get("token")
		seccode = params.Get("code")
		if validate == "" || seccode == "" {
			log.Error("validate(%+v)  ", params)
			c.JSON(nil, ecode.RequestErr)
		}
	}
	req, err := svc.Validate(c, challenge, validate, seccode, ct, successi, mid.(int64), cookie, captchaType, comargs)
	if err != nil {
		log.Error("svc.Validate(%d,%d,%s,%s,%s,%s,%s,%s) error(%+v)", mid.(int64), successi, challenge, validate, seccode, ct,
			cookie, comargs, err)
		c.JSON(nil, err)
		return
	}
	cool, err := svc.Cool(c, req.HistoryID, mid.(int64))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	cool.Hid = req.HistoryID
	cool.URL = fmt.Sprintf(model.ProPassed, req.HistoryID)
	c.JSON(cool, nil)
}

// captcha get captcha
func captcha(c *bm.Context) {
	var (
		err    error
		mid, _ = c.Get("mid")
		mobile = strings.Contains(c.Request.UserAgent(), model.MobileUserAgentFlag)
		ct     string
	)
	if mobile {
		ct = model.PlatH5
	} else {
		ct = model.PlatPC
	}
	req, err := svc.Captcha(c, mid.(int64), ct, 1)
	if err != nil {
		log.Error("svc.QueCaptcha(%d,%s) error(%+v)", mid.(int64), ct, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(req, nil)
}
