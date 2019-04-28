package v1

import (
	"context"

	"github.com/prometheus/common/log"

	"go-common/app/service/live/xanchor/conf"
	"go-common/app/service/live/xanchor/dao"
	"go-common/app/service/live/xanchor/dao/consumer"
	"go-common/library/queue/databus"
)

// XAnchorService struct
type ConsumerService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao             *consumer.Dao
	Pdao            *dao.Dao
	liveDanmuSub    *databus.Databus
	liveGiftSendSub *databus.Databus
}

//NewXAnchorService init
func NewXAnchorService(c *conf.Config) (s *ConsumerService) {
	s = &ConsumerService{
		conf:            c,
		dao:             consumer.New(c),
		Pdao:            dao.New(c),
		liveDanmuSub:    databus.New(c.LiveDanmuSub),
		liveGiftSendSub: databus.New(c.LiveDanmuSub),
	}
	return s
}

//DanmuSend  消费发送弹幕消息
func (s *ConsumerService) DanmuSend(ctx context.Context) {
	log.Info("danmu_send_notify_start")
	//接收消息
	liveDanmuMsgs := s.liveDanmuSub.Messages()
	msg, ok := <-liveDanmuMsgs
	if msg == nil || !ok {
		log.Error("xanchor_danmu_send_notify_error:msg=%v", liveDanmuMsgs)
		return
	}
	s.DanmuCountStatistics(ctx, msg)
	log.Info("danmu_send_notify_end")
}
func (s *ConsumerService) GiftSend(ctx context.Context) {
	log.Info("gift_send_notify_start")
	//接收消息
	liveGiftSendMsgs := s.liveGiftSendSub.Messages()
	msg, ok := <-liveGiftSendMsgs
	if msg == nil || !ok {
		log.Error("xanchor_figt_send_notify_error:msg=%v", liveGiftSendMsgs)
		return
	}
	s.GiftCountStatistics(ctx, msg)
	log.Info("gift_send_notify_end")
}
