package dao

import (
	"bytes"
	"context"
	xsql "database/sql"
	"fmt"
	"go-common/app/service/main/coupon/model"
	"go-common/library/database/sql"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	_couponAllowanceNoStartCheckSQL = "SELECT id,coupon_token,mid,state,start_time,expire_time,origin,ver,batch_token,order_no,amount,full_amount,ctime,mtime FROM coupon_allowance_info_%02d WHERE mid = ?  AND expire_time > ? AND state = ?;"
	_couponByOrderNOSQL             = "SELECT id,coupon_token,mid,state,start_time,expire_time,origin,ver,batch_token,order_no,amount,full_amount,ctime,mtime,remark FROM coupon_allowance_info_%02d WHERE order_no = ?;"
	_couponUsableAllowanceSQL       = "SELECT id,coupon_token,mid,state,start_time,expire_time,origin,ver,batch_token,order_no,amount,full_amount,ctime,mtime FROM coupon_allowance_info_%02d WHERE mid = ?  AND expire_time > ? AND start_time < ? AND state = ?;"
	_couponAllowanceByTokenSQL      = "SELECT id,coupon_token,mid,state,start_time,expire_time,origin,ver,batch_token,order_no,amount,full_amount,ctime,mtime FROM coupon_allowance_info_%02d WHERE coupon_token = ?;"
	_updateCouponAllowanceInUseSQL  = "UPDATE coupon_allowance_info_%02d SET state =?, order_no = ?, remark = ?, ver =ver+1  WHERE coupon_token = ? AND ver = ?;"
	_updateCouponAllowanceToUseSQL  = "UPDATE coupon_allowance_info_%02d SET state =?, order_no = ?, ver =ver+1  WHERE coupon_token = ? AND ver = ? AND state = ?;"
	_getCouponByOrderNoSQL          = "SELECT mid,coupon_token,order_no,amount,full_amount,state,ver FROM coupon_allowance_info_%02d WHERE order_no = ?"
	_addCouponAllowanceChangeLogSQL = "INSERT INTO coupon_allowance_change_log_%02d (coupon_token,order_no,mid,state,ctime, change_type) VALUES(?,?,?,?,?,?);"
	_batchAllowanceCountByMid       = "SELECT  COUNT(1) FROM  coupon_allowance_info_%02d WHERE mid = ? AND batch_token = ?;"
	_batchAddAllowanceCouponSQL     = "INSERT INTO coupon_allowance_info_%02d(coupon_token,mid,state,start_time,expire_time,origin,batch_token,amount,full_amount,ctime,app_id) VALUES "
	_addAllowanceCouponSQL          = "INSERT INTO coupon_allowance_info_%02d(coupon_token,mid,state,start_time,expire_time,origin,batch_token,amount,full_amount,app_id) VALUES (?,?,?,?,?,?,?,?,?,?)"
	_couponAllowancePageNotUsedSQL  = "SELECT id,coupon_token,mid,state,start_time,expire_time,origin,ver,batch_token,order_no,amount,full_amount,ctime,mtime,remark  FROM coupon_allowance_info_%02d WHERE mid = ? AND (state = 0 OR state = 1)  AND expire_time > ? AND start_time < ? AND ctime > ? ORDER BY id DESC"
	_couponAllowancePageUsedSQL     = "SELECT id,coupon_token,mid,state,start_time,expire_time,origin,ver,batch_token,order_no,amount,full_amount,ctime,mtime,remark  FROM coupon_allowance_info_%02d WHERE mid = ? AND state = 2 AND ctime > ?  ORDER BY id DESC "
	_couponAllowancePageExpireSQL   = "SELECT id,coupon_token,mid,state,start_time,expire_time,origin,ver,batch_token,order_no,amount,full_amount,ctime,mtime,remark  FROM coupon_allowance_info_%02d WHERE mid = ? AND state <> 2 AND expire_time < ? AND ctime > ? ORDER BY id DESC "
)

func hitAllowanceInfo(mid int64) int64 {
	return mid % 10
}

func hitAllowanceChangeLog(mid int64) int64 {
	return mid % 10
}

// ByStateAndExpireAllowances query by coupon state and expire .
func (d *Dao) ByStateAndExpireAllowances(c context.Context, mid int64, state int8, t int64) (res []*model.CouponAllowanceInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_couponAllowanceNoStartCheckSQL, hitAllowanceInfo(mid)), mid, t, state); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.CouponAllowanceInfo{}
		if err = rows.Scan(&r.ID, &r.CouponToken, &r.Mid, &r.State, &r.StartTime, &r.ExpireTime, &r.Origin, &r.Ver, &r.BatchToken,
			&r.OrderNO, &r.Amount, &r.FullAmount, &r.CTime, &r.MTime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// AllowanceByOrderNO query coupon by orderno.
func (d *Dao) AllowanceByOrderNO(c context.Context, mid int64, orderNO string) (r *model.CouponAllowanceInfo, err error) {
	var row *sql.Row
	r = &model.CouponAllowanceInfo{}
	row = d.db.QueryRow(c, fmt.Sprintf(_couponByOrderNOSQL, hitAllowanceInfo(mid)), orderNO)
	if err = row.Scan(&r.ID, &r.CouponToken, &r.Mid, &r.State, &r.StartTime, &r.ExpireTime, &r.Origin, &r.Ver, &r.BatchToken,
		&r.OrderNO, &r.Amount, &r.FullAmount, &r.CTime, &r.MTime, &r.Remark); err != nil {
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

// UsableAllowances usable allowance .
func (d *Dao) UsableAllowances(c context.Context, mid int64, state int8, t int64) (res []*model.CouponAllowanceInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_couponUsableAllowanceSQL, hitAllowanceInfo(mid)), mid, t, t, state); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.CouponAllowanceInfo{}
		if err = rows.Scan(&r.ID, &r.CouponToken, &r.Mid, &r.State, &r.StartTime, &r.ExpireTime, &r.Origin, &r.Ver, &r.BatchToken,
			&r.OrderNO, &r.Amount, &r.FullAmount, &r.CTime, &r.MTime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// AllowanceByToken query coupon by token.
func (d *Dao) AllowanceByToken(c context.Context, mid int64, token string) (r *model.CouponAllowanceInfo, err error) {
	var row *sql.Row
	r = &model.CouponAllowanceInfo{}
	row = d.db.QueryRow(c, fmt.Sprintf(_couponAllowanceByTokenSQL, hitAllowanceInfo(mid)), token)
	if err = row.Scan(&r.ID, &r.CouponToken, &r.Mid, &r.State, &r.StartTime, &r.ExpireTime, &r.Origin, &r.Ver, &r.BatchToken,
		&r.OrderNO, &r.Amount, &r.FullAmount, &r.CTime, &r.MTime); err != nil {
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

// UpdateAllowanceCouponInUse update coupon in use.
func (d *Dao) UpdateAllowanceCouponInUse(c context.Context, tx *sql.Tx, cp *model.CouponAllowanceInfo) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_updateCouponAllowanceInUseSQL, hitAllowanceInfo(cp.Mid)), cp.State, cp.OrderNO, cp.Remark, cp.CouponToken, cp.Ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UpdateAllowanceCouponToUse update coupon in use.
func (d *Dao) UpdateAllowanceCouponToUse(c context.Context, tx *sql.Tx, cp *model.CouponAllowanceInfo) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_updateCouponAllowanceToUseSQL, hitAllowanceInfo(cp.Mid)), cp.State, cp.OrderNO, cp.CouponToken, cp.Ver, model.InUse); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UpdateAllowanceCouponToUsed update coupon in used.
func (d *Dao) UpdateAllowanceCouponToUsed(c context.Context, tx *sql.Tx, cp *model.CouponAllowanceInfo) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_updateCouponAllowanceToUseSQL, hitAllowanceInfo(cp.Mid)), cp.State, cp.OrderNO, cp.CouponToken, cp.Ver, model.NotUsed); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//InsertCouponAllowanceHistory insert coupon history .
func (d *Dao) InsertCouponAllowanceHistory(c context.Context, tx *sql.Tx, l *model.CouponAllowanceChangeLog) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_addCouponAllowanceChangeLogSQL, hitAllowanceChangeLog(l.Mid)), l.CouponToken, l.OrderNO, l.Mid, l.State, l.Ctime, l.ChangeType); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// CountByAllowanceBranchToken get user count by bratch token.
func (d *Dao) CountByAllowanceBranchToken(c context.Context, mid int64, token string) (count int64, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_batchAllowanceCountByMid, hitAllowanceInfo(mid)), mid, token)
	if err = row.Scan(&count); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// GetCouponByOrderNo .
func (d *Dao) GetCouponByOrderNo(c context.Context, mid int64, orderNo string) (res *model.CouponAllowanceInfo, err error) {
	res = &model.CouponAllowanceInfo{}
	row := d.db.QueryRow(c, fmt.Sprintf(_getCouponByOrderNoSQL, hitAllowanceInfo(mid)), orderNo)
	if err = row.Scan(&res.Mid, &res.CouponToken, &res.OrderNO, &res.Amount, &res.FullAmount, &res.State, &res.Ver); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//TxAddAllowanceCoupon tx add lowance coupon
func (d *Dao) TxAddAllowanceCoupon(tx *sql.Tx, cp *model.CouponAllowanceInfo) (err error) {
	if _, err = tx.Exec(fmt.Sprintf(_addAllowanceCouponSQL, hitAllowanceInfo(cp.Mid)), cp.CouponToken, cp.Mid, cp.State, cp.StartTime, cp.ExpireTime, cp.Origin, cp.BatchToken, cp.Amount, cp.FullAmount, cp.AppID); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// BatchAddAllowanceCoupon batch add allowance coupon.
func (d *Dao) BatchAddAllowanceCoupon(c context.Context, tx *sql.Tx, mid int64, cps []*model.CouponAllowanceInfo) (a int64, err error) {
	var (
		buf bytes.Buffer
		res xsql.Result
		sql string
	)
	buf.WriteString(fmt.Sprintf(_batchAddAllowanceCouponSQL, hitAllowanceInfo(mid)))
	for _, v := range cps {
		buf.WriteString("('")
		buf.WriteString(v.CouponToken)
		buf.WriteString("',")
		buf.WriteString(strconv.FormatInt(v.Mid, 10))
		buf.WriteString(",")
		buf.WriteString(fmt.Sprintf("%d", v.State))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(v.StartTime, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(v.ExpireTime, 10))
		buf.WriteString(",")
		buf.WriteString(fmt.Sprintf("%d", v.Origin))
		buf.WriteString(",'")
		buf.WriteString(v.BatchToken)
		buf.WriteString("',")
		buf.WriteString(fmt.Sprintf("%f", v.Amount))
		buf.WriteString(",")
		buf.WriteString(fmt.Sprintf("%f", v.FullAmount))
		buf.WriteString(",'")
		buf.WriteString(v.CTime.Time().Format("2006-01-02 15:04:05"))
		buf.WriteString("',")
		buf.WriteString(strconv.FormatInt(v.AppID, 10))
		buf.WriteString("),")
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

// AllowanceList list.
func (d *Dao) AllowanceList(c context.Context, mid int64, state int8, t int64, stime time.Time) (res []*model.CouponAllowanceInfo, err error) {
	var rows *sql.Rows
	switch state {
	case model.NotUsed:
		rows, err = d.db.Query(c, fmt.Sprintf(_couponAllowancePageNotUsedSQL, hitAllowanceInfo(mid)), mid, t, t, stime)
	case model.Used:
		rows, err = d.db.Query(c, fmt.Sprintf(_couponAllowancePageUsedSQL, hitAllowanceInfo(mid)), mid, stime)
	case model.Expire:
		rows, err = d.db.Query(c, fmt.Sprintf(_couponAllowancePageExpireSQL, hitAllowanceInfo(mid)), mid, t, stime)
	default:
		return
	}
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.CouponAllowanceInfo)
		if err = rows.Scan(&r.ID, &r.CouponToken, &r.Mid, &r.State, &r.StartTime, &r.ExpireTime, &r.Origin, &r.Ver, &r.BatchToken, &r.OrderNO, &r.Amount, &r.FullAmount,
			&r.CTime, &r.MTime, &r.Remark); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}
