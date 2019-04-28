package service

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-common/app/interface/main/reply/conf"
	"go-common/app/interface/main/reply/model/adminlog"
	"go-common/app/interface/main/reply/model/drawyoo"
	"go-common/app/interface/main/reply/model/reply"
	accmdl "go-common/app/service/main/account/api"
	"go-common/app/service/main/archive/api"
	arcmdl "go-common/app/service/main/archive/model/archive"
	assmdl "go-common/app/service/main/assist/model/assist"
	filgrpc "go-common/app/service/main/filter/api/grpc/v1"
	locmdl "go-common/app/service/main/location/model"
	relmdl "go-common/app/service/main/relation/model"
	thumdl "go-common/app/service/main/thumbup/model"
	ugcpay "go-common/app/service/main/ugcpay/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	xip "go-common/library/net/ip"
	"go-common/library/net/metadata"
	"go-common/library/queue/databus/report"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"

	"github.com/mvdan/xurls"
)

var (
	_emptyReplies   = make([]*reply.Reply, 0)
	_emptyAction    = map[int64]int8{}
	_emptyCards     = make(map[int64]*accmdl.Card)
	_emptyBlackList = make(map[int64]bool)
	_emptyRelations = make(map[int64]*accmdl.RelationReply)

	_emojiCode = regexp.MustCompile(`\[[^\[+][^]]+]`)
)

// status
const (
	StatusNormal      = 1
	StatusNeedContest = 2
	StatusForbidden   = 3
)

// IsWhiteAid IsWhiteAid
func (s *Service) IsWhiteAid(aid int64, tp int8) bool {
	if tp != 1 {
		return false
	}
	for _, white := range s.aidWhiteList {
		if aid == white {
			return true
		}
	}
	return false
}

// UserBlockStatus UserBlockStatus
func (s *Service) UserBlockStatus(c context.Context, mid int64) (int, error) {
	res, err := s.dao.BlockStatus.BlockInfo(c, mid)
	if err != nil {
		return 0, err
	}
	if res.ForeverBlock || time.Now().Unix() < res.BlockUntil {
		return StatusForbidden, nil
	}
	if res.PassTest == 0 {
		return StatusNormal, nil
	}
	return StatusNeedContest, nil
}

// ValidUserStatus validate reply user status
func (s *Service) ValidUserStatus(c context.Context, profile *accmdl.Profile, isUpper bool) (err error) {
	// if myInfo.Silence == 1 {
	// 	err = ecode.UserDisabled
	// } else if myInfo.Active == 0 {
	// 	err = ecode.UserInactive
	// } else if myInfo.Moral < 60 {
	// 	err = ecode.LackOfScores
	// } else if myInfo.Rank == 5000 {
	// 	err = ecode.UserNoMember
	// } else if myInfo.Level.Cur < 1 {
	// 	err = ecode.UserLevelLow
	// }
	if profile.Silence == 1 {
		err = ecode.UserDisabled
	} else if profile.TelStatus == 0 && profile.EmailStatus == 0 {
		err = ecode.UserInactive
	} else if profile.Moral < 60 {
		err = ecode.LackOfScores
	} else if profile.Rank == 5000 && !isUpper {
		err = ecode.UserNoMember
	} else if profile.Level < 1 && !isUpper {
		err = ecode.UserLevelLow
	}
	return
}

// ValidUserAction validate reply user status for like/hate action.
func (s *Service) ValidUserAction(c context.Context, profile *accmdl.Profile) (err error) {
	if profile.Silence == 1 {
		err = ecode.UserDisabled
	} else if profile.TelStatus == 0 && profile.EmailStatus == 0 {
		err = ecode.UserInactive
	} else if profile.Moral < 60 {
		err = ecode.LackOfScores
	}
	return
}

// checkSpam detemine whether user can reply or not
func (s *Service) checkSpam(c context.Context, sub *reply.Subject, mid int64, captcha string, level int) (uri string, err error) {
	if sub.Type != reply.SubTypeBBQ && sub.Type != reply.SubTypeHuoniao {
		if level <= reply.UserLevelFirst && sub.Mid != mid {
			if captcha == "" {
				var uri string
				uri, err = s.Captcha(c, mid)
				if err != nil {
					return "", err
				}
				return uri, ecode.ReplyDeniedAsCaptcha
			} else if err = s.VerifyCaptcha(c, captcha, mid); err != nil {
				return "", err
			}
		}
	}
	recent, daily, err := s.dao.Redis.SpamReply(c, mid)
	if err != nil {
		log.Error("replyCacheDao.SpamReply(%d), err (%v)", mid, err)
		return "", err
	}
	if recent == ecode.ReplyDeniedAsCD.Code() || daily == ecode.ReplyDeniedAsCD.Code() {
		return "", ecode.ReplyDeniedAsCD
	}
	if recent == ecode.ReplyDeniedAsCaptcha.Code() || daily == ecode.ReplyDeniedAsCaptcha.Code() {
		if captcha == "" {
			uri, err := s.Captcha(c, mid)
			if err != nil {
				return "", err
			}
			return uri, ecode.ReplyDeniedAsCaptcha
		}
		if err := s.VerifyCaptcha(c, captcha, mid); err != nil {
			return "", err
		}
		s.dao.Redis.DelReplyIncr(c, mid, sub.Mid == mid)
		s.dao.Redis.DelReplySpam(c, mid)
	}
	s.dao.Databus.AddSpam(c, sub.Oid, mid, sub.Mid == mid, sub.Type)
	return "", nil
}

func (s *Service) isNormalVip(c context.Context, profile *accmdl.Profile) bool {
	return profile.Vip.Type != 0 && profile.Vip.Status == 1
}

// ContainUrls ContainUrls
func (s *Service) ContainUrls(msg string) bool {
	return xurls.Strict.FindAllString(msg, -1) != nil
}

// bigDataFilter check content by big data and find conmment garbage
func (s *Service) bigDataFilter(c context.Context, msg string) (err error) {
	if err = s.bigdata.Filter(c, msg); err != nil {
		log.Error("s.bigdata.Filter(%s) error(%v)", msg, err)
	}
	return
}

func (s *Service) isUpper(c context.Context, mid, oid int64, tp int8) bool {
	sub, err := s.getSubject(c, oid, tp)
	if err != nil {
		return false
	}
	return sub.Mid == mid
}

// CheckAssist check whether upper grant the supervision permission for user
func (s *Service) CheckAssist(c context.Context, mid, uid int64) (assisted bool, operation bool) {
	arg := &assmdl.ArgAssist{
		Mid:       mid,
		AssistMid: uid,
		Type:      1,
		RealIP:    "",
	}
	if respro, _ := s.assist.Assist(c, arg); respro == nil {
		log.Error("s.assist.Assist(%d, %d) error(%v)", mid, uid, "获取up协管关系错误")
	} else if respro.Assist == 1 {
		assisted = true
		if respro.Allow == 1 {
			operation = true
		}
	}
	return assisted, operation
}

// getAssistList fetch all assistants of user mid
func (s *Service) getAssistList(c context.Context, mid int64) (assistMap map[int64]int) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &assmdl.ArgAssists{
		Mid:    mid,
		RealIP: ip,
	}
	assistMap = make(map[int64]int)
	if response, err := s.assist.AssistIDs(c, arg); err != nil {
		log.Error("s.assist.Assists(%d) error(%v)", mid, err)
	} else {
		for _, tmp := range response {
			assistMap[tmp] = 1
		}
	}
	return
}

// checkContentFilter2 check content by word filter and minus moral when this be filtered.
func (s *Service) checkContentFilter2(c context.Context, oid, mid, rpid int64, ip, msg string, tp int8) (correct string, err error) {
	arg := &filgrpc.FilterReq{
		Message: msg,
		Area:    "reply",
		Id:      rpid,
		Oid:     oid,
		Mid:     mid,
	}
	var res *filgrpc.FilterReply
	res, err = s.filcli.Filter(c, arg)
	if err != nil {
		log.Error("s.filter.Filter(%s) error(%v)", msg, err)
		return msg, err
	}
	switch int(res.Limit) {
	case ecode.FilterHitLimitBlack.Code():
		log.Info("Reply id %d, content %q contains sensitive msg, not allowed to send out", rpid, msg)
		err = ecode.ReplyHitBlacklist

	case ecode.FilterHitRubLimit.Code():
		log.Info("Reply id %d, content %q was sent too many times, exceed allowed counts", rpid, msg)
		err = ecode.ReplyOverRateLimit

	case ecode.FilterHitStrictLimit.Code():
		if res.Level == 0 {
			log.Info("Reply id %d, content %q was limit strictly", rpid, msg)
			err = ecode.ReplyDeniedAsCaptcha
		}
	}
	if err != nil {
		return msg, err
	}

	if res.Level > 0 {
		s.cache.Do(c, func(ctx context.Context) {
			s.AddFilteredReply(ctx, rpid, oid, mid, tp, int8(res.Level), msg, time.Now())
		})
		switch res.Level {
		case 10, 20:
			err = ecode.ReplyMosaicByFilter
		case 30:
			err = ecode.ReplyDeniedByFilter
			return
		case 40:
			tmp := []rune(msg)
			if len(tmp) > 80 {
				tmp = tmp[:80]
			}
			arg := &accmdl.MoralReq{
				Mid:    mid,
				Moral:  -1,
				Oper:   "",
				Reason: "发布恶意评论: " + string(tmp),
				Remark: "云屏蔽",
				RealIp: ip,
			}
			if _, err = s.acc.AddMoral3(c, arg); err != nil {
				log.Error("s.acc.AddMoral3(%d) error(%v)", mid, err)
				return
			}
			err = ecode.ReplyDeniedByFilter
			return
		}
	}
	correct = res.Result
	return
}

// AddFilteredReply AddFilteredReply
func (s *Service) AddFilteredReply(c context.Context, rpID, oid, mid int64, tp, level int8, message string, now time.Time) (err error) {
	return s.dao.Reply.AddFilteredReply(c, rpID, oid, mid, tp, level, message, now)
}

// UseBigdata use bigdata switch
func (s *Service) UseBigdata(c context.Context, b bool, per int64) bool {
	s.useBigData = b
	return s.useBigData
}

// AddReply add a reply.
func (s *Service) AddReply(c context.Context, mid, oid int64, tp, plat int8, ats []int64, accessKey, cookie, captcha, msg, dev, ver, platform string, build int64, buvid string) (r *reply.Reply, uri string, err error) {
	var (
		rootID, parentID, dialog int64
		profile                  *accmdl.Profile
		subject                  *reply.Subject
	)
	//whitelist for test
	profile, subject, uri, err = s.validateReply(c, mid, oid, tp, captcha, msg, accessKey, cookie)
	if err != nil {
		return
	}

	// check content contain emoji code
	if emoCodes := _emojiCode.FindAllString(msg, -1); len(emoCodes) > 0 {
		if s.isNormalVip(c, profile) {
			if len(emoCodes) > conf.Conf.Reply.MaxEmoji {
				err = ecode.ReplyEmojiOverMax
				return
			}
			needRepressEmoCodes := make([]string, 0)
			for _, emoCode := range emoCodes {
				if _, ok := s.emojisM[emoCode]; !ok {
					needRepressEmoCodes = append(needRepressEmoCodes, emoCode)
				}
			}
			if len(needRepressEmoCodes) > 0 {
				msg = RepressEmotions(msg, needRepressEmoCodes)
			}
		} else {
			msg = RepressEmotions(msg, emoCodes)
		}
	}
	if err = s.SuperviseReply(c, mid, accessKey, cookie, int8(tp)); err != nil {
		return
	}
	r, err = s.persistReply(c, mid, rootID, parentID, plat, tp, ats, msg, dev, ver, captcha, platform, build, buvid, subject, dialog)
	if err == ecode.ReplyDeniedAsCaptcha {
		uri, err := s.Captcha(c, mid)
		if err != nil {
			return r, "", err
		}
		return r, uri, ecode.ReplyDeniedAsCaptcha
	}
	return
}

// RepressEmotions RepressEmotions
func RepressEmotions(msg string, emoCodes []string) string {
	for _, emoCode := range emoCodes {
		msg = repressEmotion(msg, emoCode)
	}
	return msg
}

func repressEmotion(msg, emoCode string) string {
	// replace [] to 【】
	emoCode = emoCode[1 : len(emoCode)-1]
	return strings.Replace(msg, "["+emoCode+"]", "【"+emoCode+"】", -1)
}

// AddReplyReply add reply to a root reply.
func (s *Service) AddReplyReply(c context.Context, mid, oid, rootID, parentID int64, tp, plat int8, ats []int64, accessKey, cookie, captcha, msg, dev, ver, platform string, build int64, buvid string) (r *reply.Reply, uri string, err error) {
	var dialog int64
	var profile *accmdl.Profile
	var subject *reply.Subject
	profile, subject, uri, err = s.validateReply(c, mid, oid, tp, captcha, msg, accessKey, cookie)
	if err != nil {
		return
	}
	if emoCodes := _emojiCode.FindAllString(msg, -1); len(emoCodes) > 0 {
		if s.isNormalVip(c, profile) {
			if len(emoCodes) > conf.Conf.Reply.MaxEmoji {
				err = ecode.ReplyEmojiOverMax
				return
			}
			needRepressEmoCodes := make([]string, 0)
			for _, emoCode := range emoCodes {
				if _, ok := s.emojisM[emoCode]; !ok {
					needRepressEmoCodes = append(needRepressEmoCodes, emoCode)
				}
			}
			if len(needRepressEmoCodes) > 0 {
				msg = RepressEmotions(msg, needRepressEmoCodes)
			}
		} else {
			msg = RepressEmotions(msg, emoCodes)
		}
	}
	rootReply, err := s.GetRootReply(c, oid, rootID, tp)
	if err != nil {
		return
	}
	// NOTE if the pending reply, the state is not normal
	if rootReply.IsDeleted() {
		err = ecode.ReplyDeleted
		return
	}
	if s.RelationBlocked(c, rootReply.Mid, mid) {
		err = ecode.ReplyBlacklistFilter
		return
	}
	if err = s.SuperviseReply(c, mid, accessKey, cookie, int8(tp)); err != nil {
		return
	}
	if rootID != parentID {
		var parentReply *reply.Reply
		if parentReply, err = s.GetReply(c, oid, parentID, tp); err != nil {
			return
		}
		// if parentReply.Dialog == 0 {
		// 	s.dao.Databus.RecoverFixDialogIdx(c, oid, tp, rootID)
		// }
		dialog = parentReply.Dialog
		if parentReply.Root != rootID {
			err = ecode.ReplyIllegalRoot
			return
		}
		if mid != parentReply.Mid && !parentReply.IsNormal() {
			err = ecode.ReplyNotExist
			return
		}
		if s.RelationBlocked(c, parentReply.Mid, mid) {
			err = ecode.ReplyBlacklistFilter
			return
		}
	}
	r, err = s.persistReply(c, mid, rootID, parentID, plat, tp, ats, msg, dev, ver, captcha, platform, build, buvid, subject, dialog)
	if err == ecode.ReplyDeniedAsCaptcha {
		uri, err := s.Captcha(c, mid)
		if err != nil {
			return r, "", err
		}
		return r, uri, ecode.ReplyDeniedAsCaptcha
	}
	return
}

func (s *Service) validateReply(c context.Context, mid, oid int64, tp int8, captcha, msg, accessKey, cookie string) (profile *accmdl.Profile, subject *reply.Subject, uri string, err error) {
	profile, err = s.userInfo(c, mid)
	if err != nil {
		log.Error("myinfo(%d) error(%v)", mid, err)
		return nil, nil, "", err
	}
	if tp != reply.SubTypeBBQ {
		if conf.Conf.Identification.SwitchOn && profile.Identification == 0 {
			if profile.TelStatus == 0 {
				err = ecode.UserCheckNoPhone
				return
			}
			if profile.TelStatus == 2 && profile.Identification == 0 {
				err = ecode.UserCheckInvalidPhone
				return
			}
		}
	}
	subject, err = s.Subject(c, oid, tp)
	if err != nil {
		return
	}
	if tp != reply.SubTypeBBQ && tp != reply.SubTypeHuoniao {
		if err = s.ValidUserStatus(c, profile, subject.Mid == mid); err != nil {
			log.Warn("s.ValidUserStatus(%d,%+v) error(%v)", mid, profile.Level, err)
			return
		}
	}
	if s.RelationBlocked(c, subject.Mid, mid) {
		err = ecode.ReplyBlacklistFilter
		return
	}
	if tp != reply.SubTypeBBQ && tp != reply.SubTypeHuoniao {
		if profile.Level < reply.UserLevelSnd && subject.Mid != mid {
			err = ecode.UserLevelLow
			return
		}
	}
	if mid != 165252 && mid != 10287644 {
		if uri, err = s.checkSpam(c, subject, mid, captcha, int(profile.Level)); err != nil {
			log.Error("s.checkSpam failed(%d) err is %V", mid, err)
			return
		}
	}
	if tp == reply.SubTypeArchive && subject.Mid != mid {
		var arc *api.Arc
		arc, err = s.arcSrv.Archive3(c, &arcmdl.ArgAid2{Aid: oid})
		if err != nil {
			log.Error("s.arcSrc.Archive3(%d) failed!err:=%v", oid, err)
			return
		}
		if arc.Rights.UGCPay == 1 {
			var relation *ugcpay.AssetRelationResp
			relation, err = s.ugcpay.AssetRelation(c, &ugcpay.AssetRelationReq{Mid: mid, Oid: oid, Otype: "archive"})
			if err != nil {
				log.Error("s.ugcpay.AssetRelation(%d,%d) failed!err:=%v", mid, oid, err)
				return
			}
			if relation.State != "paid" {
				err = ecode.ReplyForbidReplyNotPay
				return
			}
		}
	}
	return
}

func (s *Service) persistReply(c context.Context, mid, root, parent int64, plat, tp int8, ats []int64, msg, dev, ver, captcha, platform string, build int64, buvid string, subject *reply.Subject, dialog int64) (r *reply.Reply, err error) {
	rpID, err := s.nextID(c)
	if err != nil {
		return
	}
	// 一级子评论
	if root == parent && root != 0 {
		dialog = rpID
	} else if root != parent {
		parentRp, err := s.reply(c, mid, subject.Oid, parent, tp)
		if err != nil {
			return nil, err
		}
		dialog = parentRp.Dialog
	}
	cTime := xtime.Time(time.Now().Unix())
	ip := metadata.String(c, metadata.RemoteIP)
	port := metadata.String(c, metadata.RemotePort)
	r = &reply.Reply{
		RpID:   rpID,
		Oid:    subject.Oid,
		Type:   tp,
		Mid:    mid,
		Root:   root,
		State:  reply.ReplyStateNormal,
		Parent: parent,
		CTime:  cTime,
		Dialog: dialog,
		Content: &reply.Content{
			RpID:    rpID,
			Message: msg,
			Ats:     ats,
			IP:      xip.InetAtoN(ip),
			Plat:    plat,
			Device:  dev,
			Version: ver,
			CTime:   cTime,
		},
	}
	if s.useBigData {
		if err = s.bigDataFilter(c, msg); err != nil {
			if err == ecode.ReplyDeniedAsGarbage {
				// TODO: do not use garbage as state
				r.State = reply.ReplyStateGarbage
				r.AttrSet(reply.AttrYes, reply.ReplyAttrGarbage)
			}
		}
	}
	// if not rpID passed, then no data will be recorded
	msg, err = s.checkContentFilter2(c, r.Oid, mid, rpID, ip, msg, r.Type)
	if err != nil {
		if err != ecode.ReplyDeniedAsCaptcha && err != ecode.ReplyMosaicByFilter {
			log.Error("s.checkContentFilter2(%d, %d, msg: %s) error(%v)", mid, subject.Oid, msg, err)
			return
		}
		if err == ecode.ReplyHitBlacklist || err == ecode.ReplyOverRateLimit {
			return
		}
		if err == ecode.ReplyDeniedAsCaptcha {
			if captcha == "" {
				return
			}
			if err = s.VerifyCaptcha(c, captcha, mid); err != nil {
				return
			}
		} else {
			r.Content.Message = msg
			r.AttrSet(reply.AttrYes, reply.ReplyAttrFilter)
			r.State = reply.ReplyStateFiltered
		}
	}
	// NOTE audit pending most priority
	if subject.AttrVal(reply.SubAttrMonitor) == reply.AttrYes {
		r.State = reply.ReplyStateMonitor
	}
	if subject.AttrVal(reply.SubAttrAudit) == reply.AttrYes {
		r.State = reply.ReplyStateAudit
	}
	s.dao.Databus.AddReply(c, subject.Oid, r)
	report.User(&report.UserInfo{
		Mid:      r.Mid,
		Platform: platform,
		Build:    build,
		Buvid:    buvid,
		Business: 41,
		Type:     int(r.Type),
		Oid:      r.Oid,
		Action:   reply.ReportReplyAdd,
		Ctime:    time.Now(),
		IP:       ip + ":" + port,
		Index: []interface{}{
			r.RpID,
			r.State,
			r.State,
			strconv.FormatInt(r.Root, 10),
		},
	})
	return
}

// checkUpSpam determine user can add up.
func (s *Service) checkActionSpam(c context.Context, mid int64) (err error) {
	var ret int
	if ret, err = s.dao.Redis.SpamAction(c, mid); err != nil {
		log.Error("replyCacheDao.SpamAction(%d), err (%v)", mid, err)
	} else {
		if ret != ecode.OK.Code() {
			err = ecode.ReplyForbidAction
		}
	}
	return
}

// AddAction do act or cancel act for a reply.
func (s *Service) AddAction(c context.Context, mid, oid, rpID int64, tp, action int8, ak, ck, op, platform, buvid string, build int64) (err error) {
	if err = reply.CheckAction(action); err != nil {
		return
	}
	user, err := s.userInfo(c, mid)
	if err != nil {
		return
	}
	if err = s.ValidUserAction(c, user); err != nil {
		return
	}
	if err = s.checkActionSpam(c, mid); err != nil {
		log.Error("s.checkActionSpam(%d) err (%v)", mid, err)
		return
	}
	r, err := s.reply(c, mid, oid, rpID, tp)
	if err != nil {
		return
	}
	// NOTE if the pending reply, the state is not normal
	if mid != r.Mid && !r.IsNormal() {
		err = ecode.ReplyForbidAction
		return
	}
	if s.RelationBlocked(c, r.Mid, mid) {
		err = ecode.ReplyBlacklistFilter
		return
	}
	var (
		userLikes map[int64]int8
		act       int8
	)
	if userLikes, err = s.thumbup.HasLike(c, &thumdl.ArgHasLike{Business: "reply", MessageIDs: []int64{rpID}, Mid: mid}); err != nil {
		log.Error("s.thumbup.HasLike(%d,%d,%d) error(%v)", mid, rpID, oid, err)
		return
	}
	act = userLikes[rpID]
	now := time.Now()
	remoteIP := metadata.String(c, metadata.RemoteIP)
	var ac string
	if op == "like" {
		if (int8(act) == reply.ActionLike && action == reply.OpAdd) || (int8(act) != reply.ActionLike && action == reply.OpCancel) {
			err = ecode.ReplyActioned
			return
		}
		if action == reply.OpAdd {
			ac = reply.ReportReplyLike
			err = s.thumbup.Like(c, &thumdl.ArgLike{UpMid: r.Mid, Business: "reply", Mid: mid, MessageID: rpID, Type: thumdl.TypeLike, RealIP: remoteIP, OriginID: oid})
		} else {
			ac = reply.ReportReplyCancelLike
			err = s.thumbup.Like(c, &thumdl.ArgLike{UpMid: r.Mid, Business: "reply", Mid: mid, MessageID: rpID, Type: thumdl.TypeCancelLike, RealIP: remoteIP, OriginID: oid})
		}
		if err == nil {
			s.dao.Databus.Like(c, oid, rpID, mid, action, now.Unix())
		}
	} else {
		if (int8(act) == reply.ActionHate && action == reply.OpAdd) || (int8(act) != reply.ActionHate && action == reply.OpCancel) {
			err = ecode.ReplyActioned
			return
		}
		if action == reply.OpAdd {
			ac = reply.ReportReplyHate
			err = s.thumbup.Like(c, &thumdl.ArgLike{UpMid: r.Mid, Business: "reply", Mid: mid, MessageID: rpID, Type: thumdl.TypeDislike, RealIP: remoteIP, OriginID: oid})
		} else {
			ac = reply.ReportReplyCancelHate
			err = s.thumbup.Like(c, &thumdl.ArgLike{UpMid: r.Mid, Business: "reply", Mid: mid, MessageID: rpID, Type: thumdl.TypeCancelDislike, RealIP: remoteIP, OriginID: oid})
		}
		if err == nil {
			s.dao.Databus.Hate(c, oid, rpID, mid, action, now.Unix())
		}
	}
	if err != nil {
		if ecode.ThumbupCancelDislikeErr.Equal(err) || ecode.ThumbupCancelLikeErr.Equal(err) || ecode.ThumbupDupLikeErr.Equal(err) || ecode.ThumbupDupDislikeErr.Equal(err) {
			err = nil
			return
		}
		log.Error("thumbup (%d,%d,%d,%s,%d) failed!err:=%v", mid, oid, rpID, op, action, err)
		return
	}
	err = s.infoc.Info(mid, platform, build, buvid, 41, int(r.Type), r.Oid, ac, remoteIP, time.Now().Format("2006-01-02 15:04:05"), r.RpID, r.Mid, "", "", "", "", "")
	if err != nil {
		log.Error("infoc error (%v)", err)
	}
	return
}

func (s *Service) getIdsByRoots(c context.Context, oid int64, roots []int64, tp int8, pn, ps int) (sidsmap map[int64][]int64, ids []int64, err error) {
	var (
		start    = (pn - 1) * ps
		end      = start + ps - 1
		miss     []int64
		tmprpIDs []int64
	)
	if sidsmap, ids, miss, err = s.dao.Redis.RangeByRoots(c, roots, start, end); err != nil {
		log.Error("s.dao.Redis.RangeByRoots() err(%v)", err)
		return
	}
	if len(miss) == 0 {
		return
	}
	for _, root := range miss {
		if tmprpIDs, err = s.dao.Reply.GetIdsByRoot(c, oid, root, tp, start, ps); err != nil {
			log.Error("s.dao.Reply.GetIdsByRoot(oid %d,tp %d,root %d) err(%v)", oid, tp, root, err)
		}
		if len(tmprpIDs) != 0 {
			sidsmap[root] = tmprpIDs
			ids = append(ids, tmprpIDs...)
			s.dao.Databus.RecoverIndexByRoot(c, oid, root, tp)
		}
	}
	return
}

func (s *Service) actions(c context.Context, mid, oid int64, rpIDs []int64) (amap map[int64]int8, err error) {
	if mid == 0 {
		amap = _emptyAction
		return
	}
	amap, err = s.thumbup.HasLike(c, &thumdl.ArgHasLike{Business: "reply", MessageIDs: rpIDs, Mid: mid})
	if err != nil {
		log.Error("s.thumbup.HasLike(%d, %d) error(%v)", mid, rpIDs, err)
		return
	}
	// NOTE: may have many keys
	// if mid not action,add -1 as mark
	if len(amap) == 0 {
		amap = map[int64]int8{-1: 0}
	}
	return
}

// getAccInfo get account infos of mids
func (s *Service) getAccInfo(c context.Context, mids []int64) (cards map[int64]*accmdl.Card, err error) {
	if len(mids) == 0 {
		cards = _emptyCards
		return
	}
	var cardsReply *accmdl.CardsReply
	if cardsReply, err = s.acc.Cards3(c, &accmdl.MidsReq{Mids: mids}); err != nil {
		log.Error("s.acc.MultiInfo2(%v) error(%v)", mids, err)
		return nil, err
	}
	cards = cardsReply.Cards
	return
}

// GetBlacklist get account infos of mids
func (s *Service) GetBlacklist(c context.Context, mid int64) (blacklistMap map[int64]bool, err error) {
	if mid == 0 {
		blacklistMap = _emptyBlackList
		return
	}
	var blacksReply *accmdl.BlacksReply
	if blacksReply, err = s.acc.Blacks3(c, &accmdl.MidReq{Mid: mid}); err != nil {
		log.Error("s.acc.Blacks(%v) error(%v)", mid, err)
		return
	}
	blacklistMap = blacksReply.BlackList
	return
}

// GetAttentions get relationships whether the user(mid) follows the target reply users
func (s *Service) getAttentions(c context.Context, mid int64, targetMids []int64) (relations map[int64]*accmdl.RelationReply, err error) {
	if len(targetMids) == 0 {
		relations = _emptyRelations
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	var relationsReply *accmdl.RelationsReply
	if relationsReply, err = s.acc.Relations3(c, &accmdl.RelationsReq{Mid: mid, Owners: targetMids, RealIp: ip}); err != nil {
		log.Error("s.acc.Relations2(%v, %v) error(%v)", mid, targetMids, err)
		return
	}
	relations = relationsReply.Relations
	return
}

// Subject get normal state reply subject
func (s *Service) Subject(c context.Context, oid int64, tp int8) (*reply.Subject, error) {
	subject, err := s.getSubject(c, oid, tp)
	if err != nil {
		return nil, err
	}
	if subject.State == reply.SubStateForbid {
		return nil, ecode.ReplyForbidReply
	}
	return subject, nil
}

func (s *Service) getSubject(c context.Context, oid int64, tp int8) (*reply.Subject, error) {
	if !reply.LegalSubjectType(tp) {
		log.Error("illegal subject type: %v", tp)
		return nil, ecode.ReplyIllegalSubType
	}
	sub, err := s.dao.Mc.GetSubject(c, oid, tp)
	if err != nil {
		log.Error("replyCacheDao.GetSubject(%d, %d) error(%v)", oid, tp, err)
	}
	if sub != nil {
		return sub, nil
	}
	sub, err = s.dao.Subject.Get(c, oid, tp)
	if err != nil {
		log.Error("s.subject.Get(%d, %d) error(%v)", oid, tp, err)
	}
	if err == nil && sub != nil {
		s.dao.Mc.AddSubject(c, sub)
		return sub, nil
	}
	// fetch from remote call
	if tp != reply.SubTypeDrawyoo {
		log.Error("subject type is nether topic nor drawyoo: %v", tp)
		return nil, ecode.ReplyForbidReply
	}
	var mid int64
	if tp == reply.SubTypeDrawyoo {
		var yoo *drawyoo.Drawyoo
		if yoo, err = s.drawyoo.Info(c, oid); err != nil || yoo == nil {
			log.Warn("drawtyoo.DrawInfo(%d) not exist", oid)
			err = ecode.ReplyForbidReply
			return nil, err
		}
		mid = yoo.Mid
	}
	sub, err = s.upsertSubject(c, oid, tp, reply.SubStateNormal, mid)
	if err == ecode.ReplySubjectExist {
		sub, err = s.dao.Subject.Get(c, oid, tp)
		if err != nil {
			log.Error("s.subject.Get(%d, %d) error(%v)", oid, tp, err)
			return nil, err
		}
		return sub, nil
	}
	return sub, err
}

// upsertSubject insert or update a subject.
func (s *Service) upsertSubject(c context.Context, oid int64, tp, state int8, mid int64) (sub *reply.Subject, err error) {
	now := time.Now()
	sub = &reply.Subject{
		Oid:   oid,
		Type:  tp,
		Mid:   mid,
		State: state,
		CTime: xtime.Time(now.Unix()),
		MTime: xtime.Time(now.Unix()),
	}
	sub.ID, err = s.dao.Subject.Set(c, sub)
	if err != nil {
		log.Error("s.subject.Insert(%s) error(%v)", sub, err)
		return
	}
	if sub.ID == 0 {
		log.Warn("already have subject oid(%d) type(%d)", oid, tp)
		err = ecode.ReplySubjectExist
	}
	return
}

// setSubject insert or update a subject.
//
// Deprecated
func (s *Service) setSubject(c context.Context, oid int64, tp, state int8, mid int64) (sub *reply.Subject, err error) {
	now := time.Now()
	sub, err = s.dao.Subject.Get(c, oid, tp)
	if err != nil {
		return
	}
	if sub != nil && sub.AttrVal(reply.SubAttrFrozen) == reply.AttrYes {
		err = ecode.ReplySubjectFrozen
		return
	}
	sub = &reply.Subject{
		Oid:   oid,
		Type:  tp,
		Mid:   mid,
		State: state,
		CTime: xtime.Time(now.Unix()),
		MTime: xtime.Time(now.Unix()),
	}
	sub.ID, err = s.dao.Subject.Set(c, sub)
	if err != nil {
		log.Error("s.subject.Insert(%s) error(%v)", sub, err)
		return
	}
	if sub.ID == 0 {
		log.Warn("already have subject oid(%d) type(%d)", oid, tp)
		err = ecode.ReplySubjectExist
	}
	return
}

// Reply get reply from cache or db.
// NOTE old php api call
// TODO mobile jump
func (s *Service) Reply(c context.Context, oid int64, tp int8, rpID int64) (r *reply.Reply, err error) {
	if r, err = s.GetReply(c, oid, rpID, tp); err != nil {
		log.Error("s.reply(oid %d,rpid %d) err(%v)", oid, rpID, err)
		return
	}
	r.Content, _ = s.dao.Content.Get(c, oid, rpID)
	arg := &accmdl.MidReq{Mid: r.Mid}
	r.Member = new(reply.Member)
	var card *accmdl.CardReply
	if card, err = s.acc.Card3(c, arg); err != nil {
		log.Error("s.acc.Info2(%d) error(%v)", r.Mid, err)
		return
	}
	r.Member.Info = &reply.Info{}
	if card != nil {
		r.Member.Info.FromCard(card.Card)
	}
	return
}

// Deprecated: use GetReply instead
func (s *Service) reply(c context.Context, mid, oid, rpID int64, tp int8) (r *reply.Reply, err error) {
	if r, err = s.dao.Mc.GetReply(c, rpID); err != nil {
		log.Error("replyCacheDao.GetReply(%d, %d, %d) error(%v)", oid, rpID, tp, err)
		err = nil // NOTE ignore error
	}
	if r == nil {
		if r, err = s.dao.Reply.Get(c, oid, rpID); err != nil {
			log.Error("s.reply.GetReply(%d, %d) error(%v)", oid, rpID, err)
			return
		}
	}
	if r != nil {
		if r.Oid != oid || r.Type != tp {
			log.Warn("reply dismatches with parameter, oid: %d, rpID: %d, tp: %d, actual: %d, %d, %d", oid, rpID, tp, r.Oid, r.RpID, r.Type)
			err = ecode.RequestErr
			return
		}
		// NOTE if the pending reply, the state is audit
		if mid != r.Mid && !r.IsNormal() {
			err = ecode.ReplyNotExist
			return
		}
	} else {
		err = ecode.ReplyNotExist
	}
	return
}

// GetRootReply GetRootReply
func (s *Service) GetRootReply(c context.Context, oid, rpID int64, tp int8) (*reply.Reply, error) {
	r, err := s.GetReply(c, oid, rpID, tp)
	if err != nil {
		return nil, err
	}
	if r.Root != 0 {
		return nil, ecode.ReplyIllegalRoot
	}
	return r, nil
}

// GetReply GetReply
func (s *Service) GetReply(c context.Context, oid, rpID int64, tp int8) (*reply.Reply, error) {
	r, err := s.dao.Mc.GetReply(c, rpID)
	if err != nil {
		log.Error("replyCacheDao.GetReply(%d, %d, %d) error(%v)", oid, rpID, tp, err)
		err = nil // NOTE ignore error
	}
	if r == nil {
		r, err = s.dao.Reply.Get(c, oid, rpID)
		if err != nil {
			log.Error("s.reply.GetReply(%d, %d) error(%v)", oid, rpID, err)
			return nil, err
		}
		if r == nil {
			return nil, ecode.ReplyNotExist
		}
		if r.Oid != oid {
			log.Warn("reply dismatches with parameter, oid: %d, rpID: %d, tp: %d, actual: %d, %d, %d", oid, rpID, tp, r.Oid, r.RpID, r.Type)
			return nil, ecode.RequestErr
		}
	}
	return r, nil
}

// getReplyPos get root reply position.
func (s *Service) getReplyPos(c context.Context, sub *reply.Subject, rp *reply.Reply) (pos int) {
	if ok, _ := s.dao.Redis.ExpireIndex(c, sub.Oid, sub.Type, reply.SortByFloor); ok {
		var err error
		if pos, err = s.dao.Redis.RankIndex(c, sub.Oid, sub.Type, rp.RpID, reply.SortByFloor); err == nil && pos >= 0 {
			pos++
			return
		}
	}
	// If get position from redis failed, then calc by subject
	pos = sub.Count - rp.Floor + 1
	return
}

// getReplyPosByRoot get reply position from root reply.
func (s *Service) getReplyPosByRoot(c context.Context, rootRp *reply.Reply, rp *reply.Reply) (pos int) {
	if ok, _ := s.dao.Redis.ExpireIndexByRoot(c, rootRp.RpID); ok {
		var err error
		if pos, err = s.dao.Redis.RankIndexByRoot(c, rootRp.RpID, rp.RpID); err == nil && pos >= 0 {
			pos++
			return
		}
	}
	// If get position from redis failed, then calc by subject
	pos = rootRp.Count - rp.Floor + 1
	return
}

// Hide hide reply by upper.
func (s *Service) Hide(c context.Context, oid, mid, rpID int64, tp int8, ak, ck string) (err error) {
	now := time.Now()
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	if !s.isUpper(c, mid, oid, tp) {
		err = ecode.AccessDenied
		return
	}
	if _, err = s.reply(c, mid, oid, rpID, tp); err != nil {
		return
	}
	s.dao.Databus.Hide(c, oid, rpID, tp, now.Unix())
	return
}

// Show show reply by upper.
func (s *Service) Show(c context.Context, oid, mid, rpID int64, tp int8, ak, ck string) (err error) {
	now := time.Now()
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	if !s.isUpper(c, mid, oid, tp) {
		err = ecode.AccessDenied
		return
	}
	if _, err = s.reply(c, mid, oid, rpID, tp); err != nil {
		return
	}
	s.dao.Databus.Show(c, oid, rpID, tp, now.Unix())
	return
}

// Emojis get vip emojis
func (s *Service) Emojis(c context.Context) (emo []*reply.EmojiPackage) {
	emo = s.emojis
	return
}

// UpperAddTop add top reply by upper
func (s *Service) UpperAddTop(c context.Context, mid, oid, rpID int64, tp, act int8, platform string, build int64, buvid string) (err error) {
	var (
		ts = time.Now().Unix()
		r  *reply.Reply
	)
	if !reply.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	if !s.isUpper(c, mid, oid, tp) {
		err = ecode.AccessDenied
		return
	}
	sub, err := s.Subject(c, oid, tp)
	if err != nil {
		log.Error("s.Subject(oid %v) err(%v)", oid, err)
		return
	}
	if r, err = s.GetTop(c, sub, oid, tp, reply.ReplyAttrUpperTop); err != nil {
		log.Error("s.GetTop(%d,%d) err(%v)", oid, tp, err)
		return
	}
	if r != nil && act == 1 {
		log.Warn("oid(%d) type(%d) already have top ", oid, tp)
		err = ecode.ReplyHaveTop
		return
	}
	if r == nil && act == 0 {
		log.Warn("oid(%d) type(%d) do not have top ", oid, tp)
		err = ecode.ReplyNotExist
		return
	}
	if r != nil && r.RpID != rpID {
		log.Error("reply not exist top(%v) rpID(%v)", r.RpID, rpID)
		err = ecode.ReplyNotExist
		return
	}
	// TODO: only need reply,no not need content and user info
	if r, err = s.reply(c, mid, oid, rpID, tp); err != nil {
		log.Error("s.GetReply err (%v)", err)
		return
	}
	if r == nil {
		log.Warn("oid(%d) type(%d) rpID(%d) do not exist ", oid, tp, rpID)
		err = ecode.ReplyNotExist
		return
	}
	if r.AttrVal(reply.ReplyAttrAdminTop) == 1 {
		err = ecode.ReplyHaveTop
		return
	}
	if r.Root != 0 {
		log.Warn("oir(%d) type(%d) rpID(%d) not root reply", oid, tp, rpID)
		err = ecode.ReplyNotRootReply
		return
	}
	s.dao.Databus.UpperAddTop(c, mid, oid, rpID, ts, act, tp)
	var action = reply.ReportReplyTop
	if act == 0 {
		action = reply.ReportReplyUntop
	}
	ip := metadata.String(c, metadata.RemoteIP)
	report.User(&report.UserInfo{
		Mid:      r.Mid,
		Platform: platform,
		Build:    build,
		Buvid:    buvid,
		Business: 41,
		Type:     int(r.Type),
		Oid:      r.Oid,
		Action:   action,
		Ctime:    time.Now(),
		IP:       ip,
		Index: []interface{}{
			r.RpID,
		},
	})
	return
}

// GetTop get upperTop reply from cache or db.
func (s *Service) GetTop(c context.Context, sub *reply.Subject, oid int64, tp int8, top uint32) (r *reply.Reply, err error) {
	if (top == reply.ReplyAttrUpperTop) && sub.AttrVal(reply.SubAttrUpperTop) == 0 {
		return
	}
	if (top == reply.ReplyAttrAdminTop) && sub.AttrVal(reply.SubAttrAdminTop) == 0 {
		return
	}
	if r, err = s.dao.Mc.GetTop(c, oid, tp, top); err != nil {
		log.Error("s.dao.Mc.GetAdminTop(%d, %d) error(%v)", oid, tp, err)
		err = ecode.ServerErr
		return
	}
	// NOTE load by job ,in case  Cache penetration
	if r == nil {
		// if r, err = s.dao.Reply.GetTop(c, oid, tp, top); err != nil {
		// 	log.Error("s.dao.Reply.GetTop(%d, %d) error(%v)", oid, tp, err)
		// 	err = ecode.ServerErr
		// 	return
		// }
		// if r == nil {
		// 	err = ecode.ReplyNotExist
		// 	return
		// }
		// if r.Content, err = s.dao.Content.Get(c, oid, r.rpID); err != nil {
		// 	return
		// }
		// select {
		// case s.topRpChan <- topRpChan{oid: oid, tp: tp, rp: r}:
		// default:
		// 	log.Warn("s.replyChan is full")
		s.dao.Databus.AddTop(c, oid, tp, top)
	}
	return
}

// getRelation get account infos of mids
func (s *Service) getRelation(c context.Context, srcID, targetID int64, ip string) (uint32, error) {
	if targetID == 0 {
		return 0, nil
	}
	relMap, err := s.acc.RichRelations3(c, &accmdl.RichRelationReq{Owner: srcID, Mids: []int64{targetID}, RealIp: ip})
	if err != nil || relMap == nil {
		log.Error("s.acc.RichRelations2 sourceId(%v) targetId(%v)error(%v)", srcID, targetID, err)
		// return normal relation if remote service is down!
		return 0, nil
	}
	rel, ok := relMap.RichRelations[targetID]
	if !ok {
		// return normal relation if remote service is down!
		return 0, nil
	}
	return relmdl.Attr(uint32(rel)), nil
}

// RelationBlocked RelationBlocked
func (s *Service) RelationBlocked(c context.Context, srcMid, targetMid int64) bool {
	rel, _ := s.getRelation(c, srcMid, targetMid, "")
	return rel == relmdl.AttrBlack
}

// SuperviseReply Supervise Reply
func (s *Service) SuperviseReply(c context.Context, mid int64, ak, ck string, tp int8) (err error) {
	if conf.Conf.Supervision.Completed {
		err = ecode.ReplyUpgrading
		return
	}
	now := time.Now()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", conf.Conf.Supervision.StartTime, loc)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", conf.Conf.Supervision.EndTime, loc)
	if now.Before(startTime) || now.After(endTime) {
		err = nil
		return
	}
	if tp == reply.SubTypeDrawyoo {
		err = ecode.ReplyUpgrading
		return
	}
	if overseas, _ := s.checkOverseasUser(c); overseas {
		err = ecode.ReplyUpgrading
	}
	return
}

func (s *Service) checkOverseasUser(c context.Context) (overseas bool, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &locmdl.ArgIP{
		IP: ip,
	}
	overseas = false
	if respro, err := s.location.Info(c, arg); err != nil || respro == nil {
		log.Error("s.location.Info(%s) error(%v) or respro is nil", ip, err)
	} else {
		if !strings.EqualFold(respro.Country, conf.Conf.Supervision.Location) {
			overseas = true
		}
	}
	return
}

// FetchFans fetching a fans relation array between upper(mid) and the users(uids)
func (s *Service) FetchFans(c context.Context, uids []int64, mid int64) (fans map[int64]*reply.FansDetail, err error) {
	if fans, err = s.fans.Fetch(c, uids, mid, time.Now()); err != nil {
		log.Error("s.fans.fetch(%d, %d) error(%v)", mid, uids, err)
	}
	return
}

// PaginateUpperDeletedLogs paginating the admin logs for size of 'pageSize', and returning the number of reporting, the number of admin logs delete by administrator
func (s *Service) PaginateUpperDeletedLogs(c context.Context, oid int64, tp int, curPage, pageSize int) (logs []*adminlog.AdminLog, replyCount, reportCount, pageCount, total int64, err error) {
	var states = []int64{16, 18}
	if logs, replyCount, reportCount, pageCount, total, err = s.search.LogPaginate(c, oid, tp, states, curPage, pageSize, conf.Conf.AssistConfig.StartTime, time.Now()); err != nil {
		log.Error("s.adminlog.Paginate(%d, %d) error(%v)", oid, tp, err)
		return nil, 0, 0, 0, 0, err
	}
	var mids = make([]int64, 0)
	for _, d := range logs {
		mids = append(mids, d.AdminID, d.ReplyMid)
	}
	minfos, err := s.getAccInfo(c, mids)
	if err != nil {
		log.Error("s.getAccInfo(mids %v) err(%v)", mids, err)
		// NOTE degrade account
		err = nil
	}
	for _, d := range logs {
		if userInfo, ok := minfos[d.ReplyMid]; ok {
			rs := []rune(userInfo.Name)
			length := len(rs)
			if length >= 3 {
				d.ReplyUser = string(rs[0]) + "***" + string(rs[length-1])
			} else if length == 2 {
				d.ReplyUser = string(rs[0]) + "***"
			} else {
				d.ReplyUser = userInfo.Name
			}
			d.ReplyFacePic = userInfo.Face
		}
		if upperInfo, ok := minfos[d.AdminID]; ok {
			d.Operator = upperInfo.Name
		}
	}
	return
}

// GetReplyLogConfig get reply configuration from memocached or load a record from database by oid, type, category
func (s *Service) GetReplyLogConfig(c context.Context, sub *reply.Subject, category int8) (config *reply.Config, err error) {
	if sub.AttrVal(reply.SubAttrConfig) == 0 {
		return nil, nil
	}
	config, err = s.dao.Mc.GetReplyConfig(c, sub.Oid, sub.Type, category)
	if err != nil {
		log.Error("replyConfigCacheDao.GetReplyConfig(%d, %d, %d) error(%v)", sub.Oid, sub.Type, err)
		err = nil // NOTE ignore error
	}
	if config == nil {
		config, err = s.dao.Config.LoadConfig(c, sub.Oid, sub.Type, category)
		if err != nil {
			log.Error("s.reply.GetReply(%d, %d) error(%v)", sub.Oid, sub.Type, err)
			return nil, err
		}
		if config == nil {
			return nil, nil
		}
		if err = s.dao.Mc.AddReplyConfigCache(c, config); err != nil {
			log.Error("replyConfigCacheDao.AddReplyConfig(%v) error(%v)", config, err)
			return config, nil
		}
	}
	return
}

// VerifyCaptcha VerifyCaptcha
func (s *Service) VerifyCaptcha(c context.Context, captcha string, mid int64) error {
	token, err := s.dao.Mc.CaptchaToken(c, mid)
	if err != nil {
		return err
	}
	return s.dao.Captcha.Verify(c, token, captcha)
}

// Captcha return Captcha
func (s *Service) Captcha(c context.Context, mid int64) (string, error) {
	token, uri, err := s.dao.Captcha.Captcha(c)
	if err != nil {
		return "", err
	}
	s.dao.Mc.SetCaptchaToken(c, mid, token)
	return uri, nil
}

// Topics return topics
func (s *Service) Topics(c context.Context, mid int64, oid int64, typ int8, msg string) (topics []string, err error) {
	if s.IsBnj(oid, typ) {
		topics = []string{"拜年祭"}
		return
	}
	topics, err = s.bigdata.Topics(c, mid, oid, typ, msg)
	if err != nil {
		return
	}
	if len(topics) == 0 {
		return
	}
	messages := make(map[string]string)
	for i := range topics {
		key := strconv.FormatInt(int64(i), 10)
		messages[key] = topics[i]
	}
	topics = topics[:0]
	mf := &filgrpc.MFilterReq{
		Area:   "reply",
		MsgMap: messages,
	}
	res, err := s.filcli.MFilter(c, mf)
	if err != nil {
		log.Error("s.fil.MFilter(%v) failed!err:=%v", messages, err)
		return
	}
	for _, data := range res.RMap {
		if data.Level > 15 {
			continue
		}
		topics = append(topics, data.Result)
	}
	return
}

// IsHotReply IsHotReply
func (s *Service) IsHotReply(c context.Context, tp int8, oid, rpID int64) (isHot bool, err error) {
	rpIDs, _, err := s.dao.Redis.Range(c, oid, tp, reply.SortByLike, 0, 5)
	if err != nil {
		log.Error("s.dao.Redis.Range() error(%v)", err)
		return
	}
	rs, err := s.GetReplyByIDs(c, oid, tp, rpIDs)
	if err != nil {
		log.Error("s.GetReplyByIDs() error(%v)", err)
		return
	}
	for _, rp := range rs {
		if rpID == rp.RpID && rp.Like >= 3 {
			isHot = true
			return
		}
	}
	return
}

// HotsBatch return HotsBatch
func (s *Service) HotsBatch(c context.Context, tp, size int8, oids []int64, mid int64) (res map[int64][]*reply.Reply, err error) {
	var (
		missed     []int64
		missedSubs []int64
		m          sync.Mutex
		oidMap     = make(map[int64][]int64)
		subMap     map[int64]*reply.Subject
	)
	res = make(map[int64][]*reply.Reply, len(oids))
	if subMap, missedSubs, err = s.dao.Mc.GetMultiSubject(c, oids, tp); err != nil {
		log.Error("s.dao.Mc.GetMultiSubject() error(%v)", err)
		return
	}
	if len(missedSubs) > 0 {
		var missedSubMap map[int64]*reply.Subject
		if missedSubMap, err = s.dao.Subject.Gets(c, missedSubs, tp); err != nil {
			log.Error("s.dao.Subject.Gets() error(%v)", err)
			return
		}
		var subs []*reply.Subject
		for oid, sub := range missedSubMap {
			subMap[oid] = sub
			subs = append(subs, sub)
		}
		s.cache.Do(c, func(ctx context.Context) { s.dao.Mc.AddSubject(ctx, subs...) })
	}
	if oidMap, missed, err = s.dao.Redis.RangeByOids(c, oids, tp, reply.SortByLike, 0, size); err != nil {
		log.Error("s.dao.Redis.RangeByOids() error(%v)", err)
		return
	}
	if len(missed) > 0 {
		g, ctx := errgroup.WithContext(c)
		for _, missedOid := range missed {
			missedOid := missedOid
			g.Go(func() error {
				s.cache.Do(ctx, func(ctx context.Context) { s.dao.Databus.RecoverIndex(ctx, missedOid, tp, reply.SortByLike) })
				rpIDs, err := s.dao.Reply.GetIdsSortLike(ctx, missedOid, tp, 0, int(size))
				if err != nil {
					return err
				}
				m.Lock()
				oidMap[missedOid] = rpIDs
				m.Unlock()
				return nil
			})
		}
		if err = g.Wait(); err != nil {
			return
		}
	}
	for oid, rpIDs := range oidMap {
		var (
			rpsMap map[int64]*reply.Reply
			rps    = make([]*reply.Reply, 0, len(rpIDs))
		)
		// 通过评论ID获取评论内容等元信息
		if rpsMap, err = s.repliesMap(c, oid, tp, rpIDs); err != nil {
			return
		}
		for _, rpID := range rpIDs {
			if r, ok := rpsMap[rpID]; ok {
				rps = append(rps, r)
			}
		}
		if sub, ok := subMap[oid]; ok {
			// 点赞以及关注等等关系构建
			if err = s.buildReply(c, sub, rps, mid, false); err != nil {
				return
			}
		}
		res[oid] = rps
	}
	return
}
