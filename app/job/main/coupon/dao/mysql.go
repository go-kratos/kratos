package dao

import (
	"bytes"
	"context"
	xsql "database/sql"
	"fmt"
	"strconv"
	"time"

	"go-common/app/job/main/coupon/model"
	"go-common/library/database/sql"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_updateStateSQL        = "UPDATE `coupon_info_%02d` SET `state` = ?,`use_ver` = ?,`ver` = ? WHERE `coupon_token` = ? AND `ver` = ?;"
	_couponByTokenSQL      = "SELECT `id`,`coupon_token`,`mid`,`state`,`start_time`,`expire_time`,`origin`,`coupon_type`,`order_no`,`oid`,`remark`,`use_ver`,`ver`,`ctime`,`mtime` FROM `coupon_info_%02d` WHERE `coupon_token` = ?;"
	_couponsSQL            = "SELECT `id`,`coupon_token`,`mid`,`state`,`start_time`,`expire_time`,`origin`,`coupon_type`,`order_no`,`oid`,`remark`,`use_ver`,`ver`,`ctime`,`mtime` FROM `coupon_info_%02d` WHERE  `state` = ? AND `mtime` > ?;"
	_addCouponChangeLogSQL = "INSERT INTO `coupon_change_log_%02d` (`coupon_token`,`mid`,`state`,`ctime`) VALUES(?,?,?,?);"

	// balance coupon
	_byOrderNoSQL          = "SELECT `id`,`order_no`,`mid`,`count`,`state`,`coupon_type`,`third_trade_no`,`remark`,`tips`,`use_ver`,`ver`,`ctime`,`mtime` FROM `coupon_order` WHERE  `order_no` = ?;"
	_updateOrderSQL        = "UPDATE `coupon_order` SET `state` = ?,`use_ver` =?,`ver` = ? WHERE `order_no` = ? AND `ver`  = ?;"
	_addOrderLogSQL        = "INSERT INTO `coupon_order_log`(`order_no`,`mid`,`state`,`ctime`)VALUES(?,?,?,?);"
	_consumeCouponLogSQL   = "SELECT `id`,`order_no`,`mid`,`batch_token`,`balance`,`change_balance`,`change_type`,`ctime`,`mtime` FROM `coupon_balance_change_log_%02d` WHERE  `order_no` = ? AND `change_type` = ?;"
	_byMidAndBatchTokenSQL = "SELECT `id`,`batch_token`,`mid`,`balance`,`start_time`,`expire_time`,`origin`,`coupon_type`,`ver`,`ctime`,`mtime` FROM `coupon_balance_info_%02d` WHERE `mid` = ?  AND `batch_token` = ? ;"
	_inPayOrderSQL         = "SELECT `id`,`order_no`,`mid`,`count`,`state`,`coupon_type`,`third_trade_no`,`remark`,`tips`,`use_ver`,`ver`,`ctime`,`mtime` FROM `coupon_order` WHERE `state` = ? AND `mtime` > ?;"
	_batchUpdateBalance    = "UPDATE `coupon_balance_info_%02d` SET `ver` =`ver` + 1, `balance` = CASE id"
	_addBalanceLogSQL      = "INSERT INTO `coupon_balance_change_log_%02d`(`order_no`,`mid`,`batch_token`,`balance`,`change_balance`,`change_type`,`ctime`)VALUES "
	_couponBlancesSQL      = "SELECT `id`,`batch_token`,`mid`,`balance`,`start_time`,`expire_time`,`origin`,`coupon_type`,`ver`,`ctime`,`mtime` FROM `coupon_balance_info_%02d` WHERE `mid` = ? AND `coupon_type` = ?;"
	_updateBlanceSQL       = "UPDATE `coupon_balance_info_%02d` SET `balance` = ?,`ver` = `ver` + 1 WHERE `id` = ? AND `ver` = ?;"

	_updateUserCardSQL = "UPDATE coupon_user_card SET state=? WHERE mid=? AND coupon_token=? AND batch_token=?"
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

// UpdateCoupon update coupon in use.
func (d *Dao) UpdateCoupon(c context.Context, tx *sql.Tx, mid int64, state int8, useVer int64, ver int64, couponToken string) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_updateStateSQL, hitInfo(mid)), state, useVer, ver+1, couponToken, ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
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
		if err == xsql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// CouponList query .
func (d *Dao) CouponList(c context.Context, index int64, state int8, t time.Time) (res []*model.CouponInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_couponsSQL, hitInfo(index)), state, t); err != nil {
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

// BeginTran begin transaction.
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.db.Begin(c)
}

// ByOrderNo query order by order no.
func (d *Dao) ByOrderNo(c context.Context, orderNo string) (r *model.CouponOrder, err error) {
	var row *sql.Row
	r = &model.CouponOrder{}
	row = d.db.QueryRow(c, _byOrderNoSQL, orderNo)
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

// UpdateOrderState update order state.
func (d *Dao) UpdateOrderState(c context.Context, tx *sql.Tx, mid int64, state int8, useVer int64, ver int64, orderNo string) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_updateOrderSQL, state, useVer, ver+1, orderNo, ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
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

// ConsumeCouponLog consume coupon log.
func (d *Dao) ConsumeCouponLog(c context.Context, mid int64, orderNo string, ct int8) (rs []*model.CouponBalanceChangeLog, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_consumeCouponLogSQL, hitUserLog(mid)), orderNo, ct); err != nil {
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

// OrderInPay order in pay.
func (d *Dao) OrderInPay(c context.Context, state int8, t time.Time) (res []*model.CouponOrder, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _inPayOrderSQL, state, t); err != nil {
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

// BlanceList user balance by mid.
func (d *Dao) BlanceList(c context.Context, mid int64, ct int8) (res []*model.CouponBalanceInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_couponBlancesSQL, hitUser(mid)), mid, ct); err != nil {
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

// UpdateUserCard .
func (d *Dao) UpdateUserCard(c context.Context, mid int64, state int8, couponToken, batchToken string) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _updateUserCardSQL, state, mid, couponToken, batchToken); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
