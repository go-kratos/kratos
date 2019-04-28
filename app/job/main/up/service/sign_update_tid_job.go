package service

import (
	"go-common/app/admin/main/up/util"
	"go-common/app/admin/main/up/util/mathutil"
	"go-common/app/job/main/up/model/signmodel"
	"go-common/app/job/main/up/model/upcrmmodel"
	"go-common/library/log"
	"time"
)

// UpdateUpTidJob 检查sign_up中的状态
func (s *Service) UpdateUpTidJob(date time.Time) {
	log.Info("start run UpdateUpTidJob, date=%s", date)
	s.updateUpTidJob(date)
	log.Info("finish run UpdateUpTidJob, date=%s", date)
}

func (s *Service) updateUpTidJob(date time.Time) {
	// 获取所有mid信息
	var limit = 200
	var count = limit
	var err error
	var mids []int64
	var offset = 0
	for limit == count {
		var signList []*signmodel.SignUp
		if err = s.crmdb.GetDb().Offset(offset).Limit(limit).Find(&signList).Error; err != nil {
			log.Error("fail to get signs from sign ups, err=%v", err)
			return
		}

		count = len(signList)
		offset += count

		for _, v := range signList {
			mids = append(mids, v.Mid)
		}
	}

	mids = util.Unique(mids)

	// 从up_base_info中读取tid并更新
	// <tid, mid list>
	var tidMidMap = make(map[int64][]int64)
	for begin := 0; begin < len(mids); begin += limit {
		var end = mathutil.Min(begin+limit, len(mids))
		var baseInfoList []*upcrmmodel.UpBaseInfo
		if err = s.crmdb.GetDb().Select("mid, active_tid").Where("mid in (?)", mids[begin:end]).Limit(limit).Find(&baseInfoList).Error; err != nil {
			log.Error("fail to get signs from sign ups, err=%v", err)
			return
		}
		// 更新到sign表中
		for _, v := range baseInfoList {
			tidMidMap[v.ActiveTid] = append(tidMidMap[v.ActiveTid], v.Mid)
		}
	}

	// 更新到sign_up表中
	for k, v := range tidMidMap {
		if err = s.crmdb.GetDb().Table(signmodel.TableNameSignUp).Where("mid in (?)", v).Update("active_tid", k).Error; err != nil {
			log.Error("update sign up's active tid fail, err=%v", err)
			return
		}
	}
}
