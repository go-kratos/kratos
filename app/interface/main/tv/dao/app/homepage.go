package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/interface/main/tv/conf"
	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// Zones represents the different zones to display on the homepage recommendation area
var Zones = []string{"jp", "cn", "tv", "movie", "documentary"}

// HeaderData gets the header data from PGC API
func (d *Dao) HeaderData(ctx context.Context, appInfo *conf.TVApp) (result map[string][]*model.Card, err error) {
	var (
		res = new(struct {
			Code   int                      `json:"code"`
			Result map[string][]*model.Card `json:"result"`
		})
		bangumiURL = d.conf.Host.APIIndex
		params     = url.Values{}
	)
	params.Set("build", appInfo.Build)
	params.Set("mobi_app", appInfo.MobiApp)
	params.Set("platform", appInfo.Platform)
	if err = d.client.Get(ctx, bangumiURL, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), bangumiURL+"?"+params.Encode())
		return
	}
	result = res.Result
	// check result contain the 5 type data
	for _, v := range Zones {
		_, ok := result[v]
		if !ok {
			err = fmt.Errorf("Result Miss Data: %s", v)
			return
		}
	}
	return
}

// FollowData gets the follow data from PGC API
func (d *Dao) FollowData(ctx context.Context, appInfo *conf.TVApp, accessKey string) (result []*model.Follow, err error) {
	var (
		bangumiURL = d.conf.Host.APIFollow
		params     = url.Values{}
		res        = model.ResFollow{}
	)
	params.Set("access_key", accessKey)
	params.Set("build", appInfo.Build)
	params.Set("mobi_app", appInfo.MobiApp)
	params.Set("platform", appInfo.Platform)
	params.Set("pagesize", fmt.Sprintf("%d", d.conf.Homepage.FollowSize))
	if err = d.client.Get(ctx, bangumiURL, "", params, &res); err != nil {
		log.Error("FollowData ERROR:%v", err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), bangumiURL+"?"+params.Encode())
		log.Error("FollowData ERROR:%v", err)
		return
	}
	result = res.Result
	return
}
