package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go-common/app/service/main/ugcpay/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// TradeCreate .
func (s *Service) TradeCreate(c context.Context, platform string, mid int64, oid int64, otype string, currency string) (orderID string, payData string, err error) {
	// 1. 查询asset信息
	var (
		asset         *model.Asset
		platformPrice map[string]int64
		price         int64
		state         string
		order         *model.Order
	)
	if asset, platformPrice, err = s.AssetQuery(c, oid, otype, currency); err != nil {
		return
	}
	if asset == nil {
		err = ecode.UGCPayAssetCantBuy
		return
	}

	// 2. 判断是否已付费
	if state, err = s.AssetRelation(c, mid, oid, otype); err != nil {
		return
	}
	if state == model.AssetRelationPaid {
		err = ecode.UGCPayAssetPaid
		return
	}

	// 3. 重复下单直接返回
	// if order, err = s.dao.RawOrderUserByAsset(c, mid, oid, otype, platform, time.Now().Add(-time.Second*time.Duration(conf.Conf.Biz.Pay.OrderTTL*2))); err != nil {
	// 	return
	// }
	// if order != nil {
	// 	if order.IsPay() {
	// 		err = ecode.UGCPayAssetPaid
	// 		return
	// 	}
	// 	orderID = order.ID
	// 	payData, err = s.createPayData(c, order)
	// 	if err == nil {
	// 		// 恢复原订单成功，直接返回
	// 		return
	// 	}
	// 	// 恢复原订单失败，继续执行
	// 	log.Error("Recover Order: %+v, failed: %+v", order, err)
	// 	err = nil
	// 	order = nil
	// }

	// 4. 获得平台价格
	price = asset.PickPrice(platform, platformPrice)
	if price <= 0 {
		log.Error("Asset : oid %d , otype %s , pick price : platform %s , currency %s , price <= 0", oid, otype, platform, currency)
		err = ecode.UGCPayAssetCantBuy
		return
	}

	// 5. 检查archive付费状态
	var (
		payflag = false
	)
	switch otype {
	case model.OTypeArchive:
		if payflag, err = s.dao.ArchiveUGCPay(c, oid); err != nil {
			log.Error("s.dao.ArchiveUGCPay err : %+v", err)
			err = ecode.UGCPayDependArchiveErr
			return
		}
	}
	if !payflag {
		err = ecode.UGCPayAssetCantBuy
		return
	}

	// 6. 创建订单
	order = &model.Order{
		OrderID:  s.orderID(),
		MID:      mid,
		Biz:      model.BizAsset,
		Platform: platform,
		OID:      oid,
		OType:    otype,
		Fee:      price,
		RealFee:  asset.Price,
		Currency: currency,
		State:    model.OrderStateCreated,
	}
	orderID = order.OrderID

	// 7. 创建支付参数
	if payData, err = s.createPayData(c, order); err != nil {
		return
	}

	if order.ID, err = s.dao.InsertOrderUser(c, order); err != nil {
		return
	}

	// 8. 清理订单状态缓存
	s.cache.Save(func() {
		if theErr := s.dao.DelCacheOrderUser(context.Background(), order.OrderID); theErr != nil {
			log.Error("%+v", theErr)
			return
		}
	})
	return
}

// TradeQuery .
func (s *Service) TradeQuery(ctx context.Context, orderID string) (order *model.Order, err error) {
	order, err = s.dao.OrderUser(ctx, orderID)
	if order == nil {
		err = ecode.UGCPayOrderInvalid
		return
	}
	if order.IsPay() {
		return
	}
	if err = s.tradeConfirm(ctx, order); err != nil {
		return
	}
	return
}

// TradeCancel .
func (s *Service) TradeCancel(ctx context.Context, orderID string) (err error) {
	if err = s.updateOrder(ctx, orderID, model.OrderStateClosed, "TradeCancel"); err != nil {
		return
	}
	// 支付中心取消
	if err = s.cancelOrder(ctx, orderID); err != nil {
		return
	}
	return
}

// TradeConfirm a trade state from pay
func (s *Service) TradeConfirm(ctx context.Context, orderID string) (order *model.Order, err error) {
	// 1. 查询本地订单
	if order, err = s.dao.OrderUser(ctx, orderID); err != nil {
		return
	}
	if order == nil {
		err = ecode.UGCPayOrderInvalid
		return
	}
	if err = s.tradeConfirm(ctx, order); err != nil {
		return
	}
	return
}

func (s *Service) tradeConfirm(ctx context.Context, order *model.Order) (err error) {
	// 2. 从支付平台获取订单状态
	var (
		params    url.Values
		jsonData  string
		orders    map[string][]*model.PayOrder
		payOrders []*model.PayOrder
		ok        bool
	)
	params = s.pay.Query(order.OrderID)
	if err = s.pay.Sign(params); err != nil {
		return
	}
	if jsonData, err = s.pay.ToJSON(params); err != nil {
		return
	}
	if orders, err = s.dao.PayQuery(ctx, jsonData); err != nil {
		return
	}
	if payOrders, ok = orders[order.OrderID]; !ok || len(payOrders) == 0 {
		log.Info("TradeConfirm from pay platform not found order: %s", order.OrderID)
		return
	}

	// 3. 变更订单状态
	var (
		changed       bool
		txID          = ""
		payStatusDesc = ""
		payFailReason = ""
	)

	for _, o := range payOrders {
		changed = false
		switch o.PayStatus {
		case model.PayStateSuccess, model.PayStateFinished:
			changed = order.UpdateState(model.OrderStatePaid)
		case model.PayStateOverdue:
			changed = order.UpdateState(model.OrderStateExpired)
		case model.PayStateClosed, model.PayStatePaySuccessAndCancel, model.PayStatePayCancel:
			changed = order.UpdateState(model.OrderStateClosed)
		case model.PayStateFail:
			changed = order.UpdateState(model.OrderStateFailed)
		default:
			log.Error("TradeConfirm unknown pay status: %+v", o)
		}
		if changed {
			payStatusDesc = o.PayStatusDesc
			payFailReason = o.FailReason
			txID = strconv.FormatInt(o.TXID, 10)
		}
	}
	if !changed {
		log.Info("tradeConfirm: %+v no need to update", order)
		return
	}

	orderUpdater := func(order *model.Order) {
		order.PayID = txID
		order.PayReason = s.payReason(payStatusDesc, payFailReason)
	}
	return s.updateOrder(ctx, order.OrderID, order.State, "TradeConfirm", orderUpdater)
}

func (s *Service) payReason(desc string, failReason string) string {
	if failReason == "" {
		return desc
	}
	return desc + ":" + failReason
}

// TradeRefunds .
func (s *Service) TradeRefunds(ctx context.Context, orderIDs []string) (err error) {
	for _, id := range orderIDs {
		if err = s.TradeRefund(ctx, id); err != nil {
			return
		}
	}
	return
}

// TradeRefund .
func (s *Service) TradeRefund(ctx context.Context, orderID string) (err error) {
	var (
		order *model.Order
	)
	if order, err = s.dao.OrderUser(ctx, orderID); err != nil {
		return
	}
	if order == nil {
		err = ecode.UGCPayOrderInvalid
		return
	}
	if order.PayID == "" {
		err = ecode.UGCPayOrderNotPay
		return
	}

	var toState string
	if order.State == model.OrderStatePaid {
		toState = model.OrderStateRefunding
	} else if order.State == model.OrderStateSettled {
		toState = model.OrderStateSettledRefunding
	} else {
		err = ecode.UGCPayOrderNotPay
		return
	}

	extra := func(ctx context.Context) error {
		// 通知支付平台退款
		return s.refundOrder(ctx, order, "后台退款")
	}
	return s.updateOrderTrans(ctx, orderID, toState, "TradeRefund", extra)
}

// TradePayCallback 支付订单回调
func (s *Service) TradePayCallback(ctx context.Context, msgID int64, msgData string) (retMSG string, err error) {
	// 1. 解析msg
	retMSG = "SUCCESS"
	msg := &model.PayCallbackMSG{}
	if err = json.Unmarshal([]byte(msgData), msg); err != nil {
		log.Error("TradePayCallback decode err : %+v", errors.WithStack(err))
		err = nil
		return
	}
	var (
		order *model.Order
	)

	// 2. 获得order
	if order, err = s.dao.OrderUser(ctx, msg.OrderID); err != nil {
		log.Error("TradePayCallback s.dao.OrderUser err : %+v", err)
		err = nil
		retMSG = "FAIL"
		return
	}
	if order == nil {
		log.Error("TradePayCallback order not found err : %+v", err)
		return
	}

	// 3. 校验价格
	if msg.PayAmount != order.Fee {
		log.Error("TradePayCallback check fee not equal, msg: %+v, order: %+v", msg, order)
		return
	}

	// 4. 更新order
	orderUpdater := func(order *model.Order) {
		order.PayID = strconv.FormatInt(msg.TXID, 10)
		order.ParsePaidTime(msg.OrderPayTime)
	}
	if err = s.updateOrder(ctx, order.OrderID, model.OrderStatePaid, "TradePayCallback", orderUpdater); err != nil {
		log.Error("TradePayCallback s.updateOrder err : %+v", err)
		err = nil
		retMSG = "FAIL"
		return
	}
	return
}

// TradeRefundCallback 交易退款回调
func (s *Service) TradeRefundCallback(ctx context.Context, msgID int64, msgData string) (retMSG string, err error) {
	// 1. 解析msg
	retMSG = "SUCCESS"
	msg := &model.PayRefundCallbackMSG{}
	if err = json.Unmarshal([]byte(msgData), msg); err != nil {
		log.Error("TradeRefundCallback decode err : %+v", errors.WithStack(err))
		err = nil
		return
	}

	// 检查是否成功退款
	flag := false
	refundTime := ""
	for _, ele := range msg.List {
		if ele.IsSuccess() {
			flag = true
			refundTime = ele.RefundEndTime
			log.Info("TradeRefundCallback got refunded msg: %+v", ele)
		}
	}
	if !flag {
		log.Warn("TradeRefundCallback got msg with not refunded state: %+v", msg)
		return
	}

	// 2. 获得order
	var (
		order *model.Order
	)
	if order, err = s.dao.OrderUser(ctx, msg.OrderID); err != nil {
		log.Error("TradeRefundCallback s.dao.OrderUser err : %+v", err)
		err = nil
		retMSG = "FAIL"
		return
	}
	if order == nil {
		log.Error("TradeRefundCallback order not found err : %+v", err)
		return
	}

	// 3. 更新order
	var toState string
	switch order.State {
	case model.OrderStateRefunding, model.OrderStatePaid:
		toState = model.OrderStateRefunded
	case model.OrderStateSettledRefunding, model.OrderStateSettled:
		toState = model.OrderStateSettledRefunded
	default:
		toState = model.OrderStateDupRefunded
	}

	orderUpdater := func(order *model.Order) {
		order.PayID = strconv.FormatInt(msg.TXID, 10)
		order.ParseRefundedTime(refundTime)
	}
	if err = s.updateOrder(ctx, order.OrderID, toState, "TradeRefundCallback", orderUpdater); err != nil {
		log.Error("TradeRefundCallback s.updateOrder err : %+v", err)
		err = nil
		retMSG = "FAIL"
		return
	}
	return
}

// RechargeShellCallback 转换贝壳回调
func (s *Service) RechargeShellCallback(ctx context.Context, msgID int64, msgData string) (retMSG string, err error) {
	// 1. 解析msg
	retMSG = "SUCCESS"
	msg := &model.RechargeShellCallbackMSG{}
	if err = json.Unmarshal([]byte(msgData), msg); err != nil {
		log.Error("RechargeShellCallback decode err : %+v", errors.WithStack(err))
		err = nil
		return
	}

	// 查询转贝壳订单
	order, err := s.dao.RawOrderRechargeShell(ctx, msg.ThirdOrderNo)
	if err != nil {
		retMSG = "FAIL"
		log.Error("%+v", err)
		err = nil
		return
	}
	if order == nil {
		log.Error("RechargeShellCallback not found valid order: %+v", msg)
		return
	}

	// 决定状态流转
	state := model.RechargeShellFail
	switch msg.Status {
	case "SUCCESS":
		state = model.RechargeShellSuccess
	case "FAIL":
		state = model.RechargeShellFail
	default:
		log.Error("RechargeShellCallback got unknown statue: %+v", msg)
	}
	var (
		rechargeShellLog = &model.OrderRechargeShellLog{
			OrderID:           order.OrderID,
			FromState:         "created",
			ToState:           state,
			Desc:              msg.Status,
			BillUserMonthlyID: fmt.Sprintf("end_%d", time.Now().Unix()),
		}
	)
	order.State = state

	// 开始事务
	tx, err := s.dao.BeginTran(ctx)
	if err != nil {
		log.Error("%+v", err)
		err = nil
		return
	}
	if err = s.dao.TXUpdateOrderRechargeShell(ctx, tx, order); err != nil {
		tx.Rollback()
		log.Error("%+v", err)
		err = nil
		retMSG = "FAIL"
		return
	}
	if _, err = s.dao.TXInsertOrderRechargeShellLog(ctx, tx, rechargeShellLog); err != nil {
		tx.Rollback()
		log.Error("%+v", err)
		err = nil
		retMSG = "FAIL"
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("%+v", err)
		err = nil
		retMSG = "FAIL"
		return
	}
	return
}

func (s *Service) updateOrder(ctx context.Context, orderID string, toState string, taskName string, orderUpdaters ...func(order *model.Order)) (err error) {
	return s.updateOrderTrans(ctx, orderID, toState, taskName, nil, orderUpdaters...)
}

func (s *Service) updateOrderTrans(ctx context.Context, orderID string, toState string, taskName string, extra func(context.Context) error, orderUpdaters ...func(order *model.Order)) (err error) {
	// 需要异步处理的函数列表
	asyncFuncList := make([]func() error, 0)
	fn := func(ctx context.Context) (affected bool, err error) {
		var (
			tx    *xsql.Tx
			order *model.Order
		)
		affected = true

		// 查找 order
		// if order, err = s.dao.OrderUser(ctx, orderID); err != nil {
		// 	return
		// }
		if order, err = s.dao.RawOrderUser(ctx, orderID); err != nil {
			return
		}
		if order == nil {
			log.Error("updateOrder: %s get nil order: %s", taskName, orderID)
			err = ecode.UGCPayOrderInvalid
			return
		}
		log.Info("updateOrderTrans: orderID: %s, toState: %s, order: %+v", orderID, toState, order)
		preState := order.State

		// 通过状态机，内存变更
		changed := order.UpdateState(toState)
		if !changed {
			log.Error("updateOrder: %s, change state from order: %+v to state: %s failed", taskName, order, toState)
			return
		}
		for _, ou := range orderUpdaters {
			ou(order)
		}

		// 更新DB
		var (
			rowAffected int64
			logOrder    = &model.LogOrder{
				OrderID:   order.OrderID,
				FromState: preState,
				ToState:   order.State,
				Desc:      taskName,
			}
		)
		if tx, err = s.dao.BeginTran(ctx); err != nil {
			return
		}
		if rowAffected, err = s.dao.TXUpdateOrderUser(ctx, tx, order); err != nil {
			tx.Rollback()
			return
		}
		if rowAffected <= 0 {
			affected = false
			log.Error("updateOrder: %s, TXUpdateOrderUser: %+v rowAffected: %d <= 0", taskName, order, rowAffected)
			tx.Rollback()
			return
		}
		if _, err = s.dao.TXInsertOrderUserLog(ctx, tx, logOrder); err != nil {
			tx.Rollback()
			return
		}

		// 根据更新后 state 各种后处理
		switch order.State {
		case model.OrderStatePaid: // 支付订单完成
			var (
				relation = &model.AssetRelation{
					OID:   order.OID,
					OType: order.OType,
					MID:   order.MID,
					State: model.OrderStatePaid,
				}
				assetRelationAffected int64
			)
			if assetRelationAffected, err = s.dao.TXUpsertAssetRelation(ctx, tx, relation); err != nil {
				tx.Rollback()
				return
			}
			// 如果 asset_relation 更新失败, 认为是重复订单, 这里不对 order 的 state 做流转
			if assetRelationAffected <= 0 {
				asyncFuncList = append(asyncFuncList, func() error {
					return s.refundOrder(context.Background(), order, "重复扣费退款")
				})
			}
			// 清理付费关系缓存
			asyncFuncList = append(asyncFuncList, func() error {
				return s.dao.DelCacheAssetRelationState(context.Background(), relation.OID, relation.OType, relation.MID)
			})
		case model.OrderStateRefunded, model.OrderStateSettledRefunded: // 支付订单退款
			var (
				relation = &model.AssetRelation{
					OID:   order.OID,
					OType: order.OType,
					MID:   order.MID,
					State: "none",
				}
			)
			if _, err = s.dao.TXUpsertAssetRelation(ctx, tx, relation); err != nil {
				tx.Rollback()
				return
			}
			// 清理付费关系缓存
			asyncFuncList = append(asyncFuncList, func() error {
				return s.dao.DelCacheAssetRelationState(context.Background(), relation.OID, relation.OType, relation.MID)
			})
		default: // 默认
		}

		// 额外任务执行
		if extra != nil {
			if err = extra(ctx); err != nil {
				tx.Rollback()
				return
			}
		}

		// 提交事务
		if err = tx.Commit(); err != nil {
			return
		}
		return
	}
	if err = runCAS(ctx, fn); err != nil {
		return
	}
	asyncFuncList = append(asyncFuncList, func() error {
		return s.dao.DelCacheOrderUser(context.Background(), orderID)
	})

	// 启动需要异步处理的函数
	for _, f := range asyncFuncList {
		f := f
		s.cache.Save(func() {
			if theErr := f(); theErr != nil {
				log.Error("updateOrder: %s, err: %+v", taskName, theErr)
			}
		})
	}
	return
}

// orderID get order id
func (s *Service) orderID() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%05d", s.rnd.Int63n(99999)))
	b.WriteString(fmt.Sprintf("%03d", time.Now().UnixNano()/1e6%1000))
	b.WriteString(time.Now().Format("060102150405"))
	return b.String()
}

func (s *Service) payTitle(oid int64, fee int64) string {
	return fmt.Sprintf("付费视频观看（av%d）消费%.2fB币", oid, float64(fee)/100)
}
