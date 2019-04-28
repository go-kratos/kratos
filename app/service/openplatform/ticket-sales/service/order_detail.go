package service

import (
	"context"
	"fmt"

	"go-common/app/service/openplatform/ticket-sales/api/grpc/type"
	rpc "go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/library/log"
)

// UpdateBuyer 更新购买信息
func (s *Service) UpdateBuyer(c context.Context, req *rpc.UpBuyerRequest) (res *rpc.UpDetailResponse, err error) {
	var oi int16 = 1
	res = &rpc.UpDetailResponse{}
	res.OrderID = req.OrderID
	_, err = s.dao.UpdateDetailBuyer(c, req)
	if err != nil {
		log.Warn("更新订单详情%d失败", req.OrderID)
		res.IsUpdate = 0
		return
	}

	key := fmt.Sprintf("%s:%d", model.CacheKeyOrderDt, req.OrderID)
	s.dao.RedisDel(c, key)
	res.IsUpdate = oi
	return
}

// UpdateDelivery 更新delivery 信息
func (s *Service) UpdateDelivery(c context.Context, req *rpc.UpDeliveryRequest) (res *rpc.UpDetailResponse, err error) {
	var oi int16 = 1
	res = &rpc.UpDetailResponse{}
	res.OrderID = req.OrderID
	deli := req.DeliverDetail
	delivery := &_type.OrderDeliver{
		AddrID: deli.AddrID,
		Name:   deli.Name,
		Addr:   deli.Addr,
		Tel:    deli.Tel,
	}

	_, err = s.dao.UpdateDetailDelivery(c, delivery, req.OrderID)
	if err != nil {
		log.Warn("更新订单详情配送信息%d失败", req.OrderID)
		res.IsUpdate = 0
		return
	}
	key := fmt.Sprintf("%s:%d", model.CacheKeyOrderDt, req.OrderID)
	s.dao.RedisDel(c, key)
	res.IsUpdate = oi
	return
}
