package service

import (
	"context"

	"go-common/app/interface/main/reply/model/reply"
	"go-common/library/log"
)

const (
	_replyChanBuf = 10240
	_topRpChanBuf = 128
)

type replyChan struct {
	rps []*reply.Reply
}

type topRpChan struct {
	oid int64
	tp  int8
	rp  *reply.Reply
}

func (s *Service) cacheproc() {
	for {
		select {
		case msg := <-s.replyChan:
			if err := s.dao.Mc.AddReply(context.Background(), msg.rps...); err != nil {
				log.Error("s.mcache.AddReply error(%v)", err)
			}
		case msg := <-s.topRpChan:
			if err := s.dao.Mc.AddTop(context.Background(), msg.oid, msg.tp, msg.rp); err != nil {
				log.Error("s.mcache.AddTop error(%v)", err)
			}
		}
	}
}
