package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_replyState   = "0" // 0: open, 1: close
	_replyURL     = "http://api.bilibili.co/x/internal/v2/reply/subject/regist"
	_gameOfficial = 32708316
)

// RegReply opens eports's reply.
func (d *Dao) RegReply(c context.Context, maid, adid int64, replyType string) (err error) {
	params := url.Values{}
	params.Set("adid", strconv.FormatInt(adid, 10))
	params.Set("mid", strconv.FormatInt(_gameOfficial, 10))
	params.Set("oid", strconv.FormatInt(maid, 10))
	params.Set("type", replyType)
	params.Set("state", _replyState)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.replyClient.Post(c, _replyURL, "", params, &res); err != nil {
		log.Error("d.replyClient.Post(%s) error(%v)", _replyURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.replyClient.Post(%s) error(%v)", _replyURL+"?"+params.Encode(), err)
		err = ecode.Int(res.Code)
	}
	return
}
