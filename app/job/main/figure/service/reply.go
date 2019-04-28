package service

import (
	"context"
	"time"

	"go-common/app/job/main/figure/model"
)

// PutReplyInfo handle user reply info chenage message
func (s *Service) PutReplyInfo(c context.Context, info *model.ReplyEvent) (err error) {
	switch info.Action {
	case model.EventAdd:
		// only handle normal state reply
		if info.Reply.State == 0 {
			s.figureDao.PutReplyAct(c, info.Mid, model.ACColumnReplyAct, int64(1))
		}
	case model.EventLike:
		s.figureDao.PutReplyAct(c, info.Reply.Mid, model.ACColumnReplyLiked, int64(1))
	case model.EventLikeCancel:
		s.figureDao.PutReplyAct(c, info.Reply.Mid, model.ACColumnReplyLiked, int64(-1))
	case model.EventHate:
		s.figureDao.PutReplyAct(c, info.Reply.Mid, model.ACColumnReplyHate, int64(1))
	case model.EventHateCancel:
		s.figureDao.PutReplyAct(c, info.Reply.Mid, model.ACColumnReplyHate, int64(-1))
	case model.EventReportDel:
		s.figureDao.PutReplyAct(c, info.Report.Mid, model.ACColumnReplyReoprtPassed, int64(1))
		s.figureDao.PutReplyAct(c, info.Reply.Mid, model.ACColumnPublishReplyDeleted, int64(-1))
		s.figureDao.SetWaiteUserCache(c, info.Report.Mid, s.figureDao.Version(time.Now()))
	case model.EventReportRecover:
		s.figureDao.PutReplyAct(c, info.Report.Mid, model.ACColumnReplyReoprtPassed, int64(-1))
		s.figureDao.PutReplyAct(c, info.Reply.Mid, model.ACColumnPublishReplyDeleted, int64(1))
		s.figureDao.SetWaiteUserCache(c, info.Report.Mid, s.figureDao.Version(time.Now()))
	}
	return
}
