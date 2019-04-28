package dao

import (
	"context"
	"net/url"

	"go-common/app/interface/main/tv/conf"
	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"fmt"

	"github.com/pkg/errors"
)

// RecomData gets the recom data from PGC API
func (d *Dao) RecomData(ctx context.Context, appInfo *conf.TVApp, sid string, stype string) (result []*model.Recom, err error) {
	var (
		bangumiURL = d.conf.Host.APIRecom
		params     = url.Values{}
		response   = &model.ResponseRecom{}
	)
	params.Set("season_id", sid)
	params.Set("season_type", stype)
	params.Set("build", appInfo.Build)
	params.Set("mobi_app", appInfo.MobiApp)
	params.Set("platform", appInfo.Platform)
	log.Info("[RecomData Request] URL: %s, Params: %s", bangumiURL, params.Encode())
	if err = d.client.Get(ctx, bangumiURL, "", params, response); err != nil {
		log.Error("RecomData ERROR:%v", err)
		return
	}
	if response.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(response.Code), fmt.Sprintf("Bangumi API Error %v", response.Message))
		log.Error("RecomData ERROR:%v, URL: %s", err, bangumiURL+"?"+params.Encode())
		return
	}
	if response.Result == nil {
		err = errors.Wrap(ecode.ServerErr, "bangumi api returns empty")
		log.Error("RecomData ERROR:%v", err)
		return
	}
	result = response.Result.List
	log.Info("[RecomData] For Sid: %s, Stype: %s, Get PGC data NB: %d", sid, stype, len(result))
	return
}
