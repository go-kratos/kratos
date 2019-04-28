package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

const (
	_filterURI       = "/x/internal/filter"
	_accountMoralURI = "/api/moral/add"
	_blockUserURI    = "/api/member/blockAccountWithTime"
	_sendMsgURI      = "/api/notify/send.user.notify.do"
	_regionHotTagURI = "%s/tag/hot-info-%d.json"
	_originURI       = "http://www.bilibili.com/video/av%d/"
	_originContent   = "%s-%s"
	_archiveHotURI   = "/tag_list"
)

// Filter filter.
func (d *Dao) Filter(c context.Context, msg string) (err error) {
	uri := d.hosts.PlatformHost + _filterURI
	params := url.Values{}
	params.Set("area", "tag")
	params.Set("msg", msg)
	res := &struct {
		Code int `json:"code"`
		Data struct {
			Level int `json:"level"`
		} `json:"data"`
	}{}
	if err = d.client.Get(c, uri, metadata.String(c, metadata.RemoteIP), params, res); err != nil {
		log.Error("d.client.Get(%s) error(%v)", uri+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("d.client.Get(%s) error(%v)", uri+"?"+params.Encode(), err)
		return ecode.Int(res.Code)
	}
	if res.Data.Level > 10 {
		err = ecode.TagIsSealing
	}
	return
}

// RegionHot RegionHot.
func (d *Dao) RegionHot(c context.Context, rid int64) ([]*model.BasicTag, int64, error) {
	params := url.Values{}
	res := &struct {
		Count int64             `json:"num"`
		Code  int               `json:"code"`
		List  []*model.BasicTag `json:"list"`
	}{}
	uri := fmt.Sprintf(_regionHotTagURI, d.hosts.HotTagHost, rid)
	if err := d.client.Get(c, uri, metadata.String(c, metadata.RemoteIP), params, res); err != nil {
		log.Error("d.client.Get(%s) error(%v)", uri+"?"+params.Encode(), err)
		return nil, 0, err
	}
	if res.Code != 0 {
		log.Error("d.client.Get(%s?%s) code:%d", uri, params.Encode(), res.Code)
		return nil, 0, ecode.Int(res.Code)
	}
	return res.List, res.Count, nil
}

// ArchiveHot ArchiveHot.
func (d *Dao) ArchiveHot(c context.Context, rid int64) (checked map[int64][]*model.BasicTag, tags []string, err error) {
	uri := d.hosts.ArchiveHotHost + _archiveHotURI
	params := url.Values{}
	params.Set("typeid", fmt.Sprintf("%d", rid))
	res := &struct {
		Code int `json:"code"`
		Data struct {
			Checked map[int64][]*model.BasicTag `json:"checked"`
			Tags    []string                    `json:"tags"`
		} `json:"data"`
	}{}
	if err = d.client.Get(c, uri, metadata.String(c, metadata.RemoteIP), params, res); err != nil {
		log.Error("d.client.Get(%s) error(%v)", uri+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("d.client.Get(%s?%s) code:%d", uri, params.Encode(), res.Code)
		return nil, nil, err
	}
	return res.Data.Checked, res.Data.Tags, nil
}

// SendMsg SendMsg.
func (d *Dao) SendMsg(c context.Context, mc, title, context string, dataType int32, mids []int64) (err error) {
	uri := d.hosts.MessageHost + _sendMsgURI
	params := url.Values{}
	res := &struct {
		Code int `json:"code"`
	}{}
	params.Set("mc", mc)
	params.Set("title", title)
	params.Set("data_type", fmt.Sprintf("%d", dataType))
	params.Set("context", context)
	params.Set("mid_list", xstr.JoinInts(mids))
	if err = d.client.Post(c, uri, metadata.String(c, metadata.RemoteIP), params, res); err != nil {
		log.Error("d.client.Post(%s) error(%v)", uri+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("d.client.Post(%s?%s) code:%d", uri, params.Encode(), res.Code)
		return ecode.Int(res.Code)
	}
	return
}

// BlockUser BlockUser.
func (d *Dao) BlockUser(c context.Context, tname, title, uname, action, note string, mid, oid int64, reasonType, isNotify, blockTimeLength, blockForever int32) (err error) {
	uri := d.hosts.AccountHost + _blockUserURI
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%d", mid))
	params.Set("blockForever", fmt.Sprintf("%d", blockForever))
	params.Set("blockTimeLength", fmt.Sprintf("%d", blockTimeLength))
	params.Set("blockRemark", note)
	params.Set("operator", uname)
	params.Set("originType", fmt.Sprintf("%d", model.RptOriginTag))
	params.Set("reasonType", fmt.Sprintf("%d", reasonType))
	params.Set("originTitle", title)
	params.Set("originUrl", fmt.Sprintf(_originURI, oid))
	params.Set("originContent", fmt.Sprintf(_originContent, action, tname))
	params.Set("isNotify", fmt.Sprintf("%d", isNotify))
	res := &struct {
		Code int `json:"code"`
	}{}
	if err = d.client.Post(c, uri, metadata.String(c, metadata.RemoteIP), params, res); err != nil {
		log.Error("d.client.Get(%s) error(%v)", uri+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("d.client.Get(%s?%s) code:%d", uri, params.Encode(), res.Code)
		return ecode.Int(res.Code)
	}
	log.Warn("Block user success")
	return
}

// AddMoral AddMoral.
func (d *Dao) AddMoral(c context.Context, username, remark, reason string, addMoral, isNotify int32, mids []int64) (err error) {
	uri := d.hosts.AccountHost + _accountMoralURI
	params := url.Values{}
	params.Set("reason_type", "3")
	params.Set("addMoral", fmt.Sprintf("%d", addMoral))
	params.Set("mid", xstr.JoinInts(mids))
	params.Set("origin", "2")
	params.Set("operater", username)
	params.Set("is_notify", fmt.Sprintf("%d", isNotify))
	params.Set("remark", remark)
	params.Set("reason", reason)
	res := &struct {
		Code int `json:"code"`
	}{}
	if err = d.client.Post(c, uri, metadata.String(c, metadata.RemoteIP), params, res); err != nil {
		log.Error("d.client.Get(%s) error(%v)", uri+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("d.client.Get(%s?%s) code:%d", uri, params.Encode(), res.Code)
		return ecode.Int(res.Code)
	}
	return
}
