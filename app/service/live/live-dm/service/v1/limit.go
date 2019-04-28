package v1

import (
	"context"
	"go-common/app/service/live/live-dm/dao"
	"go-common/library/sync/errgroup"
)

func rateLimit(ctx context.Context, dm *SendMsg) error {
	lc := &dao.LimitCheckInfo{
		UID:     dm.SendMsgReq.Uid,
		RoomID:  dm.SendMsgReq.Roomid,
		Msg:     dm.SendMsgReq.Msg,
		MsgType: dm.SendMsgReq.Msgtype,
		Dao:     dm.Dmservice.dao,
		Conf:    dm.LimitConf,
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return lc.LimitPerSec(ctx)
	})
	g.Go(func() error {
		return lc.LimitSameMsg(ctx)
	})
	g.Go(func() error {
		return lc.LimitRoomPerSecond(ctx)
	})
	return g.Wait()
}
