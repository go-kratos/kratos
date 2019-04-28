package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/space/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_webTopPhotoURI = "/api/member/getTopPhoto"
	_topPhotoURI    = "/api/member/getUploadTopPhoto"
)

// WebTopPhoto getTopPhoto from space
func (d *Dao) WebTopPhoto(c context.Context, mid int64) (space *model.TopPhoto, err error) {
	var (
		params   = url.Values{}
		remoteIP = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		model.TopPhoto
	}
	if err = d.httpR.Get(c, d.webTopPhotoURL, remoteIP, params, &res); err != nil {
		log.Error("d.httpR.Get(%s,%d) error(%v)", d.webTopPhotoURL, mid, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpR.Get(%s,%d) code error(%d)", d.webTopPhotoURL, mid, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	space = &res.TopPhoto
	return
}

// TopPhoto getTopPhoto from space php.
func (d *Dao) TopPhoto(c context.Context, mid, vmid int64, platform, device string) (imgURL string, err error) {
	var (
		params   = url.Values{}
		remoteIP = metadata.String(c, metadata.RemoteIP)
	)
	if mid > 0 {
		params.Set("mid", strconv.FormatInt(mid, 10))
	}
	params.Set("vmid", strconv.FormatInt(vmid, 10))
	params.Set("platform", platform)
	if device != "" {
		params.Set("device", device)
	}
	var res struct {
		Code int `json:"code"`
		Data struct {
			ImgURL string `json:"imgUrl"`
		}
	}
	if err = d.httpR.Get(c, d.topPhotoURL, remoteIP, params, &res); err != nil {
		log.Error("d.httpR.Get(%s,%d) error(%v)", d.topPhotoURL, mid, err)
		return
	}

	if res.Code != ecode.OK.Code() {
		log.Error("d.httpR.Get(%s,%d) code error(%d)", d.topPhotoURL, mid, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	imgURL = res.Data.ImgURL
	return
}
