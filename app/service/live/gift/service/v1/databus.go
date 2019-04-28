package v1

import (
	"context"
	"go-common/app/service/live/gift/model"
	"go-common/library/log"
	"strconv"
)

//SendAddGiftMsg 投递databus
func (s *GiftService) SendAddGiftMsg(ctx context.Context, uid, giftID, giftNum, expireAt int64, source, msgID string) {
	freeGift := &model.AddFreeGift{
		UID:      uid,
		GiftID:   giftID,
		GiftNum:  giftNum,
		ExpireAt: expireAt,
		Source:   source,
		MsgID:    msgID,
	}
	sendRet := s.addGift.Send(ctx, strconv.FormatInt(uid, 10), freeGift)
	log.Info("addFreeGift,ret:%v,params:%v", sendRet, freeGift)
}
