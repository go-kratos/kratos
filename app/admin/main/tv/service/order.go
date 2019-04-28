package service

import (
	"time"

	"go-common/app/admin/main/tv/model"
	"go-common/library/log"
)

const (
	_orderTable = "tv_pay_order"
)

//OrderList user tv-vip order list
func (s *Service) OrderList(mid, pn, ps, paymentStime, paymentEtime int64, status int8, orderNo string) (data *model.OrderPageHelper, err error) {
	var (
		db     = s.dao.DB.Table(_orderTable)
		orders []*model.TvPayOrderResp
	)
	data = &model.OrderPageHelper{}

	if mid != 0 {
		db = db.Where("mid = ?", mid)
	}
	if orderNo != "" {
		db = db.Where("order_no = ?", orderNo)
	}
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	if paymentStime != 0 {
		db = db.Where("payment_time > ?", time.Unix(paymentStime, 0))
	}
	if paymentEtime != 0 {
		db = db.Where("payment_time < ?", time.Unix(paymentEtime, 0))
	}

	db.Count(&data.Total)
	if err := db.Offset((pn - 1) * ps).Limit(ps).Find(&orders).Error; err != nil {
		log.Error("OrderList %v, Err %v", orders, err)
	}
	data.Items = orders

	return
}
