package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/dm2/model"
	"go-common/app/interface/main/dm2/model/oplog"
	account "go-common/app/service/main/account/api"
	"go-common/app/service/main/archive/api"
	arcMdl "go-common/app/service/main/archive/model/archive"
	figureMdl "go-common/app/service/main/figure/model"
	filterMdl "go-common/app/service/main/filter/api/grpc/v1"
	locmdl "go-common/app/service/main/location/model"
	spyMdl "go-common/app/service/main/spy/model"
	ugcPayMdl "go-common/app/service/main/ugcpay/api/grpc/v1"
	seasonMbl "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

var (
	msgRegex = regexp.MustCompile(`^(\s|\xE3\x80\x80)*$`) // 全文仅空格

	_dateFormat = "20060102"
)

//Post dm post
func (s *Service) Post(c context.Context, dm *model.DM, aid, rnd int64) (err error) {
	// 验证主题是否存在
	sub, err := s.subject(c, dm.Type, dm.Oid)
	if err != nil {
		return
	}
	if sub.State == model.SubStateClosed {
		err = ecode.DMForbidPost
		return
	}
	if sub.Maxlimit == 0 {
		return
	}
	// 验证海外用户
	if err = s.checkOverseasUser(c); err != nil {
		return
	}
	// 验证账号信息
	myinfo, err := s.checkAccountInfo(c, dm.Mid)
	if err != nil {
		return
	}
	// 验证弹幕内容
	if err = s.checkMsg(c, dm, myinfo); err != nil {
		return
	}
	// 验证弹幕progress
	duration, err := s.videoDuration(c, aid, dm.Oid)
	if err != nil {
		return
	}
	if duration > 0 && int64(dm.Progress) > duration {
		return ecode.DMProgressTooBig
	}
	// 验证稿件信息
	arc, err := s.checkArchiveInfo(c, aid, dm.Oid, myinfo.GetRank(), dm.Content.Mode)
	if err != nil {
		return
	}
	if sub.Mid != dm.Mid && arc.Rights.UGCPay == arcMdl.AttrYes {
		if err = s.archiveUgcPay(c, dm.Mid, aid); err != nil {
			return
		}
	}
	// 验证弹幕发送速率
	if err = s.checkPubRate(c, dm, sub.Mid, rnd, myinfo.GetRank()); err != nil {
		return
	}
	// 判定用户发送速率
	if err = s.checkRateLimit(c, myinfo.GetRank(), dm.Mid); err != nil {
		return
	}
	// 验证弹幕颜色
	if err = s.checkMsgColor(sub.Mid, myinfo, dm); err != nil {
		return
	}
	// 判定弹幕字体大小
	if err = s.checkFontsize(dm, sub.Mid, myinfo.GetRank()); err != nil {
		return
	}
	// 验证mode 1、4、5、6
	if err = s.checkNormalMode(sub.Mid, myinfo, dm); err != nil {
		return
	}
	// 验证mode 7发送权限
	if err = s.checkAdvanceMode(c, dm, sub.Mid, myinfo); err != nil {
		return
	}
	// 验证特殊弹幕池(pool=1 pool=2)发送权限
	if err = s.checkSpecialPool(c, dm, sub.Mid, myinfo); err != nil {
		return
	}
	// 生成弹幕id
	if err = s.genDMID(c, dm); err != nil {
		return
	}
	// 弹幕屏蔽词过滤
	if err = s.checkFilterService(c, dm, arc); err != nil {
		return
	}
	// 检查up的全局屏蔽词
	if err = s.checkUpFilter(c, dm, sub.Mid, myinfo.GetRank()); err != nil {
		return
	}
	// 弹幕的先审后发和先发后审
	if err = s.checkMonitor(c, sub, dm); err != nil {
		return
	}
	// 垃圾弹幕过滤
	if err = s.checkUnusualAction(c, dm, myinfo); err != nil {
		return
	}
	// bnj专用 黑名单过滤
	s.checkShield(c, aid, dm)
	// 发消息给job异步落库
	if err = s.asyncAddDM(dm, arc.Aid, arc.TypeID); err != nil {
		return
	}
	remark := fmt.Sprintf("新增弹幕,ip:%s,port:%s", metadata.String(c, metadata.RemoteIP), metadata.String(c, metadata.RemotePort))
	s.OpLog(c, dm.Oid, dm.Mid, time.Now().Add(time.Second).Unix(), 1, []int64{dm.ID}, "status", "", strconv.FormatInt(int64(dm.State), 10), remark, oplog.SourcePlayer, oplog.OperatorMember)
	// 弹幕行为日志
	s.dao.ReportDmLog(c, dm)
	// 弹幕广播
	select {
	case s.broadcastChan <- &broadcast{Aid: arc.Aid, Rnd: rnd, DM: dm}:
	default:
		log.Error("broadcast channel is full")
	}
	return
}

// checkUnusualAction check unusual post dm
func (s *Service) checkUnusualAction(c context.Context, dm *model.DM, accInfo *account.Profile) (err error) {
	var (
		ip             = metadata.String(c, metadata.RemoteIP)
		userScore      *spyMdl.UserScore
		figureWithRank *figureMdl.FigureWithRank
		dmDailyLimit   *model.DailyLimiter
	)

	if !s.garbageDanmu {
		return
	}
	if userScore, err = s.spyRPC.UserScore(c, &spyMdl.ArgUserScore{
		Mid: accInfo.GetMid(),
		IP:  ip,
	}); err != nil {
		// dragrade spy service
		err = nil
		return
	}
	if figureWithRank, err = s.figureRPC.UserFigure(c, &figureMdl.ArgUserFigure{
		Mid: accInfo.GetMid(),
	}); err != nil {
		if ecode.Cause(err).Code() != ecode.FigureNotFound.Code() {
			log.Error("checkUnusualAction.UserFigure(mid:%v) error(%v)", accInfo.GetMid(), err)
			return
		}
		err = nil
		return
	}
	if accInfo.GetMid() > 50000000 && accInfo.GetLevel() <= 2 && figureWithRank.Percentage > 40 && userScore.Score < 90 {
		if dmDailyLimit, err = s.dao.GetDmDailyLimitCache(c, accInfo.GetMid()); err != nil {
			return
		}
		now := time.Now().Format(_dateFormat)
		if dmDailyLimit == nil || now > dmDailyLimit.Date {
			dmDailyLimit = &model.DailyLimiter{
				Date:  time.Now().Format(_dateFormat),
				Count: 0,
			}
		}
		dmDailyLimit.Count++
		if dmDailyLimit.Count > 5 {
			dm.State = model.StateDelete
			s.OpLog(c, dm.Oid, dm.Mid, time.Now().Unix(), 1, []int64{dm.ID}, "status", "", strconv.FormatInt(int64(dm.State), 10), "垃圾弹幕发送", oplog.SourcePlayer, oplog.OperatorMember)
			s.dao.ReportDmGarbageLog(c, dm)
		}
		if err = s.dao.SetDmDailyLimitCache(c, accInfo.GetMid(), dmDailyLimit); err != nil {
			return
		}
		return
	}
	return
}

func (s *Service) checkFontsize(dm *model.DM, upid int64, rank int32) (err error) {
	if !s.isSuperUser(rank) && dm.Content.Mode != model.ModeSpecial {
		if dm.Content.FontSize != 18 && dm.Content.FontSize != 25 {
			dm.Content.FontSize = 25
		}
	}
	return
}

func (s *Service) checkMsgColor(upid int64, profile *account.Profile, dm *model.DM) (err error) {
	if s.isSuperUser(profile.GetRank()) || dm.Mid == upid {
		return
	}
	if profile.GetLevel() <= 1 && dm.Content.Color != 0xffffff {
		err = ecode.DMMsgNoColorPerm
	}
	return
}

// 验证mode 1、4、5、6
func (s *Service) checkNormalMode(upid int64, profile *account.Profile, dm *model.DM) (err error) {
	if s.isSuperUser(profile.GetRank()) {
		return
	}
	switch dm.Content.Mode {
	case 1:
		if profile.GetLevel() < 1 && upid != dm.Mid {
			err = ecode.DMMsgNoPubPerm
		}
	case 4:
		if profile.GetLevel() < 3 && upid != dm.Mid {
			err = ecode.DMMsgNoPubBottomPerm
		}
	case 5:
		if profile.GetLevel() < 3 && upid != dm.Mid {
			err = ecode.DMMsgNoPubTopPerm
		}
	case 6:
		err = ecode.DMMsgNoPubAdvancePerm
	}
	return
}

// 验证 mode 7
func (s *Service) checkAdvanceMode(c context.Context, dm *model.DM, upid int64, profile *account.Profile) (err error) {
	if dm.Content.Mode != 7 || dm.Mid == upid || s.isSuperUser(profile.GetRank()) {
		return
	}
	if profile.GetLevel() <= 1 {
		err = ecode.DMMsgNoPubStylePerm
		return
	}
	var adv *model.AdvanceCmt
	if adv, err = s.advanceComment(c, dm.Oid, dm.Mid, "sp"); err != nil {
		return
	}
	if adv.Type != "buy" && adv.Type != "accept" {
		err = ecode.DMMsgNoPubStylePerm
	}
	return
}

// 验证 pool 1、2
func (s *Service) checkSpecialPool(c context.Context, dm *model.DM, upid int64, profile *account.Profile) (err error) {
	if (dm.Pool != 1 && dm.Pool != 2) || dm.Mid == upid || s.isSuperUser(profile.Rank) {
		return
	}
	if dm.Pool == model.PoolSubtitle {
		dm.Pool = model.PoolNormal
		log.Warn("force change(%+v) pool to normal", dm)
		return
	}
	if profile.Level <= 1 {
		err = ecode.DMMsgNoPubStylePerm
		return
	}
	var adv *model.AdvanceCmt
	if dm.Pool == model.PoolSpecial {
		if adv, err = s.advanceComment(c, dm.Oid, dm.Mid, "advance"); err != nil {
			return
		}
		if adv.Type != "buy" && adv.Type != "accept" {
			err = ecode.DMMsgNoPubStylePerm
			return
		}
		if profile.Rank <= 10000 {
			err = ecode.DMMsgNoPubStylePerm
		}
	}
	return
}

func (s *Service) checkFilterService(c context.Context, dm *model.DM, arc *api.Arc) (err error) {
	if dm.Content.Mode == 7 || dm.Content.Mode == 8 || dm.Content.Mode == 9 {
		return
	}
	var (
		sid         int32
		pid         int64
		rid         = arc.TypeID
		keys        []string
		filterReply *filterMdl.FilterReply
		seasonReply *seasonMbl.CardsInfoReply
	)
	if v, ok := s.arcTypes[int16(rid)]; ok {
		pid = int64(v)
	}
	if arc.AttrVal(arcMdl.AttrBitIsBangumi) == arcMdl.AttrYes {
		if seasonReply, err = s.seasonRPC.CardsByAids(c, &seasonMbl.SeasonAidReq{
			Aids: []int32{int32(arc.Aid)},
		}); err != nil || seasonReply == nil {
			log.Error("s.seasonRPC.CardsByAids(%d) error(%v)", arc.Aid, err) // NOTE ignore error and continue
		} else if seasonInfo, ok := seasonReply.Cards[int32(arc.Aid)]; !ok {
			log.Error("seasonReply.Cards(%d) don't exist", arc.Aid)
		} else {
			sid = seasonInfo.SeasonId
		}
	}
	if sid > 0 {
		keys = append(keys, fmt.Sprintf("season:%d", sid))
	}
	keys = append(keys, fmt.Sprintf("typeid:%d", rid))
	keys = append(keys, fmt.Sprintf("typeid:%d", pid))
	keys = append(keys, fmt.Sprintf("cid:%d", dm.Oid))
	keys = append(keys, fmt.Sprintf("aid:%d", arc.Aid))
	if filterReply, err = s.filterRPC.Filter(c, &filterMdl.FilterReq{
		Area:    "danmu",
		Message: dm.Content.Msg,
		TypeId:  int64(rid),
		Id:      dm.ID,
		Oid:     dm.Oid,
		Mid:     dm.Mid,
		Keys:    keys,
	}); err != nil {
		log.Error("checkFilterService(dm:%+v),err(%v)", dm, err)
		return
	}
	if filterReply.Level > 0 || filterReply.Limit == model.SpamBlack || filterReply.Limit == model.SpamOverflow {
		dm.State = model.StateFilter
		log.Info("filter service delete(dmid:%d,data:+%v)", dm.ID, filterReply)
		remark := filterReply.Result
		if filterReply.Limit == model.SpamBlack {
			remark = "命中反垃圾黑名单"
		}
		if filterReply.Limit == model.SpamOverflow {
			remark = "超过反垃圾限制次数"
		}
		s.OpLog(c, dm.Oid, dm.Mid, time.Now().Unix(), 1, []int64{dm.ID}, "status", "", strconv.FormatInt(int64(dm.State), 10), remark, oplog.SourcePlayer, oplog.OperatorMember)
		return
	}
	if filterReply.Ai != nil && len(filterReply.Ai.Scores) > 0 && filterReply.Ai.Scores[0] > filterReply.Ai.Threshold {
		dm.State = model.StateAiDelete
		s.OpLog(c, dm.Oid, dm.Mid, time.Now().Unix(), 1, []int64{dm.ID}, "status", "", strconv.FormatInt(int64(dm.State), 10), "ai 反垃圾屏蔽", oplog.SourcePlayer, oplog.OperatorMember)
		return
	}
	return
}

// checkUpFilter up主屏蔽词过滤
func (s *Service) checkUpFilter(c context.Context, dm *model.DM, upid int64, rank int32) (err error) {
	// 命中屏蔽词系统 或者 mode 8 9 后则不再进行校验up主屏蔽词
	if dm.State == model.StateFilter || dm.Content.Mode == model.ModeCode || dm.Content.Mode == model.ModeBAS {
		return
	}
	var (
		msg                  string
		fltModes             []int32
		texts, regexs, users []string
	)
	// up主全局屏蔽词
	filters, err := s.UpFilters(c, upid)
	if err != nil {
		return
	}
	for _, f := range filters {
		switch f.Type {
		case 0: // 文本类型
			texts = append(texts, f.Filter)
		case 1: // 正则类型
			f.Filter = strings.Replace(f.Filter, "/", "\\/", -1)
			regexs = append(regexs, f.Filter)
		case 2: // 用户黑名单
			f.Filter = strings.ToLower(f.Filter)
			users = append(users, f.Filter)
		case 4, 5, 6, 7: // 4:down,5:up,6:reverse,7:special
			fltModes = append(fltModes, int32(f.Type))
		}
	}
	// check content filter
	if dm.Content.Mode == model.ModeSpecial {
		strs := strings.Split(dm.Content.Msg, ",")
		if len(strs) < 4 || len(strs[4]) < 2 {
			log.Error("s.checkUpFilter(%s) error(spec content format err)", dm)
			return
		}
		msg = strs[4][1 : len(strs[4])-1]
	} else {
		msg = dm.Content.Msg
	}
	// 验证up设置的mode屏蔽
	for _, mode := range fltModes {
		if mode == dm.Content.Mode {
			dm.State = model.StateBlock
			s.OpLog(c, dm.Oid, dm.Mid, time.Now().Unix(), 1, []int64{dm.ID}, "status", "", strconv.FormatInt(int64(dm.State), 10), "命中类型屏蔽（"+fmt.Sprint(mode)+"）", oplog.SourcePlayer, oplog.OperatorMember)
			return
		}
	}
	// 关键字过滤
	for _, text := range texts {
		if strings.Contains(strings.ToLower(msg), strings.ToLower(text)) {
			dm.State = model.StateBlock
			s.OpLog(c, dm.Oid, dm.Mid, time.Now().Unix(), 1, []int64{dm.ID}, "status", "", strconv.FormatInt(int64(dm.State), 10), "命中关键字屏蔽（"+text+"）", oplog.SourcePlayer, oplog.OperatorMember)
			return
		}
	}
	// 正则过滤
	for _, reg := range regexs {
		rc, rErr := regexp.Compile(reg)
		if rErr != nil {
			log.Error("regexp.Compile(%s) error(%v)", reg, rErr)
			continue
		}
		if rc.MatchString(msg) {
			dm.State = model.StateBlock
			s.OpLog(c, dm.Oid, dm.Mid, time.Now().Unix(), 1, []int64{dm.ID}, "status", "", strconv.FormatInt(int64(dm.State), 10), "命中正则屏蔽（"+reg+"）", oplog.SourcePlayer, oplog.OperatorMember)
			return
		}
	}
	// 验证up设置的屏蔽用户
	hashID := model.Hash(dm.Mid, uint32(dm.Content.IP))
	for _, user := range users {
		if hashID == user {
			dm.State = model.StateBlock
			s.OpLog(c, dm.Oid, dm.Mid, time.Now().Unix(), 1, []int64{dm.ID}, "status", "", strconv.FormatInt(int64(dm.State), 10), "屏蔽黑名单屏蔽（"+hashID+"）", oplog.SourcePlayer, oplog.OperatorMember)
			return
		}
	}
	return
}

func (s *Service) checkRateLimit(c context.Context, rank int32, mid int64) (err error) {
	now := time.Now().Unix()
	ltime, ok := model.LimitPerMin[rank]
	if !ok {
		ltime = model.LimitPerMin[0]
	}
	limiter, err := s.dao.DMLimitCache(c, mid)
	if err != nil {
		return
	}
	if limiter == nil {
		limiter = &model.Limiter{Allowance: ltime, Timestamp: now}
	}
	allowance := limiter.Allowance + ((now - limiter.Timestamp) * ltime / 600)
	if allowance > ltime {
		allowance = ltime
	}
	if allowance < 1 {
		err = ecode.DMMsgPubTooFast
		allowance = 0
	} else {
		allowance--
	}
	limiter = &model.Limiter{Allowance: allowance, Timestamp: now}
	s.dao.AddDMLimitCache(c, mid, limiter) // NOTE omit error
	return
}

func (s *Service) checkAccountInfo(c context.Context, mid int64) (profile *account.Profile, err error) {
	var (
		profileReply *account.ProfileReply
	)
	if profileReply, err = s.accountRPC.Profile3(c, &account.MidReq{
		Mid: mid,
	}); err != nil {
		log.Error("accRPC.UserInfo(%v) error(%v)", mid, err)
		return
	}
	if profileReply.GetProfile().GetIdentification() == 0 && profileReply.GetProfile().GetTelStatus() == 0 {
		err = ecode.UserCheckNoPhone
		return
	}
	if profileReply.GetProfile().GetIdentification() == 0 && profileReply.GetProfile().GetTelStatus() == 2 {
		err = ecode.UserCheckInvalidPhone
		return
	}
	if profileReply.GetProfile().GetEmailStatus() == 0 && profileReply.GetProfile().GetTelStatus() == 0 {
		err = ecode.UserCheckNoPhone
		return
	}
	if profileReply.GetProfile().GetSilence() == 1 {
		err = ecode.UserDisabled
		return
	}
	if profileReply.GetProfile().GetMoral() < 60 {
		err = ecode.LackOfScores
	}
	if profile = profileReply.GetProfile(); profile == nil {
		err = ecode.UserNotExist
	}
	return
}

func (s *Service) checkMsg(c context.Context, dm *model.DM, profile *account.Profile) (err error) {
	var (
		msg          = dm.Content.Msg
		msgLen       = len([]rune(dm.Content.Msg))
		isNormalMode = s.isNormalMode(dm.Content.Mode)
	)
	if msgRegex.MatchString(msg) { // 空白弹幕
		err = ecode.DMMsgIlleagel
		return
	}
	if profile.GetRank() < 20000 && profile.GetLevel() == 1 && msgLen > 20 {
		err = ecode.DMMsgTooLongLevel1
		return
	}
	if (msgLen > model.MaxLenDefMsg && isNormalMode) || (msgLen > model.MaxLen7Msg && dm.Content.Mode == 7) {
		err = ecode.DMMsgTooLong
		return
	}
	if isNormalMode && (strings.Contains(msg, `\n`) || strings.Contains(msg, `/n`)) {
		err = ecode.DMMsgIlleagel
		return
	}
	// 校验单字符刷屏
	if msgLen == 1 {
		var count int64
		if count, err = s.dao.CharPubCnt(c, dm.Mid, dm.Oid); err != nil {
			return
		}
		if count+1 > 3 {
			err = ecode.DMMsgPubTooFast
			return
		}
		if err = s.dao.IncrCharPubCnt(c, dm.Mid, dm.Oid); err != nil {
			return
		}
	} else {
		if err = s.dao.DelCharPubCnt(c, dm.Mid, dm.Oid); err != nil {
			return
		}
	}
	return
}

func (s *Service) checkPubRate(c context.Context, dm *model.DM, upid, rnd int64, rank int32) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	cached, err := s.dao.MsgPublock(c, dm.Mid, dm.Content.Color, rnd, dm.Content.Mode, dm.Content.FontSize, ip, dm.Content.Msg)
	if err != nil {
		return
	}
	if cached {
		err = ecode.DMMsgPubTooFast
		return
	}
	s.dao.AddMsgPubLock(c, dm.Mid, dm.Content.Color, rnd, dm.Content.Mode, dm.Content.FontSize, ip, dm.Content.Msg)
	if !s.isSuperUser(rank) {
		if cached, err = s.dao.OidPubLock(c, dm.Mid, dm.Oid, ip); err != nil {
			return
		}
		if cached {
			err = ecode.DMMsgPubTooFast
			return
		}
		s.dao.AddOidPubLock(c, dm.Mid, dm.Oid, ip)
	}
	if rank <= 15000 {
		var count int64
		count, err = s.dao.PubCnt(c, dm.Mid, dm.Content.Color, dm.Content.Mode, dm.Content.FontSize, ip, dm.Content.Msg)
		if err != nil {
			return
		}
		if count >= 8 && dm.Mid != upid {
			arg := &account.MoralReq{
				Mid:    dm.Mid,
				Moral:  -0.25,
				Oper:   "",
				Reason: "恶意刷弹幕",
				Remark: "云屏蔽",
				RealIp: ip,
			}
			if _, err = s.accountRPC.AddMoral3(c, arg); err != nil {
				log.Error("s.accRPC.AddMoral2(%v) error(%v)", arg, err)
				return
			}
			err = ecode.DMMsgPubTooFast
			return
		}
		if err = s.dao.IncrPubCnt(c, dm.Mid, dm.Content.Color, dm.Content.Mode, dm.Content.FontSize, ip, dm.Content.Msg); err != nil {
			return
		}
	}
	return
}

func (s *Service) checkArchiveInfo(c context.Context, aid, oid int64, rank, mode int32) (arc *api.Arc, err error) {
	// if _, err = s.arcRPC.Video3(c, &arcMdl.ArgVideo2{
	// 	Aid:    aid,
	// 	Cid:    oid,
	// 	RealIP: metadata.String(c, metadata.RemoteIP),
	// }); err != nil {
	// 	log.Error("s.arcRPC.Video3(aid:%v,oid:%v) error(%v)", aid, oid, err)
	// 	return
	// }
	arg := &arcMdl.ArgAid2{
		Aid:    aid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if arc, err = s.arcRPC.Archive3(c, arg); err != nil {
		log.Error("s.arcRPC.Archive3(%v) error(%v)", arg, err)
		return
	}
	if arc.State < 0 && arc.State != -6 {
		err = ecode.DMArchiveIlleagel // 禁止向未审核的视频发送弹幕
		return
	}
	return
}

func (s *Service) archiveUgcPay(c context.Context, mid int64, aid int64) (err error) {
	var (
		resp *ugcPayMdl.AssetRelationResp
	)
	resp, err = s.ugcPayRPC.AssetRelation(c, &ugcPayMdl.AssetRelationReq{
		Mid:   mid,
		Oid:   aid,
		Otype: model.UgcPayTypeArchive,
	})
	if err != nil {
		log.Error("archiveUgcPay(aid:%d,mid:%d),error(%v)", aid, mid, err)
		return
	}
	if resp.State != model.UgcPayRelationStatePaid {
		err = ecode.DMNotpayForPost
		return
	}
	return
}

func (s *Service) advanceComment(c context.Context, oid, mid int64, mode string) (adv *model.AdvanceCmt, err error) {
	if adv, err = s.dao.AdvanceCmtCache(c, oid, mid, mode); err != nil {
		log.Error("s.dao.AdvanceCmtCache mid=%d oid=%d mode=%s  error(%v)", mid, oid, mode, err)
		return
	}
	if adv == nil {
		if adv, err = s.dao.AdvanceCmt(c, oid, mid, mode); err != nil {
			log.Error("s.dao.AdvanceCmt mid=%d oid=%d mode=%s  error(%v)", mid, oid, mode, err)
			return
		}
		if adv == nil {
			adv = &model.AdvanceCmt{}
		}
		if err = s.dao.AddAdvanceCmtCache(c, oid, mid, mode, adv); err != nil {
			log.Error("s.dao.AddAdvanceCmtCache mid=%d  oid=%d mode=%s  error(%v)", mid, oid, mode, err)
			return
		}
	}
	return
}

func (s *Service) genDMID(c context.Context, dm *model.DM) (err error) {
	dmid, err := s.seqRPC.ID(c, s.seqDmArg)
	if err != nil {
		return
	}
	dm.ID = dmid
	dm.Content.ID = dmid
	if dm.ContentSpe != nil {
		dm.ContentSpe.ID = dmid
	}
	return
}

// checkMonitor check the oid is monitored
func (s *Service) checkMonitor(c context.Context, sub *model.Subject, dm *model.DM) (err error) {
	if dm.State != model.StateNormal {
		return
	}
	if sub.AttrVal(model.AttrSubMonitorBefore) == model.AttrYes {
		dm.State = model.StateMonitorBefore
	} else if sub.AttrVal(model.AttrSubMonitorAfter) == model.AttrYes {
		dm.State = model.StateMonitorAfter
	} else {
		return
	}
	return
}

func (s *Service) asyncAddDM(dm *model.DM, aid int64, tid int32) (err error) {
	var (
		data []byte
		msg  = &struct {
			*model.DM
			Aid int64 `json:"aid"`
			Tid int32 `json:"tid"`
		}{
			DM:  dm,
			Aid: aid,
			Tid: tid,
		}
		c = context.TODO()
	)
	if data, err = json.Marshal(msg); err != nil {
		log.Error("json.Marshal(%+v) error(%v)", msg, err)
		return
	}
	act := &model.Action{Action: model.ActAddDM, Data: data}
	for i := 0; i < 3; i++ {
		if err = s.dao.SendAction(c, fmt.Sprint(dm.Oid), act); err != nil {
			continue
		} else {
			break
		}
	}
	return
}

// 用户是否是特权用户
func (s *Service) isSuperUser(rank int32) bool {
	return rank >= 20000
}

// 判断弹幕模式是否是普通弹幕模式
func (s *Service) isNormalMode(mode int32) bool {
	if mode == 1 || mode == 4 || mode == 5 || mode == 6 {
		return true
	}
	return false
}

func (s *Service) checkOverseasUser(c context.Context) (err error) {
	if s.conf.Supervision.Completed {
		err = ecode.ServiceUpdate
		return
	}
	now := time.Now()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", s.conf.Supervision.StartTime, loc)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", s.conf.Supervision.EndTime, loc)
	if now.Before(startTime) || now.After(endTime) {
		err = nil
		return
	}
	arg := &locmdl.ArgIP{IP: metadata.String(c, metadata.RemoteIP)}
	zone, err := s.locationRPC.Info(c, arg)
	if err != nil {
		log.Error("s.locationRPC.Info(%s) error(%v)", metadata.String(c, metadata.RemoteIP), err)
		err = nil
		return
	}
	if !strings.EqualFold(zone.Country, s.conf.Supervision.Location) {
		err = ecode.ServiceUpdate
	}
	return
}

func (s *Service) checkShield(c context.Context, aid int64, dm *model.DM) {
	if _, ok := s.aidSheild[aid]; !ok {
		return
	}
	if _, ok := s.midsSheild[dm.Mid]; !ok {
		return
	}
	// hit
	dm.State = model.StateBlock
	s.OpLog(c, dm.Oid, dm.Mid, time.Now().Unix(), 1, []int64{dm.ID}, "status", "", strconv.FormatInt(int64(dm.State), 10), "弹幕指定稿件黑名单屏蔽", oplog.SourcePlayer, oplog.OperatorMember)
}
