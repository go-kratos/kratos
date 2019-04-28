package v1

import (
	"context"
	"encoding/json"

	"go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/app/service/live/dao-anchor/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

//LiveRoomTag  房间实时标签，同步落地
func (s *ConsumerService) LiveRoomTag(ctx context.Context, msg *databus.Message) (err error) {
	defer msg.Commit()
	return s.internalLiveRoomTag(ctx, msg)
}

func (s *ConsumerService) internalLiveRoomTag(ctx context.Context, msg *databus.Message) (err error) {
	msgContent := new(model.LiveRoomTagMsg)
	if err = json.Unmarshal(msg.Value, &msgContent); err != nil {
		log.Error("LiveRoomTag_error:msg=%v;err=%v;msgInfo=%v", msg, err, msgContent)
		return
	}
	log.Info("LiveRoomTag:msgInfo=%v", string(msg.Value))
	req := &v1.RoomTagCreateReq{
		RoomId:      msgContent.RoomId,
		TagId:       msgContent.TagId,
		TagSubId:    msgContent.TagSubId,
		TagExt:      msgContent.TagExt,
		TagExpireAt: msgContent.ExpireTime,
		TagValue:    msgContent.TagValue,
	}
	resp, err := s.dao.RoomTagCreate(ctx, req)
	if err != nil || resp.AffectedRows <= 0 {
		log.Error("LiveRoomTag: msg=%v; err=%v; req=%v", msg, err, req)
	}
	return
}
