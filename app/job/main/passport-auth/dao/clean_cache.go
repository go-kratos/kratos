package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
)

// AsoCleanCache aso clean cache
func (d *Dao) AsoCleanCache(c context.Context, token, session string, mid int64) (err error) {
	log.Info("aso clean cache,mid = %d,token = %s, session = %s", mid, token, session)
	params := url.Values{}
	params.Set("token", token)
	params.Set("session", session)
	params.Set("mid", strconv.Itoa(int(mid)))
	res := &struct {
		Code int `json:"code"`
	}{}
	if err = d.httpClient.Get(c, d.c.AuthJobConfig.AsoCleanURL, "127.0.0.2", params, res); err != nil {
		log.Error("AsoCleanCache HTTP request err %v,token = %s,session = %s,mid = %d", err, token, session, mid)
		return
	}
	if res.Code != 0 {
		log.Error("AsoCleanCache server err_code %d,token = %s,session = %s,mid=%d", res.Code, token, session, mid)
		if res.Code == ecode.RequestErr.Code() {
			err = nil
			return
		}
		err = ecode.Int(res.Code)
		return
	}
	return
}
