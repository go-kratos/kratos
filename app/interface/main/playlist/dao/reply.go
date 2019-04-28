package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
)

const (
	_replyType  = "18"
	_replyState = "0" // 0: open, 1: close
	_replyURL   = "http://api.bilibili.co/x/internal/v2/reply/subject/regist"
)

// RegReply opens playlist's reply.
func (d *Dao) RegReply(c context.Context, pid, mid int64) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("oid", strconv.FormatInt(pid, 10))
	params.Set("type", _replyType)
	params.Set("state", _replyState)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.http.Post(c, _replyURL, "", params, &res); err != nil {
		PromError("reply:打开评论", "d.http.Post(%s) error(%v)", _replyURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != ecode.OK.Code() {
		PromError("reply:打开评论状态码异常", "d.http.Post(%s) error(%v)", _replyURL+"?"+params.Encode(), err)
		err = ecode.Int(res.Code)
	}
	return
}
