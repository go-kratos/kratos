package dao

import (
	"context"
	xsql "database/sql"
	"fmt"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"
	"go-common/library/time"

	"github.com/pkg/errors"
)

const (
	_vipUserInfoByMid                   = "SELECT id,mid,vip_type,vip_pay_type,vip_status,pay_channel_id,vip_start_time,vip_overdue_time,annual_vip_overdue_time,ctime,mtime,vip_recent_time,ios_overdue_time,ver FROM vip_user_info WHERE mid = ?;"
	_updateVipStatusTypeSQL             = "UPDATE `vip_user_info` SET `vip_status`=? AND `vip_type`=? WHERE `mid`=? ;"
	_vipChangeHistorySQL                = "SELECT id,mid,change_type,change_time,days,operator_id,relation_id,batch_id,remark FROM vip_user_change_history WHERE mid=? ORDER BY change_time DESC"
	_changeHistoryCountSQL              = "SELECT COUNT(1) FROM vip_user_change_history WHERE mid=?"
	_updateChannelIDSQL                 = "UPDATE vip_user_info SET pay_channel_id=?,ver=? WHERE mid=? AND ver=?"
	_dupUserDiscountSQL                 = "INSERT INTO vip_user_discount_history(mid,discount_id,order_no,status) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE order_no= VALUES(order_no) ,status = VALUES(status);"
	_updatePayTypeSQL                   = "UPDATE vip_user_info SET vip_pay_type = ?,ver=? WHERE mid =? AND ver=?"
	_selVipChangeHistoryByRelationIDSQL = "SELECT id,mid,change_type,change_time,days,operator_id,relation_id,batch_id,remark,ctime FROM vip_user_change_history WHERE relation_id=?"
	_InsertVipChangeHistory             = "INSERT INTO vip_user_change_history(mid,change_type,change_time,days,operator_id,relation_id,batch_id,remark,batch_code_id) VALUES(?,?,?,?,?,?,?,?,?)"
	_SelVipUserInfoByMid                = "SELECT id,mid,vip_type,vip_pay_type,vip_status,ver,vip_start_time,vip_overdue_time,annual_vip_overdue_time,ctime,mtime,vip_recent_time,ios_overdue_time FROM vip_user_info WHERE mid = ?"
	_InsertVipUserInfo                  = "INSERT INTO vip_user_info(mid,vip_type,vip_pay_type,vip_status,vip_start_time,vip_overdue_time,annual_vip_overdue_time,vip_recent_time) VALUES(?,?,?,?,?,?,?,?)"
	_addIosVipUserInfo                  = "INSERT INTO vip_user_info(mid,vip_type,vip_status,vip_start_time,ios_overdue_time) VALUES(?,?,?,?,?)"
	_UpdateVipUserInfoByID              = "UPDATE vip_user_info SET vip_type=?,vip_pay_type=?,vip_status=?,vip_overdue_time=?,annual_vip_overdue_time=?,vip_recent_time=?,ver=? WHERE mid=? AND ver=?"
	_updateIosUserInfoSQL               = "UPDATE vip_user_info SET ios_overdue_time=? WHERE mid = ?"
	_updateIosRenewUserInfoSQL          = "UPDATE vip_user_info SET vip_pay_type=?,ver=?,pay_channel_id=? WHERE mid = ? AND ver=?"
	_addUserDiscount                    = "INSERT IGNORE INTO vip_user_discount_history(mid,discount_id,order_no,status) VALUES(?,?,?,?)"
	//refund
	_oldUpdateVipUserInfoByID = "UPDATE vip_user_info SET vip_type=?,vip_status=?,vip_overdue_time=?,annual_vip_overdue_time=?,ver=? WHERE id=? AND ver=?"
	_oldSelVipUserInfoByMid   = "SELECT id,mid,vip_type,vip_status,ver,vip_start_time,vip_overdue_time,annual_vip_overdue_time,vip_recent_time,ios_overdue_time FROM vip_user_info WHERE mid = ?"
	_oldAddVipChangeHistory   = "INSERT INTO vip_change_history(mid,change_type,change_time,days,operator_id,relation_id,batch_id,remark,batch_code_id) VALUES(?,?,?,?,?,?,?,?,?)"
	//sync
	_syncAddUser    = "INSERT INTO vip_user_info(mid,vip_type,vip_pay_type,vip_status,vip_start_time,vip_overdue_time,annual_vip_overdue_time,vip_recent_time,ios_overdue_time,pay_channel_id,ver) VALUES(?,?,?,?,?,?,?,?,?,?,?)"
	_syncUpdateUser = "UPDATE vip_user_info SET vip_type=?,vip_pay_type=?,vip_status=?,vip_overdue_time=?,annual_vip_overdue_time=?,vip_recent_time=?,ios_overdue_time=?,pay_channel_id=?,ver=? WHERE mid=? AND ver=?"

	//clean job cache
	_cleanjobSQL = "UPDATE vip_user_info SET ver = ver + 1 WHERE mid=?"
)

//VipInfo select user info by mid.
func (d *Dao) VipInfo(c context.Context, mid int64) (r *model.VipInfoDB, err error) {
	var row = d.db.QueryRow(c, _vipUserInfoByMid, mid)
	r = new(model.VipInfoDB)
	if err = row.Scan(&r.ID, &r.Mid, &r.VipType, &r.VipPayType, &r.VipStatus, &r.PayChannelID, &r.VipStartTime, &r.VipOverdueTime, &r.AnnualVipOverdueTime,
		&r.Ctime, &r.Mtime, &r.VipRecentTime, &r.IosOverdueTime, &r.Ver); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
		} else {
			err = errors.WithStack(err)
			d.errProm.Incr("row_scan_db")
		}
	}
	return
}

// UpdateVipTypeAndStatus update vip type and status.
func (d *Dao) UpdateVipTypeAndStatus(c context.Context, mid int64, vipStatus int32, vipType int32) (ret int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _updateVipStatusTypeSQL, vipStatus, vipType, mid); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("exec_db")
		return
	}
	return res.RowsAffected()
}

//SelChangeHistoryCount .
func (d *Dao) SelChangeHistoryCount(c context.Context, arg *model.ArgChangeHistory) (count int64, err error) {
	row := d.db.QueryRow(c, _changeHistoryCountSQL, arg.Mid)

	if err = row.Scan(&count); err != nil {
		if sql.ErrNoRows == err {
			err = nil
			count = 0
		}
		err = errors.WithStack(err)
	}

	return
}

//SelChangeHistory .
func (d *Dao) SelChangeHistory(c context.Context, arg *model.ArgChangeHistory) (vcs []*model.VipChangeHistory, err error) {
	SQLStr := _vipChangeHistorySQL
	if arg.Pn > 0 && arg.Ps > 0 {
		SQLStr += fmt.Sprintf(" LIMIT %d,%d", (arg.Pn-1)*arg.Ps, arg.Ps)
	}
	var rows *sql.Rows
	if rows, err = d.db.Query(c, SQLStr, arg.Mid); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("query_db")
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipChangeHistory)
		if rows.Scan(&r.ID, &r.Mid, &r.ChangeType, &r.ChangeTime, &r.Days, &r.OperatorID, &r.RelationID, &r.BatchID, &r.Remark); err != nil {
			err = errors.WithStack(err)
			d.errProm.Incr("row_scan_db")
			vcs = nil
			return
		}
		vcs = append(vcs, r)
	}
	err = rows.Err()
	return
}

//TxUpdateChannelID .
func (d *Dao) TxUpdateChannelID(tx *sql.Tx, mid int64, payChannelID int32, ver int64, oldVer int64) (err error) {
	if _, err = tx.Exec(_updateChannelIDSQL, payChannelID, ver, mid, oldVer); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//TxDupUserDiscount .
func (d *Dao) TxDupUserDiscount(tx *sql.Tx, mid, discountID int64, orderNo string, status int8) (err error) {
	if _, err = tx.Exec(_dupUserDiscountSQL, mid, discountID, orderNo, status); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//UpdatePayType .
func (d *Dao) UpdatePayType(tx *sql.Tx, mid int64, payType int8, ver, oldVer int64) (a int64, err error) {
	var result xsql.Result
	if result, err = tx.Exec(_updatePayTypeSQL, payType, ver, mid, oldVer); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = result.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//SelVipChangeHistory select vip_change_history by relationId
func (d *Dao) SelVipChangeHistory(c context.Context, relationID string) (r *model.VipChangeHistory, err error) {
	var row = d.db.QueryRow(c, _selVipChangeHistoryByRelationIDSQL, relationID)
	r = new(model.VipChangeHistory)
	if err = row.Scan(&r.ID, &r.Mid, &r.ChangeType, &r.ChangeTime, &r.Days, &r.OperatorID, &r.RelationID, &r.BatchID, &r.Remark, &r.Ctime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
		} else {
			err = errors.WithStack(err)
			d.errProm.Incr("scan_error")
		}
	}
	return
}

//InsertVipChangeHistory insert vipChangeHistory
func (d *Dao) InsertVipChangeHistory(tx *sql.Tx, r *model.VipChangeHistory) (id int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_InsertVipChangeHistory, r.Mid, r.ChangeType, r.ChangeTime, r.Days, r.OperatorID, r.RelationID, r.BatchID, r.Remark, r.BatchCodeID); err != nil {
		err = errors.WithStack(err)
	} else {
		if id, err = res.LastInsertId(); err != nil {
			err = errors.WithStack(err)
		}
	}
	return
}

//TxSelVipUserInfo .
func (d *Dao) TxSelVipUserInfo(tx *sql.Tx, mid int64) (r *model.VipInfoDB, err error) {
	var row = tx.QueryRow(_SelVipUserInfoByMid, mid)
	r = new(model.VipInfoDB)
	if err = row.Scan(&r.ID, &r.Mid, &r.VipType, &r.VipPayType, &r.VipStatus, &r.Ver, &r.VipStartTime, &r.VipOverdueTime, &r.AnnualVipOverdueTime, &r.Ctime, &r.Mtime, &r.VipRecentTime, &r.IosOverdueTime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

// TxAddIosVipUserInfo  tx add ios vip user info.
func (d *Dao) TxAddIosVipUserInfo(tx *sql.Tx, r *model.VipInfoDB) (eff int64, err error) {
	var result xsql.Result

	if result, err = tx.Exec(_addIosVipUserInfo, r.Mid, r.VipType, r.VipStatus, r.VipStartTime, r.IosOverdueTime); err != nil {
		err = errors.WithStack(err)
		return
	}

	if eff, err = result.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//InsertVipUserInfo insert vipUserInfo
func (d *Dao) InsertVipUserInfo(tx *sql.Tx, r *model.VipInfoDB) (err error) {
	if _, err = tx.Exec(_InsertVipUserInfo, r.Mid, r.VipType, r.VipPayType, r.VipStatus, r.VipStartTime, r.VipOverdueTime, r.AnnualVipOverdueTime, r.VipRecentTime); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//TxUpdateIosUserInfo update ios time
func (d *Dao) TxUpdateIosUserInfo(tx *sql.Tx, iosTime time.Time, mid int64) (eff int64, err error) {
	var result xsql.Result
	if result, err = tx.Exec(_updateIosUserInfoSQL, iosTime, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff, err = result.RowsAffected(); err != nil {
		err = errors.WithStack(err)

	}
	return
}

//TxUpdateIosRenewUserInfo .
func (d *Dao) TxUpdateIosRenewUserInfo(tx *sql.Tx, paychannelID, ver, oldVer, mid int64, payType int8) (err error) {
	if _, err = tx.Exec(_updateIosRenewUserInfoSQL, payType, ver, paychannelID, mid, oldVer); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//UpdateVipUserInfo update vip user info by id
func (d *Dao) UpdateVipUserInfo(tx *sql.Tx, r *model.VipInfoDB, ver int64) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_UpdateVipUserInfoByID, r.VipType, r.VipPayType, r.VipStatus, r.VipOverdueTime, r.AnnualVipOverdueTime, r.VipRecentTime, r.Ver, r.Mid, ver); err != nil {
		err = errors.WithStack(err)
	} else {
		if a, err = res.RowsAffected(); err != nil {
			err = errors.WithStack(err)
		}
	}
	return
}

// TxAddUserDiscount add user discount history.
func (d *Dao) TxAddUserDiscount(tx *sql.Tx, r *model.VipUserDiscountHistory) (eff int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_addUserDiscount, r.Mid, r.DiscountID, r.OrderNo, r.Status); err != nil {
		err = errors.WithStack(err)
		return
	}
	eff, err = res.RowsAffected()
	return
}

//SyncAddUser insert vipUserInfo
func (d *Dao) SyncAddUser(tx *sql.Tx, r *model.VipInfoDB) (err error) {
	if _, err = tx.Exec(_syncAddUser, r.Mid, r.VipType, r.VipPayType, r.VipStatus, r.VipStartTime, r.VipOverdueTime, r.AnnualVipOverdueTime, r.VipRecentTime, r.IosOverdueTime, r.PayChannelID, r.Ver); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//SyncUpdateUser insert vipUserInfo
func (d *Dao) SyncUpdateUser(tx *sql.Tx, r *model.VipInfoDB, ver int64) (eff int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_syncUpdateUser, r.VipType, r.VipPayType, r.VipStatus, r.VipOverdueTime, r.AnnualVipOverdueTime, r.VipRecentTime, r.IosOverdueTime, r.PayChannelID, r.Ver, r.Mid, ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//OldTxSelVipUserInfo .
func (d *Dao) OldTxSelVipUserInfo(tx *sql.Tx, mid int64) (r *model.VipInfoDB, err error) {
	var row = tx.QueryRow(_oldSelVipUserInfoByMid, mid)
	r = new(model.VipInfoDB)
	if err = row.Scan(&r.ID, &r.Mid, &r.VipType, &r.VipStatus, &r.Ver, &r.VipStartTime, &r.VipOverdueTime, &r.AnnualVipOverdueTime, &r.VipRecentTime, &r.IosOverdueTime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

//OldTxUpdateVipUserInfo .
func (d *Dao) OldTxUpdateVipUserInfo(tx *sql.Tx, r *model.VipInfoDB, ver int64) (eff int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_oldUpdateVipUserInfoByID, r.VipType, r.VipStatus, r.VipOverdueTime, r.AnnualVipOverdueTime, r.Ver, r.ID, ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	eff, err = res.RowsAffected()
	return
}

//OldTxAddChangeHistory .
func (d *Dao) OldTxAddChangeHistory(tx *sql.Tx, r *model.VipChangeHistory) (err error) {
	if _, err = tx.Exec(_oldAddVipChangeHistory, r.Mid, r.ChangeType, r.ChangeTime, r.Days, r.OperatorID, r.RelationID, r.BatchID, r.Remark, r.BatchCodeID); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//CleanCache clean job cache.
func (d *Dao) CleanCache(c context.Context, mid int64) (err error) {
	if _, err = d.db.Exec(c, _cleanjobSQL, mid); err != nil {
		err = errors.WithStack(err)
	}
	return
}
