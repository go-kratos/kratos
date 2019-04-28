package service

import (
	"time"

	"go-common/app/admin/ep/melloi/model"
)

// QueryOrder query order by order object
func (s *Service) QueryOrder(qor *model.QueryOrderRequest) (*model.QueryOrderResponse, error) {
	return s.dao.QueryOrder(&qor.Order, qor.PageNum, qor.PageSize)
}

// UpdateOrder update perf order information
func (s *Service) UpdateOrder(order *model.Order) error {
	return s.dao.UpdateOrder(order)
}

// AddOrder create new order
func (s *Service) AddOrder(order *model.Order) error {
	order.ApplyDate = time.Now()
	order.Active = 1
	return s.dao.AddOrder(order)
}

// DelOrder delete order info by orderID
func (s *Service) DelOrder(id int64) error {
	return s.dao.DelOrder(id)
}
