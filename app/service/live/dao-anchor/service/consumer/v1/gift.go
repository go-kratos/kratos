package v1

import (
	"context"
	"encoding/json"

	"go-common/app/service/live/dao-anchor/dao"
	"go-common/app/service/live/dao-anchor/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

//消费礼物消息的业务逻辑

//GiftCountStatistics  实时礼物相关统计，异步落地
func (s *ConsumerService) GiftCount(ctx context.Context, msg *databus.Message) (err error) {
	defer msg.Commit()
	msgInfo := new(model.MessageValue)
	if err = json.Unmarshal(msg.Value, &msgInfo); err != nil {
		log.Error("GiftCountStatistics_Unmarshal_error:msg=%v;err=%v;msgInfo=%v", msg.Value, err, msgInfo)
		return
	}
	log.Info("GiftCount:msgInfo=%v", string(msg.Value))
	msgContent := new(model.GiftSendMsgContent)
	if err = json.Unmarshal([]byte(msgInfo.MsgContent), &msgContent); err != nil {
		log.Error("GiftCountStatistics_Unmarshal_error:err=%v;bodyInfo=%v", err, msgContent)
		return
	}
	bodyInfo := msgContent.Body
	roomId := bodyInfo.RoomId
	num := bodyInfo.Num
	payCoin := bodyInfo.PayCoin
	//增加礼物数
	GiftNumKeyInfo := s.dao.SRoomRecordCurrent(dao.GIFT_NUM, roomId)
	s.dao.Incr(ctx, GiftNumKeyInfo.RedisKey, num, GiftNumKeyInfo.TimeOut)
	//增加金瓜子数量
	if bodyInfo.CoinType == "gold" {
		GiftGoldNumKeyInfo := s.dao.SRoomRecordCurrent(dao.GIFT_GOLD_NUM, roomId)
		s.dao.Incr(ctx, GiftGoldNumKeyInfo.RedisKey, num, GiftGoldNumKeyInfo.TimeOut)
	}
	//增加礼物金额
	if bodyInfo.CoinType == "gold" {
		GiftGoldAmountKeyInfo := s.dao.SRoomRecordCurrent(dao.GIFT_GOLD_AMOUNT, roomId)
		s.dao.Incr(ctx, GiftGoldAmountKeyInfo.RedisKey, payCoin, GiftGoldAmountKeyInfo.TimeOut)
	}
	return
}
