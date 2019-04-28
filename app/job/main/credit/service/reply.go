package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/job/main/credit/model"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

func (s *Service) replyAllConsumer() {
	defer s.wg.Done()
	var (
		msgs = s.replyAllSub.Messages()
		err  error
		c    = context.TODO()
	)
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.replyAllSub.Message closed")
			return
		}
		msg.Commit()
		m := &model.Reply{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}
		switch m.Action {
		case model.RouteReplyReport:
			err = s.addReplyReport(c, m)
		default:
			log.Warn("replyAllConsumer unknown message action(%s)", m.Action)
		}
		if err != nil {
			log.Error("replyMessage key(%s) value(%s) partition(%d) offset(%d) commit error(%v)", msg.Key, msg.Value, msg.Partition, msg.Offset, err)
			continue
		}
		log.Info("replyMessage key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
	}
}

func (s *Service) addReplyReport(c context.Context, m *model.Reply) (err error) {
	if m.Reply == nil || m.Subject == nil || m.Report == nil {
		log.Warn("reply content(%+v) empty!", m)
		return
	}
	if m.Report.State == model.ReportStateAddJuge {
		log.Warn("rpid(%d) state(%d) is in juge", m.Report.RPID, m.Report.State)
		return
	}
	if m.Report.Type != model.SubTypeArchive {
		log.Warn("m.Report.Type(%d) model.SubTypeArchive(%d)", m.Report.Type, model.SubTypeArchive)
		return
	}
	var ca *model.AutoCaseConf
	if ca, err = s.dao.AutoCaseConf(c, model.OriginReply); err != nil {
		log.Error("s.dao.AutoCaseConf(%d) error(%v)", model.OriginReply, err)
		return
	}
	if ca == nil {
		log.Warn("otype(%d) auto acse conf is not set !", model.OriginReply)
		return
	}
	replyReportReason := model.BlockedReasonTypeByReply(m.Report.Reason)
	if _, ok := ca.Reasons[replyReportReason]; !ok {
		log.Warn("m.Report.Reason(%d) not int ca.Reasons(%+v)", m.Report.Reason, ca.Reasons)
		return
	}
	if m.Report.Score < ca.ReportScore {
		log.Warn("m.Report.Score(%d) ca.ReportScore(%d)", m.Report.Score, ca.ReportScore)
		return
	}
	if m.Reply.Like > int64(ca.Likes) {
		log.Warn("m.Reply.Like(%d) ca.Likes(%d)", m.Reply.Like, ca.Likes)
		return
	}
	var disCount, total int64
	if disCount, err = s.dao.CountCaseMID(c, m.Reply.MID, model.OriginReply); err != nil {
		log.Error("s.dao.CountBlocked(%d) err(%v)", m.Reply.MID, err)
		return
	}
	if disCount > 0 {
		log.Warn("rpid(%d) mid(%d) report in 24 hours", m.Report.RPID, m.Reply.MID)
		return
	}
	if total, err = s.dao.CountBlocked(c, m.Reply.MID, time.Now().AddDate(-1, 0, 0)); err != nil {
		log.Error("s.dao.CountBlocked(%d) err(%v)", m.Reply.MID, err)
		return
	}
	punishResult, blockedDay := model.PunishResultDays(total)
	mc := &model.Case{
		Mid:          m.Reply.MID,
		Status:       model.CaseStatusGrantStop,
		RelationID:   fmt.Sprintf("%d-%d-%d", m.Report.RPID, m.Report.Type, m.Report.OID),
		PunishResult: punishResult,
		BlockedDay:   blockedDay,
		OPID:         model.AutoOPID,
		BCtime:       m.Reply.CTime,
	}
	mc.Origin = model.Origin{
		OriginTitle:   s.replyOriginTitle(c, m.Report.OID, m.Report.Type),
		OriginContent: m.Reply.Content.Message,
		OriginType:    int64(model.OriginReply),
		OriginURL:     s.replyOriginURL(m.Report.RPID, m.Report.OID, m.Report.Type),
		ReasonType:    int64(replyReportReason),
	}
	var count int64
	if count, err = s.dao.CaseRelationIDCount(c, model.OriginReply, mc.RelationID); err != nil {
		log.Error("ss.dao.CaseRelationIDCount(%d,%s) err(%v)", model.OriginReply, mc.RelationID, err)
		return
	}
	if count > 0 {
		log.Warn("rpid(%d) state(%d) is alreadly juge", m.Report.RPID, m.Report.State)
		return
	}
	var need bool
	if need, err = s.dao.CheckFilter(c, "credit", m.Reply.Content.Message, ""); err != nil {
		log.Error("s.dao.CheckFilter(%s,%s) error(%v)", "credit", m.Reply.Content.Message, err)
		return
	}
	if need {
		log.Warn("reply(%d) message(%s) is filter", m.Report.RPID, m.Reply.Content.Message)
		return
	}
	if err = s.dao.AddBlockedCase(c, mc); err != nil {
		log.Error("s.dao.AddBlockedCase(%+v) err(%v)", mc, err)
		return
	}
	s.dao.UpReplyState(c, m.Report.OID, m.Report.RPID, m.Report.Type, model.ReportStateAddJuge)
	s.dao.UpAppealState(c, model.AppealBusinessID, m.Report.OID, m.Report.RPID)
	return
}

func (s *Service) replyOriginTitle(c context.Context, oid int64, oType int8) (title string) {
	switch oType {
	case model.SubTypeArchive:
		arg := &archive.ArgAid2{Aid: oid}
		arc, err := s.arcRPC.Archive3(c, arg)
		if err != nil {
			log.Error("s.arcRPC.Archive3(%v) error(%v)", arg, err)
			return
		}
		title = arc.Title
		return
	}
	return
}

func (s *Service) replyOriginURL(rpid, oid int64, oType int8) string {
	switch oType {
	case model.SubTypeArchive:
		return fmt.Sprintf(model.ReplyOriginURL, oid, rpid)
	}
	return ""
}

// RegReply regist reply subject.
func (s *Service) RegReply(c context.Context, table string, nwMsg []byte, oldMsg []byte) (err error) {
	var (
		replyType int8
		replyID   int64
	)
	switch table {
	case _blockedCaseTable:
		mr := &model.Case{}
		if err = json.Unmarshal(nwMsg, mr); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
			return
		}
		replyType = model.ReplyCase
		replyID = mr.ID
	case _blockedInfoTable:
		mr := &model.BlockedInfo{}
		if err = json.Unmarshal(nwMsg, mr); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
			return
		}
		replyType = model.ReplyBlocked
		replyID = mr.ID
	case _blockedPublishTable:
		mr := &model.Publish{}
		if err = json.Unmarshal(nwMsg, mr); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
			return
		}
		if mr.PStatus != model.PublishOpen {
			return
		}
		replyType = model.ReplyPublish
		replyID = mr.ID
	}
	return s.dao.RegReply(c, replyID, replyType)
}
