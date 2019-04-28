package search

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-intl/conf"
	arcdao "go-common/app/interface/main/app-intl/dao/archive"
	bgmdao "go-common/app/interface/main/app-intl/dao/bangumi"
	"go-common/app/interface/main/app-intl/model"
	"go-common/app/interface/main/app-intl/model/bangumi"
	"go-common/app/interface/main/app-intl/model/search"
	"go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"

	"github.com/pkg/errors"
)

const (
	_main     = "/main/search"
	_suggest3 = "/main/suggest/new"
)

// Dao is search dao
type Dao struct {
	client   *httpx.Client
	arcDao   *arcdao.Dao
	bgmDao   *bgmdao.Dao
	main     string
	suggest3 string
}

// New initial search dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:   httpx.NewClient(c.HTTPSearch),
		arcDao:   arcdao.New(c),
		bgmDao:   bgmdao.New(c),
		main:     c.Host.Search + _main,
		suggest3: c.Host.Search + _suggest3,
	}
	return
}

// Search app all search .
func (d *Dao) Search(c context.Context, mid, zoneid int64, mobiApp, device, platform, buvid, keyword, duration, order, filtered, fromSource, recommend string, plat int8, seasonNum, movieNum, upUserNum, uvLimit, userNum, userVideoLimit, biliUserNum, biliUserVideoLimit, rid, highlight, build, pn, ps int, now time.Time) (res *search.Search, code int, err error) {
	var (
		req *http.Request
		ip  = metadata.String(c, metadata.RemoteIP)
	)
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
	params.Set("media_bangumi_num", strconv.Itoa(seasonNum))
	params.Set("bili_user_num", strconv.Itoa(biliUserNum))
	params.Set("bili_user_vl", strconv.Itoa(biliUserVideoLimit))
	params.Set("user_num", strconv.Itoa(userNum))
	params.Set("user_video_limit", strconv.Itoa(userVideoLimit))
	params.Set("query_rec_need", recommend)
	params.Set("platform", platform)
	params.Set("duration", duration)
	params.Set("order", order)
	params.Set("search_type", "all")
	params.Set("from_source", fromSource)
	if filtered == "1" {
		params.Set("filtered", filtered)
	}
	params.Set("zone_id", strconv.FormatInt(zoneid, 10))
	params.Set("media_ft_num", strconv.Itoa(movieNum))
	params.Set("is_new_pgc", "1")
	params.Set("is_internation", "1")
	params.Set("no_display_default", "game,live_room")
	params.Set("flow_need", "1")
	params.Set("app_highlight", "media_bangumi,media_ft")
	// new request
	if req, err = d.client.NewRequest("GET", d.main, ip, params); err != nil {
		return
	}
	req.Header.Set("Buvid", buvid)
	if err = d.client.Do(c, req, res); err != nil {
		return
	}
	b, _ := json.Marshal(res)
	log.Error("wocao----%s---%s---%s", d.main+"?"+params.Encode(), buvid, b)
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
	}
	for _, flow := range res.FlowResult {
		flow.Change()
	}
	code = res.Code
	return
}

// Season2 search new season data.
func (d *Dao) Season2(c context.Context, mid int64, keyword, mobiApp, device, platform, buvid string, highlight, build, pn, ps int) (st *search.TypeSearch, err error) {
	var (
		req       *http.Request
		ip        = metadata.String(c, metadata.RemoteIP)
		seasonIDs []int64
		bangumis  map[string]*bangumi.Card
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
	params.Set("is_pgc_all", "1")
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
		err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		return
	}
	for _, v := range res.List {
		seasonIDs = append(seasonIDs, v.SeasonID)
	}
	if len(seasonIDs) > 0 {
		if bangumis, err = d.bgmDao.Card(c, mid, seasonIDs); err != nil {
			log.Error("%+v", err)
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
		ip        = metadata.String(c, metadata.RemoteIP)
		seasonIDs []int64
		bangumis  map[string]*bangumi.Card
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
	params.Set("is_pgc_all", "1")
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
		err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		return
	}
	for _, v := range res.List {
		seasonIDs = append(seasonIDs, v.SeasonID)
	}
	if len(seasonIDs) > 0 {
		if bangumis, err = d.bgmDao.Card(c, mid, seasonIDs); err != nil {
			log.Error("%+v", err)
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

// Upper search upper data.
func (d *Dao) Upper(c context.Context, mid int64, keyword, mobiApp, device, platform, buvid, filtered, order string, biliUserVL, highlight, build, userType, orderSort, pn, ps int, now time.Time) (st *search.TypeSearch, err error) {
	var (
		req   *http.Request
		avids []int64
		avm   map[int64]*api.Arc
		ip    = metadata.String(c, metadata.RemoteIP)
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
	params.Set("search_type", "bili_user")
	params.Set("bili_user_vl", strconv.Itoa(biliUserVL))
	params.Set("user_type", strconv.Itoa(userType))
	params.Set("order_sort", strconv.Itoa(orderSort))
	params.Set("order", order)
	params.Set("source_type", "0")
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
		err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		return
	}
	items := make([]*search.Item, 0, len(res.List))
	for _, v := range res.List {
		for _, vr := range v.Res {
			avids = append(avids, vr.Aid)
		}
	}
	g, ctx := errgroup.WithContext(c)
	if len(avids) != 0 {
		g.Go(func() (err error) {
			if avm, err = d.arcDao.Archives(ctx, avids); err != nil {
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
		si.FromUpUser(v, avm)
		items = append(items, si)
	}
	st = &search.TypeSearch{TrackID: res.SeID, Pages: res.Pages, Items: items}
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
		si.FromArticle(v)
		items = append(items, si)
	}
	st = &search.TypeSearch{TrackID: res.SeID, Pages: res.Pages, Items: items}
	return
}

// Channel for search channel
func (d *Dao) Channel(c context.Context, mid int64, keyword, mobiApp, platform, buvid, device, order, sType string, build, pn, ps, highlight int) (st *search.TypeSearch, err error) {
	var (
		req *http.Request
		ip  = metadata.String(c, metadata.RemoteIP)
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
		err = errors.Wrap(ecode.Int(res.Code), d.main+"?"+params.Encode())
		return
	}
	items := make([]*search.Item, 0, len(res.List))
	for _, v := range res.List {
		si := &search.Item{}
		si.FromChannel(v)
		items = append(items, si)
	}
	st = &search.TypeSearch{TrackID: res.SeID, Pages: res.Pages, Items: items}
	return
}

// Suggest3 suggest data.
func (d *Dao) Suggest3(c context.Context, mid int64, platform, buvid, term string, build, highlight int, mobiApp string, now time.Time) (res *search.Suggest3, err error) {
	var (
		req *http.Request
		ip  = metadata.String(c, metadata.RemoteIP)
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
	}
	return
}
