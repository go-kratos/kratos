package service

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	vipFrozenFlag = 1
)

// ByMid get vip userinfo by mid.
func (s *Service) ByMid(c context.Context, mid int64) (res *model.VipInfoResp, err error) {
	var (
		v *model.VipInfo
	)
	if v, err = s.VipInfo(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	s.logicVipStatusUpdate(c, v)
	s.logicAutoRenew(v)
	res = s.convertInfo(v)
	return
}

/**
 *  因目前无法获取iap实际的签约状态
 *  超过下次扣费时间+48h的那一天的24点，即is_auto_renew=0
 */
func (s *Service) logicAutoRenew(v *model.VipInfo) {
	if v.PayChannelID == model.IapPayChannelID &&
		v.VipPayType == model.AutoRenewPay &&
		!v.IosOverdueTime.Time().IsZero() &&
		EndOfDay(v.IosOverdueTime.Time().AddDate(0, 0, 2)).Before(time.Now()) {
		v.VipPayType = model.NormalPay
	}
}

//VipInfos def.
func (s *Service) VipInfos(c context.Context, mids []int64) (vMap map[int64]*model.VipInfoResp, err error) {
	vMap = make(map[int64]*model.VipInfoResp)
	if len(mids) == 0 {
		return
	}
	if len(mids) > _maxSizeUsers {
		err = ecode.RequestErr
		return
	}
	for _, v := range mids {
		var res *model.VipInfoResp
		if res, err = s.ByMid(context.TODO(), v); err != nil {
			err = nil
			continue
		}
		vMap[v] = res
	}
	return
}

// VipInfo .
func (s *Service) VipInfo(c context.Context, mid int64) (v *model.VipInfo, err error) {
	var (
		cache = true
		vdb   *model.VipInfoDB
	)
	if v, err = s.dao.VipInfoCache(c, mid); err != nil {
		log.Error("VipInfoCache(%d) err %+v", mid, err)
		cache = false
		err = nil
	}
	if v != nil && v.VipInfoDB != nil {
		return
	}
	v = new(model.VipInfo)
	if vdb, err = s.dao.VipInfo(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if vdb == nil {
		vdb = &model.VipInfoDB{
			Mid:       mid,
			VipType:   model.NotVip,
			VipStatus: model.VipStatusOverTime,
		}
	} else {
		var dh *model.VipUserDiscountHistory
		if dh, err = s.dao.DiscountSQL(c, mid, model.FirstDiscountBuyVip); err != nil {
			err = errors.WithStack(err)
			return
		}
		if dh != nil {
			v.AutoRenewed = model.IsAutoRenewed
		}
	}
	v.VipInfoDB = vdb
	if cache {
		s.cache(func() {
			if err1 := s.dao.SetVipInfoCache(context.TODO(), mid, v); err1 != nil {
				log.Error("SetVipInfoCache(%d) err %+v", mid, err1)
			}
		})
	}
	return
}

func (s *Service) logicVipStatusUpdate(c context.Context, v *model.VipInfo) (err error) {
	var (
		finalType   = v.VipType
		finalStatus = v.VipStatus
		now         = time.Now()
		val         int
	)
	if v.VipType != model.Vip && v.VipType != model.AnnualVip {
		return
	}
	if !v.AnnualVipOverdueTime.Time().IsZero() && v.AnnualVipOverdueTime.Time().Before(now) {
		finalType = model.Vip
	}
	if !v.VipOverdueTime.Time().IsZero() && v.VipOverdueTime.Time().Before(now) {
		finalStatus = model.VipStatusOverTime
	}
	if val, err = s.dao.GetVipFrozen(c, v.Mid); err != nil {
		log.Error("get vip frozen err(%+v)", err)
		err = nil
	}
	if finalStatus == model.VipStatusNotOverTime && val == vipFrozenFlag {
		finalStatus = model.VipStatusFrozen
	}
	if finalStatus != v.VipStatus || finalType != v.VipType {
		v.VipStatus = finalStatus
		v.VipType = finalType
	}
	return
}

func (s *Service) convertInfo(v *model.VipInfo) (res *model.VipInfoResp) {
	var (
		now = time.Now().Unix()
	)
	res = new(model.VipInfoResp)
	res.Mid = v.Mid
	res.PayType = int8(v.VipPayType)
	res.VipType = int8(v.VipType)
	res.PayChannelID = int32(v.PayChannelID)
	res.VipStatus = v.VipStatus
	if !v.VipOverdueTime.Time().IsZero() {
		res.VipDueDate = v.VipOverdueTime.Time().Unix()
	}
	if !v.VipOverdueTime.Time().IsZero() {
		res.VipSurplusMsec = (v.VipOverdueTime.Time().Unix() - now)
		if now > v.VipOverdueTime.Time().Unix() {
			res.VipDueMsec = (now - v.VipOverdueTime.Time().Unix())
		}
		if !v.VipStartTime.Time().IsZero() {
			res.VipTotalMsec = (v.VipOverdueTime.Time().Unix() - v.VipStartTime.Time().Unix())
			res.VipHoldMsec = (now - v.VipStartTime.Time().Unix())
		}
		if res.VipSurplusMsec > 0 {
			days := int(math.Ceil(float64(res.VipSurplusMsec) / float64(_daysecond)))
			if days > 0 && days < _remindday {
				res.DueRemark = fmt.Sprintf(_remindtxt, days)
			}
		}
	}
	if !v.VipRecentTime.Time().IsZero() {
		res.VipRecentTime = v.VipRecentTime.Time().Unix()
	}
	res.AutoRenewed = v.AutoRenewed
	return
}

// UpdateTypeAndStatus do update viptype or vipstatus
func (s *Service) UpdateTypeAndStatus(c context.Context, v *model.VipInfo) (err error) {
	if v.VipStatus == model.VipStatusFrozen {
		return
	}
	if _, err = s.dao.UpdateVipTypeAndStatus(c, v.Mid, v.VipStatus, v.VipType); err != nil {
		err = errors.WithStack(err)
		return
	}
	s.dao.DelVipInfoCache(c, v.Mid)
	return
}

//H5History .
func (s *Service) H5History(c context.Context, arg *model.ArgChangeHistory) (vh []*model.VipChangeHistoryVo, err error) {
	var (
		vcs []*model.VipChangeHistory
	)

	if vcs, err = s.dao.SelChangeHistory(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}

	if vh, err = s.fmtHistory(c, vcs); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//History .
func (s *Service) History(c context.Context, arg *model.ArgChangeHistory) (vh []*model.VipChangeHistoryVo, count int64, err error) {
	var (
		vcs []*model.VipChangeHistory
	)
	if count, err = s.dao.SelChangeHistoryCount(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}

	if vcs, err = s.dao.SelChangeHistory(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}

	if vh, err = s.fmtHistory(c, vcs); err != nil {
		err = errors.WithStack(err)
	}
	return
}

func (s *Service) fmtHistory(c context.Context, vcs []*model.VipChangeHistory) (vh []*model.VipChangeHistoryVo, err error) {
	var (
		relationIds []string
		actives     []*model.VipActiveShow
		actveMap    = make(map[string][]*model.VipActiveShow)
	)
	for _, v := range vcs {
		if len(v.RelationID) > 0 {
			relationIds = append(relationIds, v.RelationID)
		}
	}
	if actives, err = s.Actives(c, relationIds); err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, v := range actives {
		shows := actveMap[v.RelationID]
		shows = append(shows, v)
		actveMap[v.RelationID] = shows
	}

	for _, v := range vcs {
		r := new(model.VipChangeHistoryVo)
		r.ID = fmt.Sprintf("%d", v.ID)
		r.ChangeTime = v.ChangeTime.Time().Unix()
		r.ChangeType = int8(v.ChangeType)
		r.Days = int32(v.Days)
		year := r.Days / model.VipDaysYear
		month := (r.Days % model.VipDaysYear) / model.VipDaysMonth
		r.Month = int16(year*12 + month)
		r.OpenRemark = s.openRemark(v)
		r.ChangeTypeStr = model.OpenChangeMap[int8(v.ChangeType)]
		r.Remark = v.Remark
		r.Actives = actveMap[v.RelationID]
		vh = append(vh, r)
	}
	return
}

func (s *Service) openRemark(v *model.VipChangeHistory) string {
	var (
		buf          bytes.Buffer
		y, m, d, sub int
	)
	y = int(v.Days) / model.VipDaysYear
	sub = int(v.Days) % model.VipDaysYear
	m = sub / model.VipDaysMonth
	d = sub % model.VipDaysMonth
	buf.WriteString("")
	if y != 0 {
		buf.WriteString(strconv.Itoa(y))
		buf.WriteString("年")
	}
	if m != 0 {
		buf.WriteString(strconv.Itoa(m))
		buf.WriteString("个月")
	}
	if d != 0 {
		buf.WriteString(strconv.Itoa(d))
		buf.WriteString("天")
	}
	return buf.String()
}

// VipInfoBo vipinfo bo.
func (s *Service) VipInfoBo(c context.Context, mid int64) (bo *model.VipInfoBoResp, err error) {
	var (
		v *model.VipInfo
	)
	if v, err = s.VipInfo(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if v == nil {
		return
	}
	bo = new(model.VipInfoBoResp)
	bo.Mid = v.Mid
	bo.VipType = v.VipType
	bo.PayType = v.VipPayType
	bo.PayChannelID = v.PayChannelID
	bo.VipStatus = v.VipStatus
	if !v.VipStartTime.Time().IsZero() {
		bo.VipStartTime = v.VipStartTime.Time().Unix()
	}
	if !v.VipOverdueTime.Time().IsZero() {
		bo.VipOverdueTime = v.VipOverdueTime.Time().Unix()
	}
	if !v.AnnualVipOverdueTime.Time().IsZero() {
		bo.AnnualVipOverdueTime = v.AnnualVipOverdueTime.Time().Unix()
	}
	if !v.VipRecentTime.Time().IsZero() {
		bo.VipRecentTime = v.VipRecentTime.Time().Unix()
	}
	if !v.IosOverdueTime.Time().IsZero() {
		bo.IosOverdueTime = v.IosOverdueTime.Time().Unix()
	}
	bo.AutoRenewed = v.AutoRenewed
	return
}

//SurplusFrozenTime surplus time.
func (s *Service) SurplusFrozenTime(c context.Context, mid int64) (stime int64, err error) {
	var (
		frozenTime int64
		now        = time.Now()
	)
	if frozenTime, err = s.dao.FrozenTime(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if frozenTime > now.Unix() {
		stime = frozenTime - now.Unix()
	}
	return
}

//Unfrozen unfrozen .
func (s *Service) Unfrozen(c context.Context, mid int64) (err error) {
	var val int
	if val, err = s.dao.GetVipFrozen(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if val != vipFrozenFlag {
		err = ecode.VipUserUnFrozenErr
		return
	}
	if err = s.dao.DelRedisCache(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = s.dao.RemQueue(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = s.dao.CleanCache(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = s.dao.DelVipFrozen(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = s.dao.OldFrozenChange(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = s.dao.Loginout(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
