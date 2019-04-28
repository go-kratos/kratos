package dao

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/service/main/tv/internal/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	_getPayOrderByID                     = "SELECT `id`, `order_no`, `platform`, `order_type`, `mid`, `buy_months`, `product_id`, `money`, `quantity`, `refund_amount`, `status`, `third_trade_no`, `payment_money`, `payment_type`, `payment_time`, `ver`, `token`, `ctime`, `mtime` FROM `tv_pay_order` WHERE `id`=?"
	_getPayOrderByOrderNo                = "SELECT `id`, `order_no`, `platform`, `order_type`, `mid`, `buy_months`, `product_id`, `money`, `quantity`, `refund_amount`, `status`, `third_trade_no`, `payment_money`, `payment_type`, `payment_time`, `ver`, `token`, `ctime`, `mtime` FROM `tv_pay_order` WHERE `order_no`=?"
	_getPayOrdersByMid                   = "SELECT `id`, `order_no`, `platform`, `order_type`, `mid`, `buy_months`, `product_id`, `money`, `quantity`, `refund_amount`, `status`, `third_trade_no`, `payment_money`, `payment_type`, `payment_time`, `ver`, `token`, `ctime`, `mtime` FROM `tv_pay_order` WHERE `mid`=? ORDER BY `ctime` DESC LIMIT ?,?"
	_countPayOrderByMid                  = "SELECT count(*) FROM `tv_pay_order` WHERE `mid`=?"
	_getPayOrdersByMidAndStatus          = "SELECT `id`, `order_no`, `platform`, `order_type`, `mid`, `buy_months`, `product_id`, `money`, `quantity`, `refund_amount`, `status`, `third_trade_no`, `payment_money`, `payment_type`, `payment_time`, `ver`, `token`, `ctime`, `mtime` FROM `tv_pay_order` WHERE `mid`=? AND `status`=? ORDER BY `ctime` DESC LIMIT ?,?"
	_countPayOrderByMidAndStatus         = "SELECT count(*) FROM `tv_pay_order` WHERE `mid`=? AND `status`=?"
	_getPayOrdersByMidAndStatusAndCtime  = "SELECT `id`, `order_no`, `platform`, `order_type`, `mid`, `buy_months`, `product_id`, `money`, `quantity`, `refund_amount`, `status`, `third_trade_no`, `payment_money`, `payment_type`, `payment_time`, `ver`, `ctime`, `mtime`, `token` FROM `tv_pay_order` WHERE `mid`=? AND `status`=? AND `ctime`>= ? AND `ctime` <= ? ORDER BY `ctime` DESC LIMIT ?,?"
	_countPayOrderByMidAndStatusAndCtime = "SELECT count(*) FROM `tv_pay_order` WHERE `mid`=? AND `status`=? AND `ctime`>= ? AND `ctime` <= ?"

	_insertPayOrder = "INSERT INTO tv_pay_order (`order_no`, `platform`, `order_type`, `mid`, `buy_months`, `product_id`, `money`, `quantity`, `refund_amount`, `status`, `third_trade_no`, `payment_money`, `payment_type`, `payment_time`, `ver`, `token`, `app_channel`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	_updatePayOrder = "UPDATE `tv_pay_order` SET `status`=?, `payment_time`=?, `ver` = `ver` + 1 WHERE `id`=?  AND `ver`=?"

	_getUnpaidNoCallbackOrder = "SELECT `id`, `order_no`, `platform`, `order_type`, `mid`, `buy_months`, `product_id`, `money`, `quantity`, `refund_amount`, `status`, `third_trade_no`, `payment_money`, `payment_type`, `payment_time`, `ver`, `token`, `ctime`, `mtime` FROM `tv_pay_order` WHERE `status` = 1 and ctime > ? and ctime < ? order by id LIMIT ?,?"
)

// PayOrderByID quires one row from tv_pay_order.
func (d *Dao) PayOrderByID(c context.Context, id int) (po *model.PayOrder, err error) {
	row := d.db.QueryRow(c, _getPayOrderByID, id)
	po = &model.PayOrder{}
	err = row.Scan(&po.ID, &po.OrderNo, &po.Platform, &po.OrderType, &po.Mid, &po.BuyMonths, &po.ProductId, &po.Money, &po.Quantity, &po.RefundAmount, &po.Status, &po.ThirdTradeNo, &po.PaymentMoney, &po.PaymentType, &po.PaymentTime, &po.Ver, &po.Token, &po.Ctime, &po.Mtime)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Error("rows.Scan(%s) error(%v)", _getPayOrderByID, err)
		err = errors.WithStack(err)
		return nil, err
	}
	return po, nil
}

// PayOrderByOrderNo quires one row from tv_pay_order.
func (d *Dao) PayOrderByOrderNo(c context.Context, orderNo string) (po *model.PayOrder, err error) {
	row := d.db.QueryRow(c, _getPayOrderByOrderNo, orderNo)
	po = &model.PayOrder{}
	err = row.Scan(&po.ID, &po.OrderNo, &po.Platform, &po.OrderType, &po.Mid, &po.BuyMonths, &po.ProductId, &po.Money, &po.Quantity, &po.RefundAmount, &po.Status, &po.ThirdTradeNo, &po.PaymentMoney, &po.PaymentType, &po.PaymentTime, &po.Ver, &po.Token, &po.Ctime, &po.Mtime)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Error("rows.Scan(%s) error(%v)", _getPayOrderByOrderNo, err)
		err = errors.WithStack(err)
		return nil, err
	}
	return po, nil
}

// PayOrdersByMid quires rows from tv_pay_order.
func (d *Dao) PayOrdersByMid(c context.Context, mid int, pn, ps int) (res []*model.PayOrder, total int, err error) {
	res = make([]*model.PayOrder, 0)
	totalRow := d.db.QueryRow(c, _countPayOrderByMid, mid)
	if err = totalRow.Scan(&total); err != nil {
		log.Error("row.ScanCount error(%v)", err)
		err = errors.WithStack(err)
		return
	}
	rows, err := d.db.Query(c, _getPayOrdersByMid, mid, (pn-1)*ps, ps)
	if err != nil {
		log.Error("db.Query(%s) error(%v)", _getPayOrdersByMid, err)
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		po := &model.PayOrder{}
		if err = rows.Scan(&po.ID, &po.OrderNo, &po.Platform, &po.OrderType, &po.Mid, &po.BuyMonths, &po.ProductId, &po.Money, &po.Quantity, &po.RefundAmount, &po.Status, &po.ThirdTradeNo, &po.PaymentMoney, &po.PaymentType, &po.PaymentTime, &po.Ver, &po.Token, &po.Ctime, &po.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			err = errors.WithStack(err)
			return
		}
		res = append(res, po)
	}
	return
}

// PayOrdersByMidAndStatus quires rows from tv_pay_order.
func (d *Dao) PayOrdersByMidAndStatus(c context.Context, mid int, status int8, pn, ps int) (res []*model.PayOrder, total int, err error) {
	res = make([]*model.PayOrder, 0)
	totalRow := d.db.QueryRow(c, _countPayOrderByMidAndStatus, mid, status)
	if err = totalRow.Scan(&total); err != nil {
		log.Error("row.ScanCount error(%v)", err)
		err = errors.WithStack(err)
		return
	}
	rows, err := d.db.Query(c, _getPayOrdersByMidAndStatus, mid, status, (pn-1)*ps, ps)
	if err != nil {
		log.Error("db.Query(%s) error(%v)", _getPayOrdersByMidAndStatus, err)
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		po := &model.PayOrder{}
		if err = rows.Scan(&po.ID, &po.OrderNo, &po.Platform, &po.OrderType, &po.Mid, &po.BuyMonths, &po.ProductId, &po.Money, &po.Quantity, &po.RefundAmount, &po.Status, &po.ThirdTradeNo, &po.PaymentMoney, &po.PaymentType, &po.PaymentTime, &po.Ver, &po.Token, &po.Ctime, &po.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			err = errors.WithStack(err)
			return
		}
		res = append(res, po)
	}
	return
}

// PayOrdersByMidAndStatusAndCtime quires rows from tv_pay_order.
func (d *Dao) PayOrdersByMidAndStatusAndCtime(c context.Context, mid int64, status int8, from, to xtime.Time, pn, ps int) (res []*model.PayOrder, total int, err error) {
	res = make([]*model.PayOrder, 0)
	totalRow := d.db.QueryRow(c, _countPayOrderByMidAndStatusAndCtime, mid, status, from, to)
	if err = totalRow.Scan(&total); err != nil {
		log.Error("row.ScanCount error(%v)", err)
		err = errors.WithStack(err)
		return
	}
	rows, err := d.db.Query(c, _getPayOrdersByMidAndStatusAndCtime, mid, status, from, to, (pn-1)*ps, ps)
	if err != nil {
		log.Error("db.Query(%s) error(%v)", _getPayOrdersByMidAndStatusAndCtime, err)
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		po := &model.PayOrder{}
		if err = rows.Scan(&po.ID, &po.OrderNo, &po.Platform, &po.OrderType, &po.Mid, &po.BuyMonths, &po.ProductId, &po.Money, &po.Quantity, &po.RefundAmount, &po.Status, &po.ThirdTradeNo, &po.PaymentMoney, &po.PaymentType, &po.PaymentTime, &po.Ver, &po.Ctime, &po.Mtime, &po.Token); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			err = errors.WithStack(err)
			return
		}
		res = append(res, po)
	}
	return
}

// TxInsertPayOrder insert one row into tv_pay_order.
func (d *Dao) TxInsertPayOrder(ctx context.Context, tx *xsql.Tx, po *model.PayOrder) (id int64, err error) {
	var (
		res sql.Result
	)
	if res, err = tx.Exec(_insertPayOrder, po.OrderNo, po.Platform, po.OrderType, po.Mid, po.BuyMonths, po.ProductId, po.Money, po.Quantity, po.RefundAmount, po.Status, po.ThirdTradeNo, po.PaymentMoney, po.PaymentType, po.PaymentTime, po.Ver, po.Token, po.AppChannel); err != nil {
		log.Error("tx.Exec(%s) error(%v)", _insertPayOrder, err)
		err = errors.WithStack(err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// TxUpdatePayOrder updates status, third party no and payment time.
func (d *Dao) TxUpdatePayOrder(ctx context.Context, tx *xsql.Tx, po *model.PayOrder) error {
	if _, err := tx.Exec(_updatePayOrder, po.Status, xtime.Time(time.Now().Unix()), po.ID, po.Ver); err != nil {
		log.Error("tx.Exec(%s) error(%v)", _updatePayOrder, err)
		err = errors.WithStack(err)
		return err
	}
	return nil
}

//UnpaidNotCallbackOrder get orders not paid where stime < ctime < etime
func (d *Dao) UnpaidNotCallbackOrder(c context.Context, stime, etime xtime.Time, pn, ps int) (res []*model.PayOrder, err error) {
	rows, err := d.db.Query(c, _getUnpaidNoCallbackOrder, stime, etime, (pn-1)*ps, ps)
	if err != nil {
		log.Error("db.Query(%s) error(%v)", _getUnpaidNoCallbackOrder, err)
		err = errors.WithStack(err)
		return
	}
	for rows.Next() {
		po := &model.PayOrder{}
		if err = rows.Scan(&po.ID, &po.OrderNo, &po.Platform, &po.OrderType, &po.Mid, &po.BuyMonths, &po.ProductId, &po.Money, &po.Quantity, &po.RefundAmount, &po.Status, &po.ThirdTradeNo, &po.PaymentMoney, &po.PaymentType, &po.PaymentTime, &po.Ver, &po.Token, &po.Ctime, &po.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			err = errors.WithStack(err)
			return
		}
		res = append(res, po)
	}
	return
}
