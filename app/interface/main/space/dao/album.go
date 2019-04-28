package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/space/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_albumCountURI = "/link_draw/v1/doc/upload_count"
	_albumListURI  = "/link_draw/v1/doc/doc_list"
)

// AlbumCount get album count.
func (d *Dao) AlbumCount(c context.Context, mid int64) (count int64, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("uid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int              `json:"code"`
		Data model.AlbumCount `json:"data"`
	}
	if err = d.httpR.Get(c, d.albumCountURL, ip, params, &res); err != nil {
		log.Error("d.httpR.Get(%s) mid(%d) error(%v)", d.albumCountURL, mid, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpR.Get(%s) mid(%d) code(%d)", d.albumCountURL, mid, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	count = res.Data.AllCount
	return
}

// AlbumList get album list.
func (d *Dao) AlbumList(c context.Context, mid int64, pn, ps int) (list []*model.Album, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("page_num", strconv.Itoa(pn))
	params.Set("page_size", strconv.Itoa(ps))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Items []*model.Album `json:"items"`
		} `json:"data"`
	}
	if err = d.httpR.Get(c, d.albumListURL, ip, params, &res); err != nil {
		log.Error("d.httpR.Get(%s) mid(%d) error(%v)", d.albumListURL, mid, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpR.Get(%s) mid(%d) code(%d)", d.albumListURL, mid, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	list = res.Data.Items
	return
}
