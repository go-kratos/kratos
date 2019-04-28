package dao

import (
	"go-common/app/admin/ep/marthe/model"
	"go-common/library/ecode"

	pkgerr "github.com/pkg/errors"
)

// InsertTapdBugTemplate Insert Tapd Bug Template.
func (d *Dao) InsertTapdBugTemplate(tapdBugTemplate *model.TapdBugTemplate) error {
	return pkgerr.WithStack(d.db.Create(tapdBugTemplate).Error)
}

// UpdateTapdBugTemplate Update Tapd Bug Template.
func (d *Dao) UpdateTapdBugTemplate(tapdBugTemplate *model.TapdBugTemplate) error {
	return pkgerr.WithStack(d.db.Save(&tapdBugTemplate).Error)
}

// QueryTapdBugTemplate Query Tapd Bug Template.
func (d *Dao) QueryTapdBugTemplate(id int64) (tapdBugTemplate *model.TapdBugTemplate, err error) {
	tapdBugTemplate = &model.TapdBugTemplate{}
	if err = d.db.Where("id=?", id).First(&tapdBugTemplate).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			err = pkgerr.WithStack(err)
		}
	}
	return
}

// QueryTapdBugTemplateByProjectID Query Tapd Bug Template by project id.
func (d *Dao) QueryTapdBugTemplateByProjectID(projectID string) (tapdBugTemplate *model.TapdBugTemplate, err error) {
	tapdBugTemplate = &model.TapdBugTemplate{}
	if err = d.db.Where("project_id=?", projectID).First(&tapdBugTemplate).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			err = pkgerr.WithStack(err)
		}
	}
	return
}

// QueryAllTapdBugTemplates Query All Tapd Bug Templates.
func (d *Dao) QueryAllTapdBugTemplates() (tapdBugTemplates []*model.TapdBugTemplate, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.TapdBugTemplate{}).Find(&tapdBugTemplates).Error)
	return
}

// FindTapdBugTemplates Find tapd Bug Templates.
func (d *Dao) FindTapdBugTemplates(req *model.QueryTapdBugTemplateRequest) (total int64, tapdBugTemplates []*model.TapdBugTemplate, err error) {
	gDB := d.db.Model(&model.TapdBugTemplate{})

	if req.UpdateBy != "" {
		gDB = gDB.Where("update_by=?", req.UpdateBy)
	}

	if err = pkgerr.WithStack(gDB.Count(&total).Error); err != nil {
		return
	}

	err = pkgerr.WithStack(gDB.Order("ctime desc").Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Find(&tapdBugTemplates).Error)
	return
}
