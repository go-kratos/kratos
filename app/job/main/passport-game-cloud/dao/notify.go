package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
)

// NotifyGame to notify game.
func (d *Dao) NotifyGame(c context.Context, mid int64, accessToken, action string) (err error) {
	params := url.Values{}
	params.Set("modifiedAttr", action)
	params.Set("mid", strconv.FormatInt(mid, 10))
	if accessToken != "" {
		params.Set("access_token", accessToken)
	}
	params.Set("from", "passport-game-cloud-job")
	var res struct {
		Code int `json:"code"`
	}
	if err = d.gameClient.Get(c, d.delGameCacheURI, "127.0.0.1", params, &res); err != nil {
		log.Error("failed to notify game, d.gameClient.Get(%s) error(%v)", d.delGameCacheURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("failed to notify game, url(%s) err(%v)", d.delGameCacheURI+"?"+params.Encode(), err)
	}
	return
}
