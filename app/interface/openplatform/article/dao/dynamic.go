package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
)

const _Dynamic = "http://api.vc.bilibili.co/dynamic_repost/v0/dynamic_repost/view_repost"

// DynamicCount get dynamic count from api
func (d *Dao) DynamicCount(c context.Context, aid int64) (count int64, err error) {
	params := url.Values{}
	params.Set("rid", strconv.FormatInt(aid, 10))
	params.Set("type", "64")
	params.Set("offset", "0")
	var res struct {
		Code int    `json:"errno"`
		Msg  string `json:"msg"`
		Data struct {
			TotalCount int64 `json:"total_count"`
		} `json:"data"`
	}
	err = d.httpClient.Get(c, _Dynamic, "", params, &res)
	if err != nil {
		PromError("count:dynamic")
		log.Error("dynamic: d.client.Get(%s) error(%+v)", _Dynamic+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		PromError("count:dynamic接口")
		log.Error("dynamic: url(%s) res code(%d) msg: %s", _Dynamic+"?"+params.Encode(), res.Code, res.Msg)
		err = ecode.Int(res.Code)
		return
	}
	count = res.Data.TotalCount
	return
}
