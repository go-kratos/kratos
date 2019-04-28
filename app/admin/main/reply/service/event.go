package service

import (
	"context"

	"go-common/app/admin/main/reply/model"
	"go-common/library/log"
)

func (s *Service) pubEvent(c context.Context, action string, mid int64, sub *model.Subject, rp *model.Reply, rpt *model.Report) (err error) {
	if err = s.dao.PubEvent(c, action, mid, sub, rp, rpt); err != nil {
		log.Error("s.dao.PubEvent(%s,%d,%v,%v,%v) error(%v)", action, mid, sub, rp, rpt, err)
	}
	return
}

func (s *Service) pubSearchReply(c context.Context, rps map[int64]*model.Reply, newState int32) (err error) {
	if err = s.dao.UpSearchReply(c, rps, newState); err != nil {
		log.Error("s.dao.UpSearchReply(%v,%d) error(%v)", rps, newState, err)
	}
	return
}

func (s *Service) pubSearchReport(c context.Context, rpts map[int64]*model.Report, rpState *int32) (err error) {
	if err = s.dao.UpSearchReport(c, rpts, rpState); err != nil {
		log.Error("s.dao.UpSearchReport(%v,%d) error(%v)", rpts, rpState, err)
	}
	return
}
