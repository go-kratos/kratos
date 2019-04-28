package v1

import (
	"context"
	"encoding/json"

	"go-common/app/service/live/xanchor/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

//消费弹幕消息的业务逻辑

//DanmuCountStatistics  弹幕数统计
func (s *ConsumerService) DanmuCountStatistics(ctx context.Context, msg *databus.Message) (err error) {
	defer msg.Commit()
	msgInfo := new(model.DanmuSendMessage)
	if err = json.Unmarshal(msg.Value, &msgInfo); err != nil {
		log.Error("DanmuCountStatistics_error:msg=%v;err=%v;msgInfo=%v", msg, err, msgInfo)
		return
	}
	msgContent := msgInfo.MsgContent
	roomId := msgContent.RoomId
	redisKeyInfo := s.Pdao.DanmuNumKey(roomId)
	s.dao.Incr(ctx, redisKeyInfo.RedisKey, 1, redisKeyInfo.TimeOut)
	return
}
