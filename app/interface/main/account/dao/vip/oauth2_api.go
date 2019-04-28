package vip

import (
	"context"
	"net/url"

	"go-common/app/interface/main/account/model"
	"go-common/library/ecode"

	pkgerr "github.com/pkg/errors"
)

const (
	_oauth2UserInfoPath = "/oauth2/user_info"
)

//OAuth2ByCode get user info by oauth2 code.
func (d *Dao) OAuth2ByCode(c context.Context, a *model.ArgAuthCode) (data *model.OAuth2InfoResp, err error) {
	params := url.Values{}
	params.Add("code", a.Code)
	params.Add("grant_type", "authorization_code")
	var res struct {
		Code int                   `json:"code"`
		Data *model.OAuth2InfoResp `json:"data"`
	}
	if err = d.cl.get(c, d.c.Host.PassportCom, _oauth2UserInfoPath, a.IP, params, &res); err != nil {
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		err = pkgerr.Wrap(err, "dao oauth2 userinfo")
		return
	}
	data = res.Data
	return
}
