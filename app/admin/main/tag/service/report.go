package service

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_reportURI             = "\"http://www.bilibili.com/video/av%d\""
	_reportSendMsg         = "您在视频av%d举报的项目已处理,感谢反馈#{点击进入查看}{%s}"
	_deductIntegralContext = "您在视频av%d由于%s,您被举报处理%d节操值,具体原因请看#{节操记录}{\"http://account.bilibili.com/site/record?type=moral\"}"
	_deductIntegralTitle   = "您被举报处理%d节操"
	_reportFeedback        = "举报反馈"
	_msgMCVersion11        = "2_1_11"
	_msgMCVersion10        = "2_1_10"
	_punishReason          = "标签违规操作:%s"
)

var (
	_emtpyReportLog = make([]*model.ReportLog, 0)
)

// ReportList ReportList.
func (s *Service) ReportList(c context.Context, param *model.ParamReportList) (res []*model.ReportInfo, total int64, err error) {
	res = make([]*model.ReportInfo, 0, param.ParamPage.Ps)
	var (
		sqlReport []string
		arcs      map[int64]*model.SearchRes
		esTagList = new(model.MngSearchTagList)
	)
	punishLogMap := make(map[int64]map[int32]*model.ReportLog)
	noneLogMap := make(map[int64]*model.ReportLog)
	tagMap := make(map[int64]*model.Tag)
	arcMap := make(map[int64]*model.SearchRes)
	start := (param.Pn - 1) * param.Ps
	end := param.Ps
	order := "rpt.oid"
	if param.State == model.AuditNotHanleFirst || param.State == model.AuditHanledFirst {
		order = "rpt.ctime"
	}
	if param.TName != "" {
		es := &model.ESTag{
			Sort:    model.DefaultSort,
			Order:   model.DefaultOrder,
			Keyword: param.TName,
			ParamPage: model.ParamPage{
				Pn: model.DefaultPageNum,
				Ps: model.DefaultSearchNum,
			},
		}
		if esTagList, err = s.dao.ESearchTag(c, es); err != nil {
			return
		}
		for _, t := range esTagList.Result {
			param.Tids = append(param.Tids, t.ID)
		}
	}
	if len(param.Tids) != 0 {
		sqlReport = append(sqlReport, fmt.Sprintf("rpt.tid in (%s)", xstr.JoinInts(param.Tids)))
	}
	if param.RptMid > 0 {
		sqlReport = append(sqlReport, fmt.Sprintf("ru.mid=%d", param.RptMid))
	}
	if param.Reason != 0 {
		sqlReport = append(sqlReport, fmt.Sprintf("rpt.reason=%d", param.Reason))
	}
	if param.Mid > 0 {
		sqlReport = append(sqlReport, fmt.Sprintf("rpt.mid=%d", param.Mid))
	}
	if len(param.STime) > 0 && param.Oid == 0 && param.Mid == 0 && param.RptMid == 0 && param.TName == "" {
		sqlReport = append(sqlReport, fmt.Sprintf("rpt.ctime >= %q", param.STime))
	}
	if len(param.ETime) > 0 && param.Oid == 0 && param.Mid == 0 && param.RptMid == 0 && param.TName == "" {
		sqlReport = append(sqlReport, fmt.Sprintf("rpt.ctime <= %q", param.ETime))
	}
	if param.Oid > 0 {
		sqlReport = append(sqlReport, fmt.Sprintf("rpt.oid=%d ", param.Oid))
		sqlReport = append(sqlReport, fmt.Sprintf("rpt.type=%d ", param.Type))
	}
	if len(param.Rid) != 0 {
		sqlReport = append(sqlReport, fmt.Sprintf("rpt.rid in (%s)", xstr.JoinInts(param.Rid)))
	}
	if total, err = s.dao.ReportCount(c, sqlReport, order, param.State); err != nil {
		return
	}
	var (
		rpts    []*model.ReportInfo
		rptTids []int64
		rptIDs  []int64
		rptOids []int64
	)
	if rpts, rptTids, rptIDs, rptOids, err = s.dao.ReportInfoList(c, sqlReport, order, param.State, start, end); err != nil {
		return
	}
	if len(rptIDs) > 0 {
		punishLogMap, noneLogMap, _ = s.spliteReportLog(c, rptIDs)
	}
	if len(rptTids) > 0 {
		_, tagMap, _ = s.dao.Tags(c, rptTids)
	}
	if param.Type == model.ResTypeArchive && len(rptOids) > 0 {
		arcs, _, _ = s.arcInfos(c, rptOids)
	}
	for _, v := range arcs {
		if _, ok := arcMap[v.ID]; ok {
			continue
		}
		arcMap[v.ID] = v
	}
	for _, v := range rpts {
		var handleType int32 = -1
		v.Log = _emtpyReportLog
		if k, ok := arcMap[v.Oid]; ok {
			v.Title = k.Title
			v.MissionID = k.MissionID
			if k.Mid == v.Mid {
				v.MidIsUp = 1
			}
			if k.Mid == v.RptMid {
				v.RptIsUp = 1
			}
		}
		if k, ok := tagMap[v.Tid]; ok {
			v.TName = k.Name
			v.TagState = k.State
		}
		if k, ok := punishLogMap[v.ID]; ok {
			logs := make([]*model.ReportLog, 0)
			for _, v := range k {
				logs = append(logs, v)
			}
			v.Log = logs
		}
		noneLog := noneLogMap[v.ID]
		if noneLog != nil {
			handleType = noneLog.HandleType
		} else {
			if v.State != 0 {
				handleType = model.HandleNull
			}
		}
		v.Examine = handleType
		res = append(res, v)
	}
	return
}

// ReportInfo ReportInfo.
func (s *Service) ReportInfo(c context.Context, id int64) (res []*model.ReportDetail, err error) {
	var (
		rpt *model.Report
		arc *model.SearchRes
	)
	res = make([]*model.ReportDetail, 0)
	userMap := make(map[int64]*model.ReportUser)
	tagMap := make(map[int64]*model.Tag)
	if rpt, err = s.dao.Report(c, id); err != nil {
		return
	}
	if rpt == nil {
		return nil, ecode.TagReportNotExist
	}
	if arc, err = s.arcInfo(c, rpt.Oid); err != nil {
		return
	}
	if arc == nil {
		return nil, ecode.ArchiveNotExist
	}
	rpts, rptIDs, tids, _ := s.dao.ReportByOidMid(c, rpt.Mid, rpt.Oid, rpt.Type)
	if len(tids) > 0 {
		_, tagMap, _ = s.dao.Tags(c, tids)
	}
	if len(rptIDs) > 0 {
		_, userMap, _ = s.dao.ReportUsers(c, rptIDs)
	}
	for _, v := range rpts {
		r := &model.ReportDetail{
			ID:         v.ID,
			Oid:        v.Oid,
			Type:       v.Type,
			Tid:        v.Tid,
			Mid:        v.Mid,
			Action:     v.Action,
			Rid:        v.Rid,
			Title:      arc.Title,
			Reason:     v.Reason,
			IsDelMoral: v.Moral,
			Score:      v.Score,
			State:      v.State,
			CTime:      v.CTime,
			MTime:      v.MTime,
		}
		if k, ok := userMap[v.ID]; ok {
			r.RptMid = k.Mid
		}
		if k, ok := tagMap[v.Tid]; ok {
			r.TName = k.Name
			r.TagState = k.State
		}
		res = append(res, r)
	}
	return
}

// ReportHandle ReportHandle.
func (s *Service) ReportHandle(c context.Context, uname string, uid, id int64, audit, action int32) (err error) {
	var (
		state      int32
		handleType = int32(-1)
		rpt        *model.Report
		tag        *model.Tag
		rptUser    *model.ReportUser
	)
	if rpt, err = s.dao.Report(c, id); err != nil {
		return
	}
	if rpt == nil {
		return ecode.TagReportNotExist
	}
	if tag, err = s.dao.Tag(c, rpt.Tid); err != nil {
		return
	}
	if tag == nil {
		return ecode.TagNotExist
	}
	if action == model.ActionAdd {
		err = s.RelationAdd(c, tag.Name, rpt.Oid, uid, rpt.Type)
	} else {
		err = s.RelationDelete(c, tag.ID, rpt.Oid, rpt.Type)
	}
	if rpt.State == model.AuditHanledFirst || rpt.State == model.AuditHanledSecond {
		if rpt.State == model.AuditHanledFirst {
			state = model.AuditNotHanleSecond
			if action == model.ActionDel {
				handleType = model.HandleRestoreDelFirst
			} else {
				handleType = model.HandleRestoreAddFirst
			}
		} else {
			state = model.AuditHanledSecond
			if action == model.ActionDel {
				handleType = model.HandleRestoreDelSecond
			} else {
				handleType = model.HandleRestoreAddSecond
			}
		}
		s.dao.UpReportState(c, rpt.ID, state)
		l := &model.ReportLog{
			RptID:      rpt.ID,
			UserName:   uname,
			Oid:        rpt.Oid,
			Type:       rpt.Type,
			Mid:        rpt.Mid,
			Tid:        rpt.Tid,
			Rid:        rpt.Rid,
			HandleType: handleType,
		}
		s.addReportLog(c, l)
		return
	}
	if audit == model.AuditFirst {
		state = model.AuditHanledFirst
		if err == nil {
			switch action {
			case model.ActionDel:
				handleType = model.HandleDelFirst
			case model.ActionAdd:
				handleType = model.HandleAddFirst
			}
		} else {
			switch action {
			case model.ActionDel:
				handleType = model.HandleDelUserFirst
			case model.ActionAdd:
				handleType = model.HandleAddUserFirst
			}
		}
	}
	if audit == model.AuditSecond {
		state = model.AuditHanledSecond
		if err == nil {
			switch action {
			case model.ActionDel:
				handleType = model.HandleDelSecond
			case model.ActionAdd:
				handleType = model.HandleAddSecond
			}
		} else {
			switch action {
			case model.ActionDel:
				handleType = model.HandleDelUserSecond
			case model.ActionAdd:
				handleType = model.HandleAddUserSecond
			}
		}
	}
	s.dao.UpReportState(c, rpt.ID, state)
	if err == nil {
		if rptUser, err = s.dao.ReportUser(c, rpt.ID); err != nil {
			return
		}
		if rptUser == nil {
			err = nil
			return
		}
		s.dao.SendMsg(c, _msgMCVersion11, _reportFeedback, fmt.Sprintf(_reportSendMsg, rpt.Oid, fmt.Sprintf(_reportURI, rpt.Oid)), model.DataType, []int64{rptUser.Mid})
	}
	l := &model.ReportLog{
		RptID:      rpt.ID,
		UserName:   uname,
		Oid:        rpt.Oid,
		Type:       rpt.Type,
		Mid:        rpt.Mid,
		Tid:        rpt.Tid,
		Rid:        rpt.Rid,
		HandleType: handleType,
	}
	s.addReportLog(c, l)
	return
}

// ReportState ReportState.
func (s *Service) ReportState(c context.Context, id int64, state int32) (err error) {
	var r *model.Report
	if r, err = s.dao.Report(c, id); err != nil {
		return
	}
	if r == nil {
		return ecode.TagReportNotExist
	}
	if _, err = s.dao.UpReportState(c, id, state); err != nil {
		return ecode.TagOperateFail
	}
	return
}

// ReportIgnore ReportIgnore.
func (s *Service) ReportIgnore(c context.Context, uname string, audit int32, ids []int64) (err error) {
	var (
		state, handleType int32
		logs              []*model.ReportLog
		rpt               []*model.Report
		rptIDs            []int64
	)
	handleType = model.HandleWait
	switch audit {
	case model.AuditFirst:
		state = model.AuditNotHanleSecond
		handleType = model.HandleIgnoreFirst
	case model.AuditSecond:
		state = model.AuditHanledSecond
		handleType = model.HandleIgnoreSecond
	}
	if rpt, _, rptIDs, err = s.dao.Reports(c, ids); err != nil {
		return
	}
	if len(rptIDs) > 0 {
		if _, err = s.dao.UpReportsState(c, rptIDs, state); err != nil {
			return
		}
	}
	for _, v := range rpt {
		l := &model.ReportLog{
			Oid:        v.Oid,
			Type:       v.Type,
			Tid:        v.Tid,
			Mid:        v.Mid,
			RptID:      v.ID,
			Rid:        v.Rid,
			HandleType: handleType,
			UserName:   uname,
		}
		logs = append(logs, l)
	}
	s.addReportLogs(c, logs)
	return
}

// ReportDelete ReportDelete.
func (s *Service) ReportDelete(c context.Context, uname string, uid int64, ids []int64, audit int32) (err error) {
	for _, id := range ids {
		s.ReportHandle(c, uname, uid, id, audit, model.ActionDel)
	}
	return
}

// ReportPunish ReportPunish.
func (s *Service) ReportPunish(c context.Context, uname, remark, note string, ids []int64, reason, moral, isNotify, reasonType, blockTimeLength int32) (err error) {
	var (
		rptMap map[int64]*model.Report
		rptIDs []int64
	)
	if _, rptMap, rptIDs, err = s.dao.Reports(c, ids); err != nil {
		return
	}
	if rptMap == nil {
		return ecode.TagReportNotExist
	}
	if moral != 0 {
		if err = s.deductIntegral(c, uname, remark, note, reason, moral, isNotify, rptMap, rptIDs); err != nil {
			return
		}
	}
	if reasonType != 0 && len(ids) == 1 {
		r, ok := rptMap[ids[0]]
		if !ok {
			return
		}
		log.Warn("block user report(%v), blockTimeLength(%d)", r, blockTimeLength)
		if blockTimeLength > 0 || blockTimeLength == -1 {
			if err = s.blockUser(c, uname, note, r, reasonType, blockTimeLength, isNotify); err != nil {
				return
			}
		}
	}
	return
}

// ReportLogList ReportLogList.
func (s *Service) ReportLogList(c context.Context, oid, tid, mid, rid int64, tp, pn, ps int32, handleType []int64, stime, etime, username string) (res []*model.ReportLog, total int64, err error) {
	var (
		sql    string
		sqlStr []string
		tids   []int64
	)
	res = make([]*model.ReportLog, 0)
	start := (pn - 1) * ps
	end := ps
	tagMap := make(map[int64]*model.Tag)
	if oid > 0 {
		sqlStr = append(sqlStr, fmt.Sprintf("r.oid=%d", oid))
	}
	if tid > 0 {
		sqlStr = append(sqlStr, fmt.Sprintf("r.tid=%d", tid))
	}
	if tp > 0 {
		sqlStr = append(sqlStr, fmt.Sprintf("r.type=%d", tp))
	}
	if mid > 0 {
		sqlStr = append(sqlStr, fmt.Sprintf("r.mid=%d", mid))
	}
	if rid > 0 {
		sqlStr = append(sqlStr, fmt.Sprintf("r.rid=%d", rid))
	}
	if len(handleType) > 0 {
		sqlStr = append(sqlStr, fmt.Sprintf("r.handle_type in (%s)", xstr.JoinInts(handleType)))
	}
	if len(username) > 0 {
		sqlStr = append(sqlStr, fmt.Sprintf("r.username=%q ", username))
	}
	if len(stime) > 0 {
		sqlStr = append(sqlStr, fmt.Sprintf("r.ctime >= %q ", stime))
	}
	if len(etime) > 0 {
		sqlStr = append(sqlStr, fmt.Sprintf("r.ctime <= %q ", etime))
	}
	if len(sqlStr) > 0 {
		sql = fmt.Sprintf(" WHERE %s ", strings.Join(sqlStr, " AND "))
	}
	if total, err = s.dao.ReportLogCount(c, sql); err != nil {
		return
	}
	if res, tids, err = s.dao.ReportLogList(c, sql, start, end); err != nil {
		return
	}
	if len(tids) > 0 {
		_, tagMap, _ = s.dao.Tags(c, tids)
	}
	for _, v := range res {
		if v.Tid > 0 {
			if k, ok := tagMap[v.Tid]; ok {
				v.Tag = k
			}
		}
	}
	return
}

// ReportLogInfo ReportLogInfo.
func (s *Service) ReportLogInfo(c context.Context, id int64) (res []*model.ReportLog, err error) {
	var (
		tids []int64
		tags map[int64]*model.Tag
	)
	if res, tids, err = s.dao.ReportLog(c, id); err != nil {
		return
	}
	if len(res) == 0 {
		return _emtpyReportLog, nil
	}
	_, tags, _ = s.dao.Tags(c, tids)
	if len(tags) == 0 {
		return
	}
	for _, v := range res {
		if tag, ok := tags[v.Tid]; ok {
			v.Tag = tag
		}
	}
	return
}

func (s *Service) addReportLog(c context.Context, l *model.ReportLog) (err error) {
	if l == nil {
		return ecode.TagReportLogAddFailed
	}
	if l.HandleType == model.HandleDelSecond || l.HandleType == model.HandleAddSecond || l.HandleType == model.HandleIgnoreSecond {
		l.Reason = ""
	}
	affect, err := s.dao.AddReportLog(c, l)
	if err != nil {
		return
	}
	if affect <= 0 {
		err = ecode.TagReportLogAddFailed
	}
	return
}

func (s *Service) addReportLogs(c context.Context, logs []*model.ReportLog) (err error) {
	var (
		sqls   []string
		affect int64
	)
	sql := " (%d,%q,%d,%d,%d,%d,%d,%d,%q,%d,%d) "
	for _, v := range logs {
		if v.HandleType == model.HandleDelSecond || v.HandleType == model.HandleAddSecond || v.HandleType == model.HandleIgnoreSecond {
			v.Reason = ""
		}
		s := fmt.Sprintf(sql, v.RptID, v.UserName, v.Points, v.Oid, v.Type, v.Mid, v.Tid, v.Rid, v.Reason, v.HandleType, v.Notice)
		sqls = append(sqls, s)
	}
	if len(sqls) > 0 {
		if affect, err = s.dao.AddReportLogs(c, sqls); err != nil {
			return
		}
		if affect <= 0 {
			err = ecode.TagReportLogAddFailed
		}
	}
	return
}

func (s *Service) blockUser(c context.Context, uname, note string, rpt *model.Report, reasontype, blockTimeLength, isNotify int32) (err error) {
	var (
		points, blockForever int32
		title                string
		tag                  *model.Tag
	)
	action := "删除"
	if blockTimeLength == -1 {
		blockTimeLength = 0
		blockForever = 1
		points = -1
	} else {
		points = blockTimeLength
	}
	if rpt.Action == 0 {
		action = "新增"
	}
	if tag, err = s.dao.Tag(c, rpt.Tid); err != nil {
		return
	}
	if tag == nil {
		return ecode.TagNotExist
	}
	arc, _ := s.arcInfo(c, rpt.Oid)
	if arc != nil {
		title = arc.Title
	}
	if err = s.dao.BlockUser(c, tag.Name, title, uname, action, note, rpt.Mid, rpt.Oid, reasontype, isNotify, blockTimeLength, blockForever); err != nil {
		return
	}
	l := &model.ReportLog{
		RptID:      rpt.ID,
		HandleType: model.HandleBlock,
		Notice:     isNotify,
		Points:     points,
		Reason:     note,
		Rid:        rpt.Rid,
		UserName:   uname,
		Oid:        rpt.Oid,
		Type:       rpt.Type,
		Tid:        rpt.Tid,
		Mid:        rpt.Mid,
	}
	s.addReportLog(c, l)
	return
}

func (s *Service) deductIntegral(c context.Context, uname, remark, note string, reason, moral, isNotify int32, rptMap map[int64]*model.Report, rptIDs []int64) (err error) {
	var (
		mids []int64
		logs []*model.ReportLog
	)
	for _, rpt := range rptMap {
		l := &model.ReportLog{
			RptID:      rpt.ID,
			HandleType: model.HandleIntegral,
			Notice:     isNotify,
			Points:     moral,
			Reason:     note,
			Rid:        rpt.Rid,
			UserName:   uname,
			Oid:        rpt.Oid,
			Type:       rpt.Type,
			Tid:        rpt.Tid,
			Mid:        rpt.Mid,
		}
		logs = append(logs, l)
		mids = append(mids, rpt.Mid)
		if isNotify == 1 {
			s.dao.SendMsg(c, _msgMCVersion10, fmt.Sprintf(_deductIntegralTitle, moral), fmt.Sprintf(_deductIntegralContext, rpt.Oid, s.reason(reason), moral), model.DataType, mids)
		}
	}
	s.addReportLogs(c, logs)
	if moral != 0 {
		if _, err = s.dao.UpdateMorals(c, model.MoralHasDeducted, rptIDs); err != nil {
			return
		}
	}
	if remark == "" {
		remark = s.reason(reason)
	}
	if err = s.dao.AddMoral(c, uname, remark, s.reason(reason), moral, isNotify, mids); err != nil {
		return
	}
	return
}

// reportLogList 获取扣节操与封禁与其他的日志，分开存储
func (s *Service) spliteReportLog(c context.Context, rptIDs []int64) (punishMap map[int64]map[int32]*model.ReportLog, noneMap map[int64]*model.ReportLog, err error) {
	_, logMap, err := s.dao.ReportLogByRptID(c, rptIDs)
	punishMap = make(map[int64]map[int32]*model.ReportLog)
	noneMap = make(map[int64]*model.ReportLog)
	for rptID, logs := range logMap {
		var (
			strTwo *model.ReportLog
		)
		str := make(map[int32]*model.ReportLog)
		for _, v := range logs {
			if v.HandleType == model.HandleIntegral || v.HandleType == model.HandleBlock {
				if _, ok := str[v.HandleType]; !ok {
					str[v.HandleType] = v
				}
				continue
			}
			if strTwo == nil {
				strTwo = v
			}
		}
		punishMap[rptID] = str
		noneMap[rptID] = strTwo
	}
	return
}

// reason reason.
func (s *Service) reason(r int32) string {
	switch r {
	case 1:
		return fmt.Sprintf(_punishReason, "内容不相关")
	case 2:
		return fmt.Sprintf(_punishReason, "敏感信息")
	case 3:
		return fmt.Sprintf(_punishReason, "恶意攻击")
	case 4:
		return fmt.Sprintf(_punishReason, "剧透内容")
	case 5:
		return fmt.Sprintf(_punishReason, "恶意删除")
	case 6:
		return fmt.Sprintf(_punishReason, "大量违规操作")
	case 7:
		return fmt.Sprintf(_punishReason, "其他信息")
	default:
		return fmt.Sprintf(_punishReason, "")
	}
}
