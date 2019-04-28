package dao

import (
	"context"
	"database/sql"

	"go-common/app/service/main/ugcpay/model"
	xsql "go-common/library/database/sql"
)

var (
	_selectOrderUser    = `SELECT id,order_id,mid,biz,platform,oid,otype,fee,real_fee,currency,pay_id,pay_reason,pay_time,state,ctime,mtime,refund_time,version FROM order_user WHERE order_id=? LIMIT 1`
	_insertOrderUser    = "INSERT INTO order_user (order_id,mid,biz,platform,oid,otype,fee,real_fee,currency,pay_id,pay_reason,pay_time,state,refund_time,version) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_updateOrderUser    = "UPDATE order_user SET order_id=?,mid=?,biz=?,platform=?,oid=?,otype=?,fee=?,real_fee=?,currency=?,pay_id=?,pay_reason=?,pay_time=?,state=?,refund_time=?,version=version+1 WHERE order_id=? AND version=?"
	_insertLogOrderUser = "INSERT INTO log_order_user (order_id,from_state,to_state,`desc`) VALUES (?,?,?,?)"

	_selectAsset = "SELECT id,mid,oid,otype,currency,price,state,ctime,mtime FROM asset WHERE oid=? AND otype=? AND currency=? LIMIT 1"
	_upsertAsset = "INSERT INTO asset (mid,oid,otype,currency,price,state) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE mid=?,price=?,state=?"

	_selectAssetRelation = "SELECT id,oid,otype,mid,state,ctime,mtime FROM asset_relation WHERE oid=? AND otype=? AND mid=? LIMIT 1"
	_upsertAssetRelation = "INSERT INTO asset_relation (oid,otype,mid,state) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE state=?"

	_selectBillUserDaily               = "SELECT id,mid,biz,currency,`in`,`out`,ver,ctime,mtime,version FROM bill_user_daily WHERE mid=? AND biz=? AND currency=? AND ver=? LIMIT 1"
	_selectBillUserDailyListByMonthVer = "SELECT id,mid,biz,currency,`in`,`out`,ver,ctime,mtime,version FROM bill_user_daily WHERE mid=? AND biz=? AND currency=? AND month_ver=?"

	_selectAggrIncomeUser      = "SELECT id,mid,currency,pay_success,pay_error,total_in,total_out,ctime,mtime FROM aggr_income_user WHERE mid=? AND currency=? LIMIT 1"
	_selectAggrIncomeAssetList = "SELECT id,mid,currency,ver,oid,otype,pay_success,pay_error,total_in,total_out,ctime,mtime FROM aggr_income_user_asset WHERE mid=? AND currency=? AND ver=? ORDER BY oid DESC LIMIT ?"
	_selectAggrIncomeAsset     = "SELECT id,mid,currency,ver,oid,otype,pay_success,pay_error,total_in,total_out,ctime,mtime FROM aggr_income_user_asset WHERE mid=? AND currency=? AND ver=? AND oid=? AND otype=? LIMIT 1"

	_orderRechargeShell          = "SELECT id,mid,order_id,biz,amount,pay_msg,state,`ver`,ctime,mtime FROM  order_recharge_shell WHERE order_id=?"
	_updateOrderRechargeShell    = "UPDATE order_recharge_shell SET mid=?,order_id=?,biz=?,amount=?,pay_msg=?,state=?,`ver`=? WHERE order_id=?"
	_insertOrderRechargeShellLog = "INSERT INTO log_order_recharge_shell (order_id,from_state,to_state,`desc`,bill_user_monthly_id) VALUES (?,?,?,?,?)"
)

// BeginTran begin transcation.
func (d *Dao) BeginTran(c context.Context) (tx *xsql.Tx, err error) {
	return d.db.Begin(c)
}

// RawOrderUser get user order
func (d *Dao) RawOrderUser(ctx context.Context, id string) (data *model.Order, err error) {
	data = &model.Order{}
	row := d.db.Master().QueryRow(ctx, _selectOrderUser, id)
	if err = row.Scan(&data.ID, &data.OrderID, &data.MID, &data.Biz, &data.Platform, &data.OID, &data.OType, &data.Fee, &data.RealFee, &data.Currency, &data.PayID, &data.PayReason, &data.PayTime, &data.State, &data.CTime, &data.MTime, &data.RefundTime, &data.Version); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
			return
		}
		return
	}
	return
}

// InsertOrderUser is.
func (d *Dao) InsertOrderUser(ctx context.Context, data *model.Order) (id int64, err error) {
	var (
		res sql.Result
	)
	if res, err = d.db.Exec(ctx, _insertOrderUser, data.OrderID, data.MID, data.Biz, data.Platform, data.OID, data.OType, data.Fee, data.RealFee, data.Currency, data.PayID, data.PayReason, data.PayTime, data.State, data.RefundTime, data.Version); err != nil {
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		return
	}
	return
}

// TXUpdateOrderUser .
func (d *Dao) TXUpdateOrderUser(ctx context.Context, tx *xsql.Tx, data *model.Order) (affected int64, err error) {
	res, err := tx.Exec(_updateOrderUser, data.OrderID, data.MID, data.Biz, data.Platform, data.OID, data.OType, data.Fee, data.RealFee, data.Currency, data.PayID, data.PayReason, data.PayTime, data.State, data.RefundTime, data.OrderID, data.Version)
	if err != nil {
		return
	}
	affected, err = res.RowsAffected()
	return
}

// TXInsertOrderUserLog .
func (d *Dao) TXInsertOrderUserLog(ctx context.Context, tx *xsql.Tx, data *model.LogOrder) (id int64, err error) {
	var (
		res sql.Result
	)
	if res, err = tx.Exec(_insertLogOrderUser, data.OrderID, data.FromState, data.ToState, data.Desc); err != nil {
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		return
	}
	return
}

// RawAsset is
func (d *Dao) RawAsset(ctx context.Context, oid int64, otype string, currency string) (data *model.Asset, err error) {
	data = &model.Asset{}
	row := d.db.Master().QueryRow(ctx, _selectAsset, oid, otype, currency)
	if err = row.Scan(&data.ID, &data.MID, &data.OID, &data.OType, &data.Currency, &data.Price, &data.State, &data.CTime, &data.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
			return
		}
		return
	}
	return
}

// UpsertAsset is
func (d *Dao) UpsertAsset(ctx context.Context, data *model.Asset) (err error) {
	if _, err = d.db.Exec(ctx, _upsertAsset, data.MID, data.OID, data.OType, data.Currency, data.Price, data.State, data.MID, data.Price, data.State); err != nil {
		return
	}
	return
}

// RawAssetRelation is
func (d *Dao) RawAssetRelation(ctx context.Context, mid int64, oid int64, otype string) (data *model.AssetRelation, err error) {
	data = &model.AssetRelation{}
	row := d.db.Master().QueryRow(ctx, _selectAssetRelation, oid, otype, mid)
	if err = row.Scan(&data.ID, &data.OID, &data.OType, &data.MID, &data.State, &data.CTime, &data.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
			return
		}
		return
	}
	return
}

// UpsertAssetRelation is
func (d *Dao) UpsertAssetRelation(ctx context.Context, data *model.AssetRelation) (rows int64, err error) {
	var (
		result sql.Result
	)
	if result, err = d.db.Exec(ctx, _upsertAssetRelation, data.OID, data.OType, data.MID, data.State, data.State); err != nil {
		return
	}
	if rows, err = result.RowsAffected(); err != nil {
		return
	}
	return
}

// TXUpsertAssetRelation is
func (d *Dao) TXUpsertAssetRelation(ctx context.Context, tx *xsql.Tx, data *model.AssetRelation) (rows int64, err error) {
	var (
		result sql.Result
	)
	if result, err = tx.Exec(_upsertAssetRelation, data.OID, data.OType, data.MID, data.State, data.State); err != nil {
		return
	}
	if rows, err = result.RowsAffected(); err != nil {
		return
	}
	return
}

// RawAggrIncomeUser is.
func (d *Dao) RawAggrIncomeUser(ctx context.Context, mid int64, currency string) (data *model.AggrIncomeUser, err error) {
	data = &model.AggrIncomeUser{}
	row := d.db.QueryRow(ctx, _selectAggrIncomeUser, mid, currency)
	if err = row.Scan(&data.ID, &data.MID, &data.Currency, &data.PaySuccess, &data.PayError, &data.TotalIn, &data.TotalOut, &data.CTime, &data.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
			return
		}
		return
	}
	return
}

// RawAggrIncomeUserAssetList is.
func (d *Dao) RawAggrIncomeUserAssetList(ctx context.Context, mid int64, currency string, ver int64, limit int) (data []*model.AggrIncomeUserAsset, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(ctx, _selectAggrIncomeAssetList, mid, currency, ver, limit); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			d = &model.AggrIncomeUserAsset{}
		)
		if err = rows.Scan(&d.ID, &d.MID, &d.Currency, &d.Ver, &d.OID, &d.OType, &d.PaySuccess, &d.PayError, &d.TotalIn, &d.TotalOut, &d.CTime, &d.MTime); err != nil {
			return
		}
		data = append(data, d)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

// RawAggrIncomeUserAsset .
func (d *Dao) RawAggrIncomeUserAsset(ctx context.Context, mid int64, currency string, oid int64, otype string, ver int64) (data *model.AggrIncomeUserAsset, err error) {
	data = &model.AggrIncomeUserAsset{}
	row := d.db.QueryRow(ctx, _selectAggrIncomeAsset, mid, currency, ver, oid, otype)
	if err = row.Scan(&data.ID, &data.MID, &data.Currency, &data.Ver, &data.OID, &data.OType, &data.PaySuccess, &data.PayError, &data.TotalIn, &data.TotalOut, &data.CTime, &data.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
			return
		}
		return
	}
	return
}

// RawBillUserDaily is.
func (d *Dao) RawBillUserDaily(ctx context.Context, mid int64, biz string, currency string, ver int64) (data *model.Bill, err error) {
	data = &model.Bill{}
	row := d.db.QueryRow(ctx, _selectBillUserDaily, mid, biz, currency, ver)
	if err = row.Scan(&data.ID, &data.MID, &data.Biz, &data.Currency, &data.In, &data.Out, &data.Ver, &data.CTime, &data.MTime, &data.Version); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
			return
		}
		return
	}
	return
}

// RawBillUserDailyByMonthVer .
func (d *Dao) RawBillUserDailyByMonthVer(ctx context.Context, mid int64, biz string, currency string, monthVer int64) (datas []*model.Bill, err error) {
	rows, err := d.db.Query(ctx, _selectBillUserDailyListByMonthVer, mid, biz, currency, monthVer)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		data := &model.Bill{}
		if err = rows.Scan(&data.ID, &data.MID, &data.Biz, &data.Currency, &data.In, &data.Out, &data.Ver, &data.CTime, &data.MTime, &data.Version); err != nil {
			return
		}
		datas = append(datas, data)
	}
	err = rows.Err()
	return
}

// RawAccountUser is.
// func (d *Dao) RawAccountUser(ctx context.Context, mid int64, biz string, currency string) (data *model.AccountUser, err error) {
// 	data = &model.AccountUser{}
// 	row := d.db.QueryRow(ctx, _selectAccountUser, mid, biz, currency)
// 	if err = row.Scan(&data.ID, &data.Biz, &data.MID, &data.Currency, &data.Balance, &data.Ver, &data.State, &data.CTime, &data.MTime); err != nil {
// 		if err == xsql.ErrNoRows {
// 			err = nil
// 			data = nil
// 			return
// 		}
// 		return
// 	}
// 	return
// }

// RawOrderRechargeShell .
func (d *Dao) RawOrderRechargeShell(ctx context.Context, orderID string) (data *model.OrderRechargeShell, err error) {
	data = &model.OrderRechargeShell{}
	row := d.db.QueryRow(ctx, _orderRechargeShell, orderID)
	if err = row.Scan(&data.ID, &data.MID, &data.OrderID, &data.Biz, &data.Amount, &data.PayMSG, &data.State, &data.Ver, &data.CTime, &data.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			data = nil
			return
		}
		return
	}
	return
}

// TXUpdateOrderRechargeShell .
func (d *Dao) TXUpdateOrderRechargeShell(ctx context.Context, tx *xsql.Tx, data *model.OrderRechargeShell) (err error) {
	if _, err = tx.Exec(_updateOrderRechargeShell, data.MID, data.OrderID, data.Biz, data.Amount, data.PayMSG, data.State, data.Ver, data.OrderID); err != nil {
		return
	}
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
