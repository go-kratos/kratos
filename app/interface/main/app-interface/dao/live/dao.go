package live

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-interface/conf"
	"go-common/app/interface/main/app-interface/model/live"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_live        = "/AppRoom/getRoomInfo"
	_medalStatus = "/fans_medal/v1/medal/get_medal_opened"
	_appMRoom    = "/room/v1/Room/rooms_for_app_index"
	_statusInfo  = "/room/v1/Room/get_status_info_by_uids"
	_visibleInfo = "/rc/v1/Glory/get_visible"
	_usersInfo   = "/user/v3/User/getMultiple"
	_LiveByRID   = "/room/v2/Room/get_by_ids"
)

// Dao is space dao
type Dao struct {
	client      *httpx.Client
	live        string
	medalStatus string
	appMRoom    string
	statusInfo  string
	visibleInfo string
	userInfo    string
	liveByRID   string
}

// New initial space dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:      httpx.NewClient(c.HTTPLive),
		live:        c.Host.APILiveCo + _live,
		medalStatus: c.Host.APILiveCo + _medalStatus,
		appMRoom:    c.Host.APILiveCo + _appMRoom,
		statusInfo:  c.Host.APILiveCo + _statusInfo,
		visibleInfo: c.Host.APILiveCo + _visibleInfo,
		userInfo:    c.Host.APILiveCo + _usersInfo,
		liveByRID:   c.Host.APILiveCo + _LiveByRID,
	}
	return
}

// Live is space live data.
func (d *Dao) Live(c context.Context, mid int64, platform string) (live json.RawMessage, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("platform", platform)
	var res struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err = d.client.Get(c, d.live, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.live+"?"+params.Encode())
		return
	}
	live = res.Data
	return
}

// MedalStatus for live
func (d *Dao) MedalStatus(c context.Context, mid int64) (status int, err error) {
	var (
		req *http.Request
		res struct {
			Code int `json:"code"`
			Data *struct {
				MasterStatus int `json:"master_status"`
			} `json:"data"`
		}
		ip = metadata.String(c, metadata.RemoteIP)
	)
	if req, err = d.client.NewRequest("GET", d.medalStatus, ip, nil); err != nil {
		return
	}
	req.Header.Set("X-BILILIVE-UID", strconv.FormatInt(mid, 10))
	if err = d.client.Do(c, req, &res); err != nil {
		err = errors.Wrapf(err, "url(%s) header(X-BiliLive-UID:%s)", req.URL.String(), req.Header.Get("X-BiliLive-UID"))
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrapf(ecode.Int(res.Code), "url(%s) header(X-BiliLive-UID:%s)", req.URL.String(), req.Header.Get("X-BiliLive-UID"))
		return
	}
	if res.Data != nil {
		status = res.Data.MasterStatus
	}
	return
}

// AppMRoom for live
func (d *Dao) AppMRoom(c context.Context, roomids []int64) (rs map[int64]*live.Room, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("room_ids", xstr.JoinInts(roomids))
	var res struct {
		Code int          `json:"code"`
		Data []*live.Room `json:"data"`
	}
	if err = d.client.Get(c, d.appMRoom, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.appMRoom+"?"+params.Encode())
		return
	}
	rs = make(map[int64]*live.Room, len(res.Data))
	for _, r := range res.Data {
		rs[r.RoomID] = r
	}
	return
}

// StatusInfo for live
func (d *Dao) StatusInfo(c context.Context, mids []int64) (status map[int64]*live.Status, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	for _, mid := range mids {
		params.Add("uids[]", strconv.FormatInt(mid, 10))
	}
	params.Set("filter_offline", "1")
	var res struct {
		Code int                    `json:"code"`
		Data map[int64]*live.Status `json:"data"`
	}
	if err = d.client.Get(c, d.statusInfo, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.statusInfo+"?"+params.Encode())
		return
	}
	status = res.Data
	return
}

// Glory for live search
func (d *Dao) Glory(c context.Context, uid int64) (glory []*live.Glory, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(uid, 10))
	var res struct {
		Code int           `json:"code"`
		Data []*live.Glory `json:"data"`
	}
	if err = d.client.Get(c, d.visibleInfo, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.visibleInfo+"?"+params.Encode())
		return
	}
	glory = res.Data
	return
}

// UserInfo for live search
func (d *Dao) UserInfo(c context.Context, uids []int64) (userInfo map[int64]map[string]*live.Exp, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	for _, uid := range uids {
		params.Set("uids[]", strconv.FormatInt(uid, 10))
	}
	params.Set("attributes[]", "exp")
	var res struct {
		Code int                            `json:"code"`
		Data map[int64]map[string]*live.Exp `json:"data"`
	}
	if err = d.client.Get(c, d.userInfo, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.userInfo+"?"+params.Encode())
		return
	}
	userInfo = res.Data
	return
}

// LiveByRIDs get live info by room_ids.
func (d *Dao) LiveByRIDs(c context.Context, roomIDs []int64) (info map[int64]*live.RoomInfo, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	for _, id := range roomIDs {
		params.Add("ids[]", strconv.FormatInt(id, 10))
	}
	params.Add("fields[]", "roomid")
	params.Add("fields[]", "title")
	params.Add("fields[]", "cover")
	params.Add("fields[]", "user_cover")
	params.Add("fields[]", "uid")
	params.Add("fields[]", "uname")
	params.Add("fields[]", "area_v2_name")
	params.Add("fields[]", "live_status")
	params.Add("fields[]", "broadcast_type")
	params.Add("fields[]", "short_id")
	params.Add("need_broadcast_type", "1")
	var res struct {
		Code int                      `json:"code"`
		Data map[int64]*live.RoomInfo `json:"data"`
	}
	if err = d.client.Get(c, d.liveByRID, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.liveByRID+"?"+params.Encode())
		return
	}
	info = res.Data
	return
}
