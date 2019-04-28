package dao

import (
	"context"
	xsql "database/sql"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_selVipUserInfoByMid     = "select id,mid,vip_type,vip_status,vip_start_time,vip_overdue_time,annual_vip_overdue_time,wander,access_status,ctime,mtime,vip_recent_time from vip_user_info where mid = ?"
	_updateVipUserInfoByID   = "update vip_user_info set vip_type=?,vip_status=?,vip_overdue_time=?,annual_vip_overdue_time=?,vip_recent_time=? where id=?"
	_insertVipUserInfo       = "INSERT INTO vip_user_info(mid,vip_type,vip_status,vip_start_time,vip_overdue_time,annual_vip_overdue_time,wander,access_status,vip_recent_time) values(?,?,?,?,?,?,?,?,?)"
	_insertVipChangeHistory  = "INSERT INTO vip_change_history(mid,change_type,change_time,days,month,operator_id,relation_id,batch_id,remark) values(?,?,?,?,?,?,?,?,?)"
	_selVipChangeHistory     = "select id,mid,change_type,change_time,days,month,operator_id,relation_id,batch_id,remark,ctime from vip_change_history where relation_id=? AND batch_id=?"
	_selVipResourceBatchByID = "select id,pool_id,unit,count,ver,start_time,end_time,surplus_count,code_use_count,direct_use_count,ctime,mtime from vip_resource_batch where id=?"
	_updateBatchCountByID    = "update vip_resource_batch set ver=ver+1,surplus_count=surplus_count-1,direct_use_count=direct_use_count+1,code_use_count=? where id=? and surplus_count>0"
	_selResourcePoolByID     = "select id,pool_name,business_id,business_name,reason,code_expire_time,start_time,end_time,contacts,contacts_number,ctime,mtime  from vip_resource_pool where id=?"
	_selBusinessByID         = "select id,business_name,business_type,status,app_key,secret,contacts,contacts_number,ctime,mtime from vip_business_info where id=?"
	_selVipAppInfo           = "select id,type,name,app_key,purge_url,ctime,mtime from vip_app_info where type=?"
	_selLastBCoinByMid       = "select id,mid,status,give_now_status,month,amount,memo,ctime,mtime from vip_bcoin_salary where mid = ? order by month desc"
	_selAllConfig            = "select id,config_key,content from vip_config"
	_insertVipBcoinSalary    = "INSERT INTO vip_bcoin_salary(mid,status,give_now_status,month,amount,memo) values(?,?,?,?,?,?)"

	//resouce sql
	_selResourcePoolByIDSQL = "SELECT id,pool_name,business_id,reason,code_expire_time,start_time,end_time,contacts,contacts_number,ctime,mtime  FROM vip_resource_pool WHERE id=?"
	_selBusinessByIDSQL     = "SELECT id,business_name,business_type,status,app_key,secret,contacts,contacts_number,ctime,mtime FROM vip_business_info WHERE id=?"
	_selBusinessByAppkeySQL = "SELECT id,business_name,business_type,status,app_key,secret,contacts,contacts_number,ctime,mtime FROM vip_business_info WHERE app_key=?"
	_selCodeSQL             = "SELECT id,batch_code_id,status,code,mid,use_time,relation_id FROM vip_resource_code WHERE code = ?"
	_selCodesSQL            = "SELECT id,batch_code_id,status,code FROM vip_resource_code WHERE code IN ('%v')"
	_selCodeOpenedSQL       = "SELECT id,code,use_time FROM vip_resource_code WHERE status = 2 AND batch_code_id IN(%v) AND id>%d  AND use_time>='%s' AND use_time<='%s' ORDER BY id ASC LIMIT %d"
	_selBatchCodeSQL        = "SELECT id,business_id,pool_id,status,batch_name,reason,unit,count,surplus_count,price,start_time,end_time,type,limit_day,max_count FROM vip_resource_batch_code WHERE id = ?"
	_selBatchCodeByBisSQL   = "SELECT id,business_id,pool_id,status,batch_name,reason,unit,count,surplus_count,price,start_time,end_time,type,limit_day,max_count FROM vip_resource_batch_code WHERE business_id = ?"
	_selBatchCodesSQL       = "SELECT id,business_id,pool_id,status,batch_name,reason,unit,count,surplus_count,price,start_time,end_time,type,limit_day,max_count FROM vip_resource_batch_code WHERE id IN(%v)"
	_updateCodeSQL          = "UPDATE vip_resource_code SET mid=?,use_time=? WHERE id = ?"
	_updateCodeStatusSQL    = "UPDATE vip_resource_code SET status=? WHERE id = ?"
	_updateBatchCodeSQL     = "UPDATE vip_resource_batch_code SET surplus_count=? WHERE id = ?"
	_selActiveSQL           = "SELECT id,type,product_name,product_pic,relation_id,bus_id,product_detail,use_type FROM vip_active_show WHERE relation_id IN ('%v')"
	_selCodesByBmidSQL      = "SELECT code FROM vip_resource_code WHERE bmid =?"
	_selBatchCountSQL       = "SELECT COUNT(1) FROM vip_resource_code WHERE mid=? AND batch_code_id=?"

	_updateOrderCancelSQL = "UPDATE vip_pay_order SET status = ? WHERE order_no = ? AND status = ?;"
)

//OldUpdateOrderCancel order update order cancel.
func (d *Dao) OldUpdateOrderCancel(c context.Context, r *model.VipPayOrderOld) (a int64, err error) {
	var res xsql.Result
	if res, err = d.olddb.Exec(c, _updateOrderCancelSQL, r.Status, r.OrderNo, model.PAYING); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//OldInsertVipBcoinSalary insert vip bcoin salary
func (d *Dao) OldInsertVipBcoinSalary(c context.Context, r *model.VipBcoinSalary) (err error) {
	if _, err = d.olddb.Exec(c, _insertVipBcoinSalary, &r.Mid, &r.Status, &r.GiveNowStatus, &r.Month, &r.Amount, &r.Memo); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//SelAllConfig sel all config
func (d *Dao) SelAllConfig(c context.Context) (res []*model.VipConfig, err error) {
	var rows *sql.Rows
	if rows, err = d.olddb.Query(c, _selAllConfig); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipConfig)
		if err = rows.Scan(&r.ID, &r.ConfigKey, &r.Content); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

//SelVipAppInfo selVipAppInfo
func (d *Dao) SelVipAppInfo(c context.Context, t int) (res []*model.VipAppInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.olddb.Query(c, _selVipAppInfo, t); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipAppInfo)
		if err = rows.Scan(&r.ID, &r.Type, &r.Name, &r.AppKey, &r.PurgeURL, &r.Ctime, &r.Mtime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

//OldSelLastBcoin sel last bcoin by mid
func (d *Dao) OldSelLastBcoin(c context.Context, mid int64) (r *model.VipBcoinSalary, err error) {
	row := d.olddb.QueryRow(c, _selLastBCoinByMid, mid)
	r = new(model.VipBcoinSalary)
	if err = row.Scan(&r.ID, &r.Mid, &r.Status, &r.GiveNowStatus, &r.Month, &r.Amount, &r.Memo, &r.Ctime, &r.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

//SelBusiness select businessInfo by id
func (d *Dao) SelBusiness(c context.Context, id int64) (r *model.VipBusinessInfo, err error) {
	var row = d.olddb.QueryRow(c, _selBusinessByID, id)
	r = new(model.VipBusinessInfo)
	if err = row.Scan(&r.ID, &r.BusinessName, &r.BusinessType, &r.Status, &r.AppKey, &r.Secret, &r.Contacts, &r.ContactsNumber, &r.Ctime, &r.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.WithStack(err)
	}
	return
}

//SelResourcePool select resource pool by id
func (d *Dao) SelResourcePool(c context.Context, id int64) (r *model.VipResourcePool, err error) {
	var row = d.olddb.QueryRow(c, _selResourcePoolByID, id)
	r = new(model.VipResourcePool)
	if err = row.Scan(&r.ID, &r.PoolName, &r.BusinessID, &r.BusinessName, &r.Reason, &r.CodeExpireTime, &r.StartTime, &r.EndTime, &r.Contacts, &r.ContactsNumber, &r.Ctime, &r.Mtime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
			return
		}
		err = errors.WithStack(err)

	}
	return
}

//OldSelVipUserInfo select user info by mid
func (d *Dao) OldSelVipUserInfo(c context.Context, mid int64) (r *model.VipUserInfo, err error) {
	var row = d.olddb.QueryRow(c, _selVipUserInfoByMid, mid)
	r = new(model.VipUserInfo)
	if err = row.Scan(&r.ID, &r.Mid, &r.VipType, &r.VipStatus, &r.VipStartTime, &r.VipOverdueTime, &r.AnnualVipOverdueTime, &r.Wander, &r.AccessStatus, &r.Ctime, &r.Mtime, &r.VipRecentTime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

//SelVipResourceBatch select vip resource Batch by id
func (d *Dao) SelVipResourceBatch(c context.Context, id int64) (r *model.VipResourceBatch, err error) {
	var row = d.olddb.QueryRow(c, _selVipResourceBatchByID, id)
	r = new(model.VipResourceBatch)
	if err = row.Scan(&r.ID, &r.PoolID, &r.Unit, &r.Count, &r.Ver, &r.StartTime, &r.EndTime, &r.SurplusCount, &r.CodeUseCount, &r.DirectUseCount, &r.Ctime, &r.Mtime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
			return
		}
		err = errors.WithStack(err)
	}
	return
}

//OldUpdateVipUserInfo update vip user info by id
func (d *Dao) OldUpdateVipUserInfo(c context.Context, tx *sql.Tx, r *model.VipUserInfo) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_updateVipUserInfoByID, r.VipType, r.VipStatus, r.VipOverdueTime, r.AnnualVipOverdueTime, r.VipRecentTime, r.ID); err != nil {
		err = errors.WithStack(err)
	} else {
		if a, err = res.RowsAffected(); err != nil {
			err = errors.WithStack(err)
		}
	}
	return
}

//UpdateBatchCount updateBatch by Id
func (d *Dao) UpdateBatchCount(c context.Context, tx *sql.Tx, r *model.VipResourceBatch, ver int64) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_updateBatchCountByID, r.CodeUseCount, r.ID); err != nil {
		err = errors.WithStack(err)
	} else {
		if a, err = res.RowsAffected(); err != nil {
			err = errors.WithStack(err)
		}
	}
	return
}

//OldInsertVipUserInfo insert vipUserInfo.
func (d *Dao) OldInsertVipUserInfo(c context.Context, tx *sql.Tx, r *model.VipUserInfo) (err error) {
	if _, err = tx.Exec(_insertVipUserInfo, r.Mid, r.VipType, r.VipStatus, r.VipStartTime, r.VipOverdueTime, r.AnnualVipOverdueTime, r.Wander, r.AccessStatus, r.VipRecentTime); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//OldInsertVipChangeHistory insert vipChangeHistory
func (d *Dao) OldInsertVipChangeHistory(c context.Context, tx *sql.Tx, r *model.VipChangeHistory) (id int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_insertVipChangeHistory, r.Mid, r.ChangeType, r.ChangeTime, r.Days, r.Month, r.OperatorID, r.RelationID, r.BatchID, r.Remark); err != nil {
		err = errors.WithStack(err)
	} else {
		if id, err = res.LastInsertId(); err != nil {
			err = errors.WithStack(err)
		}
	}
	return
}

//OldVipchangeHistory old vip change history.
func (d *Dao) OldVipchangeHistory(c context.Context, relationID string, batchID int64) (r *model.VipChangeHistory, err error) {
	var row = d.olddb.QueryRow(c, _selVipChangeHistory, relationID, batchID)
	r = new(model.VipChangeHistory)
	if err = row.Scan(&r.ID, &r.Mid, &r.ChangeType, &r.ChangeTime, &r.Days, &r.Month, &r.OperatorID, &r.RelationID, &r.BatchID, &r.Remark, &r.Ctime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
			return
		}
		err = errors.WithStack(err)
		d.errProm.Incr("scan_error")
	}
	return
}
