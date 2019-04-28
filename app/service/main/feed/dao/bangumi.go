package dao

import (
	"context"
	"net/url"
	"strconv"

	feedmdl "go-common/app/service/main/feed/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_bangumiURL     = "http://bangumi.bilibili.co"
	_pullURL        = _bangumiURL + "/internal_api/follow_pull"
	_pullSeasonsURL = _bangumiURL + "/internal_api/follow_seasons"
)

// BangumiPull pull bangumi feed.
func (d *Dao) BangumiPull(c context.Context, mid int64, ip string) (seasonIDS []int64, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code   int             `json:"code"`
		Result []*feedmdl.Pull `json:"result"`
	}
	if err = d.httpClient.Get(c, _pullURL, ip, params, &res); err != nil {
		PromWarn("bangumi:Pull接口")
		log.Error("d.client.Get(%s) error(%v)", _pullURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		PromWarn("bangumi:Pull接口")
		log.Error("url(%s) res code(%d) or res.result(%v)", _pullURL+"?"+params.Encode(), res.Code, res.Result)
		err = ecode.Int(res.Code)
		return
	}
	for _, r := range res.Result {
		seasonIDS = append(seasonIDS, r.SeasonID)
	}
	return
}

// BangumiSeasons get bangumi info by seasonids.
func (d *Dao) BangumiSeasons(c context.Context, seasonIDs []int64, ip string) (psm map[int64]*feedmdl.Bangumi, err error) {
	if len(seasonIDs) == 0 {
		return
	}
	params := url.Values{}
	params.Set("season_ids", xstr.JoinInts(seasonIDs))
	var res struct {
		Code   int                `json:"code"`
		Result []*feedmdl.Bangumi `json:"result"`
	}
	if err = d.httpClient.Get(c, _pullSeasonsURL, ip, params, &res); err != nil {
		PromWarn("bangumi:详情接口")
		log.Error("d.client.Get(%s) error(%v)", _pullSeasonsURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		PromWarn("bangumi:详情接口")
		log.Error("url(%s) res code(%d) or res.result(%v)", _pullSeasonsURL+"?"+params.Encode(), res.Code, res.Result)
		err = ecode.Int(res.Code)
		return
	}
	psm = make(map[int64]*feedmdl.Bangumi, len(res.Result))
	for _, p := range res.Result {
		if p == nil {
			continue
		}
		psm[p.SeasonID] = p
	}
	return
}
