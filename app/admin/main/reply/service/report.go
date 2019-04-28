package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/reply/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"
)

func forbidResult(ftime int64) string {
	if ftime == -1 {
		return "永久"
	}
	return fmt.Sprintf("%d天", ftime)
}

func (s *Service) report(c context.Context, oid, rpID int64) (report *model.Report, err error) {
	if report, err = s.dao.Report(c, oid, rpID); err != nil {
		return
	}
	if report == nil {
		err = ecode.ReplyReportNotExist
	}
	return
}

func (s *Service) reports(c context.Context, oids, rpIDs []int64) (res map[int64]*model.Report, err error) {
	if res, err = s.dao.Reports(c, oids, rpIDs); err != nil {
		return
	}
	return
}

// ReportSearch return report result from search.
func (s *Service) ReportSearch(c context.Context, sp *model.SearchReportParams, page, pageSize int64) (res *model.SearchReportResult, err error) {
	if res, err = s.dao.SearchReport(c, sp, page, pageSize); err != nil {
		log.Error("s.dao.SearchReport(%+v,%d,%d) error(%v)", sp, page, pageSize, err)
		return
	}
	filterIds := map[int64]string{}
	oids := map[int64]string{}
	admins := map[int64]string{}
	titles := map[int64]string{}
	for _, r := range res.Result {
		r.OidStr = strconv.FormatInt(r.Oid, 10)
		if strings.Contains(r.Message, "*") {
			filterIds[r.ID] = ""
		} else if len(r.Attr) > 0 {
			for _, attr := range r.Attr {
				if attr == 4 {
					filterIds[r.ID] = ""
				}
			}
		}
		admins[r.AdminID] = ""
		if r.ReplyMid == r.ArcMid {
			r.IsUp = 1
		}
		oids[r.Oid] = ""
		if int32(r.Type) == model.SubTypeArchive {
			titles[r.Oid] = ""
		}
	}
	s.linkByOids(c, oids, sp.Type)
	s.titlesByOids(c, titles)
	s.dao.AdminName(c, admins)
	s.dao.FilterContents(c, filterIds)
	for i, data := range res.Result {
		res.Result[i].AdminName = admins[res.Result[i].AdminID]
		if content, ok := filterIds[data.ID]; ok && content != "" {
			data.Message = content
		}
		if content, ok := oids[data.Oid]; ok {
			data.RedirectURL = fmt.Sprintf("%s#reply%d", content, data.ID)
		}
		if int32(data.Type) == model.SubTypeArchive && data.Title == "" {
			if title := titles[data.Oid]; title != "" {
				data.Title = title
			}
		}
	}

	return
}

// ReportIgnore ignore a report.
func (s *Service) ReportIgnore(c context.Context, oids, rpIDs []int64, adminID int64, adName string, typ, audit int32, remark string, delReport bool) (err error) {
	var (
		state  int32
		op     int32
		result string
		action string
	)
	if audit == model.AuditTypeFirst {
		result = "一审忽略"
		op = model.AdminOperRptIgnore1
		state = model.ReportStateIgnore1
		action = model.ReportActionReportIgnore1
	} else if audit == model.AuditTypeSecond {
		result = "二审忽略"
		op = model.AdminOperRptIgnore2
		state = model.ReportStateIgnore2
		action = model.ReportActionReportIgnore2
	} else {
		err = ecode.RequestErr
		return
	}
	subs, err := s.subjects(c, oids, typ)
	if err != nil {
		log.Error("ReportIgnore subjects(%v,%v,%d) error(%v)", oids, typ, state, err)
		return
	}
	rps, err := s.replies(c, oids, rpIDs)
	if err != nil {
		log.Error("ReportIgnore replies(%v,%v,%d) error(%v)", oids, rpIDs, state, err)
		return
	}
	rpts, err := s.reports(c, oids, rpIDs)
	if err != nil {
		log.Error("ReportIgnore reports(%v,%v,%d) error(%v)", oids, rpIDs, state, err)
		return
	}
	now := time.Now()
	rows, err := s.dao.UpReportsState(c, oids, rpIDs, state, now)
	if err != nil {
		log.Error("ReportIgnore UpReportsState(%v,%v,%d) rows:%d error(%v)", oids, rpIDs, state, rows, err)
		return
	} else if rows == 0 {
		return
	}
	for _, rp := range rps {
		rpt, ok := rpts[rp.ID]
		if !ok {
			continue
		}
		sub, ok := subs[rpt.Oid]
		if !ok {
			continue
		}
		s.pubEvent(c, model.EventReportDel, rpt.Mid, sub, rp, rpt)
	}
	for _, rpt := range rpts {
		report.Manager(&report.ManagerInfo{
			UID:      adminID,
			Uname:    adName,
			Business: 41,
			Type:     int(typ),
			Oid:      rpt.Oid,
			Ctime:    now,
			Action:   action,
			Index: []interface{}{
				rpt.RpID,
				rpt.State,
				state,
			},
			Content: map[string]interface{}{
				"remark": remark,
			},
		})
		rpt.State = state
		rpt.MTime = xtime.Time(now.Unix())
		if delReport {
			s.cache.Do(c, func(ctx context.Context) {
				s.dao.DelReport(ctx, rpt.Oid, rpt.RpID)
			})
		}

	}
	s.addAdminLogs(c, rps, adminID, typ, model.AdminIsNew, model.AdminIsReport, op, result, remark, now)
	s.cache.Do(c, func(ctx context.Context) {
		s.pubSearchReport(ctx, rpts, nil)
	})
	return
}

// ReportDel delete reply by report.
func (s *Service) ReportDel(c context.Context, oids, rpIDs []int64, adminID, ftime int64, typ, audit, moral, reason, freason int32, notify bool, adminName, remark, content string) (err error) {
	return s.reportDel(c, oids, rpIDs, adminID, ftime, typ, audit, moral, reason, freason, notify, adminName, remark, content)
}

func (s *Service) reportDel(c context.Context, oids, rpIDs []int64, adminID, ftime int64, typ, audit, moral, reason, freason int32, notify bool, adminName, remark, content string) (err error) {
	var (
		op       int32
		rptState int32
		action   string
	)
	if audit == model.AuditTypeFirst {
		op = model.AdminOperRptDel1
		rptState = model.ReportStateDelete1
		action = model.ReportActionReportDel1
	} else if audit == model.AuditTypeSecond {
		op = model.AdminOperRptDel2
		rptState = model.ReportStateDelete2
		action = model.ReportActionReportDel2
	} else {
		err = ecode.RequestErr
		return
	}
	now := time.Now()
	rpts, err := s.reports(c, oids, rpIDs)
	if err != nil {
		log.Error("reportDel reports(%v, %v) error(%v)", oids, rpIDs, err)
		return
	}

	wg, ctx := errgroup.WithContext(c)
	for _, rept := range rpts {
		rpt := rept
		wg.Go(func() (err error) {
			isPunish := false
			var sub *model.Subject
			var rp *model.Reply
			sub, rp, err = s.delReply(ctx, rpt.Oid, rpt.RpID, model.StateDelAdmin, now)
			if err != nil {
				if ecode.ReplyDeleted.Equal(err) && rp.IsDeleted() {
					isPunish = true
					err = nil
				} else {
					log.Error("reportDel tranDel(%d, %d) error(%v)", rpt.Oid, rpt.RpID, err)
					return err
				}
			}
			s.delCache(ctx, sub, rp)
			rea := reason
			if reason == -1 {
				rea = rpt.Reason
			}
			report.Manager(&report.ManagerInfo{
				UID:      adminID,
				Uname:    adminName,
				Business: 41,
				Type:     int(typ),
				Oid:      rpt.Oid,
				Ctime:    now,
				Action:   action,
				Index: []interface{}{
					rpt.RpID,
					rpt.State,
					rptState,
				},
				Content: map[string]interface{}{
					"moral":   moral,
					"notify":  notify,
					"ftime":   ftime,
					"freason": freason,
					"reason":  rea,
					"remark":  remark,
				},
			})
			rpt.State = rptState
			rpt.MTime = xtime.Time(now.Unix())
			rpt.ReplyCtime = rp.CTime
			s.pubEvent(ctx, model.EventReportDel, rpt.Mid, sub, rp, rpt)
			s.cache.Do(c, func(ctx context.Context) {
				s.dao.DelReport(ctx, rp.Oid, rp.ID)
				s.moralAndNotify(ctx, rp, moral, notify, rp.Mid, adminID, adminName, remark, rea, freason, ftime, isPunish)
			})
			return nil
		})
	}
	if err = wg.Wait(); err != nil {
		return
	}
	var rows int64
	if reason == -1 {
		rows, err = s.dao.UpReportsState(c, oids, rpIDs, rptState, now)
	} else {
		rows, err = s.dao.UpReportsStateWithReason(c, oids, rpIDs, rptState, reason, content, now)
	}
	if err != nil || rows == 0 {
		return
	}
	s.cache.Do(c, func(ctx context.Context) {
		state := model.StateDelAdmin
		s.pubSearchReport(ctx, rpts, &state)
	})
	s.addAdminIDLogs(c, oids, rpIDs, adminID, typ, model.AdminIsNew, model.AdminIsReport, op, fmt.Sprintf("已通过举报并删除封禁%s/扣除%d节操", forbidResult(ftime), moral), remark, now)

	return
}

// ReportRecover recover a report reply.
func (s *Service) ReportRecover(c context.Context, oids, rpIDs []int64, adminID int64, typ, audit int32, remark string) (err error) {
	s.reportRecover(c, oids, rpIDs, adminID, typ, audit, remark)
	return
}

func (s *Service) reportRecover(c context.Context, oids, rpIDs []int64, adminID int64, typ, audit int32, remark string) (err error) {
	var (
		rptState int32
		op       int32
		result   string
	)
	if audit == model.AuditTypeFirst {
		result = "一审恢复评论"
		op = model.AdminOperRptRecover1
		rptState = model.ReportStateIgnore1
	} else if audit == model.AuditTypeSecond {
		result = "二审恢复评论"
		op = model.AdminOperRptRecover2
		rptState = model.ReportStateIgnore2
	} else {
		err = ecode.RequestErr
		return
	}
	now := time.Now()
	rpts, err := s.reports(c, oids, rpIDs)
	if err != nil {
		log.Error("ReportRecover reports(%v, %v) error(%v)", oids, rpIDs, err)
		return
	}
	var rps = map[int64]*model.Reply{}
	wg, ctx := errgroup.WithContext(c)
	for _, report := range rpts {
		rpt := report
		wg.Go(func() (err error) {
			sub, rp, err := s.recReply(ctx, rpt.Oid, rpt.RpID, model.StateNormal, now)
			if err != nil {
				log.Error("ReportRecover tranRecover(%d, %d) error(%v)", rpt.Oid, rpt.RpID, err)
				return
			}
			rpt.MTime = xtime.Time(now.Unix())
			rpt.State = rptState
			s.pubEvent(ctx, model.EventReportRecover, rpt.Mid, sub, rp, rpt)
			rps[rp.ID] = rp
			return
		})
	}
	if err = wg.Wait(); err != nil {
		return
	}
	rows, err := s.dao.UpReportsState(c, oids, rpIDs, rptState, now)
	if err != nil {
		return
	} else if rows == 0 {
		return
	}
	s.addAdminIDLogs(c, oids, rpIDs, adminID, typ, model.AdminIsNew, model.AdminIsReport, op, result, remark, now)
	s.cache.Do(c, func(ctx context.Context) {
		s.pubSearchReply(ctx, rps, model.StateNormal)
	})
	return
}

// ReportTransfer transfer a report.
func (s *Service) ReportTransfer(c context.Context, oids, rpIDs []int64, adminID int64, adName string, typ, audit int32, remark string) (err error) {
	var (
		state  int32
		op     int32
		result string
		action string
	)
	if audit == model.AuditTypeFirst {
		result = "二审转一审"
		op = model.AdminOperRptTransfer1
		state = model.ReportStateNew
		action = model.ReportActionReport2To1
	} else if audit == model.AuditTypeSecond {
		result = "一审转二审"
		op = model.AdminOperRptTransfer2
		state = model.ReportStateNew2
		action = model.ReportActionReport1To2
	} else {
		err = ecode.RequestErr
		return
	}
	rpts, err := s.reports(c, oids, rpIDs)
	if err != nil {
		return
	}
	now := time.Now()
	var tranOids, tranRpIDs []int64
	var tranRpt []*model.Report
	for _, rpt := range rpts {
		if rpt.AttrVal(model.ReportAttrTransferred) == model.AttrNo {
			rpt.State = state
			rpt.MTime = xtime.Time(now.Unix())
			tranOids = append(tranOids, rpt.Oid)
			tranRpIDs = append(tranRpIDs, rpt.RpID)
			tranRpt = append(tranRpt, rpt)
		}
	}
	rows, err := s.dao.UpReportsState(c, tranOids, tranRpIDs, state, now)
	if err != nil {
		log.Error("s.dao.UpReportState(%v,%v,%d) rows:%d error(%v)", oids, rpIDs, state, rows, err)
		return
	} else if rows == 0 {
		return
	}
	for _, rpt := range tranRpt {
		report.Manager(&report.ManagerInfo{
			UID:      adminID,
			Uname:    adName,
			Business: 41,
			Type:     int(typ),
			Oid:      rpt.Oid,
			Ctime:    now,
			Action:   action,
			Index: []interface{}{
				rpt.RpID,
				rpt.State,
				state,
			},
			Content: map[string]interface{}{
				"remark": remark,
			},
		})
	}
	if _, err = s.dao.UpReportsAttrBit(c, tranOids, tranRpIDs, model.ReportAttrTransferred, model.AttrYes, now); err != nil {
		log.Error("s.dao.UpReportAttrBit(%v,%v) transfered errror(%v)", tranOids, tranRpIDs, err)
		return
	}
	s.addAdminIDLogs(c, oids, rpIDs, adminID, typ, model.AdminIsNew, model.AdminIsReport, op, result, remark, now)
	s.cache.Do(c, func(ctx context.Context) {
		s.pubSearchReport(ctx, rpts, nil)
	})
	return
}

// ReportStateSet transfer a report.
func (s *Service) ReportStateSet(c context.Context, oids, rpIDs []int64, adminID int64, adname string, typ, state int32, remark string, delReport bool) (err error) {
	if state != model.ReportStateTransferred {
		return
	}
	var (
		op     int32
		result string
		action string
	)
	op = model.AdminOperRptTransferArbitration
	action = model.ReportActionReportArbitration
	if adminID == 0 {
		result = "系统自动移交至风纪委"
	} else {
		result = "管理员移交至风纪委"
	}
	rpts, err := s.reports(c, oids, rpIDs)
	if err != nil {
		return
	}
	rps, err := s.replies(c, oids, rpIDs)
	if err != nil {
		log.Error("s.replies (%v,%v) error(%v)", oids, rpIDs, err)
		return
	}
	links := make(map[int64]string, 0)
	titles := make(map[int64]string, 0)
	for _, rp := range rps {
		title, link, _ := s.TitleLink(c, rp.Oid, rp.Type)
		links[rp.ID] = link
		titles[rp.ID] = title
	}
	err = s.dao.TransferArbitration(c, rps, rpts, adminID, adname, titles, links)
	if err != nil {
		log.Error("s.dao.TransferArbitration (%v,%v) error(%v)", rps, rpts, err)
		return
	}
	mtime := time.Now()
	for _, rpt := range rpts {
		report.Manager(&report.ManagerInfo{
			UID:      adminID,
			Uname:    adname,
			Business: 41,
			Type:     int(typ),
			Oid:      rpt.Oid,
			Ctime:    mtime,
			Action:   action,
			Index: []interface{}{
				rpt.RpID,
				rpt.State,
				state,
			},
			Content: map[string]interface{}{
				"remark": remark,
			},
		})
		rpt.State = state
		rpt.MTime = xtime.Time(mtime.Unix())
		if delReport {
			s.cache.Do(c, func(ctx context.Context) {
				s.dao.DelReport(ctx, rpt.Oid, rpt.RpID)
			})
		}
	}
	rows, err := s.dao.UpReportsState(c, oids, rpIDs, state, mtime)
	if err != nil {
		log.Error("s.dao.UpReportState(%v,%v,%d) rows:%d error(%v)", oids, rpIDs, state, rows, err)
		return
	} else if rows == 0 {
		return
	}
	s.addAdminIDLogs(c, oids, rpIDs, adminID, typ, model.AdminIsNew, model.AdminIsReport, op, result, remark, mtime)
	s.cache.Do(c, func(ctx context.Context) {
		s.pubSearchReport(ctx, rpts, nil)
	})
	return
}
