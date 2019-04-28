package service

import (
	"context"
	"time"

	v1 "go-common/app/service/main/vipinfo/api"
	"go-common/app/service/main/vipinfo/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Info get vipinfo by mid.
func (s *Service) Info(c context.Context, mid int64) (res *v1.ModelInfo, err error) {
	res = new(v1.ModelInfo)
	if mid <= 0 {
		return
	}
	var vdb *model.VipUserInfo
	if vdb, err = s.dao.Info(c, mid); err != nil {
		return
	}
	if vdb == nil {
		return
	}
	s.logicVipInfo(c, vdb, time.Now())
	s.logicFrozen(c, vdb)
	res = s.convertInfo(vdb)
	return
}

// Infos get vipinfos by mids
func (s *Service) Infos(c context.Context, mids []int64) (res map[int64]*v1.ModelInfo, err error) {
	res = make(map[int64]*v1.ModelInfo, len(mids))
	if len(mids) <= 0 {
		return
	}
	if len(mids) > 100 {
		err = ecode.RequestErr
		return
	}
	var vs map[int64]*model.VipUserInfo
	if vs, err = s.dao.Infos(c, mids); err != nil {
		return
	}
	var now = time.Now()
	for _, v := range vs {
		if v == nil {
			continue
		}
		s.logicVipInfo(c, v, now)
	}
	s.logicFrozens(c, vs)
	for mid, v := range vs {
		res[mid] = s.convertInfo(v)
	}
	return
}

func (s *Service) logicVipInfo(c context.Context, v *model.VipUserInfo, now time.Time) {
	if v.VipType != model.Vip && v.VipType != model.AnnualVip {
		return
	}
	if !v.AnnualVipOverdueTime.Time().IsZero() && v.AnnualVipOverdueTime.Time().Before(now) {
		v.VipType = model.Vip
	}
	if !v.VipOverdueTime.Time().IsZero() && v.VipOverdueTime.Time().Before(now) {
		v.VipStatus = model.VipStatusOverTime
	}
	/**
	 *  因目前无法获取iap实际的签约状态
	 *  超过下次扣费时间+48h的那一天的24点，即is_auto_renew=0
	 */
	if v.PayChannelID == model.IapPayChannelID &&
		v.VipPayType == model.AutoRenewPay &&
		!v.IosOverdueTime.Time().IsZero() &&
		endOfDay(v.IosOverdueTime.Time().AddDate(0, 0, 2)).Before(time.Now()) {
		v.VipPayType = model.NormalPay
	}
}

//TODO 冻结逻辑 二期会再次进行改造(依赖增加vip-cache-job及冻结状态落库)
func (s *Service) logicFrozen(c context.Context, v *model.VipUserInfo) (err error) {
	var flag int
	if v.VipStatus == model.VipStatusNotOverTime {
		if flag, err = s.dao.CacheVipFrozen(c, v.Mid); err != nil {
			log.Error("get vip frozen err(%+v)", err)
			err = nil
		}
		if flag == 1 {
			v.VipStatus = model.VipStatusFrozen
		}
	}
	return
}

//TODO 冻结逻辑 二期会再次进行改造(依赖增加vip-cache-job及冻结状态落库)
func (s *Service) logicFrozens(c context.Context, vs map[int64]*model.VipUserInfo) (err error) {
	var (
		frozenFmap map[int64]int
		mids       = []int64{}
		flag       int
	)
	for _, v := range vs {
		if v.VipStatus == model.VipStatusNotOverTime {
			mids = append(mids, v.Mid)
		}
	}
	if len(mids) > 0 {
		if frozenFmap, err = s.dao.CacheVipFrozens(c, mids); err != nil {
			log.Error("get vip frozens err(%+v)", err)
			err = nil
		}
		if len(frozenFmap) == 0 {
			return
		}
		for mid, v := range vs {
			if flag = frozenFmap[mid]; flag == 1 {
				v.VipStatus = model.VipStatusFrozen
			}
		}
	}
	return
}

func (s *Service) convertInfo(v *model.VipUserInfo) (res *v1.ModelInfo) {
	res = new(v1.ModelInfo)
	res.VipPayType = v.VipPayType
	res.Type = v.VipType
	res.Status = v.VipStatus
	if !v.VipOverdueTime.Time().IsZero() {
		// 返回的过期时间戳与以前保持一致，单位:毫秒
		res.DueDate = v.VipOverdueTime.Time().Unix() * 1000
	}
	return
}
