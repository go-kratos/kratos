package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_token  = "/x/internal/v1/captcha/token"
	_verify = "/x/internal/v1/captcha/verify"

	_openCode       = "/pay/codeOpen"
	_passportDetail = "/intranet/acc/detail"

	token = "vip-service.token"
)

//GetPassportDetail get passport detail
func (d *Dao) GetPassportDetail(c context.Context, mid int64) (res *model.PassportDetail, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	val := url.Values{}
	val.Add("mid", strconv.FormatInt(mid, 10))
	resq := new(struct {
		Code int                   `json:"code"`
		Data *model.PassportDetail `json:"data"`
	})
	defer func() {
		log.Info("get passport detail is error url:%+v params:%+v resq:%+v err:%+v", d.passportDetail, val, resq, err)
	}()
	if err = d.client.Get(c, d.passportDetail, ip, val, resq); err != nil {
		err = errors.WithStack(err)
		return
	}
	if resq.Code != ecode.OK.Code() {
		err = ecode.Int(resq.Code)
		return
	}
	res = resq.Data
	return
}

//OpenCode open code.
func (d *Dao) OpenCode(c context.Context, mid, batchCodeID int64, unit int32, remark, code string) (data *model.CommonResq, err error) {
	data = new(model.CommonResq)
	val := url.Values{}
	val.Add("token", token)
	val.Add("remark", remark)
	val.Add("code", code)
	val.Add("mid", fmt.Sprintf("%v", mid))
	if err = d.doRequest(c, d.c.Property.VipURL, _openCode, "", val, data, d.client.Post); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//GetOpenInfo get open info.
func (d *Dao) GetOpenInfo(c context.Context, code string) (data *model.OpenCodeResp, err error) {
	data = new(model.OpenCodeResp)
	val := url.Values{}
	val.Add("token", token)
	val.Add("code", code)
	if err = d.doRequest(c, d.c.Property.VipURL, _openCode, "", val, data, d.client.Get); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//GetToken get open token.
func (d *Dao) GetToken(c context.Context, bid string, ip string) (data *model.TokenResq, err error) {
	data = new(model.TokenResq)
	val := url.Values{}
	val.Add("bid", bid)
	if err = d.doRequest(c, d.c.Property.APICoURL, _token, ip, val, data, d.client.Get); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//Verify verify code.
func (d *Dao) Verify(c context.Context, code, token, ip string) (data *model.TokenResq, err error) {
	data = new(model.TokenResq)
	val := url.Values{}
	val.Add("token", token)
	val.Add("code", code)
	if err = d.doRequest(c, d.c.Property.APICoURL, _verify, ip, val, data, d.client.Post); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (d *Dao) doRequest(c context.Context, basicURL, path, IP string, params url.Values, data interface{},
	fn func(c context.Context, uri, ip string, params url.Values, res interface{}) error) (err error) {
	var (
		url = basicURL + path
	)
	if len(IP) == 0 {
		IP = "127.0.0.1"
	}

	err = fn(c, url, IP, params, data)
	log.Info("reques url %v params:%+v result:%+v", url, params, data)
	if err != nil {
		log.Error("request error %+v", err)
		err = errors.WithStack(err)
		return
	}
	return

}
