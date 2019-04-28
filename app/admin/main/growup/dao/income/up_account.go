package income

import (
	"context"
	"fmt"

	model "go-common/app/admin/main/growup/model/income"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// count
	_upAccountCountSQL = "SELECT count(*) FROM up_account WHERE %s is_deleted = ?"
	// select
	_upAccountSQL      = "SELECT mid,total_income,total_unwithdraw_income,total_withdraw_income,withdraw_date_version,last_withdraw_time,mtime FROM up_account WHERE %s is_deleted = ? LIMIT ?,?"
	_upAccountByMIDSQL = "SELECT mid,total_income,total_unwithdraw_income,withdraw_date_version,version FROM up_account WHERE mid = ? AND is_deleted = 0"
	// update
	_breachUpAccountSQL = "UPDATE up_account SET total_income = ?, total_unwithdraw_income = ?, version = ? WHERE mid = ? AND version = ? AND is_deleted = 0"
)

// UpAccountCount get up_account count
func (d *Dao) UpAccountCount(c context.Context, query string, isDeleted int) (total int64, err error) {
	err = d.db.QueryRow(c, fmt.Sprintf(_upAccountCountSQL, query), isDeleted).Scan(&total)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

// ListUpAccount list up account bu query
func (d *Dao) ListUpAccount(c context.Context, query string, isDeleted, from, limit int) (ups []*model.UpAccount, err error) {
	ups = make([]*model.UpAccount, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_upAccountSQL, query), isDeleted, from, limit)
	if err != nil {
		log.Error("ListUpAccount d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := &model.UpAccount{}
		err = rows.Scan(&up.MID, &up.TotalIncome, &up.TotalUnwithdrawIncome, &up.TotalWithdrawIncome, &up.WithdrawDateVersion, &up.LastWithdrawTime, &up.MTime)
		if err != nil {
			log.Error("ListUpAccount rows scan error(%v)", err)
			return
		}
		ups = append(ups, up)
	}
	err = rows.Err()
	return
}

// GetUpAccount get up_account by mid
func (d *Dao) GetUpAccount(c context.Context, mid int64) (up *model.UpAccount, err error) {
	up = &model.UpAccount{}
	err = d.db.QueryRow(c, _upAccountByMIDSQL, mid).Scan(&up.MID, &up.TotalIncome, &up.TotalUnwithdrawIncome, &up.WithdrawDateVersion, &up.Version)
	return
}

// TxBreachUpAccount breach up_account
func (d *Dao) TxBreachUpAccount(tx *sql.Tx, total, unwithdraw, mid, newVersion, oldVersion int64) (rows int64, err error) {
	res, err := tx.Exec(_breachUpAccountSQL, total, unwithdraw, newVersion, mid, oldVersion)
	if err != nil {
		return
	}
	return res.RowsAffected()
}
