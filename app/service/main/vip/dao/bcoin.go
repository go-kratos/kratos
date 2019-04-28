package dao

import (
	"context"
	"time"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_vipUserBcoinSalaryList = "SELECT `id`,`mid`,`status`,`give_now_status`,`payday`,`amount`,`memo`,`ctime`,`mtime` FROM `vip_user_bcoin_salary` WHERE `mid` = ? AND `payday` > ? AND `payday` < ?;"
	_SelLastBCoinByMid      = "SELECT id,mid,`status`,give_now_status,payday,amount,memo,ctime,mtime from vip_user_bcoin_salary WHERE mid = ? ORDER BY payday DESC"
	_InsertVipBcoinSalary   = "INSERT INTO vip_user_bcoin_salary(mid,status,give_now_status,payday,amount,memo) VALUES(?,?,?,?,?,?)"
	_delBcoinSalary         = "DELETE FROM vip_bcoin_salary WHERE mid=? AND month>=?"
)

//BcoinSalaryList vip user bcoin salary list.
func (d *Dao) BcoinSalaryList(c context.Context, mid int64, start, end time.Time) (res []*model.VipBcoinSalary, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _vipUserBcoinSalaryList, mid, start, end); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("query_db")
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipBcoinSalary)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Status, &r.GiveNowStatus, &r.Month, &r.Amount, &r.Memo, &r.Ctime, &r.Mtime); err != nil {
			err = errors.WithStack(err)
			d.errProm.Incr("row_scan_db")
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

//SelLastBcoin sel last bcoin by mid.
func (d *Dao) SelLastBcoin(c context.Context, mid int64) (r *model.VipBcoinSalary, err error) {
	row := d.db.QueryRow(c, _SelLastBCoinByMid, mid)
	r = new(model.VipBcoinSalary)
	if err = row.Scan(&r.ID, &r.Mid, &r.Status, &r.GiveNowStatus, &r.Amount, &r.Memo, &r.Ctime, &r.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

//InsertVipBcoinSalary insert vip bcoin salary
func (d *Dao) InsertVipBcoinSalary(c context.Context, r *model.VipBcoinSalary) (err error) {
	if _, err = d.db.Exec(c, _InsertVipBcoinSalary, &r.Mid, &r.Status, &r.GiveNowStatus, &r.PayDay, &r.Amount, &r.Memo); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//OldTxDelBcoinSalary del bcoin salary
func (d *Dao) OldTxDelBcoinSalary(tx *sql.Tx, mid int64, month time.Time) (err error) {
	if _, err = tx.Exec(_delBcoinSalary, mid, month); err != nil {
		err = errors.WithStack(err)
	}
	return
}
