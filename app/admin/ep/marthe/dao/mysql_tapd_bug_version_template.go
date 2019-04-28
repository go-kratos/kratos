package dao

import (
	"go-common/app/admin/ep/marthe/model"
	"go-common/library/ecode"

	pkgerr "github.com/pkg/errors"
)

// InsertTapdBugVersionTemplate Insert TapdBug Version Template.
func (d *Dao) InsertTapdBugVersionTemplate(tapdBugVersionTemplate *model.TapdBugVersionTemplate) error {
	return pkgerr.WithStack(d.db.Create(tapdBugVersionTemplate).Error)
}

// UpdateTapdBugVersionTemplate Update Tapd Bug Version Template.
func (d *Dao) UpdateTapdBugVersionTemplate(tapdBugVersionTemplate *model.TapdBugVersionTemplate) error {
	return pkgerr.WithStack(d.db.Save(&tapdBugVersionTemplate).Error)
}

// QueryTapdBugVersionTemplate Query Tapd Bug Version Template.
func (d *Dao) QueryTapdBugVersionTemplate(id int64) (tapdBugVersionTemplate *model.TapdBugVersionTemplate, err error) {
	tapdBugVersionTemplate = &model.TapdBugVersionTemplate{}
	if err = d.db.Where("id=?", id).First(&tapdBugVersionTemplate).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			err = pkgerr.WithStack(err)
		}
	}
	return
}

// QueryTapdBugVersionTemplateByVersion Query Tapd Bug Version Template.
func (d *Dao) QueryTapdBugVersionTemplateByVersion(version string) (tapdBugVersionTemplate *model.TapdBugVersionTemplate, err error) {
	tapdBugVersionTemplate = &model.TapdBugVersionTemplate{}
	if err = d.db.Where("version=?", version).First(&tapdBugVersionTemplate).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			err = pkgerr.WithStack(err)
		}
	}
	return
}

// FindTapdBugVersionTemplates Find Tapd Bug Version Templates.
func (d *Dao) FindTapdBugVersionTemplates(req *model.QueryTapdBugVersionTemplateRequest) (total int64, tapdBugVersionTemplate []*model.TapdBugVersionTemplate, err error) {
	gDB := d.db.Model(&model.TapdBugVersionTemplate{})

	if req.ProjectID > 0 {
		gDB = gDB.Where("project_template_id = ?", req.ProjectID)
	}

	if req.Version != "" {
		gDB = gDB.Where("version like ?", req.Version+_wildcards)
	}

	if req.UpdateBy != "" {
		gDB = gDB.Where("update_by = ?", req.UpdateBy)
	}

	if err = pkgerr.WithStack(gDB.Count(&total).Error); err != nil {
		return
	}

	err = pkgerr.WithStack(gDB.Order("ctime desc").Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Find(&tapdBugVersionTemplate).Error)
	return
}
