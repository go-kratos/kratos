package dao

import (
	cml "go-common/app/admin/main/apm/model/canal"
	"go-common/library/ecode"
	"go-common/library/log"
)

//SetConfigID set canal_apply table conf_id
func (d *Dao) SetConfigID(id int64, addr string) (err error) {

	ups := map[string]interface{}{
		"conf_id": id,
	}
	if err = d.DBCanal.Model(&cml.Apply{}).Where("addr = ?", addr).Updates(ups).Error; err != nil {
		log.Error(" SetConfigID  error(%v)", err)
		err = ecode.SetConfigIDErr
		return
	}
	return
}

//CanalInfoCounts count master_info
func (d *Dao) CanalInfoCounts(v *cml.ConfigReq) (cnt int, err error) {
	if err = d.DBCanal.Model(&cml.Canal{}).Where("addr=?", v.Addr).Count(&cnt).Error; err != nil {
		log.Error("apmSvc.CanalInfoCounts count error(%v)", err)
		err = ecode.RequestErr
		return
	}
	return
}

//CanalInfoEdit update master_info
func (d *Dao) CanalInfoEdit(v *cml.ConfigReq) (err error) {

	ups := map[string]interface{}{
		"remark":  v.Mark,
		"cluster": v.Project,
		"leader":  v.Leader,
	}
	if err = d.DBCanal.Model(&cml.Canal{}).Where("addr=?", v.Addr).Updates(ups).Error; err != nil {
		log.Error(" CanalInfoEdit update error(%v)", err)
		err = ecode.CanalApplyUpdateErr
		return
	}
	return
}

//CanalApplyCounts count canal_apply
func (d *Dao) CanalApplyCounts(v *cml.ConfigReq) (cnt int, err error) {

	if err = d.DBCanal.Model(&cml.Apply{}).Where("addr=?", v.Addr).Count(&cnt).Error; err != nil {
		log.Error("apmSvc.CanalApplyEdit count error(%v)", err)
		err = ecode.RequestErr
		return
	}
	return
}

//CanalApplyEdit update canal_apply
func (d *Dao) CanalApplyEdit(v *cml.ConfigReq, username string) (err error) {

	ups := map[string]interface{}{
		"remark":   v.Mark,
		"operator": username,
		"state":    1,
		"cluster":  v.Project,
		"leader":   v.Leader,
	}
	if err = d.DBCanal.Model(&cml.Apply{}).Where("addr=?", v.Addr).Updates(ups).Error; err != nil {
		log.Error(" CanalApplyEdit update error(%v)", err)
		err = ecode.CanalApplyUpdateErr
		return
	}
	return
}

//CanalApplyCreate insert into canal_apply
func (d *Dao) CanalApplyCreate(v *cml.ConfigReq, username string) (err error) {
	var (
		ap = &cml.Apply{
			Addr:     v.Addr,
			Remark:   v.Mark,
			State:    1,
			Operator: username,
			Cluster:  v.Project,
			Leader:   v.Leader,
		}
	)
	if err = d.DBCanal.Create(ap).Error; err != nil {
		log.Error("apSvc.CanalApplyCreate create error(%v)", err)
		err = ecode.CanalApplyErr
		return
	}
	return
}
