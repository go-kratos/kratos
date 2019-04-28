package service

import (
	"context"
	"fmt"
	notice "go-common/app/service/bbq/notice-service/api/v1"
	"go-common/app/service/bbq/user/api"
	"go-common/app/service/bbq/user/internal/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AddLike 添加点赞
func (s *Service) AddLike(c context.Context, in *api.LikeReq) (res *api.LikeReply, err error) {
	res = new(api.LikeReply)
	mid := in.Mid
	svid := in.Opid
	upMid := in.UpMid

	// 0. 前期up主校验
	userReply, err := s.ListUserInfo(c, &api.ListUserInfoReq{UpMid: []int64{mid, upMid}})
	if err != nil {
		log.Errorw(c, "log", "get user info fail when adding like", "mid", mid, "up_mid", upMid, "err", err)
		return
	}
	if len(userReply.List) != 2 {
		log.Errorw(c, "log", "user error when adding like", "mid", mid, "up_mid", upMid, "user_info_reply", userReply.String())
		return
	}

	var num int64
	// TODO:考虑消息队列解耦等
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("add user like begin tran err(%v) mid(%d) svid(%d)", err, mid, svid)))
		err = ecode.AddUserLikeErr
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if num, err = s.dao.TxAddUserLike(tx, mid, svid); err != nil {
		err = ecode.AddUserLikeErr
		return
	}
	res.AffectedNum = num
	if num == 0 {
		return
	}

	if _, err = s.dao.TxIncrUserStatisticsField(c, tx, mid, "like_total"); err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("update user statistics like incr err(%v) mid(%d) svid(%d)", err, mid, svid)))
		return
	}
	if _, err = s.dao.TxIncrUserStatisticsField(c, tx, upMid, "rev_like_total"); err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("update user statistics rev_like incr err(%v) mid(%d) svid(%d)", err, mid, svid)))
		return
	}

	// TODO: 推送点赞给通知中心
	if upMid == mid {
		log.V(1).Infov(c, log.KV("log", "action_mid=mid"), log.KV("mid", mid))
	} else {
		notice := &notice.NoticeBase{Mid: upMid, ActionMid: mid, SvId: svid, Title: "点赞了你的作品", NoticeType: notice.NoticeTypeLike, BizType: notice.NoticeBizTypeSv}
		tmpErr := s.dao.CreateNotice(c, notice)
		if tmpErr != nil {
			log.Error("create like notice fail: notice_msg=%s", notice.String())
		}
	}

	return
}

// CancelLike 取消点赞
func (s *Service) CancelLike(c context.Context, in *api.LikeReq) (res *api.LikeReply, err error) {
	res = new(api.LikeReply)
	mid := in.Mid
	svid := in.Opid
	upMid := in.UpMid

	// 0. 前期up主校验
	userReply, err := s.ListUserInfo(c, &api.ListUserInfoReq{UpMid: []int64{mid, upMid}})
	if err != nil {
		log.Errorw(c, "log", "get user info fail when cancelling like", "mid", mid, "up_mid", upMid, "err", err)
		return
	}
	if len(userReply.List) != 2 {
		log.Errorw(c, "log", "user error when cancelling like", "mid", mid, "up_mid", upMid, "user_info_reply", userReply.String())
		return
	}

	var num int64
	// TODO:考虑消息队列解耦等
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("cancel user like begin tran err(%v) mid(%d) svid(%d)", err, mid, svid)))
		err = ecode.AddUserLikeErr
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if num, err = s.dao.TxCancelUserLike(tx, mid, svid); err != nil {
		err = ecode.AddUserLikeErr
		return
	}
	res.AffectedNum = num
	if num == 0 {
		return
	}

	if num, err = s.dao.TxDescUserStatisticsField(c, tx, mid, "like_total"); err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("update user statistics like incr err(%v) mid(%d) svid(%d)", err, mid, svid)))
		return
	} else if num == 0 {
		// 没能-1，也当做成功
		log.Errorv(c, log.KV("log", fmt.Sprintf("desc user statistics like incr err(%v) mid(%d) svid(%d)", err, mid, svid)))
		return
	}
	if num, err = s.dao.TxDescUserStatisticsField(c, tx, upMid, "rev_like_total"); err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("update user statistics rev_like incr err(%v) mid(%d) svid(%d)", err, mid, svid)))
		return
	} else if num == 0 {
		// 没能-1，也当做成功
		log.Errorv(c, log.KV("log", fmt.Sprintf("desc user statistics rev_like incr err(%v) mid(%d) svid(%d)", err, mid, svid)))
		return
	}

	return
}

// ListUserLike 用户的点赞列表
func (s *Service) ListUserLike(ctx context.Context, req *api.ListUserLikeReq) (res *api.ListUserLikeReply, err error) {
	res = new(api.ListUserLikeReply)
	// 0. 解析cursor
	cursor, cursorNext, err := model.ParseCursor(req.CursorPrev, req.CursorNext)
	if err != nil {
		return
	}

	// 1. 获取svid列表
	res.List, err = s.dao.GetUserLikeList(ctx, req.UpMid, cursorNext, cursor, model.UserListLen)
	if err != nil {
		log.Errorv(ctx, log.KV("event", "user_like_list"), log.KV("error", err))
		return
	}
	// 这里之所以跟len/2，是因为有可能同一时间多个点赞，这种情况会出现回包数量<len，但是这不代表has_more=false
	// TODO: 同样情况还有用户列表，后续考虑优化
	if len(res.List) < model.UserListLen/2 {
		res.HasMore = false
		if len(res.List) == 0 {
			return
		}
	} else {
		res.HasMore = true
	}

	return
}

// IsLike 返回是否点赞，点赞的才会返回
func (s *Service) IsLike(ctx context.Context, req *api.IsLikeReq) (res *api.IsLikeReply, err error) {
	res = new(api.IsLikeReply)
	res.List, err = s.dao.CheckUserLike(ctx, req.Mid, req.Svids)
	if err != nil {
		log.Warnv(ctx, log.KV("log", fmt.Sprintf("check user is like fail: req=%s", req.String())))
		return
	}
	return
}
