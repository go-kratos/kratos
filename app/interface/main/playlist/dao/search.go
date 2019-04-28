package dao

import (
	"context"
	"go-common/app/interface/main/playlist/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"net/url"
	"strconv"
)

const (
	_searchURL  = "/main/search"
	_searchVer  = "v3"
	_platform   = "web"
	_searchType = "video"
	_searchFrom = "web_playlist"
)

// SearchVideo get search video.
func (d *Dao) SearchVideo(c context.Context, pn, ps int, query string) (res []*model.SearchArc, count int, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("main_ver", _searchVer)
	params.Set("platform", _platform)
	params.Set("search_type", _searchType)
	params.Set("from_source", _searchFrom)
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("page", strconv.Itoa(pn))
	params.Set("keyword", query)
	params.Set("clientip", ip)
	var rs struct {
		Code       int                `json:"code"`
		Result     []*model.SearchArc `json:"result"`
		NumResults int                `json:"numResults"`
	}
	if err = d.http.Get(c, d.searchURL, ip, params, &rs); err != nil {
		log.Error("d.http.Get(%s,%v) error(%v)", d.searchURL, params, err)
		return
	}
	if rs.Code != ecode.OK.Code() {
		log.Error("d.http.Get(%s,%v) code error(%v)", d.searchURL, params, rs.Code)
		return
	}
	res = rs.Result
	count = rs.NumResults
	return
}
