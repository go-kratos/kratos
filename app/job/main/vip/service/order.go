package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/vip/model"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	_autoRenewFailTwoMsg  = "连续包月大会员今天凌晨续费又失败了，服务已暂停。如果想要再次享受连续包月大会员服务，请先取消连续包月，再去开通哦~"
	_autoRenewFailOneMsg  = "连续包月大会员今天凌晨续费失败了，%s 0点将会再次重试。"
	_deadlineAutoRenewMsg = "您的连续包月服务将在%s 0点续费。"

	_autoRenewFailTitle     = "连续包月服务续费失败"
	_deadlineAutoRenewTitle = "连续包月服务即将续费"

	_sleep = 20 * time.Millisecond

	_maxtime = 20
)

func (s *Service) handlerinsertorderproc() {
	var (
		err error
	)
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.handlerinsertorderproc panic(%v)", x)
			go s.handlerinsertorderproc()
			log.Info("service.handlerinsertorderproc recover")
		}
	}()
	for {
		order := <-s.handlerInsertOrder
		for i := 0; i < s.c.Property.Retry; i++ {
			if err = s.insertOrder(order); err == nil {
				break
			}
			log.Error("error(%+v)", err)
		}

	}
}

func (s *Service) insertOrder(r *model.VipPayOrder) (err error) {
	var aff int64
	if aff, err = s.dao.AddPayOrder(context.TODO(), r); err != nil {
		err = errors.WithStack(err)
		return
	}
	if aff != 1 {
		return
	}
	log.Info("vip_pay_order sysn data(%+v)", r)
	rlog := new(model.VipPayOrderLog)
	rlog.Mid = r.Mid
	rlog.OrderNo = r.OrderNo
	rlog.Status = r.Status
	if _, err = s.dao.AddPayOrderLog(context.TODO(), rlog); err != nil {
		log.Error("add pay order log(%+v) error(%+v)", rlog, err)
		err = nil
	}
	return
}

func (s *Service) handlerupdateorderproc() {
	var (
		err  error
		flag bool
	)
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.handlerupdateorderproc panic(%v)", x)
			go s.handlerupdateorderproc()
			log.Info("service.handlerupdateorderproc recover")
		}
	}()
	for {
		order := <-s.handlerUpdateOrder
		flag = true
		for i := 0; i < s.c.Property.Retry; i++ {
			if err = s.updatePayOrder(context.TODO(), order); err == nil {
				flag = false
				break
			}
			log.Error("error(%+v)", err)
		}
		if flag {
			s.handlerFailPayOrder <- order
		}

	}
}

func (s *Service) handlerfailpayorderproc() {
	var (
		err error
	)
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.handlerfailpayorderproc panic(%v)", x)
			go s.handlerfailpayorderproc()
			log.Info("service.handlerfailpayorderproc  recover")
		}
	}()
	for {
		order := <-s.handlerFailPayOrder
		_time := 0
		for {
			if err = s.updatePayOrder(context.TODO(), order); err == nil {
				break
			}
			log.Error("pay order error(%+v)", err)
			_time++
			if _time > _maxtime {
				break
			}
		}
	}
}

func (s *Service) handlerfailrechargeorderproc() {
	var (
		eff int64
		err error
	)
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.handlerfailrechargeorderproc panic(%v)", x)
			go s.handlerfailrechargeorderproc()
			log.Info("service.handlerfailrechargeorderproc recover")
		}
	}()
	for {
		order := <-s.handlerFailRechargeOrder
		_time := 0
		for {
			if eff, err = s.dao.UpdateRechargeOrder(context.TODO(), order); err != nil {
				log.Error("error(%+v)", err)
				break
			}
			if eff > 0 {
				break
			}
			_time++
			if _time > _maxtime {
				break
			}
			time.Sleep(_sleep)
		}
	}
}

func (s *Service) handlerupdaterechargeorderproc() {
	var (
		eff  int64
		err  error
		flag bool
	)
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.handlerupdaterechargeorderproc panic(%v)", x)
			go s.handlerupdaterechargeorderproc()
			log.Info("service.handlerupdaterechargeorderproc recover")
		}
	}()
	for {
		order := <-s.handlerRechargeOrder
		flag = true
		for i := 0; i < s.c.Property.Retry; i++ {
			if eff, err = s.dao.UpdateRechargeOrder(context.TODO(), order); err != nil {
				log.Error("error(%+v)", err)
				continue
			}
			if eff > 0 {
				log.Info("update recharge order(%+v)", order)
				flag = false
				break
			}
			time.Sleep(_sleep)
		}
		if flag {
			s.handlerFailRechargeOrder <- order
		}
	}
}

func (s *Service) updatePayOrder(c context.Context, r *model.VipPayOrder) (err error) {
	var eff int64
	if eff, err = s.dao.UpdatePayOrderStatus(c, r); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff <= 0 {
		err = fmt.Errorf("order更新未执行(%+v)", r)
		time.Sleep(_sleep)
		return
	}
	log.Info("cur pay order update order(%+v)", r)

	rlogKey := fmt.Sprintf("%v:%v", r.OrderNo, r.Status)
	if succeed := s.dao.AddTransferLock(c, rlogKey); succeed {
		rlog := new(model.VipPayOrderLog)
		rlog.Mid = r.Mid
		rlog.OrderNo = r.OrderNo
		rlog.Status = r.Status
		if _, err = s.dao.AddPayOrderLog(context.TODO(), rlog); err != nil {
			log.Error("add pay order log(%+v) error(%+v)", rlog, err)
			err = nil
		}
		return
	}
	return
}

func (s *Service) convertPayOrder(r *model.VipPayOrderOldMsg) (res *model.VipPayOrder) {
	res = new(model.VipPayOrder)
	res.Mid = r.Mid
	res.AppID = r.AppID
	res.AppSubID = r.AppSubID
	res.BuyMonths = r.BuyMonths
	res.Money = r.Money
	res.RechargeBp = r.RechargeBp
	res.OrderNo = r.OrderNo
	res.OrderType = r.OrderType
	res.PayType = r.PayType
	res.Platform = r.Platform
	res.Status = r.Status
	res.ToMid = r.Bmid
	res.Ver = r.Ver
	res.CouponMoney = r.CouponMoney
	if paymentTime, err := time.ParseInLocation("2006-01-02 15:04:05", r.PaymentTime, time.Local); err == nil {
		res.PaymentTime = xtime.Time(paymentTime.Unix())
	}
	res.Ctime = xtime.Time(parseTime(r.Ctime).Unix())
	res.Mtime = xtime.Time(parseTime(r.Mtime).Unix())
	return
}

func (s *Service) convertPayOrderByMsg(r *model.VipRechargeOrderMsg) (res *model.VipPayOrder) {
	res = new(model.VipPayOrder)
	res.Mid = r.PayMid
	res.OrderNo = r.PayOrderNo
	res.ThirdTradeNo = r.ThirdTradeNo
	res.RechargeBp = r.RechargeBp
	return
}

func convertPayOrderOldToNew(r *model.VipPayOrderOld) (res *model.VipPayOrder) {
	res = new(model.VipPayOrder)
	res.Mid = r.Mid
	res.AppID = r.AppID
	res.AppSubID = r.AppSubID
	res.BuyMonths = r.BuyMonths
	res.Money = r.Money
	res.OrderNo = r.OrderNo
	res.OrderType = r.OrderType
	res.PayType = r.PayType
	res.Platform = r.Platform
	res.Status = r.Status
	res.ToMid = r.Bmid
	res.Ver = r.Ver
	res.PaymentTime = r.PaymentTime
	res.CouponMoney = r.CouponMoney
	return
}

//HandlerPayOrder handler pay order
func (s *Service) HandlerPayOrder() (err error) {
	var (
		size     = 1000
		oldMaxID int
	)
	if oldMaxID, err = s.dao.SelOldOrderMaxID(context.TODO()); err != nil {
		log.Error("selOldOrderMaxID error(%+v)", err)
		return
	}

	page := oldMaxID / size
	if oldMaxID%size != 0 {
		page++
	}
	for i := 0; i < page; i++ {
		startID := i * size
		endID := (i + 1) * size
		if endID > oldMaxID {
			endID = oldMaxID
		}
		var (
			res              []*model.VipPayOrderOld
			batchOrder       []*model.VipPayOrder
			orderNos         []string
			oldRechargeOrder []*model.VipRechargeOrder
		)

		rechargeMap := make(map[string]*model.VipRechargeOrder)
		if res, err = s.dao.SelOldPayOrder(context.TODO(), startID, endID); err != nil {
			log.Error("selOldPayOrder(startID:%v endID:%v) error(%+v)", startID, endID, err)
			return
		}
		for _, v := range res {
			batchOrder = append(batchOrder, convertPayOrderOldToNew(v))
		}

		for _, v := range batchOrder {
			orderNos = append(orderNos, v.OrderNo)
		}
		if oldRechargeOrder, err = s.dao.SelOldRechargeOrder(context.TODO(), orderNos); err != nil {
			return
		}
		for _, v := range oldRechargeOrder {
			rechargeMap[v.PayOrderNo] = v
		}
		for _, v := range batchOrder {
			rechargeOrder := rechargeMap[v.OrderNo]
			if rechargeOrder != nil {
				v.ThirdTradeNo = rechargeOrder.ThirdTradeNo
				v.RechargeBp = rechargeOrder.RechargeBp
			}
		}
		if err = s.dao.BatchAddPayOrder(context.TODO(), batchOrder); err != nil {
			return
		}

	}
	return
}

func (s *Service) willDedutionMsg() (err error) {
	var (
		size  = 5000
		endID int
		now   time.Time
		vips  []*model.VipUserInfo
	)
	if now, err = time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local); err != nil {
		log.Error("time.ParseInLocation(%v) error(%+v)", time.Now(), err)
		return
	}
	start := now.AddDate(0, 0, 1)

	end := start.AddDate(0, 0, 3)

	if endID, err = s.dao.SelMaxID(context.TODO()); err != nil {
		return
	}
	page := endID / size
	if endID%size != 0 {
		page++
	}
	for i := 0; i < page; i++ {
		startID := i * size
		eID := (i + 1) * size
		if vips, err = s.dao.SelVipUsers(context.TODO(), startID, eID, xtime.Time(start.Unix()), xtime.Time(end.Unix())); err != nil {
			continue
		}
		for _, v := range vips {
			if v.OverdueTime.Time().Equal(start) {
				s.dao.SendMultipMsg(context.TODO(), fmt.Sprintf("%v", v.Mid),
					_autoRenewFailTwoMsg,
					_autoRenewFailTitle,
					vipWillExpiredMsgCode,
					systemNotify)
			} else if start.AddDate(0, 0, 1).Equal(v.OverdueTime.Time()) {
				s.dao.SendMultipMsg(context.TODO(), fmt.Sprintf("%v", v.Mid),
					fmt.Sprint(_autoRenewFailOneMsg, v.OverdueTime.Time().AddDate(0, 0, -1).Format("2006-01-02")),
					_autoRenewFailTitle,
					vipWillExpiredMsgCode,
					systemNotify)
			} else if start.AddDate(0, 0, 2).Equal(v.OverdueTime.Time()) {
				s.dao.SendMultipMsg(context.TODO(), fmt.Sprintf("%v", v.Mid),
					fmt.Sprint(_deadlineAutoRenewMsg, start.Format("2006-01-02")),
					_deadlineAutoRenewTitle,
					vipWillExpiredMsgCode,
					systemNotify)
			}
		}
	}
	return
}

func (s *Service) autoRenews() (err error) {
	//var (
	//	size  = 5000
	//	endID int
	//	now   time.Time
	//	vips  []*model.VipUserInfo
	//	price float64
	//)
	//if now, err = time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local); err != nil {
	//	log.Error("time.ParseInLocation(%v) error(%+v)", time.Now(), err)
	//	return
	//}
	//
	//start := now.AddDate(0, 0, 1)
	//
	//end := start.AddDate(0, 0, 3)
	//if price, err = s.vipRPC.Price(context.TODO(), 1, xmodel.DevicePC, xmodel.AutoRenew, 1); err != nil {
	//	err = errors.WithStack(err)
	//	return
	//}
	//if endID, err = s.dao.SelMaxID(context.TODO()); err != nil {
	//	return
	//}
	//
	//page := endID / size
	//if endID%size != 0 {
	//	page++
	//}
	//for i := 0; i < page; i++ {
	//	startID := i * size
	//	eID := (i + 1) * size
	//	if vips, err = s.dao.SelVipUsers(context.TODO(), startID, eID, xtime.Time(start.Unix()), xtime.Time(end.Unix())); err != nil {
	//		err = errors.WithStack(err)
	//		continue
	//	}
	//	for _, v := range vips {
	//		var params = make(map[string]interface{}, 0)
	//		if params, err = s.vipRPC.CreateOrderPlatfrom(context.TODO(), int64(v.Mid), 0, 0, 1, price, xmodel.DevicePC, 5, xmodel.AutoRenew, ""); err != nil {
	//			log.Error("CreateOrderPlatform error(%+v)", err)
	//			continue
	//		}
	//		params["payChannelId"] = v.PayChannelId
	//		params["payChannel"] = s.c.Property.PayMapping[strconv.Itoa(int(v.PayChannelId))]
	//		if err = s.dao.PayOrder(context.TODO(), params); err != nil {
	//			log.Error("handler fail orderId->%v mid:%v", params["orderId"], v.Mid)
	//			continue
	//		}
	//		log.Info("handler success orderId:%v mid:%v", params["orderId"], v.Mid)
	//	}
	//}
	return
}

// AutoRenewJob auto renew job.
//func (s *Service) autoRenewJob() {
//	defer func() {
//		if x := recover(); x != nil {
//			log.Error("service.autoRenewJob panic(%v)", x)
//			go s.autoRenewJob()
//			log.Info("service.autoRenewJob recover")
//		}
//	}()
//	log.Info("auto renew job start.................................")
//	var err error
//	if err = s.autoRenews(); err != nil {
//		log.Error("autoRenews error(%+v)", err)
//	}
//	log.Info("auto renew job end...................................")
//}

// SendMessageJob send message job.
func (s *Service) sendMessageJob() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.sendMessageJob panic(%v)", x)
			go s.sendMessageJob()
			log.Info("service.sendMessageJob recover")
		}
	}()

	log.Info("sendMessage job start .........................")
	s.willDedutionMsg()
	log.Info("sendMessage job end .........................")
}
