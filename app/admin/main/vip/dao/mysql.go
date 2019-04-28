package dao

import (
	"context"
	xsql "database/sql"
	"fmt"
	"strconv"
	"time"

	"go-common/app/admin/main/vip/model"
	"go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_SelPoolQueryCount = "SELECT COUNT(1) count from vip_resource_pool WHERE 1=1"
	_SelBusinessByID   = "SELECT id,business_name,business_type,status,app_key,secret,contacts,contacts_number,ctime,mtime FROM vip_business_info WHERE id=?"

	_SelBusinessByQuery = "SELECT id,business_name,business_type,status,app_key,secret,contacts,contacts_number,ctime,mtime FROM vip_business_info WHERE 1=1 "

	_SelPoolByQuery    = "SELECT id,pool_name,business_id,reason,code_expire_time,start_time,end_time,contacts,contacts_number FROM vip_resource_pool WHERE 1=1"
	_SelPoolByName     = "SELECT id,pool_name,business_id,reason,code_expire_time,start_time,end_time,contacts,contacts_number FROM vip_resource_pool WHERE pool_name = ?"
	_SelPoolByID       = "SELECT id,pool_name,business_id,reason,code_expire_time,start_time,end_time,contacts,contacts_number FROM vip_resource_pool WHERE id=?"
	_AddPool           = "INSERT INTO vip_resource_pool(pool_name,business_id,reason,code_expire_time,start_time,end_time,contacts,contacts_number) VALUES(?,?,?,?,?,?,?,?)"
	_UpdatePool        = "UPDATE vip_resource_pool SET pool_name=?,business_id=?,reason=?,code_expire_time=?,start_time=?,end_time=?,contacts=?,contacts_number=? WHERE id = ?"
	_SelBatchByPoolID  = "SELECT id,pool_id,unit,count,ver,start_time,end_time,surplus_count,code_use_count,direct_use_count FROM vip_resource_batch WHERE pool_id=?"
	_SelBatchByID      = "SELECT id,pool_id,unit,count,ver,start_time,end_time,surplus_count,code_use_count,direct_use_count FROM vip_resource_batch WHERE id=?"
	_AddBatch          = "INSERT INTO vip_resource_batch(pool_id,unit,count,ver,start_time,end_time,surplus_count,code_use_count,direct_use_count) VALUES(?,?,?,?,?,?,?,?,?)"
	_UpdateBatch       = "UPDATE vip_resource_batch SET count = ?,ver=?,start_time=?,end_time=?,surplus_count=? WHERE id = ? AND ver=?"
	_UseBatch          = "UPDATE vip_resource_batch SET ver=?,surplus_count = ?,direct_use_count=? WHERE id =? AND ver=?"
	_allVersionSQL     = "SELECT `id`,`platform_id`,`version`,`tip`,`operator`,`link` FROM `vip_app_version`;"
	_updateVersionSQL  = "UPDATE `vip_app_version` SET  %s WHERE `id` = ?;"
	_businessInfosSQL  = "SELECT id,business_name,business_type,status,app_key,contacts,contacts_number,ctime,mtime FROM vip_business_info WHERE 1=1 "
	_businessCountSQL  = "SELECT COUNT(1) FROM vip_business_info WHERE 1=1"
	_addBusinessSQL    = "INSERT INTO `vip_business_info` (`business_name`,`business_type`,`status`,`app_key`,`secret`,`contacts`,`contacts_number`) VALUES (?,?,?,?,?,?,?);"
	_updateBusinessSQL = "UPDATE `vip_business_info` SET `business_name` = ?,`business_type` = ?,`status` = ?,`app_key` = ?,`secret` = ?,`contacts` = ?,`contacts_number` = ? WHERE `id` = ?;"

	_allMonth = "SELECT id,month,month_type,operator,status,mtime FROM vip_month WHERE deleted = 0"
	_getMonth = "SELECT id,month,month_type,operator,status,mtime FROM vip_month WHERE id=?"

	_updateMonthStatus = "UPDATE vip_month SET status=?,operator=? WHERE id=?"
	_allMonthPrice     = "SELECT id,month_id,month_type,money,selected,first_discount_money,discount_money,start_time,end_time,remark,operator FROM vip_month_price WHERE month_id=?"
	_monthPriceSQL     = "SELECT id,month_id,month_type,money,selected,first_discount_money,discount_money,start_time,end_time,remark,operator FROM vip_month_price WHERE id=?"
	_addMonthPrice     = "INSERT INTO vip_month_price (month_id,month_type,money,first_discount_money,discount_money,start_time,end_time,remark,operator) VALUES (?,?,?,?,?,?,?,?,?)"
	_editMonthPrice    = "UPDATE vip_month_price SET month_type=?,money=?,first_discount_money=?,discount_money=?,start_time=?,end_time=?,remark=?,operator=? WHERE id=?"

	//resouce SQL
	_addBatchCodeSQL       = "INSERT INTO vip_resource_batch_code(business_id,pool_id,status,batch_name,reason,unit,count,surplus_count,price,start_time,end_time,contacts,contacts_number,type,limit_day,max_count,operator) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_updateBatchCodeSQL    = "UPDATE vip_resource_batch_code SET status=?,batch_name=?,reason=?,price=?,contacts=?,contacts_number=?,type=?,limit_day=?,max_count=?,operator=? WHERE id=?"
	_selBatchCodeIDSQL     = "SELECT id FROM vip_resource_batch_code WHERE  1=1"
	_selBatchCodeSQL       = "SELECT id,business_id,pool_id,status,type,limit_day,max_count,batch_name,reason,unit,count,surplus_count,price,start_time,end_time,contacts,contacts_number,ctime FROM vip_resource_batch_code WHERE  1=1 "
	_selBatchCodeByIDSQL   = "SELECT id,business_id,pool_id,status,type,limit_day,max_count,batch_name,reason,unit,count,surplus_count,price,start_time,end_time,contacts,contacts_number,ctime FROM vip_resource_batch_code WHERE  id = ?"
	_selBatchCodeCountSQL  = "SELECT COUNT(1) FROM vip_resource_batch_code WHERE  1=1 "
	_selBatchCodeByNameSQL = "SELECT id,business_id,pool_id,status,type,limit_day,max_count,batch_name,reason,unit,count,surplus_count,price,start_time,end_time,contacts,contacts_number,ctime FROM vip_resource_batch_code WHERE  batch_name = ?"
	_selBatchCodesSQL      = "SELECT id,business_id,pool_id,status,type,limit_day,max_count,batch_name,reason,unit,count,surplus_count,price,start_time,end_time,contacts,contacts_number,ctime FROM vip_resource_batch_code WHERE  1=1 AND id in (%v)"

	_batchAddCodeSQL = "INSERT INTO vip_resource_code(batch_code_id,status,code) VALUES"
	_updateCodeSQL   = "UPDATE vip_resource_code SET status=? WHERE id = ?"
	_selCodeSQL      = "SELECT id,batch_code_id,status,code,mid,use_time,ctime FROM vip_resource_code WHERE 1=1"
	_selCodeByIDSQL  = "SELECT id,batch_code_id,status,code,mid,use_time,ctime FROM vip_resource_code WHERE id = ?"

	//pushData
	_selPushDataCountSQL = "SELECT COUNT(1) FROM vip_push_data WHERE 1=1 "
	_selPushDataSQL      = "SELECT id,disable_type,`group_name`,title,content,push_total_count,pushed_count,progress_status,`status`,platform,link_type,link_url,error_code,expired_day_start,expired_day_end,effect_start_date,effect_end_date,push_start_time,push_end_time,operator FROM vip_push_data WHERE 1=1 "
	_selPushDataByIDSQL  = "SELECT id,disable_type,`group_name`,title,content,push_total_count,pushed_count,progress_status,`status`,platform,link_type,link_url,error_code,expired_day_start,expired_day_end,effect_start_date,effect_end_date,push_start_time,push_end_time,operator FROM vip_push_data WHERE id=?"
	_addPushDataSQL      = "INSERT INTO vip_push_data(`group_name`,title,content,push_total_count,pushed_count,progress_status,`status`,platform,link_type,link_url,expired_day_start,expired_day_end,effect_start_date,effect_end_date,push_start_time,push_end_time,operator) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_updatePushDataSQL   = "UPDATE vip_push_data SET  group_name=?,title=?,content=?,push_total_count=?,progress_status=?,platform=?,link_type=?,link_url=?,expired_day_start=?,expired_day_end=?,effect_start_date=?,effect_end_date=?,push_start_time=?,push_end_time=?,operator=? WHERE id=?"
	_delPushDataSQL      = "DELETE FROM vip_push_data WHERE id = ?"
	_disablePushDataSQL  = "UPDATE vip_push_data SET disable_type=1,progress_status=?,push_total_count=?,effect_end_date=? WHERE id=?"

	//order
	_vipOrderListSQL  = "SELECT id,order_no,app_id,platform,order_type,mid,to_mid,buy_months,money,refund_amount,status,pay_type,recharge_bp,third_trade_no,ver,payment_time,ctime,mtime,app_sub_id FROM vip_pay_order WHERE 1=1"
	_vipOrderCountSQL = "SELECT COUNT(1) FROM vip_pay_order WHERE 1=1"
	_vipOrderSQL      = "SELECT id,order_no,app_id,platform,order_type,mid,to_mid,buy_months,money,refund_amount,status,pay_type,recharge_bp,third_trade_no,ver,payment_time,ctime,mtime,app_sub_id FROM vip_pay_order WHERE order_no = ? "
	_addOrderLogSQL   = "INSERT INTO vip_pay_order_log(order_no,refund_id,refund_amount,operator,mid,status) VALUES(?,?,?,?,?,?)"

	_getUserChangeHistorySQL      = "SELECT id,mid,change_type,change_time,days,operator_id,relation_id,batch_id,remark,ctime,mtime FROM vip_user_change_history WHERE 1=1 "
	_getUserChangeHistoryCountSQL = "SELECT COUNT(1) FROM `vip_user_change_history` WHERE 1=1 "
	_defpn                        = 1
	_defps                        = 20
	_maxps                        = 100
)

func (d *Dao) joinPoolCondition(sql string, q *model.ResoucePoolBo, pn, ps int) string {
	if len(q.PoolName) > 0 {
		sql += " and pool_name like '%" + q.PoolName + "%'"
	}
	if q.BusinessID > 0 {
		sql += " and business_id = " + strconv.Itoa(q.BusinessID)
	}
	if q.StartTime > 0 {
		sql += " and start_time >= '" + q.StartTime.Time().Format("2006-01-02 15:04:05") + "'"
	}
	if q.EndTime > 0 {
		sql += " and end_time <='" + q.EndTime.Time().Format("2006-01-02 15:04:05") + "'"
	}
	if q.ID > 0 || q.ID == -1 {
		sql += " and id = " + strconv.Itoa(q.ID)
	}
	if pn > 0 {
		if ps <= 0 {
			ps = 20
		}
		offer := (q.PN - 1) * q.PS
		sql += " limit " + strconv.Itoa(offer) + ", " + strconv.Itoa(q.PS)
	}
	return sql
}

func (d *Dao) joinHistoryCondition(sql string, u *model.UserChangeHistoryReq, iscount bool) string {
	if u.Mid > 0 {
		sql += " AND mid = " + fmt.Sprintf("%d", u.Mid)
	}
	if u.ChangeType > 0 {
		sql += " AND change_type = " + fmt.Sprintf("%d", u.ChangeType)
	}
	if u.StartChangeTime > 0 {
		stime := time.Unix(u.StartChangeTime, 0).Format("2006-01-02 15:04:05")
		sql += " AND change_time >= '" + stime + "'"
	}
	if u.EndChangeTime > 0 {
		etime := time.Unix(u.EndChangeTime, 0).Format("2006-01-02 15:04:05")
		sql += " AND change_time <= '" + etime + "'"
	}
	if u.BatchID > 0 {
		sql += " AND batch_id = " + fmt.Sprintf("%d", u.BatchID)
	}
	if len(u.RelationID) > 0 {
		sql += " AND relation_id = " + fmt.Sprintf("'%v'", u.RelationID)
	}
	if iscount {
		return sql
	}
	if u.Pn == 0 {
		u.Pn = _defpn
	}
	if u.Ps == 0 {
		u.Ps = _defps
	}
	offer := (u.Pn - 1) * u.Ps
	sql += " ORDER BY id DESC LIMIT " + strconv.Itoa(offer) + ", " + strconv.Itoa(u.Ps)
	return sql
}

// SelPoolByName sel pool by name
func (d *Dao) SelPoolByName(c context.Context, name string) (r *model.VipResourcePool, err error) {
	var row *sql.Row
	if row = d.db.QueryRow(c, _SelPoolByName, name); err != nil {
		log.Error("SelPoolByName db.query() error(%v)", err)
		return
	}
	r = new(model.VipResourcePool)
	if err = row.Scan(&r.ID, &r.PoolName, &r.BusinessID, &r.Reason, &r.CodeExpireTime, &r.StartTime, &r.EndTime, &r.Contacts, &r.ContactsNumber); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
		return
	}
	return
}

// SelCountPool sel count Pool
func (d *Dao) SelCountPool(c context.Context, r *model.ResoucePoolBo) (count int, err error) {
	var row *sql.Row
	if row = d.db.QueryRow(c, d.joinPoolCondition(_SelPoolQueryCount, r, 0, 0)); err != nil {
		log.Error("SelCountPool db.query() error(%v)", err)
		return
	}
	if err = row.Scan(&count); err != nil {
		log.Error("row.scan() error(%v)", err)
		return
	}
	return
}

// SelPool sel pool by query condition
func (d *Dao) SelPool(c context.Context, r *model.ResoucePoolBo, pn, ps int) (res []*model.VipResourcePool, err error) {

	var rows *sql.Rows
	if rows, err = d.db.Query(c, d.joinPoolCondition(_SelPoolByQuery, r, pn, ps)); err != nil {
		log.Error("SelPool db.query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipResourcePool)
		if err = rows.Scan(&r.ID, &r.PoolName, &r.BusinessID, &r.Reason, &r.CodeExpireTime, &r.StartTime, &r.EndTime, &r.Contacts, &r.ContactsNumber); err != nil {
			log.Error("row.scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)

	}
	return
}

// SelPoolRow sel pool by id
func (d *Dao) SelPoolRow(c context.Context, id int) (r *model.VipResourcePool, err error) {

	var row *sql.Row
	if row = d.db.QueryRow(c, _SelPoolByID, id); err != nil {
		log.Error("SelPoolRow db.query() error(%v)", err)
		return
	}
	r = new(model.VipResourcePool)
	if err = row.Scan(&r.ID, &r.PoolName, &r.BusinessID, &r.Reason, &r.CodeExpireTime, &r.StartTime, &r.EndTime, &r.Contacts, &r.ContactsNumber); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
		return
	}
	return
}

// AddPool add pool
func (d *Dao) AddPool(c context.Context, r *model.ResoucePoolBo) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _AddPool, r.PoolName, r.BusinessID, r.Reason, r.CodeExpireTime, r.StartTime, r.EndTime, r.Contacts, r.ContactsNumber); err != nil {
		log.Error("AddPool d.db.exec(%v) error(%v)", r, err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		log.Error("AddPool RowsAffected() error(%v)", err)
		return
	}
	return
}

// UpdatePool update pool
func (d *Dao) UpdatePool(c context.Context, r *model.ResoucePoolBo) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _UpdatePool, r.PoolName, r.BusinessID, r.Reason, r.CodeExpireTime, r.StartTime, r.EndTime, r.Contacts, r.ContactsNumber, r.ID); err != nil {
		log.Error("UpdatePool d.db.exec(%v) error(%v)", r, err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		log.Error("UpdatePool RowsAffected() error(%v)", err)
		return
	}
	return
}

// SelBatchRow sel batch by id
func (d *Dao) SelBatchRow(c context.Context, id int) (r *model.VipResourceBatch, err error) {
	var row *sql.Row
	if row = d.db.QueryRow(c, _SelBatchByID, id); err != nil {
		log.Error("SelBatchRow db.query() error(%v)", err)
		return
	}
	r = new(model.VipResourceBatch)
	if err = row.Scan(&r.ID, &r.PoolID, &r.Unit, &r.Count, &r.Ver, &r.StartTime, &r.EndTime, &r.SurplusCount, &r.CodeUseCount, &r.DirectUseCount); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
		return
	}
	return
}

// SelBatchRows sel batch by poolID
func (d *Dao) SelBatchRows(c context.Context, poolID int) (res []*model.VipResourceBatch, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _SelBatchByPoolID, poolID); err != nil {
		log.Error("SelBatchRows db.query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipResourceBatch)
		if err = rows.Scan(&r.ID, &r.PoolID, &r.Unit, &r.Count, &r.Ver, &r.StartTime, &r.EndTime, &r.SurplusCount, &r.CodeUseCount, &r.DirectUseCount); err != nil {
			log.Error("row.scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)

	}
	return
}

// AddBatch add batch
func (d *Dao) AddBatch(c context.Context, r *model.ResouceBatchBo) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _AddBatch, r.PoolID, r.Unit, r.Count, 0, r.StartTime, r.EndTime, r.SurplusCount, r.CodeUseCount, r.DirectUseCount); err != nil {
		log.Error("AddBatch d.db.exec(%v) error(%v)", r, err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		log.Error("AddBatch RowsAffected() error(%v)", err)
		return
	}
	return
}

// UpdateBatch update batch data
func (d *Dao) UpdateBatch(c context.Context, r *model.VipResourceBatch, ver int) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _UpdateBatch, r.Count, r.Ver, r.StartTime, r.EndTime, r.SurplusCount, r.ID, ver); err != nil {
		log.Error("UpdateBatch d.db.exec(%v) error(%v)", r, err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		log.Error("UpdateBatch RowsAffected() error(%v)", err)
		return
	}
	return
}

// UseBatch use batch resouce
func (d *Dao) UseBatch(tx *sql.Tx, r *model.VipResourceBatch, ver int) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_UseBatch, r.Ver, r.SurplusCount, r.DirectUseCount, r.ID, ver); err != nil {
		log.Error("UseBatch d.db.exec(%v) error(%v)", r, err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		log.Error("UseBatch RowsAffected() error(%v)", err)
		return
	}
	return
}

// SelBusiness select businessInfo by id
func (d *Dao) SelBusiness(c context.Context, id int) (r *model.VipBusinessInfo, err error) {
	var row = d.db.QueryRow(c, _SelBusinessByID, id)
	r = new(model.VipBusinessInfo)
	if err = row.Scan(&r.ID, &r.BusinessName, &r.BusinessType, &r.Status, &r.AppKey, &r.Secret, &r.Contacts, &r.ContactsNumber, &r.Ctime, &r.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// SelBusinessByQuery .
func (d *Dao) SelBusinessByQuery(c context.Context, arg *model.QueryBusinessInfo) (r *model.VipBusinessInfo, err error) {
	queryStr := ""
	if len(arg.Name) > 0 {
		queryStr += fmt.Sprintf(" AND business_name = '%v' ", arg.Name)
	}
	if len(arg.Appkey) > 0 {
		queryStr += fmt.Sprintf(" AND app_key = '%v' ", arg.Appkey)
	}
	var row = d.db.QueryRow(c, _SelBusinessByQuery+queryStr)
	r = new(model.VipBusinessInfo)
	if err = row.Scan(&r.ID, &r.BusinessName, &r.BusinessType, &r.Status, &r.AppKey, &r.Secret, &r.Contacts, &r.ContactsNumber, &r.Ctime, &r.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// AllVersion all version.
func (d *Dao) AllVersion(c context.Context) (res []*model.VipAppVersion, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _allVersionSQL); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipAppVersion)
		if err = rows.Scan(&r.ID, &r.PlatformID, &r.Version, &r.Tip, &r.Operator, &r.Link); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// UpdateVersion update version.
func (d *Dao) UpdateVersion(c context.Context, v *model.VipAppVersion) (ret int64, err error) {
	var (
		sql string
		res xsql.Result
	)
	if len(v.Version) > 0 {
		sql += "`version` = '" + v.Version + "',"
	}
	if len(v.Tip) > 0 {
		sql += "`tip` = '" + v.Tip + "',"
	}
	if len(v.Link) > 0 {
		sql += "`link` = '" + v.Link + "',"
	}
	sql += "`operator` = '" + v.Operator + "'"
	if res, err = d.db.Exec(c, fmt.Sprintf(_updateVersionSQL, sql), v.ID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if ret, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// BussinessList business infos.
func (d *Dao) BussinessList(c context.Context, pn, ps, status int) (res []*model.VipBusinessInfo, err error) {
	var rows *sql.Rows
	sqlStr := _businessInfosSQL
	if status >= 0 {
		sqlStr += fmt.Sprintf(" AND status=%v", status)
	}
	if pn <= 0 {
		pn = _defpn
	}
	if pn <= 0 || pn > _maxps {
		ps = _defps
	}
	sqlStr += fmt.Sprintf(" ORDER BY id DESC limit %v,%v", (pn-1)*ps, ps)
	if rows, err = d.db.Query(c, sqlStr); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipBusinessInfo)
		if err = rows.Scan(&r.ID, &r.BusinessName, &r.BusinessType, &r.Status, &r.AppKey,
			&r.Contacts, &r.ContactsNumber, &r.Ctime, &r.Mtime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// BussinessCount bussiness info count.
func (d *Dao) BussinessCount(c context.Context, status int) (count int64, err error) {
	var row *sql.Row
	sqlstr := _businessCountSQL
	if status >= 0 {
		sqlstr += fmt.Sprintf(" AND status=%v", status)
	}
	if row = d.db.QueryRow(c, sqlstr); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("query_db")
		return
	}
	if err = row.Scan(&count); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("row_scan_db")
		return
	}
	return
}

// UpdateBusiness update business info.
func (d *Dao) UpdateBusiness(c context.Context, r *model.VipBusinessInfo) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _updateBusinessSQL, r.BusinessName, r.BusinessType, r.Status, r.AppKey, r.Secret, r.Contacts, r.ContactsNumber, r.ID); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("exec_db")
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("rows_affected_db")
		return
	}
	return
}

// AddBusiness add business info.
func (d *Dao) AddBusiness(c context.Context, r *model.VipBusinessInfo) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _addBusinessSQL, r.BusinessName, r.BusinessType, r.Status, r.AppKey, r.Secret, r.Contacts, r.ContactsNumber); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("exec_db")
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("rows_affected_db")
		return
	}
	return
}

// HistoryCount user change history.
func (d *Dao) HistoryCount(c context.Context, u *model.UserChangeHistoryReq) (count int, err error) {
	var row *sql.Row
	if row = d.db.QueryRow(c, d.joinHistoryCondition(_getUserChangeHistoryCountSQL, u, true)); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("query_db")
		return
	}
	if err = row.Scan(&count); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("row_scan_db")
		return
	}
	return
}

// HistoryList history list.
func (d *Dao) HistoryList(c context.Context, u *model.UserChangeHistoryReq) (res []*model.VipChangeHistory, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, d.joinHistoryCondition(_getUserChangeHistorySQL, u, false)); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("query_db")
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipChangeHistory)
		if err = rows.Scan(&r.ID, &r.Mid, &r.ChangeType, &r.ChangeTime, &r.Days, &r.OperatorID, &r.RelationID, &r.BatchID, &r.Remark, &r.Ctime, &r.Mtime); err != nil {
			err = errors.WithStack(err)
			d.errProm.Incr("rows_scan_db")
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}
