package service

import (
	"context"
	"encoding/json"

	"go-common/app/interface/main/history/model"
	"go-common/library/log"
)

func (s *Service) subproc() {
	for {
		msg, ok := <-s.sub.Messages()
		if !ok {
			log.Info("subproc exit")
			return
		}
		msg.Commit()
		m := &model.History{}
		if err := json.Unmarshal(msg.Value, &m); err != nil {
			log.Error("json.Unmarshal() error(%v)", err)
			continue
		}
		if m.Mid != 0 && m.Aid != 0 {
			s.add(m)
		}
	}
}

func (s *Service) add(m *model.History) {
	for j := 0; j < 3; j++ {
		err := s.dao.Add(context.Background(), m)
		if err == nil {
			return
		}
		log.Error("s.dao.Add() err:%+v", err)
	}
}
