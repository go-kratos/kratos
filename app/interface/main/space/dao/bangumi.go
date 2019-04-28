package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/space/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_build               = "0"
	_platform            = "web"
	_bangumiURI          = "/api/get_concerned_season"
	_bangumiConcernURI   = "/api/concern_season"
	_bangumiUnConcernURI = "/api/unconcern_season"
)

// BangumiList get bangumi sub list by mid.
func (d *Dao) BangumiList(c context.Context, mid int64, pn, ps int) (data []*model.Bangumi, count int, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("build", _build)
	params.Set("platform", _platform)
	var res struct {
		Code   int              `json:"code"`
		Count  string           `json:"count"`
		Result []*model.Bangumi `json:"result"`
	}
	if err = d.httpR.Get(c, d.bangumiURL, ip, params, &res); err != nil {
		log.Error("d.httpR.Get(%s,%d) error(%v)", d.bangumiURL, mid, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpR.Get(%s,%d) error(%v)", d.bangumiURL, mid, err)
		err = ecode.Int(res.Code)
		return
	}
	data = res.Result
	count, _ = strconv.Atoi(res.Count)
	return
}

// BangumiConcern bangumi concern.
func (d *Dao) BangumiConcern(c context.Context, mid, seasonID int64) (err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("season_id", strconv.FormatInt(seasonID, 10))
	params.Set("build", _build)
	params.Set("platform", _platform)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpW.Post(c, d.bangumiConcernURL, ip, params, &res); err != nil {
		log.Error("d.httpW.Post(%s,%d) error(%v)", d.bangumiConcernURL, mid, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpW.Post(%s,%d) error(%v)", d.bangumiConcernURL, mid, err)
		err = ecode.Int(res.Code)
	}
	return
}

// BangumiUnConcern bangumi cancel sub.
func (d *Dao) BangumiUnConcern(c context.Context, mid, seasonID int64) (err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("season_id", strconv.FormatInt(seasonID, 10))
	params.Set("build", _build)
	params.Set("platform", _platform)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpW.Post(c, d.bangumiUnConcernURL, ip, params, &res); err != nil {
		log.Error("d.httpW.Post(%s,%d) error(%v)", d.bangumiUnConcernURL, mid, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpW.Post(%s,%d) error(%v)", d.bangumiUnConcernURL, mid, err)
		err = ecode.Int(res.Code)
	}
	return
}
