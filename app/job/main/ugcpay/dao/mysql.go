package dao

import (
	"context"
	"math"
	"time"

	"go-common/app/job/main/ugcpay/model"
	xsql "go-common/library/database/sql"
)

const (
	_asset = `SELECT id,mid,oid,otype,currency,price,state,ctime,mtime FROM asset WHERE oid=? AND otype=? AND currency=? LIMIT 1`

	_countPaidOrderUser          = "SELECT count(1) FROM order_user WHERE pay_time BETWEEN ? AND ? AND state='paid' LIMIT 1"
	_countRefundedOrderUser      = "SELECT count(1) FROM order_user WHERE refund_time BETWEEN ? AND ? AND state='st_refunded' LIMIT 1"
	_sumPaidOrderUserRealFee     = "SELECT IFNULL(sum(real_fee),0) FROM order_user WHERE pay_time BETWEEN ? AND ? AND state='settled' LIMIT 1"
	_sumRefundedOrderUserRealFee = "SELECT IFNULL(sum(real_fee),0) FROM order_user WHERE refund_time BETWEEN ? AND ? AND state='ref_finished' LIMIT 1"
	_minIDOrderPaid              = `SELECT id FROM order_user WHERE pay_time>=? AND state='paid' ORDER BY id ASC LIMIT 1`
	_minIDOrderRefunded          = `SELECT id FROM order_user WHERE refund_time>=? AND state='st_refunded' ORDER BY id ASC LIMIT 1`
	_orderPaidList               = `SELECT id,mid,order_id,biz,platform,oid,otype,fee,real_fee,currency,pay_id,pay_reason,pay_time,state,ctime,mtime,refund_time,version FROM order_user WHERE pay_time BETWEEN ? AND ? AND state='paid' AND id>? ORDER BY ID ASC LIMIT ?`
	_orderRefundedList           = `SELECT id,mid,order_id,biz,platform,oid,otype,fee,real_fee,currency,pay_id,pay_reason,pay_time,state,ctime,mtime,refund_time,version FROM order_user WHERE refund_time BETWEEN ? AND ? AND state='st_refunded' AND id>? ORDER BY ID ASC LIMIT ?`
	_updateOrder                 = `UPDATE order_user SET mid=?,order_id=?,biz=?,platform=?,oid=?,otype=?,fee=?,real_fee=?,currency=?,pay_id=?,pay_reason=?,pay_time=?,state=?,refund_time=?,version=version+1 WHERE id=? AND version=?`
	_insertLogOrderUser          = "INSERT INTO log_order_user (order_id,from_state,to_state,`desc`) VALUES (?,?,?,?)"

	_orderBadDebt       = `SELECT id,order_id,type,state,ctime,mtime FROM order_bad_debt WHERE order_id=? ORDER BY id ASC LIMIT 1`
	_insertOrderBadDebt = `INSERT INTO order_bad_debt (order_id,type,state) VALUES (?,?,?)`
	_updateOrderBadDebt = `UPDATE order_bad_debt SET order_id=?,type=?,state=? WHERE order_id=?`

	_countDailyBillByVer      = "SELECT count(1) FROM bill_user_daily WHERE ver=? LIMIT 1"
	_countDailyBillByMonthVer = "SELECT count(1) FROM bill_user_daily WHERE month_ver=? LIMIT 1"
	_sumDailyBill             = "SELECT IFNULL(sum(`in`),0),IFNULL(sum(`out`),0) FROM bill_user_daily WHERE ver=? LIMIT 1"
	_minIDDailyBillByMonthVer = `SELECT id FROM bill_user_daily WHERE month_ver=? ORDER BY id ASC LIMIT 1`
	_dailyBillListByMonthVer  = "SELECT id,bill_id,mid,biz,currency,`in`,`out`,ver,month_ver,ctime,mtime,version FROM bill_user_daily WHERE month_ver=? AND id>? ORDER BY ID ASC LIMIT ?"
	_minIDDailyBillByVer      = `SELECT id FROM bill_user_daily WHERE ver=? ORDER BY id ASC LIMIT 1`
	_dailyBillListByVer       = "SELECT id,bill_id,mid,biz,currency,`in`,`out`,ver,month_ver,ctime,mtime,version FROM bill_user_daily WHERE ver=? AND id>? ORDER BY ID ASC LIMIT ?"
	_dailyBill                = "SELECT id,bill_id,mid,biz,currency,`in`,`out`,ver,month_ver,ctime,mtime,version FROM bill_user_daily WHERE mid=? AND biz=? AND currency=? AND ver=? LIMIT 1"
	_insertDailyBill          = "INSERT INTO bill_user_daily (bill_id,mid,biz,currency,`in`,`out`,ver,month_ver,version) VALUES (?,?,?,?,?,?,?,?,?)"
	_updateDailyBill          = "UPDATE bill_user_daily SET bill_id=?,mid=?,biz=?,currency=?,`in`=?,`out`=?,ver=?,month_ver=?,version=version+1 WHERE mid=? AND biz=? AND currency=? AND ver=? AND version=?"
	_insertDailyBillLog       = "INSERT INTO log_bill_user_daily (bill_id,from_in,to_in,from_out,to_out,order_id) VALUES (?,?,?,?,?,?)"

	_countMonthlyBillByVer = "SELECT count(1) FROM bill_user_monthly WHERE ver=? LIMIT 1"
	_minIDMonthlyBill      = `SELECT id FROM bill_user_monthly WHERE ver=? ORDER BY id ASC LIMIT 1`
	_monthlyBillList       = "SELECT id,bill_id,mid,biz,currency,`in`,`out`,ver,ctime,mtime,version FROM bill_user_monthly WHERE ver=? AND id>? ORDER BY ID ASC LIMIT ?"
	_monthlyBill           = "SELECT id,bill_id,mid,biz,currency,`in`,`out`,ver,ctime,mtime,version FROM bill_user_monthly WHERE mid=? AND biz=? AND currency=? AND ver=? LIMIT 1"
	_insertMonthlyBill     = "INSERT INTO bill_user_monthly (bill_id,mid,biz,currency,`in`,`out`,ver,version) VALUES (?,?,?,?,?,?,?,?)"
	_updateMonthlyBill     = "UPDATE bill_user_monthly SET bill_id=?,mid=?,biz=?,currency=?,`in`=?,`out`=?,ver=?,version=version+1 WHERE mid=? AND biz=? AND currency=? AND ver=? AND version=?"
	_insertMonthlyBillLog  = "INSERT INTO log_bill_user_monthly (bill_id,from_in,to_in,from_out,to_out,bill_user_daily_id) VALUES (?,?,?,?,?,?)"

	_minIDUserAccount     = `SELECT id FROM account_user WHERE mtime>=? ORDER BY id ASC LIMIT 1`
	_userAccountList      = "SELECT id,mid,biz,currency,balance,ver,state,ctime,mtime FROM account_user WHERE mtime BETWEEN ? AND ? AND id>? ORDER BY ID ASC LIMIT ?"
	_userAccount          = `SELECT id,mid,biz,currency,balance,ver,state,ctime,mtime FROM account_user WHERE mid=? AND biz=? AND currency=? LIMIT 1`
	_insertUserAccount    = `INSERT INTO account_user (mid,biz,currency,balance,ver,state) VALUES (?,?,?,?,?,?)`
	_updateUserAccount    = `UPDATE account_user SET mid=?,biz=?,currency=?,balance=?,ver=ver+1,state=? WHERE mid=? AND biz=? AND currency=? AND ver=?`
	_insertUserAccountLog = "INSERT INTO log_account_user (account_id,`from`,`to`,ver,state,name) VALUES (?,?,?,?,?,?)"

	_bizAccount          = `SELECT id,biz,currency,balance,ver,state,ctime,mtime FROM account_biz WHERE biz=? AND currency=? LIMIT 1`
	_insertBizAccount    = `INSERT INTO account_biz (biz,currency,balance,ver,state) VALUES (?,?,?,?,?)`
	_updateBizAccount    = `UPDATE account_biz SET biz=?,currency=?,balance=?,ver=ver+1,state=? WHERE biz=? AND currency=? AND ver=?`
	_insertBizAccountLog = "INSERT INTO log_account_biz (account_id,`from`,`to`,ver,state,name) VALUES (?,?,?,?,?,?)"

	_aggrIncomeUser       = "SELECT id,mid,currency,pay_success,pay_error,total_in,total_out,ctime,mtime FROM aggr_income_user WHERE mid=? AND currency=? LIMIT 1"
	_insertAggrIncomeUser = "INSERT INTO aggr_income_user (mid,currency,pay_success,pay_error,total_in,total_out) VALUES (?,?,?,?,?,?)"
	_updateAggrIncomeUser = "UPDATE aggr_income_user SET mid=?,currency=?,pay_success=?,pay_error=?,total_in=?,total_out=? WHERE mid=? AND currency=?"

	_aggrIncomeUserAsset       = "SELECT id,mid,currency,ver,oid,otype,pay_success,pay_error,total_in,total_out,ctime,mtime FROM aggr_income_user_asset WHERE mid=? AND currency=? AND ver=? AND oid=? AND otype=? LIMIT 1"
	_insertAggrIncomeUserAsset = "INSERT INTO aggr_income_user_asset (mid,currency,`ver`,oid,otype,pay_success,pay_error,total_in,total_out) VALUES (?,?,?,?,?,?,?,?,?)"
	_updateAggrIncomeUserAsset = "UPDATE aggr_income_user_asset SET mid=?,currency=?,ver=?,oid=?,otype=?,pay_success=?,pay_error=?,total_in=?,total_out=? WHERE mid=? AND currency=? AND ver=? AND oid=? AND otype=?"

	_insertOrderRechargeShell    = "INSERT INTO order_recharge_shell (mid,order_id,biz,amount,pay_msg,state,`ver`) VALUES (?,?,?,?,?,?,?)"
	_insertOrderRechargeShellLog = "INSERT INTO log_order_recharge_shell (order_id,from_state,to_state,`desc`,bill_user_monthly_id) VALUES (?,?,?,?,?)"

	_logTask            = "SELECT id,name,expect,success,failure,state,ctime,mtime FROM log_task WHERE name=? LIMIT 1"
	_insertLogTask      = "INSERT INTO log_task (name,expect,success,failure,state) VALUES (?,?,?,?,?)"
	_logTaskSuccessIncr = "UPDATE log_task SET success=success+1 WHERE name=?"
	_logTaskFailureIncr = "UPDATE log_task SET failure=failure+1 WHERE name=?"
)

// CountPaidOrderUser .
func (d *Dao) CountPaidOrderUser(ctx context.Context, beginTime, endTime time.Time) (count int64, err error) {
	row := d.db.QueryRow(ctx, _countPaidOrderUser, beginTime, endTime)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			count = 0
		}
		return
	}
	return
}

// CountRefundedOrderUser .
func (d *Dao) CountRefundedOrderUser(ctx context.Context, beginTime, endTime time.Time) (count int64, err error) {
	row := d.db.QueryRow(ctx, _countRefundedOrderUser, beginTime, endTime)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			count = 0
		}
		return
	}
	return
}

// CountDailyBillByVer .
func (d *Dao) CountDailyBillByVer(ctx context.Context, ver int64) (count int64, err error) {
	row := d.db.QueryRow(ctx, _countDailyBillByVer, ver)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			count = 0
		}
		return
	}
	return
}

// CountDailyBillByMonthVer .
func (d *Dao) CountDailyBillByMonthVer(ctx context.Context, monthVer int64) (count int64, err error) {
	row := d.db.QueryRow(ctx, _countDailyBillByMonthVer, monthVer)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			count = 0
		}
		return
	}
	return
}

// CountMonthlyBillByVer .
func (d *Dao) CountMonthlyBillByVer(ctx context.Context, ver int64) (count int64, err error) {
	row := d.db.QueryRow(ctx, _countMonthlyBillByVer, ver)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			count = 0
		}
		return
	}
	return
}

// LogTask .
func (d *Dao) LogTask(ctx context.Context, name string) (data *model.LogTask, err error) {
	row := d.db.QueryRow(ctx, _logTask, name)
	data = &model.LogTask{}
	if err = row.Scan(&data.ID, &data.Name, &data.Expect, &data.Success, &data.Failure, &data.State, &data.CTime, &data.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
		}
		return
	}
	return
}

// InsertLogTask .
func (d *Dao) InsertLogTask(ctx context.Context, data *model.LogTask) (id int64, err error) {
	result, err := d.db.Exec(ctx, _insertLogTask, data.Name, data.Expect, data.Success, data.Failure, data.State)
	if err != nil {
		return
	}
	id, err = result.LastInsertId()
	return
}

// TXIncrLogTaskSuccess .
func (d *Dao) TXIncrLogTaskSuccess(ctx context.Context, tx *xsql.Tx, name string) (rows int64, err error) {
	result, err := tx.Exec(_logTaskSuccessIncr, name)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

// IncrLogTaskFailure .
func (d *Dao) IncrLogTaskFailure(ctx context.Context, name string) (rows int64, err error) {
	result, err := d.db.Exec(ctx, _logTaskFailureIncr, name)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

// Asset .
func (d *Dao) Asset(ctx context.Context, oid int64, otype string, currency string) (data *model.Asset, err error) {
	row := d.db.QueryRow(ctx, _asset, oid, otype, currency)
	data = &model.Asset{}
	if err = row.Scan(&data.ID, &data.MID, &data.OID, &data.OType, &data.Currency, &data.Price, &data.State, &data.CTime, &data.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
		}
		return
	}
	return
}

// SumPaidOrderUserRealFee .
func (d *Dao) SumPaidOrderUserRealFee(ctx context.Context, beginTime, endTime time.Time) (sum int64, err error) {
	row := d.db.QueryRow(ctx, _sumPaidOrderUserRealFee, beginTime, endTime)
	if err = row.Scan(&sum); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			sum = 0
		}
		return
	}
	return
}

// SumRefundedOrderUserRealFee .
func (d *Dao) SumRefundedOrderUserRealFee(ctx context.Context, beginTime, endTime time.Time) (sum int64, err error) {
	row := d.db.QueryRow(ctx, _sumRefundedOrderUserRealFee, beginTime, endTime)
	if err = row.Scan(&sum); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			sum = 0
		}
		return
	}
	return
}

// SumDailyBill .
func (d *Dao) SumDailyBill(ctx context.Context, ver int64) (sumIn int64, sumOut int64, err error) {
	row := d.db.QueryRow(ctx, _sumDailyBill, ver)
	if err = row.Scan(&sumIn, &sumOut); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			sumIn = 0
			sumOut = 0
		}
		return
	}
	return
}

// MinIDOrderPaid .
func (d *Dao) MinIDOrderPaid(ctx context.Context, beginTime time.Time) (minID int64, err error) {
	row := d.db.QueryRow(ctx, _minIDOrderPaid, beginTime)
	if err = row.Scan(&minID); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			minID = math.MaxInt32
		}
		return
	}
	minID--
	return
}

// OrderPaidList order list
func (d *Dao) OrderPaidList(ctx context.Context, beginTime time.Time, endTime time.Time, fromID int64, limit int) (maxID int64, data []*model.Order, err error) {
	rows, err := d.db.Query(ctx, _orderPaidList, beginTime, endTime, fromID, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		d := &model.Order{}
		if err = rows.Scan(&d.ID, &d.MID, &d.OrderID, &d.Biz, &d.Platform, &d.OID, &d.OType, &d.Fee, &d.RealFee, &d.Currency, &d.PayID, &d.PayReason, &d.PayTime, &d.State, &d.CTime, &d.MTime, &d.RefundTime, &d.Version); err != nil {
			return
		}
		if d.ID > maxID {
			maxID = d.ID
		}
		data = append(data, d)
	}
	err = rows.Err()
	return
}

// MinIDOrderRefunded .
func (d *Dao) MinIDOrderRefunded(ctx context.Context, beginTime time.Time) (minID int64, err error) {
	row := d.db.QueryRow(ctx, _minIDOrderRefunded, beginTime)
	if err = row.Scan(&minID); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			minID = math.MaxInt32
		}
		return
	}
	minID--
	return
}

// OrderRefundedList order list
func (d *Dao) OrderRefundedList(ctx context.Context, beginTime time.Time, endTime time.Time, fromID int64, limit int) (maxID int64, data []*model.Order, err error) {
	rows, err := d.db.Query(ctx, _orderRefundedList, beginTime, endTime, fromID, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		d := &model.Order{}
		if err = rows.Scan(&d.ID, &d.MID, &d.OrderID, &d.Biz, &d.Platform, &d.OID, &d.OType, &d.Fee, &d.RealFee, &d.Currency, &d.PayID, &d.PayReason, &d.PayTime, &d.State, &d.CTime, &d.MTime, &d.RefundTime, &d.Version); err != nil {
			return
		}
		if d.ID > maxID {
			maxID = d.ID
		}
		data = append(data, d)
	}
	err = rows.Err()
	return
}

// TXUpdateOrder .
func (d *Dao) TXUpdateOrder(ctx context.Context, tx *xsql.Tx, order *model.Order) (rows int64, err error) {
	result, err := tx.Exec(_updateOrder, order.MID, order.OrderID, order.Biz, order.Platform, order.OID, order.OType, order.Fee, order.RealFee, order.Currency, order.PayID, order.PayReason, order.PayTime, order.State, order.RefundTime, order.ID, order.Version)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

// TXInsertOrderUserLog .
func (d *Dao) TXInsertOrderUserLog(ctx context.Context, tx *xsql.Tx, data *model.LogOrder) (id int64, err error) {
	result, err := tx.Exec(_insertLogOrderUser, data.OrderID, data.FromState, data.ToState, data.Desc)
	if err != nil {
		return
	}
	if id, err = result.LastInsertId(); err != nil {
		return
	}
	return
}

// MinIDDailyBillByVer .
func (d *Dao) MinIDDailyBillByVer(ctx context.Context, ver int64) (minID int64, err error) {
	row := d.db.QueryRow(ctx, _minIDDailyBillByVer, ver)
	if err = row.Scan(&minID); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			minID = math.MaxInt32
		}
		return
	}
	minID--
	return
}

// DailyBillListByVer bill list
func (d *Dao) DailyBillListByVer(ctx context.Context, ver int64, fromID int64, limit int) (maxID int64, data []*model.DailyBill, err error) {
	rows, err := d.db.Query(ctx, _dailyBillListByVer, ver, fromID, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		d := &model.DailyBill{}
		if err = rows.Scan(&d.ID, &d.BillID, &d.MID, &d.Biz, &d.Currency, &d.In, &d.Out, &d.Ver, &d.MonthVer, &d.CTime, &d.MTime, &d.Version); err != nil {
			return
		}
		if d.ID > maxID {
			maxID = d.ID
		}
		data = append(data, d)
	}
	err = rows.Err()
	return
}

// MinIDDailyBillByMonthVer .
func (d *Dao) MinIDDailyBillByMonthVer(ctx context.Context, monthVer int64) (minID int64, err error) {
	row := d.db.QueryRow(ctx, _minIDDailyBillByMonthVer, monthVer)
	if err = row.Scan(&minID); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			minID = math.MaxInt32
		}
		return
	}
	minID--
	return
}

// DailyBillListByMonthVer bill list
func (d *Dao) DailyBillListByMonthVer(ctx context.Context, monthVer int64, fromID int64, limit int) (maxID int64, data []*model.DailyBill, err error) {
	rows, err := d.db.Query(ctx, _dailyBillListByMonthVer, monthVer, fromID, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		d := &model.DailyBill{}
		if err = rows.Scan(&d.ID, &d.BillID, &d.MID, &d.Biz, &d.Currency, &d.In, &d.Out, &d.Ver, &d.MonthVer, &d.CTime, &d.MTime, &d.Version); err != nil {
			return
		}
		if d.ID > maxID {
			maxID = d.ID
		}
		data = append(data, d)
	}
	err = rows.Err()
	return
}

// TXInsertLogDailyBill .
func (d *Dao) TXInsertLogDailyBill(ctx context.Context, tx *xsql.Tx, log *model.LogBillDaily) (id int64, err error) {
	result, err := tx.Exec(_insertDailyBillLog, log.BillID, log.FromIn, log.ToIn, log.FromOut, log.ToOut, log.OrderID)
	if err != nil {
		return
	}
	id, err = result.LastInsertId()
	return
}

// TXInsertLogMonthlyBill .
func (d *Dao) TXInsertLogMonthlyBill(ctx context.Context, tx *xsql.Tx, log *model.LogBillMonthly) (id int64, err error) {
	result, err := tx.Exec(_insertMonthlyBillLog, log.BillID, log.FromIn, log.ToIn, log.FromOut, log.ToOut, log.BillUserDailyID)
	if err != nil {
		return
	}
	id, err = result.LastInsertId()
	return
}

// MinIDMonthlyBill .
func (d *Dao) MinIDMonthlyBill(ctx context.Context, ver int64) (minID int64, err error) {
	row := d.db.QueryRow(ctx, _minIDMonthlyBill, ver)
	if err = row.Scan(&minID); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			minID = math.MaxInt32
		}
		return
	}
	minID--
	return
}

// MonthlyBillList bill list
func (d *Dao) MonthlyBillList(ctx context.Context, ver int64, fromID int64, limit int) (maxID int64, data []*model.Bill, err error) {
	rows, err := d.db.Query(ctx, _monthlyBillList, ver, fromID, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		d := &model.Bill{}
		if err = rows.Scan(&d.ID, &d.BillID, &d.MID, &d.Biz, &d.Currency, &d.In, &d.Out, &d.Ver, &d.CTime, &d.MTime, &d.Version); err != nil {
			return
		}
		if d.ID > maxID {
			maxID = d.ID
		}
		data = append(data, d)
	}
	err = rows.Err()
	return
}

// DailyBill .
func (d *Dao) DailyBill(ctx context.Context, mid int64, biz string, currency string, ver int64) (data *model.DailyBill, err error) {
	row := d.db.QueryRow(ctx, _dailyBill, mid, biz, currency, ver)
	data = &model.DailyBill{}
	if err = row.Scan(&data.ID, &data.BillID, &data.MID, &data.Biz, &data.Currency, &data.In, &data.Out, &data.Ver, &data.MonthVer, &data.CTime, &data.MTime, &data.Version); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
		}
		return
	}
	return
}

// InsertDailyBill .
func (d *Dao) InsertDailyBill(ctx context.Context, bill *model.DailyBill) (id int64, err error) {
	result, err := d.db.Exec(ctx, _insertDailyBill, bill.BillID, bill.MID, bill.Biz, bill.Currency, bill.In, bill.Out, bill.Ver, bill.MonthVer, bill.Version)
	if err != nil {
		return
	}
	id, err = result.LastInsertId()
	return
}

// TXUpdateDailyBill .
func (d *Dao) TXUpdateDailyBill(ctx context.Context, tx *xsql.Tx, bill *model.DailyBill) (rows int64, err error) {
	result, err := tx.Exec(_updateDailyBill, bill.BillID, bill.MID, bill.Biz, bill.Currency, bill.In, bill.Out, bill.Ver, bill.MonthVer, bill.MID, bill.Biz, bill.Currency, bill.Ver, bill.Version)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

const (
	_upsertDeltaDailyBill = "INSERT INTO bill_user_daily (bill_id,mid,biz,currency,`in`,`out`,ver,month_ver) VALUES (?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `in`=`in`+?,`out`=`out`+?"
	_updateDeltaDailyBill = "UPDATE bill_user_daily SET `in`=`in`+?,`out`=`out`+? WHERE mid=? AND biz=? AND currency=? AND ver=?"
)

// TXUpsertDeltaDailyBill .
func (d *Dao) TXUpsertDeltaDailyBill(ctx context.Context, tx *xsql.Tx, bill *model.DailyBill, deltaIn, deltaOut int64) (rows int64, err error) {
	result, err := tx.Exec(_upsertDeltaDailyBill, bill.BillID, bill.MID, bill.Biz, bill.Currency, deltaIn, deltaOut, bill.Ver, bill.MonthVer, deltaIn, deltaOut)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

// TXUpdateDeltaDailyBill .
func (d *Dao) TXUpdateDeltaDailyBill(ctx context.Context, tx *xsql.Tx, deltaIn, deltaOut int64, mid int64, biz string, currency string, ver int64) (rows int64, err error) {
	result, err := tx.Exec(_updateDeltaDailyBill, deltaIn, deltaOut, mid, biz, currency, ver)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

// MonthlyBill .
func (d *Dao) MonthlyBill(ctx context.Context, mid int64, biz string, currency string, ver int64) (data *model.Bill, err error) {
	row := d.db.QueryRow(ctx, _monthlyBill, mid, biz, currency, ver)
	data = &model.Bill{}
	if err = row.Scan(&data.ID, &data.BillID, &data.MID, &data.Biz, &data.Currency, &data.In, &data.Out, &data.Ver, &data.CTime, &data.MTime, &data.Version); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
		}
		return
	}
	return
}

// InsertMonthlyBill .
func (d *Dao) InsertMonthlyBill(ctx context.Context, bill *model.Bill) (id int64, err error) {
	result, err := d.db.Exec(ctx, _insertMonthlyBill, bill.BillID, bill.MID, bill.Biz, bill.Currency, bill.In, bill.Out, bill.Ver, bill.Version)
	if err != nil {
		return
	}
	id, err = result.LastInsertId()
	return
}

// TXUpdateMonthlyBill .
func (d *Dao) TXUpdateMonthlyBill(ctx context.Context, tx *xsql.Tx, bill *model.Bill) (rows int64, err error) {
	result, err := tx.Exec(_updateMonthlyBill, bill.BillID, bill.MID, bill.Biz, bill.Currency, bill.In, bill.Out, bill.Ver, bill.MID, bill.Biz, bill.Currency, bill.Ver, bill.Version)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

const (
	_upsertDeltaMonthlyBill = "INSERT INTO bill_user_monthly (bill_id,mid,biz,currency,`in`,`out`,ver) VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `in`=`in`+?,`out`=`out`+?"
	_updateDeltaMonthlyBill = "UPDATE bill_user_monthly SET `in`=`in`+?,`out`=`out`+? WHERE mid=? AND biz=? AND currency=? AND ver=?"
)

// TXUpsertDeltaMonthlyBill .
func (d *Dao) TXUpsertDeltaMonthlyBill(ctx context.Context, tx *xsql.Tx, bill *model.Bill, deltaIn, deltaOut int64) (rows int64, err error) {
	result, err := tx.Exec(_upsertDeltaMonthlyBill, bill.BillID, bill.MID, bill.Biz, bill.Currency, deltaIn, deltaOut, bill.Ver, deltaIn, deltaOut)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

// TXUpdateDeltaMonthlyBill .
func (d *Dao) TXUpdateDeltaMonthlyBill(ctx context.Context, tx *xsql.Tx, deltaIn, deltaOut int64, mid int64, biz string, currency string, ver int64) (rows int64, err error) {
	result, err := tx.Exec(_updateDeltaMonthlyBill, deltaIn, deltaOut, mid, biz, currency, ver)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

// MinIDUserAccount .
func (d *Dao) MinIDUserAccount(ctx context.Context, beginTime time.Time) (minID int64, err error) {
	row := d.db.QueryRow(ctx, _minIDUserAccount, beginTime)
	if err = row.Scan(&minID); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			minID = math.MaxInt32
		}
		return
	}
	minID--
	return
}

// UserAccountList bill list
func (d *Dao) UserAccountList(ctx context.Context, beginTime time.Time, endTime time.Time, fromID int64, limit int) (maxID int64, datas []*model.UserAccount, err error) {
	rows, err := d.db.Query(ctx, _userAccountList, beginTime, endTime, fromID, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &model.UserAccount{}
		if err = rows.Scan(&data.ID, &data.MID, &data.Biz, &data.Currency, &data.Balance, &data.Ver, &data.State, &data.CTime, &data.MTime); err != nil {
			return
		}
		if data.ID > maxID {
			maxID = data.ID
		}
		datas = append(datas, data)
	}
	err = rows.Err()
	return
}

// UserAccount .
func (d *Dao) UserAccount(ctx context.Context, mid int64, biz string, currency string) (data *model.UserAccount, err error) {
	row := d.db.QueryRow(ctx, _userAccount, mid, biz, currency)
	data = &model.UserAccount{}
	if err = row.Scan(&data.ID, &data.MID, &data.Biz, &data.Currency, &data.Balance, &data.Ver, &data.State, &data.CTime, &data.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
		}
		return
	}
	return
}

// InsertUserAccount .
func (d *Dao) InsertUserAccount(ctx context.Context, account *model.UserAccount) (id int64, err error) {
	result, err := d.db.Exec(ctx, _insertUserAccount, account.MID, account.Biz, account.Currency, account.Balance, account.Ver, account.State)
	if err != nil {
		return
	}
	id, err = result.LastInsertId()
	return
}

// TXUpdateUserAccount .
func (d *Dao) TXUpdateUserAccount(ctx context.Context, tx *xsql.Tx, account *model.UserAccount) (rows int64, err error) {
	result, err := tx.Exec(_updateUserAccount, account.MID, account.Biz, account.Currency, account.Balance, account.State, account.MID, account.Biz, account.Currency, account.Ver)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

const (
	_upsertDeltaUserAccount = "INSERT INTO account_user (mid,biz,currency,balance,ver,state) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE balance=balance+?,ver=ver+1"
	_updateDeltaUserAccount = `UPDATE account_user SET balance=balance+?,ver=ver+1 WHERE mid=? AND biz=? AND currency=?`
)

// TXUpsertDeltaUserAccount .
func (d *Dao) TXUpsertDeltaUserAccount(ctx context.Context, tx *xsql.Tx, account *model.UserAccount, deltaBalance int64) (rows int64, err error) {
	result, err := tx.Exec(_upsertDeltaUserAccount, account.MID, account.Biz, account.Currency, deltaBalance, account.Ver, account.State, deltaBalance)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

// TXUpdateDeltaUserAccount .
func (d *Dao) TXUpdateDeltaUserAccount(ctx context.Context, tx *xsql.Tx, deltaBalance int64, mid int64, biz string, currency string) (rows int64, err error) {
	result, err := tx.Exec(_updateDeltaUserAccount, deltaBalance, mid, biz, currency)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

// TXInsertUserAccountLog .
func (d *Dao) TXInsertUserAccountLog(ctx context.Context, tx *xsql.Tx, accountLog *model.AccountLog) (err error) {
	_, err = tx.Exec(_insertUserAccountLog, accountLog.AccountID, accountLog.From, accountLog.To, accountLog.Ver, accountLog.State, accountLog.Name)
	return
}

// BizAccount .
func (d *Dao) BizAccount(ctx context.Context, biz string, currency string) (data *model.BizAccount, err error) {
	row := d.db.QueryRow(ctx, _bizAccount, biz, currency)
	data = &model.BizAccount{}
	if err = row.Scan(&data.ID, &data.Biz, &data.Currency, &data.Balance, &data.Ver, &data.State, &data.CTime, &data.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
		}
		return
	}
	return
}

// InsertBizAccount .
func (d *Dao) InsertBizAccount(ctx context.Context, account *model.BizAccount) (id int64, err error) {
	result, err := d.db.Exec(ctx, _insertBizAccount, account.Biz, account.Currency, account.Balance, account.Ver, account.State)
	if err != nil {
		return
	}
	id, err = result.LastInsertId()
	return
}

// TXUpdateBizAccount .
func (d *Dao) TXUpdateBizAccount(ctx context.Context, tx *xsql.Tx, account *model.BizAccount) (rows int64, err error) {
	result, err := tx.Exec(_updateBizAccount, account.Biz, account.Currency, account.Balance, account.State, account.Biz, account.Currency, account.Ver)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

const (
	_upsertDeltaBizAccount = "INSERT INTO account_biz (biz,currency,balance,ver,state) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE balance=balance+?,ver=ver+1"
	_updateDeltaBizAccount = `UPDATE account_biz SET balance=balance+?,ver=ver+1 WHERE biz=? AND currency=?`
)

// TXUpsertDeltaBizAccount .
func (d *Dao) TXUpsertDeltaBizAccount(ctx context.Context, tx *xsql.Tx, account *model.BizAccount, deltaBalance int64) (rows int64, err error) {
	result, err := tx.Exec(_upsertDeltaBizAccount, account.Biz, account.Currency, deltaBalance, account.Ver, account.State, deltaBalance)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

// TXUpdateDeltaBizAccount .
func (d *Dao) TXUpdateDeltaBizAccount(ctx context.Context, tx *xsql.Tx, deltaBalance int64, biz string, currency string) (rows int64, err error) {
	result, err := tx.Exec(_updateDeltaBizAccount, deltaBalance, biz, currency)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

// TXInsertBizAccountLog .
func (d *Dao) TXInsertBizAccountLog(ctx context.Context, tx *xsql.Tx, accountLog *model.AccountLog) (err error) {
	_, err = tx.Exec(_insertBizAccountLog, accountLog.AccountID, accountLog.From, accountLog.To, accountLog.Ver, accountLog.State, accountLog.Name)
	return
}

// AggrIncomeUser .
func (d *Dao) AggrIncomeUser(ctx context.Context, mid int64, currency string) (data *model.AggrIncomeUser, err error) {
	row := d.db.QueryRow(ctx, _aggrIncomeUser, mid, currency)
	data = &model.AggrIncomeUser{}
	if err = row.Scan(&data.ID, &data.MID, &data.Currency, &data.PaySuccess, &data.PayError, &data.TotalIn, &data.TotalOut, &data.CTime, &data.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
		}
		return
	}
	return
}

// InsertAggrIncomeUser .
func (d *Dao) InsertAggrIncomeUser(ctx context.Context, aggr *model.AggrIncomeUser) (id int64, err error) {
	result, err := d.db.Exec(ctx, _insertAggrIncomeUser, aggr.MID, aggr.Currency, aggr.PaySuccess, aggr.PayError, aggr.TotalIn, aggr.TotalOut)
	if err != nil {
		return
	}
	id, err = result.LastInsertId()
	return
}

// TXUpdateAggrIncomeUser .
func (d *Dao) TXUpdateAggrIncomeUser(ctx context.Context, tx *xsql.Tx, aggr *model.AggrIncomeUser) (rows int64, err error) {
	result, err := tx.Exec(_updateAggrIncomeUser, aggr.MID, aggr.Currency, aggr.PaySuccess, aggr.PayError, aggr.TotalIn, aggr.TotalOut, aggr.MID, aggr.Currency)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

const (
	_upsertDeltaAggrIncomeUser = "INSERT INTO aggr_income_user (mid,currency,pay_success,pay_error,total_in,total_out) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE pay_success=pay_success+?,pay_error=pay_error+?,total_in=total_in+?,total_out=total_out+?"
	_updateDeltaAggrIncomeUser = "UPDATE aggr_income_user SET pay_success=pay_success+?,pay_error=pay_error+?,total_in=total_in+?,total_out=total_out+? WHERE mid=? AND currency=?"
)

// TXUpsertDeltaAggrIncomeUser .
func (d *Dao) TXUpsertDeltaAggrIncomeUser(ctx context.Context, tx *xsql.Tx, aggr *model.AggrIncomeUser, deltaPaySuccess, deltaPayError, deltaTotalIn, deltaTotalOut int64) (rows int64, err error) {
	result, err := tx.Exec(_upsertDeltaAggrIncomeUser, aggr.MID, aggr.Currency, deltaPaySuccess, deltaPayError, deltaTotalIn, deltaTotalOut, deltaPaySuccess, deltaPayError, deltaTotalIn, deltaTotalOut)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

// TXUpdateDeltaAggrIncomeUser .
func (d *Dao) TXUpdateDeltaAggrIncomeUser(ctx context.Context, tx *xsql.Tx, deltaPaySuccess, deltaPayError, deltaTotalIn, deltaTotalOut int64, mid int64, currency string) (rows int64, err error) {
	result, err := tx.Exec(_updateDeltaAggrIncomeUser, deltaPaySuccess, deltaPayError, deltaTotalIn, deltaTotalOut, mid, currency)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

// AggrIncomeUserAsset .
func (d *Dao) AggrIncomeUserAsset(ctx context.Context, mid int64, currency string, ver int64, oid int64, otype string) (data *model.AggrIncomeUserAsset, err error) {
	row := d.db.QueryRow(ctx, _aggrIncomeUserAsset, mid, currency, ver, oid, otype)
	data = &model.AggrIncomeUserAsset{}
	if err = row.Scan(&data.ID, &data.MID, &data.Currency, &data.Ver, &data.OID, &data.OType, &data.PaySuccess, &data.PayError, &data.TotalIn, &data.TotalOut, &data.CTime, &data.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
		}
		return
	}
	return
}

// InsertAggrIncomeUserAsset .
func (d *Dao) InsertAggrIncomeUserAsset(ctx context.Context, aggr *model.AggrIncomeUserAsset) (id int64, err error) {
	result, err := d.db.Exec(ctx, _insertAggrIncomeUserAsset, aggr.MID, aggr.Currency, aggr.Ver, aggr.OID, aggr.OType, aggr.PaySuccess, aggr.PayError, aggr.TotalIn, aggr.TotalOut)
	if err != nil {
		return
	}
	id, err = result.LastInsertId()
	return
}

// TXUpdateAggrIncomeUserAsset .
func (d *Dao) TXUpdateAggrIncomeUserAsset(ctx context.Context, tx *xsql.Tx, aggr *model.AggrIncomeUserAsset) (rows int64, err error) {
	result, err := tx.Exec(_updateAggrIncomeUserAsset, aggr.MID, aggr.Currency, aggr.Ver, aggr.OID, aggr.OType, aggr.PaySuccess, aggr.PayError, aggr.TotalIn, aggr.TotalOut, aggr.MID, aggr.Currency, aggr.Ver, aggr.OID, aggr.OType)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

const (
	_upsertDeltaAggrIncomeUserAsset = "INSERT INTO aggr_income_user_asset (mid,currency,`ver`,oid,otype,pay_success,pay_error,total_in,total_out) VALUES (?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE pay_success=pay_success+?,pay_error=pay_error+?,total_in=total_in+?,total_out=total_out+?"
	_updateDeltaAggrIncomeUserAsset = "UPDATE aggr_income_user_asset SET pay_success=pay_success+?,pay_error=pay_error+?,total_in=total_in+?,total_out=total_out+? WHERE mid=? AND currency=? AND ver=? AND oid=? AND otype=?"
)

// TXUpsertDeltaAggrIncomeUserAsset .
func (d *Dao) TXUpsertDeltaAggrIncomeUserAsset(ctx context.Context, tx *xsql.Tx, aggr *model.AggrIncomeUserAsset, deltaPaySuccess, deltaPayError, deltaTotalIn, deltaTotalOut int64) (rows int64, err error) {
	result, err := tx.Exec(_upsertDeltaAggrIncomeUserAsset, aggr.MID, aggr.Currency, aggr.Ver, aggr.OID, aggr.OType, deltaPaySuccess, deltaPayError, deltaTotalIn, deltaTotalOut, deltaPaySuccess, deltaPayError, deltaTotalIn, deltaTotalOut)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

// TXUpdateDeltaAggrIncomeUserAsset .
func (d *Dao) TXUpdateDeltaAggrIncomeUserAsset(ctx context.Context, tx *xsql.Tx, deltaPaySuccess, deltaPayError, deltaTotalIn, deltaTotalOut int64, mid int64, currency string, ver int64, oid int64, otype string) (rows int64, err error) {
	result, err := tx.Exec(_updateDeltaAggrIncomeUserAsset, deltaPaySuccess, deltaPayError, deltaTotalIn, deltaTotalOut, mid, currency, ver, oid, otype)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

// OrderBadDebt .
func (d *Dao) OrderBadDebt(ctx context.Context, orderID string) (data *model.OrderBadDebt, err error) {
	row := d.db.QueryRow(ctx, _orderBadDebt, orderID)
	data = &model.OrderBadDebt{}
	if err = row.Scan(&data.ID, &data.OrderID, &data.Type, &data.State, &data.CTime, &data.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
		}
		return
	}
	return
}

// InsertOrderBadDebt .
func (d *Dao) InsertOrderBadDebt(ctx context.Context, order *model.OrderBadDebt) (id int64, err error) {
	result, err := d.db.Exec(ctx, _insertOrderBadDebt, order.OrderID, order.Type, order.State)
	if err != nil {
		return
	}
	id, err = result.LastInsertId()
	return
}

// TXUpdateOrderBadDebt .
func (d *Dao) TXUpdateOrderBadDebt(ctx context.Context, tx *xsql.Tx, order *model.OrderBadDebt) (rows int64, err error) {
	result, err := tx.Exec(_updateOrderBadDebt, order.OrderID, order.Type, order.State, order.OrderID)
	if err != nil {
		return
	}
	rows, err = result.RowsAffected()
	return
}

// TXInsertOrderRechargeShell .
func (d *Dao) TXInsertOrderRechargeShell(ctx context.Context, tx *xsql.Tx, order *model.OrderRechargeShell) (id int64, err error) {
	result, err := tx.Exec(_insertOrderRechargeShell, order.MID, order.OrderID, order.Biz, order.Amount, order.PayMSG, order.State, order.Ver)
	if err != nil {
		return
	}
	id, err = result.LastInsertId()
	return
}

// TXInsertOrderRechargeShellLog .
func (d *Dao) TXInsertOrderRechargeShellLog(ctx context.Context, tx *xsql.Tx, order *model.OrderRechargeShellLog) (id int64, err error) {
	result, err := tx.Exec(_insertOrderRechargeShellLog, order.OrderID, order.FromState, order.ToState, order.Desc, order.BillUserMonthlyID)
	if err != nil {
		return
	}
	id, err = result.LastInsertId()
	return
}
