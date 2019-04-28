package dao

import (
	"database/sql"

	"go-common/app/admin/ep/marthe/model"
	"go-common/library/ecode"

	pkgerr "github.com/pkg/errors"
)

// InsertBuglyProject Insert Bugly Project.
func (d *Dao) InsertBuglyProject(buglyProject *model.BuglyProject) error {
	return pkgerr.WithStack(d.db.Create(buglyProject).Error)
}

// UpdateBuglyProject Update Bugly Project.
func (d *Dao) UpdateBuglyProject(buglyProject *model.BuglyProject) error {
	return pkgerr.WithStack(d.db.Model(&model.BuglyProject{}).Updates(buglyProject).Error)
}

// QueryBuglyProject Query Bugly Project.
func (d *Dao) QueryBuglyProject(id int64) (buglyProject *model.BuglyProject, err error) {
	buglyProject = &model.BuglyProject{}
	err = pkgerr.WithStack(d.db.Where("id = ?", id).First(buglyProject).Error)
	return
}

// QueryBuglyProjectByName Query Bugly Project.
func (d *Dao) QueryBuglyProjectByName(projectName string) (buglyProject *model.BuglyProject, err error) {
	buglyProject = &model.BuglyProject{}
	if err = d.db.Where("project_name = ?", projectName).First(buglyProject).Error; err == ecode.NothingFound {
		err = nil
	}
	return
}

// QueryAllBuglyProjects Query All Bugly Project.
func (d *Dao) QueryAllBuglyProjects() (buglyProjects []*model.BuglyProject, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.BuglyProject{}).Find(&buglyProjects).Error)
	return
}

// FindBuglyProjects Find Bugly Project.
func (d *Dao) FindBuglyProjects(req *model.QueryBuglyProjectRequest) (total int64, buglyProject []*model.BuglyProject, err error) {
	gDB := d.db.Model(&model.BuglyProject{})

	if req.ProjectID != "" {
		gDB = gDB.Where("project_id=?", req.ProjectID)
	}
	if req.ProjectName != "" {
		gDB = gDB.Where("project_name like ?", _wildcards+req.ProjectName+_wildcards)
	}
	if req.PlatformID != "" {
		gDB = gDB.Where("platform_id=?", req.PlatformID)
	}

	if req.UpdateBy != "" {
		gDB = gDB.Where("update_by=?", req.UpdateBy)
	}

	if err = pkgerr.WithStack(gDB.Count(&total).Error); err != nil {
		return
	}

	err = pkgerr.WithStack(gDB.Order("ctime desc").Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Find(&buglyProject).Error)
	return
}

// QueryBuglyProjectList Query Bugly Project List.
func (d *Dao) QueryBuglyProjectList() (projectList []string, err error) {
	var (
		rows *sql.Rows
	)
	sql := "select DISTINCT project_name from bugly_projects"
	if rows, err = d.db.Raw(sql).Rows(); err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var ver string
		if err = rows.Scan(&ver); err != nil {
			return
		}
		projectList = append(projectList, ver)
	}
	return
}
