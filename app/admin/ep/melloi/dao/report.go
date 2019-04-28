package dao

import (
	"go-common/app/admin/ep/melloi/model"
)

// QueryReportSummarys query reportSummarys
func (d *Dao) QueryReportSummarys(reportSummary *model.ReportSummary, searchAll bool, pn, ps int32, treeNodes []string) (qrsr *model.QueryReportSuResponse, err error) {
	qrsr = &model.QueryReportSuResponse{}
	reportSummary.Active = 1
	reportSummary.Debug = -1

	gDB := d.DB.Table(model.ReportSummary{}.TableName()) //.Where("app in (?)", treeNodes)
	if reportSummary.TestName != "" && reportSummary.ScriptID != 0 {
		gDB = gDB.Where("test_name LIKE ? and active = ? and debug = ? and script_id = ?", "%"+reportSummary.TestName+"%", 1, -1, reportSummary.ScriptID)
	}
	if reportSummary.UserName != "" && reportSummary.ScriptID != 0 {
		gDB = gDB.Where("user_name LIKE ? and active = ? and debug = ? and script_id = ?", "%"+reportSummary.UserName+"%", 1, -1, reportSummary.ScriptID)
	}
	if reportSummary.TestName != "" {
		gDB = gDB.Where("test_name LIKE ? and active = ? and debug = ?", "%"+reportSummary.TestName+"%", 1, -1)
	}
	if reportSummary.UserName != "" {
		gDB = gDB.Where("user_name LIKE ? and active = ? and debug = ?", "%"+reportSummary.UserName+"%", 1, -1)
	}
	if searchAll {
		gDB = gDB.Where("type in (0, 1, 2)")
	} else {
		if reportSummary.ID == 0 {
			gDB = gDB.Where("type = ?", reportSummary.Type)
		}
	}
	if reportSummary.TestName == "" && reportSummary.UserName == "" {
		gDB = gDB.Where(reportSummary)
	}
	err = gDB.Count(&qrsr.TotalSize).
		Order("ctime desc").Offset((pn - 1) * ps).Limit(ps).Find(&qrsr.ReportSummarys).Error

	qrsr.PageSize = ps /**/
	qrsr.PageNum = pn
	return
}

// QueryReportSummarysWhiteName query reportSummarys by whiteName
func (d *Dao) QueryReportSummarysWhiteName(reportSummary *model.ReportSummary, searchAll bool, pn, ps int32) (qrsr *model.QueryReportSuResponse, err error) {
	qrsr = &model.QueryReportSuResponse{}
	reportSummary.Active = 1
	reportSummary.Debug = -1

	gDB := d.DB.Table(model.ReportSummary{}.TableName())
	if reportSummary.TestName != "" && reportSummary.ScriptID != 0 {
		gDB = gDB.Where("test_name LIKE ? and active = ? and debug = ? and script_id = ?", "%"+reportSummary.TestName+"%", 1, -1, reportSummary.ScriptID)
	}
	if reportSummary.UserName != "" && reportSummary.ScriptID != 0 {
		gDB = gDB.Where("user_name LIKE ? and active = ? and debug = ? and script_id = ?", "%"+reportSummary.UserName+"%", 1, -1, reportSummary.ScriptID)
	}
	if reportSummary.TestName != "" {
		gDB = gDB.Where("test_name LIKE ? and active = ? and debug = ?", "%"+reportSummary.TestName+"%", 1, -1)
	}
	if reportSummary.UserName != "" {
		gDB = gDB.Where("user_name LIKE ? and active = ? and debug = ?", "%"+reportSummary.UserName+"%", 1, -1)
	}
	if searchAll {
		gDB = gDB.Where("type in (0, 1, 2)")
	} else {
		if reportSummary.ID == 0 {
			gDB = gDB.Where("type = ?", reportSummary.Type)
		}
	}
	if reportSummary.TestName == "" && reportSummary.UserName == "" {
		gDB = gDB.Where(reportSummary)
	}
	err = gDB.Count(&qrsr.TotalSize).
		Order("ctime desc").Offset((pn - 1) * ps).Limit(ps).Find(&qrsr.ReportSummarys).Error

	qrsr.PageSize = ps /**/
	qrsr.PageNum = pn
	return
}

// QueryReportSurys query reportSummarys
func (d *Dao) QueryReportSurys(reportSummary *model.ReportSummary) (res []*model.ReportSummary, err error) {
	err = d.DB.Table(model.ReportSummary{}.TableName()).Where(reportSummary).Find(&res).Error
	return
}

// QueryReportSuryByID query reportSummary by id
func (d *Dao) QueryReportSuryByID(id int) (res *model.ReportSummary, err error) {
	reportSummary := model.ReportSummary{ID: id}
	res = &model.ReportSummary{}
	err = d.DB.Table(model.ReportSummary{}.TableName()).Where(reportSummary).First(&res).Error
	return
}

//CountQueryReportSummarys count queryReportSummarys
func (d *Dao) CountQueryReportSummarys(reportSummary *model.ReportSummary) (total int, err error) {
	err = d.DB.Table(model.ReportSummary{}.TableName()).Where(reportSummary).Count(&total).Error
	return
}

//AddReportSummary add reportSummary
func (d *Dao) AddReportSummary(reportSummary *model.ReportSummary) (reportSuID int, err error) {
	err = d.DB.Create(reportSummary).Error
	reportSuID = reportSummary.ID
	return
}

//UpdateReportSummary update Report
func (d *Dao) UpdateReportSummary(reportSummary *model.ReportSummary) error {
	return d.DB.Model(&model.ReportSummary{}).Where("id = ?", reportSummary.ID).Updates(reportSummary).Error
}

//UpdateReportStatusByID update report status
func (d *Dao) UpdateReportStatusByID(ID, testStatus int) error {
	return d.DB.Exec("update report_summary set test_status = ? where id = ? ", testStatus, ID).Error
}

//UpdateReportStatus update Report
func (d *Dao) UpdateReportStatus(status int) error {
	return d.DB.Model(&model.ReportSummary{}).Where("test_status = 2").Update("test_status", 1).Error
}

//UpdateReportDockByID update Report
func (d *Dao) UpdateReportDockByID(ID, dockerSum int) error {
	return d.DB.Model(&model.ReportSummary{}).Where("id =?", ID).Updates(model.ReportSummary{DockerSum: dockerSum}).Error
}

//DelReportSummary delete reportSummary
func (d *Dao) DelReportSummary(id int) error {
	return d.DB.Model(&model.ReportSummary{}).Where("ID=?", id).Update("active", 0).Error
}

//QueryReTimely query rueryReTimely
func (d *Dao) QueryReTimely(testName, beginTime, afterTime string, podNames []string) (reportTimelys []*model.ReportTimely, err error) {
	err = d.DB.Model(&model.ReportTimely{}).Where("test_name = ? and mtime >= ? and mtime <= ? and pod_name in (?)", testName, beginTime, afterTime, podNames).Find(&reportTimelys).Error
	return
}
