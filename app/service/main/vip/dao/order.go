package dao

import (
	"context"
	xsql "database/sql"
	"fmt"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_vipOrderListSQL               = "SELECT id,order_no,app_id,platform,order_type,mid,to_mid,buy_months,money,refund_amount,status,pay_type,recharge_bp,third_trade_no,ver,payment_time,ctime,mtime,app_sub_id FROM vip_pay_order WHERE mid = ? AND  status = ? ORDER BY id DESC LIMIT ?,?;"
	_vipOrderCountSQL              = "SELECT COUNT(1) FROM vip_pay_order WHERE mid = ? AND status = ? ORDER BY id desc;"
	_getOrderInfoSQL               = "SELECT id,order_no, app_id,platform,order_type,mid,to_mid,buy_months,money,refund_amount,status,pay_type,recharge_bp,third_trade_no,payment_time,ver,ctime,mtime,app_sub_id  FROM  `vip_pay_order` WHERE `order_no` = ?;"
	_discountSQL                   = "SELECT mid,discount_id,status FROM vip_user_discount_history where mid=? AND discount_id = ? AND status = 1 "
	_allPriceMapping               = "SELECT id,month_id,month_type,money,selected,first_discount_money,discount_money,start_time,end_time,remark,operator,mtime FROM vip_month_price WHERE month_id = ? AND month_type = ? LIMIT 1;"
	_insertPayOrder                = "INSERT INTO vip_pay_order(order_no,app_id,platform,order_type,mid,to_mid,buy_months,money,status,pay_type,recharge_bp,third_trade_no,app_sub_id,coupon_money,pid,user_ip)VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
	_allMonthOrderSQL              = "SELECT id,month,month_type,operator,status,deleted,mtime FROM vip_month where deleted = 0 AND status = 1 ORDER BY month "
	_updateOrderStatusSQL          = "UPDATE vip_pay_order SET status = ?,pay_type = ?,third_trade_no = ? WHERE order_no = ?;"
	_updatePayOrderSQL             = "UPDATE vip_pay_order SET pay_type=?,third_trade_no=?,recharge_bp=?,ver=? WHERE id = ? AND ver=? "
	_updateIosPayOrderSQL          = "UPDATE vip_pay_order SET pay_type=?,third_trade_no=?,recharge_bp=?,ver=?,order_type=?,status=? WHERE id = ? AND ver=? "
	_updatePayOrderStatusSQL       = "UPDATE vip_pay_order SET status = ? WHERE id = ? AND ver=?;"
	_addOrderLogSQL                = "INSERT INTO vip_pay_order_log(order_no,refund_id,refund_amount,mid,status) VALUES(?,?,?,?,?);"
	_payOrderLastSQL               = "SELECT id,order_no,app_id,platform,order_type,mid,to_mid,buy_months,money,refund_amount,status,pay_type,recharge_bp,third_trade_no,payment_time,ver,ctime,mtime,app_sub_id FROM vip_pay_order WHERE mid = %v AND  status = %v AND order_type IN (%v) ORDER BY id DESC LIMIT 1 "
	_updatePayOrderRefundAmountSQL = "UPDATE vip_pay_order SET refund_amount=?,ver=? WHERE id=? AND ver=?"
	_selOldPayOrder                = "SELECT order_no,app_id,order_type,mid,IFNULL(buy_months,0),money,IFNULL(pay_type,4),IFNULL(status,1),ver,IFNULL(platform,3),mtime,app_sub_id,pid,user_ip from vip_pay_order WHERE order_no = ?;"
	_selOrderLogSQL                = "SELECT order_no,refund_id,mid,status FROM vip_pay_order_log WHERE order_no = ? AND refund_id = ? AND status = ?"
)

// BeginTran begin tran.
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	tx, err = d.db.Begin(c)
	return
}

//OrderCount order count.
func (d *Dao) OrderCount(c context.Context, mid int64, status int8) (count int64, err error) {
	var row *sql.Row
	if row = d.db.QueryRow(c, _vipOrderCountSQL, mid, status); err != nil {
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

// OrderList order list.
func (d *Dao) OrderList(c context.Context, mid int64, status int8, pn, ps int) (res []*model.PayOrder, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _vipOrderListSQL, mid, status, (pn-1)*ps, ps); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("query_db")
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.PayOrder)
		if err = rows.Scan(&r.ID, &r.OrderNo, &r.AppID, &r.Platform, &r.OrderType, &r.Mid, &r.ToMid, &r.BuyMonths, &r.Money, &r.RefundAmount, &r.Status, &r.PayType, &r.RechargeBp,
			&r.ThirdTradeNo, &r.Ver, &r.PaymentTime, &r.Ctime, &r.Mtime, &r.AppSubID); err != nil {
			err = errors.WithStack(err)
			d.errProm.Incr("row_scan_db")
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

//OrderInfo select order by order no.
func (d *Dao) OrderInfo(c context.Context, orderNo string) (r *model.OrderInfo, err error) {
	var row = d.db.QueryRow(c, _getOrderInfoSQL, orderNo)
	r = new(model.OrderInfo)
	if err = row.Scan(&r.ID, &r.OrderNo, &r.AppID, &r.Platform, &r.OrderType, &r.Mid, &r.ToMid, &r.BuyMonths, &r.Money, &r.RefundAmount, &r.Status, &r.PayType,
		&r.RechargeBP, &r.ThirdTradeNo, &r.PaymentTime, &r.Ver, &r.Ctime, &r.Mtime, &r.AppSubID); err != nil {
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

// DiscountSQL discount sql.
func (d *Dao) DiscountSQL(c context.Context, mid int64, discountID int64) (dh *model.VipUserDiscountHistory, err error) {
	row := d.db.QueryRow(c, _discountSQL, mid, discountID)
	dh = new(model.VipUserDiscountHistory)
	if err = row.Scan(&dh.Mid, &dh.DiscountID, &dh.Status); err != nil {
		if err == sql.ErrNoRows {
			dh = nil
			err = nil
		} else {
			err = errors.WithStack(err)
			d.errProm.Incr("row_scan_db")
		}
	}
	return
}

//PriceMapping all price mapping.
func (d *Dao) PriceMapping(c context.Context, monthID int64, platform int8) (r *model.PriceMapping, err error) {
	row := d.db.QueryRow(c, _allPriceMapping, monthID, platform)
	r = new(model.PriceMapping)
	if err = row.Scan(&r.ID, &r.MonthID, &r.MonthType, &r.Money, &r.Selected, &r.FirstDiscountMoney, &r.DiscountMoney,
		&r.StartTime, &r.EndTime, &r.Remark, &r.Operator, &r.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
		} else {
			err = errors.WithStack(err)
			d.errProm.Incr("row_scan_db")
		}
	}
	return
}

// TxAddOrder tx add order.
func (d *Dao) TxAddOrder(tx *sql.Tx, p *model.PayOrder) (id int64, err error) {
	var result xsql.Result
	if result, err = tx.Exec(_insertPayOrder, &p.OrderNo, &p.AppID, &p.Platform, &p.OrderType, &p.Mid, &p.ToMid,
		&p.BuyMonths, &p.Money, &p.Status, &p.PayType, &p.RechargeBp, &p.ThirdTradeNo, &p.AppSubID, &p.CouponMoney, &p.PID, &p.UserIP); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("exec_db")
		return
	}

	if id, err = result.LastInsertId(); err != nil {
		err = errors.WithStack(err)

	}
	return
}

//AllMonthByOrder order by month.
func (d *Dao) AllMonthByOrder(c context.Context, orderStr string) (res []*model.Month, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _allMonthOrderSQL+orderStr); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("query_db")
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Month)
		if err = rows.Scan(&r.ID, &r.Month, &r.MonthType, &r.Operator, &r.Status, &r.Deleted, &r.Mtime); err != nil {
			err = errors.WithStack(err)
			d.errProm.Incr("row_scan_db")
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// TxUpdateOrderStatus update order status.
func (d *Dao) TxUpdateOrderStatus(c context.Context, tx *sql.Tx, status int8, payType string, thirdTradeNO string, orderNO string) (err error) {
	if _, err = tx.Exec(_updateOrderStatusSQL, status, payType, thirdTradeNO, orderNO); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("exec_db")
		return
	}
	return
}

//TxUpdatePayOrder .
func (d *Dao) TxUpdatePayOrder(tx *sql.Tx, o *model.OrderInfo, ver int64) (err error) {
	if _, err = tx.Exec(_updatePayOrderSQL, o.PayType, o.ThirdTradeNo, o.RechargeBP, o.Ver, o.ID, ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//TxUpdateIosPayOrder .
func (d *Dao) TxUpdateIosPayOrder(tx *sql.Tx, o *model.OrderInfo, ver int64) (err error) {
	if _, err = tx.Exec(_updateIosPayOrderSQL, o.PayType, o.ThirdTradeNo, o.RechargeBP, o.Ver, o.OrderType, o.Status, o.ID, ver); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//TxUpdatePayOrderStatus .
func (d *Dao) TxUpdatePayOrderStatus(tx *sql.Tx, status int8, id int64, ver int64) (a int64, err error) {
	var result xsql.Result
	if result, err = tx.Exec(_updatePayOrderStatusSQL, status, id, ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = result.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//TxAddOrderLog .
func (d *Dao) TxAddOrderLog(tx *sql.Tx, arg *model.VipPayOrderLog) (err error) {
	if _, err = tx.Exec(_addOrderLogSQL, arg.OrderNo, arg.RefundID, arg.RefundAmount, arg.Mid, arg.Status); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//PayOrderLast .
func (d *Dao) PayOrderLast(c context.Context, mid int64, status int8, orderTypes ...int64) (r *model.PayOrder, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_payOrderLastSQL, mid, status, xstr.JoinInts(orderTypes)))
	r = new(model.PayOrder)
	if err = row.Scan(&r.ID, &r.OrderNo, &r.AppID, &r.Platform, &r.OrderType, &r.Mid, &r.ToMid, &r.BuyMonths, &r.Money, &r.RefundAmount, &r.Status, &r.PayType,
		&r.RechargeBp, &r.ThirdTradeNo, &r.PaymentTime, &r.Ver, &r.Ctime, &r.Mtime, &r.AppSubID); err != nil {
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

//SelOldPayOrder sel old payorder
func (d *Dao) SelOldPayOrder(c context.Context, orderNo string) (r *model.VipPayOrderOld, err error) {
	var row = d.olddb.QueryRow(c, _selOldPayOrder, orderNo)
	r = new(model.VipPayOrderOld)
	if err = row.Scan(&r.OrderNo, &r.AppID, &r.OrderType, &r.Mid, &r.BuyMonths, &r.Money, &r.PayType, &r.Status, &r.Ver, &r.Platform, &r.PaymentTime, &r.AppSubID, &r.PID, &r.UserIP); err != nil {
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

// SelPayOrderLog sel pay order log.
func (d *Dao) SelPayOrderLog(c context.Context, orderNo, refundID string, status int8) (res *model.VipPayOrderLog, err error) {
	row := d.db.QueryRow(c, _selOrderLogSQL, orderNo, refundID, status)
	res = new(model.VipPayOrderLog)
	if err = row.Scan(&res.OrderNo, &res.RefundID, &res.Mid, &res.Status); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		err = errors.WithStack(err)
		d.errProm.Incr("row_scan_db")
	}
	return
}

//TxUpdatePayOrderRefundAmount update payorder refund amount
func (d *Dao) TxUpdatePayOrderRefundAmount(tx *sql.Tx, id int64, refundAmount float64, ver, oldVer int64) (err error) {
	if _, err = tx.Exec(_updatePayOrderRefundAmountSQL, refundAmount, ver, id, oldVer); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("tx_db")
	}
	return
}
