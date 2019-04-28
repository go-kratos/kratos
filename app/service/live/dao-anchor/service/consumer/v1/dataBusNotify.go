package v1

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/common/log"

	"go-common/app/service/live/dao-anchor/conf"
	"go-common/app/service/live/dao-anchor/dao"
	"go-common/library/queue/databus"
)

type attrInfo struct {
	AttrSubId   int64  `json:"attr_sub_id"`
	ContentType string `json:"content_type"`
}

//ConsumerService struct
type ConsumerService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao                   *dao.Dao
	wg                    sync.WaitGroup
	liveDanmuSub          *databus.Databus
	liveGiftSendByPaySub  *databus.Databus
	liveGiftSendByFreeSub *databus.Databus
	liveGuardBuySub       *databus.Databus
	livePopularitySub     *databus.Databus
	liveValidLiveDays     *databus.Databus
	liveRoomTag           *databus.Databus
	liveRankList          *databus.Databus
}

//NewConsumerService init  统计类数据同步落地，实时数据异步落地
func NewConsumerService(c *conf.Config) (s *ConsumerService) {
	s = &ConsumerService{
		conf:                  c,
		dao:                   dao.New(c),
		liveDanmuSub:          databus.New(c.LiveDanmuSub),
		liveGiftSendByPaySub:  databus.New(c.LiveGiftSendByPaySub),
		liveGiftSendByFreeSub: databus.New(c.LiveGiftSendByFreeSub),
		liveGuardBuySub:       databus.New(c.LiveGuardBuySub),
		livePopularitySub:     databus.New(c.LivePopularitySub),
		liveValidLiveDays:     databus.New(c.LiveValidLiveDaysSub),
		liveRoomTag:           databus.New(c.LiveRoomTagSub),
		liveRankList:          databus.New(c.LiveRankListSub),
	}
	s.wg.Add(1)
	go s.LiveMessageContent(context.TODO())
	return s
}

//LiveMessageContent ...
func (s *ConsumerService) LiveMessageContent(ctx context.Context) {
	log.Info("LiveMessageContent_notify_start")
	defer func() {
		log.Warn("LiveMessageContent exists")
		s.wg.Done()
	}()
	//接收消息
	liveRankListMsgs := s.liveRankList.Messages()
	livePopularityMsgs := s.livePopularitySub.Messages()
	liveRoomTagMsgs := s.liveRoomTag.Messages()
	validLiveDaysMsgs := s.liveValidLiveDays.Messages()
	liveDanmuSendMsgs := s.liveDanmuSub.Messages()
	liveGiftSendByPayMsgs := s.liveGiftSendByPaySub.Messages()
	liveGiftSendByFreeMsgs := s.liveGiftSendByFreeSub.Messages()
	liveGuardBuyMsgs := s.liveGuardBuySub.Messages()
	for {
		select {
		case msg, ok := <-liveRankListMsgs:
			if err := s.doConsume(ctx, msg, ok, s.LiveRankList); err != nil {
				return
			}
		case msg, ok := <-livePopularityMsgs:
			if err := s.doConsume(ctx, msg, ok, s.PopularityStatistics); err != nil {
				return
			}
		case msg, ok := <-liveRoomTagMsgs:
			if err := s.doConsume(ctx, msg, ok, s.LiveRoomTag); err != nil {
				return
			}
		case msg, ok := <-validLiveDaysMsgs:
			if err := s.doConsume(ctx, msg, ok, s.ValidDaysStatistics); err != nil {
				return
			}
		case msg, ok := <-liveDanmuSendMsgs:
			if err := s.doConsume(ctx, msg, ok, s.DanmuCountStatistics); err != nil {
				return
			}
		case msg, ok := <-liveGiftSendByPayMsgs:
			if err := s.doConsume(ctx, msg, ok, s.GiftCount); err != nil {
				return
			}
		case msg, ok := <-liveGiftSendByFreeMsgs:
			if err := s.doConsume(ctx, msg, ok, s.GiftCount); err != nil {
				return
			}
		case msg, ok := <-liveGuardBuyMsgs:
			if err := s.doConsume(ctx, msg, ok, s.GuardBuyStatistics); err != nil {
				return
			}
		default:
			time.Sleep(time.Second * 2)
			continue
		}
	}

	log.Info("LiveMessageContent_notify_end")
}

func (s *ConsumerService) doConsume(ctx context.Context, msg *databus.Message, ok bool,
	fn func(context.Context, *databus.Message) error) (err error) {
	name := getFunctionName(fn)
	if msg == nil || !ok {
		err = errors.New(fmt.Sprintf("consume error ; msg:%v, ok:%v", msg, ok))
		log.Error("%s_notify_error:msg=%v", name, msg)
		return
	}
	log.Info("%s:key=%v;offset=%v;msg=%v", name, msg.Key, msg.Offset, msg.Value)
	if s.dao.CanConsume(ctx, msg) {
		fn(ctx, msg)
	}
	return
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
