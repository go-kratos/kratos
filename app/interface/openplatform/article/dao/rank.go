package dao

import (
	"context"
	"net/url"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	_monthURL           = "/data/rank/article/all-30.json"
	_weekURL            = "/data/rank/article/all-7.json"
	_yesterDayURL       = "/data/rank/article/all-1.json"
	_beforeYesterDayURL = "/data/rank/article/all-2.json"
)

// Rank get rank from bigdata
func (d *Dao) Rank(c context.Context, cid int64, ip string) (res model.RankResp, err error) {
	var addr string
	switch cid {
	case model.RankMonth:
		addr = _monthURL
	case model.RankWeek:
		addr = _weekURL
	case model.RankYesterday:
		addr = _yesterDayURL
	case model.RankBeforeYesterday:
		addr = _beforeYesterDayURL
	default:
		err = ecode.RequestErr
		return
	}
	params := url.Values{}
	var resp struct {
		Code int `json:"code"`
		model.RankResp
	}
	if err = d.httpClient.Get(c, d.c.Article.RankHost+addr, ip, params, &resp); err != nil {
		PromError("rank:rank接口")
		log.Error("d.client.Get(%s) error(%+v)", addr+"?"+params.Encode(), err)
		return
	}
	if resp.Code != ecode.OK.Code() {
		PromError("rank:rank接口")
		log.Error("url(%s) res code(%d) or res.result(%+v)", addr+"?"+params.Encode(), resp.Code, resp)
		err = ecode.Int(resp.Code)
		return
	}
	res = resp.RankResp
	return
}
