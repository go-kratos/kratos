package service

import (
	"context"
	"database/sql"

	"go-common/app/admin/main/config/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// CreateBuild create build.
func (s *Service) CreateBuild(build *model.Build, treeID int64, env, zone string) (err error) {
	var app *model.App
	if app, err = s.AppByTree(treeID, env, zone); err != nil {
		return
	}
	build.AppID = app.ID
	return s.dao.DB.Create(build).Error
}

//UpdateTag update tag.
func (s *Service) UpdateTag(ctx context.Context, treeID int64, env, zone, name string, tag *model.Tag) (err error) {
	var (
		app   *model.App
		build *model.Build
	)
	if app, err = s.AppByTree(treeID, env, zone); err != nil {
		return
	}
	if build, err = s.BuildByName(app.ID, name); err != nil {
		if err != ecode.NothingFound {
			return
		}
		build = &model.Build{AppID: app.ID, Name: name, Mark: tag.Mark, Operator: tag.Operator}
		if err = s.dao.DB.Create(build).Error; err != nil {
			log.Error("CreateBuild(%s) error(%v)", build.Name, err)
			return
		}
	}
	tag.AppID = app.ID
	tag.BuildID = build.ID
	if err = s.dao.DB.Create(&tag).Error; err != nil {
		log.Error("CreateTag(%s) error(%v)", tag.Mark, err)
		return
	}
	if tag.Force == 1 {
		//Clear stand-alone force
		forces := []*model.Force{}
		if err = s.dao.DB.Where("app_id = ?", app.ID).Find(&forces).Error; err != nil {
			log.Error("select forces(%s) error(%v)", app.ID, err)
			return
		}
		mHosts := model.MapHosts{}
		for _, val := range forces {
			mHosts[val.Hostname] = val.IP
		}
		if len(mHosts) > 0 {
			if err = s.ClearForce(ctx, treeID, env, zone, name, mHosts); err != nil {
				log.Error("clear forces(%s) error(%v)", app.ID, err)
				return
			}
		}
	}

	tx := s.dao.DB.Begin()
	if err = tx.Model(&model.Build{ID: build.ID}).Update(map[string]interface{}{
		"tag_id":   tag.ID,
		"operator": tag.Operator,
	}).Error; err != nil {
		tx.Rollback()
		log.Error("updateTagID(%d) error(%v)", tag.ID, err)
		return
	}
	//push
	if err = s.Push(ctx, treeID, env, zone, build.Name, tag.ID); err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

//UpdateTagID update tag.
func (s *Service) UpdateTagID(ctx context.Context, env, zone, bName string, tag, TreeID int64) (err error) {
	build := new(model.Build)
	build.Name = bName
	build.TagID = tag
	var app *model.App
	if app, err = s.AppByTree(TreeID, env, zone); err != nil {
		return
	}
	tx := s.dao.DB.Begin()
	if err = tx.Model(&model.Build{}).Where("app_id = ? and name = ?", app.ID, build.Name).Update(map[string]interface{}{
		"tag_id":   build.TagID,
		"operator": build.Operator,
	}).Error; err != nil {
		tx.Rollback()
		log.Error("updateTagID(%d) error(%v)", build.TagID, err)
		return
	}
	if err = s.Push(ctx, TreeID, env, zone, build.Name, build.TagID); err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

//Builds get builds by app id.
func (s *Service) Builds(treeID int64, appName, env, zone string) (builds []*model.Build, err error) {
	var (
		app *model.App
		tag *model.Tag
	)
	if app, err = s.AppByTree(treeID, env, zone); err != nil {
		if err == ecode.NothingFound {
			if err = s.CreateApp(appName, env, zone, treeID); err == nil {
				builds = make([]*model.Build, 0)
			}
			return
		}
	}
	if builds, err = s.BuildsByApp(app.ID); err != nil {
		return
	}
	for _, build := range builds {
		if tag, err = s.Tag(build.TagID); err != nil {
			if err == ecode.NothingFound {
				err = nil
			}
		}
		build.Mark = tag.Mark
		build.Operator = tag.Operator
		build.Mtime = tag.Mtime
	}
	return
}

//BuildsByApp buildsByApp.
func (s *Service) BuildsByApp(appID int64) (builds []*model.Build, err error) {
	if err = s.dao.DB.Find(&builds, "app_id = ? ", appID).Error; err != nil {
		log.Error("BuildsByApp(%s) error(%v)", appID, err)
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

//Build get Build by build ID.
func (s *Service) Build(ID int64) (build *model.Build, err error) {
	build = new(model.Build)
	if err = s.dao.DB.First(&build, ID).Error; err != nil {
		log.Error("Build(%v) error(%v)", ID, err)
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
		}
	}
	return
}

//Delete delete Build by build ID.
func (s *Service) Delete(ID int64) (err error) {
	if err = s.dao.DB.Delete(&model.Build{}, ID).Error; err != nil {
		log.Error("Delete(%v) error(%v)", ID, err)
	}
	return
}

//BuildByName get Build by build ID.
func (s *Service) BuildByName(appID int64, name string) (build *model.Build, err error) {
	build = new(model.Build)
	if err = s.dao.DB.First(&build, "app_id = ? and name = ?", appID, name).Error; err != nil {
		log.Error("BuildByName(%s) error(%v)", name, err)
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
		}
	}
	return
}

//GetDelInfos get delete info.
func (s *Service) GetDelInfos(c context.Context, BuildID int64) (err error) {
	build := &model.Build{}
	if err = s.dao.DB.Where("id = ?", BuildID).First(build).Error; err != nil {
		log.Error("GetDelInfos BuildID(%v) error(%v)", BuildID, err)
		return
	}
	app := &model.App{}
	if err = s.dao.DB.Where("id = ?", build.AppID).First(app).Error; err != nil {
		log.Error("GetDelInfos AppID(%v) error(%v)", build.AppID, err)
		return
	}
	hosts, err := s.Hosts(c, app.TreeID, app.Name, app.Env, app.Zone)
	if err != nil {
		log.Error("GetDelInfos hosts(%v) error(%v)", hosts, err)
		return
	}
	for _, v := range hosts {
		if v.BuildVersion == build.Name {
			err = ecode.NothingFound
			return
		}
	}
	return
}

// AllBuilds ...
func (s *Service) AllBuilds(appIDS []int64) (builds []*model.Build, err error) {
	if err = s.dao.DB.Where("app_id in (?)", appIDS).Find(&builds).Error; err != nil {
		log.Error("AllBuild error(%v)", err)
	}
	return
}
