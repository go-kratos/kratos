package dao

import (
	"context"
	"fmt"
	user "go-common/app/service/bbq/user/api"
	"go-common/library/log"
)

// AddLike .
func (d *Dao) AddLike(c context.Context, mid, upMid, svid int64) (affectedNum int64, err error) {
	reply, err := d.userClient.AddLike(c, &user.LikeReq{Mid: mid, UpMid: upMid, Opid: svid})
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("add like fail: mid=%d, up_mid=%d, svid=%d", mid, upMid, svid)))
		return
	}
	affectedNum = reply.AffectedNum
	return
}

// CancelLike .
func (d *Dao) CancelLike(c context.Context, mid, upMid, svid int64) (affectedNum int64, err error) {
	reply, err := d.userClient.CancelLike(c, &user.LikeReq{Mid: mid, UpMid: upMid, Opid: svid})
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("cancel like fail: mid=%d, up_mid=%d, svid=%d", mid, upMid, svid)))
		return
	}
	affectedNum = reply.AffectedNum
	return
}

// CheckUserLike 检测用户是否点赞
func (d *Dao) CheckUserLike(c context.Context, mid int64, svids []int64) (res map[int64]bool, err error) {
	res = make(map[int64]bool)
	reply, err := d.userClient.IsLike(c, &user.IsLikeReq{Mid: mid, Svids: svids})
	if err != nil {
		log.Errorv(c, log.KV("log", "get is like info fail"))
		return
	}

	for _, svid := range reply.List {
		res[svid] = true
	}
	return
}

// UserLikeList .
func (d *Dao) UserLikeList(c context.Context, upMid int64, cursorPrev, cursorNext string) (res *user.ListUserLikeReply, err error) {
	res, err = d.userClient.ListUserLike(c, &user.ListUserLikeReq{UpMid: upMid, CursorPrev: cursorPrev, CursorNext: cursorNext})
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("get user like list fail: up_mid=%d, cursor_prev=%s, next=%s", upMid, cursorPrev, cursorNext)))
		return
	}
	return
}
