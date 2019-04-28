package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go-common/app/service/main/passport-sns/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_weiboAuthorizeUrl   = "https://api.weibo.com/oauth2/authorize"
	_weiboAccessTokenUrl = "https://api.weibo.com/oauth2/access_token"
)

// WeiboAuthorize .
func (d *Dao) WeiboAuthorize(c context.Context, appID, redirectURL, display string) (url string) {
	return fmt.Sprintf(_weiboAuthorizeUrl+"?client_id=%s&redirect_uri=%s&scope=all", appID, redirectURL)
}

// WeiboOauth2Info .
func (d *Dao) WeiboOauth2Info(c context.Context, code, redirectUrl string, app *model.SnsApps) (res *model.Oauth2Info, err error) {
	accessResp, err := d.weiboAccessToken(c, code, app.AppID, app.AppSecret, redirectUrl)
	if err != nil {
		return nil, err
	}
	res = &model.Oauth2Info{
		Token:   accessResp.Token,
		Refresh: accessResp.Refresh,
		Expires: time.Now().Unix() + accessResp.Expires,
		OpenID:  accessResp.OpenID,
		UnionID: accessResp.OpenID,
	}
	return
}

// weiboAccessToken .
func (d *Dao) weiboAccessToken(c context.Context, code, appID, appSecret, redirectUrl string) (resp *model.WeiboAccessResp, err error) {
	var (
		res    *http.Response
		params = url.Values{}
	)
	params.Set("client_id", appID)
	params.Set("client_secret", appSecret)
	params.Set("grant_type", "authorization_code")
	params.Set("code", code)
	params.Set("redirect_uri", redirectUrl)

	res, err = d.client.PostForm(_weiboAccessTokenUrl, params)
	if err != nil {
		log.Error("d.weiboAccessToken error(%+v) code(%s) appID(%s)", err, code, appID)
		return nil, err
	}
	defer res.Body.Close()
	dc := json.NewDecoder(res.Body)
	resp = new(model.WeiboAccessResp)
	if err = dc.Decode(resp); err != nil {
		return
	}
	if resp.Code == _respCodeSuccess {
		return resp, nil
	}
	log.Error("request weibo failed with code(%d) desc(%s)", resp.Code, resp.Description)
	return nil, ecode.PassportSnsRequestErr
}
