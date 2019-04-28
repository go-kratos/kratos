package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/ugcpay/dao"
	"go-common/app/job/main/ugcpay/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

type taskAccountBiz struct {
	dao        *dao.Dao
	taskPre    TaskProcess
	dayOffset  int
	namePrefix string
	tl         *taskLog
}

func (s *taskAccountBiz) Run() (err error) {
	// 检查日账单任务是否完成
	if _, finished := s.tl.checkTask(s.taskPre); !finished {
		log.Info("taskAccountBiz check task: %s not finished", s.taskPre.Name())
		return nil
	}
	var (
		ctx      = context.Background()
		finished bool
		expectFN = func(ctx context.Context) (expect int64, err error) {
			expect = 1
			return
		}
	)
	if finished, err = checkOrCreateTaskFromLog(ctx, s, s.tl, expectFN); err != nil || finished {
		return
	}
	return runTXCASTaskWithLog(ctx, s, s.tl, s.run)
}

func (s *taskAccountBiz) TTL() int32 {
	return 3600 * 2
}

func (s *taskAccountBiz) Name() string {
	return fmt.Sprintf("%s_%d", s.namePrefix, dailyBillVer(time.Now()))
}

func (s *taskAccountBiz) run(ctx context.Context, tx *xsql.Tx) (affected bool, err error) {
	var (
		timeFrom, timeTo        time.Time = dayRange(s.dayOffset)
		ver                               = dailyBillVer(timeFrom)
		bizAccount              *model.BizAccount
		bizAccountLog           *model.AccountLog
		sumPaidOrderRealfee     int64
		sumRefundedOrderRealfee int64
		sumBillDailyIn          int64
		sumBillDailyOut         int64
		bizProfit               int64
	)
	affected = true

	if sumPaidOrderRealfee, err = s.dao.SumPaidOrderUserRealFee(ctx, timeFrom, timeTo); err != nil {
		return
	}
	if sumRefundedOrderRealfee, err = s.dao.SumRefundedOrderUserRealFee(ctx, timeFrom, timeTo); err != nil {
		return
	}
	if sumBillDailyIn, sumBillDailyOut, err = s.dao.SumDailyBill(ctx, ver); err != nil {
		return
	}

	log.Info("taskAccountBiz: %s, sumPaidOrderRealfee: %d, sumRefundedOrderRealfee: %d, sumBillDailyIn: %d, sumBillDailyOut: %d", s.Name(), sumPaidOrderRealfee, sumRefundedOrderRealfee, sumBillDailyIn, sumBillDailyOut)

	if sumPaidOrderRealfee < sumBillDailyIn {
		err = errors.Errorf("taskAccountBiz find sumPaidOrderRealfee(%d) < sumBillDailyIn(%d), ver: %d", sumPaidOrderRealfee, sumBillDailyIn, ver)
		return
	}
	if sumRefundedOrderRealfee < sumBillDailyOut {
		err = errors.Errorf("taskAccountBiz find sumRefundedOrderRealfee(%d) < sumBillDailyOut(%d), ver: %d", sumRefundedOrderRealfee, sumBillDailyOut, ver)
		return
	}

	// 日收益 - 日支出
	bizProfit = (sumPaidOrderRealfee - sumBillDailyIn) - (sumRefundedOrderRealfee - sumBillDailyOut)

	// 获得 biz_account
	if bizAccount, err = s.dao.BizAccount(ctx, model.BizAsset, model.CurrencyBP); err != nil {
		return
	}
	// 初始化 biz_account
	if bizAccount == nil {
		if bizAccount, err = initBizAccount(ctx, model.BizAsset, model.CurrencyBP, s.dao); err != nil {
			return
		}
	}
	bizAccountLog = &model.AccountLog{
		AccountID: bizAccount.ID,
		Name:      s.Name(),
		From:      bizAccount.Balance,
		To:        bizAccount.Balance + bizProfit,
		Ver:       bizAccount.Ver + 1,
		State:     model.AccountStateProfit,
	}
	bizAccount.Balance = bizAccount.Balance + bizProfit

	// 更新 biz account
	rowAffected, err := s.dao.TXUpdateBizAccount(ctx, tx, bizAccount)
	if err != nil {
		tx.Rollback()
		return
	}
	if rowAffected <= 0 {
		log.Error("TXUpdateBizAccount no affected biz account: %+v", bizAccount)
		tx.Rollback()
		affected = false
		return
	}
	// 添加资金池账户 log
	err = s.dao.TXInsertBizAccountLog(ctx, tx, bizAccountLog)
	if err != nil {
		tx.Rollback()
		return
	}
	log.Info("taskAccountBiz: %+v ", bizAccount)
	return
}

func initBizAccount(ctx context.Context, biz, currency string, dao *dao.Dao) (bizAccount *model.BizAccount, err error) {
	bizAccount = &model.BizAccount{
		Biz:      biz,
		Currency: currency,
		State:    model.StateValid,
		Ver:      1,
	}
	if bizAccount.ID, err = dao.InsertBizAccount(ctx, bizAccount); err != nil {
		return
	}
	return
}
