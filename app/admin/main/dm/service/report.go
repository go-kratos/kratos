package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"go-common/app/admin/main/dm/dao"
	"go-common/app/admin/main/dm/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_searchTimeFormat = "2006-01-02 15:04:05"
)

// ChangeReportStat set dm report status by mult dmid.
func (s *Service) ChangeReportStat(c context.Context, cidDmids map[int64][]int64, state, reason, notice int8, adminID, block, blockReason, moral int64, remark, operator string) (affect int64, err error) {
	var (
		optDur   int64
		dmids    []int64
		dmidList = make([]int64, 0)
		nowTime  = time.Now()
		dmLogMap = make(map[int64][]*model.ReportLog)
		rptsMap  = make(map[int64]*model.Report)
		uptRpts  = make([]*model.UptSearchReport, 0)
	)
	for _, dmids2 := range cidDmids {
		dmids = append(dmids, dmids2...)
	}
	if state == model.StatFirstInit || state == model.StatSecondInit || state == model.StatJudgeInit {
		if rptsMap, err = s.reports(c, dmids); err != nil {
			log.Error("s.reports(cidDmids:%v) error(%v)", cidDmids, err)
			return
		}
	} else {
		if rptsMap, err = s.reportsDetail(c, dmids); err != nil {
			log.Error("s.reportsDetail(cidDmids:%v) error(%v)", cidDmids, err)
			return
		}
	}
	for cid, dmids := range cidDmids {
		if state == model.StatSecondIgnore || state == model.StatFirstIgnore {
			if err = s.dao.IgnoreReport(c, cid, dmids, state); err != nil {
				log.Error("s.dao.IgnoreReport(cid:%d, dmid:%v) error(%v)", cid, dmids, err)
				return 0, err
			}
		} else {
			if err = s.dao.ChangeReportStat(c, cid, dmids, state); err != nil {
				log.Error("s.dao.ChangeReportStat(cid:%d, dmid:%v) error(%v)", cid, dmids, err)
				return 0, err
			}
		}
		var rpts []*model.Report
		if rpts, err = s.dao.Reports(c, cid, dmids); err != nil {
			log.Error("s.dao.Reports(cid:%d, dmids:%v) error(%v)", cid, dmids, err)
			err = nil
		} else {
			for _, rpt := range rpts {
				var ctime, mtime time.Time
				ctime, err = time.Parse(time.RFC3339, rpt.Ctime)
				if err != nil {
					log.Error("strconv.RarseInt(%s) error(%v)", rpt.Ctime, err)
					err = nil
					continue
				}
				mtime, err = time.Parse(time.RFC3339, rpt.Mtime)
				if err != nil {
					log.Error("strconv.RarseInt(%s) error(%v)", rpt.Mtime, err)
					err = nil
					continue
				}
				uptRpt := &model.UptSearchReport{
					DMid:  rpt.Did,
					Ctime: ctime.Format("2006-01-02 15:04:05"),
					Mtime: mtime.Format("2006-01-02 15:04:05"),
					State: state,
				}
				uptRpts = append(uptRpts, uptRpt)
			}
		}
		for _, dmid := range dmids {
			rpt, ok := rptsMap[dmid]
			if !ok {
				err = fmt.Errorf("get report detail empty,dmid:%d", dmid)
				log.Error("s.ReportsDetail(cid:%d, dmid:%v) error(%v)", cid, dmid, err)
				continue
			}
			if isDeleteOperation(state) {
				s.addRptDelAction(rpt)
			}
			rpt.State = state
			var mtime time.Time
			if mtime, err = time.ParseInLocation("2006-01-02 15:04:05", rpt.Mtime, time.Local); err == nil {
				optDur = int64(time.Since(mtime).Seconds())
			}
			lg := &model.ReportLog{
				Did:     dmid,
				AdminID: adminID,
				Reason:  reason,
				Result:  state,
				Remark:  remark,
				Elapsed: optDur,
				Ctime:   nowTime,
				Mtime:   nowTime,
			}
			dmLogMap[dao.LogTable(rpt.Did)] = append(dmLogMap[dao.LogTable(rpt.Did)], lg)
			dmidList = append(dmidList, rpt.Did)
			// if moral > 0 {
			// 	s.reduceMoral(c, rpt.UID, moral, reason, uname, fmt.Sprintf("%s, cid:%d, dmid:%d", model.CheckStateBelong(state), cid, dmid))
			// }
			if block != 0 {
				s.blockUser(c, rpt, block, blockReason, moral, operator)
			}
			if notice == model.NoticeReporter || notice == model.NoticeAll { // 发送邮件给举报方
				if len(rpt.RptUsers) > 0 {
					s.sendMsgToReporter(c, rpt, block, blockReason, int64(reason))
				}
			}
			if notice == model.NoticePoster || notice == model.NoticeAll {
				if len(rpt.RptUsers) > 0 {
					s.sendMsgToPoster(c, rpt, block, blockReason, int64(reason))
				}
			}
		}
		// if delete or recover this danmu
		if isDeleteOperation(state) {
			tmpRemark := model.AdminRptReason[reason] + "，" + model.BlockReason[int8(blockReason)]
			if len(remark) >= 0 {
				tmpRemark = remark + "，" + tmpRemark
			}
			// s.dao.SetStateByIDs(c, model.SubTypeVideo, cid, dmids, model.StateReportDelete)
			// s.OpLog(c, cid, adminID, 1, dmids, "status", "", fmt.Sprint(model.StateReportDelete), tmpRemark, oplog.SourceManager, oplog.OperatorAdmin)
			s.editDmState(c, model.SubTypeVideo, model.StateReportDelete, cid, reason, dmids, float64(moral), adminID, operator, tmpRemark)
		}
	}
	if isDeleteOperation(state) {
		// search update ignore error
		s.uptSearchDmState(c, model.SubTypeVideo, model.StateReportDelete, cidDmids)
	}
	if !(state == model.StatSecondInit || state == model.StatFirstInit) {
		if err = s.ChangeReportUserStat(c, dmidList); err != nil {
			log.Error("s.ChangeReportUserStat(%v) error(%v)", dmidList, err)
		}
	}
	if len(uptRpts) > 0 {
		if err = s.dao.UptSearchReport(c, uptRpts); err != nil {
			log.Error("s.dao.UpSearchReport(%v) error(%v)", uptRpts, err)
			err = nil
		}
	}
	for k, v := range dmLogMap {
		if len(v) <= 0 {
			continue
		}
		if err = s.dao.AddReportLog(c, k, v); err != nil {
			log.Error("s.dao.AddReportLog(%v) error(%v)", v, err)
			return
		}
	}
	return
}

func isDeleteOperation(state int8) bool {
	if state == model.StatSecondDelete || state == model.StatFirstDelete || state == model.StatSecondAutoDelete || state == model.StatJudgeDelete {
		return true
	}
	return false
}

func (s *Service) dmState(c context.Context, cidDmids map[int64][]int64) (stateMap map[int64]int64, err error) {
	var (
		idxMap map[int64]*model.DM
		tp     = int32(1)
	)
	stateMap = make(map[int64]int64)
	for oid, dmids := range cidDmids {
		if idxMap, _, err = s.dao.IndexsByID(c, tp, oid, dmids); err != nil {
			log.Error("s.dmState(oid:%v,dmids:%v) err(%v)", oid, dmids, err)
			return
		}
		for dmid, dm := range idxMap {
			stateMap[dmid] = int64(dm.State)
		}
	}
	return
}

// ReportList2 .
func (s *Service) ReportList2(c context.Context, params *model.ReportListParams) (rtList *model.ReportList, err error) {
	var (
		aidMap   = make(map[int64]bool)
		aids     []int64
		cidDmids = make(map[int64][]int64)
		stateMap = make(map[int64]int64)
	)
	if params.Start == "" {
		now := time.Now()
		params.Start = time.Date(now.Year(), now.Month(), now.Day()-3, 0, 0, 0, 0, now.Location()).Format(_searchTimeFormat)
	}
	if params.End == "" {
		now := time.Now()
		params.End = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()).Format(_searchTimeFormat)
	}
	rptSearch, err := s.dao.SearchReport2(c, params)
	if err != nil {
		log.Error("s.dao.SearchReport2(params:%+v) error(%v)", params, err)
		return
	}
	for _, v := range rptSearch.Result {
		aidMap[v.Aid] = true
		cidDmids[v.Cid] = append(cidDmids[v.Cid], v.Did)
	}
	for aid := range aidMap {
		aids = append(aids, aid)
	}
	archives, err := s.archiveInfos(c, aids)
	if err != nil {
		log.Error("s.archives(%v) error(%v)", aids, err)
		return
	}
	if stateMap, err = s.dmState(c, cidDmids); err != nil {
		return
	}
	for _, v := range rptSearch.Result {
		v.DidStr = strconv.FormatInt(v.Did, 10)
		if arc, ok := archives[v.Aid]; ok {
			v.Title = arc.Title
		}
		if state, ok := stateMap[v.Did]; ok {
			v.Deleted = state
		}
	}
	rtList = &model.ReportList{
		Code:      rptSearch.Code,
		Order:     rptSearch.Order,
		Page:      rptSearch.Page.Num,
		PageSize:  rptSearch.Page.Size,
		PageCount: (rptSearch.Page.Total-1)/rptSearch.Page.Size + 1,
		Total:     rptSearch.Page.Total,
		Result:    rptSearch.Result,
	}
	return
}

// ReportList get report list from search
func (s *Service) ReportList(c context.Context, page, size int64, start, end, order, sort, keyword string, tid, rpID, state, upOp []int64, rt *model.Report) (rtList *model.ReportList, err error) {
	var (
		aidMap   = make(map[int64]bool)
		aids     []int64
		cidDmids = make(map[int64][]int64)
		stateMap = make(map[int64]int64)
	)
	rptSearch, err := s.dao.SearchReport(c, page, size, start, end, order, sort, keyword, tid, rpID, state, upOp, rt)
	if err != nil {
		log.Error("s.dao.SearchReport() error(%v)", err)
		return
	}
	for _, v := range rptSearch.Result {
		aidMap[v.Aid] = true
		cidDmids[v.Cid] = append(cidDmids[v.Cid], v.Did)
	}
	for aid := range aidMap {
		aids = append(aids, aid)
	}
	archives, err := s.archiveInfos(c, aids)
	if err != nil {
		log.Error("s.archives(%v) error(%v)", aids, err)
		return
	}
	if stateMap, err = s.dmState(c, cidDmids); err != nil {
		return
	}
	for _, v := range rptSearch.Result {
		v.DidStr = strconv.FormatInt(v.Did, 10)
		if arc, ok := archives[v.Aid]; ok {
			v.Title = arc.Title
		}
		if state, ok := stateMap[v.Did]; ok {
			v.Deleted = state
		}
	}
	rtList = &model.ReportList{
		Code:      rptSearch.Code,
		Order:     rptSearch.Order,
		Page:      rptSearch.Page.Num,
		PageSize:  rptSearch.Page.Size,
		PageCount: (rptSearch.Page.Total-1)/rptSearch.Page.Size + 1,
		Total:     rptSearch.Page.Total,
		Result:    rptSearch.Result,
	}
	return
}

func (s *Service) archiveInfos(c context.Context, aids []int64) (archives map[int64]*api.Arc, err error) {
	var (
		g        errgroup.Group
		mu       sync.Mutex
		l        = len(aids)
		pagesize = 50
		pagenum  = int(math.Ceil(float64(l) / float64(pagesize)))
	)
	archives = make(map[int64]*api.Arc)
	for i := 0; i < pagenum; i++ {
		start := i * pagesize
		end := (i + 1) * pagesize
		if end > l {
			end = l
		}
		g.Go(func() (err error) {
			arg := &archive.ArgAids2{Aids: aids[start:end]}
			res, err := s.arcRPC.Archives3(c, arg)
			if err != nil {
				log.Error("s.arcRPC.Archives3(%v) error(%v)", arg, err)
				return
			}
			for aid, info := range res {
				mu.Lock()
				archives[aid] = info
				mu.Unlock()
			}
			return
		})
	}
	err = g.Wait()
	return
}

// reportUsers get mult reports users
func (s *Service) reportUsers(c context.Context, dmids []int64) (rptUsers map[int64][]*model.ReportUser, err error) {
	var (
		g        errgroup.Group
		mu       sync.Mutex
		dmidsMap = map[int64][]int64{}
	)
	for _, dmid := range dmids {
		dmidsMap[dao.UserTable(dmid)] = append(dmidsMap[dao.UserTable(dmid)], dmid)
	}
	rptUsers = make(map[int64][]*model.ReportUser)
	for tableID, dmids := range dmidsMap {
		key, value := tableID, dmids
		g.Go(func() (err error) {
			userTmp, err := s.dao.ReportUsers(c, key, value, model.NoticeUnsend)
			if err != nil {
				return
			}
			for dmid, users := range userTmp {
				mu.Lock()
				rptUsers[dmid] = users
				mu.Unlock()
			}
			return
		})
	}
	err = g.Wait()
	return
}

// reportsDetail get report list from search and get user list、archive list、dm list.
func (s *Service) reportsDetail(c context.Context, dmids []int64) (res map[int64]*model.Report, err error) {
	var (
		aidMap       = make(map[int64]bool)
		aids, dmids2 []int64
	)
	reports, err := s.reports(c, dmids)
	if err != nil {
		return
	}
	for dmid, rpt := range reports {
		if _, ok := aidMap[rpt.Aid]; !ok {
			aidMap[rpt.Aid] = true
			aids = append(aids, rpt.Aid)
		}
		dmids2 = append(dmids2, dmid)
	}
	archives, err := s.archiveInfos(c, aids)
	if err != nil {
		log.Error("s.archives(%v) error(%v)", aids, err)
		return nil, err
	}
	rptUsers, err := s.reportUsers(c, dmids2)
	if err != nil {
		log.Error("s.rptUsers(%v) error(%v)", dmids2, err)
		return nil, err
	}
	res = make(map[int64]*model.Report)
	for dmid, rpt := range reports {
		if arc, ok := archives[rpt.Aid]; ok {
			rpt.Title = arc.Title
		}
		if users, ok := rptUsers[dmid]; ok {
			rpt.RptUsers = users
		} else {
			rpt.RptUsers = make([]*model.ReportUser, 0)
		}
		res[dmid] = rpt
	}
	return
}

// reports get report list by cid and dmids from search.
func (s *Service) reports(c context.Context, dmids []int64) (res map[int64]*model.Report, err error) {
	rptSearchList, err := s.dao.SearchReportByID(c, dmids)
	if err != nil || len(rptSearchList.Result) <= 0 {
		log.Error("dao.SearchReportByID(ids:%v) error(%v)", dmids, err)
		return
	}
	res = make(map[int64]*model.Report)
	for _, rpt := range rptSearchList.Result {
		res[rpt.Did] = rpt
	}
	return
}

// ReportLog get report log by dmid.
func (s *Service) ReportLog(c context.Context, dmid int64) (res []*model.ReportLog, err error) {
	if res, err = s.dao.ReportLog(c, dmid); err != nil {
		log.Error("s.dao.ReportLog(dmid:%d) error(%v)", dmid, err)
	}
	return
}

// ChangeReportUserStat change report_user data
func (s *Service) ChangeReportUserStat(c context.Context, dmids []int64) (err error) {
	var (
		dmidMap = map[int64][]int64{}
	)
	for _, v := range dmids {
		dmidMap[dao.UserTable(v)] = append(dmidMap[dao.UserTable(v)], v)
	}
	for k, v := range dmidMap {
		if _, err = s.dao.UpReportUserState(c, k, v, model.NoticeSend); err != nil {
			log.Error("s.dao.UpReportUserState(dmids:%v) error(%v)", v, err)
		}
	}
	return
}

func (s *Service) sendMsgToReporter(c context.Context, rpt *model.Report, block, blockReason, rptReason int64) {
	var (
		buf bytes.Buffer
	)
	for _, user := range rpt.RptUsers {
		buf.WriteString(fmt.Sprintf("%d,", user.UID))
	}
	buf.Truncate(buf.Len() - 1)
	m := &model.ReportMsg{
		Aid:         rpt.Aid,
		Did:         rpt.Did,
		Title:       rpt.Title,
		Msg:         rpt.Msg,
		RptReason:   int8(rptReason),
		Uids:        buf.String(),
		State:       rpt.State,
		Block:       block,
		BlockReason: int8(blockReason),
	}
	select {
	case s.msgReporterChan <- m:
	default:
		log.Error("s.msgReporterChan err, channel full(msg:%v)", m)
	}
}

func (s *Service) sendMsgToPoster(c context.Context, rpt *model.Report, block, blockReason, rptReason int64) {
	m := &model.ReportMsg{
		Aid:         rpt.Aid,
		Did:         rpt.Did,
		Title:       rpt.Title,
		Msg:         rpt.Msg,
		RptReason:   int8(rptReason),
		Uids:        fmt.Sprint(rpt.UID),
		Block:       block,
		BlockReason: int8(blockReason),
	}
	select {
	case s.msgPosterChan <- m:
	default:
		log.Error("s.msgPosterChan err, channel full(msg:%v)", m)
	}
}

func (s *Service) reduceMoral(c context.Context, uid, moral int64, reason int8, operator, remark string) {
	m := &model.ReduceMoral{
		UID:        uid,
		Moral:      moral,
		Origin:     2,
		Reason:     reason,
		ReasonType: 1,
		Operator:   operator,
		IsNotify:   0,
		Remark:     remark,
	}
	select {
	case s.reduceMoralChan <- m:
	default:
		log.Error("s.reduceMoral err, channel full(msg:%v)", m)
	}
}

func (s *Service) blockUser(c context.Context, rpt *model.Report, block, blockReason, moral int64, uname string) {
	var (
		blockEver   int64
		blockLength int64
	)
	if block == -1 {
		blockEver = 1
	} else {
		blockLength = block
	}
	m := &model.BlockUser{
		UID:             rpt.UID,
		BlockForever:    blockEver,
		BlockTimeLength: blockLength,
		BlockRemark:     fmt.Sprintf("%s, cid:%d, dmid:%d", model.CheckStateBelong(rpt.State), rpt.Cid, rpt.Did),
		Operator:        uname,
		OriginType:      2,
		Moral:           moral,
		ReasonType:      blockReason,
		OriginTitle:     rpt.Title,
		OriginContent:   rpt.Msg,
		OriginURL:       fmt.Sprintf("http://www.bilibili.com/av%d", rpt.Aid),
		IsNotify:        0,
	}
	select {
	case s.blockUserChan <- m:
	default:
		log.Error("s.blockUserChan err, channel full(msg:%v)", m)
	}
}

// DMReportJudge send report judge
func (s *Service) DMReportJudge(c context.Context, cidDmids map[int64][]int64, uid int64, uname string) (err error) {
	var (
		aids      []int64
		dmids     []int64
		rptJudges []*model.ReportJudge
	)
	for _, dmids2 := range cidDmids {
		dmids = append(dmids, dmids2...)
	}
	rpts, err := s.reportsDetail(c, dmids) // get report detail by multi dmids
	if len(rpts) == 0 {
		log.Error("dmjudge error! id:%v not exist in search", dmids)
		return
	}
	for _, rpt := range rpts {
		aids = append(aids, rpt.Aid)
	}
	arg := &archive.ArgAids2{Aids: aids}
	archs, err := s.arcRPC.Archives3(c, arg) // get archive info
	if err != nil {
		log.Error("s.arcSvc.Archives3(aids:%v) err(%v)", aids, err)
		return
	}
	if len(archs) == 0 {
		log.Error("dmjudge error! id:%v not exist in archive rpc", aids)
		err = ecode.ArchiveNotExist
		return
	}
	for _, rpt := range rpts {
		j := &model.ReportJudge{}
		arc, ok := archs[rpt.Aid]
		if !ok {
			continue
		}
		arg := &archive.ArgVideo2{
			Aid: rpt.Aid,
			Cid: rpt.Cid,
		}
		var vInfo *api.Page
		if vInfo, err = s.arcRPC.Video3(c, arg); err != nil {
			log.Error("s.arcSvc.Video3(arg:%v) err(%v)", arg, err)
			j.Page = 1
		} else {
			j.Page = int64(vInfo.Page)
		}
		j.MID = rpt.UID
		j.Operator = uname
		j.OperID = uid
		j.OContent = rpt.Msg
		j.OTitle = arc.Title
		j.OType = 2
		j.OURL = fmt.Sprintf("http://www.bilibili.com/av%d", rpt.Aid)
		j.ReasonType = int64(model.RpReasonToJudgeReason(int8(rpt.RpType)))
		j.AID = rpt.Aid
		j.OID = rpt.Cid
		j.RPID = rpt.Did
		sendTime, _ := time.Parse("2006-01-02 15:04:05", rpt.SendTime)
		j.BTime = sendTime.Unix()
		rptJudges = append(rptJudges, j)
	}
	if len(rptJudges) <= 0 {
		return
	}
	if err = s.dao.SendJudgement(c, rptJudges); err != nil {
		log.Error("s.dao.SendJudgement(data:%v) err (%v)", rptJudges, err)
	}
	_, err = s.ChangeReportStat(c, cidDmids, model.StatJudgeInit, 0, 0, uid, 0, 0, 0, "转风纪委", uname)
	if err != nil {
		log.Error("s.ChangeReportStat(id:%v) err(%v)", cidDmids, err)
		return
	}
	return
}

// JudgeResult receive judge result
func (s *Service) JudgeResult(c context.Context, cid, dmid, result int64) (err error) {
	var (
		state  int8
		remark string
	)
	res, err := s.dao.Reports(c, cid, []int64{dmid})
	if err != nil {
		log.Error("s.dao.Reports(cid:%d,dmid:%d) err(%v)", cid, dmid, err)
		return
	}
	if len(res) <= 0 {
		log.Error("dmJudge: cid:%d,dmid:%d not found", cid, dmid)
		err = ecode.RequestErr
		return
	}
	m := map[int64][]int64{
		res[0].Cid: {res[0].Did},
	}
	if result == 0 {
		state = model.StatJudgeIgnore
		remark = "风纪委处理:忽略"
	} else {
		state = model.StatJudgeDelete
		remark = "风纪委处理:删除"
	}
	_, err = s.ChangeReportStat(c, m, state, int8(res[0].RpType), 0, 0, 0, 0, 0, remark, "")
	if err != nil {
		log.Error("s.ChangeReportStat(cid:%d,dmid:%d) err(%v)", cid, dmid, err)
		return
	}
	return
}

func (s *Service) addRptDelAction(rpt *model.Report) (err error) {
	data, err := json.Marshal(rpt)
	if err != nil {
		log.Error("json.Marshal(%v) error(%v)", rpt, err)
		return
	}
	action := &model.Action{
		Oid:    rpt.Cid,
		Action: model.ActReportDel,
		Data:   data,
	}
	s.addAction(action)
	return
}
