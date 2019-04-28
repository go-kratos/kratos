package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go-common/app/job/main/ugcpay/dao"
	"go-common/app/job/main/ugcpay/model"
	"go-common/app/job/main/ugcpay/service/pay"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

type taskBillDaily struct {
	dao        *dao.Dao
	pay        *pay.Pay
	rnd        *rand.Rand
	dayOffset  int
	namePrefix string
	tl         *taskLog
}

func (s *taskBillDaily) Run() (err error) {
	var (
		ctx      = context.Background()
		finished bool
		expectFN = func(ctx context.Context) (expect int64, err error) {
			var (
				beginTime, endTime         = dayRange(s.dayOffset)
				expectPaid, expectRefunded int64
			)
			if expectPaid, err = s.dao.CountPaidOrderUser(ctx, beginTime, endTime); err != nil {
				return
			}
			if expectRefunded, err = s.dao.CountRefundedOrderUser(ctx, beginTime, endTime); err != nil {
				return
			}
			expect = expectPaid + expectRefunded
			return
		}
	)
	if finished, err = checkOrCreateTaskFromLog(ctx, s, s.tl, expectFN); err != nil || finished {
		return
	}
	return s.run(ctx)
}

func (s *taskBillDaily) TTL() int32 {
	return 3600 * 2
}

func (s *taskBillDaily) Name() string {
	return fmt.Sprintf("%s_%d", s.namePrefix, dailyBillVer(time.Now()))
}

// 日账单生成
func (s *taskBillDaily) run(ctx context.Context) (err error) {
	// 已支付成功订单入账
	paidLL := &orderPaidLL{
		limit: 1000,
		dao:   s.dao,
	}
	paidLL.beginTime, paidLL.endTime = dayRange(s.dayOffset)
	if err = runLimitedList(ctx, paidLL, time.Millisecond*5, s.runPaidOrder); err != nil {
		return
	}

	// 已退款订单入账
	refundLL := &orderRefundedLL{
		limit: 1000,
		dao:   s.dao,
	}
	refundLL.beginTime, refundLL.endTime = dayRange(s.dayOffset)
	return runLimitedList(ctx, refundLL, time.Millisecond*5, s.runRefundedOrder)
}

// 处理退款order
func (s *taskBillDaily) runRefundedOrder(ctx context.Context, ele interface{}) (err error) {
	order, ok := ele.(*model.Order)
	if !ok {
		return errors.Errorf("refundedOrderHandler convert ele: %+v failed", order)
	}
	log.Info("runRefundedOrder handle order: %+v", order)

	logOrder := &model.LogOrder{
		OrderID:   order.OrderID,
		FromState: order.State,
		ToState:   model.OrderStateRefundFinished,
	}
	order.State = model.OrderStateRefundFinished

	fn := func(ctx context.Context, tx *xsql.Tx) (affected bool, err error) {
		var (
			bill          *model.DailyBill
			asset         *model.Asset
			ver           = dailyBillVer(order.RefundTime)
			monthVer      = monthlyBillVer(order.RefundTime)
			billDailyLog  *model.LogBillDaily
			userIncome, _ = calcAssetIncome(order.RealFee) // 收入计算结果
		)
		affected = true

		// 获得订单对应的asset
		if asset, err = s.dao.Asset(ctx, order.OID, order.OType, order.Currency); err != nil {
			return
		}
		if asset == nil {
			err = errors.Errorf("dailyBillHander find invalid asset order, order: %+v", order)
			return
		}

		// 获得该mid对应的日账单
		if bill, err = s.dao.DailyBill(ctx, asset.MID, model.BizAsset, model.CurrencyBP, ver); err != nil {
			return
		}
		if bill == nil {
			if bill, err = s.initDailyBill(ctx, asset.MID, model.BizAsset, model.CurrencyBP, ver, monthVer); err != nil {
				return
			}
		}
		// 计算日账单
		billDailyLog = &model.LogBillDaily{
			BillID:  bill.BillID,
			FromIn:  bill.In,
			ToIn:    bill.In,
			FromOut: bill.Out + userIncome,
			ToOut:   bill.Out,
			OrderID: order.OrderID + "_r",
		}
		bill.Out += userIncome

		// 更新order
		rowAffected, err := s.dao.TXUpdateOrder(ctx, tx, order)
		if err != nil {
			tx.Rollback()
			return
		}
		if rowAffected <= 0 {
			tx.Rollback()
			log.Error("UpdateOrder no affected from order: %+v", order)
			affected = false
			return
		}

		// 添加 order log
		_, err = s.dao.TXInsertOrderUserLog(ctx, tx, logOrder)
		if err != nil {
			tx.Rollback()
			return
		}

		// 更新daily_bill
		rowAffected, err = s.dao.TXUpdateDailyBill(ctx, tx, bill)
		if err != nil {
			tx.Rollback()
			return
		}
		if rowAffected <= 0 {
			log.Error("TXUpdateDailyBill no affected bill: %+v", bill)
			tx.Rollback()
			affected = false
			return
		}

		// 添加 daily bill log , uk order_id
		_, err = s.dao.TXInsertLogDailyBill(ctx, tx, billDailyLog)
		if err != nil {
			tx.Rollback()
			return
		}

		// 更新 aggr
		aggrMonthlyAsset := &model.AggrIncomeUserAsset{
			MID:      bill.MID,
			Currency: bill.Currency,
			Ver:      monthVer,
			OID:      order.OID,
			OType:    order.OType,
		}
		aggrAllAsset := &model.AggrIncomeUserAsset{
			MID:      bill.MID,
			Currency: bill.Currency,
			Ver:      0,
			OID:      order.OID,
			OType:    order.OType,
		}
		aggrUser := &model.AggrIncomeUser{
			MID:      bill.MID,
			Currency: bill.Currency,
		}
		_, err = s.dao.TXUpsertDeltaAggrIncomeUserAsset(ctx, tx, aggrAllAsset, 0, 1, 0, userIncome)
		if err != nil {
			tx.Rollback()
			return
		}
		_, err = s.dao.TXUpsertDeltaAggrIncomeUserAsset(ctx, tx, aggrMonthlyAsset, 0, 1, 0, userIncome)
		if err != nil {
			tx.Rollback()
			return
		}
		_, err = s.dao.TXUpsertDeltaAggrIncomeUser(ctx, tx, aggrUser, 0, 1, 0, userIncome)
		if err != nil {
			tx.Rollback()
			return
		}

		log.Info("Settle daily bill: %+v, aggrAllAsset: %+v, aggrMonthlyAsset: %+v, aggrUser: %+v, from refunded order: %+v", bill, aggrAllAsset, aggrMonthlyAsset, aggrUser, order)
		return
	}

	return runTXCASTaskWithLog(ctx, s, s.tl, fn)
}

func (s *taskBillDaily) runPaidOrder(ctx context.Context, ele interface{}) (err error) {
	order, ok := ele.(*model.Order)
	if !ok {
		return errors.Errorf("runPaidOrder convert ele: %+v failed", order)
	}
	log.Info("runPaidOrder handle order: %+v", order)

	checkOK, payDesc, err := s.checkOrder(ctx, order) // 对支付订单对账
	if err != nil {
		return err
	}
	logOrder := &model.LogOrder{
		OrderID:   order.OrderID,
		FromState: order.State,
		Desc:      payDesc,
	}
	var fn func(context.Context, *xsql.Tx) (affected bool, err error)

	if checkOK { // 对账成功
		logOrder.ToState = model.OrderStateSettled
		order.State = model.OrderStateSettled
		fn = func(ctx context.Context, tx *xsql.Tx) (affected bool, err error) {
			var (
				bill          *model.DailyBill
				asset         *model.Asset
				ver           = dailyBillVer(order.PayTime)
				monthVer      = monthlyBillVer(order.PayTime)
				billDailyLog  *model.LogBillDaily
				userIncome, _ = calcAssetIncome(order.RealFee) // 收入计算结果
			)
			affected = true

			// 获得订单对应的asset
			if asset, err = s.dao.Asset(ctx, order.OID, order.OType, order.Currency); err != nil {
				return
			}
			if asset == nil {
				log.Error("runPaidOrder find invalid asset order, order: %+v", order)
				return
			}

			// 获得该mid对应的日账单
			if bill, err = s.dao.DailyBill(ctx, asset.MID, model.BizAsset, model.CurrencyBP, ver); err != nil {
				return
			}
			if bill == nil {
				if bill, err = s.initDailyBill(ctx, asset.MID, model.BizAsset, model.CurrencyBP, ver, monthVer); err != nil {
					return
				}
			}
			// 计算日账单
			billDailyLog = &model.LogBillDaily{
				BillID:  bill.BillID,
				FromIn:  bill.In,
				ToIn:    bill.In + userIncome,
				FromOut: bill.Out,
				ToOut:   bill.Out,
				OrderID: order.OrderID,
			}
			bill.In += userIncome

			// 更新order
			rowAffected, err := s.dao.TXUpdateOrder(ctx, tx, order)
			if err != nil {
				tx.Rollback()
				return
			}
			if rowAffected <= 0 {
				tx.Rollback()
				log.Error("UpdateOrder no affected from order: %+v", order)
				affected = false
				return
			}

			// 添加 order log
			_, err = s.dao.TXInsertOrderUserLog(ctx, tx, logOrder)
			if err != nil {
				tx.Rollback()
				return
			}

			// 更新daily_bill
			rowAffected, err = s.dao.TXUpdateDailyBill(ctx, tx, bill)
			if err != nil {
				tx.Rollback()
				return
			}
			if rowAffected <= 0 {
				log.Error("TXUpsertDeltaDailyBill no affected bill: %+v", bill)
				tx.Rollback()
				affected = false
				return
			}

			// 添加 daily bill log , uk order_id
			_, err = s.dao.TXInsertLogDailyBill(ctx, tx, billDailyLog)
			if err != nil {
				tx.Rollback()
				return
			}

			// 更新 aggr
			aggrMonthlyAsset := &model.AggrIncomeUserAsset{
				MID:      bill.MID,
				Currency: bill.Currency,
				Ver:      monthVer,
				OID:      order.OID,
				OType:    order.OType,
			}
			aggrAllAsset := &model.AggrIncomeUserAsset{
				MID:      bill.MID,
				Currency: bill.Currency,
				Ver:      0,
				OID:      order.OID,
				OType:    order.OType,
			}
			aggrUser := &model.AggrIncomeUser{
				MID:      bill.MID,
				Currency: bill.Currency,
			}
			_, err = s.dao.TXUpsertDeltaAggrIncomeUserAsset(ctx, tx, aggrAllAsset, 1, 0, userIncome, 0)
			if err != nil {
				tx.Rollback()
				return
			}
			_, err = s.dao.TXUpsertDeltaAggrIncomeUserAsset(ctx, tx, aggrMonthlyAsset, 1, 0, userIncome, 0)
			if err != nil {
				tx.Rollback()
				return
			}
			_, err = s.dao.TXUpsertDeltaAggrIncomeUser(ctx, tx, aggrUser, 1, 0, userIncome, 0)
			if err != nil {
				tx.Rollback()
				return
			}
			log.Info("taskBillDaily: %+v, aggrAllAsset: %+v, aggrMonthlyAsset: %+v, aggrUser: %+v, from paid order: %+v", bill, aggrAllAsset, aggrMonthlyAsset, aggrUser, order)
			return
		}
	} else { // 对账失败
		logOrder.ToState = model.OrderStateBadDebt
		order.State = model.OrderStateBadDebt
		fn = func(ctx context.Context, tx *xsql.Tx) (affected bool, err error) {
			var (
				orderBadDebt *model.OrderBadDebt
			)
			affected = true

			if orderBadDebt, err = s.dao.OrderBadDebt(ctx, order.OrderID); err != nil {
				return
			}
			if orderBadDebt == nil {
				if orderBadDebt, err = s.initBadDebt(ctx, order.OrderID); err != nil {
					return
				}
			}
			orderBadDebt.Type = "unknown"
			orderBadDebt.State = "failed"

			// 更新order
			rowAffected, theErr := s.dao.TXUpdateOrder(ctx, tx, order)
			if theErr != nil {
				tx.Rollback()
				return
			}
			if rowAffected <= 0 {
				tx.Rollback()
				log.Error("UpdateOrder no affected from order: %+v", order)
				affected = false
				return
			}

			// 添加order log
			_, theErr = s.dao.TXInsertOrderUserLog(ctx, tx, logOrder)
			if theErr != nil {
				tx.Rollback()
				return
			}

			// 添加坏账表
			_, err = s.dao.TXUpdateOrderBadDebt(ctx, tx, orderBadDebt)
			if err != nil {
				tx.Rollback()
				return
			}
			log.Info("Add bad debt: %+v", orderBadDebt)
			return
		}
	}
	return runTXCASTaskWithLog(ctx, s, s.tl, fn)
}

func (s *taskBillDaily) checkOrder(ctx context.Context, order *model.Order) (ok bool, payDesc string, err error) {
	ok = false
	if order == nil {
		return
	}
	if order.PayID == "" {
		log.Error("Check order found baddebt order: %+v", order)
		return
	}

	payParam := s.pay.CheckOrder(order.PayID)
	s.pay.Sign(payParam)
	payJSON, err := s.pay.ToJSON(payParam)
	if err != nil {
		return
	}
	orders, err := s.dao.PayCheckOrder(ctx, payJSON)
	if err != nil {
		return
	}
	result, ok := orders[order.PayID]
	if !ok {
		return
	}
	payDesc = result.RecoStatusDesc
	switch result.RecoStatusDesc {
	case model.PayCheckOrderStateSuccess:
		ok = true
	default:
		ok = false
	}
	return
}

func (s *taskBillDaily) initDailyBill(ctx context.Context, mid int64, biz, currency string, ver, monthVer int64) (bill *model.DailyBill, err error) {
	bill = &model.DailyBill{}
	bill.BillID = orderID(s.rnd)
	bill.MID = mid
	bill.Biz = model.BizAsset
	bill.Currency = model.CurrencyBP
	bill.In = 0
	bill.Out = 0
	bill.Ver = ver
	bill.MonthVer = monthVer
	bill.Version = 1

	if bill.ID, err = s.dao.InsertDailyBill(ctx, bill); err != nil {
		return
	}
	return
}

func (s *taskBillDaily) initBadDebt(ctx context.Context, orderID string) (data *model.OrderBadDebt, err error) {
	data = &model.OrderBadDebt{
		OrderID: orderID,
		Type:    "",
		State:   "",
	}
	if data.ID, err = s.dao.InsertOrderBadDebt(ctx, data); err != nil {
		return
	}
	return
}
