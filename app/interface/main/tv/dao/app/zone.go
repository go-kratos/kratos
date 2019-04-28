package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/interface/main/tv/conf"
	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// ChannelData gets the header data from PGC API
func (d *Dao) ChannelData(c context.Context, seasonType int, appInfo *conf.TVApp) (result []*model.Card, err error) {
	var res struct {
		Code   int           `json:"code"`
		Result []*model.Card `json:"result"`
	}
	bangumiURL := d.conf.Host.APIZone
	params := url.Values{}
	params.Set("build", appInfo.Build)
	params.Set("mobi_app", appInfo.MobiApp)
	params.Set("platform", appInfo.Platform)
	params.Set("season_type", strconv.Itoa(int(seasonType)))
	if err = d.client.Get(c, bangumiURL, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), bangumiURL+"?"+params.Encode())
		return
	}
	if len(res.Result) == 0 {
		err = ecode.TvPGCRankEmpty
		log.Error("[LoadPGCList] Zone %d, Err %v", seasonType, err)
		return
	}
	for _, v := range res.Result {
		if v.NewEP != nil {
			v.BePGC()
			result = append(result, v)
		}
	}
	if len(result) == 0 {
		err = ecode.TvPGCRankNewEPNil
		log.Error("[LoadPGCList] Zone %d, Err %v", seasonType, err)
	}
	return
}

// UgcAIData gets the ugc types rank data from AI
func (d *Dao) UgcAIData(c context.Context, tid int16) (result []*model.AIData, err error) {
	var (
		res    model.RespAI
		AIURL  = fmt.Sprintf(d.conf.Host.AIUgcType, tid)
		params = url.Values{}
	)
	if err = d.client.Get(c, AIURL, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), AIURL+"?"+params.Encode())
		return
	}
	result = res.List
	return
}
