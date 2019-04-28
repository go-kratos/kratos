package dao

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_replyType       = "12"
	_replyStateOpen  = "0" // 0: open, 1: close
	_replyStateClose = "1"
	_replyURL        = "http://api.bilibili.co/x/internal/v2/reply/subject/regist"
)

// OpenReply opens article's reply.
func (d *Dao) OpenReply(c context.Context, aid, mid int64) (err error) {
	defer func() {
		if err == nil {
			return
		}
		time.Sleep(time.Second)
		if e := d.PushReply(c, aid, mid); e != nil {
			log.Error("d.PushReply(%d,%d) error(%+v)", aid, mid, e)
		}
	}()
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("oid", strconv.FormatInt(aid, 10))
	params.Set("type", _replyType)
	params.Set("state", _replyStateOpen)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpClient.Post(c, _replyURL, "", params, &res); err != nil {
		log.Error("d.httpClient.Post(%s) error(%+v)", _replyURL+"?"+params.Encode(), err)
		PromError("reply:打开评论")
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpClient.Post(%s) code(%d)", _replyURL+"?"+params.Encode(), res.Code)
		PromError("reply:打开评论状态码异常")
		err = ecode.Int(res.Code)
	}
	return
}

// CloseReply close article's reply.
func (d *Dao) CloseReply(c context.Context, aid, mid int64) (err error) {
	defer func() {
		if err == nil {
			return
		}
		time.Sleep(time.Second)
		if e := d.PushReply(c, aid, mid); e != nil {
			log.Error("d.PushReply(%d,%d) error(%+v)", aid, mid, e)
		}
	}()
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("oid", strconv.FormatInt(aid, 10))
	params.Set("type", _replyType)
	params.Set("state", _replyStateClose)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpClient.Post(c, _replyURL, "", params, &res); err != nil {
		log.Error("d.httpClient.Post(%s) error(%+v)", _replyURL+"?"+params.Encode(), err)
		PromError("reply:打开评论")
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpClient.Post(%s) code(%d)", _replyURL+"?"+params.Encode(), res.Code)
		PromError("reply:打开评论状态码异常")
		err = ecode.Int(res.Code)
	}
	return
}
