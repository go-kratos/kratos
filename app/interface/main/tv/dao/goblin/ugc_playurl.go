package goblin

import (
	"context"
	"encoding/json"
	xhttp "net/http"
	"net/url"

	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_httpHeaderRemoteIP = "x-backend-bili-real-ip"
)

// UgcPlayurl is use for get ugc play url
func (d *Dao) UgcPlayurl(ctx context.Context, p *model.PlayURLReq) (res map[string]interface{}, resp *model.PlayURLResp, err error) {
	var (
		params = url.Values{}
		url    = d.conf.Host.UgcPlayURL
		bs     []byte
		req    *xhttp.Request
		ip     = metadata.String(ctx, metadata.RemoteIP)
	)
	res = make(map[string]interface{})
	params.Set("platform", p.Platform)
	params.Set("device", p.Device)
	params.Set("expire", p.Expire)
	params.Set("build", p.Build)
	params.Set("mid", p.Mid)
	params.Set("qn", p.Qn)
	params.Set("npcybs", p.Npcybs)
	params.Set("buvid", p.Buvid)
	params.Set("otype", "json")
	params.Set("trackPath", p.TrackPath)
	params.Set("cid", p.Cid)
	params.Set("access_key", p.AccessKey)
	params.Set("platform", "tvproj")
	if req, err = d.client.NewRequest(xhttp.MethodGet, url, ip, params); err != nil {
		return
	}
	if ip != "" { // add ip into header
		req.Header.Set(_httpHeaderRemoteIP, ip)
	}
	log.Info("ugcPlayURL Cid %d, IP %s", p.Cid, ip)
	if bs, err = d.client.Raw(ctx, req); err != nil {
		log.Error("ugcPl URL %s, Cid %d, Client Raw Err %v", url, p.Cid, err)
		return
	}
	if err = json.Unmarshal(bs, &resp); err != nil { // json unmarshal to struct, to detect error
		log.Error("ugcPl URL %s, Cid %d, Json Unmarshal %s, Err %v", url, p.Cid, string(bs), err)
		return
	}
	if resp.Code != ecode.OK.Code() {
		log.Error("ugcPl URL %s, Cid %d, Resp Code %d, Msg %s", url, p.Cid, resp.Code, resp.Message)
		err = ecode.TvVideoNotFound
		return
	}
	err = json.Unmarshal(bs, &res)
	return
}
