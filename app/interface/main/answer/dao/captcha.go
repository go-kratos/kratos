package dao

import (
	"context"
	"net/url"

	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// Captcha Captcha.
func (d *Dao) Captcha(c context.Context) (string, string, error) {
	params := url.Values{}
	params.Set("bid", "answer")
	res := &struct {
		Code int `json:"code"`
		Data struct {
			Token string `json:"token"`
			URL   string `json:"url"`
		} `json:"data"`
		Msg string `json:"message"`
		TTL int    `json:"ttl"`
	}{}
	ip := metadata.String(c, metadata.RemoteIP)
	if err := d.captcha.Get(c, d.c.Answer.CaptchaTokenURL, ip, params, res); err != nil {
		log.Error("s.captcha.Get(%s) error(%v)", d.c.Answer.CaptchaTokenURL+"?"+params.Encode(), err)
		return "", "", err
	}
	if res.Code != 0 {
		log.Error("s.captcha.Get(%s?%s) code:%d", d.c.Answer.CaptchaTokenURL, params.Encode(), res.Code)
		return "", "", ecode.Int(res.Code)
	}
	return res.Data.Token, res.Data.URL, nil
}

// Verify Verify.
func (d *Dao) Verify(c context.Context, token, code, ip string) error {
	params := url.Values{}
	params.Set("token", token)
	params.Set("code", code)
	res := &struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
		TTL  int    `json:"ttl"`
	}{}
	if err := d.captcha.Post(c, d.c.Answer.CaptchaVerifyURL, ip, params, res); err != nil {
		log.Error("s.captcha.POST(%s) error(%v)", d.c.Answer.CaptchaVerifyURL+"?"+params.Encode(), err)
		return err
	}
	if res.Code != 0 {
		log.Error("s.captcha.POST(%s?%s) code:%d", d.c.Answer.CaptchaVerifyURL, params.Encode(), res.Code)
		return ecode.Int(res.Code)
	}
	return nil
}
