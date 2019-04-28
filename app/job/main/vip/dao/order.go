package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"go-common/app/job/main/vip/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_insertPayOrder          = "INSERT IGNORE INTO vip_pay_order(order_no,app_id,platform,order_type,mid,to_mid,buy_months,money,status,pay_type,recharge_bp,third_trade_no,payment_time,ver,app_sub_id,coupon_money) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_selPayOrder             = "SELECT order_no,app_id,platform,order_type,mid,to_mid,buy_months,money,status,pay_type,third_trade_no,payment_time,ver,app_sub_id FROM vip_pay_order WHERE id>? AND id <=?"
	_selPayOrderByMidSQL     = "SELECT order_no,app_id,platform,order_type,mid,to_mid,buy_months,money,status,pay_type,third_trade_no,payment_time,ver,app_sub_id FROM vip_pay_order WHERE mid=? AND order_type=? AND status=? ORDER BY ID DESC LIMIT 1"
	_selOrderByMidSQL        = "SELECT order_no,app_id,platform,order_type,mid,to_mid,buy_months,money,status,pay_type,third_trade_no,payment_time,ver,app_sub_id FROM vip_pay_order WHERE order_no=?"
	_selPayOrderLogByMidSQL  = "SELECT order_no,status,mid FROM vip_pay_order_log WHERE mid=? AND  status=? ORDER BY ID DESC LIMIT 1"
	_selOldPayOrder          = "SELECT order_no,app_id,order_type,mid,IFNULL(buy_months,0),money,IFNULL(pay_type,4),IFNULL(status,1),ver,IFNULL(platform,3),mtime,app_sub_id,bmid,coupon_money from vip_pay_order WHERE id>? AND id<=?"
	_selOldRechargeOrder     = "SELECT pay_order_no,third_trade_no,recharge_bp FROM vip_recharge_order WHERE pay_order_no IN "
	_updatePayOrderStatusSQL = "UPDATE vip_pay_order SET "
	_updateRechageOrderSQL   = "UPDATE vip_pay_order SET recharge_bp = ?,third_trade_no = ? where order_no = ?"
	_insertPayOrderLog       = "INSERT INTO vip_pay_order_log(order_no,mid,status) VALUES(?,?,?)"
	_batchAddPayOrder        = "INSERT INTO vip_pay_order(order_no,app_id,platform,order_type,mid,to_mid,buy_months,money,status,pay_type,recharge_bp,third_trade_no,payment_time,ver,app_sub_id) VALUES"

	_selOrderMaxIDSQL    = "SELECT IFNULL(MAX(id),0) FROM vip_pay_order"
	_selOldOrderMaxIDSQL = "SELECT IFNULL(MAX(id),0) FROM vip_pay_order"
)

//SelPayOrderByMid sel payorder by mid
func (d *Dao) SelPayOrderByMid(c context.Context, mid int64, orderType, status int8) (r *model.VipPayOrder, err error) {
	row := d.db.QueryRow(c, _selPayOrderByMidSQL, mid, orderType, status)
	r = new(model.VipPayOrder)
	if err = row.Scan(&r.OrderNo, &r.AppID, &r.Platform, &r.OrderType, &r.Mid, &r.ToMid, &r.BuyMonths, &r.Money, &r.Status, &r.PayType, &r.ThirdTradeNo, &r.PaymentTime, &r.Ver, &r.AppSubID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.WithStack(err)
		d.errProm.Incr("db_scan")
		return
	}
	return
}

//SelOrderByOrderNo sel payorder by orderNo
func (d *Dao) SelOrderByOrderNo(c context.Context, orderNo string) (r *model.VipPayOrder, err error) {
	row := d.db.QueryRow(c, _selOrderByMidSQL, orderNo)
	r = new(model.VipPayOrder)
	if err = row.Scan(&r.OrderNo, &r.AppID, &r.Platform, &r.OrderType, &r.Mid, &r.ToMid, &r.BuyMonths, &r.Money, &r.Status, &r.PayType, &r.ThirdTradeNo, &r.PaymentTime, &r.Ver, &r.AppSubID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.WithStack(err)
		d.errProm.Incr("db_scan")
		return
	}
	return
}

//SelPayOrderLog sel pay order log.
func (d *Dao) SelPayOrderLog(c context.Context, mid int64, status int8) (r *model.VipPayOrderLog, err error) {
	row := d.db.QueryRow(c, _selPayOrderLogByMidSQL, mid, status)
	r = new(model.VipPayOrderLog)
	if err = row.Scan(&r.OrderNo, &r.Status, &r.Mid); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.WithStack(err)
		d.errProm.Incr("db_scan")
		return
	}
	return
}

//AddPayOrder add payorder
func (d *Dao) AddPayOrder(c context.Context, r *model.VipPayOrder) (a int64, err error) {
	var result sql.Result
	if result, err = d.db.Exec(c, _insertPayOrder, &r.OrderNo, &r.AppID, &r.Platform, &r.OrderType, &r.Mid, &r.ToMid, &r.BuyMonths, &r.Money, &r.Status, &r.PayType, &r.RechargeBp, &r.ThirdTradeNo, &r.PaymentTime, &r.Ver, &r.AppSubID, &r.CouponMoney); err != nil {
		log.Error("AddPayOrder d.db.exec(%v) error(%v)", r, err)
		return
	}
	if a, err = result.RowsAffected(); err != nil {
		log.Error("AddPayOrder result.RowsAffected() error(%v)", err)
		return
	}
	return
}

//SelOrderMaxID sel order maxID
func (d *Dao) SelOrderMaxID(c context.Context) (maxID int, err error) {
	row := d.db.QueryRow(c, _selOrderMaxIDSQL)
	if err = row.Scan(&maxID); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_scan")
		return
	}
	return
}

//SelOldOrderMaxID sel old order maxID
func (d *Dao) SelOldOrderMaxID(c context.Context) (maxID int, err error) {
	row := d.oldDb.QueryRow(c, _selOldOrderMaxIDSQL)
	if err = row.Scan(&maxID); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_scan")
		return
	}
	return
}

//BatchAddPayOrder batch add pay order
func (d *Dao) BatchAddPayOrder(c context.Context, res []*model.VipPayOrder) (err error) {
	var values []string
	if len(res) == 0 {
		return
	}
	for _, v := range res {
		value := fmt.Sprintf("('%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v')", v.OrderNo, v.AppID, v.Platform, v.OrderType, v.Mid, v.ToMid, v.BuyMonths, v.Money, v.Status, v.PayType, v.RechargeBp, v.ThirdTradeNo, v.PaymentTime.Time().Format("2006-01-02 15:04:05"), v.Ver, v.AppSubID, v.CouponMoney)
		values = append(values, value)
	}
	valuesStr := strings.Join(values, ",")
	dupStr := " ON DUPLICATE KEY UPDATE app_id = VALUES(app_id),platform=VALUES(platform),order_type=VALUES(order_type)," +
		"mid=VALUES(mid),to_mid=VALUES(to_mid),buy_months=VALUES(buy_months),money=VALUES(money),status=VALUES(status),pay_type=VALUES(pay_type),recharge_bp=VALUES(recharge_bp)," +
		"third_trade_no=VALUES(third_trade_no) ,payment_time=VALUES(payment_time) ,ver=VALUES(ver),app_sub_id=VALUES(app_sub_id),coupon_money=VALUES(coupon_money) "
	if _, err = d.db.Exec(c, _batchAddPayOrder+valuesStr+dupStr); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//UpdatePayOrderStatus update payorder status
func (d *Dao) UpdatePayOrderStatus(c context.Context, r *model.VipPayOrder) (a int64, err error) {
	var result sql.Result
	sqlStr := _updatePayOrderStatusSQL
	if r.PayType != 0 {
		sqlStr += fmt.Sprintf(" pay_type=%v ,payment_time='%v' ,", r.PayType, r.PaymentTime.Time().Format("2006-01-02 15:04:05"))
	}
	sqlStr += fmt.Sprintf(" ctime='%v',mtime='%v',ver=%v,status=%v,order_type=%v,coupon_money=%v WHERE order_no='%v'  ", r.Ctime.Time().Format("2006-01-02 15:04:05"), r.Mtime.Time().Format("2006-01-02 15:04:05"), r.Ver, r.Status, r.OrderType, r.CouponMoney, r.OrderNo)
	if result, err = d.db.Exec(c, sqlStr); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = result.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//UpdateRechargeOrder update recharge order info
func (d *Dao) UpdateRechargeOrder(c context.Context, r *model.VipPayOrder) (a int64, err error) {
	var result sql.Result
	if result, err = d.db.Exec(c, _updateRechageOrderSQL, &r.RechargeBp, &r.ThirdTradeNo, &r.OrderNo); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = result.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//AddPayOrderLog add payorderlog
func (d *Dao) AddPayOrderLog(c context.Context, r *model.VipPayOrderLog) (a int64, err error) {
	var result sql.Result
	if result, err = d.db.Exec(c, _insertPayOrderLog, &r.OrderNo, &r.Mid, &r.Status); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = result.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//SelPayOrder sel payorder
func (d *Dao) SelPayOrder(c context.Context, sID, eID int) (res []*model.VipPayOrder, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _selPayOrder, sID, eID); err != nil {
		log.Error("SelPayOrder d.db.query(sID:%v,eID:%v) error(%v)", sID, eID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipPayOrder)
		if err = rows.Scan(&r.OrderNo, &r.AppID, &r.Platform, &r.OrderType, &r.Mid, &r.ToMid, &r.BuyMonths, &r.Money, &r.Status, &r.PayType, &r.ThirdTradeNo, &r.PaymentTime, &r.Ver, &r.AppSubID); err != nil {
			log.Error("SelPayOrder rows.scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

//SelOldPayOrder sel old payorder
func (d *Dao) SelOldPayOrder(c context.Context, sID, eID int) (res []*model.VipPayOrderOld, err error) {
	var rows *xsql.Rows
	if rows, err = d.oldDb.Query(c, _selOldPayOrder, sID, eID); err != nil {
		log.Error("SelOldPayOrder d.db.query(sID:%v,eID:%v) error(%v)", sID, eID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipPayOrderOld)
		if err = rows.Scan(&r.OrderNo, &r.AppID, &r.OrderType, &r.Mid, &r.BuyMonths, &r.Money, &r.PayType, &r.Status, &r.Ver, &r.Platform, &r.PaymentTime, &r.AppSubID, &r.Bmid, &r.CouponMoney); err != nil {
			log.Error("SelOldPayOrder rows.scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

//SelOldRechargeOrder sel old rechargeOrder
func (d *Dao) SelOldRechargeOrder(c context.Context, orderNos []string) (res []*model.VipRechargeOrder, err error) {
	var inStr = strings.Join(orderNos, "','")
	inStr = "('" + inStr + "')"

	var rows *xsql.Rows
	if rows, err = d.oldDb.Query(c, _selOldRechargeOrder+inStr); err != nil {
		log.Error("SelOldRechargeOrder d.db.query(%v) error(%v)", orderNos, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipRechargeOrder)
		if err = rows.Scan(&r.PayOrderNo, &r.ThirdTradeNo, &r.RechargeBp); err != nil {
			log.Error("SelOldRechargeOrder rows.scan(%v) error(%v)", orderNos, err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}
