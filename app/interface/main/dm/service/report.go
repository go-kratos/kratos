package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/dm/model"
	"go-common/app/interface/main/dm2/model/oplog"
	arcMdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_reportLock    = 1800
	_reportLimit   = 10
	_rptAutoDelCnt = 3
)

// AddReport add dm report.
func (s *Service) AddReport(c context.Context, cid, dmid, uid int64, reason int8, content string) (id int64, err error) {
	var (
		needDelete, needHide bool
		state                int8
		score                int32
	)
	t, err := s.dao.RptBrigTime(c, uid)
	if err != nil {
		return
	}
	if time.Now().Unix()-t < _reportLock {
		err = ecode.DMReportLimit
		return
	}
	cnt, err := s.dao.RptCnt(c, uid)
	if err != nil {
		return
	}

	if cnt >= _reportLimit {
		if err = s.dao.AddRptBrig(c, uid); err != nil {
			log.Error("s.dao.AddRptBrig(%d) error(%v)", uid, err)
		}
		err = ecode.DMReportLimit
		return
	}
	dms, err := s.dms(c, model.SubTypeVideo, cid, []int64{dmid})
	if err != nil {
		return
	}
	if len(dms) == 0 || !dms[0].NeedDisplay() {
		err = ecode.DMReportNotExist
		return
	}
	if _, err = s.dao.UptRptCnt(c, uid); err != nil {
		log.Error("s.dao.UptRptCnt(%d) error(%v)", uid, err)
		err = nil
	}
	rpt, err := s.dao.Report(c, cid, dmid)
	if err != nil {
		log.Error("s.dao.Report(cid:%d,dmid:%d) error(%v)", cid, dmid, err)
		return
	}
	// score
	if score, err = s.dao.FigureInfo(c, uid); err != nil {
		log.Error("s.dao.FigureInfo(uid: %d) error(%v)", uid, err)
		return
	}
	if score < 0 || score > 100 {
		log.Error("s.AddReport.score illegal(uid: %d) score(%v)", uid, score)
		score = 0
	} else {
		score = 100 - score
	}
	if rpt != nil {
		state, err = s.checkReasonType(c, reason, cid, dmid, rpt)
		if err != nil {
			log.Error("s.checkReasonType(reason: %d) error(%v)", reason, err)
			return
		}
		/*
			以下弹幕不能自动删除与隐藏
			1. mode7、mode8、mode9弹幕
			2. 字幕弹幕
			3. 保护弹幕
			4. up主在自己的视频下发送的弹幕
		*/
		if state == model.StatSecondInit &&
			dms[0].Content.Mode != 7 &&
			dms[0].Content.Mode != 8 &&
			dms[0].Content.Mode != 9 &&
			dms[0].AttrVal(model.AttrProtect) != model.AttrYes &&
			dms[0].Pool == 0 &&
			!s.isUpper(c, cid, dms[0].Mid) {
			if (reason == model.ReportReasonAd && rpt.Reason == model.ReportReasonAd) || rpt.Count >= _rptAutoDelCnt-1 || rpt.Score+score >= 160 {
				needDelete = true
			} else {
				needHide = true
			}
		}
	} else {
		state = s.reasonType(reason)
		if state == model.StatSecondInit &&
			dms[0].Content.Mode != 7 &&
			dms[0].Content.Mode != 8 &&
			dms[0].Content.Mode != 9 &&
			dms[0].AttrVal(model.AttrProtect) != model.AttrYes &&
			dms[0].Pool == 0 &&
			!s.isUpper(c, cid, dms[0].Mid) {
			needHide = true
		}
	}
	nowTime := time.Now()
	r := &model.Report{
		Cid:     cid,
		Did:     dmid,
		UID:     uid,
		Reason:  reason,
		Content: content,
		Count:   1,
		State:   state,
		UpOP:    model.StatUpperInit,
		Score:   score,
		RpTime:  nowTime,
		Ctime:   nowTime,
		Mtime:   nowTime,
	}
	u := &model.User{
		Did:     dmid,
		UID:     uid,
		State:   model.NoticeUnsend,
		Reason:  reason,
		Content: content,
		Ctime:   nowTime,
		Mtime:   nowTime,
	}
	lastID, err := s.dao.AddReportUser(c, u)
	if err != nil {
		log.Error("s.dao.AddReportUser(%v) error(%v)", u, err)
		return
	}
	if lastID < 1 {
		err = ecode.DMReportExist
		return
	}
	if id, err = s.dao.AddReport(c, r); err != nil {
		log.Error("s.dao.AddReport(%v) error(%v)", r, err)
		return
	}
	if needDelete {
		select {
		case s.delDMReportChan <- rpt:
		default:
			log.Error("s.delDMReportChan.full(%v)", rpt)
		}
	}
	if needHide && !needDelete {
		select {
		case s.hideDMReportChan <- r:
		default:
			log.Error("s.hideDMReportChan.full(%v)", r)
		}
	}
	return
}

// EditReport edit report dm.
func (s *Service) EditReport(c context.Context, tp int32, cid, mid, dmid int64, op int8) (affect int64, err error) {
	r, err := s.dao.Report(c, cid, dmid)
	if err != nil {
		log.Error("s.dao.Report(cid:%d, dmid:%d), error(%v)", cid, dmid, err)
		return
	}
	if r == nil {
		err = ecode.DMReportNotExist
		return
	}
	if !s.isUpper(c, cid, mid) {
		err = ecode.ArchiveOwnerErr
		return
	}
	if affect, err = s.dao.UpdateReportUPOp(c, cid, dmid, op); err != nil {
		log.Error("s.dao.UpdateReportUPOp(cid:%d, dmid:%d, op:%d) error(%v)", cid, dmid, op, err)
	}
	if op == model.StatUpperDelete {
		if err = s.EditDMState(c, tp, cid, mid, model.StateDelete, oplog.SourceUp, oplog.OperatorUp, dmid); err != nil {
			log.Error("s.EditDMStat(cid:%d,dmid:%d) error(%v)", cid, dmid, err)
			return
		}
	}
	v := &model.UptSearchReport{
		DMid:  dmid,
		Upop:  op,
		Ctime: r.Ctime.Format("2006-01-02 15:04:05"),
		Mtime: time.Now().Format("2006-01-02 15:04:05"),
	}
	err = s.dao.UpdateSearchReport(c, []*model.UptSearchReport{v})
	return
}

// ReportList 获取一个用户的所有被举报的弹幕
func (s *Service) ReportList(c context.Context, mid, aid, page, size int64, upOp int8, states []int64) (res *model.RptSearchs, err error) {
	res = &model.RptSearchs{}
	var (
		aidsMap  = make(map[int64]struct{})
		aids     = []int64{}
		cidDMids = make(map[int64][]int64)
		stateMap = make(map[int64]int64)
	)
	rptSearch, err := s.dao.SearchReport(c, mid, aid, page, size, upOp, states)
	if err != nil || rptSearch == nil {
		return
	}
	for _, v := range rptSearch.Result {
		if _, ok := aidsMap[v.Aid]; !ok {
			aidsMap[v.Aid] = struct{}{}
			aids = append(aids, v.Aid)
		}
		cidDMids[v.Cid] = append(cidDMids[v.Cid], v.Did)
	}
	archives := s.archiveInfos(c, aids)
	if stateMap, err = s.dmState(c, cidDMids); err != nil {
		return
	}
	for _, v := range rptSearch.Result {
		if arc, ok := archives[v.Aid]; ok {
			v.Title = arc.Title
			v.Cover = arc.Pic
		}
		if state, ok := stateMap[v.Did]; ok {
			v.Deleted = state
		}
	}
	res = &model.RptSearchs{
		Page:      rptSearch.Page.Num,
		PageSize:  rptSearch.Page.Size,
		PageCount: (rptSearch.Page.Total-1)/rptSearch.Page.Size + 1,
		Total:     rptSearch.Page.Total,
		Result:    rptSearch.Result,
	}
	return
}

func (s *Service) isUpper(c context.Context, cid, mid int64) bool {
	sub, err := s.subject(c, 1, cid)
	if err != nil {
		return false
	}
	return sub.Mid == mid
}

// checkReasonType get state by report reason
func (s *Service) checkReasonType(c context.Context, reason int8, cid, dmid int64, rpt *model.Report) (state int8, err error) {
	reportLog, err := s.dao.ReportLog(c, dmid)
	if err != nil {
		log.Error("s.dao.ReportLog(dmid:%d) error(%v)", dmid, err)
		return
	}
	if len(reportLog) > 0 { // 如果这个举报已经在后台被处理了,那就根据举报当前的状态返回对应的一二审状态
		if rpt.State == model.StatFirstInit ||
			rpt.State == model.StatFirstDelete ||
			rpt.State == model.StatFirstIgnore {
			state = model.StatFirstInit
		} else {
			state = model.StatSecondInit
		}
	} else { // 否则就根据举报理由来返回状态
		state = s.reasonType(reason)
	}
	return
}

// ReportArchives 获取一个用户的所有被举报的稿件
func (s *Service) ReportArchives(c context.Context, mid int64, upOp int8, states []int8, pn, ps int64) (res *model.Archives, err error) {
	res = &model.Archives{}
	aids, err := s.dao.SearchReportAid(c, mid, upOp, states, pn, ps)
	if err != nil || len(aids) == 0 {
		return
	}
	archiveInfos := s.archiveInfos(c, aids)
	for aid, info := range archiveInfos {
		res.Result = append(res.Result, &struct {
			Aid   int64  `json:"aid"`
			Title string `json:"title"`
		}{Aid: aid, Title: info.Title})
	}
	return
}

func (s *Service) reasonType(reason int8) (state int8) {
	if reason == model.ReportReasonProhibited ||
		reason == model.ReportReasonPorn ||
		reason == model.ReportReasonFraud ||
		reason == model.ReportReasonAttack ||
		reason == model.ReportReasonPrivate ||
		reason == model.ReportReasonTeenagers {
		state = model.StatFirstInit
	} else {
		state = model.StatSecondInit
	}
	return
}

// deleteDMReport delete dm report.
func (s *Service) deleteDMReport(c context.Context, rpt *model.Report) (err error) {
	sub, err := s.subject(c, 1, rpt.Cid)
	if err != nil {
		return
	}
	rpt.State = model.StatSecondAutoDelete
	if err = s.EditDMState(c, 1, rpt.Cid, sub.Mid, model.StateScriptDelete, oplog.SourceUp, oplog.OperatorUp, rpt.Did); err != nil {
		log.Error("s.EditDMStat(cid:%d, state:%d, dmid:%d ) error(%v)", rpt.Cid, rpt.State, rpt.Did, err)
		return
	}
	if _, err = s.dao.UpdateReportStat(c, rpt.Cid, rpt.Did, rpt.State); err != nil {
		log.Error("s.dao.UpdateReportStat(cid:%d, state:%d, dmid:%d ) error(%v)", rpt.Cid, rpt.State, rpt.Did, err)
		return
	}
	lg := &model.RptLog{
		Did:     rpt.Did,
		AdminID: 0,
		Reason:  rpt.Reason,
		Result:  rpt.State,
		Remark:  "自动删除",
		Elapsed: int64(time.Since(rpt.Mtime).Seconds()),
		Ctime:   time.Now(),
		Mtime:   time.Now(),
	}
	if err = s.dao.AddReportLog(c, lg); err != nil {
		log.Error("s.dao.AddReportLog(log:%v) error(%v)", lg, err)
		return
	}
	users, err := s.dao.ReportUser(c, rpt.Did)
	if err != nil {
		log.Error("s.dao.ReportUser(dmid:%d) error(%v)", rpt.Did, err)
		return
	}
	if err = s.dao.SetReportUserFinished(c, rpt.Did); err != nil {
		log.Error("s.dao.SetReportUserFinished(dmid:%d) error(%v)", rpt.Did, err)
		return
	}
	ct, err := s.dao.Content(c, rpt.Cid, rpt.Did)
	if err != nil || ct == nil {
		return
	}
	arc, err := s.acvSvc.Archive3(c, &arcMdl.ArgAid2{Aid: sub.Pid})
	if err != nil {
		log.Error("s.acvSvc.Archive3(%d) error(%v)", sub.Pid, err)
		return
	}
	for _, u := range users {
		content := fmt.Sprintf(model.RptMsgTemplate, arc.Title, arc.Aid, ct.Msg, model.ReportReason[rpt.Reason])
		s.dao.SendNotify(c, model.RptMsgTitle, content, []int64{u.UID})
	}
	return
}

//hideDMReport hide reported dm
func (s *Service) hideDMReport(c context.Context, rpt *model.Report) (err error) {
	dmids := []int64{rpt.Did}
	//change dm state to hide
	_, err = s.dao.UpdateDMStat(c, 1, rpt.Cid, model.StateHide, dmids)
	if err != nil {
		log.Error("s.dao.UpdateDMStat(oid:%d state:%d dmids:%v) error(%v)", rpt.Cid, model.StateHide, rpt.Did, err)
		return
	}
	// send hideTime to databus
	time := time.Now().Unix() + 72000
	act := &model.ReportAction{
		Cid:      rpt.Cid,
		Did:      rpt.Did,
		HideTime: time,
	}
	// make sure all the hided dms are in one partition
	if err = s.dao.SendAction(context.TODO(), "1", act); err != nil {
		log.Error("databus.Send(%+v) error(%v)", act, err)
	}
	return
}

func (s *Service) dmReportProc() {
	for {
		select {
		case rpt := <-s.delDMReportChan:
			if err := s.deleteDMReport(context.TODO(), rpt); err != nil {
				log.Error("s.deleteDMReport(rpt:%v) error(%v)", rpt, err)
			}
		case rpt := <-s.hideDMReportChan:
			if err := s.hideDMReport(context.TODO(), rpt); err != nil {
				log.Error("s.hideDMReport(rpt:%v) error(%v)", rpt, err)
			}
		}
	}
}

func (s *Service) dmState(c context.Context, cidDmids map[int64][]int64) (stateMap map[int64]int64, err error) {
	var (
		idxMap map[int64]*model.DM
		tp     = int32(1)
	)
	stateMap = make(map[int64]int64)
	for oid, dmids := range cidDmids {
		if idxMap, _, err = s.dao.IndexsByID(c, tp, oid, dmids); err != nil {
			return
		}
		for dmid, dm := range idxMap {
			stateMap[dmid] = int64(dm.State)
		}
	}
	return
}
