package service

import (
	"context"
	"fmt"

	"go-common/app/admin/main/apm/model/need"
	"go-common/library/ecode"

	"go-common/library/log"
)

//NeedInfoAdd is
func (s *Service) NeedInfoAdd(c context.Context, req *need.NAddReq, username string) (err error) {

	if err = s.dao.NeedInfoAdd(req, username); err != nil {
		return
	}
	return

}

//NeedInfoList is
func (s *Service) NeedInfoList(c context.Context, arg *need.NListReq, username string) (res []*need.NInfo, count int64, err error) {
	var (
		like *need.UserLikes
	)

	if count, err = s.dao.NeedInfoCount(arg); err != nil {
		return
	}
	if count == 0 {
		return
	}
	if res, err = s.dao.NeedInfoList(arg); err != nil {
		return
	}
	for _, r := range res {
		lq := &need.Likereq{
			ReqID: r.ID,
		}
		if like, err = s.dao.GetVoteInfo(lq, username); err == nil {
			r.LikeState = like.LikeType
		}
	}
	err = nil
	return
}

//NeedInfoEdit is
func (s *Service) NeedInfoEdit(c context.Context, arg *need.NEditReq, username string) (err error) {
	var (
		ni *need.NInfo
	)
	if ni, err = s.dao.GetNeedInfo(arg.ID); err != nil {
		err = ecode.NeedInfoErr
		return
	}
	if ni.Reporter != username {
		err = ecode.AccessDenied
		return
	}
	if ni.Status != 5 && ni.Status != 1 {
		err = ecode.NeedEditErr
		return
	}
	if err = s.dao.NeedInfoEdit(arg); err != nil {
		return
	}

	return
}

//NeedInfoVerify is
func (s *Service) NeedInfoVerify(c context.Context, arg *need.NVerifyReq) (ni *need.NInfo, err error) {

	if ni, err = s.dao.GetNeedInfo(arg.ID); err != nil {
		err = ecode.NeedInfoErr
		return
	}
	if ni.Status != 1 {
		err = ecode.NeedVerifyErr
		return
	}
	if err = s.dao.NeedVerify(arg); err != nil {
		return
	}

	return
}

//NeedInfoVote is
func (s *Service) NeedInfoVote(c context.Context, arg *need.Likereq, username string) (err error) {
	var (
		nu            *need.UserLikes
		like, dislike int
	)
	if _, err = s.dao.GetNeedInfo(arg.ReqID); err != nil {
		err = ecode.NeedInfoErr
		return
	}
	if nu, err = s.dao.GetVoteInfo(arg, username); err != nil {
		if err = s.dao.AddVoteInfo(arg, username); err != nil {
			return
		}
		switch arg.LikeType {
		case need.TypeLike:
			like = 1
			dislike = 0
		case need.TypeDislike:
			like = 0
			dislike = 1
		}
	} else {
		if err = s.dao.UpdateVoteInfo(arg, username); err != nil {
			log.Error("arg.liketype:%d,db.like%+v,like:%d,dislike:%d", arg.LikeType, nu, like, dislike)
			return
		}
		if nu.LikeType != arg.LikeType {
			switch {
			case arg.LikeType == need.TypeCancel && nu.LikeType == need.TypeLike:
				like = -1
			case arg.LikeType == need.TypeCancel && nu.LikeType == need.TypeDislike:
				dislike = -1
			case arg.LikeType == need.TypeLike && nu.LikeType == need.TypeCancel:
				like = 1
			case arg.LikeType == need.TypeDislike && nu.LikeType == need.TypeCancel:
				dislike = 1
			case arg.LikeType == need.TypeLike && nu.LikeType == need.TypeDislike:
				like = 1
				dislike = -1
			case arg.LikeType == need.TypeDislike && nu.LikeType == need.TypeLike:
				like = -1
				dislike = 1
			}
		} else {
			return
		}
	}
	if err = s.dao.LikeCountsStats(arg, like, dislike); err != nil {
		if nu != nil {
			log.Error("arg.liketype:%d,db.like:%+v,like:%d,dislike:%d", arg.LikeType, nu, like, dislike)
		}
		return
	}
	return
}

//NeedVoteList is show votelist
func (s *Service) NeedVoteList(c context.Context, arg *need.Likereq) (res []*need.UserLikes, count int64, err error) {
	if count, err = s.dao.VoteInfoCounts(arg); err != nil {
		return
	}
	if count == 0 {
		return
	}
	if res, err = s.dao.VoteInfoList(arg); err != nil {
		return
	}
	return
}

//SendWeMessage send  wechat message
func (s *Service) SendWeMessage(c context.Context, title, task, result, sender string, receiver []string) (err error) {
	var (
		msg   = "[sven需求与建议来信 📨 ]\n"
		users []string
	)
	switch task {
	case need.VerifyType[need.NeedApply]:
		msg += fmt.Sprintf("%s童鞋认真地提了一份建议(%s)，快去看看吧~\n%s\n", sender, title, "http://sven.bilibili.co/#/suggestion/list")
		users = append(users, receiver...)
	case need.VerifyType[need.NeedVerify]:
		msg += fmt.Sprintf("%s童鞋的建议(%s)反馈我们已经收到啦，先发一朵小红花感谢支持！🌺 \n", sender, title)
		users = append(users, sender)
	case need.VerifyType[need.NeedReview]:
		msg += fmt.Sprintf("%s童鞋的建议(%s)审核结果是%s, %s", receiver, title, result, "%s")
		users = append(users, receiver...)
	}

	switch result {
	case need.VerifyType[need.VerifyAccept]:
		msg = fmt.Sprintf(msg, "恭喜恭喜，喝杯快乐水开心一下~ 🍻 ")
	case need.VerifyType[need.VerifyReject]:
		msg = fmt.Sprintf(msg, "不要灰心，可能是需求描述或使用姿势不够准确，还请多多支持，欢迎再来！🙇‍  🙆‍ ")
	case need.VerifyType[need.VerifyObserved]:
		msg = fmt.Sprintf(msg, "您的意见我们先保留啦，还请多多支持，欢迎再来补充！🙇‍  🙆‍ ")
	}

	if err = s.dao.SendWechatToUsers(c, users, msg); err != nil {
		log.Error("apmSvc.SendWechatMessage error(%v)", err)
		return
	}
	return
}
