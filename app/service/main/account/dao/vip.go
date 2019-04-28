package dao

import (
	"context"
	"net/url"

	v1 "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

var (
	_emptyVipInfos3 = map[int64]*v1.VipInfo{}
)

// RawVip get mid's vip info from account center by vip API.
func (d *Dao) RawVip(c context.Context, mid int64) (vip *v1.VipInfo, err error) {
	params := url.Values{}
	var res struct {
		Code int `json:"code"`
		Data struct {
			Type       int32 `json:"vipType"`
			Status     int32 `json:"vipStatus"`
			DueDate    int64 `json:"vipDueDate"`
			VipPayType int32 `json:"isAutoRenew"`
		} `json:"data"`
	}
	err = d.httpR.RESTfulGet(c, d.vipInfoURI, metadata.String(c, metadata.RemoteIP), params, &res, mid)
	if err != nil {
		err = errors.Wrap(err, "dao vip")
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		err = errors.Wrap(err, "dao vip")
		return
	}
	vip = &v1.VipInfo{
		Type:       res.Data.Type,
		Status:     res.Data.Status,
		DueDate:    res.Data.DueDate,
		VipPayType: res.Data.VipPayType,
	}
	return
}

// RawVips get multi mid's vip info from account center by vip API.
func (d *Dao) RawVips(c context.Context, mids []int64) (res map[int64]*v1.VipInfo, err error) {
	params := url.Values{}
	params.Set("idList", xstr.JoinInts(mids))
	var info struct {
		Code int `json:"code"`
		Data map[int64]*struct {
			Type    int32 `json:"vipType"`
			Status  int32 `json:"vipStatus"`
			DueDate int64 `json:"vipDueDate"`
		} `json:"data"`
	}
	err = d.httpR.Get(c, d.vipMultiInfoURI, "", params, &info)
	if err != nil {
		err = errors.Wrap(err, "dao vip")
		return
	}
	if info.Code != 0 {
		err = ecode.Int(info.Code)
		err = errors.Wrap(err, "dao vip")
		return
	}
	if len(info.Data) == 0 {
		res = _emptyVipInfos3
		return
	}
	res = make(map[int64]*v1.VipInfo, len(info.Data))
	for mid, v := range info.Data {
		if v == nil {
			continue
		}
		vip := &v1.VipInfo{
			Type:    v.Type,
			Status:  v.Status,
			DueDate: v.DueDate,
		}
		res[mid] = vip
	}
	return
}
