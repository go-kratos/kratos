package dao

import (
	"bytes"
	"context"
	xsql "database/sql"
	"fmt"
	"strconv"
	"time"

	"go-common/app/service/main/coupon/model"
	"go-common/library/database/sql"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_couponsSQL                = "SELECT id,coupon_token,mid,state,start_time,expire_time,origin,coupon_type,order_no,oid,remark,use_ver,ver,ctime,mtime FROM coupon_info_%02d WHERE mid = ? AND state = ? AND coupon_type = ? AND expire_time > ? AND start_time < ? ORDER BY expire_time;"
	_couponsNoStartCheckSQL    = "SELECT id,coupon_token,mid,state,start_time,expire_time,origin,coupon_type,order_no,oid,remark,use_ver,ver,ctime,mtime FROM coupon_info_%02d WHERE mid = ? AND state = ? AND coupon_type = ? AND expire_time > ? ORDER BY expire_time;"
	_couponByOrderNOAndTypeSQL = "SELECT id,coupon_token,mid,state,start_time,expire_time,origin,coupon_type,order_no,oid,remark,use_ver,ver,ctime,mtime FROM coupon_info_%02d WHERE order_no = ?  AND coupon_type = ?;"
	_updateCouponInUseSQL      = "UPDATE coupon_info_%02d SET state =?, order_no = ?, oid = ?, remark = ?,use_ver = ?,ver = ? WHERE coupon_token = ? AND ver = ?;"
	_addCouponChangeLogSQL     = "INSERT INTO coupon_change_log_%02d (coupon_token,mid,state,ctime) VALUES(?,?,?,?);"
	_couponByTokenSQL          = "SELECT id,coupon_token,mid,state,start_time,expire_time,origin,coupon_type,order_no,oid,remark,use_ver,ver,ctime,mtime FROM coupon_info_%02d WHERE coupon_token = ?;"
	_couponPageNotUsedSQL      = "SELECT id,coupon_token,mid,state,start_time,expire_time,origin,coupon_type,order_no,oid,remark,use_ver,ver,ctime,mtime FROM coupon_info_%02d WHERE mid = ? AND state = ? AND expire_time > ? AND start_time < ? AND ctime > ? ORDER BY mtime DESC LIMIT ?,?;"
	_couponPageUsedSQL         = "SELECT id,coupon_token,mid,state,start_time,expire_time,origin,coupon_type,order_no,oid,remark,use_ver,ver,ctime,mtime FROM coupon_info_%02d WHERE mid = ? AND state = ? AND ctime > ?  ORDER BY mtime DESC LIMIT ?,?;"
	_couponPageExpireSQL       = "SELECT id,coupon_token,mid,state,start_time,expire_time,origin,coupon_type,order_no,oid,remark,use_ver,ver,ctime,mtime FROM coupon_info_%02d WHERE mid = ? AND state = 0 AND expire_time < ? AND ctime > ? ORDER BY mtime DESC LIMIT ?,?;"
	_countNotUsedSQL           = "SELECT COUNT(1) FROM coupon_info_%02d WHERE mid = ? AND state = ? AND expire_time > ? AND start_time < ? AND ctime > ? ;"
	_countUsedSQL              = "SELECT COUNT(1) FROM coupon_info_%02d WHERE mid = ? AND state = ? AND ctime > ?  ;"
	_countExpireSQL            = "SELECT COUNT(1) FROM coupon_info_%02d WHERE mid = ? AND state = 0 AND expire_time < ? AND ctime > ? ;"
	_addCouponSQL              = "INSERT INTO coupon_info_%02d (coupon_token,mid,state,start_time,expire_time,origin,coupon_type,ctime)VALUES(?,?,?,?,?,?,?,?);"
	_updateStateSQL            = "UPDATE coupon_info_%02d SET state = ?,use_ver = ?,ver = ? WHERE coupon_token = ? AND ver = ? "
	_batchAddCouponSQL         = "INSERT INTO coupon_info_%02d (coupon_token,mid,state,start_time,expire_time,origin,coupon_type,ctime,batch_token)VALUES  "
	_batchCountByMid           = "SELECT  COUNT(1) FROM  coupon_info_%02d WHERE mid = ? AND batch_token = ?;"

	//coupon blance
	_couponBlanceNoStartCheckSQL = "SELECT id,batch_token,mid,balance,start_time,expire_time,origin,coupon_type,ver,ctime,mtime FROM coupon_balance_info_%02d WHERE mid = ?  AND expire_time > ? AND coupon_type = ? ORDER BY expire_time;"
	_couponBlanceSQL             = "SELECT id,batch_token,mid,balance,start_time,expire_time,origin,coupon_type,ver,ctime,mtime FROM coupon_balance_info_%02d WHERE mid = ?  AND expire_time > ? AND coupon_type = ? AND start_time < ? ORDER BY expire_time;"
	_orderByThirdTradeNoSQL      = "SELECT id,order_no,mid,count,state,coupon_type,third_trade_no,remark,tips,use_ver,ver,ctime,mtime FROM coupon_order WHERE  third_trade_no = ? AND coupon_type= ?;"
	_updateBlanceSQL             = "UPDATE coupon_balance_info_%02d SET balance = ?,ver = ver + 1 WHERE id = ? AND ver = ?;"
	_addOrderSQL                 = "INSERT INTO coupon_order(order_no,mid,count,state,coupon_type,third_trade_no,remark,tips,use_ver,ver,ctime)VALUES(?,?,?,?,?,?,?,?,?,?,?);"
	_addOrderLogSQL              = "INSERT INTO coupon_order_log(order_no,mid,state,ctime)VALUES(?,?,?,?);"
	_addBalanceLogSQL            = "INSERT INTO coupon_balance_change_log_%02d(order_no,mid,batch_token,balance,change_balance,change_type,ctime)VALUES "
	_batchUpdateBalance          = "UPDATE coupon_balance_info_%02d SET ver =ver + 1, balance = CASE id"
	_countBalanceNotUsed         = "SELECT COUNT(1) FROM coupon_balance_info_%02d  WHERE mid = ? AND  expire_time > ? AND start_time < ? AND coupon_type = ? AND  balance > 0 AND ctime > ? ;"
	_countUseListSQL             = "SELECT COUNT(1) FROM coupon_order WHERE mid= ? AND state= ? AND coupon_type = ? AND ctime > ? ;"
	_countBalanceExpire          = "SELECT COUNT(1) FROM coupon_balance_info_%02d  WHERE mid = ? AND  expire_time < ?  AND coupon_type = ? AND  balance > 0 AND ctime > ? ;"
	_balanceNotUsedPageSQL       = "SELECT id,batch_token,mid,balance,start_time,expire_time,origin,coupon_type,ver,ctime,mtime FROM coupon_balance_info_%02d  WHERE mid = ? AND  expire_time > ? AND start_time < ? AND coupon_type = ? AND  balance > 0 AND ctime > ? ORDER BY id DESC LIMIT ?,?;"
	_useOrderPageSQL             = "SELECT id,order_no,mid,count,state,coupon_type,third_trade_no,remark,tips,use_ver,ver,ctime,mtime FROM coupon_order WHERE mid= ? AND state= ? AND coupon_type = ? AND ctime > ?  ORDER BY id DESC LIMIT ?,?;"
	_balanceExpirePageSQL        = "SELECT id,batch_token,mid,balance,start_time,expire_time,origin,coupon_type,ver,ctime,mtime FROM coupon_balance_info_%02d  WHERE mid = ? AND  expire_time < ?  AND coupon_type = ? AND  balance > 0 AND ctime > ?  ORDER BY id DESC LIMIT ?,?;"
	_addBalanceCouponSQL         = "INSERT INTO coupon_balance_info_%02d(batch_token,mid,balance,start_time,expire_time,origin,coupon_type,ver,ctime) VALUES(?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE balance = balance + ?,ver = ver + 1 ;"
	_byMidAndBatchTokenSQL       = "SELECT id,batch_token,mid,balance,start_time,expire_time,origin,coupon_type,ver,ctime,mtime FROM coupon_balance_info_%02d WHERE mid = ?  AND batch_token = ? ;"
	_addBlanceChangeLog          = "INSERT INTO coupon_balance_change_log_%02d(order_no,mid,batch_token,balance,change_balance,change_type,ctime)VALUES(?,?,?,?,?,?,?);"
	_batchInfoSQL                = "SELECT id,app_id,name,batch_token,max_count,current_count,start_time,expire_time,expire_day,ver,ctime,limit_count,full_amount,amount,state,coupon_type FROM coupon_batch_info WHERE batch_token = ?;"
	_updateBatchSQL              = "UPDATE coupon_batch_info SET current_count = current_count + ? WHERE batch_token = ?;"
	_updateBatchLimitSQL         = "UPDATE coupon_batch_info SET current_count = current_count + ? WHERE batch_token = ? AND current_count + ? <=  max_count;"
	_grantCouponLogSQL           = "SELECT id,order_no,mid,batch_token,balance,change_balance,change_type,ctime,mtime FROM coupon_balance_change_log_%02d WHERE mid = ? AND batch_token = ? AND change_type = ?;"
	_allBatchInfoSQL             = "SELECT id,app_id,name,batch_token,max_count,current_count,start_time,expire_time,expire_day,ver,ctime,limit_count,full_amount,amount,state,coupon_type,platform_limit,product_limit_month,product_limit_renewal FROM coupon_batch_info;"
	_couponReceiveSQL            = "SELECT id,appkey,order_no,mid,coupon_token,coupon_type FROM coupon_receive_log WHERE order_no=? AND appkey=? AND coupon_type=?"
	_addReceiveSQL               = "INSERT INTO coupon_receive_log(appkey,order_no,mid,coupon_token,coupon_type) VALUES(?,?,?,?,?)"
)

func hitInfo(mid int64) int64 {
	return mid % 100
}

func hitChangeLog(mid int64) int64 {
	return mid % 100
}

func hitUser(mid int64) int64 {
	return mid % 10
}

func hitUserLog(mid int64) int64 {
	return mid % 10
}

// BeginTran begin transaction.
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.db.Begin(c)
}

// CouponList query .
func (d *Dao) CouponList(c context.Context, mid int64, state int8, ct int8, t int64) (res []*model.CouponInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_couponsSQL, hitInfo(mid)), mid, state, ct, t, t); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.CouponInfo{}
		if err = rows.Scan(&r.ID, &r.CouponToken, &r.Mid, &r.State, &r.StartTime, &r.ExpireTime, &r.Origin, &r.CouponType, &r.OrderNO, &r.Oid, &r.Remark,
			&r.UseVer, &r.Ver, &r.CTime, &r.MTime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// CouponNoStartCheckList had not check start query .
func (d *Dao) CouponNoStartCheckList(c context.Context, mid int64, state int8, ct int8, t int64) (res []*model.CouponInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_couponsNoStartCheckSQL, hitInfo(mid)), mid, state, ct, t); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.CouponInfo{}
		if err = rows.Scan(&r.ID, &r.CouponToken, &r.Mid, &r.State, &r.StartTime, &r.ExpireTime, &r.Origin, &r.CouponType, &r.OrderNO, &r.Oid, &r.Remark,
			&r.UseVer, &r.Ver, &r.CTime, &r.MTime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// BlanceNoStartCheckList had not check start query .
func (d *Dao) BlanceNoStartCheckList(c context.Context, mid int64, ct int8, t int64) (res []*model.CouponBalanceInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_couponBlanceNoStartCheckSQL, hitUser(mid)), mid, t, ct); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.CouponBalanceInfo{}
		if err = rows.Scan(&r.ID, &r.BatchToken, &r.Mid, &r.Balance, &r.StartTime, &r.ExpireTime, &r.Origin, &r.CouponType, &r.Ver, &r.CTime, &r.MTime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// ByOrderNO query coupon by orderno and type.
func (d *Dao) ByOrderNO(c context.Context, mid int64, orderNO string, ct int8) (r *model.CouponInfo, err error) {
	var row *sql.Row
	r = &model.CouponInfo{}
	row = d.db.QueryRow(c, fmt.Sprintf(_couponByOrderNOAndTypeSQL, hitInfo(mid)), orderNO, ct)
	if err = row.Scan(&r.ID, &r.CouponToken, &r.Mid, &r.State, &r.StartTime, &r.ExpireTime, &r.Origin, &r.CouponType, &r.OrderNO, &r.Oid, &r.Remark,
		&r.UseVer, &r.Ver, &r.CTime, &r.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// UpdateCouponInUse update coupon in use.
func (d *Dao) UpdateCouponInUse(c context.Context, tx *sql.Tx, cp *model.CouponInfo) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_updateCouponInUseSQL, hitInfo(cp.Mid)), cp.State, cp.OrderNO, cp.Oid, cp.Remark, cp.UseVer, cp.Ver+1,
		cp.CouponToken, cp.Ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//InsertPointHistory .
func (d *Dao) InsertPointHistory(c context.Context, tx *sql.Tx, l *model.CouponChangeLog) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_addCouponChangeLogSQL, hitChangeLog(l.Mid)), l.CouponToken, l.Mid, l.State, l.Ctime); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// CouponInfo coupon info.
func (d *Dao) CouponInfo(c context.Context, mid int64, token string) (r *model.CouponInfo, err error) {
	var row *sql.Row
	r = &model.CouponInfo{}
	row = d.db.QueryRow(c, fmt.Sprintf(_couponByTokenSQL, hitInfo(mid)), token)
	if err = row.Scan(&r.ID, &r.CouponToken, &r.Mid, &r.State, &r.StartTime, &r.ExpireTime, &r.Origin, &r.CouponType, &r.OrderNO, &r.Oid, &r.Remark,
		&r.UseVer, &r.Ver, &r.CTime, &r.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// CountByState coupon count buy state.
func (d *Dao) CountByState(c context.Context, mid int64, state int8, t int64, stime time.Time) (count int64, err error) {
	var row *sql.Row
	switch state {
	case model.NotUsed:
		row = d.db.QueryRow(c, fmt.Sprintf(_countNotUsedSQL, hitInfo(mid)), mid, state, t, t, stime)
	case model.Used:
		row = d.db.QueryRow(c, fmt.Sprintf(_countUsedSQL, hitInfo(mid)), mid, state, stime)
	case model.Expire:
		row = d.db.QueryRow(c, fmt.Sprintf(_countExpireSQL, hitInfo(mid)), mid, t, stime)
	default:
		return
	}
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

// CouponPage page.
func (d *Dao) CouponPage(c context.Context, mid int64, state int8, t int64, start int, ps int, stime time.Time) (res []*model.CouponInfo, err error) {
	var rows *sql.Rows
	switch state {
	case model.NotUsed:
		rows, err = d.db.Query(c, fmt.Sprintf(_couponPageNotUsedSQL, hitInfo(mid)), mid, state, t, t, stime, start, ps)
	case model.Used:
		rows, err = d.db.Query(c, fmt.Sprintf(_couponPageUsedSQL, hitInfo(mid)), mid, state, stime, start, ps)
	case model.Expire:
		rows, err = d.db.Query(c, fmt.Sprintf(_couponPageExpireSQL, hitInfo(mid)), mid, t, stime, start, ps)
	default:
		return
	}
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.CouponInfo{}
		if err = rows.Scan(&r.ID, &r.CouponToken, &r.Mid, &r.State, &r.StartTime, &r.ExpireTime, &r.Origin, &r.CouponType, &r.OrderNO, &r.Oid, &r.Remark,
			&r.UseVer, &r.Ver, &r.CTime, &r.MTime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// AddCoupon add coupon.
func (d *Dao) AddCoupon(c context.Context, cp *model.CouponInfo) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_addCouponSQL, hitInfo(cp.Mid)), cp.CouponToken, cp.Mid, cp.State, cp.StartTime, cp.ExpireTime, cp.Origin, cp.CouponType, cp.CTime); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// BatchAddCoupon batch add coupon.
func (d *Dao) BatchAddCoupon(c context.Context, tx *sql.Tx, mid int64, cps []*model.CouponInfo) (a int64, err error) {
	var (
		buf bytes.Buffer
		res xsql.Result
		sql string
	)
	buf.WriteString(fmt.Sprintf(_batchAddCouponSQL, hitInfo(mid)))
	for _, v := range cps {
		buf.WriteString("('")
		buf.WriteString(v.CouponToken)
		buf.WriteString("',")
		buf.WriteString(strconv.FormatInt(v.Mid, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(v.State, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(v.StartTime, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(v.ExpireTime, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(v.Origin, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(v.CouponType, 10))
		buf.WriteString(",'")
		buf.WriteString(v.CTime.Time().Format("2006-01-02 15:04:05"))
		buf.WriteString("','")
		buf.WriteString(v.BatchToken)
		buf.WriteString("'),")
	}
	sql = buf.String()
	if res, err = tx.Exec(sql[0 : len(sql)-1]); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// UpdateCoupon update coupon in use.
func (d *Dao) UpdateCoupon(c context.Context, mid int64, state int8, useVer int64, ver int64, couponToken string) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_updateStateSQL, hitInfo(mid)), state, useVer, ver+1, couponToken, ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// ByThirdTradeNo query order by third trade no.
func (d *Dao) ByThirdTradeNo(c context.Context, thirdTradeNo string, ct int8) (r *model.CouponOrder, err error) {
	var row *sql.Row
	r = &model.CouponOrder{}
	row = d.db.QueryRow(c, _orderByThirdTradeNoSQL, thirdTradeNo, ct)
	if err = row.Scan(&r.ID, &r.OrderNo, &r.Mid, &r.Count, &r.State, &r.CouponType, &r.ThirdTradeNo, &r.Remark, &r.Tips, &r.UseVer, &r.Ver, &r.Ctime, &r.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// CouponBlances query coupon blances .
func (d *Dao) CouponBlances(c context.Context, mid int64, ct int8, t int64) (res []*model.CouponBalanceInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_couponBlanceSQL, hitUser(mid)), mid, t, ct, t); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.CouponBalanceInfo{}
		if err = rows.Scan(&r.ID, &r.BatchToken, &r.Mid, &r.Balance, &r.StartTime, &r.ExpireTime, &r.Origin, &r.CouponType, &r.Ver, &r.CTime, &r.MTime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// UpdateBlance update blance.
func (d *Dao) UpdateBlance(c context.Context, tx *sql.Tx, id int64, mid int64, ver int64, balance int64) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_updateBlanceSQL, hitUser(mid)), balance, id, ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// BatchUpdateBlance batch update blance.
func (d *Dao) BatchUpdateBlance(c context.Context, tx *sql.Tx, mid int64, blances []*model.CouponBalanceInfo) (a int64, err error) {
	var (
		res xsql.Result
		buf bytes.Buffer
		ids []int64
	)
	buf.WriteString(fmt.Sprintf(_batchUpdateBalance, hitUser(mid)))
	for _, v := range blances {
		buf.WriteString(" WHEN ")
		buf.WriteString(strconv.FormatInt(v.ID, 10))
		buf.WriteString(" THEN ")
		buf.WriteString(strconv.FormatInt(v.Balance, 10))
		ids = append(ids, v.ID)
	}
	buf.WriteString(" END  WHERE `id` in (")
	buf.WriteString(xstr.JoinInts(ids))
	buf.WriteString(") AND `ver` = CASE id ")
	for _, v := range blances {
		buf.WriteString(" WHEN ")
		buf.WriteString(strconv.FormatInt(v.ID, 10))
		buf.WriteString(" THEN ")
		buf.WriteString(strconv.FormatInt(v.Ver, 10))
	}
	buf.WriteString(" END;")
	if res, err = tx.Exec(buf.String()); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// BatchInsertBlanceLog Batch Insert Balance log
func (d *Dao) BatchInsertBlanceLog(c context.Context, tx *sql.Tx, mid int64, ls []*model.CouponBalanceChangeLog) (a int64, err error) {
	var (
		buf bytes.Buffer
		res xsql.Result
		sql string
	)
	buf.WriteString(fmt.Sprintf(_addBalanceLogSQL, hitUserLog(mid)))
	for _, v := range ls {
		buf.WriteString("('")
		buf.WriteString(v.OrderNo)
		buf.WriteString("',")
		buf.WriteString(strconv.FormatInt(v.Mid, 10))
		buf.WriteString(",'")
		buf.WriteString(v.BatchToken)
		buf.WriteString("',")
		buf.WriteString(strconv.FormatInt(v.Balance, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(v.ChangeBalance, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(int(v.ChangeType)))
		buf.WriteString(",'")
		buf.WriteString(fmt.Sprintf("%v", v.Ctime.Time().Format("2006-01-02 15:04:05")))
		buf.WriteString("'),")
	}
	sql = buf.String()
	if res, err = tx.Exec(sql[0 : len(sql)-1]); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// AddOrder add order.
func (d *Dao) AddOrder(c context.Context, tx *sql.Tx, o *model.CouponOrder) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_addOrderSQL, o.OrderNo, o.Mid, o.Count, o.State, o.CouponType, o.ThirdTradeNo, o.Remark, o.Tips, o.UseVer, o.Ver, o.Ctime); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// AddOrderLog add order log.
func (d *Dao) AddOrderLog(c context.Context, tx *sql.Tx, o *model.CouponOrderLog) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_addOrderLogSQL, o.OrderNo, o.Mid, o.State, o.Ctime); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// CouponCarToonCount coupon cartoon page.
func (d *Dao) CouponCarToonCount(c context.Context, mid int64, t int64, ct int8, state int8, stime time.Time) (count int64, err error) {
	var row *sql.Row
	switch state {
	case model.NotUsed:
		row = d.db.QueryRow(c, fmt.Sprintf(_countBalanceNotUsed, hitUser(mid)), mid, t, t, ct, stime)
	case model.Used:
		row = d.db.QueryRow(c, _countUseListSQL, mid, state, ct, stime)
	case model.Expire:
		row = d.db.QueryRow(c, fmt.Sprintf(_countBalanceExpire, hitUser(mid)), mid, t, ct, stime)
	default:
		return
	}
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

// CouponNotUsedPage query coupon page .
func (d *Dao) CouponNotUsedPage(c context.Context, mid int64, ct int8, t int64, stime time.Time, pn int, ps int) (res []*model.CouponBalanceInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_balanceNotUsedPageSQL, hitUser(mid)), mid, t, t, ct, stime, (pn-1)*ps, ps); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.CouponBalanceInfo{}
		if err = rows.Scan(&r.ID, &r.BatchToken, &r.Mid, &r.Balance, &r.StartTime, &r.ExpireTime, &r.Origin, &r.CouponType, &r.Ver, &r.CTime, &r.MTime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// CouponExpirePage query coupon page .
func (d *Dao) CouponExpirePage(c context.Context, mid int64, ct int8, t int64, stime time.Time, pn int, ps int) (res []*model.CouponBalanceInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_balanceExpirePageSQL, hitUser(mid)), mid, t, ct, stime, (pn-1)*ps, ps); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.CouponBalanceInfo{}
		if err = rows.Scan(&r.ID, &r.BatchToken, &r.Mid, &r.Balance, &r.StartTime, &r.ExpireTime, &r.Origin, &r.CouponType, &r.Ver, &r.CTime, &r.MTime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// OrderUsedPage order used page.
func (d *Dao) OrderUsedPage(c context.Context, mid int64, state int8, ct int8, stime time.Time, pn int, ps int) (res []*model.CouponOrder, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _useOrderPageSQL, mid, state, ct, stime, (pn-1)*ps, ps); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.CouponOrder{}
		if err = rows.Scan(&r.ID, &r.OrderNo, &r.Mid, &r.Count, &r.State, &r.CouponType, &r.ThirdTradeNo, &r.Remark, &r.Tips, &r.UseVer,
			&r.Ver, &r.Ctime, &r.Mtime); err != nil {
			if err == sql.ErrNoRows {
				err = nil
				res = nil
				return
			}
			err = errors.WithStack(err)
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// AddBalanceCoupon add balance coupon.
func (d *Dao) AddBalanceCoupon(c context.Context, tx *sql.Tx, b *model.CouponBalanceInfo) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_addBalanceCouponSQL, hitUser(b.Mid)), b.BatchToken, b.Mid, b.Balance, b.StartTime, b.ExpireTime, b.Origin, b.CouponType, b.Ver, b.CTime,
		b.Balance); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// ByMidAndBatchToken query coupon by batch token and mid.
func (d *Dao) ByMidAndBatchToken(c context.Context, mid int64, batchToken string) (r *model.CouponBalanceInfo, err error) {
	var row *sql.Row
	r = &model.CouponBalanceInfo{}
	row = d.db.QueryRow(c, fmt.Sprintf(_byMidAndBatchTokenSQL, hitUser(mid)), mid, batchToken)
	if err = row.Scan(&r.ID, &r.BatchToken, &r.Mid, &r.Balance, &r.StartTime, &r.ExpireTime, &r.Origin, &r.CouponType, &r.Ver, &r.CTime, &r.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// AddBalanceChangeLog add coupon balance change log.
func (d *Dao) AddBalanceChangeLog(c context.Context, tx *sql.Tx, bl *model.CouponBalanceChangeLog) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_addBlanceChangeLog, hitUserLog(bl.Mid)), bl.OrderNo, bl.Mid, bl.BatchToken, bl.Balance, bl.ChangeBalance, bl.ChangeType, bl.Ctime); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// BatchInfo batch info.
func (d *Dao) BatchInfo(c context.Context, token string) (r *model.CouponBatchInfo, err error) {
	var row *sql.Row
	r = new(model.CouponBatchInfo)
	row = d.db.QueryRow(c, _batchInfoSQL, token)
	if err = row.Scan(&r.ID, &r.AppID, &r.Name, &r.BatchToken, &r.MaxCount, &r.CurrentCount, &r.StartTime, &r.ExpireTime, &r.ExpireDay, &r.Ver, &r.Ctime, &r.LimitCount,
		&r.FullAmount, &r.Amount, &r.State, &r.CouponType); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// UpdateBatchInfo update batch info.
func (d *Dao) UpdateBatchInfo(c context.Context, tx *sql.Tx, token string, count int) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_updateBatchSQL, count, token); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UpdateBatchLimitInfo update batch limit info.
func (d *Dao) UpdateBatchLimitInfo(c context.Context, tx *sql.Tx, token string, count int) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_updateBatchLimitSQL, count, token, count); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// GrantCouponLog grant coupon log.
func (d *Dao) GrantCouponLog(c context.Context, mid int64, token string, ct int8) (rs []*model.CouponBalanceChangeLog, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_grantCouponLogSQL, hitUserLog(mid)), mid, token, ct); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.CouponBalanceChangeLog{}
		if err = rows.Scan(&r.ID, &r.OrderNo, &r.Mid, &r.BatchToken, &r.Balance, &r.ChangeBalance, &r.ChangeType, &r.Ctime, &r.Mtime); err != nil {
			err = errors.WithStack(err)
			rs = nil
			return
		}
		rs = append(rs, r)
	}
	err = rows.Err()
	return
}

// AllBranchInfo query all branch info .
func (d *Dao) AllBranchInfo(c context.Context) (res []*model.CouponBatchInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _allBatchInfoSQL); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.CouponBatchInfo{}
		if err = rows.Scan(&r.ID, &r.AppID, &r.Name, &r.BatchToken, &r.MaxCount, &r.CurrentCount, &r.StartTime, &r.ExpireTime, &r.ExpireDay, &r.Ver, &r.Ctime,
			&r.LimitCount, &r.FullAmount, &r.Amount, &r.State, &r.CouponType, &r.PlatformLimit, &r.ProdLimMonth, &r.ProdLimRenewal); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// CountByBranchToken get user count by bratch token.
func (d *Dao) CountByBranchToken(c context.Context, mid int64, token string) (count int64, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_batchCountByMid, hitInfo(mid)), mid, token)
	if err = row.Scan(&count); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//ReceiveLog get receive log.
func (d *Dao) ReceiveLog(c context.Context, appkey, orderNo string, ct int8) (r *model.CouponReceiveLog, err error) {
	row := d.db.QueryRow(c, _couponReceiveSQL, orderNo, appkey, ct)
	r = new(model.CouponReceiveLog)
	if err = row.Scan(&r.ID, &r.Appkey, &r.OrderNo, &r.Mid, &r.CouponToken, &r.CouponType); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
			return
		}
		err = errors.WithStack(err)
	}
	return
}

//TxAddReceiveLog add receive log.
func (d *Dao) TxAddReceiveLog(tx *sql.Tx, rlog *model.CouponReceiveLog) (err error) {
	if _, err = tx.Exec(_addReceiveSQL, rlog.Appkey, rlog.OrderNo, rlog.Mid, rlog.CouponToken, rlog.CouponType); err != nil {
		err = errors.WithStack(err)
	}
	return
}
