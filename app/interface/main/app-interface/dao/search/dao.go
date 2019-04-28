package search

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-interface/conf"
	arcdao "go-common/app/interface/main/app-interface/dao/archive"
	bangumidao "go-common/app/interface/main/app-interface/dao/bangumi"
	livedao "go-common/app/interface/main/app-interface/dao/live"
	"go-common/app/interface/main/app-interface/model"
	bangumimdl "go-common/app/interface/main/app-interface/model/bangumi"
	livemdl "go-common/app/interface/main/app-interface/model/live"
	"go-common/app/interface/main/app-interface/model/search"
	"go-common/app/service/main/archive/api"
	locmdl "go-common/app/service/main/location/model"
	locrpc "go-common/app/service/main/location/rpc/client"
	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"

	"github.com/pkg/errors"
)

const (
	_main         = "/main/search"
	_suggest      = "/main/suggest"
	_hot          = "/main/hotword"
	_defaultWords = "/widget/getSearchDefaultWords"
	_rcmd         = "/query/recommend"
	_rcmdNoResult = "/search/recommend"
	_suggest3     = "/main/suggest/new"
	_pre          = "/search/frontpage"
)

const (
	_oldAndroid = 514000
	_oldIOS     = 6090

	_searchCodeLimitAndroid = 5250001
	_searchCodeLimitIPhone  = 6680
)

// Dao is search dao
type Dao struct {
	c            *conf.Config
	arcDao       *arcdao.Dao
	liveDao      *livedao.Dao
	bangumiDao   *bangumidao.Dao
	client       *httpx.Client
	client2      *httpx.Client
	main         string
	suggest      string
	hot          string
	defaultWords string
	rcmd         string
	rcmdNoResult string
	suggest3     string
	pre          string
	// rpc
	locRPC *locrpc.Service
}

// New initial search dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:            c,
		arcDao:       arcdao.New(c),
		liveDao:      livedao.New(c),
		bangumiDao:   bangumidao.New(c),
		client:       httpx.NewClient(c.HTTPSearch),
		client2:      httpx.NewClient(c.HTTPClient),
		main:         c.Host.Search + _main,
		suggest:      c.Host.Search + _suggest,
		hot:          c.Host.Search + _hot,
		defaultWords: c.Host.WWW + _defaultWords,
		rcmd:         c.Host.Search + _rcmd,
		rcmdNoResult: c.Host.Search + _rcmdNoResult,
		suggest3:     c.Host.Search + _suggest3,
		pre:          c.Host.Search + _pre,
		locRPC:       locrpc.New(c.LocationRPC),
	}
	return
}

// Search app all search .
func (d *Dao) Search(c context.Context, mid, zoneid int64, mobiApp, device, platform, buvid, keyword, duration, order, filtered, fromSource, recommend, parent string, plat int8, seasonNum, movieNum, upUserNum, uvLimit, userNum, userVideoLimit, biliUserNum, biliUserVideoLimit, rid, highlight, build, pn, ps, isQuery int, old bool, now time.Time, newPGC, flow bool) (res *search.Search, code int, err error) {
	var (
		req    *http.Request
		ip     = metadata.String(c, metadata.RemoteIP)
		ipInfo *locmdl.Info
	)
	if ipInfo, err = d.locRPC.Info(c, &locmdl.ArgIP{IP: ip}); err != nil {
		log.Warn("%v", err)
		err = nil
	}
	res = &search.Search{}
	params := url.Values{}
	params.Set("build", strconv.Itoa(build))
	params.Set("keyword", keyword)
	params.Set("main_ver", "v3")
	params.Set("highlight", strconv.Itoa(highlight))
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("tids", strconv.Itoa(rid))
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	if newPGC {
		params.Set("media_bangumi_num", strconv.Itoa(seasonNum))
	} else {
		params.Set("bangumi_num", strconv.Itoa(seasonNum))
		params.Set("smerge", "1")
	}
	if (plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD) {
		params.Set("is_new_user", "1")
	} else {
		if old {
			params.Set("upuser_num", strconv.Itoa(upUserNum))
			params.Set("uv_limit", strconv.Itoa(uvLimit))
		} else {
			params.Set("bili_user_num", strconv.Itoa(biliUserNum))
			params.Set("bili_user_vl", strconv.Itoa(biliUserVideoLimit))
		}
		params.Set("user_num", strconv.Itoa(userNum))
		params.Set("user_video_limit", strconv.Itoa(userVideoLimit))
		params.Set("query_rec_need", recommend)
	}
	params.Set("platform", platform)
	params.Set("duration", duration)
	params.Set("order", order)
	params.Set("search_type", "all")
	params.Set("from_source", fromSource)
	if filtered == "1" {
		params.Set("filtered", filtered)
	}
	if model.IsOverseas(plat) {
		params.Set("use_area", "1")
		params.Set("zone_id", strconv.FormatInt(zoneid, 10))
	} else if newPGC {
		params.Set("media_ft_num", strconv.Itoa(movieNum))
		params.Set("is_new_pgc", "1")
	} else {
		params.Set("movie_num", strconv.Itoa(movieNum))
	}
	if flow {
		params.Set("flow_need", "1")
	}
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.ComicAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.ComicIOS) || model.IsIPhoneB(plat) {
		params.Set("is_comic", "1")
	}
	if (model.IsAndroid(plat) && build > search.SearchLiveAllAndroid) || (model.IsIPhone(plat) && build > search.SearchLiveAllIOS) || (plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD) || model.IsIPhoneB(plat) {
		params.Set("new_live", "1")
	}
	if (model.IsAndroid(plat) && build > search.SearchTwitterAndroid) || (model.IsIPhone(plat) && build > search.SearchTwitterIOS) || model.IsIPhoneB(plat) {
		params.Set("is_twitter", "1")
	}
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.PGCHighLightAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.PGCHighLightIOS) {
		params.Set("app_highlight", "media_bangumi,media_ft")
	}
	if (model.IsAndroid(plat) && build >= search.SearchConvergeAndroid) || (model.IsIPhone(plat) && build >= search.SearchConvergeIOS) || model.IsIPhoneB(plat) {
		params.Set("card_num", "1")
	}
	params.Set("is_parents", parent)
	if (model.IsAndroid(plat) && build > search.SearchStarAndroid) || (model.IsIPhone(plat) && build > search.SearchStarIOS) || model.IsIPhoneB(plat) {
		params.Set("is_star", "1")
	}
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.PGCALLAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.PGCALLIOS) || model.IsIPhoneB(plat) {
		params.Set("is_pgc_all", "1")
	}
	if (model.IsAndroid(plat) && build > search.SearchTicketAndroid) || (model.IsIPhone(plat) && build > search.SearchTicketIOS) || model.IsIPhoneB(plat) {
		params.Set("is_ticket", "1")
	}
	if (model.IsAndroid(plat) && build > search.SearchProductAndroid) || (model.IsIPhone(plat) && build > search.SearchProductIOS) || model.IsIPhoneB(plat) {
		params.Set("is_product", "1")
	}
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.SpecialerGuideAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.SpecialerGuideIOS) {
		params.Set("is_special_guide", "1")
	}
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.SearchArticleAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.SearchArticleIOS) || model.IsIPhoneB(plat) {
		params.Set("is_article", "1")
	}
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.ChannelAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.ChannelIOS) {
		params.Set("is_tag", "1")
	}
	params.Set("is_org_query", strconv.Itoa(isQuery))
	if ipInfo != nil {
		params.Set("user_city", ipInfo.City)
	}
	// new request
	if req, err = d.client.NewRequest("GET", d.main, ip, params); err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	if err = d.client.Do(c, req, res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if (model.IsAndroid(plat) && build > _searchCodeLimitAndroid) || (model.IsIPhone(plat) && build > _searchCodeLimitIPhone) || (plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD) || model.IsIPhoneB(plat) {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		} else if res.Code != model.ForbidCode && res.Code != model.NoResultCode {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		}
	}
	for _, flow := range res.FlowResult {
		flow.Change()
	}
	code = res.Code
	return
}

// Season search season data.
func (d *Dao) Season(c context.Context, mid, zoneID int64, keyword, mobiApp, device, platform, buvid, filtered string, plat int8, build, pn, ps int, now time.Time) (st *search.TypeSearch, err error) {
	var (
		req *http.Request
		ip  = metadata.String(c, metadata.RemoteIP)
	)
	params := url.Values{}
	params.Set("main_ver", "v3")
	params.Set("platform", platform)
	params.Set("build", strconv.Itoa(build))
	params.Set("keyword", keyword)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("func", "search")
	params.Set("smerge", "1")
	params.Set("search_type", "bangumi")
	params.Set("source_type", "0")
	if filtered == "1" {
		params.Set("filtered", filtered)
	}
	if model.IsOverseas(plat) {
		params.Set("use_area", "1")
		params.Set("zone_id", strconv.FormatInt(zoneID, 10))
	}
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.SpecialerGuideAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.SpecialerGuideIOS) {
		params.Set("is_special_guide", "1")
	}
	// new request
	if req, err = d.client.NewRequest("GET", d.main, ip, params); err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	var res struct {
		Code  int               `json:"code"`
		SeID  string            `json:"seid"`
		Pages int               `json:"numPages"`
		List  []*search.Bangumi `json:"result"`
	}
	// do
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if (model.IsAndroid(plat) && build > _searchCodeLimitAndroid) || (model.IsIPhone(plat) && build > _searchCodeLimitIPhone) || (plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD) || model.IsIPhoneB(plat) {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		} else if res.Code != model.ForbidCode && res.Code != model.NoResultCode {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		}
		return
	}
	items := make([]*search.Item, 0, len(res.List))
	for _, v := range res.List {
		si := &search.Item{}
		if (model.IsAndroid(plat) && build <= _oldAndroid) || (model.IsIPhone(plat) && build <= _oldIOS) {
			si.FromSeason(v, model.GotoBangumi)
		} else {
			si.FromSeason(v, model.GotoBangumiWeb)
		}
		items = append(items, si)
	}
	st = &search.TypeSearch{TrackID: res.SeID, Pages: res.Pages, Items: items}
	return
}

// Upper search upper data.
func (d *Dao) Upper(c context.Context, mid int64, keyword, mobiApp, device, platform, buvid, filtered, order string, biliUserVL, highlight, build, userType, orderSort, pn, ps int, old bool, now time.Time) (st *search.TypeSearch, err error) {
	var (
		req     *http.Request
		plat    = model.Plat(mobiApp, device)
		avids   []int64
		avm     map[int64]*api.Arc
		roomIDs []int64
		lm      map[int64]*livemdl.RoomInfo
		ip      = metadata.String(c, metadata.RemoteIP)
	)
	params := url.Values{}
	params.Set("main_ver", "v3")
	params.Set("keyword", keyword)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("highlight", strconv.Itoa(highlight))
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("func", "search")
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("smerge", "1")
	params.Set("platform", platform)
	params.Set("build", strconv.Itoa(build))
	if old {
		params.Set("search_type", "upuser")
	} else {
		params.Set("search_type", "bili_user")
		params.Set("bili_user_vl", strconv.Itoa(biliUserVL))
		params.Set("user_type", strconv.Itoa(userType))
		params.Set("order_sort", strconv.Itoa(orderSort))
	}
	params.Set("order", order)
	params.Set("source_type", "0")
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.SpecialerGuideAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.SpecialerGuideIOS) {
		params.Set("is_special_guide", "1")
	}
	if filtered == "1" {
		params.Set("filtered", filtered)
	}
	// new request
	if req, err = d.client.NewRequest("GET", d.main, ip, params); err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	var res struct {
		Code  int            `json:"code"`
		SeID  string         `json:"seid"`
		Pages int            `json:"numPages"`
		List  []*search.User `json:"result"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if (model.IsAndroid(plat) && build > _searchCodeLimitAndroid) || (model.IsIPhone(plat) && build > _searchCodeLimitIPhone) || (plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD) || model.IsIPhoneB(plat) {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		} else if res.Code != model.ForbidCode && res.Code != model.NoResultCode {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		}
		return
	}
	items := make([]*search.Item, 0, len(res.List))
	for _, v := range res.List {
		for _, vr := range v.Res {
			avids = append(avids, vr.Aid)
		}
		roomIDs = append(roomIDs, v.RoomID)
	}
	g, ctx := errgroup.WithContext(c)
	if len(avids) != 0 {
		g.Go(func() (err error) {
			if avm, err = d.arcDao.Archives2(ctx, avids); err != nil {
				log.Error("Upper %+v", err)
				err = nil
			}
			return
		})
	}
	if len(roomIDs) > 0 {
		g.Go(func() (err error) {
			if lm, err = d.liveDao.LiveByRIDs(ctx, roomIDs); err != nil {
				log.Error("Upper %+v", err)
				err = nil
			}
			return
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
		return
	}
	for _, v := range res.List {
		si := &search.Item{}
		si.FromUpUser(v, avm, lm[v.RoomID])
		items = append(items, si)
	}
	st = &search.TypeSearch{TrackID: res.SeID, Pages: res.Pages, Items: items}
	return
}

// MovieByType search movie data from api .
func (d *Dao) MovieByType(c context.Context, mid, zoneid int64, keyword, mobiApp, device, platform, buvid, filtered string, plat int8, build, pn, ps int, now time.Time) (st *search.TypeSearch, err error) {
	var (
		req   *http.Request
		avids []int64
		avm   map[int64]*api.Arc
		ip    = metadata.String(c, metadata.RemoteIP)
	)
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("platform", platform)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("build", strconv.Itoa(build))
	params.Set("main_ver", "v3")
	params.Set("search_type", "pgc")
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("order", "totalrank")
	if filtered == "1" {
		params.Set("filtered", filtered)
	}
	if model.IsOverseas(plat) {
		params.Set("use_area", "1")
		params.Set("zone_id", strconv.FormatInt(zoneid, 10))
	}
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.SpecialerGuideAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.SpecialerGuideIOS) {
		params.Set("is_special_guide", "1")
	}
	// new request
	if req, err = d.client.NewRequest("GET", d.main, ip, params); err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	var res struct {
		Code  int             `json:"code"`
		SeID  string          `json:"seid"`
		Pages int             `json:"numPages"`
		List  []*search.Movie `json:"result"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if (model.IsAndroid(plat) && build > _searchCodeLimitAndroid) || (model.IsIPhone(plat) && build > _searchCodeLimitIPhone) || (plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD) || model.IsIPhoneB(plat) {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		} else if res.Code != model.ForbidCode && res.Code != model.NoResultCode {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		}
		return
	}
	items := make([]*search.Item, 0, len(res.List))
	for _, v := range res.List {
		if v.Type == "movie" {
			avids = append(avids, v.Aid)
		}
	}
	if len(avids) != 0 {
		if avm, err = d.arcDao.Archives2(c, avids); err != nil {
			log.Error("RecommendNoResult %+v", err)
			err = nil
		}
	}
	for _, v := range res.List {
		si := &search.Item{}
		si.FromMovie(v, avm)
		items = append(items, si)
	}
	st = &search.TypeSearch{TrackID: res.SeID, Pages: res.Pages, Items: items}
	return
}

// LiveByType search by diff type
func (d *Dao) LiveByType(c context.Context, mid, zoneid int64, keyword, mobiApp, device, platform, buvid, filtered, order, sType string, plat int8, build, pn, ps int, now time.Time) (st *search.TypeSearch, err error) {
	var (
		req *http.Request
		ip  = metadata.String(c, metadata.RemoteIP)
		// roomIDs []int64
		// lm      map[int64]*livemdl.RoomInfo
	)
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("platform", platform)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("build", strconv.Itoa(build))
	params.Set("main_ver", "v3")
	params.Set("search_type", sType)
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("order", order)
	if filtered == "1" {
		params.Set("filtered", filtered)
	}
	if model.IsOverseas(plat) {
		params.Set("use_area", "1")
		params.Set("zone_id", strconv.FormatInt(zoneid, 10))
	}
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.SpecialerGuideAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.SpecialerGuideIOS) {
		params.Set("is_special_guide", "1")
	}
	// new request
	if req, err = d.client.NewRequest("GET", d.main, ip, params); err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	var res struct {
		Code  int            `json:"code"`
		SeID  string         `json:"seid"`
		Pages int            `json:"numPages"`
		List  []*search.Live `json:"result"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if (model.IsAndroid(plat) && build > _searchCodeLimitAndroid) || (model.IsIPhone(plat) && build > _searchCodeLimitIPhone) || (plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD) || model.IsIPhoneB(plat) {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		} else if res.Code != model.ForbidCode && res.Code != model.NoResultCode {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		}
		return
	}
	// for _, v := range res.List {
	// 	roomIDs = append(roomIDs, v.RoomID)
	// }
	// if len(roomIDs) > 0 {
	// 	if lm, err = d.liveDao.LiveByRIDs(c, roomIDs, ip); err != nil {
	// 		log.Error("LiveByType %+v", err)
	// 		err = nil
	// 	}
	// }
	items := make([]*search.Item, 0, len(res.List))
	for _, v := range res.List {
		si := &search.Item{}
		si.FromLive(v, nil)
		items = append(items, si)
	}
	st = &search.TypeSearch{TrackID: res.SeID, Pages: res.Pages, Items: items}
	return
}

// Live search for live
func (d *Dao) Live(c context.Context, mid int64, keyword, mobiApp, platform, buvid, device, order, sType string, build, pn, ps int) (st *search.TypeSearch, err error) {
	var (
		req     *http.Request
		plat    = model.Plat(mobiApp, device)
		roomIDs []int64
		lm      map[int64]*livemdl.RoomInfo
		ip      = metadata.String(c, metadata.RemoteIP)
	)
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("mobi_app", mobiApp)
	params.Set("platform", platform)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("build", strconv.Itoa(build))
	params.Set("main_ver", "v3")
	params.Set("search_type", sType)
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("device", device)
	params.Set("order", order)
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.SpecialerGuideAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.SpecialerGuideIOS) {
		params.Set("is_special_guide", "1")
	}
	// new request
	if req, err = d.client.NewRequest("GET", d.main, ip, params); err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	var res struct {
		Code  int            `json:"code"`
		SeID  string         `json:"seid"`
		Pages int            `json:"numPages"`
		List  []*search.Live `json:"result"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if (model.IsAndroid(plat) && build > _searchCodeLimitAndroid) || (model.IsIPhone(plat) && build > _searchCodeLimitIPhone) || (plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD) || model.IsIPhoneB(plat) {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		} else if res.Code != model.ForbidCode && res.Code != model.NoResultCode {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		}
		return
	}
	if !model.IsAndroid(plat) || (model.IsAndroid(plat) && build > search.LiveBroadcastTypeAndroid) {
		for _, v := range res.List {
			roomIDs = append(roomIDs, v.RoomID)
		}
		if len(roomIDs) > 0 {
			if lm, err = d.liveDao.LiveByRIDs(c, roomIDs); err != nil {
				log.Error("Live %+v", err)
				err = nil
			}
		}
	}
	items := make([]*search.Item, 0, len(res.List))
	for _, v := range res.List {
		si := &search.Item{}
		si.FromLive2(v, lm[v.RoomID])
		items = append(items, si)
	}
	st = &search.TypeSearch{TrackID: res.SeID, Pages: res.Pages, Items: items}
	return
}

// LiveAll search for live version > 5.28
func (d *Dao) LiveAll(c context.Context, mid int64, keyword, mobiApp, platform, buvid, device, order, sType string, build, pn, ps int) (st *search.TypeSearchLiveAll, err error) {
	var (
		req     *http.Request
		plat    = model.Plat(mobiApp, device)
		roomIDs []int64
		lm      map[int64]*livemdl.RoomInfo
		ip      = metadata.String(c, metadata.RemoteIP)
	)
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("mobi_app", mobiApp)
	params.Set("platform", platform)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("build", strconv.Itoa(build))
	params.Set("main_ver", "v3")
	params.Set("search_type", sType)
	params.Set("page", strconv.Itoa(pn))
	params.Set("live_room_num", strconv.Itoa(ps))
	params.Set("device", device)
	params.Set("order", order)
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.SpecialerGuideAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.SpecialerGuideIOS) {
		params.Set("is_special_guide", "1")
	}
	// new request
	if req, err = d.client.NewRequest("GET", d.main, ip, params); err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	var res struct {
		Code     int    `json:"code"`
		SeID     string `json:"seid"`
		Pages    int    `json:"numPages"`
		PageInfo *struct {
			Master *search.Live `json:"live_master"`
			Room   *search.Live `json:"live_room"`
		} `json:"pageinfo"`
		List *struct {
			Master []*search.Live `json:"live_master"`
			Room   []*search.Live `json:"live_room"`
		} `json:"result"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if (model.IsAndroid(plat) && build > _searchCodeLimitAndroid) || (model.IsIPhone(plat) && build > _searchCodeLimitIPhone) || model.IsIPad(plat) || model.IsIPhoneB(plat) {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		} else if res.Code != model.ForbidCode && res.Code != model.NoResultCode {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		}
		return
	}
	st = &search.TypeSearchLiveAll{
		TrackID: res.SeID,
		Pages:   res.Pages,
		Master:  &search.TypeSearch{},
		Room:    &search.TypeSearch{},
	}
	if res.PageInfo != nil {
		if res.PageInfo.Master != nil {
			st.Master.Pages = res.PageInfo.Master.Pages
			st.Master.Total = res.PageInfo.Master.Total
		}
		if res.PageInfo.Room != nil {
			st.Room.Pages = res.PageInfo.Room.Pages
			st.Room.Total = res.PageInfo.Room.Total
		}
	}
	if res.List != nil {
		if !model.IsAndroid(plat) || (model.IsAndroid(plat) && build > search.LiveBroadcastTypeAndroid) {
			for _, v := range res.List.Master {
				roomIDs = append(roomIDs, v.RoomID)
			}
			for _, v := range res.List.Room {
				roomIDs = append(roomIDs, v.RoomID)
			}
			if len(roomIDs) > 0 {
				if lm, err = d.liveDao.LiveByRIDs(c, roomIDs); err != nil {
					log.Error("LiveAll %+v", err)
					err = nil
				}
			}
		}
		st.Master.Items = make([]*search.Item, 0, len(res.List.Master))
		for _, v := range res.List.Master {
			si := &search.Item{}
			si.FromLiveMaster(v, lm[v.RoomID])
			st.Master.Items = append(st.Master.Items, si)
		}
		st.Room.Items = make([]*search.Item, 0, len(res.List.Room))
		for _, v := range res.List.Room {
			si := &search.Item{}
			si.FromLive2(v, lm[v.RoomID])
			st.Room.Items = append(st.Room.Items, si)
		}
	}
	return
}

// ArticleByType search article.
func (d *Dao) ArticleByType(c context.Context, mid, zoneid int64, keyword, mobiApp, device, platform, buvid, filtered, order, sType string, plat int8, categoryID, build, highlight, pn, ps int, now time.Time) (st *search.TypeSearch, err error) {
	var (
		req *http.Request
		ip  = metadata.String(c, metadata.RemoteIP)
	)
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("platform", platform)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("build", strconv.Itoa(build))
	params.Set("main_ver", "v3")
	params.Set("highlight", strconv.Itoa(highlight))
	params.Set("search_type", sType)
	params.Set("category_id", strconv.Itoa(categoryID))
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("order", order)
	if filtered == "1" {
		params.Set("filtered", filtered)
	}
	if model.IsOverseas(plat) {
		params.Set("use_area", "1")
		params.Set("zone_id", strconv.FormatInt(zoneid, 10))
	}
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.SpecialerGuideAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.SpecialerGuideIOS) {
		params.Set("is_special_guide", "1")
	}
	if req, err = d.client.NewRequest("GET", d.main, ip, params); err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	var res struct {
		Code  int               `json:"code"`
		SeID  string            `json:"seid"`
		Pages int               `json:"numPages"`
		List  []*search.Article `json:"result"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if res.Code != model.ForbidCode && res.Code != model.NoResultCode {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		}
		return
	}
	items := make([]*search.Item, 0, len(res.List))
	for _, v := range res.List {
		si := &search.Item{}
		si.FromArticle(v, nil)
		items = append(items, si)
	}
	st = &search.TypeSearch{TrackID: res.SeID, Pages: res.Pages, Items: items}
	return
}

// HotSearch is hot words search data.
func (d *Dao) HotSearch(c context.Context, buvid string, mid int64, build, limit int, mobiApp, device, platform string, now time.Time) (res *search.Hot, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("main_ver", "v3")
	params.Set("actionKey", "appkey")
	params.Set("limit", strconv.Itoa(limit))
	params.Set("build", strconv.Itoa(build))
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("platform", platform)
	params.Set("userid", strconv.FormatInt(mid, 10))
	req, err := d.client.NewRequest("GET", d.hot, ip, params)
	if err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.hot+"?"+params.Encode())
	}
	return
}

// Suggest suggest data.
func (d *Dao) Suggest(c context.Context, mid int64, buvid, term string, build int, mobiApp, device string, now time.Time) (res *search.Suggest, err error) {
	plat := model.Plat(mobiApp, device)
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("main_ver", "v4")
	params.Set("func", "suggest")
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("buvid", buvid)
	params.Set("build", strconv.Itoa(build))
	params.Set("mobi_app", mobiApp)
	params.Set("bangumi_acc_num", "3")
	params.Set("special_acc_num", "3")
	params.Set("topic_acc_num", "3")
	params.Set("upuser_acc_num", "1")
	params.Set("tag_num", "10")
	params.Set("special_num", "10")
	params.Set("bangumi_num", "10")
	params.Set("upuser_num", "3")
	params.Set("suggest_type", "accurate")
	params.Set("term", term)
	res = &search.Suggest{}
	if err = d.client.Get(c, d.suggest, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if (model.IsAndroid(plat) && build > _searchCodeLimitAndroid) || (model.IsIPhone(plat) && build > _searchCodeLimitIPhone) || (plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD) || model.IsIPhoneB(plat) {
			err = errors.Wrap(ecode.Int(res.Code), d.suggest+"?"+params.Encode())
		} else if res.Code != model.ForbidCode && res.Code != model.NoResultCode {
			err = errors.Wrap(ecode.Int(res.Code), d.suggest+"?"+params.Encode())
		}
		return
	}
	if res.ResultBs == nil {
		return
	}
	switch v := res.ResultBs.(type) {
	case []interface{}:
		return
	case map[string]interface{}:
		if acc, ok := v["accurate"]; ok {
			if accm, ok := acc.(map[string]interface{}); ok && accm != nil {
				res.Result.Accurate.UpUser = accm["upuser"]
				res.Result.Accurate.Bangumi = accm["bangumi"]
			}
		}
		if tag, ok := v["tag"]; ok {
			if tags, ok := tag.([]interface{}); ok {
				for _, t := range tags {
					if tm, ok := t.(map[string]interface{}); ok && tm != nil {
						if v, ok := tm["value"]; ok {
							if vs, ok := v.(string); ok {
								res.Result.Tag = append(res.Result.Tag, &struct {
									Value string `json:"value,omitempty"`
								}{vs})
							}
						}
					}
				}
			}
		}
	}
	return
}

// Suggest2 suggest data.
func (d *Dao) Suggest2(c context.Context, mid int64, platform, buvid, term string, build int, mobiApp string, now time.Time) (res *search.Suggest2, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("main_ver", "v4")
	params.Set("suggest_type", "accurate")
	params.Set("platform", platform)
	params.Set("mobi_app", mobiApp)
	params.Set("clientip", ip)
	params.Set("build", strconv.Itoa(build))
	if mid != 0 {
		params.Set("userid", strconv.FormatInt(mid, 10))
	}
	params.Set("term", term)
	params.Set("sug_num", "10")
	params.Set("buvid", buvid)
	res = &search.Suggest2{}
	if err = d.client.Get(c, d.suggest, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.suggest+"?"+params.Encode())
	}
	return
}

// Suggest3 suggest data.
func (d *Dao) Suggest3(c context.Context, mid int64, platform, buvid, term, device string, build, highlight int, mobiApp string, now time.Time) (res *search.Suggest3, err error) {
	var (
		req  *http.Request
		plat = model.Plat(mobiApp, device)
		ip   = metadata.String(c, metadata.RemoteIP)
	)
	params := url.Values{}
	params.Set("suggest_type", "accurate")
	params.Set("platform", platform)
	params.Set("mobi_app", mobiApp)
	params.Set("clientip", ip)
	params.Set("highlight", strconv.Itoa(highlight))
	params.Set("build", strconv.Itoa(build))
	if mid != 0 {
		params.Set("userid", strconv.FormatInt(mid, 10))
	}
	params.Set("term", term)
	params.Set("sug_num", "10")
	params.Set("buvid", buvid)
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.SugDetailAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.SugDetailIOS) {
		params.Set("is_detail_info", "1")
	}
	if req, err = d.client.NewRequest("GET", d.suggest3, ip, params); err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	res = &search.Suggest3{}
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.suggest3+"?"+params.Encode())
		return
	}
	for _, flow := range res.Result {
		flow.SugChange()
	}
	return
}

// Season2 search new season data.
func (d *Dao) Season2(c context.Context, mid int64, keyword, mobiApp, device, platform, buvid string, highlight, build, pn, ps int) (st *search.TypeSearch, err error) {
	var (
		req       *http.Request
		plat      = model.Plat(mobiApp, device)
		ip        = metadata.String(c, metadata.RemoteIP)
		seasonIDs []int64
		bangumis  map[string]*bangumimdl.Card
	)
	params := url.Values{}
	params.Set("main_ver", "v3")
	params.Set("platform", platform)
	params.Set("build", strconv.Itoa(build))
	params.Set("keyword", keyword)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("search_type", "media_bangumi")
	params.Set("order", "totalrank")
	params.Set("highlight", strconv.Itoa(highlight))
	params.Set("app_highlight", "media_bangumi")
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.PGCALLAndroid) || (model.IsIPhone(plat) && build >= d.c.SearchBuildLimit.PGCALLIOS) || model.IsIPhoneB(plat) {
		params.Set("is_pgc_all", "1")
	}
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.SpecialerGuideAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.SpecialerGuideIOS) {
		params.Set("is_special_guide", "1")
	}
	if req, err = d.client.NewRequest("GET", d.main, ip, params); err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	var res struct {
		Code  int             `json:"code"`
		SeID  string          `json:"seid"`
		Total int             `json:"numResults"`
		Pages int             `json:"numPages"`
		List  []*search.Media `json:"result"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if (model.IsAndroid(plat) && build > _searchCodeLimitAndroid) || (model.IsIPhone(plat) && build > _searchCodeLimitIPhone) || (plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD) || model.IsIPhoneB(plat) {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		} else if res.Code != model.ForbidCode && res.Code != model.NoResultCode {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		}
		return
	}
	for _, v := range res.List {
		seasonIDs = append(seasonIDs, v.SeasonID)
	}
	if len(seasonIDs) > 0 {
		if bangumis, err = d.bangumiDao.Card(c, mid, seasonIDs); err != nil {
			log.Error("Season2 %+v", err)
			err = nil
		}
	}
	items := make([]*search.Item, 0, len(res.List))
	for _, v := range res.List {
		si := &search.Item{}
		si.FromMedia(v, "", model.GotoBangumi, bangumis)
		items = append(items, si)
	}
	st = &search.TypeSearch{TrackID: res.SeID, Pages: res.Pages, Total: res.Total, Items: items}
	return
}

// MovieByType2 search new movie data from api .
func (d *Dao) MovieByType2(c context.Context, mid int64, keyword, mobiApp, device, platform, buvid string, highlight, build, pn, ps int) (st *search.TypeSearch, err error) {
	var (
		req       *http.Request
		plat      = model.Plat(mobiApp, device)
		ip        = metadata.String(c, metadata.RemoteIP)
		seasonIDs []int64
		bangumis  map[string]*bangumimdl.Card
	)
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("platform", platform)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("build", strconv.Itoa(build))
	params.Set("main_ver", "v3")
	params.Set("search_type", "media_ft")
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("order", "totalrank")
	params.Set("highlight", strconv.Itoa(highlight))
	params.Set("app_highlight", "media_ft")
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.PGCALLAndroid) || (model.IsIPhone(plat) && build >= d.c.SearchBuildLimit.PGCALLIOS) || model.IsIPhoneB(plat) {
		params.Set("is_pgc_all", "1")
	}
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.SpecialerGuideAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.SpecialerGuideIOS) {
		params.Set("is_special_guide", "1")
	}
	if req, err = d.client.NewRequest("GET", d.main, ip, params); err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	var res struct {
		Code  int             `json:"code"`
		SeID  string          `json:"seid"`
		Total int             `json:"numResults"`
		Pages int             `json:"numPages"`
		List  []*search.Media `json:"result"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if (model.IsAndroid(plat) && build > _searchCodeLimitAndroid) || (model.IsIPhone(plat) && build > _searchCodeLimitIPhone) || (plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD) || model.IsIPhoneB(plat) {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		} else if res.Code != model.ForbidCode && res.Code != model.NoResultCode {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		}
		return
	}
	for _, v := range res.List {
		seasonIDs = append(seasonIDs, v.SeasonID)
	}
	if len(seasonIDs) > 0 {
		if bangumis, err = d.bangumiDao.Card(c, mid, seasonIDs); err != nil {
			log.Error("MovieByType2 %+v", err)
			err = nil
		}
	}
	items := make([]*search.Item, 0, len(res.List))
	for _, v := range res.List {
		si := &search.Item{}
		si.FromMedia(v, "", model.GotoMovie, bangumis)
		items = append(items, si)
	}
	st = &search.TypeSearch{TrackID: res.SeID, Pages: res.Pages, Total: res.Total, Items: items}
	return
}

// User search user data.
func (d *Dao) User(c context.Context, mid int64, keyword, mobiApp, device, platform, buvid, filtered, order, fromSource string, highlight, build, userType, orderSort, pn, ps int, now time.Time) (user []*search.User, err error) {
	var (
		req  *http.Request
		plat = model.Plat(mobiApp, device)
		ip   = metadata.String(c, metadata.RemoteIP)
	)
	params := url.Values{}
	params.Set("platform", platform)
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("build", strconv.Itoa(build))
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("keyword", keyword)
	params.Set("highlight", strconv.Itoa(highlight))
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("main_ver", "v3")
	params.Set("func", "search")
	params.Set("smerge", "1")
	params.Set("source_type", "0")
	params.Set("search_type", "bili_user")
	params.Set("user_type", strconv.Itoa(userType))
	params.Set("order", order)
	params.Set("order_sort", strconv.Itoa(orderSort))
	params.Set("from_source", fromSource)
	if filtered == "1" {
		params.Set("filtered", filtered)
	}
	// new request
	if req, err = d.client.NewRequest("GET", d.main, ip, params); err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	var res struct {
		Code  int            `json:"code"`
		SeID  string         `json:"seid"`
		Pages int            `json:"numPages"`
		List  []*search.User `json:"result"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if (model.IsAndroid(plat) && build > _searchCodeLimitAndroid) || (model.IsIPhone(plat) && build > _searchCodeLimitIPhone) || (plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD) || model.IsIPhoneB(plat) {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		} else if res.Code != model.ForbidCode && res.Code != model.NoResultCode {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		}
		return
	}
	user = res.List
	return
}

// Recommend is recommend search data.
func (d *Dao) Recommend(c context.Context, mid int64, build, from, show int, buvid, platform, mobiApp, device string) (res *search.RecommendResult, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("main_ver", "v3")
	params.Set("platform", platform)
	params.Set("mobi_app", mobiApp)
	params.Set("clientip", ip)
	params.Set("device", device)
	params.Set("build", strconv.Itoa(build))
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("search_type", "guess")
	params.Set("req_source", strconv.Itoa(from))
	params.Set("show_area", strconv.Itoa(show))
	req, err := d.client.NewRequest("GET", d.rcmd, ip, params)
	if err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	var rcmdRes struct {
		Code      int    `json:"code,omitempty"`
		SeID      string `json:"seid,omitempty"`
		Tips      string `json:"recommend_tips,omitempty"`
		NumResult int    `json:"numResult,omitempty"`
		Resutl    []struct {
			ID   int64  `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
			Type string `json:"type,omitempty"`
		} `json:"result,omitempty"`
	}
	if err = d.client.Do(c, req, &rcmdRes); err != nil {
		return
	}
	if rcmdRes.Code != ecode.OK.Code() {
		if rcmdRes.Code != model.ForbidCode {
			err = errors.Wrap(ecode.Int(rcmdRes.Code), d.rcmd+"?"+params.Encode())
		}
		return
	}
	res = &search.RecommendResult{TrackID: rcmdRes.SeID, Title: rcmdRes.Tips, Pages: rcmdRes.NumResult}
	for _, v := range rcmdRes.Resutl {
		item := &search.Item{}
		item.ID = v.ID
		item.Param = strconv.Itoa(int(v.ID))
		item.Title = v.Name
		item.Type = v.Type
		res.Items = append(res.Items, item)
	}
	return
}

// DefaultWords is default words search data.
func (d *Dao) DefaultWords(c context.Context, mid int64, build, from int, buvid, platform, mobiApp, device string) (res *search.DefaultWords, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("main_ver", "v3")
	params.Set("platform", platform)
	params.Set("mobi_app", mobiApp)
	params.Set("clientip", ip)
	params.Set("device", device)
	params.Set("build", strconv.Itoa(build))
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("search_type", "default")
	params.Set("req_source", strconv.Itoa(from))
	req, err := d.client.NewRequest("GET", d.rcmd, ip, params)
	if err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	var rcmdRes struct {
		Code      int    `json:"code,omitempty"`
		SeID      string `json:"seid,omitempty"`
		Tips      string `json:"recommend_tips,omitempty"`
		NumResult int    `json:"numResult,omitempty"`
		ShowFront int    `json:"show_front,omitempty"`
		Resutl    []struct {
			ID       int64  `json:"id,omitempty"`
			Name     string `json:"name,omitempty"`
			ShowName string `json:"show_name,omitempty"`
			Type     string `json:"type,omitempty"`
		} `json:"result,omitempty"`
	}
	if err = d.client.Do(c, req, &rcmdRes); err != nil {
		return
	}
	if rcmdRes.Code != ecode.OK.Code() {
		if rcmdRes.Code != model.ForbidCode {
			err = errors.Wrap(ecode.Int(rcmdRes.Code), d.rcmd+"?"+params.Encode())
		}
		return
	}
	if len(rcmdRes.Resutl) == 0 {
		err = ecode.NothingFound
		return
	}
	res = &search.DefaultWords{}
	for _, v := range rcmdRes.Resutl {
		res.Trackid = rcmdRes.SeID
		res.Param = strconv.Itoa(int(v.ID))
		res.Show = v.ShowName
		res.Word = v.Name
		res.ShowFront = rcmdRes.ShowFront
	}
	return
}

// RecommendNoResult is no result recommend search data.
func (d *Dao) RecommendNoResult(c context.Context, platform, mobiApp, device, buvid, keyword string, build, pn, ps int, mid int64) (res *search.NoResultRcndResult, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("main_ver", "v3")
	params.Set("platform", platform)
	params.Set("mobi_app", mobiApp)
	params.Set("clientip", ip)
	params.Set("device", device)
	params.Set("build", strconv.Itoa(build))
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("search_type", "video")
	params.Set("keyword", keyword)
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	req, err := d.client.NewRequest("GET", d.rcmdNoResult, ip, params)
	if err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	var (
		resTmp      *search.NoResultRcmd
		avids       []int64
		avm         map[int64]*api.Arc
		cooperation bool
	)
	if err = d.client.Do(c, req, &resTmp); err != nil {
		return
	}
	if resTmp.Code != ecode.OK.Code() {
		if resTmp.Code != model.ForbidCode {
			err = errors.Wrap(ecode.Int(resTmp.Code), d.rcmdNoResult+"?"+params.Encode())
		}
		return
	}
	res = &search.NoResultRcndResult{TrackID: resTmp.Trackid, Title: resTmp.RecommendTips, Pages: resTmp.NumResults}
	for _, v := range resTmp.Result {
		avids = append(avids, v.ID)
	}
	if len(avids) != 0 {
		if avm, err = d.arcDao.Archives2(c, avids); err != nil {
			log.Error("RecommendNoResult %+v", err)
			err = nil
		}
	}
	items := make([]*search.Item, 0, len(resTmp.Result))
	for _, v := range resTmp.Result {
		ri := &search.Item{}
		ri.FromVideo(v, avm[v.ID], cooperation)
		items = append(items, ri)
	}
	res.Items = items
	return
}

// Channel for search channel
func (d *Dao) Channel(c context.Context, mid int64, keyword, mobiApp, platform, buvid, device, order, sType string, build, pn, ps, highlight int) (st *search.TypeSearch, err error) {
	var (
		req  *http.Request
		plat = model.Plat(mobiApp, device)
		ip   = metadata.String(c, metadata.RemoteIP)
	)
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("mobi_app", mobiApp)
	params.Set("platform", platform)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("build", strconv.Itoa(build))
	params.Set("main_ver", "v3")
	params.Set("search_type", sType)
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("device", device)
	params.Set("order", order)
	params.Set("highlight", strconv.Itoa(highlight))
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.SpecialerGuideAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.SpecialerGuideIOS) {
		params.Set("is_special_guide", "1")
	}
	// new request
	if req, err = d.client.NewRequest("GET", d.main, ip, params); err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	var res struct {
		Code  int               `json:"code"`
		SeID  string            `json:"seid"`
		Pages int               `json:"numPages"`
		List  []*search.Channel `json:"result"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if res.Code != model.ForbidCode && res.Code != model.NoResultCode {
			err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		}
		return
	}
	items := make([]*search.Item, 0, len(res.List))
	for _, v := range res.List {
		si := &search.Item{}
		avm := make(map[int64]*api.Arc)
		bangumis := make(map[int32]*seasongrpc.CardInfoProto)
		lm := make(map[int64]*livemdl.RoomInfo)
		si.FromChannel(v, avm, bangumis, lm, nil)
		items = append(items, si)
	}
	st = &search.TypeSearch{TrackID: res.SeID, Pages: res.Pages, Items: items}
	return
}

// RecommendPre search at pre-page
func (d *Dao) RecommendPre(c context.Context, platform, mobiApp, device, buvid string, build, ps int, mid int64) (res *search.RecommendPreResult, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("main_ver", "v3")
	params.Set("platform", platform)
	params.Set("mobi_app", mobiApp)
	params.Set("clientip", ip)
	params.Set("device", device)
	params.Set("build", strconv.Itoa(build))
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("search_type", "discover_page")
	params.Set("pagesize", strconv.Itoa(ps))
	req, err := d.client.NewRequest("GET", d.pre, ip, params)
	if err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	var (
		resTmp    *search.RecommendPre
		avids     []int64
		avm       map[int64]*api.Arc
		seasonIDs []int32
		bangumis  map[int32]*seasongrpc.CardInfoProto
	)
	if err = d.client.Do(c, req, &resTmp); err != nil {
		return
	}
	if resTmp.Code != ecode.OK.Code() {
		if resTmp.Code != model.ForbidCode {
			err = errors.Wrap(ecode.Int(resTmp.Code), d.pre+"?"+params.Encode())
		}
		return
	}
	for _, v := range resTmp.Result {
		for _, vv := range v.List {
			if vv.Type == "video" {
				avids = append(avids, vv.ID)
			} else if vv.Type == "pgc" {
				seasonIDs = append(seasonIDs, int32(vv.ID))
			}
		}
	}
	g, ctx := errgroup.WithContext(c)
	if len(avids) != 0 {
		g.Go(func() (err error) {
			if avm, err = d.arcDao.Archives2(ctx, avids); err != nil {
				log.Error("RecommendPre avids(%v) error(%v)", avids, err)
				err = nil
			}
			return
		})
	}
	if len(seasonIDs) > 0 {
		g.Go(func() (err error) {
			if bangumis, err = d.bangumiDao.Cards(c, seasonIDs); err != nil {
				log.Error("RecommendPre seasonIDs(%v) error(%v)", seasonIDs, err)
				err = nil
			}
			return
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
		return
	}
	res = &search.RecommendPreResult{TrackID: resTmp.Trackid, Total: resTmp.NumResult}
	items := make([]*search.Item, 0, len(resTmp.Result))
	for _, v := range resTmp.Result {
		rs := &search.Item{Title: v.Query}
		for _, vv := range v.List {
			if vv.Type == "video" {
				if a, ok := avm[vv.ID]; ok {
					r := &search.Item{}
					r.FromRcmdPre(vv.ID, a, nil)
					rs.Item = append(rs.Item, r)
				}
			} else if vv.Type == "pgc" {
				if b, ok := bangumis[int32(vv.ID)]; ok {
					r := &search.Item{}
					r.FromRcmdPre(vv.ID, nil, b)
					rs.Item = append(rs.Item, r)
				}
			}
		}
		items = append(items, rs)
	}
	res.Items = items
	return
}

// Video search new archive data.
func (d *Dao) Video(c context.Context, mid int64, keyword, mobiApp, device, platform, buvid string, highlight, build, pn, ps int) (st *search.TypeSearch, err error) {
	var (
		req         *http.Request
		ip          = metadata.String(c, metadata.RemoteIP)
		plat        = model.Plat(mobiApp, device)
		avids       []int64
		avm         map[int64]*api.Arc
		cooperation bool
	)
	params := url.Values{}
	params.Set("main_ver", "v3")
	params.Set("platform", platform)
	params.Set("build", strconv.Itoa(build))
	params.Set("keyword", keyword)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("search_type", "video")
	params.Set("order", "totalrank")
	params.Set("highlight", strconv.Itoa(highlight))
	if (model.IsAndroid(plat) && build > d.c.SearchBuildLimit.SpecialerGuideAndroid) || (model.IsIPhone(plat) && build > d.c.SearchBuildLimit.SpecialerGuideIOS) {
		params.Set("is_special_guide", "1")
	}
	if req, err = d.client.NewRequest("GET", d.main, ip, params); err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	var res struct {
		Code  int             `json:"code"`
		SeID  string          `json:"seid"`
		Total int             `json:"numResults"`
		Pages int             `json:"numPages"`
		List  []*search.Video `json:"result"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		return
	}
	for _, v := range res.List {
		avids = append(avids, v.ID)
	}
	if len(avids) > 0 {
		if avm, err = d.arcDao.Archives2(c, avids); err != nil {
			log.Error("Upper %+v", err)
			err = nil
		}
	}
	items := make([]*search.Item, 0, len(res.List))
	for _, v := range res.List {
		si := &search.Item{}
		si.FromVideo(v, avm[v.ID], cooperation)
		items = append(items, si)
	}
	st = &search.TypeSearch{TrackID: res.SeID, Pages: res.Pages, Total: res.Total, Items: items}
	return
}
