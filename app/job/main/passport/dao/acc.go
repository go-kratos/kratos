package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/job/main/passport/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// SetToken set token via passport api.
func (d *Dao) SetToken(c context.Context, t *model.Token) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(t.Mid, 10))
	params.Set("appid", strconv.FormatInt(t.Appid, 10))
	params.Set("appSubid", strconv.FormatInt(t.Subid, 10))
	params.Set("accessToken", t.Token)
	params.Set("refreshToken", t.RToken)
	params.Set("tp", strconv.FormatInt(t.Type, 10))
	params.Set("createAt", strconv.FormatInt(t.CTime, 10))
	params.Set("expires", strconv.FormatInt(t.Expires, 10))
	params.Set("from", "passport-job")
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(context.TODO(), d.setTokenURI, "127.0.0.1", params, &res); err != nil {
		log.Error("d.client.Get error(%v)", err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("set token url(%s) err(%v)", d.setTokenURI+"?"+params.Encode(), err)
	}
	return
}

// DelCache del cache via passport api.
func (d *Dao) DelCache(c context.Context, accessKey string) (err error) {
	params := url.Values{}
	params.Set("access_key", accessKey)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Get(context.TODO(), d.delCacheURI, "", params, &res); err != nil {
		log.Error("d.client.Get url(%s) error(%v)", d.delCacheURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("del cache url(%s) error(%v)", d.delCacheURI+"?"+params.Encode(), err)
	}
	return
}

// NotifyGame notify game.
func (d *Dao) NotifyGame(token *model.AccessInfo, action string) (err error) {
	params := url.Values{}
	params.Set("modifiedAttr", action)
	params.Set("mid", strconv.FormatInt(token.Mid, 10))
	params.Set("access_token", token.Token)
	params.Set("from", "passport-job")
	var res struct {
		Code int `json:"code"`
	}
	if err = d.gameClient.Get(context.TODO(), d.delGameCacheURI, "127.0.0.1", params, &res); err != nil {
		log.Error("d.client.Get url(%s) error(%v)", d.delGameCacheURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("url(%s) err(%v)", d.delGameCacheURI+"?"+params.Encode(), err)
	}
	return
}
