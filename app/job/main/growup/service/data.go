package service

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/job/main/growup/model"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_updateUpAccount = "UPDATE up_account SET withdraw_date_version = '%s', total_unwithdraw_income=total_income - total_withdraw_income - exchange_income where withdraw_date_version = '%s' AND is_deleted = 0"
	_updateTagAdjust = "UPDATE tag_info SET adjust_type = 1 WHERE id = %d AND adjust_type = 0"

	// fix data
	_txUpdateAvIncome = "UPDATE av_income SET total_income = total_income %s, income = income %s  WHERE av_id = %d AND mid = %d AND date = '%s'"

	_txUpdateAvIncomeStatis       = "UPDATE av_income_statis SET total_income = total_income %s WHERE av_id = %d AND mid = %d"
	_txUpdateAvIncomeStatisIncome = "UPDATE av_income_%s_statis SET income = income %s WHERE category_id = %d AND cdate = '%s'"

	_txUpAccount     = "UPDATE up_account SET total_income = total_income %s, total_unwithdraw_income = total_unwithdraw_income %s WHERE mid = %d AND is_deleted = 0"
	_txUpIncomeTable = "UPDATE %s SET av_income = av_income %s, total_income = total_income %s, income = income %s  WHERE mid = %d AND date = '%s'"

	_txUpIncomeStatis      = "UPDATE up_income_statis SET total_income = total_income %s WHERE mid = %d"
	_txUpIncomeDailyStatis = "UPDATE up_income_daily_statis set income = income %s WHERE cdate = '%s'"

	_txTagInsertUpIncomeDate   = "INSERT INTO %s(mid,income,total_income,date) VALUES(%d,%d,%d,'%s') ON DUPLICATE KEY UPDATE total_income=VALUES(total_income),income=VALUES(income)"
	_txTagInsertUpIncomeStatis = "INSERT INTO up_income_statis(mid, total_income) VALUES(%d,%d) ON DUPLICATE KEY UPDATE total_income=VALUES(total_income)"
	_txTagInsertUpAccount      = "INSERT INTO up_account(mid,has_sign_contract,state,total_income,total_unwithdraw_income,withdraw_date_version) VALUES(%d,1,1,%d,%d,'%s') ON DUPLICATE KEY UPDATE total_income=VALUES(total_income),total_unwithdraw_income=VALUES(total_unwithdraw_income)"
	_txTagUpIncomeDailyStatis  = "UPDATE up_income_daily_statis SET ups = ups + %d WHERE cdate = '%s' AND money_section = %d"

	_txUpAvStatis = "INSERT INTO up_av_statis(mid,weekly_date,weekly_av_ids,monthly_date,monthly_av_ids) VALUES(%d,'2018-05-28','%s', '2018-05-01', '%s') ON DUPLICATE KEY UPDATE mid=VALUES(mid)"

	_txAddUpIncomeStatis = "UPDATE up_income_statis SET total_income = total_income + %d WHERE mid = %d"
	_txAddAvIncomeStatis = "UPDATE av_income_statis SET total_income = total_income + %d WHERE av_id = %d"
	_txAddUpIncome       = "UPDATE %s SET total_income = total_income + %d, income = income + %d, av_income = av_income + %d WHERE mid = %d AND date = '%s'"
	_txAddAvIncome       = "UPDATE av_income SET total_income = total_income + %d, income = income + %d WHERE av_id = %d AND date = '%s'"
	_txAddUpAccount      = "UPDATE up_account SET total_income = total_income + %d, total_unwithdraw_income = total_unwithdraw_income + %d WHERE mid = %d"

	_txUpdateAccountType = "UPDATE up_info_video SET account_type=%d WHERE mid=%d"

	_txUpdateUpAccountMoney = "UPDATE up_account SET total_income = total_income + %d, total_unwithdraw_income = total_unwithdraw_income + %d WHERE mid = %d LIMIT 1"
	_txUpdateUpBaseIncome   = "UPDATE up_income SET base_income=%d WHERE mid=%d AND date='%s'"
	_txInUpInfoPGCSQL       = "INSERT INTO up_info_bgm(mid,nickname,fans,account_type,account_state,sign_type) values(%d,%s,%d,%d,%d,%d) ON DUPLICATE KEY UPDATE account_type=VALUES(account_type)"
	_txDelAvBreachSQL       = "DELETE FROM av_breach_record WHERE id = %d LIMIT 1"
	_txUpTotalIncomeSQL     = "UPDATE %s SET total_income = total_income - %d WHERE mid = %d AND date = '%s' LIMIT 1"
	_txColumnTagSQL         = "UPDATE %s SET tag_id = %d WHERE date = '2018-08-19' AND tag_id in (%s) AND inc_charge > 0"
	_txUpdateBgmBaseIncome  = "UPDATE up_income SET bgm_base_income=bgm_income WHERE mid=%d AND date='%s' AND bgm_base_income=0"
	_txDelData              = "DELETE FROM %s LIMIT %d"
)

// DelDataLimit del up_bill
func (s *Service) DelDataLimit(c context.Context, table string, count int64) (err error) {
	if table == "" {
		return
	}
	return s.txUpdateSQL(c, fmt.Sprintf(_txDelData, table, count), count)
}

// FixBgmBaseIncome fix bgm base income
func (s *Service) FixBgmBaseIncome(c context.Context, mid int64, date string) (err error) {
	sql := fmt.Sprintf(_txUpdateBgmBaseIncome, mid, date)
	return s.txUpdateSQL(c, sql, 1)
}

// FixBaseIncome fix income
func (s *Service) FixBaseIncome(c context.Context, base int64, mid int64, date string) (err error) {
	sql := fmt.Sprintf(_txUpdateUpBaseIncome, base, mid, date)
	return s.txUpdateSQL(c, sql, 1)
}

// FixIncome fix income
func (s *Service) FixIncome(c context.Context) (err error) {
	date := time.Date(2018, 9, 10, 0, 0, 0, 0, time.Local)
	total, err := s.getAvIncome(c, date)
	if err != nil {
		log.Error("s.getAvIncome error(%v)", err)
		return
	}
	avIncomes := make([]*model.IncomeInfo, 0)
	for _, av := range total {
		if av.UploadTime.Unix() >= date.Unix() {
			avIncomes = append(avIncomes, av)
		}
	}
	if len(avIncomes) != 1768 {
		err = fmt.Errorf("get av_income(%d) != 1768", len(avIncomes))
		return
	}

	var tx *sql.Tx
	tx, err = s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}

	for _, av := range avIncomes {
		// av_income_statis
		err = s.updateSQL(tx, fmt.Sprintf(_txAddAvIncomeStatis, av.Income, av.AVID), 1)
		if err != nil {
			log.Error("s.UpdateSQL(%s) error(%v)", _txAddAvIncomeStatis, err)
			return
		}
		// up_income_statis
		err = s.updateSQL(tx, fmt.Sprintf(_txAddUpIncomeStatis, av.Income, av.MID), 1)
		if err != nil {
			log.Error("s.UpdateSQL(%s) error(%v)", _txAddUpIncomeStatis, err)
			return
		}

		// av_income
		err = s.updateSQL(tx, fmt.Sprintf(_txAddAvIncome, av.Income, av.Income, av.AVID, "2018-09-10"), 1)
		if err != nil {
			log.Error("s.UpdateSQL(%s) error(%v)", _txAddAvIncome, err)
			return
		}

		// up_income
		err = s.updateSQL(tx, fmt.Sprintf(_txAddUpIncome, "up_income", av.Income, av.Income, av.Income, av.MID, "2018-09-10"), 1)
		if err != nil {
			log.Error("s.UpdateSQL(%s) error(%v)", _txAddUpIncome, err)
			return
		}

		// up_income_weekly
		err = s.updateSQL(tx, fmt.Sprintf(_txAddUpIncome, "up_income_weekly", av.Income, av.Income, av.Income, av.MID, "2018-09-10"), 1)
		if err != nil {
			log.Error("s.UpdateSQL(%s) error(%v)", _txAddUpIncome, err)
			return
		}

		// up_income_monthly
		err = s.updateSQL(tx, fmt.Sprintf(_txAddUpIncome, "up_income_monthly", av.Income, av.Income, av.Income, av.MID, "2018-09-01"), 1)
		if err != nil {
			log.Error("s.UpdateSQL(%s) error(%v)", _txAddUpIncome, err)
			return
		}

		// up_account
		err = s.updateSQL(tx, fmt.Sprintf(_txAddUpAccount, av.Income, av.Income, av.MID), 1)
		if err != nil {
			log.Error("s.UpdateSQL(%s) error(%v)", _txAddUpAccount, err)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error")
	}
	return
}

// FixUpAvStatis fix up_av_statis
func (s *Service) FixUpAvStatis(c context.Context, count int) (err error) {
	upIncome, err := s.GetUpIncome(c, "up_income", "2018-05-31")
	if err != nil {
		log.Error("s.GetUpIncome error(%v)", err)
		return
	}

	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}

	// 301
	addCount := 0
	for _, up := range upIncome {
		if up.TotalIncome == 12700 && up.Income == 12700 {
			err = s.updateSQL(tx, fmt.Sprintf(_txUpAvStatis, up.MID, "", ""), 0)
			if err != nil {
				log.Error("s.UpdateSQL error(%v)", err)
				return
			}
			addCount++
		}
	}

	if count != addCount {
		err = fmt.Errorf("需要添加的record不匹配 %d:%d", count, addCount)
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error")
	}
	return
}

// FixUpIncome fix up_income
func (s *Service) FixUpIncome(c context.Context, date string, tagID int64, addCount, needAddIncome int) (err error) {
	upIncome, err := s.GetUpIncome(c, "up_income", date)
	if err != nil {
		log.Error("s.GetUpIncome error(%v)", err)
		return
	}
	upChargeRatio, err := s.dao.GetUpChargeRatio(c, tagID)
	if err != nil {
		log.Error("s.dao.GetUpChargeRatio error(%v)", err)
		return
	}
	for _, up := range upIncome {
		if _, ok := upChargeRatio[up.MID]; ok {
			delete(upChargeRatio, up.MID)
		}
	}
	if len(upChargeRatio) != addCount {
		err = fmt.Errorf("需要调节的up主数量不匹配 %d:%d", len(upChargeRatio), addCount)
		return
	}
	mids := make([]int64, 0)
	for mid := range upChargeRatio {
		mids = append(mids, mid)
	}
	upIncomeStatis, err := s.dao.GetUpIncomeStatis(c, mids)
	if err != nil {
		log.Error("s.dao.GetUpIncomeStatis error(%v)", err)
		return
	}

	upIncomeWeek, err := s.dao.GetUpIncomeDate(c, mids, "up_income_weekly", "2018-05-28")
	if err != nil {
		log.Error("s.dao.GetUpIncomeDate error(%v)", err)
		return
	}

	upIncomeMonth, err := s.dao.GetUpIncomeDate(c, mids, "up_income_monthly", "2018-05-01")
	if err != nil {
		log.Error("s.dao.GetUpIncomeDate error(%v)", err)
		return
	}

	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}

	// start add
	var totalAddIncome int64
	for mid, ratio := range upChargeRatio {
		totalIncome := upIncomeStatis[mid]
		weekIncome := upIncomeWeek[mid]
		monthIncome := upIncomeMonth[mid]
		// up_income
		err = s.insertIntoSQL(tx, fmt.Sprintf(_txTagInsertUpIncomeDate, "up_income", mid, ratio, totalIncome+ratio, date), 1)
		if err != nil {
			log.Error("s.UpdateSQL error(%v)", err)
			return
		}

		// up_income_weekly todo
		weekDate := "2018-05-28"
		err = s.insertIntoSQL(tx, fmt.Sprintf(_txTagInsertUpIncomeDate, "up_income_weekly", mid, ratio+weekIncome, totalIncome+ratio, weekDate), 0)
		if err != nil {
			log.Error("s.UpdateSQL error(%v)", err)
			return
		}

		// up_income_monthly todo
		monthDate := "2018-05-01"
		err = s.insertIntoSQL(tx, fmt.Sprintf(_txTagInsertUpIncomeDate, "up_income_monthly", mid, ratio+monthIncome, totalIncome+ratio, monthDate), 0)
		if err != nil {
			log.Error("s.UpdateSQL error(%v)", err)
			return
		}

		// up_income_statis
		err = s.insertIntoSQL(tx, fmt.Sprintf(_txTagInsertUpIncomeStatis, mid, ratio+totalIncome), 0)
		if err != nil {
			log.Error("s.insertIntoSQL error(%v)", err)
			return
		}

		// up_account
		err = s.insertIntoSQL(tx, fmt.Sprintf(_txTagInsertUpAccount, mid, ratio+totalIncome, ratio+totalIncome, "2018-04"), 0)
		if err != nil {
			log.Error("s.insertIntoSQL error(%v)", err)
			return
		}
		totalAddIncome += ratio
	}

	// up_income_daily_statis
	if totalAddIncome != int64(needAddIncome) {
		err = fmt.Errorf("需要调节的up主总收入不匹配 %d:%d", totalAddIncome, needAddIncome)
		tx.Rollback()
		return
	}

	add := fmt.Sprintf(" + %d", totalAddIncome)
	err = s.updateSQL(tx, fmt.Sprintf(_txUpIncomeDailyStatis, add, date), 12)
	if err != nil {
		log.Error("s.UpdateSQL error(%v)", err)
		return
	}
	err = s.updateSQL(tx, fmt.Sprintf(_txTagUpIncomeDailyStatis, len(upChargeRatio), date, 3), 1)
	if err != nil {
		log.Error("s.UpdateSQL error(%v)", err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error")
	}
	return
}

func (s *Service) getAllAvRatio(c context.Context, limit int64) (rs map[int64]*model.AvChargeRatio, err error) {
	rs = make(map[int64]*model.AvChargeRatio)
	var id int64
	for {
		var ros map[int64]*model.AvChargeRatio
		ros, id, err = s.dao.AvChargeRatio(c, id, limit)
		if err != nil {
			return
		}
		if len(ros) == 0 {
			break
		}
		for k, v := range ros {
			rs[k] = v
		}
	}
	return
}

// UpdateTagIncome update tag_Info income
func (s *Service) UpdateTagIncome(c context.Context, date string) (err error) {
	avRatio, err := s.getAllAvRatio(c, 2000)
	if err != nil {
		log.Error("s.getAllAvRatio error(%v)", err)
		return
	}
	log.Info("get avratios:%d", len(avRatio))

	updateDate := time.Date(2018, 6, 24, 0, 0, 0, 0, time.Local)
	avIncome, err := s.getAvIncome(c, updateDate)
	if err != nil {
		log.Error("s.getAvIncome error(%v)", err)
		return
	}
	log.Info("get av_income:%d ", len(avIncome))

	err = s.updateIncome(c, avIncome, avRatio)
	if err != nil {
		log.Error("s.updateIncome error(%v)", err)
		return
	}
	return
}

// GetTrueAvsIncome get true av_income
func (s *Service) GetTrueAvsIncome(c context.Context, mids []int64, date string) (avs map[int64]*model.Patch, err error) {
	avs = make(map[int64]*model.Patch)
	for _, mid := range mids {
		var av map[int64]*model.Patch
		av, err = s.AvIncomes(c, mid, date)
		if err != nil {
			return
		}
		for key, val := range av {
			avs[key] = val
		}
	}
	return
}

func (s *Service) updateIncome(c context.Context, avIncome []*model.IncomeInfo, avRatio map[int64]*model.AvChargeRatio) (err error) {
	trueAvs := make([]*model.IncomeInfo, 0)
	uploadTime := time.Date(2018, 6, 24, 0, 0, 0, 0, time.Local)
	for _, av := range avIncome {
		if !uploadTime.After(av.UploadTime) {
			if _, ok := avRatio[av.AVID]; !ok {
				trueAvs = append(trueAvs, av)
			}
		}
	}
	if len(trueAvs) != 1856 {
		err = fmt.Errorf("实际被作用稿件(%d) != 1856", len(trueAvs))
		return
	}

	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}

	for _, av := range trueAvs {
		var incIncome, categoryID int64
		incIncome, categoryID, err = s.dao.AvDailyIncCharge(c, av.AVID)
		if err != nil {
			log.Error("s.dao.AvDailyIncCharge avid(%d) error(%v)", av.MID, err)
			return
		}
		err = s.TxUpdateIncome(tx, av.AVID, av.MID, categoryID, incIncome)
		if err != nil {
			log.Error("ERROR(%v) avid(%d), mid(%d)", err, av.AVID, av.MID)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error")
	}
	return
}

// TxUpdateIncome update creative income
func (s *Service) TxUpdateIncome(tx *sql.Tx, avID, mid, categoryID int64, incIncome int64) (err error) {
	if incIncome <= 0 {
		return
	}
	var incIncomeStr string
	if incIncome > 0 {
		incIncomeStr = fmt.Sprintf("+ %d", incIncome)
	}
	date := "2018-06-24"

	// av_income
	avIncomeSQL := fmt.Sprintf(_txUpdateAvIncome, incIncomeStr, incIncomeStr, avID, mid, date)
	err = s.updateSQL(tx, avIncomeSQL, 1)
	if err != nil {
		log.Error("s.UpdateSQL error(%v)", err)
		return
	}

	// av_income_statis
	avIncomeStatisSQL := fmt.Sprintf(_txUpdateAvIncomeStatis, incIncomeStr, avID, mid)
	err = s.updateSQL(tx, avIncomeStatisSQL, 1)
	if err != nil {
		log.Error("s.UpdateSQL error(%v)", err)
		return
	}
	// av_income_daily_statis
	avIncomeDailyStatisIncome := fmt.Sprintf(_txUpdateAvIncomeStatisIncome, "daily", incIncomeStr, categoryID, date)
	err = s.updateSQL(tx, avIncomeDailyStatisIncome, 12)
	if err != nil {
		log.Error("s.UpdateSQL error(%v)", err)
		return
	}

	// av_income_weekly_statis
	avIncomeWeeklyStatisIncome := fmt.Sprintf(_txUpdateAvIncomeStatisIncome, "weekly", incIncomeStr, categoryID, "2018-06-18")
	err = s.updateSQL(tx, avIncomeWeeklyStatisIncome, 12)
	if err != nil {
		log.Error("s.UpdateSQL error(%v)", err)
		return
	}

	// av_income_monthly_statis
	avIncomeMonthlyStatisIncome := fmt.Sprintf(_txUpdateAvIncomeStatisIncome, "monthly", incIncomeStr, categoryID, "2018-06-01")
	err = s.updateSQL(tx, avIncomeMonthlyStatisIncome, 12)
	if err != nil {
		log.Error("s.UpdateSQL error(%v)", err)
		return
	}

	// up_account
	upAcccountSQL := fmt.Sprintf(_txUpAccount, incIncomeStr, incIncomeStr, mid)
	err = s.updateSQL(tx, upAcccountSQL, 1)
	if err != nil {
		log.Error("s.UpdateSQL error(%v)", err)
		return
	}
	// up_income_statis
	upIncomeStatisSQL := fmt.Sprintf(_txUpIncomeStatis, incIncomeStr, mid)
	err = s.updateSQL(tx, upIncomeStatisSQL, 1)
	if err != nil {
		log.Error("s.UpdateSQL error(%v)", err)
		return
	}

	// up_income_daily_statis
	upIncomeDailyStatisSQL := fmt.Sprintf(_txUpIncomeDailyStatis, incIncomeStr, date)
	err = s.updateSQL(tx, upIncomeDailyStatisSQL, 12)
	if err != nil {
		log.Error("s.UpdateSQL error(%v)", err)
		return
	}

	// up_income
	upIncomeSQL := fmt.Sprintf(_txUpIncomeTable, "up_income", incIncomeStr, incIncomeStr, incIncomeStr, mid, date)
	err = s.updateSQL(tx, upIncomeSQL, 1)
	if err != nil {
		log.Error("s.UpdateSQL error(%v)", err)
		return
	}

	// up_income_weekly
	upIncomeWeeklySQL := fmt.Sprintf(_txUpIncomeTable, "up_income_weekly", incIncomeStr, incIncomeStr, incIncomeStr, mid, "2018-06-18")
	err = s.updateSQL(tx, upIncomeWeeklySQL, 1)
	if err != nil {
		log.Error("s.UpdateSQL error(%v)", err)
		return
	}

	// up_income_monthly
	upIncomeMonthlySQL := fmt.Sprintf(_txUpIncomeTable, "up_income_monthly", incIncomeStr, incIncomeStr, incIncomeStr, mid, "2018-06-01")
	err = s.updateSQL(tx, upIncomeMonthlySQL, 1)
	if err != nil {
		log.Error("s.UpdateSQL error(%v)", err)
		return
	}

	// up_av_statis 不需要修改
	// up_income_withdraw 不需要修改
	return
}

// UpdateWithdraw update up_account withdraw
func (s *Service) UpdateWithdraw(c context.Context, oldDate, newDate string, count int64) (err error) {
	sql := fmt.Sprintf(_updateUpAccount, newDate, oldDate)
	return s.txUpdateSQL(c, sql, count)
}

// UpdateTagAdjust update tag adjust_type
func (s *Service) UpdateTagAdjust(c context.Context, id int64) (err error) {
	sql := fmt.Sprintf(_updateTagAdjust, id)
	return s.txUpdateSQL(c, sql, 1)
}

// UpdateAccountType update account type
func (s *Service) UpdateAccountType(c context.Context, mid int64, accType int) (err error) {
	sql := fmt.Sprintf(_txUpdateAccountType, accType, mid)
	return s.txUpdateSQL(c, sql, 1)
}

// UpdateUpAccountMoney update up_account
func (s *Service) UpdateUpAccountMoney(c context.Context, mid int64, total, unwithdraw int64) (err error) {
	sql := fmt.Sprintf(_txUpdateUpAccountMoney, total, unwithdraw, mid)
	return s.txUpdateSQL(c, sql, 1)
}

// SyncUpPGC sync pgc up from up_info_video to up_info_column
func (s *Service) SyncUpPGC(c context.Context) (err error) {
	ups, err := s.getAllUps(c, 2000)
	if err != nil {
		log.Error("s.getAllUps error(%v)", err)
		return
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}
	for _, up := range ups {
		if up.AccountType == 2 {
			sql := fmt.Sprintf(_txInUpInfoPGCSQL, up.MID, "\""+up.Nickname+"\"", up.Fans, 2, 1, 2)
			err = s.insertIntoSQL(tx, sql, 0)
			if err != nil {
				log.Error("s.UpdateSQL(%s) error(%v)", _txInUpInfoPGCSQL, err)
				return
			}
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error")
	}
	return
}

// FixAvBreach fix av_breach_record data
func (s *Service) FixAvBreach(c context.Context, mid int64, date string, count int) (err error) {
	breachs, err := s.dao.GetAvBreach(c, date, date)
	if err != nil {
		log.Error("s.dao.GetAvBreach error(%v)", err)
		return
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}

	avMap := make(map[int64]bool)
	var delCount int
	for _, b := range breachs {
		if b.MID != mid {
			continue
		}
		if avMap[b.AvID] {
			err = s.updateSQL(tx, fmt.Sprintf(_txDelAvBreachSQL, b.ID), 1)
			if err != nil {
				log.Error("s.UpdateSQL(%s) error(%v)", _txDelAvBreachSQL, err)
				return
			}
			delCount++
		} else {
			avMap[b.AvID] = true
		}
	}

	if count != delCount {
		tx.Rollback()
		log.Error("delete count error %d %d", count, delCount)
		err = fmt.Errorf("delete count error")
		return
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error")
	}
	return
}

// FixUpTotalIncome fix up_income total income
func (s *Service) FixUpTotalIncome(c context.Context, table, date string, count int) (err error) {
	upIncome, err := s.GetUpIncome(c, "up_income", date)
	if err != nil {
		log.Error("s.GetUpIncome error(%v)", err)
		return
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}
	trueCount := 0
	for _, up := range upIncome {
		err = s.updateSQL(tx, fmt.Sprintf(_txUpTotalIncomeSQL, table, up.Income, up.MID, date), 1)
		if err != nil {
			log.Error("s.UpdateSQL(%s) error(%v)", _txUpTotalIncomeSQL, err)
			return
		}
		trueCount++
	}

	if count != trueCount {
		tx.Rollback()
		log.Error("count error %d %d", count, trueCount)
		err = fmt.Errorf(" count error")
		return
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error")
	}
	return
}

// UpdateColumnTag update column tag
func (s *Service) UpdateColumnTag(c context.Context, table string, oldTag string, newTag int, count int64) (err error) {
	sql := fmt.Sprintf(_txColumnTagSQL, table, newTag, oldTag)
	return s.txUpdateSQL(c, sql, count)
}

func (s *Service) txUpdateSQL(c context.Context, sql string, count int64) (err error) {
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}

	err = s.updateSQL(tx, sql, count)
	if err != nil {
		log.Error("s.UpdateSQL(%s) error(%v)", sql, err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error")
	}
	return
}

func (s *Service) updateSQL(tx *sql.Tx, stmt string, count int64) error {
	rows, err := s.dao.UpdateDate(tx, stmt)
	if err != nil {
		tx.Rollback()
		log.Error("s.dao.UpdateDate (%s) error(%v)", stmt, err)
		return err
	}

	if count == 0 && rows <= 1 {
		return nil
	}

	if rows != count {
		tx.Rollback()
		return fmt.Errorf("%s : rows(%d) != count(%d) error", stmt, rows, count)
	}

	return nil
}

func (s *Service) insertIntoSQL(tx *sql.Tx, stmt string, count int) error {
	rows, err := s.dao.UpdateDate(tx, stmt)
	if err != nil {
		tx.Rollback()
		log.Error("s.dao.UpdateDate (%s) error(%v)", stmt, err)
		return err
	}

	if count == 0 {
		if rows > 2 {
			tx.Rollback()
			return fmt.Errorf("rows(%d) error", rows)
		}
	} else if rows != int64(count) {
		tx.Rollback()
		return fmt.Errorf("rows(%d) != count(%d) error", rows, count)
	}

	return nil
}

// SyncAvBaseIncome sync base_income to av_base_income by mid_date
func (s *Service) SyncAvBaseIncome(c context.Context, table string) (err error) {
	data, err := s.avBaseIncomes(c, table)
	if err != nil {
		return
	}
	err = s.batchUpdateUpIncome(c, data, table)
	if err != nil {
		log.Error("batch update av base income error(%v)", err)
	}
	return
}

func (s *Service) avBaseIncomes(c context.Context, table string) (data []*model.AvBaseIncome, err error) {
	var id int64
	for {
		var abs []*model.AvBaseIncome
		abs, id, err = s.dao.GetAvBaseIncome(c, table, id, 2000)
		if err != nil {
			return
		}
		if len(abs) == 0 {
			break
		}
		for _, ab := range abs {
			if ab.AvBaseIncome > 0 {
				data = append(data, ab)
			}
		}
	}
	return
}

func (s *Service) batchUpdateUpIncome(c context.Context, us []*model.AvBaseIncome, table string) (err error) {
	var (
		buff    = make([]*model.AvBaseIncome, 2000)
		buffEnd = 0
	)

	for _, u := range us {
		buff[buffEnd] = u
		buffEnd++

		if buffEnd >= 2000 {
			values := avBaseIncomeValues(buff[:buffEnd])
			buffEnd = 0
			_, err = s.dao.BatchUpdateUpIncome(c, table, values)
			if err != nil {
				return
			}
		}
	}
	if buffEnd > 0 {
		values := avBaseIncomeValues(buff[:buffEnd])
		buffEnd = 0
		_, err = s.dao.BatchUpdateUpIncome(c, table, values)
	}
	return
}

func avBaseIncomeValues(us []*model.AvBaseIncome) (values string) {
	var buf bytes.Buffer
	for _, u := range us {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(u.MID, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + u.Date.Time().Format(_layout) + "'")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.AvBaseIncome, 10))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values = buf.String()
	buf.Reset()
	return
}

// SyncCreditScore sync credit score
func (s *Service) SyncCreditScore(c context.Context) (err error) {
	m := make(map[int64]int)
	am, err := s.getCreditScore(c, "video")
	if err != nil {
		return
	}
	for mid, score := range am {
		m[mid] = score
	}
	cm, err := s.getCreditScore(c, "column")
	if err != nil {
		return
	}
	for mid, score := range cm {
		m[mid] = score
	}
	return s.batchInsertCreditScore(c, m)
}

func (s *Service) batchInsertCreditScore(c context.Context, m map[int64]int) (err error) {
	batch := make(map[int64]int)
	for mid, score := range m {
		batch[mid] = score
		if len(batch) == 2000 {
			values := creditScoreValues(batch)
			_, err = s.dao.SyncCreditScore(c, values)
			if err != nil {
				return
			}
			batch = make(map[int64]int)
		}
	}
	if len(batch) > 0 {
		values := creditScoreValues(batch)
		_, err = s.dao.SyncCreditScore(c, values)
	}
	return
}

func creditScoreValues(m map[int64]int) (values string) {
	var buf bytes.Buffer
	for mid, score := range m {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(mid, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(score))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values = buf.String()
	buf.Reset()
	return
}

func (s *Service) getCreditScore(c context.Context, table string) (m map[int64]int, err error) {
	m = make(map[int64]int)
	var id int64
	for {
		var sm map[int64]int
		sm, id, err = s.dao.GetCreditScore(c, table, id, 2000)
		if err != nil {
			return
		}
		if len(sm) == 0 {
			break
		}
		for k, v := range sm {
			m[k] = v
		}
	}
	return
}

// FixBgmIncomeStatis fix bgm income statis
func (s *Service) FixBgmIncomeStatis(c context.Context) (err error) {
	total, err := s.dao.GetBGMIncome(c)
	if err != nil {
		return
	}
	for sid, income := range total {
		_, err = s.dao.InsertBGMIncomeStatis(c, sid, income)
		if err != nil {
			return
		}
	}
	return
}
