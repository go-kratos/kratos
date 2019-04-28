package bangumi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-card/model/card/bangumi"
	"go-common/app/interface/main/app-feed/conf"
	feed "go-common/app/service/main/feed/model"
	episodegrpc "go-common/app/service/openplatform/pgc-season/api/grpc/episode/v1"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_seasons     = "/api/inner/aid_episodes_v2"
	_updates     = "/internal_api/follow_update"
	_pullSeasons = "/internal_api/follow_seasons"
	_followPull  = "/pgc/internal/moe/2018/follow/pull"
)

// Dao is show dao.
type Dao struct {
	// http client
	client *httpx.Client
	// bangumi
	seasons     string
	updates     string
	pullSeasons string
	followPull  string
	// grpc
	rpcClient episodegrpc.EpisodeClient
}

// New new a bangumi dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// http clients
		client:      httpx.NewClient(c.HTTPBangumi),
		seasons:     c.Host.Bangumi + _seasons,
		updates:     c.Host.Bangumi + _updates,
		pullSeasons: c.Host.Bangumi + _pullSeasons,
		followPull:  c.Host.APICo + _followPull,
	}
	var err error
	if d.rpcClient, err = episodegrpc.NewClient(c.PGCRPC); err != nil {
		panic(fmt.Sprintf("episodegrpc NewClientt error (%+v)", err))
	}
	return d
}

func (d *Dao) Seasons(c context.Context, aids []int64, now time.Time) (sm map[int64]*bangumi.Season, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("aids", xstr.JoinInts(aids))
	params.Set("type", "av")
	params.Set("build", "app-feed")
	params.Set("platform", "Golang")
	var res struct {
		Code   int                       `json:"code"`
		Result map[int64]*bangumi.Season `json:"result"`
	}
	if err = d.client.Get(c, d.seasons, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.seasons+"?"+params.Encode())
		return
	}
	sm = res.Result
	return
}

func (d *Dao) Updates(c context.Context, mid int64, now time.Time) (u *bangumi.Update, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code   int             `json:"code"`
		Result *bangumi.Update `json:"result"`
	}
	if err = d.client.Get(c, d.updates, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.updates+"?"+params.Encode())
		return
	}
	u = res.Result
	return
}

func (d *Dao) PullSeasons(c context.Context, seasonIDs []int64, now time.Time) (psm map[int64]*feed.Bangumi, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("season_ids", xstr.JoinInts(seasonIDs))
	var res struct {
		Code   int             `json:"code"`
		Result []*feed.Bangumi `json:"result"`
	}
	if err = d.client.Get(c, d.pullSeasons, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.pullSeasons+"?"+params.Encode())
		return
	}
	psm = make(map[int64]*feed.Bangumi, len(res.Result))
	for _, p := range res.Result {
		psm[p.SeasonID] = p
	}
	return
}

func (d *Dao) FollowPull(c context.Context, mid int64, mobiApp, device string, now time.Time) (moe *bangumi.Moe, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	var res struct {
		Code   int          `json:"code"`
		Result *bangumi.Moe `json:"result"`
	}
	if err = d.client.Get(c, d.followPull, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.followPull+"?"+params.Encode())
		return
	}
	moe = res.Result
	return
}
