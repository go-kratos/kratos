package service

import (
	"context"

	memmdl "go-common/app/service/main/member/model"
	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// BindInfoByMid bind info by mid[bilibili->third].
func (s *Service) BindInfoByMid(c context.Context, a *model.ArgBindInfo) (res *model.BindInfo, err error) {
	if a.Mid <= 0 {
		return nil, ecode.RequestErr
	}
	res = new(model.BindInfo)
	var (
		m  *memmdl.BaseInfo
		b  *model.OpenBindInfo
		oi *model.OpenInfo
	)
	if m, err = s.memRPC.Base(c, &memmdl.ArgMemberMid{
		Mid: a.Mid,
	}); err != nil {
		return
	}
	res.Account = &model.BindAccount{
		Mid:  a.Mid,
		Name: m.Name,
		Face: m.Face,
	}
	res.Outer = new(model.BindOuter)
	if b, err = s.dao.BindInfoByMid(c, a.Mid, a.AppID); err != nil {
		return
	}
	if b == nil {
		return
	}
	if oi, err = s.dao.OpenInfoByMid(c, a.Mid, a.AppID); err != nil {
		return
	}
	if oi == nil {
		return nil, ecode.VipAssociateOpenIDNotExsitErr
	}
	res.Outer.BindState = b.State
	// get ele mobile
	var data *model.EleUnionMobileResp
	if data, err = s.dao.EleUnionMobile(c, &model.ArgEleUnionMobile{
		ElemeOpenID: b.OutOpenID,
		BliOpenID:   oi.OpenID,
	}); err != nil {
		return
	}
	if data != nil {
		res.Outer.Tel = data.BlurMobile
	}
	return
}

// OpenBindByOutOpenID associate user bind by out_open_id [third -> bilibili].
func (s *Service) OpenBindByOutOpenID(c context.Context, a *model.ArgBind) (err error) {
	var (
		oi          *model.OpenInfo
		ob          *model.OpenBindInfo
		byoutOpenID *model.OpenBindInfo
		bymid       *model.OpenBindInfo
	)
	// check open_id
	if oi, err = s.dao.RawOpenInfoByOpenID(c, a.OpenID, a.AppID); err != nil {
		return
	}
	if oi == nil {
		return ecode.VipAssociateOpenIDNotExsitErr
	}
	if ob, err = s.dao.BindInfoByOutOpenIDAndMid(c, oi.Mid, a.OutOpenID, a.AppID); err != nil {
		return
	}
	if ob != nil { //已绑定
		return
	}
	// check out_ipen_id
	if byoutOpenID, err = s.checkOutOpenID(c, a.OutOpenID, a.AppID); err != nil {
		return
	}
	// check mid
	if bymid, err = s.checkMid(c, oi.Mid, a.AppID); err != nil {
		return
	}
	err = s.addOpenBind(c, oi.Mid, a.OutOpenID, a.AppID, byoutOpenID, bymid)
	return
}

// OpenBindByMid associate user bind by mid [bilibili->third].
func (s *Service) OpenBindByMid(c context.Context, a *model.ArgOpenBindByMid) (err error) {
	var (
		ob          *model.OpenBindInfo
		oi          *model.OpenInfo
		data        *model.EleUnionUpdateOpenIDResp
		byoutOpenID *model.OpenBindInfo
		bymid       *model.OpenBindInfo
	)
	// check open_id
	if oi, err = s.dao.OpenInfoByMid(c, a.Mid, a.AppID); err != nil {
		return
	}
	if oi == nil {
		return ecode.VipAssociateOpenIDNotExsitErr
	}
	if ob, err = s.dao.BindInfoByOutOpenIDAndMid(c, oi.Mid, a.OutOpenID, a.AppID); err != nil {
		return
	}
	if ob != nil { //已绑定
		return
	}
	// check out_ipen_id
	if byoutOpenID, err = s.checkOutOpenID(c, a.OutOpenID, a.AppID); err != nil {
		return
	}
	// check mid
	if bymid, err = s.checkMid(c, oi.Mid, a.AppID); err != nil {
		return
	}
	arg := &model.ArgEleUnionUpdateOpenID{
		ElemeOpenID: a.OutOpenID,
		BliOpenID:   oi.OpenID,
	}
	if data, err = s.dao.EleUnionUpdateOpenID(c, arg); err != nil {
		return
	}
	if data.Status != 1 { //非1:绑定不成功
		log.Error("call ele update_open_id mid:%d params:%+v data:%+v", a.Mid, arg, data)
		return ecode.VipAssociateThirdBindErr
	}
	// bilibili bind
	err = s.addOpenBind(c, a.Mid, a.OutOpenID, a.AppID, byoutOpenID, bymid)
	return
}

// addOpenBind add open bind.
func (s *Service) addOpenBind(c context.Context, mid int64, outOpenID string, appID int64,
	byoutOpenID *model.OpenBindInfo, bymid *model.OpenBindInfo) (err error) {
	tx, err := s.dao.StartTx(c)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		if err == nil {
			if err = tx.Commit(); err != nil {
				tx.Rollback()
			}
		} else {
			tx.Rollback()
		}
		s.dao.DelBindInfoCache(c, mid, appID)
		if byoutOpenID != nil {
			s.dao.DelBindInfoCache(c, byoutOpenID.Mid, byoutOpenID.AppID)
		}
		if bymid != nil {
			s.dao.DelBindInfoCache(c, bymid.Mid, bymid.AppID)
		}
	}()
	// 解绑
	if byoutOpenID != nil {
		if err = s.dao.TxDeleteBindInfo(tx, byoutOpenID.ID); err != nil {
			return
		}
	}
	// 解绑
	if bymid != nil {
		if err = s.dao.TxDeleteBindInfo(tx, bymid.ID); err != nil {
			return
		}
	}
	if err = s.dao.TxAddBind(tx, &model.OpenBindInfo{
		Mid:       mid,
		AppID:     appID,
		OutOpenID: outOpenID,
	}); err != nil {
		return
	}
	return
}

// checkPay check mid had pay.
func (s *Service) checkPay(c context.Context, mid int64, appID int64, ver int64) (count int64, err error) {
	if count, err = s.dao.CountAssociateGrants(c, mid, appID); err != nil {
		return
	}
	if count <= 0 {
		if count, err = s.dao.CountAssociateOrder(c, mid, appID); err != nil {
			return
		}
	}
	if count > 0 && ver > 0 {
		if s.UpdateBindState(c, &model.OpenBindInfo{
			Mid:   mid,
			State: model.AssociateBindStatePurchased,
			AppID: appID,
			Ver:   ver,
		}); err != nil {
			return
		}
	}
	return
}

// checkMid check mid had pay.
func (s *Service) checkMid(c context.Context, mid int64, appID int64) (res *model.OpenBindInfo, err error) {
	var payCount int64
	if res, err = s.dao.RawBindInfoByMid(c, mid, appID); err != nil {
		return
	}
	if res != nil {
		if res.State == model.AssociateBindStatePurchased {
			return nil, ecode.VipAssociateBindPurchasedErr
		}
		if payCount, err = s.checkPay(c, res.Mid, res.AppID, res.Ver); err != nil {
			return
		}
		if payCount > 0 {
			return nil, ecode.VipAssociateBindPurchasedErr
		}
	}
	return
}

// checkOutOpenID check mid had pay.
func (s *Service) checkOutOpenID(c context.Context, outOpenID string, appID int64) (res *model.OpenBindInfo, err error) {
	var payCount int64
	if res, err = s.dao.ByOutOpenID(c, outOpenID, appID); err != nil {
		return
	}
	if res != nil {
		if res.State == model.AssociateBindStatePurchased {
			return nil, ecode.VipAssociateBindPurchasedErr
		}
		if payCount, err = s.checkPay(c, res.Mid, res.AppID, res.Ver); err != nil {
			return
		}
		if payCount > 0 {
			return nil, ecode.VipAssociateBindPurchasedErr
		}
	}
	return
}

// UpdateBindState update bind state.
func (s *Service) UpdateBindState(c context.Context, arg *model.OpenBindInfo) (err error) {
	if err = s.dao.UpdateBindState(c, arg); err != nil {
		return
	}
	s.dao.DelBindInfoCache(c, arg.Mid, arg.AppID)
	return
}
