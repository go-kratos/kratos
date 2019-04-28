package service

import (
	"context"

	"go-common/app/interface/main/dm2/model"
	account "go-common/app/service/main/account/api"
	thumbUpApi "go-common/app/service/main/thumbup/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// ThumbupDM like or cancel like a dm
func (s *Service) ThumbupDM(c context.Context, oid, dmid, mid int64, op int8) (err error) {
	if !model.CheckThumbup(op) {
		err = ecode.RequestErr
		return
	}
	// check dm state
	idxMap, _, err := s.dao.IndexsByid(c, model.SubTypeVideo, oid, []int64{dmid})
	if err != nil {
		log.Error("s.dao.IndexsByid(tp:%d,oid:%d,dmid:%d), error(%v)", 1, oid, dmid, err)
		return
	}
	if len(idxMap) <= 0 || !model.IsDMVisible(idxMap[dmid].State) {
		err = ecode.DMNotFound
		return
	}
	reply, err := s.accountRPC.Profile3(c, &account.MidReq{
		Mid: mid,
	})
	if err != nil {
		log.Error("s.accountRPC.Profile3(arg:%+v), error(%v)", mid, err)
		return
	}
	if reply.GetProfile().GetEmailStatus() == 0 && reply.GetProfile().GetTelStatus() == 0 {
		err = ecode.UserCheckNoPhone
		return
	}
	if reply.GetProfile().GetSilence() == 1 {
		err = ecode.DMActSilence
		return
	}
	sub, err := s.subject(c, model.SubTypeVideo, oid)
	if err != nil {
		return
	}
	arg2 := &thumbUpApi.LikeReq{
		Mid:       mid,
		UpMid:     sub.Mid,
		Business:  "danmu",
		OriginID:  oid,
		MessageID: dmid,
		Action:    thumbUpApi.Action(op),
	}
	_, err = s.thumbupRPC.Like(c, arg2)
	if err != nil {
		log.Error("dmAct s.thumbupRPC.Like(arg:%+v), error(%v)", arg2, err)
		return
	}
	return
}

// ThumbupList get list
func (s *Service) ThumbupList(c context.Context, oid, mid int64, dmids []int64) (res map[int64]*model.ThumbupStat, err error) {
	var (
		statsReply *thumbUpApi.StatsReply
	)
	if statsReply, err = s.thumbupRPC.Stats(c, &thumbUpApi.StatsReq{
		Business:   "danmu",
		OriginID:   oid,
		MessageIds: dmids,
		Mid:        mid,
		IP:         metadata.String(c, metadata.RemoteIP),
	}); err != nil {
		log.Error("dmAct s.thumbupRPC.StatsWithLike(oid:%+v,dmids:%+v), error(%v)", oid, dmids, err)
		return
	}
	res = make(map[int64]*model.ThumbupStat)
	if statsReply == nil {
		return
	}
	for id, li := range statsReply.Stats {
		st := new(model.ThumbupStat)
		st.Likes = li.LikeNumber
		st.UserLike = int8(li.LikeState)
		res[id] = st
	}
	return
}
