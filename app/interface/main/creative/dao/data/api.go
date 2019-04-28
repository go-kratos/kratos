package data

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"go-common/app/interface/main/creative/model/data"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_statURL  = "/data/member/upinfo/%d.json"
	_tagv2URL = "/tag_predict/v2"
	_coverURL = "/cover_recomm"
)

// Stat get user stat play/fans/...
func (d *Dao) stat(c context.Context, mid int64, ip string) (st *data.Stat, err error) {
	var res struct {
		Code int `json:"code"`
		*data.Stat
	}
	if err = d.client.RESTfulGet(c, d.statURI, ip, url.Values{}, &res, mid); err != nil {
		log.Error("statURI url(%s) mid(%d) error(%v)", d.statURI, mid, err)
		err = ecode.CreativeDataErr
		return
	}
	if res.Code != 0 {
		log.Error("data url(%s) res(%v)", d.statURI, res)
		err = ecode.CreativeDataErr
		return
	}
	st = res.Stat
	return
}

// TagsWithChecked get predict tag with checked mark.
func (d *Dao) TagsWithChecked(c context.Context, mid int64, tid uint16, title, filename, desc, cover string, tagFrom int8) (t []*data.CheckedTag, err error) {
	params := url.Values{}
	t = make([]*data.CheckedTag, 0)
	params.Set("client_type", strconv.Itoa(int(tagFrom)))
	params.Set("typeid", strconv.Itoa(int(tid)))
	params.Set("title", title)
	params.Set("filename", filename)
	params.Set("desc", desc)
	params.Set("cover", cover)
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Tags []*data.CheckedTag `json:"tags"`
		} `json:"data"`
	}
	log.Info("TagsWithChecked url(%s)", d.tagV2URI+"?"+params.Encode())
	if err = d.client.Get(c, d.tagV2URI, "", params, &res); err != nil {
		log.Error("data url(%s) error(%v)", d.tagV2URI+"?"+params.Encode(), err)
		err = ecode.CreativeDataErr
		return
	}
	if res.Code != 0 {
		log.Error("data url(%s) res(%v)", d.tagV2URI+"?"+params.Encode(), res)
		err = ecode.CreativeDataErr
		return
	}
	t = res.Data.Tags
	return
}

// RecommendCovers get recommend covers from AI.
func (d *Dao) RecommendCovers(c context.Context, mid int64, fns []string) (cvs []string, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("filename", strings.Join(fns, ","))
	var res struct {
		Code int      `json:"code"`
		Data []string `json:"data"`
	}
	if err = d.client.Get(c, d.coverBFSURI, "", params, &res); err != nil {
		log.Error("Covers url(%s) error(%v)", d.coverBFSURI+"?"+params.Encode(), err)
		err = ecode.CreativeDataErr
		return
	}
	if res.Code != 0 {
		log.Error("Covers url(%s) res(%v)", d.coverBFSURI+"?"+params.Encode(), res)
		err = ecode.CreativeDataErr
		return
	}
	cvs = res.Data
	return
}
