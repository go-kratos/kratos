package service

import (
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/log"
)

// AddReport add report
func (s *Service) AddReport(userName string, report *model.OrderReport) (err error) {
	var qor *model.QueryOrderResponse
	report.UpdateBy = userName
	report.Active = 1

	// 更新order report
	if report.ID != 0 {
		return s.UpdateReportByID(report)
	}

	// 新增order report
	order := model.Order{ID: report.OrderID}
	if qor, err = s.dao.QueryOrder(&order, 1, 1); err != nil {
		log.Error("order_report.service get order error (%v)", err)
		return err
	}
	report.Name = qor.Orders[0].Name
	return s.dao.AddReport(report)
}

//QueryReportByOrderID  query report
func (s *Service) QueryReportByOrderID(orderID int64) (*model.OrderReport, error) {
	return s.dao.QueryReportByOrderID(orderID)
}

//UpdateReportByID update report by order_id
func (s *Service) UpdateReportByID(report *model.OrderReport) (err error) {
	err = s.dao.UpdateReportByID(report)
	return
}
