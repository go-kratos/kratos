package dao

import (
	"go-common/app/admin/ep/saga/conf"
	"go-common/app/admin/ep/saga/model"

	"github.com/jinzhu/gorm"
	pkgerr "github.com/pkg/errors"
)

// ProjectExist check the project exist or not
func (d *Dao) ProjectExist(projID int) (exist bool, err error) {
	var count int
	if err = pkgerr.WithStack(d.db.Model(&model.ProjectInfo{}).Where(&model.ProjectInfo{ProjectID: projID}).Count(&count).Error); err != nil {
		return
	}
	if count > 0 {
		exist = true
	}
	return
}

// FavoriteProjects get user's favorite projects
func (d *Dao) FavoriteProjects(userName string) (favorites []*model.ProjectFavorite, err error) {
	//err = pkgerr.WithStack(d.db.Where(&model.ProjectFavorite{UserName: userName}).Find(&favorites).Error)
	err = pkgerr.WithStack(d.db.Where("user_name = ?", userName).Find(&favorites).Error)
	return
}

// AddFavorite add favorite project for user
func (d *Dao) AddFavorite(userName string, projID int) (err error) {
	return pkgerr.WithStack(d.db.Create(&model.ProjectFavorite{UserName: userName, ProjID: projID}).Error)
}

// DelFavorite delete favorite project for user
func (d *Dao) DelFavorite(userName string, projID int) (err error) {
	return pkgerr.WithStack(d.db.Delete(model.ProjectFavorite{UserName: userName, ProjID: projID}).Error)
}

// AddProjectInfo add ProjectInfo
func (d *Dao) AddProjectInfo(projectInfo *model.ProjectInfo) (err error) {
	return pkgerr.WithStack(d.db.Create(projectInfo).Error)
}

// ProjectsInfo all the projects in saga
func (d *Dao) ProjectsInfo() (projects []*model.ProjectInfo, err error) {
	err = pkgerr.WithStack(d.db.Find(&projects).Error)
	return
}

// ProjectInfoByID query ProjectInfo by ID
func (d *Dao) ProjectInfoByID(projectID int) (projectInfo *model.ProjectInfo, err error) {
	projectInfo = &model.ProjectInfo{}
	err = pkgerr.WithStack(d.db.Where("project_id = ?", projectID).First(projectInfo).Error)
	return
}

// UpdateProjectInfo update
func (d *Dao) UpdateProjectInfo(projectID int, projectInfo *model.ProjectInfo) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.ProjectInfo{}).Where("project_id = ?", projectID).Update(projectInfo).Error)
}

// HasProjectInfo is exist project info in database.
func (d *Dao) HasProjectInfo(projectID int) (b bool, err error) {
	var size int64
	if err = pkgerr.WithStack(d.db.Model(&model.ProjectInfo{}).Where("project_id = ?", projectID).Count(&size).Error); err != nil {
		return
	}
	b = size > 0
	return
}

// QueryProjectInfo query Project Info.
func (d *Dao) QueryProjectInfo(ifPage bool, req *model.ProjectInfoRequest) (total int, projectInfo []*model.ProjectInfo, err error) {
	var (
		projectIDs = conf.Conf.Property.DefaultProject.ProjectIDs
		db         *gorm.DB
	)

	gDB := d.db.Model(&model.ProjectInfo{})
	if req.Name != "" {
		gDB = gDB.Where("name = ?", req.Name)
	}
	if req.Department != "" {
		gDB = gDB.Where("department = ?", req.Department)
	}
	if req.Business != "" {
		gDB = gDB.Where("business = ?", req.Business)
	}

	if err = pkgerr.WithStack(gDB.Count(&total).Error); err != nil {
		return
	}
	if ifPage {
		if req.PageNum == 1 {
			db = gDB.Order("name DESC").Where("project_id in (?)", projectIDs).Find(&projectInfo)
		} else {
			db = gDB.Order("name DESC").Offset((req.PageNum - 2) * req.PageSize).Limit(req.PageSize).Find(&projectInfo)
		}
	} else {
		db = gDB.Find(&projectInfo)
	}
	if db.Error != nil {
		if db.RecordNotFound() {
			err = nil
		} else {
			err = pkgerr.WithStack(db.Error)
		}
	}
	return
}

// QueryConfigInfo query saga and runner Info.
func (d *Dao) QueryConfigInfo(name, department, business, queryObject string) (total int, err error) {

	gDB := d.db.Model(&model.ProjectInfo{})
	if name != "" {
		gDB = gDB.Where("name = ?", name)
	}
	if department != "" {
		gDB = gDB.Where("department = ?", department)
	}
	if business != "" {
		gDB = gDB.Where("business = ?", business)
	}

	if queryObject == model.ObjectSaga {
		gDB = gDB.Where("Saga = ?", true)
	} else if queryObject == model.ObjectRunner {
		gDB = gDB.Where("Runner = ?", true)
	} else {
		return
	}

	err = pkgerr.WithStack(gDB.Count(&total).Error)

	return
}
