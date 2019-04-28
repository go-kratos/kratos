package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

//OrderCount order count.
func (d *Dao) OrderCount(c context.Context, arg *model.ArgPayOrder) (count int64, err error) {
	if arg.Mid == 0 && len(arg.OrderNo) == 0 {
		return
	}
	sqlStr := _vipOrderCountSQL
	if arg.Mid > 0 {
		sqlStr += fmt.Sprintf(" AND mid = %v", arg.Mid)
	}
	if len(arg.OrderNo) > 0 {
		sqlStr += fmt.Sprintf(" AND order_no = '%v'", arg.OrderNo)
	}
	if arg.Status > 0 {
		sqlStr += fmt.Sprintf(" AND status = %v", arg.Status)
	}
	var row = d.db.QueryRow(c, sqlStr)
	if err = row.Scan(&count); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("row_scan_db")
		return
	}
	return
}

// OrderList order list.
func (d *Dao) OrderList(c context.Context, arg *model.ArgPayOrder) (res []*model.PayOrder, err error) {
	if arg.Mid == 0 && len(arg.OrderNo) == 0 {
		return
	}
	sqlStr := _vipOrderListSQL
	if arg.Mid > 0 {
		sqlStr += fmt.Sprintf(" AND mid = %v", arg.Mid)
	}
	if len(arg.OrderNo) > 0 {
		sqlStr += fmt.Sprintf(" AND order_no = '%v'", arg.OrderNo)
	}
	if arg.Status > 0 {
		sqlStr += fmt.Sprintf(" AND status = %v", arg.Status)
	}
	if arg.PN < 0 {
		arg.PN = _defpn
	}
	if arg.PS < 0 || arg.PS > 100 {
		arg.PS = _defps
	}
	sqlStr += fmt.Sprintf(" ORDER BY ID DESC LIMIT %v,%v", (arg.PN-1)*arg.PS, arg.PS)
	var rows *sql.Rows
	if rows, err = d.db.Query(c, sqlStr); err != nil {
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

// SelOrder sel order by orderNo.
func (d *Dao) SelOrder(c context.Context, orderNo string) (r *model.PayOrder, err error) {
	row := d.db.QueryRow(c, _vipOrderSQL, orderNo)
	r = new(model.PayOrder)
	if err = row.Scan(&r.ID, &r.OrderNo, &r.AppID, &r.Platform, &r.OrderType, &r.Mid, &r.ToMid, &r.BuyMonths, &r.Money, &r.RefundAmount, &r.Status, &r.PayType, &r.RechargeBp,
		&r.ThirdTradeNo, &r.Ver, &r.PaymentTime, &r.Ctime, &r.Mtime, &r.AppSubID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.WithStack(err)
		d.errProm.Incr("row_scan_db")
	}
	return
}

//AddPayOrderLog add order log.
func (d *Dao) AddPayOrderLog(c context.Context, arg *model.PayOrderLog) (err error) {
	if _, err = d.db.Exec(c, _addOrderLogSQL, arg.OrderNo, arg.RefundID, arg.RefundAmount, arg.Operator, arg.Mid, arg.Status); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_exec")
	}
	return
}
