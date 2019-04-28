package space

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/job/main/app/model/space"
	"go-common/library/ecode"

	"github.com/pkg/errors"
)

const (
	_clipList  = "/clip_ext/v0/video/allForSpace"
	_albumList = "/link_draw_ex/v0/Doc/batchids"
	_audioList = "/audio/music-service-c/songs/internal/uppersongs-page"
)

func (d *Dao) ClipList(c context.Context, vmid int64, pos, size int, ip string) (cls []*space.Clip, more, offset int, err error) {
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(vmid, 10))
	params.Set("next_offset", strconv.Itoa(pos))
	params.Set("page_size", strconv.Itoa(size))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			More   int           `json:"has_more"`
			Offset int           `json:"next_offset"`
			Item   []*space.Clip `json:"items"`
		} `json:"data"`
	}
	if err = d.clientAsyn.Get(c, d.clipList, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.clipList+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		more = res.Data.More
		offset = res.Data.Offset
		cls = res.Data.Item
	}
	return
}

func (d *Dao) AlbumList(c context.Context, vmid int64, pos, size int, ip string) (als []*space.Album, more, offset int, err error) {
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(vmid, 10))
	params.Set("next_offset", strconv.Itoa(pos))
	params.Set("page_size", strconv.Itoa(size))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			More   int            `json:"has_more"`
			Offset int            `json:"next_offset"`
			Item   []*space.Album `json:"items"`
		} `json:"data"`
	}
	if err = d.clientAsyn.Get(c, d.albumList, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.albumList+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		more = res.Data.More
		offset = res.Data.Offset
		als = res.Data.Item
	}
	return
}

func (d *Dao) AudioList(c context.Context, vmid int64, pn, ps int, ip string) (aus []*space.Audio, hasNext bool, nextPage int, err error) {
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(vmid, 10))
	params.Set("pageIndex", strconv.Itoa(pn))
	params.Set("pageSize", strconv.Itoa(ps))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			NextPage    int            `json:"nextPage"`
			HasNextPage bool           `json:"hasNextPage"`
			List        []*space.Audio `json:"list"`
		} `json:"data"`
	}
	if err = d.clientAsyn.Get(c, d.audioList, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.audioList+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		aus = res.Data.List
		hasNext = res.Data.HasNextPage
		nextPage = res.Data.NextPage
	}
	return
}
