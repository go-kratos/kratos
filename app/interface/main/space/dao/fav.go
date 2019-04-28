package dao

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"go-common/app/interface/main/space/model"
	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_favFolderURI  = "/x/internal/v2/fav/folder"
	_favArchiveURI = "/x/internal/v2/fav/video"
	_favAlbumURI   = "/userext/v1/Fav/getMyFav"
	_favMovieURI   = "/follow/api/list/mine"
	_samplePage    = "1"
	_samplePs      = "1"
)

// FavFolder favorite folder list.
func (d *Dao) FavFolder(c context.Context, mid, vmid int64) (res []*favmdl.VideoFolder, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("vmid", strconv.FormatInt(vmid, 10))
	var rs struct {
		Code int                   `json:"code"`
		Data []*favmdl.VideoFolder `json:"data"`
	}
	if err = d.httpR.Get(c, d.favFolderURL, ip, params, &rs); err != nil {
		log.Error("d.http.Get(%s,%d) error(%v)", d.favFolderURL, mid, err)
		return
	}
	if rs.Code != ecode.OK.Code() {
		log.Error("d.http.Get(%s,%d) code(%d)", d.favFolderURL, mid, rs.Code)
		err = ecode.Int(rs.Code)
		return
	}
	res = rs.Data
	return
}

// LiveFavCount get live(vc or album) favorite count.
func (d *Dao) LiveFavCount(c context.Context, mid int64, favType int) (count int, err error) {
	var (
		req *http.Request
		rs  struct {
			Code int `json:"code"`
			Data struct {
				PageInfo struct {
					Page      int `json:"page"`
					PageSize  int `json:"page_size"`
					TotalPage int `json:"total_page"`
					Count     int `json:"count"`
				} `json:"pageinfo"`
			} `json:"data"`
		}
		ip = metadata.String(c, metadata.RemoteIP)
	)
	params := url.Values{}
	params.Set("biz_type", strconv.Itoa(favType))
	if req, err = d.httpR.NewRequest("GET", d.favAlbumURL, ip, params); err != nil {
		log.Error("d.httpR.NewRequest %s error(%v)", d.favAlbumURL, err)
		return
	}
	req.Header.Set("X-BILILIVE-UID", strconv.FormatInt(mid, 10))
	if err = d.httpR.Do(c, req, &rs); err != nil {
		log.Error("d.httpR.Get(%s,%d) error(%v)", d.favAlbumURL, mid, err)
		return
	}
	if rs.Code != ecode.OK.Code() {
		log.Error("d.httpR.Get(%s,%d) code(%d)", d.favAlbumURL, mid, rs.Code)
		err = ecode.Int(rs.Code)
		return
	}
	count = rs.Data.PageInfo.Count
	return
}

// MovieFavCount get movie fav count
func (d *Dao) MovieFavCount(c context.Context, mid int64, favType int) (count int, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("season_type", strconv.Itoa(favType))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("page", _samplePage)
	params.Set("pagesize", _samplePs)
	params.Set("build", _build)
	params.Set("platform", _platform)
	var (
		rs struct {
			Code  int    `json:"code"`
			Count string `json:"count"`
		}
	)
	if err = d.httpR.Get(c, d.favMovieURL, ip, params, &rs); err != nil {
		log.Error("d.httpR.Get(%s,%d) error(%v)", d.favMovieURL, mid, err)
		return
	}
	if rs.Code != ecode.OK.Code() {
		log.Error("d.httpR.Get(%s,%d) code(%d)", d.favMovieURL, mid, rs.Code)
		err = ecode.Int(rs.Code)
		return
	}
	count, _ = strconv.Atoi(rs.Count)
	return
}

// FavArchive fav archive.
func (d *Dao) FavArchive(c context.Context, mid int64, arg *model.FavArcArg) (res *favmdl.SearchArchive, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	if mid > 0 {
		params.Set("mid", strconv.FormatInt(mid, 10))
	}
	params.Set("vmid", strconv.FormatInt(arg.Vmid, 10))
	params.Set("fid", strconv.FormatInt(arg.Fid, 10))
	if arg.Tid > 0 {
		params.Set("tid", strconv.FormatInt(arg.Tid, 10))
	}
	if arg.Keyword != "" {
		params.Set("keyword", arg.Keyword)
	}
	if arg.Order != "" {
		params.Set("order", arg.Order)
	}
	params.Set("pn", strconv.Itoa(arg.Pn))
	params.Set("ps", strconv.Itoa(arg.Ps))
	var rs struct {
		Code int                   `json:"code"`
		Data *favmdl.SearchArchive `json:"data"`
	}
	if err = d.httpR.Get(c, d.favArcURL, ip, params, &rs); err != nil {
		log.Error("d.http.Get(%s,%d) error(%v)", d.favArcURL, mid, err)
		return
	}
	if rs.Code != ecode.OK.Code() {
		log.Error("d.http.Get(%s,%d) code(%d)", d.favArcURL, mid, rs.Code)
		err = ecode.Int(rs.Code)
		return
	}
	res = rs.Data
	return
}
