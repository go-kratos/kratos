package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/service/main/identify-game/api/grpc/v1"
	"go-common/app/service/main/identify-game/model"

	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_createCookieURI = "/intranet/auth/createCookie/byMid"
)

// AccessToken .
func (d *Dao) AccessToken(c context.Context, accesskey, target string) (token *model.AccessInfo, err error) {
	params := url.Values{}
	params.Set("access_key", accesskey)
	var res struct {
		Code  int `json:"code"`
		Token *struct {
			Mid        string `json:"mid"`
			AppID      int64  `json:"appid"`
			Token      string `json:"access_key"`
			CreateAt   int64  `json:"create_at"`
			UserID     string `json:"userid"`
			Name       string `json:"uname"`
			Expires    string `json:"expires"`
			Permission string `json:"permission"`
		} `json:"access_info,omitempty"`
		Data *model.AccessInfo `json:"data,omitempty"`
	}
	tokenURL := d.c.Dispatcher.Oauth[target]
	if err = d.client.Get(c, tokenURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("oauth for region %s, url(%s) error(%v)", target, tokenURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("oauth for region %s, url(%s) error(%v)", target, tokenURL+"?"+params.Encode(), err)
		return
	}
	if res.Token != nil {
		t := res.Token
		var mid int64
		if mid, err = strconv.ParseInt(t.Mid, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s, 10, 64) error(%v)", t.Mid, err)
			return
		}
		var expires int64
		if expires, err = strconv.ParseInt(t.Expires, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s, 10, 64) error(%v)", t.Expires, err)
			return
		}
		token = &model.AccessInfo{
			Mid:        mid,
			AppID:      t.AppID,
			Token:      t.Token,
			CreateAt:   t.CreateAt,
			UserID:     t.UserID,
			Name:       t.Name,
			Expires:    expires,
			Permission: t.Permission,
		}
	} else {
		token = res.Data
	}
	return
}

// RenewToken request passport renewToken .
func (d *Dao) RenewToken(c context.Context, accesskey, target string) (renewToken *model.RenewInfo, err error) {
	params := url.Values{}
	params.Set("access_key", accesskey)
	var res struct {
		Code    int   `json:"code"`
		Expires int64 `json:"expires"`
		Data    struct {
			Expires int64 `json:"expires"`
		}
	}
	renewURL := d.c.Dispatcher.RenewToken[target]
	if err = d.client.Get(c, renewURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("renewtoken for region %s, url(%s) error(%v)", target, renewURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("renewtoken for region %s, url(%s) error(%v)", target, renewURL+"?"+params.Encode(), err)
		return
	}
	expires := res.Expires
	if expires == 0 {
		expires = res.Data.Expires
	}
	renewToken = &model.RenewInfo{
		Expires: expires,
	}
	return
}

// GetCookieByMid get cookie by mid
func (d *Dao) GetCookieByMid(c context.Context, mid int64) (cookies *v1.CreateCookieReply, err error) {
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%d", mid))
	URL := d.c.Passport.Host["auth"] + _createCookieURI
	var res struct {
		Code int                   `json:"code"`
		Data *v1.CreateCookieReply `json:"data,omitempty"`
	}
	if err = d.client.Get(c, URL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("get cookies by mid(%d) from url(%s) error(%v)", mid, URL, err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("get cookies error, mid(%d), error(%v)", mid, err)
		return
	}
	cookies = res.Data
	return
}
