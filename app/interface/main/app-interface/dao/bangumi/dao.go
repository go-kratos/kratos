package bangumi

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-interface/conf"
	"go-common/app/interface/main/app-interface/model/bangumi"
	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_season     = "/api/inner/season"
	_movie      = "/internal_api/movie_aid_info"
	_bp         = "/sponsor/inner/xAjaxGetBP"
	_concern    = "/api/get_concerned_season"
	_hasFollows = "/follow/internal_api/has_follows"
	_card       = "/pgc/internal/season/search/card"
	_favDisplay = "/pgc/internal/follow/app/tab/view"
)

// Dao is bangumi dao
type Dao struct {
	client     *httpx.Client
	season     string
	movie      string
	bp         string
	concern    string
	hasFollows string
	card       string
	favDisplay string
	// grpc
	rpcClient seasongrpc.SeasonClient
}

// New bangumi dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:     httpx.NewClient(c.HTTPBangumi),
		season:     c.Host.Bangumi + _season,
		movie:      c.Host.Bangumi + _movie,
		bp:         c.Host.Bangumi + _bp,
		concern:    c.Host.Bangumi + _concern,
		hasFollows: c.Host.Bangumi + _hasFollows,
		card:       c.Host.APICo + _card,
		favDisplay: c.Host.APICo + _favDisplay,
	}
	var err error
	if d.rpcClient, err = seasongrpc.NewClient(c.PGCRPC); err != nil {
		log.Error("seasongrpc NewClientt error(%v)", err)
	}
	return
}

// Season bangumi Season.
func (d *Dao) Season(c context.Context, aid, mid int64, ip string) (s *bangumi.Season, err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("type", "av")
	params.Set("build", "app-interface")
	params.Set("platform", "Golang")
	var res struct {
		Code   int             `json:"code"`
		Result *bangumi.Season `json:"result"`
	}
	if err = d.client.Get(c, d.season, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.season+"?"+params.Encode())
		return
	}
	s = res.Result
	return
}

// BPInfo get bp info data.
func (d *Dao) BPInfo(c context.Context, aid, mid int64, ip string) (data json.RawMessage, err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("build", "app-interface")
	params.Set("platform", "Golang")
	var res struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err = d.client.Get(c, d.bp, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.bp+"?"+params.Encode())
		return
	}
	data = res.Data
	return
}

// Movie bangumi Movie
func (d *Dao) Movie(c context.Context, aid, mid int64, build int, mobiApp, device, ip string) (m *bangumi.Movie, err error) {
	params := url.Values{}
	params.Set("id", strconv.FormatInt(aid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("build", strconv.Itoa(build))
	params.Set("device", device)
	params.Set("mobi_app", mobiApp)
	params.Set("platform", "Golang")
	var res struct {
		Code   int            `json:"code"`
		Result *bangumi.Movie `json:"result"`
	}
	if err = d.client.Get(c, d.movie, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.movie+"?"+params.Encode())
		return
	}
	m = res.Result
	return
}

// Concern get concern data from api.
func (d *Dao) Concern(c context.Context, mid, vmid int64, pn, ps int) (ss []*bangumi.Season, total int, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("taid", strconv.FormatInt(vmid, 10))
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("build", "app-interface")
	params.Set("platform", "Golang")
	var res struct {
		Code   int               `json:"code"`
		Total  string            `json:"count"`
		Result []*bangumi.Season `json:"result"`
	}
	if err = d.client.Get(c, d.concern, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.concern+"?"+params.Encode())
		return
	}
	ss = res.Result
	total, _ = strconv.Atoi(res.Total)
	return
}

// HasFollows get bngumi tab.
func (d *Dao) HasFollows(c context.Context, mid int64) (has bool, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("type", "2,3,5")
	var res struct {
		Code   int `json:"code"`
		Result int `json:"result"`
	}
	if err = d.client.Get(c, d.hasFollows, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.hasFollows+"?"+params.Encode())
		return
	}
	if res.Result == 1 {
		has = true
	}
	return
}

// Card bangumi card.
func (d *Dao) Card(c context.Context, mid int64, sids []int64) (s map[string]*bangumi.Card, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("season_ids", xstr.JoinInts(sids))
	var res struct {
		Code   int                      `json:"code"`
		Result map[string]*bangumi.Card `json:"result"`
	}
	if err = d.client.Get(c, d.card, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.card+"?"+params.Encode())
		return
	}
	s = res.Result
	return
}

// FavDisplay fav tab display or not.
func (d *Dao) FavDisplay(c context.Context, mid int64) (bangumi, cinema int, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code   int `json:"code"`
		Result struct {
			Bangumi int `json:"bangumi"`
			Cinema  int `json:"cinema"`
		} `json:"result"`
	}
	if err = d.client.Get(c, d.favDisplay, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.favDisplay+"?"+params.Encode())
		return
	}
	bangumi = res.Result.Bangumi
	cinema = res.Result.Cinema
	return
}
