package relation

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/account/conf"
	"go-common/app/interface/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// Dao is
type Dao struct {
	conf         *conf.Config
	httpClient   *bm.Client
	recommendURL string
}

const (
	_recommendURI = "/main/recommend"
)

// New is
func New(conf *conf.Config) *Dao {
	return &Dao{
		conf:         conf,
		httpClient:   bm.NewClient(conf.HTTPClient.Normal),
		recommendURL: conf.Host.Search + _recommendURI,
	}
}

func paltform(device *bm.Device) string {
	if device.RawMobiApp == "" {
		return "web"
	}
	if device.IsAndroid() {
		return "android"
	}
	if device.IsIOS() {
		return "ios"
	}
	return "web"
}

func buvid(device *bm.Device) string {
	if device.RawMobiApp == "" {
		return device.Buvid3
	}
	return device.Buvid
}

// Recommend is
func (d *Dao) Recommend(ctx context.Context, mid int64, serviceArea string, mainTids string, subTids string, device *bm.Device, pagesize int64, ip string) (*model.RecommendResponse, error) {
	params := url.Values{}
	params.Set("platform", paltform(device))
	params.Set("mobi_app", device.MobiApp())
	params.Set("device", device.Device)
	params.Set("build", strconv.FormatInt(device.Build, 10))
	params.Set("clientip", ip)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("buvid", buvid(device))
	params.Set("pagesize", strconv.FormatInt(pagesize, 10))
	params.Set("rec_type", "up_rec")
	params.Set("service_area", serviceArea)
	if mainTids != "" {
		params.Set("main_tids", mainTids)
	}
	if subTids != "" {
		params.Set("sub_tids", subTids)
	}

	resp := &model.RecommendResponse{}
	if err := d.httpClient.Get(ctx, d.recommendURL, ip, params, resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		log.Error("Failed to call recommendation service with error code: %d, params: %+v, trackid(%s)", resp.Code, params, resp.TrackID)
		return nil, ecode.Int(int(resp.Code))
	}
	return resp, nil
}

// TagSuggestRecommend is
func (d *Dao) TagSuggestRecommend(ctx context.Context, mid int64, contextID string, tagname string, device *bm.Device, pagesize int64, ip string) (*model.TagSuggestRecommendResponse, error) {
	params := url.Values{}
	params.Set("platform", paltform(device))
	params.Set("mobi_app", device.MobiApp())
	params.Set("device", device.Device)
	params.Set("build", strconv.FormatInt(device.Build, 10))
	params.Set("clientip", ip)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("buvid", buvid(device))
	params.Set("pagesize", strconv.FormatInt(pagesize, 10))
	params.Set("rec_type", "tagup_rec")
	params.Set("service_area", "tag_suggest")
	params.Set("context_id", contextID)
	if tagname != "" {
		params.Set("tagname", tagname)
	}
	resp := &model.TagSuggestRecommendResponse{}
	if err := d.httpClient.Get(ctx, d.recommendURL, ip, params, resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		log.Error("Failed to call recommendation service with error code: %d, params: %+v, trackid(%s)", resp.Code, params, resp.TrackID)
		return nil, ecode.Int(int(resp.Code))
	}
	return resp, nil
}
