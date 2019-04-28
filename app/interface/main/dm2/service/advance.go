package service

import (
	"context"
	"fmt"

	"go-common/app/interface/main/dm2/model"
	account "go-common/app/service/main/account/api"
	arcMdl "go-common/app/service/main/archive/model/archive"
	coinmdl "go-common/app/service/main/coin/model"
	relation "go-common/app/service/main/relation/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// BuyAdvance 购买高级弹幕
func (s *Service) BuyAdvance(c context.Context, mid, cid int64, mode string) (err error) {
	var (
		refund    = int64(1)
		coin      = model.AdvSPCoin
		reason    = model.AdvSPCoinReason
		advPermit int8
	)
	if !s.dao.AddAdvanceLock(c, cid, mid) {
		return
	}
	defer s.dao.DelAdvanceLock(c, cid, mid)
	reply, err := s.accountRPC.Info3(c, &account.MidReq{Mid: mid})
	if err != nil {
		log.Error("s.accRPC.Info3(%d) error(%v)", mid, err)
		return
	}
	if s.isSuperUser(reply.GetInfo().GetRank()) {
		return
	}
	if mode != model.AdvSpeMode { // pool=2 弹幕
		coin = model.AdvCoin
		reason = model.AdvCoinReason
	}
	typ, err := s.dao.AdvanceType(c, cid, mid, mode)
	if err != nil {
		log.Error("dao.AdvanceType(%d,%d,%s) error(%v)", cid, mid, mode, err)
		return
	}
	sub, err := s.subject(c, model.SubTypeVideo, cid)
	if err != nil {
		return
	}
	if typ != "" { // 已有购买记录
		if typ == model.AdvTypeRequest && mid != sub.Mid {
			err = ecode.DMAdvConfirm
		} else {
			err = ecode.DMAdvBought
		}
		return
	}
	if mid != sub.Mid {
		advPermit, err = s.dao.UpperConfig(c, sub.Mid)
		if err != nil {
			return
		}
		if err = s.checkAdvancePermit(c, advPermit, sub.Mid, mid); err != nil {
			return
		}
	}
	typ = model.AdvTypeRequest
	if sub.Mid == mid {
		typ = model.AdvTypeBuy
	}
	coins, err := s.coinRPC.UserCoins(c, &coinmdl.ArgCoinInfo{Mid: mid})
	if err != nil {
		log.Error("coinRPC.UserCoins(%v) error(%v)", mid, err)
		return
	}
	if coins < float64(coin) {
		err = ecode.LackOfCoins
		return
	}
	if _, err = s.dao.BuyAdvance(c, mid, cid, sub.Mid, refund, typ, mode); err != nil {
		return
	}
	if _, err = s.coinRPC.ModifyCoin(c, &coinmdl.ArgModifyCoin{Mid: mid, Count: -float64(coin), Reason: reason, CheckZero: 1}); err != nil {
		log.Error("coinRPC.ModifyCoin(%v,%v,%v) error(%v)", mid, coin, reason, err)
		return
	}
	if err = s.dao.DelAdvCache(c, mid, cid, mode); err != nil {
		return
	}
	s.cache.Do(c, func(ctx context.Context) {
		arc, err := s.arcRPC.Archive3(ctx, &arcMdl.ArgAid2{Aid: sub.Pid})
		if err != nil {
			log.Error("s.arcRPC.Archive3(aid:%d) error(%v)", sub.Pid, err)
			return
		}
		title := "您收到了一条高级弹幕请求"
		content := fmt.Sprintf(`您的稿件《%s》收到了一条高级弹幕请求，#{立即处理}{"https://member.bilibili.com/v/#/danmu/report/advance"}`, arc.Title)
		s.dao.SendNotify(ctx, title, content, "4", []int64{sub.Mid})
	})
	return
}

func (s *Service) checkAdvancePermit(c context.Context, advPermit int8, upid, mid int64) (err error) {

	switch advPermit {
	case model.AdvPermitAll:
		return
	case model.AdvPermitFollower, model.AdvPermitAttention:
		var (
			arg      = &relation.ArgMid{Mid: upid}
			res      = make([]*relation.Following, 0)
			follower *relation.Following
		)
		if res, err = s.relRPC.Followers(c, arg); err != nil {
			log.Error("relRPC.Followers(%+v) error(%v)", arg, err)
			return
		}
		for _, v := range res {
			if v.Mid == mid {
				follower = v
				break
			}
		}
		if follower == nil {
			err = ecode.DMAdvNotAllow
			return
		}
		if advPermit == model.AdvPermitAttention && follower.Attribute != 6 { // Attribute=6为相互关注
			err = ecode.DMAdvNotAllow
			return
		}
	case model.AdvPermitForbid:
		err = ecode.DMAdvNotAllow
		return
	}
	return
}

// AdvanceState 高级弹幕状态
func (s *Service) AdvanceState(c context.Context, mid, cid int64, mode string) (state *model.AdvState, err error) {
	state = &model.AdvState{
		Accept:  true,
		Coins:   model.AdvSPCoin,
		Confirm: model.AdvStatConfirmDefault,
	}
	if mode == model.AdvMode { // pool=2 弹幕
		state.Coins = model.AdvCoin
	}
	reply, err := s.accountRPC.Info3(c, &account.MidReq{Mid: mid})
	if err != nil {
		log.Error("s.accRPC.Info3(%d) error(%v)", mid, err)
		return
	}
	if s.isSuperUser(reply.GetInfo().GetRank()) {
		state.HasBuy = true
		state.Confirm = model.AdvStatConfirmAgree
		return
	}
	typ, err := s.dao.AdvanceType(c, cid, mid, mode)
	if err != nil {
		return
	}
	switch typ {
	case model.AdvTypeAccept, model.AdvTypeBuy: // 已通过
		state.Confirm = model.AdvStatConfirmAgree
		state.HasBuy = true
	case model.AdvTypeRequest: // 正在确认中
		state.Confirm = model.AdvStatConfirmRequest
		state.HasBuy = true
	case model.AdvTypeDeny: // up 主拒绝
		state.Confirm = model.AdvStatConfirmDeny
		state.HasBuy = true
	}
	return
}

// Advances 高级弹幕申请列表
func (s *Service) Advances(c context.Context, mid int64) (res []*model.Advance, err error) {
	var (
		cidMap    = make(map[int64]bool)
		midMap    = make(map[int64]bool)
		mids      = make([]int64, 0)
		cidAidMap = make(map[int64]int64)
		aids      = make([]int64, 0)
	)
	list, err := s.dao.Advances(c, mid)
	if err != nil {
		log.Error("dao.Advances(%d) error(%v)", mid, err)
		return
	}
	if len(list) == 0 {
		return
	}
	res = make([]*model.Advance, 0, len(list))
	for _, l := range list {
		if _, ok := midMap[l.Mid]; !ok {
			midMap[l.Mid] = true
			mids = append(mids, l.Mid)
		}
		cidMap[l.Cid] = true
	}
	for cid := range cidMap { // get cids->aids
		var sub *model.Subject
		if sub, err = s.subject(c, model.SubTypeVideo, cid); err != nil {
			log.Error("s.subject(%d) error(%v)", cid, err)
			return
		}
		if sub == nil {
			err = ecode.NothingFound
			continue
		}
		if _, ok := cidAidMap[sub.Oid]; !ok {
			cidAidMap[sub.Oid] = sub.Pid
			aids = append(aids, sub.Pid)
		}
	}
	arcs, err := s.archiveInfos(c, aids) // get archiveinfos
	if err != nil {
		return
	}
	reply, err := s.accountRPC.Infos3(c, &account.MidsReq{
		Mids: mids,
	})
	if err != nil {
		log.Error("s.accRPC.Infos3(%v) error(%v)", mids, err)
		return
	}
	for _, v := range list {
		if aid, ok := cidAidMap[v.Cid]; ok {
			v.Aid = aid
		} else {
			continue
		}
		if archive, ok := arcs[v.Aid]; ok {
			v.Title = archive.Title
			v.Cover = archive.Pic
		}
		if user, ok := reply.GetInfos()[v.Mid]; ok {
			v.Uname = user.Name
		}
		res = append(res, v)
	}
	return
}

// PassAdvance 通过高级弹幕申请
func (s *Service) PassAdvance(c context.Context, mid, id int64) (err error) {
	adv, err := s.dao.Advance(c, mid, id)
	if err != nil {
		log.Error("dao.Advance(%d,%d) error(%v)", mid, id, err)
		return
	}
	if adv == nil || adv.Type == model.AdvTypeDeny {
		err = ecode.DMAdvNoFound
		return
	}
	if adv.Type == model.AdvTypeAccept {
		return
	}
	if _, err = s.dao.UpdateAdvType(c, id, model.AdvTypeAccept); err != nil {
		log.Error("dao.UpdateAdvType(%d,%d,%s) error(%v)", mid, id, model.AdvTypeAccept, err)
		return
	}
	if err = s.dao.DelAdvCache(c, adv.Mid, adv.Cid, adv.Mode); err != nil {
		log.Error("dao.DelAdvCache(%+v) error(%v)", adv, err)
	}
	return
}

// DenyAdvance 拒绝高级弹幕申请
func (s *Service) DenyAdvance(c context.Context, mid, id int64) (err error) {
	var (
		coin   float64
		reason string
		af     int64
	)
	adv, err := s.dao.Advance(c, mid, id)
	if err != nil {
		log.Error("dao.Advance(%d,%d) error(%v)", mid, id, err)
		return
	}
	if adv == nil {
		return
	}
	if len(adv.Type) == 0 || adv.Type == model.AdvTypeDeny {
		err = ecode.DMAdvNoFound
		return
	}
	if af, err = s.dao.UpdateAdvType(c, id, model.AdvTypeDeny); err != nil {
		log.Error("dao.UpdateAdvType(%d) error(%v)", id, err)
		return
	}
	if err = s.dao.DelAdvCache(c, adv.Mid, adv.Cid, adv.Mode); err != nil {
		log.Error("dao.DelAdvCache(%+v) error(%v)", adv, err)
		return
	}
	if af < 1 {
		err = ecode.DMAdvNoFound
		return
	}
	if adv.Refund == 0 {
		return
	}
	coin = model.AdvSPCoin
	reason = model.AdvSPCoinCancelReason
	if adv.Mode == model.AdvMode {
		coin = model.AdvCoin
		reason = model.AdvCoinCancelReason
	}
	if _, err = s.coinRPC.ModifyCoin(c, &coinmdl.ArgModifyCoin{Mid: adv.Mid, Count: coin, Reason: reason}); err != nil {
		log.Error("s.accRPC.AddCoin2(%v,%v,%v) error(%v)", adv.Mid, coin, reason, err)
	}
	return
}

// CancelAdvance 取消高级弹幕申请
func (s *Service) CancelAdvance(c context.Context, mid, id int64) (err error) {
	var (
		adv    *model.Advance
		coin   float64
		reason string
		af     int64
	)
	if adv, err = s.dao.Advance(c, mid, id); err != nil {
		log.Error("s.dao.Advance(%d,%d) error(%v)", mid, id, err)
		return
	}
	if adv == nil {
		err = ecode.DMAdvNoFound
		return
	}
	if af, err = s.dao.DelAdvance(c, id); err != nil {
		log.Error("s.dao.DelAdvance(%d) error(%v)", id, err)
		return
	}
	if err = s.dao.DelAdvCache(c, adv.Mid, adv.Cid, adv.Mode); err != nil {
		log.Error("s.dao.DelAdvCache(%+v) error(%v)", adv, err)
		return
	}
	if af < 1 {
		err = ecode.DMAdvNoFound
		return
	}
	if adv.Refund == 0 || adv.Type == model.AdvTypeDeny {
		return
	}
	coin = model.AdvSPCoin
	reason = model.AdvSPCoinCancelReason
	if adv.Mode == model.AdvMode {
		coin = model.AdvCoin
		reason = model.AdvCoinCancelReason
	}
	if _, err = s.coinRPC.ModifyCoin(c, &coinmdl.ArgModifyCoin{Mid: adv.Mid, Count: coin, Reason: reason}); err != nil {
		log.Error("s.accRPC.AddCoin2(%v,%v,%v) error(%v)", adv.Mid, coin, reason, err)
	}
	return
}

// UpdateAdvancePermit update advance permit.
func (s *Service) UpdateAdvancePermit(c context.Context, mid int64, advPermit int8) (err error) {
	_, err = s.dao.AddUpperConfig(c, mid, advPermit)
	return
}

// AdvancePermit get advance permission.
func (s *Service) AdvancePermit(c context.Context, mid int64) (advPermit int8, err error) {
	advPermit, err = s.dao.UpperConfig(c, mid)
	return
}
