package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/passport-sns/model"
	"go-common/library/log"
)

func (s *Service) snsLogConsume() {
	for {
		msg, ok := <-s.snsLogConsumer.Messages()
		if !ok {
			log.Error("s.snsLogConsumer.Messages closed")
			return
		}
		snsLog := &model.SnsLog{}
		if err := json.Unmarshal(msg.Value, snsLog); err != nil {
			log.Error("json.Unmarshal(%s) error(%+v)", string(msg.Value), err)
			continue
		}
		log.Info("receive msg snsLog(%+v)", snsLog)
		for {
			if _, err := s.d.AddSnsLog(context.Background(), snsLog); err != nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			break
		}
	}
}
