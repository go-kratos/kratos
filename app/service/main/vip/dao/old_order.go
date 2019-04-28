package dao

import (
	"context"

	"github.com/pkg/errors"

	"go-common/app/service/main/vip/model"
)

const (
	_insertOldPayOrder      = "INSERT INTO vip_pay_order(order_no,app_id,mid,buy_months,money,status,ver,platform,app_sub_id,bmid,order_type,coupon_money,pid,user_ip)VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
	_insertOldRechargeOrder = "INSERT INTO vip_recharge_order(app_id,pay_mid,order_no,recharge_bp,pay_order_no,status,remark,ver,third_trade_no)VALUES(?,?,?,?,?,?,?,?,?);"
)

//AddOldPayOrder add old payorder.
func (d *Dao) AddOldPayOrder(c context.Context, r *model.VipOldPayOrder) (err error) {
	if _, err = d.olddb.Exec(c, _insertOldPayOrder, &r.OrderNo, &r.AppID, &r.Mid, &r.BuyMonths, &r.Money, &r.Status, &r.Ver, &r.Platform, &r.AppSubID, &r.Bmid,
		&r.OrderType, &r.CouponMoney, &r.PID, &r.UserIP); err != nil {
		err = errors.Wrapf(err, "dao add old pay order(%+v)", r)
		return
	}
	return
}

//AddOldRechargeOrder add recharge order.
func (d *Dao) AddOldRechargeOrder(c context.Context, r *model.VipOldRechargeOrder) (err error) {
	if _, err = d.olddb.Exec(c, _insertOldRechargeOrder, &r.AppID, &r.PayMid, &r.OrderNo, &r.RechargeBp, &r.PayOrderNO, &r.Status, &r.Remark, &r.Ver, &r.ThirdTradeNO); err != nil {
		err = errors.Wrapf(err, "dao add old recharge order(%+v)", r)
		return
	}
	return
}
