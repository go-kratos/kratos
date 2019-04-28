package upcrm

import "go-common/app/admin/main/up/model/upcrmmodel"

//QueryUpBaseInfo query db
func (d *Dao) QueryUpBaseInfo(mid int64, fields string) (result upcrmmodel.UpBaseInfo, err error) {
	err = d.crmdb.Select(fields).Where("mid=?", mid).Find(&result).Error
	return
}

//QueryUpBaseInfoBatchByMid query db
func (d *Dao) QueryUpBaseInfoBatchByMid(fields string, mid ...int64) (result []upcrmmodel.UpBaseInfo, err error) {
	err = d.crmdb.Select(fields).Where("mid in(?)", mid).Find(&result).Error
	return
}

//QueryUpBaseInfoBatchByID query db
func (d *Dao) QueryUpBaseInfoBatchByID(fields string, id ...int64) (result []upcrmmodel.UpBaseInfo, err error) {
	err = d.crmdb.Select(fields).Where("id in(?)", id).Find(&result).Error
	return
}
