package dao

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"go-common/app/interface/main/space/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_liveURI        = "/room/v1/Room/getRoomInfoOld"
	_medalStatusURI = "/fans_medal/v1/medal/get_medal_opened"
)

// Live is space live data.
func (d *Dao) Live(c context.Context, mid int64, platform string) (live *model.Live, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("platform", platform)
	var res struct {
		Code int         `json:"code"`
		Data *model.Live `json:"data"`
	}
	if err = d.httpR.Get(c, d.liveURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.liveURL+"?"+params.Encode())
		return
	}
	live = res.Data
	return
}

// LiveMetal get live metal
func (d *Dao) LiveMetal(c context.Context, mid int64) (rs bool, err error) {
	var (
		req *http.Request
		ip  = metadata.String(c, metadata.RemoteIP)
	)
	if req, err = d.httpR.NewRequest(http.MethodGet, d.liveMetalURL, ip, url.Values{}); err != nil {
		log.Error("d.httpR.NewRequest(%s) error(%v)", d.liveMetalURL, err)
		return
	}
	req.Header.Set("X-BILILIVE-UID", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			MasterStatus int `json:"master_status"`
		} `json:"data"`
	}
	if err = d.httpR.Do(c, req, &res); err != nil {
		log.Error("d.httpR.Do(%s) error(%v)", d.liveMetalURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpR.Do(%s) code(%d) error", d.liveMetalURL, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	if res.Data.MasterStatus == 1 {
		rs = true
	}
	return
}
