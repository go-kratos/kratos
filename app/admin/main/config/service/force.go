package service

import (
	"context"

	"go-common/app/admin/main/config/model"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

//UpdateForce update force.
func (s *Service) UpdateForce(ctx context.Context, treeID, version int64, env, zone, build, username string, hosts map[string]string) (err error) {
	var (
		app   *model.App
		force *model.Force
		ups   map[string]interface{}
	)
	if app, err = s.AppByTree(treeID, env, zone); err != nil {
		return
	}
	tx := s.dao.DB.Begin()
	for key, val := range hosts {
		force = &model.Force{}
		force.Hostname = key
		force.AppID = app.ID
		force.IP = val
		force.Operator = username
		force.Version = version
		if err = s.dao.DB.Where("app_id = ? and hostname = ? and ip = ?", app.ID, key, val).First(&model.Force{}).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				tx.Rollback()
				log.Error("UpdateForce first error(%v)", err)
				return
			}
			//create
			if err = s.dao.DB.Create(force).Error; err != nil {
				tx.Rollback()
				log.Error("UpdateForce(%s) error(%v)", force, err)
				return
			}
		} else {
			//update
			ups = map[string]interface{}{
				"hostname": key,
				"app_id":   app.ID,
				"ip":       val,
				"operator": username,
				"version":  version,
			}
			if err = s.dao.DB.Model(&model.Force{}).Where("app_id = ? and hostname = ? and ip = ?", app.ID, key, val).Updates(ups).Error; err != nil {
				tx.Rollback()
				log.Error("UpdateForce(%s) error(%v)", force, err)
				return
			}
		}

	}
	if err = s.PushForce(ctx, treeID, env, zone, build, version, hosts, 1); err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

//ClearForce delete force.
func (s *Service) ClearForce(ctx context.Context, treeID int64, env, zone, build string, hosts map[string]string) (err error) {
	var (
		app *model.App
	)
	if app, err = s.AppByTree(treeID, env, zone); err != nil {
		return
	}
	tx := s.dao.DB.Begin()
	for key, val := range hosts {
		if err = s.dao.DB.Where("app_id = ? and hostname = ?", app.ID, key).Delete(model.Force{}).Error; err != nil {
			tx.Rollback()
			log.Error("ClearForce hostname(%s) ip(%v) error(%v)", key, val, err)
			return
		}
	}
	if err = s.PushForce(ctx, treeID, env, zone, build, 0, hosts, 0); err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}
