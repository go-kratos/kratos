package dao

import (
	"context"
	"net/url"

	"go-common/app/service/main/coupon/model"
	"go-common/library/ecode"

	"github.com/pkg/errors"
)

//CaptchaToken get captcha token.
func (d *Dao) CaptchaToken(c context.Context, bid string, ip string) (res *model.Token, err error) {
	args := url.Values{}
	args.Add("bid", bid)
	resq := new(struct {
		Code int          `json:"code"`
		Data *model.Token `json:"data"`
	})
	if err = d.client.Get(c, d.c.Property.CaptchaTokenURL, ip, args, resq); err != nil {
		err = errors.Wrapf(err, "dao captcha token do")
		return
	}
	if resq.Code != ecode.OK.Code() {
		err = ecode.Int(resq.Code)
		err = errors.Wrapf(err, "dao captcha token code")
		return
	}
	res = resq.Data
	return
}

//CaptchaVerify get captcha verify.
func (d *Dao) CaptchaVerify(c context.Context, code string, token string, ip string) (err error) {
	args := url.Values{}
	args.Add("token", token)
	args.Add("code", code)
	resq := new(struct {
		Code int `json:"code"`
	})
	if err = d.client.Post(c, d.c.Property.CaptchaVerifyURL, ip, args, resq); err != nil {
		err = errors.Wrapf(err, "dao captcha verify do")
		return
	}
	if resq.Code != ecode.OK.Code() {
		err = ecode.CouponCodeVerifyFaildErr
	}
	return
}
