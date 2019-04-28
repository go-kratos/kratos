package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"go-common/library/ecode"
	"go-common/library/log"
	"net/http"
)

const (
	_uposURL = "api/v1/task/push/audiowaveform"

	_uposBfsFmt = "subtitle/waveform_%d_%d.pcm"

	_uposCallback = "x/internal/v2/dm/subtitle/upos/callback"

	_defaultPixelDensity = 20
)

// UposReq .
type UposReq struct {
	Cid          int64  `json:"cid"`
	SaveTo       string `json:"saveto"`
	CallbackURL  string `json:"callback_url"`
	PixelDensity int32  `json:"pixel_density"`
}

// UposResp .
type UposResp struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// Upos .
func (d *Dao) Upos(c context.Context, oid int64) (saveTo string, err error) {
	var (
		req  *http.Request
		bs   []byte
		resp []byte
	)
	saveTo = fmt.Sprintf(_uposBfsFmt, oid, 1)
	params := &UposReq{
		Cid:          oid,
		SaveTo:       fmt.Sprintf("bfs://%s", saveTo),
		CallbackURL:  fmt.Sprintf("%s/%s?oid=%d", d.conf.Host.Self, _uposCallback, oid),
		PixelDensity: _defaultPixelDensity,
	}
	if bs, err = json.Marshal(&params); err != nil {
		log.Error("params(%+v),error(%v)", params, err)
		return
	}
	if req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", d.conf.Host.Upos, _uposURL), bytes.NewReader(bs)); err != nil {
		log.Error("http.NewReques.error(%v)", err)
		return
	}
	if resp, err = d.httpCli.Raw(c, req); err != nil {
		log.Error("d.httpCli.Raw.error(%v)", err)
		return
	}
	uposResp := &UposResp{}
	if err = json.Unmarshal(resp, &uposResp); err != nil {
		log.Error("params(%s),error(%v)", resp, err)
		return
	}
	if uposResp.Code != 0 {
		err = ecode.SubtitleWaveFormFailed
		log.Error("d.Upos,error(%v),info(%s)", err, uposResp.Message)
	}
	return
}
