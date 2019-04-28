package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-common/app/job/main/reply/model/reply"
	model "go-common/app/job/main/reply/model/reply"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_eventReportAdd     = "report_add"
	_eventReportDel     = "report_del"
	_eventReportIgnore  = "report_ignore"
	_eventReportRecover = "report_recover"
)

func (s *Service) actionAdmin(c context.Context, msg *consumerMsg) {
	var d struct {
		Op      string     `json:"op"`
		Adid    int64      `json:"adid"`
		AdName  string     `json:"adname"`
		Oid     int64      `json:"oid"`
		RpID    int64      `json:"rpid"`
		Mid     int64      `json:"mid"`
		Tp      int8       `json:"tp"`
		Action  uint32     `json:"action"`
		Moral   int        `json:"moral"`
		Notify  bool       `json:"notify"`
		Remark  string     `json:"remark"`
		MTime   xtime.Time `json:"mtime"`
		Ftime   int64      `json:"ftime"`
		Audit   int8       `json:"audit"`
		Reason  int8       `json:"reason"`
		Content string     `json:"content"`
		FReason int8       `json:"freason"`
		Assist  bool       `json:"assist"`
		State   int8       `json:"state"`
	}
	if err := json.Unmarshal([]byte(msg.Data), &d); err != nil {
		log.Error("json.Unmarshal() error(%v)", err)
		return
	}
	if d.Oid <= 0 || d.RpID <= 0 {
		log.Error("The structure of msg.Data(%s) was wrong", msg.Data)
		return
	}
	rp, err := s.getReply(c, d.Oid, d.RpID)
	if err != nil {
		log.Error("s.getReply failed , oid(%d), RpID(%d) err(%v)", d.Oid, d.RpID, err)
		return
	}
	if rp == nil {
		log.Error("getReply nil oid(%d) RpID(%d)", d.Oid, d.RpID)
		return
	}
	switch {
	case d.Op == "del":
		s.adminDel(c, rp, d.Adid, d.Mid, d.Ftime, d.Moral, d.AdName, d.Remark, d.MTime, d.Notify, d.Reason, d.FReason)
	case d.Op == "del_rpt":
		s.reportDel(c, rp, d.Adid, d.Mid, d.Ftime, d.Moral, d.AdName, d.Remark, d.MTime, d.Notify, d.Audit, d.Reason, d.Content, d.FReason)
	case d.Op == "del_up":
		s.userDel(c, rp, d.Mid, d.MTime, d.Remark, d.Assist)
	case d.Op == "re", d.Op == "pass":
		s.passReply(c, rp, d.MTime, d.Adid, d.Remark, d.Op)
	case d.Op == "edit":
		s.addSearchUp(c, rp.State, rp, nil)
	case d.Op == "ignore":
		s.reportIgnore(c, rp, d.Audit, d.MTime, d.Adid, d.Remark)
	case d.Op == "transfer":
		s.reportTransfer(c, rp, d.Audit, d.MTime, d.Adid, d.Remark)
	case d.Op == "stateset":
		s.reportStateSet(c, rp, d.State, d.MTime, d.Adid, d.Remark)
	case d.Op == "top_add":
		err := s.topAdd(c, rp, d.MTime, d.Action, model.SubAttrAdminTop)
		if err != nil {
			log.Error("s.topAdd(oid:%d,tp:%d err(%v))", rp.Oid, rp.Type, err)
			return
		}
		//	s.dao.Redis.AddTopOid(c, rp.Oid, rp.Type)
		s.adminLog(c, rp, d.Adid, model.AdminIsReport, model.AdminOperSubTop, "管理员置顶评论", d.Remark)
		s.addSearchUp(c, rp.State, rp, nil)
	case d.Op == "rpt_re":
		s.reportRecover(c, rp, d.Audit, d.MTime, d.Adid, d.Remark)
	}
}

func (s *Service) reportStateSet(c context.Context, rp *model.Reply, state int8, mtime xtime.Time, adid int64, remark string) (err error) {
	var (
		op       int8
		result   string
		oldState = rp.State
	)
	rpt, err := s.dao.Report.Get(c, rp.Oid, rp.RpID)
	if err != nil || rpt == nil {
		log.Error("dao.Report.GetReport(%d, %d) met error (%v)", rp.Oid, rp.RpID, err)
		return
	}

	rpt.State = state
	rpt.MTime = mtime

	op = model.AdminOperRptTransferArbitration
	if adid == 0 {
		result = "系统自动移交至风纪委"
	} else {
		result = "管理员移交至风纪委"
	}
	if _, err = s.dao.Report.Update(c, rpt); err != nil {
		log.Error("dao.Report.Update(%v) error(%v)", rpt, err)
		return
	}
	// admin log
	s.adminLog(c, rp, adid, model.AdminIsReport, op, result, remark)
	s.addSearchUp(c, oldState, rp, rpt)
	return
}

func (s *Service) reportTransfer(c context.Context, rp *model.Reply, audit int8, mtime xtime.Time, adid int64, remark string) (err error) {
	var (
		op       int8
		result   string
		oldState = rp.State
	)
	rpt, err := s.dao.Report.Get(c, rp.Oid, rp.RpID)
	if err != nil || rpt == nil {
		log.Error("dao.Report.GetReport(%d, %d) met error (%v)", rp.Oid, rp.RpID, err)
		return
	}
	if rpt.IsTransferred() {
		log.Error("report(%v) has been transfferd before!", *rpt)
		return
	}
	if audit == model.AuditTypeOne {
		rpt.State = model.ReportStateNew
		op = model.AdminOperRptTransfer1
		result = "二审转一审"
	} else if audit == model.AuditTypeTwo {
		rpt.State = model.ReportStateNewTwo
		op = model.AdminOperRptTransfer2
		result = "一审转二审"
	}
	rpt.MTime = mtime
	rpt.SetTransferred()
	if _, err = s.dao.Report.Update(c, rpt); err != nil {
		log.Error("dao.Report.Update(%v) error(%v)", rpt, err)
		return
	}
	// admin log
	s.adminLog(c, rp, adid, model.AdminIsReport, op, result, remark)
	s.addSearchUp(c, oldState, rp, rpt)

	return
}

func (s *Service) reportIgnore(c context.Context, rp *model.Reply, audit int8, mtime xtime.Time, adid int64, remark string) (err error) {
	var (
		op       int8
		result   string
		oldState = rp.State
	)
	rpt, err := s.dao.Report.Get(c, rp.Oid, rp.RpID)
	if err != nil || rpt == nil {
		log.Error("dao.Report.GetReport(%d, %d) met error (%v)", rp.Oid, rp.RpID, err)
		return
	}
	if audit == model.AuditTypeOne {
		rpt.State = model.ReportStateIgnoreOne
		op = model.AdminOperRptIgnore1
		result = "一审忽略"
	} else if audit == model.AuditTypeTwo {
		rpt.State = model.ReportStateIgnoreTwo
		op = model.AdminOperRptIgnore2
		result = "二审忽略"
	} else {
		if rpt.State == model.ReportStateNew {
			rpt.State = model.ReportStateIgnoreOne
		} else {
			rpt.State = model.ReportStateIgnoreTwo
		}
		op = model.AdminOperIgnoreReport
		result = "已忽略举报"
	}
	rpt.MTime = mtime
	if _, err = s.dao.Report.Update(c, rpt); err != nil {
		log.Error("dao.Report.Update(%v) error(%v)", rpt, err)
		return
	}
	// admin log
	s.adminLog(c, rp, adid, model.AdminIsReport, op, result, remark)
	s.addSearchUp(c, oldState, rp, rpt)

	sub, err := s.getSubject(c, rp.Oid, rp.Type)
	if err != nil || sub == nil {
		log.Error("s.getSubject(%v,%v) err(%v) or sub nil", sub.Oid, sub.Type, err)
		return
	}
	if err = s.dao.PubEvent(c, _eventReportIgnore, rpt.Mid, sub, rp, rpt); err != nil {
		return
	}
	return
}

func (s *Service) reportRecover(c context.Context, rp *model.Reply, audit int8, mtime xtime.Time, adid int64, remark string) (err error) {
	var (
		ok       bool
		op       int8
		result   string
		rootRp   *model.Reply
		oldState = rp.State
	)
	isRoot := rp.Root == 0 && rp.Parent == 0
	rpt, _ := s.dao.Report.Get(c, rp.Oid, rp.RpID)
	if err != nil || rpt == nil {
		log.Error("dao.Report.GetReport(%d, %d) met error (%v)", rp.Oid, rp.RpID, err)
		return
	}
	sub, err := s.getSubject(c, rp.Oid, rp.Type)
	if err != nil || sub == nil {
		log.Error("s.getSubject(%v,%v) err(%v) or sub nil", sub.Oid, sub.Type, err)
		return
	}
	if audit == model.AuditTypeOne {
		// 一审移除，恢复到待一审忽略，评论正常
		rpt.State = model.ReportStateIgnoreOne
		op = model.AdminOperRptRecover1
		result = "一审恢复评论"
	} else if audit == model.AuditTypeTwo {
		// 二审移除，恢复到二审忽略，评论正常
		rpt.State = model.ReportStateIgnoreTwo
		op = model.AdminOperRptRecover2
		result = "二审恢复评论"
	} else {
		log.Error("reportRecover unsupport audit: %d", audit)
		return
	}
	rpt.MTime = mtime
	if _, err = s.dao.Report.Update(c, rpt); err != nil {
		log.Error("dao.Report.Update(%v) error(%v) or row==0", rpt, err)
		return
	}
	// 只恢复管理员删除的评论，不恢复用户删除的
	if rp.State == model.ReplyStateAdminDel {
		rp.MTime = mtime
		rp.State = model.ReplyStateNormal
		if err = s.tranRecover(c, rp, model.ReplyStateNormal, isRoot); err != nil {
			log.Error("Transaction recover reply failed err(%v)", err)
			return
		}
	}
	// add cache
	if isRoot {
		if err = s.dao.Mc.AddReply(c, rp); err != nil {
			log.Error("s.dao.Mc.AddReply failed , RpID(%d),  err(%v)", rp.RpID, err)
		}
		if err = s.dao.Redis.AddIndex(c, rp.Oid, rp.Type, rpt, rp, true); err != nil {
			log.Error("s.dao.Redis.AddIndex(%d, %d) error(%v)", rp.Oid, rp.Type, err)
		}
	} else {
		if ok, err = s.dao.Redis.ExpireNewChildIndex(c, rp.Root); err == nil && ok {
			if err = s.dao.Redis.AddNewChildIndex(c, rp.Root, rp); err != nil {
				log.Error("s.dao.Redis.AddFloorIndexByRoot failed , rproot(%d),  err(%v)", rp.Root, err)
			}
		}
		if rootRp, err = s.getReplyCache(c, rp.Oid, rp.Root); err != nil {
			log.Error("s.getReply failed , oid(%d), root(%d) err(%v)", rp.Oid, rp.Root, err)
		} else if rootRp != nil {
			rootRp.RCount++
			if err = s.dao.Mc.AddReply(c, rootRp); err != nil {
				log.Error("s.dao.Mc.AddReply failed , RpID(%d) err(%v)", rootRp.RpID, err)
			}
		}
	}
	// log
	s.adminLog(c, rp, adid, model.AdminIsReport, op, result, remark)
	// notify
	s.addSearchUp(c, oldState, rp, rpt)
	s.upAcount(c, sub.Oid, sub.Type, sub.ACount, rp.CTime.Time())

	if err = s.dao.PubEvent(c, _eventReportRecover, rpt.Mid, sub, rp, rpt); err != nil {
		return
	}
	return
}

// topAdd add top reply
func (s *Service) topAdd(c context.Context, rp *model.Reply, ts xtime.Time, act uint32, tp uint32) (err error) {
	sub, err := s.getSubject(c, rp.Oid, rp.Type)
	if err != nil || sub == nil {
		log.Error("s.getsubject(%v,%v) err(%v) or sub nil", rp.Oid, rp.Type, err)
		return
	}
	if act == 1 && sub.AttrVal(tp) == 1 {
		log.Error("Repeat to add top reply(%d,%d,%d,%d) ", rp.RpID, rp.Oid, tp, sub.Attr)
		return
	}
	err = sub.TopSet(rp.RpID, tp, act)
	if err != nil {
		log.Error("sub.TopSet(%d,%d,%d) failed!err:=%v ", rp.RpID, rp.Oid, tp, err)
		return
	}
	sub.AttrSet(act, tp)
	rp.AttrSet(act, tp)
	tx, err := s.beginTran(c)
	if err != nil {
		log.Error("s.beginTran() err(%v)", err)
		return
	}
	var rows int64
	//rp.State = model.ReplyStateTop
	if rows, err = s.dao.Reply.TxUpAttr(tx, rp.Oid, rp.RpID, rp.Attr, ts.Time()); err != nil || rows == 0 {
		tx.Rollback()
		log.Error("dao.Reply.UpState(%v, %d) error(%v) or row==0", rp, rp.State, err)
		return
	}
	if rows, err = s.dao.Subject.TxUpMeta(tx, sub.Oid, sub.Type, sub.Meta, ts.Time()); err != nil || rows == 0 {
		tx.Rollback()
		log.Error("dao.TxUpMeta(oid:%d,tp:%d) err(%v) rows(%d)", sub.Oid, sub.Type, err, rows)
		return
	}
	if rows, err = s.dao.Subject.TxUpAttr(tx, sub.Oid, sub.Type, sub.Attr, ts.Time()); err != nil || rows == 0 {
		tx.Rollback()
		log.Error("dao.Upattr(oid:%d,tp:%d) err(%v) rows(%d)", sub.Oid, sub.Type, err, rows)
		return
	}
	tx.Commit()
	s.dao.Mc.AddSubject(c, sub)
	if act == 1 {
		s.dao.Redis.DelIndexBySortType(c, rp, reply.SortByCount)
		s.dao.Redis.DelIndexBySortType(c, rp, reply.SortByLike)
	} else if rp.IsNormal() {
		if ok, err := s.dao.Redis.ExpireIndex(c, sub.Oid, sub.Type, model.SortByCount); err == nil && ok {
			if err = s.dao.Redis.AddCountIndex(c, sub.Oid, sub.Type, rp); err != nil {
				log.Error("s.dao.Redis.AddCountIndex failed , oid(%d) type(%d) err(%v)", sub.Oid, sub.Type, err)
			}
		}
		if ok, err := s.dao.Redis.ExpireIndex(c, sub.Oid, sub.Type, model.SortByLike); err == nil && ok {
			rpts := make(map[int64]*reply.Report, 1)
			if rpt, _ := s.dao.Report.Get(c, rp.Oid, rp.RpID); rpt != nil {
				rpts[rp.RpID] = rpt
			}
			if err = s.dao.Redis.AddLikeIndex(c, sub.Oid, sub.Type, rpts, rp); err != nil {
				log.Error("s.dao.Redis.AddLikeIndex failed , oid(%d) type(%d) err(%v)", sub.Oid, sub.Type, err)
			}
		}
	}
	s.dao.Mc.AddReply(c, rp)
	s.dao.Mc.AddTop(c, rp)
	if act == 1 {
		s.dao.PubEvent(c, "top", 0, sub, rp, nil)
	} else if act == 0 {
		s.dao.PubEvent(c, "untop", 0, sub, rp, nil)
	}
	// 折叠评论被置顶自动取消折叠
	if rp.IsFolded() && act == 1 {
		rp.State = model.ReplyStateNormal
		s.marker.Do(c, func(ctx context.Context) {
			if _, err := s.dao.Reply.UpState(ctx, rp.Oid, rp.RpID, rp.State, time.Now()); err == nil {
				if ok, err := s.dao.Redis.ExpireIndex(c, rp.Oid, rp.Type, model.SortByFloor); err == nil && ok {
					s.dao.Redis.AddFloorIndex(c, rp.Oid, rp.Type, rp)
				}
				s.handleFolded(ctx, rp)
			}
		})
	}
	return
}

func (s *Service) userDel(c context.Context, rp *model.Reply, mid int64, mtime xtime.Time, remark string, assist bool) (err error) {
	var (
		state    int8
		sub      *model.Subject
		oldState = rp.State
		isRoot   = rp.Root == 0 && rp.Parent == 0
		isFolded = rp.IsFolded()
	)
	if rp.IsDeleted() {
		s.addSearchUp(c, oldState, rp, nil)
		return
	}
	if sub, err = s.dao.Subject.Get(c, rp.Oid, rp.Type); err != nil || sub == nil {
		log.Error("s.dao.Subject.Get(%d,%d) error(%v)", rp.Oid, rp.Type, err)
		return
	}
	if rp.Mid == mid {
		state = model.ReplyStateUserDel
	} else if assist {
		state = model.ReplyStateAssistDel
	} else if mid > 0 {
		state = model.ReplyStateUpDel
	} else {
		state = model.ReplyStateAdminDel
	}
	rp.MTime = mtime
	if err = s.tranDel(c, rp, state, sub, isRoot); err != nil {
		log.Error("reportDel tranDel(%d, %d) error(%v)", rp.Oid, rp.RpID, err)
		return
	}
	if isFolded {
		s.marker.Do(c, func(ctx context.Context) {
			s.handleFolded(ctx, rp)
		})
	}
	if err = s.clearReplyCache(c, rp); err != nil {
		log.Error("reportDel clearReplyCache(%d, %d) error(%v)", rp.Oid, rp.RpID, err)
	}
	rp.State = state
	if rp.Mid == mid {
		s.adminLog(c, rp, mid, model.AdminIsNotReport, model.AdminOperDeleteUser, "由本人删除", remark)
		s.searchDao.DelReply(c, rp.RpID, sub.Oid, mid, state)
	} else if assist {
		s.adminLog(c, rp, mid, model.AdminIsNotReport, model.AdminOperDeleteAssist, "由UP主协管员删除", remark)
		s.addAssistLog(c, sub.Mid, mid, sub.Oid, 1, 1, strconv.FormatInt(rp.RpID, 10), rp.Content.Message)
	} else if mid > 0 {
		s.adminLog(c, rp, mid, model.AdminIsNotReport, model.AdminOperDeleteUp, "由up主删除", remark)
	} else {
		s.adminLog(c, rp, 0, model.AdminIsNotReport, model.AdminOperDelete, "由系统删除", remark)
	}
	s.dao.PubEvent(c, "reply_del", rp.Mid, sub, rp, nil)
	return
}

func (s *Service) adminDel(c context.Context, rp *model.Reply, adid, mid, ftime int64, moral int, adName, remark string, mtime xtime.Time, notify bool, reason, freason int8) (err error) {
	var (
		sub      *model.Subject
		report   *model.Report
		oldState = rp.State
		isRoot   = rp.Root == 0 && rp.Parent == 0
		isFolded = rp.IsFolded()
	)
	if rp.IsDeleted() {
		s.addSearchUp(c, oldState, rp, nil)
		return
	}
	if sub, err = s.dao.Subject.Get(c, rp.Oid, rp.Type); err != nil || sub == nil {
		log.Error("s.dao.Subject.Get(%d,%d) error(%v)", rp.Oid, rp.Type, err)
		return
	}
	rp.MTime = mtime
	if err = s.tranDel(c, rp, model.ReplyStateAdminDel, sub, isRoot); err != nil {
		log.Error("reportDel tranDel(%d, %d) error(%v)", rp.Oid, rp.RpID, err)
		return
	}
	if isFolded {
		s.marker.Do(c, func(ctx context.Context) {
			s.handleFolded(ctx, rp)
		})
	}
	if err = s.clearReplyCache(c, rp); err != nil {
		log.Error("reportDel clearReplyCache(%d, %d) error(%v)", rp.Oid, rp.RpID, err)
	}
	if report, err = s.dao.Report.Get(c, rp.Oid, rp.RpID); err != nil {
		log.Error("reportDel getReport(%d,%d) error(%v)", rp.Oid, rp.RpID, err)
		return
	}
	if report != nil {
		if report.State == model.ReportStateNew || report.State == model.ReportStateNewTwo {
			if report.State == model.ReportStateNew {
				report.State = model.ReportStateDeleteOne
			} else if report.State == model.ReportStateNewTwo {
				report.State = model.ReportStateDeleteTwo
			}
			report.MTime = mtime
			if _, err = s.dao.Report.Update(c, report); err != nil {
				log.Error("reportDel updateReport(%d, %d) error(%v)", report.Oid, report.ID, err)
				return
			}
			s.addSearchUp(c, oldState, rp, report)
			s.dao.PubEvent(c, _eventReportDel, 0, sub, rp, report)
		}
	}
	rp.State = model.ReplyStateAdminDel
	// add moral and notify
	s.moralAndNotify(c, rp, moral, notify, mid, adid, adName, remark, reason, freason, ftime, false)
	// forbidden tip
	forbidDay := strconv.FormatInt(ftime, 10) + "天"
	if ftime == -1 {
		forbidDay = "永久"
	}
	s.adminLog(c, rp, adid, model.AdminIsNotReport, model.AdminOperDelete, fmt.Sprintf("已删除并封禁%s/扣除%d节操", forbidDay, moral), remark)
	s.addSearchUp(c, oldState, rp, nil)
	if report == nil {
		s.dao.PubEvent(c, "reply_del", 0, sub, rp, nil)
	}
	return
}

func (s *Service) reportDel(c context.Context, rp *model.Reply, adid, mid, ftime int64, moral int, adName, remark string, mtime xtime.Time, notify bool, audit int8, reason int8, content string, freason int8) (err error) {
	var (
		op       int8
		isPunish bool
		sub      *model.Subject
		report   *model.Report
		oldState = rp.State
	)
	if sub, err = s.dao.Subject.Get(c, rp.Oid, rp.Type); err != nil || sub == nil {
		log.Error("s.dao.Subject.Get(%d,%d) error(%v)", rp.Oid, rp.Type, err)
		return
	}
	if report, err = s.dao.Report.Get(c, rp.Oid, rp.RpID); err != nil || report == nil {
		log.Error("reportDel getReport(%d,%d) error(%v)", rp.Oid, rp.RpID, err)
		return
	}
	report.MTime = mtime
	// 一审、二审操作
	switch audit {
	case model.AuditTypeOne:
		report.State = model.ReportStateDeleteOne
		op = model.AdminOperRptDel1
	case model.AuditTypeTwo:
		report.State = model.ReportStateDeleteTwo
		op = model.AdminOperRptDel2
	default:
		report.State = model.ReportStateDelete
		op = model.AdminOperDeleteByReport
	}
	report.Reason = reason
	report.Content = content
	if _, err = s.dao.Report.Update(c, report); err != nil {
		log.Error("reportDel updateReport(%d, %d) error(%v)", report.Oid, report.ID, err)
		return
	}
	if !rp.IsDeleted() {
		isRoot := rp.Root == 0 && rp.Parent == 0
		rp.MTime = mtime
		if err = s.tranDel(c, rp, model.ReplyStateAdminDel, sub, isRoot); err != nil {
			log.Error("reportDel tranDel(%d, %d) error(%v)", rp.Oid, rp.RpID, err)
			return
		}
		if err = s.clearReplyCache(c, rp); err != nil {
			log.Error("reportDel clearReplyCache(%d, %d) error(%v)", rp.Oid, rp.RpID, err)
		}
		rp.State = model.ReplyStateAdminDel
	} else {
		isPunish = true
	}
	// add moral and notify
	s.moralAndNotify(c, rp, moral, notify, mid, adid, adName, remark, reason, freason, ftime, isPunish)
	// forbidden tip
	forbidDay := strconv.FormatInt(ftime, 10) + "天"
	if ftime == -1 {
		forbidDay = "永久"
	}
	s.adminLog(c, rp, adid, model.AdminIsReport, op, fmt.Sprintf("已通过举报删除并封禁%s/扣除%d节操", forbidDay, moral), remark)
	s.addSearchUp(c, oldState, rp, report)

	if err = s.dao.PubEvent(c, _eventReportDel, report.Mid, sub, rp, report); err != nil {
		return
	}
	return
}

func (s *Service) clearReplyCache(c context.Context, rp *model.Reply) (err error) {
	var (
		sub    *model.Subject
		rootRp *model.Reply
		isRoot = rp.Root == 0 && rp.Parent == 0
	)
	if !isRoot && rp.IsNormal() {
		// update root cache for count.
		if rootRp, err = s.getReplyCache(c, rp.Oid, rp.Root); err != nil {
			return
		}
		if rootRp != nil {
			rootRp.RCount--
			if err = s.addReplyCache(c, rootRp); err != nil {
				log.Error("s.dao.Mc.addReplyCache(%d,%d,%d) error(%v)", rp.Oid, rp.Type, rp.RpID, err)
			}
		}
	}
	if rp.IsAdminTop() {
		s.dao.Mc.DeleteTop(c, rp, model.ReplyAttrAdminTop)
	}
	if rp.IsUpTop() {
		s.dao.Mc.DeleteTop(c, rp, model.ReplyAttrUpperTop)
	}
	if err = s.dao.Mc.DeleteReply(c, rp.RpID); err != nil {
		log.Error("s.dao.Mc.DeleteReply failed , RpID(%d),  err(%v)", rp.RpID, err)
	}
	if err = s.dao.Redis.DelIndex(c, rp); err != nil {
		log.Error("s.dao.Redis.DelIndex failed , RpID(%d),  err(%v)", rp.RpID, err)
	}
	if rp.State == model.ReplyStateAudit {
		s.eraseAuditIndex(c, rp)
	}
	// update reply count
	sub, err = s.dao.Subject.Get(c, rp.Oid, rp.Type)
	if err != nil || sub == nil {
		log.Error("s.dao.Subject.Get(%d,%d) error(%v)", rp.Oid, rp.Type, err)
		return
	}
	if err = s.dao.Mc.AddSubject(c, sub); err != nil {
		log.Error("s.dao.Mc.AddSubject failed , oid(%d),  err(%v)", sub.Oid, err)
	}
	s.upAcount(c, sub.Oid, sub.Type, sub.ACount, rp.CTime.Time())
	return
}

func (s *Service) passReply(c context.Context, rp *model.Reply, mtime xtime.Time, adid int64, remark, op string) {
	var (
		err      error
		ok       bool
		rootRp   *model.Reply
		oldState = rp.State
		isRoot   = rp.Root == 0 && rp.Parent == 0
	)
	sub, err := s.dao.Subject.Get(c, rp.Oid, rp.Type)
	if err != nil || sub == nil {
		log.Error("s.dao.Subject.Get(%d,%d) error(%v)", rp.Oid, rp.Type, err)
		return
	}
	if rp.State <= model.ReplyStateHidden {
		s.addSearchUp(c, oldState, rp, nil)
		return
	}
	rp.MTime = mtime
	rp.State = model.ReplyStateNormal
	if op == "re" || oldState == model.ReplyStateAudit {
		if err = s.tranRecover(c, rp, model.ReplyStateNormal, isRoot); err != nil {
			log.Error("Transaction recover reply failed err(%v)", err)
			return
		}
	} else if op == "pass" {
		var row int64
		var tx *xsql.Tx
		tx, err = s.dao.BeginTran(c)
		if err != nil {
			return
		}
		if row, err = s.dao.Reply.TxUpState(tx, rp.Oid, rp.RpID, model.ReplyStateNormal, mtime.Time()); err != nil || row == 0 {
			tx.Rollback()
			log.Error("dao.Reply.TxUpState(%v, %d) error(%v) or row==0", rp, model.ReplyStateNormal, err)
			return
		}
		if rp.State == model.ReplyStateAudit || rp.State == model.ReplyStateMonitor {
			if _, err = s.dao.Subject.TxDecrMCount(tx, rp.Oid, rp.Type, mtime.Time()); err != nil {
				tx.Rollback()
				log.Error("dao.Reply.TxDecrMCount(%v) error(%v)", rp, err)
				return
			}
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit error(%v)", err)
			return
		}
	}
	if isRoot {
		if ok, err = s.dao.Redis.ExpireIndex(c, rp.Oid, rp.Type, model.SortByFloor); err == nil && ok {
			var min int
			min, err = s.dao.Redis.MinScore(c, rp.Oid, rp.Type, reply.SortByFloor)
			if err != nil {
				log.Error("s.dao.Redis.AddFloorIndex failed , oid(%d) type(%d) err(%v)", rp.Oid, rp.Type, err)
			} else if rp.Floor > min {
				if err = s.dao.Redis.AddFloorIndex(c, rp.Oid, rp.Type, rp); err != nil {
					log.Error("s.dao.Redis.AddFloorIndex failed , RpID(%d),  err(%v)", rp.RpID, err)
				}
			}
		}
		if ok, err = s.dao.Redis.ExpireIndex(c, rp.Oid, rp.Type, model.SortByCount); err == nil && ok {
			if err = s.dao.Redis.AddCountIndex(c, rp.Oid, rp.Type, rp); err != nil {
				log.Error("s.dao.Redis.AddCountIndex failed , RpID(%d),  err(%v)", rp.RpID, err)
			}
		}
		if ok, err = s.dao.Redis.ExpireIndex(c, rp.Oid, rp.Type, model.SortByLike); err == nil && ok {
			rpts := make(map[int64]*reply.Report, 1)
			if rpt, _ := s.dao.Report.Get(c, rp.Oid, rp.RpID); rpt != nil {
				rpts[rp.RpID] = rpt
			}
			if err = s.dao.Redis.AddLikeIndex(c, rp.Oid, rp.Type, rpts, rp); err != nil {
				log.Error("s.dao.Redis.AddLikeIndex failed , RpID(%d),  err(%v)", rp.RpID, err)
			}
		}
	} else {
		if ok, err = s.dao.Redis.ExpireDialogIndex(c, rp.Dialog); err == nil && ok {
			if err = s.dao.Redis.AddDialogIndex(c, rp.Dialog, []*model.Reply{rp}); err != nil {
				log.Error("s.dao.Redis.AddDialogINdex Error (%v)", err)
			}
		}
		if ok, err = s.dao.Redis.ExpireNewChildIndex(c, rp.Root); err == nil && ok {
			if err = s.dao.Redis.AddNewChildIndex(c, rp.Root, rp); err != nil {
				log.Error("s.dao.Redis.AddFloorIndexByRoot failed , rproot(%d),  err(%v)", rp.Root, err)
			}
		}
		if op == "re" || oldState == model.ReplyStateAudit {
			if rootRp, err = s.getReplyCache(c, rp.Oid, rp.Root); err != nil {
				log.Error("s.getReply failed , oid(%d), root(%d) err(%v)", rp.Oid, rp.Root, err)
			} else if rootRp != nil {
				rootRp.RCount++
				if err = s.addReplyCache(c, rootRp); err != nil {
					log.Error("s.addReplyCache oid(%d), rpid(%d) err(%v)", rootRp.Oid, rootRp.RpID, err)
				}
			}
		}
	}
	if err = s.addReplyCache(c, rp); err != nil {
		log.Error("s.addReplyCache oid(%d), rpid(%d) err(%v)", rp.Oid, rp.RpID, err)
	}
	if oldState == model.ReplyStateAudit {
		s.dao.Redis.DelAuditIndexs(c, rp)
	}
	// update reply count
	s.upAcount(c, sub.Oid, sub.Type, sub.ACount, rp.CTime.Time())
	if err = s.dao.Mc.AddSubject(c, sub); err != nil {
		log.Error("s.dao.Mc.AddSubject failed , oid(%d) err(%v)", sub.Oid, err)
	}
	// admin log
	if op == "re" {
		s.adminLog(c, rp, adid, model.AdminIsNotReport, model.AdminOperRecover, "已恢复评论", remark)

	} else if op == "pass" {
		s.adminLog(c, rp, adid, model.AdminIsNotReport, model.AdminOperPass, "已通过评论", remark)
	}
	s.addSearchUp(c, oldState, rp, nil)
}

func (s *Service) tranRecover(c context.Context, rp *model.Reply, state int8, isRoot bool) error {
	var (
		err       error
		rootReply *model.Reply
		count     int
		rows      int64
		tx        *xsql.Tx
	)
	if tx, err = s.beginTran(c); err != nil {
		return err
	}
	mtime := rp.MTime
	if rp, err = s.dao.Reply.GetForUpdate(tx, rp.Oid, rp.RpID); err != nil || rp == nil {
		tx.Rollback()
		return fmt.Errorf("s.dao.Reply.GetForUpdate(%d,%d) error(%v) or is nil)", rp.Oid, rp.RpID, err)
	}
	if mtime != 0 {
		rp.MTime = mtime
	}
	if rp.IsNormal() {
		tx.Rollback()
		return fmt.Errorf("reply(%d,%d) already is normal", rp.Oid, rp.RpID)
	}
	rows, err = s.dao.Reply.TxUpState(tx, rp.Oid, rp.RpID, state, rp.MTime.Time())
	if err != nil || rows == 0 {
		tx.Rollback()
		return fmt.Errorf("TxUpState error(%v) or rows(%d)", err, rows)
	}
	if isRoot {
		count = rp.RCount + 1
	} else {
		rootReply, err = s.dao.Reply.Get(c, rp.Oid, rp.Root)
		if err != nil {
			tx.Rollback()
			return err
		}
		count = 1
	}
	if isRoot {
		rows, err = s.dao.Subject.TxIncrRCount(tx, rp.Oid, rp.Type, rp.MTime.Time())
	} else {
		rows, err = s.dao.Reply.TxIncrRCount(tx, rp.Oid, rp.Root, rp.MTime.Time())
	}
	if err != nil || rows == 0 {
		tx.Rollback()
		return fmt.Errorf("TxIncrRCount error(%v) or rows(%d)", err, rows)
	}
	if isRoot || rootReply != nil && rootReply.IsNormal() {
		rows, err = s.dao.Subject.TxIncrACount(tx, rp.Oid, rp.Type, count, rp.MTime.Time())
		if err != nil || rows == 0 {
			tx.Rollback()
			return fmt.Errorf("TxIncrACount error(%v) or rows(%d)", err, rows)
		}
	}
	return tx.Commit()
}

func (s *Service) tranDel(c context.Context, rp *model.Reply, state int8, sub *model.Subject, isRoot bool) error {
	var (
		count     int
		rows      int64
		rootReply *model.Reply
		err       error
		tx        *xsql.Tx
	)
	if tx, err = s.beginTran(c); err != nil {
		return err
	}
	mtime := rp.MTime
	if rp, err = s.dao.Reply.GetForUpdate(tx, rp.Oid, rp.RpID); err != nil || rp == nil {
		tx.Rollback()
		return fmt.Errorf("s.dao.Reply.GetForUpdate(%d,%d) error(%v) or is nil)", rp.Oid, rp.RpID, err)
	}
	if mtime != 0 {
		rp.MTime = mtime
	}
	if rp.IsDeleted() || rp.AttrVal(reply.ReplyAttrAdminTop) == 1 {
		tx.Rollback()
		return fmt.Errorf("reply(%d,%d) already deleted", rp.Oid, rp.RpID)
	}
	rows, err = s.dao.Reply.TxUpState(tx, rp.Oid, rp.RpID, state, rp.MTime.Time())
	if err != nil || rows == 0 {
		tx.Rollback()
		return fmt.Errorf("error(%v) or rows(%d)", err, rows)
	}
	if rp.IsNormal() {
		if isRoot {
			count = rp.RCount + 1
			rows, err = s.dao.Subject.TxDecrACount(tx, rp.Oid, rp.Type, count, rp.MTime.Time())
			if err != nil || rows == 0 {
				tx.Rollback()
				return fmt.Errorf("error(%v) or rows(%d)", err, rows)
			}
			rows, err = s.dao.Subject.TxDecrCount(tx, rp.Oid, rp.Type, rp.MTime.Time())
			if err != nil || rows == 0 {
				tx.Rollback()
				return fmt.Errorf("SubjectTxDecrCount error(%v) or rows(%d)", err, rows)
			}
		} else {
			if rootReply, err = s.dao.Reply.GetForUpdate(tx, rp.Oid, rp.Root); err != nil {
				tx.Rollback()
				return err
			}
			if rootReply != nil {
				if rootReply.IsNormal() {
					rows, err = s.dao.Subject.TxDecrACount(tx, rp.Oid, rp.Type, 1, rp.MTime.Time())
					if err != nil || rows == 0 {
						tx.Rollback()
						return fmt.Errorf("error(%v) or rows(%d)", err, rows)
					}
				}
				rows, err = s.dao.Reply.TxDecrCount(tx, rp.Oid, rp.Root, rp.MTime.Time())
				if err != nil || rows == 0 {
					tx.Rollback()
					return fmt.Errorf("ReplyTxDecrCount error(%v) or rows(%d)", err, rows)
				}
			}
		}
	}
	if rp.AttrVal(model.ReplyAttrUpperTop) == 1 {
		rp.AttrSet(0, model.ReplyAttrUpperTop)
		sub.AttrSet(0, model.SubAttrUpperTop)
		sub.TopSet(0, model.SubAttrUpperTop, 0)
		if rows, err = s.dao.Subject.TxUpMeta(tx, sub.Oid, sub.Type, sub.Meta, rp.MTime.Time()); err != nil || rows == 0 {
			tx.Rollback()
			log.Error("dao.TxUpMeta(oid:%d,tp:%d) err(%v) rows(%d)", sub.Oid, sub.Type, err, rows)
			return fmt.Errorf("dao.TxUpMeta(oid:%d,tp:%d) err(%v) rows(%d)", sub.Oid, sub.Type, err, rows)
		}
		if rows, err = s.dao.Subject.TxUpAttr(tx, sub.Oid, sub.Type, sub.Attr, rp.MTime.Time()); err != nil || rows == 0 {
			tx.Rollback()
			return fmt.Errorf("dao.Upattr(oid:%d,tp:%d) err(%v) rows(%d)", sub.Oid, sub.Type, err, rows)
		}
		//rp.State = model.ReplyStateTop
		if rows, err = s.dao.Reply.TxUpAttr(tx, rp.Oid, rp.RpID, rp.Attr, rp.MTime.Time()); err != nil || rows == 0 {
			tx.Rollback()
			return fmt.Errorf("dao.Reply.UpState(%v, %d) error(%v) or rows(%d)", rp, rp.State, err, rows)
		}
	}
	if rp.State == model.ReplyStateMonitor || rp.State == model.ReplyStateAudit {
		if _, err = s.dao.Subject.TxDecrMCount(tx, rp.Oid, rp.Type, rp.MTime.Time()); err != nil {
			tx.Rollback()
			log.Error("dao.Reply.TxDecrMCount(%v) error(%v)", rp, err)
			return fmt.Errorf("dao.Reply.TxDecrMCount error(%v)", err)
		}
	}
	return tx.Commit()
}

func (s *Service) eraseAuditIndex(c context.Context, rp *model.Reply) (err error) {
	var rs []*model.Reply
	if rp.Root == 0 && rp.Parent == 0 {
		if rs, err = s.dao.Reply.GetsByRoot(c, rp.Oid, rp.RpID, rp.Type, model.ReplyStateAudit); err != nil {
			return
		}
	}
	if err = s.dao.Redis.DelAuditIndexs(c, append(rs, rp)...); err != nil {
		return
	}
	return
}

func (s *Service) addReplyCache(c context.Context, rp *model.Reply) (err error) {
	var isRoot = rp.Root == 0 && rp.Parent == 0
	if err = s.dao.Mc.AddReply(c, rp); err != nil {
		log.Error("s.dao.Mc.AddReply(%d,%d,%d) error(%v)", rp.Oid, rp.RpID, rp.Type, err)
	}
	if isRoot && rp.IsTop() {
		if err = s.dao.Mc.AddTop(c, rp); err != nil {
			log.Error("s.dao.Mc.AddTop(%d,%d,%d) error(%v)", rp.Oid, rp.RpID, rp.Type, err)
		}
	}
	return
}
