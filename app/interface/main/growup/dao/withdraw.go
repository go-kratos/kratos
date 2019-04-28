package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/interface/main/growup/model"

	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// get up_account withdraw count
	_upAccountCount = "SELECT count(1) FROM up_account where is_deleted = 0 AND has_sign_contract = 1 AND total_unwithdraw_income > 0 AND withdraw_date_version != ?"
	// query up_account by date
	_queryUpAccountByDate = "SELECT mid, total_unwithdraw_income, withdraw_date_version FROM up_account WHERE is_deleted = 0 AND has_sign_contract = 1 AND total_unwithdraw_income > 0 AND withdraw_date_version != ? LIMIT ?,?"
	// query up_account version by mid
	_queryUpAccountVersion = "SELECT version FROM up_account WHERE is_deleted = 0 AND mid = ?"
	// update up_account withdraw
	_updateUpAccountWithdraw = "UPDATE up_account up SET up.total_unwithdraw_income = up.total_unwithdraw_income - %d, up.total_withdraw_income = up.total_withdraw_income + %d, up.last_withdraw_time = ? WHERE up.mid = ? AND up.total_unwithdraw_income > 0 AND is_deleted = 0"
	// update up_account unwithdraw income
	_updateUpAccountUnwithdrawIncome = "UPDATE up_account up SET up.total_unwithdraw_income = up.total_income - up.total_withdraw_income - up.exchange_income, up.withdraw_date_version = ?, up.version = up.version + 1 WHERE up.is_deleted = 0 AND up.mid = ? AND up.version = ?"
	// update up_account exchange and unwithdraw income
	_updateUpAccountExchangeIncome = "UPDATE up_account SET total_unwithdraw_income = total_unwithdraw_income - %d, exchange_income = exchange_income + %d, version = version + 1 WHERE is_deleted = 0 AND mid = ? AND version = ?"

	// query up_income_withdraw by mid
	_queryUpWithdrawByMID = "SELECT withdraw_income, date_version, state, ctime FROM up_income_withdraw WHERE mid = ? AND is_deleted = 0"
	// query up_income_withdraw by mids and date_version
	_queryUpWithdrawByMIDs = "SELECT id, mid, withdraw_income, date_version, state, ctime FROM up_income_withdraw WHERE is_deleted = 0 AND mid in (%s) AND date_version = ?"
	// query up_income_withdraw by id
	_queryUpWithdrawByID = "SELECT id, mid, withdraw_income, date_version, state, ctime FROM up_income_withdraw WHERE is_deleted = 0 AND id = ?"
	// query up_income_withdraw max date_version by mid
	_queryMaxUpWithdrawDateVersion = "SELECT MAX(date_version) FROM up_income_withdraw where is_deleted = 0 AND mid = ?"
	// insert record into up_income_withdraw
	_insertUpWithdrawRecord = "INSERT INTO up_income_withdraw(mid, withdraw_income, date_version, state) VALUES(?,?,?,?)"
	// update up_income_withdraw
	_updateUpWithdrawState = "UPDATE up_income_withdraw up SET up.state = ? WHERE up.id = ? AND is_deleted = 0"
)

// GetUpAccountCount get up account withdraw count
func (d *Dao) GetUpAccountCount(c context.Context, dateVersion string) (count int, err error) {
	row := d.db.QueryRow(c, _upAccountCount, dateVersion)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			count = 0
		} else {
			log.Error("db.QueryRow(%s) error(%v)", _upAccountCount, err)
		}
	}
	return
}

// QueryUpAccountByDate query up_account by date
func (d *Dao) QueryUpAccountByDate(c context.Context, dateVersion string, from, limit int) (upAccounts []*model.UpAccount, err error) {
	upAccounts = make([]*model.UpAccount, 0)
	rows, err := d.db.Query(c, _queryUpAccountByDate, dateVersion, from, limit)
	if err != nil {
		log.Error("d.db.Query(%s) error(%v)", _queryUpAccountByDate, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		up := &model.UpAccount{}
		err = rows.Scan(&up.MID, &up.TotalUnwithdrawIncome, &up.WithdrawDateVersion)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		upAccounts = append(upAccounts, up)
	}

	err = rows.Err()
	return
}

// QueryUpWithdrawByMID query up_income_withdraw by mid
func (d *Dao) QueryUpWithdrawByMID(c context.Context, mid int64) (upWithdraws []*model.UpIncomeWithdraw, err error) {
	upWithdraws = make([]*model.UpIncomeWithdraw, 0)
	rows, err := d.db.Query(c, _queryUpWithdrawByMID, mid)
	if err != nil {
		log.Error("d.db.Query(%s) error(%v)", _queryUpWithdrawByMID, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		upWithdraw := &model.UpIncomeWithdraw{}
		err = rows.Scan(&upWithdraw.WithdrawIncome, &upWithdraw.DateVersion, &upWithdraw.State, &upWithdraw.CTime)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		upWithdraws = append(upWithdraws, upWithdraw)
	}
	err = rows.Err()
	return
}

// QueryUpWithdrawByMids query up_income_withdraw by mids
func (d *Dao) QueryUpWithdrawByMids(c context.Context, mids []int64, dateVersion string) (upWithdraws map[int64]*model.UpIncomeWithdraw, err error) {
	upWithdraws = make(map[int64]*model.UpIncomeWithdraw)
	rows, err := d.db.Query(c, fmt.Sprintf(_queryUpWithdrawByMIDs, xstr.JoinInts(mids)), dateVersion)
	if err != nil {
		log.Error("d.db.Query(%s) error(%v)", _queryUpWithdrawByMIDs, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		upWithdraw := &model.UpIncomeWithdraw{}
		err = rows.Scan(&upWithdraw.ID, &upWithdraw.MID, &upWithdraw.WithdrawIncome, &upWithdraw.DateVersion, &upWithdraw.State, &upWithdraw.CTime)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		upWithdraws[upWithdraw.MID] = upWithdraw
	}
	err = rows.Err()
	return
}

// InsertUpWithdrawRecord insert record into up_income_withdraw
func (d *Dao) InsertUpWithdrawRecord(c context.Context, upWithdraw *model.UpIncomeWithdraw) (result int64, err error) {
	res, err := d.db.Exec(c, _insertUpWithdrawRecord, upWithdraw.MID, upWithdraw.WithdrawIncome, upWithdraw.DateVersion, upWithdraw.State)
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", _insertUpWithdrawRecord, err)
		return
	}
	return res.RowsAffected()
}

// QueryUpWithdrawByID get up_income_withdraw by id
func (d *Dao) QueryUpWithdrawByID(c context.Context, id int64) (upWithdraw *model.UpIncomeWithdraw, err error) {
	upWithdraw = &model.UpIncomeWithdraw{}
	row := d.db.QueryRow(c, _queryUpWithdrawByID, id)
	err = row.Scan(&upWithdraw.ID, &upWithdraw.MID, &upWithdraw.WithdrawIncome, &upWithdraw.DateVersion, &upWithdraw.State, &upWithdraw.CTime)
	return
}

// TxUpdateUpWithdrawState update up_income_withdraw state
func (d *Dao) TxUpdateUpWithdrawState(tx *xsql.Tx, id int64, state int) (result int64, err error) {
	res, err := tx.Exec(_updateUpWithdrawState, state, id)
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", _updateUpWithdrawState)
		return
	}
	return res.RowsAffected()
}

// TxUpdateUpAccountWithdraw update up_account withdraw
func (d *Dao) TxUpdateUpAccountWithdraw(tx *xsql.Tx, mid, thirdCoin int64) (result int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_updateUpAccountWithdraw, thirdCoin, thirdCoin), time.Now(), mid)
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", _updateUpAccountWithdraw)
		return
	}
	return res.RowsAffected()
}

// TxQueryMaxUpWithdrawDateVersion query max date_version from up_income_withdraw by mid
func (d *Dao) TxQueryMaxUpWithdrawDateVersion(tx *xsql.Tx, mid int64) (dateVersion string, err error) {
	row := tx.QueryRow(_queryMaxUpWithdrawDateVersion, mid)
	if err = row.Scan(&dateVersion); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			dateVersion = ""
		} else {
			log.Error("db.QueryRow(%s) error(%v)", _queryMaxUpWithdrawDateVersion, err)
		}
	}
	return
}

// TxQueryUpAccountVersion query up_account version
func (d *Dao) TxQueryUpAccountVersion(tx *xsql.Tx, mid int64) (version int64, err error) {
	row := tx.QueryRow(_queryUpAccountVersion, mid)
	err = row.Scan(&version)
	return
}

// TxUpdateUpAccountUnwithdrawIncome update up_account unwithdraw and version
func (d *Dao) TxUpdateUpAccountUnwithdrawIncome(tx *xsql.Tx, mid int64, dateVersion string, version int64) (result int64, err error) {
	res, err := tx.Exec(_updateUpAccountUnwithdrawIncome, dateVersion, mid, version)
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", _updateUpAccountUnwithdrawIncome)
		return
	}
	return res.RowsAffected()
}

// TxUpdateUpAccountExchangeIncome update up_account unwithdraw and exchange_income
func (d *Dao) TxUpdateUpAccountExchangeIncome(tx *xsql.Tx, mid, income, version int64) (result int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_updateUpAccountExchangeIncome, income, income), mid, version)
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", _updateUpAccountExchangeIncome, err)
		return
	}
	return res.RowsAffected()
}
