package service

import (
	"context"
	"fmt"
	"time"

	model "go-common/app/interface/main/credit/model"
	accmdl "go-common/app/service/main/account/api"
	filgrpc "go-common/app/service/main/filter/api/grpc/v1"
	memmdl "go-common/app/service/main/member/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

var (
	_emptyBlockedCase = []*model.BlockedCase{}
)

// Apply user apply jury.
func (s *Service) Apply(c context.Context, mid int64) (err error) {
	var (
		num       int64
		res       *accmdl.CardReply
		idfStatus *memmdl.RealnameStatus
		arg       = &accmdl.MidReq{Mid: mid}
	)
	if res, err = s.accountClient.Card3(c, arg); err != nil {
		err = errors.Wrap(err, "s.accountClient.Card2")
		return
	}
	if res.Card.Level <= 3 {
		err = ecode.CreditLevelLimit
		return
	}
	if res.Card.Silence == 1 {
		err = ecode.CreditIsBlock
		return
	}
	if idfStatus, err = s.memRPC.RealnameStatus(c, &memmdl.ArgMemberMid{Mid: mid, RemoteIP: metadata.String(c, metadata.RemoteIP)}); err != nil {
		return
	}
	if *idfStatus != memmdl.RealnameStatusTrue {
		err = ecode.CreditIsVerify
		return
	}
	if num, err = s.dao.BlockTotalTime(c, mid, time.Now().AddDate(0, 0, -90)); err != nil {
		err = ecode.CreditIsBlock
		return
	}
	if num != 0 {
		err = ecode.CreditIsBlock
		return
	}
	juryinfo, err := s.JuryInfoCache(c, mid)
	if err != nil {
		return
	}
	if juryinfo.Status == model.JuryStatusEffect {
		err = ecode.CreditNoApply
		return
	}
	if juryinfo.Black == model.JuryBlack {
		err = ecode.CreditBlack
		return
	}
	if err = s.dao.JuryApply(c, mid, time.Now().AddDate(0, 0, model.JuryExpiredDays)); err != nil {
		return
	}
	s.addCache(func() {
		s.dao.DelJuryInfoCache(context.TODO(), mid)
		s.dao.SendSysMsg(context.TODO(), mid, model.ApplyJuryTitle, fmt.Sprintf(model.ApplyJuryContext, model.JuryExpiredDays))
	})
	return
}

// Requirement user status in apply jury.
func (s *Service) Requirement(c context.Context, mid int64) (jr *model.JuryRequirement, err error) {
	var (
		num       int64
		card      *accmdl.CardReply
		idfStatus *memmdl.RealnameStatus
		argMId    = &accmdl.MidReq{Mid: mid}
	)
	jr = &model.JuryRequirement{}
	if card, err = s.accountClient.Card3(c, argMId); err != nil {
		return
	}
	jr.Level = card.Card.Level > 3
	jr.Blocked = card.Card.Silence == 1
	if idfStatus, err = s.memRPC.RealnameStatus(c, &memmdl.ArgMemberMid{Mid: mid, RemoteIP: metadata.String(c, metadata.RemoteIP)}); err != nil {
		return
	}
	jr.Cert = *idfStatus == memmdl.RealnameStatusTrue
	if num, err = s.dao.BlockTotalTime(c, mid, time.Now().AddDate(0, 0, -90)); err != nil {
		return
	}
	jr.Rule = num == 0
	return
}

// Jury jury user info.
func (s *Service) Jury(c context.Context, mid int64) (ui *model.UserInfo, err error) {
	var (
		info     *accmdl.InfoReply
		juryinfo *model.BlockedJury
		argMId   = &accmdl.MidReq{
			Mid: mid,
		}
	)
	if juryinfo, err = s.JuryInfoCache(c, mid); err != nil {
		return
	}
	if juryinfo.MID == 0 {
		err = ecode.CreditNotJury
		return
	}
	ui = &model.UserInfo{}
	ui.Status = juryinfo.Status
	if juryinfo.VoteTotal < 3 {
		ui.RightRadio = 50
	} else {
		ui.RightRadio = juryinfo.VoteRight * 100 / juryinfo.VoteTotal
	}
	if info, err = s.accountClient.Info3(c, argMId); err != nil {
		log.Error("s.accountClient.Info3(%d) error(%v)", mid, err)
		return
	}
	ui.Uname = info.Info.Name
	ui.Face = info.Info.Face
	if juryinfo.Status == 1 {
		delta := int64(juryinfo.Expired) - time.Now().Unix()
		if delta <= 0 {
			ui.RestDays = 0
		} else {
			ui.RestDays = delta / model.OneDaySecond
			if delta%model.OneDaySecond != 0 {
				ui.RestDays = ui.RestDays + 1
			}
		}
	}
	ui.CaseTotal = juryinfo.CaseTotal
	return
}

// CaseObtain jury user obtain case list.
func (s *Service) CaseObtain(c context.Context, mid int64, pubCid int64) (cid int64, err error) {
	juryInfo, err := s.JuryInfoCache(c, mid)
	if err != nil {
		return
	}
	if juryInfo.MID == 0 {
		err = ecode.CreditNotJury
		return
	}
	if juryInfo.Status != model.JuryStatusEffect || int64(juryInfo.Expired) < time.Now().Unix() {
		err = ecode.CreditJuryExpired
		return
	}
	if juryInfo.VoteTotal > 3 && (float64(juryInfo.VoteRight)/float64(juryInfo.VoteTotal)*100 < float64(s.c.Judge.JuryRatio)) {
		err = ecode.CreditUnderVoteRate
		return
	}
	if cid, err = s.caseVoteID(c, mid, pubCid); err != nil {
		return
	}
	if cid == 0 {
		err = ecode.CreditNoCase
		return
	}
	return
}

// Vote jury user vote case.
func (s *Service) Vote(c context.Context, mid, cid int64, attr, vote, aType, aReason int8, oc string, likes, hates []int64) (err error) {
	if vote < model.VoteBanned && vote > model.VoteDel {
		err = ecode.CreditVoteNotExist
		return
	}
	var r *model.BlockedJury
	if r, err = s.JuryInfoCache(c, mid); err != nil {
		return
	}
	if r.MID == 0 {
		err = ecode.CreditNotJury
		return
	}
	var ca *model.BlockedCase
	if ca, err = s.CaseInfoCache(c, cid); err != nil {
		return
	}
	if ca.ID == 0 {
		err = ecode.CreditCaseNotExist
		return
	}
	if ca.JudgeType != model.JudgeTypeUndeal || ca.Status == model.CaseStatusDealed || ca.Status == model.CaseStatusUndealed {
		err = ecode.CreditNovote
		return
	}
	var id int64
	if id, err = s.dao.IsVote(c, mid, cid); err != nil {
		return
	}
	if id == 0 {
		err = ecode.CreditNovote
		return
	}
	var content string
	if content, err = s.filterContent(c, oc, mid); err != nil {
		return
	}
	if err = s.setVote(c, mid, cid, id, content, attr, vote); err != nil {
		return
	}
	if len(likes) > 0 {
		s.dao.AddLikes(c, likes)
	}
	if len(hates) > 0 {
		s.dao.AddHates(c, hates)
	}
	if aType != 0 && model.ReasonTypeDesc(aReason) != "" {
		s.dao.AddCaseReasonApply(c, mid, cid, aType, ca.ReasonType, aReason)
	}
	var rate int8
	if rate, err = s.dao.NewKPI(c, mid); err != nil {
		return
	}
	if s.getVoteField(vote) != "" {
		if err = s.dao.AddCaseVoteTotal(c, s.getVoteField(vote), cid, s.getVoteNum(rate)); err != nil {
			return
		}
	}
	if r.CaseTotal >= model.GuardMedalPointC && r.GantMedalID() != model.GuardMedalNone {
		s.addCache(func() {
			for i := 0; i <= 5; i++ {
				if err = s.dao.SendMedal(context.Background(), mid, r.GantMedalID()); err != nil {
					log.Error("s.dao.SendMedal(mid:%d medalid:%d) err(%v)", mid, r.GantMedalID(), err)
					continue
				}
				break
			}
		})
	}
	return
}

func (s *Service) filterContent(c context.Context, oc string, mid int64) (content string, err error) {
	var res *filgrpc.FilterReply
	if oc == "" {
		return "", nil
	}
	if res, err = s.fliClient.Filter(c, &filgrpc.FilterReq{Area: "reply", Message: oc}); err != nil {
		return
	}
	if res.Level > 0 {
		content = ""
	} else {
		content = res.Result
	}
	log.Info("fiter mid(%d):oc(%s):content(%s):res(%v)", mid, oc, content, res)
	return
}

func (s *Service) getVoteField(vote int8) string {
	switch vote {
	case 1:
		return "vote_break"
	case 2:
		return "vote_rule"
	case 4:
		return "vote_delete"
	}
	return ""
}

func (s *Service) getVoteNum(rate int8) int8 {
	switch rate {
	case model.KPILevelS:
		return s.c.Judge.VoteNum.RateS
	case model.KPILevelA:
		return s.c.Judge.VoteNum.RateA
	case model.KPILevelB:
		return s.c.Judge.VoteNum.RateB
	case model.KPILevelC:
		return s.c.Judge.VoteNum.RateC
	case model.KPILevelD:
		return s.c.Judge.VoteNum.RateD
	default:
		return 1
	}
}

func (s *Service) setVote(c context.Context, mid, cid, id int64, content string, attr, vote int8) (err error) {
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("BeginTran err(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	_, err = s.dao.SetVoteTx(tx, mid, cid, vote)
	if err != nil {
		log.Error("tx.SetVote err(%v)", err)
		return
	}
	if len(content) > 0 {
		var (
			preContent string
			state      int8
		)
		if preContent, err = s.dao.OpContentMid(c, mid); err != nil {
			log.Error("s.dao.OpContentMid(%d) err(%v)", mid, err)
			return
		}
		if preContent == content {
			state = model.OpinionStateNoOK
		}
		_, err = s.dao.AddOpinionTx(tx, cid, id, mid, content, attr, vote, state)
		if err != nil {
			log.Error("AddOpinionTx err(%v)", err)
			return
		}
	}
	return
}

// VoteInfo jury user vote info.
func (s *Service) VoteInfo(c context.Context, mid int64, cid int64) (vi *model.VoteInfo, err error) {
	vi, err = s.VoteInfoCache(c, mid, cid)
	return
}

// CaseInfo jury get case info.
func (s *Service) CaseInfo(c context.Context, cid int64) (res *model.BlockedCase, err error) {
	var ca *model.BlockedCase
	if ca, err = s.CaseInfoCache(c, cid); err != nil || ca.ID == 0 {
		return
	}
	if ca.MID <= 0 {
		return
	}
	res = &model.BlockedCase{}
	*res = *ca
	argMId := &accmdl.MidReq{
		Mid: res.MID,
	}
	info, err := s.accountClient.Info3(c, argMId)
	if err != nil {
		log.Error("s.accountClient.Info3(%d) error(%v)", res.MID, err)
		return
	}
	if info != nil {
		res.Face = info.Info.Face
		res.Uname = info.Info.Name
	}
	return
}

// JuryCase jury user case info contain vote.
func (s *Service) JuryCase(c context.Context, mid int64, cid int64) (res *model.BlockedCase, err error) {
	var (
		bj *model.BlockedJury
		vi *model.VoteInfo
		ca *model.BlockedCase
	)
	if bj, err = s.JuryInfoCache(c, mid); err != nil || bj.MID == 0 {
		return
	}
	if vi, err = s.VoteInfoCache(c, mid, cid); err != nil {
		return
	}
	if bj.Black != model.JuryWhite && vi.MID == 0 {
		err = ecode.CreditCaseLimit
		return
	} else if bj.Black == model.JuryWhite && vi.MID == 0 {
		vi = &model.VoteInfo{}
	}
	if ca, err = s.CaseInfoCache(c, cid); err != nil || ca.ID == 0 {
		err = ecode.CreditCaseNotExist
		return
	}
	res = &model.BlockedCase{}
	*res = *ca
	if res.MID > 0 {
		var info *accmdl.InfoReply
		argMId := &accmdl.MidReq{
			Mid: res.MID,
		}
		if info, err = s.accountClient.Info3(c, argMId); err != nil {
			log.Error("s.accountClient.Info3(%d) error(%v)", res.MID, err)
			err = nil
		}
		if info.Info != nil {
			res.Face = info.Info.Face
			res.Uname = info.Info.Name
		}
	}
	res.Vote = vi.Vote
	res.ID = cid
	expired := int64(vi.Expired)*1000 - time.Now().UnixNano()/1e6
	res.StatusTitle = fmt.Sprintf("封禁%d天", res.BlockedDays)
	res.Build()
	res.VoteTime = res.VoteTime * 1000
	if vi.Vote == 0 {
		res.ExpiredMillis = expired
	}
	return
}

// CaseList user case list. TODO:cache
func (s *Service) CaseList(c context.Context, mid, pn, ps int64) (res []*model.BlockedCase, err error) {
	var bs []*model.BlockedCase
	if bs, err = s.CaseListCache(c, mid, pn, ps); err != nil {
		return
	}
	var mids []int64
	for _, v := range bs {
		bc := &model.BlockedCase{}
		*bc = *v
		if bc.JudgeType == 1 {
			if bc.PunishResult == 3 {
				bc.StatusTitle = "永久封禁"
			} else if v.PunishResult == 2 {
				bc.StatusTitle = fmt.Sprintf("封禁%d天", bc.BlockedDays)
			}
		}
		bc.Build()
		bc.VoteTime = bc.VoteTime * 1000
		mids = append(mids, bc.MID)
		res = append(res, bc)
	}
	argMIds := &accmdl.MidsReq{Mids: mids}
	infos, err := s.accountClient.Infos3(c, argMIds)
	if err != nil {
		log.Error("s.accountClient.Infos(%+v) error(%v)", argMIds, err)
		return
	}
	for _, v := range res {
		if info, ok := infos.Infos[v.MID]; ok {
			v.Uname = info.Name
			v.Face = info.Face
		}
	}
	return
}

// Notice get notice info.
func (s *Service) Notice(c context.Context) (n *model.Notice, err error) {
	var mc = true
	if n, err = s.dao.NoticeInfoCache(c); err != nil {
		err = nil
		mc = false
	}
	if n != nil {
		return
	}
	if n, err = s.dao.Notice(c); err != nil {
		return
	}
	if mc && n != nil {
		s.addCache(func() {
			s.dao.SetNoticeInfoCache(context.TODO(), n)
		})
	}
	return
}

// ReasonList get reason list.
func (s *Service) ReasonList(c context.Context) (n []*model.Reason, err error) {
	var mc = true
	if n, err = s.dao.ReasonListCache(c); err != nil {
		err = nil
		mc = false
	}
	if len(n) > 0 {
		return
	}
	if n, err = s.dao.ReasonList(c); err != nil {
		return
	}
	if mc && len(n) > 0 {
		s.addCache(func() {
			s.dao.SetReasonListCache(context.TODO(), n)
		})
	}
	return
}

// KPIList get kpi list.
func (s *Service) KPIList(c context.Context, mid int64) (res []*model.KPI, err error) {
	var (
		j       *model.BlockedJury
		kpiData []*model.KPIData
		rr      *model.KPI
	)
	res = []*model.KPI{}
	if kpiData, err = s.dao.KPIList(c, mid); err != nil {
		return
	}
	if j, err = s.JuryInfoCache(c, mid); err != nil {
		return
	}
	for _, r := range kpiData {
		rr = &r.KPI
		rr.VoteTotal = r.VoteRealTotal
		rr.Number = j.ID
		rr.TermEnd = r.Day
		rr.TermStart = xtime.Time(r.Day.Time().AddDate(0, 0, -30).Unix())
		res = append(res, rr)
	}
	return
}

// VoteOpinion get vote opinion.
func (s *Service) VoteOpinion(c context.Context, cid, pn, ps int64, otype int8) (ops []*model.Opinion, count int, err error) {
	var (
		start = ps * (pn - 1)
		end   = ps*pn - 1
		ok    bool
		ids   []int64
	)
	ok, _ = s.dao.ExpireVoteIdx(c, cid, otype)
	if ok {
		count, err = s.dao.LenVoteIdx(c, cid, otype)
		if err != nil {
			return
		}
		ids, err = s.dao.VoteOpIdxCache(c, cid, start, end, otype)
		if err != nil {
			log.Error("s.VoteIdxCache err(%v)", err)
			return
		}
	} else {
		var (
			allops   []*model.Opinion
			ruleIdx  []int64
			breakIdx []int64
		)
		allops, err = s.dao.OpinionIdx(c, cid)
		if err != nil {
			return
		}
		for _, op := range allops {
			if op.Vote == model.VoteBanned || op.Vote == model.VoteDel {
				breakIdx = append(breakIdx, op.OpID)
			} else if op.Vote == model.VoteRule {
				ruleIdx = append(ruleIdx, op.OpID)
			}
		}
		s.addCache(func() {
			s.dao.LoadVoteOpIdxs(context.TODO(), cid, model.OpinionRule, ruleIdx)
			s.dao.LoadVoteOpIdxs(context.TODO(), cid, model.OpinonBreak, breakIdx)
		})
		if otype == model.OpinionRule {
			count = len(ruleIdx)
			if len(ruleIdx) > int(ps) {
				ids = ruleIdx[:ps]
			} else {
				ids = ruleIdx[:]
			}
		} else {
			count = len(breakIdx)
			if len(breakIdx) > int(ps) {
				ids = breakIdx[:ps]
			} else {
				ids = breakIdx[:]
			}
		}
	}
	ops, err = s.opinion(c, ids, false)
	return
}

// CaseOpinion get case opinion.
func (s *Service) CaseOpinion(c context.Context, cid, pn, ps int64) (ops []*model.Opinion, count int, err error) {
	var (
		start = ps * (pn - 1)
		end   = ps*pn - 1
		ids   []int64
		ok    bool
	)
	ok, _ = s.dao.ExpireCaseIdx(c, cid)
	if ok {
		count, err = s.dao.LenCaseIdx(c, cid)
		if err != nil {
			return
		}
		ids, err = s.dao.CaseOpIdxCache(c, cid, start, end)
		if err != nil {
			return
		}
	} else {
		var (
			allops []*model.Opinion
			allIdx []int64
		)
		allops, err = s.dao.OpinionCaseIdx(c, cid)
		if err != nil {
			return
		}
		count = len(allops)
		s.addCache(func() {
			s.dao.LoadCaseIdxs(context.TODO(), cid, allops)
		})
		for _, op := range allops {
			allIdx = append(allIdx, op.OpID)
		}
		switch {
		case len(allIdx) < int(start):
			ids = nil
		case len(allIdx) <= int(end):
			ids = allIdx[start:]
		default:
			ids = allIdx[start : end+1]
		}
	}
	ops, err = s.opinion(c, ids, true)
	return
}

// DelOpinion del opinion.
func (s *Service) DelOpinion(c context.Context, cid, opid int64) (err error) {
	err = s.dao.DelOpinion(c, opid)
	if err != nil {
		return
	}
	s.dao.DelCaseIdx(c, cid)
	s.dao.DelVoteIdx(c, cid)
	return
}

func (s *Service) opinion(c context.Context, ids []int64, needAcc bool) (ops []*model.Opinion, err error) {
	if len(ids) == 0 {
		return
	}
	var (
		miss   []int64
		tops   []*model.Opinion
		mids   []int64
		opsmap map[int64]*model.Opinion
		infos  *accmdl.InfosReply
	)
	opsmap, miss, err = s.dao.OpinionsCache(c, ids)
	if err != nil {
		return
	}
	if len(miss) > 0 {
		if tops, err = s.dao.Opinions(c, miss); err != nil {
			log.Error("s.dao.Opinions err(%v)", err)
			return
		}
		for _, top := range tops {
			opsmap[top.OpID] = top
		}
		s.addCache(func() {
			for _, top := range tops {
				s.dao.AddOpinionCache(context.TODO(), top)
			}
		})
	}
	if needAcc && len(opsmap) > 0 {
		for _, mop := range opsmap {
			mids = append(mids, mop.Mid)
		}
		arg := &accmdl.MidsReq{
			Mids: mids,
		}
		if infos, err = s.accountClient.Infos3(c, arg); err != nil {
			log.Error("s.accountClient.Infos err(%v)", err)
			err = nil
			// ignore error
		}
	}
	for _, opid := range ids {
		if op, ok := opsmap[opid]; ok {
			if needAcc && infos.Infos != nil {
				if info, ok := infos.Infos[op.Mid]; ok {
					op.Name = info.Name
					op.Face = info.Face
				}
			}
			if op.Attr != model.BlockedOpinionAttrOn {
				op.Mid = 0
				op.Name = ""
				op.Face = ""
			}
			ops = append(ops, op)
		}
	}
	return
}

// AddBlockedCases batch add blocked cases.
func (s *Service) AddBlockedCases(c context.Context, bc []*model.ArgJudgeCase) (err error) {
	if len(bc) > model.MaxAddCaseNum {
		err = ecode.RequestErr
		log.Error("s.AddBlockedCases maxaddCaseNum(%d) err(%v)", len(bc), err)
		return
	}
	var bcsn []*model.ArgJudgeCase
	for _, b := range bc {
		switch int8(b.OType) {
		case model.OriginReply:
			if b.RPID == 0 || b.Type == 0 || b.OID == 0 {
				err = ecode.RequestErr
				return
			}
			b.RelationID = fmt.Sprintf("%d-%d-%d", b.RPID, b.Type, b.OID)
			b.ReasonType = int64(model.BlockedReasonTypeByReply(int8(b.ReasonType)))
		case model.OriginTag:
			if b.TagID == 0 || b.AID == 0 {
				err = ecode.RequestErr
				return
			}
			b.RelationID = fmt.Sprintf("%d-%d", b.TagID, b.AID)
			b.ReasonType = int64(model.BlockedReasonTypeByTag(int8(b.ReasonType)))
		case model.OriginDM:
			if b.AID == 0 || b.RPID == 0 || b.OID == 0 || b.Page == 0 {
				err = ecode.RequestErr
				return
			}
			b.RelationID = fmt.Sprintf("%d-%d-%d-%d", b.AID, b.RPID, b.OID, b.Page)
		}
		var count int64
		if count, err = s.dao.CaseRelationIDCount(c, int8(b.OType), b.RelationID); err != nil {
			log.Error("ss.dao.CaseRelationIDCount(%d,%s) err(%v)", b.OType, b.RelationID, err)
			return
		}
		if count > 0 {
			log.Warn("otype(%d) relationID(%s) is alreadly juge", int8(b.OType), b.RelationID)
			continue
		}
		var total int64
		total, err = s.dao.BlockTotalTime(c, b.MID, time.Now().AddDate(-1, 0, 0))
		if err != nil {
			return
		}
		if total == 0 {
			b.PunishResult = model.Block7Days
			b.BlockedDays = 7
		} else if total == 1 {
			b.PunishResult = model.Block15Days
			b.BlockedDays = 15
		} else if total > 1 {
			b.PunishResult = model.BlockForever
			b.BlockedDays = 0
		}
		bcsn = append(bcsn, b)
	}
	if len(bcsn) <= 0 {
		log.Warn("no case submit!")
		return
	}
	if err = s.dao.AddBlockedCases(c, bcsn); err != nil {
		return
	}
	return
}

// CaseObtainByID obtain case by case id.
// NOTE: just for specific case.
func (s *Service) CaseObtainByID(c context.Context, mid, cid int64) (err error) {
	juryInfo, err := s.JuryInfoCache(c, mid)
	if err != nil {
		return
	}
	if juryInfo.MID == 0 {
		err = ecode.CreditNotJury
		return
	}
	if juryInfo.Status != model.JuryStatusEffect || int64(juryInfo.Expired) < time.Now().Unix() {
		err = ecode.CreditJuryExpired
		return
	}
	if juryInfo.VoteTotal > 3 && (float64(juryInfo.VoteRight)/float64(juryInfo.VoteTotal)*100 < float64(s.c.Judge.JuryRatio)) {
		err = ecode.CreditUnderVoteRate
		return
	}
	if cid, err = s.caseVoteID(c, mid, cid); err != nil {
		return
	}
	if cid == 0 {
		err = ecode.CreditNoCase
		return
	}
	return
}

// SpJuryCase get specific jury case info.
// NOTE : just for specific case for boss.
func (s *Service) SpJuryCase(c context.Context, mid int64, cid int64) (res *model.BlockedCase, err error) {
	var ca *model.BlockedCase
	if ca, err = s.CaseInfoCache(c, cid); err != nil {
		return
	}
	if ca.ID == 0 {
		err = ecode.CreditCaseNotExist
		return
	}
	if !model.IsCaseTypePublic(ca.CaseType) {
		err = ecode.CreditCaseNotExist
		return
	}
	res = &model.BlockedCase{}
	*res = *ca
	var info *accmdl.InfoReply
	if info, err = s.accountClient.Info3(c, &accmdl.MidReq{Mid: res.MID}); err != nil {
		log.Error("s.accountClient.Info3(%d) error(%+v)", res.MID, err)
		err = nil
	}
	if info.Info != nil {
		res.Uname = info.Info.Name
		res.Face = info.Info.Face
	}
	s._caseExpland(cid, res)
	s._buildVoteInfo(c, mid, cid, res)
	return
}

func (s *Service) _caseExpland(cid int64, bc *model.BlockedCase) {
	bc.ID = cid
	bc.JudgeType = 0
	// set status to dealing and result to 0.
	if bc.Status > model.CaseStatusGrantStop {
		bc.Status = model.CaseStatusDealing
	}
	bc.StatusTitle = fmt.Sprintf("封禁%d天", bc.BlockedDays)
	bc.PunishTitle = fmt.Sprintf("在%s中%s", model.OriginTypeDesc(bc.OriginType), model.ReasonTypeDesc(bc.ReasonType))
}

func (s *Service) _buildVoteInfo(c context.Context, mid, cid int64, bc *model.BlockedCase) {
	var (
		err error
		vi  *model.VoteInfo
	)
	if vi, err = s.VoteInfoCache(c, mid, cid); err == nil && vi.MID != 0 {
		bc.Vote = vi.Vote
		if vi.Vote == 0 {
			bc.ExpiredMillis = int64(vi.Expired)*1000 - time.Now().UnixNano()/1e6
		}
	}
}

// VoteInfoCache use vote cid info.
func (s *Service) VoteInfoCache(c context.Context, mid, cid int64) (vi *model.VoteInfo, err error) {
	var mc = true
	if vi, err = s.dao.VoteInfoCache(c, mid, cid); err != nil {
		err = nil
		mc = false
	}
	if vi != nil {
		return
	}
	if vi, err = s.dao.VoteInfo(c, mid, cid); err != nil {
		return
	}
	if vi == nil {
		vi = &model.VoteInfo{}
	}
	if mc {
		s.addCache(func() {
			s.dao.SetVoteInfoCache(context.TODO(), mid, cid, vi)
		})
	}
	return
}

// JuryInfos mutli get jurys info.
func (s *Service) JuryInfos(c context.Context, mids []int64) (res map[int64]*model.ResJuryerStatus, err error) {
	res = make(map[int64]*model.ResJuryerStatus, len(mids))
	mbj, err := s.dao.JuryInfos(c, mids)
	if err != nil {
		err = errors.Wrap(err, "s.dao.JuryInfos")
		return
	}
	for mid, bj := range mbj {
		res[mid] = &model.ResJuryerStatus{
			Expired: bj.Expired,
			Mid:     mid,
			Status:  bj.Status,
		}
	}
	return
}

// BatchBLKCases get batch blocked cases by ids.
func (s *Service) BatchBLKCases(c context.Context, ids []int64) (cases map[int64]*model.BlockedCase, err error) {
	if cases, err = s.dao.CaseInfoIDs(c, ids); err != nil {
		err = errors.Wrapf(err, "s.dao.CaseInfoIDs(%+v)", ids)
		return
	}
	for _, item := range cases {
		item.Build()
	}
	return
}

// JuryInfoCache JuryInfo cache .
func (s *Service) JuryInfoCache(c context.Context, mid int64) (bj *model.BlockedJury, err error) {
	var mc = true
	if bj, err = s.dao.JuryInfoCache(c, mid); err != nil {
		err = nil
		mc = false
	}
	if bj != nil {
		return
	}
	if bj, err = s.dao.JuryInfo(c, mid); err != nil {
		return
	}
	if bj == nil {
		bj = &model.BlockedJury{}
	} else {
		if bj.CaseTotal, err = s.dao.CountCaseVote(c, mid); err != nil {
			return
		}
	}
	if mc {
		s.addCache(func() {
			if err = s.dao.SetJuryInfoCache(context.TODO(), mid, bj); err != nil {
				log.Error("s.dao.SetJuryInfoCache error(%+v)", err)
			}
		})
	}
	return
}

// CaseInfoCache .
func (s *Service) CaseInfoCache(c context.Context, cid int64) (bc *model.BlockedCase, err error) {
	var mc = true
	if bc, err = s.dao.CaseInfoCache(c, cid); err != nil {
		err = nil
		mc = false
	}
	if bc != nil {
		return
	}
	if bc, err = s.dao.CaseInfo(c, cid); err != nil {
		return
	}
	if bc == nil {
		bc = &model.BlockedCase{}
	}
	if mc {
		s.addCache(func() {
			if err = s.dao.SetCaseInfoCache(context.TODO(), cid, bc); err != nil {
				log.Error("s.dao.SetCaseInfoCache error(%+v)", err)
			}
		})
	}
	return
}

// CaseListCache .
func (s *Service) CaseListCache(c context.Context, mid, pn, ps int64) (bs []*model.BlockedCase, err error) {
	if pn*ps <= 100 {
		bs, err = s.caseVoteTopCache(c, mid, pn, ps)
		return
	}
	var vids, cids []int64
	if vids, cids, err = s.dao.CaseVoteIDMID(c, mid, pn, ps); err != nil {
		return
	}
	bs, err = s.buildBlockedCase(c, vids, cids)
	return
}

// caseVoteTopCache .
func (s *Service) caseVoteTopCache(c context.Context, mid, pn, ps int64) (res []*model.BlockedCase, err error) {
	var (
		mc         = true
		vids, cids []int64
		bs         []*model.BlockedCase
		end        = pn * ps
		start      = (pn - 1) * ps
	)
	defer func() {
		bl := len(bs)
		if err == nil && bl != 0 {
			switch {
			case bl <= int(start):
				res = _emptyBlockedCase
			case bl <= int(end):
				res = bs[start:]
			default:
				res = bs[start:end]
			}
		}
	}()
	if bs, err = s.dao.CaseVoteTopCache(c, mid); err != nil {
		err = nil
		mc = false
	}
	if bs != nil {
		return
	}
	if vids, cids, err = s.dao.CaseVoteIDTop(c, mid); err != nil {
		return
	}
	if bs, err = s.buildBlockedCase(c, vids, cids); err != nil {
		return
	}
	if mc {
		s.addCache(func() {
			if err = s.dao.SetCaseVoteTopCache(context.TODO(), mid, bs); err != nil {
				log.Error("s.dao.SetCaseVoteTopCache error(%+v)", err)
			}
		})
	}
	return
}

// buildBlockedCase .
func (s *Service) buildBlockedCase(c context.Context, vids, cids []int64) (bs []*model.BlockedCase, err error) {
	var (
		ok  bool
		vo  *model.VoteInfo
		bc  *model.BlockedCase
		mvo map[int64]*model.VoteInfo
		mbc map[int64]*model.BlockedCase
	)
	if len(vids) == 0 || len(cids) == 0 {
		bs = _emptyBlockedCase
		return
	}
	if mvo, err = s.dao.CaseVotesMID(c, vids); err != nil {
		return
	}
	if mbc, err = s.dao.CaseVoteIDs(c, cids); err != nil {
		return
	}
	if mvo == nil || mbc == nil {
		bs = _emptyBlockedCase
		return
	}
	for _, cid := range cids {
		if bc, ok = mbc[cid]; !ok {
			continue
		}
		if vo, ok = mvo[cid]; !ok {
			continue
		}
		bc.Vote = vo.Vote
		bc.VoteTime = vo.Mtime
		bs = append(bs, bc)
	}
	if len(bs) == 0 {
		bs = _emptyBlockedCase
	}
	return
}
