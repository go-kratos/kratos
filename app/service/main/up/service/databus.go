package service

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"strconv"

	"go-common/app/service/main/up/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

func (s *Service) upConsumer() {
	defer s.wg.Done()
	var (
		msgs = s.upSub.Messages()
		mid  int64
		err  error
	)
	if s.c.Env == "pre" {
		return
	}
	log.Info("up-service start consume")
CLOSE:
	for msg := range msgs {
		if mid, err = strconv.ParseInt(msg.Key, 10, 64); err != nil {
			log.Error("up-service strconv.Atoi(%s) error(%v)", msg.Key, err)
		}
		if s.closeSub {
			log.Info("sub close")
			break CLOSE
		}
		s.workerQueue[s.ShardingQueueIndex(mid)] <- msg
	}
	for _, c := range s.workerQueue {
		close(c)
	}
}

//ShardingQueueIndex sharding queue index
func (s *Service) ShardingQueueIndex(mid int64) (i int) { ////注：使用校验和取模的原因是：使得获取的消息均匀散入不同的worker队列中，不直接使用模运算的原因,比如5取模，加入mid都是5的倍数话，全部消息会进入0号worker队列中。
	ch := crc32.ChecksumIEEE([]byte(strconv.FormatInt(mid, 10)))
	i = int(ch) % s.workerCount
	return
}

//Start start queue
func (s *Service) Start(c chan *databus.Message) { //队列消费
	defer s.wg.Done()
	var ctx = context.Background()
	for msg := range c {
		m := &model.Msg{}
		if err := json.Unmarshal(msg.Value, m); err != nil {
			log.Error("up-service json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			err = nil
		}
		if _, err := s.Edit(ctx, m.MID, m.IsAuthor, uint8(m.From)); err != nil {
			log.Error("s.Edit(%v,databus) error(%v)", m, err)
			err = nil
		}
		msg.Commit()
		s.upMo++
		log.Info("up-service upConsumer key(%s)|value(%s)|partition(%d)|offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
		<-s.tokenChan
	}
	fmt.Printf("up-service close s.workerQueue chan(%+v)\n", c)
}
