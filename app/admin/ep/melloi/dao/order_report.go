package dao

import (
	"go-common/app/admin/ep/melloi/model"
)

// AddReport add manual report
func (d *Dao) AddReport(report *model.OrderReport) (err error) {
	return d.DB.Model(&model.OrderReport{}).Create(report).Error
}

// QueryReportByOrderID query report  by order id
func (d *Dao) QueryReportByOrderID(orderID int64) (report *model.OrderReport, err error) {
	report = &model.OrderReport{}
	err = d.DB.Table(model.OrderReport{}.TableName()).Where("order_id=? ", orderID).Find(&report).Error
	return
}

// UpdateReportByID update report by id
func (d *Dao) UpdateReportByID(report *model.OrderReport) error {
	return d.DB.Table(model.OrderReport{}.TableName()).Where("ID=?", report.ID).Update("content", report.Content).Error
}
