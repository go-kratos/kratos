package bplus

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-interface/model/bplus"
	"go-common/library/ecode"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_favorPlus     = "/user_ex/v1/Fav/getFavList"
	_clips         = "/clip/v1/video/blist"
	_allbums       = "/link_draw/v1/Doc/photo_list_ones"
	_allClip       = "/clip_ext/v1/video/all"
	_allAlbum      = "/link_draw/v1/Doc/photo_all_ones"
	_clipDetail    = "/clip_ext/v0/video/getDetailForSpace"
	_albumDetail   = "/link_draw_ex/v0/Doc/details"
	_groupsCount   = "/link_group/v1/member/created_groups_num"
	_dynamic       = "/dynamic_svr/v0/dynamic_svr/space_intro"
	_dunamicCount  = "/dynamic_svr/v0/dynamic_svr/space_dy_num"
	_dynamicDetail = "/dynamic_detail/v0/Dynamic/details"
)

// DynamicCount return dynamic count
func (d *Dao) DynamicCount(c context.Context, mid int64) (count int64, err error) {
	params := url.Values{}
	params.Set("uids", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Items []struct {
				Mid int64 `json:"uid"`
				Num int64 `json:"num"`
			} `json:"items"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.dynamicCount, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.dynamicCount+"?"+params.Encode())
		return
	}
	if len(res.Data.Items) == 0 {
		return
	}
	for _, item := range res.Data.Items {
		if item.Mid != mid {
			continue
		}
		count = item.Num
		break
	}
	return
}

// FavClips get fav from B+ api.
func (d *Dao) FavClips(c context.Context, mid int64, accessKey, actionKey, device, mobiApp, platform string, build, pn, ps int) (cs *bplus.Clips, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("access_key", accessKey)
	params.Set("actionKey", actionKey)
	params.Set("build", strconv.Itoa(build))
	params.Set("device", device)
	params.Set("mobi_app", mobiApp)
	params.Set("platform", platform)
	params.Set("biz_type", strconv.Itoa(bplus.CLIPS))
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	var res struct {
		Code    int          `json:"code"`
		Msg     string       `json:"msg"`
		Message string       `json:"message"`
		Data    *bplus.Clips `json:"data"`
	}
	if err = d.client.Get(c, d.favorPlus, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.favorPlus+"?"+params.Encode())
		return
	}
	cs = res.Data
	return
}

// FavAlbums get fav from B+ api.
func (d *Dao) FavAlbums(c context.Context, mid int64, accessKey, actionKey, device, mobiApp, platform string, build, pn, ps int) (as *bplus.Albums, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("access_key", accessKey)
	params.Set("actionKey", actionKey)
	params.Set("build", strconv.Itoa(build))
	params.Set("device", device)
	params.Set("mobi_app", mobiApp)
	params.Set("platform", platform)
	params.Set("biz_type", strconv.Itoa(bplus.ALBUMS))
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	var res struct {
		Code    int           `json:"code"`
		Msg     string        `json:"msg"`
		Message string        `json:"message"`
		Data    *bplus.Albums `json:"data"`
	}
	if err = d.client.Get(c, d.favorPlus, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.favorPlus+"?"+params.Encode())
		return
	}
	as = res.Data
	return
}

// Clips .
func (d *Dao) Clips(c context.Context, vmid int64, pos, size int) (cs []*bplus.Clip, more, offset int, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(vmid, 10))
	params.Set("next_offset", strconv.Itoa(pos))
	params.Set("page_size", strconv.Itoa(size))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			More   int           `json:"has_more"`
			Offset int           `json:"next_offset"`
			Item   []*bplus.Clip `json:"items"`
		}
	}
	if err = d.client.Get(c, d.clips, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.clips+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		cs = res.Data.Item
		more = res.Data.More
		offset = res.Data.Offset
	}
	return
}

// Albums get album list form api .
func (d *Dao) Albums(c context.Context, vmid int64, pos, size int) (as []*bplus.Album, more, offset int, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(vmid, 10))
	params.Set("next_offset", strconv.Itoa(pos))
	params.Set("page_size", strconv.Itoa(size))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			More   int            `json:"has_more"`
			Offset int            `json:"next_offset"`
			Item   []*bplus.Album `json:"items"`
		}
	}
	if err = d.client.Get(c, d.albums, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.albums+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		as = res.Data.Item
		more = res.Data.More
		offset = res.Data.Offset
	}
	return
}

// AllClip .
func (d *Dao) AllClip(c context.Context, vmid int64, size int) (cs []*bplus.Clip, count int, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(vmid, 10))
	params.Set("page_size", strconv.Itoa(size))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Count int           `json:"total_count"`
			Item  []*bplus.Clip `json:"items"`
		}
	}
	if err = d.client.Get(c, d.allClip, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.allClip+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		cs = res.Data.Item
		count = res.Data.Count
	}
	return
}

// AllAlbum .
func (d *Dao) AllAlbum(c context.Context, vmid int64, size int) (as []*bplus.Album, count int, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(vmid, 10))
	params.Set("page_size", strconv.Itoa(size))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Count int            `json:"total_count"`
			Item  []*bplus.Album `json:"items"`
		}
	}
	if err = d.client.Get(c, d.allAlbum, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.allAlbum+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		as = res.Data.Item
		count = res.Data.Count
	}
	return
}

// ClipDetail .
func (d *Dao) ClipDetail(c context.Context, ids []int64) (cs map[int64]*bplus.Clip, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("ids", xstr.JoinInts(ids))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Item []*bplus.Clip `json:"items"`
		}
	}
	if err = d.client.Get(c, d.clipDetail, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.clipDetail+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		cs = make(map[int64]*bplus.Clip, len(res.Data.Item))
		for _, clip := range res.Data.Item {
			cs[clip.ID] = clip
		}
	}
	return
}

// AlbumDetail .
func (d *Dao) AlbumDetail(c context.Context, vmid int64, ids []int64) (as map[int64]*bplus.Album, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(vmid, 10))
	params.Set("ids", xstr.JoinInts(ids))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Item []*bplus.Album `json:"items"`
		}
	}
	if err = d.client.Get(c, d.albumDetail, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.albumDetail+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		as = make(map[int64]*bplus.Album, len(res.Data.Item))
		for _, album := range res.Data.Item {
			as[album.ID] = album
		}
	}
	return
}

// GroupsCount .
func (d *Dao) GroupsCount(c context.Context, mid, vmid int64) (count int, err error) {
	var (
		req *http.Request
		ip  = metadata.String(c, metadata.RemoteIP)
	)
	params := url.Values{}
	params.Set("master_uid", strconv.FormatInt(vmid, 10))
	if req, err = d.client.NewRequest(http.MethodGet, d.groupsCount, ip, params); err != nil {
		return
	}
	req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Num int `json:"num"`
		}
	}
	if err = d.client.Do(c, req, &res); err != nil {
		err = errors.Wrapf(err, "url(%s) header(X-BiliLive-UID:%s)", req.URL.String(), req.Header.Get("X-BiliLive-UID"))
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrapf(ecode.Int(res.Code), "url(%s) header(X-BiliLive-UID:%s)", req.URL.String(), req.Header.Get("X-BiliLive-UID"))
		return
	}
	if res.Data != nil {
		count = res.Data.Num
	}
	return
}

// Dynamic .
func (d *Dao) Dynamic(c context.Context, uid int64) (has bool, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(uid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Reuslt int `json:"result"`
		}
	}
	if err = d.client.Get(c, d.dynamic, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.dynamic+"?"+params.Encode())
		return
	}
	if res.Data != nil && res.Data.Reuslt == 1 {
		has = true
	}
	return
}

// DynamicDetails get dynamic details by ids.
func (d *Dao) DynamicDetails(c context.Context, ids []int64, from string) (details map[int64]*bplus.Detail, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("from", from)
	for _, id := range ids {
		params.Add("dynamic_ids[]", strconv.FormatInt(id, 10))
	}
	var res struct {
		Code int `json:"code"`
		Data *struct {
			List []*bplus.Detail `json:"list"`
		} `json:"data"`
	}
	details = make(map[int64]*bplus.Detail)
	if err = d.client.Get(c, d.dynamicDetail, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.dynamicDetail+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		for _, detail := range res.Data.List {
			if detail.ID != 0 {
				details[detail.ID] = detail
			}
		}
	}
	return
}
