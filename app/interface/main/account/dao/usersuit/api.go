package usersuit

import (
	"context"
	"net/url"
	"strconv"

	usmdl "go-common/app/service/main/usersuit/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_groupInfo    = "/x/internal/pendant/groupInfo"
	_entryGroup   = "/x/internal/pendant/entryGroup"
	_vipGroup     = "/x/internal/pendant/vipGroup"
	_pendantInfo  = "/x/internal/pendant/pendantByID"
	_packageInfo  = "/x/internal/pendant/package"
	_orderInfo    = "/x/internal/pendant/order"
	_orderHistory = "/x/internal/pendant/orderHistory"
)

// Group get pendant group info.
func (d *Dao) Group(c context.Context, ip string) (groups []*usmdl.PendantGroupInfo, err error) {
	var res struct {
		Code int                       `json:"code"`
		Data []*usmdl.PendantGroupInfo `json:"data"`
	}
	if err = d.http.Get(c, d.groupURL, ip, nil, &res); err != nil {
		log.Error("d.http.Get(%s) error(%v)", d.groupURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		log.Error("d.http.Get(%s) error(%v)", d.groupURL, err)
		return
	}
	groups = res.Data
	return
}

// GroupEntry get entry group.
func (d *Dao) GroupEntry(c context.Context, ip string) (group *usmdl.PendantGroupInfo, err error) {
	var res struct {
		Code int                     `json:"code"`
		Data *usmdl.PendantGroupInfo `json:"data"`
	}
	if err = d.http.Get(c, d.entryGroupURL, ip, nil, &res); err != nil {
		log.Error("d.http.Get(%s) error(%v)", d.entryGroupURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		log.Error("d.http.Get(%s) error(%v)", d.entryGroupURL, err)
		return
	}
	group = res.Data
	return
}

// GroupVip get vip group.
func (d *Dao) GroupVip(c context.Context, ip string) (group *usmdl.PendantGroupInfo, err error) {
	var res struct {
		Code int                     `json:"code"`
		Data *usmdl.PendantGroupInfo `json:"data"`
	}
	if err = d.http.Get(c, d.vipGroupURL, ip, nil, &res); err != nil {
		log.Error("d.http.Get(%s) error(%v)", d.vipGroupURL, err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.http.Get(%s) error(%v)", d.vipGroupURL, err)
		return
	}
	group = res.Data
	return
}

// Pendant get pendant info.
func (d *Dao) Pendant(c context.Context, pid int64, ip string) (pendant *usmdl.Pendant, err error) {
	var res struct {
		Code int            `json:"code"`
		Data *usmdl.Pendant `json:"data"`
	}
	params := url.Values{}
	params.Set("pid", strconv.FormatInt(pid, 10))
	if err = d.http.Get(c, d.pendantURL, ip, params, &res); err != nil {
		log.Error("d.http.Get(%s) error(%v)", d.pendantURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		log.Error("d.http.Get(%s) error(%v)", d.pendantURL, err)
		return
	}
	pendant = res.Data
	return
}

// Packages get package info.
func (d *Dao) Packages(c context.Context, mid int64, ip string) (pkg []*usmdl.PendantPackage, err error) {
	var res struct {
		Code int                     `json:"code"`
		Data []*usmdl.PendantPackage `json:"data"`
	}
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	if err = d.http.Get(c, d.packageURL, ip, params, &res); err != nil {
		log.Error("d.http.Get(%s) error(%v)", d.packageURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		log.Error("d.http.Get(%s) error(%v)", d.packageURL, err)
		return
	}
	pkg = res.Data
	return
}

// Order order pendant.
func (d *Dao) Order(c context.Context, mid, pid, expires int64, moneyType int8, ip string) (payInfo *usmdl.PayInfo, err error) {
	params := url.Values{}
	params.Set("moneytype", strconv.Itoa(int(expires)))
	params.Set("expires", strconv.FormatInt(expires, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("pid", strconv.FormatInt(pid, 10))
	var res struct {
		Code int            `json:"code"`
		Data *usmdl.PayInfo `json:"data"`
	}
	if err = d.http.Post(c, d.orderURL, ip, params, &res); err != nil {
		log.Error("d.http.Post(%d) url(%s) error(%v)", mid, d.orderURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		log.Error("d.http.Post(%d) url(%s) error(%v)", mid, d.orderURL, err)
		return
	}
	payInfo = res.Data
	return
}

// OrderHistory get order history by mid.
func (d *Dao) OrderHistory(c context.Context, mid, page int64, payType int8, orderID, ip string) (hs []*usmdl.PendantOrderInfo, count map[string]int64, err error) {
	params := url.Values{}
	params.Set("orderID", orderID)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("page", strconv.FormatInt(page, 10))
	params.Set("payType", strconv.Itoa(int(payType)))
	var res struct {
		Code   int                       `json:"code"`
		Orders []*usmdl.PendantOrderInfo `json:"data"`
		Count  map[string]int64          `json:"count"`
	}
	if err = d.http.Get(c, d.orderHistory, ip, params, &res); err != nil {
		log.Error("d.http.Get(%d) url(%s) error(%v)", mid, d.orderHistory, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		log.Error("d.http.Get(%s) error(%v)", d.orderHistory, err)
		return
	}
	hs = res.Orders
	count = res.Count
	return
}
