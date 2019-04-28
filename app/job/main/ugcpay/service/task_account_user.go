package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/ugcpay/conf"
	"go-common/app/job/main/ugcpay/dao"
	"go-common/app/job/main/ugcpay/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

type taskAccountUser struct {
	dao        *dao.Dao
	taskPre    TaskProcess
	dayOffset  int
	namePrefix string
	tl         *taskLog
}

func (s *taskAccountUser) Run() (err error) {
	// 检查日账单任务是否完成
	if _, finished := s.tl.checkTask(s.taskPre); !finished {
		log.Info("taskAccountUser check task: %s not finished", s.taskPre.Name())
		return nil
	}
	var (
		ctx      = context.Background()
		finished bool
		expectFN = func(ctx context.Context) (expect int64, err error) {
			var (
				beginTime, _ = dayRange(s.dayOffset)
				ver          = dailyBillVer(beginTime)
			)
			if expect, err = s.dao.CountDailyBillByVer(ctx, ver); err != nil {
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

func (s *taskAccountUser) TTL() int32 {
	return 3600 * 2
}

func (s *taskAccountUser) Name() string {
	return fmt.Sprintf("%s_%d", s.namePrefix, dailyBillVer(time.Now()))
}

func (s *taskAccountUser) run(ctx context.Context) (err error) {
	ll := &dailyBillLLByVer{
		limit: 1000,
		dao:   s.dao,
	}
	beginTime, _ := dayRange(s.dayOffset)
	ll.ver = dailyBillVer(beginTime)
	return runLimitedList(ctx, ll, time.Millisecond*2, s.runDailyBill)
}

func (s *taskAccountUser) runDailyBill(ctx context.Context, ele interface{}) (err error) {
	dailyBill, ok := ele.(*model.DailyBill)
	if !ok {
		err = errors.Errorf("taskAccountUser convert ele: %+v failed", dailyBill)
		return
	}
	log.Info("taskAccountUser handle dailyBill: %+v", dailyBill)

	fn := func(ctx context.Context, tx *xsql.Tx) (affected bool, err error) {
		var (
			account    *model.UserAccount
			accountLog *model.AccountLog
			userProfit = dailyBill.In - dailyBill.Out
		)
		affected = true

		// 获得该 mid 的 account
		if account, err = s.dao.UserAccount(ctx, dailyBill.MID, dailyBill.Biz, dailyBill.Currency); err != nil {
			return
		}
		// 初始化 biz_account
		if account == nil {
			if account, err = s.initUserAccount(ctx, dailyBill.MID, dailyBill.Biz, dailyBill.Currency); err != nil {
				return
			}
		}
		// 虚拟账户平账，低于一定阈值由虚拟账户转出
		if userProfit < 0 {
			bizRefund, userRefund := calcRefundFee(account.Balance, -userProfit, conf.Conf.Biz.AccountUserMin)
			userProfit = -userRefund

			if bizRefund > 0 {
				// 获得 biz_account
				var bizAccount *model.BizAccount
				if bizAccount, err = s.dao.BizAccount(ctx, model.BizAsset, model.CurrencyBP); err != nil {
					return
				}
				// 初始化 biz_account
				if bizAccount == nil {
					if bizAccount, err = initBizAccount(ctx, model.BizAsset, model.CurrencyBP, s.dao); err != nil {
						return
					}
				}
				bizAccountLog := &model.AccountLog{
					AccountID: bizAccount.ID,
					Name:      s.Name(),
					From:      bizAccount.Balance,
					To:        bizAccount.Balance - bizRefund,
					Ver:       bizAccount.Ver + 1,
					State:     model.AccountStateLoss,
				}
				bizAccount.Balance = bizAccount.Balance - bizRefund
				// 更新 biz account
				var rowAffected int64
				if rowAffected, err = s.dao.TXUpdateBizAccount(ctx, tx, bizAccount); err != nil {
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
				if err = s.dao.TXInsertBizAccountLog(ctx, tx, bizAccountLog); err != nil {
					tx.Rollback()
					return
				}
			}
		}

		accountLog = &model.AccountLog{
			AccountID: account.ID,
			Name:      fmt.Sprintf("%s_%d", s.Name(), dailyBill.MID),
			From:      account.Balance,
			To:        account.Balance + userProfit,
			Ver:       account.Ver + 1,
			State:     model.AccountStateIncome,
		}
		account.Balance = account.Balance + userProfit

		// 更新 user account
		rowAffected, err := s.dao.TXUpdateUserAccount(ctx, tx, account)
		if err != nil {
			tx.Rollback()
			return
		}
		if rowAffected <= 0 {
			log.Error("TXUpdateUserAccount no affected user account: %+v", account)
			tx.Rollback()
			affected = false
			return
		}
		// 添加资金池账户 log
		err = s.dao.TXInsertUserAccountLog(ctx, tx, accountLog)
		if err != nil {
			tx.Rollback()
			return
		}
		log.Info("taskAccountUser: %+v ", account)
		return
	}
	return runTXCASTaskWithLog(ctx, s, s.tl, fn)
}

func (s *taskAccountUser) initUserAccount(ctx context.Context, mid int64, biz, currency string) (account *model.UserAccount, err error) {
	account = &model.UserAccount{}
	account.MID = mid
	account.Biz = biz
	account.Currency = currency
	account.State = model.StateValid
	account.Ver = 1

	if account.ID, err = s.dao.InsertUserAccount(ctx, account); err != nil {
		return
	}
	return
}

// 虚拟账户平账
func calcRefundFee(balance int64, loss int64, minBalance int64) (bizRefund int64, userRefund int64) {
	if balance-loss >= minBalance {
		userRefund = loss
		bizRefund = 0
		return
	}
	if balance > minBalance {
		userRefund = balance - minBalance
	}
	bizRefund = loss - userRefund
	return
}
