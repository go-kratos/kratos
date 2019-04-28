package bangumi

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-intl/conf"
	"go-common/app/interface/main/app-intl/model/bangumi"
	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_pgc            = "/pgc/internal/season/appview"
	_movie          = "/internal_api/movie_aid_info"
	_seasonidAidURL = "/api/inner/archive/seasonid2aid"
	_card           = "/pgc/internal/season/search/card"
)

// Dao is bangumi dao
type Dao struct {
	client         *httpx.Client
	pgc            string
	movie          string
	seasonidAidURL string
	card           string
	// grpc
	rpcClient seasongrpc.SeasonClient
}

// New bangumi dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:         httpx.NewClient(c.HTTPBangumi),
		pgc:            c.Host.APICo + _pgc,
		movie:          c.Host.Bangumi + _movie,
		seasonidAidURL: c.Host.Bangumi + _seasonidAidURL,
		card:           c.Host.APICo + _card,
	}
	var err error
	if d.rpcClient, err = seasongrpc.NewClient(nil); err != nil {
		panic(errors.WithMessage(err, "panic by seasongrpc"))
	}
	return
}

// PGC bangumi Season .
func (d *Dao) PGC(c context.Context, aid, mid int64, build int, mobiApp, device string) (s *bangumi.Season, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("build", strconv.Itoa(build))
	params.Set("device", device)
	params.Set("mobi_app", mobiApp)
	params.Set("platform", "Golang")
	var res struct {
		Code   int             `json:"code"`
		Result *bangumi.Season `json:"result"`
	}
	if err = d.client.Get(c, d.pgc, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.pgc+"?"+params.Encode())
		return
	}
	s = res.Result
	return
}

// Movie bangumi Movie
func (d *Dao) Movie(c context.Context, aid, mid int64, build int, mobiApp, device string) (m *bangumi.Movie, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
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

// SeasonidAid moive_id by aid
func (d *Dao) SeasonidAid(c context.Context, moiveID int64, now time.Time) (data map[int64]int64, err error) {
	params := url.Values{}
	params.Set("build", "app-api")
	params.Set("platform", "Golang")
	params.Set("season_id", strconv.FormatInt(moiveID, 10))
	params.Set("season_type", "2")
	var res struct {
		Code   int             `json:"code"`
		Result map[int64]int64 `json:"result"`
	}
	if err = d.client.Get(c, d.seasonidAidURL, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.seasonidAidURL+"?"+params.Encode())
		return
	}
	data = res.Result
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
