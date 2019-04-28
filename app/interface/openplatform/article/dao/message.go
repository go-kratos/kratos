package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
)

var _notify = "4"

// SendMessage .
func (d *Dao) SendMessage(c context.Context, tid, mid, aid int64, title, msg string) (err error) {
	params := url.Values{}
	params.Set("mid_list", strconv.FormatInt(mid, 10))
	params.Set("title", title)
	params.Set("mc", d.c.Message.MC)
	params.Set("data_type", _notify)
	params.Set("context", msg)
	params.Set("notify_type", strconv.FormatInt(tid, 10))
	params.Set("res_id", strconv.FormatInt(aid, 10))
	var res struct {
		Code int `json:"code"`
	}
	err = d.messageHTTPClient.Post(c, d.c.Message.URL, "", params, &res)
	if err != nil {
		PromError("message:send接口")
		log.Error("d.client.Post(%s) error(%+v)", d.c.Message.URL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		PromError("message:send接口")
		log.Error("url(%s) res code(%d)", d.c.Message.URL+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
		return
	}
	log.Info("发送点赞消息通知 (%s) error(%+v)", d.c.Message.URL+"?"+params.Encode(), err)
	return
}
