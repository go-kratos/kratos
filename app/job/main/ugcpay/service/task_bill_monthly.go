package service

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"time"

	"go-common/app/job/main/ugcpay/dao"
	"go-common/app/job/main/ugcpay/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

type taskBillMonthly struct {
	dao         *dao.Dao
	rnd         *rand.Rand
	monthOffset int
	namePrefix  string
	tl          *taskLog
}

func (s *taskBillMonthly) Run() (err error) {
	var (
		ctx      = context.Background()
		finished bool
		expectFN = func(ctx context.Context) (expect int64, err error) {
			var (
				beginTime, _ = monthRange(s.monthOffset)
				monthVer     = monthlyBillVer(beginTime)
			)
			if expect, err = s.dao.CountDailyBillByMonthVer(ctx, monthVer); err != nil {
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

func (s *taskBillMonthly) TTL() int32 {
	return 3600 * 2
}

func (s *taskBillMonthly) Name() string {
	return fmt.Sprintf("%s_%d", s.namePrefix, monthlyBillVer(time.Now()))
}

// 月账单生成
func (s *taskBillMonthly) run(ctx context.Context) (err error) {
	ll := &dailyBillLLByMonthVer{
		limit: 1000,
		dao:   s.dao,
	}
	beginTime, _ := monthRange(s.monthOffset)
	ll.monthVer = monthlyBillVer(beginTime)
	return runLimitedList(ctx, ll, time.Millisecond*2, s.runDailyBill)
}

func (s *taskBillMonthly) runDailyBill(ctx context.Context, ele interface{}) (err error) {
	dailyBill, ok := ele.(*model.DailyBill)
	if !ok {
		return errors.Errorf("taskBillMonthly convert ele: %+v failed", dailyBill)
	}
	log.Info("taskBillMonthly start handle daily biil: %+v", dailyBill)

	fn := func(ctx context.Context, tx *xsql.Tx) (affected bool, err error) {
		var (
			monthlyBill    *model.Bill
			monthVer       = dailyBill.MonthVer
			monthlyBillLog *model.LogBillMonthly
		)
		affected = true

		// 获得该 mid 的 daily_bill
		if monthlyBill, err = s.dao.MonthlyBill(ctx, dailyBill.MID, model.BizAsset, model.CurrencyBP, monthVer); err != nil {
			return
		}
		if monthlyBill == nil {
			if monthlyBill, err = s.initMonthlyBill(ctx, dailyBill.MID, dailyBill.Biz, dailyBill.Currency, dailyBill.MonthVer); err != nil {
				return
			}
		}
		monthlyBillLog = &model.LogBillMonthly{
			BillID:          monthlyBill.BillID,
			FromIn:          monthlyBill.In,
			ToIn:            monthlyBill.In + dailyBill.In,
			FromOut:         monthlyBill.Out,
			ToOut:           monthlyBill.Out + dailyBill.Out,
			BillUserDailyID: dailyBill.BillID,
		}
		monthlyBill.In += dailyBill.In
		monthlyBill.Out += dailyBill.Out

		// 添加 monthly bill log , uk : daily_bill_id
		_, err = s.dao.TXInsertLogMonthlyBill(ctx, tx, monthlyBillLog)
		if err != nil {
			tx.Rollback()
			return
		}

		// 更新 monthly bill
		_, err = s.dao.TXUpdateMonthlyBill(ctx, tx, monthlyBill)
		if err != nil {
			tx.Rollback()
			return
		}
		log.Info("taskBillMonthly: %+v,from daily bill: %+v", monthlyBill, dailyBill)
		return
	}
	return runTXCASTaskWithLog(ctx, s, s.tl, fn)
}

func (s *taskBillMonthly) initMonthlyBill(ctx context.Context, mid int64, biz, currency string, ver int64) (data *model.Bill, err error) {
	data = &model.Bill{
		BillID:   orderID(s.rnd),
		MID:      mid,
		Biz:      biz,
		Currency: currency,
		In:       0,
		Out:      0,
		Ver:      ver,
		Version:  1,
	}
	if data.ID, err = s.dao.InsertMonthlyBill(ctx, data); err != nil {
		return
	}
	return
}

func orderID(rnd *rand.Rand) string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%05d", rnd.Int63n(99999)))
	b.WriteString(fmt.Sprintf("%03d", time.Now().UnixNano()/1e6%1000))
	b.WriteString(time.Now().Format("060102150405"))
	return b.String()
}
