package bnj

import (
	"context"
	"net/url"

	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_opt        = "1004"
	_platform   = "web"
	_broadURL   = "/x/internal/broadcast/push/all"
	_messageURL = "/api/notify/send.user.notify.do"
	_notify     = "4"
)

// PushAll  broadcast push all
func (d *Dao) PushAll(c context.Context, msg string) (err error) {
	params := url.Values{}
	params.Set("operation", _opt)
	params.Set("platform", _platform)
	params.Set("message", msg)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.broadcastURL, "", params, &res); err != nil {
		log.Error("PushAll url(%s) error(%v)", d.broadcastURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
	}
	return
}

// SendMessage send system notify.
func (d *Dao) SendMessage(c context.Context, mids []int64, mc, title, msg string) (err error) {
	params := url.Values{}
	params.Set("mid_list", xstr.JoinInts(mids))
	params.Set("title", title)
	params.Set("mc", mc)
	params.Set("data_type", _notify)
	params.Set("context", msg)
	var res struct {
		Code int `json:"code"`
	}
	err = d.client.Post(c, d.messageURL, "", params, &res)
	if err != nil {
		log.Error("SendMessage d.client.Post(%s) error(%+v)", d.messageURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("SendMessage url(%s) res code(%d)", d.messageURL+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
	}
	return
}
