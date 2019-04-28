package service

import (
	"context"
	"time"

	"go-common/app/service/main/card/model"
	vipmol "go-common/app/service/main/vip/model"
	"go-common/library/ecode"
)

// UserCard get user card info.
func (s *Service) UserCard(c context.Context, mid int64) (res *model.UserCard, err error) {
	if mid <= 0 {
		return
	}
	var (
		e   *model.UserEquip
		cd  *model.Card
		ok  bool
		now = time.Now().Unix()
	)
	if e, err = s.dao.Equip(c, mid); err != nil {
		return
	}
	if ok, cd = s.checkEffective(e, now); !ok {
		return
	}
	res = &model.UserCard{
		Mid:          mid,
		ID:           cd.ID,
		CardURL:      cd.CardURL,
		BigCradURL:   cd.BigCradURL,
		CardType:     cd.CardType,
		Name:         cd.Name,
		ExpireTime:   e.ExpireTime,
		CardTypeName: model.CardTypeNameMap[cd.CardType],
	}
	return
}

// UserCards get user card infos.
func (s *Service) UserCards(c context.Context, mids []int64) (res map[int64]*model.UserCard, err error) {
	if len(mids) <= 0 {
		return
	}
	var (
		es  map[int64]*model.UserEquip
		cd  *model.Card
		ok  bool
		now = time.Now().Unix()
	)
	if es, err = s.dao.Equips(c, mids); err != nil {
		return
	}
	res = make(map[int64]*model.UserCard, len(es))
	for _, e := range es {
		if ok, cd = s.checkEffective(e, now); !ok {
			continue
		}
		res[e.Mid] = &model.UserCard{
			Mid:          e.Mid,
			ID:           cd.ID,
			CardURL:      cd.CardURL,
			CardType:     cd.CardType,
			Name:         cd.Name,
			ExpireTime:   e.ExpireTime,
			CardTypeName: model.CardTypeNameMap[cd.CardType],
		}
	}
	return
}

// check card effective
func (s *Service) checkEffective(e *model.UserEquip, now int64) (b bool, cd *model.Card) {
	var ok bool
	if e == nil {
		return
	}
	if cd, ok = s.cardmap[e.CardID]; !ok {
		return
	}
	if cd.CardType == model.CardTypeVip && e.ExpireTime < now {
		return
	}
	if _, ok = s.cardgroupmap[cd.GroupID]; !ok {
		return
	}
	b = true
	return
}

// Equip user equip card.
func (s *Service) Equip(c context.Context, arg *model.ArgEquip) (err error) {
	var (
		cd *model.Card
		ok bool
	)
	if cd, ok = s.cardmap[arg.CardID]; !ok {
		err = ecode.CardNotEffectiveErr
		return
	}
	if _, ok = s.cardgroupmap[cd.GroupID]; !ok {
		err = ecode.CardNotEffectiveErr
		return
	}
	e := new(model.UserEquip)
	e.CardID = arg.CardID
	e.Mid = arg.Mid
	if cd.CardType == model.CardTypeVip {
		var v *vipmol.VipInfoResp
		if v, err = s.vipRPC.VipInfo(c, &vipmol.ArgRPCMid{Mid: arg.Mid}); err != nil {
			return
		}
		if v == nil || v.VipStatus == vipmol.Expire || v.VipType == int8(vipmol.NotVip) {
			err = ecode.CardEquipNotVipErr
			return
		}
		e.ExpireTime = v.VipDueDate
	}
	if err = s.dao.CardEquip(c, e); err != nil {
		return
	}
	s.dao.DelCacheEquip(c, e.Mid)
	return
}

// DemountEquip delete equip.
func (s *Service) DemountEquip(c context.Context, mid int64) (err error) {
	if mid <= 0 {
		return
	}
	if err = s.dao.DeleteEquip(c, mid); err != nil {
		return
	}
	s.dao.DelCacheEquip(c, mid)
	return
}
