package dao

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/admin/main/reply/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_msgReplyDel     = "1_1_2"
	_msgReportAccept = "1_3_1"
	_msgTypeSystem   = 4

	// api
	_apiMsgSend = "http://message.bilibili.co/api/notify/send.user.notify.do"
)

// SendReplyDelMsg send delete reply message.
func (d *Dao) SendReplyDelMsg(c context.Context, mid int64, title, msg string, now time.Time) (err error) {
	return d.sendMsg(c, _msgReplyDel, title, msg, _msgTypeSystem, 0, []int64{mid}, "", now.Unix())
}

// SendReportAcceptMsg send report message.
func (d *Dao) SendReportAcceptMsg(c context.Context, mid int64, title, msg string, now time.Time) (err error) {
	return d.sendMsg(c, _msgReportAccept, title, msg, _msgTypeSystem, 0, []int64{mid}, "", now.Unix())
}

func (d *Dao) sendMsg(c context.Context, mc, title, msg string, typ int, pub int64, mids []int64, info string, ts int64) (err error) {
	params := url.Values{}
	params.Set("type", "json")
	params.Set("source", "1")
	params.Set("mc", mc)
	params.Set("title", title)
	params.Set("data_type", strconv.Itoa(typ))
	params.Set("context", msg)
	params.Set("mid_list", xstr.JoinInts(mids))
	params.Set("publisher", strconv.FormatInt(pub, 10))
	params.Set("ext_info", info)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpClient.Post(c, _apiMsgSend, "", params, &res); err != nil {
		log.Error("sendMsg error(%v) params(%v)", err, params)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = model.ErrMsgSend
		log.Error("sendMsg failed(%v) error(%v)", _apiMsgSend+"?"+params.Encode(), res.Code)
	}
	log.Info("sendMsg(mc:%s title:%s msg:%s type:%d mids:%v error(%v)", mc, title, msg, typ, mids, err)
	return
}
