package service

import (
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"go-common/app/admin/main/apm/model/app"

	"github.com/jinzhu/gorm"
)

// AppAdd Appadd
func (s *Service) AppAdd(c *bm.Context, username string, AppTreeID int64, AppID string, Limit int64) (err error) {
	a := &app.App{}
	b := &app.Auth{}
	tx := s.DB.Begin()
	var sqlLogs []*map[string]interface{}
	if err = s.DB.Where("app_tree_id = ?", AppTreeID).First(a).Error; err == gorm.ErrRecordNotFound {
		//新加
		a = &app.App{
			AppTreeID: AppTreeID,
			AppID:     AppID,
			Limit:     Limit,
		}

		if err = tx.Create(a).Error; err != nil {
			log.Error("s.appAdd create error(%v)", err)
			tx.Rollback()
			return
		}
		sqlLog := &map[string]interface{}{
			"SQLType": "add",
			"Content": a,
		}
		sqlLogs = append(sqlLogs, sqlLog)
		aa := &app.App{}
		if a.AppID != "main.common-arch.msm-service" {
			if err = tx.Where("app_id=?", "main.common-arch.msm-service").First(aa).Error; err != nil {
				log.Error("s.appAdd not find main.common-arch.msm-service error(%v)", err)
				tx.Rollback()
				return
			}
			//查询授权
			if err = tx.Where("service_tree_id=? and app_tree_id=?", aa.AppTreeID, a.AppTreeID).First(b).Error; err != nil {
				//创建msm授权
				b = &app.Auth{
					ServiceTreeID: aa.AppTreeID,
					ServiceID:     aa.AppID,
					AppTreeID:     a.AppTreeID,
					AppID:         a.AppID,
					RPCMethod:     "ALL",
					HTTPMethod:    "ALL",
					Quota:         10000000,
				}
				if err = tx.Create(b).Error; err != nil {
					log.Error("s.appAdd main.common-arch.msm-service create error(%v)", err)
					tx.Rollback()
					return
				}
				sqlLog := &map[string]interface{}{
					"SQLType": "add",
					"Content": b,
				}
				sqlLogs = append(sqlLogs, sqlLog)
			}
		}
	} else if err != nil {
		log.Error("s.appAdd app_tree_id first error(%v)", err)
		tx.Rollback()
		return
	} else {
		//更新
		ups := map[string]interface{}{
			"app_id": AppID,
		}
		if err = tx.Model(a).Where("app_tree_id = ?", AppTreeID).Updates(ups).Error; err != nil {
			log.Error("s.appEdit updates error(%v)", err)
			tx.Rollback()
			return
		}
		sqlLog := &map[string]interface{}{
			"SQLType": "update",
			"Where":   "app_tree_id = ?",
			"Value1":  AppTreeID,
			"Update":  ups,
			"Old":     "",
		}
		sqlLogs = append(sqlLogs, sqlLog)
		var (
			services []*app.Auth
		)
		if err = tx.Where("app_tree_id = ?", AppTreeID).Find(&services).Error; err == nil {
			for _, v := range services {
				ups = map[string]interface{}{
					"app_id": AppID,
				}
				if err = tx.Model(&app.Auth{}).Where("id = ?", v.ID).Updates(ups).Error; err != nil {
					log.Error("s.appEdit auth2 app_tree_id updates error(%v)", err)
					tx.Rollback()
					return
				}
				sqlLog := &map[string]interface{}{
					"SQLType": "update",
					"Where":   "id = ?",
					"Value1":  v.ID,
					"Update":  ups,
					"Old":     v,
				}
				sqlLogs = append(sqlLogs, sqlLog)
			}
		}
		if err = tx.Where("service_tree_id=?", AppTreeID).Find(&services).Error; err == nil {
			for _, v := range services {
				ups = map[string]interface{}{
					"service_id": AppID,
				}
				if err = tx.Model(&app.Auth{}).Where("id=?", v.ID).Updates(ups).Error; err != nil {
					log.Error("s.appEdit auth2 service_tree_id updates error(%v)", err)
					tx.Rollback()
					return
				}
				sqlLog := &map[string]interface{}{
					"SQLType": "update",
					"Where":   "id = ?",
					"Value1":  v.ID,
					"Update":  ups,
					"Old":     v,
				}
				sqlLogs = append(sqlLogs, sqlLog)
			}
		}
	}
	tx.Commit()
	s.SendLog(*c, username, 0, 5, 0, "apmSvc.appAdd", sqlLogs)
	return
}
