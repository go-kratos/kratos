package dao

import (
	"context"
	"database/sql"
	"encoding/json"

	"go-common/app/common/openplatform/encoding"
	"go-common/app/service/openplatform/ticket-sales/api/grpc/type"
	rpc "go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	"go-common/library/log"
)

const (
	// 更新购买人信息
	_updateDetailBuyer = "update order_detail set buyer = ?,tel = ?,personal_id = ? where order_id = ?"
	// 更新配送信息
	_updateDetailDelivery = "update order_detail set deliver_detail = ? where order_id = ?"
)

//UpdateDetailBuyer 更新购买人信息
func (d *Dao) UpdateDetailBuyer(c context.Context, req *rpc.UpBuyerRequest) (num int64, err error) {

	var res sql.Result
	PersonIDEnc, _ := encoding.Encrypt(req.Buyers.PersonalID, d.c.Encrypt)
	TelEnc, _ := encoding.Encrypt(req.Buyers.Tel, d.c.Encrypt)
	if res, err = d.db.Exec(c, _updateDetailBuyer, req.Buyers.Name, TelEnc, string(PersonIDEnc), req.OrderID); err != nil {
		log.Warn("更新订单详情%d失败", req.OrderID)
		return
	}
	return res.RowsAffected()
}

//UpdateDetailDelivery 更新配送信息
func (d *Dao) UpdateDetailDelivery(c context.Context, arg *_type.OrderDeliver, orderID int64) (num int64, err error) {
	var res sql.Result
	arg.Tel, _ = encoding.Encrypt(arg.Tel, d.c.Encrypt)
	DeliverDetail, err := json.Marshal(arg)

	if err != nil {
		log.Warn("更新订单详情%d失败", orderID)
	}

	if res, err = d.db.Exec(c, _updateDetailDelivery, string(DeliverDetail), orderID); err != nil {
		log.Warn("更新订单详情%d失败", orderID)
		return
	}
	return res.RowsAffected()
}
