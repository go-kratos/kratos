package upcrm

import "go-common/app/admin/main/up/model/upcrmmodel"

//QueryPlayInfo query db
func (d *Dao) QueryPlayInfo(mid int64, busiType []int) (result []upcrmmodel.UpPlayInfo, err error) {
	err = d.crmdb.Where("mid=? and business_type in (?)", mid, busiType).Find(&result).Error
	return
}

// QueryPlayInfoBatch query db
func (d *Dao) QueryPlayInfoBatch(mid []int64, busiType int) (result []*upcrmmodel.UpPlayInfo, err error) {
	err = d.crmdb.Where("mid in (?) and business_type = ?", mid, busiType).Find(&result).Error
	return
}
