package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"fmt"
	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// hdlArchiveMessage deal with archive action
func (s *Service) hdlArchiveMessage(action string, nwMsg []byte, oldMsg []byte) {
	var (
		err    error
		arc    = &archive.Archive{}
		oldArc = &archive.Archive{}
	)
	if action != _updateAct {
		return
	}
	if err = json.Unmarshal(nwMsg, arc); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", nwMsg, err)
		return
	}
	if err = json.Unmarshal(oldMsg, oldArc); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", oldMsg, err)
		return
	}
	if arc.TypeID != oldArc.TypeID {
		s.hdlMoveType(arc, oldArc)
	}
	if arc.Round != oldArc.Round {
		s.hdlRoundFlow(arc, oldArc)
	}
	if arc.State != oldArc.State && arc.State == archive.StateForbidUpDelete {
		s.arc.DelDispatchByAid(context.TODO(), arc.ID)
	}
}

// hdlMoveType deal with archive move typeid
func (s *Service) hdlMoveType(arc *archive.Archive, oldArc *archive.Archive) {
	if _, ok := archive.ReportArchiveRound[arc.Round]; !ok {
		return
	}
	s.arcMoveTypeCache.Lock()
	defer s.arcMoveTypeCache.Unlock()
	if _, ok := s.arcMoveTypeCache.Data[arc.Round]; !ok {
		s.arcMoveTypeCache.Data[arc.Round] = make(map[int16]map[string]int)
	}
	if _, ok := s.arcMoveTypeCache.Data[arc.Round][oldArc.TypeID]; !ok {
		s.arcMoveTypeCache.Data[arc.Round][oldArc.TypeID] = make(map[string]int)
	}
	if _, ok := s.arcMoveTypeCache.Data[arc.Round][arc.TypeID]; !ok {
		s.arcMoveTypeCache.Data[arc.Round][arc.TypeID] = make(map[string]int)
	}
	s.arcMoveTypeCache.Data[arc.Round][oldArc.TypeID]["out"]++
	s.arcMoveTypeCache.Data[arc.Round][arc.TypeID]["in"]++
}

// hdlRoundFlow deal with archive round flow
func (s *Service) hdlRoundFlow(arc *archive.Archive, oldArc *archive.Archive) {
	var (
		oper    *archive.Oper
		newOper *archive.Oper
		err     error
	)
	if _, ok := archive.ReportArchiveRound[oldArc.Round]; !ok {
		return
	}
	if oper, err = s.arc.LastRoundOper(context.TODO(), oldArc.ID, oldArc.Round); err != nil {
		log.Error("s.LastRoundOper(%d,%d) 获取archive_oper记录失败 error(%v)", oldArc.ID, oldArc.Round, err)
		return
	}
	if newOper, err = s.arc.NextRoundOper(context.TODO(), oper.ID, oldArc.ID, oldArc.Round); err != nil {
		log.Error("s.NextRoundOper(%d,%d,%d) 获取archive_oper记录失败 error(%v)", oper.ID, oldArc.ID, oldArc.Round, err)
		return
	}
	s.arcRoundFlowCache.Lock()
	defer s.arcRoundFlowCache.Unlock()
	if _, ok := s.arcRoundFlowCache.Data[oldArc.Round]; !ok {
		s.arcRoundFlowCache.Data[oldArc.Round] = make(map[int64]map[string]int)
	}
	if _, ok := s.arcRoundFlowCache.Data[oldArc.Round][oldArc.ID]; !ok {
		s.arcRoundFlowCache.Data[oldArc.Round][oldArc.ID] = make(map[string]int)
	}
	s.arcRoundFlowCache.Data[oldArc.Round][oldArc.ID]["take_time"] = int(newOper.MTime.Unix() - oper.MTime.Unix())
	s.arcRoundFlowCache.Data[oldArc.Round][oldArc.ID]["uid"] = int(newOper.UID)
	log.Info("s.arcRoundFlowCache.Data: %v", s.arcRoundFlowCache.Data)
}

// hdlMoveTypeCount check and write archive move type stats to db
func (s *Service) hdlMoveTypeCount() {
	var (
		report *archive.Report
		err    error
		bs     []byte
		ctime  = time.Now()
		mtime  = ctime
	)
	s.arcMoveTypeCache.Lock()
	defer s.arcMoveTypeCache.Unlock()
	if len(s.arcMoveTypeCache.Data) < 1 {
		log.Info("s.hdlMoveTypeCount() 统计内容为空，忽略：%v", s.arcMoveTypeCache.Data)
		return
	}
	if report, err = s.arc.ReportLast(context.TODO(), archive.ReportTypeArcMoveType); err != nil {
		log.Error("s.arc.ReportLast(%d) error(%v)", archive.ReportTypeArcMoveType, err)
		return
	}
	if report != nil && time.Now().Unix()-report.CTime.Unix() < 5*60 {
		log.Info("s.arc.ReportLast(%d) 距离上一次写入还没过5分钟!", archive.ReportTypeArcMoveType)
		return
	}
	if bs, err = json.Marshal(s.arcMoveTypeCache.Data); err != nil {
		log.Error("json.Marshal(%v) error(%v)", s.arcMoveTypeCache.Data, err)
		return
	}
	if _, err = s.arc.ReportAdd(context.TODO(), archive.ReportTypeArcMoveType, string(bs), ctime, mtime); err != nil {
		log.Error("s.arc.ReportAdd(%d,%s,%v,%v) error(%v)", archive.ReportTypeArcMoveType, string(bs), ctime, mtime, err)
		return
	}
	s.arcMoveTypeCache.Data = make(map[int8]map[int16]map[string]int)
}

// hdlRoundFlowCount check and write archive round flow stats to db
func (s *Service) hdlRoundFlowCount() {
	var (
		report *archive.Report
		err    error
		bs     []byte
		ctime  = time.Now()
		mtime  = ctime
	)
	s.arcRoundFlowCache.Lock()
	defer s.arcRoundFlowCache.Unlock()
	if len(s.arcRoundFlowCache.Data) < 1 {
		log.Info("s.hdlRoundFlowCount() 统计内容为空，忽略：%v", s.arcRoundFlowCache.Data)
		return
	}
	if report, err = s.arc.ReportLast(context.TODO(), archive.ReportTypeArcRoundFlow); err != nil {
		log.Error("s.arc.ReportLast(%d) error(%v)", archive.ReportTypeArcRoundFlow, err)
		return
	}
	if report != nil && time.Now().Unix()-report.CTime.Unix() < 5*60 {
		log.Info("s.arc.ReportLast(%d) 距离上一次写入还没过5分钟!", archive.ReportTypeArcRoundFlow)
		return
	}
	if bs, err = json.Marshal(s.arcRoundFlowCache.Data); err != nil {
		log.Error("json.Marshal(%v) error(%v)", s.arcRoundFlowCache.Data, err)
		return
	}
	if _, err = s.arc.ReportAdd(context.TODO(), archive.ReportTypeArcRoundFlow, string(bs), ctime, mtime); err != nil {
		log.Error("s.arc.ReportAdd(%d,%s,%v,%v) error(%v)", archive.ReportTypeArcRoundFlow, string(bs), ctime, mtime, err)
		return
	}
	s.arcRoundFlowCache.Data = make(map[int8]map[int64]map[string]int)
}

// MoveType get archive move type stats by typeid
func (s *Service) MoveType(c context.Context, stime, etime time.Time) (reports []*archive.Report, err error) {
	if reports, err = s.arc.Reports(c, archive.ReportTypeArcMoveType, stime, etime); err != nil {
		log.Error("s.arc.Reports(%d) err(%v)", archive.ReportTypeArcMoveType, err)
		return
	}
	return
}

// RoundFlow get archive round flow take time records
func (s *Service) RoundFlow(c context.Context, stime, etime time.Time) (reports []*archive.Report, err error) {
	if reports, err = s.arc.Reports(c, archive.ReportTypeArcRoundFlow, stime, etime); err != nil {
		log.Error("s.arc.Reports(%d) err(%v)", archive.ReportTypeArcRoundFlow, err)
		return
	}
	return
}

func (s *Service) arcUpdateproc(k int) {
	defer s.waiter.Done()
	for {
		var (
			upInfo *archive.UpInfo
			ok     bool
		)
		if upInfo, ok = <-s.arcUpChs[k]; !ok {
			log.Info("s.arcUpChs[k] closed", k)
			return
		}
		go s.hdlExcitation(upInfo.Nw, upInfo.Old)
		go s.hdlMonitorArc(upInfo.Nw, upInfo.Old)
		s.trackArchive(upInfo.Nw, upInfo.Old)
		go s.arcStateChange(upInfo.Nw, upInfo.Old, true)
	}
}

func (s *Service) putArcChan(action string, nwMsg []byte, oldMsg []byte) {
	var (
		err      error
		chanSize = int64(s.c.ChanSize)
	)
	nw := &archive.Archive{}
	if err = json.Unmarshal(nwMsg, nw); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", nwMsg, err)
		return
	}
	switch action {
	case _insertAct:
		s.arcUpChs[nw.ID%chanSize] <- &archive.UpInfo{Nw: nw, Old: nil}
	case _updateAct:
		old := &archive.Archive{}
		if err = json.Unmarshal(oldMsg, old); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", oldMsg, err)
			return
		}
		s.arcUpChs[nw.ID%chanSize] <- &archive.UpInfo{Nw: nw, Old: old}
	}
}

// secondRound 接收到databus的second_round消息。
func (s *Service) secondRound(c context.Context, m *archive.VideoupMsg) (err error) {
	var (
		a *archive.Archive
	)
	if a, err = s.arc.ArchiveByAid(c, m.Aid); err != nil || a.ID <= 0 {
		log.Error("secondRound s.arc.ArchiveByAid error(%v)/not found aid(%d)", err, m.Aid)
		return
	}
	s.dealFromList(c, m)
	s.dealMissionTag(c, m, a)
	//开评论逻辑判断
	s.arcStateChange(a, nil, true)
	if archive.NormalState(a.State) {
		s.adminBindTag(c, a.Mid, a.ID, a.Tag, a.TypeID)
	}
	//邮件发送开关
	if m.SendEmail {
		if m.AdminChange && !s.isPGC(a.ID) {
			s.sendMail(c, a, nil)
		}
		s.sendArchivePrivateEmail(c, a)
	}

	return
}

// dealFromList 处理from list流程
func (s *Service) dealFromList(c context.Context, m *archive.VideoupMsg) (err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			log.Error("s.dealFromList panic(%v)", pErr)
		}
	}()
	switch m.FromList {
	case archive.FromListHotReview: //热门回查
		var (
			has   = false
			state = archive.RecheckStateWait
		)
		//查询flow_design中是否存在禁止
		if has, err = s.arc.HasFlowGroup(c, archive.FlowPoolRecheck, archive.FlowGroupIDHot, m.Aid); err != nil {
			log.Error("s.updateRecheckState(%d,%d,%d,%d) error(%v)", archive.TypeHotRecheck, archive.FlowPoolRecheck, archive.FlowGroupIDHot, m.Aid, err)
			return
		}
		if has {
			state = archive.RecheckStateForbid
		} else {
			state = archive.RecheckStateNoForbid
		}
		s.updateRecheckState(c, archive.TypeHotRecheck, m.Aid, state)
	case archive.FromListExcitation: //激励回查
		s.updateRecheckState(c, archive.TypeExcitationRecheck, m.Aid, archive.RecheckStateNoForbid)
	default:
		log.Warn("Unknown message from_list (%s)", m.FromList)
	}
	return
}

// dealMissionTag 处理活动tag
func (s *Service) dealMissionTag(c context.Context, m *archive.VideoupMsg, a *archive.Archive) (err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			log.Error("s.dealMissionTag panic(%v)", pErr)
		}
	}()
	if m.MissionID != 0 { //消息里的mission_id是修改之前的
		addit, err := s.arc.Addit(c, m.Aid)
		if err != nil {
			log.Error("s.arc.Addit(%d) error(%v)", m.Aid, err)
		} else if addit == nil {
			log.Warn("s.arc.Addit(%d) warn(addit is nil)", m.Aid)
		} else if addit.MissionID == 0 {
			//取消活动资格，去掉活动tag
			tags, err := s.removeMissionTags(c, a)
			if err == nil {
				a.Tag = strings.Join(tags, ",")
			}
		}
	}
	return
}

// updateRecheckState 回查提交时的事件
func (s *Service) updateRecheckState(c context.Context, tp int, aid int64, state int8) (err error) {

	//修改archive_recheck的state
	if err = s.arc.UpdateRecheckState(c, tp, aid, state); err != nil {
		return
	}
	a, err := s.arc.ArchiveByAid(c, aid)
	if err != nil {
		log.Error("s.arc.ArchiveByAid error(%v)", err)
		err = nil
		return
	}
	tpStr := archive.RecheckType(tp)
	if tpStr != "" {
		s.arc.AddArchiveOper(c, aid, a.Attribute, a.TypeID, a.State, a.Round, 0, "", "已"+tpStr)
	}
	return
}

// addHotRecheck get hot archive from data api, and insert to archive_recheck table
func (s *Service) addHotRecheck() (err error) {
	var (
		c    = context.TODO()
		aids []int64
	)
	if aids, err = s.dataDao.HotArchive(c); err != nil {
		log.Error("s.addHotRecheck() s.dataDao.HotArchive() error(%v)", err)
		return
	}
	if err = s.arc.AddRecheckAids(c, archive.TypeHotRecheck, aids, true); err != nil {
		log.Error("s.addHotRecheck() s.arc.AddRecheckAids error(%v)", err)
		return
	}
	return
}

func (s *Service) addArchive(c context.Context, m *archive.VideoupMsg) (err error) {
	var (
		a                       *archive.Archive
		addit                   *archive.Addit
		tx                      *sql.Tx
		channelDiff, operRemark string
		operConts               []string
	)

	if a, err = s.arc.ArchiveByAid(c, m.Aid); err != nil || a.ID <= 0 {
		log.Error("addArchive s.arc.ArchiveByAid error(%v)/not found aid(%d)", err, m.Aid)
		return
	}
	//同步到tag服务方，以便在前台显示
	if err = s.upBindTag(c, a.Mid, m.Aid, a.Tag, a.TypeID); err != nil {
		return
	}
	if addit, err = s.arc.Addit(c, m.Aid); err != nil {
		log.Error("modifyArchive s.arc.Addit error(%v) aid(%d)", err, m.Aid)
		return
	}
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("addArchive s.arc.BeginTran error(%v)", err)
		return
	}
	//非活动的ugc稿件
	if addit == nil || (!addit.IsPGC() && addit.MissionID <= 0) {
		if channelDiff, operRemark, err = s.txAddChannelReview(c, tx, m.Aid); err != nil {
			log.Error("addArchive s.txAddChannelReview(%d) error(%v)", m.Aid, err)
			tx.Rollback()
			return
		}
		if channelDiff != "" {
			operConts = append(operConts, channelDiff)
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("addArchive tx.Commit() error(%v) aid(%d)", err, m.Aid)
		return
	}
	if len(operConts) > 0 && operRemark != "" {
		s.arc.AddArchiveOper(c, m.Aid, a.Attribute, a.TypeID, a.State, a.Round, 0, strings.Join(operConts, ","), operRemark)
	}
	return
}

func (s *Service) modifyArchive(c context.Context, m *archive.VideoupMsg) (err error) {
	var (
		a                       *archive.Archive
		addit                   *archive.Addit
		tx                      *sql.Tx
		channelDiff, operRemark string
		operConts               []string
	)

	if a, err = s.arc.ArchiveByAid(c, m.Aid); err != nil || a.ID <= 0 {
		log.Error("modifyArchive s.arc.ArchiveByAid error(%v)/not found aid(%d)", err, m.Aid)
		return
	}
	//tag修改或分区修改时，同步到tag服务方，以便在前台显示,即使失败也不影响后续
	if m.TagChange || m.ChangeTypeID {
		s.upBindTag(c, a.Mid, m.Aid, a.Tag, a.TypeID)
	}

	if addit, err = s.arc.Addit(c, m.Aid); err != nil {
		log.Error("modifyArchive s.arc.Addit error(%v) aid(%d)", err, m.Aid)
		return
	}
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("modifyArchive s.arc.BeginTran error(%v)", err)
		return
	}
	//新增视频 且 非活动的ugc稿件
	if m.AddVideos && (addit == nil || (!addit.IsPGC() && addit.MissionID <= 0)) {
		log.Info("begin to check channel review aid(%d)", m.Aid)
		if channelDiff, operRemark, err = s.txAddChannelReview(c, tx, m.Aid); err != nil {
			log.Error("modifyArchive s.txAddChannelReview(%d) error(%v)", m.Aid, err)
			tx.Rollback()
			return
		}
		if channelDiff != "" {
			operConts = append(operConts, channelDiff)
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("modifyArchive tx.Commit() error(%v) aid(%d)", err, m.Aid)
		return
	}
	if len(operConts) > 0 && operRemark != "" {
		s.arc.AddArchiveOper(c, m.Aid, a.Attribute, a.TypeID, a.State, a.Round, 0, strings.Join(operConts, ","), operRemark)
	}
	return
}

func (s *Service) autoOpen(c context.Context, m *archive.VideoupMsg) (err error) {
	var (
		a *archive.Archive
	)
	if a, err = s.arc.ArchiveByAid(c, m.Aid); err != nil || a.ID <= 0 {
		log.Error("autoOpen s.arc.ArchiveByAid error(%v)/not found aid(%d)", err, m.Aid)
		return
	}
	s.adminBindTag(c, a.Mid, a.ID, a.Tag, a.TypeID)
	return
}

func (s *Service) delayOpen(c context.Context, m *archive.VideoupMsg) (err error) {
	var (
		a *archive.Archive
	)
	if a, err = s.arc.ArchiveByAid(c, m.Aid); err != nil || a.ID <= 0 {
		log.Error("autoOpen s.arc.ArchiveByAid error(%v)/not found aid(%d)", err, m.Aid)
		return
	}
	s.adminBindTag(c, a.Mid, a.ID, a.Tag, a.TypeID)
	return
}

func (s *Service) postFirstRound(c context.Context, m *archive.VideoupMsg) (err error) {
	var (
		v *archive.Video
		a *archive.Archive
	)
	if v, err = s.arc.NewVideo(c, m.Filename); err != nil || v == nil {
		log.Error("postFirstRound s.arc.NewVideo error(%v)/not found filename(%d)", err, m.Filename)
		return
	}
	if a, err = s.arc.ArchiveByAid(c, m.Aid); err != nil || a.ID <= 0 {
		log.Error("postFirstRound s.arc.ArchiveByAid error(%v)/not found aid(%d)", err, m.Aid)
		return
	}
	if a.State == archive.StateForbidUpDelete {
		log.Warn("postFirstRound archive(%d) filename(%s) state(%d) is deleted", a.ID, v.Filename, a.State)
		return
	}

	if m.AdminChange && !s.isPGC(a.ID) {
		s.sendMail(c, a, v)
	}
	s.sendVideoPrivateEmail(c, a, v)
	return
}

func (s *Service) isPGC(aid int64) (is bool) {
	is = false
	if ad, _ := s.arc.Addit(context.TODO(), aid); ad != nil && ad.IsPGC() {
		is = true
	}
	return
}

func (s *Service) arcStateChange(nw *archive.Archive, old *archive.Archive, canOpen bool) (err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			log.Error("s.arcStateChange panic(%v)", pErr)
		}
	}()
	if nw == nil {
		return
	}

	oldValue := 0
	if old != nil {
		oldValue = isOpenReplyState(old.State)
	}
	switchVal := isOpenReplyState(nw.State) - oldValue
	//关评论
	if switchVal < 0 && !canOpen {
		s.arcReply(context.TODO(), nw, archive.ReplyOff)
	}
	//开评论
	if switchVal > 0 && canOpen {
		s.arcReply(context.TODO(), nw, archive.ReplyOn)
	}

	return
}

//removeMissionTags 删除活动tag
func (s *Service) removeMissionTags(c context.Context, a *archive.Archive) (tags []string, err error) {
	tags = strings.Split(a.Tag, ",")

	for i := 0; i < len(tags); i++ {
		if _, ok := s.missTagsCache[tags[i]]; ok {
			tags = append(tags[:i], tags[i+1:]...)
			i--
			continue
		}
	}
	tagStr := strings.Join(tags, ",")
	if err = s.adminBindTag(c, a.Mid, a.ID, tagStr, a.TypeID); err != nil {
		log.Error("removeMissionTags(%v) s.adminBindTag() error(%v)", a, err)
		return
	}
	if _, err = s.arc.UpTag(c, a.ID, tagStr); err != nil {
		log.Error("s.arc.UpTag(%d,%s) error(%v)", a.ID, tagStr, err)
		err = nil
	}
	if _, err = s.arc.AddArchiveOper(c, a.ID, a.Attribute, a.TypeID, a.State, a.Round, 0, fmt.Sprintf("[Tag]从[%s]设为[%s]", a.Tag, tagStr), "因被取消活动资格"); err != nil {
		log.Error("s.arc.AddArchiveOper() archive(%v) error(%v)", a, err)
		err = nil
	}
	return
}
