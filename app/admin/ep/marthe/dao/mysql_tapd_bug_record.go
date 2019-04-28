package dao

import (
	"go-common/app/admin/ep/marthe/model"

	pkgerr "github.com/pkg/errors"
)

// InsertTapdBugRecord Insert Tapd Bug Insert Log.
func (d *Dao) InsertTapdBugRecord(tapdBugRecord *model.TapdBugRecord) error {
	return pkgerr.WithStack(d.db.Create(tapdBugRecord).Error)
}

// UpdateTapdBugRecord Update Tapd Bug Insert Log.
func (d *Dao) UpdateTapdBugRecord(tapdBugRecord *model.TapdBugRecord) error {
	return pkgerr.WithStack(d.db.Save(&tapdBugRecord).Error)
}

// QueryTapdBugRecordByProjectIDAndStatus Query Tapd Bug Record By Project ID and status
func (d *Dao) QueryTapdBugRecordByProjectIDAndStatus(projectID int64, status int) (tapdBugRecords []*model.TapdBugRecord, err error) {
	err = pkgerr.WithStack(d.db.Where("project_template_id = ? and status = ?", projectID, status).Find(&tapdBugRecords).Error)
	return
}

// QueryTapdBugRecordByStatus Query Tapd Bug Record By and status
func (d *Dao) QueryTapdBugRecordByStatus(status int) (tapdBugRecords []*model.TapdBugRecord, err error) {
	err = pkgerr.WithStack(d.db.Where("status = ?", status).Find(&tapdBugRecords).Error)
	return
}

// FindBugRecords Find Bug Records.
func (d *Dao) FindBugRecords(req *model.QueryBugRecordsRequest) (total int64, tapdBugRecords []*model.TapdBugRecord, err error) {
	gDB := d.db.Model(&model.TapdBugRecord{})

	if req.ProjectTemplateID > 0 {
		gDB = gDB.Where("project_template_id=?", req.ProjectTemplateID)
	}

	if req.VersionTemplateID > 0 {
		gDB = gDB.Where("version_template_id=?", req.VersionTemplateID)
	}

	if req.Operator != "" {
		gDB = gDB.Where("operator=?", req.Operator)
	}

	if req.Status > 0 {
		gDB = gDB.Where("status=?", req.Status)
	}

	if err = pkgerr.WithStack(gDB.Count(&total).Error); err != nil {
		return
	}

	err = pkgerr.WithStack(gDB.Order("ctime desc").Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Find(&tapdBugRecords).Error)
	return
}
