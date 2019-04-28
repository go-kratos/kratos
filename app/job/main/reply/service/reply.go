package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-common/library/database/elastic"
	"net/url"
	"regexp"
	"strings"
	"time"

	"go-common/app/job/main/reply/conf"
	"go-common/app/job/main/reply/model/reply"
	model "go-common/app/job/main/reply/model/reply"
	accmdl "go-common/app/service/main/account/api"
	assmdl "go-common/app/service/main/assist/model/assist"
	relmdl "go-common/app/service/main/relation/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
)

var (
	_atReg    = regexp.MustCompile(`@([^\s^:^,^@]+)`)
	_topicReg = regexp.MustCompile(`#([^\n^@^#^\x{1F000}-\x{1F02F}^\x{1F0A0}-\x{1F0FF}^\x{1F100}-\x{1F64F}^\x{1F680}-\x{1F6FF}^\x{1F910}-\x{1F96B}^\x{1F980}-\x{1F9E0}]{1,32})#`)
	_urlReg   = regexp.MustCompile(`(((http:\/\/|https:\/\/)[a-z0-9A-Z]+\.(bilibili|biligame)\.com[a-z0-9A-Z\/\.\$\*\?~=#!%@&-]*)|((http:\/\/|https:\/\/)(acg|b23)\.tv[a-z0-9A-Z\/\.\$\*\?~=#!@&]*))`)
	_avReg    = regexp.MustCompile(`#(cv\d+)|#(av\d+)|#(vc\d+)`)

	searchHTTPClient        *xhttp.Client
	errReplyContentNotFound = errors.New("reply content not found")
)

const (
	_appIDReply  = "reply"
	_appIDReport = "replyreport"
	timeFormat   = "2006-01-02 15:03:04"

	// event
	_eventReply      = "reply"
	_eventHate       = "hate"
	_eventLike       = "like"
	_eventLikeCancel = "like_cancel"
	_eventHateCancel = "hate_cancel"
)

func (s *Service) beginTran(c context.Context) (*xsql.Tx, error) {
	return s.dao.BeginTran(c)
}

func (s *Service) actionAdd(c context.Context, msg *consumerMsg) {
	var rp *model.Reply
	if err := json.Unmarshal([]byte(msg.Data), &rp); err != nil {
		log.Error("json.Unmarshal() error(%v)", err)
		return
	}
	if rp.RpID == 0 || rp.Oid == 0 || rp.Content == nil {
		log.Error("The structure of reply(%s) from rpCh was wrong", msg.Data)
		return
	}
	if rp.Root == 0 && rp.Parent == 0 {
		s.addReply(c, rp)
	} else {
		s.addReplyReply(c, rp)
	}
}

func (s *Service) tranAdd(c context.Context, rp *model.Reply, is bool) (err error) {
	tx, err := s.beginTran(c)
	if err != nil {
		log.Error("reply(%s) beginTran error(%v)", rp, err)
		return
	}
	var rows int64
	defer func() {
		if err == nil && rows == 0 {
			err = errors.New("sql: transaction add reply failed")
		}
	}()
	if is {
		if rp.IsNormal() {
			rows, err = s.dao.Subject.TxIncrCount(tx, rp.Oid, rp.Type, rp.CTime.Time())
			if err != nil || rows == 0 {
				tx.Rollback()
				log.Error("dao.Subject.TxIncrCount(%v) error(%v) or rows==0", rp, err)
				return
			}
		} else {
			rows, err = s.dao.Subject.TxIncrFCount(tx, rp.Oid, rp.Type, rp.CTime.Time())
			if err != nil || rows == 0 {
				tx.Rollback()
				log.Error("dao.Subject.TxIncrCount(%v) error(%v) or rows==0", rp, err)
				return
			}
		}
	} else {
		var rootReply *model.Reply
		if rootReply, err = s.dao.Reply.GetForUpdate(tx, rp.Oid, rp.Root); err != nil {
			tx.Rollback()
			return err
		}
		if rootReply.IsDeleted() {
			return fmt.Errorf("the root reply is deleted(%d,%d,%d)", rp.Oid, rp.Type, rp.Root)
		}
		if rp.IsNormal() {
			rows, err = s.dao.Reply.TxIncrCount(tx, rp.Oid, rp.Root, rp.CTime.Time())
			if err != nil || rows == 0 {
				tx.Rollback()
				log.Error("dao.Reply.TxIncrCount(%v) error(%v) or rows==0", rp, err)
				return
			}
			rows, err = s.dao.Subject.TxIncrACount(tx, rp.Oid, rp.Type, 1, rp.CTime.Time())
			if err != nil || rows == 0 {
				tx.Rollback()
				log.Error("dao.Subject.TxIncrACount(%v) error(%v) or rows==0", rp, err)
				return
			}
		} else {
			rows, err = s.dao.Reply.TxIncrFCount(tx, rp.Oid, rp.Root, rp.CTime.Time())
			if err != nil || rows == 0 {
				tx.Rollback()
				log.Error("dao.Reply.TxIncrCount(%v) error(%v) or rows==0", rp, err)
				return
			}
		}
	}
	if rp.State == model.ReplyStateAudit || rp.State == model.ReplyStateMonitor {
		if rows, err = s.dao.Subject.TxIncrMCount(tx, rp.Oid, rp.Type, rp.CTime.Time()); err != nil || rows == 0 {
			tx.Rollback()
			log.Error("dao.Subject.TxIncrMCount(%v) error(%v) or rows==0", rp, err)
			return
		}
	}
	rows, err = s.dao.Content.TxInsert(tx, rp.Oid, rp.Content)
	if err != nil || rows == 0 {
		tx.Rollback()
		log.Error("dao.Content.TxInContent(%v) error(%v) or rows==0", rp, err)
		return
	}
	rows, err = s.dao.Reply.TxInsert(tx, rp)
	if err != nil || rows == 0 {
		tx.Rollback()
		log.Error("dao.Reply.TxInReply(%v) error(%v) or rows==0", rp, err)
		return
	}
	return tx.Commit()
}

func (s *Service) regTopic(c context.Context, msg string) (topics []string) {
	msg = _urlReg.ReplaceAllString(msg, "")
	msg = _avReg.ReplaceAllString(msg, "#")

	ss := _topicReg.FindAllStringSubmatch(msg, -1)
	if len(ss) == 0 {
		return
	}
	for _, nns := range ss {
		if len(nns) == 2 {
			topic := strings.TrimSpace(nns[1])
			if len(topic) > 0 {
				topics = append(topics, topic)
			}
		}
		if len(topics) >= 5 {
			break
		}
	}
	return
}

func (s *Service) regAt(c context.Context, msg string, over, self int64) (ats []int64) {
	var err error
	ss := _atReg.FindAllStringSubmatch(msg, 10)
	if len(ss) == 0 {
		return
	}
	names := make([]string, 0, len(ss))
	for _, nns := range ss {
		if len(nns) == 2 {
			names = append(names, nns[1])
		}
	}
	if len(names) == 0 {
		return
	}
	us, err := s.accSrv.InfosByName3(c, &accmdl.NamesReq{Names: names})
	if err != nil {
		log.Error("s.accSrv.InfosByName2 failed, err(%v)", err)
		return
	}
	ats = make([]int64, 0, len(us.Infos))
	for mid := range us.Infos {
		if mid != over && mid != self {
			ats = append(ats, mid)
		}
	}
	if len(ats) == 0 {
		return
	}
	ats = s.getFilterBlacklist(c, self, ats)
	return
}

func (s *Service) addReply(c context.Context, rp *model.Reply) {
	var (
		err error
		ok  bool
	)
	sub, err := s.getSubject(c, rp.Oid, rp.Type)
	if err != nil || sub == nil {
		log.Error("s.getSubject failed , oid(%d,%d) err(%v)", rp.Oid, rp.Type, err)
		return
	}
	if sub == nil {
		log.Error("get subject is nil oid(%d) type(%d)", rp.Oid, rp.Type)
		return
	}
	// init some field
	if rp.IsNormal() {
		sub.RCount = sub.RCount + 1
		sub.ACount = sub.ACount + 1
	}
	sub.Count = sub.Count + 1
	rp.Floor = sub.Count
	rp.MTime = rp.CTime
	rp.Content.RpID = rp.RpID
	rp.Content.CTime = rp.CTime
	rp.Content.MTime = rp.MTime
	if len(rp.Content.Ats) == 0 {
		rp.Content.Ats = s.regAt(c, rp.Content.Message, 0, rp.Mid)
	}
	rp.Content.Topics = s.regTopic(c, rp.Content.Message)
	// begin transaction
	if err = s.tranAdd(c, rp, true); err != nil {
		log.Error("Transaction add reply(%v) error(%v)", rp, err)
		return
	}
	// add cache
	if err = s.dao.Mc.AddSubject(c, sub); err != nil {
		log.Error("s.dao.Mc.AddSubject failed , oid(%d) err(%v)", sub.Oid, err)
	}
	if err = s.dao.Mc.AddReply(c, rp); err != nil {
		log.Error("s.dao.Mc.AddReply failed , RpID(%d) err(%v)", rp.RpID, err)
	}
	if rp.IsNormal() {
		// update reply count
		s.upAcount(c, sub.Oid, sub.Type, sub.ACount, rp.CTime.Time())
		// add index cache
		if ok, err = s.dao.Redis.ExpireIndex(c, sub.Oid, sub.Type, model.SortByFloor); err == nil && ok {
			if err = s.dao.Redis.AddFloorIndex(c, sub.Oid, sub.Type, rp); err != nil {
				log.Error("s.dao.Redis.AddFloorIndex failed , oid(%d) type(%d) err(%v)", sub.Oid, sub.Type, err)
			}
		}
		if ok, err = s.dao.Redis.ExpireIndex(c, sub.Oid, sub.Type, model.SortByCount); err == nil && ok {
			if err = s.dao.Redis.AddCountIndex(c, sub.Oid, sub.Type, rp); err != nil {
				log.Error("s.dao.Redis.AddCountIndex failed , oid(%d) type(%d) err(%v)", sub.Oid, sub.Type, err)
			}
		}
		if ok, err = s.dao.Redis.ExpireIndex(c, sub.Oid, sub.Type, model.SortByLike); err == nil && ok {
			rpts := make(map[int64]*reply.Report, 1)
			if rpt, _ := s.dao.Report.Get(c, rp.Oid, rp.RpID); rpt != nil {
				rpts[rp.RpID] = rpt
			}
			if err = s.dao.Redis.AddLikeIndex(c, sub.Oid, sub.Type, rpts, rp); err != nil {
				log.Error("s.dao.Redis.AddLikeIndex failed , oid(%d) type(%d) err(%v)", sub.Oid, sub.Type, err)
			}
		}
		s.notifyReply(c, sub, rp)
	} else if rp.State == model.ReplyStateAudit {
		if err = s.dao.Redis.AddAuditIndex(c, rp); err != nil {
			log.Error("s.dao.Redis.AddAUditIndex(%d,%d,%d) error(%v)", rp.Oid, rp.RpID, rp.Type, err)
		}
	}
	if err = s.dao.PubEvent(c, _eventReply, rp.Mid, sub, rp, nil); err != nil {
		return
	}
}

func (s *Service) addReplyReply(c context.Context, rp *model.Reply) {
	var (
		err error
		ok  bool
	)
	sub, err := s.getSubject(c, rp.Oid, rp.Type)
	if err != nil || sub == nil {
		log.Error("s.getSubject failed , oid(%d,%d) err(%v)", rp.Oid, rp.Type, err)
		return
	}
	// NOTE:depend on db,do not get from cache
	rootRp, err := s.getReply(c, rp.Oid, rp.Root)
	if err != nil {
		log.Error("s.getReply failed , oid(%d), root(%d) err(%v)", rp.Oid, rp.Root, err)
		return
	}
	if rootRp == nil {
		log.Error("get reply is nil oid(%d) type(%d) rpid(%d)", rp.Oid, rp.Type, rp.Root)
		return
	}
	var parentRp *model.Reply
	if rp.Root != rp.Parent {
		parentRp, err = s.getReply(c, rp.Oid, rp.Parent)
		if err != nil {
			log.Error("s.getReply failed , oid(%d), parent(%d) err(%v)", rp.Oid, rp.Parent, err)
			return
		}
		if parentRp == nil {
			log.Error("get reply is nil oid(%d) type(%d) rpid(%d)", rp.Oid, rp.Type, rp.Parent)
			return
		}
		if parentRp.Dialog == 0 {
			log.Warn("Dialog Need Migration oid(%d) type(%d) rootID(%d)", rp.Oid, rp.Type, rootRp.RpID)
			// s.setDialogByRoot(context.Background(), rp.Oid, rp.Type, rp.Root)
		}
		rp.Dialog = parentRp.Dialog
	} else {
		parentRp = rootRp
		if rp.Dialog != rp.RpID {
			rp.Dialog = rp.RpID
		}
	}
	// init some field
	if rp.IsNormal() {
		sub.ACount = sub.ACount + 1
		rootRp.RCount = rootRp.RCount + 1
	}
	rootRp.Count = rootRp.Count + 1
	rootRp.MTime = rp.CTime
	rp.Floor = rootRp.Count
	rp.MTime = rp.CTime
	rp.Content.RpID = rp.RpID
	rp.Content.CTime = rp.CTime
	rp.Content.MTime = rp.MTime
	if len(rp.Content.Ats) == 0 {
		rp.Content.Ats = s.regAt(c, rp.Content.Message, 0, rp.Mid)
	}
	rp.Content.Topics = s.regTopic(c, rp.Content.Message)
	// begin transaction
	if err = s.tranAdd(c, rp, false); err != nil {
		log.Error("Transaction add reply(%v) error(%v)", rp, err)
		return
	}
	// add cache
	if err = s.dao.Mc.AddSubject(c, sub); err != nil {
		log.Error("s.dao.Mc.AddSubject failed , oid(%d), err(%v)", sub.Oid, err)
	}
	if err = s.dao.Mc.AddReply(c, rp); err != nil {
		log.Error("s.dao.Mc.AddReply failed , RpID(%d), err(%v)", rp.RpID, err)
	}
	if err = s.dao.Mc.AddReply(c, rootRp); err != nil {
		log.Error("s.dao.Mc.AddReply failed , RpID(%d), err(%v)", rootRp.RpID, err)
	}
	if rootRp.IsTop() {
		if err = s.dao.Mc.AddTop(c, rootRp); err != nil {
			log.Error("s.dao.Mc.AddReply failed , RpID(%d), err(%v)", rootRp.RpID, err)
		}
	} else if rootRp.IsNormal() {
		if ok, err = s.dao.Redis.ExpireIndex(c, rp.Oid, rp.Type, model.SortByCount); err == nil && ok {
			s.dao.Redis.AddCountIndex(c, rp.Oid, rp.Type, rootRp)
		}
	}
	if rp.IsNormal() {
		// update reply count
		s.upAcount(c, sub.Oid, sub.Type, sub.ACount, rp.CTime.Time())
		// add index cache
		if ok, err = s.dao.Redis.ExpireNewChildIndex(c, rootRp.RpID); err == nil && ok {
			if err = s.dao.Redis.AddNewChildIndex(c, rootRp.RpID, rp); err != nil {
				log.Error("s.dao.Redis.AddFloorIndexByRoot failed , RpID(%d), err(%v)", rootRp.RpID, err)
			}
		}
		// add dialog cache
		if rp.Dialog != 0 {
			if ok, err = s.dao.Redis.ExpireDialogIndex(c, rp.Dialog); err == nil && ok {
				rps := []*model.Reply{rp}
				if err = s.dao.Redis.AddDialogIndex(c, rp.Dialog, rps); err != nil {
					log.Error("s.dao.Redis.AddDialogIndex failed , RpID(%d), Dialog(%d), Floor(%d) err(%v)", rp.RpID, rp.Dialog, rp.Floor, err)
				}
			}
		}
		s.notifyReplyReply(c, sub, rootRp, parentRp, rp)
	} else if rp.State == model.ReplyStateAudit {
		if err = s.dao.Redis.AddAuditIndex(c, rp); err != nil {
			log.Error("s.dao.Redis.AddAUditIndex(%d,%d,%d) error(%v)", rp.Oid, rp.RpID, rp.Type, err)
		}
	}
	if err = s.dao.PubEvent(c, _eventReply, rp.Mid, sub, rp, nil); err != nil {
		return
	}
}

func (s *Service) addTopCache(c context.Context, msg *consumerMsg) {
	var (
		err error
		sub *model.Subject
		rp  *model.Reply
	)
	var d struct {
		Oid int64  `json:"oid"`
		Tp  int8   `json:"tp"`
		Top uint32 `json:"top"`
	}
	if err = json.Unmarshal([]byte(msg.Data), &d); err != nil {
		log.Error("json.Unmarshal() error(%v)", err)
		return
	}
	if rp, err = s.dao.Mc.GetTop(c, d.Oid, d.Tp, d.Top); err != nil {
		log.Error("s.dao.Mc.GetTop(oid %v,top %v) err(%v)", d.Oid, d.Top, err)
		return
	} else if rp == nil {
		if rp, err = s.dao.Reply.GetTop(c, d.Oid, d.Tp, d.Top); err != nil || rp == nil {
			log.Error("s.dao.Reply.GetTop(%d, %d) error(%v)", d.Oid, d.Tp, err)
			return
		}
		if rp.Content, err = s.dao.Content.Get(c, d.Oid, rp.RpID); err != nil {
			return
		}
		s.dao.Mc.AddTop(c, rp)
		sub, err = s.dao.Subject.Get(c, d.Oid, d.Tp)
		if err != nil {
			log.Error("s.dao.Subject.Get(%d, %d) error(%v)", d.Oid, d.Tp, err)
			return
		}
		err = sub.TopSet(rp.RpID, d.Top, 1)
		if err != nil {
			return
		}
		_, err = s.dao.Subject.UpMeta(c, d.Oid, d.Tp, sub.Meta, time.Now())
		if err != nil {
			log.Error("s.dao.Subject.UpMeta(%d,%d,%d) failed!err:=%v ", rp.RpID, rp.Oid, d.Tp, err)
			return
		}
		s.dao.Mc.AddSubject(c, sub)
	}
}

func (s *Service) actionRpt(c context.Context, msg *consumerMsg) {
	var (
		err error
		ok  bool
	)
	var d struct {
		Oid  int64 `json:"oid"`
		RpID int64 `json:"rpid"`
		Tp   int8  `json:"tp"`
	}
	if err = json.Unmarshal([]byte(msg.Data), &d); err != nil {
		log.Error("json.Unmarshal() error(%v)", err)
		return
	}
	rp, err := s.getReplyCache(c, d.Oid, d.RpID)
	if err != nil {
		log.Error("s.getReply failed , oid(%d), RpID(%d) err(%v)", d.Oid, d.RpID, err)
		return
	}
	if rp == nil {
		return
	}
	sub, err := s.getSubject(c, d.Oid, d.Tp)
	if err != nil || sub == nil {
		log.Error("s.getSubject failed , oid(%d),tp(%d), RpID(%d) err(%v)", d.Oid, d.Tp, d.RpID, err)
		return
	}
	// update like index
	if rp.Root == 0 && rp.Parent == 0 && !rp.IsDeleted() {
		if ok, err = s.dao.Redis.ExpireIndex(c, d.Oid, rp.Type, model.SortByLike); err == nil && ok {
			rpts := make(map[int64]*reply.Report, 1)
			if rpt, _ := s.dao.Report.Get(c, rp.Oid, rp.RpID); rpt != nil {
				rpts[rp.RpID] = rpt
			}
			if err = s.dao.Redis.AddLikeIndex(c, d.Oid, rp.Type, rpts, rp); err != nil {
				log.Error("s.dao.Redis.AddLikeIndex(%d, %d) error(%v)", d.Oid, rp.Type, err)
			}
		}
	}
	report, err := s.dao.Report.Get(c, d.Oid, d.RpID)
	if err != nil || report == nil {
		log.Error("dao.Report.GetReport(%d, %d) met error (%v)", rp.Oid, rp.RpID, err)
		return
	}
	if err = s.dao.PubEvent(c, _eventReportAdd, report.Mid, sub, rp, report); err != nil {
		return
	}
}

func (s *Service) setLike(c context.Context, cmsg *StatMsg) {
	var (
		event string
	)
	rp, err := s.getReply(c, cmsg.Oid, cmsg.ID)
	if err != nil || rp == nil || rp.Content == nil {
		log.Error("s.getReply(%d, %d) reply:%+v error(%v)", cmsg.Oid, cmsg.ID, rp, err)
		return
	}
	sub, err := s.getSubject(c, rp.Oid, rp.Type)
	if err != nil || sub == nil {
		log.Error("s.getSubject failed , oid(%d) type(%d) err(%v)", rp.Oid, rp.Type, err)
		return
	}
	_, err = s.dao.Reply.UpLike(c, cmsg.Oid, cmsg.ID, cmsg.Count, cmsg.DislikeCount, time.Now())
	if err != nil {
		log.Error("s.dao.Reply.UpLike (%v) failed!err:=%v", cmsg, err)
		return
	}
	if cmsg.Count > rp.Like {
		event = _eventLike
		var max int
		if max, err = s.dao.Redis.MaxLikeCnt(c, rp.RpID); err == nil && cmsg.Count > max {
			if err = s.dao.Redis.SetMaxLikeCnt(c, rp.RpID, int64(cmsg.Count)); err == nil {
				rp.Like = cmsg.Count
				rp.Hate = cmsg.DislikeCount
				s.notifyLike(c, cmsg.Mid, rp)
			}
		}
	} else if cmsg.DislikeCount > rp.Hate {
		event = _eventHate
	} else if cmsg.Count < rp.Like {
		event = _eventLikeCancel
	} else {
		event = _eventHateCancel
	}
	rp.Like = cmsg.Count
	rp.Hate = cmsg.DislikeCount
	s.dao.Mc.AddReply(c, rp)
	if rp.AttrVal(model.ReplyAttrAdminTop) == 1 || rp.AttrVal(model.ReplyAttrUpperTop) == 1 {
		s.dao.Mc.AddTop(c, rp)
		return
	}
	// if have root, then update root's index
	if rp.Root == 0 && rp.IsNormal() {
		var ok bool
		if ok, err = s.dao.Redis.ExpireIndex(c, rp.Oid, rp.Type, model.SortByLike); err == nil && ok {
			rpts := make(map[int64]*reply.Report, 1)
			if rpt, _ := s.dao.Report.Get(c, rp.Oid, rp.RpID); rpt != nil {
				rpts[rp.RpID] = rpt
			}
			if err = s.dao.Redis.AddLikeIndex(c, rp.Oid, rp.Type, rpts, rp); err != nil {
				log.Error("s.dao.Redis.AddLikeIndex(%d, %d) error(%v)", rp.Oid, rp.Type, err)
			}
		}
	}
	if err = s.dao.PubEvent(c, event, cmsg.Mid, sub, rp, nil); err != nil {
		return
	}
}

func (s *Service) adminLog(c context.Context, rp *model.Reply, adid int64, isreport, state int8, result, remark string) {
	// admin log
	s.dao.Admin.UpIsNotNew(c, rp.RpID, time.Now())
	s.dao.Admin.Insert(c, adid, rp.Oid, rp.RpID, rp.Type, result, remark, model.AdminIsNew, isreport, state, time.Now())
}

// getSubject get reply subject from  mysql  .
// NOTE : note get from mc,count must depend on mysql
func (s *Service) getSubject(c context.Context, oid int64, tp int8) (sub *model.Subject, err error) {
	if sub, err = s.dao.Subject.Get(c, oid, tp); err != nil {
		log.Error("dao.Subject.Get(%d, %d) error(%v)", oid, tp, err)
	}
	return
}

func (s *Service) getReply(c context.Context, oid, RpID int64) (rp *model.Reply, err error) {
	if rp, err = s.dao.Reply.Get(c, oid, RpID); err != nil {
		log.Error("s.dao.Reply.Get(%d, %d) error(%v)", oid, RpID, err)
		return
	} else if rp == nil {
		return
	}
	if rp.Content, err = s.dao.Content.Get(c, rp.Oid, rp.RpID); err != nil {
		log.Error("s.dao.Content.Get(%d,%d) error(%v)", rp.Oid, rp.RpID, err)
	} else if rp.Content == nil {
		err = errReplyContentNotFound
	}
	return
}

func (s *Service) getReplyCache(c context.Context, oid, RpID int64) (rp *model.Reply, err error) {
	if rp, err = s.dao.Mc.GetReply(c, RpID); err != nil {
		log.Error("replyCacheDao.GetReply(%d, %d) error(%v)", oid, RpID, err)
	}
	if rp != nil {
		return
	}
	if rp, err = s.dao.Reply.Get(c, oid, RpID); err != nil {
		log.Error("dao.Reply.GetReply(%d, %d) error(%v)", oid, RpID, err)
	}
	if rp != nil {
		rp.Content, _ = s.dao.Content.Get(c, rp.Oid, rp.RpID)
		// NOTE  not add member info to cache
	}
	return
}

func (s *Service) upAcount(c context.Context, oid int64, tp int8, count int, now time.Time) {
	s.statDao.Send(c, tp, oid, count)
}

func (s *Service) callSearchUp(c context.Context, res map[string]*searchFlush) (err error) {
	var (
		rps  []*searchFlush
		rpts []*searchFlush
	)
	for _, r := range res {
		if r.Report != nil {
			rpts = append(rpts, r)
		} else {
			rps = append(rps, r)
		}
	}
	if len(rps) > 0 {
		err = s.callSearch(c, rps, false)
	}
	if len(rpts) > 0 {
		err = s.callSearch(c, rpts, true)
	}
	return
}

// callSearch update reply or report info to ES search.
func (s *Service) callSearch(c context.Context, params []*searchFlush, isRpt bool) (err error) {
	var (
		b      []byte
		ms     []map[string]interface{}
		p      = url.Values{}
		urlStr = conf.Conf.Host.Search + "/api/reply/internal/update"
		res    struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}
	)
	// 更新搜索ES数据字段
	if isRpt {
		// report
		p.Set("appid", _appIDReport)
		for _, p := range params {
			m := make(map[string]interface{})
			m["id"] = fmt.Sprintf("%d_%d_%d", p.Reply.RpID, p.Reply.Oid, p.Reply.Type)
			m["reply_state"] = fmt.Sprintf("%d", p.Reply.State)
			m["reason"] = fmt.Sprintf("%d", p.Report.Reason)
			m["content"] = p.Report.Content
			m["state"] = fmt.Sprintf("%d", p.Report.State)
			m["mtime"] = p.Report.MTime.Time().Format(timeFormat)
			m["index_time"] = p.Report.CTime.Time().Format(timeFormat)
			if p.Report.Attr == 1 {
				m["attr"] = []int{1}
			} else {
				m["attr"] = []int{}
			}
			ms = append(ms, m)
		}
		if b, err = json.Marshal(ms); err != nil {
			log.Error("json.Marshal(%v) error(%v)", ms, err)
			return
		}
		p.Set("val", string(b))
		if err = searchHTTPClient.Post(c, urlStr, "", p, &res); err != nil {
			log.Error("xhttp.Post(%s) failed error(%v)", urlStr+"?"+p.Encode(), err)
		}
		log.Info("updateSearch: %s post:%s ret:%v", urlStr, p.Encode(), res)
	} else {
		// reply
		var rps = make(map[int64]*model.Reply)
		for _, p := range params {
			rps[p.Reply.RpID] = p.Reply
		}
		err = s.UpSearchReply(c, rps)
	}

	return
}

// UpSearchReply update search reply index.
func (s *Service) UpSearchReply(c context.Context, rps map[int64]*model.Reply) (err error) {
	if len(rps) <= 0 {
		return
	}
	stales := s.es.NewUpdate("reply_list")
	for _, rp := range rps {
		m := make(map[string]interface{})
		m["id"] = rp.RpID
		m["state"] = rp.State
		m["mtime"] = rp.MTime.Time().Format("2006-01-02 15:04:05")
		m["oid"] = rp.Oid
		m["type"] = rp.Type
		if rp.Content != nil {
			m["message"] = rp.Content.Message
		}
		stales = stales.AddData(s.es.NewUpdate("reply_list").IndexByTime("reply_list", elastic.IndexTypeWeek, rp.CTime.Time()), m)
	}
	err = stales.Do(c)
	if err != nil {
		log.Error("upSearchReply update stales(%s) failed!err:=%v", stales.Params(), err)
		return
	}
	log.Info("upSearchReply:stale:%s ret:%+v", stales.Params(), err)
	return
}

// getBlackListRelation check if the source user blacklisted the target user
func (s *Service) getBlackListRelation(c context.Context, srcID, targetID int64) (rel bool) {
	relMap, err := s.accSrv.RichRelations3(c, &accmdl.RichRelationReq{Owner: srcID, Mids: []int64{targetID}, RealIp: ""})
	if err != nil {
		log.Error("s.acc.RichRelations2 sourceId(%v) targetId(%v)error(%v)", srcID, targetID, err)
		err = nil
		return false
	}
	if len(relMap.RichRelations) == 0 {
		return false
	}
	if rel, ok := relMap.RichRelations[targetID]; ok && relmdl.Attr(uint32(rel)) == relmdl.AttrBlack {
		return true
	}
	return false
}

// getFilterBlacklist filters the user list that the mid user can notify message for
func (s *Service) getFilterBlacklist(c context.Context, mid int64, targetIds []int64) (filterIds []int64) {
	filterIds = make([]int64, 0, len(targetIds))
	for _, tmp := range targetIds {
		if !s.getBlackListRelation(c, tmp, mid) {
			filterIds = append(filterIds, tmp)
		}

	}
	return
}

func (s *Service) addAssistLog(c context.Context, mid, uid, subjectID, typeID, action int64, objectID, content string) (err error) {
	if len(content) > 50 {
		content = substr2(content, 0, 50) + "..."
	}
	arg := &assmdl.ArgAssistLogAdd{
		Mid:       mid,
		AssistMid: uid,
		Type:      1,
		Action:    1,
		SubjectID: subjectID,
		ObjectID:  objectID,
		Detail:    content,
		RealIP:    "",
	}
	if err = s.assistSrv.AssistLogAdd(c, arg); err != nil {
		log.Error("s.assistSrv.Assist(%d, %d, %d, %d, %d) error(%v)", mid, uid, subjectID, typeID, action, err)
	}
	return
}

func substr2(str string, start int, subLength int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		start = 0
	}

	if subLength < 0 || subLength > length {
		subLength = length
	}

	return string(rs[start:subLength])
}
