package space

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-interface/conf"
	"go-common/app/interface/main/app-interface/model/space"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_space             = "/x/internal/space/setting"
	_video             = "/api/member/getVideo"
	_uploadTopPhotoURL = "/api/member/getUploadTopPhoto"
	_report            = "/api/report/add"
	_blacklist         = "/x/internal/space/blacklist"
)

// Dao is space dao
type Dao struct {
	client     *httpx.Client
	clientSync *httpx.Client
	space      string
	video      string
	report     string
	// space api
	uploadTop string
	blacklist string
}

// New initial space dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:     httpx.NewClient(c.HTTPClient),
		clientSync: httpx.NewClient(c.HTTPWrite),
		space:      c.Host.APICo + _space,
		video:      c.Host.Space + _video,
		report:     c.Host.Space + _report,
		uploadTop:  c.Host.Space + _uploadTopPhotoURL,
		blacklist:  c.Host.APICo + _blacklist,
	}
	return
}

// Setting get setting data from  api.
func (d *Dao) Setting(c context.Context, mid int64) (setting *space.Setting, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Privacy *space.Setting `json:"privacy"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.space, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.space+"?"+params.Encode())
	}
	setting = res.Data.Privacy
	return
}

// SpaceMob space mobile
func (d *Dao) SpaceMob(c context.Context, mid, vmid int64, platform, device string) (us *space.Mob, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("vmid", strconv.FormatInt(vmid, 10))
	params.Set("platform", platform)
	params.Set("device", device)
	var res struct {
		Code int        `json:"code"`
		Data *space.Mob `json:"data"`
	}
	if err = d.client.Get(c, d.uploadTop, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.uploadTop+"?"+params.Encode())
		return
	}
	us = res.Data
	return
}

// Report report
func (d *Dao) Report(c context.Context, mid int64, reason, ak string) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("access_key", ak)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("reason", reason)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.report, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.report+"?"+params.Encode())
	}
	return
}

// Blacklist is.
func (d *Dao) Blacklist(c context.Context) (list map[int64]struct{}, err error) {
	var res struct {
		Code int                `json:"code"`
		Data map[int64]struct{} `json:"data"`
	}
	if err = d.clientSync.Get(c, d.blacklist, "", nil, &res); err != nil {
		err = errors.Wrap(ecode.Int(res.Code), d.blacklist)
		return
	}
	b, _ := json.Marshal(res)
	log.Error("Blacklist url(%s) response(%s)", d.blacklist, b)
	list = res.Data
	return
}
