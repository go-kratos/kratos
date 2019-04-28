package dao

import (
	"context"
	"fmt"
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	user "go-common/app/service/bbq/user/api"
	"go-common/library/log"
)

// ModifyRelation .
func (d *Dao) ModifyRelation(c context.Context, mid, upMid int64, action int32) (res *user.ModifyRelationReply, err error) {
	res, err = d.userClient.ModifyRelation(c, &user.ModifyRelationReq{Mid: mid, UpMid: upMid, Action: action})
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("modify relation fail: mid=%d, up_mid=%d, action=%d", mid, upMid, action)))
	}
	return
}

// UserRelationList .
func (d *Dao) UserRelationList(c context.Context, userReq *user.ListRelationUserInfoReq, relationType int32) (res *v1.UserRelationListResponse, err error) {
	res = new(v1.UserRelationListResponse)

	var reply *user.ListUserInfoReply
	switch relationType {
	case user.Follow:
		reply, err = d.userClient.ListFollowUserInfo(c, userReq)
	case user.Fan:
		reply, err = d.userClient.ListFanUserInfo(c, userReq)
	case user.Black:
		reply, err = d.userClient.ListBlackUserInfo(c, userReq)
	}
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("user relation list fail: req=%s", userReq.String())))
		return
	}

	res.HasMore = reply.HasMore
	for _, v := range reply.List {
		userInfo := &v1.UserInfo{UserBase: *v.UserBase, FollowState: v.FollowState, CursorValue: v.CursorValue}
		res.List = append(res.List, userInfo)
	}

	return
}

// FetchFollowList 获取mid的所有关注up主
func (d *Dao) FetchFollowList(c context.Context, mid int64) (upMid []int64, err error) {

	res, err := d.userClient.ListFollow(c, &user.ListRelationReq{Mid: mid})
	if err != nil {
		log.Errorv(c, log.KV("log", "fetch follow list fail"))
		return
	}
	upMid = res.List
	return
}
