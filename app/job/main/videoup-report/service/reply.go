package service

import (
	"context"
	"time"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/app/job/main/videoup-report/model/email"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

func (s *Service) arcReply(c context.Context, a *archive.Archive, replySwitch int64) (err error) {
	if replySwitch != archive.ReplyOn && replySwitch != archive.ReplyOff {
		log.Error("arcReply aid(%d) archive replySwitch(%d) 's state is unknow!", a.ID, replySwitch)
		return
	}

	replyState, _ := s.dataDao.CheckReply(c, a.ID)

	//删除之前的重试机会，避免延迟重试覆盖最新结果
	//s.removeRetry(c, a.ID, email.RetryActionReply)
	if replySwitch == archive.ReplyOn {
		err = s.openReply(c, a, replyState)
	} else {
		err = s.closeReply(c, a, replyState)
	}
	return
}

func (s *Service) openReply(c context.Context, a *archive.Archive, oldState int64) (err error) {
	if a == nil {
		return
	}
	if err = s.dataDao.OpenReply(c, a.ID, a.Mid); err != nil {
		log.Error("openReply s.dataDao.OpenReply(%d,%d) error(%v)", a.ID, a.Mid, err)
		s.addRetry(c, a.ID, email.RetryActionReply, archive.ReplyOn, oldState)
		return
	}
	if oldState == archive.ReplyOn {
		return
	}
	if oldState != archive.ReplyOn && oldState != archive.ReplyOff {
		oldState = archive.ReplyDefault
	}

	s.logJob(a, archive.ReplyDesc[archive.ReplyOn], archive.ReplyDesc[oldState])
	return
}

func (s *Service) closeReply(c context.Context, a *archive.Archive, oldState int64) (err error) {
	if a == nil {
		return
	}
	if err = s.dataDao.CloseReply(c, a.ID, a.Mid); err != nil {
		log.Error("closeReply s.dataDao.CloseReply(%d,%d) error(%v)", a.ID, a.Mid, err)
		s.addRetry(c, a.ID, email.RetryActionReply, archive.ReplyOff, oldState)
		return
	}
	if oldState == archive.ReplyOff {
		return
	}
	if oldState != archive.ReplyOn && oldState != archive.ReplyOff {
		oldState = archive.ReplyDefault
	}

	s.logJob(a, archive.ReplyDesc[archive.ReplyOff], archive.ReplyDesc[oldState])
	return
}

func (s *Service) logJob(a *archive.Archive, action string, oldAction string) {
	info := &report.ManagerInfo{
		Uname:    "videoup-job",
		UID:      399,
		Business: archive.LogBusJob,
		Type:     archive.LogTypeReply,
		Oid:      a.ID,
		Action:   action,
		Ctime:    time.Now(),
		Index:    []interface{}{a.State},
		Content:  map[string]interface{}{"old": oldAction},
	}
	log.Info("logJob (%+v)", info)
	report.Manager(info)
}

func isOpenReplyState(state int) (val int) {
	if archive.NormalState(state) || state == archive.StateForbidFixed {
		val = 1
	}
	return
}
