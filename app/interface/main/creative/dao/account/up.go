package account

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_picUpInfoURL   = "/link_draw_ex/v0/doc/check"
	_blinkUpInfoURL = "/clip_ext/v0/video/have"
	_upInfoURL      = "/x/internal/uper/info"
)

// Pic pic return value
type Pic struct {
	Has int `json:"has_doc"`
}

// Blink blink return value
type Blink struct {
	Has int `json:"has"`
}

// Pic get pic up info.
func (d *Dao) Pic(c context.Context, mid int64, ip string) (has int, err error) {
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data Pic `json:"data"`
	}
	err = d.fastClient.Get(c, d.picUpInfoURL, ip, params, &res)
	if err != nil {
		log.Error("d.fastClient.Get(%s) error(%v)", d.picUpInfoURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("Pic url(%s) error(%v)", d.picUpInfoURL+"?"+params.Encode(), err)
		err = ecode.Int(res.Code)
		return
	}
	has = res.Data.Has
	return
}

// Blink get BLink up info.
func (d *Dao) Blink(c context.Context, mid int64, ip string) (has int, err error) {
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int   `json:"code"`
		Data Blink `json:"data"`
	}
	err = d.fastClient.Get(c, d.blinkUpInfoURL, ip, params, &res)
	if err != nil {
		log.Error("d.fastClient.Get(%s) error(%v)", d.blinkUpInfoURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("Blink url(%s) error(%v)", d.blinkUpInfoURL+"?"+params.Encode(), err)
		err = ecode.Int(res.Code)
		return
	}
	has = res.Data.Has
	return
}
