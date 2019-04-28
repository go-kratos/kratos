package archive

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

// ValidateQn validate qn
func (d *Dao) ValidateQn(c context.Context, qn int) (allow bool) {
	_, allow = d.playerQn[qn]
	return
}

// VipQn vip qn
func (d *Dao) VipQn(c context.Context, qn int) (isVipQn bool) {
	_, isVipQn = d.playerVipQn[qn]
	return
}

// PlayerInfos cid with player info
func (d *Dao) PlayerInfos(c context.Context, cids []int64, qn int, platform, ip string, fnver, fnval, forceHost int) (pm map[uint32]*archive.BvcVideoItem, err error) {
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
	if req, err = d.playerClient.NewRequest("GET", d.c.PlayerAPI, ip, params); err != nil {
		return
	}
	res := new(archive.BvcResponseMsg)
	if err = d.playerClient.PB(c, req, res); err != nil {
		return
	}
	if int(res.Code) != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(int(res.Code)), d.c.PlayerAPI+params.Encode())
		return
	}
	pm = res.Data
	return
}

// PGCPlayerInfos cid with pgc player info
func (d *Dao) PGCPlayerInfos(c context.Context, aids []int64, platform, ip, session string, fnval, fnver int) (pgcm map[int64]*archive.PlayerInfo, err error) {
	params := url.Values{}
	params.Set("aids", xstr.JoinInts(aids))
	params.Set("mobi_app", platform)
	params.Set("ip", ip)
	params.Set("fnver", strconv.Itoa(fnver))
	params.Set("fnval", strconv.Itoa(fnval))
	params.Set("session", session)
	res := struct {
		Code   int                          `json:"code"`
		Result map[int64]*archive.PGCPlayer `json:"result"`
	}{}
	if err = d.playerClient.Get(c, d.c.PGCPlayerAPI, ip, params, &res); err != nil {
		return
	}
	if res.Code != 0 {
		err = errors.Wrap(ecode.Int(res.Code), d.c.PGCPlayerAPI+params.Encode())
		return
	}
	pgcm = make(map[int64]*archive.PlayerInfo)
	for _, v := range res.Result {
		if v.PlayerInfo != nil {
			pgcm[v.PlayerInfo.Cid] = v.PlayerInfo
		}
	}
	return
}
