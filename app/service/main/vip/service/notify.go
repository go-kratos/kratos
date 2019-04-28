package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	_decimal       = 100
	_remarkBuy     = "充值开通"
	_remarkGift    = "还有赠送"
	_remarkAutoBuy = "自动续费订单"
	_refundSuccess = "REFUND_SUCCESS"
	_yearMonth     = 12
)

// PayNotify pay notify.
func (s *Service) PayNotify(c context.Context, p *model.PayCallBackResult) (err error) {
	var (
		o  *model.OrderInfo
		hv *model.HandlerVip
	)
	if o, err = s.OrderInfo(c, p.OutTradeNO); err != nil {
		log.Error("%+v", err)
		return
	}
	if o == nil {
		err = fmt.Errorf("order not exist(%s)", p.OutTradeNO)
		return
	}
	if p.Bp == 0 || p.Bp == o.Money {
		err = fmt.Errorf("pay amount no matching (%f)", o.Money)
		return
	}
	n := &model.PayNotifyContentOld{
		OrderID: p.OutTradeNO,
		TradeNO: p.TradeNO,
	}
	if hv, err = s.dealWithPayOrder(c, o, n); err != nil {
		log.Error("%+v", err)
		return
	}
	if hv != nil {
		s.asyncBcoin(func() {
			s.ProcesserHandler(context.TODO(), hv, "127.0.0.1")
		})
		s.cache(func() {
			s.dao.DelVipInfoCache(context.TODO(), int64(hv.Mid))
		})
	}
	return
}

// PayNotify2 pay notify2 for new pay platform.
func (s *Service) PayNotify2(c context.Context, n *model.PayNotifyContent) (err error) {
	var (
		o        *model.OrderInfo
		hv       *model.HandlerVip
		vip      *model.VipInfoDB
		vipCache *model.VipInfo
		val      int
		ip       = metadata.String(c, metadata.RemoteIP)
	)
	n.ExpiredTime = n.ExpiredTime / 1e3
	if s.c.PayConf.CustomerID != n.CustomerID {
		err = fmt.Errorf("biz customer id not matching (%v)", n.CustomerID)
		return
	}
	if n.OrderID == "" {
		err = fmt.Errorf("order no is not exist")
		return
	}
	if o, err = s.OrderInfo(c, n.OrderID); err != nil {
		log.Error("%+v", err)
		return
	}
	if o == nil {
		err = fmt.Errorf("order not exist(%s)", n.OrderID)
		return
	}

	if o.Status == model.SUCCESS && o.OrderType != model.IOSRenew {
		err = ecode.VipOrderAlreadyHandlerErr
		return
	}

	if vipCache, err = s.VipInfo(c, o.Mid); err != nil {
		err = errors.WithStack(err)
		return
	}

	if o.OrderType == model.IOSRenew && vipCache.IosOverdueTime.Time().Unix() >= n.ExpiredTime {
		err = ecode.VipOrderAlreadyHandlerErr
		return
	}

	if o.OrderType != model.IOSRenew && (o.Money == 0 || (o.Money*float64(_decimal)) != float64(n.PayAmount)) {
		err = fmt.Errorf("pay amount no matching (%d)", n.PayAmount)
		return
	}
	o.ThirdTradeNo = fmt.Sprintf("%v", n.TxID)
	o.RechargeBP = float64(n.PayAmount) / float64(_decimal)
	o.PayType = n.PayChannel
	if n.PayStatus == "SUCCESS" {
		if o.OrderType == model.IOSRenew {
			if hv, err = s.dealWithIapOrder(c, n, o); err != nil {
				err = errors.WithStack(err)
				return
			}
		} else {
			if hv, err = s.dealWithSuccessOrder(c, o, n); err != nil {
				err = errors.WithStack(err)
				return
			}
		}

		if val, err = s.dao.GetSignVip(c, hv.Mid); err != nil {
			log.Error("%+v", err)
			err = nil
		}
		if vip, err = s.dao.VipInfo(context.TODO(), o.Mid); err != nil {
			err = errors.WithStack(err)
			return
		}
		if o.OrderType == model.AutoRenew && int8(vip.VipPayType) == model.AutoRenew && val == 0 {
			s.asyncBcoin(func() {
				s.dao.SendMultipMsg(context.TODO(), strconv.FormatInt(o.Mid, 10),
					fmt.Sprintf(_autoRenewMsg, vip.VipOverdueTime.Time().Format("2006-01-02"), vip.VipOverdueTime.Time().AddDate(0, 0, -1).Format("2006-01-02")),
					_autoRenewTitle,
					model.VipOpenMsgCode,
					ip,
					model.VipSystemNotify)
				s.addBcoin(context.TODO(), o.Mid, hv.VipUser, ip)
			})
			s.cache(func() {
				s.dao.DelVipInfoCache(context.TODO(), o.Mid)
				s.dao.DelVipInfoCache(context.TODO(), hv.ToMid)
			})

		} else {
			s.asyncBcoin(func() {
				s.ProcesserHandler(context.TODO(), hv, ip)
			})
		}
	} else {
		if err = s.dealWithFailOrder(c, o); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	return
}

func (s *Service) dealWithIapOrder(c context.Context, n *model.PayNotifyContent, po *model.OrderInfo) (hv *model.HandlerVip, err error) {
	var (
		tx    *sql.Tx
		vip   *model.VipInfoDB
		eff   int64
		ver   int64
		pay   *model.PayOrder
		bo    = new(model.VipChangeBo)
		order = new(model.OrderInfo)
	)
	if tx, err = s.dao.StartTx(c); err != nil {
		err = errors.WithStack(err)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	if vip, err = s.dao.TxSelVipUserInfo(tx, po.Mid); err != nil {
		err = errors.WithStack(err)
		return
	}

	if vip == nil {
		vi := new(model.VipInfoDB)
		vi.IosOverdueTime = xtime.Time(n.ExpiredTime)
		vi.VipStartTime = xtime.Time(time.Now().Unix())
		vi.VipStatus = model.Expire
		vi.VipType = model.Vip
		vi.Mid = po.Mid
		if eff, err = s.dao.TxAddIosVipUserInfo(tx, vi); err != nil {
			err = errors.WithStack(err)
			return
		}

		if eff <= 0 {
			err = ecode.VipOrderAlreadyHandlerErr
			return
		}
	} else {

		if eff, err = s.dao.TxUpdateIosUserInfo(tx, xtime.Time(n.ExpiredTime), po.Mid); err != nil {
			err = errors.WithStack(err)
			return
		}

		if eff <= 0 {
			err = ecode.VipOrderAlreadyHandlerErr
			return
		}
	}

	order.Status = model.SUCCESS
	order.PayType = po.PayType
	order.Ver = po.Ver + 1
	order.OrderNo = po.OrderNo
	order.ID = po.ID
	order.ThirdTradeNo = po.ThirdTradeNo
	order.RechargeBP = po.RechargeBP
	ver = po.Ver

	if po.Status == model.SUCCESS {
		if pay, err = s.AddPayOrder(c, po.Mid, po.ToMid, po.BuyMonths, po.Money, int(po.Platform), model.AutoRenew, po.AppID, po.AppSubID); err != nil {
			err = errors.WithStack(err)
			return
		}
		order.OrderType = model.AutoRenew
		order.OrderNo = pay.OrderNo
		order.ID = pay.ID
		order.Ver = pay.Ver + 1
		ver = pay.Ver
	}

	if err = s.dao.TxUpdateIosPayOrder(tx, order, ver); err != nil {
		err = errors.WithStack(err)
		return
	}

	bo.Mid = po.Mid
	bo.RelationID = order.OrderNo
	bo.Months = po.BuyMonths
	bo.ChangeTime = xtime.Time(time.Now().Unix())
	bo.ChangeType = model.Recharge
	bo.Remark = _remarkAutoBuy

	if hv, err = s.UpdateVipWithHistory(c, tx, bo); err != nil {
		return
	}

	if err = s.dao.TxUpdateIosRenewUserInfo(tx, int64(n.PayChannelID), hv.VipUser.Ver+1, hv.VipUser.Ver, po.Mid, model.AuoRenewVip); err != nil {
		err = errors.WithStack(err)
		return
	}

	if err = s.dao.TxDupUserDiscount(tx, po.Mid, model.FirstDiscountBuyVip, po.OrderNo, model.Used); err != nil {
		err = errors.WithStack(err)
		return
	}
	olog := new(model.VipPayOrderLog)
	olog.Mid = po.Mid
	olog.OrderNo = order.OrderNo
	olog.Status = model.SUCCESS
	if err = s.dao.TxAddOrderLog(tx, olog); err != nil {
		err = errors.WithStack(err)
		return
	}
	olog = new(model.VipPayOrderLog)
	olog.Mid = po.Mid
	olog.OrderNo = po.OrderNo
	olog.Status = model.Sign
	if err = s.dao.TxAddOrderLog(tx, olog); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//PaySignNotify def.
func (s *Service) PaySignNotify(c context.Context, n *model.PaySignNotify) (err error) {
	if n.CustomerID != s.c.PayConf.CustomerID {
		err = fmt.Errorf("biz customerID not match")
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	s.dao.SetSignVip(c, n.UID, 2*60)
	if "ADD" == n.ChangeType {
		if err = s.dealWithSign(c, n.UID, model.AuoRenewVip); err != nil {
			err = errors.WithStack(err)
			return
		}
		go func() {
			s.dao.SendMultipMsg(context.TODO(),
				fmt.Sprintf("%v", n.UID),
				_openAutoRenewMsg,
				_openAutoRenewTitle, model.VipOpenMsgCode,
				ip, model.VipSystemNotify)
		}()
	} else if "DELETE" == n.ChangeType {
		if err = s.dealWithSign(c, n.UID, model.General); err != nil {
			err = errors.WithStack(err)
			return
		}
		var vip *model.VipInfo
		if vip, err = s.VipInfo(c, n.UID); err != nil {
			log.Error("vip info error(%v)", err)
		}

		go func() {
			s.dao.SendMultipMsg(context.TODO(),
				fmt.Sprintf("%v", n.UID),
				fmt.Sprintf(_cancelAutoRenewMsg, vip.VipOverdueTime.Time().Format("2006-01-02")),
				_cancelAutoRenewTitle, model.VipOpenMsgCode,
				ip, model.VipSystemNotify)
		}()
	}
	return
}

func (s *Service) dealWithPayOrder(c context.Context, o *model.OrderInfo, n *model.PayNotifyContentOld) (hv *model.HandlerVip, err error) {
	var tx *sql.Tx

	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("%+v", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback %+v", err)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit %+v", err)
		}
	}()
	if o.Status == model.SUCCESS {
		return
	}

	o.ThirdTradeNo = fmt.Sprintf("%v", n.TradeNO)
	o.PayType = n.PayChannel

	if err = s.dao.TxUpdateOrderStatus(c, tx, model.SUCCESS, o.PayType, o.ThirdTradeNo, o.OrderNo); err != nil {
		log.Error("%+v", err)
		return
	}
	bo := new(model.VipChangeBo)
	if o.ToMid == 0 {
		bo.Mid = o.Mid
		bo.Remark = model.RemarkBuy
	} else {
		bo.Mid = o.ToMid
		bo.Remark = model.RemarkGift
	}
	bo.ChangeType = model.Recharge
	bo.RelationID = o.OrderNo
	bo.Months = o.BuyMonths
	bo.Days = int64(o.BuyMonths) * model.VipDaysMonth
	if hv, err = s.UpdateVipWithHistory(c, tx, bo); err != nil {
		log.Error("%+v", err)
		return
	}
	return
}

func (s *Service) dealWithSuccessOrder(c context.Context, po *model.OrderInfo, n *model.PayNotifyContent) (hv *model.HandlerVip, err error) {
	var (
		a  int64
		tx *sql.Tx
		bo = new(model.VipChangeBo)
	)
	if tx, err = s.dao.StartTx(c); err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			if err = tx.Commit(); err != nil {
				log.Error("dealWithSuccessOrder commit() error(%+v)", err)
				tx.Rollback()
			}
		}
	}()
	if a, err = s.dao.TxUpdatePayOrderStatus(tx, model.SUCCESS, po.ID, po.Ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a != 1 {
		err = ecode.VipOrderAlreadyHandlerErr
		return
	}
	ver := po.Ver
	po.Ver++
	if err = s.dao.TxUpdatePayOrder(tx, po, ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	bo.Mid = po.Mid
	bo.Days = int64(po.BuyMonths) * model.VipDaysMonth
	bo.Months = po.BuyMonths
	bo.ChangeTime = xtime.Time(time.Now().Unix())
	bo.ChangeType = model.Recharge
	bo.Remark = _remarkBuy
	bo.RelationID = po.OrderNo
	if po.OrderType == model.AutoRenew {
		bo.Remark = _remarkAutoBuy
	}
	if po.ToMid != 0 {
		bo.Remark = _remarkGift
		bo.Mid = po.ToMid
	}

	if hv, err = s.UpdateVipWithHistory(c, tx, bo); err != nil {
		return
	}
	ver = hv.VipUser.Ver
	if po.OrderType == model.AutoRenew {
		if err = s.dao.TxUpdateChannelID(tx, po.Mid, n.PayChannelID, ver+1, ver); err != nil {
			err = errors.WithStack(err)
			return
		}
		if err = s.dao.TxDupUserDiscount(tx, po.Mid, model.FirstDiscountBuyVip, po.OrderNo, model.Used); err != nil {
			err = errors.WithStack(err)
			return
		}
	}

	hv.Mid = po.Mid
	hv.ToMid = po.ToMid
	olog := new(model.VipPayOrderLog)
	olog.Mid = po.Mid
	olog.OrderNo = po.OrderNo
	olog.Status = model.SUCCESS
	if err = s.dao.TxAddOrderLog(tx, olog); err != nil {
		err = errors.WithStack(err)
		return
	}

	if po.OrderType == model.AutoRenew {
		olog = new(model.VipPayOrderLog)
		olog.Status = model.Sign
		olog.OrderNo = po.OrderNo
		olog.Mid = po.Mid
		if err = s.dao.TxAddOrderLog(tx, olog); err != nil {
			err = errors.WithStack(err)
		}
	}

	return
}

func (s *Service) dealWithFailOrder(c context.Context, po *model.OrderInfo) (err error) {
	var (
		tx *sql.Tx
	)
	if tx, err = s.dao.StartTx(c); err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			if err = tx.Commit(); err != nil {
				log.Error("dealWithFailOrder commit() error(%+v)", err)
				tx.Rollback()
			}
		}
	}()
	if _, err = s.dao.TxUpdatePayOrderStatus(tx, model.FAILED, po.ID, po.Ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	olog := new(model.VipPayOrderLog)
	olog.Status = model.FAILED
	olog.Mid = po.Mid
	olog.OrderNo = po.OrderNo
	if err = s.dao.TxAddOrderLog(tx, olog); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//dealWithSign deal with sign
func (s *Service) dealWithSign(c context.Context, mid int64, status int8) (err error) {
	var (
		a  int64
		tx *sql.Tx
		vb *model.VipInfoDB
	)
	if tx, err = s.dao.StartTx(c); err != nil {
		err = errors.WithStack(err)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	if vb, err = s.dao.TxSelVipUserInfo(tx, mid); err != nil {
		err = errors.WithStack(err)
		return
	}

	if vb == nil {
		err = fmt.Errorf("not exit this user")
		return
	}

	if a, err = s.dao.UpdatePayType(tx, mid, status, vb.Ver+1, vb.Ver); err != nil {
		err = errors.WithStack(err)
		return
	}

	if a != 1 {
		err = fmt.Errorf("operate fail")
		return
	}
	return
}

//RefundNotify pay refund notify.
func (s *Service) RefundNotify(c context.Context, arg *model.PayRefundNotify) (err error) {
	var (
		order    *model.OrderInfo
		orderLog *model.VipPayOrderLog
	)
	if order, err = s.dao.OrderInfo(c, arg.OrderID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if order == nil {
		err = ecode.VipOrderNoErr
		return
	}
	if order.Status != model.SUCCESS {
		err = ecode.VipOrderStatusPayingErr
		return
	}
	if len(arg.BatchRefundList) > 0 && _refundSuccess == arg.BatchRefundList[0].RefundStatus {
		if orderLog, err = s.dao.SelPayOrderLog(c, order.OrderNo, arg.BatchRefundList[0].CustomerRefundID, model.REFUNDED); err != nil {
			err = errors.WithStack(err)
			return
		}
		if orderLog != nil {
			log.Info("current refund already reduce order->%+v refund->%+v", order, arg.BatchRefundList[0])
			return
		}
		key := fmt.Sprintf("lock:%v:%v:%v", order.OrderNo, arg.BatchRefundList[0].CustomerRefundID, model.REFUNDED)
		if success := s.dao.AddTransferLock(c, key); success {
			if err = s.dealWithRefund(c, order, arg.BatchRefundList[0]); err != nil {
				err = errors.WithStack(err)
				s.dao.DelCache(c, key)
				return
			}
			s.cache(func() {
				s.dao.DelVipInfoCache(context.TODO(), order.Mid)
			})
			ip := metadata.String(c, metadata.RemoteIP)
			go func() {
				s.dao.SendMultipMsg(context.TODO(), fmt.Sprintf("%v", order.Mid), fmt.Sprintf(_vipRefundMsgContext, order.PaymentTime.Time().Format("2006-01-02"), order.BuyMonths),
					_refundTitle,
					model.VipOpenMsgCode,
					ip,
					model.VipSystemNotify)
			}()
			s.dao.DelCache(context.TODO(), key)
		}
	} else {
		if err = s.dealWithRefundFail(c, order, arg.BatchRefundList[0].CustomerRefundID); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	return
}

func (s *Service) dealWithRefund(c context.Context, order *model.OrderInfo, arg *model.PayRefundList) (err error) {
	var (
		tx          *sql.Tx
		oldTx       *sql.Tx
		user        *model.VipInfoDB
		totalAmount = s.yuanToFen(order.Money)
		vch         = new(model.VipChangeHistory)
		eff         int64
	)
	if tx, err = s.dao.StartTx(c); err != nil {
		err = errors.WithStack(err)
		return
	}
	if oldTx, err = s.dao.OldStartTx(c); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			oldTx.Rollback()
			return
		}
		tx.Commit()
		oldTx.Commit()
	}()
	olog := new(model.VipPayOrderLog)
	olog.Mid = order.Mid
	olog.OrderNo = order.OrderNo
	olog.RefundID = arg.CustomerRefundID
	olog.RefundAmount = s.fenToYuan(float64(arg.RefundAmount))
	olog.Status = model.REFUNDED
	if err = s.dao.TxAddOrderLog(tx, olog); err != nil {
		err = errors.WithStack(err)
		return
	}
	if user, err = s.dao.OldTxSelVipUserInfo(oldTx, order.Mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	days := int64(order.BuyMonths) * model.VipDaysMonth
	if order.BuyMonths%_yearMonth == 0 {
		days = model.VipDaysYear
	}
	refundDay := int(math.Floor(float64(float64(arg.RefundAmount)/float64(totalAmount)*float64(days)) + 0.5))
	now := time.Now()
	avtime := user.AnnualVipOverdueTime.Time().AddDate(0, 0, -refundDay)
	overTime := user.VipOverdueTime.Time().AddDate(0, 0, -refundDay)
	user.AnnualVipOverdueTime = xtime.Time(avtime.Unix())
	user.VipOverdueTime = xtime.Time(overTime.Unix())
	if user.VipType == model.AnnualVip {
		delMonth := avtime
		if now.After(avtime) {
			delMonth = now
		}
		delMonth.AddDate(0, 1, 1-delMonth.Day())
		if err = s.dao.OldTxDelBcoinSalary(oldTx, order.Mid, delMonth); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	if now.After(avtime) {
		user.VipType = model.Vip
		user.AnnualVipOverdueTime = xtime.Time(now.Unix())
	}
	if now.After(overTime) {
		user.VipStatus = model.Expire
		user.VipOverdueTime = xtime.Time(now.Unix())
	}
	ver := user.Ver
	user.Ver++
	if eff, err = s.dao.OldTxUpdateVipUserInfo(oldTx, user, ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff <= 0 {
		err = ecode.VipRefundErr
		return
	}
	vch.ChangeTime = xtime.Time(time.Now().Unix())
	vch.ChangeType = model.SystemDeduction
	vch.Days = int64(-refundDay)
	vch.Mid = order.Mid
	vch.RelationID = order.OrderNo
	vch.Remark = "系统扣减"
	if err = s.dao.OldTxAddChangeHistory(oldTx, vch); err != nil {
		return
	}
	ver = order.Ver
	order.RefundAmount += s.fenToYuan(float64(arg.RefundAmount))
	if err = s.dao.TxUpdatePayOrderRefundAmount(tx, order.ID, order.RefundAmount, ver+1, ver); err != nil {
		err = errors.WithStack(err)
	}
	//TODO DELETE VIEW COUPON
	return
}

func (s *Service) dealWithRefundFail(c context.Context, arg *model.OrderInfo, refundID string) (err error) {
	var tx *sql.Tx
	if tx, err = s.dao.StartTx(c); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	olog := new(model.VipPayOrderLog)
	olog.Status = model.REFUNDFAIL
	olog.Mid = arg.Mid
	olog.OrderNo = arg.OrderNo
	olog.RefundID = refundID
	if err = s.dao.TxAddOrderLog(tx, olog); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
