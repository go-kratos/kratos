package dao

import (
	"context"
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"
	user "go-common/app/service/bbq/user/api"
	accountv1 "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
)

//GetUserBProfile 获取用户全量b站信息
func (d *Dao) GetUserBProfile(c context.Context, mid int64) (res *accountv1.ProfileReply, err error) {
	req := &accountv1.MidReq{
		Mid:    mid,
		RealIp: "",
	}
	res, err = d.accountClient.Profile3(c, req)
	return
}

// Login .
func (d *Dao) Login(c context.Context, userBase *user.UserBase) (res *user.UserBase, err error) {
	res, err = d.userClient.Login(c, userBase)
	if err != nil {
		log.Errorv(c, log.KV("log", "login fail"))
	}
	return
}

// BatchUserInfo 提供批量获取UserInfo的方法
// 由于user service返回的结构和video的回包不同，因此这里进行映射，返回video-c的结构，避免外部使用方多次映射
func (d *Dao) BatchUserInfo(c context.Context, visitorMID int64, MIDs []int64, needDesc, needStat, needFollowState bool) (res map[int64]*v1.UserInfo, err error) {
	res = make(map[int64]*v1.UserInfo)
	if len(MIDs) == 0 {
		return
	}
	if len(MIDs) > model.BatchUserLen {
		err = ecode.BatchUserTooLong
		return
	}

	userReq := &user.ListUserInfoReq{Mid: visitorMID, UpMid: MIDs, NeedDesc: needDesc, NeedStat: needStat, NeedFollowState: needFollowState}
	reply, err := d.userClient.ListUserInfo(c, userReq)
	if err != nil {
		log.Errorv(c, log.KV("log", "get user info fail: req=%s"+userReq.String()))
		return
	}

	for _, userInfo := range reply.List {
		newUserInfo := &v1.UserInfo{UserBase: *userInfo.UserBase}
		if userInfo.UserStat != nil {
			newUserInfo.UserStat = *userInfo.UserStat
		}
		newUserInfo.FollowState = userInfo.FollowState
		res[userInfo.UserBase.Mid] = newUserInfo
	}

	return
}

//JustGetUserBase 只取UserBase，不要其他
func (d *Dao) JustGetUserBase(c context.Context, mids []int64) (res map[int64]*user.UserBase, err error) {
	res = make(map[int64]*user.UserBase)
	userInfos, err := d.BatchUserInfo(c, 0, mids, false, false, false)
	if err != nil {
		log.Warnv(c, log.KV("log", "get user info fail"))
		return
	}

	for mid, userInfo := range userInfos {
		res[mid] = &userInfo.UserBase
	}

	return
}

// EditUserBase .
func (d *Dao) EditUserBase(c context.Context, userBase *user.UserBase) (err error) {
	_, err = d.userClient.UserEdit(c, userBase)
	if err != nil {
		log.Warnw(c, "log", "edit user base fail: req="+userBase.String(), "err", err.Error())
		return
	}
	return
}

// PhoneCheck .
func (d *Dao) PhoneCheck(c context.Context, mid int64) (telStatus int32, err error) {
	req := &user.PhoneCheckReq{Mid: mid}
	res, err := d.userClient.PhoneCheck(c, req)
	if err != nil {
		log.Errorw(c, "log", "call phone check fail", "err", err, "mid", mid)
		return
	}
	telStatus = res.TelStatus
	return
}
