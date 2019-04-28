package service

import (
	"context"
	"time"

	"go-common/library/log"
)

func (s *Service) flushproc() {
	for {
		time.Sleep(10 * time.Millisecond)
		aids, ok := <-s.aidsChan
		if !ok {
			log.Error("s.aidsChan is closed")
			break
		}
		s.upViewCache(aids)
	}
	s.waiter.Done()
}

func (s *Service) flushConsumeproc() {
	defer s.waiter.Done()
	maxAid, err := s.arcRPC.MaxAID(context.Background())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	var aids []int64
	for aid := maxAid; aid > 0; aid-- {
		if s.closed {
			close(s.aidsChan)
			break
		}
		aids = append(aids, aid)
		if len(aids) >= 20 {
			s.aidsChan <- aids
			aids = []int64{}
		}
	}
	s.aidsChan <- aids
}
