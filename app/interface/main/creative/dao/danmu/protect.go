package danmu

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/creative/model/danmu"
	dmMdl "go-common/app/interface/main/dm/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_dmProtectApplyListURI      = "/x/internal/dm/up/protect/apply/list"
	_dmProtectApplyStatusURI    = "/x/internal/dm/up/protect/apply/status"
	_dmProtectApplyVideoListURI = "/x/internal/dm/up/protect/apply/video/list"
)

// ProtectApplyList fn
func (d *Dao) ProtectApplyList(c context.Context, mid, page int64, aidStr, sort, ip string) (result *danmu.ApplyList, err error) {
	result = &danmu.ApplyList{}
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("aid", aidStr)
	params.Set("page", strconv.FormatInt(page, 10))
	params.Set("sort", sort)
	var res struct {
		Code int                    `json:"code"`
		Data *danmu.ApplyListFromDM `json:"data"`
	}
	if err = d.client.Get(c, d.dmProtectApplyListURL, ip, params, &res); err != nil {
		err = ecode.CreativeDanmuErr
		log.Error("d.ProtectApplyList.Get(%s,%s,%s) err(%v)", d.dmProtectApplyListURL, ip, params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.ProtectApplyList.Get(%s,%s,%s) err(%v)|code(%d)", d.dmProtectApplyListURL, ip, params.Encode(), err, res.Code)
		return
	}
	for _, v := range res.Data.List {
		v.IDStr = strconv.FormatInt(v.ID, 10)
	}
	result = &danmu.ApplyList{
		Pager: res.Data.Pager,
		List:  res.Data.List,
	}
	return
}

// ProtectApplyVideoList fn
func (d *Dao) ProtectApplyVideoList(c context.Context, mid int64, ip string) (result []*dmMdl.Video, err error) {
	result = []*dmMdl.Video{}
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int            `json:"code"`
		Data []*dmMdl.Video `json:"data"`
	}
	if err = d.client.Get(c, d.dmProtectApplyVideoListURL, ip, params, &res); err != nil {
		err = ecode.CreativeDanmuErr
		log.Error("d.ProtectApplyVideoList.Get(%s,%s,%s) err(%v)", d.dmProtectApplyVideoListURL, ip, params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.ProtectApplyVideoList.Get(%s,%s,%s) err(%v)|code(%d)", d.dmProtectApplyVideoListURL, ip, params.Encode(), err, res.Code)
		return
	}
	result = res.Data
	return
}

// ProtectOper fn
func (d *Dao) ProtectOper(c context.Context, mid, status int64, ids, ip string) (err error) {
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("status", strconv.FormatInt(status, 10))
	params.Set("ids", ids)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.dmProtectApplyStatusURL, ip, params, &res); err != nil {
		err = ecode.CreativeDanmuErr
		log.Error("d.ProtectApply.Post(%s,%s,%s) err(%v)", d.dmProtectApplyStatusURL, ip, params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.ProtectApply.Post(%s,%s,%s) err(%v)|code(%d)", d.dmProtectApplyStatusURL, ip, params.Encode(), err, res.Code)
		return
	}
	return
}
