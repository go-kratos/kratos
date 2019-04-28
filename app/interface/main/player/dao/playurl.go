package dao

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"go-common/app/interface/main/player/model"
	"go-common/library/ecode"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Playurl get playurl data.
func (d *Dao) Playurl(c context.Context, mid int64, arg *model.PlayurlArg, playurl, token string) (data *model.PlayurlRes, err error) {
	params := url.Values{}
	params.Set("cid", strconv.FormatInt(arg.Cid, 10))
	params.Set("avid", strconv.FormatInt(arg.Aid, 10))
	params.Set("qn", strconv.Itoa(arg.Qn))
	if arg.Type != "" {
		params.Set("type", arg.Type)
	}
	if arg.MaxBackup > 0 {
		params.Set("max_backup", strconv.Itoa(arg.MaxBackup))
	}
	if arg.Npcybs > 0 {
		params.Set("npcybs", strconv.Itoa(arg.Npcybs))
	}
	if mid > 0 {
		params.Set("mid", strconv.FormatInt(mid, 10))
	}
	if arg.Platform != "" {
		params.Set("platform", arg.Platform)
	}
	if arg.Buvid != "" {
		params.Set("buvid", arg.Buvid)
	}
	if arg.Resolution != "" {
		params.Set("resolution", arg.Resolution)
	}
	if arg.Model != "" {
		params.Set("model", arg.Model)
	}
	if arg.Build > 0 {
		params.Set("build", strconv.Itoa(arg.Build))
	}
	params.Set("fnver", strconv.Itoa(arg.Fnver))
	if arg.Fnval > 0 {
		params.Set("fnval", strconv.Itoa(arg.Fnval))
	}
	if arg.Session != "" {
		params.Set("session", arg.Session)
	}
	if arg.HTML5 > 0 && arg.H5GoodQuality > 0 {
		params.Set("h5_good_quality", strconv.Itoa(arg.H5GoodQuality))
	}
	params.Set("otype", "json")
	var req *http.Request
	if req, err = d.client.NewRequest(http.MethodGet, playurl, metadata.String(c, metadata.RemoteIP), params); err != nil {
		return
	}
	if token != "" {
		req.Header.Set("X-BVC-FINGERPRINT", token)
	}
	var res struct {
		Code int `json:"code"`
		*model.PlayurlRes
	}
	if err = d.client.Do(c, req, &res); err != nil {
		err = errors.Wrapf(err, "d.client.Do(%s) (%+v)", playurl+"?"+params.Encode(), arg)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), playurl+"?"+params.Encode())
	}
	data = res.PlayurlRes
	return
}
