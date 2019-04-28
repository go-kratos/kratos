package wechat

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"go-common/app/interface/main/web-goblin/model/wechat"
	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_accessTokenKey = "a_t_k"
	_accessTokenURI = "/cgi-bin/token"
	_qrcodeURI      = "/wxa/getwxacodeunlimit"
)

// RawAccessToken get wechat access token.
func (d *Dao) RawAccessToken(c context.Context) (data *wechat.AccessToken, err error) {
	params := url.Values{}
	params.Set("grant_type", "client_credential")
	params.Set("appid", d.c.Wechat.AppID)
	params.Set("secret", d.c.Wechat.Secret)
	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, d.wxAccessTokenURL+"?"+params.Encode(), nil); err != nil {
		log.Error("AccessToken http.NewRequest error(%v)", err)
		return
	}
	var res struct {
		Errcode     int    `json:"errcode"`
		Errmsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIN   int64  `json:"expires_in"`
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("AccessToken d.httpClient.Do error(%v)", err)
		return
	}
	if res.Errcode != ecode.OK.Code() {
		log.Error("AccessToken errcode error(%d) msg(%s)", res.Errcode, res.Errmsg)
		err = ecode.RequestErr
		return
	}
	data = &wechat.AccessToken{AccessToken: res.AccessToken, ExpiresIn: res.ExpiresIN}
	return
}

// Qrcode get qrcode.
func (d *Dao) Qrcode(c context.Context, accessToken, arg string) (qrcode []byte, err error) {
	var (
		req     *http.Request
		bs      []byte
		jsonErr error
	)
	params := url.Values{}
	params.Set("access_token", accessToken)
	if req, err = http.NewRequest(http.MethodPost, d.wxQrcodeURL+"?"+params.Encode(), strings.NewReader(arg)); err != nil {
		log.Error("Qrcode http.NewRequest error(%v)", err)
		return
	}
	if bs, err = d.httpClient.Raw(c, req); err != nil {
		log.Error("Qrcode d.httpClient.Do error(%v)", err)
		return
	}
	var res struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}
	if jsonErr = json.Unmarshal(bs, &res); jsonErr == nil && res.Errcode != ecode.OK.Code() {
		log.Error("Qrcode errcode error(%d) msg(%s)", res.Errcode, res.Errmsg)
		err = ecode.RequestErr
		return
	}
	qrcode = bs
	return
}

// CacheAccessToken cache access token
func (d *Dao) CacheAccessToken(c context.Context) (data *wechat.AccessToken, err error) {
	var (
		value []byte
		key   = _accessTokenKey
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if value, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("CacheAccessToken conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	data = new(wechat.AccessToken)
	if err = json.Unmarshal(value, &data); err != nil {
		log.Error("CacheAccessToken json.Unmarshal(%v) error(%v)", value, err)
	}
	return
}

// AddCacheAccessToken add access token cache
func (d *Dao) AddCacheAccessToken(c context.Context, data *wechat.AccessToken) (err error) {
	var (
		bs   []byte
		key  = _accessTokenKey
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(data); err != nil {
		log.Error("AddCacheAccessToken json.Marshal(%v) error (%v)", data, err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET, %s, %s) error(%v)", key, string(bs), err)
		return
	}
	expire := data.ExpiresIn - 60
	if err = conn.Send("EXPIRE", key, expire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, expire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}
