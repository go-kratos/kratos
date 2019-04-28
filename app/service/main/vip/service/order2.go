package service

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	col "go-common/app/service/main/coupon/model"
	memmdl "go-common/app/service/main/member/model"
	"go-common/app/service/main/vip/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

//CreateOrder2 create order v2.
func (s *Service) CreateOrder2(c context.Context, a *model.ArgCreateOrder2) (r *model.CreateOrderRet, o *model.PayOrder, err error) {
	var (
		p                           *model.VipPirce
		tx                          *xsql.Tx
		plat                        = orderPlat(a.Platform, a.Device, a.MobiApp, a.Build)
		pargs                       map[string]interface{}
		dprice, oprice, couponMoney float64
		id                          int64
	)
	r = new(model.CreateOrderRet)
	if a.Bmid > 0 {
		// give friend can not use coupon.
		a.CouponToken = ""
		a.PanelType = model.PanelTypeFriend
	}
	if a.CouponToken != "" {
		//FIXME 代金券限制平台上线后可以删除
		if a.PanelType == "ele" {
			log.Warn("illegal create order arg:%+v", a)
			err = ecode.CouPonPlatformNotSupportErr
			return
		}
		//FIXME 代金券限制平台上线后可以删除 end
		if err = s.CancelUseCoupon(c, a.Mid, a.CouponToken, IPStr(a.IP)); err != nil {
			if err == ecode.CouPonStateCanNotCancelErr {
				err = nil
			} else {
				return
			}
		}
	}
	if p, err = s.VipPriceV2(c, &model.ArgPriceV2{
		Mid:       a.Mid,
		Month:     int16(a.Month),
		SubType:   a.OrderType,
		Token:     a.CouponToken,
		Platform:  a.Platform,
		PanelType: a.PanelType,
		MobiApp:   a.MobiApp,
		Device:    a.Device,
		Build:     a.Build,
	}); err != nil {
		err = errors.WithStack(err)
		return
	}
	if p == nil || p.Panel == nil {
		err = ecode.VipOrderPirceErr
		return
	}
	// 商品限制
	if err = s.ProductLimit(c, &model.ArgProductLimit{
		Mid:       a.Mid,
		PanelType: a.PanelType,
		Months:    a.Month,
	}); err != nil {
		return
	}
	dprice = p.Panel.DPrice
	oprice = p.Panel.OPrice
	if p.Coupon != nil && p.Coupon.Amount >= 0 {
		couponMoney = p.Coupon.Amount
		dprice = s.floatRound(dprice-couponMoney, _defround)
	}
	if dprice <= 0 {
		err = ecode.VipOrderPirceErr
		return
	}
	if tx, err = s.dao.StartTx(c); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		if err == nil {
			if err = tx.Commit(); err != nil {
				tx.Rollback()
			}
		} else {
			tx.Rollback()
		}
	}()
	o = s.convertOrder(a, plat, dprice, couponMoney, a.Bmid)
	o.PID = p.Panel.Id
	if id, err = s.dao.TxAddOrder(tx, o); err != nil {
		err = errors.WithStack(err)
		return
	}
	if id > 1 {
		if err = s.dao.TxAddOrderLog(tx, &model.VipPayOrderLog{Mid: a.Mid, OrderNo: o.OrderNo, Status: o.Status}); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	if pargs, err = s.createPayParams(c, o, p.Panel, a, plat); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a.CouponToken != "" {
		if err = s.couRPC.UseAllowance(c, &col.ArgUseAllowance{
			Mid:            o.Mid,
			CouponToken:    a.CouponToken,
			Remark:         model.CouponUseRemark,
			OrderNO:        o.OrderNo,
			Price:          p.Panel.DPrice,
			Platform:       a.Platform,
			PanelType:      a.PanelType,
			MobiApp:        a.MobiApp,
			Device:         a.Device,
			Build:          a.Build,
			ProdLimMonth:   int8(a.Month),
			ProdLimRenewal: model.MapProdLlimRenewal[a.OrderType],
		}); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	r.Dprice = dprice
	r.Oprice = oprice
	r.PayParam = pargs
	r.CouponMoney = couponMoney
	r.UserIP = IPStr(a.IP)
	r.PID = o.PID
	return
}

// CreateQrCodeOrder create qrcode order.
func (s *Service) CreateQrCodeOrder(c context.Context, a *model.ArgCreateOrder2) (res *model.PayQrCodeRet, err error) {
	var (
		qr      *model.PayQrCode
		r       *model.CreateOrderRet
		orderNo string
	)
	defer func() {
		if err != nil && a.CouponToken != "" {
			s.CancelUseCoupon(c, a.Mid, a.CouponToken, IPStr(a.IP))
		}
	}()
	res = new(model.PayQrCodeRet)
	if r, _, err = s.CreateOrder2(c, a); err != nil {
		err = errors.WithStack(err)
		return
	}
	orderNo = r.PayParam["orderId"].(string)
	if qr, err = s.dao.PayQrCode(c, a.Mid, orderNo, r.PayParam); err != nil {
		err = errors.WithStack(err)
		return
	}
	data := &model.PayQrCodeResp{
		CodeURL:     qr.CodeURL,
		ExpiredTime: qr.ExpiredTime,
		Amount:      r.Dprice,
		SaveAmount:  s.floatRound(r.Oprice-r.Dprice, _defround),
		OrderNo:     orderNo,
	}
	if a.OrderType == model.AutoRenewPayType {
		data.Tip = model.QrAutoRenewTip
	} else {
		data.Tip = model.QrTip
	}
	res.PayQrCodeResp = data
	res.Dprice = r.Dprice
	res.CouponMoney = r.CouponMoney
	res.UserIP = r.UserIP
	res.PID = r.PID
	return
}

func (s *Service) convertOrder(a *model.ArgCreateOrder2, plat int, price float64, couponMoney float64, bmid int64) (o *model.PayOrder) {
	o = &model.PayOrder{
		OrderNo:      s.orderID(),
		AppID:        a.AppID,
		Platform:     int8(plat),
		OrderType:    a.OrderType,
		AppSubID:     a.AppSubID,
		Mid:          a.Mid,
		BuyMonths:    int16(a.Month),
		Money:        price,
		Status:       model.PAYING,
		RechargeBp:   price,
		ThirdTradeNo: "",
		CouponMoney:  couponMoney,
		ToMid:        bmid,
		UserIP:       a.IP,
	}
	if o.UserIP == nil {
		o.UserIP = []byte{}
	}
	if a.OrderType == model.AutoRenewPayType && (plat == model.PlatfromIOS ||
		plat == model.PlatfromIPAD || plat == model.PlatfromIPADHD || plat == model.PlatfromIOSBLUE) {
		o.OrderType = model.IapAutoRenewPayType
	}
	return
}

func orderPlat(platform, device, mobiApp string, build int64) (p int) {
	p = model.PlatformByName[platform]
	if model.DeviceIapdName == device && model.MobiAppIpadName == mobiApp {
		p = model.PlatfromIPADHD
	}
	if mobiApp == "android_i" {
		p = model.PlatfromANDROIDI
	}
	if mobiApp == "iphone_b" || (mobiApp == "iphone" && build > 7000 && build < 8000) {
		p = model.PlatfromIOSBLUE
	}
	if p == 0 {
		p = model.PlatfromPC
	}
	return
}

func (s *Service) floatRound(f float64, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}

func (s *Service) createPayParams(c context.Context, o *model.PayOrder, p *model.VipPanelInfo, a *model.ArgCreateOrder2, plat int) (r map[string]interface{}, err error) {
	var (
		m   *memmdl.BaseInfo
		now = time.Now().Unix() * _millis
		u   *url.URL
	)
	r = make(map[string]interface{})
	r["customerId"] = s.c.PayConf.CustomerID
	r["serviceType"] = model.ServiceTypeNormal
	r["originalAmount"] = s.yuanToFen(p.OPrice)
	r["createIp"] = a.IP
	r["deviceType"] = a.Dtype
	if p.PdID != "" {
		r["productId"] = p.PdID
	} else {
		r["productId"] = s.c.PayConf.ProductID
	}
	r["notifyUrl"] = s.c.PayConf.OrderNotifyURL
	r["orderCreateTime"] = now
	r["timestamp"] = now
	r["version"] = s.c.PayConf.Version
	r["payAmount"] = s.yuanToFen(o.Money)
	r["orderId"] = o.OrderNo
	r["orderExpire"] = s.c.PayConf.OrderExpire
	r["feeType"] = "CNY"
	r["traceId"] = model.UUID4()
	r["showTitle"] = model.NormalShowTitle
	r["showContent"] = fmt.Sprintf(model.ShowContent, o.BuyMonths)
	if plat == model.PlatfromPC || plat == model.PlatfromPUBLIC {
		if a.ReturnURL != "" {
			if u, err = url.Parse(a.ReturnURL); err != nil || !u.IsAbs() {
				a.ReturnURL = ""
			}
		}
		if a.ReturnURL != "" {
			r["returnUrl"] = fmt.Sprintf("%s&type=payCallBack&order_no=%s", a.ReturnURL, o.OrderNo)
		} else {
			r["returnUrl"] = fmt.Sprintf("%s&order_no=%s", s.c.PayConf.ReturnURL, o.OrderNo)
		}
		r["defaultChoose"] = "bp"
		r["uid"] = o.Mid
		if plat == model.PlatfromPUBLIC {
			r["serviceType"] = model.ServiceTypePublic
		}
	}
	switch {
	case plat == model.PlatfromIPAD || plat == model.PlatfromIOS || plat == model.PlatfromIPADHD || plat == model.PlatfromIOSBLUE:
		r["serviceType"] = model.ServiceTypeIap
		if o.OrderType == model.IapAutoRenewPayType {
			r["showTitle"] = model.AutoRenewShowTitle
			r["subscribeType"] = model.PaySubTypeAuto
		}
	case o.OrderType == int8(model.AutoRenewPay):
		r["showTitle"] = model.AutoRenewShowTitle
		r["signUrl"] = s.c.PayConf.SignNotifyURL
		r["planId"] = s.c.PayConf.PlanID
		r["serviceType"] = model.ServiceTypeAuto
		if m, err = s.retryGetMemberInfo(c, o.Mid); err == nil && m != nil {
			r["displayAccount"] = m.Name
		} else {
			log.Error("s.memRPC.Base(%d) err(%+v)", o.Mid, err)
			err = nil
			r["displayAccount"] = o.Mid
		}
		r["uid"] = o.Mid
	}
	if plat == model.PlatfromANDROIDI {
		r["serviceType"] = model.ServiceTypeInternational
	}
	r["signType"] = s.c.PayConf.SignType
	// sign
	r["sign"] = s.signPayPlatform(r, s.dao.PaySign)
	return
}
