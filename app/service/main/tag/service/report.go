package service

import (
	"context"
	"time"

	v1 "go-common/app/service/main/tag/api"
	"go-common/app/service/main/tag/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	xtime "go-common/library/time"
)

func (s *Service) checkRepeatReport(c context.Context, oid, tid, mid, rptMid int64, typ, action int32) (rptID int64, err error) {
	rptMap, err := s.dao.Report(c, oid, typ)
	if err != nil {
		return
	}
	for _, rpt := range rptMap {
		if rpt.Tid == tid && rpt.Mid == rptMid && rpt.Action == action {
			if rpt.Completed() {
				err = ecode.TagRptNotRptPassed
				return
			}
			rptID = rpt.ID
			break
		}
	}
	if rptID <= 0 {
		return
	}
	rptUsers, err := s.dao.ReportUser(c, rptID)
	if b, ok := rptUsers[mid]; ok && b {
		err = ecode.TagArcTagRpted
	}
	return
}

// AddReport add user report of resource.
func (s *Service) AddReport(c context.Context, arg *model.AddReportReq) (res *v1.AddReportReply, err error) {
	res = new(v1.AddReportReply)
	rtMap, err := s.resTags(c, arg.Oid, arg.Type)
	if err != nil {
		return
	}
	rt, ok := rtMap[arg.Tid]
	if !ok || rt.State != model.ResStateNormal {
		err = ecode.TagResTagNotExist
		return
	}
	if rt.Role == model.ResRoleAdmin {
		err = ecode.TagAdminOpCanNotTpt
		return
	}
	rptID, err := s.checkRepeatReport(c, arg.Oid, arg.Tid, arg.Mid, rt.Mid, arg.Type, model.LogActionAdd)
	if err != nil {
		return
	}
	now := time.Now()
	rpt := &model.Report{
		Oid:    arg.Oid,
		Type:   arg.Type,
		Tid:    arg.Tid,
		Mid:    rt.Mid,
		TypeID: arg.PartID,
		Action: model.LogActionAdd,
		Count:  1,
		Score:  arg.Score,
		Reason: arg.ReasonID,
		CTime:  xtime.Time(now.Unix()),
		MTime:  xtime.Time(now.Unix()),
	}
	urpt := &model.ReportUser{
		Mid: arg.Mid,
	}
	if rptID <= 0 {
		urpt.Attr = 1
	}
	if whiteUserMap, ok := s.whiteUserMap.Load().(map[int64]struct{}); ok {
		if _, ok := whiteUserMap[arg.Mid]; ok {
			urpt.SetManager()
		}
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		return
	}
	if urpt.RptID, err = s.dao.TxAddReport(tx, rpt); err != nil {
		tx.Rollback()
		return
	}
	if _, err = s.dao.TxAddUserReport(tx, urpt); err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	return
}

// ReportAction .
func (s *Service) ReportAction(c context.Context, oid, logID, mid int64, typ, partID, reason, score int32, content, ip string) (err error) {
	_, err = s.checkResType(c, oid, typ, ip)
	if err != nil {
		return
	}
	rl, err := s.dao.ResourceLog(c, oid, logID, typ)
	if err != nil {
		return
	}
	err = s.addReport(c, rl, mid, partID, reason, score, content)
	return
}

func (s *Service) addReport(c context.Context, rl *model.ResourceLog, mid int64, partID, reason, score int32, content string) (err error) {
	rs, err := s.dao.ReportAndUser(c, rl.Oid, mid, rl.Tid, rl.Type, rl.Action)
	if err != nil {
		return
	}
	urpt := &model.ReportUser{
		Mid: mid,
	}
	if len(rs) == 0 {
		rpt := &model.Report{
			Oid:     rl.Oid,
			Type:    rl.Type,
			Tid:     rl.Tid,
			Mid:     rl.Mid, // action操作人
			TypeID:  partID,
			Count:   1,
			Score:   score,
			Action:  rl.Action,
			Reason:  reason,
			Content: content,
		}
		urpt.Attr = urpt.Attr | 1 // 第一举报人
		var tx *sql.Tx
		tx, err = s.dao.BeginTran(c)
		if err != nil {
			return
		}
		var rptID int64
		if rptID, err = s.dao.TxAddReport(tx, rpt); err != nil {
			tx.Rollback()
			return
		}
		if rptID == 0 {
			tx.Rollback()
			return
		}
		urpt.RptID = rptID
		if _, err = s.dao.TxAddUserReport(tx, urpt); err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}
	urpt.RptID = rs[0].ID
	rptMap, err := s.dao.ReportUser(c, urpt.RptID)
	if err != nil {
		return
	}
	if _, ok := rptMap[mid]; ok {
		err = ecode.TagArcTagRpted
		return
	}
	_, err = s.dao.AddUserReport(c, urpt)
	return
}
