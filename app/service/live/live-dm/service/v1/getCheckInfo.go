package v1

import (
	"context"

	"go-common/library/sync/errgroup"
)

//getCheckMsg 获取弹幕检查信息
func getCheckMsg(ctx context.Context, sdm *SendMsg) error {

	g, ctxn := errgroup.WithContext(ctx)

	//获取弹幕配置
	g.Go(func() error {
		uid := sdm.SendMsgReq.Uid
		roomid := sdm.SendMsgReq.Roomid
		return sdm.DMconf.Get(ctxn, uid, roomid, sdm.Dmservice.conf)
	})
	//获取用户等级经验  TODO
	g.Go(func() error {
		return sdm.UserInfo.Get(ctxn, sdm.SendMsgReq.Uid)
	})
	//获取用户老爷等级
	g.Go(func() error {
		return sdm.UserInfo.GetVipInfo(ctxn, sdm.SendMsgReq.Uid)
	})
	//获取房间,大航海,房管信息
	g.Go(func() error {
		return getRoomConf(ctxn, sdm)
	})
	//获取用户绑定信息
	g.Go(func() error {
		return sdm.UserBindInfo.Get(ctxn, sdm.SendMsgReq.Uid)
	})
	//获取用户粉丝勋章信息
	g.Go(func() error {
		return sdm.UserInfo.MedalInfo.GetFansMedal(ctxn, sdm.SendMsgReq.Uid)
	})

	//获取弹幕真实分
	g.Go(func() error {
		return sdm.UserScore.GetMsgScore(ctxn, sdm.SendMsgReq.Msg)
	})
	//获取用户真实分
	g.Go(func() error {
		return sdm.UserScore.GetUserScore(ctx, sdm.SendMsgReq.Uid)
	})

	//获取房管信息
	g.Go(func() error {
		return sdm.UserInfo.IsRoomAdmin(ctx, sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid)
	})

	return g.Wait()
}

//GetRoomConf 获取房间,大航海,房管信息
func getRoomConf(ctx context.Context, sdm *SendMsg) error {
	if err := sdm.RoomConf.Get(ctx, sdm.SendMsgReq.Roomid); err != nil {
		return err
	}
	//获取大航海信息
	ruid := sdm.RoomConf.UID
	return sdm.UserInfo.GetPrivilegeType(ctx, sdm.SendMsgReq.Uid, ruid)
}
