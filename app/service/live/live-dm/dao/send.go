package dao

import (
	"context"
	"fmt"
	"strconv"

	broadcasrtService "go-common/app/service/live/broadcast-proxy/api/v1"
	"go-common/app/service/live/live-dm/model"
	roomService "go-common/app/service/live/room/api/liverpc/v1"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_adminMsgHistoryCache = "cache:last10_roomadminmsg:"
	_msgHistoryCache      = "cache:last10_roommsg:"
)

//SaveHistory  弹幕历史存入redis
func (d *Dao) SaveHistory(ctx context.Context, hm string, adm bool, rid int64) {
	var conn = d.redis.Get(ctx)
	defer conn.Close()
	if adm {
		admKey := fmt.Sprintf("%s%d", _adminMsgHistoryCache, rid)
		if err := conn.Send("LPUSH", admKey, hm); err != nil {
			log.Error("DM:  SaveHistory  LPUSH err: %v", err)
		}
		if err := conn.Send("LLEN", admKey); err != nil {
			log.Error("DM:  SaveHistory  LLEN err: %v", err)
			return
		}
		if err := conn.Flush(); err != nil {
			log.Error("DM:  SaveHistory  Flush err: %v", err)
			return
		}

		if _, err := conn.Receive(); err != nil {
			log.Error("DM:  SaveHistory  LPUSH err: %v", err)
		}
		count, err := redis.Int64(conn.Receive())
		if err != nil {
			log.Error("DM:  SaveHistory  LPUSH  LLEN err: %v", err)
			return
		}

		if count > 15 {
			err := conn.Send("LTRIM", admKey, 0, 9)
			if err != nil {
				log.Error("DM:  SaveHistory LTRIM err: %v", err)
			}
		}
		if err := conn.Send("EXPIRE", admKey, 86400); err != nil {
			log.Error("DM:  SaveHistory EXPIRE err: %v", err)
		}
		if err := conn.Flush(); err != nil {
			log.Error("DM:  SaveHistory Flush err: %v", err)
		}
	}

	userKey := fmt.Sprintf("%s%d", _msgHistoryCache, rid)
	if err := conn.Send("LPUSH", userKey, hm); err != nil {
		log.Error("DM:  SaveHistory LPUSH err: %v", err)
	}
	if err := conn.Send("LLEN", userKey); err != nil {
		log.Error("DM:  SaveHistory LLEN err: %v", err)
		return
	}
	if err := conn.Flush(); err != nil {
		log.Error("DM:  SaveHistory Flush err: %v", err)
		return
	}

	if _, err := conn.Receive(); err != nil {
		log.Error("DM:  SaveHistory Receive  LPUSH err: %v", err)
	}
	count, err := redis.Int64(conn.Receive())
	if err != nil {
		log.Error("DM:  SaveHistory Int64 err: %v", err)
		return
	}
	if count > 15 {
		if err := conn.Send("LTRIM", userKey, 0, 9); err != nil {
			log.Error("DM:  SaveHistory  LTRIM err: %v", err)
		}
	}

	if err := conn.Send("EXPIRE", userKey, 86400); err != nil {
		log.Error("DM:  SaveHistory  EXPIRE err: %v", err)
	}

	if err := conn.Flush(); err != nil {
		log.Error("DM:  SaveHistory  Flush err: %v", err)
	}
}

//IncrDMNum 弹幕条数
func IncrDMNum(ctx context.Context, rid int64, mode int64) {
	req := &roomService.RoomIncrDanmuSendNumReq{
		RoomId: rid,
		Mode:   mode,
	}
	resp, err := RoomServiceClient.V1Room.IncrDanmuSendNum(ctx, req)
	if err != nil {
		log.Error("DM: IncrDMNum err: %v", err)
		return
	}
	if resp.Code != 0 {
		log.Error("DM: IncrDMNum err code: %d", resp.Code)
		return
	}
}

//SendBroadCast 发送弹幕(http)
func SendBroadCast(ctx context.Context, sm string, rid int64) error {
	err := LiveBroadCastClient.PushBroadcast(ctx, rid, 0, sm)
	if err != nil {
		log.Error("DM: SendBroadCast err: %v", err)
		return err
	}
	return nil
}

//SendBroadCastGrpc 调用GRPC发送弹幕
func SendBroadCastGrpc(ctx context.Context, sm string, rid int64) error {
	req := &broadcasrtService.RoomMessageRequest{
		RoomId:  int32(rid),
		Message: sm,
	}
	_, err := BcastClient.DanmakuClient.RoomMessage(ctx, req)
	if err != nil {
		log.Error("DM: SendBroadCastGrpc err: %v", err)
		return err
	}
	return nil
}

//SendBNDatabus 拜年祭制定房间投递到databus
func SendBNDatabus(ctx context.Context, uid int64, info *model.BNDatabus) {
	uids := strconv.FormatInt(uid, 10)
	if err := bndatabus.Send(ctx, uids, info); err != nil {
		log.Error("[service.live-dm.v1.bndatabus] send error(%v), record(%v)", err, info)
	}
}
