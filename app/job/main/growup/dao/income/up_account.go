package income

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/library/log"

	model "go-common/app/job/main/growup/model/income"
)

const (
	// select
	_upAccountsSQL = "SELECT id,mid,has_sign_contract,state,total_income,total_unwithdraw_income,total_withdraw_income,last_withdraw_time,version,allowance_state,nick_name,withdraw_date_version,is_deleted FROM up_account WHERE id > ? ORDER BY id LIMIT ?"
	_upAccountSQL  = "SELECT total_income,total_unwithdraw_income,version,withdraw_date_version FROM up_account WHERE mid=?"

	// for batch insert
	_insertUpAccountSQL = "INSERT INTO up_account(mid,has_sign_contract,total_income,total_unwithdraw_income,withdraw_date_version,version) VALUES %s ON DUPLICATE KEY UPDATE total_income=VALUES(total_income),total_unwithdraw_income=VALUES(total_unwithdraw_income),version=VALUES(version)"

	_updateUpAccountSQL = "UPDATE up_account SET total_income=?,total_unwithdraw_income=?,version=? WHERE mid=? AND version=?"
)

// UpAccounts batch read up account
func (d *Dao) UpAccounts(c context.Context, id int64, limit int64) (m map[int64]*model.UpAccount, last int64, err error) {
	rows, err := d.db.Query(c, _upAccountsSQL, id, limit)
	if err != nil {
		log.Error("d.db.UpAccounts error(%v)", err)
		return
	}
	m = make(map[int64]*model.UpAccount)
	defer rows.Close()
	for rows.Next() {
		ua := &model.UpAccount{}
		err = rows.Scan(&last, &ua.MID, &ua.HasSignContract, &ua.State, &ua.TotalIncome, &ua.TotalUnwithdrawIncome, &ua.TotalWithdrawIncome, &ua.LastWithdrawTime, &ua.Version, &ua.AllowanceState, &ua.Nickname, &ua.WithdrawDateVersion, &ua.IsDeleted)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		m[ua.MID] = ua
	}
	return
}

// UpAccount get up account by mid
func (d *Dao) UpAccount(c context.Context, mid int64) (a *model.UpAccount, err error) {
	row := d.db.QueryRow(c, _upAccountSQL, mid)
	a = &model.UpAccount{}
	if err = row.Scan(&a.TotalIncome, &a.TotalUnwithdrawIncome, &a.Version, &a.WithdrawDateVersion); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("dao.UpAccount error(%v)", err)
		}
	}
	return
}

// InsertUpAccount batch insert up account
func (d *Dao) InsertUpAccount(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_insertUpAccountSQL, values))
	if err != nil {
		log.Error("d.db.Exec InsertUpAccount error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// UpdateUpAccount update up account by mid and version instead batch update
func (d *Dao) UpdateUpAccount(c context.Context, mid, ver, totalIncome, totalUnwithdrawIncome int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _updateUpAccountSQL, totalIncome, totalUnwithdrawIncome, ver+1, mid, ver)
	if err != nil {
		log.Error("d.db.Exec UpdateUpAccount error(%v)", err)
		return
	}
	return res.RowsAffected()
}
