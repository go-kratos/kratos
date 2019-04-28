package service

import (
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/main/config/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

//Apm apm.
func (s *Service) Apm(treeID int64, name, apmname, username string) (err error) {
	sns := []*model.ServiceName{}
	if err = s.DBApm.Where("name=?", apmname).Find(&sns).Error; err != nil {
		log.Error("svr.service.Apm name(%d) apmname(%d) error(%v)", name, apmname, err)
		return
	}
	for _, val := range sns {
		if err = s.ApmBuild(val, name, username, treeID); err != nil {
			log.Error("svr.service.ApmBuild val(%d) apmname(%d) error(%v)", val, apmname, err)
		}
	}
	return
}

//ApmBuild apmBuild.
func (s *Service) ApmBuild(val *model.ServiceName, name, username string, treeID int64) (err error) {
	bvs := []*model.BuildVersion{}
	mtime := time.Now().Unix() - (3600 * 24 * 60)
	if err = s.DBApm.Where("service_id=? and mtime > ?", val.ID, mtime).Find(&bvs).Error; err != nil {
		log.Error("svr.service.ApmBuild val(%v) id(%d) error(%v)", val, val.ID, err)
		return
	}
	if len(bvs) <= 0 {
		err = ecode.NothingFound
		return
	}
	for _, v := range bvs {
		var ver string
		scvs := []*model.ServiceConfigValue{}
		if err = s.DBApm.Where("config_id=?", v.ConfigID).Find(&scvs).Error; err != nil {
			log.Error("svr.service.ServiceConfigValue val(%d) ConfigID(%d) error(%v)", val, v.ConfigID, err)
			return
		}
		if len(scvs) <= 0 {
			err = ecode.NothingFound
			return
		}
		var env string
		switch val.Environment {
		case 10:
			env = "dev"
		case 11:
			env = "fat1"
		case 13:
			env = "uat"
		case 14:
			env = "pre"
		case 3:
			env = "prod"
		default:
			continue
		}
		version := strings.Split(v.Version, "-")
		if len(version) != 3 {
			continue
		}
		ver = version[1] + "-" + version[2]
		zone := version[0]
		switch zone {
		case "shylf":
			zone = "sh001"
		case "hzxs":
			zone = "sh001"
		case "shsb":
			zone = "sh001"
		default:
			continue
		}
		app := &model.App{
			Name:   name,
			TreeID: treeID,
			Env:    env,
			Zone:   zone,
			Token:  val.Token,
		}
		var tx *gorm.DB
		if tx = s.DB.Begin(); err != nil {
			log.Error("begin tran error(%v)", err)
			return
		}
		if err = tx.Where("tree_id=? and env=? and zone=?", treeID, env, zone).Find(&app).Error; err != nil {
			if err = tx.Create(&app).Error; err != nil {
				log.Error("svr.service.addapp create error(%v)", err)
				tx.Rollback()
				return
			}
		}
		configIds := ""
		for _, vv := range scvs {
			config := &model.Config{
				AppID:    app.ID,
				Name:     vv.Name,
				Comment:  vv.Config,
				From:     0,
				State:    2,
				Mark:     "一键迁移",
				Operator: username,
			}
			if err = tx.Create(&config).Error; err != nil {
				log.Error("svr.service.addconfig create error(%v)", err)
				tx.Rollback()
				return
			}
			if len(configIds) > 0 {
				configIds += ","
			}
			configIds += fmt.Sprint(config.ID)
		}
		tag := &model.Tag{
			AppID:     app.ID,
			BuildID:   0,
			ConfigIDs: configIds,
			Mark:      v.Remark,
			Operator:  username,
		}
		if err = tx.Create(&tag).Error; err != nil {
			log.Error("svr.service.addtag create error(%v)", err)
			tx.Rollback()
			return
		}
		buildNew := &model.Build{
			AppID:    app.ID,
			Name:     ver,
			TagID:    tag.ID,
			Mark:     v.Remark,
			Operator: username,
		}
		if err = tx.Create(&buildNew).Error; err != nil {
			log.Error("svr.service.addbuild create error(%v)", err)
			tx.Rollback()
			return
		}
		ups := map[string]interface{}{
			"build_id": buildNew.ID,
		}
		if err = tx.Model(tag).Where("id = ?", tag.ID).Updates(ups).Error; err != nil {
			log.Error("svr.service.edittag updates error(%v)", err)
			tx.Rollback()
			return
		}
		tx.Commit()
	}
	return
}
