package mcndao

import (
	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/library/log"
)

//GetActiveTid get tid from crm database
func (d *Dao) GetActiveTid(mids []int64) (res map[int64]int64, err error) {
	var infoList []*mcnmodel.UpBaseInfo
	err = d.mcndb.Select("mid, active_tid").Where("mid in (?) and business_type=1", mids).Find(&infoList).Error
	if err != nil {
		log.Error("fail to get active_tid from crm, err=%s", err)
		return
	}
	res = make(map[int64]int64, len(infoList))
	for _, v := range infoList {
		res[v.Mid] = v.ActiveTid
	}

	return
}

//GetUpBaseInfo get up base info from crm database
func (d *Dao) GetUpBaseInfo(fields string, mids []int64) (res map[int64]*mcnmodel.UpBaseInfo, err error) {
	var infoList []*mcnmodel.UpBaseInfo
	err = d.mcndb.Select(fields).Where("mid in (?) and business_type=1", mids).Find(&infoList).Error
	if err != nil {
		log.Error("fail to get active_tid from crm, err=%s", err)
		return
	}
	res = make(map[int64]*mcnmodel.UpBaseInfo, len(infoList))
	for _, v := range infoList {
		res[v.Mid] = v
	}

	return
}
