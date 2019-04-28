package service

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	memmdl "go-common/app/service/main/member/model"
	"go-common/app/service/main/vip/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

//CreateOrder def.
func (s *Service) CreateOrder(c context.Context, arg *model.ArgCreateOrder, ip string) (pp map[string]interface{}, err error) {
	var (
		price    float64
		payOrder *model.PayOrder
		sign     string
		dh       *model.VipUserDiscountHistory
		status   int8
	)

	if arg.OrderType != model.AutoRenew && arg.OrderType != model.General {
		err = ecode.RequestErr
		return
	}
	if dh, err = s.dao.DiscountSQL(c, arg.Mid, model.FirstDiscountBuyVip); err != nil {
		err = errors.WithStack(err)
		return
	}
	if dh != nil {
		status = dh.Status
	}

	if price, err = s.Price(c, int16(arg.Months), model.PayPlatform[model.PlatformByName[arg.Platform]], arg.OrderType, status); err != nil {
		err = errors.WithStack(err)
		return
	}
	if payOrder, err = s.AddPayOrder(c, arg.Mid, 0, arg.Months, price, model.PlatformByName[arg.Platform], arg.OrderType, arg.AppID, arg.AppSubID); err != nil {
		err = errors.WithStack(err)
		return
	}

	if pp, err = s.CreatePayPlatform(c, payOrder, int(arg.DType), arg.OrderType, ip); err != nil {
		err = errors.WithStack(err)
		return
	}
	sign = s.signPayPlatform(pp, s.dao.PaySign)
	pp["sign"] = sign
	return
}

// Price get proce by month and platform.
func (s *Service) Price(c context.Context, ms int16, platfrom int8, mt int8, firstDiscountStatus int8) (price float64, err error) {
	var (
		ok   bool
		m    *model.Month
		maps *model.PriceMapping
	)
	key := s.monthkey(ms, mt)
	if m, ok = s.months[key]; !ok {
		err = ecode.VipMonthsNotFoundErr
		return
	}
	if maps, err = s.dao.PriceMapping(c, m.ID, platfrom); err != nil {
		err = errors.WithStack(err)
		return
	}
	if maps == nil {
		err = ecode.VipMonthErr
		return
	}
	if b := s.ifDiscount(maps); b {
		price = maps.DiscountMoney
	} else {
		price = maps.Money
	}
	if firstDiscountStatus == 0 && mt == model.AutoRenew {
		price = maps.FirstDiscountMoney
	}
	if price == 0 {
		err = ecode.VipMonthPriceErr
		return
	}
	return
}

// OrderList order list.
func (s *Service) OrderList(c context.Context, mid int64, pn, ps int) (res []*model.PayOrderResp, count int64, err error) {
	var (
		orders []*model.PayOrder
	)
	if count, err = s.dao.OrderCount(c, mid, model.SUCCESS); err != nil {
		err = errors.WithStack(err)
		return
	}
	if count == 0 {
		return
	}
	if pn == 0 {
		pn = _defpn
	}
	if ps == 0 {
		ps = _defps
	}
	if orders, err = s.dao.OrderList(c, mid, model.SUCCESS, pn, ps); err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, o := range orders {
		r := &model.PayOrderResp{
			OrderNo:   o.OrderNo,
			BuyMonths: o.BuyMonths,
			Money:     o.Money,
			Status:    o.Status,
			Ctime:     o.Ctime,
		}
		res = append(res, r)
	}
	return
}

// OrderInfo get orderinfo by orderno.
func (s *Service) OrderInfo(c context.Context, orderNo string) (o *model.OrderInfo, err error) {
	if orderNo == "" {
		return
	}
	if o, err = s.dao.OrderInfo(c, orderNo); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//AddPayOrder add pay order.
func (s *Service) AddPayOrder(c context.Context, mid, bmid int64, month int16, price float64, platform int, orderType int8, appID int64, appSubID string) (pay *model.PayOrder, err error) {
	var (
		tx *xsql.Tx
		id int64
	)
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
	pay = new(model.PayOrder)
	pay.OrderNo = s.orderID()
	pay.AppID = appID
	pay.Ver = 1
	pay.Mid = mid
	pay.Money = price
	pay.Status = model.PAYING
	pay.Platform = int8(platform)
	pay.BuyMonths = month
	pay.OrderType = orderType
	pay.AppSubID = appSubID
	pay.ToMid = bmid
	pay.RechargeBp = pay.Money
	if id, err = s.dao.TxAddOrder(tx, pay); err != nil {
		err = errors.WithStack(err)
		return
	}
	pay.ID = id
	olog := new(model.VipPayOrderLog)
	olog.Status = pay.Status
	olog.OrderNo = pay.OrderNo
	olog.Mid = mid
	if err = s.dao.TxAddOrderLog(tx, olog); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// CreatePayPlatform .
func (s *Service) CreatePayPlatform(c context.Context, pay *model.PayOrder, dtype int, orderType int8, ip string) (payPlatform map[string]interface{}, err error) {
	var (
		member *memmdl.BaseInfo
	)
	payPlatform = make(map[string]interface{})
	payPlatform["customerId"] = s.c.PayConf.CustomerID
	payPlatform["serviceType"] = 0
	payPlatform["originalAmount"] = s.yuanToFen(pay.Money)
	payPlatform["payAmount"] = s.yuanToFen(pay.Money)
	payPlatform["deviceType"] = int8(dtype)
	payPlatform["productId"] = s.c.PayConf.ProductID
	payPlatform["notifyUrl"] = s.c.PayConf.OrderNotifyURL
	payPlatform["orderCreateTime"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	payPlatform["timestamp"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	payPlatform["version"] = s.c.PayConf.Version
	payPlatform["orderId"] = pay.OrderNo
	payPlatform["orderExpire"] = _defOrderExpire
	payPlatform["traceId"] = uuid.NewV1().String()
	payPlatform["createIp"] = ip
	payPlatform["feeType"] = "CNY"
	if pay.Platform == model.DevicePC {
		payPlatform["returnUrl"] = s.c.PayConf.ReturnURL + "&order_no=" + pay.OrderNo
	}
	if orderType == model.AutoRenew {
		payPlatform["showTitle"] = "购买大会员连续包月"
		payPlatform["signUrl"] = s.c.PayConf.SignNotifyURL
		if member, err = s.memInfoRetry(c, pay.Mid); err != nil {
			err = errors.WithStack(err)
			return
		}

		payPlatform["planId"] = s.c.PayConf.PlanID
		payPlatform["uid"] = pay.Mid
		payPlatform["serviceType"] = _autoPayserviceType
		payPlatform["displayAccount"] = member.Name
	} else {
		payPlatform["showTitle"] = "购买大会员"
	}
	payPlatform["signType"] = "MD5"
	return
}

func (s *Service) signPayPlatform(pay map[string]interface{}, fn func(params map[string]string, token string) string) (query string) {
	params := make(map[string]string)
	for k, v := range pay {
		params[k] = fmt.Sprintf("%v", v)
	}
	query = fn(params, s.c.PayConf.Token)
	return
}

func (s *Service) yuanToFen(money float64) int64 {
	return int64(s.floatRound(money*100, _defround))
}

func (s *Service) fenToYuan(money float64) float64 {
	return float64(s.floatRound(money/100, _defround))
}

// orderID get order id
func (s *Service) orderID() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%05d", s.r.Int63n(99999)))
	b.WriteString(fmt.Sprintf("%03d", time.Now().UnixNano()/1e6%1000))
	b.WriteString(time.Now().Format("060102150405"))
	return b.String()
}

// IfDiscount Is it during a discount.
func (s *Service) ifDiscount(m *model.PriceMapping) (b bool) {
	var (
		now   = time.Now()
		start bool
		end   bool
	)
	if m == nil || m.DiscountMoney == 0 {
		return
	}
	if m.StartTime.Time().IsZero() && m.EndTime.Time().IsZero() {
		b = true
		return
	}
	if m.StartTime.Time().IsZero() || now.After(m.StartTime.Time()) {
		start = true
	}
	if m.EndTime.Time().IsZero() || now.Before(m.EndTime.Time()) {
		end = true
	}
	if start && end {
		b = true
	}
	return
}

func (s *Service) memInfoRetry(c context.Context, mid int64) (member *memmdl.BaseInfo, err error) {
	for i := 1; i <= _retryTimes; i++ {
		member, err = s.memRPC.Base(c, &memmdl.ArgMemberMid{Mid: mid})
		if member == nil || err != nil {
			err = errors.WithStack(err)
			continue
		}
	}
	return
}

//CreateOldOrder create old order.
func (s *Service) CreateOldOrder(c context.Context, arg *model.ArgOldPayOrder) (err error) {
	var (
		tx *xsql.Tx
		id int64
	)
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
	o := new(model.PayOrder)
	o.OrderNo = arg.OrderNo
	o.AppID = arg.AppID
	o.Platform = arg.Platform
	o.OrderType = arg.OrderType
	o.AppSubID = arg.AppSubID
	o.Mid = arg.Mid
	o.ToMid = arg.ToMid
	o.BuyMonths = arg.BuyMonths
	o.Money = arg.Money
	o.Status = arg.Status
	o.PayType = arg.PayType
	o.RechargeBp = arg.RechargeBp
	o.ThirdTradeNo = arg.ThirdTradeNo
	o.UserIP = []byte{}
	if id, err = s.dao.TxAddOrder(tx, o); err != nil {
		err = errors.WithStack(err)
		return
	}
	o.ID = id
	olog := new(model.VipPayOrderLog)
	olog.Mid = arg.Mid
	olog.OrderNo = o.OrderNo
	olog.Status = o.Status
	if err = s.dao.TxAddOrderLog(tx, olog); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//OrderMng .
func (s *Service) OrderMng(c context.Context, mid int64) (order *model.OrderMng, err error) {
	var (
		vip      *model.VipInfo
		payOrder *model.PayOrder
		member   *memmdl.BaseInfo
		now      = time.Now()
	)
	order = new(model.OrderMng)
	if vip, err = s.VipInfo(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	s.logicAutoRenew(vip)
	if vip.VipPayType != 1 {
		return
	}
	if vip.VipOverdueTime.Time().Before(now) {
		return
	}
	payChannel := s.c.Property.PayChannelMapping[strconv.Itoa(int(vip.PayChannelID))]
	if payOrder, err = s.dao.PayOrderLast(c, mid, model.SUCCESS, int64(model.OtherRenew), int64(model.IOSRenew)); err != nil {
		err = errors.WithStack(err)
		return
	}
	if payOrder == nil {
		return
	}

	if vip.PayChannelID == model.IapPayChannelID {
		order.NextDedutionDate = vip.IosOverdueTime.Time().Format("2006-01-02")
	} else if vip.VipOverdueTime.Time().After(now) {
		order.NextDedutionDate = vip.VipOverdueTime.Time().AddDate(0, 0, -2).Format("2006-01-02")
	}

	if member, err = s.memInfoRetry(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}

	order.ExpireDate = vip.VipOverdueTime.Time().Format("2006-01-02")
	order.IsAutoRenew = int8(vip.VipPayType)
	if vip.PayChannelID == model.IapPayChannelID {
		payChannel = "苹果应用内购买"
	}
	order.PayType = payChannel
	order.Mid = vip.Mid
	order.ChannelID = int32(vip.PayChannelID)
	order.AutoRenewLoop = strconv.Itoa(int(payOrder.BuyMonths)) + "个月"
	order.PriceTip = fmt.Sprintf("%.2f元", payOrder.Money)
	order.Username = member.Name
	return
}

//Rescision pay rescision.
func (s *Service) Rescision(c context.Context, mid int64, deviceType int32) (err error) {
	var (
		vip    *model.VipInfo
		params = make(map[string]interface{})
	)
	if vip, err = s.VipInfo(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if vip.PayChannelID == int32(_iosPaychannel) {
		err = ecode.Int(1001)
		return
	}

	if vip.VipPayType != model.AuoRenewVip {
		err = ecode.VipRenewTypeErr
		return
	}

	payChannel := s.c.Property.PayChannelMapping[strconv.Itoa(int(vip.PayChannelID))]

	if len(payChannel) == 0 {
		err = ecode.VipPayChannelNotExitErr
		return
	}

	params["customerId"] = s.c.PayConf.CustomerID
	params["traceId"] = uuid.NewV1().String()
	params["timestamp"] = time.Now().UnixNano() / 1e6
	params["signType"] = "MD5"
	params["version"] = "1.0"
	params["planId"] = s.c.PayConf.PlanID
	params["remark"] = "我要解约"
	params["uid"] = mid
	params["payChannel"] = payChannel
	params["deviceType"] = deviceType
	params["serviceType"] = _autoPayserviceType
	params["payChannelId"] = vip.PayChannelID

	if err = s.dao.PayRecission(c, params, metadata.String(c, metadata.RemoteIP)); err != nil {
		err = errors.WithStack(err)
	}
	return
}
