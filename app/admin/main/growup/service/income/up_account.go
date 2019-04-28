package income

import (
	"context"
	"fmt"

	model "go-common/app/admin/main/growup/model/income"

	"go-common/library/database/sql"
	"go-common/library/log"
)

// GetUpAccount get up_account
func (s *Service) GetUpAccount(c context.Context, query string, isDeleted int) (ups []*model.UpAccount, err error) {
	from, limit := 0, 2000
	for {
		var up []*model.UpAccount
		up, err = s.ListUpAccount(c, query, isDeleted, from, limit)
		if err != nil {
			return
		}
		ups = append(ups, up...)
		if len(up) < limit {
			break
		}
		from += limit
	}
	return
}

// TxUpAccountBreach up_account breach
func (s *Service) TxUpAccountBreach(c context.Context, tx *sql.Tx, mid, preMonthBreach, thisMonthBreach int64) error {
	if mid <= 0 {
		return fmt.Errorf("请输入正确的mid,请确认")
	}
	if preMonthBreach < 0 {
		return fmt.Errorf("未提现金额必须大于等于零,请确认")
	}
	if thisMonthBreach < 0 {
		return fmt.Errorf("当前月未提现金额必须大于等于零,请确认")
	}

	times := 0
	for {
		upAccount, err := s.dao.GetUpAccount(c, mid)
		if err != nil {
			log.Error("s.dao.GetUpAccount error(%v)", err)
			return err
		}

		total := upAccount.TotalIncome - preMonthBreach - thisMonthBreach
		unwithdraw := upAccount.TotalUnwithdrawIncome - preMonthBreach
		if total < 0 {
			log.Info("up_account(%d) total(%d) < 0", mid, total)
			total = 0
		}
		if unwithdraw < 0 {
			log.Info("up_account(%d) total_unwithdraw_income(%d) < 0", mid, unwithdraw)
			unwithdraw = 0
		}

		rows, err := s.dao.TxBreachUpAccount(tx, total, unwithdraw, mid, upAccount.Version+1, upAccount.Version)
		if err != nil {
			tx.Rollback()
			log.Error("s.dao.TxBreachUpAccount error(%v)", err)
			return err
		}
		if rows == 1 {
			break
		}

		times++
		if times >= 5 {
			tx.Rollback()
			return fmt.Errorf("更新up主金额错误")
		}
	}
	return nil
}

// UpAccountCount get up_account count
func (s *Service) UpAccountCount(c context.Context, query string, isDeleted int) (total int64, err error) {
	if query != "" {
		query += " AND"
	}
	return s.dao.UpAccountCount(c, query, isDeleted)
}

// ListUpAccount list up account bu query
func (s *Service) ListUpAccount(c context.Context, query string, isDeleted, from, limit int) (ups []*model.UpAccount, err error) {
	if query != "" {
		query += " AND"
	}
	return s.dao.ListUpAccount(c, query, isDeleted, from, limit)
}
