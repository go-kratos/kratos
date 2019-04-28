package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/reply/conf"
	"go-common/app/interface/main/reply/model/reply"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/queue/databus/report"
	xtime "go-common/library/time"
)

// AdminGetSubject get subject state
func (s *Service) AdminGetSubject(c context.Context, oid int64, tp int8) (sub *reply.Subject, err error) {
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	if sub, err = s.dao.Subject.Get(c, oid, tp); err != nil {
		log.Error("s.dao.Subject.Get(%d,%d) error(%v)", oid, tp, err)
		return
	}
	if sub == nil {
		err = ecode.NothingFound
	}
	return
}

// AdminSubRegist register  subject of new archive
func (s *Service) AdminSubRegist(c context.Context, oid, mid int64, tp, state int8, appkey string) (err error) {
	var has bool
	tps, ok := conf.Conf.AppkeyType[appkey]
	if !ok || len(tps) == 0 {
		err = ecode.ReplyIllegalSubType
		return
	}
	if tp != 0 {
		for _, t := range tps {
			if t == tp {
				has = true
				break
			}
		}
	} else {
		has = true
		tp = tps[0]
	}
	if !has {
		err = ecode.ReplyIllegalSubType
		return
	}
	now := time.Now()
	sub := &reply.Subject{
		Oid:   oid,
		Type:  tp,
		Mid:   mid,
		State: state,
		CTime: xtime.Time(now.Unix()),
		MTime: xtime.Time(now.Unix()),
	}
	if sub.ID, err = s.dao.Subject.Insert(c, sub); err != nil {
		log.Error("s.dao.Subject.Insert(%v) error(%v)", sub, err)
	}
	s.cache.Do(c, func(ctx context.Context) {
		if err = s.dao.Mc.DeleteSubject(ctx, oid, tp); err != nil {
			log.Error("s.dao.Mc.DeleteSubject(%d, %d) state:%d error(%v)", oid, tp, state, err)
		}
	})
	return
}

// AdminSubjectState change subject state by admin.
func (s *Service) AdminSubjectState(c context.Context, adid, oid, mid int64, tp, state int8, remark string) (err error) {
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	now := time.Now()
	if err = reply.CheckSubState(state); err != nil {
		log.Error("checkstate err(%v)", err)
		return
	}
	if mid <= 0 {
		var sub *reply.Subject
		sub, err = s.dao.Subject.Get(c, oid, tp)
		if err != nil {
			log.Error("s.dao.subject (%d,%d)error(%v)", oid, tp, err)
			return
		}
		if sub == nil {
			err = ecode.NothingFound
			return
		}
		mid = sub.Mid
	}
	if _, err = s.setSubject(c, oid, tp, state, mid); err != nil {
		log.Error("s.addSubject(%d, %d, %d, %d) error(%v)", oid, tp, state, mid, err)
		return
	}
	s.cache.Do(c, func(ctx context.Context) {
		if err = s.dao.Mc.DeleteSubject(ctx, oid, tp); err != nil {
			log.Error("s.dao.Mc.DeleteSubject(%d, %d) state:%d error(%v)", oid, tp, state, err)
		}
	})
	s.dao.Admin.Insert(c, adid, oid, 0, tp, fmt.Sprintf("修改主题状态为: %d", state), remark, reply.AdminIsNotNew, reply.AdminIsNotReport, reply.AdminOperSubState, now)
	return
}

// AdminSubjectMid set the subject mid info.
func (s Service) AdminSubjectMid(c context.Context, adid, mid, oid int64, tp int8, remark string) (err error) {
	// check subject
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	now := time.Now()
	if _, err = s.dao.Subject.UpMid(c, mid, oid, tp, now); err != nil {
		log.Error("replySubDao.UpMid(%d, %d, %d) error(%v)", mid, oid, tp, err)
		return
	}
	s.cache.Do(c, func(ctx context.Context) {
		if err = s.dao.Mc.DeleteSubject(ctx, oid, tp); err != nil {
			log.Error("s.dao.Mc.DeleteSubject(%d, %d) mid:%d error(%v)", oid, tp, mid, err)
		}
	})
	s.dao.Admin.Insert(c, adid, oid, 0, tp, fmt.Sprintf("修改主题mid为: %d", mid), remark, reply.AdminIsNotNew, reply.AdminIsNotReport, reply.AdminOperSubMid, now)
	return
}

// Delete delete reply by upper or self.
func (s *Service) Delete(c context.Context, mid, oid, rpID int64, tp int8, ak, ck, platform string, build int64, buvid string) (err error) {
	// check subject
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	rp, _ := s.reply(c, mid, oid, rpID, tp)
	if rp == nil {
		return
	}
	subject, err := s.getSubject(c, oid, tp)
	if err != nil {
		return
	}
	var (
		assist, operation bool
		state             = reply.ReplyStateUserDel
	)
	// check permission，only upper and self can del
	if !s.IsWhiteAid(subject.Oid, subject.Type) {
		assist, operation = s.CheckAssist(c, subject.Mid, mid)
		if !operation {
			assist = false
		}
		if assist {
			state = reply.ReplyStateAssistDel
		}
	}
	if subject.Mid == mid {
		state = reply.ReplyStateUpDel
	}
	if subject.Mid == mid || mid == rp.Mid || assist {
		if rp.IsDeleted() {
			s.dao.Redis.DelIndex(c, rp)
			err = ecode.ReplyDeleted
			return
		} else if rp.AttrVal(reply.ReplyAttrAdminTop) == 1 {
			err = ecode.ReplyDelTopForbidden
			return
		}
	} else {
		err = ecode.AccessDenied
		return
	}
	s.dao.Databus.Delete(c, mid, oid, rpID, time.Now().Unix(), tp, assist)
	remoteIP := metadata.String(c, metadata.RemoteIP)
	report.User(&report.UserInfo{
		Mid:      rp.Mid,
		Platform: platform,
		Build:    build,
		Buvid:    buvid,
		Business: 41,
		Type:     int(rp.Type),
		Oid:      rp.Oid,
		Action:   reply.ReportReplyDel,
		Ctime:    time.Now(),
		IP:       remoteIP,
		Index: []interface{}{
			rp.RpID,
			rp.State,
			state,
		},
	})
	return
}

// AdminEdit edit reply content by admin.
func (s *Service) AdminEdit(c context.Context, adid, oid, rpID int64, tp int8, msg, remark string) (err error) {
	now := time.Now()
	var rp *reply.Reply
	// check subject
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	if rp, err = s.reply(c, 0, oid, rpID, tp); err != nil {
		return
	}
	if rp.IsDeleted() {
		err = ecode.ReplyDeleted
		return
	}
	if _, err = s.dao.Content.UpMessage(c, oid, rpID, msg, now); err != nil {
		log.Error("s.content.UpMessage(%d, %d, %s, %v), err is (%v)", oid, rpID, msg, now, err)
		return
	}
	if err = s.dao.Mc.DeleteReply(c, rpID); err != nil {
		log.Error("s.dao.Mc.DeleteReply(%d, %d, %s, %v), err is (%v)", oid, rpID, msg, now, err)
	}
	// admin log
	if _, err = s.dao.Admin.UpIsNotNew(c, rpID, now); err != nil {
		log.Error("s.admin.UpIsNotNew(%d, %d, %s, %v), err is (%v)", oid, rpID, msg, now, err)
	}
	if _, err = s.dao.Admin.Insert(c, adid, oid, rpID, tp, "已修改评论内容", remark, reply.AdminIsNew, reply.AdminIsNotReport, reply.AdminOperEdit, now); err != nil {
		log.Error("s.admin.Insert(%d, %d, %s, %v), err is (%v)", oid, rpID, msg, now, err)
	}
	// dao.Kafka
	s.dao.Databus.AdminEdit(c, oid, rpID, tp)
	return
}

// AdminDelete delete reply by admin.
func (s *Service) AdminDelete(c context.Context, adid, oid, rpID, ftime int64, tp int8, moral int, notify bool, adname, remark string, reason, freason int8) (err error) {
	// check subject
	now := time.Now()
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	rp, _ := s.reply(c, 0, oid, rpID, tp)
	if rp == nil {
		log.Error("s.Reply(oid:%v,tp:%v,:rpID:%v)", oid, tp, rpID, err)
		return
	} else if rp.AttrVal(reply.ReplyAttrAdminTop) == 1 {
		err = ecode.ReplyDelTopForbidden
		return
	}
	s.dao.Databus.AdminDelete(c, adid, oid, rpID, ftime, moral, notify, adname, remark, now.Unix(), tp, reason, freason)
	return
}

// AdminPass recover reply by admin.
func (s *Service) AdminPass(c context.Context, adid, oid, rpID int64, tp int8, remark string) (err error) {
	// check subject
	now := time.Now()
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	s.dao.Databus.AdminPass(c, adid, oid, rpID, remark, now.Unix(), tp)
	return
}

// AdminRecover recover reply by admin.
func (s *Service) AdminRecover(c context.Context, adid, oid, rpID int64, tp int8, remark string) (err error) {
	// check subject
	now := time.Now()
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	s.dao.Databus.AdminRecover(c, adid, oid, rpID, remark, now.Unix(), tp)
	return
}

// AdminDeleteByReport delete report reply by admin.
func (s *Service) AdminDeleteByReport(c context.Context, adid, oid, rpID, ftime int64, tp int8, moral int, notify bool, adname, remark string, audit int8, reason int8, content string, freason int8) (err error) {
	var (
		rp  *reply.Reply
		now = time.Now()
	)
	// check subject
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	rpt, err := s.dao.Report.Get(c, oid, rpID)
	if err != nil {
		log.Error("s.report.GetReport(%d, %d) error(%v)", oid, rpID, err)
		return
	}
	if rpt == nil {
		err = ecode.ReplyReportNotExist
		return
	}
	if rp, err = s.dao.Reply.Get(c, oid, rpID); err != nil {
		log.Error("s.reply.GetReply(%d, %d) error(%v)", oid, rpID, err)
		return
	} else if rp == nil {
		err = ecode.ReplyNotExist
		return
	} else if rp.AttrVal(reply.ReplyAttrAdminTop) == 1 {
		err = ecode.ReplyDelTopForbidden
		return
	}
	s.dao.Databus.AdminDeleteByReport(c, adid, oid, rpID, rpt.Mid, ftime, moral, notify, adname, remark, now.Unix(), tp, audit, reason, content, freason)
	return
}

// AdminReportStateSet set report state by admin.
func (s *Service) AdminReportStateSet(c context.Context, adid, oid, rpID int64, tp, state int8) (err error) {
	now := time.Now()
	// check subject
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	rpt, err := s.dao.Report.Get(c, oid, rpID)
	if err != nil {
		log.Error("s.report.GetReport(%d, %d) met error (%v)", oid, rpID, err)
		return
	}
	if rpt == nil {
		err = ecode.ReplyReportNotExist
		return
	}
	// dao.Kafka
	s.dao.Databus.AdminStateSet(c, adid, oid, rpID, now.Unix(), tp, state)
	return
}

// AdminReportTransfer transfer report by admin.
func (s *Service) AdminReportTransfer(c context.Context, adid, oid, rpID int64, tp, audit int8) (err error) {
	now := time.Now()
	// check subject
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	rpt, err := s.dao.Report.Get(c, oid, rpID)
	if err != nil {
		log.Error("s.report.GetReport(%d, %d) met error (%v)", oid, rpID, err)
		return
	}
	if rpt == nil {
		err = ecode.ReplyReportNotExist
		return
	}
	// dao.Kafka
	s.dao.Databus.AdminTransfer(c, adid, oid, rpID, now.Unix(), tp, audit)
	return
}

// AdminReportIgnore ignore report by admin.
func (s *Service) AdminReportIgnore(c context.Context, adid, oid, rpID int64, tp, audit int8, remark string) (err error) {
	now := time.Now()
	// check subject
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	rpt, err := s.dao.Report.Get(c, oid, rpID)
	if err != nil {
		log.Error("s.report.GetReport(%d, %d) met error (%v)", oid, rpID, err)
		return
	}
	if rpt == nil {
		err = ecode.ReplyReportNotExist
		return
	}
	// dao.Kafka
	s.dao.Databus.AdminIgnore(c, adid, oid, rpID, now.Unix(), tp, audit)
	return
}

// AdminAddTop add top reply by admin
func (s *Service) AdminAddTop(c context.Context, adid, oid, rpID int64, tp, act int8) (err error) {
	var (
		ts = time.Now().Unix()
		r  *reply.Reply
	)
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	sub, err := s.Subject(c, oid, tp)
	if err != nil {
		log.Error("s.Subject(oid %v) err(%v)", oid, err)
		return
	}
	if r, err = s.GetTop(c, sub, oid, tp, reply.ReplyAttrAdminTop); err != nil {
		log.Error("s.GetTop(%d,%d) err(%v)", oid, tp, err)
		return
	}
	if r != nil && act == 1 {
		log.Warn("oid(%d) type(%d) already have top ", oid, tp)
		err = ecode.ReplyHaveTop
		return
	}
	if r == nil && act == 0 {
		log.Warn("oid(%d) type(%d) do not have top ", oid, tp)
		err = ecode.ReplyNotExist
		return
	}
	// TODO: only need reply,no not need content and user info
	if r, err = s.reply(c, 0, oid, rpID, tp); err != nil {
		log.Error("s.GetReply err (%v)", err)
		return
	}
	if r == nil {
		log.Warn("oid(%d) type(%d) rpID(%d) do not exist ", oid, tp, rpID)
		err = ecode.ReplyNotExist
		return
	}
	if r.AttrVal(reply.ReplyAttrUpperTop) == 1 {
		err = ecode.ReplyHaveTop
		return
	}
	if r.Root != 0 {
		log.Warn("oir(%d) type(%d) rpID(%d) not root reply", oid, tp, rpID)
		err = ecode.ReplyNotRootReply
		return
	}
	s.dao.Databus.AdminAddTop(c, adid, oid, rpID, ts, act, tp)
	return
}

// AdminReportRecover recover report by admin.
func (s *Service) AdminReportRecover(c context.Context, adid, oid, rpID int64, tp, audit int8, remark string) (err error) {
	// check subject
	now := time.Now()
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	s.dao.Databus.AdminReportRecover(c, adid, oid, rpID, remark, now.Unix(), tp, audit)
	return
}
