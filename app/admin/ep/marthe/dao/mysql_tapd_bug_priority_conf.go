package dao

import (
	"go-common/app/admin/ep/marthe/model"

	pkgerr "github.com/pkg/errors"
)

// InsertTapdBugPriorityConf Insert Tapd Bug Priority Conf.
func (d *Dao) InsertTapdBugPriorityConf(tapdBugPriorityConf *model.TapdBugPriorityConf) error {
	return pkgerr.WithStack(d.db.Create(tapdBugPriorityConf).Error)
}

// UpdateTapdBugPriorityConf Update Tapd Bug Priority Conf.
func (d *Dao) UpdateTapdBugPriorityConf(tapdBugPriorityConf *model.TapdBugPriorityConf) error {
	return pkgerr.WithStack(d.db.Save(&tapdBugPriorityConf).Error)
}

// QueryTapdBugPriorityConfsByProjectTemplateIdAndStatus Query Tapd Bug Priority Confs By Project TemplateId And tatus.
func (d *Dao) QueryTapdBugPriorityConfsByProjectTemplateIdAndStatus(projectTemplateID int64, status int) (tapdBugPriorityConfs []*model.TapdBugPriorityConf, err error) {
	err = pkgerr.WithStack(d.db.Where("project_template_id = ? and status = ?", projectTemplateID, status).Find(&tapdBugPriorityConfs).Error)
	return
}

// FindTapdBugPriorityConfs Find Tapd Bug Priority Confs.
func (d *Dao) FindTapdBugPriorityConfs(req *model.QueryTapdBugPriorityConfsRequest) (total int64, tapdBugPriorityConfs []*model.TapdBugPriorityConf, err error) {
	gDB := d.db.Model(&model.TapdBugPriorityConf{})

	if req.ProjectTemplateID > 0 {
		gDB = gDB.Where("project_template_id=?", req.ProjectTemplateID)
	}

	if req.UpdateBy != "" {
		gDB = gDB.Where("update_by=?", req.UpdateBy)
	}

	if req.Status > 0 {
		gDB = gDB.Where("status=?", req.Status)
	}

	if err = pkgerr.WithStack(gDB.Count(&total).Error); err != nil {
		return
	}

	err = pkgerr.WithStack(gDB.Order("ctime desc").Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Find(&tapdBugPriorityConfs).Error)
	return
}
