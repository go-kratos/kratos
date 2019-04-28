package v1

import (
	"context"

	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

func getDMconfig(ctx context.Context, sdm *SendMsg) error {
	g, ctxn := errgroup.WithContext(ctx)
	if md, ok := metadata.FromContext(ctxn); ok {
		md[metadata.Mid] = sdm.SendMsgReq.Uid
	}
	//获取用户的昵称颜色
	g.Go(func() error {
		return sdm.UserInfo.GetUnameColor(ctxn, sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid)
	})
	//获取特殊勋章
	g.Go(func() error {
		return sdm.UserInfo.MedalInfo.GetSpeicalMedal(ctxn, sdm.SendMsgReq.Uid, sdm.RoomConf.UID)
	})
	//获取用户等级RANK
	g.Go(func() error {
		return sdm.UserInfo.GetUserLevelRank(ctxn, sdm.SendMsgReq.Uid)
	})
	//获取用户头衔
	g.Go(func() error {
		return sdm.TitleConf.GetCommentTitle(ctxn)
	})
	//获取勋章对应主播的昵称
	if sdm.UserInfo.MedalInfo.RUID != 0 {
		g.Go(func() error {
			return sdm.UserInfo.MedalInfo.GetMedalanchorName(ctxn, sdm.UserInfo.MedalInfo.RUID)
		})
		//获取勋章对应主播的房间号
		g.Go(func() error {
			return sdm.UserInfo.MedalInfo.GetMedalRoomid(ctxn, sdm.UserInfo.MedalInfo.RUID)
		})
	}
	//获取用户气泡
	g.Go(func() error {
		return sdm.UserInfo.GetUserBubble(ctxn, sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid, sdm.SendMsgReq.Bubble, sdm.UserInfo.PrivilegeType)
	})
	return g.Wait()

}
