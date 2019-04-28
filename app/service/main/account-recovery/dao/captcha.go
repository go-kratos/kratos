package dao

import (
	"context"
	"net/url"

	"go-common/app/service/main/account-recovery/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// GetToken get open token.
func (d *Dao) GetToken(c context.Context, bid string) (res *model.TokenResq, err error) {
	params := url.Values{}
	params.Add("bid", bid)
	if err = d.httpClient.Get(c, d.c.CaptchaConf.TokenURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("GetToken HTTP request err(%v)", err)
		return
	}
	if res.Code != 0 {
		log.Error("GetToken service return err(%v)", res.Code)
		err = ecode.Int(int(res.Code))
		return
	}
	return
}

// Verify verify code.
func (d *Dao) Verify(c context.Context, code, token string) (ok bool, err error) {
	params := url.Values{}
	params.Add("token", token)
	params.Add("code", code)
	res := new(struct {
		Code int `json:"code"`
	})
	if err = d.httpClient.Post(c, d.c.CaptchaConf.VerifyURL, metadata.String(c, metadata.RemoteIP), params, res); err != nil {
		log.Error("Verify HTTP request err(%v)", err)
		return
	}
	if res.Code != 0 {
		log.Error("Verify service return err(%v)", res.Code)
		err = ecode.Int(res.Code)
		return
	}
	return true, nil
}
