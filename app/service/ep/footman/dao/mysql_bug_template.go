package dao

import (
	"go-common/app/service/ep/footman/model"
)

// FindBugTemplates Find Bug Templates.
func (d *Dao) FindBugTemplates(projectID string) (bugTemplate *model.BugTemplate, err error) {
	bugTemplate = &model.BugTemplate{}
	err = d.db.Where("project_id = ?", projectID).First(bugTemplate).Error
	return
}
