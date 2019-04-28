package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/growup/model"

	"go-common/library/log"
)

var (
	_withdrawing     = 1 // 处理中
	_withdrawSuccess = 2 // 提现成功
	// _withdrawFail    = 3 // 提现失败
)

// GetWithdraw get up withdraw
func (s *Service) GetWithdraw(c context.Context, dateVersion string, from, limit int) (count int, withdrawVos []*model.WithdrawVo, err error) {
	count, upAccounts, err := s.UpWithdraw(c, dateVersion, from, limit)
	if err != nil {
		log.Error("s.UpWithdraw error(%v)", err)
		return
	}

	mids := make([]int64, len(upAccounts))
	for i, up := range upAccounts {
		mids[i] = up.MID
	}

	withdrawVos = make([]*model.WithdrawVo, 0)
	if len(mids) == 0 {
		return
	}

	upIncomeWithdrawMap, err := s.dao.QueryUpWithdrawByMids(c, mids, dateVersion)
	if err != nil {
		log.Error("s.dao.QueryUpWithdrawByMids error(%v)", err)
		return
	}

	for _, up := range upAccounts {
		if upIncomeWithdraw, ok := upIncomeWithdrawMap[up.MID]; ok && upIncomeWithdraw.State == _withdrawing {
			vo := &model.WithdrawVo{
				MID:          up.MID,
				ThirdCoin:    float64(up.TotalUnwithdrawIncome) * float64(0.01),
				ThirdOrderNo: strconv.FormatInt(upIncomeWithdraw.ID, 10),
				CTime:        time.Unix(int64(upIncomeWithdraw.CTime), 0).Format("2006-01-02 15:04:05"),
				NotifyURL:    "http://up-profit.bilibili.co/allowance/api/x/internal/growup/up/withdraw/success",
			}

			withdrawVos = append(withdrawVos, vo)
		}
	}

	return
}

// UpWithdraw get up withdraw
func (s *Service) UpWithdraw(c context.Context, dateVersion string, from, limit int) (count int, upAccounts []*model.UpAccount, err error) {
	count, err = s.dao.GetUpAccountCount(c, dateVersion)
	if err != nil {
		log.Error("s.dao.GetUpAccountCount error(%v)", err)
		return
	}
	if count <= 0 {
		return
	}

	upAccounts, err = s.dao.QueryUpAccountByDate(c, dateVersion, from, limit)
	if err != nil {
		log.Error("s.dao.QueryUpAccountByDate error(%v)", err)
		return
	}
	if len(upAccounts) == 0 {
		return
	}

	mids := make([]int64, len(upAccounts))
	for i, up := range upAccounts {
		mids[i] = up.MID
	}

	// get up_income_withdraw by mids and date
	upIncomeWithdrawMap, err := s.dao.QueryUpWithdrawByMids(c, mids, dateVersion)
	if err != nil {
		log.Error("s.dao.QueryUpWithdrawByMids error(%v)", err)
		return
	}

	for _, up := range upAccounts {
		if _, ok := upIncomeWithdrawMap[up.MID]; !ok {
			upIncomeWithdraw := &model.UpIncomeWithdraw{
				MID:            up.MID,
				WithdrawIncome: up.TotalUnwithdrawIncome,
				DateVersion:    dateVersion,
				State:          _withdrawing,
			}

			err = s.InsertUpWithdrawRecord(c, upIncomeWithdraw)
			if err != nil {
				log.Error("s.InsertUpWithdrawRecord error(%v)", err)
				return
			}
		}
	}
	return
}

// InsertUpWithdrawRecord insert up_withdraw_income record
func (s *Service) InsertUpWithdrawRecord(c context.Context, upIncomeWithdraw *model.UpIncomeWithdraw) (err error) {
	result, err := s.dao.InsertUpWithdrawRecord(c, upIncomeWithdraw)
	if err != nil {
		log.Error("s.dao.InsertUpIncomeWithdraw error(%v)", err)
		return
	}
	if result < 1 {
		log.Error("s.dao.InsertUpIncomeWithdraw error mid(%d), dateVersion(%s)", upIncomeWithdraw.MID, upIncomeWithdraw.DateVersion)
		return
	}
	return
}

// WithdrawSuccess withdraw success callback
func (s *Service) WithdrawSuccess(c context.Context, orderNo int64, tradeStatus int) (err error) {
	upWithdraw, err := s.dao.QueryUpWithdrawByID(c, orderNo)
	if err != nil {
		log.Error("s.dao.QueryUpWithdrawByID error(%v)", err)
		return
	}

	if tradeStatus != _withdrawSuccess {
		log.Info("param tradeStatus(%d) != withdraw success(2)", tradeStatus)
		return
	}

	if upWithdraw.State == _withdrawSuccess {
		log.Info("withdraw has successed already")
		return
	}

	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}

	// update up_income_withdraw state
	rows, err := s.dao.TxUpdateUpWithdrawState(tx, orderNo, _withdrawSuccess)
	if err != nil {
		tx.Rollback()
		log.Error("s.dao.UpdateUpWithdrawState error(%v)", err)
		return
	}
	if rows != 1 {
		tx.Rollback()
		log.Error("s.dao.UpdateUpWithdrawState Update withdraw record error id(%d)", orderNo)
		return
	}

	// update up_account withdraw
	rows, err = s.dao.TxUpdateUpAccountWithdraw(tx, upWithdraw.MID, upWithdraw.WithdrawIncome)
	if err != nil {
		tx.Rollback()
		log.Error("s.dao.UpdateUpAccountWithdraw error(%v)", err)
		return
	}
	if rows != 1 {
		tx.Rollback()
		log.Error("s.dao.UpdateUpAccountWithdraw Update up account record error id(%d)", orderNo)
		return
	}

	maxUpWithdrawDateVersion, err := s.dao.TxQueryMaxUpWithdrawDateVersion(tx, upWithdraw.MID)
	if err != nil {
		tx.Rollback()
		log.Error("s.dao.QueryMaxUpWithdrawDateVersion error(%v)", err)
		return
	}

	time := 0
	var version int64
	for {
		version, err = s.dao.TxQueryUpAccountVersion(tx, upWithdraw.MID)
		if err != nil {
			tx.Rollback()
			log.Error("s.dao.QueryUpAccountVersion error(%v)", err)
			return
		}
		if maxUpWithdrawDateVersion == "" {
			maxUpWithdrawDateVersion = upWithdraw.DateVersion
		}

		rows, err = s.dao.TxUpdateUpAccountUnwithdrawIncome(tx, upWithdraw.MID, maxUpWithdrawDateVersion, version)
		if err != nil {
			tx.Rollback()
			log.Error("s.dao.UpdateUpAccountUnwithdrawIncome error(%v)", err)
			return
		}
		if rows == 1 {
			if err = tx.Commit(); err != nil {
				log.Error("tx.Commit error")
				return err
			}
			break
		}

		time++
		if time >= 10 {
			tx.Rollback()
			log.Info("try to synchronize unwithdraw income 10 times error mid(%d)", upWithdraw.MID)
			err = fmt.Errorf("try to synchronize unwithdraw income 10 times error mid(%d)", upWithdraw.MID)
			break
		}
	}

	return
}

// WithdrawDetail get withdraw detail
func (s *Service) WithdrawDetail(c context.Context, mid int64) (upWithdraws []*model.UpIncomeWithdraw, err error) {
	return s.dao.QueryUpWithdrawByMID(c, mid)
}
