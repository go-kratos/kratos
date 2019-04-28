package creative

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"go-common/app/interface/main/videoup/conf"
	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	_setWatermark   = "/x/internal/creative/watermark/set"
	_uploadMaterial = "/x/internal/creative/upload/material"
)

// SetWatermark fn
func (d *Dao) SetWatermark(c context.Context, mid int64, state, ty, pos int8, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("state", strconv.Itoa(int(state)))
	params.Set("type", strconv.Itoa(int(ty)))
	params.Set("position", strconv.Itoa(int(pos)))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpW.Post(c, d.setWatermarkURL, ip, params, &res); err != nil {
		log.Error("d.httpW.Post(%s) error(%v)", d.setWatermarkURL+"?"+params.Encode(), err)
		return
	}
	log.Info("SetWatermark url(%s) code(%d)", d.setWatermarkURL+"?"+params.Encode(), res.Code)
	if res.Code != 0 {
		log.Error("url(%s) code(%d)", d.setWatermarkURL+"?"+params.Encode(), res.Code)
	}
	return
}

// UploadMaterial fn
func (d *Dao) UploadMaterial(c context.Context, editors []*archive.Editor, aid, mid int64, ip string) (err error) {
	params := url.Values{}
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	mh := md5.Sum([]byte(params.Encode() + d.c.App.Secret))
	params.Set("sign", hex.EncodeToString(mh[:]))
	var (
		uri = d.uploadMaterialURL + "?" + params.Encode()
	)
	bs, err := json.Marshal(editors)
	if err != nil {
		log.Error("UploadMaterial json.Marshal error(%+v)|editor(%+v)", err, editors)
		return
	}
	req, err := http.NewRequest("POST", uri, bytes.NewReader(bs))
	if err != nil {
		log.Error("UploadMaterial http.NewRequest error(%v) | uri(%s) bs(%+v)", err, uri, bs)
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = d.httpW.Do(c, req, &res); err != nil {
		log.Error("UploadMaterial do error(%v)|uri(%s)", err, uri)
		return
	}
	log.Info("UploadMaterial url(%s) code(%d)", uri, res.Code)
	if res.Code != 0 {
		log.Error("url(%s) code(%d)", uri, res.Code)
	}
	return
}
