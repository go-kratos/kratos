package archive

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-view/model/view"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_realteURL     = "/recsys/related"
	_commercialURL = "/x/internal/creative/arc/commercial"
	_relateRecURL  = "/recommand"
	_playURL       = "/playurl/batch"
)

// RelateAids get relate by aid
func (d *Dao) RelateAids(c context.Context, aid int64) (aids []int64, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("key", strconv.FormatInt(aid, 10))
	var res struct {
		Code int `json:"code"`
		Data []*struct {
			Value string `json:"value"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.realteURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.realteURL+"?"+params.Encode())
		return
	}
	if len(res.Data) != 0 {
		if aids, err = xstr.SplitInts(res.Data[0].Value); err != nil {
			err = errors.Wrap(err, res.Data[0].Value)
		}
	}
	return
}

// Commercial is
func (d *Dao) Commercial(c context.Context, aid int64) (gameID int64, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			GameID int64 `json:"game_id"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.commercialURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.commercialURL+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		gameID = res.Data.GameID
	}
	return
}

// NewRelateAids relate online recommend 在线实时推荐
func (d *Dao) NewRelateAids(c context.Context, aid, mid int64, build, parentMode int, buvid, from string, plat int8) (rec []*view.NewRelateRec, userFeature, returnCode string, dalaoExp int, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("from", "2")
	params.Set("cmd", "related")
	params.Set("timeout", "100")
	params.Set("plat", strconv.Itoa(int(plat)))
	params.Set("build", strconv.Itoa(build))
	params.Set("buvid", buvid)
	params.Set("from_av", strconv.FormatInt(aid, 10))
	params.Set("request_cnt", "40")
	params.Set("source_page", from)
	params.Set("parent_mode", strconv.Itoa(parentMode))
	params.Set("need_dalao", "1")
	var res struct {
		Code        int                  `json:"code"`
		Data        []*view.NewRelateRec `json:"data"`
		UserFeature string               `json:"user_feature"`
		DalaoExp    int                  `json:"dalao_exp"`
	}
	log.Warn("dalaotest url(%s)", d.relateRecURL+"?"+params.Encode())
	if err = d.client.Get(c, d.relateRecURL, ip, params, &res); err != nil {
		returnCode = "500"
		return
	}
	returnCode = strconv.Itoa(res.Code)
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.relateRecURL+"?"+params.Encode())
		return
	}
	dalaoExp = res.DalaoExp
	userFeature = res.UserFeature
	rec = res.Data
	return
}

// PlayerInfos cid with player info
func (d *Dao) PlayerInfos(c context.Context, cids []int64, qn, fnver, fnval, forceHost int, platform string) (pm map[uint32]*archive.BvcVideoItem, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("cid", xstr.JoinInts(cids))
	params.Set("qn", strconv.Itoa(qn))
	params.Set("platform", platform)
	params.Set("uip", ip)
	params.Set("layout", "pb")
	params.Set("fnver", strconv.Itoa(fnver))
	params.Set("fnval", strconv.Itoa(fnval))
	params.Set("force_host", strconv.Itoa(forceHost))
	var req *http.Request
	if req, err = d.client.NewRequest("GET", d.playURL, ip, params); err != nil {
		return
	}
	res := new(archive.BvcResponseMsg)
	if err = d.httpClient.PB(c, req, res); err != nil {
		return
	}
	if int(res.Code) != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(int(res.Code)), d.playURL+params.Encode())
		return
	}
	pm = res.Data
	return
}
