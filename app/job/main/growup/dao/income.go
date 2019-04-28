package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/growup/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// select
	_avTagRatioSQL       = "SELECT id, tag_id, av_id FROM av_charge_ratio WHERE id > ? ORDER BY id LIMIT ?"
	_avIncomeInfoSQL     = "SELECT av_id, mid, income, date FROM av_income WHERE av_id = ? AND date = ? AND is_deleted = 0"
	_tagAVTotalIncomeSQL = "SELECT total_income, date FROM up_tag_income WHERE tag_id = ? AND av_id = ? AND is_deleted = 0"
	_avIncomeDateSQL     = "SELECT id, av_id, mid, tag_id, income, total_income, date FROM av_income where id > ? LIMIT ?"
	_upAccuntSQL         = "SELECT mid, total_income, total_unwithdraw_income, withdraw_date_version FROM up_account WHERE withdraw_date_version = ? AND ctime < ? AND total_unwithdraw_income > 0 AND is_deleted = 0 LIMIT ?,?"
	_upWithdrawSQL       = "SELECT mid, withdraw_income FROM up_income_withdraw WHERE date_version = ? AND state = 2 LIMIT ?,?"
	_upIncomeSQL         = "SELECT id, mid, av_count, av_income, column_count, column_income, bgm_count, bgm_income, income, tax_money, total_income, date FROM %s WHERE id > ? AND date = ? ORDER BY id LIMIT ?"

	_upTotalIncomeSQL       = "SELECT id, total_income, is_deleted FROM up_account WHERE id > ? ORDER BY id LIMIT ?"
	_upDateIncomeSQL        = "SELECT id, mid, income, total_income, is_deleted FROM up_income WHERE id > ? AND date = ? ORDER BY id LIMIT ?"
	_avDateIncomeSQL        = "SELECT id, av_id, mid, tag_id, income, base_income, total_income,tax_money,upload_time,date,is_deleted FROM av_income WHERE id > ? AND date = ? AND is_deleted = 0 ORDER BY id LIMIT ?"
	_getUpTotalIncomeCntSQL = "SELECT count(*) FROM up_account WHERE total_income > 0 AND is_deleted = 0"
	_avIncomeStatisCount    = "SELECT count(*) FROM av_income_statis"

	// insert
	_insertTagIncomeSQL = "INSERT INTO up_tag_income(tag_id, mid, av_id, income, total_income, date) VALUES %s ON DUPLICATE KEY UPDATE tag_id = values(tag_id), mid = values(mid), av_id = values(av_id), income = values(income), total_income = values(total_income), date = values(date)"
)

// GetAvTagRatio get av tag info from av_charge_ratio.
func (d *Dao) GetAvTagRatio(c context.Context, from, limit int64) (infos []*model.ActivityAVInfo, err error) {
	rows, err := d.db.Query(c, _avTagRatioSQL, from, limit)
	if err != nil {
		log.Error("dao.GetAvTagRatio query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.ActivityAVInfo{}
		if err = rows.Scan(&a.MID, &a.TagID, &a.AVID); err != nil {
			log.Error("dao.GetAvTagRatio scan error(%v)", err)
			return
		}
		infos = append(infos, a)
	}
	return
}

// GetAvIncomeInfo get av income from av_income.
func (d *Dao) GetAvIncomeInfo(c context.Context, avID int64, date time.Time) (info *model.TagAvIncome, err error) {
	info = new(model.TagAvIncome)
	row := d.db.QueryRow(c, _avIncomeInfoSQL, avID, date)
	err = row.Scan(&info.AVID, &info.MID, &info.Income, &info.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			info = nil
			return
		}
		log.Error("dao.GetAvInfoInfo scan error(%v)", err)
	}
	return
}

// TxInsertTagIncome insert tag_income.
func (d *Dao) TxInsertTagIncome(tx *sql.Tx, sql string) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_insertTagIncomeSQL, sql))
	if err != nil {
		log.Error("dao.TxInsertTagIncome exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// GetTagAvTotalIncome get av total_income from up_tag_income.
func (d *Dao) GetTagAvTotalIncome(c context.Context, tagID, avID int64) (infos []*model.AvIncome, err error) {
	rows, err := d.db.Query(c, _tagAVTotalIncomeSQL, tagID, avID)
	if err != nil {
		log.Error("dao.GetTagAvTotalIncome query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.AvIncome{}
		if err = rows.Scan(&a.TotalIncome, &a.Date); err != nil {
			log.Error("dao.GetTagAvTotalIncome scan error(%v)", err)
			return
		}
		infos = append(infos, a)
	}
	return
}

// ListAvIncome list av income by query
func (d *Dao) ListAvIncome(c context.Context, id int64, limit int) (avIncome []*model.AvIncome, err error) {
	avIncome = make([]*model.AvIncome, 0)
	rows, err := d.db.Query(c, _avIncomeDateSQL, id, limit)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		list := &model.AvIncome{}
		err = rows.Scan(&list.ID, &list.AvID, &list.MID, &list.TagID, &list.Income, &list.TotalIncome, &list.Date)
		if err != nil {
			log.Error("ListAvIncome rows scan error(%v)", err)
			return
		}
		avIncome = append(avIncome, list)
	}

	err = rows.Err()
	return
}

// ListUpAccount list up_acoount by date
func (d *Dao) ListUpAccount(c context.Context, withdrawDate, ctime string, from, limit int) (upAct []*model.UpAccount, err error) {
	upAct = make([]*model.UpAccount, 0)
	rows, err := d.db.Query(c, _upAccuntSQL, withdrawDate, ctime, from, limit)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		list := &model.UpAccount{}
		err = rows.Scan(&list.MID, &list.TotalIncome, &list.TotalUnwithdrawIncome, &list.WithdrawDateVersion)
		if err != nil {
			log.Error("ListUpAccount rows scan error(%v)", err)
			return
		}
		upAct = append(upAct, list)
	}
	err = rows.Err()
	return
}

// ListUpIncome list up_income_? by date
func (d *Dao) ListUpIncome(c context.Context, table, date string, id int64, limit int) (um []*model.UpIncome, err error) {
	um = make([]*model.UpIncome, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_upIncomeSQL, table), id, date, limit)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		list := &model.UpIncome{}
		err = rows.Scan(&list.ID, &list.MID, &list.AvCount, &list.AvIncome, &list.ColumnCount, &list.ColumnIncome, &list.BgmCount, &list.BgmIncome, &list.Income, &list.TaxMoney, &list.TotalIncome, &list.Date)
		if err != nil {
			log.Error("ListUpIncome rows scan error(%v)", err)
			return
		}
		um = append(um, list)
	}
	err = rows.Err()
	return
}

// ListUpWithdraw list up_withdraw_income by date
func (d *Dao) ListUpWithdraw(c context.Context, date string, from, limit int) (ups map[int64]int64, err error) {
	ups = make(map[int64]int64)
	rows, err := d.db.Query(c, _upWithdrawSQL, date, from, limit)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid, income int64
		err = rows.Scan(&mid, &income)
		if err != nil {
			log.Error("ListUpWithdraw rows scan error(%v)", err)
			return
		}
		ups[mid] = income
	}
	err = rows.Err()
	return
}

// GetUpTotalIncome get up totalIncome.
func (d *Dao) GetUpTotalIncome(c context.Context, from, limit int64) (infos []*model.MIDInfo, err error) {
	infos = make([]*model.MIDInfo, 0)
	rows, err := d.db.Query(c, _upTotalIncomeSQL, from, limit)
	if err != nil {
		log.Error("dao.GetUpTotalIncome query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		info := &model.MIDInfo{}
		if err = rows.Scan(&info.ID, &info.TotalIncome, &info.IsDeleted); err != nil {
			log.Error("dao.GetUpTotalIncome scan error(%v)", err)
			return
		}
		infos = append(infos, info)
	}
	return
}

// GetUpIncome get up date income.
func (d *Dao) GetUpIncome(c context.Context, date time.Time, from, limit int64) (infos []*model.MIDInfo, err error) {
	infos = make([]*model.MIDInfo, 0)
	rows, err := d.db.Query(c, _upDateIncomeSQL, from, date, limit)
	if err != nil {
		log.Error("dao.GetUpIncome query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		info := &model.MIDInfo{}
		if err = rows.Scan(&info.ID, &info.MID, &info.Income, &info.TotalIncome, &info.IsDeleted); err != nil {
			log.Error("dao.GetUpIncome scan error(%v)", err)
			return
		}
		infos = append(infos, info)
	}
	return
}

// GetAvIncome get av income info from av_income.
func (d *Dao) GetAvIncome(c context.Context, date time.Time, id, limit int64) (infos []*model.IncomeInfo, err error) {
	infos = make([]*model.IncomeInfo, 0)
	rows, err := d.db.Query(c, _avDateIncomeSQL, id, date, limit)
	if err != nil {
		log.Error("dao.GetIncome query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		info := &model.IncomeInfo{}
		if err = rows.Scan(&info.ID, &info.AVID, &info.MID, &info.TagID, &info.Income, &info.BaseIncome, &info.TotalIncome, &info.TaxMoney, &info.UploadTime, &info.Date, &info.IsDeleted); err != nil {
			log.Error("dao.GetAvIncome scan error(%v)", err)
			return
		}
		infos = append(infos, info)
	}
	return
}

// GetUpTotalIncomeCnt get up t-2 total income > 0 upcnt
func (d *Dao) GetUpTotalIncomeCnt(c context.Context) (upCnt int, err error) {
	row := d.db.QueryRow(c, _getUpTotalIncomeCntSQL)
	if err = row.Scan(&upCnt); err != nil {
		log.Error("growup-job dao.GetUpTotalIncomeCnt scan error(%v)", err)
	}
	return
}

// GetAvStatisCount get av_income_statis count
func (d *Dao) GetAvStatisCount(c context.Context) (cnt int, err error) {
	err = d.db.QueryRow(c, _avIncomeStatisCount).Scan(&cnt)
	return
}
