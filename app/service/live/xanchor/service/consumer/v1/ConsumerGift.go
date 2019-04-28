package v1

import (
	"context"
	"encoding/json"
	"go-common/app/service/live/xanchor/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

//消费礼物消息的业务逻辑

//GiftCountStatistics  礼物相关统计
func (s *ConsumerService) GiftCountStatistics(ctx context.Context, msg *databus.Message) (err error) {
	defer msg.Commit()
	msgInfo := new(model.GiftSendMessage)
	if err = json.Unmarshal(msg.Value, &msgInfo); err != nil {
		log.Error("GiftCountStatistics_error:msg=%v;err=%v;msgInfo=%v", msg, err, msgInfo)
		return
	}
	msgContent := msgInfo.MsgContent
	roomId := msgContent.RoomId
	num := msgContent.Num
	payCoin := msgContent.PayCoin
	//增加礼物数
	GiftNumKeyInfo := s.Pdao.GiftNumKey(roomId)
	s.dao.Incr(ctx, GiftNumKeyInfo.RedisKey, num, GiftNumKeyInfo.TimeOut)
	//增加金瓜子数量
	if msgContent.CoinType == "gold" {
		GiftGoldNumKeyInfo := s.Pdao.GiftGoldNumKey(roomId)
		s.dao.Incr(ctx, GiftGoldNumKeyInfo.RedisKey, num, GiftGoldNumKeyInfo.TimeOut)
	}
	//增加礼物金额
	if msgContent.CoinType == "gold" {
		GiftGoldAmountKeyInfo := s.Pdao.GiftGoldAmountKey(roomId)
		s.dao.Incr(ctx, GiftGoldAmountKeyInfo.RedisKey, payCoin, GiftGoldAmountKeyInfo.TimeOut)
	}
	return
}
