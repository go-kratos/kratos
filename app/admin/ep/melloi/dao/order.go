package dao

import (
	"go-common/app/admin/ep/melloi/model"
)

// QueryOrder query order info: ps : pageNumber pn: pageSize
func (d *Dao) QueryOrder(order *model.Order, pn int32, ps int32) (qor *model.QueryOrderResponse, err error) {
	qor = &model.QueryOrderResponse{}
	err = d.DB.Table(model.Order{}.TableName()).Where(model.Order{
		ID: order.ID, Name: order.Name, Broker: order.Broker, Type: order.Type, TestType: order.TestType,
		Project: order.Project, Department: order.Department, App: order.App, Status: order.Status, UpdateBy: order.UpdateBy,
		Active: 1, Handler: order.Handler}).Count(&qor.TotalSize).Offset((pn - 1) * ps).Limit(ps).Order("id desc").Find(&qor.Orders).Error
	qor.PageSize = ps
	qor.PageNum = pn
	return
}

// AddOrder add order by order object
func (d *Dao) AddOrder(order *model.Order) error {
	return d.DB.Table(model.Order{}.TableName()).Create(order).Error
}

// UpdateOrder update order by order object
func (d *Dao) UpdateOrder(order *model.Order) error {
	return d.DB.Model(&model.Order{}).Update(order).Where("ID=?", order.ID).Error
}

// DelOrder delete order by orderID
func (d *Dao) DelOrder(id int64) error {
	return d.DB.Model(&model.Order{}).Where("ID=?", id).Update("active", -1).Error
}
