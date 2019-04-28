package dao

import (
	"context"
	"net/url"

	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_msgURL = "/api/notify/send.user.notify.do"
)

// Message send message.
func (d *Dao) Message(c context.Context, title, msg string, mids []int64) (err error) {
	return d.RawMessage(c, "2_2_2", title, msg, mids)
}

// RawMessage send message with mc.
func (d *Dao) RawMessage(c context.Context, mc string, title, msg string, mids []int64) (err error) {
	params := url.Values{}
	params.Set("type", "json")
	params.Set("source", "2")
	params.Set("mc", mc)
	params.Set("title", title)
	params.Set("data_type", "4")
	params.Set("context", msg)
	params.Set("mid_list", xstr.JoinInts(mids))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpClient.Post(c, d.msgURL, "", params, &res); err != nil {
		err = errors.Wrap(err, "dao send message")
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrapf(err, "message send failed,mid(%v)", mids)
		return
	}
	log.Info("sendmessage mc:%s, mids:%v, title:%s, msg:%s", mc, mids, title, msg)
	return
}
