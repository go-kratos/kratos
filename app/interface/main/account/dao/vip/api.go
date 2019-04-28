package vip

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/interface/main/account/model"
	vipmol "go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_vipCodeVerify = "/x/internal/vip/code/verify"
	_vipCodeOpen   = "/x/internal/vip/code/open"
	_viptips       = "/x/internal/vip/tips"
	_couponCancel  = "/x/internal/vip/coupon/cancel"
	_vipCodeOpened = "/x/internal/vip/code/opened"

	// vip with java
	_vipInfo = "/internal/v1/user/"
)

//CodeVerify code verify
func (d *Dao) CodeVerify(c context.Context) (token *model.Token, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	var tokenResq *model.TokenResq
	val := url.Values{}
	if err = d.client.Get(c, d.codeVerifyURL, ip, val, &tokenResq); err != nil {
		err = errors.WithStack(err)
		return
	}
	if tokenResq.Code == int64(ecode.OK.Code()) {
		token = tokenResq.Data
	}
	return
}

//CodeOpen http code open
func (d *Dao) CodeOpen(c context.Context, mid int64, code, token, verify string) (data *model.ResourceCodeResq, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	val := url.Values{}
	val.Add("mid", fmt.Sprintf("%v", mid))
	val.Add("token", token)
	val.Add("verify", verify)
	val.Add("code", code)
	defer func() {
		log.Info("qingqiu url(%+v) params(%+v) return(%+v)", d.codeOpenURL, val, data)
	}()
	if err = d.client.Post(c, d.codeOpenURL, ip, val, &data); err != nil {
		err = errors.WithStack(err)
		return
	}

	if data.Code > int64(ecode.OK.Code()) {
		err = ecode.Int(int(data.Code))
		return
	}
	return
}

// Info get vip info
func (d *Dao) Info(c context.Context, mid int64, ip string) (info *model.VIPInfo, err error) {
	var res struct {
		Code int            `json:"code"`
		Data *model.VIPInfo `json:"data"`
	}
	if err = d.client.Get(c, d.infoURL+strconv.FormatInt(mid, 10), ip, nil, &res); err != nil {
		log.Error("d.client.Get(%d) error(%v)", mid, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		log.Error("d.client.Get(%d) error(%v)", mid, err)
		return
	}
	info = res.Data
	return
}

// Tips tips info.
func (d *Dao) Tips(c context.Context, version int64, platform string) (data *vipmol.TipsResp, err error) {
	params := url.Values{}
	params.Add("version", fmt.Sprintf("%v", version))
	params.Add("platform", platform)
	var resp struct {
		Code int              `json:"code"`
		Data *vipmol.TipsResp `json:"data"`
	}
	if err = d.client.Get(c, d.tipsURL, "", params, &resp); err != nil {
		err = errors.Errorf("vip tips d.httpClient.Do(%s) error(%+v)", d.tipsURL+"?"+params.Encode(), err)
		return
	}
	if resp.Code != ecode.OK.Code() {
		err = errors.Errorf("vip tips url(%s) res(%+v) err(%+v)", d.tipsURL+"?"+params.Encode(), resp, ecode.Int(resp.Code))
	}
	data = resp.Data
	return
}

// CancelUseCoupon cancel use coupon.
func (d *Dao) CancelUseCoupon(c context.Context, arg *vipmol.ArgCancelUseCoupon) (err error) {
	params := url.Values{}
	params.Add("mid", fmt.Sprintf("%d", arg.Mid))
	params.Add("coupon_token", arg.CouponToken)
	var resp struct {
		Code int `json:"code"`
	}
	if err = d.clientSlow.Post(c, d.cancelCouponURL, "", params, &resp); err != nil {
		err = errors.Errorf("vip cancel coupon d.httpClient.Do(%s) error(%+v)", d.cancelCouponURL+"?"+params.Encode(), err)
		return
	}
	if resp.Code != ecode.OK.Code() {
		err = ecode.Int(resp.Code)
	}
	return
}

//CodeOpeneds sel code opened data
func (d *Dao) CodeOpeneds(c context.Context, arg *model.CodeInfoReq, ip string) (resp []*vipmol.CodeInfoResp, err error) {
	val := url.Values{}
	val.Add("bis_appkey", arg.Appkey)
	val.Add("bis_sign", arg.Sign)
	val.Add("bis_ts", strconv.FormatInt(arg.Ts.Time().Unix(), 10))
	val.Add("start_time", strconv.FormatInt(arg.StartTime.Time().Unix(), 10))
	val.Add("end_time", strconv.FormatInt(arg.EndTime.Time().Unix(), 10))
	val.Add("cursor", strconv.FormatInt(arg.Cursor, 10))
	rep := new(struct {
		Code int                    `json:"code"`
		Data []*vipmol.CodeInfoResp `json:"data"`
	})
	defer func() {
		log.Info("vip code opened url:%+v params:%+v return:%+v", d.codeOpenedURL, val, rep)
	}()
	if err = d.client.Get(c, d.codeOpenedURL, ip, val, rep); err != nil {
		err = errors.Errorf("vip code opened url:%+v params:%+v return:%+v,err:%+v", d.codeOpenedURL, val, rep, err)
		return
	}
	if rep.Code != ecode.OK.Code() {
		err = ecode.Int(rep.Code)
		return
	}
	resp = rep.Data
	return
}
