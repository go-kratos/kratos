package creative

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-view/conf"
	"go-common/app/interface/main/app-view/model/creative"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_special = "/x/internal/uper/special"
	_follow  = "/x/internal/uper/switch"
	_bgm     = "/x/internal/creative/archive/bgm"
	_points  = "/x/internal/creative/video/viewpoints"
)

// Dao is space dao
type Dao struct {
	client  *httpx.Client
	special string
	follow  string
	bgm     string
	points  string
}

// New initial space dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:  httpx.NewClient(c.HTTPClient),
		special: c.Host.APICo + _special,
		follow:  c.Host.APICo + _follow,
		bgm:     c.Host.APICo + _bgm,
		points:  c.Host.APICo + _points,
	}
	return
}

// FollowSwitch get auto follow switch .
func (d *Dao) FollowSwitch(c context.Context, vmid int64) (s *creative.FollowSwitch, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(vmid, 10))
	params.Set("from", "0")
	var res struct {
		Code int                    `json:"code"`
		Data *creative.FollowSwitch `json:"data"`
	}
	ip := metadata.String(c, metadata.RemoteIP)
	if err = d.client.Get(c, d.follow, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.follow+"?"+params.Encode())
		return
	}
	s = res.Data
	return
}

// Bgm get archive bgm
func (d *Dao) Bgm(c context.Context, aid, cid int64) (bgm []*creative.Bgm, err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("cid", strconv.FormatInt(cid, 10))
	var res struct {
		Code int             `json:"code"`
		Data []*creative.Bgm `json:"data"`
	}
	ip := metadata.String(c, metadata.RemoteIP)
	if err = d.client.Get(c, d.bgm, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.bgm+"?"+params.Encode())
		return
	}
	bgm = res.Data
	return
}

// Points get video points
func (d *Dao) Points(c context.Context, aid, cid int64) (points []*creative.Points, err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("cid", strconv.FormatInt(cid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Points []*creative.Points `json:"points"`
		} `json:"data"`
	}
	ip := metadata.String(c, metadata.RemoteIP)
	if err = d.client.Get(c, d.points, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.points+"?"+params.Encode())
		return
	}
	points = res.Data.Points
	return
}

// Special is
func (d *Dao) Special(c context.Context) (midsM map[int64]struct{}, err error) {
	params := url.Values{}
	params.Set("group_id", "20")
	var res struct {
		Code int `json:"code"`
		Data []struct {
			Mid int64 `json:"mid"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.special, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.special+"?"+params.Encode())
		return
	}
	midsM = make(map[int64]struct{})
	for _, l := range res.Data {
		midsM[l.Mid] = struct{}{}
	}
	return
}
