package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/space/model"
	"go-common/library/ecode"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_audioCntURI       = "/x/internal/v1/audio/personal/audio-cnt"
	_audioCardURI      = "/x/internal/v1/audio/privilege/mcard"
	_audioUpperCertURI = "/audio/music-service-c/internal/upper-cert"
)

// AudioCard get audio card info.
func (d *Dao) AudioCard(c context.Context, mid ...int64) (cardm map[int64]*model.AudioCard, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", xstr.JoinInts(mid))
	var res struct {
		Code int                        `json:"code"`
		Data map[int64]*model.AudioCard `json:"data"`
	}
	if err = d.httpR.Get(c, d.audioCardURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.audioCardURL+"?"+params.Encode())
		return
	}
	cardm = res.Data
	return
}

// AudioUpperCert get audio upper cert.
func (d *Dao) AudioUpperCert(c context.Context, uid int64) (cert *model.AudioUpperCert, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("uid", strconv.FormatInt(uid, 10))
	var res struct {
		Code int                   `json:"code"`
		Data *model.AudioUpperCert `json:"data"`
	}
	if err = d.httpR.Get(c, d.audioUpperCertURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.audioUpperCertURL+"?"+params.Encode())
		return
	}
	cert = res.Data
	return
}

// AudioCnt get audio cnt.
func (d *Dao) AudioCnt(c context.Context, mid int64) (count int, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Song int `json:"song"`
		} `json:"data"`
	}
	if err = d.httpR.Get(c, d.audioCntURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.audioCntURL+"?"+params.Encode())
		return
	}
	count = res.Data.Song
	return
}
