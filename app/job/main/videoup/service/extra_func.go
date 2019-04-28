package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/videoup/model/archive"
	"go-common/app/job/main/videoup/model/message"
	"go-common/app/job/main/videoup/model/redis"
	accApi "go-common/app/service/main/account/api"
	"go-common/library/conf/env"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
	"strings"
)

func (s *Service) archiveVideo(c context.Context, filename string) (v *archive.Video, a *archive.Archive, err error) {
	if v, err = s.arc.NewVideo(c, filename); err != nil {
		log.Error("s.arc.Video(%s) error(%v)", filename, err)
		return
	}
	if v == nil {
		log.Error("s.arc.Video(%s) video is nil", filename)
		err = fmt.Errorf("video(%s) is not exists", filename)
		return
	}
	if a, err = s.arc.Archive(c, v.Aid); err != nil {
		log.Error("s.arc.Archive(%d) filename(%s) error(%v)", v.Aid, filename, err)
		return
	}
	if a == nil {
		log.Error("s.arc.Archive(%s) archive(%d) is nil", filename, v.Aid)
		err = fmt.Errorf("archive(%d) filename(%s) is not exists", v.Aid, filename)
	}
	return
}

func (s *Service) archiveVideoByAid(c context.Context, filename string, aid int64) (v *archive.Video, a *archive.Archive, err error) {
	if v, err = s.arc.NewVideoByAid(c, filename, aid); err != nil {
		log.Error("s.arc.Video(%s) error(%v)", filename, err)
		return
	}
	if v == nil {
		log.Error("s.arc.Video(%s) video is nil", filename)
		err = fmt.Errorf("video(%s) is not exists", filename)
		return
	}
	if a, err = s.arc.Archive(c, aid); err != nil {
		log.Error("s.arc.Archive(%d) filename(%s) error(%v)", aid, filename, err)
		return
	}
	if a == nil {
		log.Error("s.arc.Archive(%s) archive(%d) is nil", filename, aid)
		err = fmt.Errorf("archive(%d) filename(%s) is not exists", aid, filename)
	}
	return
}

func (s *Service) isPorder(a *archive.Archive) bool {
	if a == nil {
		return false
	}
	return a.AttrVal(archive.AttrBitIsPorder) == archive.AttrYes
}

func (s *Service) isUGCPay(a *archive.Archive) bool {
	if a == nil {
		return false
	}
	return a.AttrVal(archive.AttrBitUGCPay) == archive.AttrYes
}

func (s *Service) isStaff(a *archive.Archive) bool {
	if a == nil {
		return false
	}
	return a.AttrVal(archive.AttrBitSTAFF) == archive.AttrYes
}

func (s *Service) canDo(mid int64) bool {
	return !s.c.Debug || (s.c.DebugMid == mid && mid > 0)
}

func (s *Service) archiveState(c context.Context, a *archive.Archive, v *archive.Video, ad *archive.AuditParam) (state int8, access int16, attr archive.Attr, forbidAttr *archive.ForbidAttr) {
	// videos
	var (
		vs  []*archive.Video
		err error
	)
	if vs, err = s.arc.NewVideos(c, a.Aid); err != nil {
		log.Error("s.arc.Videos(%d) error(%v)", a.Aid, err)
		return
	}
	if len(vs) == 0 {
		state = archive.StateForbidWait
		log.Warn("archive(%d) have no videos", a.Aid)
		return
	}
	var (
		newState  int8
		newAccess int16
	)
	forbidAttr, _ = s.arc.Forbid(c, a.Aid)
	//聚合状态和属性
	for _, tv := range vs {
		if v != nil && tv.Filename == v.Filename {
			tv = v // NOTE: v maybe change by tran begin, so use v.
		}
		if tv.Status == archive.VideoStatusDelete {
			continue
		}
		if tv.Status == archive.VideoStatusLock {
			newState = archive.StateForbidLock
			break
		}
		if tv.Status == archive.VideoStatusRecicle || newState == archive.StateForbidRecicle {
			newState = archive.StateForbidRecicle
			continue
		}
		if tv.Status == archive.VideoStatusXcodeFail || newState == archive.StateForbidXcodeFail {
			newState = archive.StateForbidXcodeFail
			continue
		}
		if tv.Status == archive.VideoStatusSubmit || newState == archive.StateForbidSubmit {
			newState = archive.StateForbidSubmit
			continue
		}
		if tv.Status == archive.VideoStatusWait || newState == archive.StateForbidSubmit {
			newState = archive.StateForbidSubmit
			continue
		}
		if tv.XcodeState <= archive.VideoXcodeHDFinish || newState == archive.StateForbidWait {
			newState = archive.StateForbidSubmit
			continue
		}
		if tv.XcodeState <= archive.VideoDispatchRunning || newState == archive.StateForbidSubmit {
			newState = archive.StateForbidSubmit
			continue
		}
		if tv.Status == archive.VideoStatusAccess || newAccess == archive.AccessMember {
			newState = archive.StateOpen
			newAccess = archive.AccessMember
		} else if tv.Status == archive.VideoStatusOpen {
			newState = archive.StateOpen
		}
		// attr
		if tv.AttrVal(archive.AttrBitNoRank) == archive.AttrYes || a.AttrVal(archive.AttrBitNoRank) == archive.AttrYes {
			attr.Set(archive.AttrYes, archive.AttrBitNoRank)
			forbidAttr.SetAttr(archive.ForbidRank, archive.AttrYes, archive.ForbidRankMain)
		}
		if tv.AttrVal(archive.AttrBitNoDynamic) == archive.AttrYes || a.AttrVal(archive.AttrBitNoDynamic) == archive.AttrYes {
			attr.Set(archive.AttrYes, archive.AttrBitNoDynamic)
			forbidAttr.SetAttr(archive.ForbidDynamic, archive.AttrYes, archive.ForbidDynamicMain)
		}
		if tv.AttrVal(archive.AttrBitOverseaLock) == archive.AttrYes || a.AttrVal(archive.AttrBitOverseaLock) == archive.AttrYes {
			attr.Set(archive.AttrYes, archive.AttrBitOverseaLock)
			forbidAttr.SetAttr(archive.ForbidShow, archive.AttrYes, archive.ForbidShowOversea)
		}
		if tv.AttrVal(archive.AttrBitNoRecommend) == archive.AttrYes || a.AttrVal(archive.AttrBitNoRecommend) == archive.AttrYes {
			attr.Set(archive.AttrYes, archive.AttrBitNoRecommend)
			forbidAttr.SetAttr(archive.ForbidRecommend, archive.AttrYes, archive.ForbidRecommendMain)
		}
		if tv.AttrVal(archive.AttrBitNoSearch) == archive.AttrYes || a.AttrVal(archive.AttrBitNoSearch) == archive.AttrYes {
			attr.Set(archive.AttrYes, archive.AttrBitNoSearch)
			// NOTE: not search forbit
		}
		if tv.AttrVal(archive.AttrNoPushBplus) == archive.AttrYes || a.AttrVal(archive.AttrNoPushBplus) == archive.AttrYes {
			attr.Set(archive.AttrYes, archive.AttrNoPushBplus)
		}
		if tv.AttrVal(archive.AttrBitParentMode) == archive.AttrYes || a.AttrVal(archive.AttrBitParentMode) == archive.AttrYes {
			attr.Set(archive.AttrYes, archive.AttrBitParentMode)
		}
	}
	if newState == archive.StateOpen {
		if s.hadPassed(c, a.Aid) {
			newState = archive.StateForbidFixed
			if ad != nil && ad.IsAudit {
				newState = archive.StateOpen
			}
		} else if s.isAuditType(a.TypeID) { // 多P或指定分区
			newState = archive.StateForbidWait
		} else if addit, _ := s.arc.Addit(c, a.Aid); addit != nil && (addit.OrderID > 0 || addit.UpFrom == archive.UpFromPGC || addit.UpFrom == archive.UpFromPGCSecret || addit.MissionID > 0) { // NOTE: order || up_from pgc || mission
			newState = archive.StateForbidWait
		} else if hisCnt, _ := s.arc.HistoryCount(c, a.Aid); hisCnt > 1 { // modified before dispatch finish
			newState = archive.StateForbidWait
		} else if delay, err := s.arc.Delay(c, a.Aid); err == nil && delay != nil {
			if delay.DTime.Before(time.Now()) {
				newState = archive.StateForbidWait
			} else if s.isPorder(a) && !s.isAuditType(a.TypeID) { //私单+定时+非特殊分区->newState = -1
				newState = archive.StateForbidWait
			} else if delay.Type == archive.DelayTypeForUser {
				newState = archive.StateForbidUserDelay
			}
		} else if s.isPorder(a) {
			newState = archive.StateForbidWait
		} else if s.isUGCPay(a) {
			//付费稿件在付费待审完成后job才可以自动开放
			newState = archive.StateForbidWait
		}
		if newState == archive.StateOpen && !s.isWhite(a.Mid) && !s.isBlack(a.Mid) {
			if pfl, _ := s.profile(c, a.Mid); pfl != nil && pfl.Follower < int64(s.fansCache) && s.isRoundType(a.TypeID) && !s.isMission(c, a.Aid) {
				newState = archive.StateOrange // NOTE: auto open must
			}
		}
	}
	state = newState
	access = newAccess
	return
}

func (s *Service) archiveRound(c context.Context, a *archive.Archive) (round int8) {
	if archive.NormalState(a.State) {
		//定时发布 或者  自动过审逻辑
		isAuditType := s.isAuditType(a.TypeID)
		if addit, _ := s.arc.Addit(c, a.Aid); addit != nil && (addit.OrderID > 0 || addit.UpFrom == archive.UpFromPGC || addit.UpFrom == archive.UpFromCoopera ||
			addit.MissionID > 0 || (s.isPorder(a) && isAuditType)) {
			round = archive.RoundEnd // NOTE: maybe -40 -> 0   商单稿件||pgc稿件|| 合作方嵌套稿件 || 活动||私单稿件 进round = 99
		} else if addit != nil && addit.UpFrom == archive.UpFromPGCSecret {
			round = archive.RoundTriggerClick // Note: 机密PGC通过后进机密回查
		} else if (a.Round == archive.RoundAuditThird || a.Round == archive.RoundAuditSecond) && isAuditType {
			round = archive.RoundEnd // NOTE: 特殊分区二三审进round = 99
		} else if s.isWhite(a.Mid) || s.isBlack(a.Mid) {
			round = archive.RoundReviewSecond
		} else if pfl, _ := s.profile(c, a.Mid); pfl != nil && pfl.Follower >= int64(s.fansCache) {
			round = archive.RoundReviewSecond
		} else if pfl != nil && pfl.Follower < int64(s.fansCache) && s.isRoundType(a.TypeID) {
			round = archive.RoundReviewFirst // NOTE: if audit type, state must not open!!! so cannot execute here...
		} else {
			round = archive.RoundReviewFirstWaitTrigger
			if pfl == nil {
				log.Info("archive(%d) card(%d) is nil", a.Aid, a.Mid)
			} else {
				log.Info("archive(%d) card(%d) fans(%d) little than config(%d)", a.Aid, a.Mid, pfl.Follower, s.fansCache)
			}
		}
	} else if a.State == archive.StateForbidWait {
		if addit, _ := s.arc.Addit(c, a.Aid); addit != nil && (addit.OrderID > 0 || addit.UpFrom == archive.UpFromPGC || addit.UpFrom == archive.UpFromPGCSecret || addit.UpFrom == archive.UpFromCoopera || addit.MissionID > 0) {
			//二审待审  商单、pgc,活动
			round = archive.RoundAuditSecond
		} else if s.isAuditType(a.TypeID) {
			//指定分区待审逻辑
			// 如果修改未过审  history >1  and not passed (用户修改未过审)走二审
			hadPassed := s.hadPassed(c, a.Aid)
			hasEdit, _ := s.arc.HistoryCount(c, a.Aid)
			if !hadPassed && hasEdit > 1 {
				round = archive.RoundAuditSecond
			} else if count, _ := s.arc.NewVideoCount(c, a.Aid); count == 1 {
				//如果 新增投稿 多p 进二审 单P 进三审
				round = archive.RoundAuditThird
			} else {
				round = archive.RoundAuditSecond
			}
		} else {
			hasEdit, _ := s.arc.HistoryCount(c, a.Aid)
			if s.isPorder(a) && hasEdit <= 1 {
				round = archive.RoundReviewFlow
			} else if s.isUGCPay(a) {
				//付费待审  s -1 r 24
				round = archive.RoundAuditUGCPayFlow
			} else {
				round = archive.RoundAuditSecond
			}
		}
	} else if a.State == archive.StateForbidFixed {
		round = archive.RoundAuditSecond
	} else if a.State == archive.StateForbidUserDelay {
		round = a.Round
	} else if a.State == archive.StateForbidLater {
		round = a.Round
	} else {
		round = archive.RoundBegin
		// NOTE: user delete?? admin later???
	}
	return
}

func (s *Service) tranRound(c context.Context, tx *sql.Tx, a *archive.Archive) (round int8, err error) {
	round = s.archiveRound(c, a)
	if _, err = s.arc.TxUpRound(tx, a.Aid, round); err != nil {
		log.Error("s.arc.TxUpRound(%d, %d) error(%v)", a.Aid, round, err)
		return
	}
	return
}

func (s *Service) tranArchiveOper(tx *sql.Tx, a *archive.Archive) (err error) {
	if _, err = s.arc.TxArchiveOper(tx, a.Aid, a.TypeID, a.State, a.Round, a.Attribute, archive.FirstRoundID, ""); err != nil {
		log.Error("s.arc.TxArchiveOper error(%v)", err)
	}
	return
}

func (s *Service) tranVideo(c context.Context, tx *sql.Tx, a *archive.Archive, v *archive.Video) (err error) {
	//up xcode_state
	if _, err = s.arc.TxUpXcodeState(tx, v.Filename, v.XcodeState); err != nil {
		log.Error("s.arc.TxUpXcodeState(%s, %d) error(%v)", v.Filename, v.XcodeState, err)
		return
	}
	if _, err = s.arc.TxUpVideoXState(tx, v.Filename, v.XcodeState); err != nil {
		log.Error("s.arc.TxUpVideoXState(%s, %d) error(%v)", v.Filename, v.XcodeState, err)
		return
	}
	log.Info("archive(%d) filename(%s) upXcodeState(%d)", a.Aid, v.Filename, v.XcodeState)
	var stCh bool // NOTE: status changed
	// check xcode state
	if v.XcodeState == archive.VideoXcodeSDFail || v.XcodeState == archive.VideoXcodeHDFail {
		if _, err = s.arc.TxUpFailCode(tx, v.Filename, v.FailCode); err != nil {
			log.Error("s.arc.TxUpFailCode(%s, %d, %d) error(%v)", v.Filename, v.FailCode, err)
			return
		}
		if _, err = s.arc.TxUpVideoFailCode(tx, v.Filename, v.FailCode); err != nil {
			log.Error("s.arc.TxUpVideoFailCode(%s, %d, %d) error(%v)", v.Filename, v.FailCode, err)
			return
		}
		log.Info("archive(%d) filename(%s) upFailCode(%d)", a.Aid, v.Filename, v.FailCode)
		stCh = true
	} else if v.XcodeState == archive.VideoXcodeSDFinish {
		if _, err = s.arc.TxUpPlayurl(tx, v.Filename, v.Playurl); err != nil {
			log.Error("s.arc.TxUpPlayurl(%s, %s) error(%v)", v.Filename, v.Playurl, err)
			return
		}
		if _, err = s.arc.TxUpVideoPlayurl(tx, v.Filename, v.Playurl); err != nil {
			log.Error("s.arc.TxUpVideoPlayurl(%s, %s) error(%v)", v.Filename, v.Playurl, err)
			return
		}
		log.Info("archive(%d) filename(%s) upPlayurl(%s)", a.Aid, v.Filename, v.Playurl)
		if _, err = s.arc.TxUpVideoDuration(tx, v.Filename, v.Duration); err != nil {
			log.Error("s.arc.TxUpVideoDuration(%s, %d) error(%v)", v.Filename, v.Duration, err)
			return
		}
		if _, err = s.arc.TxUpVDuration(tx, v.Filename, v.Duration); err != nil {
			log.Error("s.arc.TxUpVDuration(%s, %d) error(%v)", v.Filename, v.Duration, err)
			return
		}
		log.Info("archive(%d) filename(%s) upVdoDuration(%d)", a.Aid, v.Filename, v.Duration)
		stCh = true
	} else if v.XcodeState == archive.VideoXcodeHDFinish {
		if _, err = s.arc.TxUpResolutions(tx, v.Filename, v.Resolutions); err != nil {
			log.Error("s.arc.TxUpResolutions(%s, %s) error(%v)", v.Filename, v.Resolutions, err)
			return
		}
		if _, err = s.arc.TxUpVideoResolutionsAndDimensions(tx, v.Filename, v.Resolutions, v.Dimensions); err != nil {
			log.Error("s.arc.TxUpVideoResolutions(%s,%s, %s) error(%v)", v.Filename, v.Resolutions, v.Dimensions, err)
			return
		}
		log.Info("archive(%d) filename(%s) upResolution(%s)", a.Aid, v.Filename, v.Resolutions)
		if _, err = s.arc.TxUpVideoDuration(tx, v.Filename, v.Duration); err != nil {
			log.Error("s.arc.TxUpVideoDuration(%s, %d) error(%v)", v.Filename, v.Duration, err)
			return
		}
		if _, err = s.arc.TxUpVDuration(tx, v.Filename, v.Duration); err != nil {
			log.Error("s.arc.TxUpVDuration(%s, %d) error(%v)", v.Filename, v.Duration, err)
			return
		}
		log.Info("archive(%d) filename(%s) upVdoDuration(%d)", a.Aid, v.Filename, v.Duration)
		if _, err = s.arc.TxUpFilesize(tx, v.Filename, v.Filesize); err != nil {
			log.Error("s.arc.TxUpFilesize(%s, %d) error(%v)", v.Filename, v.Filesize, err)
			return
		}
		if _, err = s.arc.TxUpVideoFilesize(tx, v.Filename, v.Filesize); err != nil {
			log.Error("s.arc.TxUpVideoFilesize(%s, %d) error(%v)", v.Filename, v.Filesize, err)
			return
		}
		log.Info("archive(%d) filename(%s) upFilesize(%d)", a.Aid, v.Filename, v.Filesize)
	}
	//else if v.XcodeState == archive.VideoDispatchRunning || v.XcodeState == archive.VideoDispatchFinish {
	//	// TODO ???
	//}
	if !stCh {
		return
	}
	if _, err = s.arc.TxUpStatus(tx, v.Filename, v.Status); err != nil {
		log.Error("s.arc.TxUpStatus(%s, %d) error(%v) or rows==0", v.Filename, v.Status, err)
		return
	}
	if v.Status == archive.VideoStatusDelete {
		if _, err = s.arc.TxUpRelationStatus(tx, v.Cid, archive.StateForbidUpDelete); err != nil {
			log.Error("s.arc.TxUpRelationStatus(%d, %d) error(%v) or rows==0", v.Cid, archive.StateOpen, err)
			return
		}
	} else {
		if _, err = s.arc.TxUpVideoStatus(tx, v.Filename, v.Status); err != nil {
			log.Error("s.arc.TxUpVideoStatus(%s, %d) error(%v) or rows==0", v.Filename, v.Status, err)
			return
		}
	}
	// NOTE: reset relation back to active for data consistent. -100 still -100.
	if v.Cid > 0 && v.Status != archive.VideoStatusDelete {
		if _, err = s.arc.TxUpRelationStatus(tx, v.Cid, archive.StateOpen); err != nil {
			log.Error("s.arc.TxUpRelationStatus(%d, %d) error(%v) or rows==0", v.Cid, archive.StateOpen, err)
			return
		}
	}
	log.Info("archive(%d) filename(%s) upStatus(%d)", a.Aid, v.Filename, v.Status)
	var reason string
	if v.Status == archive.VideoStatusXcodeFail {
		if v.XcodeState == archive.VideoXcodeSDFail {
			reason = "转码失败：" + archive.XcodeFailMsgs[v.FailCode]
		} else if v.XcodeState == archive.VideoXcodeHDFail {
			reason = "转码失败：" + archive.XcodeFailMsgs[v.FailCode]
		}
	} else if v.Status == archive.VideoStatusOpen {
		reason = "生产组稿件一转成功"
	} else if v.Status == archive.VideoStatusWait {
		reason = "一转成功"
	}
	if _, err = s.arc.TxAddAudit(tx, v.ID, a.Aid, reason); err != nil {
		log.Error("s.arc.TxAddAudit(%d, %d) filename(%s) reason(%s) error(%v)", v.ID, a.Aid, v.Filename, reason, err)
		return
	}
	log.Info("archive(%d) filename(%s) addAudit reason(%s)", a.Aid, v.Filename, reason)
	return
}

func (s *Service) tranArchive(c context.Context, tx *sql.Tx, a *archive.Archive, v *archive.Video, ad *archive.AuditParam) (change bool, err error) {
	// start archive
	if a.NotAllowUp() {
		log.Warn("archive(%d) filename(%s) state(%d) not allow update", a.Aid, v.Filename, a.State)
		return
	}
	var (
		state, access, attr, forbidAttr = s.archiveState(c, a, v, ad)
		now                             = time.Now()
	)
	if state == a.State {
		log.Warn("archive(%d) filename(%s) newState(%d)==oldState(%d)", a.Aid, v.Filename, state, a.State)
	} else {
		change = true
		firstPass := false
		// archive
		if firstPass, err = s.txUpArcState(c, tx, a.Aid, state); err != nil {
			log.Error("s.txUpArcState(%d, %d) filename(%s) error(%v)", a.Aid, state, v.Filename, err)
			return
		}
		a.State = state
		log.Info("archive(%d) filename(%s) upState(%d)", a.Aid, v.Filename, a.State)
		if firstPass {
			if _, err = s.arc.TxUpPTime(tx, a.Aid, now); err != nil {
				log.Error("s.arc.TxUpPTime(%d, %d) error(%v)", a.Aid, now.Unix(), err)
				return
			}
			a.PTime = xtime.Time(now.Unix())
			log.Info("archive(%d) filename(%s) upPTime(%d)", a.Aid, v.Filename, a.PTime)
		}
	}
	if a.Access != access {
		if _, err = s.arc.TxUpAccess(tx, a.Aid, access); err != nil {
			log.Error("s.arc.TxUpAccess(%d, %d) filename(%s) error(%v)", a.Aid, access, v.Filename, err)
			return
		}
		a.Access = access
		log.Info("archive(%d) filename(%s) upAccess(%d)", a.Aid, v.Filename, a.Access)
	}
	if err = s.tranSumDuration(c, tx, a); err != nil {
		log.Info("s.tranSumDuration error(%v)", err)
		return
	}
	if a.Attribute != (a.Attribute | int32(attr)) {
		if _, err = s.arc.TxUpAttr(tx, a.Aid, attr); err != nil {
			log.Error("s.arc.TxUpAttr(%d, %d) filename(%s) error(%v)", a.Aid, attr, v.Filename, err)
			return
		}
		a.WithAttr(attr)
		log.Info("archive(%d) filename(%s) upAttribute(%d)", a.Aid, v.Filename, a.Attribute)
	}
	if _, err = s.arc.TxUpForbid(tx, forbidAttr); err != nil {
		log.Error("s.arc.TxUpForbid(%+v) error(%v)", forbidAttr, err)
		return
	}
	log.Info("archive(%d) filename(%s) forbidAttr(%+v)", a.Aid, v.Filename, forbidAttr)
	return
}

func (s *Service) tranSumDuration(c context.Context, tx *sql.Tx, a *archive.Archive) (err error) {
	var sum int64
	if sum, err = s.arc.NewSumDuration(c, a.Aid); err != nil {
		log.Error("s.arc.SumDuration(%d) error(%v)", a.Aid, err)
		err = nil
	} else if sum > 0 && a.Duration != sum {
		if _, err = s.arc.TxUpArcDuration(tx, a.Aid, sum); err != nil {
			log.Error("s.arc.TxUpArcDuration(%d, %d) error(%v)", a.Aid, sum, err)
			return
		}
		a.Duration = sum
		log.Info("archive(%d) upArcDuration(%d)", a.Aid, a.Duration)
	}
	return
}

func (s *Service) tranArcCover(c context.Context, tx *sql.Tx, a *archive.Archive, v *archive.Video) (err error) {
	if a.Cover != "" {
		return
	}
	// NOTE: first round need view archive cover, delete when three covers select for user
	var cvs []string
	//从ai处获取封面,若失败/没有封面信息，则直接返回
	if cvs, err = s.arc.AICover(c, v.Filename); err != nil || len(cvs) == 0 {
		log.Error("a.arc.AICover(aid(%d) filename(%s)) got covers from AI error(%v), cvs(%v) ", a.Aid, v.Filename, err, cvs)
		err = nil
		return
	}
	log.Info("s.arc.AICover(aid(%d), filename(%s)) got covers from AI: cvs(%v)", a.Aid, v.Filename, cvs)
	a.Cover = cvs[0] //ai cover只取第一个元素
	if _, err = s.arc.TxUpCover(tx, a.Aid, strings.Replace(a.Cover, "https:", "", -1)); err != nil {
		log.Error("s.arc.TxUpCover(%d, %s) filename(%s) error(%v)", a.Aid, v.Filename, a.Cover, err)
		return
	}
	log.Info("archive(%d) filename(%s) upCover(%s)", a.Aid, v.Filename, a.Cover)
	return
}

func (s *Service) hadPassed(c context.Context, aid int64) (had bool) {
	id, err := s.arc.GetFirstPassByAID(c, aid)
	if err != nil {
		log.Error("hadPassed s.arc.GetFirstPassByAID error(%v) aid(%d)", err, aid)
		return
	}

	had = id > 0
	return
}

func (s *Service) profile(c context.Context, mid int64) (p *accApi.ProfileStatReply, err error) {
	if p, err = s.accRPC.ProfileWithStat3(c, &accApi.MidReq{Mid: mid}); err != nil {
		p = nil
		log.Error("s.accRPC.ProfileWithStat3(%d) error(%v)", mid, err)
	}
	return
}

func (s *Service) changeMission(c context.Context, a *archive.Archive, missionID int64) (err error) {
	if missionID > 0 {
		s.activity.UpVideo(c, a, missionID)
		return
	}
	var addit *archive.Addit
	if addit, err = s.arc.Addit(c, a.Aid); err != nil {
		log.Error("s.arc.Addit(%d) error(%v)", a.Aid, err)
		return
	}
	if addit == nil {
		return
	}
	if addit.MissionID > 0 {
		s.activity.AddVideo(c, a, addit.MissionID)
	}
	return
}

func (s *Service) unBindMission(c context.Context, a *archive.Archive, missionID int64) (err error) {
	if missionID == 0 {
		var addit *archive.Addit
		if addit, err = s.arc.Addit(c, a.Aid); err != nil {
			log.Error("s.arc.Addit(%d) error(%v)", a.Aid, err)
			return
		}
		if addit == nil || addit.MissionID == 0 {
			return
		}
		missionID = addit.MissionID
	}
	s.activity.UpVideo(c, a, missionID)
	return
}

func (s *Service) cidsByAid(c context.Context, aid int64) (cids []int64, err error) {
	var (
		vs []*archive.Video
	)
	if vs, err = s.arc.NewVideos(c, aid); err != nil {
		log.Error("archive(%d) cidsByAid s.arc.Videos error(%v)", aid, err)
		return
	}
	for _, v := range vs {
		cids = append(cids, v.Cid)
	}
	return
}

func (s *Service) syncBVC(c context.Context, a *archive.Archive) (err error) {
	if env.DeployEnv == env.DeployEnvFat1 || env.DeployEnv == env.DeployEnvDev {
		log.Info("archive(%d) syncBVC stop for dev/fat1 env", a.Aid)
		return
	}
	var cids []int64
	if cids, err = s.cidsByAid(c, a.Aid); err != nil {
		log.Error("archive(%d) second_round cidsByAid error(%v)", a.Aid, err)
		s.syncRetry(c, a.Aid, a.Mid, redis.ActionForBvcCapable, "", "")
		return
	}
	var retryCids, okCids, noCids []int64
	for _, cid := range cids {
		//pgc付费 单独播放通道
		//ugc 付费 也是单独播放通道
		//ugc 通道普通视频播放控制
		if a.AttrVal(archive.AttrBitBadgepay) == archive.AttrYes || a.AttrVal(archive.AttrBitIsBangumi) == archive.AttrYes || a.AttrVal(archive.AttrBitUGCPay) == archive.AttrYes {
			noCids = append(noCids, cid)
			continue
		}
		var count int
		// playable videos num
		if count, err = s.arc.NewVideoCountCapable(c, cid); err != nil {
			log.Error("archive(%d) syncBVC  checkCids cid(%d) s.arc.VideoCountCapable error(%v)", a.Aid, cid, err)
			retryCids = append(retryCids, cid)
			continue
		}
		if count == 0 {
			noCids = append(noCids, cid)
		} else if count == 1 {
			if !a.IsNormal() {
				noCids = append(noCids, cid)
			} else {
				okCids = append(okCids, cid)
			}
		} else {
			if a.IsForbid() {
				aids, _ := s.arc.ValidAidByCid(c, cid)
				fcnt := 0
				for _, aid := range aids {
					var arc *archive.Archive
					if arc, err = s.arc.Archive(c, aid); err != nil {
						log.Error("syncBVC get archive (%d) cid(%d) error(%v)", aid, cid, err)
						break
					}
					if arc != nil && arc.State >= 0 {
						fcnt++
						break
					}
				}
				log.Info("syncBVC ValidAidByCid cid(%d) aids(%v) fcnt(%d)", cid, aids, fcnt)
				// num of cids used in other open state archive
				if err != nil {
					log.Error("checkCids cid(%d) error(%v)", cid, err)
					retryCids = append(retryCids, cid)
					continue
				} else if fcnt == 0 {
					noCids = append(noCids, cid)
					continue
				}
				// NOTE: when fcnt>0, means cid has normal archive.
			}
			okCids = append(okCids, cid)
		}
	}
	var flagRetry = false
	if len(retryCids) != 0 {
		log.Warn("syncBVC aid(%d) cids(%v) need retry again", a.Aid, retryCids)
		flagRetry = true
	}
	if len(okCids) != 0 {
		if err = s.bvc.VideoCapable(c, a.Aid, okCids, message.CanPlay); err != nil {
			log.Error("syncBVC  aid(%d) cids(%v) s.bvc.VideoCapable error(%v)", a.Aid, okCids, err)
			flagRetry = true
		}
	}
	if len(noCids) != 0 {
		if err = s.bvc.VideoCapable(c, a.Aid, noCids, message.CanNotPlay); err != nil {
			log.Error("syncBVC aid(%d) cids(%v) s.bvc.VideoCapable error(%v)", a.Aid, noCids, err)
			flagRetry = true
		}
	}
	if flagRetry {
		s.syncRetry(c, a.Aid, a.Mid, redis.ActionForBvcCapable, "", "")
	}
	return
}

//txAddFirstPass 添加第一次过审记录
func (s *Service) txAddFirstPass(c context.Context, tx *sql.Tx, aid int64, state int8) (firstPass bool, err error) {
	if !archive.NormalState(state) || s.hadPassed(c, aid) {
		return
	}

	if err = s.arc.AddFirstPass(tx, aid); err != nil {
		log.Error("txAddFirstPass error(%v) aid(%d)", err, aid)
		return
	}
	firstPass = true
	return
}

//txUpArcState 更新稿件的state并联动添加第一次过审记录
func (s *Service) txUpArcState(c context.Context, tx *sql.Tx, aid int64, state int8) (firstPass bool, err error) {
	if _, err = s.arc.TxUpState(tx, aid, state); err != nil {
		log.Error("txUpArcState s.arc.TxUpState error(%v) aid(%d) state(%d)", err, aid, state)
		return
	}

	if firstPass, err = s.txAddFirstPass(c, tx, aid, state); err != nil {
		log.Error("txUpArcState s.txAddFirstPass error(%v) aid(%d) state(%d)", err, aid, state)
		return
	}
	return
}

// IsUpperFirstPass 是否UP主第一次过审稿件
func (s *Service) IsUpperFirstPass(c context.Context, mid, aid int64) (is bool, err error) {
	is = true
	sMap, err := s.arc.UpperArcStateMap(c, mid)
	if err != nil {
		log.Error("s.arc.UpperArcStateMap(%d,%d) error(%v)", mid, aid, err)
		return
	}
	delete(sMap, aid) //剔除当前稿件
	var aids []int64
	for k, v := range sMap {
		aids = append(aids, k)
		if archive.NormalState(v) {
			//如果有其它过审的稿件，那么当前UP主肯定不是第一次过审
			is = false
			break
		}
	}
	if is {
		//查询first_pass表
		var count int
		count, err = s.arc.FirstPassCount(c, aids)
		if err != nil {
			log.Error("s.arc.FirstPassCount(%v) error(%v)", aids, err)
			return
		}
		is = count == 0
	}
	return
}
