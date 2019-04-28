package danmu

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/creative/model/danmu"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	// api
	_getDmPurchaseListURI   = "/x/internal/dm/adv/list"
	_setDmPurchasePassURI   = "/x/internal/dm/adv/pass"
	_setDmPurchaseDenyURI   = "/x/internal/dm/adv/deny"
	_setDmPurchaseCancelURI = "/x/internal/dm/adv/cancel"
)

// GetAdvDmPurchases fn
func (d *Dao) GetAdvDmPurchases(c context.Context, mid int64, ip string) (danmus []*danmu.AdvanceDanmu, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int                   `json:"code"`
		Data []*danmu.AdvanceDanmu `json:"data"`
	}
	if err = d.client.Get(c, d.advDmPurchaseListURL, ip, params, &res); err != nil {
		log.Error("d.ListAdvanceDm.Get(%s,%s,%s) err(%v)", d.advDmPurchaseListURL, ip, params.Encode(), err)
		err = ecode.CreativeDanmuErr
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.ListAdvanceDm.Get(%s,%s,%s) err(%v)|code（%d）", d.advDmPurchaseListURL, ip, params.Encode(), err, res.Code)
		return
	}
	danmus = res.Data
	return
}

// PassAdvDmPurchase fn
func (d *Dao) PassAdvDmPurchase(c context.Context, mid, id int64, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("id", strconv.FormatInt(id, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.advDmPurchasePassURL, ip, params, &res); err != nil {
		log.Error("d.advDmPurchasePass.Post(%s,%s,%s) err(%v)", d.advDmPurchasePassURL, ip, params.Encode(), err)
		err = ecode.CreativeDanmuErr
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.advDmPurchasePass.Post(%s,%s,%s) err(%v)|code(%d)", d.advDmPurchasePassURL, ip, params.Encode(), err, res.Code)
		return
	}
	return
}

// DenyAdvDmPurchase fn
func (d *Dao) DenyAdvDmPurchase(c context.Context, mid, id int64, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("id", strconv.FormatInt(id, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.advDmPurchaseDenyURL, ip, params, &res); err != nil {
		log.Error("d.advDmPurchaseDeny.Post(%s,%s,%s) err(%v)", d.advDmPurchaseDenyURL, ip, params.Encode(), err)
		err = ecode.CreativeDanmuErr
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.advDmPurchaseDeny.Post(%s,%s,%s) err(%v)|code(%d)", d.advDmPurchaseDenyURL, ip, params.Encode(), err, res.Code)
		return
	}
	return
}

// CancelAdvDmPurchase fn
func (d *Dao) CancelAdvDmPurchase(c context.Context, mid, id int64, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("id", strconv.FormatInt(id, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.advDmPurchaseCancelURL, ip, params, &res); err != nil {
		log.Error("d.advDmPurchaseCancel.Post(%s,%s,%s) err(%v)", d.advDmPurchaseCancelURL, ip, params.Encode(), err)
		err = ecode.CreativeDanmuErr
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.advDmPurchaseCancel.Post(%s,%s,%s) err(%v)|code(%d)", d.advDmPurchaseCancelURL, ip, params.Encode(), err, res.Code)
		return
	}
	return
}
