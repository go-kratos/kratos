package dao

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"go-common/app/interface/main/space/model"
	"go-common/library/ecode"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_groupsCountURI = "/link_group/v1/member/created_groups_num"
	_dynamicListURI = "/dynamic_svr/v0/dynamic_svr/co_space_history"
	_dynamicCntURI  = "/dynamic_svr/v0/dynamic_svr/space_dy_num"
	_dynamicURI     = "/dynamic_svr/v1/dynamic_svr/get_dynamic_detail"
)

// GroupsCount .
func (d *Dao) GroupsCount(c context.Context, mid, vmid int64) (count int, err error) {
	var (
		req *http.Request
		ip  = metadata.String(c, metadata.RemoteIP)
	)
	params := url.Values{}
	params.Set("master_uid", strconv.FormatInt(vmid, 10))
	if req, err = d.httpR.NewRequest(http.MethodGet, d.groupsCountURL, ip, params); err != nil {
		return
	}
	req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Num int `json:"num"`
		}
	}
	if err = d.httpR.Do(c, req, &res); err != nil {
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

// DynamicCnt dynamic count.
func (d *Dao) DynamicCnt(c context.Context, mid int64) (cnt int64, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("uids", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Items []*struct {
				UID int64 `json:"uid"`
				Num int64 `json:"num"`
			} `json:"items"`
		}
	}
	if err = d.httpR.Get(c, d.dynamicCntURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.dynamicCntURL+"?"+params.Encode())
		return
	}
	if len(res.Data.Items) > 0 && res.Data.Items[0].UID == mid {
		cnt = res.Data.Items[0].Num
	}
	return
}

// DynamicList .
func (d *Dao) DynamicList(c context.Context, mid, vmid, dyID int64, qn, page int) (data *model.DyList, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	if mid > 0 {
		params.Set("visitor_uid", strconv.FormatInt(mid, 10))
	}
	params.Set("host_uid", strconv.FormatInt(vmid, 10))
	params.Set("offset_dynamic_id", strconv.FormatInt(dyID, 10))
	params.Set("qn", strconv.Itoa(qn))
	params.Set("page", strconv.Itoa(page))
	var res struct {
		Code int           `json:"code"`
		Data *model.DyList `json:"data"`
	}
	if err = d.httpR.Get(c, d.dynamicListURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.dynamicListURL+"?"+params.Encode())
		return
	}
	data = res.Data
	return
}

// Dynamic .
func (d *Dao) Dynamic(c context.Context, mid, dynamicID int64, qn int) (data *model.DyCard, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	if mid > 0 {
		params.Set("uid", strconv.FormatInt(mid, 10))
	}
	params.Set("dynamic_id", strconv.FormatInt(dynamicID, 10))
	params.Set("qn", strconv.Itoa(qn))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Card *model.DyCard `json:"card"`
		} `json:"data"`
	}
	if err = d.httpR.Get(c, d.dynamicURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.dynamicURL+"?"+params.Encode())
		return
	}
	data = res.Data.Card
	return
}
