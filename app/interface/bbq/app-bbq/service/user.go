package service

import (
	"context"
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"
	user "go-common/app/service/bbq/user/api"
	"go-common/library/ecode"
	"go-common/library/log"
)

// UserBase 获取用户信息
func (s *Service) UserBase(c context.Context, mid int64) (res *v1.UserInfo, err error) {
	res = new(v1.UserInfo)
	var userInfos map[int64]*v1.UserInfo
	userInfos, err = s.dao.BatchUserInfo(c, mid, []int64{mid}, true, false, false)
	if err != nil {
		log.Errorv(c, log.KV("event", "user_base failed"), log.KV("err", err))
		err = ecode.GetUserBaseErr
		return
	}
	userInfo, exists := userInfos[mid]
	if exists && userInfo != nil {
		res = userInfo
	} else {
		log.Errorv(c, log.KV("event", "UserBase empty"), log.KV("err", err))
		err = ecode.GetUserBaseErr
	}
	return
}

//SpaceUserProfile ...
func (s *Service) SpaceUserProfile(c context.Context, mid int64, upMid int64) (res *v1.UserInfo, err error) {
	if upMid == 0 {
		err = ecode.UPMIDNotExists
		return
	}

	res = &v1.UserInfo{}
	userInfos, err := s.dao.BatchUserInfo(c, mid, []int64{upMid}, true, true, true)
	if err != nil {
		log.Errorv(c, log.KV("event", "batch_user_info"), log.KV("err", err))
		return
	}
	res, exists := userInfos[upMid]
	if !exists {
		err = ecode.UPMIDNotExists
		log.Errorv(c, log.KV("event", "user_info_not_found"), log.KV("up_mid", upMid), log.KV("err", err))
		return
	}

	return
}

// UserEdit 完善用户信息
//              该请求需要保证请求的mid已经存在，如果不存在，该接口会返回失败
func (s *Service) UserEdit(c context.Context, userBase *user.UserBase) (res *v1.NumResponse, err error) {
	res = new(v1.NumResponse)
	if err = s.dao.EditUserBase(c, userBase); err != nil {
		return
	}
	return
}

//AddUserLike 点赞
func (s *Service) AddUserLike(c context.Context, mid, svid int64) (res *v1.NumResponse, err error) {
	res = new(v1.NumResponse)
	var upMid, num int64

	videoBase, err := s.dao.VideoBase(c, mid, svid)
	if err != nil {
		log.Warnw(c, "log", "get video base fail", "svid", svid)
		return
	}

	upMid = videoBase.Mid

	likeNum, err := s.dao.AddLike(c, mid, upMid, svid)
	if err != nil {
		log.Warnv(c, log.KV("log", "add like fail"))
		return
	}
	if likeNum == 0 {
		log.Infow(c, "log", "repeated like", "mid", mid, "svid", svid)
		return
	}

	// TODO:考虑消息队列解耦等
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("add user like begin tran err(%v) mid(%d) svid(%d)", err, mid, svid)
		err = ecode.AddUserLikeErr
		return
	}
	if num, err = s.dao.TxIncrVideoStatisticsLike(tx, svid); err != nil {
		err = ecode.AddUserLikeErr
		tx.Rollback()
		return
	}
	if num == 0 {
		log.Errorw(c, "log", "incr video statistics fail maybe due to the missing in video_statistics table", "svid", svid)
		videoStatistics := &model.VideoStatistics{SVID: svid, Like: 1}
		if num, err = s.dao.TxAddVideoStatistics(tx, videoStatistics); err != nil || num == 0 {
			err = ecode.AddUserLikeErr
			tx.Rollback()
			return
		}
		res.Num = 1
	}
	tx.Commit()

	return
}

//CancelUserLike 取消点赞
func (s *Service) CancelUserLike(c context.Context, mid, svid int64) (res *v1.NumResponse, err error) {
	res = new(v1.NumResponse)
	var upMid, num int64
	videoBase, err := s.dao.VideoBase(c, mid, svid)
	if ecode.IsCancelSvLikeAvailable(err) {
		err = nil
		log.Infow(c, "log", "allow cancel like when video unreachable", "svid", svid, "mid", mid)
	} else if err != nil {
		log.Warnw(c, "log", "get video base fail", "svid", svid)
		return
	}
	upMid = videoBase.Mid

	likeNum, err := s.dao.CancelLike(c, mid, upMid, svid)
	if err != nil {
		log.Warnv(c, log.KV("log", "cancel like fail"))
		return
	}
	if likeNum == 0 {
		log.Infow(c, "log", "repeated cancel like", "mid", mid, "svid", svid)
		return
	}

	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("cancel user like begin tran err(%v) mid(%d) svid(%d)", err, mid, svid)
		err = ecode.CancelUserLikeErr
		return
	}
	if num, err = s.dao.TxDecrVideoStatisticsLike(tx, svid); err != nil {
		err = ecode.CancelUserLikeErr
		tx.Rollback()
		return
	}
	if num == 0 {
		log.Errorw(c, "log", "desc video statistics fail maybe due to the missing in video_statistics table", "svid", svid)
		videoStatistics := &model.VideoStatistics{SVID: svid, Like: 0}
		if num, err = s.dao.TxAddVideoStatistics(tx, videoStatistics); err != nil || num == 0 {
			err = ecode.AddUserLikeErr
			tx.Rollback()
			return
		}
		res.Num = 1
	}
	tx.Commit()
	return
}

// UserLikeList 用户点赞列表
// TODO: 这段代码和SpaceSvList基本一样，后面抽出共同代码
func (s *Service) UserLikeList(c context.Context, req *v1.SpaceSvListRequest) (res v1.SpaceSvListResponse, err error) {
	// 0.前期校验
	// 这里就不校验up主是否存在
	upMid := req.UpMid
	if upMid == 0 {
		err = ecode.ReqParamErr
		log.Errorv(c, log.KV("event", "req_param"), log.KV("up_mid", 0))
		return
	}
	res.List = make([]*v1.SvDetail, 0)

	// 1. 获取svid列表
	likeReply, err := s.dao.UserLikeList(c, upMid, req.CursorPrev, req.CursorNext)
	if err != nil {
		log.Errorv(c, log.KV("event", "user_like_list"), log.KV("error", err))
		return
	}
	res.HasMore = likeReply.HasMore

	var svids []int64
	// svid, LikeSv
	likeSvs := make(map[int64]*user.LikeSv)
	for _, likeSv := range likeReply.List {
		svids = append(svids, likeSv.Svid)
		likeSvs[likeSv.Svid] = likeSv
	}

	// 2.获取sv详情
	detailMap, err := s.getVideoDetail(c, req.MID, req.Qn, req.Device, svids, true)
	if err != nil {
		log.Warnv(c, log.KV("event", "get video detail fail"))
		return
	}
	for _, svID := range svids {
		item, exists := detailMap[svID]
		if !exists {
			item = new(v1.SvDetail)
			item.SVID = svID
			// TODO: 操蛋的客户端，这里等版本收敛之后可以删除
			item.IsLike = true
			item.Play.SVID = svID
			item.Play.FileInfo = []*v1.FileInfo{
				{},
			}
			item.Play.SupportQuality = make([]int64, 0)
		}
		if likeSv, exists := likeSvs[svID]; exists {
			item.CursorValue = likeSv.CursorValue
		}
		res.List = append(res.List, item)
	}

	return
}

// ModifyRelation 关注、取关、拉黑、取消拉黑
func (s *Service) ModifyRelation(c context.Context, mid, upMid int64, action int32) (res *user.ModifyRelationReply, err error) {
	res, err = s.dao.ModifyRelation(c, mid, upMid, action)
	if err != nil {
		log.Warnv(c, log.KV("log", "modify reltaion fail"))
	}
	return
}

// UserFollowList 关注列表
// 注意：这里出现3种mid，visitor/host/host关注的mid，按顺序缩写为V、H、HF，其中根据H获取HF，然后计算V和HF的关注关系
func (s *Service) UserFollowList(c context.Context, req *user.ListRelationUserInfoReq) (res *v1.UserRelationListResponse, err error) {
	res, err = s.dao.UserRelationList(c, req, user.Follow)
	return
}

// UserFanList 粉丝列表
func (s *Service) UserFanList(c context.Context, req *user.ListRelationUserInfoReq) (res *v1.UserRelationListResponse, err error) {
	res, err = s.dao.UserRelationList(c, req, user.Fan)
	return
}

// UserBlackList 黑名单列表
func (s *Service) UserBlackList(c context.Context, req *user.ListRelationUserInfoReq) (res *v1.UserRelationListResponse, err error) {
	// 强制只能访问自己的拉黑名单
	req.UpMid = req.Mid
	res, err = s.dao.UserRelationList(c, req, user.Black)
	return
}

//Login .
func (s *Service) Login(c context.Context, req *user.UserBase) (res *user.UserBase, err error) {
	res, err = s.dao.Login(c, req)
	return
}

//PhoneCheck ..
func (s *Service) PhoneCheck(c context.Context, mid int64) (err error) {
	telStatus, err := s.dao.PhoneCheck(c, mid)
	if err != nil {
		return
	}
	if telStatus == 0 {
		err = ecode.BBQNoBindPhone
		return
	}
	return
}
