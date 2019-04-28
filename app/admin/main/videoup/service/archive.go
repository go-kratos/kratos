package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/app/admin/main/videoup/model/manager"
	"go-common/app/admin/main/videoup/model/oversea"
	tagrpc "go-common/app/interface/main/tag/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/time"
	"go-common/library/xstr"
)

// Submit second_round submit,update archive_addit\archive_delay\archive.
func (s *Service) Submit(c context.Context, ap *archive.ArcParam) (err error) {
	var (
		tx                                    *sql.Tx
		a                                     *archive.Archive
		addit                                 *archive.Addit
		operConts, flowConts                  []string
		oldMissionID, missionID, descFormatID int64
		mUser                                 *manager.User
		porderConts                           []string
		porder                                = &archive.Porder{}
		rel                                   = &oversea.ArchiveRelation{}
	)
	if a, err = s.arc.Archive(c, ap.Aid); err != nil || a == nil {
		log.Error("s.arc.Archive(%d) error(%v) or a==nil", ap.Aid, err)
		return
	}
	if addit, _ = s.arc.Addit(c, ap.Aid); addit != nil {
		missionID = addit.MissionID
		descFormatID = addit.DescFormatID
	}
	forbid, _ := s.arc.Forbid(c, ap.Aid)
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	//私单稿件 审核后台允许 私单只填写 流量TAG
	if ap.GroupID > 0 {
		porderConts, porder = s.diffPorder(c, a.Aid, ap)
		operConts = append(operConts, porderConts...)
		if err = s.TxUpPorder(tx, ap); err != nil {
			tx.Rollback()
			return
		}
		//自动匹配流量tag
		ap.OnFlowID = ap.GroupID
		if err = s.txUpFlowID(tx, ap); err != nil {
			tx.Rollback()
			return
		}
	}
	s.sendPorderLog(c, ap, porderConts, porder, a)
	//流量属性控制，比如：频道禁止
	if flowConts, err = s.txBatchUpFlowsState(c, tx, a.Aid, ap.UID, ap.FlowAttribute); err != nil {
		log.Error("s.txBatchUpFlowsState error(%v)", err)
		tx.Rollback()
		return
	}
	if len(flowConts) > 0 {
		operConts = append(operConts, flowConts...)
	}

	//审核更换流量TAG 或者新增私单 私单二期 flow_design.pool=2 聚合到 archive_forbid + archive.attr 由前端merge
	var delayCont string
	if ap.State, delayCont, err = s.txUpArcDelay(c, tx, ap.Aid, ap.Mid, ap.State, ap.Delay, ap.DTime); err != nil {
		tx.Rollback()
		log.Error("s.arc.txUpArcDelay(%d,%d,%d,%d,%v,%d) error(%v)", ap.Aid, ap.Mid, ap.TypeID, ap.State, ap.Delay, ap.DTime, err)
		return
	}
	if delayCont != "" {
		operConts = append(operConts, delayCont)
	}
	ap.Round = s.archiveRound(c, a, ap.Aid, a.Mid, a.TypeID, a.Round, ap.State, ap.CanCelMission)
	if ap.Access, err = s.txUpArcMainState(c, tx, ap.Aid, ap.Forward, ap.TypeID, ap.Access, ap.State, ap.Round, ap.RejectReason); err != nil {
		tx.Rollback()
		log.Error("s.arc.txUpArcMainState(%d,%d,%d,%d,%d,%d,%s) error(%v)", ap.Aid, ap.Forward, ap.TypeID, ap.Access, ap.State, ap.Round, ap.RejectReason, err)
		return
	}
	ap.PTime = s.archivePtime(c, ap.Aid, ap.State, ap.PTime)
	// access、cancel_mission、cover、reject_reason、attr、source、redirecturl、forward、state、pubtime、copyright、mtime、round、title、typeid、content、delay
	if _, err = s.arc.TxUpArchive(tx, ap.Aid, ap.Title, ap.Content, ap.Cover, ap.Note, ap.Copyright, ap.PTime); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpArchive(%d,%s,%s,%s,%s,%d,%d) error(%v)", ap.Aid, ap.Title, ap.Content, ap.Cover, ap.Note, ap.Copyright, ap.PTime, err)
		return
	}
	log.Info("aid(%d) update archive title(%s) content(%s) cover(%s) copyright(%d) ,ptime(%d), round(%d), state(%d)", ap.Aid, ap.Title, ap.Content, ap.Cover, ap.Copyright, ap.PTime, ap.Round, ap.State)
	// cancel activity
	if ap.CanCelMission {
		oldMissionID = missionID
		missionID = 0
	}
	desc := ""
	if descFormatID > 0 {
		desc = ap.Content
	}
	if _, err = s.arc.TxUpAddit(tx, ap.Aid, missionID, ap.Source, desc, ap.Dynamic); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpAddit(%d,%d,%s,%s,%s) error(%v)", ap.Aid, missionID, ap.Source, desc, ap.Dynamic, err)
		return
	}
	if ap.PolicyID > 1 {
		ap.Attrs.LimitArea = 1
	} else {
		ap.Attrs.LimitArea = 0
	}
	attrs, forbidAttrs := s.archiveAttr(c, ap, true)
	var attrConts []string
	if attrConts, err = s.TxUpArchiveAttr(c, tx, a, ap.Aid, ap.UID, attrs, forbidAttrs, ap.URL); err != nil {
		tx.Rollback()
		log.Error("s.TxUpArchiveAttr(%d) error(%v)", ap.Aid, err)
		return
	}
	operConts = append(operConts, attrConts...)
	log.Info("aid(%d) update archive attribute(%+v)", ap.Aid, attrs)
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	if ap.PolicyID != 0 {
		if rel, err = s.oversea.PolicyRelation(c, ap.Aid); err != nil {
			log.Error("s.oversea.PolicyRelation(%d) error(%v)", ap.Aid, err)
			return
		}
		if rel == nil || rel.GroupID != ap.PolicyID {
			if _, err = s.oversea.UpPolicyRelation(c, ap.Aid, ap.PolicyID); err != nil {
				log.Error("s.oversea.UpPolicyRelation(%d,%d) error(%v)", ap.Aid, ap.PolicyID, err)
				return
			}
			operConts = append(operConts, fmt.Sprintf("[地区展示]应用策略组ID[%d]", ap.PolicyID))
		}
	}
	log.Info("aid(%d) end second_round submit pro", ap.Aid)
	// is send email
	var isChanged = true
	if (((ap.State == archive.StateOpen && ap.Access == archive.AccessDefault) || ap.State == archive.StateForbidUserDelay) && !ap.AdminChange) ||
		(ap.State == archive.StateForbidWait) {
		isChanged = false
	}
	archiveOperConts, changeTypeID, changeCopyright, changeTitle, changeCover := s.diffArchiveOper(ap, a, addit, forbid)

	operConts = append(operConts, archiveOperConts...)
	oper := &archive.ArcOper{Aid: ap.Aid, UID: ap.UID, TypeID: ap.TypeID, State: archive.AccessState(ap.State, ap.Access), Round: ap.Round, Attribute: a.Attribute, Remark: ap.Note}
	oper.Content = strings.Join(operConts, "，")
	if ap.ApplyUID != 0 {
		mUser, _ = s.mng.User(c, ap.ApplyUID)
		oper.Content = "[通过" + mUser.NickName + "(" + strconv.FormatInt(ap.ApplyUID, 10) + ")申请的工单]" + oper.Content
	}
	s.addArchiveOper(c, oper)
	// databus
	s.busSecondRound(ap.Aid, oldMissionID, ap.Notify, isChanged, changeTypeID, changeCopyright, changeTitle, changeCover, ap.FromList, ap)
	go s.busSecondRoundUpCredit(ap.Aid, 0, ap.Mid, ap.UID, ap.State, ap.Round, ap.ReasonID, ap.RejectReason)
	// log
	s.sendArchiveLog(c, ap, operConts, a)
	return
}

// BatchArchive  batch async archive audit.
func (s *Service) BatchArchive(c context.Context, aps []*archive.ArcParam, action string) (err error) {
	var mp = &archive.MultSyncParam{}
	var ok bool
	for _, ap := range aps {
		mp.Action = action
		mp.ArcParam = ap
		if ok, err = s.busCache.PushMultSync(c, mp); err != nil {
			log.Error("s.busCache.PushMultSync(ap(%+v)) error(%v)", ap, err)
			return
		}
		if !ok {
			log.Warn("s.busCache.PushMultSync(ap(%+v))", ap)
			continue
		}
	}
	return
}

func (s *Service) dealArchive(c context.Context, ap *archive.ArcParam) (err error) {
	var a *archive.Archive
	if a, err = s.arc.Archive(c, ap.Aid); err != nil || a == nil {
		log.Error("s.arc.Archive(%d) error(%v) or a==nil", ap.Aid, err)
		return
	}
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	log.Info("aid(%d) start BatchSubmit pro params is (%+v)", ap.Aid, ap)
	if ap.State, _, err = s.txUpArcDelay(c, tx, ap.Aid, ap.Mid, ap.State, true, ap.DTime); err != nil {
		tx.Rollback()
		log.Error("s.arc.txUpArcDelay(%d,%d,%d,%d,%v,%d) error(%v)", ap.Aid, ap.Mid, ap.TypeID, ap.State, true, ap.DTime, err)
		return
	}
	ap.Round = s.archiveRound(c, a, ap.Aid, a.Mid, a.TypeID, a.Round, ap.State, ap.CanCelMission)
	if ap.Access, err = s.txUpArcMainState(c, tx, ap.Aid, ap.Forward, ap.TypeID, ap.Access, ap.State, ap.Round, ap.RejectReason); err != nil {
		tx.Rollback()
		log.Error("s.arc.txUpArcMainState(%d,%d,%d,%d,%d,%d,%s) error(%v)", ap.Aid, ap.Forward, ap.TypeID, ap.Access, ap.State, ap.Round, ap.RejectReason, err)
		return
	}
	ap.PTime = s.archivePtime(c, ap.Aid, ap.State, a.PTime) // NOTE: for batch no ptime...
	if _, err = s.arc.TxUpArcPTime(tx, ap.Aid, ap.PTime); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpArcPTime(%d,%d) error(%v)", ap.Aid, ap.PTime, err)
		return
	}
	log.Info("aid(%d) update archive pubtime(%d)", ap.Aid, ap.PTime)
	if _, err = s.arc.TxUpArcNote(tx, ap.Aid, ap.Note); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpArcNote(%d,%s) error(%v)", ap.Aid, ap.Note, err)
		return
	}
	log.Info("aid(%d) update archive note(%s)", ap.Aid, ap.Note)
	if ap.FlagCopyright {
		if _, err = s.arc.TxUpArcCopyRight(tx, ap.Aid, ap.Copyright); err != nil {
			tx.Rollback()
			log.Error("s.arc.TxUpArcCopyRight(%d,%d) error(%v)", ap.Aid, ap.Copyright, err)
			return
		}
		log.Info("aid(%d) update archive Copyright(%d)", ap.Aid, ap.Copyright)
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	log.Info("aid(%d) end BatchSubmit pro", ap.Aid)
	var isChanged = true
	if (((ap.State == archive.StateOpen && ap.Access == archive.AccessDefault) || ap.State == archive.StateForbidUserDelay) && !ap.AdminChange) ||
		(ap.State == archive.StateForbidWait) {
		isChanged = false
	}
	oper := &archive.ArcOper{Aid: ap.Aid, UID: ap.UID, TypeID: ap.TypeID, State: archive.AccessState(ap.State, ap.Access), Round: ap.Round, Attribute: a.Attribute, Remark: ap.Note}
	operConts := s.diffBatchArchiveOper(ap, a)
	oper.Content = strings.Join(operConts, "，")
	s.addArchiveOper(c, oper)
	s.busSecondRound(ap.Aid, 0, ap.Notify, isChanged, false, false, false, false, ap.FromList, ap)
	go s.busSecondRoundUpCredit(ap.Aid, 0, ap.Mid, ap.UID, ap.State, ap.Round, ap.ReasonID, ap.RejectReason)
	// log
	ap.CTime = a.CTime
	s.sendArchiveLog(c, ap, operConts, a)
	return
}

func (s *Service) dealAttrs(c context.Context, ap *archive.ArcParam) (err error) {
	var a *archive.Archive
	if a, err = s.arc.Archive(c, ap.Aid); err != nil || a == nil {
		log.Error("s.arc.Archive(%d) error(%v) or a==nil", ap.Aid, err)
		return
	}
	//付费稿件不支持批量属性修改 因为涉及到价格设置一致性问题
	if s.isUGCPay(a) {
		log.Info("dealAttrs skip UGCPay(%d) error(%v)", ap.Aid, err)
		return
	}
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	log.Info("aid(%d) start BatchAttr pro", ap.Aid)
	var operConts, flowConts []string
	attrs, forbidAttrs := s.archiveAttr(c, ap, false)
	if operConts, err = s.TxUpArchiveAttr(c, tx, a, ap.Aid, ap.UID, attrs, forbidAttrs, ""); err != nil {
		tx.Rollback()
		log.Error("s.arc.txUpArcAttrs(%d) error(%v)", ap.Aid, err)
		return
	}
	log.Info("aid(%d) update archive attribute( %+v )", ap.Aid, attrs)
	if _, err = s.arc.TxUpArcNote(tx, ap.Aid, ap.Note); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpArcNote(%d,%s) error(%v)", ap.Aid, ap.Note, err)
		return
	}
	log.Info("aid(%d) update archive note(%s)", ap.Aid, ap.Note)
	//流量属性控制，比如：频道禁止
	if flowConts, err = s.txBatchUpFlowsState(c, tx, a.Aid, ap.UID, ap.FlowAttribute); err != nil {
		log.Error("s.txBatchUpFlowsState error(%v)", err)
		tx.Rollback()
		return
	}
	if len(flowConts) > 0 {
		operConts = append(operConts, flowConts...)
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	log.Info("aid(%d) end BatchAttr pro", ap.Aid)
	oper := &archive.ArcOper{Aid: ap.Aid, UID: ap.UID, TypeID: a.TypeID, State: archive.AccessState(a.State, a.Access), Round: a.Round, Attribute: a.Attribute, Remark: ap.Note}
	oper.Content = strings.Join(operConts, "，")
	s.addArchiveOper(c, oper)
	s.busSecondRound(ap.Aid, 0, ap.Notify, ap.AdminChange, false, false, false, false, ap.FromList, ap)
	// log
	ap.CTime = a.CTime
	s.sendArchiveLog(c, ap, operConts, a)
	return
}

func (s *Service) dealArchiveSecondRound(c context.Context, ap *archive.ArcParam) (err error) {
	s.busSecondRound(ap.Aid, 0, false, false, false, false, false, false, ap.FromList, ap)
	return
}

func (s *Service) dealTypeID(c context.Context, ap *archive.ArcParam) (err error) {
	var a *archive.Archive
	if a, err = s.arc.Archive(c, ap.Aid); err != nil || a == nil {
		log.Error("s.arc.Archive(%d) error(%v) or a==nil", ap.Aid, err)
		return
	}
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	log.Info("aid(%d) start BatchType pro", ap.Aid)
	if _, err = s.arc.TxUpArcTypeID(tx, ap.Aid, ap.TypeID); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpArcTypeID(%d,%d) error(%v)", ap.Aid, ap.TypeID, err)
		return
	}
	log.Info("aid(%d) update archive type_id(%d)", ap.Aid, ap.TypeID)
	if _, err = s.arc.TxUpArcNote(tx, ap.Aid, ap.Note); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpArcNote(%d,%s) error(%v)", ap.Aid, ap.Note, err)
		return
	}
	log.Info("aid(%d) update archive note(%s)", ap.Aid, ap.Note)
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	log.Info("aid(%d) end BatchType pro", ap.Aid)
	var changeTypeID bool
	oper := &archive.ArcOper{Aid: ap.Aid, UID: ap.UID, TypeID: ap.TypeID, State: archive.AccessState(a.State, a.Access), Round: a.Round, Attribute: a.Attribute, Remark: ap.Note}
	oper.Content, changeTypeID = s.diffTypeID(ap.TypeID, a.TypeID, a.State)
	var operConts []string
	operConts = append(operConts, oper.Content)
	s.addArchiveOper(c, oper)
	if ap.ForceSync {
		s.busArchiveForceSync(ap.Aid)
	} else {
		s.busSecondRound(ap.Aid, 0, ap.Notify, ap.AdminChange, changeTypeID, false, false, false, ap.FromList, ap)
	}
	// log
	ap.CTime = a.CTime
	s.sendArchiveLog(c, ap, operConts, a)
	return
}

// // BatchZlimit batche modify zlimit.
// func (s *Service) BatchZlimit(c context.Context, ap *archive.ArcParam) (err error) {
// 	var (
// 		tx   *sql.Tx
// 		aps  []*archive.ArcParam
// 		aids = xstr.JoinInts(ap.Aids)
// 	)
// 	if tx, err = s.mng.BeginTran(c); err != nil {
// 		log.Error("s.arc.BeginTran() error(%v)", err)
// 		return
// 	}
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 			log.Error("wocao jingran recover le error(%v)", r)
// 		}
// 	}()
// 	log.Info("aids(%s) start BatchZlimit pro", aids)
// 	for _, aid := range ap.Aids {
// 		if _, err = s.mng.TxAddUpArea(tx, aid, ap.GroupID); err != nil {
// 			tx.Rollback()
// 			log.Error("s.arc.TxAddUpArea(%d,%d) error(%v)", aid, ap.GroupID, err)
// 			return
// 		}
// 		log.Info("aid(%d) update archive gid(%d)", aid, ap.GroupID)
// 		ap.Aid = aid
// 		aps = append(aps, ap)
// 	}
// 	if err = tx.Commit(); err != nil {
// 		log.Error("tx.Commit() error(%v)", err)
// 		return
// 	}
// 	log.Info("aids(%s) end BatchZlimit pro", aids)
// 	s.busSecondRound(aps)
// 	return
// }

// UpAuther update owner.
func (s *Service) UpAuther(c context.Context, ap *archive.ArcParam) (err error) {
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	log.Info("aid(%d) start up author", ap.Aid)
	if _, err = s.arc.TxUpArcAuthor(tx, ap.Aid, ap.Mid, ap.Author); err != nil {
		tx.Rollback()
		log.Error("s.Auther s.dede.TxUpArcAuthor(%d,%d,%s) error(%v)", ap.Aid, ap.Mid, ap.Author, err)
		return
	}
	log.Info("aid(%d) update archive mid(%d) author(%s)", ap.Aid, ap.Mid, ap.Author)
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	log.Info("aid(%d) end up author", ap.Aid)
	// databus
	s.busSecondRound(ap.Aid, 0, ap.Notify, false, false, false, false, false, ap.FromList, ap)
	return
}

// UpAccess update access.
func (s *Service) UpAccess(c context.Context, ap *archive.ArcParam) (err error) {
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	log.Info("aid(%d) start up access", ap.Aid)
	if _, err = s.arc.TxUpArcAccess(tx, ap.Aid, ap.Access); err != nil {
		tx.Rollback()
		log.Error("s.Access s.dede.TxUpArcAccess(%d,%d) error(%v)", ap.Aid, ap.Access, err)
		return
	}
	log.Info("aid(%d) update archive access(%d)", ap.Aid, ap.Access)
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	log.Info("aid(%d) end up access", ap.Aid)
	// databus
	s.busSecondRound(ap.Aid, 0, ap.Notify, ap.AdminChange, false, false, false, false, ap.FromList, ap)
	return
}

// UpArchiveAttr update archive attr by aid.
func (s *Service) UpArchiveAttr(c context.Context, aid, uid int64, attrs map[uint]int32, forbidAttrs map[string]map[uint]int32, redirectURL string) (err error) {
	var a *archive.Archive
	if a, err = s.arc.Archive(c, aid); err != nil || a == nil {
		log.Error("s.arc.Archive(%d) error(%v) or a==nil", aid, err)
		return
	}
	log.Info("aid(%d) begin tran change attribute", aid)
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("begin tran error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	var conts []string
	if conts, err = s.txUpArcAttrs(tx, a, attrs, redirectURL); err != nil {
		tx.Rollback()
		return
	}
	var tmpCs []string
	if tmpCs, err = s.txUpArcForbidAttrs(c, tx, a, forbidAttrs); err != nil {
		tx.Rollback()
		return
	}
	conts = append(conts, tmpCs...)
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	log.Info("aid(%d) end tran change attribute", aid)
	if len(conts) > 0 {
		s.arc.AddArcOper(c, a.Aid, uid, a.Attribute, a.TypeID, int16(a.State), a.Round, 1, strings.Join(conts, "，"), "")
	}
	// NOTE: send modify_archive for sync dede.
	s.busModifyArchive(aid, false, false)
	return
}

// UpArcDtime update archive dtime by aid.
func (s *Service) UpArcDtime(c context.Context, aid int64, dtime time.Time) (err error) {
	if delay, _ := s.arc.Delay(c, aid, archive.DelayTypeForUser); delay == nil {
		err = ecode.NothingFound
		return
	}
	log.Info("aid(%d) dtime(%d)  begin tran change archive delaytime", aid, dtime)
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("begin tran error(%v)", err)
		return
	}
	if _, err = s.arc.TxUpDelayDtime(tx, aid, archive.DelayTypeForUser, dtime); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpDelayDtime() error(%v)", err)
		return
	}
	if _, err = s.arc.TxUpArcMtime(tx, aid); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpArcMtime() error(%v)", err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	log.Info("aid(%d) dtime(%d) end tran change archive delaytime", aid, dtime)
	return
}

/**
 * archive map
 * return map[int64]*archive.ChannelReviewInfo, error/nil
 */
func (s *Service) checkChannelReview(c context.Context, arcs map[int64]*archive.Archive) (res map[int64]*archive.ChannelReviewInfo, err error) {
	var reviewRes map[int64]*tagrpc.ResChannelCheckBack
	res = map[int64]*archive.ChannelReviewInfo{}
	aids := []int64{}

	//检查是否开放浏览
	for _, a := range arcs {
		if a == nil {
			continue
		}

		if archive.NormalState(a.State) {
			aids = append(aids, a.Aid)
		}
		res[a.Aid] = &archive.ChannelReviewInfo{AID: a.Aid}
	}
	if len(aids) <= 0 {
		return
	}

	//检查是否有回查数据
	if _, aids, err = s.arc.RecheckIDByAID(c, archive.TypeChannelRecheck, aids); err != nil || len(aids) <= 0 {
		return
	}
	//实时查询是否变更了频道，频道是否需要回查
	if reviewRes, err = s.tag.CheckChannelReview(c, aids); err != nil {
		log.Error("checkChannelReview s.arc.CheckChannelReview error(%v) aids(%v)", err, aids)
		err = nil
	}
	for _, aid := range aids {
		res[aid].CanOperRecheck = true
		if reviewRes != nil && reviewRes[aid] != nil {
			res[aid].NeedReview = reviewRes[aid].CheckBack == 1
			cids := []int64{}
			for chid := range reviewRes[aid].Channels {
				cids = append(cids, chid)
			}
			res[aid].ChannelIDs = xstr.JoinInts(cids)
		}
	}
	return
}

//UpArcTag 保存tag
func (s *Service) UpArcTag(c context.Context, uid int64, pm *archive.TagParam) (err error) {
	var (
		arc                     *archive.Archive
		checkRes                map[int64]*archive.ChannelReviewInfo
		needReview, fromChannel bool
		canOperRecheck          bool
	)
	//archive check
	if arc, err = s.arc.Archive(c, pm.AID); err != nil {
		log.Error("UpArcTag s.arc.Archive(%d) error(%v) params(%+v)", pm.AID, err, pm)
		return
	}
	if arc == nil {
		err = ecode.NothingFound
		return
	}
	//check whether need channel review: channel changes every 3h, used to notice
	fromChannel = strings.TrimSpace(pm.FromChannelReview) == "1"
	if fromChannel {
		if checkRes, err = s.checkChannelReview(c, map[int64]*archive.Archive{pm.AID: arc}); err != nil {
			log.Error("UpArcTag s.checkChannelReview(%d) error(%v) params(%+v)", pm.AID, err, pm)
			return
		}
		if checkRes != nil && checkRes[pm.AID] != nil {
			needReview = checkRes[pm.AID].NeedReview
			canOperRecheck = checkRes[pm.AID].CanOperRecheck
		}
	}
	if err = s.saveTag(c, uid, arc, pm.Tags, "", canOperRecheck, nil); err != nil {
		log.Error("UpArcTag s.saveTag error(%v) uid(%d) canoperrecheck(%v) params(%+v)", err, uid, canOperRecheck, pm)
		return
	}
	if fromChannel && !needReview {
		err = ecode.VideoupChannelReviewNotTrigger
	}
	return
}

//saveTag 更新tag，可能触发频道回查
func (s *Service) saveTag(c context.Context, uid int64, arc *archive.Archive, tags, note string, canOperRecheck bool, ap *archive.ArcParam) (err error) {
	var (
		tx               *sql.Tx
		tagChange        bool
		remark, flowDiff string
		content          []string
	)
	//compare tag
	aid := arc.Aid
	arc.Tag = strings.Join(Slice2String(SliceUnique(Slice2Interface(strings.Split(arc.Tag, ",")))), ",")
	tags = strings.Join(Slice2String(SliceUnique(Slice2Interface(strings.Split(tags, ",")))), ",")
	tagChange = arc.Tag != tags

	//update db by transaction
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("saveTag s.arc.BeginTran error(%v)", err)
		return
	}
	if tagChange {
		if _, err = s.arc.TxUpTag(tx, aid, tags); err != nil {
			log.Error("saveTag s.arc.TxUpTag(%d,%s) error(%v)", aid, tags, err)
			tx.Rollback()
			return
		}
		content = append(content, archive.Operformat(archive.OperTypeTag, arc.Tag, tags, archive.OperStyleOne))
	}

	if canOperRecheck {
		if err = s.arc.TxUpRecheckState(tx, archive.TypeChannelRecheck, aid, archive.RecheckStateDone); err != nil {
			log.Error("saveTag s.arc.TxUpRecheckState(%d) error(%v)", aid, err)
			tx.Rollback()
			return
		}
		if _, flowDiff, err = s.txAddOrUpdateFlowState(c, tx, aid, archive.FlowGroupNoChannel, uid, archive.PoolArcForbid, archive.FlowDelete, "频道回查取消频道禁止"); err != nil {
			log.Error("saveTag s.txAddOrUpdateFlowState(%d,%d) error(%v)", aid, uid, err)
			tx.Rollback()
			return
		}
		if flowDiff != "" {
			content = append(content, flowDiff)
		}
		remark = "频道已回查"
	}
	if err = tx.Commit(); err != nil {
		log.Error("saveTag tx.Commit error(%v)", err)
		return
	}
	log.Info("saveTag aid(%d) origin old-tags(%s) new-tags(%s) canoperrecheck(%v)", aid, arc.Tag, tags, canOperRecheck)

	//log tag change/channel review history
	remark = remark + note
	if len(content) > 0 || remark != "" {
		oper := &archive.ArcOper{
			Aid:       aid,
			UID:       uid,
			TypeID:    arc.TypeID,
			State:     archive.AccessState(arc.State, arc.Access),
			Round:     arc.Round,
			Content:   strings.Join(content, ","),
			Attribute: arc.Attribute,
			Remark:    remark,
		}
		s.arc.AddArcOper(c, oper.Aid, oper.UID, oper.Attribute, oper.TypeID, oper.State, oper.Round, 0, oper.Content, oper.Remark)
	}

	//同步tag服务
	if tagChange {
		if ap != nil && ap.IsUpBind {
			err = s.upBind(c, aid, arc.Mid, tags, arc.TypeID)
		} else {
			err = s.adminBind(c, aid, arc.Mid, tags, arc.TypeID)
		}
	}
	//仅同步隐藏tag
	if !tagChange && ap != nil && ap.SyncHiddenTag {
		err = s.upBind(c, aid, arc.Mid, arc.Tag, arc.TypeID)
	}

	return
}

//BatchUpTag batch update archive tag
func (s *Service) BatchUpTag(c context.Context, uid int64, pm *archive.BatchTagParam) (warning string, err error) {
	var (
		arcList         map[int64]*archive.Archive
		reviewResp      map[int64]*archive.ChannelReviewInfo
		ok, fromChannel bool
		warningID       []int64
		originTag       string
	)
	//get archive list
	if arcList, err = s.arc.Archives(c, pm.AIDs); err != nil {
		log.Error("batchTag vdaSvc.GetArchives error(%v) params(%+v)", err, pm)
		return
	}
	fromChannel = strings.TrimSpace(pm.FromList) == archive.FromListChannelReview
	//batch check whether need channel review, if not mark aid to noneedaids
	if fromChannel {
		if reviewResp, err = s.checkChannelReview(c, arcList); err != nil {
			log.Error("batchTag s.checkChannelReview error(%v) params(%+v)", err, pm)
			return
		}
	}

	//save each archive tag
	mp := &archive.MultSyncParam{
		ArcParam: &archive.ArcParam{
			UID:           uid,
			Note:          pm.Note,
			FromList:      pm.FromList,
			IsUpBind:      pm.IsUpBind,
			SyncHiddenTag: pm.SyncHiddenTag,
		},
	}
	for id, arc := range arcList {
		if fromChannel && (reviewResp == nil || reviewResp[id] == nil || !reviewResp[id].NeedReview) {
			warningID = append(warningID, id)
		}
		originTag = arc.Tag
		//tag change
		if pm.Action != "" {
			arc.Tag = StringHandler(arc.Tag, pm.Tags, ",", pm.Action == "delete")
		}
		//批量修改tag
		mp.Action = archive.ActionArchiveTag
		//频道回查逻辑
		if reviewResp != nil && reviewResp[id] != nil && reviewResp[id].CanOperRecheck {
			mp.Action = archive.ActionArchiveTagRecheck
		}
		//非频道回查的tag未变更/频道回查tag未变更且无回查记录，提前返回
		if originTag == arc.Tag && ((!fromChannel && !pm.SyncHiddenTag) || (fromChannel && mp.Action == archive.ActionArchiveTag)) {
			continue
		}
		mp.ArcParam.Aid = id
		mp.ArcParam.Tag = arc.Tag

		log.Info("BatchUpTag begin to pushmultsync action(%s) arcparam(%+v)", mp.Action, mp.ArcParam)
		if ok, err = s.busCache.PushMultSync(c, mp); err != nil {
			log.Error("BatchUpTag s.busCache.PushMultSync(%d,%s,%s,%v) error(%v)", id, arc.Tag, pm.Note, pm.FromList, err)
			return
		}
		if !ok {
			log.Warn("BatchUpTag s.busCache.PushMultSync(%d,%s,%s,%v)", id, arc.Tag, pm.Note, pm.FromList)
			continue
		}
	}

	if len(warningID) > 0 {
		warning = fmt.Sprintf("稿件 %s 不需要频道回查", xstr.JoinInts(warningID))
	}
	return
}

func (s *Service) dealTag(c context.Context, canOperRecheck bool, ap *archive.ArcParam) (err error) {
	var arc *archive.Archive
	if arc, err = s.arc.Archive(c, ap.Aid); err != nil {
		log.Error(" UpArcTag s.arc.Archive(%d) error(%v)", ap.Aid, err)
		return
	}
	if arc == nil {
		err = ecode.NothingFound
		return
	}

	if err = s.saveTag(c, ap.UID, arc, ap.Tag, ap.Note, canOperRecheck, ap); err != nil {
		log.Error("dealTag s.saveTag(%d,%d,%s,%s,%v) error(%v)", ap.UID, ap.Aid, ap.Tag, ap.Note, canOperRecheck, err)
	}
	return
}

func (s *Service) adminBind(c context.Context, aid, mid int64, tags string, typeID int16) (err error) {
	log.Info("before sync tag for aid(%d) tags(%s)", aid, tags)
	typeName := ""
	if tp, ok := s.typeCache[typeID]; ok && tp != nil {
		typeName = tp.Name
		if tp, err = s.TypeTopParent(typeID); err != nil {
			log.Error("adminBind s.TypeTopParent(%d) error(%v) aid(%d)", typeID, err, aid)
			err = nil
		} else if tp != nil {
			typeName = fmt.Sprintf("%s,%s", typeName, tp.Name)
		}
	}
	if err = s.tag.AdminBind(c, aid, mid, tags, typeName, ""); err != nil {
		log.Error("adminBind sync tag error(%v) aid(%+v) tags(%s) typename(%s)", err, aid, tags, typeName)
		return
	}
	log.Info("end sync tag for aid(%d) tags(%s) successfully", aid, tags)
	return
}

func (s *Service) upBind(c context.Context, aid, mid int64, tags string, typeID int16) (err error) {
	log.Info("upBind before sync tag for aid(%d) tags(%s)", aid, tags)
	typeName := ""
	if tp, ok := s.typeCache[typeID]; ok && tp != nil {
		typeName = tp.Name
		if tp, err = s.TypeTopParent(typeID); err != nil {
			log.Error("upBind s.TypeTopParent(%d) error(%v) aid(%d)", typeID, err, aid)
			err = nil
		} else if tp != nil {
			typeName = fmt.Sprintf("%s,%s", typeName, tp.Name)
		}
	}
	if err = s.tag.UpBind(c, aid, mid, tags, typeName, ""); err != nil {
		log.Error("upBind sync tag error(%v) aid(%+v) tags(%s) typename(%s)", err, aid, tags, typeName)
		return
	}
	log.Info("upBind end sync tag for aid(%d) tags(%s) successfully", aid, tags)
	return
}

//GetChannelInfo get channel info & hit_rules & need review
func (s *Service) GetChannelInfo(c context.Context, aids []int64) (info map[int64]*archive.ChannelInfo, err error) {
	info = make(map[int64]*archive.ChannelInfo, len(aids))
	for _, aid := range aids {
		info[aid] = &archive.ChannelInfo{}
	}

	res, err := s.tag.CheckChannelReview(c, aids)
	if err != nil || res == nil {
		log.Error("GetChannelInfo s.tag.GetChannelInfo error(%v)/resp=nil aids(%v)", err, aids)
		return
	}
	if len(res) <= 0 {
		return
	}

	for aid := range info {
		if res[aid] == nil {
			continue
		}
		chs := []*archive.Channel{}
		for _, ch := range res[aid].Channels {
			if ch == nil {
				continue
			}

			chs = append(chs, &archive.Channel{
				TID:         ch.Tid,
				Tname:       ch.TName,
				HitRules:    ch.HitRules,
				HitTagNames: ch.HitTNames,
			})
		}
		info[aid].Channels = chs
		info[aid].CheckBack = res[aid].CheckBack
	}
	return
}

// AITrack .
func (s *Service) AITrack(c context.Context, aid []int64) (aids string, err error) {
	return s.data.ArchiveRelated(c, aid)
}

//ChannelNamesByAids .
func (s *Service) ChannelNamesByAids(c context.Context, aids []int64) (aidMap map[int64][]string) {
	var (
		size       = len(aids)
		maxSize    = 1000
		sliceLimit = 50
		cnt        = 0
		aidSli     = [][]int64{}
		sli        = []int64{}
		grp        = errgroup.Group{}
		chlist     chan map[int64][]string
	)
	if size > maxSize {
		size = maxSize
	}
	aidMap = make(map[int64][]string, size)

	//去重分组
	for _, aid := range aids {
		if _, exist := aidMap[aid]; exist {
			continue
		}
		aidMap[aid] = []string{}

		if cnt >= sliceLimit {
			cnt = 0
			aidSli = append(aidSli, sli)
			sli = []int64{}
		}
		cnt++
		sli = append(sli, aid)
	}
	if cnt > 0 {
		aidSli = append(aidSli, sli)
	}

	chlist = make(chan map[int64][]string, size)

	//batch get channel names
	for i, aids := range aidSli {
		if i >= maxSize {
			break
		}
		aidtmp := aids
		grp.Go(func() error {
			resp, err := s.tag.CheckChannelReview(context.TODO(), aidtmp)
			if err != nil {
				log.Error("ChannelNamesByAids s.tag.CheckChannelReview error(%v) aids(%v)", err, aidtmp)
				return nil
			}

			tnames := map[int64][]string{}
			for aid, rp := range resp {
				if rp == nil && rp.Channels == nil {
					continue
				}
				if _, exist := tnames[aid]; !exist {
					tnames[aid] = []string{}
				}
				for _, ch := range rp.Channels {
					if ch == nil {
						continue
					}

					tnames[aid] = append(tnames[aid], ch.TName)
				}
			}
			chlist <- tnames
			return nil
		})
	}
	grp.Wait()
	close(chlist)

	for tnames := range chlist {
		for aid, tname := range tnames {
			aidMap[aid] = tname
		}
	}
	return
}
