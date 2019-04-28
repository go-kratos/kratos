package service

import (
	"context"
	"encoding/json"
	"fmt"
	notice "go-common/app/service/bbq/notice-service/api/v1"
	"go-common/app/service/bbq/user/api"
	"go-common/app/service/bbq/user/internal/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/time"
)

// ModifyRelation .
func (s *Service) ModifyRelation(c context.Context, arg *api.ModifyRelationReq) (res *api.ModifyRelationReply, err error) {
	mid := arg.Mid
	upMid := arg.UpMid
	res = new(api.ModifyRelationReply)
	if mid == upMid {
		err = ecode.FollowMyselfErr
		return
	}

	// 0. 前期up主校验
	userReply, err := s.ListUserInfo(c, &api.ListUserInfoReq{UpMid: []int64{mid, upMid}})
	if err != nil {
		log.Errorw(c, "log", "get user info fail when modifying relation", "mid", mid, "up_mid", upMid, "err", err)
		return
	}
	if len(userReply.List) != 2 {
		log.Errorw(c, "log", "user error when modifying relation", "mid", mid, "up_mid", upMid, "user_info_reply", userReply.String())
		return
	}

	// 1. 执行relation修改
	switch arg.Action {
	case api.FollowAdd:
		res.FollowState, err = s.addUserFollow(c, mid, upMid)
	case api.FollowCancel:
		res.FollowState, err = s.cancelUserFollow(c, mid, upMid)
	case api.BlackAdd:
		res.FollowState, err = s.addUserBlack(c, mid, upMid)
	case api.BlackCancel:
		res.FollowState, err = s.cancelUserBlack(c, mid, upMid)
	default:
		err = ecode.ReqParamErr
		return
	}
	return
}

// ListFollowUserInfo 关注列表
// 注意：这里出现3种mid，visitor/host/host关注的mid，按顺序缩写为V、H、HF，其中根据H获取HF，然后计算V和HF的关注关系
func (s *Service) ListFollowUserInfo(ctx context.Context, in *api.ListRelationUserInfoReq) (res *api.ListUserInfoReply, err error) {
	res, err = s.relationUserInfoList(ctx, in, model.FollowListType)
	return
}

// ListFanUserInfo .
func (s *Service) ListFanUserInfo(ctx context.Context, in *api.ListRelationUserInfoReq) (res *api.ListUserInfoReply, err error) {
	res, err = s.relationUserInfoList(ctx, in, model.FanListType)
	return
}

// ListBlackUserInfo .
func (s *Service) ListBlackUserInfo(ctx context.Context, in *api.ListRelationUserInfoReq) (res *api.ListUserInfoReply, err error) {
	// 强制只能访问自己的拉黑名单
	in.UpMid = in.Mid
	res, err = s.relationUserInfoList(ctx, in, model.BlackListType)
	return
}

// ListFollow 返回全部关注数组
func (s *Service) ListFollow(ctx context.Context, in *api.ListRelationReq) (res *api.ListRelationReply, err error) {
	res = new(api.ListRelationReply)
	res.List, err = s.dao.FetchFollowList(ctx, in.Mid)
	return
}

// ListBlack 返回全部黑名单数组
func (s *Service) ListBlack(ctx context.Context, in *api.ListRelationReq) (res *api.ListRelationReply, err error) {
	res = new(api.ListRelationReply)
	res.List, err = s.dao.FetchBlackList(ctx, in.Mid)
	return
}

// relationUserInfoList 关注列表
// 注意：这里出现3种mid，visitor/host/host关注的mid，按顺序缩写为V、H、HF，其中根据H获取HF，然后计算V和HF的关注关系
// @Param
//              relationType: 1-follow; 2-fan; 4-black
func (s *Service) relationUserInfoList(c context.Context, req *api.ListRelationUserInfoReq, listType model.UserListType) (res *api.ListUserInfoReply, err error) {
	res = new(api.ListUserInfoReply)
	visitorMID := req.Mid
	hostMID := req.UpMid
	// 0.前期校验
	// 这里就不校验up主是否存在
	if hostMID == 0 {
		err = ecode.ReqParamErr
		log.Errorv(c, log.KV("event", "req_param"), log.KV("up_mid", 0))
		return
	}
	// parseCursor，该接口不支持向前查找，所以cursorPrev填空
	cursor, _, err := model.ParseCursor("", req.CursorNext)
	if err != nil {
		return
	}

	// 1. 获取用户列表
	var MID2TimeMap map[int64]time.Time
	var followedMIDs []int64
	// 关注列表和粉丝列表的区别只在于获取mid列表
	switch listType {
	case model.FollowListType:
		MID2TimeMap, followedMIDs, err = s.dao.FetchPartFollowList(c, hostMID, cursor, model.UserListLen)
	case model.FanListType:
		MID2TimeMap, followedMIDs, err = s.dao.FetchPartFanList(c, hostMID, cursor, model.UserListLen)
	case model.BlackListType:
		MID2TimeMap, followedMIDs, err = s.dao.FetchPartBlackList(c, hostMID, cursor, model.UserListLen)
	default:
		err = ecode.BBQSystemErr
		log.Errorv(c, log.KV("event", "error_list_type"), log.KV("list_type", listType))
		return
	}
	if err != nil {
		log.Errorv(c, log.KV("event", "fetch_relation_list"), log.KV("error", err))
		return
	}
	log.V(1).Infov(c, log.KV("event", "get_relation_list"), log.KV("rsp_size", len(followedMIDs)))
	// 这里之所以跟len/2，是因为有可能同一时间多个关注，这种情况会出现回包数量<len，但是这不代表has_more=false
	// TODO: like也有该情况，后续考虑优化，能出意外的就在粉丝列表这里
	if len(followedMIDs) < model.UserListLen/2 {
		res.HasMore = false
		if len(followedMIDs) == 0 {
			return
		}
	} else {
		res.HasMore = true
	}

	// 2. 获取user info
	userInfos, err := s.batchUserInfo(c, visitorMID, followedMIDs, &api.ListUserInfoConf{NeedDesc: false, NeedStat: false, NeedFollowState: true})
	if err != nil {
		log.Errorv(c, log.KV("event", "batch_user_info"), log.KV("err", err))
		return
	}

	// 3. form rsp
	for _, mid := range followedMIDs {
		userInfo, exists := userInfos[mid]
		if exists {
			res.List = append(res.List, userInfo)
		} else {
			log.Errorv(c, log.KV("event", "get_use_info"), log.KV("mid", mid))
		}
	}

	// 4. 后处理，为每个item添加cursor值
	var itemCursor model.CursorValue
	for _, item := range res.List {
		var mtime time.Time
		mtime, exists := MID2TimeMap[item.UserBase.Mid]
		if !exists {
			log.Errorv(c, log.KV("event", "relation_id_not_found"), log.KV("relation_mid", item.UserBase.Mid))
		} else {
			itemCursor.CursorID = item.UserBase.Mid
			itemCursor.CursorTime = mtime
		}
		jsonStr, _ := json.Marshal(itemCursor) // marshal的时候相信库函数，不做err判断
		item.CursorValue = string(jsonStr)
	}

	return
}

//addUserFollow 关注
func (s *Service) addUserFollow(c context.Context, mid, upMid int64) (followState int8, err error) {
	// TODO: 后续改成Userbase
	userBases, err := s.dao.UserBase(c, []int64{mid, upMid})
	if err != nil {
		log.Errorv(c, log.KV("event", "user_base"))
		return
	}
	// 两个人的mid是否都存在
	if len(userBases) != 2 {
		log.Errorv(c, log.KV("event", "mid_not_found"), log.KV("mid", mid), log.KV("up_mid", upMid), log.KV("rsp_size", len(userBases)))
		err = ecode.BBQSystemErr
		return
	}
	// 关注是否达到上限
	userStats, err := s.dao.RawBatchUserStatistics(c, []int64{mid})
	if err != nil {
		log.Errorv(c, log.KV("event", "user_statistics"), log.KV("mid", mid), log.KV("err", err))
		return
	}
	// 这种情况下为user_statistics中没有该条记录，可以选择继续下面流程
	if userStat, exist := userStats[mid]; !exist {
		log.Errorv(c, log.KV("event", "user_statistics_not_found"), log.KV("mid", mid))
	} else if userStat.Follow >= model.MaxFollowListLen {
		err = ecode.UserFollowLimitErr
		return
	}

	tx, err := s.dao.BeginTran(c)
	if err != nil {
		err = ecode.AddUserFollowErr
		log.Error("add user follow begin tran err(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			log.V(1).Infov(c, log.KV("event", "rollback"))
		} else {
			tx.Commit()
			log.V(1).Infov(c, log.KV("event", "commit"))
			// 相互关系
			followState = 1
			MIDMap := s.dao.IsFollow(c, upMid, []int64{mid})
			if MIDMap != nil {
				_, exists := MIDMap[mid]
				if exists {
					followState |= 2
				}
			}
			notice := &notice.NoticeBase{
				Mid: upMid, ActionMid: mid, NoticeType: notice.NoticeTypeFan, Title: "关注了你",
				BizType: notice.NoticeBizTypeUser}
			tmpErr := s.dao.CreateNotice(c, notice)
			if tmpErr != nil {
				log.Error("create follow notice fail: notice_msg=%s", notice.String())
			}
		}
	}()

	var num1, num2, cancelBlackNum int64
	if cancelBlackNum, err = s.dao.TxCancelUserBlack(c, tx, mid, upMid); err != nil {
		log.Errorv(c, log.KV("event", "cancel_user_black"), log.KV("err", err), log.KV("mid", mid), log.KV("up_mid", upMid))
		err = ecode.UserBlackErr
		return
	}
	if cancelBlackNum == 1 {
		err = ecode.UserAlreadyBlackFollowErr
		return
	}
	if num1, err = s.dao.TxAddUserFollow(c, tx, mid, upMid); err != nil {
		err = ecode.AddUserFollowErr
		return
	}
	if num2, err = s.dao.TxAddUserFan(tx, upMid, mid); err != nil {
		err = ecode.AddUserFollowErr
		return
	}
	if num1 == 0 && num2 == 0 {
		log.Infov(c, log.KV("log", fmt.Sprintf("already follow: mid=%d, up_mid=%d", mid, upMid)))
		return
	}
	// 关注和粉丝不匹配的情况，当做关注成功，但是统计数计不变化吧
	if num1 == 0 || num2 == 0 {
		log.Errorv(c, log.KV("log", "consistency error with follow and fan"), log.KV("event", "fatal"), log.KV("up_mid", upMid))
		return
	}

	// 更新关注数和粉丝数
	if _, err = s.dao.TxIncrUserStatisticsField(c, tx, mid, "follow_total"); err != nil {
		log.Errorv(c, log.KV("event", "incr_user_statistic"), log.KV("field", "follow_total"), log.KV("mid", mid), log.KV("err", err))
		err = ecode.AddUserFollowErr
		return
	}
	if _, err = s.dao.TxIncrUserStatisticsField(c, tx, upMid, "fan_total"); err != nil {
		log.Errorv(c, log.KV("event", "incr_user_statistic"), log.KV("field", "fan_total"), log.KV("up_mid", upMid), log.KV("err", err))
		err = ecode.AddUserFollowErr
	}

	return
}

//CancelUserFollow 取消关注
func (s *Service) cancelUserFollow(c context.Context, mid, upMid int64) (followState int8, err error) {
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		err = ecode.CancelUserFollowErr
		log.Error("cancel user follow begin tran err(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var num1, num2 int64
	if num1, err = s.dao.TxCancelUserFollow(c, tx, mid, upMid); err != nil {
		log.Error("cancel user follow err(%v)", err)
		err = ecode.CancelUserFollowErr
		return
	}
	if num2, err = s.dao.TxCancelUserFan(tx, upMid, mid); err != nil {
		log.Error("cancel user fan err(%v)", err)
		err = ecode.CancelUserFollowErr
		return
	}
	if num1 == 0 && num2 == 0 {
		return
	}
	if num1 == 0 || num2 == 0 {
		log.Errorv(c, log.KV("log", "consistency error with follow and fan"), log.KV("event", "fatal"), log.KV("up_mid", upMid))
		return
	}

	if _, err = s.dao.TxDecrUserStatisticsFollow(tx, mid); err != nil {
		log.Error("update user statistics follow decr err(%v) mid(%d) upmid(%d)", err, mid, upMid)
		return
	}
	if _, err = s.dao.TxDecrUserStatisticsFan(tx, upMid); err != nil {
		log.Error("update user statistics fan decr err(%v) mid(%d) upmid(%d)", err, mid, upMid)
		return
	}

	return
}

//AddUserBlack 拉黑
func (s *Service) addUserBlack(c context.Context, mid, upMid int64) (followState int8, err error) {
	// 0. 前期校验
	// 是否存在mid和up_mid
	// TODO: 后续改成Userbase
	userBases, err := s.dao.UserBase(c, []int64{mid, upMid})
	if err != nil {
		log.Errorv(c, log.KV("event", "user_base"))
		return
	}
	// 两个人的mid是否都存在
	if len(userBases) != 2 {
		log.Errorv(c, log.KV("event", "mid_not_found"), log.KV("mid", mid), log.KV("up_mid", upMid), log.KV("rsp_size", len(userBases)))
		err = ecode.BBQSystemErr
		return
	}
	// 黑名单数量是否满足要求
	userStats, err := s.dao.RawBatchUserStatistics(c, []int64{mid})
	if err != nil {
		log.Errorv(c, log.KV("event", "user_statistics"), log.KV("mid", mid), log.KV("err", err))
		return
	}
	// 这种情况下为user_statistics中没有该条记录，可以选择继续下面流程
	if userStat, exist := userStats[mid]; !exist {
		log.Errorv(c, log.KV("event", "user_statistics_not_found"), log.KV("mid", mid))
	} else if userStat.Black >= model.MaxBlacklistLen {
		err = ecode.UserBlackLimitErr
		return
	}

	// 若存在关注，则直接取消
	_, err = s.cancelUserFollow(c, mid, upMid)
	if err != nil {
		log.Errorv(c, log.KV("log", "black fail due to cancel user follow fail"))
		err = ecode.UserBlackErr
		return
	}

	// 1. 插入黑名单
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		err = ecode.BBQSystemErr
		log.Errorv(c, log.KV("event", "begin_transaction"), log.KV("mid", mid), log.KV("err", err))
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			log.V(1).Infov(c, log.KV("event", "rollback"))
		} else {
			tx.Commit()
			log.V(1).Infov(c, log.KV("event", "commit"))
		}
	}()

	var addBlackNum int64
	if addBlackNum, err = s.dao.TxAddUserBlack(c, tx, mid, upMid); err != nil {
		log.Errorv(c, log.KV("event", "add_user_black"), log.KV("err", err), log.KV("mid", mid), log.KV("up_mid", upMid))
		err = ecode.UserBlackErr
		return
	}
	if addBlackNum == 0 {
		log.Infov(c, log.KV("log", fmt.Sprintf("already black: mid=%d, up_mid=%d", mid, upMid)))
		return
	}
	if _, err = s.dao.TxIncrUserStatisticsField(c, tx, mid, "black_total"); err != nil {
		log.Errorv(c, log.KV("event", "incr_user_statistic"), log.KV("field", "black_total"), log.KV("mid", mid), log.KV("err", err))
		err = ecode.UserBlackErr
	}
	return
}

//CancelUserBlack 取消拉黑
func (s *Service) cancelUserBlack(c context.Context, mid, upMid int64) (followState int8, err error) {
	// 0. 前期校验
	// 是否存在mid和up_mid
	// TODO: 后续改成Userbase
	userBases, err := s.dao.UserBase(c, []int64{mid, upMid})
	if err != nil {
		log.Errorv(c, log.KV("event", "user_base"))
		return
	}
	// 两个人的mid是否都存在
	if len(userBases) != 2 {
		log.Errorv(c, log.KV("event", "mid_not_found"), log.KV("mid", mid), log.KV("up_mid", upMid), log.KV("rsp_size", len(userBases)))
		err = ecode.BBQSystemErr
		return
	}

	// 1. 取消黑名单
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		err = ecode.BBQSystemErr
		log.Errorv(c, log.KV("event", "begin_transaction"), log.KV("mid", mid), log.KV("err", err))
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			log.V(1).Infov(c, log.KV("event", "rollback"))
		} else {
			tx.Commit()
			log.V(1).Infov(c, log.KV("event", "commit"))
		}
	}()

	var cancelBlackNum int64
	if cancelBlackNum, err = s.dao.TxCancelUserBlack(c, tx, mid, upMid); err != nil {
		log.Errorv(c, log.KV("event", "cancel_user_black"), log.KV("err", err), log.KV("mid", mid), log.KV("up_mid", upMid))
		err = ecode.UserBlackErr
		return
	}
	if cancelBlackNum == 0 {
		return
	}

	if _, err := s.dao.TxDescUserStatisticsField(c, tx, mid, "black_total"); err != nil {
		log.Errorv(c, log.KV("event", "desc_user_statistic"), log.KV("field", "black_total"), log.KV("mid", mid), log.KV("err", err))
		err = ecode.UserBlackErr
	}
	return
}
