package recommend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_rcmd           = "/pegasus/feed/%d"
	_hot            = "/data/rank/reco-tmzb.json"
	_group          = "/group_changes/pegasus.json"
	_top            = "/feed/tag/top"
	_followModeList = "/data/rank/others/followmode_whitelist.json"
)

// Recommend is.
func (d *Dao) Recommend(c context.Context, plat int8, buvid string, mid int64, build, loginEvent, parentMode, recsysMode int, zoneID int64, group int, interest, network string, style int, column model.ColumnStatus, flush int, autoplay string, now time.Time) (rs []*ai.Item, userFeature json.RawMessage, respCode int, newUser bool, err error) {
	if mid == 0 && buvid == "" {
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	uri := fmt.Sprintf(d.rcmd, group)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("buvid", buvid)
	params.Set("plat", strconv.Itoa(int(plat)))
	params.Set("build", strconv.Itoa(build))
	params.Set("login_event", strconv.Itoa(loginEvent))
	params.Set("zone_id", strconv.FormatInt(zoneID, 10))
	params.Set("interest", interest)
	params.Set("network", network)
	if column > -1 {
		params.Set("column", strconv.Itoa(int(column)))
	}
	params.Set("style", strconv.Itoa(style))
	params.Set("flush", strconv.Itoa(flush))
	params.Set("autoplay_card", autoplay)
	params.Set("parent_mode", strconv.Itoa(parentMode))
	params.Set("recsys_mode", strconv.Itoa(recsysMode))
	var res struct {
		Code        int             `json:"code"`
		NewUser     bool            `json:"new_user"`
		UserFeature json.RawMessage `json:"user_feature"`
		Data        []*ai.Item      `json:"data"`
	}
	if err = d.client.Get(c, uri, ip, params, &res); err != nil {
		respCode = ecode.ServerErr.Code()
		return
	}
	code := ecode.Int(res.Code)
	if !code.Equal(ecode.OK) {
		respCode = res.Code
		err = errors.Wrap(code, uri+"?"+params.Encode())
		return
	}
	rs = res.Data
	userFeature = res.UserFeature
	newUser = res.NewUser
	return
}

// Hots is.
func (d *Dao) Hots(c context.Context) (aids []int64, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	var res struct {
		Code int `json:"code"`
		List []struct {
			Aid int64 `json:"aid"`
		} `json:"list"`
	}
	if err = d.clientAsyn.Get(c, d.hot, ip, nil, &res); err != nil {
		return
	}
	code := ecode.Int(res.Code)
	if !code.Equal(ecode.OK) {
		err = errors.Wrap(code, d.hot)
		return
	}
	for _, list := range res.List {
		if list.Aid != 0 {
			aids = append(aids, list.Aid)
		}
	}
	return
}

// TagTop is.
func (d *Dao) TagTop(c context.Context, mid, tid int64, rn int) (aids []int64, err error) {
	params := url.Values{}
	params.Set("src", "2")
	params.Set("pn", "1")
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("tag", strconv.FormatInt(tid, 10))
	params.Set("rn", strconv.Itoa(rn))
	var res struct {
		Code int     `json:"code"`
		Data []int64 `json:"data"`
	}
	if err = d.client.Get(c, d.top, "", params, &res); err != nil {
		return
	}
	code := ecode.Int(res.Code)
	if !code.Equal(ecode.OK) {
		err = errors.Wrap(code, d.top+"?"+params.Encode())
		return
	}
	aids = res.Data
	return
}

// Group is.
func (d *Dao) Group(c context.Context) (gm map[int64]int, err error) {
	err = d.clientAsyn.Get(c, d.group, "", nil, &gm)
	return
}

// FollowModeList is.
func (d *Dao) FollowModeList(c context.Context) (list map[int64]struct{}, err error) {
	var res struct {
		Code int     `json:"code"`
		Data []int64 `json:"data"`
	}
	if err = d.clientAsyn.Get(c, d.followModeList, "", nil, &res); err != nil {
		return
	}
	code := ecode.Int(res.Code)
	if !code.Equal(ecode.OK) {
		err = errors.Wrap(code, d.followModeList)
		return
	}
	b, _ := json.Marshal(&res)
	log.Warn("FollowModeList param(%s) res(%s)", b, d.followModeList)
	list = make(map[int64]struct{}, len(res.Data))
	for _, mid := range res.Data {
		list[mid] = struct{}{}
	}
	return
}
