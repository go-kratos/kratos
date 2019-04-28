package dao

import (
	"context"
	xsql "database/sql"
	"fmt"
	"strings"

	"go-common/app/job/main/vip/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/time"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_getVipStatus           = "SELECT id,mid,vip_status FROM vip_user_info WHERE mid=?"
	_dupUserDiscountHistory = "INSERT INTO vip_user_discount_history(mid,discount_id,order_no,status) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE order_no= VALUES(order_no) ,status = VALUES(status)"
	_insertUserInfo         = "INSERT ignore INTO vip_user_info(mid,vip_type,vip_pay_type,vip_status,vip_start_time,vip_recent_time,vip_overdue_time,annual_vip_overdue_time,pay_channel_id,ios_overdue_time,ver) VALUES(?,?,?,?,?,?,?,?,?,?,?) "

	_updateUserInfo           = "UPDATE vip_user_info SET vip_type=?,vip_pay_type=?,vip_status=?,vip_start_time=?,vip_recent_time=?,vip_overdue_time=?,annual_vip_overdue_time=?,pay_channel_id=?,ios_overdue_time=?,ver=? WHERE mid=? AND ver=? "
	_SelUserDiscountHistorys  = "SELECT mid,discount_id FROM vip_user_discount_history WHERE id>? AND id<=? AND discount_id=?"
	_SelVipUserInfos          = "SELECT mid,vip_type,vip_pay_type,vip_status,vip_start_time,vip_recent_time,vip_overdue_time,annual_vip_overdue_time FROM  vip_user_info WHERE id>? and id<=? "
	_SelOldVipUserInfos       = "SELECT mid,vip_type,is_auto_renew,auto_renewed,IFNULL(vip_status,0),vip_recent_time,TIMESTAMP(IFNULL(vip_start_time,'2016-01-01 00:00:00')),TIMESTAMP(IFNULL(vip_overdue_time,'2016-01-01 00:00:00')),TIMESTAMP(IFNULL(annual_vip_overdue_time,'2016-01-01 00:00:00')),ver,ios_overdue_time,pay_channel_id,ctime,mtime FROM vip_user_info WHERE id>? and id<=?"
	_SelVipList               = "SELECT id,mid,vip_type,vip_status,vip_overdue_time,annual_vip_overdue_time,is_auto_renew,pay_channel_id FROM vip_user_info WHERE id>? AND id <=? AND vip_status != 0 AND annual_vip_overdue_time<=?"
	_selVipUsersSQL           = "SELECT id,mid,vip_type,vip_status,vip_overdue_time,annual_vip_overdue_time FROM vip_user_info WHERE id>? AND id <=? AND vip_pay_type = 1 AND vip_overdue_time>=? AND vip_overdue_time<?"
	_selVipUserInfoSQL        = "SELECT mid FROM vip_user_info WHERE id>? AND id <=? AND vip_status = ? AND vip_overdue_time>=? AND vip_overdue_time<?"
	_UpdateVipUserInfoByID    = "UPDATE vip_user_info SET pay_channel_id=?,vip_pay_type=?,vip_recent_time=?,vip_overdue_time=?,annual_vip_overdue_time=?, vip_type=?, vip_status=? WHERE mid=?"
	_updateVipStatus          = "UPDATE vip_user_info SET vip_status=? WHERE mid=?"
	_updateVipUserSQL         = "UPDATE vip_user_info SET vip_status=?,vip_type=?,is_auto_renew=? WHERE mid=?"
	_SelMaxID                 = "SELECT IFNULL(MAX(id),0) id FROM vip_user_info"
	_SelEffectiveScopeVipList = "SELECT `id`,`mid`,`ver`,`vip_type`,`vip_pay_type`,`pay_channel_id`,`vip_status`,`vip_start_time`,`vip_recent_time`,`vip_overdue_time`,`annual_vip_overdue_time`,`ios_overdue_time`,`ctime` FROM `vip_user_info` WHERE id>? AND id <=?;"

	_selDiscountMaxID = "SELECT IFNULL(MAX(id),0) id FROM vip_user_discount_history"
	_selUserInfoMaxID = "SELECT IFNULL(MAX(id),0) id FROM vip_user_info"

	//Vip change history
	_selOldMaxIDSQL               = "SELECT IFNULL(MAX(id),0) FROM vip_change_history"
	_selMaxIDSQL                  = "SELECT IFNULL(MAX(id),0) FROM vip_user_change_history"
	_selOldChangeHistorySQL       = "SELECT id,mid,change_type,change_time,days,month,operator_id,relation_id,batch_id,IFNULL(remark,''),batch_code_id FROM vip_change_history WHERE id>? AND id<=?"
	_selOldChangeHistoryByMidsSQL = "SELECT id,mid,change_type,change_time,days,month,operator_id,relation_id,batch_id,IFNULL(remark,''),batch_code_id FROM vip_change_history WHERE mid IN (%v)"
	_selChangeHistorySQL          = "SELECT id,mid,change_type,change_time,days,operator_id,relation_id,batch_id,remark,batch_code_id FROM vip_user_change_history WHERE id>? AND id<=?"
	_selChangeHistoryByMidsSQL    = "SELECT id,mid,change_type,change_time,days,operator_id,relation_id,batch_id,remark,batch_code_id FROM vip_user_change_history WHERE mid IN (%v)"
	_addChangeHistoryBatchSQL     = "INSERT INTO vip_user_change_history(mid,change_type,change_time,days,operator_id,relation_id,batch_id,remark,batch_code_id) VALUES"

	// old vip db
	_selOldVipUserInfo = "SELECT mid,vip_type,is_auto_renew,auto_renewed,IFNULL(vip_status,0),vip_recent_time,TIMESTAMP(IFNULL(vip_start_time,'2016-01-01 00:00:00')),TIMESTAMP(IFNULL(vip_overdue_time,'2016-01-01 00:00:00')),annual_vip_overdue_time,ios_overdue_time,ver,pay_channel_id FROM vip_user_info WHERE mid=?"

	_selVipByMidsSQL = "SELECT id,mid,vip_type,vip_status,vip_recent_time,vip_start_time,vip_overdue_time,annual_vip_overdue_time,vip_pay_type,pay_channel_id,ios_overdue_time,ctime,mtime,ver FROM vip_user_info WHERE mid IN (%s)"

	_selVipUserInfoByMid = "SELECT id,mid,ver,vip_type,vip_pay_type,vip_status,vip_start_time,vip_overdue_time,annual_vip_overdue_time,vip_recent_time FROM vip_user_info WHERE mid = ?"
)

//AddChangeHistoryBatch batch add change history
func (d *Dao) AddChangeHistoryBatch(res []*model.VipChangeHistory) (err error) {
	var values []string
	if len(res) <= 0 {
		return
	}
	for _, v := range res {
		value := fmt.Sprintf("('%v','%v','%v','%v','%v','%v','%v','%v','%v')", v.Mid, v.ChangeType, v.ChangeTime.Time().Format("2006-01-02 15:04:05"), v.Days, v.OperatorID, v.RelationID, v.BatchID, v.Remark, v.BatchCodeID)
		values = append(values, value)
	}
	valuesStr := strings.Join(values, ",")
	if _, err = d.db.Exec(context.TODO(), _addChangeHistoryBatchSQL+valuesStr); err != nil {
		log.Error("AddChangeHistoryBatch d.db.exec(%v),error(%v)", valuesStr, err)
		return
	}
	return
}

//SelChangeHistory .
func (d *Dao) SelChangeHistory(c context.Context, startID, endID int64) (res []*model.VipChangeHistory, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _selChangeHistorySQL, startID, endID); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipChangeHistory)
		if err = rows.Scan(&r.ID, &r.Mid, &r.ChangeType, &r.ChangeTime, &r.Days, &r.OperatorID, &r.RelationID, &r.BatchID, &r.Remark, &r.BatchCodeID); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

//SelOldChangeHistory .
func (d *Dao) SelOldChangeHistory(c context.Context, startID, endID int64) (res []*model.VipChangeHistory, err error) {
	var rows *sql.Rows
	if rows, err = d.oldDb.Query(c, _selOldChangeHistorySQL, startID, endID); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipChangeHistory)
		if err = rows.Scan(&r.ID, &r.Mid, &r.ChangeType, &r.ChangeTime, &r.Days, &r.Month, &r.OperatorID, &r.RelationID, &r.BatchID, &r.Remark, &r.BatchCodeID); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return

}

//SelChangeHistoryMaps .
func (d *Dao) SelChangeHistoryMaps(c context.Context, mids []int64) (res map[int64][]*model.VipChangeHistory, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_selChangeHistoryByMidsSQL, xstr.JoinInts(mids))); err != nil {
		err = errors.WithStack(err)
		return
	}
	res = make(map[int64][]*model.VipChangeHistory)
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipChangeHistory)
		if err = rows.Scan(&r.ID, &r.Mid, &r.ChangeType, &r.ChangeTime, &r.Days, &r.OperatorID, &r.RelationID, &r.BatchID, &r.Remark, &r.BatchCodeID); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		histories := res[r.Mid]
		histories = append(histories, r)
		res[r.Mid] = histories
	}
	err = rows.Err()
	return
}

//SelOldChangeHistoryMaps .
func (d *Dao) SelOldChangeHistoryMaps(c context.Context, mids []int64) (res map[int64][]*model.VipChangeHistory, err error) {
	var rows *sql.Rows
	if rows, err = d.oldDb.Query(c, fmt.Sprintf(_selOldChangeHistoryByMidsSQL, xstr.JoinInts(mids))); err != nil {
		err = errors.WithStack(err)
		return
	}
	res = make(map[int64][]*model.VipChangeHistory)

	defer rows.Close()
	for rows.Next() {
		r := new(model.VipChangeHistory)
		if err = rows.Scan(&r.ID, &r.Mid, &r.ChangeType, &r.ChangeTime, &r.Days, &r.Month, &r.OperatorID, &r.RelationID, &r.BatchID, &r.Remark, &r.BatchCodeID); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		histories := res[r.Mid]
		histories = append(histories, r)
		res[r.Mid] = histories
	}
	err = rows.Err()
	return

}

//SelOldChangeHistoryMaxID .
func (d *Dao) SelOldChangeHistoryMaxID(c context.Context) (id int64, err error) {
	var row = d.oldDb.QueryRow(c, _selOldMaxIDSQL)
	if err = row.Scan(&id); err != nil {
		log.Error("SelMaxID row.Scan() error(%v)", err)
		return
	}
	return
}

//SelChangeHistoryMaxID .
func (d *Dao) SelChangeHistoryMaxID(c context.Context) (id int64, err error) {
	var row = d.db.QueryRow(c, _selMaxIDSQL)
	if err = row.Scan(&id); err != nil {
		log.Error("SelMaxID row.Scan() error(%v)", err)
		return
	}
	return
}

// VipStatus get user vip status.
func (d *Dao) VipStatus(c context.Context, mid int64) (res *model.VipUserInfo, err error) {
	row := d.db.QueryRow(c, _getVipStatus, mid)
	res = new(model.VipUserInfo)
	if err = row.Scan(&res.ID, &res.Mid, &res.Status); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res = nil
		} else {
			err = errors.Wrapf(err, "d.VipStatus(%d)", mid)
		}
	}
	return
}

//SelMaxID get max id
func (d *Dao) SelMaxID(c context.Context) (mID int, err error) {
	var row = d.db.QueryRow(c, _SelMaxID)
	if err = row.Scan(&mID); err != nil {
		log.Error("SelMaxID row.Scan() error(%v)", err)
		return
	}
	return
}

//UpdateVipUserInfo update vip user info by id
func (d *Dao) UpdateVipUserInfo(tx *sql.Tx, r *model.VipUserInfo) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_UpdateVipUserInfoByID, r.PayChannelID, r.PayType, r.RecentTime, r.OverdueTime, r.AnnualVipOverdueTime, r.Type, r.Status, r.Mid); err != nil {
		log.Error("UpdateVipUserInfo: db.Exec(%v) error(%v)", r, err)
		return
	}

	if a, err = res.RowsAffected(); err != nil {
		log.Error("UpdateVipUserInfo: res.RowsAffected error(%v)", err)
	}
	return
}

// UpdateVipUser update vip user info
func (d *Dao) UpdateVipUser(c context.Context, mid int64, status, vipType int8, payType int8) (eff int64, err error) {
	var res xsql.Result
	if res, err = d.oldDb.Exec(c, _updateVipUserSQL, status, vipType, payType, mid); err != nil {
		err = errors.Wrapf(err, "d.UpdateVipStatus(%d)", mid)
		return
	}
	if eff, err = res.RowsAffected(); err != nil {
		err = errors.Wrapf(err, "d.UpdateVipStatus(%d) res.RowsAffected", mid)
	}
	return
}

//UpdateVipStatus update vip status info by mid
func (d *Dao) UpdateVipStatus(c context.Context, mid int64, status int) (eff int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _updateVipStatus, status, mid); err != nil {
		err = errors.Wrapf(err, "d.UpdateVipStatus(%d)", mid)
		return
	}
	if eff, err = res.RowsAffected(); err != nil {
		err = errors.Wrapf(err, "d.UpdateVipStatus(%d) res.RowsAffected", mid)
	}
	return
}

// DupUserDiscountHistory add user discount history.
func (d *Dao) DupUserDiscountHistory(tx *sql.Tx, r *model.VipUserDiscountHistory) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_dupUserDiscountHistory, r.Mid, r.DiscountID, r.OrderNo, r.Status); err != nil {
		err = errors.WithStack(err)
		return
	}

	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//AddUserInfo add user info.
func (d *Dao) AddUserInfo(tx *sql.Tx, r *model.VipUserInfo) (eff int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_insertUserInfo, r.Mid, r.Type, r.PayType, r.Status, r.StartTime, r.RecentTime, r.OverdueTime, r.AnnualVipOverdueTime, r.PayChannelID, r.IosOverdueTime, r.Ver); err != nil {
		log.Error("AddUserDiscountHistory d.db.exec(%v) error(%v)", r, err)
		return
	}
	if eff, err = res.RowsAffected(); err != nil {
		log.Error("AddUserDiscountHistory RowsAffected(%v)", err)
	}
	return
}

//SelVipUserInfo sel vipuser info
func (d *Dao) SelVipUserInfo(c context.Context, mid int64) (r *model.VipUserInfo, err error) {
	var row = d.db.QueryRow(c, _selVipUserInfoByMid, mid)
	r = new(model.VipUserInfo)
	if err = row.Scan(&r.ID, &r.Mid, &r.Ver, &r.Type, &r.PayType, &r.Status, &r.StartTime, &r.OverdueTime, &r.AnnualVipOverdueTime, &r.RecentTime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

//UpdateUserInfo add user info
func (d *Dao) UpdateUserInfo(tx *sql.Tx, r *model.VipUserInfo) (eff int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_updateUserInfo, r.Type, r.PayType, r.Status, r.StartTime, r.RecentTime, r.OverdueTime, r.AnnualVipOverdueTime, r.PayChannelID, r.IosOverdueTime, r.Ver, r.Mid, r.OldVer); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//SelUserDiscountMaxID sel discount maxID
func (d *Dao) SelUserDiscountMaxID(c context.Context) (maxID int, err error) {
	var row = d.db.QueryRow(c, _selDiscountMaxID)
	if err = row.Scan(&maxID); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_scan")
		return
	}
	return
}

//SelUserDiscountHistorys sel user discount hsitorys
func (d *Dao) SelUserDiscountHistorys(c context.Context, sID, eID, discountID int) (res []*model.VipUserDiscountHistory, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _SelUserDiscountHistorys, sID, eID, discountID); err != nil {
		log.Error("SelUserDiscountHistorys d.db.query(sID:%d,eID:%d,discountID:%d) error(%+v)", sID, eID, discountID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipUserDiscountHistory)
		if err = rows.Scan(&r.Mid, &r.DiscountID); err != nil {
			log.Error("SelUserDiscountHistorys rows.scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

//SelOldUserInfoMaxID sel old userinfo maxID
func (d *Dao) SelOldUserInfoMaxID(c context.Context) (maxID int, err error) {
	var row = d.oldDb.QueryRow(c, _selUserInfoMaxID)
	if err = row.Scan(&maxID); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_scan")
		return
	}
	return
}

// SelUserInfoMaxID select userinfo maxID.
func (d *Dao) SelUserInfoMaxID(c context.Context) (maxID int, err error) {
	var row = d.db.QueryRow(c, _selUserInfoMaxID)
	if err = row.Scan(&maxID); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_scan")
		return
	}
	return
}

//SelUserInfos sel user infos
func (d *Dao) SelUserInfos(c context.Context, sID, eID int) (res []*model.VipUserInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _SelVipUserInfos, sID, eID); err != nil {
		log.Error("SelUserDiscountHistorys d.db.query(sID:%v,eID:%v) error(%v)", sID, eID, err)
		return
	}
	for rows.Next() {
		r := new(model.VipUserInfo)
		if err = rows.Scan(&r.Mid, &r.Type, &r.PayType, &r.Status, &r.StartTime, &r.RecentTime, &r.OverdueTime, &r.AnnualVipOverdueTime); err != nil {
			log.Error("SelUserDiscountHistorys rows.scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

//SelVipList sel vipuserinfo list
func (d *Dao) SelVipList(c context.Context, id, endID int, ot string) (res []*model.VipUserInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.oldDb.Query(c, _SelVipList, id, endID, ot); err != nil {
		log.Error("SelVipList db.query(id:%v,ot:%v,endID:%v) error(%v)", id, ot, endID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipUserInfo)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Type, &r.Status, &r.OverdueTime, &r.AnnualVipOverdueTime, &r.PayType, &r.PayChannelID); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

//SelVipUsers .
func (d *Dao) SelVipUsers(c context.Context, id, endID int, startTime, endTime time.Time) (res []*model.VipUserInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _selVipUsersSQL, id, endID, startTime, endTime); err != nil {
		log.Error("SelVipList db.query(id:%v,endID:%v,st:%v,et:%v) error(%v)", id, endID, startTime, endTime, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipUserInfo)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Type, &r.Status, &r.OverdueTime, &r.AnnualVipOverdueTime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

//SelVipUserInfos sel vipuser info datas
func (d *Dao) SelVipUserInfos(c context.Context, id, endID int, startTime, endTime time.Time, status int) (res []int64, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _selVipUserInfoSQL, id, endID, status, startTime.Time().Format("2006-01-02"), endTime.Time().Format("2006-01-02")); err != nil {
		log.Error("SelVipList db.query(id:%v,endID:%v,st:%v,et:%v) error(%v)", id, endID, startTime, endTime, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r int64
		if err = rows.Scan(&r); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

//OldVipInfo get old user by mid.
func (d *Dao) OldVipInfo(c context.Context, mid int64) (r *model.VipUserInfoOld, err error) {
	var row = d.oldDb.QueryRow(c, _selOldVipUserInfo, mid)
	r = new(model.VipUserInfoOld)
	if err = row.Scan(&r.Mid, &r.Type, &r.IsAutoRenew, &r.AutoRenewed, &r.Status, &r.RecentTime, &r.StartTime, &r.OverdueTime,
		&r.AnnualVipOverdueTime, &r.IosOverdueTime, &r.Ver, &r.PayChannelID); err != nil {
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

// SelEffectiveScopeVipList get vip list by id scope.
func (d *Dao) SelEffectiveScopeVipList(c context.Context, id, endID int) (res []*model.VipInfoDB, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _SelEffectiveScopeVipList, id, endID); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipInfoDB)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Ver, &r.Type, &r.PayType, &r.PayChannelID, &r.Status, &r.StartTime, &r.RecentTime, &r.OverdueTime, &r.AnnualVipOverdueTime,
			&r.IosOverdueTime, &r.Ctime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// SelVipByIds select vip by ids .
func (d *Dao) SelVipByIds(c context.Context, ids []int64) (res map[int64]*model.VipUserInfo, err error) {
	var (
		rows  *sql.Rows
		idStr = xstr.JoinInts(ids)
	)
	res = make(map[int64]*model.VipUserInfo, len(ids))
	if rows, err = d.db.Query(c, fmt.Sprintf(_selVipByMidsSQL, idStr)); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipUserInfo)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Type, &r.Status, &r.RecentTime, &r.StartTime, &r.OverdueTime, &r.AnnualVipOverdueTime, &r.PayType, &r.PayChannelID, &r.IosOverdueTime, &r.Ctime, &r.Mtime, &r.Ver); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res[r.Mid] = r
	}
	return
}

//SelOldUserInfoMaps sel old user info map.
func (d *Dao) SelOldUserInfoMaps(c context.Context, sID, eID int) (res map[int64]*model.VipUserInfoOld, err error) {
	var (
		rows *sql.Rows
	)
	res = make(map[int64]*model.VipUserInfoOld, eID-sID)
	if rows, err = d.oldDb.Query(c, _SelOldVipUserInfos, sID, eID); err != nil {
		err = errors.WithStack(err)
		return
	}
	for rows.Next() {
		r := new(model.VipUserInfoOld)
		if err = rows.Scan(&r.Mid, &r.Type, &r.IsAutoRenew, &r.AutoRenewed, &r.Status, &r.RecentTime, &r.StartTime, &r.OverdueTime, &r.AnnualVipOverdueTime, &r.Ver, &r.IosOverdueTime, &r.PayChannelID, &r.Ctime, &r.Mtime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res[r.Mid] = r
	}
	return
}
