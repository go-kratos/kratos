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

// 结算贝壳任务
type taskRechargeShell struct {
	dao         *dao.Dao
	pay         *pay.Pay
	rnd         *rand.Rand
	monthOffset int
	namePrefix  string
	tl          *taskLog
}

func (s *taskRechargeShell) Run() (err error) {
	var (
		ctx      = context.Background()
		finished bool
		expectFN = func(ctx context.Context) (expect int64, err error) {
			var (
				beginTime, _ = monthRange(s.monthOffset)
				monthVer     = monthlyBillVer(beginTime)
			)
			if expect, err = s.dao.CountMonthlyBillByVer(ctx, monthVer); err != nil {
				return
			}
			return
		}
	)
	if finished, err = checkOrCreateTaskFromLog(ctx, s, s.tl, expectFN); err != nil || finished {
		return
	}
	return s.run(ctx)
}

func (s *taskRechargeShell) TTL() int32 {
	return 3600 * 2
}

func (s *taskRechargeShell) Name() string {
	return fmt.Sprintf("%s_%d", s.namePrefix, monthlyBillVer(time.Now()))
}

func (s *taskRechargeShell) run(ctx context.Context) (err error) {
	ll := &monthlyBillLL{
		limit: 1000,
		dao:   s.dao,
	}
	beginTime, _ := monthRange(s.monthOffset)
	ll.ver = monthlyBillVer(beginTime)
	return runLimitedList(ctx, ll, time.Millisecond*2, s.runMonthlyBill)
}

func (s *taskRechargeShell) runMonthlyBill(ctx context.Context, ele interface{}) (err error) {
	monthlyBill, ok := ele.(*model.Bill)
	if !ok {
		return errors.Errorf("runMonthlyBill convert ele: %+v failed", monthlyBill)
	}
	log.Info("taskRechargeShell start handle monthly bill: %+v", monthlyBill)

	var fn func(ctx context.Context, tx *xsql.Tx) (affected bool, err error)

	if monthlyBill.In-monthlyBill.Out > 0 {
		fn = s.fnRechargeShell(ctx, monthlyBill)
	} else {
		fn = s.fnRecordRecharge(ctx, monthlyBill)
	}
	if err = runTXCASTaskWithLog(ctx, s, s.tl, fn); err != nil {
		return
	}
	return
}

func (s *taskRechargeShell) fnRecordRecharge(ctx context.Context, monthlyBill *model.Bill) func(ctx context.Context, tx *xsql.Tx) (affected bool, err error) {
	return func(ctx context.Context, tx *xsql.Tx) (affected bool, err error) {
		affected = true
		var (
			readyRecharge = monthlyBill.In - monthlyBill.Out
		)

		// 记录贝壳订单
		orderRechargeShell := &model.OrderRechargeShell{
			MID:     monthlyBill.MID,
			OrderID: orderID(s.rnd),
			Biz:     model.BizAsset,
			Amount:  readyRecharge,
			State:   "finished",
			Ver:     monthlyBill.Ver,
		}
		var (
			orderRechargeShellLog = &model.OrderRechargeShellLog{
				OrderID:           orderRechargeShell.OrderID,
				FromState:         "finished",
				ToState:           "finished",
				BillUserMonthlyID: monthlyBill.BillID,
			}
		)
		_, err = s.dao.TXInsertOrderRechargeShell(ctx, tx, orderRechargeShell)
		if err != nil {
			tx.Rollback()
			return
		}

		// 插入 order_recharge_shell_log, uk: bill_monthly_bill_id
		_, err = s.dao.TXInsertOrderRechargeShellLog(ctx, tx, orderRechargeShellLog)
		if err != nil {
			tx.Rollback()
			return
		}

		log.Info("fnRecordRecharge : %+v, orderRechargeShell: %+v", monthlyBill, orderRechargeShell)
		return
	}
}

func (s *taskRechargeShell) fnRechargeShell(ctx context.Context, monthlyBill *model.Bill) func(ctx context.Context, tx *xsql.Tx) (affected bool, err error) {
	return func(ctx context.Context, tx *xsql.Tx) (affected bool, err error) {
		affected = true
		var (
			account       *model.UserAccount
			readyRecharge = monthlyBill.In - monthlyBill.Out
			accountLog    *model.AccountLog
		)

		// 获得该 mid 的 虚拟账户
		if account, err = s.dao.UserAccount(ctx, monthlyBill.MID, model.BizAsset, model.CurrencyBP); err != nil {
			return
		}
		if account == nil {
			err = errors.Errorf("runMonthlyBill not found valid user_account, monthly_bill: %+v", monthlyBill)
			return
		}
		// 检查虚拟账户余额是否足够
		if account.Balance-readyRecharge < 0 {
			err = errors.Errorf("runMonthlyBill failed, account.Balance - readyRecharge < 0 !!!! account: %+v, monthly bill: %+v", account, monthlyBill)
			return
		}
		accountLog = &model.AccountLog{
			AccountID: account.ID,
			Name:      s.Name(),
			From:      account.Balance,
			To:        account.Balance - readyRecharge,
			Ver:       account.Ver + 1,
			State:     model.AccountStateWithdraw,
		}
		account.Balance -= readyRecharge

		// 扣减虚拟账户余额
		rowAffected, err := s.dao.TXUpdateUserAccount(ctx, tx, account)
		if err != nil {
			tx.Rollback()
			return
		}
		if rowAffected <= 0 {
			log.Error("TXUpdateUserAccount no affected user account: %+v", account)
			affected = false
			tx.Rollback()
			return
		}
		err = s.dao.TXInsertUserAccountLog(ctx, tx, accountLog)
		if err != nil {
			tx.Rollback()
			return
		}

		// 开始转贝壳
		orderRechargeShell := &model.OrderRechargeShell{
			MID:     monthlyBill.MID,
			OrderID: orderID(s.rnd),
			Biz:     model.BizAsset,
			Amount:  readyRecharge,
			State:   "created",
			Ver:     monthlyBill.Ver,
		}
		var (
			orderRechargeShellLog = &model.OrderRechargeShellLog{
				OrderID:           orderRechargeShell.OrderID,
				FromState:         "created",
				ToState:           "created",
				BillUserMonthlyID: monthlyBill.BillID,
			}
		)
		_, err = s.dao.TXInsertOrderRechargeShell(ctx, tx, orderRechargeShell)
		if err != nil {
			tx.Rollback()
			return
		}

		// 插入 order_recharge_shell_log, uk: bill_monthly_bill_id
		_, err = s.dao.TXInsertOrderRechargeShellLog(ctx, tx, orderRechargeShellLog)
		if err != nil {
			tx.Rollback()
			return
		}

		// 请求支付中心转贝壳
		_, payJSON, err := s.pay.RechargeShell(orderRechargeShell.OrderID, orderRechargeShell.MID, orderRechargeShell.Amount, orderRechargeShell.Amount)
		if err != nil {
			tx.Rollback()
			return
		}
		if err = s.dao.PayRechargeShell(ctx, payJSON); err != nil {
			tx.Rollback()
			return
		}

		log.Info("fnRechargeShell: %+v, account: %+v, orderRechargeShell: %+v", monthlyBill, account, orderRechargeShell)
		return
	}
}
