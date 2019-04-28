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
	_tagSubURI       = "/x/internal/tag/subscribe/add"
	_tagCancelSubURI = "/x/internal/tag/subscribe/cancel"
	_subTagListURI   = "/x/internal/tag/subscribe/tags"
)

// TagSub subscribe tag.
func (d *Dao) TagSub(c context.Context, mid, tid int64) (err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("tag_id", strconv.FormatInt(tid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpW.Post(c, d.tagSubURL, ip, params, &res); err != nil {
		log.Error("tag: d.httpW.Post url(%s) mid(%d) tag_id(%d) error(%v)", d.tagSubURL, mid, tid, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("tag: d.http.Do code error(%d)", res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// TagCancelSub cancel subscribe tag.
func (d *Dao) TagCancelSub(c context.Context, mid, tid int64) (err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("tag_id", strconv.FormatInt(tid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpW.Post(c, d.tagCancelSubURL, ip, params, &res); err != nil {
		log.Error("tag: d.httpW.Post url(%s) mid(%d) tag_id(%d) error(%v)", d.tagCancelSubURL, mid, tid, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("tag: d.http.Do code error(%d)", res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// TagSubList get tag subscribe list by mid.
func (d *Dao) TagSubList(c context.Context, mid int64, pn, ps int) (rs []*model.Tag, total int, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("vmid", strconv.FormatInt(mid, 10))
	params.Set("pn", strconv.Itoa(pn))
	params.Set("ps", strconv.Itoa(ps))
	var res struct {
		Code  int          `json:"code"`
		Total int          `json:"total"`
		Data  []*model.Tag `json:"data"`
	}
	if err = d.httpR.Get(c, d.tagSubListURL, ip, params, &res); err != nil {
		log.Error("d.http.Get(%s,%d) error(%v)", d.tagSubListURL, mid, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.http.Get(%s,%d) error(%v)", d.tagSubListURL, mid, err)
		err = ecode.Int(res.Code)
		return
	}
	rs = res.Data
	total = res.Total
	return
}
