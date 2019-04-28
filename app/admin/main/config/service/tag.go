package service

import (
	"database/sql"
	"strconv"
	"strings"

	"go-common/app/admin/main/config/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// CreateTag create App.
func (s *Service) CreateTag(tag *model.Tag, treeID int64, env, zone string) (err error) {
	var app *model.App
	if app, err = s.AppByTree(treeID, env, zone); err != nil {
		return
	}
	tag.AppID = app.ID
	return s.dao.DB.Create(tag).Error
}

//LastTags get tags by app name.
func (s *Service) LastTags(treeID int64, env, zone, bName string) (tags []*model.Tag, err error) {
	var (
		app   *model.App
		build *model.Build
	)
	if app, err = s.AppByTree(treeID, env, zone); err != nil {
		return
	}
	if build, err = s.BuildByName(app.ID, bName); err != nil {
		return
	}
	if err = s.dao.DB.Where("app_id = ? and build_id = ?", app.ID, build.ID).Order("id desc").Limit(10).Find(&tags).Error; err != nil {
		log.Error("Tags(%v) error(%v)", app.ID, err)
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

//TagsByBuild get tags by app name.
func (s *Service) TagsByBuild(appName, env, zone, name string, ps, pn, treeID int64) (tagPager *model.TagConfigPager, err error) {
	var (
		app     *model.App
		build   *model.Build
		tags    []*model.Tag
		confIDs []int64
		confs   []*model.Config
		total   int64
	)
	tagConfigs := make([]*model.TagConfig, 0)
	if app, err = s.AppByTree(treeID, env, zone); err != nil {
		if err == ecode.NothingFound {
			err = s.CreateApp(appName, env, zone, treeID)
			return
		}
	}
	if build, err = s.BuildByName(app.ID, name); err != nil {
		if err == ecode.NothingFound {
			err = nil
		}
		return
	}
	if err = s.dao.DB.Where("app_id = ? and build_id =?", app.ID, build.ID).Order("id desc").Offset((pn - 1) * ps).Limit(ps).Find(&tags).Error; err != nil {
		log.Error("TagsByBuild() findTags() error(%v)", err)
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	tmp := make(map[int64]struct{})
	for _, tag := range tags {
		tmpIDs := strings.Split(tag.ConfigIDs, ",")
		for _, tmpID := range tmpIDs {
			var id int64
			if id, err = strconv.ParseInt(tmpID, 10, 64); err != nil {
				log.Error("strconv.ParseInt() error(%v)", err)
				return
			}
			if _, ok := tmp[id]; !ok {
				tmp[id] = struct{}{}
				confIDs = append(confIDs, id)
			}
		}
	}
	if confs, err = s.ConfigsByIDs(confIDs); err != nil {
		return
	}
	for _, tag := range tags {
		tagConfig := new(model.TagConfig)
		tagConfig.Tag = tag
		tagConfigs = append(tagConfigs, tagConfig)
		tmpIDs := strings.Split(tag.ConfigIDs, ",")
		for _, tmpID := range tmpIDs {
			var id int64
			if id, err = strconv.ParseInt(tmpID, 10, 64); err != nil {
				log.Error("strconv.ParseInt() error(%v)", err)
				return
			}
			for _, conf := range confs {
				if id != conf.ID { //judge config is in build.
					continue
				}
				tagConfig.Confs = append(tagConfig.Confs, conf)
			}
		}
	}
	if err = s.dao.DB.Where("app_id = ? and build_id =?", app.ID, build.ID).Model(&model.Tag{}).Count(&total).Error; err != nil {
		log.Error("TagsByBuild() count() error(%v)", err)
		return
	}
	tagPager = &model.TagConfigPager{Ps: ps, Pn: pn, Total: total, Items: tagConfigs}
	return
}

//LastTagByAppBuild get tags by app and build.
func (s *Service) LastTagByAppBuild(appID, buildID int64) (tag *model.Tag, err error) {
	var (
		tags []*model.Tag
	)
	if err = s.dao.DB.Where("app_id = ? and build_id =?", appID, buildID).Order("id desc").Limit(2).Find(&tags).Error; err != nil {
		log.Error("LastTagByAppBuild() error(%v)", err)
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	if len(tags) != 2 {
		err = ecode.NothingFound
		return
	}
	return tags[1], nil
}

//Tag get tag by id.
func (s *Service) Tag(ID int64) (tag *model.Tag, err error) {
	tag = new(model.Tag)
	if err = s.dao.DB.First(&tag, ID).Error; err != nil {
		log.Error("Tag(%v) error(%v)", ID, err)
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
		}
	}
	return
}

//TagByAppBuild ...
func (s *Service) TagByAppBuild(appID, buildID int64) (tag *model.Tag, err error) {
	tag = new(model.Tag)
	if err = s.dao.DB.Where("app_id = ? and build_id =?", appID, buildID).Order("id desc").First(&tag).Error; err != nil {
		log.Error("TagByAppBuild() error(%v)", err)
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	return
}

// TagByAppBuildLastConfig ...
func (s *Service) TagByAppBuildLastConfig(appID, buildID, tagID int64, cids []int64) (configID int64, err error) {
	tags := []*model.Tag{}
	if err = s.dao.DB.Where("app_id = ? and build_id = ? and id != ?", appID, buildID, tagID).Order("id desc").Find(&tags).Error; err != nil {
		log.Error("TagByAppBuildLastConfig() error(%v)", err)
		return
	}
	var id int64
	for _, v := range tags {
		tmpIDs := strings.Split(v.ConfigIDs, ",")
		for _, tmpID := range tmpIDs {
			id, err = strconv.ParseInt(tmpID, 10, 64)
			if err != nil {
				log.Error("strconv.ParseInt(%s) error(%v)", tmpID, err)
				continue
			}
			for _, vv := range cids {
				if vv == id {
					configID = id
					return
				}
			}
		}
	}
	return
}

//RollBackTag ...
func (s *Service) RollBackTag(tagID int64) (tag *model.Tag, err error) {
	tag = &model.Tag{}
	row := s.dao.DB.Select("`app_id`,`build_id`,`config_ids`,`force`,`mark`,`operator`").Where("id=?", tagID).Model(&model.Tag{}).Row()
	if err = row.Scan(&tag.AppID, &tag.BuildID, &tag.ConfigIDs, &tag.Force, &tag.Mark, &tag.Operator); err != nil {
		log.Error("RollBackTag(%v) err(%v)", tagID, err)
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
		}
	}
	return
}

// GetConfigIDS ...
func (s *Service) GetConfigIDS(tagIDS []int64) (tags []*model.Tag, err error) {
	if err = s.dao.DB.Where("id in (?)", tagIDS).Find(&tags).Error; err != nil {
		log.Error("GetConfigIDS err(%v)", err)
	}
	return
}

// LastTasConfigDiff ...
func (s *Service) LastTasConfigDiff(tagID, appID, buildID int64) (tag *model.Tag, err error) {
	tag = new(model.Tag)
	if err = s.dao.DB.Where("id < ? and app_id = ? and build_id = ?", tagID, appID, buildID).Order("id desc").First(tag).Error; err != nil {
		log.Error("LastTasConfigDiff() error(%v)", err)
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	return
}
