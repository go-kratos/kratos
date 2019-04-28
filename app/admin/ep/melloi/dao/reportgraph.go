package dao

import (
	"go-common/app/admin/ep/melloi/model"
)

//AddReportGraph add reportGraph
func (d *Dao) AddReportGraph(reportGraph *model.ReportGraph) error {
	return d.DB.Create(reportGraph).Error
}

//QueryReportGraph query reportGraph
func (d *Dao) QueryReportGraph(testNameNicks []string) (reportGraphs []model.ReportGraph, err error) {
	err = d.DB.Model(&model.ReportGraph{}).Where(" test_name_nick in (?) ", testNameNicks).Order("elapsd_time asc").Find(&reportGraphs).Error
	return
}
