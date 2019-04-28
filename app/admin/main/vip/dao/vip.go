package dao

import (
	"context"
	xsql "database/sql"

	"go-common/app/admin/main/vip/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/time"
)

// const .
const (
	_selVipUserInfoByMid    = "SELECT id,mid,vip_type,vip_pay_type,vip_status,vip_start_time,vip_overdue_time,annual_vip_overdue_time,ctime,mtime,vip_recent_time FROM vip_user_info WHERE mid = ?"
	_UpdateVipUserInfoByID  = "UPDATE vip_user_info SET vip_type=?,vip_status=?,vip_overdue_time=?,annual_vip_overdue_time=?,vip_recent_time=? WHERE id=?"
	_InsertVipChangeHistory = "INSERT INTO vip_user_change_history(mid,change_type,change_time,days,operator_id,relation_id,batch_id,remark) VALUES(?,?,?,?,?,?,?,?)"
	_DelBcoinSalary         = "DELETE FROM vip_user_bcoin_salary WHERE status = 0 AND mid = ? and payday>=?"
)

// DelBcoinSalary del bacoin salary
func (d *Dao) DelBcoinSalary(tx *sql.Tx, mid int64, month time.Time) (err error) {
	if _, err = tx.Exec(_DelBcoinSalary, mid, month); err != nil {
		log.Error("InertVipBcoinSalary.exec(mid:%v,month:%v) error(%+v)", mid, month, err)
	}
	return
}

// SelVipUserInfo select user info by mid
func (d *Dao) SelVipUserInfo(c context.Context, mid int64) (r *model.VipUserInfo, err error) {
	var row = d.db.QueryRow(c, _selVipUserInfoByMid, mid)
	r = new(model.VipUserInfo)
	if err = row.Scan(&r.ID, &r.Mid, &r.VipType, &r.VipPayType, &r.VipStatus, &r.VipStartTime, &r.VipOverdueTime, &r.AnnualVipOverdueTime, &r.Ctime, &r.Mtime, &r.VipRecentTime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}

	}
	return
}

// UpdateVipUserInfo update vip user info by id
func (d *Dao) UpdateVipUserInfo(tx *sql.Tx, r *model.VipUserInfo) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_UpdateVipUserInfoByID, r.VipType, r.VipStatus, r.VipOverdueTime, r.AnnualVipOverdueTime, r.VipRecentTime, r.ID); err != nil {
		log.Error("UpdateVipUserInfo: db.Exec(%v) error(%v)", r, err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		log.Error("UpdateVipUserInfo: res.RowsAffected error(%v)", err)
		return
	}

	return
}

// InsertVipChangeHistory insert vipChangeHistory
func (d *Dao) InsertVipChangeHistory(tx *sql.Tx, r *model.VipChangeHistory) (id int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_InsertVipChangeHistory, r.Mid, r.ChangeType, r.ChangeTime, r.Days, r.OperatorID, r.RelationID, r.BatchID, r.Remark); err != nil {
		log.Error("InsertVipChangeHistory db.Exec(%v) error(%v)", r, err)
	} else {
		if id, err = res.LastInsertId(); err != nil {
			log.Error("InsertVipChangeHistory LastInsertId() error(%v)", err)
		}
	}
	return
}
