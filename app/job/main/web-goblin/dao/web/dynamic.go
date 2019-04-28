package web

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_opt      = "1002"
	_platform = "web"
)

// PushAll  broadcast push all
func (d *Dao) PushAll(c context.Context, msg string, ip string) (err error) {
	var res struct {
		Code int `json:"code"`
	}
	params := url.Values{}
	params.Set("operation", _opt)
	params.Set("speed", strconv.Itoa(d.c.Rule.BroadFeed))
	params.Set("platform", _platform)
	params.Set("message", msg)
	if err = d.http.Post(c, d.broadcastURL, ip, params, &res); err != nil {
		log.Error("PushAll url(%s) error(%v)", d.broadcastURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
	}
	return
}
