package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/web/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_abServerURI = "/abserver/v1/web/match-exp"
)

// AbServer get ab server info.
func (d *Dao) AbServer(c context.Context, mid int64, platform int, channel, buvid string) (res model.AbServer, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("channel", channel)
	params.Set("buvid", buvid)
	params.Set("platform", strconv.Itoa(platform))
	var data struct {
		Code int `json:"errorCode"`
		model.AbServer
	}
	if err = d.httpR.Get(c, d.abServerURL, ip, params, &data); err != nil {
		log.Error("AbServer(%s) mid(%d) channel(%s) buvid(%s) error(%v)", d.abServerURL, mid, channel, buvid, err)
		return
	}
	if data.Code != ecode.OK.Code() {
		log.Error("AbServer(%s) mid(%d) channel(%s) buvid(%s) code error(%d)", d.abServerURL, mid, channel, buvid, data.Code)
		return
	}
	res = data.AbServer
	return
}
